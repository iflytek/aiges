package curator

import (
	"sync"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeleteBuilderTestSuite struct {
	mockContainerTestSuite
}

func TestDeleteBuilder(t *testing.T) {
	suite.Run(t, new(DeleteBuilderTestSuite))
}

func (s *DeleteBuilderTestSuite) TestDelete() {
	s.With(func(client CuratorFramework, conn *mockConn) {
		conn.On("Delete", "/node", AnyVersion).Return(nil).Once()

		assert.NoError(s.T(), client.Delete().ForPath("/node"))
	})
}

func (s *DeleteBuilderTestSuite) TestDeleteWithVersion() {
	s.With(func(client CuratorFramework, conn *mockConn, version int32) {
		conn.On("Delete", "/node", version).Return(zk.ErrBadVersion).Once()

		assert.EqualError(s.T(), client.Delete().WithVersion(version).ForPath("/node"), zk.ErrBadVersion.Error())
	})
}

func (s *DeleteBuilderTestSuite) TestNamespace() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn) {
		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Delete", "/parent/child", AnyVersion).Return(nil).Once()

		assert.NoError(s.T(), client.Delete().ForPath("/child"))
	})
}

func (s *DeleteBuilderTestSuite) TestBackground() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup) {
		ctxt := "context"

		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Delete", "/parent/child", AnyVersion).Return(nil).Once()

		assert.NoError(s.T(), client.Delete().InBackgroundWithCallbackAndContext(
			func(client CuratorFramework, event CuratorEvent) error {
				defer wg.Done()

				assert.Equal(s.T(), DELETE, event.Type())
				assert.Equal(s.T(), "/child", event.Path())
				assert.NoError(s.T(), event.Err())
				assert.Equal(s.T(), "child", event.Name())
				assert.Equal(s.T(), ctxt, event.Context())

				return nil
			}, ctxt).ForPath("/child"))
	})
}

func (s *DeleteBuilderTestSuite) TestDeletingChildren() {
	s.With(func(client CuratorFramework, conn *mockConn) {
		conn.On("Delete", "/parent", AnyVersion).Return(zk.ErrNotEmpty).Once()
		conn.On("Children", "/parent").Return([]string{"child"}, nil, nil).Once()
		conn.On("Children", "/parent/child").Return([]string{}, nil, nil).Once()
		conn.On("Delete", "/parent/child", AnyVersion).Return(nil).Once()
		conn.On("Delete", "/parent", AnyVersion).Return(nil).Once()

		assert.NoError(s.T(), client.Delete().DeletingChildrenIfNeeded().ForPath("/parent"))
	})
}
