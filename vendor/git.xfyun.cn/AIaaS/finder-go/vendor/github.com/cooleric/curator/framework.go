package curator

import (
	"fmt"
	"log"
	"time"

	"github.com/cooleric/go-zookeeper/zk"
)

const (
	DEFAULT_SESSION_TIMEOUT    = 60 * time.Second
	DEFAULT_CONNECTION_TIMEOUT = 15 * time.Second
	DEFAULT_CLOSE_WAIT         = 1 * time.Second
)

// Zookeeper framework-style client
type CuratorFramework interface {
	// Start the client.
	// Most mutator methods will not work until the client is started
	Start() error

	// Stop the client
	Close() error

	// Returns the state of this instance
	State() State

	// Return true if the client is started, not closed, etc.
	Started() bool

	// Start a create builder
	Create() CreateBuilder

	// Start a delete builder
	Delete() DeleteBuilder

	// Start an exists builder
	CheckExists() CheckExistsBuilder

	// Start a get data builder
	GetData() GetDataBuilder

	// Start a set data builder
	SetData() SetDataBuilder

	// Start a get children builder
	GetChildren() GetChildrenBuilder

	// Start a get ACL builder
	GetACL() GetACLBuilder

	// Start a set ACL builder
	SetACL() SetACLBuilder

	// Start a transaction builder
	InTransaction() Transaction

	// Perform a sync on the given path - syncs are always in the background
	DoSync(path string, backgroundContextObject interface{})

	//  Start a sync builder. Note: sync is ALWAYS in the background even if you don't use one of the background() methods
	Sync() SyncBuilder

	// Returns the listenable interface for the Connect State
	ConnectionStateListenable() ConnectionStateListenable

	// Returns the listenable interface for events
	CuratorListenable() CuratorListenable

	// Returns the listenable interface for unhandled errors
	UnhandledErrorListenable() UnhandledErrorListenable

	// Returns a facade of the current instance that does _not_ automatically pre-pend the namespace to all paths
	NonNamespaceView() CuratorFramework

	// Returns a facade of the current instance that uses the specified namespace
	// or no namespace if newNamespace is empty.
	UsingNamespace(newNamespace string) CuratorFramework

	// Return the current namespace or "" if none
	Namespace() string

	// Return the managed zookeeper client
	ZookeeperClient() CuratorZookeeperClient

	// Allocates an ensure path instance that is namespace aware
	NewNamespaceAwareEnsurePath(path string) EnsurePath

	// Block until a connection to ZooKeeper is available.
	BlockUntilConnected() error

	// Block until a connection to ZooKeeper is available or the maxWaitTime has been exceeded
	BlockUntilConnectedTimeout(maxWaitTime time.Duration) error
}

// Create a new client with default session timeout and default connection timeout
func NewClient(connString string, retryPolicy RetryPolicy) CuratorFramework {
	return NewClientTimeout(connString, DEFAULT_SESSION_TIMEOUT, DEFAULT_CONNECTION_TIMEOUT, retryPolicy)
}

// Create a new client
func NewClientTimeout(connString string, sessionTimeout, connectionTimeout time.Duration, retryPolicy RetryPolicy) CuratorFramework {
	builder := &CuratorFrameworkBuilder{
		ConnectionTimeout: connectionTimeout,
		SessionTimeout:    sessionTimeout,
		RetryPolicy:       retryPolicy,
	}

	return builder.ConnectString(connString).Build()
}

type CuratorFrameworkBuilder struct {
	AuthInfos           []AuthInfo          // the connection authorization
	ZookeeperDialer     ZookeeperDialer     // the zookeeper dialer to use
	EnsembleProvider    EnsembleProvider    // the list ensemble provider.
	DefaultData         []byte              // the data to use when PathAndBytesable.ForPath(String) is used.
	Namespace           string              // as ZooKeeper is a shared space, users of a given cluster should stay within a pre-defined namespace
	SessionTimeout      time.Duration       // the session timeout
	ConnectionTimeout   time.Duration       // the connection timeout
	MaxCloseWait        time.Duration       // the time to wait during close to wait background tasks
	RetryPolicy         RetryPolicy         // the retry policy to use
	CompressionProvider CompressionProvider // the compression provider
	AclProvider         ACLProvider         // the provider for ACLs
	CanBeReadOnly       bool                // allow ZooKeeper client to enter read only mode in case of a network partition.
}

