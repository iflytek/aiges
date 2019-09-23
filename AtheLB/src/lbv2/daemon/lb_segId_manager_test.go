package daemon

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestSegIdManager_getMin(t *testing.T) {
	var concurrent int64 = 10
	var wg sync.WaitGroup
	var cnt, cntTmp int64 = 1e4, 0
	var sm sync.Map
	for ix := int64(0); ix < concurrent; ix++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for atomic.LoadInt64(&cntTmp) < atomic.LoadInt64(&cnt) {
				sm.Store(segIdManagerInst.getMin(), true)
				atomic.AddInt64(&cntTmp, 1)
			}
		}()
	}
	wg.Wait()
	smLen := func() (rst int64) {
		sm.Range(func(key, value interface{}) bool {
			atomic.AddInt64(&rst, 1)
			return true
		})
		return
	}()
	std.Printf("cntTmp:%d\n", cntTmp)
	std.Printf("smLen:%d\n", smLen)
}

//49.67
