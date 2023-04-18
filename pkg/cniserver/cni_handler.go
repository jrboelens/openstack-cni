package cniserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/containernetworking/cni/pkg/types"
	. "github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/jboelensns/openstack-cni/pkg/util"
)

var CommandAdd = "ADD"
var CommandDel = "DEL"
var CommandCheck = "CHECK"

// CniHandler handles /cni related requests
type CniHandler struct {
	Cni     CommandHandler
	Metrics *Metrics
}

// HandleRequest validates and handlers /cni related requests
func (me *CniHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	me.Metrics.cniRequestCount.Inc()
	cmd := readBodyIntoJson[util.CniCommand](w, r)
	if cmd == nil {
		me.Metrics.cniRequestInvalidCount.Inc()
		return
	}

	if err := me.validateCommand(*cmd); err != nil {
		Log().Error().Str("cmd", cmd.String()).AnErr("err", err).Msg("failed to validate request")
		w.WriteHeader(http.StatusBadRequest)
		me.Metrics.cniRequestInvalidCount.Inc()
		return
	}

	me.HandleCommand(w, *cmd)
}

// HandleCommand handlers ADD/DEL/CHECK CNI command requests
func (me *CniHandler) HandleCommand(w http.ResponseWriter, cmd util.CniCommand) {
	switch cmd.Command {
	case CommandAdd:
		result, err := me.Cni.Add(cmd)
		if err != nil {
			me.Metrics.cniAddFailureCount.Inc()
			AddStrings(Log().Error(), cmd.ForLog()).Err(err).Msg("failed to handle /cni ADD")
			cerr := NewErrorResult(err, "error during ADD", fmt.Sprintf("command=%s", cmd))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(cerr)
			return
		} else {
			me.Metrics.cniAddSuccessCount.Inc()
			w.Write(asJson(result))
		}
		return
	case CommandDel:
		if err := me.Cni.Del(cmd); err != nil {
			me.Metrics.cniDelFailureCount.Inc()
			AddStrings(Log().Error(), cmd.ForLog()).Err(err).Msg("failed to handle /cni DEL")
			cerr := NewErrorResult(err, "error during DEL", fmt.Sprintf("containerid=%s ifname=%s", cmd.ContainerID, cmd.IfName))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(cerr)
			return
		}
		me.Metrics.cniDelSuccessCount.Inc()
		w.WriteHeader(http.StatusNoContent)
		return
	case CommandCheck:
		if err := me.Cni.Check(cmd); err != nil {
			me.Metrics.cniCheckFailureCount.Inc()
			AddStrings(Log().Error(), cmd.ForLog()).Err(err).Msg("failed to handle /cni CHECK")
			cerr := NewErrorResult(err, "error during CHECK", fmt.Sprintf("containerid=%s ifname=%s", cmd.ContainerID, cmd.IfName))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(cerr)
			return
		}
		me.Metrics.cniCheckSuccessCount.Inc()
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// validateCommand ensures that a command is valid
func (me *CniHandler) validateCommand(cmd util.CniCommand) error {
	if cmd.Command == "" ||
		cmd.ContainerID == "" ||
		cmd.IfName == "" ||
		len(cmd.StdinData) == 0 {
		return ErrBadCommand
	}

	if cmd.Command == CommandAdd || cmd.Command == CommandCheck {
		if cmd.Netns == "" {
			return ErrBadCommand
		}
	}
	return nil
}

var ErrBadCommand = fmt.Errorf("bad command")

// NewErrorResult creates a new error result and marshals it as json
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

func readBodyIntoJson[T any](w http.ResponseWriter, r *http.Request) *T {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		Log().Error().Str("body", string(body)).AnErr("err", err).Msg("failed to read body")
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	var t T
	if err := json.Unmarshal(body, &t); err != nil {
		Log().Error().Str("body", string(body)).AnErr("err", err).Msg("failed to unmarshal body")
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	return &t
}
