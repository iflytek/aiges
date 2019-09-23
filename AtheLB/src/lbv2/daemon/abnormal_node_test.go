package daemon

import (
	"testing"
	"time"
)

func Test_abnormalWindow(t *testing.T) {
	win := newAbnormalNodeWindow(time.Second, 3)
	win.setAbnormalNode("127.0.0.1:1")
	win.setAbnormalNode("127.0.0.1:2")
	win.setAbnormalNode("127.0.0.1:3")
	t.Logf("ts:%v,1s:%v\n", time.Now(), win.getStats())
	t.Logf("ts:%v,2s:%v\n", time.Now(), win.getStats())
	time.Sleep(time.Second * 3)
	t.Logf("ts:%v,1s:%v\n", time.Now(), win.getStats())
	t.Logf("ts:%v,2s:%v\n", time.Now(), win.getStats())
}
