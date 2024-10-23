package cniserver_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	"github.com/jboelensns/openstack-cni/pkg/fixtures/mocks"
	"github.com/jboelensns/openstack-cni/pkg/openstack"

	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	. "github.com/pepinns/go-hamcrest"
)

func Test_PortReaper(t *testing.T) {
	hostname, err := os.Hostname()
	Assert(t).That(err, IsNil())
	t.Run("port reaper attempts to delete ports", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			WithPortReaper(t, client, func(reaper *cniserver.PortReaper) {
				serverId := "myId"
				mock.GetServerByNameFunc = func(name string) (*servers.Server, error) {
					return &servers.Server{ID: serverId}, nil
				}
				mock.DeletePortFunc = func(portId string) error { return nil }
				mock.GetPortsByTagsFunc = func(tags []string) ([]ports.Port, error) {
					return []ports.Port{{Status: "DOWN", Tags: NeutronTags()}}, nil
				}

				WithTempDir(t, func(dir string) {
					reaper.Opts.MountedProcDir = path.Join(dir, "proc")
					os.MkdirAll(reaper.Opts.MountedProcDir, 0755)
					err = reaper.Reap(hostname)
					Assert(t).That(err, IsNil())
					Assert(t).That(len(mock.DeletePortCalls()), Equals(1))
				})
			})
		})
	})

	t.Run("will not reap a new port", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			WithPortReaper(t, client, func(reaper *cniserver.PortReaper) {
				port := ports.Port{CreatedAt: time.Now()}
				mock.DeletePortFunc = func(portId string) error { return nil }
				err = reaper.ReapPort(port)
				Assert(t).That(err, IsNil())
				Assert(t).That(len(mock.DeletePortCalls()), Equals(0))
			})
		})
	})

	t.Run("will not reap a port without our tags", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			WithPortReaper(t, client, func(reaper *cniserver.PortReaper) {
				port := ports.Port{CreatedAt: time.Now()}
				mock.DeletePortFunc = func(portId string) error { return nil }
				err = reaper.ReapPort(port)
				Assert(t).That(err, IsNil())
				Assert(t).That(len(mock.DeletePortCalls()), Equals(0))
			})
		})
	})

	t.Run("will not reap a port without /proc mounted at /host/proc", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			WithPortReaper(t, client, func(reaper *cniserver.PortReaper) {
				mock.GetServerByNameFunc = func(name string) (*servers.Server, error) {
					return nil, nil
				}
				WithTempDir(t, func(dir string) {
					// this directory won't exist
					reaper.Opts.MountedProcDir = path.Join(dir, "proc")
					err = reaper.Reap(hostname)
					Assert(t).That(err, IsNil())
					Assert(t).That(len(mock.DeletePortCalls()), Equals(0))
				})
			})
		})
	})

	t.Run("will reap an old port", func(t *testing.T) {
		WithMockClient(t, func(mock *mocks.OpenstackClientMock, client openstack.OpenstackClient) {
			WithPortReaper(t, client, func(reaper *cniserver.PortReaper) {
				port := ports.Port{Status: "DOWN", Tags: NeutronTags(), CreatedAt: time.Now().Add(-(time.Second * 6000))}
				mock.DeletePortFunc = func(portId string) error { return nil }
				err = reaper.ReapPort(port)
				Assert(t).That(err, IsNil())
				Assert(t).That(len(mock.DeletePortCalls()), Equals(1))
			})
		})
	})
}

func Test_PortReaperIntegration(t *testing.T) {
	t.Run("port reaper attempts to delete ports", func(t *testing.T) {
		WithTestConfig(t, func(cfg TestingConfig) {
			WithOpenstackClient(t, func(client openstack.OpenstackClient) {
				// create a port with a network namespace that doesn't exist for my machine
				cmd := NewTestData().CniCommand()
				context := CniContextFromConfig(t, cfg, cmd)
				cachedClient := openstack.NewCachedClient(client, time.Second*5)

				WithPortReaperWithNoMinPortAge(t, cachedClient, func(reaper *cniserver.PortReaper) {
					WithMountedProcDir(t, reaper, func() {
						pm := openstack.NewPortManager(cachedClient)
						opts := openstack.SetupPortOptsFromContext(context)
						opts.Tags = cniserver.NewPortTags(context.Command)
						opts.SkipPortAttach = true

						_, err := pm.SetupPort(opts)
						Assert(t).That(err, IsNil(), "failed to setup port")

						_, err = cachedClient.GetPortByTags(opts.Tags.AsStringSlice())
						Assert(t).That(err, IsNil(), "failed get port by tags %s", opts.Tags.String())

						// run the reaper
						Assert(t).That(reaper.Reap(cfg.Hostname), IsNil())

						// ensure the port is deleted
						_, err = client.GetPortByTags(opts.Tags.AsStringSlice())
						Assert(t).That(err, Equals(openstack.ErrPortNotFound))
					})
				})
			})
		})
	})

	t.Run("port reaper doesn't delete ports when SkipReaper=true", func(t *testing.T) {
		WithTestConfig(t, func(cfg TestingConfig) {
			WithOpenstackClient(t, func(client openstack.OpenstackClient) {
				// create a port with a network namespace that doesn't exist for my machine
				cmd := NewTestData().CniCommand()
				context := CniContextFromConfig(t, cfg, cmd)
				cachedClient := openstack.NewCachedClient(client, time.Second*5)

				WithPortReaperWithNoMinPortAge(t, cachedClient, func(reaper *cniserver.PortReaper) {
					WithMountedProcDir(t, reaper, func() {
						pm := openstack.NewPortManager(cachedClient)
						opts := openstack.SetupPortOptsFromContext(context)
						opts.Tags = cniserver.NewPortTags(context.Command)
						reaper.Opts.SkipDelete = true

						setupResult, err := pm.SetupPort(opts)
						Assert(t).That(err, IsNil(), "failed to setup port")

						p1, err := cachedClient.GetPortByTags(opts.Tags.AsStringSlice())
						Assert(t).That(err, IsNil(), "failed get port by tags %s", opts.Tags.String())

						// run the reaper
						Assert(t).That(reaper.Reap(cfg.Hostname), IsNil())

						// ensure the port is still present
						p2, err := client.GetPortByTags(opts.Tags.AsStringSlice())
						Assert(t).That(err, IsNil())
						Assert(t).That(p1.ID, Equals(p2.ID))

						// detach the port so it will be reaped
						cachedClient.DetachPort(setupResult.Port.ID, setupResult.Server.ID)

						// now reap it
						reaper.Opts.SkipDelete = false
						Assert(t).That(reaper.Reap(cfg.Hostname), IsNil())

						// ensure it's gone
						p, err := client.GetPortByTags(opts.Tags.AsStringSlice())
						Assert(t).That(p, Equals(nil))
						Assert(t).That(err, Equals(openstack.ErrPortNotFound))
					})
				})
			})
		})
	})
}
