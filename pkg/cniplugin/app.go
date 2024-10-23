package cniplugin

import (
	"fmt"
	"io"
	"os"

	"github.com/go-chi/httplog"
	"github.com/jboelensns/openstack-cni/pkg/cniclient"
	"github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

// App represents the cniplugin
type App struct {
	config Config
}

// NewApp creates a new App from configuration
func NewApp(config Config) *App {
	return &App{config: config}
}

// Run starts the cniplugin
func (me *App) Run() error {
	podName := os.Getenv("K8S_POD_NAME")
	name := fmt.Sprintf("openstack-cni (%s)", podName)

	// optionally create a logfile
	var output io.Writer
	var err error
	if me.config.LogFileName != "" {
		output, err = os.OpenFile(
			me.config.LogFileName,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create logfile %s %s", me.config.LogFileName, err))
		}
	}

	// setup the logging
	opts := httplog.Options{LogLevel: me.config.LogLevel}
	logging.SetupLogging(name, opts, output)

	clientOpts := &cniclient.ClientOpts{
		BaseUrl:        me.config.BaseUrl,
		RequestTimeout: me.config.RequestTimeout,
		LogFileName:    me.config.LogFileName,
	}

	// create a new cniclient
	client, err := cniclient.New(clientOpts)
	if err != nil {
		logging.Error("failed to create cni client", err)
		os.Exit(1)
	}

	// setup and create the plugin
	nw := NewNetworking(util.NewNetlinkWrapper())
	cni := NewCni(client, nw)
	return cni.Invoke()
}
