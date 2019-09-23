package conf

/*
	配置管理模块(全局);
	1. 管理服务框架配置;
	2. 管理用户自定义配置;
	3. 用户配置section作扁平化处理合并&字符串化处理;
	4. 用户配置以map[string]string托管及传递;
*/
import (
	"errors"
	"frame"
	xsf "git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"reflect"
	"strconv"
)

var SvcVersion string
var (
	// 服务框架通用配置
	Mock         bool   // mock开关; 开启后调用引擎空接口;
	EngSub       string // 服务类型
	SessMode     bool   // 会话模式
	Licence      int    // 并发授权数, 会话模式使用;
	NumaNode     int    // numa节点id, -1不设置亲和性;
	DelSessRt    int    // 会话管理器异步处理协程数;
	RealTimeRead bool   // 引擎同步读接口：是否实时读(写事件last读 || 边写边读);
	WrapperAsync bool   // 插件同步或异步模式：false同步, true异步;
	Catch        bool   // 异常捕获开关
	CatchDump    bool   // 异常捕获dump开关

	// 数据异步下行rabbitMQ信息
	RabbitHost	 string
	RabbitUser	 string
	RabbitPass	 string
	RabbitQueue	 string
	NrtDBUrl	 string	// 数据库状态接口地址

	// 用户自定义引擎配置
	UsrCfg     string            // 用户配置文件名
	UsrCfgData map[string]string // map<usrCfgKey, usrCfgVal>
)

func Construct(cfg *utils.Configure) (err error) {
	// 读取框架配置
	Licence, err = cfg.GetInt(sectionFc, maxLic)
	if err != nil {
		return
	}
	// 框架主配置
	err = secParseGes(cfg)
	if err != nil {
		return
	}
	// 用户自定义配置
	err = getUsrConfig(cfg)
	if err != nil {
		return
	}
	// 异步下行相关配置
	err = secParseDownAsync(cfg)
	return
}

func getUsrConfig(cfg *utils.Configure) (err error) {
	// 拉取用户自定义配置
	if len(UsrCfg) != 0 {
		cfgOpt := &utils.CfgOption{}
		cm := utils.CfgMode(*xsf.Mode)
		if cm == utils.Centre {
			// 使用远端配置;
			utils.WithCfgURL(*CmdCompanionUrl)(cfgOpt)
			utils.WithCfgPrj(*CmdProject)(cfgOpt)
			utils.WithCfgGroup(*CmdGroup)(cfgOpt)
			utils.WithCfgService(*CmdService)(cfgOpt)
		}
		utils.WithCfgName(UsrCfg)(cfgOpt)
		usrCfg, err := utils.NewCfg(cm, cfgOpt)
		if err != nil {
			return frame.ErrorGetUsrConfigure
		}

		// 遍历用户配置section(value类型->string)
		// 多个section时,扁平化处理配置项合并写入usrCfgData
		// eg:
		// [sec]
		// k1 = v1
		// k2 = v2
		// map存储格式如下<sec.k1, v1>, <sec.k2, v2>
		UsrCfgData = make(map[string]string)
		usrSecs := usrCfg.GetSecs()
		for _, sec := range usrSecs {
			secData := usrCfg.GetSection(sec)
			if secData != nil {
				kv, ok := secData.(map[string]interface{})
				if !ok {
					return frame.ErrorInvalidUsrCfg
				}
				for key, value := range kv {
					var valStr string
					switch value.(type) {
					case string:
						valStr = value.(string)
					case int:
						valStr = strconv.Itoa(value.(int))
					case int64:
						valStr = strconv.FormatInt(value.(int64), 10)
					case uint:
						valStr = strconv.FormatUint(uint64(value.(uint)), 10)
					case uint64:
						valStr = strconv.FormatUint(value.(uint64), 10)
					case bool:
						valStr = strconv.FormatBool(value.(bool))
					case float64:
						valStr = strconv.FormatFloat(value.(float64), 'f', -1, 64)
					case float32:
						valStr = strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)
					default:
						return errors.New("invalid user configure, type/sec/key " + reflect.TypeOf(value).String() + sec + key)
					}
					UsrCfgData[sec+"."+key] = valStr
				}
			}
		}
	} else {
		// 若用户配置为空则读取wrapper section;
		err = secParseWrapper(cfg)
	}
	return
}

