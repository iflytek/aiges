package metricCollector

import (
	"sync"
	"time"
)

type MetricCollector interface {
	Update(MetricResult)
	Reset()
}

var Registry = metricCollectorRegistry{
	lock: &sync.RWMutex{},
	registry: []func(name string, win time.Duration) MetricCollector{
		newDefaultMetricCollector,
	},
}

type metricCollectorRegistry struct {
	lock     *sync.RWMutex
	registry []func(name string, win time.Duration) MetricCollector
}

func (m *metricCollectorRegistry) InitializeMetricCollectors(name string, win time.Duration) []MetricCollector {
	m.lock.RLock()
	defer m.lock.RUnlock()
	metrics := make([]MetricCollector, len(m.registry))
	for i, metricCollectorInitializer := range m.registry {
		metrics[i] = metricCollectorInitializer(name, win)
	}
	return metrics
}
func (m *metricCollectorRegistry) Register(initMetricCollector func(string, time.Duration) MetricCollector) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.registry = append(m.registry, initMetricCollector)
}

type MetricResult struct {
	Attempts                float64
	Errors                  float64
	Successes               float64
	Failures                float64
	Rejects                 float64
	ShortCircuits           float64
	Timeouts                float64
	FallbackSuccesses       float64
	FallbackFailures        float64
	ContextCanceled         float64
	ContextDeadlineExceeded float64
	TotalDuration           time.Duration
	RunDuration             time.Duration
	ConcurrencyInUse        float64
}
