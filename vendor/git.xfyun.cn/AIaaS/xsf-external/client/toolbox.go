package xsf

import (
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"os"
)

var logSidGeneratorInst utils.LogSidGenerator

var dbgLoggerStd = newDbsLoggerStd("debug=>", false)

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
func (d *DbgLoggerStd) recF(format string, v ...interface{}) {
	if !d.able {
		return
	}
	d.logger.Printf(format, v...)
}
func (d *DbgLoggerStd) recLn(v ...interface{}) {
	if !d.able {
		return
	}
	d.logger.Println(v...)
}
