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
type delayWin struct {
	timeSlices []delayBucket

	timeSliceSize int64

	timePerSlice time.Duration

	winSize int64

	cursor int64

	preTs time.Time

	winDur time.Duration
}

func newDelayWindow(timePerSlice time.Duration, winSize int64) *delayWin {
	slidingWindowInst := delayWin{}
	slidingWindowInst.Init(timePerSlice, winSize)
	return &slidingWindowInst
}

func (s *delayWin) reset() {

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
func (s *delayWin) Init(timePerSlice time.Duration, winSize int64) {
	s.timePerSlice = timePerSlice
	s.winSize = winSize
	s.timeSliceSize = winSize*2 + 1

	s.timeSlices = make([]delayBucket, s.timeSliceSize)
	s.preTs = time.Now()
	s.winDur = s.timePerSlice * time.Duration(s.winSize)
	go s.reset()
}

func (s *delayWin) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

func (s *delayWin) setDur(dur int64) {
	var index = s.locationIndex()
	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, s.locationIndex())

	if oldCursor == index {
		atomic.AddInt64(&s.timeSlices[index].sum, dur)
		atomic.AddInt64(&s.timeSlices[index].cnt, 1)
		if dur > atomic.LoadInt64(&s.timeSlices[index].max) {
			atomic.StoreInt64(&s.timeSlices[index].max, dur)
		}
		if dur < atomic.LoadInt64(&s.timeSlices[index].min) {
			atomic.StoreInt64(&s.timeSlices[index].min, dur)
		}

	} else {
		atomic.AddInt64(&s.timeSlices[index].sum, dur)
		atomic.AddInt64(&s.timeSlices[index].cnt, 1)
		atomic.StoreInt64(&s.timeSlices[index].max, dur)
		atomic.StoreInt64(&s.timeSlices[index].min, dur)

		s.clearBetween(oldCursor, index)
	}

	s.preTs = time.Now()
}

//入参为tps的s
func (s *delayWin) getStats() (max, min, avg int64) {
	var index = s.locationIndex()

	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, index)

	if oldCursor != index {
		atomic.StoreInt64(&s.timeSlices[index].sum, 0)
		atomic.StoreInt64(&s.timeSlices[index].max, 0)
		atomic.StoreInt64(&s.timeSlices[index].min, 0)

		s.clearBetween(oldCursor, index)
	}

	sum := int64(0)
	cnt := int64(0)
	max = math.MinInt64
	min = math.MaxInt64
	for i := int64(0); i < s.winSize; i++ {
		bucketTmp := s.timeSlices[(index-i+s.timeSliceSize)%s.timeSliceSize]
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

	return
}

func (s *delayWin) clearBetween(fromIndex, toIndex int64) {
	for index := (fromIndex + 1) % s.timeSliceSize; index != toIndex; index = (index + 1) % s.timeSliceSize {
		atomic.StoreInt64(&s.timeSlices[index].sum, 0)
		atomic.StoreInt64(&s.timeSlices[index].cnt, 0)
		atomic.StoreInt64(&s.timeSlices[index].max, 0)
		atomic.StoreInt64(&s.timeSlices[index].min, 0)
	}
}
