// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fixtures

import (
	"github.com/jboelensns/openstack-cni/pkg/util"
	"github.com/vishvananda/netlink"
	"sync"
)

// Ensure, that NetlinkWrapperMock does implement util.NetlinkWrapper.
// If this is not the case, regenerate this file with moq.
var _ util.NetlinkWrapper = &NetlinkWrapperMock{}

// NetlinkWrapperMock is a mock implementation of util.NetlinkWrapper.
//
//	func TestSomethingThatUsesNetlinkWrapper(t *testing.T) {
//
//		// make and configure a mocked util.NetlinkWrapper
//		mockedNetlinkWrapper := &NetlinkWrapperMock{
//			AddrAddFunc: func(link netlink.Link, addr *netlink.Addr) error {
//				panic("mock out the AddrAdd method")
//			},
//			GetNetNsIdByPathFunc: func(namespace string) (int, error) {
//				panic("mock out the GetNetNsIdByPath method")
//			},
//			GetNetNsIdByPidFunc: func(pid int) (int, error) {
//				panic("mock out the GetNetNsIdByPid method")
//			},
//			LinkByNameFunc: func(ifname string) (netlink.Link, error) {
//				panic("mock out the LinkByName method")
//			},
//			LinkSetDownFunc: func(link netlink.Link) error {
//				panic("mock out the LinkSetDown method")
//			},
//			LinkSetNameFunc: func(link netlink.Link, name string) error {
//				panic("mock out the LinkSetName method")
//			},
//			LinkSetNsFdFunc: func(link netlink.Link, fd int) error {
//				panic("mock out the LinkSetNsFd method")
//			},
//			LinkSetUpFunc: func(link netlink.Link) error {
//				panic("mock out the LinkSetUp method")
//			},
//		}
//
//		// use mockedNetlinkWrapper in code that requires util.NetlinkWrapper
//		// and then make assertions.
//
//	}
type NetlinkWrapperMock struct {
	// AddrAddFunc mocks the AddrAdd method.
	AddrAddFunc func(link netlink.Link, addr *netlink.Addr) error

	// GetNetNsIdByPathFunc mocks the GetNetNsIdByPath method.
	GetNetNsIdByPathFunc func(namespace string) (int, error)

	// GetNetNsIdByPidFunc mocks the GetNetNsIdByPid method.
	GetNetNsIdByPidFunc func(pid int) (int, error)

	// LinkByNameFunc mocks the LinkByName method.
	LinkByNameFunc func(ifname string) (netlink.Link, error)

	// LinkSetDownFunc mocks the LinkSetDown method.
	LinkSetDownFunc func(link netlink.Link) error

	// LinkSetNameFunc mocks the LinkSetName method.
	LinkSetNameFunc func(link netlink.Link, name string) error

	// LinkSetNsFdFunc mocks the LinkSetNsFd method.
	LinkSetNsFdFunc func(link netlink.Link, fd int) error

	// LinkSetUpFunc mocks the LinkSetUp method.
	LinkSetUpFunc func(link netlink.Link) error

	// calls tracks calls to the methods.
	calls struct {
		// AddrAdd holds details about calls to the AddrAdd method.
		AddrAdd []struct {
			// Link is the link argument value.
			Link netlink.Link
			// Addr is the addr argument value.
			Addr *netlink.Addr
		}
		// GetNetNsIdByPath holds details about calls to the GetNetNsIdByPath method.
		GetNetNsIdByPath []struct {
			// Namespace is the namespace argument value.
			Namespace string
		}
		// GetNetNsIdByPid holds details about calls to the GetNetNsIdByPid method.
		GetNetNsIdByPid []struct {
			// Pid is the pid argument value.
			Pid int
		}
		// LinkByName holds details about calls to the LinkByName method.
		LinkByName []struct {
			// Ifname is the ifname argument value.
			Ifname string
		}
		// LinkSetDown holds details about calls to the LinkSetDown method.
		LinkSetDown []struct {
			// Link is the link argument value.
			Link netlink.Link
		}
		// LinkSetName holds details about calls to the LinkSetName method.
		LinkSetName []struct {
			// Link is the link argument value.
			Link netlink.Link
			// Name is the name argument value.
			Name string
		}
		// LinkSetNsFd holds details about calls to the LinkSetNsFd method.
		LinkSetNsFd []struct {
			// Link is the link argument value.
			Link netlink.Link
			// Fd is the fd argument value.
			Fd int
		}
		// LinkSetUp holds details about calls to the LinkSetUp method.
		LinkSetUp []struct {
			// Link is the link argument value.
			Link netlink.Link
		}
	}
	lockAddrAdd          sync.RWMutex
	lockGetNetNsIdByPath sync.RWMutex
	lockGetNetNsIdByPid  sync.RWMutex
	lockLinkByName       sync.RWMutex
	lockLinkSetDown      sync.RWMutex
	lockLinkSetName      sync.RWMutex
	lockLinkSetNsFd      sync.RWMutex
	lockLinkSetUp        sync.RWMutex
}