// Apply the current values and build a new CuratorFramework
func (b *CuratorFrameworkBuilder) Build() CuratorFramework {
	if b.EnsembleProvider == nil {
		panic("missed ensemble provider")
	}

	builder := *b

	if builder.SessionTimeout == 0 {
		builder.SessionTimeout = DEFAULT_SESSION_TIMEOUT
	}
	if builder.ConnectionTimeout == 0 {
		builder.ConnectionTimeout = DEFAULT_CONNECTION_TIMEOUT
	}
	if builder.MaxCloseWait == 0 {
		builder.MaxCloseWait = DEFAULT_CLOSE_WAIT
	}
	if builder.CompressionProvider == nil {
		builder.CompressionProvider = NewGzipCompressionProvider()
	}
	if builder.AclProvider == nil {
		builder.AclProvider = NewDefaultACLProvider()
	}

	return newCuratorFramework(&builder)
}

// Set the list of servers to connect to.
func (b *CuratorFrameworkBuilder) ConnectString(connectString string) *CuratorFrameworkBuilder {
	b.EnsembleProvider = &FixedEnsembleProvider{connectString}

	return b
}

// Add connection authorization
func (b *CuratorFrameworkBuilder) Authorization(scheme string, auth []byte) *CuratorFrameworkBuilder {
	b.AuthInfos = append(b.AuthInfos, AuthInfo{scheme, auth})

	return b
}

// Add compression provider
func (b *CuratorFrameworkBuilder) Compression(name string) *CuratorFrameworkBuilder {
	if provider, exists := CompressionProviders[name]; exists {
		b.CompressionProvider = provider
	}

	return b
}

type curatorFramework struct {
	client                  *curatorZookeeperClient
	stateManager            *connectionStateManager
	state                   State
	listeners               CuratorListenable
	unhandledErrorListeners UnhandledErrorListenable
	defaultData             []byte
	namespace               *namespaceImpl
	namespaceFacadeCache    *namespaceFacadeCache
	fixForNamespace         func(path string, isSequential bool) string
	unfixForNamespace       func(path string) string
	retryPolicy             RetryPolicy
	compressionProvider     CompressionProvider
	aclProvider             ACLProvider
}

func newCuratorFramework(b *CuratorFrameworkBuilder) *curatorFramework {
	c := &curatorFramework{
		listeners:               &curatorListenerContainer{},
		unhandledErrorListeners: &UnhandledErrorListenerContainer{},
		defaultData:             b.DefaultData,
		retryPolicy:             b.RetryPolicy,
		compressionProvider:     b.CompressionProvider,
		aclProvider:             b.AclProvider,
	}

	watcher := NewWatcher(func(event *zk.Event) {
		c.processEvent(&curatorEvent{
			eventType:    WATCHED,
			err:          event.Err,
			path:         c.unfixForNamespace(event.Path),
			watchedEvent: event,
		})
	})

	c.client = NewCuratorZookeeperClient(b.ZookeeperDialer, b.EnsembleProvider, b.SessionTimeout, b.ConnectionTimeout, watcher, b.RetryPolicy, b.CanBeReadOnly, b.AuthInfos)
	c.stateManager = newConnectionStateManager(c)
	c.namespace = newNamespace(c, b.Namespace)
	c.namespaceFacadeCache = newNamespaceFacadeCache(c)
	c.fixForNamespace = c.namespace.fixForNamespace
	c.unfixForNamespace = c.namespace.unfixForNamespace

	return c
}

func (c *curatorFramework) Start() error {
	if !c.state.Change(LATENT, STARTED) {
		return fmt.Errorf("Cannot be started more than once")
	} else if err := c.stateManager.Start(); err != nil {
		return fmt.Errorf("fail to start state manager, %s", err)
	} else if err := c.client.Start(); err != nil {
		return fmt.Errorf("fail to start client, %s", err)
	}

	return nil
}

func (c *curatorFramework) Close() error {
	if !c.state.Change(STARTED, STOPPED) {
		return nil
	}

	evt := &curatorEvent{eventType: CLOSING}

	c.listeners.ForEach(func(listener interface{}) {
		listener.(CuratorListener).EventReceived(c, evt)
	})

	c.listeners.Clear()
	c.unhandledErrorListeners.Clear()
	c.stateManager.Close()

	return c.client.Close()
}

func (c *curatorFramework) State() State {
	return c.state.Value()
}

func (c *curatorFramework) Started() bool {
	return c.State() == STARTED
}

func (c *curatorFramework) Create() CreateBuilder {
	c.state.Check(STARTED, "instance must be started before calling Create")

	return &createBuilder{client: c, acling: acling{aclProvider: c.aclProvider}}
}

func (c *curatorFramework) Delete() DeleteBuilder {
	c.state.Check(STARTED, "instance must be started before calling Delete")

	return &deleteBuilder{client: c, version: AnyVersion}
}

func (c *curatorFramework) CheckExists() CheckExistsBuilder {
	c.state.Check(STARTED, "instance must be started before calling CheckExists")

	return &checkExistsBuilder{client: c}
}

func (c *curatorFramework) GetData() GetDataBuilder {
	c.state.Check(STARTED, "instance must be started before calling GetData")

	return &getDataBuilder{client: c}
}

