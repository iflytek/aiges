package finder

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"

	common "git.xfyun.cn/AIaaS/finder-go/common"
	errors "git.xfyun.cn/AIaaS/finder-go/errors"
	"git.xfyun.cn/AIaaS/finder-go/log"
	"git.xfyun.cn/AIaaS/finder-go/route"
	"git.xfyun.cn/AIaaS/finder-go/storage"
	"git.xfyun.cn/AIaaS/finder-go/utils/serviceutil"
	"git.xfyun.cn/AIaaS/finder-go/utils/stringutil"
	"strings"
)

type ServiceFinder struct {
	locker            sync.Mutex
	rootPath          string
	config            *common.BootConfig
	handler           common.ServiceChangedHandler
	storageMgr        storage.StorageManager
	usedService       map[string]*common.Service
	subscribedService map[string]*common.Service
	serviceZkData     map[string]*ServiceZkData
	mutex             sync.Mutex
}

type ServiceZkData struct {
	ServiceName string
	ApiVersion  string
	//所有的提供者 key是addr
	ProviderList map[string]*common.ServiceInstance
	Config       *common.ServiceConfig
	Route        *common.ServiceRoute
}

func NewServiceFinder(root string, bc *common.BootConfig, sm storage.StorageManager) *ServiceFinder {
	finder := &ServiceFinder{
		locker:            sync.Mutex{},
		rootPath:          root,
		config:            bc,
		storageMgr:        sm,
		usedService:       make(map[string]*common.Service, 0),
		subscribedService: make(map[string]*common.Service, 0),
		serviceZkData:     make(map[string]*ServiceZkData, 0),
	}

	return finder
}

func (f *ServiceFinder) RegisterServiceWithAddr(addr string, version string) error {
	if f.storageMgr == nil {
		return errors.NewFinderError(errors.ZkConnectionLoss)
	}
	log.Log.Debugf("RegisterServiceWithAddr : addr-> %s %s %s", addr, " version->", version)
	return f.registerService(addr, version)
}
func (f *ServiceFinder) RegisterService(version string) error {
	if f.storageMgr == nil {
		return errors.NewFinderError(errors.ZkConnectionLoss)
	}
	log.Log.Debugf("RegisterServiceWithAddr : version-> %s %s %s", version)
	return f.registerService(f.config.MeteData.Address, version)
}
func (f *ServiceFinder) UnRegisterService(version string) error {
	if f.storageMgr == nil {
		return errors.NewFinderError(errors.ZkConnectionLoss)
	}
	servicePath := fmt.Sprintf("%s/%s/%s/provider/%s", f.rootPath, f.config.MeteData.Service, version, f.config.MeteData.Address)
	return f.storageMgr.RemoveInRecursive(servicePath)
}

func (f *ServiceFinder) UnRegisterServiceWithAddr(version string, addr string) error {
	if f.storageMgr == nil {
		return errors.NewFinderError(errors.ZkConnectionLoss)
	}
	servicePath := fmt.Sprintf("%s/%s/%s/provider/%s", f.rootPath, f.config.MeteData.Service, version, addr)

	return f.storageMgr.RemoveInRecursive(servicePath)
}

func (f *ServiceFinder) UseService(serviceItems []common.ServiceSubscribeItem) (map[string]*common.Service, error) {
	var err error
	if len(serviceItems) == 0 {
		err = errors.NewFinderError(errors.ServiceMissItem)
		return nil, err
	}

	f.locker.Lock()
	defer f.locker.Unlock()

	serviceList := make(map[string]*common.Service)
	for _, item := range serviceItems {
		//这个usedService 是作何用处？
		serviceId := item.ServiceName + "_" + item.ApiVersion

		//测试用
		servicePath := fmt.Sprintf("%s/%s/%s", f.rootPath, item.ServiceName, item.ApiVersion)
		log.Log.Infof(" useservice: %s %s %s", servicePath)
		serviceList[serviceId], err = f.getService(servicePath, item)
		//存入缓存文件
		if serviceList[serviceId] == nil || serviceList[serviceId].ProviderList == nil || len(serviceList[serviceId].ProviderList) == 0 {
			log.Log.Debugf("the service is null ")
			service, err := GetServiceFromCache(f.config.CachePath, item)
			if err != nil || service == nil {
				log.Log.Infof("query service from cache err：%v", err)
				f.subscribedService[serviceId] = &common.Service{ServiceName: item.ServiceName, ApiVersion: item.ApiVersion}
				continue
			} else {
				var tempServer = service.Dumplication()
				serviceList[serviceId] = &tempServer
				f.subscribedService[serviceId] = service
			}
		}

		err = CacheService(f.config.CachePath, serviceList[serviceId])
		if err != nil {
			log.Log.Errorf("CacheService failed")
		}

	}

	return serviceList, err
}

