package utils

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

var pprofOnce sync.Once

func pprofSrv(ip, port string) {
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%v:%v", ip, port), nil); err != nil {
			panic(err)
		}
	}()
}
func StartPProf(ip, port string) {
	pprofOnce.Do(func() {
		pprofSrv(ip, port)
	})
}
