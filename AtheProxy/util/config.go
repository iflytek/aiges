package util

import (
	"fmt"
	"config"
	"consts"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"strings"
)

var CfgOption *utils.CfgOption
//加载配置文件
func LoadConfig() {
	var co = &utils.CfgOption{}
	cfgName := consts.CONFIG_FILE
	if config.UseCfgCentre == 0 {
		cfgName = consts.PREFIX + consts.CONFIG_FILE
	}

	utils.WithCfgName(cfgName)(co)
	if config.UseCfgCentre == 1 {
		//从云端拉取配置
		utils.WithCfgURL(config.CompanionUrl)(co)
		utils.WithCfgPrj(config.Project)(co)
		utils.WithCfgGroup(config.Group)(co)
		utils.WithCfgService(config.Service)(co)
		utils.WithCfgVersion(config.Version)(co)
	}

	//读取配置文件
	cfg, err := utils.NewCfg(utils.CfgMode(config.UseCfgCentre), co)
	if err != nil {
		fmt.Println("load config error: ", err)
		panic(err)
	}

	CfgOption = co

	//读取log参数
	getLogParams(cfg)

	//读取统一服务参数
	getCommonParams(cfg)
}

/**
获取atmos公共参数配置
 */
func getCommonParams(cfg *utils.Configure) {
	debugSwitch, e := cfg.GetInt("atmos-common", "debugSwitch")
	if e == nil && debugSwitch == consts.OPEN {
		fmt.Println("getParams", "debugSwitch:", debugSwitch)
		OpenPprof()
	}

	subTypes, e := cfg.GetString("atmos-common", "subTypes")
	if e == nil {
		subs := strings.Split(subTypes, ",")
		config.SUB_TYPES = subTypes
		for _, sub := range subs {
			if sub == "" {
				continue
			}
			config.SUB_MAP[sub] = GetParamsBySub(sub, cfg)
			fmt.Println("sub:", config.SUB_MAP[sub])
		}
	}
}

/**
	根据sub读取参数
 */
func GetParamsBySub(sub string, cfg *utils.Configure) *config.AtmosConfig {

	title := "atmos-" + sub
	lb, e := cfg.GetString(title, "lb")
	if e != nil {
		fmt.Println(sub, "get lb err:", e)
	}

	defaultTimeout, e := cfg.GetInt(title, consts.DEFAULT_TIMEOUT)
	if e != nil {
		fmt.Println(sub, "get defaultTimeout err:", e)
		defaultTimeout = config.TIMEOUT
	}

	engineTimeout, e := cfg.GetInt(title, "engineTimeout")
	if e != nil {
		fmt.Println(sub, "get engineTimeout err:", e)
		engineTimeout = config.TIMEOUT
	}

	onceTimeout, e := cfg.GetInt(title, "onceTimeout")
	if e != nil {
		fmt.Println(sub, "get onceTimeout err:", e)
		onceTimeout = config.TIMEOUT
	}

	getEngineResultTimeout, e := cfg.GetInt(title, "getEngineResultTimeout")
	if e != nil {
		fmt.Println(sub, "get getEngineResultTimeout err:", e)
		getEngineResultTimeout = config.TIMEOUT * 3
	}

	engineRetry, e := cfg.GetInt(title, "engineRetry")
	if e != nil {
		engineRetry = config.ENGINE_RETRY
	}

	enableMock, e := cfg.GetInt(title, "enableMock")
	if e != nil {
		enableMock = consts.CLOSE
	}
	defaultRoute, e := cfg.GetString(title, "defaultRoute")
	if e != nil {
		defaultRoute = ""
	}

	//引擎参数过滤
	engineSkipParmMap := make(map[string]string)
	engineSkipParm, e := cfg.GetString(title, "engineSkipParm")
	if e == nil {
		parms := strings.Split(engineSkipParm, ",")
		for _, v := range parms {
			engineSkipParmMap[v] = ""
		}
	}

	//config.SUB_MAP[sub] =
	atmosConfig :=
		&config.AtmosConfig{
			sub,
			lb,
			defaultTimeout,
			onceTimeout,
			engineTimeout,
			getEngineResultTimeout,
			engineRetry,
			defaultRoute,
			enableMock,
			engineSkipParmMap,
		}
	return atmosConfig
}

//读取log参数
func getLogParams(cfg *utils.Configure) {
	var (
		v int
		s string
		e error
	)
	s, e = cfg.GetString("log", "level")
	if e == nil {
		config.LogLevel = s
	}
	s, e = cfg.GetString("log", "file")
	if e == nil {
		config.LogFile = s
	}
	v, e = cfg.GetInt("log", "size")
	if e == nil {
		config.LogSize = v
	}
	v, e = cfg.GetInt("log", "count")
	if e == nil {
		config.LogCount = v
	}
	v, e = cfg.GetInt("log", "die")
	if e == nil {
		config.LogDie = v
	}
	v, e = cfg.GetInt("log", "async")
	if e == nil && v == 0 {
		config.LogAsync = false
	}
	v, e = cfg.GetInt("log", "cache")
	if e == nil {
		config.LogCache = v
	}
	v, e = cfg.GetInt("log", "batch")
	if e == nil {
		config.LogBatch = v
	}
	v, e = cfg.GetInt("log", "caller")
	if e == nil && v == 1 {
		config.LogCaller = true
	}

}
