package cache

import (
	"errors"
	"testing"

	"github.com/curator-go/curator"
	"github.com/tevino/abool"
)

func TestRace(t *testing.T) {
	curator.NewClient("127.0.0.1:2181", nil).Start()
}

func TestSelector(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	// Create znodes
	tt.Create().ForPath("/root")
	tt.Create().ForPath("/root/n1-a")
	tt.Create().ForPath("/root/n1-b")
	tt.Create().ForPath("/root/n1-b/n2-a")
	tt.Create().ForPath("/root/n1-b/n2-b")
	tt.Create().ForPath("/root/n1-b/n2-b/n3-a")
	tt.Create().ForPath("/root/n1-c")
	tt.Create().ForPath("/root/n1-d")
	// Create selector
	selector := NewTreeCacheSelector(
		func(fullPath string) bool {
			return fullPath != "/root/n1-b/n2-b"
		},
		func(fullPath string) bool {
			return fullPath != "/root/n1-c"
		},
	)
	// Create TreeCache
	cache := tt.NewTreeCache("/root", selector)
	Assert(t, cache.Start() == nil)
	// Waiting for events
	tt.AssertEvent(TreeCacheEventNodeAdded, "/root")
	tt.AssertEvent(TreeCacheEventNodeAdded, "/root/n1-a")
	tt.AssertEvent(TreeCacheEventNodeAdded, "/root/n1-b")
	tt.AssertEvent(TreeCacheEventNodeAdded, "/root/n1-d")
	tt.AssertEvent(TreeCacheEventNodeAdded, "/root/n1-b/n2-a")
	tt.AssertEvent(TreeCacheEventNodeAdded, "/root/n1-b/n2-b")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
}

func TestStartup(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/test")
	tt.Create().ForPathWithData("/test/1", []byte("one"))
	tt.Create().ForPathWithData("/test/2", []byte("two"))
	tt.Create().ForPathWithData("/test/3", []byte("three"))
	tt.Create().ForPathWithData("/test/2/sub", []byte("two-sub"))
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil)
	Assert(t, cache.Start() == nil)
	// Wait for events
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/1", "one")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/2", "two")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/3", "three")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/2/sub", "two-sub")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
	// Check children
	tt.AssertChildren("/test", "1", "2", "3")

	tt.AssertChildrenEmpty("/test/1")
	tt.AssertChildren("/test/2", "sub")
	tt.AssertChildrenNotExist("/test/non-exist")
}

func TestCreateParentsDisabled(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create TreeCache with parent auto-create disabled
	cache := tt.NewTreeCache("/one/two/three", nil)
	Assert(t, cache.Start() == nil)
	// Wait for Initialized then nothing
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
	// Check the test path does not created
	tt.AssertNotExist("/one/two/three")
}

func TestCreateParentsEnabled(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create TreeCache with parent auto-create enabled
	cache := tt.NewTreeCache("/one/two/three", nil).SetCreateParentNodes(true)
	Assert(t, cache.Start() == nil)
	// Wait for events
	tt.AssertEvent(TreeCacheEventNodeAdded, "/one/two/three")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
	// Check the path is created automatically
	tt.AssertExist("/one/two/three")
}

func TestStartEmpty(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	cache := tt.NewTreeCache("/test", nil)
	Assert(t, cache.Start() == nil)
	// Assert Initialized
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Create znode
	tt.Create().ForPath("/test")
	// Assert added then nothing
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertNoMoreEvents()
}

func TestStartEmptyDeeper(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create TreeCache
	cache := tt.NewTreeCache("/test/foo/bar", nil)
	Assert(t, cache.Start() == nil)
	// Assert Initialized
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Create ancestor
	tt.Create().CreatingParentsIfNeeded().ForPath("/test/foo")
	// Assert nothing
	tt.AssertNoMoreEvents()
	// Create root
	tt.Create().ForPath("/test/foo/bar")
	// Assert Added then nothing
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test/foo/bar")
	tt.AssertNoMoreEvents()
}

