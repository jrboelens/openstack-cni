package cniserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jboelensns/openstack-cni/pkg/cnistate"
	. "github.com/jboelensns/openstack-cni/pkg/logging"
)

type StateHandler struct {
	state cnistate.State
}

func NewStateHandler(state cnistate.State) *StateHandler {
	return &StateHandler{state}
}

func (me *StateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerId")
	ifname := chi.URLParam(r, "ifname")
	info, err := me.state.Get(containerId, ifname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := me.state.Delete(containerId, ifname); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (me *StateHandler) Get(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerId")
	ifname := chi.URLParam(r, "ifname")
	info, err := me.state.Get(containerId, ifname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if info == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(asJson(info))
}

func (me *StateHandler) Set(w http.ResponseWriter, r *http.Request) {
	info := readBodyIntoJson[cnistate.IfaceInfo](w, r)
	if info == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := me.state.Set(info); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
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
