package curator

type syncBuilder struct {
	client        *curatorFramework
	backgrounding backgrounding
}

func (b *syncBuilder) ForPath(givenPath string) (string, error) {
	adjustedPath := b.client.fixForNamespace(givenPath, false)

	if b.backgrounding.inBackground {
		go b.pathInBackground(adjustedPath, givenPath)

		return givenPath, nil
	} else {
		return b.pathInForeground(adjustedPath)
	}
}

func (b *syncBuilder) pathInBackground(path string, givenPath string) {
	tracer := b.client.ZookeeperClient().StartTracer("syncBuilder.pathInBackground")

	defer tracer.Commit()

	syncPath, err := b.pathInForeground(path)

	if b.backgrounding.callback != nil {
		event := &curatorEvent{
			eventType: SYNC,
			err:       err,
			path:      b.client.unfixForNamespace(syncPath),
			context:   b.backgrounding.context,
		}

		if err != nil {
			event.path = givenPath
		}

		event.name = GetNodeFromPath(event.path)

		b.backgrounding.callback(b.client, event)
	}
}

func (b *syncBuilder) pathInForeground(path string) (string, error) {
	zkClient := b.client.ZookeeperClient()

	result, err := zkClient.NewRetryLoop().CallWithRetry(func() (interface{}, error) {
		if conn, err := zkClient.Conn(); err != nil {
			return nil, err
		} else {
			return conn.Sync(path)
		}
	})

	syncPath, _ := result.(string)

	return b.client.unfixForNamespace(syncPath), err
}

func (b *syncBuilder) InBackground() SyncBuilder {
	b.backgrounding = backgrounding{inBackground: true}

	return b
}

func (b *syncBuilder) InBackgroundWithContext(context interface{}) SyncBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context}

	return b
}

func (b *syncBuilder) InBackgroundWithCallback(callback BackgroundCallback) SyncBuilder {
	b.backgrounding = backgrounding{inBackground: true, callback: callback}

	return b
}

func (b *syncBuilder) InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) SyncBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context, callback: callback}

	return b
}
