package curator

import (
	"github.com/cooleric/go-zookeeper/zk"
)

type CreateBuilder interface {
	// PathAndBytesable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) (string, error)

	// Commit the currently building operation using the given path and data
	ForPathWithData(path string, payload []byte) (string, error)

	// ParentsCreatable[T]
	//
	// Causes any parent nodes to get created if they haven't already been
	CreatingParentsIfNeeded() CreateBuilder

	// CreateModable[T]
	//
	// Set a create mode - the default is CreateMode.PERSISTENT
	WithMode(mode CreateMode) CreateBuilder

	// ACLable[T]
	//
	// Set an ACL list
	WithACL(acls ...zk.ACL) CreateBuilder

	// Compressible[T]
	//
	// Cause the data to be compressed using the configured compression provider
	Compressed() CreateBuilder

	// Backgroundable[T]
	//
	// Perform the action in the background
	InBackground() CreateBuilder

	// Perform the action in the background
	InBackgroundWithContext(context interface{}) CreateBuilder

	// Perform the action in the background
	InBackgroundWithCallback(callback BackgroundCallback) CreateBuilder

	// Perform the action in the background
	InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) CreateBuilder
}

type CheckExistsBuilder interface {
	// Pathable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) (*zk.Stat, error)

	// Watchable[T]
	//
	// Have the operation set a watch
	Watched() CheckExistsBuilder

	// Set a watcher for the operation
	UsingWatcher(watcher Watcher) CheckExistsBuilder

	// Backgroundable[T]
	//
	// Perform the action in the background
	InBackground() CheckExistsBuilder

	// Perform the action in the background
	InBackgroundWithContext(context interface{}) CheckExistsBuilder

	// Perform the action in the background
	InBackgroundWithCallback(callback BackgroundCallback) CheckExistsBuilder

	// Perform the action in the background
	InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) CheckExistsBuilder
}

type DeleteBuilder interface {
	// Pathable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) error

	// ChildrenDeletable[T]
	//
	// Will also delete children if they exist.
	DeletingChildrenIfNeeded() DeleteBuilder

	// Versionable[T]
	//
	// Use the given version (the default is -1)
	WithVersion(version int32) DeleteBuilder

	// Backgroundable[T]
	//
	// Perform the action in the background
	InBackground() DeleteBuilder

	// Perform the action in the background
	InBackgroundWithContext(context interface{}) DeleteBuilder

	// Perform the action in the background
	InBackgroundWithCallback(callback BackgroundCallback) DeleteBuilder

	// Perform the action in the background
	InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) DeleteBuilder
}

type GetDataBuilder interface {
	// Pathable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) ([]byte, error)

	// Decompressible[T]
	//
	// Cause the data to be de-compressed using the configured compression provider
	Decompressed() GetDataBuilder

	// Statable[T]
	//
	// Have the operation fill the provided stat object
	StoringStatIn(stat *zk.Stat) GetDataBuilder

	// Watchable[T]
	//
	// Have the operation set a watch
	Watched() GetDataBuilder

	// Set a watcher for the operation
	UsingWatcher(watcher Watcher) GetDataBuilder

	// Backgroundable[T]
	//
	// Perform the action in the background
	InBackground() GetDataBuilder

	// Perform the action in the background
	InBackgroundWithContext(context interface{}) GetDataBuilder

	// Perform the action in the background
	InBackgroundWithCallback(callback BackgroundCallback) GetDataBuilder

	// Perform the action in the background
	InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) GetDataBuilder
}

type SetDataBuilder interface {
	// PathAndBytesable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) (*zk.Stat, error)

	// Commit the currently building operation using the given path and data
	ForPathWithData(path string, payload []byte) (*zk.Stat, error)

	// Versionable[T]
	//
	// Use the given version (the default is -1)
	WithVersion(version int32) SetDataBuilder

	// Compressible[T]
	//
	// Cause the data to be compressed using the configured compression provider
	Compressed() SetDataBuilder

	// Backgroundable[T]
	//
	// Perform the action in the background
	InBackground() SetDataBuilder

	// Perform the action in the background
	InBackgroundWithContext(context interface{}) SetDataBuilder

	// Perform the action in the background
	InBackgroundWithCallback(callback BackgroundCallback) SetDataBuilder

	// Perform the action in the background
	InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) SetDataBuilder
}

