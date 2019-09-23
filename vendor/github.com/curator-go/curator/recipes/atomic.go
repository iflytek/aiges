package recipes

import (
	"bytes"
	"time"

	"github.com/curator-go/curator"
	"github.com/samuel/go-zookeeper/zk"
)

// Debugging stats about operations
type AtomicStats struct {
	//  the number of optimistic locks used to perform the operation
	OptimisticTries int

	// the number of mutex locks used to perform the operation
	PromotedTries int

	// the time spent trying the operation with optimistic locks
	OptimisticTime time.Duration

	// the time spent trying the operation with mutex locks
	PromotedTime time.Duration
}

// Abstracts a value returned from one of the Atomics
type AtomicValue interface {
	// MUST be checked.
	// Returns true if the operation succeeded. If false is returned,
	// the operation failed and the atomic was not updated.
	Succeeded() bool

	// Returns the value of the counter prior to the operation
	PreValue() []byte

	// Returns the value of the counter after to the operation
	PostValue() []byte

	// Returns debugging stats about the operation
	Stats() *AtomicStats
}

type DistributedAtomicValue interface {
	// Returns the current value of the counter.
	Get() (AtomicValue, error)

	// Atomically sets the value to the given updated value
	// if the current value == the expected value.
	// Remember to always check AtomicValue.Succeeded().
	CompareAndSet(expectedValue, newValue []byte) (AtomicValue, error)

	// Attempt to atomically set the value to the given value.
	// Remember to always check AtomicValue.Succeeded().
	TrySet(newValue []byte) (AtomicValue, error)

	// Forcibly sets the value of the counter without any guarantees of atomicity.
	ForceSet(newValue []byte) error

	// Atomic values are initially set to the equivalent of <code>NULL</code> in a database.
	// Use this method to initialize the value.
	// The value will be set if and only iff the node does not exist.
	Initialize(value []byte) (bool, error)
}

type DistributedAtomicNumber interface {
	// Add 1 to the current value and return the new value information.
	// Remember to always check AtomicValue.Succeeded().
	Increment() (AtomicValue, error)

	// Subtract 1 from the current value and return the new value information.
	// Remember to always check AtomicValue.Succeeded().
	Decrement() (AtomicValue, error)

	// Add delta to the current value and return the new value information.
	// Remember to always check AtomicValue.Succeeded().
	Add(delta []byte) (AtomicValue, error)

	// Subtract delta from the current value and return the new value information.
	// Remember to always check AtomicValue.Succeeded().
	Subtract(delta []byte) (AtomicValue, error)
}

type mutableAtomicValue struct {
	preValue, postValue []byte
	succeeded           bool
	stats               AtomicStats
}

func (v *mutableAtomicValue) Succeeded() bool { return v.succeeded }

func (v *mutableAtomicValue) PreValue() []byte { return v.preValue }

func (v *mutableAtomicValue) PostValue() []byte { return v.postValue }

func (v *mutableAtomicValue) Stats() *AtomicStats { return &v.stats }

type PromotedToLock struct {
	lockPath    string
	maxLockTime time.Duration
	retryPolicy curator.RetryPolicy
}

type distributedAtomicValue struct {
	client         curator.CuratorFramework
	path           string
	retryPolicy    curator.RetryPolicy
	promotedToLock *PromotedToLock
	mutex          InterProcessLock
}

func NewDistributedAtomicValue(client curator.CuratorFramework, path string, retryPolicy curator.RetryPolicy) (DistributedAtomicValue, error) {
	return NewDistributedAtomicValueWithLock(client, path, retryPolicy, nil)
}

func NewDistributedAtomicValueWithLock(client curator.CuratorFramework, path string, retryPolicy curator.RetryPolicy, promotedToLock *PromotedToLock) (DistributedAtomicValue, error) {
	if err := curator.ValidatePath(path); err != nil {
		return nil, err
	}

	v := &distributedAtomicValue{
		client:         client,
		path:           path,
		retryPolicy:    retryPolicy,
		promotedToLock: promotedToLock,
	}

	if promotedToLock != nil {
		if m, err := NewInterProcessMutex(client, promotedToLock.lockPath); err != nil {
			return nil, err
		} else {
			v.mutex = m
		}
	}

	return v, nil
}

func (v *distributedAtomicValue) Get() (AtomicValue, error) {
	var result mutableAtomicValue

	if _, err := v.currentValue(&result, nil); err != nil {
		return nil, err
	}

	result.postValue = result.preValue
	result.succeeded = true

	return &result, nil
}

