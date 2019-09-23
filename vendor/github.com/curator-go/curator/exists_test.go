package curator

import (
	"sync"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CheckExistsBuilderTestSuite struct {
	mockContainerTestSuite
}

func TestCheckExistsBuilder(t *testing.T) {
	suite.Run(t, new(CheckExistsBuilderTestSuite))
}

func (s *CheckExistsBuilderTestSuite) TestCheckExists() {
	s.With(func(client CuratorFramework, conn *mockConn, stat *zk.Stat) {
		conn.On("Exists", "/node").Return(true, stat, nil).Once()

		stat2, err := client.CheckExists().ForPath("/node")

		assert.Equal(s.T(), stat, stat2)
		assert.NoError(s.T(), err)
	})

	s.With(func(client CuratorFramework, conn *mockConn) {
		conn.On("Exists", "/node").Return(false, nil, nil).Once()

		stat, err := client.CheckExists().ForPath("/node")

		assert.Nil(s.T(), stat)
		assert.Nil(s.T(), err)
	})

	s.With(func(client CuratorFramework, conn *mockConn) {
		conn.On("Exists", "/node").Return(false, nil, zk.ErrAPIError).Once()

		stat, err := client.CheckExists().ForPath("/node")

		assert.Nil(s.T(), stat)
		assert.EqualError(s.T(), err, zk.ErrAPIError.Error())
	})
}

func (s *CheckExistsBuilderTestSuite) TestNamespace() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn) {
		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Exists", "/parent/child").Return(false, nil, nil).Once()

		stat, err := client.CheckExists().ForPath("/child")

		assert.Nil(s.T(), stat)
		assert.Nil(s.T(), err)
	})
}

func (s *CheckExistsBuilderTestSuite) TestBackground() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup, stat *zk.Stat) {
		ctxt := "context"

		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Exists", "/parent/child").Return(true, stat, nil).Once()

		stat, err := client.CheckExists().InBackgroundWithCallbackAndContext(
			func(client CuratorFramework, event CuratorEvent) error {
				defer wg.Done()

				assert.Equal(s.T(), EXISTS, event.Type())
				assert.Equal(s.T(), "/child", event.Path())
				assert.NotNil(s.T(), event.Stat())
				assert.NoError(s.T(), event.Err())
				assert.Equal(s.T(), "child", event.Name())
				assert.Equal(s.T(), ctxt, event.Context())

				return nil
			}, ctxt).ForPath("/child")

		assert.Nil(s.T(), stat)
		assert.NoError(s.T(), err)
	})
}

func (s *CheckExistsBuilderTestSuite) TestWatcher() {
	s.With(func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup) {
		events := make(chan zk.Event)

		defer close(events)

		conn.On("ExistsW", "/node").Return(true, &zk.Stat{}, events, nil).Once()

		stat, err := client.CheckExists().UsingWatcher(NewWatcher(func(event *zk.Event) {
			defer wg.Done()

			assert.NotNil(s.T(), event)
			assert.Equal(s.T(), zk.EventNodeDeleted, event.Type)
			assert.Equal(s.T(), "/node", event.Path)

		})).ForPath("/node")

		assert.NotNil(s.T(), stat)
		assert.NoError(s.T(), err)

		events <- zk.Event{
			Type: zk.EventNodeDeleted,
			Path: "/node",
		}
	})
}
