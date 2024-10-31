package cniplugin

import (
	"fmt"
	"time"

	"github.com/jboelensns/openstack-cni/pkg/util"
	"github.com/joho/godotenv"
)

type Config struct {
	BaseUrl            string
	RequestTimeout     time.Duration
	LogFileName        string
	LogLevel           string
	WaitForUdev        bool
	WaitForUdevPrefix  string
	WaitForUdevDelay   time.Duration
	WaitForUdevTimeout time.Duration
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

	timeout, err := time.ParseDuration(fmt.Sprintf("%ss", util.Getenv("CNI_REQUEST_TIMEOUT", "120")))
	if err != nil {
		return config, err
	}

	waitForUdevDelay, err := time.ParseDuration(fmt.Sprintf("%sms", util.Getenv("CNI_WAIT_FOR_UDEV_DELAY_MS", "100")))
	if err != nil {
		return config, err
	}

	waitForUdevTimeout, err := time.ParseDuration(fmt.Sprintf("%sms", util.Getenv("CNI_WAIT_FOR_UDEV_TIMEOUT_MS", "5000")))
	if err != nil {
		return config, err
	}
	return Config{
		BaseUrl:            util.Getenv("CNI_API_URL", "http://127.0.0.1:4242"),
		RequestTimeout:     timeout,
		LogFileName:        util.Getenv("CNI_LOG_FILENAME", ""),
		LogLevel:           util.Getenv("CNI_LOG_LEVEL", "info"),
		WaitForUdev:        util.GetenvAsBool("CNI_WAIT_FOR_UDEV", DefaultCniOpts().WaitForUdev),
		WaitForUdevPrefix:  util.Getenv("CNI_WAIT_FOR_UDEV_PREFIX", DefaultCniOpts().WaitForUdevPrefix),
		WaitForUdevDelay:   waitForUdevDelay,
		WaitForUdevTimeout: waitForUdevTimeout,
	}, nil
}
