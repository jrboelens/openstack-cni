// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"github.com/jboelensns/openstack-cni/pkg/cnistate"
	"sync"
)

// Ensure, that StateMock does implement cnistate.State.
// If this is not the case, regenerate this file with moq.
var _ cnistate.State = &StateMock{}

// StateMock is a mock implementation of cnistate.State.
//
//	func TestSomethingThatUsesState(t *testing.T) {
//
//		// make and configure a mocked cnistate.State
//		mockedState := &StateMock{
//			DeleteFunc: func(containerId string, ifname string) error {
//				panic("mock out the Delete method")
//			},
//			GetFunc: func(containerId string, ifname string) (*cnistate.IfaceInfo, error) {
//				panic("mock out the Get method")
//			},
//			SetFunc: func(ifaceInfo *cnistate.IfaceInfo) error {
//				panic("mock out the Set method")
//			},
//		}
//
//		// use mockedState in code that requires cnistate.State
//		// and then make assertions.
//
//	}
type StateMock struct {
	// DeleteFunc mocks the Delete method.
	DeleteFunc func(containerId string, ifname string) error

	// GetFunc mocks the Get method.
	GetFunc func(containerId string, ifname string) (*cnistate.IfaceInfo, error)

	// SetFunc mocks the Set method.
	SetFunc func(ifaceInfo *cnistate.IfaceInfo) error

	// calls tracks calls to the methods.
	calls struct {
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// ContainerId is the containerId argument value.
			ContainerId string
			// Ifname is the ifname argument value.
			Ifname string
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// ContainerId is the containerId argument value.
			ContainerId string
			// Ifname is the ifname argument value.
			Ifname string
		}
		// Set holds details about calls to the Set method.
		Set []struct {
			// IfaceInfo is the ifaceInfo argument value.
			IfaceInfo *cnistate.IfaceInfo
		}
	}
	lockDelete sync.RWMutex
	lockGet    sync.RWMutex
	lockSet    sync.RWMutex
}

// Delete calls DeleteFunc.
func (mock *StateMock) Delete(containerId string, ifname string) error {
	if mock.DeleteFunc == nil {
		panic("StateMock.DeleteFunc: method is nil but State.Delete was just called")
	}
	callInfo := struct {
		ContainerId string
		Ifname      string
	}{
		ContainerId: containerId,
		Ifname:      ifname,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(containerId, ifname)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//
//	len(mockedState.DeleteCalls())
func (mock *StateMock) DeleteCalls() []struct {
	ContainerId string
	Ifname      string
} {
	var calls []struct {
		ContainerId string
		Ifname      string
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *StateMock) Get(containerId string, ifname string) (*cnistate.IfaceInfo, error) {
	if mock.GetFunc == nil {
		panic("StateMock.GetFunc: method is nil but State.Get was just called")
	}
	callInfo := struct {
		ContainerId string
		Ifname      string
	}{
		ContainerId: containerId,
		Ifname:      ifname,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(containerId, ifname)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//
//	len(mockedState.GetCalls())
func (mock *StateMock) GetCalls() []struct {
	ContainerId string
	Ifname      string
} {
	var calls []struct {
		ContainerId string
		Ifname      string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// Set calls SetFunc.
func (mock *StateMock) Set(ifaceInfo *cnistate.IfaceInfo) error {
	if mock.SetFunc == nil {
		panic("StateMock.SetFunc: method is nil but State.Set was just called")
	}
	callInfo := struct {
		IfaceInfo *cnistate.IfaceInfo
	}{
		IfaceInfo: ifaceInfo,
	}
	mock.lockSet.Lock()
	mock.calls.Set = append(mock.calls.Set, callInfo)
	mock.lockSet.Unlock()
	return mock.SetFunc(ifaceInfo)
}

// SetCalls gets all the calls that were made to Set.
// Check the length with:
//
//	len(mockedState.SetCalls())
func (mock *StateMock) SetCalls() []struct {
	IfaceInfo *cnistate.IfaceInfo
} {
	var calls []struct {
		IfaceInfo *cnistate.IfaceInfo
	}
	mock.lockSet.RLock()
	calls = mock.calls.Set
	mock.lockSet.RUnlock()
	return calls
}