func TestDepth0(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/test")
	tt.Create().ForPathWithData("/test/1", []byte("one"))
	tt.Create().ForPathWithData("/test/2", []byte("two"))
	tt.Create().ForPathWithData("/test/3", []byte("three"))
	tt.Create().ForPathWithData("/test/2/sub", []byte("two-sub"))
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil).SetMaxDepth(0)
	Assert(t, cache.Start() == nil)
	// Assert root added and nothing
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
	// Check root's children is empty
	tt.AssertChildrenEmpty("/test")
	// Check nodes that should be ignored by MaxDepth
	tt.AssertDataNotExist("/test/1")
	tt.AssertChildrenNotExist("/test/1")
	tt.AssertDataNotExist("/test/non_exist")
}

func TestDepth1(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	tt.Create().ForPath("/test")
	tt.Create().ForPathWithData("/test/1", []byte("one"))
	tt.Create().ForPathWithData("/test/2", []byte("two"))
	tt.Create().ForPathWithData("/test/3", []byte("three"))
	tt.Create().ForPathWithData("/test/2/sub", []byte("two-sub"))
	// Create TreeCache with MaxDepth 1
	cache := tt.NewTreeCache("/test", nil).SetMaxDepth(1)
	Assert(t, cache.Start() == nil)
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/1", "one")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/2", "two")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/3", "three")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()

	tt.AssertChildren("/test", "1", "2", "3")
	tt.AssertChildrenEmpty("/test/1")
	tt.AssertChildrenEmpty("/test/2")
	tt.AssertDataNotExist("/test/2/sub")
	tt.AssertChildrenNotExist("/test/2/sub")
	tt.AssertChildrenNotExist("/test/non_exist")
}

func TestDepth1Deeper(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/test")
	tt.Create().ForPath("/test/foo")
	tt.Create().ForPath("/test/foo/bar")
	tt.Create().ForPathWithData("/test/foo/bar/1", []byte("one"))
	tt.Create().ForPathWithData("/test/foo/bar/2", []byte("two"))
	tt.Create().ForPathWithData("/test/foo/bar/3", []byte("three"))
	tt.Create().ForPathWithData("/test/foo/bar/2/sub", []byte("two-sub"))
	// Create TreeCache
	cache := tt.NewTreeCache("/test/foo/bar", nil).SetMaxDepth(1)
	Assert(t, cache.Start() == nil)
	// Assert level 1 events
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test/foo/bar")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo/bar/1", "one")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo/bar/2", "two")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo/bar/3", "three")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
}

func TestAsyncInitialPopulation(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	tt.Create().ForPath("/test")
	tt.Create().ForPathWithData("/test/one", []byte("hey there"))
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil)
	Assert(t, cache.Start() == nil)
	// Assert events
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/one", "hey there")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
}

func TestFromRoot(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/test")
	tt.Create().ForPathWithData("/test/one", []byte("hey there"))
	// Create TreeCache
	cache := tt.NewTreeCache("/", IgnoreBuiltinTreeCacheSelector)
	Assert(t, cache.Start() == nil)
	// Assert node added
	tt.AssertEvent(TreeCacheEventNodeAdded, "/")
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/one", "hey there")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
	// Assert children & Data
	tt.AssertChildren("/", "test")
	tt.AssertChildren("/test", "one")
	tt.AssertChildrenEmpty("/test/one")
	tt.AssertData("/test/one", "hey there")
}

func TestFromRootWithDepth(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/test")
	tt.Create().ForPathWithData("/test/one", []byte("hey there"))
	// Create TreeCache with depth set to 1
	cache := tt.NewTreeCache("/", IgnoreBuiltinTreeCacheSelector).SetMaxDepth(1)
	Assert(t, cache.Start() == nil)
	// Assert first level events
	tt.AssertEvent(TreeCacheEventNodeAdded, "/")
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
	// Assert children & Data
	tt.AssertChildren("/", "test")
	tt.AssertChildrenEmpty("/test")
	tt.AssertDataNotExist("/test/one")
	tt.AssertChildrenNotExist("/test/one")
}

