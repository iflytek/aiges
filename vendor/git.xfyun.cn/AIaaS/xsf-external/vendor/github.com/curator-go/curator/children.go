package curator

import (
	"github.com/samuel/go-zookeeper/zk"
)

type getChildrenBuilder struct {
	client        *curatorFramework
	backgrounding backgrounding
	stat          *zk.Stat
	watching      watching
}

func (b *getChildrenBuilder) ForPath(givenPath string) ([]string, error) {
	adjustedPath := b.client.fixForNamespace(givenPath, false)

	if b.backgrounding.inBackground {
		go b.pathInBackground(adjustedPath, givenPath)

		return nil, nil
	}

	if children, err := b.pathInForeground(adjustedPath); err != nil {
		return nil, err
	} else {
		return children, err
	}
}

func (b *getChildrenBuilder) pathInBackground(adjustedPath, givenPath string) {
	tracer := b.client.ZookeeperClient().StartTracer("getChildrenBuilder.pathInBackground")

	defer tracer.Commit()

	children, err := b.pathInForeground(adjustedPath)

	if b.backgrounding.callback != nil {
		event := &curatorEvent{
			eventType: CHILDREN,
			err:       err,
			path:      b.client.unfixForNamespace(adjustedPath),
			children:  children,
			stat:      b.stat,
			context:   b.backgrounding.context,
		}

		if err != nil {
			event.path = givenPath
		}

		event.name = GetNodeFromPath(event.path)

		b.backgrounding.callback(b.client, event)
	}
}

func (b *getChildrenBuilder) pathInForeground(path string) ([]string, error) {
	zkClient := b.client.ZookeeperClient()

	result, err := zkClient.NewRetryLoop().CallWithRetry(func() (interface{}, error) {
		if conn, err := zkClient.Conn(); err != nil {
			return nil, err
		} else {
			var children []string
			var stat *zk.Stat
			var events <-chan zk.Event
			var err error

			if b.watching.watched || b.watching.watcher != nil {
				children, stat, events, err = conn.ChildrenW(path)

				if events != nil && b.watching.watcher != nil {
					go NewWatchers(b.watching.watcher).Watch(events)
				}
			} else {
				children, stat, err = conn.Children(path)
			}

			if stat != nil {
				if b.stat != nil {
					*b.stat = *stat
				} else {
					b.stat = stat
				}
			}

			return children, err
		}
	})

	children, _ := result.([]string)

	return children, err
}

func (b *getChildrenBuilder) StoringStatIn(stat *zk.Stat) GetChildrenBuilder {
	b.stat = stat

	return b
}

func (b *getChildrenBuilder) Watched() GetChildrenBuilder {
	b.watching.watched = true

	return b
}

func (b *getChildrenBuilder) UsingWatcher(watcher Watcher) GetChildrenBuilder {
	b.watching.watcher = b.client.getNamespaceWatcher(watcher)

	return b
}

func (b *getChildrenBuilder) InBackground() GetChildrenBuilder {
	b.backgrounding = backgrounding{inBackground: true}

	return b
}

func (b *getChildrenBuilder) InBackgroundWithContext(context interface{}) GetChildrenBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context}

	return b
}

func (b *getChildrenBuilder) InBackgroundWithCallback(callback BackgroundCallback) GetChildrenBuilder {
	b.backgrounding = backgrounding{inBackground: true, callback: callback}

	return b
}

func (b *getChildrenBuilder) InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) GetChildrenBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context, callback: callback}

	return b
}
