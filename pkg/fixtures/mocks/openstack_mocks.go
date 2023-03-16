// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"sync"
)

// Ensure, that OpenstackClientMock does implement openstack.OpenstackClient.
// If this is not the case, regenerate this file with moq.
var _ openstack.OpenstackClient = &OpenstackClientMock{}

// OpenstackClientMock is a mock implementation of openstack.OpenstackClient.
//
//	func TestSomethingThatUsesOpenstackClient(t *testing.T) {
//
//		// make and configure a mocked openstack.OpenstackClient
//		mockedOpenstackClient := &OpenstackClientMock{
//			AssignPortFunc: func(portId string, serverId string) (*attachinterfaces.Interface, error) {
//				panic("mock out the AssignPort method")
//			},
//			ClientsFunc: func() *openstack.ApiClients {
//				panic("mock out the Clients method")
//			},
//			CreatePortFunc: func(opts ports.CreateOpts) (*ports.Port, error) {
//				panic("mock out the CreatePort method")
//			},
//			DeletePortFunc: func(portId string) error {
//				panic("mock out the DeletePort method")
//			},
//			DetachPortFunc: func(portId string, serverId string) error {
//				panic("mock out the DetachPort method")
//			},
//			GetNetworkByNameFunc: func(name string) (*networks.Network, error) {
//				panic("mock out the GetNetworkByName method")
//			},
//			GetPortFunc: func(portId string) (*ports.Port, error) {
//				panic("mock out the GetPort method")
//			},
//			GetPortByIpFunc: func(ip string) (*ports.Port, error) {
//				panic("mock out the GetPortByIp method")
//			},
//			GetProjectByNameFunc: func(name string) (*projects.Project, error) {
//				panic("mock out the GetProjectByName method")
//			},
//			GetSecurityGroupByNameFunc: func(name string, projectId string) (*groups.SecGroup, error) {
//				panic("mock out the GetSecurityGroupByName method")
//			},
//			GetServerByNameFunc: func(name string) (*servers.Server, error) {
//				panic("mock out the GetServerByName method")
//			},
//			GetSubnetFunc: func(id string) (*subnets.Subnet, error) {
//				panic("mock out the GetSubnet method")
//			},
//			GetSubnetByNameFunc: func(name string, networkId string) (*subnets.Subnet, error) {
//				panic("mock out the GetSubnetByName method")
//			},
//		}
//
//		// use mockedOpenstackClient in code that requires openstack.OpenstackClient
//		// and then make assertions.
//
//	}
type OpenstackClientMock struct {
	// AssignPortFunc mocks the AssignPort method.
	AssignPortFunc func(portId string, serverId string) (*attachinterfaces.Interface, error)

	// ClientsFunc mocks the Clients method.
	ClientsFunc func() *openstack.ApiClients

	// CreatePortFunc mocks the CreatePort method.
	CreatePortFunc func(opts ports.CreateOpts) (*ports.Port, error)

	// DeletePortFunc mocks the DeletePort method.
	DeletePortFunc func(portId string) error

	// DetachPortFunc mocks the DetachPort method.
	DetachPortFunc func(portId string, serverId string) error

	// GetNetworkByNameFunc mocks the GetNetworkByName method.
	GetNetworkByNameFunc func(name string) (*networks.Network, error)

	// GetPortFunc mocks the GetPort method.
	GetPortFunc func(portId string) (*ports.Port, error)

	// GetPortByIpFunc mocks the GetPortByIp method.
	GetPortByIpFunc func(ip string) (*ports.Port, error)

	// GetProjectByNameFunc mocks the GetProjectByName method.
	GetProjectByNameFunc func(name string) (*projects.Project, error)

	// GetSecurityGroupByNameFunc mocks the GetSecurityGroupByName method.
	GetSecurityGroupByNameFunc func(name string, projectId string) (*groups.SecGroup, error)

	// GetServerByNameFunc mocks the GetServerByName method.
	GetServerByNameFunc func(name string) (*servers.Server, error)

	// GetSubnetFunc mocks the GetSubnet method.
	GetSubnetFunc func(id string) (*subnets.Subnet, error)

	// GetSubnetByNameFunc mocks the GetSubnetByName method.
	GetSubnetByNameFunc func(name string, networkId string) (*subnets.Subnet, error)

	// calls tracks calls to the methods.
	calls struct {
		// AssignPort holds details about calls to the AssignPort method.
		AssignPort []struct {
			// PortId is the portId argument value.
			PortId string
			// ServerId is the serverId argument value.
			ServerId string
		}
		// Clients holds details about calls to the Clients method.
		Clients []struct {
		}
		// CreatePort holds details about calls to the CreatePort method.
		CreatePort []struct {
			// Opts is the opts argument value.
			Opts ports.CreateOpts
		}
		// DeletePort holds details about calls to the DeletePort method.
		DeletePort []struct {
			// PortId is the portId argument value.
			PortId string
		}
		// DetachPort holds details about calls to the DetachPort method.
		DetachPort []struct {
			// PortId is the portId argument value.
			PortId string
			// ServerId is the serverId argument value.
			ServerId string
		}
		// GetNetworkByName holds details about calls to the GetNetworkByName method.
		GetNetworkByName []struct {
			// Name is the name argument value.
			Name string
		}
		// GetPort holds details about calls to the GetPort method.
		GetPort []struct {
			// PortId is the portId argument value.
			PortId string
		}
		// GetPortByIp holds details about calls to the GetPortByIp method.
		GetPortByIp []struct {
			// IP is the ip argument value.
			IP string
		}
		// GetProjectByName holds details about calls to the GetProjectByName method.
		GetProjectByName []struct {
			// Name is the name argument value.
			Name string
		}
		// GetSecurityGroupByName holds details about calls to the GetSecurityGroupByName method.
		GetSecurityGroupByName []struct {
			// Name is the name argument value.
			Name string
			// ProjectId is the projectId argument value.
			ProjectId string
		}
		// GetServerByName holds details about calls to the GetServerByName method.
		GetServerByName []struct {
			// Name is the name argument value.
			Name string
		}
		// GetSubnet holds details about calls to the GetSubnet method.
		GetSubnet []struct {
			// ID is the id argument value.
			ID string
		}
		// GetSubnetByName holds details about calls to the GetSubnetByName method.
		GetSubnetByName []struct {
			// Name is the name argument value.
			Name string
			// NetworkId is the networkId argument value.
			NetworkId string
		}
	}
	lockAssignPort             sync.RWMutex
	lockClients                sync.RWMutex
	lockCreatePort             sync.RWMutex
	lockDeletePort             sync.RWMutex
	lockDetachPort             sync.RWMutex
	lockGetNetworkByName       sync.RWMutex
	lockGetPort                sync.RWMutex
	lockGetPortByIp            sync.RWMutex
	lockGetProjectByName       sync.RWMutex
	lockGetSecurityGroupByName sync.RWMutex
	lockGetServerByName        sync.RWMutex
	lockGetSubnet              sync.RWMutex
	lockGetSubnetByName        sync.RWMutex
}

