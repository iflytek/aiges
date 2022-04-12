/*
* @file	workersMultipleQueues.go
* @brief	simple goroutine pool
* @author	sqjian
* @version	1.0
* @date		2017.11.14
*/

/*
      +----------+   +----------+   +----------+
      | receiver |   | receiver |   | receiver |
      +-----+----+   +-----+----+   +-----+----+
            |              |              |
            |              |              |
            |              |              |
            |              |              |
            |       +------v-------+      |
            +-----> |  dispatcher  | <----+
                    +------+-------+
                           |
       +------------+------v------+------------+
       |            |             |            |
+------v---+  +-----v----+  +-----v----+  +----v-----+
|task queue|  |task queue|  |task queue|  |task queue|
+----+-----+  +----+-----+  +-----+----+  +----+-----+
     |             |              |            |
     |             |              |            |
     |             |              |            |
 +---v----+    +---v----+    +----v---+    +---v----+
 | worker |    | worker |    | worker |    | worker |
 +--------+    +--------+    +--------+    +--------+

*/
package utils

import (
	"sync"
	"strings"
	"strconv"
	"runtime"
)

type WorkerEx interface {
	Task()
}

type PoolEx struct {
	workMap map[int]chan WorkerEx
	wg      sync.WaitGroup
}

func NewEx(maxGoroutines int) *PoolEx {
	workMapTmp := map[int]chan WorkerEx{}
	for i := 0; i < maxGoroutines; i++ {
		workMapTmp[i] = make(chan WorkerEx)
	}
	p := PoolEx{
		workMap: workMapTmp,
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func(ix int) {
			for w := range p.workMap[ix] {
				w.Task()
			}
			p.wg.Done()
		}(i)
	}

	return &p
}

func (p *PoolEx) Run(w WorkerEx, id int) {
	id = id % len(p.workMap)
	p.workMap[id] <- w
}

// Shutdown 等待所有 goroutine 停止工作
func (p *PoolEx) Shutdown() {
	for _, v := range p.workMap {
		close(v)
	}
	p.wg.Wait()
}
//获取
func GetId() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		return -1
	}
	return id
}
