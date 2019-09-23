package curator

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"
	"unicode"

	"github.com/cooleric/go-zookeeper/zk"
)

const (
	PATH_SEPARATOR = "/"
)

type PathAndNode struct {
	Path, Node string
}

// Given a full path, return the node name. i.e. "/one/two/three" will return "three"
func GetNodeFromPath(path string) string {
	if idx := strings.LastIndex(path, PATH_SEPARATOR); idx < 0 {
		return path
	} else if idx+1 >= len(path) {
		return ""
	} else {
		return path[idx+1:]
	}
}

// Given a full path, return the the individual parts, without slashes.
func SplitPath(path string) (*PathAndNode, error) {
	if idx := strings.LastIndex(path, PATH_SEPARATOR); idx < 0 {
		return &PathAndNode{path, ""}, nil
	} else if idx > 0 {
		return &PathAndNode{path[:idx], path[idx+1:]}, nil
	} else {
		return &PathAndNode{PATH_SEPARATOR, path[idx+1:]}, nil
	}
}

// Given a parent and a child node, join them in the given path
func JoinPath(parent string, children ...string) string {
	path := new(bytes.Buffer)

	if len(parent) > 0 {
		if !strings.HasPrefix(parent, PATH_SEPARATOR) {
			path.WriteString(PATH_SEPARATOR)
		}

		if strings.HasSuffix(parent, PATH_SEPARATOR) {
			path.WriteString(parent[:len(parent)-1])
		} else {
			path.WriteString(parent)
		}
	}

	for _, child := range children {
		if len(child) == 0 || child == PATH_SEPARATOR {
			if path.Len() == 0 {
				path.WriteString(PATH_SEPARATOR)
			}
		} else {
			path.WriteString(PATH_SEPARATOR)

			if strings.HasPrefix(child, PATH_SEPARATOR) {
				child = child[1:]
			}

			if strings.HasSuffix(child, PATH_SEPARATOR) {
				child = child[:len(child)-1]
			}

			path.WriteString(child)
		}
	}

	return path.String()
}

var (
	invalidCharaters = &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0000, Hi: 0x001f, Stride: 1},
			{Lo: 0x007f, Hi: 0x009F, Stride: 1},
			{Lo: 0xd800, Hi: 0xf8ff, Stride: 1},
			{Lo: 0xfff0, Hi: 0xffff, Stride: 1},
		},
	}
)

// Validate the provided znode path string
func ValidatePath(path string) error {
	if len(path) == 0 {
		return errors.New("Path cannot be null")
	}

	if !strings.HasPrefix(path, PATH_SEPARATOR) {
		return errors.New("Path must start with / character")
	}

	if len(path) == 1 {
		return nil
	}

	if strings.HasSuffix(path, PATH_SEPARATOR) {
		return errors.New("Path must not end with / character")
	}

	lastc := '/'

	for i, c := range path {
		if i == 0 {
			continue
		} else if c == 0 {
			return fmt.Errorf("null character not allowed @ %d", i)
		} else if c == '/' && lastc == '/' {
			return fmt.Errorf("empty node name specified @ %d", i)
		} else if c == '.' && lastc == '.' {
			if path[i-2] == '/' && (i+1 == len(path) || path[i+1] == '/') {
				return fmt.Errorf("relative paths not allowed @ %d", i)
			}
		} else if c == '.' {
			if path[i-1] == '/' && (i+1 == len(path) || path[i+1] == '/') {
				return fmt.Errorf("relative paths not allowed @ %d", i)
			}
		} else if unicode.In(c, invalidCharaters) {
			return fmt.Errorf("invalid charater @ %d", i)
		}

		lastc = c
	}

	return nil
}

