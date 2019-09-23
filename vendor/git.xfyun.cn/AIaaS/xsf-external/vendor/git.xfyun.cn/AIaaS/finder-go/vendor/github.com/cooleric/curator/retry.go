package curator

import (
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/cooleric/go-zookeeper/zk"
)

// Abstraction for retry policies to sleep
type RetrySleeper interface {
	// Sleep for the given time
	SleepFor(time time.Duration) error
}

// Abstracts the policy to use when retrying connections
type RetryPolicy interface {
	// Called when an operation has failed for some reason.
	// This method should return true to make another attempt.
	AllowRetry(retryCount int, elapsedTime time.Duration, sleeper RetrySleeper) bool
}

type defaultRetrySleeper struct {
}

var DefaultRetrySleeper RetrySleeper = &defaultRetrySleeper{}

func (s *defaultRetrySleeper) SleepFor(d time.Duration) error {
	time.Sleep(d)

	return nil
}

// Mechanism to perform an operation on Zookeeper that is safe against disconnections and "recoverable" errors.
type RetryLoop interface {
	// creates a retry loop calling the given proc and retrying if needed
	CallWithRetry(proc func() (interface{}, error)) (interface{}, error)
}

type retryLoop struct {
	done         bool
	retryCount   int
	startTime    time.Time
	retryPolicy  RetryPolicy
	retrySleeper RetrySleeper
	tracer       TracerDriver
}

func newRetryLoop(retryPolicy RetryPolicy, tracer TracerDriver) *retryLoop {
	return &retryLoop{
		startTime:   time.Now(),
		retryPolicy: retryPolicy,
		tracer:      tracer,
	}
}

// return true if the given Zookeeper result code is retry-able
func (l *retryLoop) ShouldRetry(err error) bool {
	if err == zk.ErrSessionExpired || err == zk.ErrSessionMoved {
		return true
	}

	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout() || netErr.Temporary()
	}

	return false
}

func (l *retryLoop) CallWithRetry(proc func() (interface{}, error)) (interface{}, error) {
	for {
		if ret, err := proc(); err == nil || !l.ShouldRetry(err) {
			return ret, err
		} else {
			l.retryCount++

			if sleeper := l.retrySleeper; sleeper == nil {
				sleeper = DefaultRetrySleeper
			} else {
				if !l.retryPolicy.AllowRetry(l.retryCount, time.Now().Sub(l.startTime), sleeper) {
					l.tracer.AddCount("retries-disallowed", 1)

					return ret, err
				} else {
					l.tracer.AddCount("retries-allowed", 1)
				}
			}
		}
	}

	return nil, nil
}

type SleepingRetry struct {
	RetryPolicy

	N            int
	getSleepTime func(retryCount int, elapsedTime time.Duration) time.Duration
}

func (r *SleepingRetry) AllowRetry(retryCount int, elapsedTime time.Duration, sleeper RetrySleeper) bool {
	if retryCount < r.N {
		if err := sleeper.SleepFor(r.getSleepTime(retryCount, elapsedTime)); err != nil {
			return false
		}

		return true
	}

	return false
}

// Retry policy that retries a max number of times
type RetryNTimes struct {
	SleepingRetry
}

func NewRetryNTimes(n int, sleepBetweenRetries time.Duration) *RetryNTimes {
	return &RetryNTimes{
		SleepingRetry: SleepingRetry{
			N:            n,
			getSleepTime: func(retryCount int, elapsedTime time.Duration) time.Duration { return sleepBetweenRetries },
		},
	}
}

// A retry policy that retries only once
type RetryOneTime struct {
	RetryNTimes
}

func NewRetryOneTime(sleepBetweenRetry time.Duration) *RetryOneTime {
	return &RetryOneTime{
		*NewRetryNTimes(1, sleepBetweenRetry),
	}
}

const (
	MAX_RETRIES_LIMIT               = 29
	DEFAULT_MAX_SLEEP time.Duration = time.Duration(math.MaxInt32 * int64(time.Second))
)

// Retry policy that retries a set number of times with increasing sleep time between retries
type ExponentialBackoffRetry struct {
	SleepingRetry
}

func NewExponentialBackoffRetry(baseSleepTime time.Duration, maxRetries int, maxSleep time.Duration) *ExponentialBackoffRetry {
	if maxRetries > MAX_RETRIES_LIMIT {
		maxRetries = MAX_RETRIES_LIMIT
	}

	return &ExponentialBackoffRetry{
		SleepingRetry: SleepingRetry{
			N: maxRetries,
			getSleepTime: func(retryCount int, elapsedTime time.Duration) time.Duration {
				sleepTime := time.Duration(int64(baseSleepTime) * rand.Int63n(1<<uint(retryCount)))

				if sleepTime > maxSleep {
					sleepTime = maxSleep
				}

				return sleepTime
			},
		}}
}

// A retry policy that retries until a given amount of time elapses
type RetryUntilElapsed struct {
	SleepingRetry

	maxElapsedTime time.Duration
}

func NewRetryUntilElapsed(maxElapsedTime, sleepBetweenRetries time.Duration) *RetryUntilElapsed {
	return &RetryUntilElapsed{
		SleepingRetry: SleepingRetry{
			N:            math.MaxInt64,
			getSleepTime: func(retryCount int, elapsedTime time.Duration) time.Duration { return sleepBetweenRetries },
		},
		maxElapsedTime: maxElapsedTime,
	}
}

func (r *RetryUntilElapsed) AllowRetry(retryCount int, elapsedTime time.Duration, sleeper RetrySleeper) bool {
	return elapsedTime < r.maxElapsedTime && r.SleepingRetry.AllowRetry(retryCount, elapsedTime, sleeper)
}
