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
	"fmt"
	"github.com/xfyun/aiges/frame"
	"github.com/xfyun/xsf/server"
	"github.com/xfyun/xsf/utils"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// 业务缺省计量字段配置;
var (
	meterIatList = "ent;caller.appid"
	meterTtsList = "vcn"
	meterIseList = "category"
)
var SvcVersion string
var (
	// 服务框架通用配置
	EngSub                   string           // 服务类型(eg:iat)
	Licence                  int              // 并发授权数, 会话模式使用;
	DelSessRt                int              // 会话管理器异步处理协程数;
	RealTimeRead             bool             // 引擎同步读接口：是否实时读(写事件last读 || 边写边读);
	WrapperAsync             bool             // 插件同步或异步模式：false同步, true异步;
	HttpRetry                int      = 1     // http下载重试,缺省3次
	GrayLabel                bool             // 集群节点灰度状态标记
	WrapperTrace             bool             // 插件回调trace日志开关
	HeaderPass               []string         // 可放行至wrapper的header参数
	WrapperDelayDetectPeriod int              //框架判断引擎接口卡死的周期,单位 秒
	StorageData              bool     = false //默认不保存数据
	AsyncRelease             bool     = true  //是否异步释放会话

	// pprof
	PProfAble bool = false
	PProfHost string
	PProfPort int = 1234 // default pprof port

	// 数据异步下行rabbitMQ信息
	RabbitHost  string
	RabbitUser  string
	RabbitPass  string
	RabbitQueue string
	RabbitRetry int    = 3
	NrtDBUrl    string // 数据库状态接口地址

	// 数据处理模块
	ReSampleAble bool = true

	// 用户自定义引擎配置
	UsrCfg     string            // 用户配置文件名
	UsrCfgData map[string]string // map<usrCfgKey, usrCfgVal>
	// GRPC Python解释器
	PythonCmd string = "cc"
)

func Construct(cfg *utils.Configure) (err error) {
	// 读取框架配置
	if Licence, err = cfg.GetInt(sectionFc, maxLic); err != nil {
		return
	}
	// 框架主配置
	if err = secParseGes(cfg); err != nil {
		return
	}
	// pprof设置
	if err = secParsePProf(cfg); err != nil {
		return
	}
	// 用户自定义配置
	if err = getUsrConfig(cfg); err != nil {
		return
	}
	// 数据处理相关配置
	if err = secParseDP(cfg); err != nil {
		return
	}
	// 异步下行相关配置
	if err = secParseDownAsync(cfg); err != nil {
		return
	}
	// 环境变量配置化
	if err = parseLoaderEnv(); err != nil {
		return
	}

	return
}

func secParseDP(cfg *utils.Configure) (err error) {
	if able, err := cfg.GetBool(sectionDp, dpRsAble); err == nil {
		ReSampleAble = able
	}
	return
}

func secParsePProf(cfg *utils.Configure) (err error) {
	if able, err := cfg.GetBool(sectionPProf, pprofAble); err == nil {
		PProfAble = able // default false
	}

	if PProfAble {
		if port, err := cfg.GetInt(sectionPProf, pprofPort); err == nil {
			PProfPort = port // default port:1234
		}

		if PProfHost, _ = cfg.GetString(sectionPProf, pporfHost); len(PProfHost) == 0 {
			PProfHost = xsf.GetNetaddr()
		}
		go func() {
			log.Println(http.ListenAndServe(PProfHost+":"+strconv.Itoa(PProfPort), nil))
		}()
	}

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

	WrapperAsync, err = cfg.GetBool(sectionAiges, wrapperMode)
	if err != nil {
		// default sync wrapper
		WrapperAsync, err = false, nil
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
	// http download retry times
	if retry, err := cfg.GetInt(sectionAiges, httpRetry); err == nil {
		HttpRetry = retry
	}

	WrapperTrace, err = cfg.GetBool(sectionAiges, wrapperTrace)
	if err != nil {
		WrapperTrace, err = true, nil
	}
	GrayLabel, err = cfg.GetBool(sectionAiges, grayMark)
	if err != nil {
		GrayLabel, err = false, nil
	}

	if pass, err := cfg.GetString(sectionAiges, headerPass); err == nil {
		HeaderPass = strings.Split(pass, ";")
	}
	// 本地列表读取/差异化部分迁移至线下，保障上线流程自动化的一致性;
	if passfile, err := ioutil.ReadFile(localPassFile); err == nil {
		pass := strings.Split(string(passfile), "\n")
		HeaderPass = append(HeaderPass, pass...)
	}
	fmt.Println("header pass list:", HeaderPass)

	//默认卡死检测 延迟为60s不出结果
	if delayDetectPeriod, err := cfg.GetInt(sectionAiges, wrapperDelayDetectPeriod); err == nil {
		if delayDetectPeriod <= 0 {
			WrapperDelayDetectPeriod = 60
		} else {
			WrapperDelayDetectPeriod = delayDetectPeriod
		}
	} else {
		WrapperDelayDetectPeriod = 60
	}
	if storage, err := cfg.GetBool(sectionAiges, storageData); err == nil {
		StorageData = storage
	}

	if asyncOr, err := cfg.GetBool(sectionAiges, asyncRelease); err == nil {
		AsyncRelease = asyncOr
	}
	// 存在引擎无需服务配置;
	UsrCfg, _ = cfg.GetString(sectionAiges, usrCfgName)

	if pcmd, err := cfg.GetString(sectionAiges, pythonPluginCmd); err == nil {
		PythonCmd = pcmd
	}
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
	NrtDBUrl, _ = cfg.GetString(sectionAsyncDown, nrtDBUrl)
	if retry, err := cfg.GetInt(sectionAsyncDown, rabRetry); err == nil {
		RabbitRetry = retry
	}
	return
}
