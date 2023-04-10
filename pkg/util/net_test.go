package util_test

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"

	. "github.com/pepinns/go-hamcrest"
)

func Test_IpParsing(t *testing.T) {
	t.Run("from string", func(t *testing.T) {
		address := "198.18.182.36/24"
		ip, err := util.GetIpNetFromAddress(address)
		Assert(t).That(err, IsNil())
		Assert(t).That(ip.String(), Equals(address))
	})

	t.Run("from port result", func(t *testing.T) {
		pr := &openstack.SetupPortResult{
			Subnet: &subnets.Subnet{
				ID:   "MYID",
				CIDR: "198.18.182.0/24",
			},
			Port: &ports.Port{
				FixedIPs: []ports.IP{
					{SubnetID: "MYID", IPAddress: "198.18.182.36"},
				},
			},
		}
		address := "198.18.182.36/24"

		ip, err := pr.GetIp()
		Assert(t).That(err, IsNil())
		Assert(t).That(ip.String(), Equals(address))
	})
}
