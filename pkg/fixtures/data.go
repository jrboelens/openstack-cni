package fixtures

import (
	"net"
	"time"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	currentcni "github.com/containernetworking/cni/pkg/types/040"
	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

type TestData struct {
}

func NewTestData() *TestData {
	return &TestData{}
}

func (me *TestData) SkelArgs() *skel.CmdArgs {
	cmd := me.CniCommand()

	return &skel.CmdArgs{
		ContainerID: cmd.ContainerID,
		Netns:       cmd.Netns,
		IfName:      cmd.IfName,
		Args:        cmd.Args,
		Path:        cmd.Path,
		StdinData:   cmd.StdinData,
	}
}

func (me *TestData) CniCommand() util.CniCommand {
	return util.CniCommand{
		Command:     "ADD",
		ContainerID: "3369ae15e741d31e8616906642c3ca309291e7776e2fff3bb8d379e642e056a8",
		Netns:       "/proc/4242/net/ns",
		IfName:      "eth37",
		Args:        "FOO=BAR",
		Path:        "/opt/cni/bin:/opt/cni/bin",
		StdinData:   me.Stdin(),
	}
}

func (me *TestData) Stdin() []byte {
	return []byte(`{
        "cniVersion": "0.3.1",
        "type": "openstack-cni",
        "name": "service-ingress",
        "network": "devint-dp-compute-internal",
        "asecurity_groups": ["dp_default", "default"]
        }
	`)
}

func (me *TestData) CniResult() *currentcni.Result {
	getIpNet := func(ipWithCidr string) net.IPNet {
		ipNet, _ := util.GetIpNetFromAddress(ipWithCidr)
		return *ipNet
	}
	zero := 0

	return &currentcni.Result{
		CNIVersion: "0.3.1",
		Interfaces: []*currentcni.Interface{
			{
				Name:    "ens3",
				Mac:     "02:42:d9:1f:22:9d",
				Sandbox: "/proc/4237/net/ns"},
		},
		IPs: []*currentcni.IPConfig{
			{
				Version:   "4",
				Interface: &zero,
				Address:   getIpNet("192.168.1.42/32"),
				Gateway:   net.ParseIP("192.168.0.1")},
		},
		Routes: []*types.Route{
			{
				Dst: getIpNet("192.168.1.100/32"),
				GW:  net.ParseIP("192.168.0.1")},
		},
		DNS: types.DNS{
			Nameservers: []string{"1.1.1.1", "8.8.8.8"},
			// Domain:      "",
			// Search:      []string{},
			// Options:     []string{},
		},
	}
}

func PortReaperOpts() cniserver.PortReaperOpts {
	return cniserver.PortReaperOpts{Interval: time.Second * 300, MinPortAge: time.Second * 300}
}

func NeutronTags() []string {
	return []string{"foo=bar", "openstack-cni=true", "netns=/proc/1234/ns"}
}
