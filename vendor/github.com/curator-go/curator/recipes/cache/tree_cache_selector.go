package cache

import "strings"

// TreeCacheSelector controls which nodes a TreeCache processes.
// When iterating over the children of a parent node, a given node's children
// are queried only if TraverseChildren() returns true.
//
// When caching the list of nodes for a parent node, a given node is
// stored only if AcceptChild() returns true
type TreeCacheSelector interface {
	// TraverseChildren returns true if children of this path should be cached.
	// i.e. if false is returned, this node is not queried to
	// determine if it has children or not
	TraverseChildren(fullPath string) bool

	// AcceptChild returns true if this node should be returned from the cache
	AcceptChild(fullPath string) bool
}

// treeCacheSelectorPrototype contains functions for TreeCacheSelector implementation
type treeCacheSelectorPrototype struct {
	traverseChildren func(string) bool
	acceptChild      func(string) bool
}

// TraverseChildren calls inner traverseChildren()
func (s treeCacheSelectorPrototype) TraverseChildren(fullPath string) bool {
	return s.traverseChildren(fullPath)
}

// AcceptChild calls inner acceptChild()
func (s treeCacheSelectorPrototype) AcceptChild(fullPath string) bool {
	return s.acceptChild(fullPath)
}

// NewTreeCacheSelector creates a new TreeCacheSelector with given functions
func NewTreeCacheSelector(traverseChildren, acceptChild func(string) bool) TreeCacheSelector {
	return &treeCacheSelectorPrototype{traverseChildren, acceptChild}
}

// DefaultTreeCacheSelector returns true for all methods
var DefaultTreeCacheSelector = NewTreeCacheSelector(
	func(p string) bool { return true },
	func(p string) bool { return true },
)

// IgnoreBuiltinTreeCacheSelector ignores path starts with /zookeeper
// This could be useful if you use / as your root
var IgnoreBuiltinTreeCacheSelector = NewTreeCacheSelector(
	func(path string) bool {
		return !strings.HasPrefix(path, "/zookeeper")
	},
	func(path string) bool {
		return !strings.HasPrefix(path, "/zookeeper")
	},
)
