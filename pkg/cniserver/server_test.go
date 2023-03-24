package cniserver_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/containernetworking/cni/pkg/types"
	currentcni "github.com/containernetworking/cni/pkg/types/040"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/jboelensns/openstack-cni/pkg/cniclient"
	"github.com/jboelensns/openstack-cni/pkg/cnistate"
	"github.com/jboelensns/openstack-cni/pkg/fixtures"
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

func Test_State(t *testing.T) {
	t.Run("GET /state: 200 when state is found", func(t *testing.T) {
		ifaceInfo := fixtures.NewTestData().IfaceInfo()
		state := &mocks.StateMock{}
		state.GetFunc = func(containerId, ifname string) (*cnistate.IfaceInfo, error) {
			return ifaceInfo, nil
		}

		opts := &ServerOpts{State: state}
		WithServerOpts(t, opts, func(fix *ServerFixture) {
			result, err := fix.CniClient().GetState(ifaceInfo.ContainerId, ifaceInfo.Ifname)

			Assert(t).That(err, IsNil())
			Assert(t).That(result, Equals(ifaceInfo))
		})
	})

	t.Run("GET /state: 404 when state doesn't exist", func(t *testing.T) {
		ifaceInfo := fixtures.NewTestData().IfaceInfo()
		state := &mocks.StateMock{}
		state.GetFunc = func(containerId, ifname string) (*cnistate.IfaceInfo, error) {
			return nil, nil
		}

		opts := &ServerOpts{State: state}
		WithServerOpts(t, opts, func(fix *ServerFixture) {
			url := fmt.Sprintf("%s/%s/%s", fix.Url("/state"), ifaceInfo.ContainerId, ifaceInfo.Ifname)
			resp, err := fix.Client().Get(url, nil)

			Assert(t).That(err, IsNil())
			Assert(t).That(resp.StatusCode, Equals(http.StatusNotFound))
		})
	})

	t.Run("DELETE /state: 204 when state does exist", func(t *testing.T) {
		ifaceInfo := fixtures.NewTestData().IfaceInfo()
		state := &mocks.StateMock{}
		state.GetFunc = func(containerId, ifname string) (*cnistate.IfaceInfo, error) {
			return ifaceInfo, nil
		}
		state.DeleteFunc = func(containerId, ifname string) error { return nil }

		opts := &ServerOpts{State: state}
		WithServerOpts(t, opts, func(fix *ServerFixture) {
			err := fix.CniClient().DeleteState(ifaceInfo.ContainerId, ifaceInfo.Ifname)
			Assert(t).That(err, IsNil())
		})
	})

	t.Run("DELETE /state: 404 when state doesn't exist", func(t *testing.T) {
		state := &mocks.StateMock{}
		state.GetFunc = func(containerId, ifname string) (*cnistate.IfaceInfo, error) {
			return nil, nil
		}

		opts := &ServerOpts{State: state}
		WithServerOpts(t, opts, func(fix *ServerFixture) {
			err := fix.CniClient().DeleteState("foo", "bar")
			Assert(t).That(err, Equals(cniclient.ErrStateNotFound))
		})
	})

	t.Run("POST /state: 204 on success", func(t *testing.T) {
		ifaceInfo := fixtures.NewTestData().IfaceInfo()
		state := &mocks.StateMock{}
		state.SetFunc = func(ifaceInfo *cnistate.IfaceInfo) error { return nil }

		opts := &ServerOpts{State: state}
		WithServerOpts(t, opts, func(fix *ServerFixture) {
			err := fix.CniClient().SetState(ifaceInfo)
			Assert(t).That(err, IsNil())
		})
	})

	t.Run("POST /state: 500 on failure", func(t *testing.T) {
		ifaceInfo := fixtures.NewTestData().IfaceInfo()
		state := &mocks.StateMock{}
		state.SetFunc = func(ifaceInfo *cnistate.IfaceInfo) error { return errors.New("BOOM") }

		opts := &ServerOpts{State: state}
		WithServerOpts(t, opts, func(fix *ServerFixture) {
			err := fix.CniClient().SetState(ifaceInfo)
			Assert(t).That(err, Contains("500"))
		})
	})
}
