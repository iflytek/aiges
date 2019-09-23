package recipes

import (
	"fmt"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/curator-go/curator"
	"github.com/samuel/go-zookeeper/zk"
)

const LockPrefix = "lock-"

type InterProcessLock interface {
	// Acquire the mutex - blocking until it's available.
	// Each call to acquire must be balanced by a call to Release()
	Acquire() (bool, error)

	// Acquire the mutex - blocks until it's available or the given time expires.
	AcquireTimeout(expires time.Duration) (bool, error)

	// Perform one release of the mutex.
	Release() error

	// Returns true if the mutex is acquired by a go-routine in this process
	IsAcquiredInThisProcess() bool
}

type RevocationListener interface {
	// Called when a revocation request has been received.
	// You should release the lock as soon as possible. Revocation is cooperative.
	RevocationRequested(forLock InterProcessMutex)
}

// Specifies locks that can be revoked
type Revocable interface {
	// Make the lock revocable.
	// Your listener will get called when another process/thread wants you to release the lock. Revocation is cooperative.
	MakeRevocable(listener RevocationListener)
}

type LockInternalsSorter interface {
	FixForSorting(str, lockName string) string
}

type PredicateResults struct {
	GetsTheLock bool
	PathToWatch string
}

type LockInternalsDriver interface {
	LockInternalsSorter

	GetsTheLock(client curator.CuratorFramework, children []string, sequenceNodeName string, maxLeases int) (*PredicateResults, error)

	CreatesTheLock(client curator.CuratorFramework, path string, lockNodeBytes []byte) (string, error)
}

type StandardLockInternalsDriver struct{}

func NewStandardLockInternalsDriver() *StandardLockInternalsDriver {
	return &StandardLockInternalsDriver{}
}

func (d *StandardLockInternalsDriver) FixForSorting(str, lockName string) string {
	if idx := strings.LastIndex(str, lockName); idx >= 0 {
		idx += len(lockName)

		if idx <= len(str) {
			return str[idx:]
		} else {
			return ""
		}
	}

	return str
}

func (d *StandardLockInternalsDriver) GetsTheLock(client curator.CuratorFramework, children []string, sequenceNodeName string, maxLeases int) (*PredicateResults, error) {
	for i, child := range children {
		if child == sequenceNodeName {
			var pathToWatch string

			getsTheLock := i < maxLeases

			if !getsTheLock {
				pathToWatch = children[i-maxLeases]
			}

			return &PredicateResults{GetsTheLock: getsTheLock, PathToWatch: pathToWatch}, nil
		}
	}

	return nil, zk.ErrNoNode
}

func (d *StandardLockInternalsDriver) CreatesTheLock(client curator.CuratorFramework, path string, lockNodeBytes []byte) (string, error) {
	if lockNodeBytes == nil {
		return client.Create().CreatingParentsIfNeeded().WithMode(curator.EPHEMERAL_SEQUENTIAL).ForPath(path)
	} else {
		return client.Create().CreatingParentsIfNeeded().WithMode(curator.EPHEMERAL_SEQUENTIAL).ForPathWithData(path, lockNodeBytes)
	}
}

// A re-entrant mutex that works across processes. Uses Zookeeper to hold the lock.
// All processes that use the same lock path will achieve an inter-process critical section.
// Further, this mutex is "fair" - each user will get the mutex in the order requested (from ZK's point of view)
type InterProcessMutex struct {
	basePath      string
	internals     *lockInternals
	lockPath      string
	lockCount     int32
	LockNodeBytes []byte
}

func NewInterProcessMutex(client curator.CuratorFramework, path string) (*InterProcessMutex, error) {
	return NewInterProcessMutexWithDriver(client, path, NewStandardLockInternalsDriver())
}

func NewInterProcessMutexWithDriver(client curator.CuratorFramework, path string, driver LockInternalsDriver) (*InterProcessMutex, error) {
	if err := curator.ValidatePath(path); err != nil {
		return nil, err
	}

	if internals, err := newLockInternals(client, driver, path, LockPrefix, 1); err != nil {
		return nil, err
	} else {
		return &InterProcessMutex{
			basePath:  path,
			internals: internals,
		}, nil
	}
}

func (m *InterProcessMutex) Acquire() (bool, error) {
	if locked, err := m.internalLock(-1); err != nil {
		return false, err
	} else if !locked {
		return false, fmt.Errorf("Lost connection while trying to acquire lock: %s", m.basePath)
	} else {
		return true, err
	}
}

func (m *InterProcessMutex) AcquireTimeout(expires time.Duration) (bool, error) {
	return m.internalLock(expires)
}

