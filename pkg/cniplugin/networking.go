package cniplugin

import (
	"fmt"
	"net"

	currentcni "github.com/containernetworking/cni/pkg/types/040"
	"github.com/jboelensns/openstack-cni/pkg/util"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// NetworkInterface represents a local network interface (e.g ens3)
type NetworkInterface struct {
	Name     string
	DestName string
	Address  *net.IPNet
}

//go:generate moq -pkg mocks -out ../fixtures/mocks/cniplugin_mocks.go . Networking

// Networking provides the ability to manipulate a network interface
type Networking interface {
	Configure(namespace string, iface *NetworkInterface) error
}

type networking struct {
	nl util.NetlinkWrapper
}

// NewNetworking returns a new Networking
func NewNetworking(nl util.NetlinkWrapper) *networking {
	return &networking{nl: nl}
}

// Configure moves an existing network interface into a new network namespace with the provided IP address and name
func (me *networking) Configure(namespace string, iface *NetworkInterface) error {
	// TODO: Drop in debugging

	// Find the link by interface name
	link, err := me.nl.LinkByName(iface.Name)
	if err != nil {
		return fmt.Errorf("failed to LinkByName ns=%s iface=%s dest_iface=%s addr=%s e=%w", namespace, iface.Name, iface.DestName, iface.Address, err)
	}

	// Find the destination namespace's fd by its path
	nsFd, err := me.nl.GetNetNsIdByPath(namespace)
	if err != nil {
		return fmt.Errorf("netlink failed to GetNetNsIdByPath ns=%s iface=%s dest_iface=%s addr=%s e=%w", namespace, iface.Name, iface.DestName, iface.Address, err)
	}

	// Move the link into the desination namespace
	if err := me.nl.LinkSetNsFd(link, nsFd); err != nil {
		return fmt.Errorf("netlink failed to LinkSetNsFd ns=%s iface=%s dest_iface=%s addr=%s e=%w", namespace, iface.Name, iface.DestName, iface.Address, err)
	}

	// Save our namespace so we can flip back to it once we're done
	oldNs, err := netns.Get()
	if err != nil {
		return fmt.Errorf("netlink failed to Get namespace ns=%s iface=%s dest_iface=%s addr=%s e=%w", namespace, iface.Name, iface.DestName, iface.Address, err)
	}
	defer oldNs.Close()

	// set ourselves into the destination namespace
	if err := netns.Set(netns.NsHandle(nsFd)); err != nil {
		return fmt.Errorf("netlink failed to Set namespace ns=%s iface=%s dest_iface=%s addr=%s e=%w", namespace, iface.Name, iface.DestName, iface.Address, err)
	}
	// when we're done we need to enter our original namespace
	defer netns.Set(oldNs)

	// bring the interface down before we configure it
	if err := me.nl.LinkSetDown(link); err != nil {
		return fmt.Errorf("netlink failed to LinkSetDown ns=%s iface=%s dest_iface=%s addr=%s e=%w", namespace, iface.Name, iface.DestName, iface.Address, err)
	}

	// set the name of the link
	if err := me.nl.LinkSetName(link, iface.DestName); err != nil {
		return fmt.Errorf("netlink failed to LinkSetName ns=%s iface=%s dest_iface=%s addr=%s e=%w", namespace, iface.Name, iface.DestName, iface.Address, err)
	}

	// set the IP on the interface
	ipAddr := &netlink.Addr{IPNet: iface.Address, Label: ""}
	if err := me.nl.AddrAdd(link, ipAddr); err != nil {
		return fmt.Errorf("netlink failed to AddrAdd ns=%s iface=%s dest_iface=%s addr=%s e=%w", namespace, iface.Name, iface.DestName, iface.Address, err)
	}

	// bring the interface up
	if err := me.nl.LinkSetUp(link); err != nil {
		return fmt.Errorf("netlink failed to LinkSetup ns=%s iface=%s dest_iface=%s addr=%s e=%w", namespace, iface.Name, iface.DestName, iface.Address, err)
	}
	return nil
}

// ConfigureInterface sets up the interfaces with the correct name, network namesapce and ip address
func (me *Cni) ConfigureInterface(cmd util.CniCommand, result *currentcni.Result) error {

	mac := result.Interfaces[0].Mac
	ifaceName, err := GetIfaceNameByMac(mac)
	if err != nil {
		return err
	}

	iface := &NetworkInterface{
		Name:     ifaceName,
		DestName: result.Interfaces[0].Name,
		Address:  &result.IPs[0].Address,
	}

	err = me.nw.Configure(cmd.Netns, iface)
	if err != nil {
		return fmt.Errorf("failed to configure interface %w", err)
	}
	return err
}

// GetIfaceNameByMac returns the name of an interface matching the given MAC address
func GetIfaceNameByMac(mac string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("error finding interface for mac=%s e=%w", mac, err)
	}

	for _, iface := range ifaces {
		if iface.HardwareAddr.String() == mac {
			return iface.Name, nil
		}
	}

	return "", fmt.Errorf("failed to find interface for %s", mac)
}
