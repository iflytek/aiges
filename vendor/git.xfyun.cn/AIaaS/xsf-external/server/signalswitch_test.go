package xsf

import (
	"fmt"
	"sync"
	"testing"
)

type killed1 struct {
}

func (k *killed1) Closeout() {
	fmt.Println("be killed1.")
}

type killed2 struct {
}

func (k *killed2) Closeout() {
	fmt.Println("be killed2.")
}
func TestSignalHandle(t *testing.T) {
	addKillerCheck(killerNormalPriority, "killed1", &killed1{})
	addKillerCheck(killerNormalPriority, "killed2", &killed2{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		signalHandle()
	}()
	wg.Wait()

}
