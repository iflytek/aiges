package curator

import (
	"github.com/samuel/go-zookeeper/zk"
)

type deleteBuilder struct {
	client                   *curatorFramework
	backgrounding            backgrounding
	deletingChildrenIfNeeded bool
	version                  int32
}

func (b *deleteBuilder) ForPath(givenPath string) error {
	adjustedPath := b.client.fixForNamespace(givenPath, false)

	if b.backgrounding.inBackground {
		go b.pathInBackground(adjustedPath, givenPath)

		return nil
	} else {
		return b.pathInForeground(adjustedPath, givenPath)
	}
}

func (b *deleteBuilder) pathInBackground(path string, givenPath string) {
	tracer := b.client.ZookeeperClient().StartTracer("deleteBuilder.pathInBackground")

	defer tracer.Commit()

	err := b.pathInForeground(path, givenPath)

	if b.backgrounding.callback != nil {
		event := &curatorEvent{
			eventType: DELETE,
			err:       err,
			path:      b.client.unfixForNamespace(path),
			context:   b.backgrounding.context,
		}

		if err != nil {
			event.path = givenPath
		}

		event.name = GetNodeFromPath(event.path)

		b.backgrounding.callback(b.client, event)
	}
}

func (b *deleteBuilder) pathInForeground(path string, givenPath string) error {
	zkClient := b.client.ZookeeperClient()

	_, err := zkClient.NewRetryLoop().CallWithRetry(func() (interface{}, error) {
		conn, err := zkClient.Conn()

		if err == nil {
			err = conn.Delete(path, b.version)

			if err == zk.ErrNotEmpty && b.deletingChildrenIfNeeded {
				err = DeleteChildren(conn, path, true)
			}
		}

		return nil, err
	})

	return err
}

func (b *deleteBuilder) DeletingChildrenIfNeeded() DeleteBuilder {
	b.deletingChildrenIfNeeded = true

	return b
}

func (b *deleteBuilder) WithVersion(version int32) DeleteBuilder {
	b.version = version

	return b
}

func (b *deleteBuilder) InBackground() DeleteBuilder {
	b.backgrounding = backgrounding{inBackground: true}

	return b
}

func (b *deleteBuilder) InBackgroundWithContext(context interface{}) DeleteBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context}

	return b
}

func (b *deleteBuilder) InBackgroundWithCallback(callback BackgroundCallback) DeleteBuilder {
	b.backgrounding = backgrounding{inBackground: true, callback: callback}

	return b
}

func (b *deleteBuilder) InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) DeleteBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context, callback: callback}

	return b
}
