package cache

import "github.com/curator-go/curator"

// TreeCacheListener represents listener for TreeCache changes
type TreeCacheListener interface {
	// Called when a change has occurred
	ChildEvent(client curator.CuratorFramework, event TreeCacheEvent) error
}

// childEventCallback is the callback type of ChildEvent within TreeCacheListener
type childEventCallback func(curator.CuratorFramework, TreeCacheEvent) error

// treeCacheListenerPrototype is the internal implementation of TreeCacheListener
type treeCacheListenerPrototype struct {
	childEvent childEventCallback
}

// ChildEvent is called when a change has occurred
func (l *treeCacheListenerPrototype) ChildEvent(client curator.CuratorFramework, event TreeCacheEvent) error {
	return l.childEvent(client, event)
}

// NewTreeCacheListener creates TreeCacheListener with given function
func NewTreeCacheListener(cb childEventCallback) TreeCacheListener {
	return &treeCacheListenerPrototype{cb}
}
