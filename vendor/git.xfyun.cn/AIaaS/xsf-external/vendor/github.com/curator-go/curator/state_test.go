package curator

import (
	"sync"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestHandleHolder(t *testing.T) {
	ensembleProvider := &mockEnsembleProvider{log: t.Logf}
	zookeeperDialer := &mockZookeeperDialer{log: t.Logf}
	zookeeperConnection := &mockConn{log: t.Logf}
	events := make(chan zk.Event)
	watcher := NewWatcher(func(event *zk.Event) {})

	// create new connection holder
	h := &handleHolder{
		zookeeperDialer:  zookeeperDialer,
		ensembleProvider: ensembleProvider,
		watcher:          watcher,
		sessionTimeout:   15 * time.Second,
	}

	assert.Equal(t, "", h.getConnectionString())
	assert.False(t, h.hasNewConnectionString())
	conn, err := h.getZookeeperConnection()
	assert.Nil(t, conn)
	assert.NoError(t, err)

	// close and reset helper
	assert.NoError(t, h.closeAndReset())

	assert.NotNil(t, h.Helper())
	assert.Equal(t, "", h.getConnectionString())
	assert.IsType(t, (*zookeeperFactory)(nil), h.Helper())

	// get and create connection
	ensembleProvider.On("ConnectionString").Return("connStr").Once()
	zookeeperDialer.On("Dial", "connStr", h.sessionTimeout, h.canBeReadOnly).Return(zookeeperConnection, events, nil).Once()

	conn, err = h.getZookeeperConnection()

	assert.NotNil(t, conn)
	assert.NoError(t, err)
	assert.Equal(t, "connStr", h.getConnectionString())
	assert.NotNil(t, h.Helper())
	assert.IsType(t, (*zookeeperCache)(nil), h.Helper())

	// close connection
	zookeeperConnection.On("Close").Return(nil).Once()

	assert.NoError(t, h.closeAndClear())

	close(events)

	ensembleProvider.AssertExpectations(t)
	zookeeperDialer.AssertExpectations(t)
	zookeeperConnection.AssertExpectations(t)
}

type ConnectionStateTestSuite struct {
	suite.Suite

	ensembleProvider  *mockEnsembleProvider
	zookeeperDialer   *mockZookeeperDialer
	conn              *mockConn
	tracer            *mockTracerDriver
	sessionTimeout    time.Duration
	connectionTimeout time.Duration
	canBeReadOnly     bool
	events            chan zk.Event
	watcher           Watcher
	sessionEvents     []*zk.Event
	state             *connectionState
	connStrTimes      int
	dialTimes         int
	connCloseTimes    int
}

func TestConnectionState(t *testing.T) {
	suite.Run(t, &ConnectionStateTestSuite{
		sessionTimeout:    15 * time.Second,
		connectionTimeout: 5 * time.Second,
	})
}

func (s *ConnectionStateTestSuite) SetupTest() {
	s.ensembleProvider = &mockEnsembleProvider{log: s.T().Logf}
	s.zookeeperDialer = &mockZookeeperDialer{log: s.T().Logf}
	s.conn = &mockConn{log: s.T().Logf}
	s.tracer = &mockTracerDriver{log: s.T().Logf}
	s.events = make(chan zk.Event)
	s.sessionEvents = nil
	s.connStrTimes = 1
	s.dialTimes = 1
	s.connCloseTimes = 1

	s.watcher = NewWatcher(func(event *zk.Event) {
		s.sessionEvents = append(s.sessionEvents, event)
	})

	// create connection
	s.state = newConnectionState(s.zookeeperDialer, s.ensembleProvider, s.sessionTimeout, s.connectionTimeout, s.watcher, s.tracer, s.canBeReadOnly)

	assert.NotNil(s.T(), s.state)
	assert.False(s.T(), s.state.Connected())
}

func (s *ConnectionStateTestSuite) Start() {
	// start connection
	s.ensembleProvider.On("Start").Return(nil).Once()
	s.ensembleProvider.On("ConnectionString").Return("connStr").Times(s.connStrTimes)
	s.zookeeperDialer.On("Dial", "connStr", s.sessionTimeout, s.canBeReadOnly).Return(s.conn, s.events, nil).Times(s.dialTimes)
	s.conn.On("Close").Return().Times(s.connCloseTimes)

	assert.NoError(s.T(), s.state.Start())
	assert.False(s.T(), s.state.Connected())
}

func (s *ConnectionStateTestSuite) Close() {
	// close connection
	s.ensembleProvider.On("Close").Return(nil).Once()

	assert.NoError(s.T(), s.state.Close())
}

func (s *ConnectionStateTestSuite) TearDownTest() {
	s.ensembleProvider.AssertExpectations(s.T())
	s.zookeeperDialer.AssertExpectations(s.T())
	s.conn.AssertExpectations(s.T())
	s.tracer.AssertExpectations(s.T())

	s.sessionEvents = nil

	close(s.events)
}

func (s *ConnectionStateTestSuite) TestConnectFailed() {
	// start connection
	s.ensembleProvider.On("Start").Return(nil).Once()
	s.ensembleProvider.On("ConnectionString").Return("connStr").Times(2)
	s.zookeeperDialer.On("Dial", "connStr", s.sessionTimeout, s.canBeReadOnly).Return(nil, nil, zk.ErrAPIError).Times(2)

	instanceIndex := s.state.InstanceIndex()

	assert.Equal(s.T(), s.state.Start(), zk.ErrAPIError)
	assert.False(s.T(), s.state.Connected())
	assert.Equal(s.T(), instanceIndex+1, s.state.InstanceIndex())

	conn, err := s.state.Conn()

	assert.Nil(s.T(), conn)
	assert.Equal(s.T(), zk.ErrAPIError, err)
	assert.Equal(s.T(), instanceIndex+1, s.state.InstanceIndex())

	// close connection
	s.ensembleProvider.On("Close").Return(nil).Once()

	assert.NoError(s.T(), s.state.Close())
	assert.Equal(s.T(), instanceIndex+1, s.state.InstanceIndex())
}

func (s *ConnectionStateTestSuite) TestConnectionTimeout() {
	s.connStrTimes = 2

	s.Start()
	defer s.Close()

	instanceIndex := s.state.InstanceIndex()

	// force to connect timeout
	s.state.connectionStart.Store(time.Now().Add(-s.connectionTimeout * 2))

	s.tracer.On("AddCount", "connections-timed-out", 1).Return().Once()

	conn, err := s.state.Conn()

	assert.Nil(s.T(), conn)
	assert.Equal(s.T(), ErrConnectionLoss, err)
	assert.Equal(s.T(), instanceIndex, s.state.InstanceIndex())
}

func (s *ConnectionStateTestSuite) TestSessionTimeout() {
	s.connStrTimes = 3
	s.dialTimes = 2
	s.connCloseTimes = 2

	s.Start()
	defer s.Close()

	instanceIndex := s.state.InstanceIndex()

	// force to session timeout
	s.state.connectionStart.Store(time.Now().Add(-s.sessionTimeout * 2))

	s.tracer.On("AddCount", "session-timed-out", 1).Return().Once()

	conn, err := s.state.Conn()

	assert.NotNil(s.T(), conn)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), instanceIndex+1, s.state.InstanceIndex())
}

