/*
	aiges工具箱,用于实现部分通用基础操作;
	包含：服务器资源利用率获取,系统资源配置获取等操作;
*/
package utils

import (
	"bytes"
	"encoding/binary"
	"github.com/xfyun/xsf/utils"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var CommonLogger *utils.Logger

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

func StorageDecodeData(pcm []byte, sid string) (err error) {
	fp, err := os.Create(sid)
	if err != nil {
		return
	}
	defer fp.Close()
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, pcm)
	fp.Write(buf.Bytes())
	return nil
}

func SetCommonLogger(xsflog *utils.Logger) {
	CommonLogger = xsflog
}