// AssignPort calls AssignPortFunc.
func (mock *OpenstackClientMock) AssignPort(portId string, serverId string) (*attachinterfaces.Interface, error) {
	if mock.AssignPortFunc == nil {
		panic("OpenstackClientMock.AssignPortFunc: method is nil but OpenstackClient.AssignPort was just called")
	}
	callInfo := struct {
		PortId   string
		ServerId string
	}{
		PortId:   portId,
		ServerId: serverId,
	}
	mock.lockAssignPort.Lock()
	mock.calls.AssignPort = append(mock.calls.AssignPort, callInfo)
	mock.lockAssignPort.Unlock()
	return mock.AssignPortFunc(portId, serverId)
}

// AssignPortCalls gets all the calls that were made to AssignPort.
// Check the length with:
//
//	len(mockedOpenstackClient.AssignPortCalls())
func (mock *OpenstackClientMock) AssignPortCalls() []struct {
	PortId   string
	ServerId string
} {
	var calls []struct {
		PortId   string
		ServerId string
	}
	mock.lockAssignPort.RLock()
	calls = mock.calls.AssignPort
	mock.lockAssignPort.RUnlock()
	return calls
}

// Clients calls ClientsFunc.
func (mock *OpenstackClientMock) Clients() *openstack.ApiClients {
	if mock.ClientsFunc == nil {
		panic("OpenstackClientMock.ClientsFunc: method is nil but OpenstackClient.Clients was just called")
	}
	callInfo := struct {
	}{}
	mock.lockClients.Lock()
	mock.calls.Clients = append(mock.calls.Clients, callInfo)
	mock.lockClients.Unlock()
	return mock.ClientsFunc()
}

