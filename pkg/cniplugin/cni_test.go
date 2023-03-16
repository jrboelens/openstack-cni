package cniplugin_test

import (
	"net"
	"testing"

	currentcni "github.com/containernetworking/cni/pkg/types/040"
	"github.com/jboelensns/openstack-cni/pkg/cniplugin"
	"github.com/jboelensns/openstack-cni/pkg/cnistate"
	"github.com/jboelensns/openstack-cni/pkg/fixtures"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	"github.com/jboelensns/openstack-cni/pkg/fixtures/mocks"
	"github.com/jboelensns/openstack-cni/pkg/util"
	. "github.com/pepinns/go-hamcrest"
)

func Test_Cni(t *testing.T) {
	WithCniState(t, func(state cnistate.State) {
		testData := fixtures.NewTestData()

		t.Run("can execute an add", func(t *testing.T) {
			cniHandler := &mocks.CommandHandlerMock{}
			networking := &mocks.NetworkingMock{}
			sopts := &ServerOpts{CniHandler: cniHandler, Networking: networking, State: state}

			WithServerOpts(t, sopts, func(fix *ServerFixture) {
				cniclient := fix.CniClient()
				// provide a meaningful result back from the http server
				cniHandler.AddFunc = func(cmd util.CniCommand) (*currentcni.Result, error) {
					result := testData.CniResult()

					// setup the mac as the mac of an interface on our machine so the lookup doesn't fail
					ifaces, err := net.Interfaces()
					Assert(t).That(err, IsNil())
					result.Interfaces[0].Mac = ifaces[0].HardwareAddr.String()

					return result, nil
				}

				// skip doing any of the netlink configuration
				networking.ConfigureFunc = func(namespace string, iface *cniplugin.NetworkInterface) error {
					return nil
				}

				args := testData.SkelArgs()

				cni := cniplugin.NewCni(cniclient, networking)
				err := cni.Add(args)
				Assert(t).That(err, IsNil())

				Assert(t).That(cniHandler.AddCalls(), HasLen(1))
				Assert(t).That(networking.ConfigureCalls(), HasLen(1))
			})
		})

		t.Run("can execute a delete", func(t *testing.T) {
			cniHandler := &mocks.CommandHandlerMock{}
			networking := &mocks.NetworkingMock{}
			sopts := &ServerOpts{CniHandler: cniHandler, Networking: networking, State: state}

			WithServerOpts(t, sopts, func(fix *ServerFixture) {
				cniclient := fix.CniClient()
				cniHandler.DelFunc = func(cmd util.CniCommand) error {
					return nil
				}
				args := testData.SkelArgs()

				// setup the state so we can delete it
				Assert(t).That(state.Set(testData.IfaceInfo()), IsNil())

				cni := cniplugin.NewCni(cniclient, networking)
				err := cni.Del(args)
				Assert(t).That(err, IsNil())

				Assert(t).That(cniHandler.DelCalls(), HasLen(1))
			})
		})
	})
}
