package cniserver

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jboelensns/openstack-cni/pkg/util"
)

type Config struct {
	ListenAddr   string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewConfig() Config {
	listenUrl := util.Getenv("CNI_API_URL", "http://127.0.0.1:4242")
	url, err := url.Parse(listenUrl)
	if err != nil {
		panic(fmt.Sprintf("invalid configuration CNI_API_URL=%s err=%s", listenUrl, err))
	}

	return Config{
		ListenAddr:   fmt.Sprintf(":%s", url.Port()),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
