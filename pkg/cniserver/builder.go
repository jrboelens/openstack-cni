package cniserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Deps represents dependencies for the application
// instances of this structure are created by the Builder
type Deps struct {
	cniHandler CommandHandler
	osClient   openstack.OpenstackClient
	metrics    *Metrics
	portReaper *PortReaper
	restServer *http.Server
}

// CniHandler returns the CommandHandler
func (me *Deps) CniHandler() CommandHandler {
	return me.cniHandler
}

// Metrics returns the Metrics
func (me *Deps) Metrics() *Metrics {
	return me.metrics
}

// OpenstackClient returns the OpenstackClient
func (me *Deps) OpenstackClient() openstack.OpenstackClient {
	return me.osClient
}

// PortReaper returns a PortReapder
func (me *Deps) PortReaper() *PortReaper {
	return me.portReaper
}

// RestServer returns an http.server
func (me *Deps) RestServer() *http.Server {
	return me.restServer
}

// Builder provides the ability to produce Deps instances using the builder pattern
type Builder struct {
	config      Config
	cniHandler  CommandHandler
	osClient    openstack.OpenstackClient
	metrics     *Metrics
	restServer  *http.Server
	portReaper  *PortReaper
	portCounter *PortCounter
}

// NewBuilder creates a new Builder
func NewBuilder(config Config) *Builder {
	return &Builder{config: config}
}

// WithCniHandler sets the current to CommandHandler to cniHandler
func (me *Builder) WithCniHandler(cniHandler CommandHandler) *Builder {
	me.cniHandler = cniHandler
	return me
}

// WithMetrics sets the current Metrics
func (me *Builder) WithMetrics(metrics *Metrics) *Builder {
	me.metrics = metrics
	return me
}

// WithOpenstackClient sets the current to OpenstackClient to client
func (me *Builder) WithOpenstackClient(client openstack.OpenstackClient) *Builder {
	me.osClient = client
	return me
}

// WithPortReaper sets the PortReaper
func (me *Builder) WithPortReaper(reaper *PortReaper) *Builder {
	me.portReaper = reaper
	return me
}

// WithOpenstackClient sets the current to OpenstackClient to client
func (me *Builder) WithRestServer(server *http.Server) *Builder {
	me.restServer = server
	return me
}

// Build creates a Dep
func (me *Builder) Build() (*Deps, error) {
	// build the default os factory if we don't have one
	if me.osClient == nil {
		var err error
		me.osClient, err = openstack.NewOpenstackClient()
		if err != nil {
			return nil, fmt.Errorf("failed to build openstack client err=%w", err)
		}
	}

	me.osClient = openstack.NewCachedClient(me.osClient, getEnvDuration("CNI_CACHE_TTL", "300s"))

	// build the default cni handler if we don't have one
	if me.cniHandler == nil {
		var err error
		pm := openstack.NewPortManager(me.osClient)
		me.cniHandler, err = NewCniCommandHandler(pm), nil
		if err != nil {
			return nil, fmt.Errorf("failed to build cni command handler err=%w", err)
		}
	}

	if me.portCounter == nil {
		me.portCounter = &PortCounter{me.osClient}
	}

	if me.metrics == nil {
		registry := prometheus.NewRegistry()
		me.metrics = NewMetrics(registry, me.portCounter.Count)
	}

	if me.portReaper == nil {
		me.portReaper = &PortReaper{
			Opts: PortReaperOpts{
				Interval:       me.config.ReapInterval,
				MinPortAge:     me.config.MinPortAge,
				MountedProcDir: "/host/proc",
			},
			OsClient: me.osClient,
			Metrics:  me.metrics,
		}
	}

	if me.restServer == nil {
		router := chi.NewRouter()
		router.Use(middleware.Logger)
		router.Get("/health", (&HealthHandler{me.osClient}).HandleRequest)
		router.Get("/ping", PingHandler)
		router.Post("/cni", (&CniHandler{me.cniHandler, me.metrics}).HandleRequest)
		router.Get("/metrics", promhttp.HandlerFor(me.metrics.Registry(), promhttp.HandlerOpts{Registry: me.metrics.Registry()}).ServeHTTP)

		me.restServer = &http.Server{
			Addr:         me.config.ListenAddr,
			Handler:      router,
			ReadTimeout:  me.config.ReadTimeout,
			WriteTimeout: me.config.WriteTimeout,
		}
	}

	return &Deps{
		cniHandler: me.cniHandler,
		osClient:   me.osClient,
		metrics:    me.metrics,
		portReaper: me.portReaper,
		restServer: me.restServer,
	}, nil
}
