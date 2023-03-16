package main

import (
	"os"

	"github.com/go-chi/httplog"
	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	"github.com/jboelensns/openstack-cni/pkg/logging"
)

func main() {
	opts := httplog.Options{LogLevel: "info"}
	logging.SetupLogging("openstack-cni-daemon", opts)

	if err := cniserver.Run(); err != nil {
		logging.Error("failed to run server", err)
		os.Exit(1)
	}
}
