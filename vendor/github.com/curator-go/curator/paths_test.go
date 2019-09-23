package curator

import (
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
)

func TestGetNodeFromPath(t *testing.T) {
	assert.Equal(t, "child", GetNodeFromPath("child"))

	assert.Equal(t, "", GetNodeFromPath("/child/"))

	assert.Equal(t, "child", GetNodeFromPath("/child"))
	assert.Equal(t, "child", GetNodeFromPath("/parent/child"))
}

func TestSplitPath(t *testing.T) {
	p, err := SplitPath("test")

	assert.NoError(t, err)
	assert.Equal(t, p.Path, "test")
	assert.Equal(t, p.Node, "")

	p, err = SplitPath("/test/hello")

	assert.NoError(t, err)
	assert.Equal(t, p.Path, "/test")
	assert.Equal(t, p.Node, "hello")

	p, err = SplitPath("/hello")

	assert.NoError(t, err)
	assert.Equal(t, p.Path, "/")
	assert.Equal(t, p.Node, "hello")
}

func TestJoinPath(t *testing.T) {
	assert.Equal(t, JoinPath("parent", "child"), "/parent/child")
	assert.Equal(t, JoinPath("parent/", "child"), "/parent/child")
	assert.Equal(t, JoinPath("/parent/", "child"), "/parent/child")

	assert.Equal(t, JoinPath("", "child"), "/child")
	assert.Equal(t, JoinPath("parent", "", "child"), "/parent/child")
	assert.Equal(t, JoinPath("parent", "child/"), "/parent/child")
	assert.Equal(t, JoinPath("parent", "/child"), "/parent/child")
	assert.Equal(t, JoinPath("parent", "child1", "child2"), "/parent/child1/child2")
}

func TestValidatePath(t *testing.T) {
	assert.EqualError(t, ValidatePath(""), "Path cannot be null")

	assert.EqualError(t, ValidatePath("test"), "Path must start with / character")

	assert.EqualError(t, ValidatePath("/test/"), "Path must not end with / character")

	assert.EqualError(t, ValidatePath("/\x00"), "null character not allowed @ 1")

	assert.EqualError(t, ValidatePath("//test"), "empty node name specified @ 1")

	assert.EqualError(t, ValidatePath("/.."), "relative paths not allowed @ 2")
	assert.EqualError(t, ValidatePath("/../test"), "relative paths not allowed @ 2")
	assert.EqualError(t, ValidatePath("/."), "relative paths not allowed @ 1")
	assert.EqualError(t, ValidatePath("/./test"), "relative paths not allowed @ 1")

	assert.EqualError(t, ValidatePath("/\u0010"), "invalid charater @ 1")
	assert.EqualError(t, ValidatePath("/\u007f"), "invalid charater @ 1")
	assert.EqualError(t, ValidatePath("/\uf805"), "invalid charater @ 1")
	assert.EqualError(t, ValidatePath("/\ufff0"), "invalid charater @ 1")
}

func TestMakeDirs(t *testing.T) {
	// skip exists `parent` and create `child`
	conn := &mockConn{}

	conn.On("Exists", "/parent").Return(true, nil, nil).Once()
	conn.On("Exists", "/parent/child").Return(false, nil, nil).Once()
	conn.On("Create", "/parent/child", []byte{}, int32(PERSISTENT), OPEN_ACL_UNSAFE).Return("", nil).Once()

	assert.NoError(t, MakeDirs(conn, "/parent/child/node", false, nil))

	conn.AssertExpectations(t)

	// fail to create `parent`
	conn = &mockConn{}

	conn.On("Exists", "/parent").Return(true, nil, zk.ErrAPIError).Once()

	assert.EqualError(t, MakeDirs(conn, "/parent/child/node", false, nil), zk.ErrAPIError.Error())

	conn.AssertExpectations(t)

	// create `child` which exists
	conn = &mockConn{}

	conn.On("Exists", "/parent").Return(true, nil, nil).Once()
	conn.On("Exists", "/parent/child").Return(false, nil, nil).Once()
	conn.On("Create", "/parent/child", []byte{}, int32(PERSISTENT), OPEN_ACL_UNSAFE).Return("", zk.ErrNodeExists).Once()

	assert.NoError(t, MakeDirs(conn, "/parent/child/node", false, nil))

	conn.AssertExpectations(t)

	// create `child` with default ACLs
	conn = &mockConn{}
	acls := &mockACLProvider{}

	conn.On("Exists", "/parent").Return(true, nil, nil).Once()
	conn.On("Exists", "/parent/child").Return(false, nil, nil).Once()

	acls.On("GetAclForPath", "/parent/child").Return([]zk.ACL{}).Once()
	acls.On("GetDefaultAcl").Return(zk.AuthACL(zk.PermAdmin)).Once()

	conn.On("Create", "/parent/child", []byte{}, int32(PERSISTENT), zk.AuthACL(zk.PermAdmin)).Return("", nil).Once()

	assert.NoError(t, MakeDirs(conn, "/parent/child/node", false, acls))

	conn.AssertExpectations(t)
	acls.AssertExpectations(t)
}

func TestDeleteChildren(t *testing.T) {
	// Delete children
	conn := &mockConn{}

	conn.On("Children", "/parent").Return([]string{"child1", "child2"}, nil, nil).Once()
	conn.On("Children", "/parent/child1").Return(nil, nil, nil).Once()
	conn.On("Children", "/parent/child2").Return(nil, nil, nil).Once()
	conn.On("Delete", "/parent/child1", AnyVersion).Return(nil).Once()
	conn.On("Delete", "/parent/child2", AnyVersion).Return(zk.ErrNoNode).Once()

	assert.NoError(t, DeleteChildren(conn, "/parent", false))

	conn.AssertExpectations(t)

	// Children failed
	conn = &mockConn{}

	conn.On("Children", "/parent").Return(nil, nil, zk.ErrNoNode).Once()

	assert.Equal(t, DeleteChildren(conn, "/parent", false), zk.ErrNoNode)

	conn.AssertExpectations(t)
}

func TestEnsurePath(t *testing.T) {
	helper := &mockEnsurePathHelper{log: t.Logf}

	ensure := NewEnsurePathWithAclAndHelper("/parent/child", nil, helper)

	assert.NotNil(t, ensure)
	assert.True(t, ensure.makeLastNode)

	ensure2 := ensure.ExcludingLast()

	assert.NotNil(t, ensure2)

	client := &mockCuratorZookeeperClient{log: t.Logf}

	helper.On("Ensure", client, "/parent/child", true).Return(nil).Once()

	assert.NoError(t, ensure.Ensure(client))

	helper.On("Ensure", client, "/parent/child", false).Return(nil).Once()

	assert.NoError(t, ensure2.Ensure(client))

	helper.AssertExpectations(t)
	client.AssertExpectations(t)
}
