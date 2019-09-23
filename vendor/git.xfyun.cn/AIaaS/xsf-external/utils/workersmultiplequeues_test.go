package utils

import (
	"fmt"
	"testing"
)

type taskDemoEx struct{}

func (m *taskDemoEx) Task() {
	fmt.Printf("Goroutine ID : %v\n", GetId())
}

func TestWorkersEx(t *testing.T) {
	np := taskDemo{}
	p := NewEx(20)
	for cnt := 0; cnt < 1000; cnt++ {
		p.Run(&np, cnt)
	}
	p.Shutdown()
}
