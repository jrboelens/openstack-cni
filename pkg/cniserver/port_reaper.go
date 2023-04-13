package cniserver

import (
	"os"
	"strings"
	"time"

	. "github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

type PortReaper struct {
	opts     PortReaperOpts
	client   openstack.OpenstackClient
	hostname string
	done     func()
}

type PortReaperOpts struct {
	Interval time.Duration
	Client   openstack.OpenstackClient
}

func NewPortReaper(opts PortReaperOpts) *PortReaper {
	hostname, _ := os.Hostname()
	return &PortReaper{
		opts:     opts,
		client:   opts.Client,
		hostname: hostname,
	}
}

func (me *PortReaper) Start() {
	if me.done == nil {
		me.done = Repeat(me.opts.Interval, func() {
			if err := me.Reap(me.hostname); err != nil {
				Log().Err(err).Str("hostname", me.hostname).Msg("error reaping ports")
			}
		})
	}
}

func (me *PortReaper) Stop() {
	if me.done != nil {
		me.done()
	}
}

// break this out into another struct/func or inject it
func (me *PortReaper) Reap(hostname string) error {
	Log().Info().Msg("reaping ports")
	// lookup the server
	server, err := me.client.GetServerByName(hostname)
	if err != nil {
		return err
	}

	// list all ports for a host
	ports, err := me.client.GetPortsByDeviceId(server.ID)
	if err != nil {
		return err
	}

	for _, port := range ports {
		// skip new ports

		netNs := ""
		for _, tag := range port.Tags {
			if strings.HasPrefix(tag, "netns=") {
				netNs = strings.Split(tag, "=")[1]
			}
		}

		if !util.DirExists(netNs) {
			Log().Info().Str("portId", port.ID).Msg("attempting to reap port")
			if err := me.client.DeletePort(port.ID); err != nil {
				return err
			}
			Log().Info().Str("portId", port.ID).Msg("succesfully reaped port")
		}
	}

	return nil
}

func Repeat(d time.Duration, fn func()) (closer func()) {
	ticker := time.NewTicker(d)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fn()
			case <-done:
				return
			}
		}
	}()

	return func() {
		close(done)
	}
}
