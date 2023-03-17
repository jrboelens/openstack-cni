package util

import (
	"net"
)

func GetIpNetFromAddress(ipaddress string) (*net.IPNet, error) {
	ip, cidr, err := net.ParseCIDR(ipaddress)
	if err != nil {
		return nil, err
	}
	return &net.IPNet{IP: ip, Mask: cidr.Mask}, nil
}
