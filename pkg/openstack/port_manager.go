package openstack

import (
	"fmt"
	"net"

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

	if opts.SecurityGroups != nil && len(*opts.SecurityGroups) > 0 {
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
		sgIds := make([]string, len(*opts.SecurityGroups), len(*opts.SecurityGroups))
		for i, sgName := range *opts.SecurityGroups {
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
		opts.SecurityGroups = &sgIds
	}

	// create a port
	portOpts := me.setupPortOpts(opts, result)

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
	// account for non-default port create options
	extraCreateOpts := opts.CreateExtraPortOpts()
	result.Port, err = me.client.CreatePort(portOpts, &extraCreateOpts)
	if err != nil {
		return result, err
	}
	log = log.With().Str("port_id", result.Port.ID).Logger()
	log.Info().Msg("created port")

	// add tags to the port
	log.Info().Str("tags", opts.Tags.String()).Msg("adding tags to port")
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

func (me *PortManager) setupPortOpts(opts SetupPortOpts, result *SetupPortResult) ports.CreateOpts {
	t := true
	portOpts := ports.CreateOpts{
		Description:    opts.PortDescription,
		DeviceID:       opts.DeviceId,
		DeviceOwner:    opts.DeviceOwner,
		MACAddress:     opts.MacAddress,
		Name:           opts.PortName,
		NetworkID:      result.Network.ID,
		SecurityGroups: opts.SecurityGroups,
		TenantID:       opts.TenantId,
		ValueSpecs:     opts.ValueSpecs,
	}
	if opts.AdminStateUp != nil {
		portOpts.AdminStateUp = opts.AdminStateUp
	} else {
		portOpts.AdminStateUp = &t
	}
	if opts.AllowedAddressPairs != nil {
		portOpts.AllowedAddressPairs = make([]ports.AddressPair, len(opts.AllowedAddressPairs), len(opts.AllowedAddressPairs))
		for i := range opts.AllowedAddressPairs {
			portOpts.AllowedAddressPairs[i] = ports.AddressPair{
				IPAddress:  opts.AllowedAddressPairs[i].IpAddress,
				MACAddress: opts.AllowedAddressPairs[i].MacAddress,
			}
		}
	}
	return portOpts
}

func (me *PortManager) TeardownPort(opts TearDownPortOpts) error {
	log := Log().With().Str("command", "DEL").Str("hostname", opts.Hostname).Str("tags", opts.Tags.String()).Logger()

	// lookup port by tags
	log.Info().Msg("looking up port by tags")
	port, err := me.client.GetPortByTags(opts.Tags.AsStringSlice())
	if err != nil {
		return err
	}
	if port == nil {
		return fmt.Errorf("failed to find port by tags %s", opts.Tags)
	}
	log = log.With().Str("port_id", port.ID).Logger()
	log.Info().Msg("found port by tags")

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
		log = log.With().Str("server_id", server.ID).Logger()
		log.Info().Msg("found server")

		log.Info().Msg("detaching port")
		err = me.client.DetachPort(port.ID, server.ID)
		if err != nil {
			return err
		}
		log.Info().Msg("detached port")
	}

	log.Info().Msg("deleting port")
	if err := me.client.DeletePort(port.ID); err != nil {
		return err
	}
	log.Info().Msg("deleted port")
	return nil
}

type SetupPortOpts struct {
	AdminStateUp        *bool
	AllowedAddressPairs []util.AddressPair
	DeviceId            string
	DeviceOwner         string
	Hostname            string
	MacAddress          string
	NetworkName         string
	PortDescription     string
	PortName            string
	ProjectName         string
	SecurityGroups      *[]string
	SubnetName          string
	SkipPortAttach      bool
	Tags                NeutronTags
	TenantId            string
	ValueSpecs          *map[string]string
	// extra options
	PortSecurityEnabled *bool
	HostID              string
	VNICType            string
	Profile             map[string]interface{}
}

func (me *SetupPortOpts) CreateExtraPortOpts() ExtraCreatePortOpts {
	return ExtraCreatePortOpts{
		PortSecurityEnabled: me.PortSecurityEnabled,
		HostID:              me.HostID,
		VNICType:            me.VNICType,
		Profile:             me.Profile,
	}
}

func SetupPortOptsFromContext(context util.CniContext) SetupPortOpts {
	return SetupPortOpts{
		AdminStateUp:        context.CniConfig.AdminStateUp,
		AllowedAddressPairs: context.CniConfig.AllowedAddressPairs,
		DeviceId:            context.CniConfig.DeviceId,
		DeviceOwner:         context.CniConfig.DeviceOwner,
		Hostname:            context.Hostname,
		MacAddress:          context.CniConfig.MacAddress,
		NetworkName:         context.CniConfig.Network,
		PortDescription:     context.CniConfig.PortDescription,
		PortName:            context.CniConfig.PortName,
		ProjectName:         context.CniConfig.ProjectName,
		SecurityGroups:      context.CniConfig.SecurityGroups,
		SubnetName:          context.CniConfig.SubnetName,
		TenantId:            context.CniConfig.TenantId,
		ValueSpecs:          context.CniConfig.ValueSpecs,
		PortSecurityEnabled: context.CniConfig.PortSecurityEnabled,
		HostID:              context.CniConfig.HostID,
		VNICType:            context.CniConfig.VNICType,
		Profile:             context.CniConfig.Profile,
	}
}

type TearDownPortOpts struct {
	Hostname       string
	Tags           NeutronTags
	SkipPortDetach bool
}

// SetupPortResult contains information gathered while setting up a port
type SetupPortResult struct {
	Server     *servers.Server
	Network    *networks.Network
	Subnet     *subnets.Subnet
	Port       *ports.Port
	Attachment *attachinterfaces.Interface
}

// GetIp returns an IPNet created from teh first FixedIP
func (me *SetupPortResult) GetIp() (*net.IPNet, error) {
	ip := net.ParseIP(me.Port.FixedIPs[0].IPAddress)
	_, cidr, err := net.ParseCIDR(me.Subnet.CIDR)
	if err != nil {
		return nil, err
	}
	return &net.IPNet{IP: ip, Mask: cidr.Mask}, nil
}
