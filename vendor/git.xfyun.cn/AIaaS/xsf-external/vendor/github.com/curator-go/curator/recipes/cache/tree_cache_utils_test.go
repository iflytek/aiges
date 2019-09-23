package cache

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/curator-go/curator"
	"github.com/samuel/go-zookeeper/zk"
)

// this file contains utils for testing of TreeCache.

// Assert asserts the given value is True for testing.
func Assert(t *testing.T, v bool) {
	if !v {
		_, fileName, line, _ := runtime.Caller(1)
		t.Errorf("\n Unexcepted: %s:%d", fileName, line)
	}
}

// KeysString asserts m is of type map[string]interface{}
// Returns a sorted slice of keys
func KeysString(m interface{}) []string {
	keyInterfaces := reflect.ValueOf(m).MapKeys()
	keys := make([]string, 0, len(keyInterfaces))
	for _, k := range keyInterfaces {
		keys = append(keys, k.String())
	}
	sort.Sort(sort.StringSlice(keys))
	return keys
}

// CacheTreeTester builds envirouments for TreeCache testing
type CacheTreeTester struct {
	sync.RWMutex
	curator.CuratorFramework
	t             *testing.T
	cache         *TreeCache
	events        []TreeCacheEvent
	listenerAdded bool
	zkCluster     *zk.TestCluster
	serverAddr    string
}

// assertEventTimeout is the timeout for asserting events
const assertEventTimeout = time.Second * 3

// NewTreeCacheTesterWithCluster creates CacheTreeTester with given ZK cluster
func NewTreeCacheTesterWithCluster(t *testing.T, cluster *zk.TestCluster) *CacheTreeTester {
	addr := fmt.Sprintf("127.0.0.1:%d", cluster.Servers[0].Port)
	t.Logf("Test server: %s", addr)
	return &CacheTreeTester{
		CuratorFramework: curator.NewClient(addr, nil),
		t:                t,
		events:           make([]TreeCacheEvent, 0),
		listenerAdded:    false,
		zkCluster:        cluster,
		serverAddr:       addr,
	}
}

// NewTreeCacheTester creates CacheTreeTester
func NewTreeCacheTester(t *testing.T) *CacheTreeTester {
	t.Logf("Starting test cluster")
	cluster, err := zk.StartTestCluster(1, nil, ioutil.Discard)
	if err != nil {
		t.Fatalf("Launch test cluster error: %s", err)
	}
	return NewTreeCacheTesterWithCluster(t, cluster)
}

// Start starts the inner client
func (ct *CacheTreeTester) Start() *CacheTreeTester {
	if err := ct.CuratorFramework.Start(); err != nil {
		ct.t.Fatal("Failed to start client: ", err)
	}
	if err := ct.CuratorFramework.BlockUntilConnectedTimeout(time.Second * 5); err != nil {
		ct.t.Fatal("Failed to connect: ", err)
	}
	// Remove anything under /
	children, err := ct.CuratorFramework.GetChildren().ForPath("/")
	if err != nil {
		ct.t.Fatal("Failed to ls /: ", err)
	}
	for _, child := range children {
		if child == "zookeeper" {
			continue
		}
		ct.t.Logf("Deleting: /%s", child)
		err = ct.CuratorFramework.Delete().DeletingChildrenIfNeeded().ForPath("/" + child)
		if err != nil {
			ct.t.Fatalf("Failed to delete path: /%s", child)
		}
	}
	return ct
}

// Stop should be called after finishing test
func (ct *CacheTreeTester) Stop() {
	if ct.CuratorFramework == nil {
		return
	}
	ct.zkCluster.Stop()
	ct.CloseCacheTree()
	ct.CuratorFramework.Close()
	ct.CuratorFramework = nil
}

// Client returns the internal CuratorFramework
func (ct *CacheTreeTester) Client() curator.CuratorFramework {
	return ct.CuratorFramework
}

// SetClient sets the internal CuratorFramework
func (ct *CacheTreeTester) SetClient(client curator.CuratorFramework) {
	ct.CuratorFramework = client
}

// ZKCluster returns the zk.TestCluster used inside
func (ct *CacheTreeTester) ZKCluster() *zk.TestCluster {
	return ct.zkCluster
}

// NewTreeCacheWithoutAttach creates TreeCache with logger set
func (ct *CacheTreeTester) NewTreeCacheWithoutAttach(root string, selector TreeCacheSelector) *TreeCache {
	c := NewTreeCache(ct.CuratorFramework, root, selector)
	c.SetLogger(newTestingLogger(ct.t))
	return c
}

// NewTreeCache creates TreeCache then attach it
func (ct *CacheTreeTester) NewTreeCache(root string, selector TreeCacheSelector) *TreeCache {
	c := ct.NewTreeCacheWithoutAttach(root, selector)
	ct.AttachCache(c)
	return c
}

