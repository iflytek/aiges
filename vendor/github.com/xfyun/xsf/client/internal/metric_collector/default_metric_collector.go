package metricCollector

import (
	"github.com/xfyun/xsf/client/internal/rolling"
	"sync"
	"time"
)

type DefaultMetricCollector struct {
	win                     time.Duration
	mutex                   *sync.RWMutex
	numRequests             *rolling.Number
	errors                  *rolling.Number
	successes               *rolling.Number
	failures                *rolling.Number
	rejects                 *rolling.Number
	shortCircuits           *rolling.Number
	timeouts                *rolling.Number
	contextCanceled         *rolling.Number
	contextDeadlineExceeded *rolling.Number
	fallbackSuccesses       *rolling.Number
	fallbackFailures        *rolling.Number
	totalDuration           *rolling.Timing
	runDuration             *rolling.Timing
}

func newDefaultMetricCollector(name string, win time.Duration) MetricCollector {
	m := &DefaultMetricCollector{}
	m.win = win
	m.mutex = &sync.RWMutex{}
	m.Reset()
	return m
}
func (d *DefaultMetricCollector) NumRequests() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.numRequests
}
func (d *DefaultMetricCollector) Errors() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.errors
}
func (d *DefaultMetricCollector) Successes() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.successes
}
func (d *DefaultMetricCollector) Failures() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.failures
}
func (d *DefaultMetricCollector) Rejects() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.rejects
}
func (d *DefaultMetricCollector) ShortCircuits() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.shortCircuits
}
func (d *DefaultMetricCollector) Timeouts() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.timeouts
}
func (d *DefaultMetricCollector) FallbackSuccesses() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.fallbackSuccesses
}
func (d *DefaultMetricCollector) ContextCanceled() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.contextCanceled
}
func (d *DefaultMetricCollector) ContextDeadlineExceeded() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.contextDeadlineExceeded
}
func (d *DefaultMetricCollector) FallbackFailures() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.fallbackFailures
}
func (d *DefaultMetricCollector) TotalDuration() *rolling.Timing {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.totalDuration
}
func (d *DefaultMetricCollector) RunDuration() *rolling.Timing {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.runDuration
}
func (d *DefaultMetricCollector) Update(r MetricResult) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	d.numRequests.Increment(r.Attempts)
	d.errors.Increment(r.Errors)
	d.successes.Increment(r.Successes)
	d.failures.Increment(r.Failures)
	d.rejects.Increment(r.Rejects)
	d.shortCircuits.Increment(r.ShortCircuits)
	d.timeouts.Increment(r.Timeouts)
	d.fallbackSuccesses.Increment(r.FallbackSuccesses)
	d.fallbackFailures.Increment(r.FallbackFailures)
	d.contextCanceled.Increment(r.ContextCanceled)
	d.contextDeadlineExceeded.Increment(r.ContextDeadlineExceeded)
	d.totalDuration.Add(r.TotalDuration)
	d.runDuration.Add(r.RunDuration)
}
func (d *DefaultMetricCollector) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.numRequests = rolling.NewNumber(d.win)
	d.errors = rolling.NewNumber(d.win)
	d.successes = rolling.NewNumber(d.win)
	d.rejects = rolling.NewNumber(d.win)
	d.shortCircuits = rolling.NewNumber(d.win)
	d.failures = rolling.NewNumber(d.win)
	d.timeouts = rolling.NewNumber(d.win)
	d.fallbackSuccesses = rolling.NewNumber(d.win)
	d.fallbackFailures = rolling.NewNumber(d.win)
	d.contextCanceled = rolling.NewNumber(d.win)
	d.contextDeadlineExceeded = rolling.NewNumber(d.win)
	d.totalDuration = rolling.NewTiming(d.win)
	d.runDuration = rolling.NewTiming(d.win)
}
