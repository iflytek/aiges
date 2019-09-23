package finder

import "log"

type ServiceChangedHandler interface {

	//服务实例上的配置信息发生变化
	OnServiceInstanceConfigChanged(name string,apiVersion string, addr string, config *ServiceInstanceConfig) bool
	//服务整体配置信息发生变化
	OnServiceConfigChanged(name string,apiVersion string,  config *ServiceConfig) bool
	//服务实例发生变化
	OnServiceInstanceChanged(name string, apiVersion string, eventList []*ServiceInstanceChangedEvent) bool
}

type ConfigChangedHandler interface {
	OnConfigFileChanged(config *Config) bool
	OnError(errInfo ConfigErrInfo)
}

type InternalServiceChangedHandler interface {
	OnServiceInstanceConfigChanged(name string, addr string, data []byte)
	OnServiceConfigChanged(name string, data []byte)
	OnServiceInstanceChanged(name string, addrList []string)
}

type InternalConfigChangedHandler interface {
	OnConfigFileChanged(name string, data []byte)
}



func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

type DefaultLogger struct {
}