func (s *ConnectionStateTestSuite) TestBackgroundException() {
	s.tracer.On("AddCount", "background-exceptions", 1).Return().Times(2)

	// deque from empty queue
	assert.NoError(s.T(), s.state.dequeBackgroundException())

	// enque and deque
	s.state.queueBackgroundException(ErrConnectionLoss)

	assert.Equal(s.T(), ErrConnectionLoss, s.state.dequeBackgroundException())
	assert.NoError(s.T(), s.state.dequeBackgroundException())

	// enque too many errors
	s.tracer.On("AddCount", "connection-drop-background-error", 1).Return().Once()

	s.state.queueBackgroundException(zk.ErrAPIError)

	for i := 0; i < MAX_BACKGROUND_ERRORS; i++ {
		s.state.queueBackgroundException(ErrConnectionLoss)
	}

	assert.Equal(s.T(), ErrConnectionLoss, s.state.dequeBackgroundException())
}

func (s *ConnectionStateTestSuite) TestConnected() {
	s.connStrTimes = 2

	s.Start()
	defer s.Close()

	instanceIndex := s.state.InstanceIndex()

	// get the connection
	conn, err := s.state.Conn()

	assert.NotNil(s.T(), conn)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), instanceIndex, s.state.InstanceIndex())

	// receive a session event
	s.tracer.On("AddTime", "connection-state-parent-process", mock.AnythingOfType("Duration")).Return().Once()

	s.events <- zk.Event{
		Type:  zk.EventSession,
		State: zk.StateHasSession,
	}

	time.Sleep(100 * time.Microsecond)

	assert.True(s.T(), s.state.Connected())
}