type GetChildrenBuilder interface {
	// Pathable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) ([]string, error)

	// Statable[T]
	//
	// Have the operation fill the provided stat object
	StoringStatIn(stat *zk.Stat) GetChildrenBuilder

	// Watchable[T]
	//
	// Have the operation set a watch
	Watched() GetChildrenBuilder

	// Set a watcher for the operation
	UsingWatcher(watcher Watcher) GetChildrenBuilder

	// Backgroundable[T]
	//
	// Perform the action in the background
	InBackground() GetChildrenBuilder

	// Perform the action in the background
	InBackgroundWithContext(context interface{}) GetChildrenBuilder

	// Perform the action in the background
	InBackgroundWithCallback(callback BackgroundCallback) GetChildrenBuilder

	// Perform the action in the background
	InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) GetChildrenBuilder
}

type GetACLBuilder interface {
	// Pathable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) ([]zk.ACL, error)

	// Statable[T]
	//
	// Have the operation fill the provided stat object
	StoringStatIn(stat *zk.Stat) GetACLBuilder

	// Backgroundable[T]
	//
	// Perform the action in the background
	InBackground() GetACLBuilder

	// Perform the action in the background
	InBackgroundWithContext(context interface{}) GetACLBuilder

	// Perform the action in the background
	InBackgroundWithCallback(callback BackgroundCallback) GetACLBuilder

	// Perform the action in the background
	InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) GetACLBuilder
}

type SetACLBuilder interface {
	// Pathable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) (*zk.Stat, error)

	// ACLable[T]
	//
	// Set an ACL list
	WithACL(acls ...zk.ACL) SetACLBuilder

	// Versionable[T]
	//
	// Use the given version (the default is -1)
	WithVersion(version int32) SetACLBuilder

	// Backgroundable[T]
	//
	// Perform the action in the background
	InBackground() SetACLBuilder

	// Perform the action in the background
	InBackgroundWithContext(context interface{}) SetACLBuilder

	// Perform the action in the background
	InBackgroundWithCallback(callback BackgroundCallback) SetACLBuilder

	// Perform the action in the background
	InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) SetACLBuilder
}

type SyncBuilder interface {
	// Pathable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) (string, error)

	// Backgroundable[T]
	//
	// Perform the action in the background
	InBackground() SyncBuilder

	// Perform the action in the background
	InBackgroundWithContext(context interface{}) SyncBuilder

	// Perform the action in the background
	InBackgroundWithCallback(callback BackgroundCallback) SyncBuilder

	// Perform the action in the background
	InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) SyncBuilder
}

type TransactionCreateBuilder interface {
	// PathAndBytesable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) TransactionBridge

	// Commit the currently building operation using the given path and data
	ForPathWithData(path string, payload []byte) TransactionBridge

	// CreateModable[T]
	//
	// Set a create mode - the default is CreateMode.PERSISTENT
	WithMode(mode CreateMode) TransactionCreateBuilder

	// ACLable[T]
	//
	// Set an ACL list
	WithACL(acls ...zk.ACL) TransactionCreateBuilder

	// Compressible[T]
	//
	// Cause the data to be compressed using the configured compression provider
	Compressed() TransactionCreateBuilder
}

type TransactionDeleteBuilder interface {
	// Pathable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) TransactionBridge

	// Versionable[T]
	//
	// Use the given version (the default is -1)
	WithVersion(version int32) TransactionDeleteBuilder
}

type TransactionSetDataBuilder interface {
	// PathAndBytesable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) TransactionBridge

	// Commit the currently building operation using the given path and data
	ForPathWithData(path string, payload []byte) TransactionBridge

	// Versionable[T]
	//
	// Use the given version (the default is -1)
	WithVersion(version int32) TransactionSetDataBuilder

	// Compressible[T]
	//
	// Cause the data to be compressed using the configured compression provider
	Compressed() TransactionSetDataBuilder
}

type TransactionCheckBuilder interface {
	// Pathable[T]
	//
	// Commit the currently building operation using the given path
	ForPath(path string) TransactionBridge

	// Versionable[T]
	//
	// Use the given version (the default is -1)
	WithVersion(version int32) TransactionCheckBuilder
}
