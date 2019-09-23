package curator

import (
	"sync"
	"time"
)

// Mechanism for timing methods and recording counters
type TracerDriver interface {
	// Record the given trace event
	AddTime(name string, d time.Duration)

	// Add to a named counter
	AddCount(name string, increment int)
}

type Tracer interface {
	Commit()
}

type defaultTracerDriver struct {
	TracerDriver

	logger   func(fmt string, args ...interface{})
	lock     sync.Mutex
	counters map[string]int
}

func newDefaultTracerDriver() *defaultTracerDriver {
	return &defaultTracerDriver{counters: make(map[string]int)}
}

func (d *defaultTracerDriver) AddTime(name string, time time.Duration) {
	if d.logger != nil {
		d.logger("Trace %s: %s", name, time)
	}
}

func (d *defaultTracerDriver) AddCount(name string, increment int) {
	d.lock.Lock()

	value, _ := d.counters[name]

	value += increment

	d.counters[name] = value

	d.lock.Unlock()

	if d.logger != nil {
		d.logger("Counter %s: %d + %d", name, value-increment, increment)
	}
}

// Utility to time a method or portion of code
type timeTracer struct {
	name      string
	driver    TracerDriver
	startTime time.Time
}

// Create and start a timer
func newTimeTracer(name string, driver TracerDriver) *timeTracer {
	return &timeTracer{
		name:      name,
		driver:    driver,
		startTime: time.Now(),
	}
}

// Record the elapsed time
func (t *timeTracer) Commit() {
	t.CommitAt(time.Now())
}

// Record the elapsed time
func (t *timeTracer) CommitAt(tm time.Time) {
	t.driver.AddTime(t.name, tm.Sub(t.startTime))
}
