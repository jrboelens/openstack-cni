package cniserver

import (
	"fmt"
	"path"
	"strconv"
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
	ProcMount  string
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
	log.Info().Str("tags", strings.Join(portTags, ",")).Msg("searching for tagged ports")
	ports, err := me.OsClient.GetPortsByTags(portTags)
	if err != nil {
		return err
	}
	if len(ports) > 0 {
		log.Info().Int("port_count", len(ports)).Msg("found tagged ports")
	} else {
		log.Info().Msg("did not find tagged ports")
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

	// grab the netns tag from the port
	portTags := NewPortTags(port)
	if portTags.Netns == "" {
		log.Info().Msg("skipping port delete.. empty netns")
		return nil
	}

	// validate netns
	procName := "proc"
	netns_parts := strings.Split(portTags.Netns, "/")
	if len(netns_parts) < 5 {
		log.Info().Str("netns", portTags.Netns).Msg("skipping port delete.. invalid netns")
		return nil
	}
	if netns_parts[1] != procName {
		log.Info().Str("netns", portTags.Netns).Msg("skipping port delete.. netns is not in /proc")
		return nil
	}
	hostPid := netns_parts[2]
	_, err := strconv.ParseInt(hostPid, 10, 64)
	if err != nil {
		log.Info().Str("netns", portTags.Netns).Str("host_pid", hostPid).Msg("skipping port delete.. netns contains invalid pid")
		return nil
	}

	// ensure the configured proc mount exists
	log.Info().Str("proc_mount", me.Opts.ProcMount).Msg("checking proc_mount existence")
	exists, err := util.DirExists(me.Opts.ProcMount)
	if err != nil {
		log.Info().AnErr("err", err).Str("proc_mount", me.Opts.ProcMount).Msg("skipping port delete.. error checking if proc_mount exists")
		return nil
	}
	if !exists {
		log.Info().Str("proc_mount", me.Opts.ProcMount).Msg("skipping port delete.. proc_mount doesn't exist")
		return nil
	}

	// ensure we're actually looking at the proc mount by looking for pid #1
	pid1Path := path.Join(me.Opts.ProcMount, "1")
	log.Info().Str("pid1_path", pid1Path).Msg("checking pid1_path existence")
	exists, err = util.DirExists(pid1Path)
	if err != nil {
		log.Info().AnErr("err", err).Str("pid1_path", pid1Path).Msg("skipping port delete.. error checking if pid1_path exists")
		return nil
	}
	if !exists {
		log.Info().Str("pid1_path", pid1Path).Msg("skipping port delete.. pid1_path doesn't exist")
		return nil
	}

	// check to see if the network namespace exists in the mounted /host/proc
	hostNetns := path.Join(me.Opts.ProcMount, hostPid, "ns", "net")
	log.Info().Str("host_ns", hostNetns).Msg("checking netns existence")
	exists, err = util.FileExists(hostNetns)
	if err != nil {
		log.Info().AnErr("err", err).Str("host_netns", hostNetns).Msg("skipping port delete.. error checking if host netns exists")
		return nil
	}
	if exists {
		log.Info().Str("host_ns", hostNetns).Msg("skipping port delete.. host netns exists")
		return nil
	}

	log.Info().Msg("attempting to reap port")
	if err := me.OsClient.DeletePort(port.ID); err != nil {
		return err
	}
	log.Info().Msg("successfully reaped port")
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
