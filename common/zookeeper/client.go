package zookeeper

import (
	"path/filepath"
	"sort"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/pkg/errors"
)

type Client struct {
	*zk.Conn
}

func New(servers []string) (*Client, error) {
	c, _, err := zk.Connect(servers, 10*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "fail to connect zk")
	}

	return &Client{c}, nil
}

func (c *Client) EnableLogging(enable bool) {
	c.SetLogger(logger{enable})
}

func (c *Client) walkNodes(path string, paths *[]string) error {
	var s string
	l, stat, err := c.Children(path)
	if err != nil {
		return err
	}

	if stat.NumChildren == 0 {
		*paths = append(*paths, path)
		return nil
	}

	for _, key := range l {
		if path == "/" {
			s = "/" + key
		} else {
			s = path + "/" + key
		}

		_, stat, err := c.Exists(s)
		if err != nil {
			return err
		}

		if stat.NumChildren == 0 {
			*paths = append(*paths, s)
		} else {
			if err := c.walkNodes(s, paths); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetZnodes get zk nodes in the path
func (c *Client) GetZnodes(path string) ([]string, error) {
	paths := make([]string, 0)

	_, _, err := c.Exists(path)
	if err != nil {
		return paths, err
	}

	err = c.walkNodes(path, &paths)
	if err != nil {
		return paths, err
	}

	sort.Strings(paths)

	return paths, nil
}

func (c *Client) DefaultCreate(path string, data []byte) (string, error) {
	return c.Create(path, data, 0, zk.WorldACL(zk.PermAll))
}

// ForceCreate force create multi-level node by acl
func (c *Client) ForceCreate(path string, data []byte, flags int32, acl []zk.ACL) error {
	return c.forceCreate(path, path, data, flags, acl)
}

func (c *Client) forceCreate(path, srcPath string, data []byte, flags int32, acl []zk.ACL) error {
	if path == "/" {
		return nil
	}

	exist, _, err := c.Exists(path)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	p := filepath.Dir(path)
	if err := c.forceCreate(p, srcPath, data, flags, acl); err != nil {
		return err
	}

	var d []byte
	if path == srcPath {
		d = data
	}

	_, err = c.Create(path, d, flags, acl)
	if err == zk.ErrNodeExists {
		return nil
	}

	return err
}

// ForceDelete force delete multi-level node
func (c *Client) ForceDelete(path string) error {
	return c.forceDelete(path)
}

func (c *Client) forceDelete(path string) error {
	if path == "/" {
		return zk.ErrInvalidPath
	}

	var s string
	l, stat, err := c.Children(path)
	if err != nil {
		return err
	}

	if stat.NumChildren == 0 {
		return c.Delete(path, stat.Version)
	}

	for _, key := range l {
		s = path + "/" + key

		if err := c.forceDelete(s); err != nil {
			return err
		}

		exist, _, err := c.Exists(s)
		if err != nil {
			return err
		}

		if !exist {
			continue
		}

		if err := c.Delete(s, stat.Version); err != nil {
			return err
		}
	}

	_, stat, err = c.Exists(path)
	if err != nil {
		return err
	}

	if stat.NumChildren == 0 {
		return c.Delete(path, stat.Version)
	}

	return nil
}
