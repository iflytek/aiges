package conf

import (
	"fmt"
	"git.xfyun.cn/AIaaS/finder-go"
	common "git.xfyun.cn/AIaaS/finder-go/common"
	"github.com/BurntSushi/toml"
	"os"
	"strings"
	"time"
	"schemas"
	"common/ratelimit2"
)


var configHandler   []func(string,[]byte)bool

func AddConfigChangerHander(f func(string,[]byte)bool)  {
	configHandler = append(configHandler,f)
}
var(
	findlerMamager *finder.FinderManager
	configChangeHandler common.ConfigChangedHandler
)
//集成配置中心与服务发现
func InitCentra() (b []byte) {
	cachePath, err := os.Getwd()

	if err != nil {
		return
	}

	//缓存信息的存放路径
	cachePath += "/findercache"
	config := common.BootConfig{
		//companion地址
		CompanionUrl: Centra.CompanionUrl,
		//缓存路径
		CachePath: cachePath,
		//是否缓存服务信息
		CacheService: true,
		//是否缓存配置信息
		CacheConfig:   true,
		ExpireTimeout: 5 * time.Second,
		MeteData: &common.ServiceMeteData{
			Project: Centra.Project,
			Group:   Centra.Group,
			Service: Centra.Service,
			Version: Centra.Version,
			Address: "",
		},
	}

	f, err := finder.NewFinderWithLogger(config, nil)

	if err != nil {
		ConsoleError("init finder manager error:"+err.Error())
		os.Exit(1)
	}
	configChangeHandler = &ConfigChangedHandle{}
	findlerMamager = f
	subRes, err := f.ConfigFinder.UseAndSubscribeConfig([]string{APP_CONFIG}, configChangeHandler)
	if err != nil {
		panic("subscribe file err:"+err.Error())
	}

	conf := subRes[APP_CONFIG]
	b = conf.File
	//加载routeMapping
	if _, err := toml.Decode(string(b), &Conf); err != nil {
		ConsoleError("cannot load app.config:"+err.Error())
		os.Exit(11)
	}
	Conf.Init()
	//订阅并且加载限流配置
	limitConf, err := f.ConfigFinder.UseAndSubscribeConfig([]string{LimitConf}, configChangeHandler)
	if err != nil {
		ConsoleWarn("subscribe limit config error:"+err.Error())
	}else{
		if lf,ok:=limitConf[LimitConf];ok && lf !=nil{
			err:=ratelimit2.LoadConfigCache(lf.File)
			if err !=nil{
				ConsoleWarn("load limit config error:"+err.Error())
			}
		}

	}
	err =schemas.LoadMapping(f,Conf.Schema.Services,configChangeHandler)

	if err !=nil{
		ConsoleError("load mapping error:"+err.Error())
		os.Exit(10)
	}

	return

}

func ConsoleError(v interface{})  {
	fmt.Println("ERROR:",time.Now().Format(time.RFC3339),v)
}

func ConsoleWarn(v interface{})  {
	fmt.Println("WARN:",time.Now().Format(time.RFC3339),v)
}



// ConfigChangedHandle ConfigChangedHandle
type ConfigChangedHandle struct {
}

// OnConfigFileChanged OnConfigFileChanged
func (s *ConfigChangedHandle) OnConfigFileChanged(config *common.Config) bool {
	if strings.HasSuffix(config.Name, ".toml") {
		fmt.Println(config.Name, " has changed:\r\n", string(config.File), " \r\n 解析后的map为 ：", config.ConfigMap)
	} else {
		fmt.Println(config.Name, " has changed:\r\n", string(config.File))
	}
	if config.Name == APP_CONFIG {
		_,err:=toml.Decode(string(config.File), &Conf)
		if err !=nil{
			ConsoleError("reload app.conf error:"+err.Error())
			return false
		}
		Conf.Init()
		err = schemas.LoadMapping(findlerMamager,Conf.Schema.Services,configChangeHandler)
		if err !=nil{
			ConsoleError("reload schema error:"+err.Error())
			return false
		}
	}


	if configHandler!=nil{
		for _,f:=range configHandler{
			if !f(config.Name,config.File){
				return false
			}
		}

	}
	return true
}

func (s *ConfigChangedHandle) OnError(errInfo common.ConfigErrInfo) {
	fmt.Println("配置文件出错：", errInfo)
}

type ServiceChangedHandle struct {
}

// OnServiceInstanceConfigChanged OnServiceInstanceConfigChanged
func (s *ServiceChangedHandle) OnServiceInstanceConfigChanged(name string, apiVersion string, instance string, config *common.ServiceInstanceConfig) bool {

	fmt.Println("服务实例配置信息更改开始，服务名：", name, "  版本号：", apiVersion, "  提供者实例为：", instance)
	fmt.Println("----当前配置为:  ", config.IsValid, "  ", config.UserConfig)
	fmt.Println("服务实例配置信息更改结束, 服务名：", name, "  版本号：", apiVersion, "  提供者实例为：", instance)
	config.IsValid = false
	config.UserConfig = "aasasasasasasa"
	config = nil
	return true
}

// OnServiceConfigChanged OnServiceConfigChanged
func (s *ServiceChangedHandle) OnServiceConfigChanged(name string, apiVersion string, config *common.ServiceConfig) bool {
	fmt.Println("服务配置信息更改开始，服务名：", name, "  版本号：", apiVersion)
	fmt.Println("-----当前配置为: ", config.JsonConfig)
	fmt.Println("服务配置信息更改结束, 服务名：", name, "  版本号：", apiVersion)
	config.JsonConfig = "zyssss"
	config = nil
	return true
}

// OnServiceInstanceChanged OnServiceInstanceChanged
func (s *ServiceChangedHandle) OnServiceInstanceChanged(name string, apiVersion string, eventList []*common.ServiceInstanceChangedEvent) bool {
	fmt.Println("服务实例变化通知开始, 服务名：", name, "  版本号：", apiVersion)
	for eventIndex, e := range eventList {
		for index, inst := range e.ServerList {
			if e.EventType == common.INSTANCEREMOVE {
				fmt.Println("----服务提供者节点减少事件 ：", e.ServerList)
				fmt.Println("-----------减少的服务提供者节点信息:  ")
				fmt.Println("----------------------- 地址: ", inst.Addr)
				fmt.Println("----------------------- 是否有效: ", inst.Config.IsValid)
				fmt.Println("----------------------- 配置: ", inst.Config.UserConfig)

			} else {
				fmt.Println("----服务提供者节点增加事件 ：", e.ServerList)
				fmt.Println("-----------增加的服务提供者节点信息:  ")
				fmt.Println("----------------------- 地址: ", inst.Addr)
				fmt.Println("----------------------- 是否有效: ", inst.Config.IsValid)
				fmt.Println("----------------------- 配置: ", inst.Config.UserConfig)

			}
			e.ServerList[index].Addr = "zy_tet"
			e.ServerList[index].Config = &common.ServiceInstanceConfig{}
		}
		eventList[eventIndex] = nil
	}

	fmt.Println("服务实例变化通知结束, 服务名：", name, "  版本号：", apiVersion)
	return true
}
