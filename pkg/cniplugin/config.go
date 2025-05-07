package cniplugin

import (
	"time"

	"github.com/jboelensns/openstack-cni/pkg/util"
	"github.com/joho/godotenv"
)

type Config struct {
	BaseUrl              string
	RequestTimeout       time.Duration
	LogFileName          string
	LogLevel             string
	WaitForUdev          bool
	WaitForUdevPrefix    string
	WaitForUdevDelay     time.Duration
	WaitForUdevTimeout   time.Duration
	EnableNetlinkExtAck  bool
	EnableNetlinkRetry   bool
	NetlinkRetryMax      time.Duration
	NetlinkRetryInterval time.Duration
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

	return Config{
		BaseUrl:              util.Getenv("CNI_API_URL", "http://127.0.0.1:4242"),
		RequestTimeout:       util.GetenvAsDuration("CNI_REQUEST_TIMEOUT", time.Second*120),
		LogFileName:          util.Getenv("CNI_LOG_FILENAME", ""),
		LogLevel:             util.Getenv("CNI_LOG_LEVEL", "info"),
		WaitForUdev:          util.GetenvAsBool("CNI_WAIT_FOR_UDEV", DefaultCniOpts().WaitForUdev),
		WaitForUdevPrefix:    util.Getenv("CNI_WAIT_FOR_UDEV_PREFIX", DefaultCniOpts().WaitForUdevPrefix),
		WaitForUdevDelay:     util.GetenvAsDuration("CNI_WAIT_FOR_UDEV_DELAY_MS", time.Millisecond*5000),
		WaitForUdevTimeout:   util.GetenvAsDuration("CNI_WAIT_FOR_UDEV_TIMEOUT_MS", time.Millisecond*5000),
		EnableNetlinkExtAck:  util.GetenvAsBool("CNI_ENABLE_NETLINK_EXT_ACK", false),
		EnableNetlinkRetry:   util.GetenvAsBool("CNI_ENABLE_NETLINK_RETRY", false),
		NetlinkRetryMax:      util.GetenvAsDuration("CNI_NETLINK_RETRY_MAX_MS", time.Millisecond*5000),
		NetlinkRetryInterval: util.GetenvAsDuration("CNI_NETLINK_RETRY_INTERVAL_MS", time.Millisecond*150),
	}, nil
}
