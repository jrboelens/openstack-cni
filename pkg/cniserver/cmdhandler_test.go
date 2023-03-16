package cniserver

import (
	"testing"

	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"

	. "github.com/pepinns/go-hamcrest"
)

func Test_CreateResultFromPortResult(t *testing.T) {
	t.Run("returns an error with incomplete data", func(t *testing.T) {
		portResult := &openstack.SetupPortResult{}
		cmd := util.CniCommand{}

		result, err := CreateResultFromPortResult(portResult, cmd)
		Assert(t).That(result, IsNil())
		Assert(t).That(err, Equals(ErrIncompletePortResult))
	})
}
