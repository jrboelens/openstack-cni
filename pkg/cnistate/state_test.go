package cnistate_test

import (
	"encoding/json"
	"testing"

	"github.com/jboelensns/openstack-cni/pkg/cnistate"
	"github.com/jboelensns/openstack-cni/pkg/fixtures"
	. "github.com/pepinns/go-hamcrest"
)

func TestState(t *testing.T) {
	data := fixtures.NewTestData()

	t.Run("can write, read and delete a state file", func(t *testing.T) {
		fixtures.WithCniState(t, func(state cnistate.State) {
			info := data.IfaceInfo()

			Assert(t).That(state.Set(info), IsNil())

			newInfo, err := state.Get(info.ContainerId, info.Ifname)
			Assert(t).That(err, IsNil())
			Assert(t).That(newInfo.IpAddress, Equals(info.IpAddress))

			Assert(t).That(state.Delete(info.ContainerId, info.Ifname), IsNil())

			newInfo, err = state.Get(info.ContainerId, info.Ifname)
			Assert(t).That(err, IsNil())
			Assert(t).That(newInfo, IsNil())
		})
	})

	t.Run("returns a nil result when the file doesn't exist", func(t *testing.T) {
		fixtures.WithCniState(t, func(state cnistate.State) {
			info := data.IfaceInfo()
			info.ContainerId = "dead1"

			result, err := state.Get(info.ContainerId, info.Ifname)
			Assert(t).That(err, IsNil())
			Assert(t).That(result, IsNil())
		})
	})

	t.Run("fails the file contains bogus data", func(t *testing.T) {
		fixtures.WithTempDir(t, func(dir string) {
			info := data.IfaceInfo()
			info.ContainerId = "f00ba7"

			state := cnistate.NewState(dir)
			err := state.SetRaw(state.Filename(info.ContainerId, info.Ifname), []byte("NOTJSON"))
			Assert(t).That(err, IsNil())

			_, err = state.Get(info.ContainerId, info.Ifname)
			Assert(t).That(IsType[*json.SyntaxError](err), IsTrue())
		})
	})

	t.Run("fails when the state directory doesn't exist", func(t *testing.T) {
		info := data.IfaceInfo()
		state := cnistate.NewState("/tmp/this/directory/does/not/exist")
		Assert(t).That(state.Set(info), Not(IsNil()))
	})
}

func IsType[T any](i any) bool {
	_, ok := i.(T)
	return ok
}
