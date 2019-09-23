package server

import (
	"testing"
	"fmt"
	"sync"
	"time"
)
var total = 10000000
func TestNewSession(t *testing.T) {
	var s *Session
	for i:=0;i<total;i++{
		//go func() {
			s = &Session{}
		//}()

	}
	fmt.Println(s)
	time.Sleep(1*time.Second)
}

var pool = sync.Pool{}

func TestSession_Close(t *testing.T) {
	pool.New = func() interface{} {
		return &Session{
			Status:0,
		}
	}

	s:=pool.Get().(*Session)
	s.Status = 2
	pool.Put(s)
	fmt.Println(pool.Get().(*Session).Status)
	time.Sleep(1*time.Second)
}