func (f *ServiceFinder) UseAndSubscribeService(serviceItems []common.ServiceSubscribeItem, handler common.ServiceChangedHandler) (map[string]common.Service, error) {

	var err error
	if len(serviceItems) == 0 {
		err = errors.NewFinderError(errors.ServiceMissItem)
		return nil, err
	}

	f.locker.Lock()
	defer f.locker.Unlock()
	f.handler = handler
	serviceList := make(map[string]common.Service)

	if f.storageMgr == nil {
		if !f.config.CacheService {
			log.Log.Infof(" [ UseAndSubscribeService ] not use cache")
			return nil, nil
		}
		log.Log.Infof(" [ UseAndSubscribeService ] get service from cache")

		//说明zk信息目前有误，暂时使用缓存数据
		for _, item := range serviceItems {
			serviceId := item.ServiceName + "_" + item.ApiVersion
			service, err := GetServiceFromCache(f.config.CachePath, item)
			if err != nil {
				log.Log.Infof("query service from cache err：", err)
				f.subscribedService[serviceId] = &common.Service{ServiceName: item.ServiceName, ApiVersion: item.ApiVersion}

			} else {
				serviceList[serviceId] = service.Dumplication()
				f.subscribedService[serviceId] = service
			}

		}
		return serviceList, nil
	}
	for _, item := range serviceItems {
		serviceId := item.ServiceName + "_" + item.ApiVersion
		if service, ok := f.subscribedService[serviceId]; ok {
			serviceList[serviceId] = service.Dumplication()
			continue
		}
		servicePath := fmt.Sprintf("%s/%s/%s", f.rootPath, item.ServiceName, item.ApiVersion)
		service, err := f.getServiceWithWatcher(servicePath, item, handler)
		if err != nil {
			log.Log.Infof(" [ UseAndSubscribeService ] subscribe service %v, version %v , err :%v ", item.ServiceName, item.ApiVersion, err)
			continue
		}
		if service == nil {
			continue
		}
		serviceList[serviceId] = service.Dumplication()
		f.subscribedService[serviceId] = service

		err = f.registerConsumer(item, f.config.MeteData.Address)
		if err != nil {
			log.Log.Errorf("registerConsumer failed, %s", err)
		}
		CacheService(f.config.CachePath, f.subscribedService[serviceId])
	}
	return serviceList, nil
}

func (f *ServiceFinder) UnSubscribeService(name string) error {
	var err error
	if len(name) == 0 {
		//	err = errors.NewFinderError(errors.ServiceMissName)
		return err
	}
	f.locker.Lock()
	defer f.locker.Unlock()
	delete(f.subscribedService, name)
	return nil
}

