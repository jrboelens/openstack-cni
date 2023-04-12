package cniserver

import (
	"time"

	"github.com/jboelensns/openstack-cni/pkg/openstack"
)

// Deps represents dependencies for the application
// instances of this structure are created by the Builder
type Deps struct {
	cniHandler CommandHandler
	osClient   openstack.OpenstackClient
}

// CniHandler returns the CommandHandler
func (me *Deps) CniHandler() CommandHandler {
	return me.cniHandler
}

// OpenstackClient returns the OpenstackClient
func (me *Deps) OpenstackClient() openstack.OpenstackClient {
	return me.osClient
}

// Builder provides the ability to produce Deps instances using the builder pattern
type Builder struct {
	cniHandler CommandHandler
	osClient   openstack.OpenstackClient
}

// NewBuilder creates a new Builder
func NewBuilder() *Builder {
	return &Builder{}
}

// WithCniHandler sets the current to CommandHandler to cniHandler
func (me *Builder) WithCniHandler(cniHandler CommandHandler) *Builder {
	me.cniHandler = cniHandler
	return me
}

// WithOpenstackClient sets the current to OpenstackClient to client
func (me *Builder) WithOpenstackClient(client openstack.OpenstackClient) *Builder {
	me.osClient = client
	return me
}

// Build creates a Dep
func (me *Builder) Build() (*Deps, error) {
	// build the default os factory if we don't have one
	if me.osClient == nil {
		var err error
		me.osClient, err = openstack.NewOpenstackClient()
		if err != nil {
			return nil, err
		}
	}
	me.osClient = openstack.NewCachedClient(me.osClient, time.Second*300)

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
