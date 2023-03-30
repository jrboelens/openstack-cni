package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	. "github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

func NewPortManager(client OpenstackClient) *PortManager {
	return &PortManager{client}
}

// PortManager provides the ability to execute various compound port actions
type PortManager struct {
	client OpenstackClient
}

// SetupPort creates a new port and assigns it to a server
func (me *PortManager) SetupPort(opts SetupPortOpts) (*SetupPortResult, error) {
	log := Log().With().Str("command", "ADD").Str("hostname", opts.Hostname).Str("networkName", opts.NetworkName).Str("projectName", opts.ProjectName).Str("portName", opts.PortName).Logger()
	result := &SetupPortResult{}
	var err error

	// look up the server
	log.Info().Msg("looking up server")
	result.Server, err = me.client.GetServerByName(opts.Hostname)
	if err != nil {
		return result, err
	}
	if result.Server == nil {
		return result, fmt.Errorf("failed to find server named %s", opts.Hostname)
	}
	log.Info().Msg("found server")

	// find the network by name
	log.Info().Msg("looking up network")
	result.Network, err = me.client.GetNetworkByName(opts.NetworkName)
	if err != nil {
		return result, err
	}
	if result.Network == nil {
		return result, fmt.Errorf("failed to find network named %s", opts.NetworkName)
	}
	log.Info().Msg("found network")

	if len(opts.SecurityGroups) > 0 {
		projectId := ""
		// we need the projectId in order to look up the security groups
		if len(opts.ProjectName) > 0 {
			log.Info().Msg("looking up project")
			project, err := me.client.GetProjectByName(opts.ProjectName)
			if err != nil {
				return nil, err
			}
			if project == nil {
				return result, fmt.Errorf("failed to find project named %s", opts.ProjectName)
			}
			log.Info().Msg("found project")
			projectId = project.ID
		}

		// if security groups were specified, look them up
		sgIds := make([]string, len(opts.SecurityGroups), len(opts.SecurityGroups))
		for i, sgName := range opts.SecurityGroups {
			log.Info().Str("sgName", sgName).Msg("looking up security group")
			sg, err := me.client.GetSecurityGroupByName(sgName, projectId)
			if err != nil {
				return nil, fmt.Errorf("failed to lookup security group named %s", sgName)
			}
			if sg == nil {
				return result, fmt.Errorf("failed to find security group named %s", sgName)
			}
			log.Info().Str("sgName", sgName).Msg("found security group")
			sgIds[i] = sg.ID
		}
		opts.SecurityGroups = sgIds
	}

	// create a port
	t := true
	portOpts := ports.CreateOpts{
		NetworkID:      result.Network.ID,
		Name:           opts.PortName,
		AdminStateUp:   &t,
		SecurityGroups: &opts.SecurityGroups,
	}

	// optionally include the subnet when creating the port
	if opts.SubnetName != "" {
		log.Info().Str("subnetName", opts.SubnetName).Str("networkId", portOpts.NetworkID).Msg("looking up subnet")
		subnet, err := me.client.GetSubnetByName(opts.SubnetName, portOpts.NetworkID)
		if err != nil {
			return result, err
		}
		if subnet == nil {
			return result, fmt.Errorf("failed to find subnet named %s in network %s", opts.SubnetName, portOpts.NetworkID)
		}
		log.Info().Str("subnetName", opts.SubnetName).Str("networkId", portOpts.NetworkID).Msg("found subnet")
		portOpts.FixedIPs = []FixedIP{{SubnetID: subnet.ID}}
	}

	log.Info().Msg("creating port")
	result.Port, err = me.client.CreatePort(portOpts)
	if err != nil {
		return result, err
	}
	log.Info().Msg("created port")

	// add tags to the port
	log.Info().Msg("adding tags to port")
	if len(opts.Tags.Tags) > 0 {
		tagger := NewNeutronTagger(me.client.Clients().NetworkClient, Ports)
		if err := tagger.SetAll(result.Port.ID, opts.Tags); err != nil {
			return result, err
		}
		log.Info().Msg("added tags to port")
	}

	// lookup the subnet that the port came from
	log.Info().Str("subnetId", result.Port.FixedIPs[0].SubnetID).Msg("looking up subnet by id")
	result.Subnet, err = me.client.GetSubnet(result.Port.FixedIPs[0].SubnetID)
	if err != nil {
		return result, err
	}
	if result.Subnet == nil {
		return result, fmt.Errorf("failed to find subnet with ID %s", result.Port.FixedIPs[0].SubnetID)
	}
	log.Info().Msg("found subnet by id")

	if !opts.SkipPortAttach {
		// assign the port to the VM
		log.Info().Str("portId", result.Port.ID).Str("serverId", result.Server.ID).Msg("assigning port to server")
		result.Attachment, err = me.client.AssignPort(result.Port.ID, result.Server.ID)
		if err != nil {
			return result, err
		}
		log.Info().Str("portId", result.Port.ID).Str("serverId", result.Server.ID).Msg("assigned port to server")
	}

	return result, nil
}

func (me *PortManager) TeardownPort(opts TearDownPortOpts) error {
	log := Log().With().Str("command", "DEL").Str("hostname", opts.Hostname).Str("ipaddress", opts.IpAddress).Logger()

	// lookup port by ip
	log.Info().Msg("looking up port by ipaddress")
	port, err := me.client.GetPortByIp(opts.IpAddress)
	if err != nil {
		return err
	}
	if port == nil {
		return fmt.Errorf("failed to find port by IP Address %s", opts.IpAddress)
	}
	log.Info().Msg("found port by ipaddress")

	if !opts.SkipPortDetach {
		// look up the server
		log.Info().Msg("looking up server")
		server, err := me.client.GetServerByName(opts.Hostname)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("failed to find server by name %s", opts.Hostname)
		}
		log.Info().Msg("found server")

		log.Info().Str("portId", port.ID).Str("serverId", server.ID).Msg("detaching port")
		err = me.client.DetachPort(port.ID, server.ID)
		if err != nil {
			return err
		}
		log.Info().Str("portId", port.ID).Str("serverId", server.ID).Msg("detached port")
	}

	log.Info().Str("portId", port.ID).Msg("deleting port")
	if err := me.client.DeletePort(port.ID); err != nil {
		return err
	}
	log.Info().Str("portId", port.ID).Msg("deleted port")
	return nil
}

type SetupPortOpts struct {
	Hostname       string
	NetworkName    string
	PortName       string
	ProjectName    string
	SecurityGroups []string
	SubnetName     string
	SkipPortAttach bool
	Tags           NeutronTags
}

func SetupPortOptsFromContext(context util.CniContext) SetupPortOpts {
	return SetupPortOpts{
		Hostname:       context.Hostname,
		NetworkName:    context.CniConfig.Network,
		PortName:       context.CniConfig.PortName,
		ProjectName:    context.CniConfig.ProjectName,
		SecurityGroups: context.CniConfig.SecurityGroups,
		SubnetName:     context.CniConfig.SubnetName,
	}
}

type TearDownPortOpts struct {
	IpAddress      string
	Hostname       string
	SkipPortDetach bool
}

func TearDownPortOptsFromContext(hostname, ipAddress string) TearDownPortOpts {
	return TearDownPortOpts{
		IpAddress: ipAddress,
		Hostname:  hostname,
	}
}

// SetupPortResult contains information gathered while setting up a port
type SetupPortResult struct {
	Server     *servers.Server
	Network    *networks.Network
	Subnet     *subnets.Subnet
	Port       *ports.Port
	Attachment *attachinterfaces.Interface
}
