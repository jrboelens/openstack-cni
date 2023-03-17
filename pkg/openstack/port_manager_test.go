package openstack_test

import (
	"testing"
	"time"

	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"

	. "github.com/pepinns/go-hamcrest"
)

func Test_PortManager(t *testing.T) {
	WithTestConfig(t, func(cfg TestingConfig) {
		cmd := util.CniCommand{StdinData: []byte("{}")}
		context := CniContextFromConfig(t, cfg, cmd)

		realClient, err := openstack.NewOpenstackClient()
		Assert(t).That(err, IsNil())

		client := openstack.NewCachedClient(realClient, time.Second*5)

		t.Run("can setup a port and tear it down", func(t *testing.T) {
			SetupAndTeardownPort(t, context, client)
		})

		t.Run("can setup a port with a subnet and tear it down", func(t *testing.T) {
			context.CniConfig.SubnetName = cfg.SubnetName
			SetupAndTeardownPort(t, context, client)
		})
	})
}

func SetupAndTeardownPort(t *testing.T, context util.CniContext, client openstack.OpenstackClient) {
	pm := openstack.NewPortManager(client)
	opts := openstack.SetupPortOptsFromContext(context)

	results, err := pm.SetupPort(opts)
	Assert(t).That(err, IsNil(), "failed to setup port")
	ipAddress := results.Attachment.FixedIPs[0].IPAddress

	_, err = client.GetPortByIp(ipAddress)
	Assert(t).That(err, IsNil(), "failed get port by ip %s", ipAddress)

	err = pm.TeardownPort(openstack.TearDownPortOptsFromContext(context.Hostname, ipAddress))
	Assert(t).That(err, IsNil(), "failed teardown port")

	_, err = client.GetPort(results.Port.ID)
	if err == nil {
		t.Errorf("expected port to be gone %s", ipAddress)
	}

	results, err = pm.SetupPort(opts)
	Assert(t).That(err, IsNil(), "failed to setup port")
	ipAddress = results.Attachment.FixedIPs[0].IPAddress

	_, err = client.GetPortByIp(ipAddress)
	Assert(t).That(err, IsNil(), "failed get port by ip %s", ipAddress)

	err = pm.TeardownPort(openstack.TearDownPortOptsFromContext(context.Hostname, ipAddress))
	Assert(t).That(err, IsNil(), "failed teardown port")

	_, err = client.GetPort(results.Port.ID)
	if err == nil {
		t.Errorf("expected port to be gone %s", ipAddress)
	}
}
