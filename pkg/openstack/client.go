package openstack

import (
	"fmt"
	"strings"

	gc_os "github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

//go:generate moq -pkg mocks -out ../fixtures/mocks/openstack_mocks.go . OpenstackClient

type OpenstackClient interface {
	AssignPort(portId, serverId string) (*attachinterfaces.Interface, error)
	CreatePort(opts ports.CreateOpts, extraOpts *ExtraCreatePortOpts) (*ports.Port, error)
	DeletePort(portId string) error
	DetachPort(portId, serverId string) error
	Clients() *ApiClients
	GetNetworkByName(name string) (*networks.Network, error)
	GetPort(portId string) (*ports.Port, error)
	GetPortsByDeviceId(deviceId string) ([]ports.Port, error)
	GetPortByTags(tags []string) (*ports.Port, error)
	GetPortsByTags(tags []string) ([]ports.Port, error)
	GetProjectByName(name string) (*projects.Project, error)
	GetServerByName(name string) (*servers.Server, error)
	GetSecurityGroupByName(name, projectId string) (*groups.SecGroup, error)
	GetSubnet(id string) (*subnets.Subnet, error)
	GetSubnetByName(name, networkId string) (*subnets.Subnet, error)
}

// openstackClient exposes various Openstack API functionality in a single location
type openstackClient struct {
	clients *ApiClients
}

var _ OpenstackClient = &openstackClient{}

// NewOpenstackClient creates a client for interacting with Openstack APIs
func NewOpenstackClient() (*openstackClient, error) {
	if err := util.ReadConfigIntoEnv(); err != nil {
		return nil, err
	}

	authOpts, err := gc_os.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	apiClients, err := NewApiClients(authOpts)
	if err != nil {
		return nil, err
	}

	return &openstackClient{apiClients}, nil
}

// AssignPort attaches a port to a server
func (me *openstackClient) AssignPort(portId, serverId string) (*attachinterfaces.Interface, error) {
	opts := attachinterfaces.CreateOpts{PortID: portId}
	result := attachinterfaces.Create(me.clients.ComputeClient, serverId, opts)
	return result.Extract()
}

func (me *openstackClient) Clients() *ApiClients {
	return me.clients
}

// CreatePort creates a neutron port inside of the specified network
func (me *openstackClient) CreatePort(opts ports.CreateOpts, extraOpts *ExtraCreatePortOpts) (*ports.Port, error) {
	// optionally include port security
	var finalOpts ports.CreateOptsBuilder
	finalOpts = opts

	if extraOpts != nil {
		// See: https://github.com/gophercloud/gophercloud/blob/main/openstack/networking/v2/extensions/portsecurity/doc.go
		if extraOpts.HasPortSecurity() {
			finalOpts = portsecurity.PortCreateOptsExt{
				CreateOptsBuilder:   finalOpts,
				PortSecurityEnabled: extraOpts.PortSecurityEnabled,
			}
		}
		// https://github.com/gophercloud/gophercloud/blob/v1/openstack/networking/v2/extensions/portsbinding/doc.go
		if extraOpts.HasPortBindings() {
			finalOpts = portsbinding.CreateOptsExt{
				VNICType: extraOpts.VNICType,
				HostID:   extraOpts.HostID,
				Profile:  extraOpts.Profile,
			}
		}

	}
	return ports.Create(me.clients.NetworkClient, finalOpts).Extract()
}

// DeletePort deletes the port
func (me *openstackClient) DeletePort(portId string) error {
	result := ports.Delete(me.clients.NetworkClient, portId)
	return result.ExtractErr()
}

// Detach port removes a port's relationship from a server
func (me *openstackClient) DetachPort(portId, serverId string) error {
	result := attachinterfaces.Delete(me.clients.ComputeClient, serverId, portId)
	return result.ExtractErr()
}

var ErrServerNotFound = fmt.Errorf("server not found")

// GetServer returns a single server based on a server name
func (me *openstackClient) GetServerByName(name string) (*servers.Server, error) {
	listOpts := servers.ListOpts{Name: regexName(name), Limit: 1}
	allPages, err := servers.List(me.clients.ComputeClient, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		return nil, err
	}

	if len(allServers) == 0 {
		return nil, ErrServerNotFound
	}
	return &allServers[0], nil
}

var ErrNetworkNotFound = fmt.Errorf("network not found")

// GetServer returns a single network based on a network name
func (me *openstackClient) GetNetworkByName(name string) (*networks.Network, error) {
	listOpts := networks.ListOpts{Name: name, Limit: 1}
	allPages, err := networks.List(me.clients.NetworkClient, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return nil, err
	}

	if len(allNetworks) == 0 {
		return nil, ErrNetworkNotFound
	}
	return &allNetworks[0], nil
}

// GetPort returns a single port based on an ID
func (me *openstackClient) GetPort(portId string) (*ports.Port, error) {
	result := ports.Get(me.clients.NetworkClient, portId)
	return result.Extract()
}

var ErrPortNotFound = fmt.Errorf("port not found")

// GetPortByTags returns a single port based on matching tags
func (me *openstackClient) GetPortByTags(tags []string) (*ports.Port, error) {
	tagsStr := strings.Join(tags, ",")
	listOpts := ports.ListOpts{Tags: tagsStr}
	return me.getPort(listOpts)
}

// GetPortsByTags returns a single port based on matching tags
func (me *openstackClient) GetPortsByTags(tags []string) ([]ports.Port, error) {
	tagsStr := strings.Join(tags, ",")
	listOpts := ports.ListOpts{Tags: tagsStr}
	return me.getPorts(listOpts)
}

// GetPortsByDeviceId returns all ports based device id (server id)
func (me *openstackClient) GetPortsByDeviceId(deviceId string) ([]ports.Port, error) {
	listOpts := ports.ListOpts{DeviceID: deviceId}
	allPages, err := ports.List(me.clients.NetworkClient, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	return ports.ExtractPorts(allPages)
}

func (me *openstackClient) getPorts(listOpts ports.ListOpts) ([]ports.Port, error) {
	allPages, err := ports.List(me.clients.NetworkClient, listOpts).AllPages()
	if err != nil {
		return nil, err
	}
	return ports.ExtractPorts(allPages)
}

func (me *openstackClient) getPort(listOpts ports.ListOpts) (*ports.Port, error) {
	ports, err := me.getPorts(listOpts)
	if err != nil {
		return nil, err
	}

	if len(ports) == 0 {
		return nil, ErrPortNotFound
	}
	return &ports[0], nil
}

var ErrProjectNotFound = fmt.Errorf("project not found")

// GetProjectByName returns a project based on name
func (me *openstackClient) GetProjectByName(name string) (*projects.Project, error) {
	listOpts := projects.ListOpts{Name: name}

	allPages, err := projects.List(me.clients.IdentityClient, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		return nil, err
	}

	if len(allProjects) == 0 {
		return nil, ErrProjectNotFound
	}
	return &allProjects[0], nil
}

var ErrSecurityGroupNotFound = fmt.Errorf("security group not found")

// GetSecurityGroupByName returns a single port based on an IpAddress
func (me *openstackClient) GetSecurityGroupByName(name, projectId string) (*groups.SecGroup, error) {
	listOpts := groups.ListOpts{Name: name, ProjectID: projectId}

	allPages, err := groups.List(me.clients.NetworkClient, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		return nil, err
	}

	if len(allGroups) == 0 {
		return nil, ErrSecurityGroupNotFound
	}
	return &allGroups[0], nil
}

// GetSubnet return a single subnet based on a subnet UUID
func (me *openstackClient) GetSubnet(id string) (*subnets.Subnet, error) {
	result := subnets.Get(me.clients.NetworkClient, id)
	return result.Extract()
}

var ErrSubnetNotFound = fmt.Errorf("subnet not found")

// GetSubnetByName returns a project based on name
func (me *openstackClient) GetSubnetByName(name, networkId string) (*subnets.Subnet, error) {
	listOpts := subnets.ListOpts{Name: name, NetworkID: networkId}

	allPages, err := subnets.List(me.clients.NetworkClient, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := subnets.ExtractSubnets(allPages)
	if err != nil {
		return nil, err
	}

	if len(all) == 0 {
		return nil, ErrSubnetNotFound
	}
	return &all[0], nil
}

type FixedIP struct {
	SubnetID  string `json:"subnet_id,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
}

// regexName turns a name into a regex that matches ONLY the name
// some of the Openstack APIs use a regex as name
// this allows for easy exact matching
func regexName(name string) string {
	return fmt.Sprintf("^%s$", name)
}

type ExtraCreatePortOpts struct {
	// PortSecurityEnabled toggles port security on a port.
	PortSecurityEnabled *bool `json:"port_security_enabled,omitempty"`

	// The ID of the host where the port is allocated.
	HostID string `json:"binding:host_id,omitempty"`

	// The virtual network interface card (vNIC) type that is bound to the
	// neutron port.
	VNICType string `json:"binding:vnic_type,omitempty"`

	// A dictionary that enables the application running on the specified
	// host to pass and receive virtual network interface (VIF) port-specific
	// information to the plug-in.
	Profile map[string]interface{} `json:"binding:profile,omitempty"`
}

func (me *ExtraCreatePortOpts) HasPortSecurity() bool {
	return me.PortSecurityEnabled != nil
}

func (me *ExtraCreatePortOpts) HasPortBindings() bool {
	return me.HostID != "" || me.VNICType != "" || len(me.Profile) > 0
}
