# Framework

The Curator Framework is a high-level API that greatly simplifies using ZooKeeper. It adds many features that build on ZooKeeper and handles the complexity of managing connections to the ZooKeeper cluster and retrying operations. Some of the features are:

- Automatic connection management:
  * There are potential error cases that require ZooKeeper clients to recreate a connection and/or retry operations. Curator automatically and transparently (mostly) handles these cases.
- Cleaner API:
  * simplifies the raw ZooKeeper methods, events, etc.
  * provides a modern, fluent interface
- Recipe implementations (see Recipes):
  * Leader election
  * Shared lock
  * Path cache and watcher
  * Distributed Queue
  * Distributed Priority Queue
  * ...
  
## Allocating a Curator Framework Instance

CuratorFrameworks are allocated using the CuratorFrameworkBuilder which provides both factory methods and a builder for creating instances. IMPORTANT: CuratorFramework instances are fully thread-safe. You should share one CuratorFramework per ZooKeeper cluster in your application.

The factory methods [NewClient()](http://godoc.org/github.com/curator-go/curator#NewClient) provide a simplified way of creating an instance. The Builder gives control over all parameters. Once you have a CuratorFramework instance, you must call the [Start()](http://godoc.org/github.com/curator-go/curator#CuratorFramework.Start) method. At the end of your application, you should call [Close()](http://godoc.org/github.com/curator-go/curator#CuratorFramework.Close).

## CuratorFramework API

The CuratorFramework uses a Fluent-style interface. Operations are constructed using builders returned by the CuratorFramework instance. When strung together, the methods form sentence-like statements. e.g.

```
client.Create().ForPathWithData("/head", []byte{});
client.Delete().InBackground().ForPath("/head");
client.Create().WithMode(curator.EPHEMERAL_SEQUENTIAL).ForPath("/head/child", []byte{});
client.GetData().Watched().InBackground().ForPath("/test");
```

## Methods

- [Create()](http://godoc.org/github.com/curator-go/curator#CuratorFramework.Create)	Begins a create operation. Call additional methods (mode or background) and finalize the operation by calling ForPath() or ForPathWithData()
- [Delete()](http://godoc.org/github.com/curator-go/curator#CuratorFramework.Delete)	Begins a delete operation. Call additional methods (version or background) and finalize the operation by calling ForPath()
- [CheckExists()](http://godoc.org/github.com/curator-go/curator#CuratorFramework.CheckExists)	Begins an operation to check that a ZNode exists. Call additional methods (watch or background) and finalize the operation by calling ForPath()
- [GetData()](http://godoc.org/github.com/curator-go/curator#CuratorFramework.GetData)	Begins an operation to get a ZNode's data. Call additional methods (watch, background or get stat) and finalize the operation by calling ForPath()
- [SetData()](http://godoc.org/github.com/curator-go/curator#CuratorFramework.SetData)	Begins an operation to set a ZNode's data. Call additional methods (version or background) and finalize the operation by calling ForPath() or ForPathWithData()
- [GetChildren()](http://godoc.org/github.com/curator-go/curator#CuratorFramework.GetChildren)	Begins an operation to get a ZNode's list of children ZNodes. Call additional methods (watch, background or get stat) and finalize the operation by calling ForPath()
- [InTransaction()](http://godoc.org/github.com/curator-go/curator#CuratorFramework.InTransaction)	Begins an atomic ZooKeeper transaction. Combine Create, SetData, Check, and/or Delete operations and then Commit() as a unit.

## Notifications

Notifications for background operations and watches are published via the ClientListener interface. You register listeners with the CuratorFramework instance using the addListener() method. The listener implements two methods:

- [EventReceived()](http://godoc.org/github.com/curator-go/curator/#CuratorListener)	A background operation has completed or a watch has triggered. Examine the given event for details

## CuratorEvent

The CuratorEvent object is a super-set POJO that can hold every type of background notification and triggered watch. The useful fields of ClientEvent depend on the type of event which is exposed via the Type() method.

- **CREATE**	    Err(), Path(), Data()
- **DELETE**	    Err(), Path()
- **EXISTS**	    Err(), Path(), Stat()
- **GETDATA**    Err(), Path(), Stat(), Data()
- **SETDATA**    Err(), Path(), Stat()
- **CHILDREN**   Err(), Path(), Stat(), Children()
- **WATCHED**	 WatchedEvent()

## Namespaces

Because a ZooKeeper cluster is a shared environment, it's vital that a namespace convention is observed so that various applications that use a given cluster don't use conflicting ZK paths.

The CuratorFramework has a concept of a "namespace". You set the namespace when creating a CuratorFramework instance (via the CuratorFrameworkBuilder). The CuratorFramework will then prepend the namespace to all paths when one of its APIs is called. i.e.

```
client := &CuratorFrameworkBuilder{Namespace: "MyApp"} ... Build()
 ...
client.Create().ForPathWithData("/test", data)
// node was actually written to: "/MyApp/test"
```