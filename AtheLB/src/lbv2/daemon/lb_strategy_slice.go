package daemon

import (
	"fmt"
	"sync/atomic"
	"time"
)

type SubSvcItem struct {
	addr      string
	timestamp int64 //纳秒,节点最近一次上报的时间

	bestInst  int64
	idleInst  int64
	totalInst int64
}

func (s *SubSvcItem) String() string {
	return fmt.Sprintf("addr:%s, timestamp:%d, bestInst:%d, idleInst:%d, totalInst:%d",
		s.addr, s.timestamp, s.bestInst, s.idleInst, s.totalInst)
}

func (s *SubSvcItem) set(total, best, idle, dead int64) {
	atomic.StoreInt64(&s.timestamp, time.Now().UnixNano())
	atomic.StoreInt64(&s.totalInst, total)
	atomic.StoreInt64(&s.bestInst, best)
	atomic.StoreInt64(&s.idleInst, idle)
	//atomic.StoreInt64(&s.dead, dead)
}

//预授
func (s *SubSvcItem) preAuthorization(delta int64) {
	atomic.AddInt64(&s.idleInst, delta)
}

//func (s *SubSvcItem) isDead() int64 {
//	return atomic.LoadInt64(&s.dead)
//}

//func (s *SubSvcItem) disposed() {
//	atomic.StoreInt64(&s.dead, 1)
//}

func (s *SubSvcItem) FInit() *SubSvcItem {
	s.addr = ""
	s.bestInst = 0
	s.idleInst = 0
	s.totalInst = 0
	return s
}

type SubSvcItemSlice []*SubSvcItem

func (l SubSvcItemSlice) Len() int {
	return len(l)
}
func (l SubSvcItemSlice) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l SubSvcItemSlice) Less(i, j int) bool {
	return l[i].idleInst < l[j].idleInst
}
