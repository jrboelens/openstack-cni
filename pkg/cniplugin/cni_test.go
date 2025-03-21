package cniplugin_test

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	currentcni "github.com/containernetworking/cni/pkg/types/040"
	"github.com/go-chi/httplog"
	"github.com/jboelensns/openstack-cni/pkg/cniplugin"
	"github.com/jboelensns/openstack-cni/pkg/fixtures"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	"github.com/jboelensns/openstack-cni/pkg/fixtures/mocks"
	"github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
	. "github.com/pepinns/go-hamcrest"
)

func Test_Cni(t *testing.T) {
	testData := fixtures.NewTestData()

	getLocalMac := func(t *testing.T) string {
		t.Helper()
		ifaces, err := net.Interfaces()
		Assert(t).That(err, IsNil())
		for _, iface := range ifaces {
			// return the first interface with a MAC
			if iface.HardwareAddr == nil {
				continue
			}
			return iface.HardwareAddr.String()
		}
		// if none of the interfaces had a MAC default to the first interface
		return ifaces[0].HardwareAddr.String()
	}
	getLocalMacAddr := func(t *testing.T) net.HardwareAddr {
		t.Helper()
		mac, err := net.ParseMAC(getLocalMac(t))
		Assert(t).That(err, IsNil())
		return mac
	}

	t.Run("can execute an add", func(t *testing.T) {
		cniHandler := &mocks.CommandHandlerMock{}
		networking := &mocks.NetworkingMock{}
		sopts := &ServerOpts{CniHandler: cniHandler, Networking: networking}

		WithServerOpts(t, sopts, func(fix *ServerFixture) {
			cniclient := fix.CniClient()
			// provide a meaningful result back from the http server
			cniHandler.AddFunc = func(cmd util.CniCommand) (*currentcni.Result, error) {
				result := testData.CniResult()

				// setup the mac as the mac of an interface on our machine so the lookup doesn't fail
				result.Interfaces[0].Mac = getLocalMac(t)

				return result, nil
			}

			networking.GetIfaceByMacFunc = func(mac string) (*net.Interface, error) {
				return &net.Interface{}, nil
			}

			// skip doing any of the netlink configuration
			networking.ConfigureFunc = func(namespace string, iface *cniplugin.NetworkInterface) error {
				return nil
			}

			args := testData.SkelArgs()

			cni := cniplugin.NewCni(cniclient, networking, cniplugin.DefaultCniOpts())
			err := cni.Add(args)
			Assert(t).That(err, IsNil())

			Assert(t).That(cniHandler.AddCalls(), HasLen(1))
			Assert(t).That(networking.ConfigureCalls(), HasLen(1))
		})
	})

	t.Run("retrying GetIfaceByMacFunc works", func(t *testing.T) {
		cniHandler := &mocks.CommandHandlerMock{}
		networking := &mocks.NetworkingMock{}
		sopts := &ServerOpts{CniHandler: cniHandler, Networking: networking}

		WithServerOpts(t, sopts, func(fix *ServerFixture) {
			cniclient := fix.CniClient()
			// provide a meaningful result back from the http server
			result := testData.CniResult()
			// setup the mac as the mac of an interface on our machine so the lookup doesn't fail
			result.Interfaces[0].Mac = getLocalMac(t)
			cniHandler.AddFunc = func(cmd util.CniCommand) (*currentcni.Result, error) {
				return result, nil
			}

			// simiulate failure to look up the mac address for some amount of time
			count := 0
			networking.GetIfaceByMacFunc = func(mac string) (*net.Interface, error) {
				if count < 3 {
					time.Sleep(time.Duration(time.Millisecond * 50))
					count += 1
					return nil, fmt.Errorf("interface not found")
				}
				return &net.Interface{
					HardwareAddr: getLocalMacAddr(t),
				}, nil
			}

			// skip doing any of the netlink configuration
			networking.ConfigureFunc = func(namespace string, iface *cniplugin.NetworkInterface) error {
				return nil
			}

			args := testData.SkelArgs()

			cniOpts := cniplugin.DefaultCniOpts()
			cni := cniplugin.NewCni(cniclient, networking, cniOpts)
			err := cni.Add(args)
			Assert(t).That(err, IsNil())

			Assert(t).That(cniHandler.AddCalls(), HasLen(1))
			Assert(t).That(networking.ConfigureCalls(), HasLen(1))
		})
	})

	t.Run("can execute a delete", func(t *testing.T) {
		logging.SetupLogging("openstack-cni-daemon", httplog.DefaultOptions, os.Stderr)
		cniHandler := &mocks.CommandHandlerMock{}
		networking := &mocks.NetworkingMock{}
		sopts := &ServerOpts{CniHandler: cniHandler, Networking: networking}

		WithServerOpts(t, sopts, func(fix *ServerFixture) {
			cniclient := fix.CniClient()
			cniHandler.DelFunc = func(cmd util.CniCommand) error {
				return nil
			}
			args := testData.SkelArgs()

			cni := cniplugin.NewCni(cniclient, networking, cniplugin.DefaultCniOpts())
			err := cni.Del(args)
			Assert(t).That(err, IsNil())

			Assert(t).That(cniHandler.DelCalls(), HasLen(1))
		})
	})

	t.Run("waitForUdev defaults to true", func(t *testing.T) {
		cfg, err := cniplugin.LoadConfig()
		Assert(t).That(err, IsNil())
		Assert(t).That(cfg.WaitForUdev, IsTrue())
	})
}