func (m *InterProcessMutex) Release() error {
	if !m.IsAcquiredInThisProcess() {
		return fmt.Errorf("You do not own the lock: %s", m.basePath)
	}

	count := atomic.AddInt32(&m.lockCount, -1)

	switch {
	case count > 0:
		return nil
	case count < 0:
		return fmt.Errorf("Lock count has gone negative for lock: %s", m.basePath)
	default:
		return m.internals.releaseLock(m.lockPath)
	}
}

func (m *InterProcessMutex) IsAcquiredInThisProcess() bool {
	return atomic.LoadInt32(&m.lockCount) > 0
}

func (m *InterProcessMutex) internalLock(expires time.Duration) (bool, error) {
	if m.IsAcquiredInThisProcess() {
		// re-entering
		atomic.AddInt32(&m.lockCount, 1)

		return true, nil
	}

	if lockPath, err := m.internals.attemptLock(expires, m.LockNodeBytes); err != nil {
		return false, err
	} else if len(lockPath) > 0 {
		m.lockPath = lockPath

		atomic.StoreInt32(&m.lockCount, 1)

		return true, nil
	}

	return false, nil
}

type lockInternals struct {
	client    curator.CuratorFramework
	driver    LockInternalsDriver
	basePath  string
	lockName  string
	lockPath  string
	maxLeases int
}

func newLockInternals(client curator.CuratorFramework, driver LockInternalsDriver, basePath, lockName string, maxLeases int) (*lockInternals, error) {
	if err := curator.ValidatePath(basePath); err != nil {
		return nil, err
	}

	return &lockInternals{
		client:    client,
		driver:    driver,
		basePath:  basePath,
		lockName:  lockName,
		lockPath:  curator.JoinPath(basePath, lockName),
		maxLeases: maxLeases,
	}, nil
}

func (l *lockInternals) attemptLock(waitTime time.Duration, lockNodeBytes []byte) (string, error) {
	startTime := time.Now()
	retryCount := 0

	for {
		var ourPath string
		var err error

		if ourPath, err = l.driver.CreatesTheLock(l.client, l.lockPath, lockNodeBytes); err == nil {
			if hasTheLock, err := l.internalLockLoop(startTime, waitTime, ourPath); err == nil {
				if hasTheLock {
					return ourPath, nil
				} else {
					return "", nil
				}
			}
		}

		if err == zk.ErrNoNode {
			retryCount++

			if l.client.ZookeeperClient().RetryPolicy().AllowRetry(retryCount, time.Now().Sub(startTime), curator.DefaultRetrySleeper) {
				continue
			}
		}

		if err != nil {
			return "", err
		}
	}
}

func (l *lockInternals) releaseLock(path string) error {
	return l.deleteOurPath(path)
}

func (l *lockInternals) deleteOurPath(path string) error {
	if err := l.client.Delete().ForPath(path); err == zk.ErrNoNode {
		return nil // ignore - already deleted (possibly expired session, etc.)
	} else {
		return err
	}
}

func (l *lockInternals) internalLockLoop(startTime time.Time, waitTime time.Duration, path string) (haveTheLock bool, err error) {
	var doDelete bool

	for l.client.State() == curator.STARTED && !haveTheLock {
		if children, err := l.getSortedChildren(); err != nil {
			break
		} else {
			sequenceNodeName := path[len(l.basePath)+1:]

			if results, err := l.driver.GetsTheLock(l.client, children, sequenceNodeName, l.maxLeases); err != nil {
				break
			} else if results.GetsTheLock {
				haveTheLock = true

				break
			} else {
				previousSequencePath := curator.JoinPath(l.basePath, results.PathToWatch)

				c := make(chan error)

				t := time.NewTimer(waitTime - time.Now().Sub(startTime))

				l.client.GetData().UsingWatcher(curator.NewWatcher(func(event *zk.Event) {
					c <- event.Err
				})).ForPath(previousSequencePath)

				select {
				case err := <-c:
					if err != nil && err != zk.ErrNoNode {
						break
					}
				case <-t.C:
				}
			}
		}
	}

	if err != nil || doDelete {
		l.deleteOurPath(path)
	}

	return haveTheLock, err
}

type ChildrenSorter struct {
	children []string
	less     func(lhs, rhs string) bool
}

func (s ChildrenSorter) Len() int {
	return len(s.children)
}

func (s ChildrenSorter) Less(i, j int) bool {
	return s.less(s.children[i], s.children[j])
}

func (s ChildrenSorter) Swap(i, j int) { s.children[i], s.children[j] = s.children[j], s.children[i] }

func (l *lockInternals) getSortedChildren() ([]string, error) {
	if children, err := l.client.GetChildren().ForPath(l.basePath); err != nil {
		return nil, err
	} else {
		sort.Sort(ChildrenSorter{children, func(lhs, rhs string) bool {
			return l.driver.FixForSorting(lhs, l.lockName) < l.driver.FixForSorting(rhs, l.lockName)
		}})

		return children, nil
	}
}
