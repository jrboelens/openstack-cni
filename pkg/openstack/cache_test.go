package openstack_test

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	"github.com/jboelensns/openstack-cni/pkg/fixtures/mocks"

	"github.com/jboelensns/openstack-cni/pkg/openstack"
	. "github.com/pepinns/go-hamcrest"
)

func Test_Cache(t *testing.T) {

	t.Run("items are cached and can expire", func(t *testing.T) {
		expiry := time.Millisecond * 25
		WithMockClientWithExpiry(t, expiry, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			networkName := "my-network"
			mock.GetNetworkByNameFunc = func(name string) (*networks.Network, error) {
				return &networks.Network{Name: name}, nil
			}

			Assert(t).That(mock.GetNetworkByNameCalls(), HasLen(0))
			invoke := func(calls int) {
				network, err := client.GetNetworkByName(networkName)
				Assert(t).That(err, IsNil())
				Assert(t).That(network.Name, Equals(networkName))
				Assert(t).That(mock.GetNetworkByNameCalls(), HasLen(calls))
			}

			// call back to back in order to ensure the cache is hit
			invoke(1)
			invoke(1)
			// wait longer than the expiry and make sure the cache is missed
			time.Sleep(expiry * 2)
			invoke(2)
		})

	})

	t.Run("GetNetworkByName is cached", func(t *testing.T) {
		// Covered in "items are cached and can expire"
	})

	t.Run("GetPort is cached", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			portId := "dead"
			mock.GetPortFunc = func(portId string) (*ports.Port, error) {
				return &ports.Port{ID: portId}, nil
			}

			Assert(t).That(mock.GetPortCalls(), HasLen(0))
			invoke := func(calls int) {
				port, err := client.GetPort(portId)
				Assert(t).That(err, IsNil())
				Assert(t).That(port.ID, Equals(portId))
				Assert(t).That(mock.GetPortCalls(), HasLen(calls))
			}

			invoke(1)
			invoke(1)
		})
	})

	t.Run("GetPortByTags is cached", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			tags := []string{"foo=bar", "this=that"}
			mock.GetPortByIpFunc = func(ip string) (*ports.Port, error) {
				return &ports.Port{
					Tags: tags,
				}, nil
			}

			Assert(t).That(mock.GetPortByIpCalls(), HasLen(0))
			invoke := func(calls int) {
				port, err := client.GetPortByTags(tags)
				Assert(t).That(err, IsNil())
				Assert(t).That(port.Tags, Equals(tags))
				Assert(t).That(mock.GetPortByIpCalls(), HasLen(calls))
			}

			invoke(1)
			invoke(1)
		})
	})

	t.Run("GetProjectByName is cached", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			name := "myname"
			mock.GetProjectByNameFunc = func(name string) (*projects.Project, error) {
				return &projects.Project{Name: name}, nil
			}

			Assert(t).That(mock.GetProjectByNameCalls(), HasLen(0))
			invoke := func(calls int) {
				project, err := client.GetProjectByName(name)
				Assert(t).That(err, IsNil())
				Assert(t).That(project.Name, Equals(name))
				Assert(t).That(mock.GetProjectByNameCalls(), HasLen(calls))
			}

			invoke(1)
			invoke(1)
		})
	})

	t.Run("GetServerByName is cached", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			name := "myname"
			mock.GetServerByNameFunc = func(name string) (*servers.Server, error) {
				return &servers.Server{Name: name}, nil
			}

			Assert(t).That(mock.GetServerByNameCalls(), HasLen(0))
			invoke := func(calls int) {
				server, err := client.GetServerByName(name)
				Assert(t).That(err, IsNil())
				Assert(t).That(server.Name, Equals(name))
				Assert(t).That(mock.GetServerByNameCalls(), HasLen(calls))
			}

			invoke(1)
			invoke(1)
		})
	})

	t.Run("GetSecurityGroupByName is cached", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			name := "myname"
			projectId := "projId"
			mock.GetSecurityGroupByNameFunc = func(name, projectId string) (*groups.SecGroup, error) {
				return &groups.SecGroup{Name: name}, nil
			}

			Assert(t).That(mock.GetSecurityGroupByNameCalls(), HasLen(0))
			invoke := func(calls int) {
				sg, err := client.GetSecurityGroupByName(name, projectId)
				Assert(t).That(err, IsNil())
				Assert(t).That(sg.Name, Equals(name))
				Assert(t).That(mock.GetSecurityGroupByNameCalls(), HasLen(calls))
			}

			invoke(1)
			invoke(1)
		})
	})

	t.Run("GetSubnet is cached", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			id := "subnetId"
			mock.GetSubnetFunc = func(id string) (*subnets.Subnet, error) {
				return &subnets.Subnet{ID: id}, nil
			}

			Assert(t).That(mock.GetSubnetCalls(), HasLen(0))
			invoke := func(calls int) {
				subnet, err := client.GetSubnet(id)
				Assert(t).That(err, IsNil())
				Assert(t).That(subnet.ID, Equals(id))
				Assert(t).That(mock.GetSubnetCalls(), HasLen(calls))
			}

			invoke(1)
			invoke(1)
		})
	})

	t.Run("GetSubnetByName is cached", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			name := "myname"
			id := "networkid"
			mock.GetSubnetByNameFunc = func(name, networkId string) (*subnets.Subnet, error) {
				return &subnets.Subnet{Name: name}, nil
			}

			Assert(t).That(mock.GetSubnetByNameCalls(), HasLen(0))
			invoke := func(calls int) {
				server, err := client.GetSubnetByName(name, id)
				Assert(t).That(err, IsNil())
				Assert(t).That(server.Name, Equals(name))
				Assert(t).That(mock.GetSubnetByNameCalls(), HasLen(calls))
			}

			invoke(1)
			invoke(1)
		})
	})
}
