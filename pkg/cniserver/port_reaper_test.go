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
	t.Run("skeleton", func(t *testing.T) {
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
			mock.DeletePortFunc = func(portId string) error {
				return nil
			}

			hostname, err := os.Hostname()
			Assert(t).That(err, IsNil())
			opts := cniserver.PortReaperOpts{time.Second * 300, client}
			reaper := cniserver.NewPortReaper(opts)
			err = reaper.Reap(hostname)
			Assert(t).That(err, IsNil())

			Assert(t).That(len(mock.DeletePortCalls()), Equals(1))
		})
	})
}
