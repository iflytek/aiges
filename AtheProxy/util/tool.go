package util

import (
	"log"
	"runtime"
	"net/http"
	"config"
	"consts"
	"strings"
)

//是否打开调试端口
func OpenPprof() {
	go func() {
		runtime.SetBlockProfileRate(1)
		log.Fatal(http.ListenAndServe(":8089", nil))
	}()
}

//是否需要显示info日志
func IsNeedShowInfoLog() bool {
	return strings.EqualFold(config.LogLevel, consts.LOG_LEVEL_INFO) || strings.EqualFold(config.LogLevel, consts.LOG_LEVEL_DEBUG)
}

//是否展示debug日志
func IsNeedShowDebugLog() bool {
	return strings.EqualFold(config.LogLevel, consts.LOG_LEVEL_DEBUG)
}