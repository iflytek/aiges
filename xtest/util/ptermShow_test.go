package util

import (
	"go.uber.org/atomic"
	"testing"
)

func TestProgressShow(t *testing.T) {
	cnt := atomic.NewInt64(1000000000)
	go func() {
		for i := 0; i < 1000000000; i++ {
			cnt.Add(-1)
		}
	}()
	ProgressShow(cnt)

}
