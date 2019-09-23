package curator

// Abstraction that provides the ZooKeeper connection string
type EnsembleProvider interface {
	// Curator will call this method when CuratorZookeeperClient.Start() is called
	Start() error

	// Curator will call this method when CuratorZookeeperClient.Close() is called
	Close() error

	// Return the current connection string to use
	ConnectionString() string
}

// Standard ensemble provider that wraps a fixed connection string
type FixedEnsembleProvider struct {
	connectString string // The connection string to use
}

func NewFixedEnsembleProvider(connectString string) *FixedEnsembleProvider {
	return &FixedEnsembleProvider{connectString}
}

func (p *FixedEnsembleProvider) Start() error { return nil }

func (p *FixedEnsembleProvider) Close() error { return nil }

func (p *FixedEnsembleProvider) ConnectionString() string { return p.connectString }
