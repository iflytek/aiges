package curator

import (
	"sync"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetChildrenBuilderTestSuite struct {
	mockContainerTestSuite
}

func TestGetChildrenBuilder(t *testing.T) {
	suite.Run(t, new(GetChildrenBuilderTestSuite))
}

func (s *GetChildrenBuilderTestSuite) TestGetChildren() {
	s.With(func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, stat *zk.Stat) {
		conn.On("Children", "/parent").Return([]string{"child"}, stat, nil).Once()

		var parentStat zk.Stat

		children, err := client.GetChildren().StoringStatIn(&parentStat).ForPath("/parent")

		assert.Equal(s.T(), []string{"child"}, children)
		assert.NoError(s.T(), err)
	})
}

func (s *GetChildrenBuilderTestSuite) TestNamespace() {
	s.WithNamespace("parent", func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, stat *zk.Stat, acls []zk.ACL) {
		conn.On("Exists", "/parent").Return(false, nil, nil).Once()
		conn.On("Create", "/parent", []byte{}, int32(PERSISTENT), OPEN_ACL_UNSAFE).Return("/parent", nil).Once()
		conn.On("Children", "/parent/child").Return([]string{"node"}, stat, nil).Once()

		var parentStat zk.Stat

		children, err := client.GetChildren().StoringStatIn(&parentStat).ForPath("/child")

		assert.Equal(s.T(), []string{"node"}, children)
		assert.NoError(s.T(), err)
	})
}

func (s *GetChildrenBuilderTestSuite) TestBackground() {
	s.With(func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup, stat *zk.Stat) {
		ctxt := "context"
		children := []string{"child"}

		conn.On("Children", "/parent").Return(children, stat, nil).Once()

		_, err := client.GetChildren().InBackgroundWithCallbackAndContext(
			func(client CuratorFramework, event CuratorEvent) error {
				defer wg.Done()

				assert.Equal(s.T(), CHILDREN, event.Type())
				assert.Equal(s.T(), "/parent", event.Path())
				assert.Equal(s.T(), stat, event.Stat())
				assert.NoError(s.T(), event.Err())
				assert.Equal(s.T(), "parent", event.Name())
				assert.Equal(s.T(), children, event.Children())
				assert.Equal(s.T(), ctxt, event.Context())

				return nil
			}, ctxt).ForPath("/parent")

		assert.NoError(s.T(), err)
	})
}

func (s *GetChildrenBuilderTestSuite) TestWatcher() {
	s.With(func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup, data []byte, stat *zk.Stat) {
		events := make(chan zk.Event)

		defer close(events)

		conn.On("ChildrenW", "/parent").Return([]string{"child"}, stat, events, nil).Once()

		children, err := client.GetChildren().UsingWatcher(NewWatcher(func(event *zk.Event) {
			defer wg.Done()

			assert.NotNil(s.T(), event)
			assert.Equal(s.T(), zk.EventNodeChildrenChanged, event.Type)
			assert.Equal(s.T(), "/parent", event.Path)

		})).ForPath("/parent")

		assert.Equal(s.T(), []string{"child"}, children)
		assert.NoError(s.T(), err)

		events <- zk.Event{
			Type: zk.EventNodeChildrenChanged,
			Path: "/parent",
		}
	})
}
