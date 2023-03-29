package openstack_test

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	. "github.com/jboelensns/openstack-cni/pkg/fixtures"
	"github.com/jboelensns/openstack-cni/pkg/openstack"

	. "github.com/pepinns/go-hamcrest"
)

func Test_Tagger(t *testing.T) {
	WithTestConfig(t, func(cfg TestingConfig) {
		WithOpenstackClient(t, func(client openstack.OpenstackClient) {
			tags := openstack.NeutronTags{
				Tags: []openstack.NeutronTag{"foo", "bar", "zilla"},
			}
			tagger := openstack.NewNeutronTagger(client.Clients().NetworkClient, openstack.Ports)
			t.Run("can execute all tagging operations on a port", func(t *testing.T) {
				WithPort(t, cfg, client, func(port *ports.Port) {
					// delete all of the tags
					Assert(t).That(tagger.DeleteAll(port.ID), IsNil())

					// make sure they are gone
					existingTags, err := tagger.GetAll(port.ID)
					Assert(t).That(err, IsNil())
					Assert(t).That(existingTags.Tags, HasLen(0))

					// add some tags
					Assert(t).That(tagger.SetAll(port.ID, tags), IsNil())

					// make sure the were added
					existingTags, err = tagger.GetAll(port.ID)
					Assert(t).That(err, IsNil())
					Assert(t).That(existingTags.Tags, AllOf(
						HasLen(3),
						Contains("foo"),
						Contains("bar"),
						Contains("zilla"),
					))

					// delete the tags again
					Assert(t).That(tagger.DeleteAll(port.ID), IsNil())

					// add a single tag
					tag := "zilla"
					Assert(t).That(tagger.Create(port.ID, tag), IsNil())

					// check for the existence of a tag
					exists, err := tagger.Exists(port.ID, tag)
					Assert(t).That(err, IsNil())
					Assert(t).That(exists, IsTrue())

					// delete the tag
					Assert(t).That(tagger.Delete(port.ID, tag), IsNil())

					// make sure the tag is gone
					exists, err = tagger.Exists(port.ID, tag)
					Assert(t).That(err, IsNil())
					Assert(t).That(exists, IsFalse())
				})
			})
		})
	})
}

func WithPort(t *testing.T, cfg TestingConfig, client openstack.OpenstackClient, callback func(*ports.Port)) {
	network, err := client.GetNetworkByName(cfg.NetworkName)
	Assert(t).That(err, IsNil())

	project, err := client.GetProjectByName(cfg.ProjectName)
	Assert(t).That(err, IsNil())

	port, err := client.CreatePort(ports.CreateOpts{
		NetworkID:  network.ID,
		Name:       "openstack-cni-unit-test",
		ProjectID:  project.ID,
		ValueSpecs: nil,
	})
	Assert(t).That(err, IsNil())

	defer func() {
		Assert(t).That(client.DeletePort(port.ID), IsNil())
	}()
	callback(port)
}