func secParseGes(cfg *utils.Configure) (err error) {
	EngSub, err = cfg.GetString(sectionAiges, engSub)
	if err != nil {
		EngSub = defaultEngSud
	}

	var mock int
	Mock = false // 缺省：关闭mock
	mock, err = cfg.GetInt(sectionAiges, gesMock)
	if err == nil && mock != 0 {
		Mock = true
	}

	SessMode, err = cfg.GetBool(sectionAiges, sessMode) // TODO 考虑与框架fc::router合并;
	if err != nil {
		SessMode, err = true, nil // 缺省：会话模式
	}

	WrapperAsync, err = cfg.GetBool(sectionAiges, wrapperMode)
	if err != nil {
		// default sync wrapper
		WrapperAsync, err = false, nil
	}

	NumaNode, err = cfg.GetInt(sectionAiges, numaNode)
	if err != nil {
		// 缺省：不设置cpu亲和性
		NumaNode, err = defaultNumaNode, nil
	}

	DelSessRt, err = cfg.GetInt(sectionAiges, sessGort) // TODO 考虑与框架fc::best合并;
	if err != nil {
		// 缺省：与服务授权保持一致
		DelSessRt, err = Licence, nil
	}

	RealTimeRead, err = cfg.GetBool(sectionAiges, realTimeRlt)
	if err != nil {
		// 缺省：实时读/实时返回
		RealTimeRead, err = true, nil
	}

	Catch, err = cfg.GetBool(sectionAiges, catchSwitch)
	if err != nil {
		Catch, err = false, nil
	}
	CatchDump, err = cfg.GetBool(sectionAiges, catchDump)
	if err != nil {
		CatchDump, err = false, nil
	}
	// 存在引擎无需服务配置;
	UsrCfg, _ = cfg.GetString(sectionAiges, usrCfgName)
	return
}

func secParseWrapper(cfg *utils.Configure) (err error) {
	secData := cfg.GetSection(sectionWrapper)
	if secData != nil {
		kv, ok := secData.(map[string]interface{})
		if ok {
			UsrCfgData = make(map[string]string)
			for key, value := range kv {
				var valStr string
				switch value.(type) {
				case string:
					valStr = value.(string)
				case int:
					valStr = strconv.Itoa(value.(int))
				case int64:
					valStr = strconv.FormatInt(value.(int64), 10)
				case uint:
					valStr = strconv.FormatUint(uint64(value.(uint)), 10)
				case uint64:
					valStr = strconv.FormatUint(value.(uint64), 10)
				case bool:
					valStr = strconv.FormatBool(value.(bool))
				case float64:
					valStr = strconv.FormatFloat(value.(float64), 'f', -1, 64)
				case float32:
					valStr = strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)
				default:
					return errors.New("invalid wrapper configure, type/key " + reflect.TypeOf(value).String() + key)
				}
				UsrCfgData[key] = valStr
			}
		}
	}
	return
}

func secParseDownAsync(cfg *utils.Configure) (err error) {
	// rabbit mq
	RabbitHost, _ = cfg.GetString(sectionAsyncDown, rabHost)
	RabbitUser, _ = cfg.GetString(sectionAsyncDown, rabUser)
	RabbitPass, _ = cfg.GetString(sectionAsyncDown, rabPass)
	RabbitQueue, _ = cfg.GetString(sectionAsyncDown, rabQueue)
	NrtDBUrl,_ = cfg.GetString(sectionAsyncDown, nrtDBUrl)
	return
}
