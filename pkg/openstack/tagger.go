package openstack

import (
	"strings"

	"github.com/gophercloud/gophercloud"
)

// NeutronTagger provides the ability to manipulate tags on all Neutron resources
// https://docs.openstack.org/neutron/latest/contributor/internals/tag.html
type NeutronTagger struct {
	networkClient *gophercloud.ServiceClient
	resource_type NeutronResourceType
}

type NeutronResourceType uint64

const (
	Tags NeutronResourceType = iota
	FloatingIps
	Networks
	NetworkSegmentRanges
	Policies
	Ports
	Routers
	SecurityGroups
	SubnetPools
	Subnets
	Trunks
)

func (me NeutronResourceType) String() string {
	switch me {
	case Tags:
		return "tags"
	case FloatingIps:
		return "floatingips"
	case Networks:
		return "networks"
	case NetworkSegmentRanges:
		return "network_segment_ranges"
	case Policies:
		return "policies"
	case Ports:
		return "ports"
	case Routers:
		return "routers"
	case SecurityGroups:
		return "security_groups"
	case SubnetPools:
		return "subnetpools"
	case Subnets:
		return "subnets"
	case Trunks:
		return "trunks"
	}
	return "unknown"
}

// NewNeutronTagger creates a new instance of a *NeutronTagger
func NewNeutronTagger(client *gophercloud.ServiceClient, resource_type NeutronResourceType) *NeutronTagger {
	return &NeutronTagger{client, resource_type}
}

// NeutronTag represents a single tag
type NeutronTag string

// NeutronTags represents a set of tags
type NeutronTags struct {
	Tags []NeutronTag `json:"tags"`
}

// Create creates a single tag for network resource
func (me *NeutronTagger) Create(id, tag string) error {
	responseBody := make(map[string]any)
	opts := &gophercloud.RequestOpts{OkCodes: []int{200, 201}}
	_, _, err := gophercloud.ParseResponse(me.networkClient.Put(me.url(id, tag), nil, &responseBody, opts))
	if err != nil {
		return err
	}

	return nil
}

// SetAll overwrites existing tags for a network resource with the provided tags
func (me *NeutronTagger) SetAll(id string, tags NeutronTags) error {
	responseBody := make(map[string]any)
	opts := &gophercloud.RequestOpts{OkCodes: []int{200, 201}}
	_, _, err := gophercloud.ParseResponse(me.networkClient.Put(me.url(id, ""), tags, &responseBody, opts))
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a single tag from a network resource
func (me *NeutronTagger) Delete(id, tag string) (err error) {
	opts := &gophercloud.RequestOpts{OkCodes: []int{204}}
	_, err = me.networkClient.Delete(me.url(id, tag), opts)
	return
}

// DeleteAll deletes all tags for a network resource
func (me *NeutronTagger) DeleteAll(id string) (err error) {
	opts := &gophercloud.RequestOpts{OkCodes: []int{204}}
	_, err = me.networkClient.Delete(me.url(id, ""), opts)
	return
}

// Exists returns true if a tag exists for a network resource; otherwise false
func (me *NeutronTagger) Exists(id, tag string) (bool, error) {
	opts := &gophercloud.RequestOpts{OkCodes: []int{204}}
	var tags NeutronTags
	_, _, err := gophercloud.ParseResponse(me.networkClient.Get(me.url(id, tag), &tags, opts))
	if err != nil {
		if strings.Contains(err.Error(), "could not be found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetAll retruns all tags for a network resource
func (me *NeutronTagger) GetAll(id string) (tags NeutronTags, err error) {
	opts := &gophercloud.RequestOpts{OkCodes: []int{200, 201}}
	_, _, err = gophercloud.ParseResponse(me.networkClient.Get(me.url(id, ""), &tags, opts))
	return
}

// url provides a simple way to create service urls for each tag related action
func (me *NeutronTagger) url(resId, tagId string) string {
	parts := []string{me.resource_type.String()}
	if resId != "" {
		parts = append(parts, resId)
	}
	parts = append(parts, "tags")
	if tagId != "" {
		parts = append(parts, tagId)
	}
	return me.networkClient.ServiceURL(parts...)
}
