package openstack

import (
	"context"
	"fmt"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

var _ OpenstackClient = &CachedClient{}

type CachedClient struct {
	OpenstackClient
	Expiration time.Duration
	cash       *cache.Cache[string, any]
	cancelFunc context.CancelFunc
}

func NewCachedClient(client OpenstackClient, expiration time.Duration) OpenstackClient {
	ctx, cancel := context.WithCancel(context.Background())
	cash := cache.NewContext[string, any](ctx)
	return &CachedClient{OpenstackClient: client, Expiration: expiration, cash: cash, cancelFunc: cancel}
}

func (me *CachedClient) Stop() {
	me.cancelFunc()
}

func (me *CachedClient) GetNetworkByName(name string) (*networks.Network, error) {
	return getValue[networks.Network](me.cash, makeKey("GetNetworkByName", name), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetNetworkByName(name)
	})
}

func (me *CachedClient) GetPort(portId string) (*ports.Port, error) {
	return getValue[ports.Port](me.cash, makeKey("GetPort", portId), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetPort(portId)
	})
}

func (me *CachedClient) GetPortByTags(tags []string) (*ports.Port, error) {
	return getValue[ports.Port](me.cash, makeKey("GetPortByTags", strings.Join(tags, ",")), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetPortByTags(tags)
	})
}

func (me *CachedClient) GetProjectByName(name string) (*projects.Project, error) {
	return getValue[projects.Project](me.cash, makeKey("GetProjectByName", name), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetProjectByName(name)
	})
}

func (me *CachedClient) GetServerByName(name string) (*servers.Server, error) {
	return getValue[servers.Server](me.cash, makeKey("GetServerByName", name), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetServerByName(name)
	})
}

func (me *CachedClient) GetSecurityGroupByName(name string, projectId string) (*groups.SecGroup, error) {
	return getValue[groups.SecGroup](me.cash, makeKey("GetSecurityGroupsByName", name), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetSecurityGroupByName(name, projectId)
	})
}

func (me *CachedClient) GetSubnet(id string) (*subnets.Subnet, error) {
	return getValue[subnets.Subnet](me.cash, makeKey("GetSubnet", id), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetSubnet(id)
	})
}

func (me *CachedClient) GetSubnetByName(name, networkId string) (*subnets.Subnet, error) {
	return getValue[subnets.Subnet](me.cash, makeKey("GetSubnetByName", fmt.Sprintf("%s--%s", name, networkId)), me.Expiration, func() (any, error) {
		return me.OpenstackClient.GetSubnetByName(name, networkId)
	})
}

func makeKey(parts ...string) string {
	return strings.Join(parts, "|")
}

func getValue[T any](store *cache.Cache[string, any], cacheKey string, expiraton time.Duration, fn func() (any, error)) (*T, error) {
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
