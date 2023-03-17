package fixtures

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/jboelensns/openstack-cni/pkg/cniclient"
	"github.com/jboelensns/openstack-cni/pkg/cniplugin"
	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	"github.com/jboelensns/openstack-cni/pkg/cnistate"
	"github.com/jboelensns/openstack-cni/pkg/fixtures/mocks"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
	"github.com/jboelensns/openstack-cni/pkg/util"

	. "github.com/pepinns/go-hamcrest"
)

type ServerFixture struct {
	BaseUrl    string
	t          *testing.T
	app        *cniserver.App
	cniHandler cniserver.CommandHandler
	openstack  openstack.OpenstackClient
	state      cnistate.State
	networking cniplugin.Networking
	cfg        cniserver.Config
}

func NewServerFixture(t *testing.T, opts *ServerOpts) *ServerFixture {
	var cniHandler cniserver.CommandHandler = &mocks.CommandHandlerMock{}
	var osClient openstack.OpenstackClient = &mocks.OpenstackClientMock{}
	var state cnistate.State = &mocks.StateMock{}
	var networking cniplugin.Networking = &mocks.NetworkingMock{}
	if opts != nil {
		if opts.CniHandler != nil {
			cniHandler = opts.CniHandler
		}
		if opts.OpenstackClient != nil {
			osClient = opts.OpenstackClient
		}
		if opts.State != nil {
			state = opts.State
		}
		if opts.Networking != nil {
			networking = opts.Networking
		}
	}

	// read in configs into our environment
	// these configs drive the openstack clients
	err := util.ReadConfigIntoEnv()
	Assert(t).That(err, IsNil())

	return &ServerFixture{
		BaseUrl:    "http://0.0.0.0",
		t:          t,
		cniHandler: cniHandler,
		openstack:  osClient,
		state:      state,
		networking: networking,
	}
}

func (me *ServerFixture) Config() cniserver.Config {
	return me.cfg
}

func (me *ServerFixture) Client() *Client {
	return &Client{}
}

func (me *ServerFixture) CniClient() *cniclient.Client {
	return &cniclient.Client{
		Opts: cniclient.ClientOpts{
			BaseUrl:        me.Url(""),
			RequestTimeout: time.Second * 5,
		},
	}
}

func (me *ServerFixture) Ping(t *testing.T) {
	t.Helper()
	resp, err := me.Client().Get(me.Url("/ping"), nil)
	Assert(t).That(err, IsNil())
	Assert(t).That(resp.StatusCode, Equals(200))
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	Assert(t).That(err, IsNil())
	Assert(t).That(b, Equals("PONG"))
}

func (me *ServerFixture) Openstack() openstack.OpenstackClient {
	return me.openstack
}

func (me *ServerFixture) Start(t *testing.T) {
	t.Helper()

	deps, err := cniserver.NewBuilder().
		WithCniHandler(me.cniHandler).
		WithOpenstackClient(me.openstack).
		WithState(me.state).
		Build()
	Assert(t).That(err, IsNil())

	me.cfg = cniserver.NewConfig()
	me.cfg.ListenAddr = me.GetListenAddr(me.GetPort())
	app, err := cniserver.NewApp(me.cfg, cniserver.SetupRoutes(deps))
	Assert(t).That(err, IsNil())
	me.app = app
	go func() {
		if err := app.Run(); err != nil {
			panic(fmt.Sprintf("failed to start app in fixture %s", err))
		}
	}()
	me.waitForStart()
}

func (me *ServerFixture) Stop(t *testing.T) {
	t.Helper()
	if me.app != nil {
		me.app.Shutdown(context.Background())
	}
}

func (me *ServerFixture) TestData() *TestData {
	return &TestData{}
}

func (me *ServerFixture) Assert(t *testing.T) *Assertions {
	return &Assertions{t}
}

func (me *ServerFixture) Url(path string) string {
	_, port, _ := net.SplitHostPort(me.Config().ListenAddr)
	return fmt.Sprintf("%s:%s%s", me.BaseUrl, port, path)
}

func (me *ServerFixture) waitForStart() {
	for d := time.Now().Add(5 * time.Second); time.Now().Before(d); {
		resp, err := me.Client().Get(me.Url("/ping"), nil)
		if err == nil && resp.StatusCode == 200 {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func (me *ServerFixture) GetListenAddr(port int) string {
	return fmt.Sprintf(":%d", me.GetPort())
}

func (me *ServerFixture) GetPort() int {
	return me.GetAvailablePort(":0")
}

func (me *ServerFixture) GetAvailablePort(addrs ...string) int {
	var lastErr error
	for _, addr := range addrs {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			lastErr = fmt.Errorf("failed to acquire listen address %s", addr)
			continue
		}
		defer listener.Close()
		return listener.Addr().(*net.TCPAddr).Port
	}
	if lastErr != nil {
		panic(lastErr)
	}
	return 0
}

func CniContextFromConfig(t *testing.T, cfg TestingConfig, cmd util.CniCommand) util.CniContext {
	context, err := util.NewCniContext(cmd)
	Assert(t).That(err, IsNil())
	context.Hostname = cfg.Hostname
	context.CniConfig.Network = cfg.NetworkName
	context.CniConfig.PortName = cfg.PortName
	context.CniConfig.ProjectName = cfg.ProjectName
	context.CniConfig.SecurityGroups = cfg.SecurityGroups
	return context
}
