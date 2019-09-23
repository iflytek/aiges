package curator

import (
	"log"
	"sync/atomic"
	"unsafe"
)

// A Closeable is a source or destination of data that can be closed.
type Closeable interface {
	// Closes this and releases any system resources associated with it.
	Close() error
}

func CloseQuietly(closeable Closeable) (err error) {
	defer func() {
		if v := recover(); v != nil {
			log.Printf("panic when closing %s, %v", closeable, v)

			err, _ = v.(error)
		}
	}()

	if err = closeable.Close(); err != nil {
		log.Printf("fail to close %s, %s", closeable, err)
	}

	return
}

type AtomicBool int32

const (
	FALSE AtomicBool = iota
	TRUE
)

func NewAtomicBool(b bool) AtomicBool {
	if b {
		return TRUE
	}

	return FALSE
}

func (b *AtomicBool) CompareAndSwap(oldValue, newValue bool) bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(b)), int32(NewAtomicBool(oldValue)), int32(NewAtomicBool(newValue)))
}

func (b *AtomicBool) Load() bool {
	return atomic.LoadInt32((*int32)(unsafe.Pointer(b))) != int32(FALSE)
}

func (b *AtomicBool) Swap(v bool) bool {
	var n AtomicBool

	if v {
		n = TRUE
	}

	return atomic.SwapInt32((*int32)(unsafe.Pointer(b)), int32(n)) != int32(FALSE)
}

func (b *AtomicBool) Set(v bool) { b.Swap(v) }

type State int32

const (
	LATENT  State = iota // Start() has not yet been called
	STARTED              // Start() has been called
	STOPPED              // Close() has been called
)

func (s *State) Change(oldState, newState State) bool {
	return atomic.CompareAndSwapInt32((*int32)(s), int32(oldState), int32(newState))
}

func (s *State) Value() State {
	return State(atomic.LoadInt32((*int32)(s)))
}

func (s State) Check(state State, msg string) {
	if s != state {
		panic(msg)
	}
}
