package xsf

import (
	"sync"
	"sync/atomic"
	"time"
)

type errCodeBucket struct {
	sync.RWMutex
	errCodeMap map[int64]*int64
}

func (b *errCodeBucket) clear() {
	//todo gc问题?
	b.Lock()
	b.errCodeMap = make(map[int64]*int64)
	b.Unlock()
}
func (b *errCodeBucket) setErrCode(errCode int64) {
	b.RLock()
	_, exist := b.errCodeMap[errCode]

	if exist {
		atomic.AddInt64(b.errCodeMap[errCode], 1)
		b.RUnlock()
		return
	}
	b.RUnlock()

	v := int64(1)
	b.Lock()
	if b.errCodeMap == nil {
		b.errCodeMap = make(map[int64]*int64)
	}
	b.errCodeMap[errCode] = &v
	b.Unlock()

}

type errCodeWindow struct {
	//循环队列
	timeSlices []errCodeBucket

	//队列的总长度
	timeSliceSize int64

	//每个时间片的时长
	timePerSlice time.Duration

	//窗口长度
	winSize int64

	//当前所使用的时间片位置
	cursor int64

	//上一次访问的时间戳
	preTs time.Time
	//窗口总时长
	winDur time.Duration
}

func newErrCodeWindow(timePerSlice time.Duration, winSize int64) *errCodeWindow {
	slidingWindowInst := errCodeWindow{}
	slidingWindowInst.Init(timePerSlice, winSize)
	return &slidingWindowInst
}

//清理闲时的节点数据
func (s *errCodeWindow) reset() {

	ticker := time.NewTicker(s.winDur)

	for {
		select {
		case <-ticker.C:
			{
				if time.Now().Sub(s.preTs) < s.winDur {
					//访问间隔没有超过所有区间
					continue
				}
				for index := 0; index < len(s.timeSlices); index++ {
					s.timeSlices[index].clear()
				}
			}
		}
	}
}
func (s *errCodeWindow) Init(timePerSlice time.Duration, winSize int64) {
	s.timePerSlice = timePerSlice
	s.winSize = winSize
	// 保证存储在至少两个window
	s.timeSliceSize = winSize*2 + 1

	s.timeSlices = make([]errCodeBucket, s.timeSliceSize)
	s.preTs = time.Now()
	s.winDur = s.timePerSlice * time.Duration(s.winSize)
	go s.reset()
}

func (s *errCodeWindow) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

//对时间片做增加操作，并返回窗口中所有的计数总和
func (s *errCodeWindow) setErrCode(errCode int64) {
	var index = s.locationIndex()

	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, s.locationIndex())

	if oldCursor == index {
		// 在当前时间片里继续
		s.timeSlices[index].setErrCode(errCode)

	} else {
		s.timeSlices[index].setErrCode(errCode)

		// 清零，访问量不大时会有时间片跳跃的情况
		s.clearBetween(oldCursor, index)
	}

	s.preTs = time.Now()
}

//入参为tps的s
func (s *errCodeWindow) getStats(unit time.Duration) map[int64]int64 {
	var index = s.locationIndex()

	// cursor不等于index，将cursor设置为index
	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, index)

	if oldCursor != index {
		// 可能有其他goroutine已经置过，问题不大
		s.timeSlices[index].clear()

		// 清零，访问量不大时会有时间片跳跃的情况
		s.clearBetween(oldCursor, index)
	}

	rst := make(map[int64]int64)
	for i := int64(0); i < s.winSize; i++ {
		s.merge(rst, s.timeSlices[(index-i+s.timeSliceSize)%s.timeSliceSize].errCodeMap)
	}

	return rst
}
func (s *errCodeWindow) merge(m1 map[int64]int64, m2 map[int64]*int64) {
	for k, v := range m2 {
		m1[k] += *v
	}
}

//将fromIndex~toIndex之间的时间片计数都清零
//极端情况下，当循环队列已经走了超过1个timeSliceSize以上，这里的清零并不能如期望的进行
func (s *errCodeWindow) clearBetween(fromIndex, toIndex int64) {
	for index := (fromIndex + 1) % s.timeSliceSize; index != toIndex; index = (index + 1) % s.timeSliceSize {
		s.timeSlices[index].clear()
	}
}
