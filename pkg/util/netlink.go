package util

import (
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
	"github.com/vishvananda/netns"
)

//go:generate moq -pkg mocks -out ../fixtures/mocks/util_mocks.go . NetlinkWrapper

type NetLinkWrapperOpts struct {
	ErrorMessageReporting bool
}

// NetlinkWrapper allows us to test without actually using netlink
type NetlinkWrapper interface {
	AddrAdd(link netlink.Link, addr *netlink.Addr) error
	AddrReplace(link netlink.Link, addr *netlink.Addr) error
	GetNetNsIdByPath(namespace string) (int, error)
	GetNetNsIdByPid(pid int) (int, error)
	LinkByIndex(index int) (netlink.Link, error)
	LinkByName(ifname string) (netlink.Link, error)
	LinkSetDown(link netlink.Link) error
	LinkSetName(link netlink.Link, name string) error
	LinkSetNsFd(link netlink.Link, fd int) error
	LinkSetUp(link netlink.Link) error
}

func NewNetlinkWrapper() *netlinkWrapper {
	return &netlinkWrapper{}
}

func NewNetlinkWrapperWithOpts(opts NetLinkWrapperOpts) *netlinkWrapper {
	w := &netlinkWrapper{}
	if opts.ErrorMessageReporting {
		w.ErrorMessageReporting(opts.ErrorMessageReporting)
	}
	return w
}

type netlinkWrapper struct {
}

// ErrorMessageReporting enables and disables NETLINK_EXT_ACK error reporting
func (me *netlinkWrapper) ErrorMessageReporting(enable bool) {
	nl.EnableErrorMessageReporting = enable
}

func (me *netlinkWrapper) GetNetNsIdByPath(namespace string) (int, error) {
	handle, err := netns.GetFromPath(namespace)
	if err != nil {
		return 0, err
	}
	return int(handle), nil
}

func (me *netlinkWrapper) AddrAdd(link netlink.Link, addr *netlink.Addr) error {
	return netlink.AddrAdd(link, addr)
}

func (me *netlinkWrapper) AddrReplace(link netlink.Link, addr *netlink.Addr) error {
	return netlink.AddrReplace(link, addr)
}

func (me *netlinkWrapper) GetNetNsIdByPid(pid int) (int, error) {
	return netlink.GetNetNsIdByPid(pid)
}

func (me *netlinkWrapper) LinkByIndex(index int) (netlink.Link, error) {
	return netlink.LinkByIndex(index)
}

func (me *netlinkWrapper) LinkByName(ifname string) (netlink.Link, error) {
	return netlink.LinkByName(ifname)
}

func (me *netlinkWrapper) LinkSetDown(link netlink.Link) error {
	return netlink.LinkSetDown(link)
}

func (me *netlinkWrapper) LinkSetName(link netlink.Link, name string) error {
	return netlink.LinkSetName(link, name)
}

func (me *netlinkWrapper) LinkSetNsFd(link netlink.Link, fd int) error {
	return netlink.LinkSetNsFd(link, fd)
}

func (me *netlinkWrapper) LinkSetUp(link netlink.Link) error {
	return netlink.LinkSetUp(link)
}
