package cnistate

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/jboelensns/openstack-cni/pkg/util"
)

//go:generate moq -pkg mocks -out ../fixtures/mocks/cnistate_mocks.go . State

type State interface {
	Delete(containerId, ifname string) error
	Get(containerId, ifname string) (*IfaceInfo, error)
	Set(*IfaceInfo) error
}

func NewState(baseDir string) *state {
	return &state{baseDir}
}

type state struct {
	baseDir string
}

func (me *state) Delete(containerId, ifname string) error {
	filename := me.Filename(containerId, ifname)
	if util.FileExists(filename) {
		return os.Remove(filename)
	}
	return nil
}

// Get returns interface info based on a container Id and interface name
// if the state is not found the resulting *IfaceInfo is nil
func (me *state) Get(containerId, ifname string) (*IfaceInfo, error) {
	filename := me.Filename(containerId, ifname)
	if !util.FileExists(filename) {
		return nil, nil
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var info IfaceInfo
	if err := json.Unmarshal(b, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (me *state) Set(info *IfaceInfo) error {
	b, err := util.ToJson(info)
	if err != nil {
		return err
	}

	return me.SetRaw(me.Filename(info.ContainerId, info.Ifname), b)
}

func (me *state) Filename(containerId, ifname string) string {
	return path.Join(me.baseDir, fmt.Sprintf("%s-%s.json", containerId, ifname))
}

func (me *state) SetRaw(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

type IfaceInfo struct {
	ContainerId string `json:"container_id,omitempty"`
	Ifname      string `json:"ifname,omitempty"`
	Netns       string `json:"netns,omitempty"`
	IpAddress   string `json:"ip_address,omitempty"`
	PodName     string `json:"pod_name,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	CniResult   any    `json:"CniResult,omitempty"`
}

func GetStateBaseDir() string {
	return util.Getenv("CNI_STATE_DIR", "/host/etc/cni/net.d/openstack-cni-state")
}
