package framework

import (
	"github.com/curator-go/curator"
	"github.com/samuel/go-zookeeper/zk"
)

func Create(client curator.CuratorFramework, path string, payload []byte) (string, error) {
	// this will create the given ZNode with the given data
	return client.Create().ForPathWithData(path, payload)
}

func CreateEphemeral(client curator.CuratorFramework, path string, payload []byte) (string, error) {
	// this will create the given EPHEMERAL ZNode with the given data
	return client.Create().WithMode(curator.EPHEMERAL).ForPathWithData(path, payload)
}

func CheckExists(client curator.CuratorFramework, path string) (*zk.Stat, error) {
	watcher := curator.NewWatcher(func(event *zk.Event) {
		// examine event for details
	})

	return client.CheckExists().UsingWatcher(watcher).ForPath(path)
}

func Delete(client curator.CuratorFramework, path string, version int32) error {
	return client.Delete().DeletingChildrenIfNeeded().WithVersion(version).ForPath(path)
}

func GetData(client curator.CuratorFramework, path string) ([]byte, error) {
	var stat zk.Stat

	watcher := curator.NewWatcher(func(event *zk.Event) {
		// examine event for details
	})

	return client.GetData().StoringStatIn(&stat).Decompressed().UsingWatcher(watcher).ForPath(path)
}

func SetData(client curator.CuratorFramework, path string, payload []byte) (*zk.Stat, error) {
	// set data for the given node
	return client.SetData().Compressed().ForPathWithData(path, payload)
}

func SetDataAsync(client curator.CuratorFramework, path string, payload []byte) (*zk.Stat, error) {
	listener := curator.NewCuratorListener(func(client curator.CuratorFramework, event curator.CuratorEvent) error {
		// examine event for details
		return nil
	})

	// this is one method of getting event/async notifications
	client.CuratorListenable().AddListener(listener)

	// set data for the given node asynchronously. The completion notification
	// is done via the CuratorListener.
	return client.SetData().InBackground().ForPathWithData(path, payload)
}

func SetDataAsyncWithCallback(client curator.CuratorFramework, callback curator.BackgroundCallback, path string, payload []byte) (*zk.Stat, error) {
	// this is another method of getting notification of an async completion
	return client.SetData().InBackgroundWithCallback(callback).ForPathWithData(path, payload)
}

func WatchedGetChildren(client curator.CuratorFramework, path string) ([]string, error) {
	// Get children and set a watcher on the node.
	// The watcher notification will come through the CuratorListener (see SetDataAsync() above).
	return client.GetChildren().Watched().ForPath(path)
}

func UsingWatcherGetChildren(client curator.CuratorFramework, path string, watcher curator.Watcher) ([]string, error) {
	// Get children and set the given watcher on the node.
	return client.GetChildren().UsingWatcher(watcher).ForPath(path)
}
