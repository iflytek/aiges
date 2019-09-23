package xsf

import (
	"sync/atomic"
	"time"
)


//todo 后续可把窗口部分合并下
type healthBucket struct {
	count    int64
	errCount int64
}

func (b *healthBucket) clear() {
	atomic.StoreInt64(&b.count, 0)
	atomic.StoreInt64(&b.errCount, 0)
}
func (b *healthBucket) setErrCode(errCode int64) {
	if 0 != errCode {
		atomic.AddInt64(&b.errCount, 1)
	}
	atomic.AddInt64(&b.count, 1)
}

type healthWin struct {
	timeSlices []healthBucket

	timeSliceSize int64

	timePerSlice time.Duration

	winSize int64

	cursor int64

	preTs time.Time

	winDur time.Duration
}

func newHealthWin(timePerSlice time.Duration, winSize int64) *healthWin {
	healthWindowInst := healthWin{}
	healthWindowInst.Init(timePerSlice, winSize)
	return &healthWindowInst
}

func (s *healthWin) reset() {

	ticker := time.NewTicker(s.winDur)

	for {
		select {
		case <-ticker.C:
			{
				if time.Now().Sub(s.preTs) < s.winDur {
					continue
				}
				for index := 0; index < len(s.timeSlices); index++ {
					s.timeSlices[index].clear()
				}
			}
		}
	}
}
func (s *healthWin) Init(timePerSlice time.Duration, winSize int64) {
	s.timePerSlice = timePerSlice
	s.winSize = winSize
	s.timeSliceSize = winSize*2 + 1

	s.timeSlices = make([]healthBucket, s.timeSliceSize)
	s.preTs = time.Now()
	s.winDur = s.timePerSlice * time.Duration(s.winSize)
	go s.reset()
}

func (s *healthWin) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

func (s *healthWin) setErrCode(errCode int64) {
	var index = s.locationIndex()

	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, s.locationIndex())

	if oldCursor == index {
		s.timeSlices[index].setErrCode(errCode)

	} else {
		s.timeSlices[index].setErrCode(errCode)

		s.clearBetween(oldCursor, index)
	}

	s.preTs = time.Now()
}

func (s *healthWin) getStats() (errCount, count int64) {
	var index = s.locationIndex()

	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, index)

	if oldCursor != index {
		s.timeSlices[index].clear()

		s.clearBetween(oldCursor, index)
	}

	for i := int64(0); i < s.winSize; i++ {
		bucketTmp := s.timeSlices[(index-i+s.timeSliceSize)%s.timeSliceSize]
		errCount += bucketTmp.errCount
		count += bucketTmp.count
	}
	return
}

func (s *healthWin) clearBetween(fromIndex, toIndex int64) {
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
	timeSlices []latencyBucket

	timeSliceSize int64

	timePerSlice time.Duration

	winSize int64

	cursor int64

	preTs time.Time

	winDur time.Duration
}

func newLatencyWin(timePerSlice time.Duration, winSize int64) *latencyWin {
	latencyWindowInst := latencyWin{}
	latencyWindowInst.Init(timePerSlice, winSize)
	return &latencyWindowInst
}

func (s *latencyWin) reset() {

	ticker := time.NewTicker(s.winDur)

	for {
		select {
		case <-ticker.C:
			{
				if time.Now().Sub(s.preTs) < s.winDur {
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
	s.timeSliceSize = winSize*2 + 1

	s.timeSlices = make([]latencyBucket, s.timeSliceSize)
	s.preTs = time.Now()
	s.winDur = s.timePerSlice * time.Duration(s.winSize)
	go s.reset()
}

func (s *latencyWin) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

func (s *latencyWin) addDur(dur int64) {
	var index = s.locationIndex()
	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, s.locationIndex())

	if oldCursor == index {
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

		s.clearBetween(oldCursor, index)
	}

	s.preTs = time.Now()
}

func (s *latencyWin) getStats() (max, min, sum, count int64) {
	var index = s.locationIndex()

	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, index)

	if oldCursor != index {
		atomic.StoreInt64(&s.timeSlices[index].sum, 0)
		atomic.StoreInt64(&s.timeSlices[index].count, 0)
		atomic.StoreInt64(&s.timeSlices[index].max, 0)
		atomic.StoreInt64(&s.timeSlices[index].min, 0)

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
		if 0 == bucketTmp.min || 0 == bucketTmp.max {
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

func (s *latencyWin) clearBetween(fromIndex, toIndex int64) {
	for index := (fromIndex + 1) % s.timeSliceSize; index != toIndex; index = (index + 1) % s.timeSliceSize {
		atomic.StoreInt64(&s.timeSlices[index].sum, 0)
		atomic.StoreInt64(&s.timeSlices[index].count, 0)
		atomic.StoreInt64(&s.timeSlices[index].max, 0)
		atomic.StoreInt64(&s.timeSlices[index].min, 0)
	}
}
