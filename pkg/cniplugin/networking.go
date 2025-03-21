package cniplugin

import (
	"fmt"
	"net"
	"strings"
	"time"

	currentcni "github.com/containernetworking/cni/pkg/types/040"
	"github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// NetworkInterface represents a local network interface (e.g ens3)
type NetworkInterface struct {
	Index    int
	DestName string
	Address  *net.IPNet
}

//go:generate moq -pkg mocks -out ../fixtures/mocks/cniplugin_mocks.go . Networking

// Networking provides the ability to manipulate a network interface
type Networking interface {
	Configure(namespace string, iface *NetworkInterface) error
	GetIfaceByMac(mac string) (*net.Interface, error)
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
	logger := logging.Log().With().
		Int("iface_index", iface.Index).
		Str("namespace", namespace).Str("dest_iface", iface.DestName).
		Str("addr", iface.Address.IP.String()).Logger()

	// Find the link by interface name
	logger.Info().Msg("calling netlink.LinkByIndex")
	link, err := me.nl.LinkByIndex(iface.Index)
	if err != nil {
		return fmt.Errorf("failed to LinkByIndex ns=%s iface_index=%d dest_iface=%s addr=%s e=%w", namespace, iface.Index, iface.DestName, iface.Address, err)
	}

	linkAttrs := link.Attrs()
	if linkAttrs != nil {
		logger = logger.With().Str("iface", linkAttrs.Name).
			Str("mac", linkAttrs.HardwareAddr.String()).Logger()
		logger.Info().Msg("found link")
	} else {
		logger.Info().Str("type", link.Type()).Str("attrs", "nil").Msg("found link")
	}

	// Find the destination namespace's fd by its path
	logger.Info().Msg("calling netlink.GetNetNsIdByPath")
	nsFd, err := me.nl.GetNetNsIdByPath(namespace)
	if err != nil {
		return fmt.Errorf("netlink failed to GetNetNsIdByPath ns=%s iface=%s iface_index=%d dest_iface=%s addr=%s e=%w", namespace, linkAttrs.Name, iface.Index, iface.DestName, iface.Address, err)
	}

	// Move the link into the desination namespace
	logger.Info().Msg("calling netlink.LinkSetNsFd")
	if err := me.nl.LinkSetNsFd(link, nsFd); err != nil {
		return fmt.Errorf("netlink failed to LinkSetNsFd ns=%s iface=%s iface_index=%d dest_iface=%s addr=%s e=%w", namespace, linkAttrs.Name, iface.Index, iface.DestName, iface.Address, err)
	}

	// Save our namespace so we can flip back to it once we're done
	logger.Info().Msg("calling netlink.Get")
	oldNs, err := netns.Get()
	if err != nil {
		return fmt.Errorf("netlink failed to Get namespace ns=%s iface=%s iface_index=%d dest_iface=%s addr=%s e=%w", namespace, linkAttrs.Name, iface.Index, iface.DestName, iface.Address, err)
	}
	defer oldNs.Close()

	// set ourselves into the destination namespace
	logger.Info().Msg("calling netlink.NsHandle")
	if err := netns.Set(netns.NsHandle(nsFd)); err != nil {
		return fmt.Errorf("netlink failed to Set namespace ns=%s iface=%s iface_index=%d dest_iface=%s addr=%s e=%w", namespace, linkAttrs.Name, iface.Index, iface.DestName, iface.Address, err)
	}
	// when we're done we need to enter our original namespace
	defer netns.Set(oldNs)

	// bring the interface down before we configure it
	logger.Info().Msg("calling netlink.LinkSetDown")
	if err := me.nl.LinkSetDown(link); err != nil {
		return fmt.Errorf("netlink failed to LinkSetDown ns=%s iface=%s iface_index=%d dest_iface=%s addr=%s e=%w", namespace, linkAttrs.Name, iface.Index, iface.DestName, iface.Address, err)
	}

	// set the name of the link
	logger.Info().Msg("calling netlink.LinkSetName")
	if err := me.nl.LinkSetName(link, iface.DestName); err != nil {
		return fmt.Errorf("netlink failed to LinkSetName ns=%s iface=%s iface_index=%d dest_iface=%s addr=%s e=%w", namespace, linkAttrs.Name, iface.Index, iface.DestName, iface.Address, err)
	}

	// set the IP on the interface
	logger.Info().Msg("calling netlink.AddrAdd")
	ipAddr := &netlink.Addr{IPNet: iface.Address, Label: ""}
	if err := me.nl.AddrReplace(link, ipAddr); err != nil {
		return fmt.Errorf("netlink failed to AddrReplace ns=%s iface=%s iface_index=%d dest_iface=%s addr=%s e=%w", namespace, linkAttrs.Name, iface.Index, iface.DestName, iface.Address, err)
	}

	// bring the interface up
	logger.Info().Msg("calling netlink.LinkSetUp")
	if err := me.nl.LinkSetUp(link); err != nil {
		return fmt.Errorf("netlink failed to LinkSetup ns=%s iface=%s iface_index=%d dest_iface=%s addr=%s e=%w", namespace, linkAttrs.Name, iface.Index, iface.DestName, iface.Address, err)
	}
	return nil
}

// GetIfaceByMac returns an interface matching the given MAC address
func (me *networking) GetIfaceByMac(mac string) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("error finding interface for mac=%s e=%w", mac, err)
	}

	for _, iface := range ifaces {
		if iface.HardwareAddr.String() == mac {
			return &iface, nil
		}
	}

	return nil, fmt.Errorf("failed to find interface for %s", mac)
}

