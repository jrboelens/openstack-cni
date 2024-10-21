package cniplugin

import (
	"fmt"
	"time"

	"github.com/jboelensns/openstack-cni/pkg/util"
	"github.com/joho/godotenv"
)

type Config struct {
	BaseUrl        string
	RequestTimeout time.Duration
	LogFileName    string
	LogLevel       string
}

func LoadConfig() (Config, error) {
	var config Config
	// attempt to read config file
	configFile := util.Getenv("CNI_CONFIG_FILE", "/etc/cni/net.d/openstack-cni.conf")
	exists, err := util.FileExists(configFile)
	if err != nil {
		return config, err
	}
	if exists {
		if err := godotenv.Load(configFile); err != nil {
			return config, err
		}
	}

	timeout, err := time.ParseDuration(fmt.Sprintf("%ss", util.Getenv("CNI_REQUEST_TIMEOUT", "60")))
	if err != nil {
		return config, err
	}
	return Config{
		BaseUrl:        util.Getenv("CNI_API_URL", "http://127.0.0.1:4242"),
		RequestTimeout: timeout,
		LogFileName:    util.Getenv("CNI_LOG_FILENAME", ""),
		LogLevel:       util.Getenv("CNI_LOG_LEVEL", "info"),
	}, nil
}