func (v *distributedAtomicValue) ForceSet(newValue []byte) (err error) {
	if _, err = v.client.SetData().ForPathWithData(v.path, newValue); err == zk.ErrNoNode {
		if _, err = v.client.Create().ForPathWithData(v.path, newValue); err == zk.ErrNodeExists {
			_, err = v.client.SetData().ForPathWithData(v.path, newValue)
		}
	}

	return
}

func (v *distributedAtomicValue) CompareAndSet(expectedValue, newValue []byte) (AtomicValue, error) {
	var result mutableAtomicValue
	var stat zk.Stat

	if createIt, err := v.currentValue(&result, &stat); err != nil {
		return nil, err
	} else if !createIt && bytes.Equal(expectedValue, result.preValue) {
		if _, err := v.client.SetData().WithVersion(stat.Version).ForPathWithData(v.path, newValue); err == nil {
			result.succeeded = true
			result.postValue = newValue
		} else if err == zk.ErrBadVersion || err == zk.ErrNoNode {
			result.succeeded = false
		} else {
			return nil, err
		}
	} else {
		result.succeeded = false
	}

	return &result, nil
}

func (v *distributedAtomicValue) TrySet(newValue []byte) (AtomicValue, error) {
	var result mutableAtomicValue

	if err := v.tryOptimistic(&result, newValue); err != nil {
		return nil, err
	} else if !result.succeeded && v.mutex != nil {
		if err := v.tryWithMutex(&result, newValue); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (v *distributedAtomicValue) Initialize(value []byte) (bool, error) {
	if _, err := v.client.Create().ForPath(v.path); err == nil {
		return true, nil
	} else if err == zk.ErrNodeExists {
		return false, nil
	} else {
		return false, err
	}
}

func (v *distributedAtomicValue) currentValue(result *mutableAtomicValue, stat *zk.Stat) (bool, error) {
	if data, err := v.client.GetData().StoringStatIn(stat).ForPath(v.path); err == nil {
		result.preValue = data

		return false, nil
	} else if err == zk.ErrNoNode {
		result.preValue = nil

		return true, nil
	} else {
		return false, err
	}
}

func (v *distributedAtomicValue) tryOptimistic(result *mutableAtomicValue, newValue []byte) error {
	startTime := time.Now()

	defer func() {
		result.stats.OptimisticTime = time.Now().Sub(startTime)
	}()

	for {
		result.stats.OptimisticTries++

		if success, err := v.tryOnce(result, newValue); err != nil {
			return err
		} else if success {
			result.succeeded = true

			break
		} else if !v.retryPolicy.AllowRetry(result.stats.OptimisticTries, time.Now().Sub(startTime), curator.DefaultRetrySleeper) {
			break
		}
	}

	return nil
}

func (v *distributedAtomicValue) tryOnce(result *mutableAtomicValue, newValue []byte) (bool, error) {
	var stat zk.Stat

	if createIt, err := v.currentValue(result, &stat); err != nil {
		return false, err
	} else {
		var err error

		if createIt {
			_, err = v.client.Create().ForPathWithData(v.path, newValue)
		} else {
			_, err = v.client.SetData().WithVersion(stat.Version).ForPathWithData(v.path, newValue)
		}

		if err == nil {
			result.postValue = newValue

			return true, nil
		} else if err == zk.ErrNodeExists || err == zk.ErrBadVersion || err == zk.ErrNoNode {
			return false, nil
		} else {
			return false, err
		}
	}
}

func (v *distributedAtomicValue) tryWithMutex(result *mutableAtomicValue, newValue []byte) error {
	startTime := time.Now()

	defer func() {
		result.stats.PromotedTime = time.Now().Sub(startTime)
	}()

	if locked, err := v.mutex.AcquireTimeout(v.promotedToLock.maxLockTime); err != nil {
		return err
	} else if locked {
		defer v.mutex.Release()

		for {
			result.stats.PromotedTries++

			if success, err := v.tryOnce(result, newValue); err != nil {
				return err
			} else if success {
				result.succeeded = true

				break
			} else if !v.promotedToLock.retryPolicy.AllowRetry(result.stats.PromotedTries, time.Now().Sub(startTime), curator.DefaultRetrySleeper) {
				break
			}
		}
	}

	return nil
}
