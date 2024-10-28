package openstack

import (
	"context"
	"fmt"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

var _ OpenstackClient = &CachedClient{}

type CachedClient struct {
	OpenstackClient OpenstackClient
	Expiration      time.Duration
	cash            *cache.Cache[string, any]
	cancelFunc      context.CancelFunc
}

func NewCachedClient(client OpenstackClient, expiration time.Duration) OpenstackClient {
	ctx, cancel := context.WithCancel(context.Background())
	cash := cache.NewContext[string, any](ctx)
	return &CachedClient{OpenstackClient: client,
		Expiration: expiration,
		cash:       cash,
		cancelFunc: cancel,
	}
}

func (me *CachedClient) Stop() {
	me.cancelFunc()
}

// AssignPort attaches a port to a server
func (me *CachedClient) AssignPort(portId, serverId string) (*attachinterfaces.Interface, error) {
	return me.OpenstackClient.AssignPort(portId, serverId)
}

func (me *CachedClient) Clients() *ApiClients {
	return me.OpenstackClient.Clients()
}

// CreatePort creates a neutron port inside of the specified network
func (me *CachedClient) CreatePort(opts ports.CreateOpts) (*ports.Port, error) {
	return me.OpenstackClient.CreatePort(opts)
}

// DeletePort deletes the port
func (me *CachedClient) DeletePort(portId string) error {
	// find all of the ports with the portId and delete them from the cache
	//
	// Note: this isn't he most efficient way to go about it, but it's significantly easier
	// implement and understand
	keys := me.cash.Keys()
	for _, key := range keys {
		// only look at keys having to do with Ports
		if !strings.HasPrefix(key, "GetPort") {
			continue
		}

		// if we find a port with a matching id, delete it
		// if we find a slice of ports, deleting the one that has a matchind id
		val, ok := me.cash.Get(key)
		if !ok {
			continue
		}

		switch v := val.(type) {
		case *ports.Port:
			if v.ID == portId {
				me.cash.Delete(key)
			}
		case []ports.Port:
			for _, port := range v {
				if port.ID == portId {
					me.cash.Delete(key)
				}
			}
		}
	}

	return me.OpenstackClient.DeletePort(portId)
}

// Detach port removes a port's relationship from a server
func (me *CachedClient) DetachPort(portId, serverId string) error {
	return me.OpenstackClient.DetachPort(portId, serverId)
}

func (me *CachedClient) GetNetworkByName(name string) (*networks.Network, error) {
	return getPtrValue[networks.Network](me.cash, makeKey("GetNetworkByName", name), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetNetworkByName(name)
	})
}

func (me *CachedClient) GetPort(portId string) (*ports.Port, error) {
	return getPtrValue[ports.Port](me.cash, makeKey("GetPort", portId), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetPort(portId)
	})
}

func (me *CachedClient) GetPortsByDeviceId(deviceId string) ([]ports.Port, error) {
	return getValue[[]ports.Port](me.cash, makeKey("GetPortsByDeviceId", deviceId), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetPortsByDeviceId(deviceId)
	})
}

func (me *CachedClient) GetPortByTags(tags []string) (*ports.Port, error) {
	return me.OpenstackClient.GetPortByTags(tags)
}

func (me *CachedClient) GetPortsByTags(tags []string) ([]ports.Port, error) {
	return me.OpenstackClient.GetPortsByTags(tags)
}

func (me *CachedClient) GetProjectByName(name string) (*projects.Project, error) {
	return getPtrValue[projects.Project](me.cash, makeKey("GetProjectByName", name), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetProjectByName(name)
	})
}

func (me *CachedClient) GetServerByName(name string) (*servers.Server, error) {
	return getPtrValue[servers.Server](me.cash, makeKey("GetServerByName", name), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetServerByName(name)
	})
}

func (me *CachedClient) GetSecurityGroupByName(name string, projectId string) (*groups.SecGroup, error) {
	return getPtrValue[groups.SecGroup](me.cash, makeKey("GetSecurityGroupsByName", name), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetSecurityGroupByName(name, projectId)
	})
}

func (me *CachedClient) GetSubnet(id string) (*subnets.Subnet, error) {
	return getPtrValue[subnets.Subnet](me.cash, makeKey("GetSubnet", id), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetSubnet(id)
	})
}

func (me *CachedClient) GetSubnetByName(name, networkId string) (*subnets.Subnet, error) {
	return getPtrValue[subnets.Subnet](me.cash, makeKey("GetSubnetByName", fmt.Sprintf("%s--%s", name, networkId)), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetSubnetByName(name, networkId)
	})
}

func makeKey(parts ...string) string {
	return strings.Join(parts, "|")
}

func getPtrValue[T any](store *cache.Cache[string, any], cacheKey string, expiraton time.Duration, fn func() (any, error)) (*T, error) {
	cachedVal, found := store.Get(cacheKey)
	if found {
		return cachedVal.(*T), nil
	}

	val, err := fn()
	if err != nil {
		return nil, err
	}

	store.Set(cacheKey, val, cache.WithExpiration(expiraton))
	return val.(*T), nil
}

func getValue[T any](store *cache.Cache[string, any], cacheKey string, expiraton time.Duration, fn func() (any, error)) (T, error) {
	cachedVal, found := store.Get(cacheKey)
	if found {
		return cachedVal.(T), nil
	}

	val, err := fn()
	if err != nil {
		return *new(T), err
	}

	store.Set(cacheKey, val, cache.WithExpiration(expiraton))
	return val.(T), nil
}
