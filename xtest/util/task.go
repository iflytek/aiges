package util

import (
	"sync"
	"time"
)

type ScheduledTaskPool struct {
	Size         int            // 任务个数
	TimeStopChan chan bool      // // 通知定时任务结束协程
	wg           sync.WaitGroup // 同步原语
}

func NewScheduledTaskPool() ScheduledTaskPool {
	return ScheduledTaskPool{
		Size:         0,
		TimeStopChan: make(chan bool, 100),
		wg:           sync.WaitGroup{},
	}
}

// Start 启动一个定时任务 jbzhou5
func (stp *ScheduledTaskPool) Start(d time.Duration, f func()) {
	stp.Size++
	stp.wg.Add(1)
	go func() {
		ticker := time.NewTicker(d)
		for {
			select {
			case <-ticker.C:
				f()
			case <-stp.TimeStopChan:
				stp.wg.Done()
				return
			}
		}
	}()
}

// Stop 结束定时任务
func (stp *ScheduledTaskPool) Stop() {
	for i := 0; i < stp.Size<<1; i++ {
		stp.TimeStopChan <- true
	}
	stp.wg.Wait()
}
