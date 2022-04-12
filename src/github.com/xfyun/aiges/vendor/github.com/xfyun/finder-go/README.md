# 配置中心概述
1 配置中心提供两个重要的功能！统一配置管理 和服务发现

## 统一配置存储
1. 配置中心存放了服务组件必要的启动配置文件。 把原本组件启动时依赖的配置文件统一存放到配置中心中，便于统一维护和管理。
2. 组件启动时可以从配置中心订阅自己的配置文件，当订阅的配置文件被修改时，组件也会收到对应的通知。

## 服务发现
1. 服务发现功能类似dns。 可以通过服务发现功能，使用目标服务的serviceName 来订阅目标服务。获取目标服务的地址，从而使用。并且当目标服务的地址有变动时，
订阅该服务的组件也会收到通知。从而达到一个高可用性。

## 使用方式： 通过集成sdk来使用。 本工程为配置中心golang的sdk （finder-go）。

1.首先需要在配置中心管理控制台创建服务的配置.：流程如下2,3,4
2.找到服务所在集群，并在集群下创建服务和服务版本号。
3. 上传自己的配置文件
4. 开始订阅配置，使用方式如下：
## 配置中心SDK finder-go 使用说明

1. 创建config

```
config := common.BootConfig{
		//companion地址
		CompanionUrl: conf.CompanionUrl,
		//缓存路径
		CachePath: "",
		//是否缓存服务
		CacheService: true,
		//是否缓存配置信息
		CacheConfig:   true,
		ExpireTimeout: 5 * time.Second,
		MeteData: &common.ServiceMeteData{
			Project: "test",   // 项目名称
			Group:   "test", // 项目集群，只用同一个集群下面的服务才能通过服务发现，发现彼此的地址
			Service: "test", // 自己组件的名称
			Version: "1.0.1",  // 自己组件的版本号
			Address: "127.0.0.1:1221",
		},
	}
```
2. 创建finder

```
f, err := finder.NewFinderWithLogger(config, nil)

```
3. 使用接口

```
//实现handler接口 订阅服务 返回的map的key是ServiceName + "_" + ApiVersion。 vlaue是Service实例，主要取ProviderList。代表当前服务的提供者列表
serviceList, err := f.ServiceFinder.UseAndSubscribeServic(subscri, handler)
//注册服务
f.ServiceFinder.RegisterServiceWithAddr(addr, apiVersion)

//获取配置文件
configFiles, err := f.ConfigFinder.UseAndSubscribeConfig(name, &handler)

```

* handler接口说明

```
type ConfigChangedHandler interface {
    //配置文件发生改变后的回调
	OnConfigFileChanged(config *Config) bool
	OnError(errInfo ConfigErrInfo)
}

type ServiceChangedHandler interface {

	//服务实例上的配置信息发生变化 回调接口
	OnServiceInstanceConfigChanged(name string,apiVersion string, addr string, config *ServiceInstanceConfig) bool
	//服务整体配置信息发生变化 回调接口
	OnServiceConfigChanged(name string,apiVersion string,  config *ServiceConfig) bool
	//服务实例发生变化回调接口，ServiceInstanceChangedEvent中的EventType代表增加提供者还是减少，ServerList代表增加的实例或者减少的实例
	OnServiceInstanceChanged(name string, apiVersion string, eventList []*ServiceInstanceChangedEvent) bool
}

```

###  查询所有服务

```

f.ServiceFinder.QueryService("AIaaS", "dx")

//返回值分析
map[string][]common.ServiceInfo 
key : serviceName
value : 列表 common.ServiceInfo 

type ServiceInfo struct {
	ApiVersion   string   //版本号
	ProviderList []string //所有地址，可能为空
}
```

### 流程
* 参加[流程图](https://github.com/xfyun/finder-go/v3/blob/master/%E9%85%8D%E7%BD%AE%E4%B8%AD%E5%BF%83%E6%B5%81%E7%A8%8B.png)

### 2.1.19 更新：

1. 集成时需要新增configChangeHandler 的接口实现函数：
````go
func (s *ConfigChangedHandle) OnConfigFilesAdded(configs map[string]*common.Config) bool {


	return true
}

func (s *ConfigChangedHandle) OnConfigFilesRemoved(configNames []string) bool {


	return true
}

````

2. github.com/xfyun/finder-go/v3/common 包名由原来错误的finder 修正为common。


###  c语言支持，执行该脚本：
1. [create.sh](./cgo/create.sh) 生成libfinder.so  和libfinder.h 
2. 使用时包含两个头文件 libfinder.h 和 [config_center.h](./cgo/config_center.h)
3. 使用demo见[test_config.c](./cgo/example/test_config.c) 和 [test_service.c](./cgo/example/test_service.c) 


