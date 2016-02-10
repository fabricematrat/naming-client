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

func (s *suite) TestCreateAndRead(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	value, err := client.Read("foo")
	c.Assert(err, gc.IsNil)
	c.Assert(*value, gc.Equals, namingclient.Model)
}

func (s *suite) TestDelete(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	err = client.Delete("foo")
	c.Assert(err, gc.IsNil)
	value, err := client.Read("foo")
	c.Assert(err, gc.ErrorMatches, "key does not exist foo")
	c.Assert(value, gc.IsNil)
}

func (s *suite) TestReadUnknownKey(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	value, err := client.Read("nope")
	c.Assert(err, gc.ErrorMatches, "key does not exist nope")
	c.Assert(value, gc.IsNil)
}

func (s *suite) TestCreateExistingFails(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.ErrorMatches, "key exists foo")
}

func (s *suite) TestReadDirectoryShouldFail(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("mydir/foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	value, err := client.Read("mydir")
	c.Assert(err, gc.ErrorMatches, "cannot read directory mydir")
	c.Assert(value, gc.IsNil)
}

func (s *suite) TestUpdate(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Create("foo", namingclient.Model)
	c.Assert(err, gc.IsNil)
	err = client.Update("foo", namingclient.Model, namingclient.Charms)
	value, err := client.Read("foo")
	c.Assert(err, gc.IsNil)
	c.Assert(*value, gc.Equals, namingclient.Charms)
}

func (s *suite) TestUpdateNoKey(c *gc.C) {
	client, err := namingclient.NewClient(s.addr...)
	c.Assert(err, gc.IsNil)
	err = client.Update("foo", namingclient.Model, namingclient.Charms)
	c.Assert(err, gc.ErrorMatches, "key does not exist foo")
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
	c.Assert(err, gc.ErrorMatches, "not a directory foo")
	c.Assert(value, gc.IsNil)
}
