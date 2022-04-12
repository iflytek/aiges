package common

import "time"

type ServiceMeteData struct {
	Project string
	Group   string
	Service string
	Version string
	Address string
}

func (meta *ServiceMeteData) Check() bool {
	if len(meta.Group) == 0 || len(meta.Project) == 0 /**|| len(meta.Version) == 0 || len(meta.Service) == 0 **/{
		return false
	}
	return true
}

type BootConfig struct {
	CompanionUrl  string
	CachePath     string
	CacheConfig   bool
	CacheService  bool
	ExpireTimeout time.Duration
	MeteData      *ServiceMeteData
}

type StorageInfo struct {
	Addr            []string
	ConfigRootPath  string
	ServiceRootPath string
	ZkNodePath      string
}

type Config struct {
	//文件名
	Name string
	//文件内容
	File []byte
	//toml文件解析后的数据
	ConfigMap map[string]interface{}
}
type ConfigErrInfo struct {
	FileName string
	ErrCode  int
	ErrMsg   string
}
type ServiceInstanceConfig struct {
	IsValid    bool `json:"is_valid"`
	UserConfig string
}

type ConsumerInstanceConfig struct {
	IsValid bool `json:"is_valid"`
}

type ServiceInstanceChangedEvent struct {
	EventType  InstanceChangedEventType
	ServerList []*ServiceInstance
}

type ServiceInstance struct {
	Addr   string
	Config *ServiceInstanceConfig
}

func (ints *ServiceInstance) Dumplication() *ServiceInstance {
	var dump ServiceInstance
	dump.Addr = ints.Addr
	dump.Config = &ServiceInstanceConfig{IsValid: ints.Config.IsValid, UserConfig: ints.Config.UserConfig}
	return &dump
}

type ServiceConfig struct {
	JsonConfig string
}
type RouteItem struct {
	Id       string
	Consumer []string
	Provider []string
	Only     string
	Name     string
}
type ServiceRoute struct {
	RouteItem []*RouteItem
}
type Service struct {
	ServiceName  string
	ApiVersion   string
	ProviderList []*ServiceInstance
	Config       *ServiceConfig
}

type ServiceInfo struct {
	ApiVersion   string
	ProviderList []string
}

func (service *Service) Dumplication() Service {
	var dump Service
	dump.ApiVersion = service.ApiVersion
	dump.ServiceName = service.ServiceName
	var providerList = make([]*ServiceInstance, len(service.ProviderList))
	for index, provider := range service.ProviderList {
		providerList[index] = &ServiceInstance{Addr: provider.Addr, Config: &ServiceInstanceConfig{IsValid: provider.Config.IsValid, UserConfig: provider.Config.UserConfig}}
	}
	dump.ProviderList = providerList
	dump.Config = &ServiceConfig{JsonConfig: service.Config.JsonConfig}
	return dump
}

type ServiceSubscribeItem struct {
	ServiceName string
	ApiVersion  string
}
type ConfigFeedback struct {
	PushID       string
	ServiceMete  *ServiceMeteData
	Config       string
	UpdateTime   int64
	UpdateStatus int
	LoadTime     int64
	LoadStatus   int
	GrayGroupId  string
}

type ServiceFeedback struct {
	PushID          string
	ServiceMete     *ServiceMeteData
	Provider        string
	ProviderVersion string
	UpdateTime      int64
	UpdateStatus    int
	LoadTime        int64
	LoadStatus      int
	Type            int
}
