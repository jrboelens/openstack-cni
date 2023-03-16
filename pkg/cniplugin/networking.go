package cniplugin

import (
	"net"

	"github.com/jboelensns/openstack-cni/pkg/util"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type NetworkInterface struct {
	Name     string
	DestName string
	Address  *net.IPNet
}

//go:generate moq -pkg mocks -out ../fixtures/mocks/cniplugin_mocks.go . Networking

type Networking interface {
	Configure(namespace string, iface *NetworkInterface) error
}

type networking struct {
	nl util.NetlinkWrapper
}

func NewNetworking(nl util.NetlinkWrapper) *networking {
	return &networking{nl: nl}
}

func (me *networking) Configure(namespace string, iface *NetworkInterface) error {
	// Find the link by interface name
	link, err := me.nl.LinkByName(iface.Name)
	if err != nil {
		return err
	}

	// Find the destination namespace's fd by its path
	nsFd, err := me.nl.GetNetNsIdByPath(namespace)
	if err != nil {
		return err
	}

	// Move the link into the desination namespace
	if err := me.nl.LinkSetNsFd(link, nsFd); err != nil {
		return err
	}

	// Save our namespace so we can flip back to it once we're done
	oldNs, err := netns.Get()
	if err != nil {
		return err
	}
	defer oldNs.Close()

	// set ourselves into the destination namespace
	if err := netns.Set(netns.NsHandle(nsFd)); err != nil {
		return err
	}
	// when we're done we need to enter our original namespace
	defer netns.Set(oldNs)

	// bring the interface down before we configure it
	if err := me.nl.LinkSetDown(link); err != nil {
		return err
	}

	// set the name of the link
	if err := me.nl.LinkSetName(link, iface.DestName); err != nil {
		return err
	}

	// set the IP on the interface
	ipAddr := &netlink.Addr{IPNet: iface.Address, Label: ""}
	if err := me.nl.AddrAdd(link, ipAddr); err != nil {
		return err
	}

	// bring the interface up
	if err := me.nl.LinkSetUp(link); err != nil {
		return err
	}
	return nil
}
