package config

import (
	"consts"
)

// xrpc配置
var (
	UseCfgCentre = 1
	CompanionUrl = "http://companion.xfyun.iflytek:6868"
	Project      = "AIaaS"
	Group        = "aitest"
	Service      = "atmos"
	Version      = "1.0.0"
	ServiceHost  = consts.LOCAL_IP
	ServicePort  = consts.DEFAULT_PORT
)

// 日志配置
var (
	LogFile   = "/log/server/atmos.log"
	LogLevel  = "warn"
	LogSize   = 100
	LogCount  = 20
	LogDie    = 10
	LogAsync  = true
	LogCache  = -1
	LogBatch  = 16000
	LogCaller = false
)

type AtmosConfig struct {
	Sub                    string
	Lb                     string
	DefaultTimeout         int
	OnceTimeout            int
	EngineTimeout          int
	GetEngineResultTimeout int
	EngineRetry            int
	DefaultRoute           string
	EnableMock             int
	EngineSkipParmMap      map[string]string
}

//统一服务配置
var (
	SUB_MAP   = make(map[string]*AtmosConfig)
	SUB       = "svc"
	SUB_TYPES = "svc"
	LB        = "svc_lb"
	//默认超时时间3000毫秒
	TIMEOUT = 3000

	//单次调用引擎接口超时时间
	ENGINE_TIMEOUT = 3000

	//单位秒，获取引擎结果的时候，最大等待时间(可能包含多次引擎接口调用 )
	GET_ENGINE_RESULT_TIMEOUT = 15

	//0表示关，1表示开
	RECOVER_SWITCH = consts.OPEN

	ENGINE_RETRY = 3

	//0表示关，1表示开
	DEBUG_SWITCH = consts.CLOSE
)