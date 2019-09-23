package curator

import (
	"github.com/samuel/go-zookeeper/zk"
)

type getDataBuilder struct {
	client        *curatorFramework
	backgrounding backgrounding
	decompress    bool
	stat          *zk.Stat
	watching      watching
}

func (b *getDataBuilder) ForPath(givenPath string) ([]byte, error) {
	adjustedPath := b.client.fixForNamespace(givenPath, false)

	if b.backgrounding.inBackground {
		go b.pathInBackground(adjustedPath, givenPath)

		return nil, nil
	}

	if payload, err := b.pathInForeground(adjustedPath); err != nil {
		return nil, err
	} else {
		return payload, err
	}
}

func (b *getDataBuilder) pathInBackground(adjustedPath, givenPath string) {
	tracer := b.client.ZookeeperClient().StartTracer("getDataBuilder.pathInBackground")

	defer tracer.Commit()

	data, err := b.pathInForeground(adjustedPath)

	if b.backgrounding.callback != nil {
		event := &curatorEvent{
			eventType: GET_DATA,
			err:       err,
			path:      b.client.unfixForNamespace(adjustedPath),
			data:      data,
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

func (b *getDataBuilder) pathInForeground(path string) ([]byte, error) {
	zkClient := b.client.ZookeeperClient()

	result, err := zkClient.NewRetryLoop().CallWithRetry(func() (interface{}, error) {
		if conn, err := zkClient.Conn(); err != nil {
			return nil, err
		} else {
			var data []byte
			var stat *zk.Stat
			var events <-chan zk.Event
			var err error

			if b.watching.watched || b.watching.watcher != nil {
				data, stat, events, err = conn.GetW(path)

				if events != nil && b.watching.watcher != nil {
					go NewWatchers(b.watching.watcher).Watch(events)
				}
			} else {
				data, stat, err = conn.Get(path)
			}

			if stat != nil {
				if b.stat != nil {
					*b.stat = *stat
				} else {
					b.stat = stat
				}
			}

			if b.decompress {
				if payload, err := b.client.compressionProvider.Decompress(path, data); err != nil {
					return nil, err
				} else {
					data = payload
				}
			}

			return data, err
		}
	})

	data, _ := result.([]byte)

	return data, err
}

func (b *getDataBuilder) Decompressed() GetDataBuilder {
	b.decompress = true

	return b
}

func (b *getDataBuilder) StoringStatIn(stat *zk.Stat) GetDataBuilder {
	b.stat = stat

	return b
}

func (b *getDataBuilder) Watched() GetDataBuilder {
	b.watching.watched = true

	return b
}

func (b *getDataBuilder) UsingWatcher(watcher Watcher) GetDataBuilder {
	b.watching.watcher = b.client.getNamespaceWatcher(watcher)

	return b
}

func (b *getDataBuilder) InBackground() GetDataBuilder {
	b.backgrounding = backgrounding{inBackground: true}

	return b
}

func (b *getDataBuilder) InBackgroundWithContext(context interface{}) GetDataBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context}

	return b
}

func (b *getDataBuilder) InBackgroundWithCallback(callback BackgroundCallback) GetDataBuilder {
	b.backgrounding = backgrounding{inBackground: true, callback: callback}

	return b
}

func (b *getDataBuilder) InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) GetDataBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context, callback: callback}

	return b
}

type setDataBuilder struct {
	client        *curatorFramework
	backgrounding backgrounding
	version       int32
	compress      bool
}

func (b *setDataBuilder) ForPath(path string) (*zk.Stat, error) {
	return b.ForPathWithData(path, b.client.defaultData)
}

func (b *setDataBuilder) ForPathWithData(givenPath string, payload []byte) (*zk.Stat, error) {
	if b.compress {
		if data, err := b.client.compressionProvider.Compress(givenPath, payload); err != nil {
			return nil, err
		} else {
			payload = data
		}
	}

	adjustedPath := b.client.fixForNamespace(givenPath, false)

	if b.backgrounding.inBackground {
		go b.pathInBackground(adjustedPath, payload, givenPath)

		return nil, nil
	} else {
		return b.pathInForeground(adjustedPath, payload)
	}
}

func (b *setDataBuilder) pathInBackground(path string, payload []byte, givenPath string) {
	tracer := b.client.ZookeeperClient().StartTracer("setDataBuilder.pathInBackground")

	defer tracer.Commit()

	stat, err := b.pathInForeground(path, payload)

	if b.backgrounding.callback != nil {
		event := &curatorEvent{
			eventType: SET_DATA,
			err:       err,
			path:      b.client.unfixForNamespace(path),
			data:      payload,
			stat:      stat,
			context:   b.backgrounding.context,
		}

		if err != nil {
			event.path = givenPath
		}

		event.name = GetNodeFromPath(event.path)

		b.backgrounding.callback(b.client, event)
	}
}

func (b *setDataBuilder) pathInForeground(path string, payload []byte) (*zk.Stat, error) {
	zkClient := b.client.ZookeeperClient()

	result, err := zkClient.NewRetryLoop().CallWithRetry(func() (interface{}, error) {
		if conn, err := zkClient.Conn(); err != nil {
			return nil, err
		} else {
			return conn.Set(path, payload, b.version)
		}
	})

	stat, _ := result.(*zk.Stat)

	return stat, err
}

func (b *setDataBuilder) WithVersion(version int32) SetDataBuilder {
	b.version = version

	return b
}

func (b *setDataBuilder) Compressed() SetDataBuilder {
	b.compress = true

	return b
}

func (b *setDataBuilder) InBackground() SetDataBuilder {
	b.backgrounding = backgrounding{inBackground: true}

	return b
}

func (b *setDataBuilder) InBackgroundWithContext(context interface{}) SetDataBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context}

	return b
}

func (b *setDataBuilder) InBackgroundWithCallback(callback BackgroundCallback) SetDataBuilder {
	b.backgrounding = backgrounding{inBackground: true, callback: callback}

	return b
}

func (b *setDataBuilder) InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) SetDataBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context, callback: callback}

	return b
}