func TestWithNamespace(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/outer")
	tt.Create().ForPath("/outer/foo")
	tt.Create().ForPath("/outer/test")
	tt.Create().ForPathWithData("/outer/test/one", []byte("hey there"))
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil)
	// Set namespace
	cache.client = cache.client.UsingNamespace("outer")
	Assert(t, cache.Start() == nil)
	// Assert events with namespaced path
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/one", "hey there")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
	// Assert namespaced children and data
	tt.AssertChildren("/test", "one")
	tt.AssertChildrenEmpty("/test/one")
	tt.AssertData("/test/one", "hey there")
}

func TestWithNamespaceAtRoot(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/outer")
	tt.Create().ForPath("/outer/foo")
	tt.Create().ForPath("/outer/test")
	tt.Create().ForPathWithData("/outer/test/one", []byte("hey there"))
	// Create TreeCache
	cache := tt.NewTreeCache("/", nil)
	// Set namespace
	cache.client = cache.client.UsingNamespace("outer")
	Assert(t, cache.Start() == nil)
	// Assert namespaced events
	tt.AssertEvent(TreeCacheEventNodeAdded, "/")
	tt.AssertEvent(TreeCacheEventNodeAdded, "/foo")
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/one", "hey there")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
	tt.AssertChildren("/", "foo", "test")
	tt.AssertChildrenEmpty("/foo")
	tt.AssertChildren("/test", "one")
	tt.AssertChildrenEmpty("/test/one")
	tt.AssertData("/test/one", "hey there")
}

func TestSyncInitialPopulation(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil)
	Assert(t, cache.Start() == nil)
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Assert events
	tt.Create().ForPath("/test")
	tt.Create().ForPathWithData("/test/one", []byte("hey there"))
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/one", "hey there")
	tt.AssertNoMoreEvents()
}

func TestChildrenInitialized(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/test")
	tt.Create().ForPathWithData("/test/1", []byte("1"))
	tt.Create().ForPathWithData("/test/2", []byte("2"))
	tt.Create().ForPathWithData("/test/3", []byte("3"))
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil)
	Assert(t, cache.Start() == nil)
	// Assert events
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/1", "1")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/2", "2")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/3", "3")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	tt.AssertNoMoreEvents()
}

func TestUpdateWhenNotCachingData(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/test")
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil).SetCacheData(false)
	Assert(t, cache.Start() == nil)
	// Assert events
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Create znode with data
	tt.Create().ForPathWithData("/test/foo", []byte("first"))
	// Assert event with data
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo", "first")
	// Create znode with data
	tt.SetData().ForPathWithData("/test/foo", []byte("something new"))
	// Assert event with data
	tt.AssertEventWithData(TreeCacheEventNodeUpdated, "/test/foo", "something new")
	tt.AssertNoMoreEvents()
	// Assert empty data because we're not caching data
	tt.AssertDataEmpty("/test/foo")
}

func TestDeleteThenCreate(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/test")
	tt.Create().ForPathWithData("/test/foo", []byte("one"))
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil)
	Assert(t, cache.Start() == nil)
	// Assert added
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo", "one")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Delete znode
	tt.Delete().ForPath("/test/foo")
	// Assert removed
	tt.AssertEventWithData(TreeCacheEventNodeRemoved, "/test/foo", "one")
	// Create znode
	tt.Create().ForPathWithData("/test/foo", []byte("two"))
	// Assert added
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo", "two")
	// Delete znode
	tt.Delete().ForPath("/test/foo")
	// Assert removed
	tt.AssertEventWithData(TreeCacheEventNodeRemoved, "/test/foo", "two")
	// Create znode
	tt.Create().ForPathWithData("/test/foo", []byte("two"))
	// Assert added
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo", "two")
	tt.AssertNoMoreEvents()
}

func TestDeleteThenCreateRoot(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znodes
	tt.Create().ForPath("/test")
	tt.Create().ForPathWithData("/test/foo", []byte("one"))
	// Create TreeCache
	cache := tt.NewTreeCache("/test/foo", nil)
	Assert(t, cache.Start() == nil)
	// Assert added
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo", "one")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Delete znode
	tt.Delete().ForPath("/test/foo")
	// Assert removed
	tt.AssertEventWithData(TreeCacheEventNodeRemoved, "/test/foo", "one")
	// Create znode
	tt.Create().ForPathWithData("/test/foo", []byte("two"))
	// Assert added
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo", "two")
	// Delete znode
	tt.Delete().ForPath("/test/foo")
	// Assert removed
	tt.AssertEventWithData(TreeCacheEventNodeRemoved, "/test/foo", "two")
	// Create znode
	tt.Create().ForPathWithData("/test/foo", []byte("two"))
	// Assert added
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo", "two")
	tt.AssertNoMoreEvents()
}

