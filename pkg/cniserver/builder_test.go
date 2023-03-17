package cniserver_test

import (
	"testing"

	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	. "github.com/pepinns/go-hamcrest"
)

func Test_builder(t *testing.T) {
	WithStateDir(t, func(dir string) {
		WithTestConfig(t, func(cfg TestingConfig) {
			t.Run("build creates proper dependencies by default", func(t *testing.T) {
				deps, err := cniserver.NewBuilder().Build()
				Assert(t).That(err, IsNil())
				Assert(t).That(deps, Not(IsNil()))
				Assert(t).That(deps.CniHandler(), Not(IsNil()))
				Assert(t).That(deps.OpenstackClient(), Not(IsNil()))
				Assert(t).That(deps.State(), Not(IsNil()))
			})
		})
	})
}