// ClientsCalls gets all the calls that were made to Clients.
// Check the length with:
//
//	len(mockedOpenstackClient.ClientsCalls())
func (mock *OpenstackClientMock) ClientsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClients.RLock()
	calls = mock.calls.Clients
	mock.lockClients.RUnlock()
	return calls
}

// CreatePort calls CreatePortFunc.
func (mock *OpenstackClientMock) CreatePort(opts ports.CreateOpts) (*ports.Port, error) {
	if mock.CreatePortFunc == nil {
		panic("OpenstackClientMock.CreatePortFunc: method is nil but OpenstackClient.CreatePort was just called")
	}
	callInfo := struct {
		Opts ports.CreateOpts
	}{
		Opts: opts,
	}
	mock.lockCreatePort.Lock()
	mock.calls.CreatePort = append(mock.calls.CreatePort, callInfo)
	mock.lockCreatePort.Unlock()
	return mock.CreatePortFunc(opts)
}

// CreatePortCalls gets all the calls that were made to CreatePort.
// Check the length with:
//
//	len(mockedOpenstackClient.CreatePortCalls())
func (mock *OpenstackClientMock) CreatePortCalls() []struct {
	Opts ports.CreateOpts
} {
	var calls []struct {
		Opts ports.CreateOpts
	}
	mock.lockCreatePort.RLock()
	calls = mock.calls.CreatePort
	mock.lockCreatePort.RUnlock()
	return calls
}

// DeletePort calls DeletePortFunc.
func (mock *OpenstackClientMock) DeletePort(portId string) error {
	if mock.DeletePortFunc == nil {
		panic("OpenstackClientMock.DeletePortFunc: method is nil but OpenstackClient.DeletePort was just called")
	}
	callInfo := struct {
		PortId string
	}{
		PortId: portId,
	}
	mock.lockDeletePort.Lock()
	mock.calls.DeletePort = append(mock.calls.DeletePort, callInfo)
	mock.lockDeletePort.Unlock()
	return mock.DeletePortFunc(portId)
}

// DeletePortCalls gets all the calls that were made to DeletePort.
// Check the length with:
//
//	len(mockedOpenstackClient.DeletePortCalls())
func (mock *OpenstackClientMock) DeletePortCalls() []struct {
	PortId string
} {
	var calls []struct {
		PortId string
	}
	mock.lockDeletePort.RLock()
	calls = mock.calls.DeletePort
	mock.lockDeletePort.RUnlock()
	return calls
}

// DetachPort calls DetachPortFunc.
func (mock *OpenstackClientMock) DetachPort(portId string, serverId string) error {
	if mock.DetachPortFunc == nil {
		panic("OpenstackClientMock.DetachPortFunc: method is nil but OpenstackClient.DetachPort was just called")
	}
	callInfo := struct {
		PortId   string
		ServerId string
	}{
		PortId:   portId,
		ServerId: serverId,
	}
	mock.lockDetachPort.Lock()
	mock.calls.DetachPort = append(mock.calls.DetachPort, callInfo)
	mock.lockDetachPort.Unlock()
	return mock.DetachPortFunc(portId, serverId)
}

// DetachPortCalls gets all the calls that were made to DetachPort.
// Check the length with:
//
//	len(mockedOpenstackClient.DetachPortCalls())
func (mock *OpenstackClientMock) DetachPortCalls() []struct {
	PortId   string
	ServerId string
} {
	var calls []struct {
		PortId   string
		ServerId string
	}
	mock.lockDetachPort.RLock()
	calls = mock.calls.DetachPort
	mock.lockDetachPort.RUnlock()
	return calls
}

// GetNetworkByName calls GetNetworkByNameFunc.
func (mock *OpenstackClientMock) GetNetworkByName(name string) (*networks.Network, error) {
	if mock.GetNetworkByNameFunc == nil {
		panic("OpenstackClientMock.GetNetworkByNameFunc: method is nil but OpenstackClient.GetNetworkByName was just called")
	}
	callInfo := struct {
		Name string
	}{
		Name: name,
	}
	mock.lockGetNetworkByName.Lock()
	mock.calls.GetNetworkByName = append(mock.calls.GetNetworkByName, callInfo)
	mock.lockGetNetworkByName.Unlock()
	return mock.GetNetworkByNameFunc(name)
}

