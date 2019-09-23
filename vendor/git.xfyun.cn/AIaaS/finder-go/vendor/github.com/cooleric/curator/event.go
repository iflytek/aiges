package curator

import (
	"fmt"

	"github.com/cooleric/go-zookeeper/zk"
)

type CuratorEventType int

const (
	CREATE   CuratorEventType = iota // CuratorFramework.Create() -> Err(), Path(), Data()
	DELETE                           // CuratorFramework.Delete() -> Err(), Path()
	EXISTS                           // CuratorFramework.CheckExists() -> Err(), Path(), Stat()
	GET_DATA                         // CuratorFramework.GetData() -> Err(), Path(), Stat(), Data()
	SET_DATA                         // CuratorFramework.SetData() -> Err(), Path(), Stat()
	CHILDREN                         // CuratorFramework.GetChildren() -> Err(), Path(), Stat(), Children()
	SYNC                             // CuratorFramework.Sync() -> Err(), Path()
	GET_ACL                          // CuratorFramework.GetACL() -> Err(), Path()
	SET_ACL                          // CuratorFramework.SetACL() -> Err(), Path()
	WATCHED                          // Watchable.UsingWatcher() -> WatchedEvent()
	CLOSING                          // Event sent when client is being closed
)

var CuratorEventTypeNames = []string{"CREATE", "DELETE", "EXISTS", "GET_DATA", "SET_DATA", "CHILDREN", "SYNC", "GET_ACL", "SET_ACL", "WATCHED", "CLOSING"}

func (t CuratorEventType) String() string {
	if int(t) < len(CuratorEventTypeNames) {
		return CuratorEventTypeNames[int(t)]
	}

	return fmt.Sprintf("Type #%d", int(t))
}

// A super set of all the various Zookeeper events/background methods.
type CuratorEvent interface {
	// check here first - this value determines the type of event and which methods will have valid values
	Type() CuratorEventType

	// "rc" from async callbacks
	Err() error

	// the path
	Path() string

	// the context object passed to Backgroundable.InBackground(interface{})
	Context() interface{}

	// any stat
	Stat() *zk.Stat

	// any data
	Data() []byte

	// any name
	Name() string

	// any children
	Children() []string

	// any ACL list or null
	ACLs() []zk.ACL

	WatchedEvent() *zk.Event
}

type curatorEvent struct {
	eventType    CuratorEventType
	err          error
	path         string
	name         string
	children     []string
	context      interface{}
	stat         *zk.Stat
	data         []byte
	watchedEvent *zk.Event
	acls         []zk.ACL
}

func (e *curatorEvent) Type() CuratorEventType { return e.eventType }

func (e *curatorEvent) Err() error { return e.err }

func (e *curatorEvent) Path() string { return e.path }

func (e *curatorEvent) Context() interface{} { return e.context }

func (e *curatorEvent) Stat() *zk.Stat { return e.stat }

func (e *curatorEvent) Data() []byte { return e.data }

func (e *curatorEvent) Name() string { return e.name }

func (e *curatorEvent) Children() []string { return e.children }

func (e *curatorEvent) ACLs() []zk.ACL { return e.acls }

func (e *curatorEvent) WatchedEvent() *zk.Event { return e.watchedEvent }
