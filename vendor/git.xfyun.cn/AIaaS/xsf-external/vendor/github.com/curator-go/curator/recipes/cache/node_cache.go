package cache

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"unsafe"

	"github.com/curator-go/curator"
	"github.com/samuel/go-zookeeper/zk"
)

type NodeCacheListener interface {
	// Called when a change has occurred
	NodeChanged() error
}

type NodeCacheListenable interface {
	curator.Listenable /* [T] */

	AddListener(listener NodeCacheListener)

	RemoveListener(listener NodeCacheListener)
}

type NodeCacheListenerContainer struct {
	*curator.ListenerContainer
}

func (c *NodeCacheListenerContainer) AddListener(listener NodeCacheListener) {
	c.Add(listener)
}

func (c *NodeCacheListenerContainer) RemoveListener(listener NodeCacheListener) {
	c.Remove(listener)
}

// A utility that attempts to keep the data from a node locally cached.
// This class will watch the node, respond to update/create/delete events, pull down the data, etc.
// You can register a listener that will get notified when changes occur.
type NodeCache struct {
	client                  curator.CuratorFramework
	path                    string
	dataIsCompressed        bool
	ensurePath              curator.EnsurePath
	state                   curator.State
	isConnected             curator.AtomicBool
	data                    *ChildData
	connectionStateListener curator.ConnectionStateListener
	watcher                 curator.Watcher
	backgroundCallback      curator.BackgroundCallback
	listeners               *NodeCacheListenerContainer
}

func NewNodeCache(client curator.CuratorFramework, path string, dataIsCompressed bool) *NodeCache {
	c := &NodeCache{
		client:           client,
		path:             path,
		dataIsCompressed: dataIsCompressed,
		ensurePath:       client.NewNamespaceAwareEnsurePath(path).ExcludingLast(),
		listeners:        &NodeCacheListenerContainer{},
	}

	c.connectionStateListener = curator.NewConnectionStateListener(func(client curator.CuratorFramework, newState curator.ConnectionState) {
		if newState.Connected() {
			if c.isConnected.CompareAndSwap(false, true) {
				if err := c.reset(); err != nil {
					panic(fmt.Errorf("Trying to reset after reconnection, %s", err))
				}
			}
		} else {
			c.isConnected.Set(false)
		}
	})

	c.watcher = curator.NewWatcher(func(event *zk.Event) {
		c.reset()
	})

	c.backgroundCallback = func(client curator.CuratorFramework, event curator.CuratorEvent) error {
		return c.processBackgroundResult(event)
	}

	return c
}

// Start the cache. The cache is not started automatically. You must call this method.
func (c *NodeCache) Start() error {
	return c.StartAndInitalize(false)
}

// Same as Start() but gives the option of doing an initial build
func (c *NodeCache) StartAndInitalize(buildInitial bool) error {
	if !c.state.Change(curator.LATENT, curator.STARTED) {
		return fmt.Errorf("Cannot be started more than once")
	} else if err := c.ensurePath.Ensure(c.client.ZookeeperClient()); err != nil {
		return err
	}

	c.client.ConnectionStateListenable().AddListener(c.connectionStateListener)

	if buildInitial {
		if err := c.internalRebuild(); err != nil {
			return err
		}
	}

	return c.reset()
}

func (c *NodeCache) Close() error {
	if c.state.Change(curator.STARTED, curator.STOPPED) {
		c.listeners.Clear()
	}

	c.client.ConnectionStateListenable().RemoveListener(c.connectionStateListener)

	return nil
}

func (c *NodeCache) NodeCacheListenable() NodeCacheListenable {
	return c.listeners
}

func (c *NodeCache) internalRebuild() error {
	var stat zk.Stat

	builder := c.client.GetData()

	if c.dataIsCompressed {
		builder.Decompressed()
	}

	if data, err := builder.StoringStatIn(&stat).ForPath(c.path); err == nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&c.data)), unsafe.Pointer(&ChildData{c.path, &stat, data}))
	} else if err == zk.ErrNoNode {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&c.data)), nil)
	} else {
		return err
	}

	return nil
}

func (c *NodeCache) reset() error {
	if c.state.Value() == curator.STARTED && c.isConnected.Load() {
		_, err := c.client.CheckExists().UsingWatcher(c.watcher).InBackgroundWithCallback(c.backgroundCallback).ForPath(c.path)

		return err
	}

	return nil
}

func (c *NodeCache) processBackgroundResult(event curator.CuratorEvent) error {
	switch event.Type() {
	case curator.GET_DATA:
		if event.Err() == nil {
			c.setNewData(&ChildData{c.path, event.Stat(), event.Data()})
		}
	case curator.EXISTS:
		if event.Err() == zk.ErrNoNode {
			c.setNewData(nil)
		} else if event.Err() == nil {
			builder := c.client.GetData()

			if c.dataIsCompressed {
				builder.Decompressed()
			}

			builder.UsingWatcher(c.watcher).InBackgroundWithCallback(c.backgroundCallback).ForPath(c.path)
		}
	}

	return nil
}

func (c *NodeCache) setNewData(newData *ChildData) {
	previousData := (*ChildData)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.data)), unsafe.Pointer(newData)))

	if !reflect.DeepEqual(previousData, newData) {
		c.listeners.ForEach(func(listener interface{}) {
			listener.(NodeCacheListener).NodeChanged()
		})
	}
}

type RefreshMode int

const (
	STANDARD RefreshMode = iota
	FORCE_GET_DATA_AND_STAT
	POST_INITIALIZED
)

// A utility that attempts to keep all data from all children of a ZK path locally cached.
// This class will watch the ZK path, respond to update/create/delete events, pull down the data, etc.
// You can register a listener that will get notified when changes occur.
type PathChildrenCache struct {
	client                  curator.CuratorFramework
	path                    string
	cacheData               bool
	dataIsCompressed        bool
	ensurePath              curator.EnsurePath
	state                   curator.State
	connectionStateListener curator.ConnectionStateListener
	isConnected             curator.AtomicBool
}

func NewPathChildrenCache(client curator.CuratorFramework, path string, cacheData, dataIsCompressed bool) *PathChildrenCache {
	c := &PathChildrenCache{
		client:           client,
		path:             path,
		cacheData:        cacheData,
		dataIsCompressed: dataIsCompressed,
		ensurePath:       client.NewNamespaceAwareEnsurePath(path),
	}

	c.connectionStateListener = curator.NewConnectionStateListener(func(client curator.CuratorFramework, newState curator.ConnectionState) {
		if newState.Connected() {
			if c.isConnected.CompareAndSwap(false, true) {
				/*
					if err := c.reset(); err != nil {
						panic(fmt.Errorf("Trying to reset after reconnection, %s", err))
					}
				*/
			}
		} else {
			c.isConnected.Set(false)
		}
	})

	return c
}

func (c *PathChildrenCache) RefreshMode(mode RefreshMode) {
	c.ensurePath.Ensure(c.client.ZookeeperClient())
	/*
		c.client.GetChildren().UsingWatcher(c.childrenWatcher).InBackground(func(client CuratorFramework, event CuratorEvent) error {
			if c.state.Value() == STOPPED {

			}
		}).ForPath(c.path)
	*/
}