func (f *ServiceFinder) QueryServiceWatch(project, group string, handler common.ServiceChangedHandler) (map[string][]common.ServiceInfo, error) {
	if len(project) == 0 || len(group) == 0 {
		return nil, errors.NewFinderError(errors.InvalidParam)
	}

	rootPath := "/polaris/service/" + fmt.Sprintf("%x", md5.Sum([]byte(project+group)))
	var serMap = make(map[string][]common.ServiceInfo)
	pC := NewQueryServiceCallback(handler, f)
	//sC:=NewQueryServiceCallback(handler,WATCH_SERVICE,f)
	//vC:=NewQueryServiceCallback(handler,WATCH_VERSION,f)

	//watch所有的server
	if sers, err := f.storageMgr.GetChildrenWithWatch(rootPath, &pC); err != nil {
		return nil, err
	} else {
		for _, ser := range sers {
			pC.serviceCache = append(pC.serviceCache, ser)
			if vers, err := f.storageMgr.GetChildrenWithWatch(rootPath+"/"+ser, &pC); err == nil {
				for _, ver := range vers {
					pC.versionCache[ser] = append(pC.versionCache[ser], ver)
					var item common.ServiceInfo
					item.ApiVersion = ver
					if providers, err := f.storageMgr.GetChildrenWithWatch(rootPath+"/"+ser+"/"+ver+"/provider", &pC); err == nil {
						item.ProviderList = providers
					}
					var finalprovider []string
					for _, provider := range item.ProviderList {
						data, err := f.storageMgr.GetDataWithWatch(rootPath+"/"+ser+"/"+ver+"/provider/"+provider, &pC)
						serviceInstance := new(common.ServiceInstance)
						//解析数据
						if data == nil || len(data) == 0 || err != nil {
							//获取数据为空
							log.Log.Infof("get data from %v is empty :", rootPath+"/"+ser+"/"+ver+"/provider/"+provider)
							serviceInstance.Config = getDefaultServiceInstanceConfig()
						} else {
							//获取的提供者配置数据不为空
							var item []byte
							_, item, err = common.DecodeValue(data)
							if err != nil {
								log.Log.Infof("service instance data is %v,unmarsh err: %v", string(data), err)
								//使用默认的配置
								serviceInstance.Config = getDefaultServiceInstanceConfig()
							} else {
								serviceInstance.Config = serviceutil.ParseServiceConfigData(item)
							}

						}
						if serviceInstance != nil && serviceInstance.Config.IsValid {
							finalprovider = append(finalprovider, provider)
						}
					}
					item.ProviderList = finalprovider
					serMap[ser] = append(serMap[ser], item)
				}
			}
		}
	}
	pC.provciderCache = serMap
	return serMap, nil
}
func (f *ServiceFinder) QueryService(project, group string) (map[string][]common.ServiceInfo, error) {
	if len(project) == 0 || len(group) == 0 {
		return nil, errors.NewFinderError(errors.InvalidParam)
	}
	rootPath := "/polaris/service/" + fmt.Sprintf("%x", md5.Sum([]byte(project+group)))
	var serMap = make(map[string][]common.ServiceInfo)
	if sers, err := f.storageMgr.GetChildren(rootPath); err != nil {
		return nil, err
	} else {
		for _, ser := range sers {
			if vers, err := f.storageMgr.GetChildren(rootPath + "/" + ser); err == nil {
				for _, ver := range vers {
					var item common.ServiceInfo
					item.ApiVersion = ver
					if providers, err := f.storageMgr.GetChildren(rootPath + "/" + ser + "/" + ver + "/provider"); err == nil {
						item.ProviderList = providers
					}
					var finalprovider []string
					for _, provider := range item.ProviderList {
						pc, _ := getServiceInstance(f.storageMgr, rootPath+"/"+ser+"/"+ver+"/provider", provider, nil)
						if pc != nil && pc.Config.IsValid {
							finalprovider = append(finalprovider, provider)
						}
					}
					item.ProviderList = finalprovider
					serMap[ser] = append(serMap[ser], item)
				}
			}
		}
	}
	return serMap, nil
}
func (f *ServiceFinder) registerService(addr string, apiVersion string) error {
	if stringutil.IsNullOrEmpty(addr) {
		err := errors.NewFinderError(errors.ServiceMissAddr)
		return err
	}
	if stringutil.IsNullOrEmpty(apiVersion) {
		log.Log.Infof("[registerService] apiversion not exist")
		return errors.NewFinderError(errors.ServiceMissApiVersion)
	}
	//目前不考虑目录不存在的情况
	path := fmt.Sprintf("%s/%s/%s/provider/%s", f.rootPath, f.config.MeteData.Service, apiVersion, addr)
	log.Log.Debugf("registerService -> path -> %s", path)
	err := f.storageMgr.SetTempPath(path)
	if err != nil {
		log.Log.Infof("registerService err %v", err)
		return err
	}
	go pushService(f.config.CompanionUrl, f.config.MeteData.Project, f.config.MeteData.Group, f.config.MeteData.Service, apiVersion)
	//if err != nil {
	//	logger.Error("RegisterService->registerService:", err)
	//}
	return nil
}

