package xsf

import (
	"container/heap"
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

func newItem(valueIn string, priorityIn int) *Item {
	return &Item{value: func() atomic.Value { var value atomic.Value; value.Store(valueIn); return value }(), priority: int64(priorityIn)}
}

type Item struct {
	value    atomic.Value
	priority int64
	index    int64
}

func (I *Item) set(value string, priority int64) {
	I.value.Store(value)
	atomic.StoreInt64(&I.priority, priority)
}

type Items []*Item

func (I Items) Len() int {
	return len(I)
}
func (I Items) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
func (I Items) Less(i, j int) bool {
	return I[i].priority < I[j].priority
}
func (I *Items) Delete(i int) {
	if i >= len(*I) {
		return
	}
	for j := i; j < len(*I)-1; j++ {
		(*I)[j] = (*I)[j+1]
	}
	*I = (*I)[:len(*I)-1]
}

type LoadCollectorInterface interface {
	/*
		item 为收纳结构，value和priority分别为addr和优先级
	*/
	update(item *Item, value string, priority int) bool
	load(n int) []string
	store(item *Item) bool
	delete(item *Item) bool
	getItem(value string) *Item
}

type Queue struct {
	m     map[string]*Item //addr、item
	mRwMu sync.RWMutex

	s     Items
	sRwMu sync.RWMutex
}

func newQueue(itemsIn []*Item) *Queue {
	var rst Queue
	rst.m = make(map[string]*Item)
	for _, v := range itemsIn {
		rst.m[v.value.Load().(string)] = v
	}
	rst.s = itemsIn
	return &rst
}
func (q *Queue) getItem(value string) *Item {
	q.mRwMu.RLock()
	defer q.mRwMu.RUnlock()
	return q.m[value]
}
func (q *Queue) update(item *Item, value string, priority int) bool {
	q.mRwMu.RLock()
	itemTmp, itemTmpOk := q.m[value]
	q.mRwMu.RUnlock()
	if !itemTmpOk {
		return false
	}
	itemTmp.value.Store(value)
	dbgLoggerStd.Printf("fn:update,value:%v,priority:%v\n", value, priority)
	atomic.StoreInt64(&(itemTmp.priority), int64(priority))
	return true
}
func (q *Queue) delete(item *Item) bool {
	q.mRwMu.Lock()
	defer q.mRwMu.Unlock()
	q.sRwMu.Lock()
	defer q.sRwMu.Unlock()

	delete(q.m, item.value.Load().(string))

	value := item.value.Load().(string)
	var i int
	var j *Item
	for i, j = range q.s {
		if j.value.Load().(string) == value {
			continue
		}
	}
	if i != 0 {
		q.s.Delete(i)
	}
	return true
}
func (q *Queue) load(n int) []string {
	dbgLoggerStd.Printf("fn:load,n:%v\n", n)
	if n != 0 {
		return q.topK(n)
	}
	return q.all()
}
func (q *Queue) all() []string {
	var rst []string
	q.sRwMu.Lock()
	fuck.Println("the q.sRwMu.Lock")
	sort.Sort(q.s)
	for _, v := range q.s {
		rst = append(rst, v.value.Load().(string))
	}
	q.sRwMu.Unlock()
	return rst
}

func (q *Queue) topK(k int) []string {
	var rst []string
	signMap := make(map[atomic.Value]struct{})
	q.sRwMu.RLock()
	for i := 0; i < k; i++ {
		var maxValue atomic.Value
		var maxPriority int64 = math.MinInt64
		for _, v := range q.s {
			/*
				1、初始时随机选择
			*/
			if 0 == v.priority {
				_, ok := signMap[v.value]
				if ok {
					continue
				}
				maxPriority = v.priority
				maxValue = v.value
				break
			}
			/*
				1、v.priority > maxPriority正常情况确保最大值
			*/
			if v.priority > maxPriority {
				_, ok := signMap[v.value]
				if ok {
					continue
				}
				maxPriority = v.priority
				maxValue = v.value
			}
		}
		maxValueString, maxValueStringOk := maxValue.Load().(string)
		if maxValueStringOk {
			signMap[maxValue] = struct{}{}
			rst = append(rst, maxValueString)
		} else {
		}
		dbgLoggerStd.Printf("fn:topK,maxValue is nil\n")
	}
	q.sRwMu.RUnlock()

	//dbgLoggerStd.Printf("fn:topK,k:%v,topKRst:%v,allRst:%v\n", k, rst, func() []map[string]int64 {
	//	var tmp []map[string]int64
	//	q.sRwMu.RLock()
	//	defer q.sRwMu.RUnlock()
	//	for _, v := range q.s {
	//		m := map[string]int64{
	//			v.value.Load().(string): v.priority,
	//		}
	//		tmp = append(tmp, m)
	//	}
	//	return tmp
	//}())
	return rst
}
func (q *Queue) store(item *Item) bool {
	value := item.value.Load().(string)
	q.mRwMu.RLock()
	itemTmp, itemOk := q.m[value]
	q.mRwMu.RUnlock()
	if !itemOk {
		q.mRwMu.Lock()
		q.m[item.value.Load().(string)] = item
		q.mRwMu.Unlock()

		q.sRwMu.Lock()
		q.s = append(q.s, item)
		q.sRwMu.Unlock()
	} else {
		itemTmp.set(value, item.priority)
	}
	return true
}

/*
each variable is used
*/
type PriorityQueue struct {
	data []*Item
	rwMu sync.RWMutex
}

func (pq *PriorityQueue) Len() int {
	pq.rwMu.RLock()
	defer pq.rwMu.RUnlock()
	return len(pq.data)
}

func (pq *PriorityQueue) Less(i, j int) bool {
	pq.rwMu.RLock()
	defer pq.rwMu.RUnlock()
	return pq.data[i].priority < pq.data[j].priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	pq.rwMu.Lock()
	defer pq.rwMu.Unlock()
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
	pq.data[i].index = int64(i)
	pq.data[j].index = int64(j)
}

func (pq *PriorityQueue) Push(x interface{}) {
	pq.rwMu.Lock()
	defer pq.rwMu.Unlock()
	n := len(pq.data)
	item := x.(*Item)
	item.index = int64(n)
	pq.data = append(pq.data, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	pq.rwMu.RLock()
	defer pq.rwMu.RUnlock()
	old := pq
	n := len(old.data)
	item := old.data[n-1]
	item.index = -1
	pq.data = old.data[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, value string, priority int) bool {
	if nil == item {
		return false
	}
	item.value.Store(value)
	item.priority = int64(priority)
	heap.Fix(pq, int(item.index))
	return true
}
func (pq *PriorityQueue) load(n int) []string {
	if n != 0 {
		return pq.topK(n)
	}
	return pq.all()
}
func (pq *PriorityQueue) topK(k int) []string {
	pq.rwMu.RLock()
	defer pq.rwMu.RUnlock()

	ix := 0
	var rst []string
	for _, v := range pq.data {
		rst = append(rst, v.value.Load().(string))
		ix++
		if ix >= k {
			break
		}
	}
	return rst
}
func (pq *PriorityQueue) all() []string {
	pq.rwMu.RLock()
	defer pq.rwMu.RUnlock()

	var rst []string
	for _, v := range pq.data {
		rst = append(rst, v.value.Load().(string))
	}
	return rst
}
func (pq *PriorityQueue) getItem(value string) *Item {
	pq.rwMu.RLock()
	defer pq.rwMu.RUnlock()
	for _, v := range pq.data {
		if v.value.Load().(string) == value {
			return v
		}
	}
	return nil
}
func (pq *PriorityQueue) store(item *Item) bool {
	pq.Push(item)
	return true
}
func (pq *PriorityQueue) delete(item *Item) bool {
	rst := pq.load(0)
	hasFlag := false
	value := item.value.Load().(string)
	for _, v := range rst {
		if v == value {
			hasFlag = true
			continue
		}
	}
	if !hasFlag {
		return true
	}
	heap.Remove(pq, int(item.index))
	return true
}

type cellCalculatorCell struct {
	errCode int64
	dur     int64
	vCpu    int64
}
type cellCalculator struct {
	healthData  *healthWindow
	latencyData *latencyWin
	vCpu        int64
}

func newCellCalculator(timePerSlice time.Duration, winSize int64) *cellCalculator {
	healthWindowInst := newHealthWindow(timePerSlice, winSize)
	latencyWindowInst := newLatencyWindow(timePerSlice, winSize)
	cellCalculatorInst := cellCalculator{
		healthData:  healthWindowInst,
		latencyData: latencyWindowInst,
	}
	return &cellCalculatorInst
}

//load=(1-health)*(100*cpu-latency)
func (c *cellCalculator) sync(cell cellCalculatorCell) {
	c.healthData.setErrCode(cell.errCode)
	c.latencyData.addDur(cell.dur)
	if 0 == cell.vCpu {
		return
	}
	c.vCpu = cell.vCpu
}
func (c *cellCalculator) calc() float64 {

	var health float64
	healthErrCount, healthCount := c.healthData.getStats()
	if 0 == healthCount {
		health = 0
	} else {
		health = float64(healthErrCount) / float64(healthCount)
	}

	var latency float64
	_, _, latencySum, latencyCount := c.latencyData.getStats()
	if 0 == latencyCount {
		latency = 0
	} else {
		latency = float64(latencySum) / float64(latencyCount)
	}

	dbgLoggerStd.Printf("fn:calc,vCpu:%v,health:%v,latency:%v\n", c.vCpu, health, latency)
	calcLoad := (1 - health) * (100*float64(c.vCpu) - latency)

	return calcLoad
}

type LoadCalculator struct {
	winSize      int64
	timePerSlice time.Duration

	cellCalculatorMap     map[string]*cellCalculator
	cellCalculatorMapRwMu sync.RWMutex
}

func newLoadCalculator(timePerSlice time.Duration, winSize int64) *LoadCalculator {
	LoadCalculatorInst := LoadCalculator{
		winSize:           winSize,
		timePerSlice:      timePerSlice,
		cellCalculatorMap: make(map[string]*cellCalculator),
	}
	return &LoadCalculatorInst
}

func (l *LoadCalculator) sync(target string, cell cellCalculatorCell) {
	l.cellCalculatorMapRwMu.RLock()
	cellCalculatorInst, cellCalculatorInstOk := l.cellCalculatorMap[target]
	l.cellCalculatorMapRwMu.RUnlock()

	if !cellCalculatorInstOk {
		l.cellCalculatorMapRwMu.Lock()
		cellCalculatorInst = newCellCalculator(l.timePerSlice, l.winSize)
		l.cellCalculatorMap[target] = cellCalculatorInst
		l.cellCalculatorMapRwMu.Unlock()
	}
	cellCalculatorInst.sync(cell)
}
func (l *LoadCalculator) syncWithLoad(target string, cell cellCalculatorCell) int64 {
	l.cellCalculatorMapRwMu.RLock()
	cellCalculatorInst, cellCalculatorInstOk := l.cellCalculatorMap[target]
	l.cellCalculatorMapRwMu.RUnlock()

	if !cellCalculatorInstOk {
		l.cellCalculatorMapRwMu.Lock()
		fuck.Println("the l.cellCalculatorMapRwMu.Lock")
		cellCalculatorInst = newCellCalculator(l.timePerSlice, l.winSize)
		l.cellCalculatorMap[target] = cellCalculatorInst
		l.cellCalculatorMapRwMu.Unlock()
	}
	dbgLoggerStd.Printf("fn:syncWithLoad,cell:%#v\n", cell)
	cellCalculatorInst.sync(cell)

	calcLoad := cellCalculatorInst.calc()
	load := int64(calcLoad * 100)

	//保留两位小数
	return load
}
func (l *LoadCalculator) calc(target string) float64 {
	l.cellCalculatorMapRwMu.RLock()
	cellCalculatorInst, cellCalculatorInstOk := l.cellCalculatorMap[target]
	l.cellCalculatorMapRwMu.RUnlock()
	if cellCalculatorInstOk {
		return cellCalculatorInst.calc()
	}
	return -1
}
