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
	Interval   time.Duration
	MinPortAge time.Duration
	SkipDelete bool
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
	log := Log().With().Str("hostname", hostname).Logger()
	log.Info().Msg("attempting reaping ports")

	// list all openstack cni ports for the host using tags
	portTags := NewPortKeyTags()
	log.Info().Str("tags", strings.Join(portTags, ",")).Msg("searching for reapable ports")
	ports, err := me.OsClient.GetPortsByTags(portTags)
	if err != nil {
		return err
	}
	if len(ports) > 0 {
		log.Info().Int("port_count", len(ports)).Msg("found repable ports")
	} else {
		log.Info().Msg("did not find reapable ports")
	}

	for _, port := range ports {
		if me.Opts.SkipDelete {
			log.Info().Str("port_id", port.ID).Msg("reaping disabled, skipping port")
			continue
		}
		if err := me.ReapPort(port); err != nil {
			log.Err(err).Str("port_id", port.ID).Msg("failed to reap port")
			me.Metrics.reapFailureCount.Inc()
			continue
		}
	}

	return nil
}

// Reap deletes any ports whose network namespaces no longer exist
func (me *PortReaper) ReapPort(port ports.Port) error {
	log := Log().With().Str("port_id", port.ID).Str("status", port.Status).Str("tags", strings.Join(port.Tags, ",")).Str("created_at", port.CreatedAt.String()).Logger()
	log.Info().Msg("attempting to reap port")

	// skip ports that aren't tagged with our special identifying tag
	if !HasOpenstackCniTag(port.Tags) {
		log.Info().Msg("skipping port delete.. missing openstack-cni=true tag")
		return nil
	}
	// skip ports that were created recently
	diff := time.Now().Sub(port.CreatedAt)
	if diff <= me.Opts.MinPortAge {
		log.Info().Msg(fmt.Sprintf("skipping port delete.. port is too new %s > %s", diff, me.Opts.MinPortAge))
		return nil
	}

	// only delete DOWN ports
	if port.Status != "DOWN" {
		log.Info().Msg("skipping port delete.. port status is not DOWN")
		return nil
	}

	// only delete detached ports
	if port.DeviceID != "" {
		log.Info().Str("device_id", port.DeviceID).Msg("skipping port delete.. still attached")
		return nil
	}

	log.Info().Str("port_id", port.ID).Msg("attempting to reap port")
	if err := me.OsClient.DeletePort(port.ID); err != nil {
		return err
	}
	log.Info().Str("port_id", port.ID).Msg("successfully reaped port")
	me.Metrics.reapSuccessCount.Inc()
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
