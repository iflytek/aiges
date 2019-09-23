package cache

// This is a port of Apache Curator's TreeCache recipe.
// See also: http://curator.apache.org/curator-recipes/tree-cache.html
//
// Repo: https://git-wip-us.apache.org/repos/asf/curator.git
// Commit 6cebfc13ccdcd9cb2a5b576fa369b957e651752a
// Message: CURATOR-337: Do not attempt to release a mutex unless it has actually been acquired

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/curator-go/curator"
	"github.com/tevino/abool"
)

// Logger provides customized logging within TreeCache.
type Logger interface {
	Printf(string, ...interface{})
	Debugf(string, ...interface{})
}

// DummyLogger is a Logger does nothing.
type DummyLogger struct{}

// Printf does nothing.
func (l DummyLogger) Printf(string, ...interface{}) {}

// Debugf does nothing.
func (l DummyLogger) Debugf(string, ...interface{}) {}

// TreeCacheListenable represents a container of TreeCacheListener(s).
type TreeCacheListenable interface {
	curator.Listenable

	AddListener(TreeCacheListener)
	RemoveListener(TreeCacheListener)
}

// TreeCacheListenerContainer is a container of TreeCacheListener.
type TreeCacheListenerContainer struct {
	curator.ListenerContainer
}

// AddListener adds a listener to the container.
func (c *TreeCacheListenerContainer) AddListener(listener TreeCacheListener) {
	c.Add(listener)
}

// RemoveListener removes a listener to the container.
func (c *TreeCacheListenerContainer) RemoveListener(listener TreeCacheListener) {
	c.Remove(listener)
}

// TreeCache is a a utility that attempts to keep all data from all children of a ZK path locally cached.
// It will watch the ZK path, respond to update/create/delete events, pull down the data, etc.
// You can register a listener that will get notified when changes occur.
//
// NOTE: It's not possible to stay transactionally in sync. Users of this class must
// be prepared for false-positives and false-negatives. Additionally, always use the version number
// when updating data to avoid overwriting another process' change.
type TreeCache struct {
	// Tracks the number of outstanding background requests in flight. The first time this count reaches 0, we publish the initialized event.
	outstandingOps          uint64
	isInitialized           *abool.AtomicBool
	root                    *TreeNode
	client                  curator.CuratorFramework
	cacheData               bool
	maxDepth                int
	selector                TreeCacheSelector
	listeners               TreeCacheListenerContainer
	errorListeners          curator.UnhandledErrorListenerContainer
	state                   curator.State
	connectionStateListener curator.ConnectionStateListener
	logger                  Logger
	createParentNodes       bool
}

// NewTreeCache creates a TreeCache for the given client and path with default options.
//
// If the client is namespaced, all operations on the resulting TreeCache will be in terms of
// the namespace, including all published events.  The given path is the root at which the
// TreeCache will watch and explore.  If no node exists at the given path, the TreeCache will
// be initially empty.
func NewTreeCache(client curator.CuratorFramework, root string, selector TreeCacheSelector) *TreeCache {
	if selector == nil {
		selector = DefaultTreeCacheSelector
	}
	tc := &TreeCache{
		isInitialized: abool.New(),
		client:        client,
		maxDepth:      math.MaxInt32,
		cacheData:     true,
		selector:      selector,
		state:         curator.LATENT,
		logger:        &DummyLogger{},
	}
	tc.root = NewTreeNode(tc, root, nil)
	tc.connectionStateListener = curator.NewConnectionStateListener(
		func(client curator.CuratorFramework, newState curator.ConnectionState) {
			tc.handleStateChange(newState)
		})
	return tc
}

// Start starts the TreeCache.
// The cache is not started automatically. You must call this method.
func (tc *TreeCache) Start() error {
	if !tc.state.Change(curator.LATENT, curator.STARTED) {
		return errors.New("already started")
	}
	if tc.createParentNodes {
		_, err := tc.client.Create().CreatingParentsIfNeeded().ForPath(tc.root.path)
		if err != nil {
			return fmt.Errorf("Failed to create parents: %s", err)
		}
	}
	tc.client.ConnectionStateListenable().AddListener(tc.connectionStateListener)
	if tc.client.ZookeeperClient().Connected() {
		tc.root.wasCreated()
	}
	return nil
}

// SetCacheData sets whether or not to cache byte data per node, default true.
// NOTE: When this set to false, the events still contain data of znode
// but you can't query them by TreeCache.CurrentData/CurrentChildren
func (tc *TreeCache) SetCacheData(cacheData bool) *TreeCache {
	tc.cacheData = cacheData
	return tc
}

// SetMaxDepth sets the maximum depth to explore/watch.
// Set to 0 will watch only the root node.
// Set to 1 will watch the root node and its immediate children.
// Default to math.MaxInt32.
func (tc *TreeCache) SetMaxDepth(depth int) *TreeCache {
	tc.maxDepth = depth
	return tc
}

// SetCreateParentNodes sets whether to auto-create parent nodes for the cached path.
// By default, TreeCache does not do this.
// Note: Parent nodes is only created when Start() is called.
func (tc *TreeCache) SetCreateParentNodes(yes bool) *TreeCache {
	tc.createParentNodes = yes
	return tc
}

