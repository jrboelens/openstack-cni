package cniserver

import (
	"github.com/jboelensns/openstack-cni/pkg/cnistate"
	"github.com/jboelensns/openstack-cni/pkg/openstack"
)

type Deps struct {
	cniHandler CommandHandler
	osClient   openstack.OpenstackClient
	state      cnistate.State
}

func (me *Deps) CniHandler() CommandHandler {
	return me.cniHandler
}

func (me *Deps) OpenstackClient() openstack.OpenstackClient {
	return me.osClient
}

func (me *Deps) State() cnistate.State {
	return me.state
}

type Builder struct {
	cniHandler CommandHandler
	osClient   openstack.OpenstackClient
	state      cnistate.State
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (me *Builder) WithCniHandler(cniHandler CommandHandler) *Builder {
	me.cniHandler = cniHandler
	return me
}

func (me *Builder) WithOpenstackClient(client openstack.OpenstackClient) *Builder {
	me.osClient = client
	return me
}

func (me *Builder) WithState(state cnistate.State) *Builder {
	me.state = state
	return me
}

func (me *Builder) Build() (*Deps, error) {
	// build the default os factory if we don't have one
	if me.osClient == nil {
		var err error
		me.osClient, err = openstack.NewOpenstackClient()
		if err != nil {
			return nil, err
		}
	}

	// build the default cni handler if we don't have one
	if me.cniHandler == nil {
		var err error
		pm := openstack.NewPortManager(me.osClient)
		state := cnistate.NewState(cnistate.GetStateBaseDir())
		me.cniHandler, err = NewCniCommandHandler(pm, state), nil
		if err != nil {
			return nil, err
		}
	}

	if me.state == nil {
		me.state = cnistate.NewState(cnistate.GetStateBaseDir())
	}

	return &Deps{
		cniHandler: me.cniHandler,
		osClient:   me.osClient,
		state:      me.state,
	}, nil
}
