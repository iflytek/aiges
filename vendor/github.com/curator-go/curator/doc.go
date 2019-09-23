/*
Curator-go is a Golang porting for Curator, make it easy to access Zookeeper


Learn ZooKeeper


Curator-go users are assumed to know ZooKeeper. A good place to start is http://zookeeper.apache.org/doc/trunk/zookeeperStarted.html


Using Curator


The Curator-go are available from github.com.

	$ go get github.com/curator-go/curator

You can easily include Curator-go into your code.

	import (
	    "github.com/curator-go/curator"
	)


Getting a Connection


Curator uses Fluent Style. If you haven't used this before, it might seem odd so it's suggested that you familiarize yourself with the style.

Curator connection instances (CuratorFramework) are allocated from the CuratorFrameworkBuilder. You only need one CuratorFramework object for each ZooKeeper cluster you are connecting to:

	curator.NewClient(connString, retryPolicy)

This will create a connection to a ZooKeeper cluster using default values. The only thing that you need to specify is the retry policy. For most cases, you should use:

	retryPolicy := curator.NewExponentialBackoffRetry(time.Second, 3, 15*time.Second)

	client := curator.NewClient(connString, retryPolicy)

	client.Start()
	defer client.Close()

The client must be started (and closed when no longer needed).


Calling ZooKeeper Directly


Once you have a CuratorFramework instance, you can make direct calls to ZooKeeper in a similar way to using the raw ZooKeeper object provided in the ZooKeeper distribution. E.g.:

	client.Create().ForPathWithData(path, payload)

The benefit here is that Curator manages the ZooKeeper connection and will retry operations if there are connection problems.


Recipes


Distributed Lock


	lock := curator.NewInterProcessMutex(client, lockPath)

	if ( lock.Acquire(maxWait, waitUnit) )
	{
	    defer lock.Release()

	    // do some work inside of the critical section here
	}


Leader Election


	listener := curator.NewLeaderSelectorListener(func(CuratorFramework client) error {
	    // this callback will get called when you are the leader
	    // do whatever leader work you need to and only exit
	    // this method when you want to relinquish leadership
	}))

	selector := curator.NewLeaderSelector(client, path, listener)
	selector.AutoRequeue()  // not required, but this is behavior that you will probably expect
	selector.Start()


Generic API


Curator provides generic API for builder

	type Pathable[T] interface {
	    // Commit the currently building operation using the given path
	    ForPath(path string) (T, error)
	}

	type PathAndBytesable[T] interface {
	    Pathable[T]

	    // Commit the currently building operation using the given path and data
	    ForPathWithData(path string, payload []byte) (T, error)
	}

	type Compressible[T] interface {
	    // Cause the data to be compressed using the configured compression provider
	    Compressed() T
	}

	type Decompressible[T] interface {
	    // Cause the data to be de-compressed using the configured compression provider
	    Decompressed() T
	}

	type CreateModable[T] interface {
	    // Set a create mode - the default is CreateMode.PERSISTENT
	    WithMode(mode CreateMode) T
	}

	type ACLable[T] interface {
	    // Set an ACL list
	    WithACL(acl ...zk.ACL) T
	}

	type Versionable[T] interface {
	    // Use the given version (the default is -1)
	    WithVersion(version int) T
	}

	type Statable[T] interface {
	    // Have the operation fill the provided stat object
	    StoringStatIn(*zk.Stat) T
	}

	type ParentsCreatable[T] interface {
	    // Causes any parent nodes to get created if they haven't already been
	    CreatingParentsIfNeeded() T
	}

	type ChildrenDeletable[T] interface {
	    // Will also delete children if they exist.
	    DeletingChildrenIfNeeded() T
	}

	type Watchable[T] interface {
	    // Have the operation set a watch
	    Watched() T

	    // Set a watcher for the operation
	    UsingWatcher(watcher Watcher) T
	}

	type Backgroundable[T] interface {
	    // Perform the action in the background
	    InBackground() T

	    // Perform the action in the background
	    InBackgroundWithContext(context interface{}) T

	    // Perform the action in the background
	    InBackgroundWithCallback(callback BackgroundCallback) T

	    // Perform the action in the background
	    InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) T
	}
*/
package curator
