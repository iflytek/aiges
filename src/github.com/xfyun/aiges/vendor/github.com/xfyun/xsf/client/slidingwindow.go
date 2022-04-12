package xsf

import (
	"sync/atomic"
	"time"
)

type healthBucket struct {
	count    int64
	errCount int64
}

func (b *healthBucket) clear() {
	atomic.StoreInt64(&b.count, 0)
	atomic.StoreInt64(&b.errCount, 0)
}
func (b *healthBucket) setErrCode(errCode int64) {
	if errCode != 0 {
		atomic.AddInt64(&b.errCount, 1)
	}
	atomic.AddInt64(&b.count, 1)
}

type healthWindow struct {
	//循环队列
	timeSlices []healthBucket

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

func newHealthWindow(timePerSlice time.Duration, winSize int64) *healthWindow {
	healthWindowInst := healthWindow{}
	healthWindowInst.Init(timePerSlice, winSize)
	return &healthWindowInst
}

//清理闲时的节点数据
func (s *healthWindow) reset() {

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
func (s *healthWindow) Init(timePerSlice time.Duration, winSize int64) {
	s.timePerSlice = timePerSlice
	s.winSize = winSize
	// 保证存储在至少两个window
	s.timeSliceSize = winSize*2 + 1

	s.timeSlices = make([]healthBucket, s.timeSliceSize)
	s.preTs = time.Now()
	s.winDur = s.timePerSlice * time.Duration(s.winSize)
	go s.reset()
}

func (s *healthWindow) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

func (s *healthWindow) setErrCode(errCode int64) {
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
func (s *healthWindow) getStats() (errCount, count int64) {
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

	for i := int64(0); i < s.winSize; i++ {
		bucketTmp := s.timeSlices[(index-i+s.timeSliceSize)%s.timeSliceSize]
		errCount += bucketTmp.errCount
		count += bucketTmp.count
	}
	return
}

//将fromIndex~toIndex之间的时间片计数都清零
//极端情况下，当循环队列已经走了超过1个timeSliceSize以上，这里的清零并不能如期望的进行
func (s *healthWindow) clearBetween(fromIndex, toIndex int64) {
	for index := (fromIndex + 1) % s.timeSliceSize; index != toIndex; index = (index + 1) % s.timeSliceSize {
		s.timeSlices[index].clear()
	}
}

type latencyBucket struct {
	max   int64
	min   int64
	sum   int64
	count int64
}
type latencyWin struct {
	//循环队列
	timeSlices []latencyBucket

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

func newLatencyWindow(timePerSlice time.Duration, winSize int64) *latencyWin {
	latencyWindowInst := latencyWin{}
	latencyWindowInst.Init(timePerSlice, winSize)
	return &latencyWindowInst
}

//清理闲时的节点数据
func (s *latencyWin) reset() {

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
					atomic.StoreInt64(&s.timeSlices[index].sum, 0)
					atomic.StoreInt64(&s.timeSlices[index].max, 0)
					atomic.StoreInt64(&s.timeSlices[index].min, 0)
				}
			}
		}
	}
}
func (s *latencyWin) Init(timePerSlice time.Duration, winSize int64) {
	s.timePerSlice = timePerSlice
	s.winSize = winSize
	// 保证存储在至少两个window
	s.timeSliceSize = winSize*2 + 1

	s.timeSlices = make([]latencyBucket, s.timeSliceSize)
	s.preTs = time.Now()
	s.winDur = s.timePerSlice * time.Duration(s.winSize)
	go s.reset()
}

func (s *latencyWin) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

//对时间片做增加操作，并返回窗口中所有的计数总和
func (s *latencyWin) addDur(dur int64) {
	var index = s.locationIndex()
	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, s.locationIndex())

	if oldCursor == index {
		// 在当前时间片里继续
		atomic.AddInt64(&s.timeSlices[index].sum, dur)
		atomic.AddInt64(&s.timeSlices[index].count, 1)
		if dur > atomic.LoadInt64(&s.timeSlices[index].max) {
			atomic.StoreInt64(&s.timeSlices[index].max, dur)
		}
		if dur < atomic.LoadInt64(&s.timeSlices[index].min) {
			atomic.StoreInt64(&s.timeSlices[index].min, dur)
		}
	} else {
		atomic.AddInt64(&s.timeSlices[index].sum, dur)
		atomic.AddInt64(&s.timeSlices[index].count, 1)
		atomic.StoreInt64(&s.timeSlices[index].max, dur)
		atomic.StoreInt64(&s.timeSlices[index].min, dur)

		// 清零，访问量不大时会有时间片跳跃的情况
		s.clearBetween(oldCursor, index)
	}

	s.preTs = time.Now()
}

//入参为tps的s
//todo new object,much newObject
func (s *latencyWin) getStats() (max, min, sum, count int64) {
	var index = s.locationIndex()

	// cursor不等于index，将cursor设置为index
	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, index)

	if oldCursor != index {
		// 可能有其他goroutine已经置过1，问题不大
		atomic.StoreInt64(&s.timeSlices[index].sum, 0)
		atomic.StoreInt64(&s.timeSlices[index].count, 0)
		atomic.StoreInt64(&s.timeSlices[index].max, 0)
		atomic.StoreInt64(&s.timeSlices[index].min, 0)

		// 清零，访问量不大时会有时间片跳跃的情况
		s.clearBetween(oldCursor, index)
	}

	sum = 0
	count = 0
	max = 0
	min = 0
	for i := int64(0); i < s.winSize; i++ {
		bucketTmp := s.timeSlices[(index-i+s.timeSliceSize)%s.timeSliceSize]
		sum += atomic.LoadInt64(&bucketTmp.sum)
		count += atomic.LoadInt64(&bucketTmp.count)
		if bucketTmp.min == 0 || bucketTmp.max == 0 {
			continue
		}
		if max < atomic.LoadInt64(&bucketTmp.max) {
			max = atomic.LoadInt64(&bucketTmp.max)
		}
		if min > atomic.LoadInt64(&bucketTmp.min) {
			min = atomic.LoadInt64(&bucketTmp.min)
		}

	}

	return
}

//将fromIndex~toIndex之间的时间片计数都清零
//极端情况下，当循环队列已经走了超过1个timeSliceSize以上，这里的清零并不能如期望的进行
func (s *latencyWin) clearBetween(fromIndex, toIndex int64) {
	for index := (fromIndex + 1) % s.timeSliceSize; index != toIndex; index = (index + 1) % s.timeSliceSize {
		atomic.StoreInt64(&s.timeSlices[index].sum, 0)
		atomic.StoreInt64(&s.timeSlices[index].count, 0)
		atomic.StoreInt64(&s.timeSlices[index].max, 0)
		atomic.StoreInt64(&s.timeSlices[index].min, 0)
	}
}
