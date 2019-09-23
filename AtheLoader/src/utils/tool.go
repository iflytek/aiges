/*
	aiges工具箱,用于实现部分通用基础操作;
	包含：服务器资源利用率获取,系统资源配置获取等操作;
*/
package utils

import (
	"runtime"
	"strconv"
	"strings"
)

func GetGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		return -1
	}
	return id
}
