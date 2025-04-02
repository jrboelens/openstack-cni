package cniserver

import (
	. "github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

type PortCounter struct {
	OsClient openstack.OpenstackClient
}

// Count determines how many openstack-cni ports are currently assigned to the host
func (me *PortCounter) Count() float64 {
	hostname, _ := util.GetHostname()
	Log().Info().Str("hostname", hostname).Msg("counting ports")
	// lookup the server
	server, err := me.OsClient.GetServerByName(hostname)
	if err != nil {
		Log().Err(err).Str("hostname", hostname).Msg("failed to GetServerByName while counting ports")
		return 0
	}

	// list all ports for a host
	ports, err := me.OsClient.GetPortsByDeviceId(server.ID)
	if err != nil {
		Log().Err(err).Str("deviceId", server.ID).Msg("failed to GetPortsByDeviceId while counting ports")
		return 0
	}

	var count float64 = 0
	for _, port := range ports {
		if !HasOpenstackCniTag(port.Tags) {
			continue
		}
		count++
	}

	Log().Info().Str("hostname", hostname).Float64("count", count).Msg("found ports")
	return count
}
