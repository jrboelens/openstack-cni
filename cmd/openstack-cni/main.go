package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-chi/httplog"
	"github.com/jboelensns/openstack-cni/pkg/cniclient"
	"github.com/jboelensns/openstack-cni/pkg/cniplugin"
	"github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

func main() {
	cfg, err := cniplugin.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config %s", err))
	}

	podName := os.Getenv("K8S_POD_NAME")
	name := fmt.Sprintf("openstack-cni (%s)", podName)

	// TODO <.> refactor everything below this line into a PluginBuilder object that can be called outside of main

	// optionally create a logfile
	var output io.Writer
	if cfg.LogFileName != "" {
		output, err = os.OpenFile(
			cfg.LogFileName,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create logfile %s %s", cfg.LogFileName, err))
		}
	}

	// setup the logging
	opts := httplog.Options{LogLevel: cfg.LogLevel}
	logging.SetupLogging(name, opts, output)

	clientOpts := &cniclient.ClientOpts{
		BaseUrl:        cfg.BaseUrl,
		RequestTimeout: cfg.RequestTimeout,
		LogFileName:    cfg.LogFileName,
	}

	// create a new cniclient
	client, err := cniclient.New(clientOpts)
	if err != nil {
		logging.Error("failed to create cni client", err)
		os.Exit(1)
	}

	// setup and create the plugin
	nw := cniplugin.NewNetworking(util.NewNetlinkWrapper())
	cni := cniplugin.NewCni(client, nw)
	if err := cni.Invoke(); err != nil {
		os.Exit(1)
	}
}
