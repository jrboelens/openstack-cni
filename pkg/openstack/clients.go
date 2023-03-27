package openstack

import (
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

// ApiClients providers easy access to various Openstack API endpoints
type ApiClients struct {
	authOpts       gophercloud.AuthOptions
	region         string
	ComputeClient  *gophercloud.ServiceClient
	IdentityClient *gophercloud.ServiceClient
	NetworkClient  *gophercloud.ServiceClient
	ProviderClient *gophercloud.ProviderClient
}

// NewApiClients creates a new ApiClients based on AuthOptions
func NewApiClients(opts gophercloud.AuthOptions) (*ApiClients, error) {
	region := os.Getenv("OS_REGION_NAME")
	if region == "" {
		region = "RegionOne"
	}

	opts.AllowReauth = true
	clients := &ApiClients{authOpts: opts}
	var err error

	// setup the provider client
	clients.ProviderClient, err = openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	// setup the compute / nova client
	clients.ComputeClient, err = openstack.NewComputeV2(clients.ProviderClient, gophercloud.EndpointOpts{
		Name: "nova",
		// Type:   "compute",
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	// setup the identity / keystone client
	clients.IdentityClient, err = openstack.NewIdentityV3(clients.ProviderClient, gophercloud.EndpointOpts{
		Name: "keystone",
		// Type:   "identity",
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	// setup the network/neutron client
	clients.NetworkClient, err = openstack.NewNetworkV2(clients.ProviderClient, gophercloud.EndpointOpts{
		Name: "neutron",
		// Type:   "network",
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	return clients, nil
}
