package fixtures

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/jboelensns/openstack-cni/pkg/cniplugin"
	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	"github.com/jboelensns/openstack-cni/pkg/fixtures/mocks"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"
	. "github.com/pepinns/go-hamcrest"
)

func WithServer(t *testing.T, callback func(fix *ServerFixture)) {
	WithServerOpts(t, nil, callback)
}

func WithServerOpts(t *testing.T, opts *ServerOpts, callback func(fix *ServerFixture)) {
	t.Helper()
	fix := NewServerFixture(t, opts)
	fix.Start(t)
	defer fix.Stop(t)
	callback(fix)
}

type ServerOpts struct {
	CniHandler      cniserver.CommandHandler
	OpenstackClient openstack.OpenstackClient
	Networking      cniplugin.Networking
}

type TestingConfig struct {
	EnableOpenstackTests string
	Hostname             string
	NetworkName          string
	PortName             string
	ProjectName          string
	SecurityGroups       *[]string
	SubnetName           string
}

func WithTestConfig(t *testing.T, callback func(cfg TestingConfig)) {
	Assert(t).That(util.LoadEnvConfig("testing.conf", "../../testing.conf"), IsNil())

	cfg := TestingConfig{
		EnableOpenstackTests: Getenv("OS_TESTS", "0"),
		Hostname:             Getenv("OS_VM_NAME", ""),
		NetworkName:          Getenv("OS_NETWORK_NAME", ""),
		PortName:             Getenv("OS_PORT_NAME", "openstack-cni-test"),
		ProjectName:          Getenv("OS_PROJECT_NAME", ""),
		SubnetName:           Getenv("OS_SUBNET_NAME", ""),
	}
	sgs := strings.Split(Getenv("OS_SECURITY_GROUPS", ""), ";")
	if len(sgs) > 0 {
		cfg.SecurityGroups = &sgs
	}

	if cfg.EnableOpenstackTests == "1" {
		Assert(t).That(cfg.NetworkName, Not(Equals("")), "missing OS_NETWORK_NAME")
		Assert(t).That(cfg.PortName, Not(Equals("")), "missing OS_PORT_NAME")
		Assert(t).That(cfg.PortName, Not(Equals("")), "missing OS_PROJECT_NAME")
		Assert(t).That(cfg.PortName, Not(Equals("")), "missing OS_PORT_NAME")
		callback(cfg)
	}
}

func WithMockClient(t *testing.T, callback func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient)) {
	WithMockClientWithExpiry(t, time.Second*5, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
		callback(mock, client)
	})
}

func WithMockClientWithExpiry(t *testing.T, expiry time.Duration, callback func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient)) {
	t.Helper()
	mock := &mocks.OpenstackClientMock{}
	client := openstack.NewCachedClient(mock, expiry)
	callback(mock, client)
}

func WithTempDir(t *testing.T, callback func(dir string)) {
	t.Helper()
	baseDir, err := os.MkdirTemp(os.TempDir(), "openstack-cni")
	Assert(t).That(err, IsNil())

	defer func() {
		Assert(t).That(os.RemoveAll(baseDir), IsNil())
	}()
	callback(baseDir)
}

func WithOpenstackClient(t *testing.T, callback func(client openstack.OpenstackClient)) {
	t.Helper()
	client, err := openstack.NewOpenstackClient()
	Assert(t).That(err, IsNil())
	callback(client)
}

func WithPortReaper(t *testing.T, client openstack.OpenstackClient, callback func(reaper *cniserver.PortReaper)) {
	t.Helper()
	reaper := &cniserver.PortReaper{
		Opts:     PortReaperOpts(),
		OsClient: client,
		Metrics:  Metrics(),
	}
	callback(reaper)
}

func WithPortReaperWithNoMinPortAge(t *testing.T, client openstack.OpenstackClient, callback func(reaper *cniserver.PortReaper)) {
	t.Helper()
	reaper := &cniserver.PortReaper{
		Opts:     PortReaperOpts(),
		OsClient: client,
		Metrics:  Metrics(),
	}
	reaper.Opts.MinPortAge = 0
	callback(reaper)
}

func WithMountedProcDir(t *testing.T, reaper *cniserver.PortReaper, callback func()) {
	t.Helper()
	WithTempDir(t, func(dir string) {
		reaper.Opts.MountedProcDir = path.Join(dir, "proc")
		Assert(t).That(os.MkdirAll(reaper.Opts.MountedProcDir, 0755), IsNil())
		callback()
	})
}

func Getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