// AddrAdd calls AddrAddFunc.
func (mock *NetlinkWrapperMock) AddrAdd(link netlink.Link, addr *netlink.Addr) error {
	if mock.AddrAddFunc == nil {
		panic("NetlinkWrapperMock.AddrAddFunc: method is nil but NetlinkWrapper.AddrAdd was just called")
	}
	callInfo := struct {
		Link netlink.Link
		Addr *netlink.Addr
	}{
		Link: link,
		Addr: addr,
	}
	mock.lockAddrAdd.Lock()
	mock.calls.AddrAdd = append(mock.calls.AddrAdd, callInfo)
	mock.lockAddrAdd.Unlock()
	return mock.AddrAddFunc(link, addr)
}

// AddrAddCalls gets all the calls that were made to AddrAdd.
// Check the length with:
//
//	len(mockedNetlinkWrapper.AddrAddCalls())
func (mock *NetlinkWrapperMock) AddrAddCalls() []struct {
	Link netlink.Link
	Addr *netlink.Addr
} {
	var calls []struct {
		Link netlink.Link
		Addr *netlink.Addr
	}
	mock.lockAddrAdd.RLock()
	calls = mock.calls.AddrAdd
	mock.lockAddrAdd.RUnlock()
	return calls
}

// GetNetNsIdByPath calls GetNetNsIdByPathFunc.
func (mock *NetlinkWrapperMock) GetNetNsIdByPath(namespace string) (int, error) {
	if mock.GetNetNsIdByPathFunc == nil {
		panic("NetlinkWrapperMock.GetNetNsIdByPathFunc: method is nil but NetlinkWrapper.GetNetNsIdByPath was just called")
	}
	callInfo := struct {
		Namespace string
	}{
		Namespace: namespace,
	}
	mock.lockGetNetNsIdByPath.Lock()
	mock.calls.GetNetNsIdByPath = append(mock.calls.GetNetNsIdByPath, callInfo)
	mock.lockGetNetNsIdByPath.Unlock()
	return mock.GetNetNsIdByPathFunc(namespace)
}

// GetNetNsIdByPathCalls gets all the calls that were made to GetNetNsIdByPath.
// Check the length with:
//
//	len(mockedNetlinkWrapper.GetNetNsIdByPathCalls())
func (mock *NetlinkWrapperMock) GetNetNsIdByPathCalls() []struct {
	Namespace string
} {
	var calls []struct {
		Namespace string
	}
	mock.lockGetNetNsIdByPath.RLock()
	calls = mock.calls.GetNetNsIdByPath
	mock.lockGetNetNsIdByPath.RUnlock()
	return calls
}

// GetNetNsIdByPid calls GetNetNsIdByPidFunc.
func (mock *NetlinkWrapperMock) GetNetNsIdByPid(pid int) (int, error) {
	if mock.GetNetNsIdByPidFunc == nil {
		panic("NetlinkWrapperMock.GetNetNsIdByPidFunc: method is nil but NetlinkWrapper.GetNetNsIdByPid was just called")
	}
	callInfo := struct {
		Pid int
	}{
		Pid: pid,
	}
	mock.lockGetNetNsIdByPid.Lock()
	mock.calls.GetNetNsIdByPid = append(mock.calls.GetNetNsIdByPid, callInfo)
	mock.lockGetNetNsIdByPid.Unlock()
	return mock.GetNetNsIdByPidFunc(pid)
}

// GetNetNsIdByPidCalls gets all the calls that were made to GetNetNsIdByPid.
// Check the length with:
//
//	len(mockedNetlinkWrapper.GetNetNsIdByPidCalls())
func (mock *NetlinkWrapperMock) GetNetNsIdByPidCalls() []struct {
	Pid int
} {
	var calls []struct {
		Pid int
	}
	mock.lockGetNetNsIdByPid.RLock()
	calls = mock.calls.GetNetNsIdByPid
	mock.lockGetNetNsIdByPid.RUnlock()
	return calls
}

// LinkByName calls LinkByNameFunc.
func (mock *NetlinkWrapperMock) LinkByName(ifname string) (netlink.Link, error) {
	if mock.LinkByNameFunc == nil {
		panic("NetlinkWrapperMock.LinkByNameFunc: method is nil but NetlinkWrapper.LinkByName was just called")
	}
	callInfo := struct {
		Ifname string
	}{
		Ifname: ifname,
	}
	mock.lockLinkByName.Lock()
	mock.calls.LinkByName = append(mock.calls.LinkByName, callInfo)
	mock.lockLinkByName.Unlock()
	return mock.LinkByNameFunc(ifname)
}

// LinkByNameCalls gets all the calls that were made to LinkByName.
// Check the length with:
//
//	len(mockedNetlinkWrapper.LinkByNameCalls())
func (mock *NetlinkWrapperMock) LinkByNameCalls() []struct {
	Ifname string
} {
	var calls []struct {
		Ifname string
	}
	mock.lockLinkByName.RLock()
	calls = mock.calls.LinkByName
	mock.lockLinkByName.RUnlock()
	return calls
}