func (f *ServiceFinder) registerConsumer(service common.ServiceSubscribeItem, addr string) error {
	if stringutil.IsNullOrEmpty(addr) {
		err := errors.NewFinderError(errors.ServiceMissAddr)
		log.Log.Errorf("registerConsumer: %s", err)
		return err
	}

	parentPath := fmt.Sprintf("%s/%s/%s/consumer", f.rootPath, service.ServiceName, service.ApiVersion)
	err := f.register(parentPath, addr)
	if err != nil {
		log.Log.Errorf("registerConsumer->register: %s", err)
		return err
	}

	return nil
}
func (f *ServiceFinder) getServiceInstanceByAddrList(providerAddrList []string, rootPath string, handler *ServiceChangedCallback) []*common.ServiceInstance {
	var serviceInstanceList = make([]*common.ServiceInstance, 0)
	for _, providerAddr := range providerAddrList {
		log.Log.Debugf(" [ getServiceInstanceByAddrList] providerAddr: %s %s %s", providerAddr, " rootPath :", rootPath)
		service, err := getServiceInstance(f.storageMgr, rootPath, providerAddr, handler)
		if err != nil || service == nil {
			continue
		}
		serviceInstanceList = append(serviceInstanceList, service)
	}
	return serviceInstanceList
}
func (f *ServiceFinder) register(parentPath string, addr string) error {
	log.Log.Infof("call register func")
	servicePath := parentPath + "/" + addr
	log.Log.Infof("servicePath:  %s", servicePath)
	return f.storageMgr.SetTempPath(servicePath)
}

func getDefaultServiceItemConfig(addr string) ([]byte, error) {
	defaultServiceInstanceConfig := common.ServiceInstanceConfig{
		IsValid: true,
	}

	data, err := json.Marshal(defaultServiceInstanceConfig)
	if err != nil {
		log.Log.Errorf("%s", err)
		return nil, err
	}

	var encodedData []byte
	encodedData, err = common.EncodeValue("", data)
	if err != nil {
		log.Log.Errorf("%s", err)
		return nil, err
	}

	return encodedData, nil
}

func getDefaultConsumerItemConfig(addr string) ([]byte, error) {
	defaultConsumeInstanceConfig := common.ConsumerInstanceConfig{
		IsValid: true,
	}

	data, err := json.Marshal(defaultConsumeInstanceConfig)
	if err != nil {
		log.Log.Errorf("%s", err)
		return nil, err
	}

	var encodedData []byte
	encodedData, err = common.EncodeValue("", data)
	if err != nil {
		log.Log.Errorf("%s", err)
		return nil, err
	}

	return encodedData, nil
}

func getServiceInstance(sm storage.StorageManager, path string, addr string, callback *ServiceChangedCallback) (*common.ServiceInstance, error) {
	var data []byte
	var err error
	if callback != nil {
		data, err = sm.GetDataWithWatch(path+"/"+addr, callback)
	} else {
		data, err = sm.GetData(path + "/" + addr)
	}
	if err != nil {
		log.Log.Infof("get data from %v err %v:", path+"/"+addr, err)
		//TODO 是否需要返回默认的
		return nil, err
	}
	serviceInstance := new(common.ServiceInstance)
	//解析数据
	if data == nil || len(data) == 0 {
		//获取数据为空
		log.Log.Infof("get data from %v is empty :", path+"/"+addr)
		serviceInstance.Config = getDefaultServiceInstanceConfig()
	} else {
		//获取的提供者配置数据不为空
		var item []byte
		_, item, err = common.DecodeValue(data)
		if err != nil {
			log.Log.Infof("service instance data is %v,unmarsh err: %v", string(data), err)
			//使用默认的配置
			serviceInstance.Config = getDefaultServiceInstanceConfig()
		} else {
			serviceInstance.Config = serviceutil.ParseServiceConfigData(item)
		}

	}
	serviceInstance.Addr = addr
	return serviceInstance, nil
}
func getDefaultServiceInstanceConfig() *common.ServiceInstanceConfig {
	serviceInstanceConfig := &common.ServiceInstanceConfig{}
	serviceInstanceConfig.IsValid = true
	serviceInstanceConfig.UserConfig = ""
	return serviceInstanceConfig
}

