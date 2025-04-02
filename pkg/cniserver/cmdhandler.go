package cniserver

import (
	"fmt"
	"net"

	"github.com/containernetworking/cni/pkg/types"
	currentcni "github.com/containernetworking/cni/pkg/types/040"
	. "github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

//go:generate moq -pkg mocks -out ../fixtures/mocks/cniserver_mocks.go . CommandHandler

// CommandHandler provides the ability to handle CNI commands
type CommandHandler interface {
	// Add handlers ADD commands
	Add(cmd util.CniCommand) (*currentcni.Result, error)
	// Check handlers DEL commands
	Del(cmd util.CniCommand) error
	// Check handlers CHECK commands (NOT IMPLEMENTED)
	Check(cmd util.CniCommand) error
}

var _ CommandHandler = &commandHandler{}

// NewCniCommandHandler creates a new CommandHandler
func NewCniCommandHandler(pm *openstack.PortManager) *commandHandler {
	return &commandHandler{pm}
}

type commandHandler struct {
	pm *openstack.PortManager
}

func (me *commandHandler) Add(cmd util.CniCommand) (*currentcni.Result, error) {
	context, err := util.NewCniContext(cmd)
	if err != nil {
		return nil, err
	}

	opts := openstack.SetupPortOptsFromContext(context)
	opts.Tags = NewPortTagsFromCommand(cmd).NeutronTags()
	portResult, err := me.pm.SetupPort(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to setup port %w", err)
	}

	return NewCniResult(portResult, cmd)
}

func (me *commandHandler) Del(cmd util.CniCommand) error {
	log := Log().With().Str("cmd", cmd.String()).Logger()
	context, err := util.NewCniContext(cmd)
	if err != nil {
		log.Error().AnErr("err", err).Msg("failed to build context")
		return nil
	}

	opts := openstack.TearDownPortOpts{Hostname: context.Hostname, Tags: NewPortTagsFromCommand(cmd).NeutronTags()}
	if err := me.pm.TeardownPort(opts); err != nil {
		log.Error().Str("hostname", context.Hostname).Str("tags", opts.Tags.String()).AnErr("err", err).Msg("failed to teardown port")
		return nil
	}
	return nil
}

func (me *commandHandler) Check(cmd util.CniCommand) error {
	return nil
}

var ErrIncompletePortResult = fmt.Errorf("Incomplete port result")

// NewCniResult creates a new Result from the combination of a SetupPortResult and CniCommand
func NewCniResult(portResult *openstack.SetupPortResult, cmd util.CniCommand) (*currentcni.Result, error) {
	if portResult.Attachment == nil || portResult.Network == nil ||
		portResult.Port == nil || portResult.Subnet == nil {

		return nil, ErrIncompletePortResult
	}

	ipnet, err := portResult.GetIp()
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
			Dst: NewIpNet(route.DestinationCIDR),
			GW:  net.ParseIP(route.NextHop),
		})
	}

	return result, nil
}

// NewIpNet creates a new IPNet from a string containing an IP and prefix (e.g. "10.1.2.3/24")
func NewIpNet(cidr string) net.IPNet {
	theip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return net.IPNet{IP: theip, Mask: ipnet.Mask}
}