// GetNetworkByNameCalls gets all the calls that were made to GetNetworkByName.
// Check the length with:
//
//	len(mockedOpenstackClient.GetNetworkByNameCalls())
func (mock *OpenstackClientMock) GetNetworkByNameCalls() []struct {
	Name string
} {
	var calls []struct {
		Name string
	}
	mock.lockGetNetworkByName.RLock()
	calls = mock.calls.GetNetworkByName
	mock.lockGetNetworkByName.RUnlock()
	return calls
}

// GetPort calls GetPortFunc.
func (mock *OpenstackClientMock) GetPort(portId string) (*ports.Port, error) {
	if mock.GetPortFunc == nil {
		panic("OpenstackClientMock.GetPortFunc: method is nil but OpenstackClient.GetPort was just called")
	}
	callInfo := struct {
		PortId string
	}{
		PortId: portId,
	}
	mock.lockGetPort.Lock()
	mock.calls.GetPort = append(mock.calls.GetPort, callInfo)
	mock.lockGetPort.Unlock()
	return mock.GetPortFunc(portId)
}

// GetPortCalls gets all the calls that were made to GetPort.
// Check the length with:
//
//	len(mockedOpenstackClient.GetPortCalls())
func (mock *OpenstackClientMock) GetPortCalls() []struct {
	PortId string
} {
	var calls []struct {
		PortId string
	}
	mock.lockGetPort.RLock()
	calls = mock.calls.GetPort
	mock.lockGetPort.RUnlock()
	return calls
}

// GetPortByIp calls GetPortByIpFunc.
func (mock *OpenstackClientMock) GetPortByIp(ip string) (*ports.Port, error) {
	if mock.GetPortByIpFunc == nil {
		panic("OpenstackClientMock.GetPortByIpFunc: method is nil but OpenstackClient.GetPortByIp was just called")
	}
	callInfo := struct {
		IP string
	}{
		IP: ip,
	}
	mock.lockGetPortByIp.Lock()
	mock.calls.GetPortByIp = append(mock.calls.GetPortByIp, callInfo)
	mock.lockGetPortByIp.Unlock()
	return mock.GetPortByIpFunc(ip)
}

// GetPortByIpCalls gets all the calls that were made to GetPortByIp.
// Check the length with:
//
//	len(mockedOpenstackClient.GetPortByIpCalls())
func (mock *OpenstackClientMock) GetPortByIpCalls() []struct {
	IP string
} {
	var calls []struct {
		IP string
	}
	mock.lockGetPortByIp.RLock()
	calls = mock.calls.GetPortByIp
	mock.lockGetPortByIp.RUnlock()
	return calls
}

// GetProjectByName calls GetProjectByNameFunc.
func (mock *OpenstackClientMock) GetProjectByName(name string) (*projects.Project, error) {
	if mock.GetProjectByNameFunc == nil {
		panic("OpenstackClientMock.GetProjectByNameFunc: method is nil but OpenstackClient.GetProjectByName was just called")
	}
	callInfo := struct {
		Name string
	}{
		Name: name,
	}
	mock.lockGetProjectByName.Lock()
	mock.calls.GetProjectByName = append(mock.calls.GetProjectByName, callInfo)
	mock.lockGetProjectByName.Unlock()
	return mock.GetProjectByNameFunc(name)
}

// GetProjectByNameCalls gets all the calls that were made to GetProjectByName.
// Check the length with:
//
//	len(mockedOpenstackClient.GetProjectByNameCalls())
func (mock *OpenstackClientMock) GetProjectByNameCalls() []struct {
	Name string
} {
	var calls []struct {
		Name string
	}
	mock.lockGetProjectByName.RLock()
	calls = mock.calls.GetProjectByName
	mock.lockGetProjectByName.RUnlock()
	return calls
}

// GetSecurityGroupByName calls GetSecurityGroupByNameFunc.
func (mock *OpenstackClientMock) GetSecurityGroupByName(name string, projectId string) (*groups.SecGroup, error) {
	if mock.GetSecurityGroupByNameFunc == nil {
		panic("OpenstackClientMock.GetSecurityGroupByNameFunc: method is nil but OpenstackClient.GetSecurityGroupByName was just called")
	}
	callInfo := struct {
		Name      string
		ProjectId string
	}{
		Name:      name,
		ProjectId: projectId,
	}
	mock.lockGetSecurityGroupByName.Lock()
	mock.calls.GetSecurityGroupByName = append(mock.calls.GetSecurityGroupByName, callInfo)
	mock.lockGetSecurityGroupByName.Unlock()
	return mock.GetSecurityGroupByNameFunc(name, projectId)
}