func (f *ServiceFinder) getService(servicePath string, serviceItem common.ServiceSubscribeItem) (*common.Service, error) {
	var service = &common.Service{ServiceName: serviceItem.ServiceName, ApiVersion: serviceItem.ApiVersion, ProviderList: make([]*common.ServiceInstance, 0)}
	var serviceZkData = &ServiceZkData{ServiceName: serviceItem.ServiceName, ApiVersion: serviceItem.ApiVersion, ProviderList: make(map[string]*common.ServiceInstance)}
	f.serviceZkData[serviceItem.ServiceName+"_"+serviceItem.ApiVersion] = serviceZkData
	var providerPath = servicePath + "/provider"
	var confPath = servicePath + "/conf"
	var routePath = servicePath + "/route"
	//先找provider路径下的数据

	providerList, err := f.storageMgr.GetChildren(providerPath)
	if err != nil {
		if strings.Compare("zk: node does not exist", err.Error()) == 0 {
			//节点不存在，则新建之
			err := f.storageMgr.SetPath(providerPath)
			if err != nil {
				log.Log.Infof("[ GetChildrenWithWatch ] create node: %s", providerPath)
			}
			return nil, err
		}
		log.Log.Infof("query service provider from %v, err: %v", providerPath, err)
		return nil, err
	}
	if len(providerList) == 0 {
		log.Log.Infof("[ getServiceWithWatcher ] current not have service provider")
	}
	for _, providerAddr := range providerList {
		serviceInstance, err := getServiceInstance(f.storageMgr, providerPath, providerAddr, nil)
		if err != nil {
			//TODO 当data为nil的时候，会返回错误。。这里要处理一下
			log.Log.Infof("query instance info err，path=: %v , err: %v", providerPath+"/"+providerAddr, err)
			// todo
			continue
		}
		serviceZkData.ProviderList[serviceInstance.Addr] = serviceInstance
		//如果该提供者被禁用了，则跳过
		if serviceInstance.Config != nil && !serviceInstance.Config.IsValid {
			continue
		}
		service.ProviderList = append(service.ProviderList, serviceInstance)
	}

	//获取config下的信息
	confData, err := f.storageMgr.GetData(confPath)
	if err != nil {

		log.Log.Infof("query config data err from %v , err %v ", confPath, err)
		if strings.Compare(common.ZK_NODE_DOSE_NOT_EXIST, err.Error()) == 0 {
			log.Log.Infof("create node: %v", confPath)
			f.storageMgr.SetPath(confPath)
		}
		service.Config = &common.ServiceConfig{JsonConfig: ""}

	} else if len(confData) == 0 {
		service.Config = &common.ServiceConfig{JsonConfig: ""}
		log.Log.Infof("query data from path: %v ,data is empty ", confPath)
	} else {
		_, fData, err := common.DecodeValue(confData)
		if err != nil {
			log.Log.Infof("parse data errr  %s", err)
		}
		service.Config = &common.ServiceConfig{JsonConfig: string(fData)}
		serviceZkData.Config = &common.ServiceConfig{JsonConfig: string(fData)}
	}

	//获取route数据
	routeData, err := f.storageMgr.GetData(routePath)
	if err != nil {
		log.Log.Infof("get route data from path %v ,err %v", routePath, err)
		if strings.Compare(common.ZK_NODE_DOSE_NOT_EXIST, err.Error()) == 0 {
			log.Log.Infof("create node: %s", routePath)
			f.storageMgr.SetPath(routePath)
		}
		serviceZkData.Route = &common.ServiceRoute{RouteItem: []*common.RouteItem{}}
	} else if routeData != nil && len(routeData) == 0 {
		log.Log.Infof("get route data from path %v is empty", routePath)
		serviceZkData.Route = &common.ServiceRoute{RouteItem: []*common.RouteItem{}}
	} else {
		_, fData, err := common.DecodeValue(routeData)
		if err != nil {
			log.Log.Infof("parse route data err %s ", err)
			serviceZkData.Route = &common.ServiceRoute{RouteItem: []*common.RouteItem{}}
		} else {
			serviceZkData.Route = route.ParseRouteData(fData)
		}
		//使用route进行过滤数据
		service.ProviderList = route.FilterServiceByRouteData(serviceZkData.Route, f.config.MeteData.Address, service.ProviderList)
	}

	return service, nil
}