func (s *ConnectionStateTestSuite) TestNewConnectionString() {
	s.connStrTimes = 3
	s.dialTimes = 2
	s.connCloseTimes = 2

	s.Start()
	defer s.Close()

	instanceIndex := s.state.InstanceIndex()

	// receive StateHasSession event
	s.tracer.On("AddTime", "connection-state-parent-process", mock.AnythingOfType("Duration")).Return().Once()
	s.tracer.On("AddCount", "connection-string-changed", 1).Return().Once()

	s.state.zooKeeper.Helper().(*zookeeperCache).connnectString = "anotherStr"

	s.events <- zk.Event{
		Type:  zk.EventSession,
		State: zk.StateHasSession,
	}

	time.Sleep(100 * time.Microsecond)

	assert.Equal(s.T(), instanceIndex+1, s.state.InstanceIndex())
	assert.False(s.T(), s.state.Connected())
}

func (s *ConnectionStateTestSuite) TestExpiredSession() {
	s.connStrTimes = 3
	s.dialTimes = 2
	s.connCloseTimes = 2

	s.Start()
	defer s.Close()

	instanceIndex := s.state.InstanceIndex()

	// get the connection
	conn, err := s.state.Conn()

	assert.NotNil(s.T(), conn)
	assert.NoError(s.T(), err)

	// receive StateHasSession event
	s.tracer.On("AddTime", "connection-state-parent-process", mock.AnythingOfType("Duration")).Return().Twice()

	s.events <- zk.Event{
		Type:  zk.EventSession,
		State: zk.StateHasSession,
	}

	time.Sleep(100 * time.Microsecond)

	assert.True(s.T(), s.state.Connected())

	// receive StateExpired event
	s.tracer.On("AddCount", "session-expired", 1).Return().Once()

	s.events <- zk.Event{
		Type:  zk.EventSession,
		State: zk.StateExpired,
	}

	time.Sleep(100 * time.Microsecond)

	assert.Equal(s.T(), instanceIndex+1, s.state.InstanceIndex())
	assert.False(s.T(), s.state.Connected())
}

func (s *ConnectionStateTestSuite) TestParentWatcher() {
	s.connStrTimes = 4

	s.Start()
	defer s.Close()

	// get the connection
	conn, err := s.state.Conn()

	assert.NotNil(s.T(), conn)
	assert.NoError(s.T(), err)

	// receive a session event
	s.tracer.On("AddTime", "connection-state-parent-process", mock.AnythingOfType("Duration")).Return().Twice()

	assert.Nil(s.T(), s.sessionEvents)

	s.events <- zk.Event{
		Type:  zk.EventSession,
		State: zk.StateConnecting,
	}

	time.Sleep(100 * time.Microsecond)

	assert.Equal(s.T(), 1, len(s.sessionEvents))

	assert.Equal(s.T(), s.state.RemoveParentWatcher(s.watcher), s.watcher)

	s.events <- zk.Event{
		Type:  zk.EventSession,
		State: zk.StateConnected,
	}

	time.Sleep(100 * time.Microsecond)

	assert.Equal(s.T(), s.state.AddParentWatcher(s.watcher), s.watcher)

	s.events <- zk.Event{
		Type:  zk.EventSession,
		State: zk.StateDisconnected,
	}

	time.Sleep(100 * time.Microsecond)

	assert.Equal(s.T(), 2, len(s.sessionEvents))
	assert.Equal(s.T(), zk.StateConnecting, s.sessionEvents[0].State)
	assert.Equal(s.T(), zk.StateDisconnected, s.sessionEvents[1].State)
}

type ConnectionStateManagerTestSuite struct {
	suite.Suite

	client         *mockCuratorFramework
	state          *connectionStateManager
	receivedStates []ConnectionState
}

func TestConnectionStateManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectionStateManagerTestSuite))
}

func (s *ConnectionStateManagerTestSuite) SetupTest() {
	s.client = &mockCuratorFramework{}
	s.state = newConnectionStateManager(s.client)
	s.state.Listenable().AddListener(NewConnectionStateListener(func(client CuratorFramework, newState ConnectionState) {
		s.receivedStates = append(s.receivedStates, newState)
	}))
}

func (s *ConnectionStateManagerTestSuite) TearDownTest() {
	s.client.AssertExpectations(s.T())

	s.receivedStates = nil
}

func (s *ConnectionStateManagerTestSuite) TestPostState() {
	assert.NoError(s.T(), s.state.Start())

	for i := 0; i < STATE_QUEUE_SIZE; i++ {
		s.state.postState(CONNECTED)
	}

	for i := 0; i < STATE_QUEUE_SIZE; i++ {
		s.state.postState(RECONNECTED)
	}

	state := <-s.state.events

	assert.Equal(s.T(), RECONNECTED, state)

	s.state.Close()

	s.state.postState(LOST)
}

