package util_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jboelensns/openstack-cni/pkg/util"

	. "github.com/pepinns/go-hamcrest"
)

func Test_Backoff(t *testing.T) {
	getOpts := func() util.RetryOpts {
		return util.RetryOpts{
			BackoffMax:      time.Millisecond * 50,
			BackoffInterval: time.Millisecond * 5,
			MaxWaitTime:     time.Millisecond * 200,
			RetryErrors:     false,
		}
	}

	boomErr := fmt.Errorf("boom")

	t.Run("Backoff does not retry if RetryErrors = false", func(t *testing.T) {
		count := 0
		opts := getOpts()
		opts.RetryErrors = false
		r, err := util.Backoff(opts, func() (*int, error) {
			count += 1
			return nil, boomErr
		})
		Assert(t).That(count, Equals(1))
		Assert(t).That(err, Equals(boomErr))
		Assert(t).That(r, IsNil())
	})

	t.Run("Backoff retries errors when RetryErrors = true", func(t *testing.T) {
		count := 0
		opts := getOpts()
		opts.RetryErrors = true
		r, err := util.Backoff(opts, func() (*int, error) {
			count += 1
			return nil, boomErr
		})
		Assert(t).That(count, GreaterThan(1))
		Assert(t).That(err, Equals(boomErr))
		Assert(t).That(r, IsNil())
	})

	t.Run("Backoff retries and eventually gets a result", func(t *testing.T) {
		count := 0
		exCount := 3
		opts := getOpts()
		opts.RetryErrors = true
		r, err := util.Backoff(opts, func() (*int, error) {
			fmt.Printf("COUNT %d\n", count)
			count += 1
			if count == exCount {
				return &count, nil
			}
			return nil, boomErr
		})
		Assert(t).That(count, Equals(3))
		Assert(t).That(err, IsNil())
		Assert(t).That(r, Equals(&exCount))
	})
}
