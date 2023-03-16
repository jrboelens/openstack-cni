package util

import (
	"fmt"
	"net"
)

func GetIfaceNameByMac(mac string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.HardwareAddr.String() == mac {
			return iface.Name, nil
		}
	}

	return "", fmt.Errorf("failed to find interface for %s", mac)
}

func GetIpNetFromAddress(ipaddress string) (*net.IPNet, error) {
	ip, cidr, err := net.ParseCIDR(ipaddress)
	if err != nil {
		return nil, err
	}
	return &net.IPNet{IP: ip, Mask: cidr.Mask}, nil
}

func IpnetFromIp(ip string) net.IPNet {
	ipaddr := net.ParseIP(ip)
	return net.IPNet{IP: ipaddr, Mask: ipaddr.DefaultMask()}
}
