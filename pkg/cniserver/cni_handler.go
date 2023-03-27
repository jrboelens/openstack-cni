package cniserver

import (
	"fmt"
	"net/http"

	"github.com/containernetworking/cni/pkg/types"
	. "github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

var CommandAdd = "ADD"
var CommandDel = "DEL"
var CommandCheck = "CHECK"

type CniHandler struct {
	Cni CommandHandler
}

func (me *CniHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	cmd := readBodyIntoJson[util.CniCommand](w, r)
	if cmd == nil {
		return
	}

	if err := me.validateRequest(*cmd); err != nil {
		Log().Error().Str("cmd", cmd.String()).AnErr("err", err).Msg("failed to validate request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	me.HandleCommand(w, *cmd)
}

func (me *CniHandler) HandleCommand(w http.ResponseWriter, cmd util.CniCommand) {
	switch cmd.Command {
	case CommandAdd:
		result, err := me.Cni.Add(cmd)
		if err != nil {
			AddStrings(Log().Error(), cmd.ForLog()).Err(err).Msg("failed to handle /cni ADD")
			cerr := NewErrorResult(err, "error during ADD", fmt.Sprintf("command=%s", cmd))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(cerr)
			return
		} else {
			w.Write(asJson(result))
		}
		return
	case CommandDel:
		if err := me.Cni.Del(cmd); err != nil {
			AddStrings(Log().Error(), cmd.ForLog()).Err(err).Msg("failed to handle /cni DEL")
			cerr := NewErrorResult(err, "error during DEL", fmt.Sprintf("containerid=%s ifname=%s", cmd.ContainerID, cmd.IfName))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(cerr)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	case CommandCheck:
		me.Cni.Check(cmd)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (me *CniHandler) validateRequest(cmd util.CniCommand) error {
	if cmd.Command == "" ||
		cmd.ContainerID == "" ||
		cmd.IfName == "" ||
		len(cmd.StdinData) == 0 {
		return ErrBadRequest
	}

	if cmd.Command == CommandAdd || cmd.Command == CommandCheck {
		if cmd.Netns == "" {
			return ErrBadRequest
		}
	}
	return nil
}

var ErrBadRequest = fmt.Errorf("bad request")

func NewErrorResult(err error, msg, details string) []byte {
	return asJson(types.NewError(types.ErrInternal, msg, err.Error()))
}

func asJson(i any) []byte {
	b, err := util.ToJson(i)
	if err != nil {
		panic(fmt.Sprintf("failed to json marshal bytes %#v", i))
	}
	return b
}
