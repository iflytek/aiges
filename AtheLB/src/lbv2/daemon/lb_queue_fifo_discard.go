package daemon

import (
	"github.com/Workiva/go-datastructures/queue"
	"time"
)

const TM = time.Millisecond
const defaultQueueSize = 1000

var LoopQueueInst LoopQueue

func init() {
	LoopQueueInst.Init(defaultQueueSize)
}

type LoopQueue struct {
	qu *queue.RingBuffer
}

func (l *LoopQueue) Init(size uint64) LbErr {
	if size <= 0 {
		return ErrLbLoopQueueSize
	}
	l.qu = queue.NewRingBuffer(size)
	return nil
}

/*
队列满时：读写指针同时向前移动，如同滑动窗口
*/
func (l *LoopQueue) Put(item interface{}) {
	var ok bool
	ok, _ = l.qu.Offer(item)
	if !ok {
		/*
			此处仅尝试三次，因为在大并发写入的情况下，丢失的只是老数据，符合我们的预期
		*/
		for retry := 0; retry < 3; retry++ {
			l.qu.Poll(TM)
			ok, _ = l.qu.Offer(item)
			if ok {
				break
			}
		}
	}
}
func (l *LoopQueue) Get() (interface{}, error) {
	return l.qu.Poll(TM)
}
