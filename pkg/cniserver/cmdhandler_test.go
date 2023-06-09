package cniserver_test

import (
	"testing"

	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"

	. "github.com/pepinns/go-hamcrest"
)

func Test_CreateResultFromPortResult(t *testing.T) {
	t.Run("returns an error with incomplete data", func(t *testing.T) {
		portResult := &openstack.SetupPortResult{}
		cmd := util.CniCommand{}

		result, err := cniserver.NewCniResult(portResult, cmd)
		Assert(t).That(result, IsNil())
		Assert(t).That(err, Equals(cniserver.ErrIncompletePortResult))
	})
}

func Test_CmdHandler(t *testing.T) {
	t.Run("can add and delete using a handler", func(t *testing.T) {
		WithTestConfig(t, func(cfg TestingConfig) {
			deps, err := cniserver.NewBuilder(cniserver.NewConfig()).Build()
			Assert(t).That(err, IsNil())

			cmd := NewTestData().CniCommand()
			results, err := deps.CniHandler().Add(cmd)
			Assert(t).That(err, IsNil())
			Assert(t).That(results, Not(IsNil()))

			// ensure the port exists
			port, err := deps.OpenstackClient().GetPortByTags(cniserver.NewPortTags(cmd).AsStringSlice())
			Assert(t).That(err, IsNil())
			Assert(t).That(port, Not(IsNil()))

			// issue a delete
			Assert(t).That(deps.CniHandler().Del(cmd), IsNil())

			// ensure the port's gone
			port, perr := deps.OpenstackClient().GetPortByTags(cniserver.NewPortTags(cmd).AsStringSlice())
			Assert(t).That(perr, Equals(openstack.ErrPortNotFound))
		})
	})
}

func Test_IpNetFromCidr(t *testing.T) {
	t.Run("can parse an IP with prefix and return a proper IPNet", func(t *testing.T) {
		ipnet := cniserver.NewIpNet("1.2.3.4/24")
		Assert(t).That(ipnet.IP, Equals("1.2.3.4"))
		Assert(t).That(ipnet.Mask, Equals("ffffff00"))
	})
}

func Test_NewPortTags(t *testing.T) {
	t.Run("ensure container ids are shortened to the first 12 chars", func(t *testing.T) {
		cmd := NewTestData().CniCommand()
		cmd.ContainerID = "12345678901234567890"
		tags := cniserver.NewPortTags(cmd)
		Assert(t).That(tags.Tags[0], Equals("containerid=123456789012"))
	})

	t.Run("ensure a short container id doesn't fail.. this shouldn't be possible in k8s but we'll play it safe", func(t *testing.T) {
		cmd := NewTestData().CniCommand()
		cmd.ContainerID = "123"
		tags := cniserver.NewPortTags(cmd)
		Assert(t).That(tags.Tags[0], Equals("containerid=123"))
	})
}
