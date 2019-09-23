package curator

import (
	"github.com/samuel/go-zookeeper/zk"
)

type checkExistsBuilder struct {
	client        *curatorFramework
	backgrounding backgrounding
	watching      watching
}

func (b *checkExistsBuilder) ForPath(givenPath string) (*zk.Stat, error) {
	adjustedPath := b.client.fixForNamespace(givenPath, false)

	if b.backgrounding.inBackground {
		go b.pathInBackground(adjustedPath)

		return nil, nil
	} else {
		return b.pathInForeground(adjustedPath)
	}
}

func (b *checkExistsBuilder) pathInBackground(path string) {
	tracer := b.client.ZookeeperClient().StartTracer("checkExistsBuilder.pathInBackground")

	defer tracer.Commit()

	stat, err := b.pathInForeground(path)

	if b.backgrounding.callback != nil {
		event := &curatorEvent{
			eventType: EXISTS,
			err:       err,
			path:      b.client.unfixForNamespace(path),
			stat:      stat,
			name:      GetNodeFromPath(path),
			context:   b.backgrounding.context,
		}

		b.backgrounding.callback(b.client, event)
	}
}

func (b *checkExistsBuilder) pathInForeground(path string) (*zk.Stat, error) {
	zkClient := b.client.ZookeeperClient()

	result, err := zkClient.NewRetryLoop().CallWithRetry(func() (interface{}, error) {
		if conn, err := zkClient.Conn(); err != nil {
			return nil, err
		} else {
			var exists bool
			var stat *zk.Stat
			var events <-chan zk.Event
			var err error

			if b.watching.watched || b.watching.watcher != nil {
				exists, stat, events, err = conn.ExistsW(path)

				if events != nil && b.watching.watcher != nil {
					go NewWatchers(b.watching.watcher).Watch(events)
				}
			} else {
				exists, stat, err = conn.Exists(path)
			}

			if err != nil {
				return nil, err
			} else if !exists {
				return nil, nil
			} else {
				return stat, nil
			}
		}
	})

	stat, _ := result.(*zk.Stat)

	return stat, err
}

func (b *checkExistsBuilder) Watched() CheckExistsBuilder {
	b.watching.watched = true

	return b
}

func (b *checkExistsBuilder) UsingWatcher(watcher Watcher) CheckExistsBuilder {
	b.watching.watcher = b.client.getNamespaceWatcher(watcher)

	return b
}

func (b *checkExistsBuilder) InBackground() CheckExistsBuilder {
	b.backgrounding = backgrounding{inBackground: true}

	return b
}

func (b *checkExistsBuilder) InBackgroundWithContext(context interface{}) CheckExistsBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context}

	return b
}

func (b *checkExistsBuilder) InBackgroundWithCallback(callback BackgroundCallback) CheckExistsBuilder {
	b.backgrounding = backgrounding{inBackground: true, callback: callback}

	return b
}

func (b *checkExistsBuilder) InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) CheckExistsBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context, callback: callback}

	return b
}
