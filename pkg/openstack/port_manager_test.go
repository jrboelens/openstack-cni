package openstack_test

import (
	"testing"
	"time"

	"github.com/jboelensns/openstack-cni/pkg/cniserver"
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

		t.Run("can setup a port with all options and tear it down", func(t *testing.T) {
			context.CniConfig.SubnetName = cfg.SubnetName
			context.CniConfig.PortDescription = "description"
			f := false
			context.CniConfig.AdminStateUp = &f
			context.CniConfig.MacAddress = "52:54:00:28:ea:16"
			// This cannot be tested without a sepcific device id
			// context.CniConfig.DeviceId = "4be2ed0a-23c4-4c5b-91b3-eedce17b3de2"
			context.CniConfig.DeviceOwner = "compute:nova"
			context.CniConfig.TenantId = "67f06cc9d851455f94fc0380233ab86c"
			context.CniConfig.AllowedAddressPairs = []util.AddressPair{{IpAddress: "1.1.1.1", MacAddress: "52:54:00:28:ea:16"}}
			// This cannot be tested unless openstack is setup to accept specific value spec pairs
			// context.CniConfig.ValueSpecs = &map[string]string{
			// 	"foo": "bar",
			// }

			SetupAndTeardownPort(t, context, client)
		})
	})
}

func SetupAndTeardownPort(t *testing.T, context util.CniContext, client openstack.OpenstackClient) {
	pm := openstack.NewPortManager(client)
	opts := openstack.SetupPortOptsFromContext(context)
	opts.Tags = cniserver.NewPortTags(context.Command)

	results, err := pm.SetupPort(opts)
	Assert(t).That(err, IsNil(), "failed to setup port")

	_, err = client.GetPortByTags(opts.Tags.AsStringSlice())
	Assert(t).That(err, IsNil(), "failed get port by tags %s", opts.Tags.String())

	tdOpts := openstack.TearDownPortOpts{Hostname: context.Hostname, Tags: cniserver.NewPortTags(context.Command)}
	err = pm.TeardownPort(tdOpts)
	Assert(t).That(err, IsNil(), "failed teardown port")

	_, err = client.GetPort(results.Port.ID)
	if err == nil {
		t.Errorf("expected port to be gone with tags %s", tdOpts.Tags.String())
	}

	results, err = pm.SetupPort(opts)
	Assert(t).That(err, IsNil(), "failed to setup port")

	_, err = client.GetPortByTags(opts.Tags.AsStringSlice())
	Assert(t).That(err, IsNil(), "failed get port by tags %s", opts.Tags.String())

	tdOpts = openstack.TearDownPortOpts{Hostname: context.Hostname, Tags: cniserver.NewPortTags(context.Command)}
	err = pm.TeardownPort(tdOpts)
	Assert(t).That(err, IsNil(), "failed teardown port")

	_, err = client.GetPort(results.Port.ID)
	if err == nil {
		t.Errorf("expected port to be gone with tags %s", tdOpts.Tags.String())
	}
}
