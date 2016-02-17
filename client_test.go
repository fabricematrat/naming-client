// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package namingclient_test

import (
	"testing"

	"github.com/coreos/etcd/integration"
	"github.com/CanonicalLtd/naming-client"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) { gc.TestingT(t) }

type suite struct {
	addr []string
}

var _ = gc.Suite(&suite{})

func (s *suite) SetUpTest(c *gc.C) {
	// Start etcd for testing use.
	cluster := integration.NewCluster(nil, 1)
	cluster.Launch(nil)
	s.addr = []string{cluster.URL(0)}
}

func (s *suite) TestCreateAndReadModel(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	value, err := client.Read("foo")
	c.Assert(err, gc.IsNil)
	c.Assert(*value, gc.Equals, namingclient.Model)
}

func (s *suite) TestCreateAndReadCharm(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Charm)
	c.Assert(err, gc.IsNil)
	value, err := client.Read("foo")
	c.Assert(err, gc.IsNil)
	c.Assert(*value, gc.Equals, namingclient.Charm)
}

func (s *suite) TestCreateAndReadPage(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Page)
	c.Assert(err, gc.IsNil)
	value, err := client.Read("foo")
	c.Assert(err, gc.IsNil)
	c.Assert(*value, gc.Equals, namingclient.Page)
}

func (s *suite) TestDelete(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	err = client.Delete("foo")
	c.Assert(err, gc.IsNil)
	value, err := client.Read("foo")
	c.Assert(err, gc.ErrorMatches, "key 'foo' does not exist")
	c.Assert(value, gc.IsNil)
}

func (s *suite) TestDeleteUnknownKey(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Delete("foo")
	c.Assert(err, gc.ErrorMatches, "key 'foo' does not exist")
}

func (s *suite) TestReadUnknownKey(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	value, err := client.Read("nope")
	c.Assert(err, gc.ErrorMatches, "key 'nope' does not exist")
	c.Assert(value, gc.IsNil)
}

func (s *suite) TestCreateExistingFails(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.ErrorMatches, "key 'foo' exists")
}

func (s *suite) TestReadDirectoryShouldFail(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("mydir/foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	value, err := client.Read("mydir")
	c.Assert(err, gc.ErrorMatches, "cannot read directory 'mydir'")
	c.Assert(value, gc.IsNil)
}

func (s *suite) TestUpdate(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	err = client.Update("foo", namingclient.Model, namingclient.Charm)
	value, err := client.Read("foo")
	c.Assert(err, gc.IsNil)
	c.Assert(*value, gc.Equals, namingclient.Charm)
}

func (s *suite) TestUpdateNoKey(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Update("foo", namingclient.Model, namingclient.Charm)
	c.Assert(err, gc.ErrorMatches, "key 'foo' does not exist")
}

func (s *suite) TestCreateAndListDirectory(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo/bar", namingclient.Model)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo/baz", namingclient.Model)
	c.Assert(err, gc.IsNil)
	value, err := client.List("foo")
	c.Assert(err, gc.IsNil)
	c.Assert(value, jc.DeepEquals, map[string]namingclient.ServiceType{
		"/foo/bar": namingclient.ServiceType("model"),
		"/foo/baz": namingclient.ServiceType("model"),
	})
}

func (s *suite) TestListNonDirectoryFails(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	value, err := client.List("foo")
	c.Assert(err, gc.ErrorMatches, "'foo' is not a directory")
	c.Assert(value, gc.IsNil)
}

func (s *suite) TestListUnknownDirectory(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	value, err := client.List("foo")
	c.Assert(err, gc.ErrorMatches, "key 'foo' does not exist")
	c.Assert(value, gc.IsNil)
}
