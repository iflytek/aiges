package curator

import (
	"sync"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CreateBuilderTestSuite struct {
	mockContainerTestSuite
}

func TestCreateBuilder(t *testing.T) {
	suite.Run(t, new(CreateBuilderTestSuite))
}

func (s *CreateBuilderTestSuite) TestCreate() {
	s.With(func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, acls []zk.ACL) {
		conn.On("Create", "/node", builder.DefaultData, int32(EPHEMERAL), acls).Return("/node", nil).Once()

		path, err := client.Create().WithMode(EPHEMERAL).WithACL(acls...).ForPath("/node")

		assert.Equal(s.T(), "/node", path)
		assert.NoError(s.T(), err)
	})
}

func (s *CreateBuilderTestSuite) TestNamespace() {
	s.WithNamespace("parent", func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, acls []zk.ACL) {
		conn.On("Exists", "/parent").Return(false, nil, nil).Once()
		conn.On("Create", "/parent", []byte{}, int32(PERSISTENT), OPEN_ACL_UNSAFE).Return("/parent", nil).Once()
		conn.On("Create", "/parent/child", builder.DefaultData, int32(EPHEMERAL), acls).Return("/parent/child", nil).Once()

		path, err := client.Create().WithMode(EPHEMERAL).WithACL(acls...).ForPath("/child")

		assert.Equal(s.T(), "/child", path)
		assert.NoError(s.T(), err)
	})
}

func (s *CreateBuilderTestSuite) TestBackground() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup, data []byte, acls []zk.ACL) {
		ctxt := "context"

		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("Create", "/parent/child", data, int32(PERSISTENT), acls).Return("", zk.ErrAPIError).Once()

		path, err := client.Create().WithACL(acls...).InBackgroundWithCallbackAndContext(
			func(client CuratorFramework, event CuratorEvent) error {
				defer wg.Done()

				assert.Equal(s.T(), CREATE, event.Type())
				assert.Equal(s.T(), "/child", event.Path())
				assert.Equal(s.T(), data, event.Data())
				assert.Equal(s.T(), acls, event.ACLs())
				assert.EqualError(s.T(), event.Err(), zk.ErrAPIError.Error())
				assert.Equal(s.T(), "child", event.Name())
				assert.Equal(s.T(), ctxt, event.Context())

				return nil
			}, ctxt).ForPathWithData("/child", data)

		assert.Equal(s.T(), "/child", path)
		assert.NoError(s.T(), err)
	})
}

func (s *CreateBuilderTestSuite) TestCompression() {
	s.With(func(client CuratorFramework, conn *mockConn, compress *mockCompressionProvider, data []byte, acls []zk.ACL) {
		compressedData := []byte("compressedData")

		compress.On("Compress", "/node", data).Return(compressedData, nil).Once()
		conn.On("Create", "/node", compressedData, int32(PERSISTENT), acls).Return("/node", nil).Once()

		path, err := client.Create().Compressed().WithACL(acls...).ForPathWithData("/node", data)

		assert.Equal(s.T(), "/node", path)
		assert.NoError(s.T(), err)
	})
}

func (s *CreateBuilderTestSuite) TestCreateParents() {
	s.With(func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, data []byte, aclProvider *mockACLProvider, acls []zk.ACL) {
		aclProvider.On("GetAclForPath", "/parent/child").Return(READ_ACL_UNSAFE).Once()
		conn.On("Create", "/parent/child", data, int32(PERSISTENT), READ_ACL_UNSAFE).Return("", zk.ErrNoNode).Once()

		conn.On("Exists", "/parent").Return(false, nil, nil).Once()
		aclProvider.On("GetAclForPath", "/parent").Return(CREATOR_ALL_ACL).Once()
		conn.On("Create", "/parent", []byte{}, int32(PERSISTENT), CREATOR_ALL_ACL).Return("", zk.ErrAPIError).Once()

		path, err := client.Create().CreatingParentsIfNeeded().ForPathWithData("/parent/child", data)

		assert.Equal(s.T(), "", path)
		assert.Equal(s.T(), err, zk.ErrAPIError)
	})
}
