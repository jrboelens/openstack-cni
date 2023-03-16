package cniserver

import (
	"fmt"
	"net"

	"github.com/containernetworking/cni/pkg/types"
	currentcni "github.com/containernetworking/cni/pkg/types/040"
	"github.com/jboelensns/openstack-cni/pkg/cnistate"
	. "github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

//go:generate moq -pkg mocks -out ../fixtures/mocks/cniserver_mocks.go . CommandHandler

type CommandHandler interface {
	Add(cmd util.CniCommand) (*currentcni.Result, error)
	Del(cmd util.CniCommand) error
	Check(cmd util.CniCommand) error
}

var _ CommandHandler = &commandHandler{}

func NewCniCommandHandler(pm *openstack.PortManager, state cnistate.State) *commandHandler {
	return &commandHandler{pm, state}
}

type commandHandler struct {
	pm    *openstack.PortManager
	state cnistate.State
}

func (me *commandHandler) Add(cmd util.CniCommand) (*currentcni.Result, error) {
	context, err := util.NewCniContext(cmd)
	if err != nil {
		return nil, err
	}

	portResult, err := me.pm.SetupPort(openstack.SetupPortOptsFromContext(context))
	if err != nil {
		return nil, fmt.Errorf("failed to setup port %w", err)
	}

	return CreateResultFromPortResult(portResult, cmd)
}

func (me *commandHandler) Del(cmd util.CniCommand) error {
	log := Log().With().Str("cmd", cmd.String()).Logger()
	context, err := util.NewCniContext(cmd)
	if err != nil {
		log.Error().AnErr("err", err).Msg("failed to build context")
		return nil
	}

	info, err := me.state.Get(cmd.ContainerID, cmd.IfName)
	if err != nil {
		log.Error().AnErr("err", err).Msg("failed to get cni state")
		return nil
	}
	if info == nil {
		log.Info().Msg("state not found")
		return nil
	}

	if err := me.pm.TeardownPort(openstack.TearDownPortOptsFromContext(context.Hostname, info.IpAddress)); err != nil {
		log.Error().Str("hostname", context.Hostname).Str("ip", info.IpAddress).AnErr("err", err).Msg("failed to teardown port")
		return nil
	}
	return nil
}

func (me *commandHandler) Check(cmd util.CniCommand) error {
	return nil
}

var ErrIncompletePortResult = fmt.Errorf("Incomplete port result")

func CreateResultFromPortResult(portResult *openstack.SetupPortResult, cmd util.CniCommand) (*currentcni.Result, error) {
	if portResult.Attachment == nil || portResult.Network == nil ||
		portResult.Port == nil || portResult.Subnet == nil {

		return nil, ErrIncompletePortResult
	}

	ipnet, err := GetIPFromPortResult(portResult)
	if err != nil {
		return nil, err
	}
	zero := 0

	result := &currentcni.Result{
		CNIVersion: currentcni.ImplementedSpecVersion,
		Interfaces: []*currentcni.Interface{
			{
				Name:    cmd.IfName,
				Mac:     portResult.Attachment.MACAddr,
				Sandbox: cmd.Netns,
			},
		},
		IPs: []*currentcni.IPConfig{
			{
				Interface: &zero,
				Address:   *ipnet,
				Gateway:   net.ParseIP(portResult.Subnet.GatewayIP),
				Version:   "4",
			},
		},
		Routes: make([]*types.Route, 0, 0),
		DNS: types.DNS{
			Nameservers: portResult.Subnet.DNSNameservers,
			// Domain:      "",         // NEED
			// Search:      []string{}, //NEED
			// Options:     []string{}, //NEED
		},
	}

	// add host routes
	for _, route := range portResult.Subnet.HostRoutes {
		result.Routes = append(result.Routes, &types.Route{
			Dst: util.IpnetFromIp(route.DestinationCIDR),
			GW:  net.ParseIP(route.NextHop),
		})
	}

	return result, nil
}

func GetIPFromPortResult(portResult *openstack.SetupPortResult) (*net.IPNet, error) {
	ip := net.ParseIP(portResult.Port.FixedIPs[0].IPAddress)
	_, cidr, err := net.ParseCIDR(portResult.Subnet.CIDR)
	if err != nil {
		return nil, err
	}
	return &net.IPNet{IP: ip, Mask: cidr.Mask}, nil
}