// TODO: Uncomment this, if a session expired error could be simulated
//
// func TestKilledSession(t *testing.T) {
// 	// Create tester
// 	tt := NewTreeCacheTester(t).Start()
// 	defer tt.Stop()
// 	// Create znode
// 	tt.Create().ForPath("/test")
// 	// Create TreeCache
// 	cache := tt.NewTreeCache("/test", nil)
// 	Assert(t, cache.Start() == nil)
// 	// Assert added
// 	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
// 	tt.AssertEvent(TreeCacheEventInitialized, "")
// 	// Create znode
// 	tt.Create().ForPathWithData("/test/foo", []byte("foo"))
// 	// Assert added
// 	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/foo", "foo")
// 	// Create ephemeral znode
// 	tt.Create().WithMode(curator.EPHEMERAL).ForPathWithData("/test/me", []byte("data"))
// 	// Assert added
// 	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/me", "data")
// 	// Kill session
// 	tt.KillSession()
// 	// Assert connection events
// 	tt.AssertEvent(TreeCacheEventConnSuspended, "")
// 	//tt.AssertEvent(TreeCacheEventConnLost, "")
// 	tt.AssertEvent(TreeCacheEventConnReconnected, "")
// 	// Assert ephemeral node removed
// 	tt.AssertEventWithData(TreeCacheEventNodeRemoved, "/test/me", "data")
// 	tt.AssertEvent(TreeCacheEventInitialized, "")
// 	tt.AssertNoMoreEvents()
// }

func TestBasics(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	defer tt.Stop()
	// Create znode
	tt.Create().ForPath("/test")
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil)
	Assert(t, cache.Start() == nil)
	// Assert added
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Assert Children
	tt.AssertChildrenEmpty("/test")
	tt.AssertChildrenNotMatch("/t")
	tt.AssertChildrenNotExist("/testing")
	// Create znode with data
	tt.Create().ForPathWithData("/test/one", []byte("hey there"))
	// Assert added
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/one", "hey there")
	// Assert data
	tt.AssertChildren("/test", "one")
	tt.AssertData("/test/one", "hey there")
	// Assert children
	tt.AssertChildrenEmpty("/test/one")
	tt.AssertChildrenNotExist("/test/o")
	tt.AssertChildrenNotExist("/test/onely")
	// Set znode data
	tt.SetData().ForPathWithData("/test/one", []byte("sup!"))
	// Assert updated
	tt.AssertEventWithData(TreeCacheEventNodeUpdated, "/test/one", "sup!")
	// Assert children
	tt.AssertChildren("/test", "one")
	// Assert new data
	tt.AssertData("/test/one", "sup!")
	// Delete znode
	tt.Delete().ForPath("/test/one")
	// Assert removed
	tt.AssertEventWithData(TreeCacheEventNodeRemoved, "/test/one", "sup!")
	// Assert children
	tt.AssertChildrenEmpty("/test")
	tt.AssertNoMoreEvents()
}

