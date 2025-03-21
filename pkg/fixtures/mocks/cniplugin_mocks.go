// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"github.com/jboelensns/openstack-cni/pkg/cniplugin"
	"net"
	"sync"
)

// Ensure, that NetworkingMock does implement cniplugin.Networking.
// If this is not the case, regenerate this file with moq.
var _ cniplugin.Networking = &NetworkingMock{}

// NetworkingMock is a mock implementation of cniplugin.Networking.
//
//	func TestSomethingThatUsesNetworking(t *testing.T) {
//
//		// make and configure a mocked cniplugin.Networking
//		mockedNetworking := &NetworkingMock{
//			ConfigureFunc: func(namespace string, iface *cniplugin.NetworkInterface) error {
//				panic("mock out the Configure method")
//			},
//			GetIfaceByMacFunc: func(mac string) (*net.Interface, error) {
//				panic("mock out the GetIfaceByMac method")
//			},
//		}
//
//		// use mockedNetworking in code that requires cniplugin.Networking
//		// and then make assertions.
//
//	}
type NetworkingMock struct {
	// ConfigureFunc mocks the Configure method.
	ConfigureFunc func(namespace string, iface *cniplugin.NetworkInterface) error

	// GetIfaceByMacFunc mocks the GetIfaceByMac method.
	GetIfaceByMacFunc func(mac string) (*net.Interface, error)

	// calls tracks calls to the methods.
	calls struct {
		// Configure holds details about calls to the Configure method.
		Configure []struct {
			// Namespace is the namespace argument value.
			Namespace string
			// Iface is the iface argument value.
			Iface *cniplugin.NetworkInterface
		}
		// GetIfaceByMac holds details about calls to the GetIfaceByMac method.
		GetIfaceByMac []struct {
			// Mac is the mac argument value.
			Mac string
		}
	}
	lockConfigure     sync.RWMutex
	lockGetIfaceByMac sync.RWMutex
}

// Configure calls ConfigureFunc.
func (mock *NetworkingMock) Configure(namespace string, iface *cniplugin.NetworkInterface) error {
	if mock.ConfigureFunc == nil {
		panic("NetworkingMock.ConfigureFunc: method is nil but Networking.Configure was just called")
	}
	callInfo := struct {
		Namespace string
		Iface     *cniplugin.NetworkInterface
	}{
		Namespace: namespace,
		Iface:     iface,
	}
	mock.lockConfigure.Lock()
	mock.calls.Configure = append(mock.calls.Configure, callInfo)
	mock.lockConfigure.Unlock()
	return mock.ConfigureFunc(namespace, iface)
}

// ConfigureCalls gets all the calls that were made to Configure.
// Check the length with:
//
//	len(mockedNetworking.ConfigureCalls())
func (mock *NetworkingMock) ConfigureCalls() []struct {
	Namespace string
	Iface     *cniplugin.NetworkInterface
} {
	var calls []struct {
		Namespace string
		Iface     *cniplugin.NetworkInterface
	}
	mock.lockConfigure.RLock()
	calls = mock.calls.Configure
	mock.lockConfigure.RUnlock()
	return calls
}

// GetIfaceByMac calls GetIfaceByMacFunc.
func (mock *NetworkingMock) GetIfaceByMac(mac string) (*net.Interface, error) {
	if mock.GetIfaceByMacFunc == nil {
		panic("NetworkingMock.GetIfaceByMacFunc: method is nil but Networking.GetIfaceByMac was just called")
	}
	callInfo := struct {
		Mac string
	}{
		Mac: mac,
	}
	mock.lockGetIfaceByMac.Lock()
	mock.calls.GetIfaceByMac = append(mock.calls.GetIfaceByMac, callInfo)
	mock.lockGetIfaceByMac.Unlock()
	return mock.GetIfaceByMacFunc(mac)
}

// GetIfaceByMacCalls gets all the calls that were made to GetIfaceByMac.
// Check the length with:
//
//	len(mockedNetworking.GetIfaceByMacCalls())
func (mock *NetworkingMock) GetIfaceByMacCalls() []struct {
	Mac string
} {
	var calls []struct {
		Mac string
	}
	mock.lockGetIfaceByMac.RLock()
	calls = mock.calls.GetIfaceByMac
	mock.lockGetIfaceByMac.RUnlock()
	return calls
}
