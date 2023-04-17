package cniserver

import (
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	. "github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

type PortReaper struct {
	opts   PortReaperOpts
	client openstack.OpenstackClient
	done   func()
}

type PortReaperOpts struct {
	Interval   time.Duration
	MinPortAge time.Duration
}

func NewPortReaper(client openstack.OpenstackClient, opts PortReaperOpts) *PortReaper {
	return &PortReaper{
		opts:   opts,
		client: client,
	}
}

func (me *PortReaper) Start() {
	hostname, _ := util.GetHostname()
	if me.done == nil {
		me.done = Repeat(me.opts.Interval, func() {
			if err := me.Reap(hostname); err != nil {
				Log().Err(err).Str("hostname", hostname).Msg("error reaping ports")
			}
		})
	}
}

func (me *PortReaper) Stop() {
	if me.done != nil {
		me.done()
	}
}

// Reap deletes any ports whose network namespaces no longer exist
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
		if err := me.ReapPort(port); err != nil {
			Log().Err(err).Str("portId", port.ID).Msg("failed to reap port")
		}
	}

	return nil
}

// Reap deletes any ports whose network namespaces no longer exist
func (me *PortReaper) ReapPort(port ports.Port) error {
	// skip ports that were created recently
	if time.Now().Sub(port.CreatedAt) <= me.opts.MinPortAge {
		return nil
	}

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
		Log().Info().Str("portId", port.ID).Msg("successfully reaped port")
	}
	return nil
}

// Repeat executes the fn function after each duration
// Executing the returned closer function will prevent repetition from occuring
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
