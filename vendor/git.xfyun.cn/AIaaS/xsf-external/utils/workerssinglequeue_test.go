package utils

import (
	"fmt"
	"time"
	"sync/atomic"
	"testing"
)

var count int64

type taskDemo struct{}

func (m *taskDemo) Task() {
	atomic.AddInt64(&count, 1)
	fmt.Printf("NO.%v\t->\t任务开始\n", count)
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("NO.%v\t->\t任务结束\n", count)
}

func TestWorkers(t *testing.T) {
	np := taskDemo{}
	p := New(20)
	for cnt := 0; cnt < 1000; cnt++ {
		p.Run(&np)
	}
	p.Shutdown()
}
