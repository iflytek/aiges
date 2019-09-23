package curator

import (
	"sync"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetDataBuilderTestSuite struct {
	mockContainerTestSuite
}

func TestGetDataBuilder(t *testing.T) {
	suite.Run(t, new(GetDataBuilderTestSuite))
}

func (s *GetDataBuilderTestSuite) TestGetData() {
	s.With(func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, compress *mockCompressionProvider, data []byte, stat *zk.Stat) {
		conn.On("Get", "/node").Return([]byte("compressed(data)"), stat, nil).Once()
		compress.On("Decompress", "/node", []byte("compressed(data)")).Return(data, nil).Once()

		var stat2 zk.Stat

		data2, err := client.GetData().StoringStatIn(&stat2).Decompressed().ForPath("/node")

		assert.Equal(s.T(), data, data2)
		assert.Equal(s.T(), stat, &stat2)
		assert.NoError(s.T(), err)
	})
}

func (s *GetDataBuilderTestSuite) TestNamespace() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn, data []byte, stat *zk.Stat) {
		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Get", "/parent/child").Return(data, stat, nil).Once()

		data2, err := client.GetData().ForPath("/child")

		assert.Equal(s.T(), data, data2)
		assert.NoError(s.T(), err)
	})
}

func (s *GetDataBuilderTestSuite) TestBackground() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup, data []byte, stat *zk.Stat) {
		ctxt := "context"

		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Get", "/parent/child").Return(data, stat, nil).Once()

		_, err := client.GetData().InBackgroundWithCallbackAndContext(
			func(client CuratorFramework, event CuratorEvent) error {
				defer wg.Done()

				assert.Equal(s.T(), GET_DATA, event.Type())
				assert.Equal(s.T(), "/child", event.Path())
				assert.Equal(s.T(), data, event.Data())
				assert.NoError(s.T(), event.Err())
				assert.Equal(s.T(), "child", event.Name())
				assert.Equal(s.T(), ctxt, event.Context())

				return nil
			}, ctxt).ForPath("/child")

		assert.NoError(s.T(), err)
	})
}

func (s *GetDataBuilderTestSuite) TestWatcher() {
	s.With(func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup, data []byte, stat *zk.Stat) {
		events := make(chan zk.Event)

		defer close(events)

		conn.On("GetW", "/node").Return(data, stat, events, nil).Once()
		data2, err := client.GetData().UsingWatcher(NewWatcher(func(event *zk.Event) {
			defer wg.Done()

			assert.NotNil(s.T(), event)
			assert.Equal(s.T(), zk.EventNodeDataChanged, event.Type)
			assert.Equal(s.T(), "/node", event.Path)

		})).ForPath("/node")

		assert.Equal(s.T(), data, data2)
		assert.NoError(s.T(), err)

		events <- zk.Event{
			Type: zk.EventNodeDataChanged,
			Path: "/node",
		}
	})
}

type SetDataBuilderTestSuite struct {
	mockContainerTestSuite
}

func TestSetDataBuilder(t *testing.T) {
	suite.Run(t, new(SetDataBuilderTestSuite))
}

func (s *SetDataBuilderTestSuite) TestSetData() {
	s.With(func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, compress *mockCompressionProvider, data []byte, version int32, stat *zk.Stat) {
		compress.On("Compress", "/node", data).Return([]byte("compressed(data)"), nil).Once()
		conn.On("Set", "/node", []byte("compressed(data)"), version).Return(stat, nil).Once()

		stat2, err := client.SetData().WithVersion(version).Compressed().ForPathWithData("/node", data)

		assert.Equal(s.T(), stat, stat2)
		assert.NoError(s.T(), err)
	})
}

func (s *SetDataBuilderTestSuite) TestNamespace() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn, data []byte, stat *zk.Stat) {
		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Set", "/parent/child", data, AnyVersion).Return(stat, nil).Once()

		stat2, err := client.SetData().ForPathWithData("/child", data)

		assert.Equal(s.T(), stat, stat2)
		assert.NoError(s.T(), err)
	})
}

func (s *SetDataBuilderTestSuite) TestBackground() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup, data []byte, stat *zk.Stat) {
		ctxt := "context"

		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Set", "/parent/child", data, AnyVersion).Return(stat, zk.ErrAPIError).Once()

		_, err := client.SetData().InBackgroundWithCallbackAndContext(
			func(client CuratorFramework, event CuratorEvent) error {
				defer wg.Done()

				assert.Equal(s.T(), SET_DATA, event.Type())
				assert.Equal(s.T(), "/child", event.Path())
				assert.Equal(s.T(), data, event.Data())
				assert.Equal(s.T(), stat, event.Stat())
				assert.EqualError(s.T(), event.Err(), zk.ErrAPIError.Error())
				assert.Equal(s.T(), "child", event.Name())
				assert.Equal(s.T(), ctxt, event.Context())

				return nil
			}, ctxt).ForPathWithData("/child", data)

		assert.NoError(s.T(), err)
	})
}
