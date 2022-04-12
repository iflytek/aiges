package xsf

import (
	"math"
	"sync/atomic"
	"time"
)

type delayBucket struct {
	max int64
	min int64
	sum int64
	cnt int64
}
type delayWindow struct {
	//循环队列
	timeSlices []delayBucket

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

func newDelayWindow(timePerSlice time.Duration, winSize int64) *delayWindow {
	slidingWindowInst := delayWindow{}
	slidingWindowInst.Init(timePerSlice, winSize)
	return &slidingWindowInst
}

//清理闲时的节点数据
func (s *delayWindow) reset() {

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
					atomic.StoreInt64(&s.timeSlices[index].cnt, 0)
				}
			}
		}
	}
}
func (s *delayWindow) Init(timePerSlice time.Duration, winSize int64) {
	s.timePerSlice = timePerSlice
	s.winSize = winSize
	// 保证存储在至少两个window
	s.timeSliceSize = winSize*2 + 1

	s.timeSlices = make([]delayBucket, s.timeSliceSize)
	s.preTs = time.Now()
	s.winDur = s.timePerSlice * time.Duration(s.winSize)
	go s.reset()
}

func (s *delayWindow) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

//对时间片做增加操作，并返回窗口中所有的计数总和
func (s *delayWindow) setDur(dur int64) {
	var index = s.locationIndex()
	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, s.locationIndex())

	if oldCursor == index {
		// 在当前时间片里继续
		atomic.AddInt64(&s.timeSlices[index].sum, dur)
		atomic.AddInt64(&s.timeSlices[index].cnt, 1)
		if dur > atomic.LoadInt64(&s.timeSlices[index].max) {
			atomic.StoreInt64(&s.timeSlices[index].max, dur)
		}
		if dur < atomic.LoadInt64(&s.timeSlices[index].min) {
			atomic.StoreInt64(&s.timeSlices[index].min, dur)
		}

	} else {
		atomic.StoreInt64(&s.timeSlices[index].sum, dur)
		atomic.StoreInt64(&s.timeSlices[index].cnt, 1)
		atomic.StoreInt64(&s.timeSlices[index].max, dur)
		atomic.StoreInt64(&s.timeSlices[index].min, dur)

		// 清零，访问量不大时会有时间片跳跃的情况
		s.clearBetween(oldCursor, index)
	}

	s.preTs = time.Now()
}

//入参为tps的s
func (s *delayWindow) getStats() (max, min, avg, qps int64) {
	var index = s.locationIndex()

	// cursor不等于index，将cursor设置为index
	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, index)

	if oldCursor != index {
		// 可能有其他goroutine已经置过1，问题不大
		atomic.StoreInt64(&s.timeSlices[index].sum, 0)
		atomic.StoreInt64(&s.timeSlices[index].max, 0)
		atomic.StoreInt64(&s.timeSlices[index].min, 0)
		atomic.StoreInt64(&s.timeSlices[index].cnt, 0)

		// 清零，访问量不大时会有时间片跳跃的情况
		s.clearBetween(oldCursor, index)
	}

	sum := int64(0)
	cnt := int64(0)
	qpsCnt := int64(0)
	max = math.MinInt64
	min = math.MaxInt64
	for i := int64(0); i < s.winSize; i++ {
		bucketTmp := s.timeSlices[(index-i+s.timeSliceSize)%s.timeSliceSize]
		qpsCnt += bucketTmp.cnt
		if bucketTmp.min == 0 || bucketTmp.max == 0 {
			continue
		}
		if max < atomic.LoadInt64(&bucketTmp.max) {
			max = atomic.LoadInt64(&bucketTmp.max)
		}
		if min > atomic.LoadInt64(&bucketTmp.min) {
			min = atomic.LoadInt64(&bucketTmp.min)
		}

		sum += atomic.LoadInt64(&bucketTmp.sum)
		cnt += atomic.LoadInt64(&bucketTmp.cnt)
	}
	if 0 == cnt {
		avg = 0
	} else {
		avg = sum / cnt
	}

	if winDurSecs := s.winDur.Milliseconds() / 1e3; winDurSecs == 0 {
		qps = 0
	} else {
		qps = qpsCnt / winDurSecs
	}
	return
}

//将fromIndex~toIndex之间的时间片计数都清零
//极端情况下，当循环队列已经走了超过1个timeSliceSize以上，这里的清零并不能如期望的进行
func (s *delayWindow) clearBetween(fromIndex, toIndex int64) {
	for index := (fromIndex + 1) % s.timeSliceSize; index != toIndex; index = (index + 1) % s.timeSliceSize {
		atomic.StoreInt64(&s.timeSlices[index].sum, 0)
		atomic.StoreInt64(&s.timeSlices[index].cnt, 0)
		atomic.StoreInt64(&s.timeSlices[index].max, 0)
		atomic.StoreInt64(&s.timeSlices[index].min, 0)
	}
}
