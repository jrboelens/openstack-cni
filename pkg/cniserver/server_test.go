package cniserver_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/containernetworking/cni/pkg/types"
	currentcni "github.com/containernetworking/cni/pkg/types/040"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	"github.com/jboelensns/openstack-cni/pkg/fixtures/mocks"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"
	. "github.com/pepinns/go-hamcrest"
)

func Test_Ping(t *testing.T) {
	WithServer(t, func(fix *ServerFixture) {
		fix.Ping(t)
	})
}

func Test_Cni_Errors(t *testing.T) {
	WithServer(t, func(fix *ServerFixture) {
		t.Run("/cni returns 405 for GET", func(t *testing.T) {
			resp, err := fix.Client().Get(fix.Url("/cni"), nil)
			Assert(t).That(err, IsNil())
			Assert(t).That(resp.StatusCode, Equals(405))
		})
		t.Run("/cni returns 400 when called with bad data", func(t *testing.T) {
			resp, err := fix.Client().Post(fix.Url("/cni"), []byte("NOT JSON DATA"), DefaultDoOpts())

			Assert(t).That(err, IsNil())
			Assert(t).That(resp.StatusCode, Equals(400))
		})
	})
}

func Test_Cni_Add(t *testing.T) {
	t.Run("/cni returns 500 with an error json when add fails", func(t *testing.T) {
		cniHandler := &mocks.CommandHandlerMock{}
		cniHandler.AddFunc = func(cmd util.CniCommand) (*currentcni.Result, error) {
			return nil, fmt.Errorf("BOOM")
		}

		opts := &ServerOpts{CniHandler: cniHandler}
		WithServerOpts(t, opts, func(fix *ServerFixture) {
			cmd := fix.TestData().CniCommand()
			resp, err := fix.CniClient().CniCommand(cmd)
			fix.Assert(t).CniErrorHasCode(resp, err, types.ErrUnknown)
		})
	})

	t.Run("/cni returns 400 when invalid data is posted", func(t *testing.T) {
		WithServer(t, func(fix *ServerFixture) {
			cmd := fix.TestData().CniCommand()
			cmd.ContainerID = ""
			resp, err := fix.CniClient().CniCommand(cmd)

			Assert(t).That(err, IsNil())
			Assert(t).That(resp.StatusCode, Equals(http.StatusBadRequest))
		})
	})

	t.Run("/cni returns 200 with result json when add succeeds", func(t *testing.T) {
		inResult := NewTestData().CniResult()

		cniHandler := &mocks.CommandHandlerMock{}
		cniHandler.AddFunc = func(cmd util.CniCommand) (*currentcni.Result, error) {
			return inResult, nil
		}

		opts := &ServerOpts{CniHandler: cniHandler}
		WithServerOpts(t, opts, func(fix *ServerFixture) {
			cmd := fix.TestData().CniCommand()
			resp, err := fix.CniClient().CniCommand(cmd)

			result := fix.Assert(t).IsCniResult(resp, err)
			Assert(t).That(result.CNIVersion, Equals(inResult.CNIVersion))
			Assert(t).That(result, Equals(inResult))
		})
	})
}

func Test_Health(t *testing.T) {
	t.Run("/health returns 200 when healthy", func(t *testing.T) {
		osClient := &mocks.OpenstackClientMock{}
		osClient.GetServerByNameFunc = func(name string) (*servers.Server, error) {
			return nil, openstack.ErrServerNotFound
		}

		opts := &ServerOpts{OpenstackClient: osClient}
		WithServerOpts(t, opts, func(fix *ServerFixture) {
			resp, err := fix.Client().Get(fix.Url("/health"), nil)
			result := fix.Assert(t).IsGoodHealthResult(resp, err)

			Assert(t).That(result.Checks, HasLen(1))
			for _, check := range result.Checks {
				if check.Name == "openstack" {
					Assert(t).That(check.IsHealthy, IsTrue())
					Assert(t).That(check.Error, Equals(""))
				}
			}
		})
	})

	t.Run("/health returns 500 when failing", func(t *testing.T) {
		osClient := &mocks.OpenstackClientMock{}
		osClient.GetServerByNameFunc = func(name string) (*servers.Server, error) {
			return nil, errors.New("BOOM")
		}

		opts := &ServerOpts{OpenstackClient: osClient}
		WithServerOpts(t, opts, func(fix *ServerFixture) {
			resp, err := fix.Client().Get(fix.Url("/health"), nil)
			result := fix.Assert(t).IsSickHealthResult(resp, err)

			Assert(t).That(result.Checks, HasLen(1))
			for _, check := range result.Checks {
				if check.Name == "openstack" {
					Assert(t).That(check.IsHealthy, IsFalse())
					Assert(t).That(check.Error, Not(Equals("")))
				}
			}
		})
	})
}
