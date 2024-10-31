package cniplugin

import (
	"fmt"
	"time"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	currentcni "github.com/containernetworking/cni/pkg/types/040"
	cniversion "github.com/containernetworking/cni/pkg/version"
	"github.com/jboelensns/openstack-cni/pkg/cniclient"
	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	"github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
	"github.com/rs/zerolog"
)

type CniOpts struct {
	WaitForUdev       bool
	WaitForUdevPrefix string
	WaitForUdevDelay  time.Duration
	WaitForUdevTimeout  time.Duration
}

func DefaultCniOpts() CniOpts {
	return CniOpts{
		WaitForUdev:       true,
		WaitForUdevPrefix: "eth",
		WaitForUdevDelay:  100 * time.Millisecond,
		WaitForUdevTimeout:  5000 * time.Millisecond,
	}
}

// Cni provides methods with the ability accept CNI spec data, make requests to the openstack-cni-daemon and return the results
type Cni struct {
	Opts   CniOpts
	client *cniclient.Client
	nw     Networking
}

// NewCni returns a new Cni
func NewCni(client *cniclient.Client, nw Networking, opts CniOpts) *Cni {
	return &Cni{
		Opts:   opts,
		client: client,
		nw:     nw,
	}
}

// Add handles ADD CNI commands
func (me *Cni) Add(args *skel.CmdArgs) error {
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
	cmd := cniCommandFromSkelArgs(cniserver.CommandCheck, args)
	_, err := me.client.HandleResponse(me.client.CniCommand(cmd))
	return err
}

// Del handles DEL CNI commands
func (me *Cni) Del(args *skel.CmdArgs) error {
	cmd := cniCommandFromSkelArgs(cniserver.CommandDel, args)
	_, err := me.client.HandleResponse(me.client.CniCommand(cmd))
	return err
}

func argLogContext(l zerolog.Context, args *skel.CmdArgs) zerolog.Logger {
	return l.Str("container_id", args.ContainerID).Str("ns", args.Netns).Str("iface", args.IfName).Str("args", args.Args).Str("path", args.Path).Logger()
}

// Invoke invokes the CNI plugin skeletons using its own methods
func (me *Cni) Invoke() error {
	err := skel.PluginMainWithError(
		func(args *skel.CmdArgs) error {
			log := argLogContext(logging.Log().With(), args)
			log.Info().Msg("received ADD")
			err := me.Add(args)
			if err != nil {
				logging.Error(fmt.Sprintf("error invoking CNI ADD for args=%s", args), err)
			} else {
				log.Info().Msg("successful ADD")
			}

			return err
		},
		func(args *skel.CmdArgs) error {
			log := argLogContext(logging.Log().With(), args)
			log.Info().Msg("received CHECK")
			err := me.Check(args)
			if err != nil {
				logging.Error(fmt.Sprintf("error invoking CNI CHECK for args=%s", args), err)
			} else {
				log.Info().Msg("successful CHECK")
			}
			return err
		},
		func(args *skel.CmdArgs) error {
			log := argLogContext(logging.Log().With(), args)
			log.Info().Msg("received DEL")
			err := me.Del(args)
			if err != nil {
				logging.Error(fmt.Sprintf("error invoking CNI DEL for args=%s", args), err)
			} else {
				log.Info().Msg("successful DEL")
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
