package xsf

import (
	"encoding/json"
	"github.com/xfyun/xsf/client/internal/rolling"
	"time"
)

type StreamHandler struct {
	parent *CircuitBreaker
	Name   string
}

func newStreamHandler(parent *CircuitBreaker, name string) *StreamHandler {
	s := &StreamHandler{}
	s.Name = name
	s.parent = parent
	return s
}

func (sh *StreamHandler) collect() ([][]byte, error) {
	var dataSet [][]byte
	metricsBytes, err := sh.publishMetrics(sh.parent)
	if err != nil {
		return dataSet, err
	}
	dataSet = append(dataSet, metricsBytes)

	threadPoolsBytes, err := sh.publishThreadPools(sh.parent.executorPool)
	if err != nil {
		return dataSet, err
	}
	dataSet = append(dataSet, threadPoolsBytes)
	return dataSet, err
}
func (sh *StreamHandler) getSettings(name string) *Settings {
	return sh.parent.parent.parent.getSettings(name)
}
func (sh *StreamHandler) publishMetrics(cb *CircuitBreaker) ([]byte, error) {
	now := time.Now()
	reqCount := cb.metrics.Requests().Sum(now)
	errCount := cb.metrics.DefaultCollector().Errors().Sum(now)
	errPct := cb.metrics.ErrorPercent(now)

	eventBytes, err := json.Marshal(&streamCmdMetric{
		Type:           "HystrixCommand",
		Name:           cb.Name,
		Group:          cb.Name,
		Time:           currentTime(),
		ReportingHosts: 1,

		RequestCount:       uint32(reqCount),
		ErrorCount:         uint32(errCount),
		ErrorPct:           uint32(errPct),
		CircuitBreakerOpen: func() bool { isOpen, _ := cb.isOpen(); return isOpen }(),

		RollingCountSuccess:            uint32(cb.metrics.DefaultCollector().Successes().Sum(now)),
		RollingCountFailure:            uint32(cb.metrics.DefaultCollector().Failures().Sum(now)),
		RollingCountThreadPoolRejected: uint32(cb.metrics.DefaultCollector().Rejects().Sum(now)),
		RollingCountShortCircuited:     uint32(cb.metrics.DefaultCollector().ShortCircuits().Sum(now)),
		RollingCountTimeout:            uint32(cb.metrics.DefaultCollector().Timeouts().Sum(now)),
		RollingCountFallbackSuccess:    uint32(cb.metrics.DefaultCollector().FallbackSuccesses().Sum(now)),
		RollingCountFallbackFailure:    uint32(cb.metrics.DefaultCollector().FallbackFailures().Sum(now)),

		LatencyTotal:       generateLatencyTimings(cb.metrics.DefaultCollector().TotalDuration()),
		LatencyTotalMean:   cb.metrics.DefaultCollector().TotalDuration().Mean(),
		LatencyExecute:     generateLatencyTimings(cb.metrics.DefaultCollector().RunDuration()),
		LatencyExecuteMean: cb.metrics.DefaultCollector().RunDuration().Mean(),

		RollingStatsWindow:         10000,
		ExecutionIsolationStrategy: "THREAD",

		CircuitBreakerEnabled:                true,
		CircuitBreakerForceClosed:            false,
		CircuitBreakerForceOpen:              cb.forceOpen,
		CircuitBreakerErrorThresholdPercent:  uint32(sh.getSettings(cb.Name).ErrorPercentThreshold),
		CircuitBreakerSleepWindow:            uint32(sh.getSettings(cb.Name).SleepWindow.Seconds() * 1000),
		CircuitBreakerRequestVolumeThreshold: uint32(sh.getSettings(cb.Name).RequestVolumeThreshold),
	})
	return eventBytes, err
}

func (sh *StreamHandler) publishThreadPools(pool *executorPool) ([]byte, error) {
	now := time.Now()

	eventBytes, err := json.Marshal(&streamThreadPoolMetric{
		Type:           "HystrixThreadPool",
		Name:           pool.Name,
		ReportingHosts: 1,

		CurrentActiveCount:        uint32(pool.ActiveCount()),
		CurrentTaskCount:          0,
		CurrentCompletedTaskCount: 0,

		RollingCountThreadsExecuted: uint32(pool.Metrics.Executed.Sum(now)),
		RollingMaxActiveThreads:     uint32(pool.Metrics.MaxActiveRequests.Max(now)),

		CurrentPoolSize:        uint32(pool.Max),
		CurrentCorePoolSize:    uint32(pool.Max),
		CurrentLargestPoolSize: uint32(pool.Max),
		CurrentMaximumPoolSize: uint32(pool.Max),

		RollingStatsWindow:          10000,
		QueueSizeRejectionThreshold: 0,
		CurrentQueueSize:            0,
	})
	return eventBytes, err
}

func generateLatencyTimings(r *rolling.Timing) streamCmdLatency {
	return streamCmdLatency{
		Timing0:   r.Percentile(0),
		Timing25:  r.Percentile(25),
		Timing50:  r.Percentile(50),
		Timing75:  r.Percentile(75),
		Timing90:  r.Percentile(90),
		Timing95:  r.Percentile(95),
		Timing99:  r.Percentile(99),
		Timing995: r.Percentile(99.5),
		Timing100: r.Percentile(100),
	}
}

