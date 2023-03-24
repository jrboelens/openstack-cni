package main

import (
	"os"

	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	"github.com/jboelensns/openstack-cni/pkg/logging"
)

func main() {
	if err := cniserver.Run(); err != nil {
		logging.Error("failed to run server", err)
		os.Exit(1)
	}
}
