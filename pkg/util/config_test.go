package util_test

import (
	"testing"

	"github.com/jboelensns/openstack-cni/pkg/util"

	. "github.com/pepinns/go-hamcrest"
)

func Test_ParsingArgs(t *testing.T) {
	t.Run("returns an empty map with bad data", func(t *testing.T) {
		badArgs := []string{
			"",
			"ONE=2;TWO",
			"NOT;GOOD!AT@ALL=",
		}
		for _, args := range badArgs {
			mappy := util.ParseCniArgs(args)
			Assert(t).That(mappy, HasLen(0))
		}
	})
	t.Run("succeeds with a simply key value", func(t *testing.T) {
		args := "FOO=BAR"
		mappy := util.ParseCniArgs(args)
		Assert(t).That(mappy, HasItem("FOO", Equals("BAR")))
	})
	t.Run("succeeds with many values", func(t *testing.T) {
		args := "ONE=1;TWO=2;THREE=3"
		mappy := util.ParseCniArgs(args)
		Assert(t).That(mappy, AllOf(
			HasItem("ONE", Equals("1")),
			HasItem("TWO", Equals("2")),
			HasItem("THREE", Equals("3")),
		))
	})
}
