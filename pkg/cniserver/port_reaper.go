package cniserver

import (
	"fmt"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	. "github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

type PortReaper struct {
	Opts     PortReaperOpts
	OsClient openstack.OpenstackClient
	Metrics  *Metrics
	done     func()
}

type PortReaperOpts struct {
	Interval       time.Duration
	MinPortAge     time.Duration
	MountedProcDir string
	SkipDelete     bool
}

func (me *PortReaper) Start() {
	hostname, _ := util.GetHostname()
	if me.done == nil {
		me.done = Repeat(me.Opts.Interval, func() {
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

	// make sure that /host/proc was mounted
	exists, err := util.DirExists(me.Opts.MountedProcDir)
	if err != nil {
		return err
	}
	if !exists {
		Log().Warn().Str("dir", me.Opts.MountedProcDir).
			Msg("skipping port reaping. could not find mounted proc directory")
		return nil
	}

	// list all openstack cni ports for the host using tags
	ports, err := me.OsClient.GetPortsByTags(NewPortKeyTags())
	if err != nil {
		return err
	}

	for _, port := range ports {
		if me.Opts.SkipDelete {
			Log().Info().Str("portId", port.ID).Msg("reaping disabled, skipping port")
			continue
		}
		Log().Error().Str("port", fmt.Sprintf("%v", port)).Msg("would've reaped port")
		if err := me.ReapPort(port); err != nil {
			Log().Err(err).Str("portId", port.ID).Msg("failed to reap port")
			me.Metrics.reapFailureCount.Inc()
			continue
		}
	}

	return nil
}

// Reap deletes any ports whose network namespaces no longer exist
func (me *PortReaper) ReapPort(port ports.Port) error {
	// skip ports that aren't tagged with our special identifying tag
	if !HasOpenstackCniTag(port.Tags) {
		return nil
	}
	// skip ports that were created recently
	if time.Now().Sub(port.CreatedAt) <= me.Opts.MinPortAge {
		return nil
	}

	netNs := ""
	for _, tag := range port.Tags {
		if strings.HasPrefix(tag, "netns=") {
			netNs = strings.Split(tag, "=")[1]
		}
	}
	if netNs == "" {
		return nil
	}
	netNs = strings.Replace(netNs, "/proc", me.Opts.MountedProcDir, 1)

	exists, err := util.DirExists(netNs)
	if err != nil {
		return err
	}
	if !exists {
		Log().Info().Str("portId", port.ID).Msg("attempting to reap port")
		if err := me.OsClient.DeletePort(port.ID); err != nil {
			return err
		}
		Log().Info().Str("portId", port.ID).Msg("successfully reaped port")
		me.Metrics.reapSuccessCount.Inc()
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
