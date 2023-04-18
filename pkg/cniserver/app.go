package cniserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/httplog"
	"github.com/hashicorp/go-multierror"
	. "github.com/jboelensns/openstack-cni/pkg/logging"
)

// App represents the application running the http server
type App struct {
	config  Config
	server  *http.Server
	reaper  *PortReaper
	metrics *Metrics
}

// NewApp creates a new App from configuration
func NewApp(config Config, server *http.Server, reaper *PortReaper) (*App, error) {
	return &App{
		config: config,
		server: server,
		reaper: reaper,
	}, nil
}

// Run starts the http server and blocks
func (me *App) Run() error {
	Log().Info().Str("duration", me.config.ReapInterval.String()).Msg("starting port reaper")
	me.reaper.Start()
	Log().Info().Str("addr", me.server.Addr).Msg("starting http server")
	if err := me.server.ListenAndServe(); err != http.ErrServerClosed {
		Error("failed to start http server", err)
		return err
	}
	return nil
}

// Shutdown signals the http server to shutdown
func (me *App) Shutdown(ctx context.Context) error {
	Log().Info().Msg("shutting port reaper")
	me.reaper.Stop()
	Log().Info().Msg("shut down port reaper")
	Log().Info().Msg("shutting down http server")
	defer func() {
		Log().Info().Msg("shut down http server")
	}()
	return me.server.Shutdown(ctx)
}

func BuildApp() (*App, error) {
	SetupLogging("openstack-cni-daemon", httplog.DefaultOptions)
	Log().Info().Msg("preparing http server")

	config := NewConfig()
	deps, err := NewBuilder(config).Build()
	if err != nil {
		Error("failed to build dependencies", err)
		return nil, err
	}
	app, err := NewApp(config, deps.RestServer(), deps.PortReaper())
	if err != nil {
		Log().Error().Str("addr", app.config.ListenAddr).AnErr("err", err).Msg("failed to initialize server")
		return nil, err
	}
	return app, err
}

// Run builds up the dependencies for the application, creates the application and runs it
func Run() error {
	app, err := BuildApp()
	if err != nil {
		return err
	}

	go app.HandleSignals()
	return app.Run()
}

// HandleSignals setups signal handling
func (me *App) HandleSignals() {
	go func() {
		err := HandleSignals(context.Background(), me.Shutdown)
		if err != nil {
			Error("error handling signals", err)
		}
	}()
}

// PingHandler handles /ping requests
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}

// HandleSignals handlers SIGINT and SIGINT signals
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
