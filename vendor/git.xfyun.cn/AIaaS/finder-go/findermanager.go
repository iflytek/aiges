package finder

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	common "git.xfyun.cn/AIaaS/finder-go/common"
	companion "git.xfyun.cn/AIaaS/finder-go/companion"
	errors "git.xfyun.cn/AIaaS/finder-go/errors"
	log "git.xfyun.cn/AIaaS/finder-go/log"
	"git.xfyun.cn/AIaaS/finder-go/storage"
	"git.xfyun.cn/AIaaS/finder-go/utils/arrayutil"
	"git.xfyun.cn/AIaaS/finder-go/utils/fileutil"
	"git.xfyun.cn/AIaaS/finder-go/utils/netutil"
	"git.xfyun.cn/AIaaS/finder-go/utils/stringutil"
	"sync"
)

var (
	hc *http.Client
)

const VERSION = "2.0.18"

type zkAddrChangeCallback struct {
	path string
	fm   *FinderManager
}

func (callback *zkAddrChangeCallback) ChildDeleteCallBack(path string) {

}
func (callback *zkAddrChangeCallback) ChildrenChangedCallback(path string, node string, children []string) {

}
func recoverFunc() {
	if err := recover(); err != nil {
		log.Log.Debugf("recover  %v", err)
	}
}
func (callback *zkAddrChangeCallback) Process(path string, node string) {
	defer recoverFunc()
	callback.fm.ServiceFinder.mutex.Lock()
	defer callback.fm.ServiceFinder.mutex.Unlock()
	log.Log.Debugf("handler zk_node_path  change  %v ", path)
	fm := callback.fm
	var tempPath sync.Map //恢复临时路径
	if fm.storageMgr != nil {
		tempPath = fm.storageMgr.GetTempPaths()
	}
	storageMgr, storageCfg, err := initStorageMgr(fm.config)
	if err != nil {
		log.Log.Errorf("init storage err , now retry :  %v", err)
		go watchStorageInfo(fm)
	} else {
		go watchZkInfo(fm)
		fm.storageMgr.Destroy()
		fm.storageMgr = storageMgr
		fm.ConfigFinder.storageMgr = storageMgr
		fm.ConfigFinder.rootPath = storageCfg.ConfigRootPath
		fm.ConfigFinder.config = fm.config
		fm.ServiceFinder.storageMgr = storageMgr
		fm.ServiceFinder.rootPath = storageCfg.ServiceRootPath
		fm.storageMgr.SetTempPaths(tempPath)
		fm.storageMgr.RecoverTempPaths() //恢复临时路径
		if len(fm.ServiceFinder.subscribedService) != 0 {
			//TODO ReGetServiceInfo会把之前的缓存数据给清了
			ReGetServiceInfo(fm)
		}
		if len(fm.ConfigFinder.fileSubscribe) != 0 {
			ReGetConfigInfo(fm)
		}
	}
}

func (callback *zkAddrChangeCallback) DataChangedCallback(path string, node string, data []byte) {

}

