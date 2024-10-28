package main

import (
	"fmt"
	"os"

	"github.com/jboelensns/openstack-cni/pkg/cniplugin"
)

func main() {
	config, err := cniplugin.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config %s", err))
	}

	if err := cniplugin.NewApp(config).Run(); err != nil {
		os.Exit(1)
	}
}
