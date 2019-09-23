package xsf

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"os"
	"strings"
)

type ToolBoxServer struct {
}

func (t *ToolBoxServer) Cmdserver(ctx context.Context, in *utils.Request) (*utils.Response, error) {
	queryMap, headersMap := make(map[string]string), make(map[string]string)
	if err := json.Unmarshal([]byte(in.Query), &queryMap); nil != err {
		return &utils.Response{}, fmt.Errorf("query:%v,err:%v", in.Query, err)
	}

	if err := json.Unmarshal([]byte(in.Headers), &headersMap); nil != err {
		return &utils.Response{}, fmt.Errorf("headers:%v,err:%v", in.Headers, err)
	}
	if "GET" != strings.ToUpper(headersMap["method"]) {
		return &utils.Response{}, errors.New("don't support the method")
	}
	buf := bytes.NewBuffer(nil)
	cmdServerRouter(queryMap["cmd"], queryMap, buf)
	return &utils.Response{Body: buf.String()}, nil
}

var loggerStd = (&utils.LoggerStderr{}).Init("")
var dbgLoggerStd = newDbsLoggerStd("debuging ", false)

func init() {
	xsfDbg := os.Getenv("XSF-DEBUG")
	if "1" == xsfDbg {
		dbgLoggerStd = newDbsLoggerStd("debug=>", true)
	}
}

type DbgLoggerStd struct {
	logger *utils.LoggerStderr
	able   bool
}

func newDbsLoggerStd(prefix string, able bool) *DbgLoggerStd {
	l := DbgLoggerStd{
		logger: (&utils.LoggerStderr{}).Init(prefix),
		able:   able,
	}
	return &l
}
func (d *DbgLoggerStd) Printf(format string, v ...interface{}) {
	if !d.able {
		return
	}
	d.logger.Printf(format, v...)
}
func (d *DbgLoggerStd) Println(v ...interface{}) {
	if !d.able {
		return
	}
	d.logger.Println(v...)
}
