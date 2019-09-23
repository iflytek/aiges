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

type errCodeWin struct {
	timeSlices []errCodeBucket

	timeSliceSize int64

	timePerSlice time.Duration

	winSize int64

	cursor int64

	preTs time.Time

	winDur time.Duration
}

func newErrCodeWindow(timePerSlice time.Duration, winSize int64) *errCodeWin {
	slidingWindowInst := errCodeWin{}
	slidingWindowInst.Init(timePerSlice, winSize)
	return &slidingWindowInst
}

func (s *errCodeWin) reset() {

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
func (s *errCodeWin) Init(timePerSlice time.Duration, winSize int64) {
	s.timePerSlice = timePerSlice
	s.winSize = winSize
	s.timeSliceSize = winSize*2 + 1

	s.timeSlices = make([]errCodeBucket, s.timeSliceSize)
	s.preTs = time.Now()
	s.winDur = s.timePerSlice * time.Duration(s.winSize)
	go s.reset()
}

func (s *errCodeWin) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

func (s *errCodeWin) setErrCode(errCode int64) {
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

//入参为tps的s
func (s *errCodeWin) getStats(unit time.Duration) map[int64]int64 {
	var index = s.locationIndex()

	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, index)

	if oldCursor != index {
		s.timeSlices[index].clear()

		s.clearBetween(oldCursor, index)
	}

	rst := make(map[int64]int64)
	for i := int64(0); i < s.winSize; i++ {
		bucketTmp := s.timeSlices[(index-i+s.timeSliceSize)%s.timeSliceSize]
		s.merge(rst, bucketTmp.errCodeMap)
	}

	return rst
}
func (s *errCodeWin) merge(m1 map[int64]int64, m2 map[int64]*int64) {
	for k, v := range m2 {
		m1[k] += *v
	}
}

func (s *errCodeWin) clearBetween(fromIndex, toIndex int64) {
	for index := (fromIndex + 1) % s.timeSliceSize; index != toIndex; index = (index + 1) % s.timeSliceSize {
		s.timeSlices[index].clear()
	}
}
