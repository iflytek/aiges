package framework

import (
	"time"

	"github.com/curator-go/curator"
)

func CreateSimple(connString string) curator.CuratorFramework {
	// these are reasonable arguments for the ExponentialBackoffRetry.
	// the first retry will wait 1 second,
	// the second will wait up to 2 seconds,
	// the third will wait up to 4 seconds.
	retryPolicy := curator.NewExponentialBackoffRetry(time.Second, 3, 15*time.Second)

	// The simplest way to get a CuratorFramework instance. This will use default values.
	// The only required arguments are the connection string and the retry policy
	return curator.NewClient(connString, retryPolicy)
}

func CreateWithOptions(connString string, retryPolicy curator.RetryPolicy, connectionTimeout, sessionTimeout time.Duration) curator.CuratorFramework {
	// using the CuratorFrameworkBuilder gives fine grained control over creation options.
	builder := &curator.CuratorFrameworkBuilder{
		ConnectionTimeout: connectionTimeout,
		SessionTimeout:    sessionTimeout,
		RetryPolicy:       retryPolicy,
	}

	return builder.ConnectString(connString).Authorization("digest", []byte("user:pass")).Build()
}
