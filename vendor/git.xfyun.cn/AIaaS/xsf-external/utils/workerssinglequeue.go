/*
* @file	workersSingleQueue.go
* @brief	simple goroutine pool (refers to fasthttp)
* @author	sqjian
* @version	1.0
* @date		2017.11.4
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
        +-----> |  task queue  | <----+
                +------+-------+
                       |
                       |
    +-------------+----v------+------------+
    |             |           |            |
    |             |           |            |
    |             |           |            |
    |             |           |            |
    |             |           |            |
    |             |           |            |
+---v----+   +----v---+   +---v----+   +---v----+
| worker |   | worker |   | worker |   | worker |
+--------+   +--------+   +--------+   +--------+
*/
package utils

import (
	"sync"
)

type Worker interface {
	Task()
}

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown 等待所有 goroutine 停止工作
func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
