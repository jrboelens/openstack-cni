package cniplugin

import (
	"fmt"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	currentcni "github.com/containernetworking/cni/pkg/types/040"
	cniversion "github.com/containernetworking/cni/pkg/version"
	"github.com/jboelensns/openstack-cni/pkg/cniclient"
	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	"github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

// Cni provides methods with the ability accept CNI spec data, make requests to the openstack-cni-daemon and return the results
type Cni struct {
	client *cniclient.Client
	nw     Networking
}

// NewCni returns a new Cni
func NewCni(client *cniclient.Client, nw Networking) *Cni {
	return &Cni{
		client: client,
		nw:     nw,
	}
}

// Add handles ADD CNI commands
func (me *Cni) Add(args *skel.CmdArgs) error {
	logging.Log().Info().Str("container_id", args.ContainerID).Str("ns", args.Netns).Str("iface", args.IfName).Str("args", args.Args).Str("path", args.Path).Msg("received ADD")
	var netConf types.NetConf
	if err := util.FromJson(args.StdinData, &netConf); err != nil {
		return err
	}

	cmd := cniCommandFromSkelArgs(cniserver.CommandAdd, args)
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

// Check handles CHECK CNI commands
func (me *Cni) Check(args *skel.CmdArgs) error {
	logging.Log().Info().Str("container_id", args.ContainerID).Str("ns", args.Netns).Str("iface", args.IfName).Str("args", args.Args).Str("path", args.Path).Msg("received CHECK")
	cmd := cniCommandFromSkelArgs(cniserver.CommandCheck, args)
	_, err := me.client.HandleResponse(me.client.CniCommand(cmd))
	return err
}

// Del handles DEL CNI commands
func (me *Cni) Del(args *skel.CmdArgs) error {
	logging.Log().Info().Str("container_id", args.ContainerID).Str("ns", args.Netns).Str("iface", args.IfName).Str("args", args.Args).Str("path", args.Path).Msg("received DEL")
	cmd := cniCommandFromSkelArgs(cniserver.CommandDel, args)
	_, err := me.client.HandleResponse(me.client.CniCommand(cmd))
	return err
}

// Invoke invokes the CNI plugin skeletons using its own methods
func (me *Cni) Invoke() error {
	err := skel.PluginMainWithError(
		func(args *skel.CmdArgs) error {
			err := me.Add(args)
			if err != nil {
				logging.Error(fmt.Sprintf("error invoking CNI ADD for args=%s", args), err)
			}
			return err
		},
		func(args *skel.CmdArgs) error {
			err := me.Check(args)
			if err != nil {
				logging.Error(fmt.Sprintf("error invoking CNI CHECK for args=%s", args), err)
			}
			return err
		},
		func(args *skel.CmdArgs) error {
			err := me.Del(args)
			if err != nil {
				logging.Error(fmt.Sprintf("error invoking CNI DEL for args=%s", args), err)
			}
			return err
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

func cniCommandFromSkelArgs(cmdStr string, args *skel.CmdArgs) util.CniCommand {
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