func (f *ServiceFinder) getServiceWithWatcher(servicePath string, serviceItem common.ServiceSubscribeItem, handler common.ServiceChangedHandler) (*common.Service, error) {
	var service = &common.Service{ServiceName: serviceItem.ServiceName, ApiVersion: serviceItem.ApiVersion, ProviderList: make([]*common.ServiceInstance, 0)}

	var serviceZkData = &ServiceZkData{ServiceName: serviceItem.ServiceName, ApiVersion: serviceItem.ApiVersion, ProviderList: make(map[string]*common.ServiceInstance)}
	log.Log.Infof("zk data %v ", f.serviceZkData[serviceItem.ServiceName+"_"+serviceItem.ApiVersion])
	//if f.serviceZkData[serviceItem.ServiceName+"_"+serviceItem.ApiVersion] == nil {
	f.serviceZkData[serviceItem.ServiceName+"_"+serviceItem.ApiVersion] = serviceZkData
	//}
	var providerPath = servicePath + "/provider"
	var confPath = servicePath + "/conf"
	var routePath = servicePath + "/route"
	//先找provider路径下的数据
	callback := NewServiceChangedCallback(serviceItem, SERVICE_INSTANCE_CHANGED, f, handler)
	//获取数据的时候添加子节点变更的Watcher
	providerList, err := f.storageMgr.GetChildrenWithWatch(providerPath, &callback)

	//TODO 提供者为空的情况
	if err != nil {
		if strings.Compare("zk: node does not exist", err.Error()) == 0 {
			//节点不存在，则新建之
			err := f.storageMgr.SetPath(providerPath)
			if err != nil {
				log.Log.Infof("[ GetChildrenWithWatch ] creade node: %v", providerPath)
			}
			return nil, err
		}
		log.Log.Infof("query service provider from %v, err: %v", providerPath, err)
		return nil, err
	}
	if len(providerList) == 0 {
		log.Log.Infof(" [ getServiceWithWatcher ]current  service provider is emtpy %v", "")
	}
	for _, providerAddr := range providerList {
		proiderCallBack := NewServiceChangedCallback(serviceItem, SERVICE_INSTANCE_CONFIG_CHANGED, f, handler)
		serviceInstance, err := getServiceInstance(f.storageMgr, providerPath, providerAddr, &proiderCallBack)
		if err != nil {
			//TODO 当data为nil的时候，会返回错误。。这里要处理一下
			log.Log.Infof("query service instance err: %v, path: %v", err, providerPath+"/"+providerAddr)
			// todo
			continue
		}
		serviceZkData.ProviderList[serviceInstance.Addr] = serviceInstance
		//如果该提供者被禁用了，则跳过
		if serviceInstance.Config != nil && !serviceInstance.Config.IsValid {
			continue
		}
		service.ProviderList = append(service.ProviderList, serviceInstance)
	}

	//获取config下的信息
	confCallBack := NewServiceChangedCallback(serviceItem, SERVICE_CONFIG_CHANGED, f, handler)
	confData, err := f.storageMgr.GetDataWithWatch(confPath, &confCallBack)
	if err != nil {

		log.Log.Infof("query config data err from %v, err: %v", confPath, err)
		if strings.Compare(common.ZK_NODE_DOSE_NOT_EXIST, err.Error()) == 0 {
			log.Log.Infof("create node: %s", confPath)
			f.storageMgr.SetPath(confPath)
		}
		service.Config = &common.ServiceConfig{JsonConfig: ""}

	} else if len(confData) == 0 {
		service.Config = &common.ServiceConfig{JsonConfig: ""}
		log.Log.Infof("get config data from path: %v is emtpy", confPath)
	} else {
		_, fData, err := common.DecodeValue(confData)
		if err != nil {
			log.Log.Infof("parse data err %v", err)
		}
		service.Config = &common.ServiceConfig{JsonConfig: string(fData)}
		serviceZkData.Config = &common.ServiceConfig{JsonConfig: string(fData)}
	}
	//获取route数据
	routeCallBack := NewServiceChangedCallback(serviceItem, SERVICE_ROUTE_CHANGED, f, handler)
	routeData, err := f.storageMgr.GetDataWithWatch(routePath, &routeCallBack)
	if err != nil {
		log.Log.Infof("query route data from path: %v %v", routePath, err)
		if strings.Compare(common.ZK_NODE_DOSE_NOT_EXIST, err.Error()) == 0 {
			log.Log.Infof("create node: %s", routePath)
			f.storageMgr.SetPath(routePath)
		}
		serviceZkData.Route = &common.ServiceRoute{RouteItem: []*common.RouteItem{}}

	} else if routeData != nil && len(routeData) == 0 {
		log.Log.Infof("query route data from path: %v empty", routePath)
		serviceZkData.Route = &common.ServiceRoute{RouteItem: []*common.RouteItem{}}
	} else {
		_, fData, err := common.DecodeValue(routeData)
		if err != nil {
			log.Log.Infof("parse err %s", err)
			serviceZkData.Route = &common.ServiceRoute{RouteItem: []*common.RouteItem{}}
		} else {
			serviceZkData.Route = route.ParseRouteData(fData)
		}
		service.ProviderList = route.FilterServiceByRouteData(serviceZkData.Route, f.config.MeteData.Address, service.ProviderList)

	}

	return service, nil
}
