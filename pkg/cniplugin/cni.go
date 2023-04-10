package cniplugin

import (
	"fmt"
	"net"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	currentcni "github.com/containernetworking/cni/pkg/types/040"
	cniversion "github.com/containernetworking/cni/pkg/version"
	"github.com/jboelensns/openstack-cni/pkg/cniclient"
	"github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

type Cni struct {
	client *cniclient.Client
	nw     Networking
}

func NewCni(client *cniclient.Client, nw Networking) *Cni {
	return &Cni{
		client: client,
		nw:     nw,
	}
}

func (me *Cni) Add(args *skel.CmdArgs) error {
	var netConf types.NetConf
	if err := util.FromJson(args.StdinData, &netConf); err != nil {
		return err
	}

	cmd := CniCommandFromSkelArgs("ADD", args)
	body, err := me.client.HandleResponse(me.client.CniCommand(cmd))
	if err != nil {
		return err
	}

	var result currentcni.Result
	if err := util.FromJson(body, &result); err != nil {
		return err
	}

	if err := me.ConfigureInterface(cmd, &result); err != nil {
		return err
	}

	finalResult, err := result.GetAsVersion(netConf.CNIVersion)
	if err != nil {
		return err
	}

	return finalResult.Print()
}

func (me *Cni) Check(args *skel.CmdArgs) error {
	cmd := CniCommandFromSkelArgs("CHECK", args)
	_, err := me.client.HandleResponse(me.client.CniCommand(cmd))
	return err
}

func (me *Cni) Del(args *skel.CmdArgs) error {
	cmd := CniCommandFromSkelArgs("DEL", args)
	_, err := me.client.HandleResponse(me.client.CniCommand(cmd))
	return err
}

func (me *Cni) Invoke() error {
	err := skel.PluginMainWithError(
		func(args *skel.CmdArgs) error {
			return me.Add(args)
		},
		func(args *skel.CmdArgs) error {
			return me.Check(args)
		},
		func(args *skel.CmdArgs) error {
			return me.Del(args)
		},
		cniversion.All,
		"openstack CNI plugin that plumbs neutron ports into containers")

	if err != nil {
		if perr := err.Print(); perr != nil {
			logging.Error("error writing CNI error JSON to stdout", err)
		}
		return err
	}
	return nil
}

func (me *Cni) ConfigureInterface(cmd util.CniCommand, result *currentcni.Result) error {

	mac := result.Interfaces[0].Mac
	ifaceName, err := GetIfaceNameByMac(mac)
	if err != nil {
		return err
	}

	iface := &NetworkInterface{
		Name:     ifaceName,
		DestName: result.Interfaces[0].Name,
		Address:  &result.IPs[0].Address,
	}

	err = me.nw.Configure(cmd.Netns, iface)
	if err != nil {
		return fmt.Errorf("failed to configure interface %w", err)
	}
	return err
}

func CniCommandFromSkelArgs(cmdStr string, args *skel.CmdArgs) util.CniCommand {
	return util.CniCommand{
		Command:     cmdStr,
		ContainerID: args.ContainerID,
		Netns:       args.Netns,
		IfName:      args.IfName,
		Args:        args.Args,
		Path:        args.Path,
		StdinData:   args.StdinData,
	}
}

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