// GetSecurityGroupByNameCalls gets all the calls that were made to GetSecurityGroupByName.
// Check the length with:
//
//	len(mockedOpenstackClient.GetSecurityGroupByNameCalls())
func (mock *OpenstackClientMock) GetSecurityGroupByNameCalls() []struct {
	Name      string
	ProjectId string
} {
	var calls []struct {
		Name      string
		ProjectId string
	}
	mock.lockGetSecurityGroupByName.RLock()
	calls = mock.calls.GetSecurityGroupByName
	mock.lockGetSecurityGroupByName.RUnlock()
	return calls
}

// GetServerByName calls GetServerByNameFunc.
func (mock *OpenstackClientMock) GetServerByName(name string) (*servers.Server, error) {
	if mock.GetServerByNameFunc == nil {
		panic("OpenstackClientMock.GetServerByNameFunc: method is nil but OpenstackClient.GetServerByName was just called")
	}
	callInfo := struct {
		Name string
	}{
		Name: name,
	}
	mock.lockGetServerByName.Lock()
	mock.calls.GetServerByName = append(mock.calls.GetServerByName, callInfo)
	mock.lockGetServerByName.Unlock()
	return mock.GetServerByNameFunc(name)
}

// GetServerByNameCalls gets all the calls that were made to GetServerByName.
// Check the length with:
//
//	len(mockedOpenstackClient.GetServerByNameCalls())
func (mock *OpenstackClientMock) GetServerByNameCalls() []struct {
	Name string
} {
	var calls []struct {
		Name string
	}
	mock.lockGetServerByName.RLock()
	calls = mock.calls.GetServerByName
	mock.lockGetServerByName.RUnlock()
	return calls
}

// GetSubnet calls GetSubnetFunc.
func (mock *OpenstackClientMock) GetSubnet(id string) (*subnets.Subnet, error) {
	if mock.GetSubnetFunc == nil {
		panic("OpenstackClientMock.GetSubnetFunc: method is nil but OpenstackClient.GetSubnet was just called")
	}
	callInfo := struct {
		ID string
	}{
		ID: id,
	}
	mock.lockGetSubnet.Lock()
	mock.calls.GetSubnet = append(mock.calls.GetSubnet, callInfo)
	mock.lockGetSubnet.Unlock()
	return mock.GetSubnetFunc(id)
}

// GetSubnetCalls gets all the calls that were made to GetSubnet.
// Check the length with:
//
//	len(mockedOpenstackClient.GetSubnetCalls())
func (mock *OpenstackClientMock) GetSubnetCalls() []struct {
	ID string
} {
	var calls []struct {
		ID string
	}
	mock.lockGetSubnet.RLock()
	calls = mock.calls.GetSubnet
	mock.lockGetSubnet.RUnlock()
	return calls
}

// GetSubnetByName calls GetSubnetByNameFunc.
func (mock *OpenstackClientMock) GetSubnetByName(name string, networkId string) (*subnets.Subnet, error) {
	if mock.GetSubnetByNameFunc == nil {
		panic("OpenstackClientMock.GetSubnetByNameFunc: method is nil but OpenstackClient.GetSubnetByName was just called")
	}
	callInfo := struct {
		Name      string
		NetworkId string
	}{
		Name:      name,
		NetworkId: networkId,
	}
	mock.lockGetSubnetByName.Lock()
	mock.calls.GetSubnetByName = append(mock.calls.GetSubnetByName, callInfo)
	mock.lockGetSubnetByName.Unlock()
	return mock.GetSubnetByNameFunc(name, networkId)
}

// GetSubnetByNameCalls gets all the calls that were made to GetSubnetByName.
// Check the length with:
//
//	len(mockedOpenstackClient.GetSubnetByNameCalls())
func (mock *OpenstackClientMock) GetSubnetByNameCalls() []struct {
	Name      string
	NetworkId string
} {
	var calls []struct {
		Name      string
		NetworkId string
	}
	mock.lockGetSubnetByName.RLock()
	calls = mock.calls.GetSubnetByName
	mock.lockGetSubnetByName.RUnlock()
	return calls
}