func init() {
	hc = &http.Client{
		Transport: &http.Transport{
			Dial: func(nw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(1 * time.Second)
				c, err := net.DialTimeout(nw, addr, time.Second*1)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

}

// FinderManager for controll all
type FinderManager struct {
	config        *common.BootConfig
	storageMgr    storage.StorageManager
	ConfigFinder  *ConfigFinder
	ServiceFinder *ServiceFinder
}

func checkCachePath(path string) (string, error) {
	if stringutil.IsNullOrEmpty(path) {
		p, err := os.Getwd()
		if err == nil {
			p += (fileutil.GetSystemSeparator() + common.DefaultCacheDir)
			path = p
		} else {
			return path, err
		}
	}

	return path, nil
}

func createCacheDir(path string) error {
	exist, err := fileutil.ExistPath(path)
	if err == nil && !exist {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func checkConfig(c *common.BootConfig) {
	if c.ExpireTimeout <= 0 {
		c.ExpireTimeout = 3 * time.Second
	}
}

func getStorageInfo(config *common.BootConfig) (*common.StorageInfo, error) {
	url := config.CompanionUrl + fmt.Sprintf("/finder/query_zk_info?project=%s&group=%s&service=%s&version=%s", config.MeteData.Project, config.MeteData.Group, config.MeteData.Service, config.MeteData.Version)
	info, err := companion.GetStorageInfo(hc, url)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func checkAddr(n []string, o []string) bool {
	vchanged := false
	for _, nv := range o {
		if !arrayutil.Contains(nv, o) {
			vchanged = true
		}
	}

	return vchanged
}

func onZkInfoChanged(smr storage.StorageManager) {
	// todo.
}

func getStorageConfig(config *common.BootConfig) (*storage.StorageConfig, error) {
	checkConfig(config)
	info, err := getStorageInfo(config)
	if err != nil {
		return nil, err
	}
	storageConfig := &storage.StorageConfig{
		Name:   "zookeeper",
		Params: make(map[string]string),
	}
	storageConfig.Params["servers"] = strings.Join(info.Addr, ",")
	storageConfig.Params["session_timeout"] = strconv.FormatInt(int64(config.ExpireTimeout/time.Millisecond), 10)
	storageConfig.Params["zk_node_path"] = info.ZkNodePath
	storageConfig.ConfigRootPath = info.ConfigRootPath
	storageConfig.ServiceRootPath = info.ServiceRootPath
	return storageConfig, nil
}

func initStorageMgr(config *common.BootConfig) (storage.StorageManager, *storage.StorageConfig, error) {
	storageConfig, err := getStorageConfig(config)
	if err != nil {
		log.Log.Errorf("get storage config err : %v", err)
		return nil, nil, err
	}
	log.Log.Debugf("storage info ： %v ", storageConfig.Params)
	storageMgr, err := storage.NewManager(storageConfig)
	if err != nil {
		log.Log.Errorf("[ initStorageMgr ] NewManager: %s", err)
		return nil, storageConfig, err
	}
	err = storageMgr.Init()
	if err != nil {
		log.Log.Errorf(" storage init err %v", err)
		return nil, storageConfig, err
	}

	return storageMgr, storageConfig, nil
}

// NewFinder for creating an instance
func newFinder(config common.BootConfig) (*FinderManager, error) {
	log.Log = log.NewDefaultLogger()
	if stringutil.IsNullOrEmpty(config.CompanionUrl) {
		err := errors.NewFinderError(errors.MissCompanionUrl)
		return nil, err
	}

	if stringutil.IsNullOrEmpty(config.MeteData.Address) {
		localIP, err := netutil.GetLocalIP(config.CompanionUrl)
		if err != nil {
			log.Log.Errorf("%s", err)
			return nil, err
		}
		config.MeteData.Address = localIP
	}

	// 检查缓存路径，如果传入cachePath是空，则使用默认路径
	p, err := checkCachePath(config.CachePath)
	if err != nil {
		return nil, err
	}

	// 创建缓存目录
	err = createCacheDir(p)
	if err != nil {
		return nil, err
	}
	config.CachePath = p
	// 初始化finder
	fm := new(FinderManager)
	fm.config = &config
	// 初始化zk
	var storageCfg *storage.StorageConfig
	fm.storageMgr, storageCfg, err = initStorageMgr(fm.config)
	if err != nil {
		log.Log.Infof("初始化zk信息出错，开启新的goroutine 去不断尝试")
		fm.ConfigFinder = NewConfigFinder("", fm.config, nil)
		fm.ServiceFinder = NewServiceFinder("", fm.config, nil)
		//return nil, err
	} else {
		fm.ConfigFinder = NewConfigFinder(storageCfg.ConfigRootPath, fm.config, fm.storageMgr)
		fm.ServiceFinder = NewServiceFinder(storageCfg.ServiceRootPath, fm.config, fm.storageMgr)
	}

	return fm, nil
}

func NewFinderWithLogger(config common.BootConfig, logger log.Logger) (*FinderManager, error) {
	if func(obj interface{}) bool {
		if obj == nil {
			return true
		}
		type eface struct {
			rtype unsafe.Pointer
			data  unsafe.Pointer
		}
		return (*eface)(unsafe.Pointer(&obj)).data == nil
	}(logger) {
		log.Log = log.NewDefaultLogger()
	} else {
		log.Log = logger
	}
	log.Log.Infof("current version : %v " + VERSION)
	if stringutil.IsNullOrEmpty(config.CompanionUrl) {
		err := errors.NewFinderError(errors.MissCompanionUrl)
		return nil, err
	}
	if !config.MeteData.Check() {
		err := errors.NewFinderError(errors.InvalidParam)
		return nil, err
	}
	if stringutil.IsNullOrEmpty(config.MeteData.Address) {
		localIP, err := netutil.GetLocalIP(config.CompanionUrl)
		if err != nil {
			log.Log.Errorf("%s", err)
			return nil, err
		}
		config.MeteData.Address = localIP
	}

	// 检查缓存路径，如果传入cachePath是空，则使用默认路径
	p, err := checkCachePath(config.CachePath)
	if err != nil {
		return nil, err
	}
	// 创建缓存目录
	err = createCacheDir(p)
	if err != nil {
		return nil, err
	}
	config.CachePath = p
	// 初始化finder
	fm := new(FinderManager)
	fm.config = &config
	// 初始化zk
	var storageCfg *storage.StorageConfig
	fm.storageMgr, storageCfg, err = initStorageMgr(fm.config)
	if err != nil {
		log.Log.Infof("init zk err retry")
		fm.ConfigFinder = NewConfigFinder("", fm.config, nil)
		fm.ServiceFinder = NewServiceFinder("", fm.config, nil)
		go manitorStorage(fm)
		go watchStorageInfo(fm)
		return fm, nil
	} else {
		fm.ConfigFinder = NewConfigFinder(storageCfg.ConfigRootPath, fm.config, fm.storageMgr)
		fm.ServiceFinder = NewServiceFinder(storageCfg.ServiceRootPath, fm.config, fm.storageMgr)
	}
	//创建一个goroutine来执行监听zk地址的数据
	go manitorStorage(fm)
	go watchZkInfo(fm)
	return fm, nil
}
func manitorStorage(fm *FinderManager) {
	log.Log.Infof("manitorStorage")
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if fm.storageMgr == nil {
				continue
			}
			storageConfig, err := getStorageConfig(fm.config)
			if err != nil {
				log.Log.Errorf("[ manitorStorage ] getStorageConfig: %s", err)
				continue
			}
			if storageConfig.Params["servers"] == fm.storageMgr.GetServerAddr() {
				continue
			}
			log.Log.Infof("zk addr is change %s,%s,%s", storageConfig.Params["servers"], " ----> ", fm.storageMgr.GetServerAddr())
			go watchStorageInfo(fm)
		}
	}
}
func watchZkInfo(fm *FinderManager) {

	zkNodePath, err := fm.storageMgr.GetZkNodePath()
	if err != nil {
		log.Log.Errorf("zk node path is err,%v",err)
	}
	fm.storageMgr.GetDataWithWatchV2(zkNodePath, &zkAddrChangeCallback{path: zkNodePath, fm: fm})
}

func watchStorageInfo(fm *FinderManager) {
	defer recoverFunc()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	var tempPathMap sync.Map
	if fm.storageMgr != nil {
		tempPathMap = fm.storageMgr.GetTempPaths()
	}
	var storageChange bool
	for {
		select {
		case <-ticker.C:
			fm.ServiceFinder.mutex.Lock()
			storageMgr, storageCfg, err := initStorageMgr(fm.config)
			if err != nil {
				storageChange = false
				log.Log.Infof("init zk err %v ", err)
				//if fm.storageMgr != nil {
				//	fm.storageMgr.Destroy()
				//}
			} else {
				storageChange = true
				if fm.storageMgr != nil {
					fm.storageMgr.Destroy()
				}
				fm.storageMgr = storageMgr
				fm.ConfigFinder.storageMgr = storageMgr
				fm.ConfigFinder.rootPath = storageCfg.ConfigRootPath
				fm.ConfigFinder.config = fm.config
				fm.ServiceFinder.storageMgr = storageMgr
				fm.ServiceFinder.rootPath = storageCfg.ServiceRootPath
			}
			fm.ServiceFinder.mutex.Unlock()
		}
		if storageChange {
			fm.ServiceFinder.mutex.Lock()
			go watchZkInfo(fm)
			fm.storageMgr.SetTempPaths(tempPathMap)
			fm.storageMgr.RecoverTempPaths()
			if len(fm.ServiceFinder.subscribedService) != 0 {
				log.Log.Debugf("retry query service info，%v", fm.ServiceFinder.subscribedService)
				//重新拉取所有订阅服务的信息
				ReGetServiceInfo(fm)
			}
			if len(fm.ConfigFinder.fileSubscribe) != 0 {
				log.Log.Debugf("retry query config info，%v", fm.ConfigFinder.fileSubscribe)
				ReGetConfigInfo(fm)
			}
			fm.ServiceFinder.mutex.Unlock()
			break
		}
	}
}
func ReGetConfigInfo(fm *FinderManager) {
	handler := fm.ConfigFinder.handler
	var fileSubscribe []string
	for _, file := range fm.ConfigFinder.fileSubscribe {
		fileSubscribe = append(fileSubscribe, file)
	}
	fm.ConfigFinder.fileSubscribe = []string{}
	fm.ConfigFinder.grayConfig.Range(func(key, value interface{}) bool {
		fm.ConfigFinder.grayConfig.Delete(key)
		return true
	})

	fm.ConfigFinder.usedConfig.Range(func(key, value interface{}) bool {
		fm.ConfigFinder.usedConfig.Delete(key)
		return true
	})
	fileMap, err := fm.ConfigFinder.UseAndSubscribeConfig(fileSubscribe, handler)
	if err != nil {
		log.Log.Errorf("query info err %v", err)
	}
	for _, fileData := range fileMap {
		var config = common.Config{Name: fileData.Name, File: fileData.File, ConfigMap: fileData.ConfigMap}
		if handler != nil {
			handler.OnConfigFileChanged(&config)

		}
	}
}

func ReGetServiceInfo(fm *FinderManager) {

	for key,value:=range fm.ServiceFinder.serviceZkData{
		log.Log.Infof("all zk data  %v ,%v",key,value.ProviderList)
	}

	for _, value := range fm.ServiceFinder.subscribedService {
		var item = common.ServiceSubscribeItem{ServiceName: value.ServiceName, ApiVersion: value.ApiVersion}
		servicePath := fmt.Sprintf("%s/%s/%s", fm.ServiceFinder.rootPath, item.ServiceName, item.ApiVersion)
		prevService:=fm.ServiceFinder.subscribedService[value.ServiceName+"_"+value.ApiVersion]
		service, err := fm.ServiceFinder.getServiceWithWatcher(servicePath, item, fm.ServiceFinder.handler)
		if err != nil {
			log.Log.Errorf("query info err %s", err)
		}
		if service == nil {
			log.Log.Infof("query service is empty ")
		}
		if prevService==nil{
			prevService, err = GetServiceFromCache(fm.config.CachePath, item)
		}
		ChangeEvent(prevService, service, fm.ServiceFinder.handler)
		for key,value:=range fm.ServiceFinder.serviceZkData{
			log.Log.Infof("all zk data event %v ,%v",key,value.ProviderList)
		}
		fm.ServiceFinder.subscribedService[value.ServiceName+"_"+value.ApiVersion] = service
		if service != nil {
			CacheService(fm.config.CachePath, service)
		}

	}
}

func ChangeEvent(prevService *common.Service, currService *common.Service, handler common.ServiceChangedHandler) {
	if handler == nil {
		return
	}
	if prevService == nil {
		handler.OnServiceConfigChanged(currService.ServiceName, currService.ApiVersion, &common.ServiceConfig{JsonConfig: currService.Config.JsonConfig})
		eventList := providerChangeEvent([]*common.ServiceInstance{}, currService.ProviderList)
		log.Log.Infof("prevService is nil : %v",eventList)
		if len(eventList) == 0 {
			return
		}

		handler.OnServiceInstanceChanged(currService.ServiceName, currService.ApiVersion, eventList)
		return
	}
	if currService == nil {
		handler.OnServiceConfigChanged(prevService.ServiceName, prevService.ApiVersion, &common.ServiceConfig{JsonConfig: prevService.Config.JsonConfig})
		eventList := providerChangeEvent(prevService.ProviderList, []*common.ServiceInstance{})
		if len(eventList) == 0 {
			return
		}
		handler.OnServiceInstanceChanged(prevService.ServiceName, prevService.ApiVersion, eventList)
		return
	}
	prevConfig := prevService.Config
	currConfig := currService.Config
	if prevConfig.JsonConfig != currConfig.JsonConfig {
		handler.OnServiceConfigChanged(currService.ServiceName, currService.ApiVersion, &common.ServiceConfig{JsonConfig: currConfig.JsonConfig})
	}
	eventList := providerChangeEvent(prevService.ProviderList, currService.ProviderList)
	if len(eventList) == 0 {
		return
	}
	handler.OnServiceInstanceChanged(currService.ServiceName, currService.ApiVersion, eventList)
	return
}

func providerChangeEvent(prevProviderList, currProviderList []*common.ServiceInstance) []*common.ServiceInstanceChangedEvent {
	var eventList []*common.ServiceInstanceChangedEvent
	if len(prevProviderList) == 0 && len(currProviderList) == 0 {
		return nil
	}
	if len(prevProviderList) == 0 {
		var changeList []*common.ServiceInstance
		for _, provider := range currProviderList {
			changeList = append(changeList, provider.Dumplication())
		}
		event := common.ServiceInstanceChangedEvent{EventType: common.INSTANCEADDED, ServerList: changeList}
		eventList = append(eventList, &event)
		return eventList
	}
	if len(currProviderList) == 0 {
		var changeList []*common.ServiceInstance
		for _, provider := range prevProviderList {
			changeList = append(changeList, provider.Dumplication())
		}
		event := common.ServiceInstanceChangedEvent{EventType: common.INSTANCEREMOVE, ServerList: changeList}
		eventList = append(eventList, &event)
		return eventList
	}
	var addServerList []*common.ServiceInstance
	//TODO 后续优化
	var providerMap = make(map[string]*common.ServiceInstance)
	for _, prevProvider := range prevProviderList {
		providerMap[prevProvider.Addr] = prevProvider
	}
	for _, currProvider := range currProviderList {
		if _, ok := providerMap[currProvider.Addr]; !ok {
			addServerList = append(addServerList, currProvider.Dumplication())
		} else {
			delete(providerMap, currProvider.Addr)
		}
	}
	var removeServerList []*common.ServiceInstance
	for _, provider := range providerMap {
		removeServerList = append(removeServerList, provider.Dumplication())
	}
	removeEvent := common.ServiceInstanceChangedEvent{EventType: common.INSTANCEREMOVE, ServerList: removeServerList}
	eventList = append(eventList, &removeEvent)
	addEvent := common.ServiceInstanceChangedEvent{EventType: common.INSTANCEADDED, ServerList: addServerList}
	eventList = append(eventList, &addEvent)
	return eventList

}

func DestroyFinder(finder *FinderManager) {
	finder.storageMgr.Destroy()
	// todo
}

func onCfgUpdateEvent(c common.Config) int {
	return errors.ConfigSuccess
}
