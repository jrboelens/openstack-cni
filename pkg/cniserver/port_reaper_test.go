package cniserver_test

import (
	"os"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	"github.com/jboelensns/openstack-cni/pkg/fixtures/mocks"
	"github.com/jboelensns/openstack-cni/pkg/openstack"

	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	. "github.com/pepinns/go-hamcrest"
)

func Test_PortReaper(t *testing.T) {
	hostname, err := os.Hostname()
	Assert(t).That(err, IsNil())
	t.Run("port reaper attempts to delete ports", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			serverId := "myId"
			mock.GetServerByNameFunc = func(name string) (*servers.Server, error) {
				return &servers.Server{ID: serverId}, nil
			}
			mock.GetPortsByDeviceIdFunc = func(deviceId string) ([]ports.Port, error) {
				return []ports.Port{
					{Tags: []string{"foo=bar", "netns=/proc/1234/ns"}},
				}, nil
			}
			mock.DeletePortFunc = func(portId string) error { return nil }

			reaper := cniserver.NewPortReaper(client, PortReaperOpts())
			err = reaper.Reap(hostname)
			Assert(t).That(err, IsNil())

			Assert(t).That(len(mock.DeletePortCalls()), Equals(1))
		})
	})

	t.Run("will not reap a new port", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			reaper := cniserver.NewPortReaper(client, PortReaperOpts())
			port := ports.Port{CreatedAt: time.Now()}
			mock.DeletePortFunc = func(portId string) error { return nil }
			err = reaper.ReapPort(port)
			Assert(t).That(err, IsNil())
			Assert(t).That(len(mock.DeletePortCalls()), Equals(0))
		})
	})

	t.Run("will reap an old port", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			reaper := cniserver.NewPortReaper(client, PortReaperOpts())
			port := ports.Port{CreatedAt: time.Now().Add(-(time.Second * 6000))}
			mock.DeletePortFunc = func(portId string) error { return nil }
			err = reaper.ReapPort(port)
			Assert(t).That(err, IsNil())
			Assert(t).That(len(mock.DeletePortCalls()), Equals(1))
		})
	})
}

func Test_PortReaperIntegration(t *testing.T) {
	t.Run("port reaper attempts to delete ports", func(t *testing.T) {
		WithTestConfig(t, func(cfg TestingConfig) {
			// create a port with a network namespace that doesn't exist for my machine
			cmd := NewTestData().CniCommand()
			context := CniContextFromConfig(t, cfg, cmd)

			realClient, err := openstack.NewOpenstackClient()
			Assert(t).That(err, IsNil())

			client := openstack.NewCachedClient(realClient, time.Second*5)
			pm := openstack.NewPortManager(client)
			opts := openstack.SetupPortOptsFromContext(context)
			opts.Tags = cniserver.NewPortTags(context.Command)

			_, err = pm.SetupPort(opts)
			Assert(t).That(err, IsNil(), "failed to setup port")

			_, err = client.GetPortByTags(opts.Tags.AsStringSlice())
			Assert(t).That(err, IsNil(), "failed get port by tags %s", opts.Tags.String())

			// run the reaper
			reaper := cniserver.NewPortReaper(client, cniserver.PortReaperOpts{})
			//TODO: MOVE host name to an option rather than an argument
			Assert(t).That(reaper.Reap(cfg.Hostname), IsNil())

			// ensure the port is deleted
			_, err = client.GetPortByTags(opts.Tags.AsStringSlice())
			Assert(t).That(err, Equals(openstack.ErrPortNotFound))
		})
	})
}
