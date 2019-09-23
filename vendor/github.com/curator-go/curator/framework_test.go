package curator

import (
	"fmt"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
)

func StartNewTestingClient(t *testing.T) CuratorFramework {
	zkCluster, err := zk.StartTestCluster(1, nil, nil)
	assert.NoError(t, err)
	var c = NewClient(fmt.Sprintf("127.0.0.1:%d", zkCluster.Servers[0].Port), nil)
	c.ConnectionStateListenable().AddListener(NewConnectionStateListener(
		func(client CuratorFramework, newState ConnectionState) {
			t.Log("New State: ", newState)
		}))
	assert.NoError(t, c.Start())
	return c
}

func TestConnectionStateConnected(t *testing.T) {
	var c = StartNewTestingClient(t)
	assert.NoError(t, c.ZookeeperClient().BlockUntilConnectedOrTimedOut())
	assert.NoError(t, c.BlockUntilConnected())
}
