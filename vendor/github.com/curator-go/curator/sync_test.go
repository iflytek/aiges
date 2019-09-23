package curator

import (
	"sync"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SyncBuilderTestSuite struct {
	mockContainerTestSuite
}

func TestSyncBuilder(t *testing.T) {
	suite.Run(t, new(SyncBuilderTestSuite))
}

func (s *SyncBuilderTestSuite) TestSync() {
	s.With(func(client CuratorFramework, conn *mockConn, stat *zk.Stat) {
		conn.On("Sync", "/node").Return("/node", nil).Once()

		path, err := client.Sync().ForPath("/node")

		assert.Equal(s.T(), "/node", path)
		assert.NoError(s.T(), err)
	})
}

func (s *SyncBuilderTestSuite) TestNamespace() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn) {
		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Sync", "/parent/child").Return("", zk.ErrNoNode).Once()

		path, err := client.Sync().ForPath("/child")

		assert.Equal(s.T(), "/child", path)
		assert.EqualError(s.T(), err, zk.ErrNoNode.Error())
	})
}

func (s *SyncBuilderTestSuite) TestBackground() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup) {
		ctxt := "context"

		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Sync", "/parent/child").Return("/parent/child", nil).Once()

		path, err := client.Sync().InBackgroundWithCallbackAndContext(
			func(client CuratorFramework, event CuratorEvent) error {
				defer wg.Done()

				assert.Equal(s.T(), SYNC, event.Type())
				assert.Equal(s.T(), "/child", event.Path())
				assert.NoError(s.T(), event.Err())
				assert.Equal(s.T(), "child", event.Name())
				assert.Equal(s.T(), ctxt, event.Context())

				return nil
			}, ctxt).ForPath("/child")

		assert.Equal(s.T(), "/child", path)
		assert.NoError(s.T(), err)
	})
}
