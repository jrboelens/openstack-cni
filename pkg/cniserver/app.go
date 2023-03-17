package cniserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
	"github.com/hashicorp/go-multierror"
	. "github.com/jboelensns/openstack-cni/pkg/logging"
)

type App struct {
	config Config
	server *http.Server
}

func SetupRoutes(deps *Deps) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/health", (NewHealthHandler(deps.OpenstackClient())).HandleRequest)
	router.Get("/ping", PingHandler)
	router.Post("/cni", (&CniHandler{deps.CniHandler()}).HandleRequest)

	// state
	stateHandler := &StateHandler{deps.State()}
	router.Get("/state/{containerId}/{ifname}", stateHandler.Get)
	router.Delete("/state/{containerId}/{ifname}", stateHandler.Delete)
	router.Post("/state", stateHandler.Set)
	return router
}

func NewApp(config Config, mux http.Handler) (*App, error) {
	return &App{
		config: config,
		server: &http.Server{
			Addr:         config.ListenAddr,
			Handler:      mux,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
	}, nil
}

func (me *App) Run() error {
	Log().Info().Str("addr", me.server.Addr).Msg("starting http server")
	if err := me.server.ListenAndServe(); err != http.ErrServerClosed {
		Error("failed to start http server", err)
		return err
	}

	return nil
}

func (me *App) Shutdown(ctx context.Context) error {
	Log().Info().Msg("shutting down server")
	defer func() {
		Log().Info().Msg("shut down server")
	}()
	return me.server.Shutdown(ctx)
}

func Run() error {
	SetupLogging("openstack-cni-daemon", httplog.DefaultOptions)
	Log().Info().Msg("preparing http server")

	deps, err := NewBuilder().Build()
	if err != nil {
		Error("failed to build dependencies", err)
		return err
	}
	app, err := NewApp(NewConfig(), SetupRoutes(deps))
	if err != nil {
		Log().Error().Str("addr", app.config.ListenAddr).AnErr("err", err).Msg("failed to initialize server")
		return err
	}

	go app.HandleSignals()
	return app.Run()
}

func (me *App) HandleSignals() {
	go func() {
		err := HandleSignals(context.Background(), me.Shutdown)
		if err != nil {
			Error("error handling signals", err)
		}
	}()
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}

func HandleSignals(ctx context.Context, callbacks ...func(context.Context) error) error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	var errs *multierror.Error
	for _, callback := range callbacks {
		err := callback(ctx)
		if err != nil {
			errs = multierror.Append(err)
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}
