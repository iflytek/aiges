package xsf

import (
	"fmt"
	"sync"
	"time"
)

type CircuitBreaker struct {
	parent                 *CircuitBreakers
	Name                   string
	open                   bool
	forceOpen              bool
	mutex                  *sync.RWMutex
	openedOrLastTestedTime int64

	metrics       *metricExchange
	executorPool  *executorPool
	streamHandler *StreamHandler
}
type CircuitBreakers struct {
	parent               *Circuit
	circuitBreakers      map[string]*CircuitBreaker
	circuitBreakersMutex *sync.RWMutex
}

func (c *CircuitBreakers) init(parent *Circuit) {
	c.parent = parent
	c.circuitBreakersMutex = &sync.RWMutex{}
	c.circuitBreakers = make(map[string]*CircuitBreaker)
}

func (c *CircuitBreakers) getCircuitBreaker(name string) (*CircuitBreaker, bool, error) {
	c.circuitBreakersMutex.RLock()
	_, ok := c.circuitBreakers[name]
	if !ok {
		c.circuitBreakersMutex.RUnlock()
		c.circuitBreakersMutex.Lock()
		defer c.circuitBreakersMutex.Unlock()

		if cb, ok := c.circuitBreakers[name]; ok {
			return cb, false, nil
		}
		c.circuitBreakers[name] = newCircuitBreaker(c, name)
	} else {
		defer c.circuitBreakersMutex.RUnlock()
	}

	return c.circuitBreakers[name], !ok, nil
}

func (c *CircuitBreakers) flush() {
	c.circuitBreakersMutex.Lock()
	defer c.circuitBreakersMutex.Unlock()

	for name, cb := range c.circuitBreakers {
		cb.metrics.Reset()
		cb.executorPool.Metrics.Reset()
		delete(c.circuitBreakers, name)
	}
}

func newCircuitBreaker(parent *CircuitBreakers, name string) *CircuitBreaker {
	c := &CircuitBreaker{}
	c.Name = name
	c.parent = parent
	c.mutex = &sync.RWMutex{}
	c.metrics = newMetricExchange(c, name)
	c.executorPool = newExecutorPool(c, name)
	c.streamHandler = newStreamHandler(c, name)

	return c
}
func (c *CircuitBreaker) collectMetrics() ([][]byte, error) {
	return c.streamHandler.collect()
}
func (c *CircuitBreaker) toggleForceOpen(toggle bool) error {
	c.forceOpen = toggle
	return nil
}
func (c *CircuitBreaker) getSettings(name string) *Settings {
	return c.parent.parent.getSettings(name)
}
func (c *CircuitBreaker) GetExecutorPool() *executorPool {
	return c.getExecutorPool()
}
func (c *CircuitBreaker) getExecutorPool() *executorPool {
	return c.executorPool
}

func (c *CircuitBreaker) isOpen() (bool, error) {
	c.mutex.RLock()
	open, forceOpen := c.open, c.forceOpen
	c.mutex.RUnlock()

	if open || forceOpen {
		return true, fmt.Errorf("open:%v,forceOpen:%v", open, forceOpen)
	}

	if requests, requestVolumeThreshold := uint64(c.metrics.Requests().Sum(time.Now())),
		c.getSettings(c.Name).RequestVolumeThreshold;
		requests >= requestVolumeThreshold {
		return true, fmt.Errorf("the requests(%d) exceed the RequestVolumeThreshold(%d)", requests, requestVolumeThreshold)
	}

	if isHealthy, isHealthyErr := c.metrics.IsHealthy(time.Now()); !isHealthy {
		c.setOpen()
		return true, isHealthyErr
	}

	return false, nil
}

func (c *CircuitBreaker) AllowRequest() (bool, error) {
	return c.allowRequest()
}
func (c *CircuitBreaker) allowRequest() (bool, error) {
	isOpen, isOpenErr := c.isOpen()
	return !isOpen, isOpenErr
}

func (c *CircuitBreaker) setOpen() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.open {
		return
	}

	c.openedOrLastTestedTime = time.Now().UnixNano()
	c.open = true
}

func (c *CircuitBreaker) setClose() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.open {
		return
	}

	c.open = false
	c.metrics.Reset()
}

func (c *CircuitBreaker) ReportEvent(eventTypes []string, start time.Time, runDuration time.Duration) error {
	return c.reportEvent(eventTypes, start, runDuration)
}
func (c *CircuitBreaker) reportEvent(eventTypes []string, start time.Time, runDuration time.Duration) error {
	if len(eventTypes) == 0 {
		return fmt.Errorf("no event types sent for metrics")
	}

	c.mutex.RLock()
	o := c.open
	c.mutex.RUnlock()
	if eventTypes[0] == "success" && o {
		c.setClose()
	}

	var concurrencyInUse float64
	if c.executorPool.Max > 0 {
		concurrencyInUse = float64(c.executorPool.ActiveCount()) / float64(c.executorPool.Max)
	}

	select {
	case c.metrics.Updates <- &commandExecution{
		Types:            eventTypes,
		Start:            start,
		RunDuration:      runDuration,
		ConcurrencyInUse: concurrencyInUse,
	}:
	default:
		fmt.Printf("metrics channel (%v) is at capacity\n", c.Name) //todo not complete
		//return CircuitError{Message: fmt.Sprintf("metrics channel (%v) is at capacity", c.Name)}
		return nil
	}

	return nil
}
