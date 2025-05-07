package util

import (
	"time"

	"github.com/cloudflare/backoff"
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

type RetryOpts struct {
	BackoffMax      time.Duration
	BackoffInterval time.Duration
	MaxWaitTime     time.Duration
	RetryErrors     bool
}

func DefaultRetryOpts() RetryOpts {
	return RetryOpts{
		BackoffMax:      time.Duration(time.Second * 5),
		BackoffInterval: time.Duration(time.Millisecond * 150),
		MaxWaitTime:     time.Duration(time.Second * 5),
		RetryErrors:     true,
	}
}

type netlinkWithRetry struct {
	nl   NetlinkWrapper
	opts RetryOpts
}

// NewNetlinkWithRetry wraps a NetlinkWrapper with exponential backoff behavior
func NewNetlinkWithRetry(nl NetlinkWrapper, opts RetryOpts) *netlinkWithRetry {
	return &netlinkWithRetry{opts: opts, nl: nl}
}

func (me *netlinkWithRetry) GetNetNsIdByPath(namespace string) (int, error) {
	return execWithResult(me.opts, func() (int, error) {
		handle, err := netns.GetFromPath(namespace)
		if err != nil {
			return 0, err
		}
		return int(handle), nil
	})
}

func (me *netlinkWithRetry) AddrAdd(link netlink.Link, addr *netlink.Addr) error {
	return execWithError(me.opts, func() error { return netlink.AddrAdd(link, addr) })
}

func (me *netlinkWithRetry) AddrReplace(link netlink.Link, addr *netlink.Addr) error {
	return execWithError(me.opts, func() error { return netlink.LinkSetDown(link) })
}

func (me *netlinkWithRetry) GetNetNsIdByPid(pid int) (int, error) {
	return execWithResult(me.opts, func() (int, error) {
		return netlink.GetNetNsIdByPid(pid)
	})
}

func (me *netlinkWithRetry) LinkByIndex(index int) (netlink.Link, error) {
	return execWithResult(me.opts, func() (netlink.Link, error) {
		return netlink.LinkByIndex(index)
	})
}

func (me *netlinkWithRetry) LinkByName(ifname string) (netlink.Link, error) {
	return execWithResult(me.opts, func() (netlink.Link, error) {
		return netlink.LinkByName(ifname)
	})
}

func (me *netlinkWithRetry) LinkSetDown(link netlink.Link) error {
	return execWithError(me.opts, func() error { return netlink.LinkSetDown(link) })
}

func (me *netlinkWithRetry) LinkSetName(link netlink.Link, name string) error {
	return execWithError(me.opts, func() error { return netlink.LinkSetName(link, name) })
}

func (me *netlinkWithRetry) LinkSetNsFd(link netlink.Link, fd int) error {
	return execWithError(me.opts, func() error { return netlink.LinkSetNsFd(link, fd) })
}

func (me *netlinkWithRetry) LinkSetUp(link netlink.Link) error {
	return execWithError(me.opts, func() error { return netlink.LinkSetUp(link) })
}

// execWithError provides a convienent way to backoff functions that return an error
func execWithError(opts RetryOpts, callback func() error) error {
	_, e := Backoff(opts, func() (struct{}, error) {
		e := callback()
		return *new(struct{}), e
	})
	return e
}

// execWithError provides a convienent way to backoff functions that return an error and a result
func execWithResult[R any](opts RetryOpts, callback func() (R, error)) (R, error) {
	return Backoff(opts, func() (R, error) { return callback() })
}

func Backoff[R any](opts RetryOpts, callback func() (R, error)) (R, error) {
	start_time := time.Now()
	bo := backoff.New(opts.BackoffMax, opts.BackoffInterval)
	var err error
	var r R
	for {
		r, err = callback()
		if !opts.RetryErrors {
			if err != nil {
				return r, err
			}
		}
		if err == nil {
			return r, err
		}

		now := time.Now()
		if now.Sub(start_time) > opts.MaxWaitTime {
			return r, err
		}

		<-time.After(bo.Duration())
	}
}