// Make sure all the nodes in the path are created
func MakeDirs(conn ZookeeperConnection, path string, makeLastNode bool, aclProvider ACLProvider) error {
	if err := ValidatePath(path); err != nil {
		return err
	}

	pos := 1 // skip first slash, root is guaranteed to exist

	for pos < len(path) {
		if idx := strings.Index(path[pos+1:], PATH_SEPARATOR); idx == -1 {
			if makeLastNode {
				pos = len(path)
			} else {
				return nil
			}
		} else {
			pos += idx + 1
		}

		subPath := path[:pos]

		if exists, _, err := conn.Exists(subPath); err != nil {
			return err
		} else if !exists {
			var acls []zk.ACL

			if aclProvider != nil {
				if acls = aclProvider.GetAclForPath(subPath); len(acls) == 0 {
					acls = aclProvider.GetDefaultAcl()
				}
			}

			if acls == nil {
				acls = OPEN_ACL_UNSAFE
			}

			if _, err := conn.Create(subPath, []byte{}, int32(PERSISTENT), acls); err != nil && err != zk.ErrNodeExists {
				return err
			}
		}
	}

	return nil
}

// Recursively deletes children of a node.
func DeleteChildren(conn ZookeeperConnection, path string, deleteSelf bool) error {
	if err := ValidatePath(path); err != nil {
		return err
	}

	if children, _, err := conn.Children(path); err != nil {
		return err
	} else {
		for _, child := range children {
			if err := DeleteChildren(conn, JoinPath(path, child), true); err != nil {
				return err
			}
		}
	}

	if deleteSelf {
		if err := conn.Delete(path, -1); err != nil {
			switch err {
			case zk.ErrNotEmpty:
				return DeleteChildren(conn, path, true)
			case zk.ErrNoNode:
				return nil
			default:
				return err
			}
		}
	}

	return nil
}

type EnsurePath interface {
	// First time, synchronizes and makes sure all nodes in the path are created.
	// Subsequent calls with this instance are NOPs.
	Ensure(client CuratorZookeeperClient) error

	// Returns a view of this EnsurePath instance that does not make the last node.
	ExcludingLast() EnsurePath
}

type EnsurePathHelper interface {
	Ensure(client CuratorZookeeperClient, path string, makeLastNode bool) error
}

type ensurePathHelper struct {
	owner   *ensurePath
	lock    sync.Mutex
	started bool
}

func (h *ensurePathHelper) Ensure(client CuratorZookeeperClient, path string, makeLastNode bool) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if !h.started {
		_, err := client.NewRetryLoop().CallWithRetry(func() (interface{}, error) {
			if conn, err := client.Conn(); err != nil {
				return nil, err
			} else if err := MakeDirs(conn, path, makeLastNode, h.owner.aclProvider); err != nil {
				return nil, err
			} else {
				return nil, nil
			}
		})

		h.started = true

		h.owner.helper = nil

		return err
	}

	return nil
}

// Utility to ensure that a particular path is created.
type ensurePath struct {
	path         string
	aclProvider  ACLProvider
	makeLastNode bool
	helper       EnsurePathHelper
}

func NewEnsurePath(path string) *ensurePath {
	return NewEnsurePathWithAclAndHelper(path, nil, nil)
}

func NewEnsurePathWithAcl(path string, aclProvider ACLProvider) *ensurePath {
	return NewEnsurePathWithAclAndHelper(path, aclProvider, nil)
}

func NewEnsurePathWithAclAndHelper(path string, aclProvider ACLProvider, helper EnsurePathHelper) *ensurePath {
	p := &ensurePath{
		path:         path,
		aclProvider:  aclProvider,
		makeLastNode: true,
	}

	if helper == nil {
		p.helper = &ensurePathHelper{owner: p}
	} else {
		p.helper = helper
	}

	return p
}

func (p *ensurePath) ExcludingLast() EnsurePath {
	return &ensurePath{
		path:         p.path,
		aclProvider:  p.aclProvider,
		makeLastNode: false,
		helper:       p.helper,
	}
}

func (p *ensurePath) Ensure(client CuratorZookeeperClient) error {
	if p.helper != nil {
		return p.helper.Ensure(client, p.path, p.makeLastNode)
	}

	return nil
}