// ConfigureInterface sets up the interfaces with the correct name, network namesapce and ip address
func (me *Cni) ConfigureInterface(cmd util.CniCommand, result *currentcni.Result) error {
	mac := result.Interfaces[0].Mac

	// ensure that if udev rules are in use they have had time to run
	// this accounts accounts for a race condition between nova/neutron creation/attachment and the interface showing up on the host
	start := time.Now()
	for {
		logger := logging.Log().With().Str("prefix", me.Opts.WaitForUdevPrefix).Logger()
		// return an error if we've passed our WaitUdevTimeout
		if me.Opts.WaitForUdev {
			if time.Now().Sub(start) >= me.Opts.WaitForUdevTimeout {
				logger.Error().Str("timeout", me.Opts.WaitForUdevTimeout.String()).Msg("reached udev wait timeout")
				return fmt.Errorf("reached udev wait timeout mac=%s", mac)
			}
		}
		iface, err := me.nw.GetIfaceByMac(mac)
		if err != nil {
			// If we're waiting for udev log an error and retry; otherwise return the error
			if me.Opts.WaitForUdev {
				logger.Error().Str("mac", mac).Err(err).Msg("failed to find mac address")
				time.Sleep(me.Opts.WaitForUdevDelay)
				continue
			} else {
				return fmt.Errorf("failed to find interface by mac %s %w", mac, err)
			}
		}
		logger = logger.With().Str("iface", iface.Name).Str("mac", mac).Logger()

		// test to see if we found a valid prefix
		if me.Opts.WaitForUdev {
			logger.Info().Msg("waiting for interface with valid prefix")
			if strings.HasPrefix(iface.Name, me.Opts.WaitForUdevPrefix) {
				logger.Info().Msg("found interface name matching disallowed udev prefix... waiting")
				time.Sleep(me.Opts.WaitForUdevDelay)
				continue
			}
			logger.Info().Msg("found valid interface name")
		}

		// ensure that eth0 is not used
		if iface.Name == "eth0" {
			return fmt.Errorf("failed to configure interface. eth0 is an invalid name")
		}

		netIface := &NetworkInterface{
			Index:    iface.Index,
			DestName: result.Interfaces[0].Name,
			Address:  &result.IPs[0].Address,
		}

		err = me.nw.Configure(cmd.Netns, netIface)
		if err != nil {
			return fmt.Errorf("failed to configure interface %w", err)
		}
		return err
	}
}
