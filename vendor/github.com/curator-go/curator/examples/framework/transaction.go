package framework

import (
	"fmt"

	"github.com/curator-go/curator"
)

func CreateTransaction(client curator.CuratorFramework) error {
	// this example shows how to use ZooKeeper's new transactions
	results, err := client.InTransaction().Create().ForPathWithData("/a/path", []byte("some data")).
		And().SetData().ForPathWithData("/another/path", []byte("other data")).
		And().Delete().ForPath("/yet/another/path").
		Commit() // IMPORTANT! The transaction is not submitted until commit() is called

	for _, result := range results {
		fmt.Printf("%s - %s", result.ForPath, result.Type)
	}

	return err
}

// These next four methods show how to use Curator's transaction APIs in a more traditional - one-at-a-time - manner

func StartTransaction(client curator.CuratorFramework) curator.Transaction {
	// start the transaction builder
	return client.InTransaction()
}

func AddCreateToTransaction(transaction curator.Transaction) curator.TransactionFinal {
	// add a create operation
	return transaction.Create().ForPathWithData("/a/path", []byte("some data")).And()
}

func AddDeleteToTransaction(transaction curator.Transaction) curator.TransactionFinal {
	// add a create operation
	return transaction.Delete().ForPath("/yet/another/path").And()
}

func CommitTransaction(transaction curator.TransactionFinal) ([]curator.TransactionResult, error) {
	// commit the transaction
	return transaction.Commit()
}
