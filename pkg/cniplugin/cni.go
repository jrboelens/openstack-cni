package cniplugin

import (
	"fmt"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	currentcni "github.com/containernetworking/cni/pkg/types/040"
	cniversion "github.com/containernetworking/cni/pkg/version"
	"github.com/jboelensns/openstack-cni/pkg/cniclient"
	"github.com/jboelensns/openstack-cni/pkg/cnistate"
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
	if _, err := me.client.HandleResponse(me.client.CniCommand(cmd)); err != nil {
		return err
	}
	err := me.client.DeleteState(cmd.ContainerID, cmd.IfName)
	if err == cniclient.ErrStateNotFound {
		return nil
	}
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
	ifaceName, err := util.GetIfaceNameByMac(mac)
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

	context, err := util.NewCniContext(cmd)
	if err != nil {
		return err
	}

	info := &cnistate.IfaceInfo{
		ContainerId: cmd.ContainerID,
		Ifname:      iface.DestName,
		Netns:       cmd.Netns,
		IpAddress:   iface.Address.IP.String(),
		PodName:     context.GetArg("K8S_POD_NAME"),
		Namespace:   context.GetArg("K8S_POD_NAMESPACE"),
	}

	return me.client.SetState(info)
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
