package daemon

import (
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

func pprofSrv(port int) {
	go func() {
		const localIp = "0.0.0.0:"
		if err := http.ListenAndServe(
			localIp+strconv.Itoa(port),
			nil,
		); nil != err {
			panic(err)
		}
	}()
}