// LinkSetDown calls LinkSetDownFunc.
func (mock *NetlinkWrapperMock) LinkSetDown(link netlink.Link) error {
	if mock.LinkSetDownFunc == nil {
		panic("NetlinkWrapperMock.LinkSetDownFunc: method is nil but NetlinkWrapper.LinkSetDown was just called")
	}
	callInfo := struct {
		Link netlink.Link
	}{
		Link: link,
	}
	mock.lockLinkSetDown.Lock()
	mock.calls.LinkSetDown = append(mock.calls.LinkSetDown, callInfo)
	mock.lockLinkSetDown.Unlock()
	return mock.LinkSetDownFunc(link)
}

// LinkSetDownCalls gets all the calls that were made to LinkSetDown.
// Check the length with:
//
//	len(mockedNetlinkWrapper.LinkSetDownCalls())
func (mock *NetlinkWrapperMock) LinkSetDownCalls() []struct {
	Link netlink.Link
} {
	var calls []struct {
		Link netlink.Link
	}
	mock.lockLinkSetDown.RLock()
	calls = mock.calls.LinkSetDown
	mock.lockLinkSetDown.RUnlock()
	return calls
}

// LinkSetName calls LinkSetNameFunc.
func (mock *NetlinkWrapperMock) LinkSetName(link netlink.Link, name string) error {
	if mock.LinkSetNameFunc == nil {
		panic("NetlinkWrapperMock.LinkSetNameFunc: method is nil but NetlinkWrapper.LinkSetName was just called")
	}
	callInfo := struct {
		Link netlink.Link
		Name string
	}{
		Link: link,
		Name: name,
	}
	mock.lockLinkSetName.Lock()
	mock.calls.LinkSetName = append(mock.calls.LinkSetName, callInfo)
	mock.lockLinkSetName.Unlock()
	return mock.LinkSetNameFunc(link, name)
}

// LinkSetNameCalls gets all the calls that were made to LinkSetName.
// Check the length with:
//
//	len(mockedNetlinkWrapper.LinkSetNameCalls())
func (mock *NetlinkWrapperMock) LinkSetNameCalls() []struct {
	Link netlink.Link
	Name string
} {
	var calls []struct {
		Link netlink.Link
		Name string
	}
	mock.lockLinkSetName.RLock()
	calls = mock.calls.LinkSetName
	mock.lockLinkSetName.RUnlock()
	return calls
}

// LinkSetNsFd calls LinkSetNsFdFunc.
func (mock *NetlinkWrapperMock) LinkSetNsFd(link netlink.Link, fd int) error {
	if mock.LinkSetNsFdFunc == nil {
		panic("NetlinkWrapperMock.LinkSetNsFdFunc: method is nil but NetlinkWrapper.LinkSetNsFd was just called")
	}
	callInfo := struct {
		Link netlink.Link
		Fd   int
	}{
		Link: link,
		Fd:   fd,
	}
	mock.lockLinkSetNsFd.Lock()
	mock.calls.LinkSetNsFd = append(mock.calls.LinkSetNsFd, callInfo)
	mock.lockLinkSetNsFd.Unlock()
	return mock.LinkSetNsFdFunc(link, fd)
}

// LinkSetNsFdCalls gets all the calls that were made to LinkSetNsFd.
// Check the length with:
//
//	len(mockedNetlinkWrapper.LinkSetNsFdCalls())
func (mock *NetlinkWrapperMock) LinkSetNsFdCalls() []struct {
	Link netlink.Link
	Fd   int
} {
	var calls []struct {
		Link netlink.Link
		Fd   int
	}
	mock.lockLinkSetNsFd.RLock()
	calls = mock.calls.LinkSetNsFd
	mock.lockLinkSetNsFd.RUnlock()
	return calls
}

// LinkSetUp calls LinkSetUpFunc.
func (mock *NetlinkWrapperMock) LinkSetUp(link netlink.Link) error {
	if mock.LinkSetUpFunc == nil {
		panic("NetlinkWrapperMock.LinkSetUpFunc: method is nil but NetlinkWrapper.LinkSetUp was just called")
	}
	callInfo := struct {
		Link netlink.Link
	}{
		Link: link,
	}
	mock.lockLinkSetUp.Lock()
	mock.calls.LinkSetUp = append(mock.calls.LinkSetUp, callInfo)
	mock.lockLinkSetUp.Unlock()
	return mock.LinkSetUpFunc(link)
}

// LinkSetUpCalls gets all the calls that were made to LinkSetUp.
// Check the length with:
//
//	len(mockedNetlinkWrapper.LinkSetUpCalls())
func (mock *NetlinkWrapperMock) LinkSetUpCalls() []struct {
	Link netlink.Link
} {
	var calls []struct {
		Link netlink.Link
	}
	mock.lockLinkSetUp.RLock()
	calls = mock.calls.LinkSetUp
	mock.lockLinkSetUp.RUnlock()
	return calls
}
