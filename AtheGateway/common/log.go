package common

import (
	"git.xfyun.cn/AIaaS/xsf-external/client"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
//	"conf"
	"github.com/gin-gonic/gin"
	"time"
	"schemas"
	"fmt"
)

var (
	Logger *xsf.Logger
)

func InitXsfLog(path string, level string, logsize int, lognum int, caller bool, batch int,async bool)error {
	if logsize == 0 {
		logsize = 100
	}
	if lognum == 0 {
		lognum = 10
	}
	var err error
	var logger *xsf.Logger
	if batch != 0 {
		logger, err = utils.NewLocalLog(utils.SetCaller(caller),utils.SetAsync(async), utils.SetLevel(level), utils.SetBatchSize(batch), utils.SetFileName(path), utils.SetMaxAge(30), utils.SetMaxSize(logsize), utils.SetMaxBackups(lognum))
	} else {
		logger, err = utils.NewLocalLog(utils.SetCaller(caller),utils.SetAsync(async), utils.SetLevel(level), utils.SetFileName(path), utils.SetMaxAge(30), utils.SetMaxSize(logsize), utils.SetMaxBackups(lognum))
	}
	if err != nil {
		fmt.Println("logger init failed .....................")
		return err
	}
	Logger = logger
	schemas.Logger = logger
	fmt.Println("logger inited........................")
	return nil
}

func Loggers()gin.HandlerFunc  {
	return func(context *gin.Context) {
		start:=time.Now()
		context.Next()
		end:=time.Now()
		Logger.Infof("clientIp:%s | time:%s | path:%s cost:%v",context.ClientIP(),end.Format(time.RFC3339),context.Request.URL.Path,end.Sub(start))
	}
}

