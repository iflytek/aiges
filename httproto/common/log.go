package common

import (
	"github.com/xfyun/xsf/client"
	"github.com/xfyun/xsf/utils"
	"sync/atomic"
	"unsafe"
)

//var (
//	Logger *xsf.Logger
//)

func InitXsfLog(path string, level string, logsize int, lognum int, caller bool, batch int, async bool) (*xsf.Logger, error) {
	if logsize == 0 {
		logsize = 100
	}
	if lognum == 0 {
		lognum = 10
	}
	var err error
	var logger *xsf.Logger
	if batch != 0 {
		logger, err = utils.NewLocalLog(utils.SetCaller(caller), utils.SetAsync(async), utils.SetLevel(level), utils.SetBatchSize(batch), utils.SetFileName(path), utils.SetMaxAge(30), utils.SetMaxSize(logsize), utils.SetMaxBackups(lognum))
	} else {
		logger, err = utils.NewLocalLog(utils.SetCaller(caller), utils.SetAsync(async), utils.SetLevel(level), utils.SetFileName(path), utils.SetMaxAge(30), utils.SetMaxSize(logsize), utils.SetMaxBackups(lognum))
	}
	if err != nil {
		return nil, err
	}
	return logger, nil
}

var loggerInst unsafe.Pointer

func GetLoggerInstance() *xsf.Logger {
	return (*xsf.Logger)(atomic.LoadPointer(&loggerInst))
}

func UpdateLogger(path string, level string, logsize int, lognum int, caller bool, batch int, async bool) error {
	logger, err := InitXsfLog(path, level, logsize, lognum, caller, batch, async)
	if err != nil {
		return err
	}
	atomic.StorePointer(&loggerInst, unsafe.Pointer(logger))
	return nil
}