func (s *ConnectionStateManagerTestSuite) TestStateChange() {
	// return false before StateManager.Start
	assert.False(s.T(), s.state.AddStateChange(CONNECTED))
	assert.False(s.T(), s.state.SetToSuspended())

	assert.NoError(s.T(), s.state.Start())

	defer s.state.Close()

	assert.NotNil(s.T(), s.state.Listenable())
	assert.False(s.T(), s.state.Connected())

	// return false when newState same to current
	assert.False(s.T(), s.state.AddStateChange(s.state.currentConnectionState))

	// automatic broadcast CONNECTED on the first time
	assert.True(s.T(), s.state.AddStateChange(RECONNECTED))

	time.Sleep(100 * time.Millisecond)

	assert.Equal(s.T(), RECONNECTED, s.state.currentConnectionState)
	assert.Equal(s.T(), []ConnectionState{CONNECTED}, s.receivedStates)
	assert.True(s.T(), s.state.Connected())

	// broadcast RECONNECTED on the second time
	assert.True(s.T(), s.state.AddStateChange(CONNECTED))

	time.Sleep(100 * time.Millisecond)

	assert.Equal(s.T(), CONNECTED, s.state.currentConnectionState)
	assert.Equal(s.T(), []ConnectionState{CONNECTED, CONNECTED}, s.receivedStates)
	assert.True(s.T(), s.state.Connected())

	// set to suspend
	assert.True(s.T(), s.state.SetToSuspended())
	assert.Equal(s.T(), SUSPENDED, s.state.currentConnectionState)

	time.Sleep(100 * time.Millisecond)

	assert.False(s.T(), s.state.SetToSuspended())
	assert.Equal(s.T(), []ConnectionState{CONNECTED, CONNECTED, SUSPENDED}, s.receivedStates)
}

func (s *ConnectionStateManagerTestSuite) TestBlockUntilConnected() {
	var wc sync.WaitGroup

	assert.NoError(s.T(), s.state.Start())

	defer s.state.Close()

	wc.Add(1)

	go func() {
		defer wc.Done()

		assert.NoError(s.T(), s.state.BlockUntilConnected(0))
	}()

	time.Sleep(100 * time.Millisecond)

	s.state.AddStateChange(CONNECTED)

	wc.Wait()

	assert.Equal(s.T(), CONNECTED, s.state.currentConnectionState)
}

func (s *ConnectionStateManagerTestSuite) TestBlockUntilConnectedWithTimeout() {
	var wc sync.WaitGroup

	assert.NoError(s.T(), s.state.Start())

	defer s.state.Close()

	wc.Add(1)

	go func() {
		defer wc.Done()

		assert.NoError(s.T(), s.state.BlockUntilConnected(time.Second))
	}()

	time.Sleep(100 * time.Millisecond)

	s.state.AddStateChange(CONNECTED)

	wc.Wait()

	assert.Equal(s.T(), CONNECTED, s.state.currentConnectionState)
}

func (s *ConnectionStateManagerTestSuite) TestBlockUntilConnectedTimeouted() {
	var wc sync.WaitGroup

	assert.NoError(s.T(), s.state.Start())

	defer s.state.Close()

	wc.Add(1)

	go func() {
		defer wc.Done()

		assert.Equal(s.T(), ErrTimeout, s.state.BlockUntilConnected(100*time.Millisecond))
	}()

	wc.Wait()

	assert.Equal(s.T(), UNKNOWN, s.state.currentConnectionState)
}

func assertEvent(t *testing.T, e *zk.Event, ch chan *zk.Event, timeout time.Duration) {
	select {
	case <-time.After(timeout):
		t.Fatal("Waiting for event timed out: ", e)
	case evt := <-ch:
		assert.Equal(t, e, evt)
	}
}

func createTestingWatcher() (chan *zk.Event, Watcher) {
	var ch = make(chan *zk.Event)
	var watcher = NewWatcher(func(event *zk.Event) {
		ch <- event
	})
	return ch, watcher
}

func TestProcessingMultipleWatchers(t *testing.T) {
	var ch1, w1 = createTestingWatcher()
	var ch2, w2 = createTestingWatcher()
	var evt = &zk.Event{Type: zk.EventSession}
	var state = newConnectionState(nil, nil, time.Second, time.Second, nil, newDefaultTracerDriver(), false)

	state.AddParentWatcher(w1)
	state.AddParentWatcher(w2)
	state.process(evt)

	assertEvent(t, evt, ch1, time.Second)
	assertEvent(t, evt, ch2, time.Second)
}
