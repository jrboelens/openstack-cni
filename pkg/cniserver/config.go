package cniserver

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jboelensns/openstack-cni/pkg/util"
)

// Config is used to configure the application
type Config struct {
	ListenAddr   string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ReapInterval time.Duration
	MinPortAge   time.Duration
	SkipReaping  bool
}

// NewConfig creates a new default Config
func NewConfig() Config {
	listenUrl := util.Getenv("CNI_API_URL", "http://127.0.0.1:4242")
	url, err := url.Parse(listenUrl)
	if err != nil {
		panic(fmt.Sprintf("invalid configuration CNI_API_URL=%s err=%s", listenUrl, err))
	}

	return Config{
		ListenAddr:   url.Host,
		ReadTimeout:  getEnvDuration("CNI_READ_TIMEOUT", "10s"),
		WriteTimeout: getEnvDuration("CNI_WRITE_TIMEOUT", "10s"),
		ReapInterval: getEnvDuration("CNI_REAP_INTERVAL", "300s"),
		MinPortAge:   getEnvDuration("CNI_MIN_PORT_AGE", "300s"),
		SkipReaping:  getEnvBool("CNI_SKIP_REAPING", "false"),
	}
}

func getEnvDuration(name, defVal string) time.Duration {
	envStr := util.Getenv(name, defVal)
	duration, err := time.ParseDuration(envStr)
	if err != nil {
		panic(fmt.Sprintf("invalid configuration %s=%s err=%s", name, envStr, err))
	}
	return duration
}

func getEnvBool(name, defVal string) bool {
	envStr := util.Getenv(name, defVal)
	b, err := strconv.ParseBool(envStr)
	if err != nil {
		panic(fmt.Sprintf("invalid configuration %s=%s err=%s", name, envStr, err))
	}
	return b
}
