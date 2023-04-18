package cniserver_test

import (
	"testing"

	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	. "github.com/pepinns/go-hamcrest"
)

func Test_builder(t *testing.T) {
	WithTestConfig(t, func(cfg TestingConfig) {
		t.Run("build creates proper dependencies by default", func(t *testing.T) {
			config := cniserver.NewConfig()
			deps, err := cniserver.NewBuilder(config).Build()
			Assert(t).That(err, IsNil())
			Assert(t).That(deps, Not(IsNil()))
			Assert(t).That(deps.CniHandler(), Not(IsNil()))
			Assert(t).That(deps.OpenstackClient(), Not(IsNil()))
		})
	})
}
