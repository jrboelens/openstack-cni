package cniserver_test

import (
	"testing"

	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	"github.com/jboelensns/openstack-cni/pkg/cnistate"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"

	. "github.com/pepinns/go-hamcrest"
)

func Test_CreateResultFromPortResult(t *testing.T) {
	t.Run("returns an error with incomplete data", func(t *testing.T) {
		portResult := &openstack.SetupPortResult{}
		cmd := util.CniCommand{}

		result, err := cniserver.CreateResultFromPortResult(portResult, cmd)
		Assert(t).That(result, IsNil())
		Assert(t).That(err, Equals(cniserver.ErrIncompletePortResult))
	})
}

func Test_CmdHandler(t *testing.T) {
	t.Run("can add and delete using a handler", func(t *testing.T) {
		WithCniState(t, func(state cnistate.State) {
			WithTestConfig(t, func(cfg TestingConfig) {
				deps, err := cniserver.NewBuilder().WithState(state).Build()
				Assert(t).That(err, IsNil())

				cmd := NewTestData().CniCommand()
				results, err := deps.CniHandler().Add(cmd)
				Assert(t).That(err, IsNil())
				Assert(t).That(results, Not(IsNil()))

				state.Set(&cnistate.IfaceInfo{
					ContainerId: cmd.ContainerID,
					Ifname:      cmd.IfName,
					Netns:       cmd.Netns,
					IpAddress:   results.IPs[0].Address.IP.String(),
					PodName:     "POD",
					Namespace:   "NAMESPACE",
				})

				// ensure the port exists
				port, err := deps.OpenstackClient().GetPortByIp(results.IPs[0].Address.IP.String())
				Assert(t).That(err, IsNil())
				Assert(t).That(port, Not(IsNil()))

				// issue a delete
				Assert(t).That(deps.CniHandler().Del(cmd), IsNil())

				// ensure the port's gone
				port, perr := deps.OpenstackClient().GetPortByIp(results.IPs[0].Address.IP.String())
				Assert(t).That(perr, Equals(openstack.ErrPortNotFound))
			})
		})
	})
}

func Test_IpNetFromCidr(t *testing.T) {
	t.Run("can parse an IP with prefix and return a proper IPNet", func(t *testing.T) {
		ipnet := cniserver.IpnetFromCidr("1.2.3.4/24")
		Assert(t).That(ipnet.IP, Equals("1.2.3.4"))
		Assert(t).That(ipnet.Mask, Equals("ffffff00"))

	})
}
