package zookeeper

import (
	"bytes"
	"strings"
	"unicode"

	errors "git.xfyun.cn/AIaaS/finder-go/errors"
	"github.com/cooleric/go-zookeeper/zk"
)

const (
	DEFAULT_VERSION       = -1
	PATH_SEPARATOR        = "/"
	PERSISTENT            = 0
	PERSISTENT_SEQUENTIAL = zk.FlagSequence
	EPHEMERAL             = zk.FlagEphemeral
	EPHEMERAL_SEQUENTIAL  = zk.FlagEphemeral + zk.FlagSequence
)

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

func makeDirs(conn *zk.Conn, path string, lastNode bool) error {
	if err := validatePath(path); err != nil {
		return err
	}

	pos := 1 // skip first slash, root is guaranteed to exist
	for pos < len(path) {
		if idx := strings.Index(path[pos+1:], PATH_SEPARATOR); idx == -1 {
			if lastNode {
				pos = len(path)
			} else {
				return nil
			}
		} else {
			pos += idx + 1
		}

		subPath := path[:pos]
		if exist, _, err := conn.Exists(subPath); err != nil {
			return err
		} else if !exist {
			if _, err := conn.Create(subPath, []byte{}, PERSISTENT, zk.WorldACL(zk.PermAll)); err != nil && err != zk.ErrNodeExists {
				return err
			}
		}
	}

	return nil
}

// Recursively deletes children of a node.
func recursiveDelete(conn *zk.Conn, path string, deleteSelf bool) error {
	if err := validatePath(path); err != nil {
		return err
	}

	if children, _, err := conn.Children(path); err != nil {
		return err
	} else {
		for _, child := range children {
			if err := recursiveDelete(conn, joinPath(path, child), true); err != nil {
				return err
			}
		}
	}

	if deleteSelf {
		if err := conn.Delete(path, DEFAULT_VERSION); err != nil {
			switch err {
			case zk.ErrNotEmpty:
				return recursiveDelete(conn, path, true)
			case zk.ErrNoNode:
				return nil
			default:
				return err
			}
		}
	}

	return nil
}

// Given a parent and a child node, join them in the given path
func joinPath(parent string, children ...string) string {
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

// Validate the provided znode path string
func validatePath(path string) error {
	if len(path) == 0 {
		return errors.NewFinderError(errors.ZkPathCannotNil)
	}

	if !strings.HasPrefix(path, PATH_SEPARATOR) {
		return errors.NewFinderError(errors.ZkPathPrefixIllegal)
	}

	if len(path) == 1 {
		return nil
	}

	if strings.HasSuffix(path, PATH_SEPARATOR) {
		return errors.NewFinderError(errors.ZkPathSuffixIllegal)
	}

	lastc := '/'

	for i, c := range path {
		if i == 0 {
			continue
		} else if c == 0 {
			return errors.NewFinderError(errors.ZkPathNullCharacterNotAllowed)
		} else if c == '/' && lastc == '/' {
			return errors.NewFinderError(errors.ZkPathEmptyNodeNameNotAllowed)
		} else if c == '.' && lastc == '.' {
			if path[i-2] == '/' && (i+1 == len(path) || path[i+1] == '/') {
				return errors.NewFinderError(errors.ZkPathRelativePathNotAllowed)
			}
		} else if c == '.' {
			if path[i-1] == '/' && (i+1 == len(path) || path[i+1] == '/') {
				return errors.NewFinderError(errors.ZkPathRelativePathNotAllowed)
			}
		} else if unicode.In(c, invalidCharaters) {
			return errors.NewFinderError(errors.ZkPathInvalidCharacterNotAllowed)
		}

		lastc = c
	}

	return nil
}

func getNodeFromPath(path string) string {
	if idx := strings.LastIndex(path, PATH_SEPARATOR); idx < 0 {
		return path
	} else if idx+1 >= len(path) {
		return ""
	} else {
		return path[idx+1:]
	}
}
