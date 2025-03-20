package util

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/jboelensns/openstack-cni/pkg/logging"
	"github.com/joho/godotenv"
)

// CniContext represents all of the configuration required for the Cni Handler to function
// stdin data has been unmarshaled into CniContext
// args have been parsed and loaded into Args
type CniContext struct {
	Command   CniCommand
	Args      map[string]string
	CniConfig CniConfig
	Hostname  string
}

// NewCniContext creates a new CniContext instance out of a CniCommand
func NewCniContext(cmd CniCommand) (CniContext, error) {
	cniConfig, err := NewCniConfig(cmd.StdinData)
	if err != nil {
		logging.Log().Error().Str("stdindata", string(cmd.StdinData)).Err(err).Msg("failed to create cni context")
		return CniContext{}, err
	}

	hostname, _ := GetHostname()

	return CniContext{
		Command:   cmd,
		Args:      ParseCniArgs(cmd.Args),
		CniConfig: cniConfig,
		Hostname:  hostname,
	}, nil
}

// GetArg returns a value from the Args map
// asking for a non-existent key yields ""
func (me *CniContext) GetArg(name string) string {
	val, found := me.Args[name]
	if !found {
		return ""
	}
	return val
}

// CniCommand contains all of the data required for a CNI command
type CniCommand struct {
	Command     string `json:"command,omitempty"`
	ContainerID string `json:"container_id,omitempty"`
	Netns       string `json:"netns,omitempty"`
	IfName      string `json:"ifname,omitempty"`
	Args        string `json:"args,omitempty"`
	Path        string `json:"path,omitempty"`
	StdinData   []byte `json:"stdindata,omitempty"`
}

func (me CniCommand) String() string {
	return fmt.Sprintf("command=%q,container_id=%q,netns=%q,ifname=%q,args=%q,path=%q,stdindata=%q",
		me.Command, me.ContainerID, me.Netns, me.IfName, me.Args, me.Path, me.StdinData)
}

func (me CniCommand) ForLog() [][]string {
	return [][]string{
		{"cmd", me.Command},
		{"args", me.Args},
		{"container", me.ContainerID},
		{"ifname", me.IfName},
		{"netns", me.Netns},
		{"path", me.Path},
		{"stdindata", string(me.StdinData)},
	}
}

// CniConfig represents the config section of the NetworkAttachmentDefinition
/*
	For Example:
	spec:
    config: '{
      "type": "openstack-cni",
      "cniVersion": "0.3.1",
      "name": "service-ingress",
      "network": "compute-internal",
      "security_groups": ["dp_default", "default"]
      "enable_port_security": true
    }'
*/
type CniConfig struct {
	*types.NetConf
	AllowedAddressPairs []AddressPair `json:"allowed_address_pairs,omitempty"`
	AdminStateUp        *bool         `json:"admin_state_up,omitempty"`
	DeviceId            string        `json:"device_id,omitempty"`
	DeviceOwner         string        `json:"device_owner,omitempty"`
	// FixedIPs string `json:"fixed_ips,omitempty"`
	MacAddress      string             `json:"mac_address,omitempty"`
	Network         string             `json:"network,omitempty"`
	PortDescription string             `json:"port_description,omitempty"`
	PortName        string             `json:"port_name,omitempty"`
	ProjectName     string             `json:"project_name,omitempty"`
	SecurityGroups  *[]string          `json:"security_groups,omitempty"`
	SubnetName      string             `json:"subnet_name,omitempty"`
	TenantId        string             `json:"tenant_id,omitempty"`
	ValueSpecs      *map[string]string `json:"value_specs,omitempty"`
	// PortSecurityEnabled toggles port security on a port.
	PortSecurityEnabled *bool `json:"port_security_enabled,omitempty"`
	// The ID of the host where the port is allocated.
	HostID string `json:"binding:host_id,omitempty"`
	// The virtual network interface card (vNIC) type that is bound to the
	// neutron port.
	VNICType string `json:"binding:vnic_type,omitempty"`
	// A dictionary that enables the application running on the specified
	// host to pass and receive virtual network interface (VIF) port-specific
	// information to the plug-in.
	Profile map[string]interface{} `json:"binding:profile,omitempty"`
}

func NewCniConfig(bytes []byte) (CniConfig, error) {
	conf := &CniConfig{}
	if err := json.Unmarshal(bytes, conf); err != nil {
		return *conf, fmt.Errorf("Failed to load config data, error = %+v", err)
	}
	if conf.PortName == "" {
		conf.PortName = Getenv("OS_PORT_NAME", "openstack-cni")
	}
	if conf.ProjectName == "" {
		conf.ProjectName = Getenv("OS_PROJECT_NAME", "")
	}
	if conf.DeviceOwner == "" {
		conf.DeviceOwner = Getenv("CNI_PORT_DEVICE_OWNER", "")
	}
	// port security is enabled by default
	t := true
	if conf.PortSecurityEnabled == nil {
		conf.PortSecurityEnabled = &t
	}

	return *conf, nil
}

type AddressPair struct {
	IpAddress  string `json:"ip_address,omitempty"`
	MacAddress string `json:"mac_address,omitempty"`
}

// ParseCniArgs parses a key value pair such as "IgnoreUnknown=true;K8S_POD_NAMESPACE=lightning;"
// an empty map is returned if bad data is encountered
func ParseCniArgs(args string) map[string]string {
	mappy := make(map[string]string)
	pairs := strings.Split(args, ";")
	if len(pairs) == 0 {
		return mappy
	}
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts)%2 != 0 {
			return make(map[string]string)
		}
		for i := 0; i < len(parts); i += 2 {
			mappy[parts[i]] = parts[i+1]
		}
	}
	return mappy
}

func GetHostname() (vmName string, err error) {
	vmName = os.Getenv("OS_VM_NAME")
	if vmName == "" {
		vmName, err = os.Hostname()
	}

	return
}

func LoadEnvConfig(filenames ...string) error {
	for _, file := range filenames {
		exists, err := FileExists(file)
		if err != nil {
			return err
		}
		if exists {
			if err := godotenv.Load(file); err != nil {
				return err
			}
		}
	}
	return nil
}

func ReadConfigIntoEnv() error {
	return LoadEnvConfig(Getenv("CNI_CONFIG_FILE", "config.conf"))
}