func (c *curatorFramework) SetData() SetDataBuilder {
	c.state.Check(STARTED, "instance must be started before calling SetData")

	return &setDataBuilder{client: c, version: AnyVersion}
}

func (c *curatorFramework) GetChildren() GetChildrenBuilder {
	c.state.Check(STARTED, "instance must be started before calling GetChildren")

	return &getChildrenBuilder{client: c}
}

func (c *curatorFramework) GetACL() GetACLBuilder {
	c.state.Check(STARTED, "instance must be started before calling GetACL")

	return &getACLBuilder{client: c}
}

func (c *curatorFramework) SetACL() SetACLBuilder {
	c.state.Check(STARTED, "instance must be started before calling SetACL")

	return &setACLBuilder{client: c, version: AnyVersion, acling: acling{aclProvider: c.aclProvider}}
}

func (c *curatorFramework) InTransaction() Transaction {
	c.state.Check(STARTED, "instance must be started before calling InTransaction")

	return &curatorTransaction{client: c}
}

func (c *curatorFramework) DoSync(path string, context interface{}) {
	c.Sync().InBackgroundWithContext(context).ForPath(path)
}

func (c *curatorFramework) Sync() SyncBuilder {
	c.state.Check(STARTED, "instance must be started before calling this method")

	return &syncBuilder{client: c}
}

func (c *curatorFramework) ConnectionStateListenable() ConnectionStateListenable {
	return c.stateManager.Listenable()
}

func (c *curatorFramework) CuratorListenable() CuratorListenable {
	return c.listeners
}

func (c *curatorFramework) UnhandledErrorListenable() UnhandledErrorListenable {
	return c.unhandledErrorListeners
}

func (c *curatorFramework) processEvent(event CuratorEvent) {
	if event.Type() == WATCHED {
		c.validateConnection(event.WatchedEvent().State)
	}

	c.listeners.ForEach(func(l interface{}) {
		tracer := c.client.StartTracer("EventListener")

		defer tracer.Commit()

		if err := l.(CuratorListener).EventReceived(c, event); err != nil {
			c.logError(fmt.Errorf("Event listener threw exception, %s", err))
		}
	})
}

func (c *curatorFramework) validateConnection(state zk.State) {
	switch state {
	case zk.StateDisconnected:
		c.suspendConnection()

	case zk.StateExpired:
		c.stateManager.AddStateChange(LOST)

	case zk.StateConnected:
		c.stateManager.AddStateChange(RECONNECTED)

	case zk.StateConnectedReadOnly:
		c.stateManager.AddStateChange(READ_ONLY)
	}
}

func (c *curatorFramework) suspendConnection() {
	if !c.stateManager.SetToSuspended() {
		return
	}

	go c.doSyncForSuspendedConnection(c.client.InstanceIndex())
}

func (c *curatorFramework) doSyncForSuspendedConnection(instanceIndex int64) {
	tracer := c.client.StartTracer("BackgroundSyncImpl")

	defer tracer.Commit()

	if conn, err := c.client.Conn(); err == nil {
		if path, err := conn.Sync("/"); err == nil {
			c.processEvent(&curatorEvent{
				eventType: SYNC,
				err:       err,
				path:      path,
			})

			return
		}
	}

	if instanceIndex < 0 || instanceIndex == c.client.InstanceIndex() {
		c.stateManager.AddStateChange(LOST)
	} else {
		go c.doSyncForSuspendedConnection(-1)
	}
}

func (c *curatorFramework) logError(err error) {
	log.Printf("error: %s", err)

	c.unhandledErrorListeners.ForEach(func(listener interface{}) {
		listener.(UnhandledErrorListener).UnhandledError(err)
	})
}

func (c *curatorFramework) NonNamespaceView() CuratorFramework {
	return c.UsingNamespace("")
}

func (c *curatorFramework) UsingNamespace(newNamespace string) CuratorFramework {
	c.state.Check(STARTED, "instance must be started before calling this method")

	return c.namespaceFacadeCache.Get(newNamespace)
}

func (c *curatorFramework) Namespace() string {
	return c.namespace.namespace
}

func (c *curatorFramework) getNamespaceWatcher(watcher Watcher) Watcher {
	return watcher
}

func (c *curatorFramework) ZookeeperClient() CuratorZookeeperClient {
	return c.client
}

func (c *curatorFramework) NewNamespaceAwareEnsurePath(path string) EnsurePath {
	return NewEnsurePathWithAcl(c.fixForNamespace(path, false), c.aclProvider)
}

func (c *curatorFramework) BlockUntilConnected() error {
	return c.BlockUntilConnectedTimeout(0)
}

func (c *curatorFramework) BlockUntilConnectedTimeout(maxWaitTime time.Duration) error {
	return c.stateManager.BlockUntilConnected(maxWaitTime)
}
