package curator

import (
	"github.com/cooleric/go-zookeeper/zk"
)

var (
	ErrConnectionClosed        = zk.ErrConnectionClosed
	ErrUnknown                 = zk.ErrUnknown
	ErrAPIError                = zk.ErrAPIError
	ErrNoNode                  = zk.ErrNoNode
	ErrNoAuth                  = zk.ErrNoAuth
	ErrBadVersion              = zk.ErrBadVersion
	ErrNoChildrenForEphemerals = zk.ErrNoChildrenForEphemerals
	ErrNodeExists              = zk.ErrNodeExists
	ErrNotEmpty                = zk.ErrNotEmpty
	ErrSessionExpired          = zk.ErrSessionExpired
	ErrInvalidACL              = zk.ErrInvalidACL
	ErrAuthFailed              = zk.ErrAuthFailed
	ErrClosing                 = zk.ErrClosing
	ErrNothing                 = zk.ErrNothing
	ErrSessionMoved            = zk.ErrSessionMoved
)

var (
	EventNodeCreated         = zk.EventNodeCreated
	EventNodeDeleted         = zk.EventNodeDeleted
	EventNodeDataChanged     = zk.EventNodeDataChanged
	EventNodeChildrenChanged = zk.EventNodeChildrenChanged
)

const AnyVersion int32 = -1

type CreateMode int32

const (
	PERSISTENT            CreateMode = 0
	PERSISTENT_SEQUENTIAL            = zk.FlagSequence
	EPHEMERAL                        = zk.FlagEphemeral
	EPHEMERAL_SEQUENTIAL             = zk.FlagEphemeral + zk.FlagSequence
)

func (m CreateMode) IsSequential() bool { return (m & zk.FlagSequence) == zk.FlagSequence }
func (m CreateMode) IsEphemeral() bool  { return (m & zk.FlagEphemeral) == zk.FlagEphemeral }

// Called when the async background operation completes
type BackgroundCallback func(client CuratorFramework, event CuratorEvent) error

type backgrounding struct {
	inBackground bool
	context      interface{}
	callback     BackgroundCallback
}

type watching struct {
	watcher Watcher
	watched bool
}
