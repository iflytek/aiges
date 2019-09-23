package set

import "sync"

var pool = sync.Pool{}

type Set struct {
	items     map[interface{}]struct{}
	lock      sync.RWMutex
	flattened []interface{}
}

func (set *Set) Add(items ...interface{}) {
	set.lock.Lock()
	defer set.lock.Unlock()

	set.flattened = nil
	for _, item := range items {
		set.items[item] = struct{}{}
	}
}

func (set *Set) Remove(items ...interface{}) {
	set.lock.Lock()
	defer set.lock.Unlock()

	set.flattened = nil
	for _, item := range items {
		delete(set.items, item)
	}
}

func (set *Set) Flatten() []interface{} {
	set.lock.Lock()
	defer set.lock.Unlock()

	if nil != set.flattened {
		return set.flattened
	}

	set.flattened = make([]interface{}, 0, len(set.items))
	for item := range set.items {
		set.flattened = append(set.flattened, item)
	}
	return set.flattened
}

func (set *Set) Exists(item interface{}) bool {
	set.lock.RLock()

	_, ok := set.items[item]

	set.lock.RUnlock()

	return ok
}

func (set *Set) Len() int64 {
	set.lock.RLock()

	size := int64(len(set.items))

	set.lock.RUnlock()

	return size
}

func (set *Set) Clear() {
	set.lock.Lock()

	set.items = map[interface{}]struct{}{}

	set.lock.Unlock()
}

func New(items ...interface{}) *Set {
	set := pool.Get().(*Set)
	for _, item := range items {
		set.items[item] = struct{}{}
	}

	if len(items) > 0 {
		set.flattened = nil
	}

	return set
}

func init() {
	pool.New = func() interface{} {
		return &Set{
			items: make(map[interface{}]struct{}, 10),
		}
	}
}
