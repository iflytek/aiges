package curator

import (
	"errors"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type infof func(format string, args ...interface{})

type mockCloseable struct {
	mock.Mock

	crash bool
}

func (c *mockCloseable) Close() error {
	if c.crash {
		panic(errors.New("panic"))
	}

	return c.Called().Error(0)
}

type mockTracerDriver struct {
	mock.Mock

	log infof
}

func (t *mockTracerDriver) AddTime(name string, d time.Duration) {
	if t.log != nil {
		t.log("TracerDriver.AddTime(name=\"%s\", d=%v)", name, d)
	}

	t.Called(name, d)
}

func (t *mockTracerDriver) AddCount(name string, increment int) {
	if t.log != nil {
		t.log("TracerDriver.AddCount(name=\"%s\", increment=%d)", name, increment)
	}

	t.Called(name, increment)
}

type mockRetrySleeper struct {
	mock.Mock

	log infof
}

func (s *mockRetrySleeper) SleepFor(time time.Duration) error {
	return s.Called(time).Error(0)
}

type mockRetryPolicy struct {
	mock.Mock

	log infof
}

func (r *mockRetryPolicy) AllowRetry(retryCount int, elapsedTime time.Duration, sleeper RetrySleeper) bool {
	args := r.Called(retryCount, elapsedTime, sleeper)

	allow := args.Bool(0)

	if r.log != nil {
		r.log("RetryPolicy.AllowRetry(retryCount=%d, elapsedTime=%v, sleeper=%p) allow=%v", retryCount, elapsedTime, sleeper, allow)
	}

	return allow
}

type mockEnsembleProvider struct {
	mock.Mock

	log infof
}

func (p *mockEnsembleProvider) Start() error {
	err := p.Called().Error(0)

	if p.log != nil {
		p.log("EnsembleProvider.Start() error=%v", err)
	}

	return err
}

func (p *mockEnsembleProvider) Close() error {
	err := p.Called().Error(0)

	if p.log != nil {
		p.log("EnsembleProvider.Close() error=%v", err)
	}

	return err
}

func (p *mockEnsembleProvider) ConnectionString() string {
	connStr := p.Called().String(0)

	if p.log != nil {
		p.log("EnsembleProvider.ConnectionString() \"%v\"", connStr)
	}

	return connStr
}

type mockConn struct {
	mock.Mock

	log        infof
	operations []interface{}
}

func (c *mockConn) AddAuth(scheme string, auth []byte) error {
	args := c.Called(scheme, auth)
	err := args.Error(0)

	if c.log != nil {
		c.log("ZookeeperConnection.AddAuth(scheme=\"%s\", auth=[]byte(\"%s\")) error=%v", scheme, auth, err)
	}

	return err
}

func (c *mockConn) Close() {
	if c.log != nil {
		c.log("ZookeeperConnection.Close()")
	}

	c.Called()
}

func (c *mockConn) Create(path string, data []byte, flags int32, acls []zk.ACL) (string, error) {
	args := c.Called(path, data, flags, acls)

	createPath := args.String(0)
	err := args.Error(1)

	if c.log != nil {
		c.log("ZookeeperConnection.Create(path=\"%s\", data=[]byte(\"%s\"), flags=%d, alcs=%v) (createdPath=\"%s\", error=%v)", path, data, flags, acls, createPath, err)
	}

	return createPath, err
}

func (c *mockConn) Exists(path string) (bool, *zk.Stat, error) {
	args := c.Called(path)

	exists := args.Bool(0)
	stat, _ := args.Get(1).(*zk.Stat)
	err := args.Error(2)

	if c.log != nil {
		c.log("ZookeeperConnection.Exists(path=\"%s\")(exists=%v, stat=%v, error=%v)", path, exists, stat, err)
	}

	return exists, stat, err
}

func (c *mockConn) ExistsW(path string) (bool, *zk.Stat, <-chan zk.Event, error) {
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

func (c *mockConn) Delete(path string, version int32) error {
	args := c.Called(path, version)

	err := args.Error(0)

	if c.log != nil {
		c.log("ZookeeperConnection.Delete(path=\"%s\", version=%d) error=%v", path, version, err)
	}

	return err
}

func (c *mockConn) Get(path string) ([]byte, *zk.Stat, error) {
	args := c.Called(path)

	data, _ := args.Get(0).([]byte)
	stat, _ := args.Get(1).(*zk.Stat)
	err := args.Error(2)

	if c.log != nil {
		c.log("ZookeeperConnection.Get(path=\"%s\")(data=%v, stat=%v, error=%v)", path, data, stat, err)
	}

	return data, stat, err
}

func (c *mockConn) GetW(path string) ([]byte, *zk.Stat, <-chan zk.Event, error) {
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

func (c *mockConn) Set(path string, data []byte, version int32) (*zk.Stat, error) {
	args := c.Called(path, data, version)

	stat, _ := args.Get(0).(*zk.Stat)
	err := args.Error(1)

	if c.log != nil {
		c.log("ZookeeperConnection.Set(path=\"%s\", data=%v, version=%d) (stat=%v, error=%v)", path, data, version, stat, err)
	}

	return stat, err
}

func (c *mockConn) Children(path string) ([]string, *zk.Stat, error) {
	args := c.Called(path)

	children, _ := args.Get(0).([]string)
	stat, _ := args.Get(1).(*zk.Stat)
	err := args.Error(2)

	if c.log != nil {
		c.log("ZookeeperConnection.Children(path=\"%s\")(children=%v, stat=%v, error=%v)", path, children, stat, err)
	}

	return children, stat, err
}

func (c *mockConn) ChildrenW(path string) ([]string, *zk.Stat, <-chan zk.Event, error) {
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

func (c *mockConn) GetACL(path string) ([]zk.ACL, *zk.Stat, error) {
	args := c.Called(path)

	acls, _ := args.Get(0).([]zk.ACL)
	stat, _ := args.Get(1).(*zk.Stat)
	err := args.Error(2)

	if c.log != nil {
		c.log("ZookeeperConnection.GetACL(path=\"%s\")(acls=%v, stat=%v, error=%v)", path, acls, stat, err)
	}

	return acls, stat, err
}

func (c *mockConn) SetACL(path string, acls []zk.ACL, version int32) (*zk.Stat, error) {
	args := c.Called(path, acls, version)

	stat, _ := args.Get(0).(*zk.Stat)
	err := args.Error(1)

	if c.log != nil {
		c.log("ZookeeperConnection.SetACL(path=\"%s\", acls=%v, version=%d) (stat=%v, error=%v)", path, acls, version, stat, err)
	}

	return stat, err
}

func (c *mockConn) Multi(ops ...interface{}) ([]zk.MultiResponse, error) {
	c.operations = append(c.operations, ops...)

	args := c.Called(ops)

	res, _ := args.Get(0).([]zk.MultiResponse)
	err := args.Error(1)

	if c.log != nil {
		c.log("ZookeeperConnection.Multi(ops=%v)(responses=%v, error=%v)", ops, res, err)
	}

	return res, err
}

func (c *mockConn) Sync(path string) (string, error) {
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

	log infof
}

func (d *mockZookeeperDialer) Dial(connString string, sessionTimeout time.Duration, canBeReadOnly bool) (ZookeeperConnection, <-chan zk.Event, error) {
	args := d.Called(connString, sessionTimeout, canBeReadOnly)

	conn, _ := args.Get(0).(ZookeeperConnection)
	events, _ := args.Get(1).(chan zk.Event)
	err := args.Error(2)

	if d.log != nil {
		d.log("ZookeeperDialer.Dial(connectString=\"%s\", sessionTimeout=%v, canBeReadOnly=%v)(conn=%p, events=%v, error=%v)", connString, sessionTimeout, canBeReadOnly, conn, events, err)
	}

	return conn, events, err
}

type mockCompressionProvider struct {
	mock.Mock

	log infof
}

func (p *mockCompressionProvider) Compress(path string, data []byte) ([]byte, error) {
	args := p.Called(path, data)

	compressedData, _ := args.Get(0).([]byte)
	err := args.Error(1)

	if p.log != nil {
		p.log("CompressionProvider.Compress(path=\"%s\", data=[]byte(\"%s\"))(compressedData=[]byte(\"%s\"), error=%v)", path, data, compressedData, err)
	}

	return compressedData, err
}

func (p *mockCompressionProvider) Decompress(path string, compressedData []byte) ([]byte, error) {
	args := p.Called(path, compressedData)

	data, _ := args.Get(0).([]byte)
	err := args.Error(1)

	if p.log != nil {
		p.log("CompressionProvider.Decompress(path=\"%s\", compressedData=[]byte(\"%s\"))(data=[]byte(\"%s\"), error=%v)", path, compressedData, data, err)
	}

	return data, err
}

type mockACLProvider struct {
	mock.Mock

	log infof
}

func (p *mockACLProvider) GetDefaultAcl() []zk.ACL {
	args := p.Called()

	acls, _ := args.Get(0).([]zk.ACL)

	if p.log != nil {
		p.log("ACLProvider.GetDefaultAcl()(acls=%v)", acls)
	}

	return acls
}

func (p *mockACLProvider) GetAclForPath(path string) []zk.ACL {
	args := p.Called(path)

	acls, _ := args.Get(0).([]zk.ACL)

	if p.log != nil {
		p.log("ACLProvider.GetAclForPath(path=\"%s\")(acls=%v)", path, acls)
	}

	return acls
}

type mockEnsurePath struct {
	mock.Mock

	log infof
}

func (e *mockEnsurePath) Ensure(client CuratorZookeeperClient) error {
	args := e.Mock.Called(client)

	err := args.Error(0)

	if e.log != nil {
		e.log("EnsurePath.Ensure(client=%p) error=%v", client, err)
	}

	return err
}

func (e *mockEnsurePath) ExcludingLast() EnsurePath {
	args := e.Mock.Called()

	ret, _ := args.Get(0).(EnsurePath)

	if e.log != nil {
		e.log("EnsurePath.ExcludingLast() EnsurePath=%p", ret)
	}

	return ret
}

type mockEnsurePathHelper struct {
	mock.Mock

	log infof
}

func (h *mockEnsurePathHelper) Ensure(client CuratorZookeeperClient, path string, makeLastNode bool) error {
	args := h.Called(client, path, makeLastNode)

	err := args.Error(0)

	if h.log != nil {
		h.log("EnsurePathHelper.Ensure(client=%p, path=\"%s\", makeLastNode=%v) error=%v", client, path, makeLastNode, err)
	}

	return err
}

type mockCuratorZookeeperClient struct {
	mock.Mock

	log infof
}

func (c *mockCuratorZookeeperClient) Conn() (ZookeeperConnection, error) {
	args := c.Called()

	conn, _ := args.Get(0).(ZookeeperConnection)
	err := args.Error(1)

	if c.log != nil {
		c.log("CuratorZookeeperClient.Conn() conn=%v", conn)
	}

	return conn, err
}

func (c *mockCuratorZookeeperClient) RetryPolicy() RetryPolicy {
	retryPolicy := c.Called().Get(0).(RetryPolicy)

	if c.log != nil {
		c.log("CuratorZookeeperClient.RetryPolicy() retryPolicy=%v", retryPolicy)
	}

	return retryPolicy
}

func (c *mockCuratorZookeeperClient) NewRetryLoop() RetryLoop {
	retryLoop, _ := c.Called().Get(0).(RetryLoop)

	if c.log != nil {
		c.log("CuratorZookeeperClient.NewRetryLoop() retryLoop=%v", retryLoop)
	}

	return retryLoop
}

func (c *mockCuratorZookeeperClient) Connected() bool {
	connected := c.Called().Bool(0)

	if c.log != nil {
		c.log("CuratorZookeeperClient.Connected() connected=%v", connected)
	}

	return connected
}

func (c *mockCuratorZookeeperClient) BlockUntilConnectedOrTimedOut() error {
	err := c.Called().Error(0)

	if c.log != nil {
		c.log("CuratorZookeeperClient.BlockUntilConnectedOrTimedOut() error=%v", err)
	}

	return err
}

func (c *mockCuratorZookeeperClient) Start() error {
	err := c.Called().Error(0)

	if c.log != nil {
		c.log("CuratorZookeeperClient.Start() error=%v", err)
	}

	return err
}

func (c *mockCuratorZookeeperClient) Close() error {
	err := c.Called().Error(0)

	if c.log != nil {
		c.log("CuratorZookeeperClient.Close() error=%v", err)
	}

	return err
}

func (c *mockCuratorZookeeperClient) StartTracer(name string) Tracer {
	tracer, _ := c.Called(name).Get(0).(Tracer)

	if c.log != nil {
		c.log("CuratorZookeeperClient.StartTracer(name=\"%s\") tracer=%v", name, tracer)
	}

	return tracer
}

type mockCuratorFramework struct {
	mock.Mock

	log infof
}

func (c *mockCuratorFramework) Start() error {
	err := c.Called().Error(0)

	if c.log != nil {
		c.log("CuratorFramework.Start() error=%v", err)
	}

	return err
}

func (c *mockCuratorFramework) Close() error {
	err := c.Called().Error(0)

	if c.log != nil {
		c.log("CuratorFramework.Close() error=%v", err)
	}

	return err
}

func (c *mockCuratorFramework) State() State {
	state, _ := c.Called().Get(0).(State)

	if c.log != nil {
		c.log("CuratorFramework.State() State=%v", state)
	}

	return state
}

func (c *mockCuratorFramework) Started() bool {
	started := c.Called().Bool(0)

	if c.log != nil {
		c.log("CuratorFramework.Started() started=%v", started)
	}

	return started
}

func (c *mockCuratorFramework) Create() CreateBuilder {
	builder, _ := c.Called().Get(0).(CreateBuilder)

	if c.log != nil {
		c.log("CuratorFramework.Create() CreateBuilder=%v", builder)
	}

	return builder
}

func (c *mockCuratorFramework) Delete() DeleteBuilder {
	builder, _ := c.Called().Get(0).(DeleteBuilder)

	if c.log != nil {
		c.log("CuratorFramework.Delete() DeleteBuilder=%v", builder)
	}

	return builder
}

func (c *mockCuratorFramework) CheckExists() CheckExistsBuilder {
	builder, _ := c.Called().Get(0).(CheckExistsBuilder)

	if c.log != nil {
		c.log("CuratorFramework.CheckExists() CheckExistsBuilder=%v", builder)
	}

	return builder
}

func (c *mockCuratorFramework) GetData() GetDataBuilder {
	builder, _ := c.Called().Get(0).(GetDataBuilder)

	if c.log != nil {
		c.log("CuratorFramework.GetData() GetDataBuilder=%v", builder)
	}

	return builder
}

func (c *mockCuratorFramework) SetData() SetDataBuilder {
	builder, _ := c.Called().Get(0).(SetDataBuilder)

	if c.log != nil {
		c.log("CuratorFramework.SetData() SetDataBuilder=%v", builder)
	}

	return builder
}

func (c *mockCuratorFramework) GetChildren() GetChildrenBuilder {
	builder, _ := c.Called().Get(0).(GetChildrenBuilder)

	if c.log != nil {
		c.log("CuratorFramework.GetChildren() GetChildrenBuilder=%v", builder)
	}

	return builder
}

func (c *mockCuratorFramework) GetACL() GetACLBuilder {
	builder, _ := c.Called().Get(0).(GetACLBuilder)

	if c.log != nil {
		c.log("CuratorFramework.GetACL() GetACLBuilder=%v", builder)
	}

	return builder
}

func (c *mockCuratorFramework) SetACL() SetACLBuilder {
	builder, _ := c.Called().Get(0).(SetACLBuilder)

	if c.log != nil {
		c.log("CuratorFramework.SetACL() SetACLBuilder=%v", builder)
	}

	return builder
}

func (c *mockCuratorFramework) InTransaction() Transaction {
	transaction, _ := c.Called().Get(0).(Transaction)

	if c.log != nil {
		c.log("CuratorFramework.InTransaction() Transaction=%v", transaction)
	}

	return transaction
}

func (c *mockCuratorFramework) DoSync(path string, backgroundContextObject interface{}) {
	c.Called(path, backgroundContextObject)

	if c.log != nil {
		c.log("CuratorFramework.DoSync(path=\"%s\", backgroundContextObject=%v)", path, backgroundContextObject)
	}
}

func (c *mockCuratorFramework) Sync() SyncBuilder {
	builder, _ := c.Called().Get(0).(SyncBuilder)

	if c.log != nil {
		c.log("CuratorFramework.Sync() SyncBuilder=%v", builder)
	}

	return builder
}

func (c *mockCuratorFramework) ConnectionStateListenable() ConnectionStateListenable {
	listenable, _ := c.Called().Get(0).(ConnectionStateListenable)

	if c.log != nil {
		c.log("CuratorFramework.ConnectionStateListenable() Listenable=%v", listenable)
	}

	return listenable
}

func (c *mockCuratorFramework) CuratorListenable() CuratorListenable {
	listenable, _ := c.Called().Get(0).(CuratorListenable)

	if c.log != nil {
		c.log("CuratorFramework.CuratorListenable() Listenable=%v", listenable)
	}

	return listenable
}

func (c *mockCuratorFramework) UnhandledErrorListenable() UnhandledErrorListenable {
	listenable, _ := c.Called().Get(0).(UnhandledErrorListenable)

	if c.log != nil {
		c.log("CuratorFramework.UnhandledErrorListenable() Listenable=%v", listenable)
	}

	return listenable
}

func (c *mockCuratorFramework) NonNamespaceView() CuratorFramework {
	framework, _ := c.Called().Get(0).(CuratorFramework)

	if c.log != nil {
		c.log("CuratorFramework.NonNamespaceView() Framework=%v", framework)
	}

	return framework
}

func (c *mockCuratorFramework) UsingNamespace(newNamespace string) CuratorFramework {
	framework, _ := c.Called(newNamespace).Get(0).(CuratorFramework)

	if c.log != nil {
		c.log("CuratorFramework.NonNamespaceView(newNamespace=\"%s\") Framework=%v", newNamespace, framework)
	}

	return framework
}

func (c *mockCuratorFramework) Namespace() string {
	namespace := c.Called().String(0)

	if c.log != nil {
		c.log("CuratorFramework.Namespace() Namespace=%v", namespace)
	}

	return namespace
}

func (c *mockCuratorFramework) ZookeeperClient() CuratorZookeeperClient {
	client, _ := c.Called().Get(0).(CuratorZookeeperClient)

	if c.log != nil {
		c.log("CuratorFramework.ZookeeperClient() Client=%v", client)
	}

	return client
}

func (c *mockCuratorFramework) NewNamespaceAwareEnsurePath(path string) EnsurePath {
	ensure, _ := c.Called(path).Get(0).(EnsurePath)

	if c.log != nil {
		c.log("CuratorFramework.NewNamespaceAwareEnsurePath(path=\"%s\") EnsurePath=%v", path, ensure)
	}

	return ensure
}

func (c *mockCuratorFramework) BlockUntilConnected() error {
	err := c.Called().Error(0)

	if c.log != nil {
		c.log("CuratorFramework.BlockUntilConnected() error=%v", err)
	}

	return err
}

func (c *mockCuratorFramework) BlockUntilConnectedTimeout(maxWaitTime time.Duration) error {
	err := c.Called(maxWaitTime).Error(0)

	if c.log != nil {
		c.log("CuratorFramework.BlockUntilConnectedTimeout(maxWaitTime=%v) error=%v", maxWaitTime, err)
	}

	return err
}

type mockContainer struct {
	builder *CuratorFrameworkBuilder
}

func newMockContainer() *mockContainer {
	return &mockContainer{
		builder: &CuratorFrameworkBuilder{
			SessionTimeout:    DEFAULT_SESSION_TIMEOUT,
			ConnectionTimeout: DEFAULT_CONNECTION_TIMEOUT,
			MaxCloseWait:      DEFAULT_CLOSE_WAIT,
			DefaultData:       []byte("default"),
		},
	}
}

func (c *mockContainer) Prepare(callback func(builder *CuratorFrameworkBuilder)) *mockContainer {
	callback(c.builder)

	return c
}

func (c *mockContainer) WithNamespace(namespace string) *mockContainer {
	c.builder.Namespace = namespace

	return c
}

func (c *mockContainer) Test(t *testing.T, callback interface{}) {
	var client CuratorFramework
	var events chan zk.Event
	var wg *sync.WaitGroup

	zookeeperConnection := &mockConn{log: t.Logf}
	zookeeperDialer := &mockZookeeperDialer{log: t.Logf}
	ensembleProvider := &mockEnsembleProvider{}
	compressionProvider := &mockCompressionProvider{log: t.Logf}
	retryPolicy := &mockRetryPolicy{log: t.Logf}
	aclProvider := &mockACLProvider{log: t.Logf}

	data := []byte("data")
	version := rand.Int31()
	stat := &zk.Stat{Version: version, Mtime: time.Now().Unix()}
	acls := zk.AuthACL(zk.PermRead)

	if c.builder.ZookeeperDialer == nil {
		c.builder.ZookeeperDialer = zookeeperDialer
	}

	if c.builder.EnsembleProvider == nil {
		c.builder.EnsembleProvider = ensembleProvider
	}

	if c.builder.CompressionProvider == nil {
		c.builder.CompressionProvider = compressionProvider
	}

	if c.builder.RetryPolicy == nil {
		c.builder.RetryPolicy = retryPolicy
	}

	if c.builder.AclProvider == nil {
		c.builder.AclProvider = aclProvider
	}

	fn := reflect.TypeOf(callback)

	assert.Equal(t, reflect.Func, fn.Kind())

	args := make([]reflect.Value, fn.NumIn())

	for i := 0; i < fn.NumIn(); i++ {
		switch argType := fn.In(i); argType {
		case reflect.TypeOf(c.builder):
			args[i] = reflect.ValueOf(c.builder)

		case reflect.TypeOf((*CuratorFramework)(nil)).Elem():
			client = c.builder.Build()
			args[i] = reflect.ValueOf(client)

		case reflect.TypeOf((*ZookeeperConnection)(nil)).Elem(), reflect.TypeOf(zookeeperConnection):
			args[i] = reflect.ValueOf(zookeeperConnection)

		case reflect.TypeOf((*ZookeeperDialer)(nil)).Elem(), reflect.TypeOf(zookeeperDialer):
			args[i] = reflect.ValueOf(zookeeperDialer)

		case reflect.TypeOf((*EnsembleProvider)(nil)).Elem(), reflect.TypeOf(ensembleProvider):
			args[i] = reflect.ValueOf(ensembleProvider)

		case reflect.TypeOf((*ZookeeperDialer)(nil)).Elem(), reflect.TypeOf(compressionProvider):
			args[i] = reflect.ValueOf(compressionProvider)

		case reflect.TypeOf((*RetryPolicy)(nil)).Elem(), reflect.TypeOf(retryPolicy):
			args[i] = reflect.ValueOf(retryPolicy)

		case reflect.TypeOf((*ACLProvider)(nil)).Elem(), reflect.TypeOf(aclProvider):
			args[i] = reflect.ValueOf(aclProvider)

		case reflect.TypeOf(events):
			events = make(chan zk.Event)
			args[i] = reflect.ValueOf(events)

		case reflect.TypeOf(wg):
			wg = new(sync.WaitGroup)
			args[i] = reflect.ValueOf(wg)

		case reflect.TypeOf(data):
			args[i] = reflect.ValueOf(data)

		case reflect.TypeOf(version):
			args[i] = reflect.ValueOf(version)

		case reflect.TypeOf(stat):
			args[i] = reflect.ValueOf(stat)

		case reflect.TypeOf(acls):
			args[i] = reflect.ValueOf(acls)

		default:
			t.Errorf("unknown arg type: %s", fn.In(i))
		}
	}

	if client != nil {
		if c.builder.EnsembleProvider == ensembleProvider {
			ensembleProvider.On("ConnectionString").Return("connStr").Once()
			ensembleProvider.On("Start").Return(nil).Once()
			ensembleProvider.On("Close").Return(nil).Once()
		}

		if c.builder.ZookeeperDialer == zookeeperDialer {
			zookeeperDialer.On("Dial", mock.AnythingOfType("string"), c.builder.SessionTimeout, c.builder.CanBeReadOnly).Return(zookeeperConnection, events, nil).Once()
		}

		assert.NoError(t, client.Start())
	}

	if wg != nil {
		wg.Add(1)
	}

	reflect.ValueOf(callback).Call(args)

	if wg != nil {
		wg.Wait()
	}

	if client != nil {
		if c.builder.ZookeeperDialer == zookeeperDialer {
			zookeeperConnection.On("Close").Return().Once()
		}

		assert.NoError(t, client.Close())
	}

	if events != nil {
		close(events)
	}

	zookeeperConnection.AssertExpectations(t)
	zookeeperDialer.AssertExpectations(t)
	ensembleProvider.AssertExpectations(t)
	compressionProvider.AssertExpectations(t)
	retryPolicy.AssertExpectations(t)
	aclProvider.AssertExpectations(t)
}

type mockContainerTestSuite struct {
	suite.Suite
}

func (s *mockContainerTestSuite) With(callback interface{}) {
	newMockContainer().Test(s.T(), callback)
}

func (s *mockContainerTestSuite) WithNamespace(namespace string, callback interface{}) {
	newMockContainer().WithNamespace(namespace).Test(s.T(), callback)
}

func (s *mockContainerTestSuite) WithPrepare(prepare func(*CuratorFrameworkBuilder), callback interface{}) {
	newMockContainer().Prepare(prepare).Test(s.T(), callback)
}