// SetLogger sets the inner Logger of TreeCache.
func (tc *TreeCache) SetLogger(l Logger) *TreeCache {
	tc.logger = l
	return tc
}

// Stop stops the cache.
func (tc *TreeCache) Stop() {
	if tc.state.Change(curator.STARTED, curator.STOPPED) {
		tc.client.ConnectionStateListenable().RemoveListener(tc.connectionStateListener)
		tc.listeners.Clear()
		tc.root.wasDeleted()
	}
}

// Listenable returns the cache listeners.
func (tc *TreeCache) Listenable() TreeCacheListenable {
	return &tc.listeners
}

// UnhandledErrorListenable allows catching unhandled errors in asynchornous operations.
func (tc *TreeCache) UnhandledErrorListenable() curator.UnhandledErrorListenable {
	return &tc.errorListeners
}

// ErrNodeNotFound indicates a node can not be found.
var ErrNodeNotFound = errors.New("node not found")

// ErrRootNotMatch indicates the root path does not match.
var ErrRootNotMatch = errors.New("root path not match")

// ErrNodeNotLive indicates the state of node is not LIVE.
var ErrNodeNotLive = errors.New("node state is not LIVE")

// findNode finds the node which matches the given path.
// ErrRootNotMatch is returned if the given path doesn't share a same root with
// the TreeCache.
// ErrNodeNotFound is returned if the given path can not be found.
func (tc *TreeCache) findNode(path string) (*TreeNode, error) {
	if !strings.HasPrefix(path, tc.root.path) {
		return nil, ErrRootNotMatch
	}

	path = strings.TrimPrefix(path, tc.root.path)
	current := tc.root
	for _, part := range strings.Split(path, "/") {
		if part == "" {
			continue
		}
		next, exists := current.FindChild(part)
		if !exists {
			return nil, ErrNodeNotFound
		}
		current = next
	}

	return current, nil
}

// CurrentChildren returns the current set of children at the given full path, mapped by child name.
// There are no guarantees of accuracy; this is merely the most recent view of the data.
// If there is no node at this path, ErrNodeNotFound is returned.
func (tc *TreeCache) CurrentChildren(fullPath string) (map[string]*ChildData, error) {
	node, err := tc.findNode(fullPath)
	if err != nil {
		return nil, err
	}
	if node.state.Load() != NodeStateLIVE {
		return nil, ErrNodeNotLive
	}

	children := node.Children()
	m := make(map[string]*ChildData, len(children))
	for child, childNode := range children {
		// Double-check liveness after retreiving data.
		childData := childNode.ChildData()
		if childData != nil && childNode.state.Load() == NodeStateLIVE {
			m[child] = childData
		}
	}

	// Double-check liveness after retreiving children.
	if node.state.Load() != NodeStateLIVE {
		return nil, ErrNodeNotLive
	}
	return m, nil
}

// CurrentData returns the current data for the given full path.
// There are no guarantees of accuracy. This is merely the most recent view of the data.
// If there is no node at the given path, ErrNodeNotFound is returned.
func (tc *TreeCache) CurrentData(fullPath string) (*ChildData, error) {
	node, err := tc.findNode(fullPath)
	if err != nil {
		return nil, err
	}
	if node.state.Load() != NodeStateLIVE {
		return nil, ErrNodeNotLive
	}

	return node.ChildData(), nil
}

// callListeners calls all listeners with given event.
// Error is handled by handleException().
func (tc *TreeCache) callListeners(evt TreeCacheEvent) {
	tc.listeners.ForEach(func(listener interface{}) {
		if err := listener.(TreeCacheListener).ChildEvent(tc.client, evt); err != nil {
			tc.handleException(err)
		}
	})
}

// handleException sends an exception to any listeners, or else log the error if there are none.
func (tc *TreeCache) handleException(e error) {
	if tc.errorListeners.Len() == 0 {
		tc.logger.Printf("%s", e)
		return
	}
	tc.errorListeners.ForEach(func(listener interface{}) {
		listener.(curator.UnhandledErrorListener).UnhandledError(e)
	})
}

func (tc *TreeCache) handleStateChange(newState curator.ConnectionState) {
	switch newState {
	case curator.SUSPENDED:
		tc.publishEvent(TreeCacheEventConnSuspended, nil)
	case curator.LOST:
		tc.isInitialized.UnSet()
		tc.publishEvent(TreeCacheEventConnLost, nil)
	case curator.CONNECTED:
		tc.root.wasCreated()
	case curator.RECONNECTED:
		if err := tc.root.wasReconnected(); err == nil {
			tc.publishEvent(TreeCacheEventConnReconnected, nil)
		}
	}
}

// publishEvent publish an event with given type and data to all listeners.
func (tc *TreeCache) publishEvent(tp TreeCacheEventType, data *ChildData) {
	if tc.state.Value() != curator.STOPPED {
		evt := TreeCacheEvent{Type: tp, Data: data}
		tc.logger.Debugf("publishEvent: %v", evt)
		go tc.callListeners(evt)
	}
}
