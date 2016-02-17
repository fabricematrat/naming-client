// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package namingclient

import (
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"gopkg.in/errgo.v1"
)

// ServiceType represents the type of the service.
type ServiceType string

const (
	// The name type for charm.
	Charm ServiceType = "charm"
	// The name type for model.
	Model ServiceType = "model"
	// The name type for page.
	Page ServiceType = "page"
)

// Client is a naming service client.
type Client struct {
	kapi    client.KeysAPI
	context context.Context
}

// NewClient creates a naming service client.
func NewClient(addr ...string) (*Client, error) {
	namingclient := &Client{}
	cfg := client.Config{
		Endpoints:               addr,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	namingclient.context = context.Background()
	namingclient.kapi = client.NewKeysAPI(c)
	return namingclient, nil
}

// Create creates a new name with given ServiceType value in the naming server.
func (c *Client) Create(name string, value ServiceType) error {
	_, err := c.kapi.Create(c.context, name, string(value))
	if err != nil && err.(client.Error).Code == client.ErrorCodeNodeExist {
		return errgo.Newf("key '%s' exists", name)
	}
	return err
}

// Update updates a name with given ServiceType value in the naming server.
func (c *Client) Update(name string, oldServiceType, newServiceType ServiceType) error {
	sopt := &client.SetOptions{
		PrevValue: string(oldServiceType),
		PrevExist: client.PrevExist,
	}
	_, err := c.kapi.Set(c.context, name, string(newServiceType), sopt)
	if err != nil && err.(client.Error).Code == client.ErrorCodeKeyNotFound {
		return errgo.Newf("key '%s' does not exist", name)
	}
	return err
}

// Read reads a name from the naming server, returning its ServiceType.
func (c *Client) Read(name string) (*ServiceType, error) {
	quorum := &client.GetOptions{Quorum: true}
	response, err := c.kapi.Get(c.context, name, quorum)
	if err != nil && err.(client.Error).Code == client.ErrorCodeKeyNotFound {
		return nil, errgo.Newf("key '%s' does not exist", name)
	} else if err != nil {
		return nil, err
	}
	if response.Node.Dir {
		return nil, errgo.Newf("cannot read directory '%s'", name)
	}

	value := response.Node.Value
	st := ServiceType(value)
	return &st, nil
}

// Delete deletes a name from the naming server. This will allow the name to be
// used again.
// CAUTION, THIS RECURSES DIRECTORIES!
func (c *Client) Delete(name string) error {
	_, err := c.kapi.Delete(c.context, name, nil)
	if err != nil && err.(client.Error).Code == client.ErrorCodeKeyNotFound {
		return errgo.Newf("key '%s' does not exist", name)
	} else if err != nil {
		return err
	}
	return nil
}

// List lists returns a map of string to service type.
func (c *Client) List(namespace string) (map[string]ServiceType, error) {
	response, err := c.kapi.Get(c.context, namespace, nil)
	if err != nil && err.(client.Error).Code == client.ErrorCodeKeyNotFound {
		return nil, errgo.Newf("key '%s' does not exist", namespace)
	} else if err != nil {
		return nil, err
	}
	if !response.Node.Dir {
		return nil, errgo.Newf("'%s' is not a directory", namespace)
	}
	nodes := response.Node.Nodes
	results := make(map[string]ServiceType, len(nodes))
	for _, inode := range nodes {
		if inode.Dir {
			continue
		}
		results[inode.Key] = ServiceType(inode.Value)
	}
	return results, nil
}
