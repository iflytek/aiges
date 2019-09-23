package recipes

import (
	"testing"
	"time"

	"github.com/curator-go/curator"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/mock"
)

type logFunc func(format string, args ...interface{})

type mockZookeeperConnection struct {
	mock.Mock

	log        logFunc
	operations []interface{}
}

func (c *mockZookeeperConnection) AddAuth(scheme string, auth []byte) error {
	args := c.Called(scheme, auth)
	err := args.Error(0)

	if c.log != nil {
		c.log("ZookeeperConnection.AddAuth(scheme=\"%s\", auth=[]byte(\"%s\")) error=%v", scheme, auth, err)
	}

	return err
}

func (c *mockZookeeperConnection) Close() {
	if c.log != nil {
		c.log("ZookeeperConnection.Close()")
	}

	c.Called()
}

func (c *mockZookeeperConnection) Create(path string, data []byte, flags int32, acls []zk.ACL) (string, error) {
	args := c.Called(path, data, flags, acls)

	createPath := args.String(0)
	err := args.Error(1)

	if c.log != nil {
		c.log("ZookeeperConnection.Create(path=\"%s\", data=[]byte(\"%s\"), flags=%d, alcs=%v) (createdPath=\"%s\", error=%v)", path, data, flags, acls, createPath, err)
	}

	return createPath, err
}

func (c *mockZookeeperConnection) Exists(path string) (bool, *zk.Stat, error) {
	args := c.Called(path)

	exists := args.Bool(0)
	stat, _ := args.Get(1).(*zk.Stat)
	err := args.Error(2)

	if c.log != nil {
		c.log("ZookeeperConnection.Exists(path=\"%s\")(exists=%v, stat=%v, error=%v)", path, exists, stat, err)
	}

	return exists, stat, err
}

func (c *mockZookeeperConnection) ExistsW(path string) (bool, *zk.Stat, <-chan zk.Event, error) {
	args := c.Called(path)

	exists := args.Bool(0)
	stat, _ := args.Get(1).(*zk.Stat)
	events, _ := args.Get(2).(chan zk.Event)
	err := args.Error(3)

	if c.log != nil {
		c.log("ZookeeperConnection.ExistsW(path=\"%s\")(exists=%v, stat=%v, events=%v, error=%v)", path, exists, stat, events, err)
	}

	return exists, stat, events, err
}

func (c *mockZookeeperConnection) Delete(path string, version int32) error {
	args := c.Called(path, version)

	err := args.Error(0)

	if c.log != nil {
		c.log("ZookeeperConnection.Delete(path=\"%s\", version=%d) error=%v", path, version, err)
	}

	return err
}

func (c *mockZookeeperConnection) Get(path string) ([]byte, *zk.Stat, error) {
	args := c.Called(path)

	data, _ := args.Get(0).([]byte)
	stat, _ := args.Get(1).(*zk.Stat)
	err := args.Error(2)

	if c.log != nil {
		c.log("ZookeeperConnection.Get(path=\"%s\")(data=%v, stat=%v, error=%v)", path, data, stat, err)
	}

	return data, stat, err
}

func (c *mockZookeeperConnection) GetW(path string) ([]byte, *zk.Stat, <-chan zk.Event, error) {
	args := c.Called(path)

	data, _ := args.Get(0).([]byte)
	stat, _ := args.Get(1).(*zk.Stat)
	events, _ := args.Get(2).(chan zk.Event)
	err := args.Error(3)

	if c.log != nil {
		c.log("ZookeeperConnection.GetW(path=\"%s\")(data=%v, stat=%v, events=%p, error=%v)", path, data, stat, err)
	}

	return data, stat, events, err
}

func (c *mockZookeeperConnection) Set(path string, data []byte, version int32) (*zk.Stat, error) {
	args := c.Called(path, data, version)

	stat, _ := args.Get(0).(*zk.Stat)
	err := args.Error(1)

	if c.log != nil {
		c.log("ZookeeperConnection.Set(path=\"%s\", data=%v, version=%d) (stat=%v, error=%v)", path, data, version, stat, err)
	}

	return stat, err
}

func (c *mockZookeeperConnection) Children(path string) ([]string, *zk.Stat, error) {
	args := c.Called(path)

	children, _ := args.Get(0).([]string)
	stat, _ := args.Get(1).(*zk.Stat)
	err := args.Error(2)

	if c.log != nil {
		c.log("ZookeeperConnection.Children(path=\"%s\")(children=%v, stat=%v, error=%v)", path, children, stat, err)
	}

	return children, stat, err
}

func (c *mockZookeeperConnection) ChildrenW(path string) ([]string, *zk.Stat, <-chan zk.Event, error) {
	args := c.Called(path)

	children, _ := args.Get(0).([]string)
	stat, _ := args.Get(1).(*zk.Stat)
	events, _ := args.Get(2).(chan zk.Event)
	err := args.Error(3)

	if c.log != nil {
		c.log("ZookeeperConnection.ChildrenW(path=\"%s\")(children=%v, stat=%v, events=%v, error=%v)", path, children, stat, events, err)
	}

	return children, stat, events, err
}

func (c *mockZookeeperConnection) GetACL(path string) ([]zk.ACL, *zk.Stat, error) {
	args := c.Called(path)

	acls, _ := args.Get(0).([]zk.ACL)
	stat, _ := args.Get(1).(*zk.Stat)
	err := args.Error(2)

	if c.log != nil {
		c.log("ZookeeperConnection.GetACL(path=\"%s\")(acls=%v, stat=%v, error=%v)", path, acls, stat, err)
	}

	return acls, stat, err
}

