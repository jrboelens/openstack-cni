package cniserver

import (
	"github.com/jboelensns/openstack-cni/pkg/openstack"
)

type Deps struct {
	cniHandler CommandHandler
	osClient   openstack.OpenstackClient
}

func (me *Deps) CniHandler() CommandHandler {
	return me.cniHandler
}

func (me *Deps) OpenstackClient() openstack.OpenstackClient {
	return me.osClient
}

type Builder struct {
	cniHandler CommandHandler
	osClient   openstack.OpenstackClient
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
		me.cniHandler, err = NewCniCommandHandler(pm), nil
		if err != nil {
			return nil, err
		}
	}

	return &Deps{
		cniHandler: me.cniHandler,
		osClient:   me.osClient,
	}, nil
}
