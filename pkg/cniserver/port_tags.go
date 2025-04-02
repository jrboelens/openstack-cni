package cniserver

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

const OPENSTACK_CNI_TAG = "openstack-cni=true"
const TAG_CONTAINERID_NAME = "containerid"
const TAG_IFNAME_NAME = "ifname"
const TAG_NETNS_NAME = "netns"
const TAG_HOST_NAME = "host"

type PortTags struct {
	ContainerId string
	IfName      string
	Netns       string
	Host        string
}

// NewPortTags returns a new PortTags
func NewPortTags(port ports.Port) PortTags {
	return PortTags{
		ContainerId: GetPortTagWithDefault(TAG_CONTAINERID_NAME, "", port.Tags),
		IfName:      GetPortTagWithDefault(TAG_IFNAME_NAME, "", port.Tags),
		Netns:       GetPortTagWithDefault(TAG_NETNS_NAME, "", port.Tags),
		Host:        NewHostTag(),
	}
}

// NewPortTagsFromCommand creates a PortTags including container, interface and namespace data
func NewPortTagsFromCommand(cmd util.CniCommand) PortTags {
	containerId := cmd.ContainerID
	if len(containerId) > 12 {
		containerId = containerId[0:12]
	}

	return PortTags{
		ContainerId: containerId,
		IfName:      cmd.IfName,
		Netns:       cmd.Netns,
		Host:        NewHostTag(),
	}
}

func (me PortTags) NeutronTags() openstack.NeutronTags {
	return openstack.NewNeutronTags(
		fmt.Sprintf("%s=%s", TAG_CONTAINERID_NAME, me.ContainerId),
		fmt.Sprintf("%s=%s", TAG_IFNAME_NAME, me.IfName),
		fmt.Sprintf("%s=%s", TAG_NETNS_NAME, me.Netns),
		fmt.Sprintf(OPENSTACK_CNI_TAG),
		fmt.Sprintf("%s=%s", TAG_HOST_NAME, me.Host),
	)
}

func NewHostTag() string {
	hostname, _ := util.GetHostname()
	return fmt.Sprintf("%s=%s", TAG_HOST_NAME, hostname)
}

func NewPortKeyTags() []string {
	return []string{
		OPENSTACK_CNI_TAG,
		NewHostTag(),
	}
}

func HasOpenstackCniTag(tags []string) bool {
	for _, tag := range tags {
		if tag == OPENSTACK_CNI_TAG {
			return true
		}
	}
	return false
}

// Returns a port tag by name
func GetPortTag(name string, tags []string) *string {
	for _, tag := range tags {
		parts := strings.Split(tag, "=")
		if len(parts) > 1 {
			if parts[0] == name {
				return &parts[1]
			}
		}
	}
	return nil
}

// Returns a port tag by name or the default value
func GetPortTagWithDefault(name string, def string, tags []string) string {
	value := GetPortTag(name, tags)
	if value != nil {
		return *value
	}
	return def
}