// AttachCache sets then inits the inner cache of CacheTreeTester
func (ct *CacheTreeTester) AttachCache(cache *TreeCache) {
	if ct.cache == nil {
		ct.cache = cache
		ct.addListener()
	}
}

// getSessionID gets session ID from a curator.ZookeeperConnection
// Sadly the sessionID is a private field
func getSessionID(conn curator.ZookeeperConnection) int64 {
	rawConn := conn.(*zk.Conn)
	connVal := reflect.ValueOf(*rawConn)
	sessionID := connVal.FieldByName("sessionID")
	return sessionID.Int()
}

// changeSessionID changes the internal sessionID field inside a zk.Conn to given value
// NOTE: this may not work on all platforms, but sessionID is a private field, so sad
// NOTE: don't use this in production, it's for testing only
func changeSessionID(conn *zk.Conn, id int64) {
	pointerVal := reflect.ValueOf(conn)
	val := reflect.Indirect(pointerVal)
	sessionID := val.FieldByName("sessionID")
	ptrToSessionID := unsafe.Pointer(sessionID.UnsafeAddr())
	realPtrToSessionID := (*int64)(ptrToSessionID)
	*realPtrToSessionID = id
}

// KillSession tries to simulate a ZK session expired error
// See: http://wiki.apache.org/hadoop/ZooKeeper/FAQ#A4
//
// In the case of testing we want to cause a problem, so to explicitly expire a session an
// application connects to ZooKeeper, saves the session id and password, creates another
// ZooKeeper handle with that id and password, and then closes the new handle. Since both
// handles reference the same session, the close on second handle will invalidate the session
// causing a SESSION_EXPIRED on the first handle
func (ct *CacheTreeTester) KillSession() error {
	// Get current ZK connection
	conn, err := ct.CuratorFramework.ZookeeperClient().Conn()
	if err != nil {
		return err
	}
	// Create a duplicate connection with the same SessionID
	rawConn, evtUpdate, err := zk.Connect([]string{ct.serverAddr}, time.Second, func(c *zk.Conn) {
		sessionID := getSessionID(conn)
		changeSessionID(c, sessionID)
	})
	if err != nil {
		return err
	}
	for e := range evtUpdate {
		if e.State == zk.StateConnected {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	rawConn.Close()
	ct.t.Log("Session killed")
	return nil
}

// CloseCacheTree closes the cache and dettach it from the tester
func (ct *CacheTreeTester) CloseCacheTree() {
	if ct.cache == nil {
		return
	}
	ct.cache.Stop()
	ct.Lock()
	defer ct.Unlock()
	ct.events = make([]TreeCacheEvent, 0)
	ct.listenerAdded = false
	ct.cache = nil
}

// addListener adds listener to inner cache so incoming events are send to inner chan
func (ct *CacheTreeTester) addListener() {
	if ct.listenerAdded || ct.cache == nil {
		return
	}
	listener := NewTreeCacheListener(func(client curator.CuratorFramework, evt TreeCacheEvent) error {
		ct.t.Logf("Receiving event: %v", evt)
		ct.Lock()
		defer ct.Unlock()
		ct.events = append(ct.events, evt)
		return nil
	})
	ct.cache.Listenable().AddListener(listener)
	ct.listenerAdded = true
}

func eventMatch(e TreeCacheEvent, expectedType TreeCacheEventType, expectedPath string, expectedData []byte) bool {
	if e.Type != expectedType {
		return false
	}
	// The event must be TreeCacheEventInitialized if expectedPath is empty
	if expectedPath != "" {
		if e.Data == nil {
			return false
		}
		if e.Data.Path() != expectedPath {
			return false
		}
		if !bytes.Equal(e.Data.Data(), expectedData) {
			return false
		}
	}
	return true
}

// findEvent finds and pop an event from inner event set that fulfills the given condition
func (ct *CacheTreeTester) findEvent(expectedType TreeCacheEventType, expectedPath string, expectedData []byte) (evt TreeCacheEvent, found bool) {
	deadline := time.Now().Add(assertEventTimeout)
	for time.Now().Before(deadline) {
		ct.Lock()
		for i, e := range ct.events {
			if !eventMatch(e, expectedType, expectedPath, expectedData) {
				continue
			}
			// Remove e from events
			ct.events = append(ct.events[:i], ct.events[i+1:]...)
			ct.Unlock()
			evt = e
			found = true
			return
		}
		ct.Unlock()
	}
	ct.t.Fatalf("Waiting for event timed out: %s %s '%s'",
		expectedType, expectedPath, expectedData)
	return
}

// AssertDataNotExist asserts given path not exist with nil Data
func (ct *CacheTreeTester) AssertChildren(path string, expectedChildren ...string) {
	children, err := ct.cache.CurrentChildren(path)
	keys := KeysString(children)
	ct.t.Logf("AssertChildren: %s %v = %v %v", path, keys, expectedChildren, err)
	Assert(ct.t, err == nil)
	Assert(ct.t, reflect.DeepEqual(keys, expectedChildren))
}

// AssertChildrenEmpty asserts given path exists with zero children
func (ct *CacheTreeTester) AssertChildrenEmpty(path string) {
	children, err := ct.cache.CurrentChildren(path)
	ct.t.Logf("AssertChildrenEmpty: %s '%v' %v", path, children, err)
	Assert(ct.t, err == nil)
	Assert(ct.t, len(children) == 0)
}

// AssertChildrenError asserts CurrentChildren with given path will return an error of given type
func (ct *CacheTreeTester) AssertChildrenError(path string, errExpected error) {
	children, err := ct.cache.CurrentChildren(path)
	ct.t.Logf("AssertChildrenNotExist: %s '%v' %v", path, children, err)
	Assert(ct.t, err == errExpected)
	Assert(ct.t, len(children) == 0)
}

// AssertChildrenNotExist asserts given path not exist with zero children
func (ct *CacheTreeTester) AssertChildrenNotExist(path string) {
	ct.AssertChildrenError(path, ErrNodeNotFound)
}

// AssertChildrenNotMatch asserts CurrentChildren with given path will result in ErrRootNotMatch
func (ct *CacheTreeTester) AssertChildrenNotMatch(path string) {
	ct.AssertChildrenError(path, ErrRootNotMatch)
}

func (ct *CacheTreeTester) AssertData(path, expectedData string) {
	data, err := ct.cache.CurrentData(path)
	ct.t.Logf("AssertData: %s %s = %s %v", path, data, expectedData, err)
	Assert(ct.t, err == nil)
	Assert(ct.t, data != nil)
	Assert(ct.t, string(data.Data()) == expectedData)
}

// AssertDataEmpty asserts given path exists with nil Data
func (ct *CacheTreeTester) AssertDataEmpty(path string) {
	data, err := ct.cache.CurrentData(path)
	ct.t.Logf("AssertDataEmpty: %s '%s' %v", path, data, err)
	Assert(ct.t, err == nil)
	Assert(ct.t, data != nil)
	Assert(ct.t, data.Data() == nil)
}

// AssertDataNotExist asserts given path not exist with nil Data
func (ct *CacheTreeTester) AssertDataNotExist(path string) {
	data, err := ct.cache.CurrentData(path)
	ct.t.Logf("AssertDataNotExist: %s '%s' %v", path, data, err)
	Assert(ct.t, err == ErrNodeNotFound)
	Assert(ct.t, data == nil)
}

// AssertExist asserts given path exists
func (ct *CacheTreeTester) AssertExist(path string) {
	stat, err := ct.CuratorFramework.CheckExists().ForPath(path)
	ct.t.Logf("AssertExist: %s", path)
	Assert(ct.t, err == nil)
	Assert(ct.t, stat != nil)
}

// AssertNotExist asserts given path not exist
func (ct *CacheTreeTester) AssertNotExist(path string) {
	stat, err := ct.CuratorFramework.CheckExists().ForPath(path)
	ct.t.Logf("AssertNotExist: %s", path)
	Assert(ct.t, err == nil)
	Assert(ct.t, stat == nil)
}

// AssertEvent asserts an event's comming
func (ct *CacheTreeTester) AssertEvent(expectedType TreeCacheEventType, expectedPath string) {
	ct.RLock()
	ct.t.Logf("AssertEvent: [%s %s] in %v", expectedType, expectedPath, ct.events)
	ct.RUnlock()
	_, found := ct.findEvent(expectedType, expectedPath, nil)
	Assert(ct.t, found)
}

// AssertEvent asserts an event's comming with given data
func (ct *CacheTreeTester) AssertEventWithData(expectedType TreeCacheEventType, expectedPath string, expectedData string) {
	ct.RLock()
	ct.t.Logf("AssertEvent: [%s %s %v] in %v", expectedType, expectedPath, expectedData, ct.events)
	ct.RUnlock()
	_, found := ct.findEvent(expectedType, expectedPath, []byte(expectedData))
	Assert(ct.t, found)
}

// AssertNoMoreEvents asserts the inner event set is empty
func (ct *CacheTreeTester) AssertNoMoreEvents() {
	ct.RLock()
	defer ct.RUnlock()
	ct.t.Logf("AssertNoMoreEvents: %v", ct.events)
	Assert(ct.t, len(ct.events) == 0)
}

// testingLogger is a Logger used during testing
type testingLogger struct {
	*testing.T
}

// newTestingLogger creates a testingLogger
func newTestingLogger(t *testing.T) *testingLogger {
	return &testingLogger{t}
}

// Printf outputs log by calling T.Logf
func (l *testingLogger) Printf(format string, args ...interface{}) {
	l.T.Logf(format, args...)
}

// Debugf outputs log by calling T.Logf
func (l *testingLogger) Debugf(format string, args ...interface{}) {
	l.T.Logf(format, args...)
}
