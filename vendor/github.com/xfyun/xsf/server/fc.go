package xsf

import (
	"log"
	"sync"
)

type fcDelay struct {
	rwmu    sync.RWMutex
	smDelay []func() error
}

var fcDelayInst fcDelay
var fcDelayOnce sync.Once

func (s *fcDelay) add(in func() error) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()
	s.smDelay = append(s.smDelay, in)
}
func (s *fcDelay) exec() {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()
	for _, v := range s.smDelay {
		fcDelayOnce.Do(func() {
			loggerStd.Println("about to call fc delay task")
		})
		if err := v(); err != nil {
			log.Panic(err)
		}
	}
}
