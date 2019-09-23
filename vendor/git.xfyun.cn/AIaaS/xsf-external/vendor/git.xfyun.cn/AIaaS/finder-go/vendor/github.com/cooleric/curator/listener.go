package curator

import (
	"sync"
)

type ConnectionStateListener interface {
	// Called when there is a state change in the connection
	StateChanged(client CuratorFramework, newState ConnectionState)
}

// Receives notifications about errors and background events
type CuratorListener interface {
	// Called when a background task has completed or a watch has triggered
	EventReceived(client CuratorFramework, event CuratorEvent) error
}

type UnhandledErrorListener interface {
	// Called when an exception is caught in a background thread, handler, etc.
	UnhandledError(err error)
}

type connectionStateListenerCallback func(client CuratorFramework, newState ConnectionState)

type connectionStateListenerStub struct {
	callback connectionStateListenerCallback
}

func NewConnectionStateListener(callback connectionStateListenerCallback) ConnectionStateListener {
	return &connectionStateListenerStub{callback}
}

func (l *connectionStateListenerStub) StateChanged(client CuratorFramework, newState ConnectionState) {
	l.callback(client, newState)
}

type curatorListenerCallback func(client CuratorFramework, event CuratorEvent) error

type curatorListenerStub struct {
	callback curatorListenerCallback
}

func NewCuratorListener(callback curatorListenerCallback) CuratorListener {
	return &curatorListenerStub{callback}
}

func (l *curatorListenerStub) EventReceived(client CuratorFramework, event CuratorEvent) error {
	return l.callback(client, event)
}

type unhandledErrorCallback func(error)

func (cb unhandledErrorCallback) UnhandledError(err error) {
	cb(err)
}

// NewUnhandledErrorListener creates an UnhandledErrorListener with given callback
func NewUnhandledErrorListener(callback func(error)) UnhandledErrorListener {
	return unhandledErrorCallback(callback)
}

// Abstracts a listenable object
type Listenable /* [T] */ interface {
	Len() int

	Clear()

	ForEach(callback func(interface{}))
}

type ConnectionStateListenable interface {
	Listenable /* [T] */

	AddListener(listener ConnectionStateListener)

	RemoveListener(listener ConnectionStateListener)
}

type CuratorListenable interface {
	Listenable /* [T] */

	AddListener(listener CuratorListener)

	RemoveListener(listener CuratorListener)
}

type UnhandledErrorListenable interface {
	Listenable /* [T] */

	AddListener(listener UnhandledErrorListener)

	RemoveListener(listener UnhandledErrorListener)
}

type ListenerContainer struct {
	lock      sync.RWMutex
	listeners []interface{}
}

func (c *ListenerContainer) Add(listener interface{}) {
	if c != nil {
		c.lock.Lock()

		c.listeners = append(c.listeners, listener)

		c.lock.Unlock()
	}
}

func (c *ListenerContainer) Remove(listener interface{}) {
	if c == nil {
		return
	}

	c.lock.Lock()

	for i, l := range c.listeners {
		if l == listener {
			c.listeners = append(c.listeners[:i], c.listeners[i+1:]...)
			break
		}
	}

	c.lock.Unlock()
}

func (c *ListenerContainer) Len() int {
	if c == nil {
		return 0
	}

	return len(c.listeners)
}

func (c *ListenerContainer) Clear() {
	if c == nil {
		return
	}

	c.lock.Lock()

	c.listeners = nil

	c.lock.Unlock()
}

func (c *ListenerContainer) ForEach(callback func(interface{})) {
	if c == nil {
		return
	}

	c.lock.RLock()

	for _, listener := range c.listeners {
		callback(listener)
	}

	c.lock.RUnlock()
}

type connectionStateListenerContainer struct {
	ListenerContainer
}

func (c *connectionStateListenerContainer) AddListener(listener ConnectionStateListener) {
	c.Add(listener)
}

func (c *connectionStateListenerContainer) RemoveListener(listener ConnectionStateListener) {
	c.Remove(listener)
}

type curatorListenerContainer struct {
	ListenerContainer
}

func (c *curatorListenerContainer) AddListener(listener CuratorListener) {
	c.Add(listener)
}

func (c *curatorListenerContainer) RemoveListener(listener CuratorListener) {
	c.Remove(listener)
}

type UnhandledErrorListenerContainer struct {
	ListenerContainer
}

func (c *UnhandledErrorListenerContainer) AddListener(listener UnhandledErrorListener) {
	c.Add(listener)
}

func (c *UnhandledErrorListenerContainer) RemoveListener(listener UnhandledErrorListener) {
	c.Remove(listener)
}
