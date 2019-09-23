package curator

import (
	"sync"

	"github.com/samuel/go-zookeeper/zk"
)

type Watcher interface {
	process(event *zk.Event)
}

type simpleWatcher struct {
	Func func(event *zk.Event)
}

func NewWatcher(fn func(event *zk.Event)) Watcher {
	return &simpleWatcher{fn}
}

func (w *simpleWatcher) process(event *zk.Event) {
	w.Func(event)
}

type Watchers struct {
	lock     sync.Mutex
	watchers []Watcher
}

func NewWatchers(watchers ...Watcher) *Watchers {
	return &Watchers{watchers: watchers}
}

func (w *Watchers) Len() int { return len(w.watchers) }

func (w *Watchers) Add(watcher Watcher) Watcher {
	w.lock.Lock()

	w.watchers = append(w.watchers, watcher)

	w.lock.Unlock()

	return watcher
}

func (w *Watchers) Remove(watcher Watcher) Watcher {
	w.lock.Lock()
	defer w.lock.Unlock()

	for i, v := range w.watchers {
		if v == watcher {
			w.watchers = append(w.watchers[:i], w.watchers[i+1:]...)

			return watcher
		}
	}

	return nil
}

func (w *Watchers) Fire(event *zk.Event) {
	for _, watcher := range w.watchers {
		if watcher != nil {
			go watcher.process(event)
		}
	}
}

func (w *Watchers) Watch(events <-chan zk.Event) {
	for {
		if event, ok := <-events; !ok {
			break
		} else {
			w.Fire(&event)
		}
	}
}
