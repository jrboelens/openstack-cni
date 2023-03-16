package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
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
	result := &SetupPortResult{}
	var err error

	// look up the server
	result.Server, err = me.client.GetServerByName(opts.Hostname)
	if err != nil {
		return result, err
	}
	if result.Server == nil {
		return result, fmt.Errorf("failed to find server named %s", opts.Hostname)
	}

	// find the network by name
	result.Network, err = me.client.GetNetworkByName(opts.NetworkName)
	if err != nil {
		return result, err
	}
	if result.Network == nil {
		return result, fmt.Errorf("failed to find network named %s", opts.NetworkName)
	}

	if len(opts.SecurityGroups) > 0 {
		projectId := ""
		// we need the projectId in order to look up the security groups
		if len(opts.ProjectName) > 0 {
			project, err := me.client.GetProjectByName(opts.ProjectName)
			if err != nil {
				return nil, err
			}
			if project == nil {
				return result, fmt.Errorf("failed to find project named %s", opts.ProjectName)
			}
			projectId = project.ID
		}

		// if security groups were specified, look them up
		sgIds := make([]string, len(opts.SecurityGroups), len(opts.SecurityGroups))
		for i, sgName := range opts.SecurityGroups {
			sg, err := me.client.GetSecurityGroupByName(sgName, projectId)
			if err != nil {
				return nil, fmt.Errorf("failed to lookup security group named %s", sgName)
			}
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
		subnet, err := me.client.GetSubnetByName(opts.SubnetName, portOpts.NetworkID)
		if err != nil {
			return result, err
		}
		portOpts.FixedIPs = []FixedIP{{SubnetID: subnet.ID}}
	}

	result.Port, err = me.client.CreatePort(portOpts)
	if err != nil {
		return result, err
	}
	if result.Port == nil {
		return result, fmt.Errorf("failed to find port named %s", opts.PortName)
	}

	// lookup the subnet that the port came from
	result.Subnet, err = me.client.GetSubnet(result.Port.FixedIPs[0].SubnetID)
	if err != nil {
		return result, err
	}
	if result.Subnet == nil {
		return result, fmt.Errorf("failed to find subnet with ID %s", result.Port.FixedIPs[0].SubnetID)
	}

	if !opts.SkipPortAttach {
		// assign the port to the VM
		result.Attachment, err = me.client.AssignPort(result.Port.ID, result.Server.ID)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func (me *PortManager) TeardownPort(opts TearDownPortOpts) error {

	// lookup port by ip
	port, err := me.client.GetPortByIp(opts.IpAddress)
	if err != nil {
		return err
	}
	if port == nil {
		return fmt.Errorf("failed to find port by IP Address %s", opts.IpAddress)
	}

	if !opts.SkipPortDetach {
		// look up the server
		server, err := me.client.GetServerByName(opts.Hostname)
		if err != nil {
			return err
		}

		if server == nil {
			return fmt.Errorf("failed to find server by name %s", opts.Hostname)
		}

		err = me.client.DetachPort(port.ID, server.ID)
		if err != nil {
			return err
		}
	}

	return me.client.DeletePort(port.ID)
}

type SetupPortOpts struct {
	Hostname       string
	NetworkName    string
	PortName       string
	ProjectName    string
	SecurityGroups []string
	SubnetName     string
	SkipPortAttach bool
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