func (c *mockZookeeperConnection) SetACL(path string, acls []zk.ACL, version int32) (*zk.Stat, error) {
	args := c.Called(path, acls, version)

	stat, _ := args.Get(0).(*zk.Stat)
	err := args.Error(1)

	if c.log != nil {
		c.log("ZookeeperConnection.SetACL(path=\"%s\", acls=%v, version=%d) (stat=%v, error=%v)", path, acls, version, stat, err)
	}

	return stat, err
}

func (c *mockZookeeperConnection) Multi(ops ...interface{}) ([]zk.MultiResponse, error) {
	c.operations = append(c.operations, ops...)

	args := c.Called(ops)

	res, _ := args.Get(0).([]zk.MultiResponse)
	err := args.Error(1)

	if c.log != nil {
		c.log("ZookeeperConnection.Multi(ops=%v)(responses=%v, error=%v)", ops, res, err)
	}

	return res, err
}

func (c *mockZookeeperConnection) Sync(path string) (string, error) {
	args := c.Called(path)
	p := args.String(0)
	err := args.Error(1)

	if c.log != nil {
		c.log("ZookeeperConnection.Sync(path=\"%s\")(path=\"%s\", error=%v)", path, p, err)
	}

	return path, err
}

type mockZookeeperDialer struct {
	mock.Mock

	log logFunc
}

func (d *mockZookeeperDialer) Dial(connString string, sessionTimeout time.Duration, canBeReadOnly bool) (curator.ZookeeperConnection, <-chan zk.Event, error) {
	args := d.Called(connString, sessionTimeout, canBeReadOnly)

	conn, _ := args.Get(0).(curator.ZookeeperConnection)
	events, _ := args.Get(1).(chan zk.Event)
	err := args.Error(2)

	if d.log != nil {
		d.log("ZookeeperDialer.Dial(connectString=\"%s\", sessionTimeout=%v, canBeReadOnly=%v)(conn=%p, events=%v, error=%v)", connString, sessionTimeout, canBeReadOnly, conn, events, err)
	}

	return conn, events, err
}

type mockRetryPolicy struct {
	mock.Mock

	log logFunc
}

func (r *mockRetryPolicy) AllowRetry(retryCount int, elapsedTime time.Duration, sleeper curator.RetrySleeper) bool {
	args := r.Called(retryCount, elapsedTime, sleeper)

	allow := args.Bool(0)

	if r.log != nil {
		r.log("RetryPolicy.AllowRetry(retryCount=%d, elapsedTime=%v, sleeper=%p) allow=%v", retryCount, elapsedTime, sleeper, allow)
	}

	return allow
}

type mockLockInternalsDriver struct {
	mock.Mock

	log logFunc
}

func (d *mockLockInternalsDriver) FixForSorting(str, lockName string) string {
	suffix := d.Called(str, lockName).String(0)

	if d.log != nil {
		d.log("LockInternalsDriver.FixForSorting(str=\"%s\", lockName=\"%s\") path=\"%s\"", str, lockName, suffix)
	}

	return suffix
}

func (d *mockLockInternalsDriver) GetsTheLock(client curator.CuratorFramework, children []string, sequenceNodeName string, maxLeases int) (*PredicateResults, error) {
	args := d.Called(client, children, sequenceNodeName, maxLeases)

	ret, _ := args.Get(0).(*PredicateResults)
	err := args.Error(1)

	if d.log != nil {
		d.log("LockInternalsDriver.FixForSorting(client=%p, children=%v, sequenceNodeName=\"%s\", maxLeases=%d) (results=%v, err=%v)", client, children, sequenceNodeName, maxLeases, ret, err)
	}

	return ret, err
}

func (d *mockLockInternalsDriver) CreatesTheLock(client curator.CuratorFramework, path string, lockNodeBytes []byte) (string, error) {
	args := d.Called(client, path, lockNodeBytes)

	str := args.String(0)
	err := args.Error(1)

	if d.log != nil {
		d.log("LockInternalsDriver.FixForSorting(client=%p, path=\"%s\", lockNodeBytes=%v) (path=%s, err=%v)", client, path, lockNodeBytes, str, err)
	}

	return str, err
}

type mockBuilder struct {
	conn        *mockZookeeperConnection
	events      chan zk.Event
	dialer      *mockZookeeperDialer
	builder     *curator.CuratorFrameworkBuilder
	retryPolicy *mockRetryPolicy
	driver      *mockLockInternalsDriver
}

func newMockBuilder(t *testing.T) *mockBuilder {
	conn := &mockZookeeperConnection{log: t.Logf}

	dialer := &mockZookeeperDialer{log: t.Logf}
	builder := &curator.CuratorFrameworkBuilder{ZookeeperDialer: dialer}

	return &mockBuilder{
		conn:        conn,
		events:      make(chan zk.Event),
		dialer:      dialer,
		builder:     builder,
		retryPolicy: &mockRetryPolicy{log: t.Logf},
		driver:      &mockLockInternalsDriver{log: t.Logf},
	}
}

func (b *mockBuilder) Build() curator.CuratorFramework {
	b.builder.ConnectString("connStr")

	b.dialer.On("Dial", "connStr", curator.DEFAULT_SESSION_TIMEOUT, b.builder.CanBeReadOnly).Return(b.conn, b.events, nil).Once()

	return b.builder.Build()
}

func (b *mockBuilder) Check(t *testing.T) {
	b.conn.AssertExpectations(t)
	b.dialer.AssertExpectations(t)
	b.retryPolicy.AssertExpectations(t)
	b.driver.AssertExpectations(t)
}
