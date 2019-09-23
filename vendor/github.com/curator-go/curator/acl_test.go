package curator

import (
	"sync"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetAclBuilderTestSuite struct {
	mockContainerTestSuite
}

func TestGetAclBuilder(t *testing.T) {
	suite.Run(t, new(GetAclBuilderTestSuite))
}

func (s *GetAclBuilderTestSuite) TestGetACL() {
	s.With(func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, stat *zk.Stat) {
		conn.On("GetACL", "/node").Return(READ_ACL_UNSAFE, stat, nil).Once()

		var nodeStat zk.Stat

		acls, err := client.GetACL().StoringStatIn(&nodeStat).ForPath("/node")

		assert.Equal(s.T(), acls, READ_ACL_UNSAFE)
		assert.NoError(s.T(), err)
	})
}

func (s *GetAclBuilderTestSuite) TestNamespace() {
	s.WithNamespace("parent", func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, stat *zk.Stat, acls []zk.ACL) {
		conn.On("Exists", "/parent").Return(false, nil, nil).Once()
		conn.On("Create", "/parent", []byte{}, int32(PERSISTENT), OPEN_ACL_UNSAFE).Return("/parent", nil).Once()
		conn.On("GetACL", "/parent/child").Return(READ_ACL_UNSAFE, stat, nil).Once()

		var nodeStat zk.Stat

		acls, err := client.GetACL().StoringStatIn(&nodeStat).ForPath("/child")

		assert.Equal(s.T(), acls, READ_ACL_UNSAFE)
		assert.NoError(s.T(), err)
	})
}

func (s *GetAclBuilderTestSuite) TestBackground() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn, wg *sync.WaitGroup, acls []zk.ACL, stat *zk.Stat) {
		ctxt := "context"

		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		conn.On("GetACL", "/parent/child").Return(acls, stat, nil).Once()

		_, err := client.GetACL().InBackgroundWithCallbackAndContext(
			func(client CuratorFramework, event CuratorEvent) error {
				defer wg.Done()

				assert.Equal(s.T(), GET_ACL, event.Type())
				assert.Equal(s.T(), "/child", event.Path())
				assert.Equal(s.T(), stat, event.Stat())
				assert.NoError(s.T(), event.Err())
				assert.Equal(s.T(), "child", event.Name())
				assert.Equal(s.T(), acls, event.ACLs())
				assert.Equal(s.T(), ctxt, event.Context())

				return nil
			}, ctxt).ForPath("/child")

		assert.NoError(s.T(), err)
	})
}

type SetAclBuilderTestSuite struct {
	mockContainerTestSuite
}

func TestSetAclBuilder(t *testing.T) {
	suite.Run(t, new(SetAclBuilderTestSuite))
}

func (s *SetAclBuilderTestSuite) TestGetACL() {
	s.With(func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, acls []zk.ACL, version int32, stat *zk.Stat) {
		conn.On("SetACL", "/node", acls, version).Return(stat, nil).Once()

		nodeStat, err := client.SetACL().WithACL(acls...).WithVersion(version).ForPath("/node")

		assert.Equal(s.T(), stat, nodeStat)
		assert.NoError(s.T(), err)
	})
}

func (s *SetAclBuilderTestSuite) TestNamespace() {
	s.WithNamespace("parent", func(builder *CuratorFrameworkBuilder, client CuratorFramework, conn *mockConn, version int32, stat *zk.Stat, acls []zk.ACL) {
		conn.On("Exists", "/parent").Return(false, nil, nil).Once()
		conn.On("Create", "/parent", []byte{}, int32(PERSISTENT), OPEN_ACL_UNSAFE).Return("/parent", nil).Once()
		conn.On("SetACL", "/parent/child", acls, version).Return(stat, nil).Once()

		nodeStat, err := client.SetACL().WithACL(acls...).WithVersion(version).ForPath("/child")

		assert.Equal(s.T(), stat, nodeStat)
		assert.NoError(s.T(), err)
	})
}

func (s *SetAclBuilderTestSuite) TestBackground() {
	s.WithNamespace("parent", func(client CuratorFramework, conn *mockConn, aclProvider *mockACLProvider, wg *sync.WaitGroup, acls []zk.ACL, version int32, stat *zk.Stat) {
		ctxt := "context"

		conn.On("Exists", "/parent").Return(true, nil, nil).Once()
		aclProvider.On("GetAclForPath", "/parent/child").Return(acls).Twice()
		conn.On("SetACL", "/parent/child", acls, version).Return(stat, nil).Once()

		_, err := client.SetACL().WithVersion(version).InBackgroundWithCallbackAndContext(
			func(client CuratorFramework, event CuratorEvent) error {
				defer wg.Done()

				assert.Equal(s.T(), SET_ACL, event.Type())
				assert.Equal(s.T(), "/child", event.Path())
				assert.Equal(s.T(), stat, event.Stat())
				assert.NoError(s.T(), event.Err())
				assert.Equal(s.T(), "child", event.Name())
				assert.Equal(s.T(), acls, event.ACLs())
				assert.Equal(s.T(), ctxt, event.Context())

				return nil
			}, ctxt).ForPath("/child")

		assert.NoError(s.T(), err)
	})
}