func TestBasicsOnTwoCaches(t *testing.T) {
	// Create testers
	tt := NewTreeCacheTester(t).Start()
	tt2 := NewTreeCacheTesterWithCluster(t, tt.ZKCluster()).Start()
	defer tt.Stop()
	defer tt2.Stop()
	// Create TreeCaches
	cache := tt.NewTreeCache("/test", nil)
	cache2 := tt2.NewTreeCache("/test", nil)
	// Starts both TreeCaches
	cache.Start()
	cache2.Start()
	// Create znode
	tt.Create().ForPath("/test")
	// Assert added
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Assert added on second
	tt2.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt2.AssertEvent(TreeCacheEventInitialized, "")
	// Create znode with data
	tt.Create().ForPathWithData("/test/one", []byte("hey there"))
	// Assert added
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/one", "hey there")
	tt.AssertData("/test/one", "hey there")
	// Assert added on second
	tt2.AssertEventWithData(TreeCacheEventNodeAdded, "/test/one", "hey there")
	tt2.AssertData("/test/one", "hey there")
	// Set znode data
	tt.SetData().ForPathWithData("/test/one", []byte("sup!"))
	// Assert updated
	tt.AssertEventWithData(TreeCacheEventNodeUpdated, "/test/one", "sup!")
	tt.AssertData("/test/one", "sup!")
	// Assert updated on second
	tt2.AssertEventWithData(TreeCacheEventNodeUpdated, "/test/one", "sup!")
	tt2.AssertData("/test/one", "sup!")
	// Delete znode
	tt.Delete().ForPath("/test/one")
	// Assert removed
	tt.AssertEventWithData(TreeCacheEventNodeRemoved, "/test/one", "sup!")
	tt.AssertDataNotExist("/test/one")
	// Assert removed on second
	tt2.AssertEventWithData(TreeCacheEventNodeRemoved, "/test/one", "sup!")
	tt2.AssertDataNotExist("/test/one")
	// Assert no more
	tt.AssertNoMoreEvents()
	tt2.AssertNoMoreEvents()
}

func TestDeleteNodeAfterCloseDoesntCallExecutor(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	// Create znode
	tt.Create().ForPath("/test")
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil)
	cache.Start()
	// Assert added
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Create znode with data
	tt.Create().ForPathWithData("/test/one", []byte("hey there"))
	// Assert added
	tt.AssertEventWithData(TreeCacheEventNodeAdded, "/test/one", "hey there")
	tt.AssertData("/test/one", "hey there")
	// Close TreeCache
	cache.Stop()
	// Assert no more events after closing
	tt.AssertNoMoreEvents()
	t.Log("TreeCache Closed")
	// Delete znode
	tt.Delete().ForPath("/test/one")
	// Assert no more events because TreeCache was closed
	tt.AssertNoMoreEvents()
}

func TestServerNotStartedYet(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	// Stop the existing server.
	tt.ZKCluster().Servers[0].Srv.Stop()
	// Shutdown the existing client and re-create it started.
	tt.Client().Close()
	client := curator.NewClient(tt.serverAddr, nil)
	tt.SetClient(client)
	client.Start()
	// Start the client disconnected.
	cache := tt.NewTreeCache("/test", nil)
	cache.Start()
	// Assert nothing
	tt.AssertNoMoreEvents()
	// Now restart the server
	tt.ZKCluster().Servers[0].Srv.Start()
	// Assert initialized
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Create znode
	_, err := tt.Create().ForPath("/test")
	Assert(t, err == nil)
	// Assert added
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertNoMoreEvents()
}

func TestErrorListener(t *testing.T) {
	// Create tester
	tt := NewTreeCacheTester(t).Start()
	// Create znodes
	tt.Create().ForPath("/test")
	// Create TreeCache
	cache := tt.NewTreeCache("/test", nil)
	// Create a error for testing
	expectedError := errors.New("test error")
	// Register a listener that fail the test
	cache.Listenable().AddListener(NewTreeCacheListener(
		func(client curator.CuratorFramework, evt TreeCacheEvent) error {
			if evt.Type == TreeCacheEventNodeUpdated {
				return expectedError
			}
			return nil
		}))
	errHandlerCalled := abool.New()
	cache.UnhandledErrorListenable().AddListener(
		curator.NewUnhandledErrorListener(func(e error) {
			Assert(t, e == expectedError)
			Assert(t, errHandlerCalled.SetToIf(false, true))
		}))
	// Start TreeCache
	cache.Start()
	// Assert added
	tt.AssertEvent(TreeCacheEventNodeAdded, "/test")
	tt.AssertEvent(TreeCacheEventInitialized, "")
	// Set znode data
	tt.SetData().ForPathWithData("/test", []byte("hey there"))
	// Assert updated
	tt.AssertEventWithData(TreeCacheEventNodeUpdated, "/test", "hey there")
	tt.AssertNoMoreEvents()
	Assert(t, errHandlerCalled.IsSet())
}
