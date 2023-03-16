package main

import (
	"fmt"
	"os"

	"github.com/go-chi/httplog"
	"github.com/jboelensns/openstack-cni/pkg/cniclient"
	"github.com/jboelensns/openstack-cni/pkg/cniplugin"
	 "github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

func main() {
	podName := os.Getenv("K8S_POD_NAME")
	name := fmt.Sprintf("openstack-cni (%s)", podName)
	opts := httplog.Options{LogLevel: "info"}
	logging.SetupLogging(name, opts)

	client, err := cniclient.New(nil)
	if err != nil {
		logging.Error("failed to create cni client", err)
		os.Exit(1)
	}

	nw := cniplugin.NewNetworking(util.NewNetlinkWrapper())
	cni := cniplugin.NewCni(client, nw)
	if err := cni.Invoke(); err != nil {
		os.Exit(1)
	}
}
