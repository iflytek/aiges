package util

import (
	"fmt"
	"testing"
	"time"
)

func TestScheduledTaskPool(t *testing.T) {
	st := NewScheduledTaskPool()   // 新建定时任务池
	st.Start(time.Second, func() { // 启动一个1s的定时任务
		fmt.Println("Start")
	})

	time.Sleep(10 * time.Second)
	st.Stop() // 关闭定时任务
}
