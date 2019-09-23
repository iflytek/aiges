package curator

import (
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRetryLoop(t *testing.T) {
	d := 3 * time.Second
	p := NewRetryNTimes(3, d)
	sleeper := &mockRetrySleeper{}
	tracer := &mockTracerDriver{}

	retryLoop := newRetryLoop(p, tracer)

	assert.NotNil(t, retryLoop)
	assert.Equal(t, 0, retryLoop.retryCount)

	retryLoop.retrySleeper = sleeper

	errors := []error{zk.ErrSessionExpired, zk.ErrSessionMoved, nil}

	sleeper.On("SleepFor", d).Return(nil).Times(2)
	tracer.On("AddCount", "retries-allowed", 1).Return().Twice()

	_, err := retryLoop.CallWithRetry(func() (interface{}, error) {
		return nil, errors[retryLoop.retryCount]
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, retryLoop.retryCount)

	sleeper.AssertExpectations(t)
	tracer.AssertExpectations(t)

	// retry loop failed
	retryLoop = newRetryLoop(p, nil)

	_, err = retryLoop.CallWithRetry(func() (interface{}, error) {
		return nil, zk.ErrClosing
	})

	assert.EqualError(t, err, zk.ErrClosing.Error())
}

func TestRetryNTimes(t *testing.T) {
	d := 3 * time.Second
	p := NewRetryNTimes(3, d)
	s := &mockRetrySleeper{}

	assert.NotNil(t, p)

	s.On("SleepFor", d).Return(nil).Times(3)

	assert.True(t, p.AllowRetry(0, 0, s))
	assert.True(t, p.AllowRetry(1, 0, s))
	assert.True(t, p.AllowRetry(2, 0, s))
	assert.False(t, p.AllowRetry(3, 0, s))

	s.AssertExpectations(t)
}

func TestRetryOneTime(t *testing.T) {
	d := 3 * time.Second
	p := NewRetryOneTime(d)
	s := &mockRetrySleeper{}

	assert.NotNil(t, p)

	s.On("SleepFor", d).Return(nil).Once()

	assert.True(t, p.AllowRetry(0, 0, s))
	assert.False(t, p.AllowRetry(1, 0, s))

	s.AssertExpectations(t)
}

func TestExponentialBackoffRetry(t *testing.T) {
	d := 3 * time.Second
	p := NewExponentialBackoffRetry(d, 3, 9*time.Second)
	s := &mockRetrySleeper{}

	assert.NotNil(t, p)

	s.On("SleepFor", mock.AnythingOfType("Duration")).Return(nil).Times(3)

	assert.True(t, p.AllowRetry(0, 0, s))
	assert.True(t, p.AllowRetry(1, 0, s))
	assert.True(t, p.AllowRetry(2, 0, s))
	assert.False(t, p.AllowRetry(3, 0, s))

	assert.True(t, s.Calls[0].Arguments.Get(0).(time.Duration) < 1*d)
	assert.True(t, s.Calls[1].Arguments.Get(0).(time.Duration) < 2*d)
	assert.True(t, s.Calls[2].Arguments.Get(0).(time.Duration) < 4*d)

	s.AssertExpectations(t)
}

func TestRetryUntilElapsed(t *testing.T) {
	d := 3 * time.Second
	p := NewRetryUntilElapsed(3*d, d)
	s := &mockRetrySleeper{}

	assert.NotNil(t, p)

	s.On("SleepFor", d).Return(nil).Times(3)

	assert.True(t, p.AllowRetry(0, 0, s))
	assert.True(t, p.AllowRetry(0, d*1, s))
	assert.True(t, p.AllowRetry(0, d*2, s))
	assert.False(t, p.AllowRetry(0, d*3, s))

	s.AssertExpectations(t)
}
