package curator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloseQuietly(t *testing.T) {
	// No Error
	c := &mockCloseable{}

	c.On("Close").Return(nil).Once()

	assert.Nil(t, CloseQuietly(c))

	c.AssertExpectations(t)

	// Return Error
	c = &mockCloseable{}

	c.On("Close").Return(errors.New("test")).Once()

	assert.EqualError(t, CloseQuietly(c), "test")

	c.AssertExpectations(t)

	// Panic
	c = &mockCloseable{crash: true}

	assert.EqualError(t, CloseQuietly(c), "panic")

	c.AssertNotCalled(t, "Close")
}

func TestAtomicBool(t *testing.T) {
	var b AtomicBool

	assert.False(t, b.Load())

	b = NewAtomicBool(true)

	assert.True(t, b.Load())

	assert.False(t, b.CompareAndSwap(false, true))

	assert.True(t, b.Swap(false))

	assert.False(t, b.Load())

	assert.True(t, b.CompareAndSwap(false, true))

	assert.True(t, b.Load())

	assert.True(t, b.Swap(false))
	assert.False(t, b.Swap(true))
	assert.True(t, b.Load())
}

func TestState(t *testing.T) {
	var s State

	assert.Equal(t, LATENT, s)
	assert.Equal(t, LATENT, s.Value())

	assert.True(t, s.Change(LATENT, STARTED))

	assert.Equal(t, STARTED, s)

	assert.False(t, s.Change(LATENT, STARTED))

	defer func() {
		assert.Equal(t, "panic", recover())
	}()

	s.Check(STOPPED, "panic")

	t.Fail()
}
