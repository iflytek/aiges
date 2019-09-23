package util

import (
	"config"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

var SugarLog *utils.Logger

//初始化日志
func InitLogger() {
	logger, err := utils.NewLocalLog(
		utils.SetCaller(config.LogCaller),
		utils.SetLevel(config.LogLevel),
		utils.SetFileName(config.LogFile),
		utils.SetMaxSize(config.LogSize),
		utils.SetMaxBackups(config.LogCount),
		utils.SetMaxAge(config.LogDie),
		utils.SetAsync(config.LogAsync),
		utils.SetCacheMaxCount(config.LogCache),
		utils.SetBatchSize(config.LogBatch))
	if err != nil {
		panic(err)
	}
	SugarLog = logger
}
