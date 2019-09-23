package recipes

import (
	"time"

	"github.com/curator-go/curator"
	"github.com/fanliao/go-promise"
)

// Utility class to allow execution of logic once a ZooKeeper connection becomes available.
type AfterConnectionEstablished struct {
	Client  curator.CuratorFramework
	Timeout time.Duration
}

// Spawns a new new background thread that will block
// until a connection is available and then execute the 'runAfterConnection' logic
func (c *AfterConnectionEstablished) Future() *promise.Future {
	p := promise.NewPromise()

	go func() {
		if err := c.Client.BlockUntilConnectedTimeout(c.Timeout); err != nil {
			p.Reject(err)
		} else {
			p.Resolve(p)
		}
	}()

	return p.Future
}
