package server

import (
	"sync"
	"time"
)

var cache Cache

type Cache struct {
	sessionMg sync.Map
}

func Insert(k string, v *Session) {
	cache.sessionMg.Store(k, v)
}

func Remove(k string) {
	cache.sessionMg.Delete(k)
}

func RemoveAll(sids []string)  {
	for _,v:=range sids{
		Remove(v)
	}
}
func getCurrentSessionNum()int  {
	size:=0
	cache.sessionMg.Range(func(key, value interface{}) bool {
		size++
		return true
	})

	return size
}
func Get(k string) *Session {
	v, ok := cache.sessionMg.Load(k)
	if ok {
		return v.(*Session)
	}
	return nil
}
func InitCache(scanInterver time.Duration) {
	go checkConnTimeOut(scanInterver)
}

func checkConnTimeOut(scanInterver time.Duration) {
	for range time.Tick(scanInterver) {
		cache.sessionMg.Range(func(key, value interface{}) bool {
			s,ok := value.(*Session)
			if ok && s.checkSessionTimeOut() {
				s.Close()
				Remove(key.(string))
			}
			return true
		})
	}
}
