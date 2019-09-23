package curator

import (
	"github.com/samuel/go-zookeeper/zk"
)

// Transactional/atomic operations.
//
// The general form for this interface is:
//
// 		curator.InTransaction().operation().arguments().ForPath(...).
//              And().more-operations.
//              And().Commit()
//
// Here's an example that creates two nodes in a transaction
//
//		curator.InTransaction().
//              Create().ForPathWithData("/path-one", path-one-data).
//              And().Create().ForPathWithData("/path-two", path-two-data).
//              And().Commit()
//
// <b>Important:</b> the operations are not submitted until CuratorTransactionFinal.Commit() is called.
//
type Transaction interface {
	// Start a create builder in the transaction
	Create() TransactionCreateBuilder

	// Start a delete builder in the transaction
	Delete() TransactionDeleteBuilder

	// Start a set data builder in the transaction
	SetData() TransactionSetDataBuilder

	// Start a check builder in the transaction
	Check() TransactionCheckBuilder
}

// Transaction operation types
type OperationType int

const (
	OP_CREATE OperationType = iota
	OP_DELETE
	OP_SET_DATA
	OP_CHECK
)

// Holds the result of one transactional operation
type TransactionResult struct {
	Type       OperationType
	ForPath    string
	ResultPath string
	ResultStat *zk.Stat
}

// Adds commit to the transaction interface
type TransactionFinal interface {
	Transaction

	// Commit all added operations as an atomic unit and return results for the operations.
	// One result is returned for each operation added.
	// Further, the ordering of the results matches the ordering that the operations were added.
	Commit() ([]TransactionResult, error)
}

// Syntactic sugar to make the fluent interface more readable
type TransactionBridge interface {
	TransactionFinal

	And() TransactionFinal
}

type curatorTransaction struct {
	client     *curatorFramework
	operations []interface{}
}

func (t *curatorTransaction) Create() TransactionCreateBuilder {
	return &transactionCreateBuilder{transaction: t, acling: acling{aclProvider: t.client.aclProvider}}
}

func (t *curatorTransaction) Delete() TransactionDeleteBuilder {
	return &transactionDeleteBuilder{transaction: t, version: AnyVersion}
}

func (t *curatorTransaction) SetData() TransactionSetDataBuilder {
	return &transactionSetDataBuilder{transaction: t, version: AnyVersion}
}

func (t *curatorTransaction) Check() TransactionCheckBuilder {
	return &transactionCheckBuilder{transaction: t, version: AnyVersion}
}

func (t *curatorTransaction) And() TransactionFinal {
	return t
}

func (t *curatorTransaction) Commit() ([]TransactionResult, error) {
	zkClient := t.client.ZookeeperClient()

	result, err := zkClient.NewRetryLoop().CallWithRetry(func() (interface{}, error) {
		if conn, err := zkClient.Conn(); err != nil {
			return nil, err
		} else {
			return conn.Multi(t.operations...)
		}
	})

	var results []TransactionResult

	if responses, ok := result.([]zk.MultiResponse); ok {
		for i, res := range responses {
			switch req := t.operations[i].(type) {
			case *zk.CreateRequest:
				results = append(results, TransactionResult{
					Type:       OP_CREATE,
					ForPath:    req.Path,
					ResultPath: t.client.unfixForNamespace(res.String),
				})
			case *zk.DeleteRequest:
				results = append(results, TransactionResult{
					Type:    OP_DELETE,
					ForPath: req.Path,
				})
			case *zk.SetDataRequest:
				results = append(results, TransactionResult{
					Type:       OP_SET_DATA,
					ForPath:    req.Path,
					ResultStat: res.Stat,
				})
			case *zk.CheckVersionRequest:
				results = append(results, TransactionResult{
					Type:    OP_CHECK,
					ForPath: req.Path,
				})
			}
		}
	}

	return results, err
}

type transactionCreateBuilder struct {
	transaction *curatorTransaction
	createMode  CreateMode
	compress    bool
	acling      acling
}

func (b *transactionCreateBuilder) ForPath(path string) TransactionBridge {
	return b.ForPathWithData(path, b.transaction.client.defaultData)
}

func (b *transactionCreateBuilder) ForPathWithData(path string, payload []byte) TransactionBridge {
	var data []byte

	if b.compress {
		data, _ = b.transaction.client.compressionProvider.Compress(path, payload)
	} else {
		data = payload
	}

	b.transaction.operations = append(b.transaction.operations, &zk.CreateRequest{
		Path:  b.transaction.client.fixForNamespace(path, false),
		Data:  data,
		Acl:   b.acling.getAclList(path),
		Flags: int32(b.createMode),
	})

	return b.transaction
}

func (b *transactionCreateBuilder) WithMode(mode CreateMode) TransactionCreateBuilder {
	b.createMode = mode

	return b
}

func (b *transactionCreateBuilder) WithACL(acls ...zk.ACL) TransactionCreateBuilder {
	b.acling.aclList = acls

	return b
}

func (b *transactionCreateBuilder) Compressed() TransactionCreateBuilder {
	b.compress = true

	return b
}

type transactionDeleteBuilder struct {
	transaction *curatorTransaction
	version     int32
}

func (b *transactionDeleteBuilder) ForPath(path string) TransactionBridge {
	b.transaction.operations = append(b.transaction.operations, &zk.DeleteRequest{
		Path:    b.transaction.client.fixForNamespace(path, false),
		Version: b.version,
	})

	return b.transaction
}

func (b *transactionDeleteBuilder) WithVersion(version int32) TransactionDeleteBuilder {
	b.version = version

	return b
}

type transactionSetDataBuilder struct {
	transaction *curatorTransaction
	version     int32
	compress    bool
}

func (b *transactionSetDataBuilder) ForPath(path string) TransactionBridge {
	return b.ForPathWithData(path, b.transaction.client.defaultData)
}

func (b *transactionSetDataBuilder) ForPathWithData(path string, payload []byte) TransactionBridge {
	var data []byte

	if b.compress {
		data, _ = b.transaction.client.compressionProvider.Compress(path, payload)
	} else {
		data = payload
	}

	b.transaction.operations = append(b.transaction.operations, &zk.SetDataRequest{
		Path:    b.transaction.client.fixForNamespace(path, false),
		Data:    data,
		Version: b.version,
	})

	return b.transaction
}

func (b *transactionSetDataBuilder) WithVersion(version int32) TransactionSetDataBuilder {
	b.version = version

	return b
}

func (b *transactionSetDataBuilder) Compressed() TransactionSetDataBuilder {
	b.compress = true

	return b
}

type transactionCheckBuilder struct {
	transaction *curatorTransaction
	version     int32
}

func (b *transactionCheckBuilder) ForPath(path string) TransactionBridge {
	b.transaction.operations = append(b.transaction.operations, &zk.CheckVersionRequest{
		Path:    b.transaction.client.fixForNamespace(path, false),
		Version: b.version,
	})

	return b.transaction
}

func (b *transactionCheckBuilder) WithVersion(version int32) TransactionCheckBuilder {
	b.version = version

	return b
}
