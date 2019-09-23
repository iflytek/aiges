package daemon

import (
	"sync"
	"sync/atomic"
	"time"
)

type nodeBucket struct {
	sync.RWMutex
	nodeMap map[string]struct{}
}

func (b *nodeBucket) clear() {
	b.Lock()
	b.nodeMap = make(map[string]struct{})
	b.Unlock()
}
func (b *nodeBucket) setNode(node string) {
	b.RLock()
	_, exist := b.nodeMap[node]

	if exist {
		b.RUnlock()
		return
	}
	b.RUnlock()

	b.Lock()
	if nil == b.nodeMap {
		b.nodeMap = make(map[string]struct{})
	}
	b.nodeMap[node] = struct{}{}
	b.Unlock()

}

type abnormalNodeWindow struct {
	//循环队列
	timeSlices []nodeBucket

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

const (
	nodeWindowsBoundary = "_"
	timePerSlice        = time.Second
	defWinSize          = 1
)

func newAbnormalNodeWindow(nodeDur time.Duration) *abnormalNodeWindow {
	winSize := nodeDur.Seconds() / timePerSlice.Seconds()
	if 0 == winSize || winSize < 1 {
		winSize = defWinSize
	}
	std.Printf("timePerSlice:%v,nodeDur:%v,winSize(defWinSize:%v):%v,\n", timePerSlice, nodeDur, defWinSize, winSize)
	slidingWindowInst := abnormalNodeWindow{}
	slidingWindowInst.Init(timePerSlice, int64(winSize))
	return &slidingWindowInst
}

//清理闲时的节点数据
func (s *abnormalNodeWindow) reset() {

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
func (s *abnormalNodeWindow) Init(timePerSlice time.Duration, winSize int64) {
	s.timePerSlice = timePerSlice
	s.winSize = winSize
	// 保证存储在至少两个window
	s.timeSliceSize = winSize*2 + 1

	s.timeSlices = make([]nodeBucket, s.timeSliceSize)
	s.preTs = time.Now()
	s.winDur = s.timePerSlice * time.Duration(s.winSize)
	go s.reset()
}

func (s *abnormalNodeWindow) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

//对时间片做增加操作，并返回窗口中所有的计数总和
func (s *abnormalNodeWindow) setAbnormalNode(node string) {
	var index = s.locationIndex()

	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, s.locationIndex())

	if oldCursor == index {
		// 在当前时间片里继续
		s.timeSlices[index].setNode(node)

	} else {
		s.timeSlices[index].setNode(node)

		// 清零，访问量不大时会有时间片跳跃的情况
		s.clearBetween(oldCursor, index)
	}

	s.preTs = time.Now()
}

//入参为tps的s
func (s *abnormalNodeWindow) getStats() []string {

	toSlice := func(m map[string]struct{}) []string {
		var rst []string
		for k := range m {
			rst = append(rst, k)
		}
		return rst
	}

	merge := func(m1 map[string]struct{}, m2 map[string]struct{}) {
		for k, v := range m2 {
			m1[k] = v
		}
	}

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

	rst := make(map[string]struct{})
	for i := int64(0); i < s.winSize; i++ {
		bucketTmp := s.timeSlices[(index-i+s.timeSliceSize)%s.timeSliceSize]
		merge(rst, bucketTmp.nodeMap)
	}

	return toSlice(rst)
}

func (s *abnormalNodeWindow) clearBetween(fromIndex, toIndex int64) {
	for index := (fromIndex + 1) % s.timeSliceSize; index != toIndex; index = (index + 1) % s.timeSliceSize {
		s.timeSlices[index].clear()
	}
}