type streamCmdMetric struct {
	Type           string `json:"type"`
	Name           string `json:"name"`
	Group          string `json:"group"`
	Time           int64  `json:"currentTime"`
	ReportingHosts uint32 `json:"reportingHosts"`

	RequestCount       uint32 `json:"requestCount"`
	ErrorCount         uint32 `json:"errorCount"`
	ErrorPct           uint32 `json:"errorPercentage"`
	CircuitBreakerOpen bool   `json:"isCircuitBreakerOpen"`

	RollingCountCollapsedRequests  uint32 `json:"rollingCountCollapsedRequests"`
	RollingCountExceptionsThrown   uint32 `json:"rollingCountExceptionsThrown"`
	RollingCountFailure            uint32 `json:"rollingCountFailure"`
	RollingCountFallbackFailure    uint32 `json:"rollingCountFallbackFailure"`
	RollingCountFallbackRejection  uint32 `json:"rollingCountFallbackRejection"`
	RollingCountFallbackSuccess    uint32 `json:"rollingCountFallbackSuccess"`
	RollingCountResponsesFromCache uint32 `json:"rollingCountResponsesFromCache"`
	RollingCountSemaphoreRejected  uint32 `json:"rollingCountSemaphoreRejected"`
	RollingCountShortCircuited     uint32 `json:"rollingCountShortCircuited"`
	RollingCountSuccess            uint32 `json:"rollingCountSuccess"`
	RollingCountThreadPoolRejected uint32 `json:"rollingCountThreadPoolRejected"`
	RollingCountTimeout            uint32 `json:"rollingCountTimeout"`

	CurrentConcurrentExecutionCount uint32 `json:"currentConcurrentExecutionCount"`

	LatencyExecuteMean uint32           `json:"latencyExecute_mean"`
	LatencyExecute     streamCmdLatency `json:"latencyExecute"`
	LatencyTotalMean   uint32           `json:"latencyTotal_mean"`
	LatencyTotal       streamCmdLatency `json:"latencyTotal"`

	CircuitBreakerRequestVolumeThreshold             uint32 `json:"propertyValue_circuitBreakerRequestVolumeThreshold"`
	CircuitBreakerSleepWindow                        uint32 `json:"propertyValue_circuitBreakerSleepWindowInMilliseconds"`
	CircuitBreakerErrorThresholdPercent              uint32 `json:"propertyValue_circuitBreakerErrorThresholdPercentage"`
	CircuitBreakerForceOpen                          bool   `json:"propertyValue_circuitBreakerForceOpen"`
	CircuitBreakerForceClosed                        bool   `json:"propertyValue_circuitBreakerForceClosed"`
	CircuitBreakerEnabled                            bool   `json:"propertyValue_circuitBreakerEnabled"`
	ExecutionIsolationStrategy                       string `json:"propertyValue_executionIsolationStrategy"`
	ExecutionIsolationThreadTimeout                  uint32 `json:"propertyValue_executionIsolationThreadTimeoutInMilliseconds"`
	ExecutionIsolationThreadInterruptOnTimeout       bool   `json:"propertyValue_executionIsolationThreadInterruptOnTimeout"`
	ExecutionIsolationThreadPoolKeyOverride          string `json:"propertyValue_executionIsolationThreadPoolKeyOverride"`
	ExecutionIsolationSemaphoreMaxConcurrentRequests uint32 `json:"propertyValue_executionIsolationSemaphoreMaxConcurrentRequests"`
	FallbackIsolationSemaphoreMaxConcurrentRequests  uint32 `json:"propertyValue_fallbackIsolationSemaphoreMaxConcurrentRequests"`
	RollingStatsWindow                               uint32 `json:"propertyValue_metricsRollingStatisticalWindowInMilliseconds"`
	RequestCacheEnabled                              bool   `json:"propertyValue_requestCacheEnabled"`
	RequestLogEnabled                                bool   `json:"propertyValue_requestLogEnabled"`
}

type streamCmdLatency struct {
	Timing0   uint32 `json:"0"`
	Timing25  uint32 `json:"25"`
	Timing50  uint32 `json:"50"`
	Timing75  uint32 `json:"75"`
	Timing90  uint32 `json:"90"`
	Timing95  uint32 `json:"95"`
	Timing99  uint32 `json:"99"`
	Timing995 uint32 `json:"99.5"`
	Timing100 uint32 `json:"100"`
}

type streamThreadPoolMetric struct {
	Type           string `json:"type"`
	Name           string `json:"name"`
	ReportingHosts uint32 `json:"reportingHosts"`

	CurrentActiveCount        uint32 `json:"currentActiveCount"`
	CurrentCompletedTaskCount uint32 `json:"currentCompletedTaskCount"`
	CurrentCorePoolSize       uint32 `json:"currentCorePoolSize"`
	CurrentLargestPoolSize    uint32 `json:"currentLargestPoolSize"`
	CurrentMaximumPoolSize    uint32 `json:"currentMaximumPoolSize"`
	CurrentPoolSize           uint32 `json:"currentPoolSize"`
	CurrentQueueSize          uint32 `json:"currentQueueSize"`
	CurrentTaskCount          uint32 `json:"currentTaskCount"`

	RollingMaxActiveThreads     uint32 `json:"rollingMaxActiveThreads"`
	RollingCountThreadsExecuted uint32 `json:"rollingCountThreadsExecuted"`

	RollingStatsWindow          uint32 `json:"propertyValue_metricsRollingStatisticalWindowInMilliseconds"`
	QueueSizeRejectionThreshold uint32 `json:"propertyValue_queueSizeRejectionThreshold"`
}

func currentTime() int64 {
	return time.Now().UnixNano() / int64(1000000)
}
