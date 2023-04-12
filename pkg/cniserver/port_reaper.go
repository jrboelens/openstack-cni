package cniserver

import (
	"strings"

	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

type PortReaper struct {
	client openstack.OpenstackClient
}

func NewPortReaper(client openstack.OpenstackClient) *PortReaper {
	return &PortReaper{client}
}

func (me *PortReaper) Reap(hostname string) error {
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
			if err := me.client.DeletePort(port.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
