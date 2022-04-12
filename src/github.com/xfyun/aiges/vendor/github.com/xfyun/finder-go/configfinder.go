package finder

import (
	"fmt"
	"strings"
	"sync"

	common "github.com/xfyun/finder-go/common"
	errors "github.com/xfyun/finder-go/errors"
	"github.com/xfyun/finder-go/log"
	"github.com/xfyun/finder-go/storage"
	"github.com/xfyun/finder-go/utils/fileutil"
)

var (
	configEventPrefix = "config_"
)

type ConfigFinder struct {
	locker           sync.Mutex
	rootPath         string
	currentWatchPath string
	config           *common.BootConfig
	storageMgr       storage.StorageManager
	handler          common.ConfigChangedHandler
	usedConfig       sync.Map
	fileSubscribe    []string
	grayConfig       sync.Map
}

func NewConfigFinder(root string, bc *common.BootConfig, sm storage.StorageManager) *ConfigFinder {

	finder := &ConfigFinder{
		locker:     sync.Mutex{},
		rootPath:   root,
		config:     bc,
		storageMgr: sm,
		usedConfig: sync.Map{},
	}

	return finder
}

func (f *ConfigFinder) UseConfig(name []string) (map[string]*common.Config, error) {
	cfg, err := f.useConfig(name)
	m := f.config.MeteData
	if err != nil {
		return nil, fmt.Errorf("subscribe config error companion:%s,  path:/%s/%s/%s/%s files:%v err:%w", f.config.CompanionUrl, m.Project, m.Group, m.Service, m.Version, name, err)
	}
	return cfg, nil
}

func (f *ConfigFinder) UseAndSubscribeConfig(name []string, handler common.ConfigChangedHandler) (map[string]*common.Config, error) {
	cfg, err := f.useAndSubscribeConfig(name, handler)
	m := f.config.MeteData
	if err != nil {
		return nil, fmt.Errorf("subscribe config error companion:%s, path:/%s/%s/%s/%s files:%v err:%w", f.config.CompanionUrl, m.Project, m.Group, m.Service, m.Version, name, err)
	}
	return cfg, nil
}

func (f *ConfigFinder) UseAndSubscribeWithPrefix(prefix string, handler common.ConfigChangedHandler) (map[string]*common.Config, error) {
	cfg, err := f.useAndSubscribeWithPrefix(prefix, handler)

	m := f.config.MeteData
	if err != nil {
		return nil, fmt.Errorf("subscribe config error companion:%s, path: /%s/%s/%s/%s files:%v err:%w", f.config.CompanionUrl, m.Project, m.Group, m.Service, m.Version, prefix, err)
	}
	return cfg, nil
}

// UseConfig for 订阅相关配置文件
func (f *ConfigFinder) useConfig(name []string) (map[string]*common.Config, error) {
	if len(name) == 0 {
		err := errors.NewFinderError(errors.ConfigMissName)
		return nil, err
	}

	f.locker.Lock()
	defer f.locker.Unlock()
	if f.storageMgr == nil {
		log.Log.Infof("zk init err")
		return nil, errors.NewFinderError(errors.ZkGetInfoError)
	}
	err := GetGrayConfigData(f, f.rootPath, nil)
	if err != nil {
		log.Log.Infof("query gray info err: %s", err)
		return nil, err
	}
	configFiles := make(map[string]*common.Config)
	for _, n := range name {
		if c, ok := f.usedConfig.Load(n); !ok {
			//先获取gray的数据，用于判断订阅的配置是否在灰度组中
			basePath := f.rootPath
			if groupId, ok := f.grayConfig.Load(f.config.MeteData.Address); ok {
				basePath += "/gray/" + groupId.(string)
			}
			//真正获取数据
			data, err := f.storageMgr.GetData(basePath + "/" + n)
			if err != nil {
				//出错 从配置文件中取
				onUseConfigErrorWithCache(configFiles, n, f.config.CachePath, err)
			} else {
				_, fData, err := common.DecodeValue(data)
				if err != nil {
					//出错 从配置文件中获取
					onUseConfigErrorWithCache(configFiles, n, f.config.CachePath, err)
				} else {
					var config *common.Config
					if fileutil.IsTomlFile(n) {
						tomlConfig := fileutil.ParseTomlFile(fData)
						config = &common.Config{Name: n, File: fData, ConfigMap: tomlConfig}
					} else {
						config = &common.Config{Name: n, File: fData}
					}
					configFiles[n] = config
					//存到缓存
					err = CacheConfig(f.config.CachePath, config)
					if err != nil {
						log.Log.Errorf("CacheConfig: %s", err)
					}
				}
			}
		} else {
			// todo
			if config, ok := c.(common.Config); ok {
				configFiles[n] = &config
			} else {
				// get config from cache
				configFiles[n] = getCachedConfig(n, f.config.CachePath)
			}
		}
	}

	return configFiles, nil
}

// UseAndSubscribeConfig for
//新增监控灰度组的Watch
func (f *ConfigFinder) useAndSubscribeConfig(name []string, handler common.ConfigChangedHandler) (map[string]*common.Config, error) {
	if len(name) == 0 {
		err := errors.NewFinderError(errors.ConfigMissName)
		return nil, err
	}
	f.locker.Lock()
	defer f.locker.Unlock()
	f.handler = handler
	configFiles := make(map[string]*common.Config)
	if f.storageMgr == nil {
		if f.config.CacheConfig {
			log.Log.Infof("init zk err,use cache")
			for _, n := range name {
				f.fileSubscribe = append(f.fileSubscribe, n)
				configFiles[n] = getCachedConfig(n, f.config.CachePath)
			}
			return configFiles, nil
		} else {
			log.Log.Infof("init zk err ,not use cache ,exit")
			return nil, fmt.Errorf("finder is nil") // 此处应该返回error
		}
	}
	log.Log.Debugf("subscribe file ：%v", name)
	//先查看灰度组的设置

	callback := NewConfigChangedCallback(f.config.MeteData.Address, CONFIG_CHANGED, f.rootPath, handler, f.config, f.storageMgr, f)

	err := GetGrayConfigData(f, f.rootPath, &callback)
	if err != nil {
		log.Log.Infof("get gray config err %v", err)
		return nil, err
	}

	if groupId, ok := f.grayConfig.Load(f.config.MeteData.Address); ok {
		if ok := f.checkFileExist(f.rootPath+"/gray/"+groupId.(string), name); !ok {
			log.Log.Infof("file not exist,path: %v", f.rootPath+"/gray/"+groupId.(string))
			return nil, errors.NewFinderError(errors.ConfigFileNotExist)
		}
	} else {
		if ok := f.checkFileExist(f.rootPath, name); !ok {
			log.Log.Infof("file not exist,path: %v", f.rootPath)
			return nil, errors.NewFinderError(errors.ConfigFileNotExist)
		}
	}

	consumerPath := f.rootPath + "/consumer"
	if groupId, ok := f.grayConfig.Load(f.config.MeteData.Address); ok {
		//如果在灰度组。则进行注册到灰度组中
		consumerPath += "/gray/" + groupId.(string) + "/" + f.config.MeteData.Address
		f.storageMgr.SetTempPath(consumerPath)
	} else {
		consumerPath += "/normal/" + f.config.MeteData.Address
		f.storageMgr.SetTempPath(consumerPath)
	}

	path := ""
	for _, n := range name {
		f.fileSubscribe = append(f.fileSubscribe, n)

		basePath := f.rootPath
		if groupId, ok := f.grayConfig.Load(f.config.MeteData.Address); ok {
			basePath += "/gray/" + groupId.(string)
		}
		callback := NewConfigChangedCallback(n, CONFIG_CHANGED, f.rootPath, handler, f.config, f.storageMgr, f)

		//根据获取的灰度组设置的结果，来到特定的节点获取配置文件数据
		path = basePath + "/" + n
		data, err := f.storageMgr.GetDataWithWatchV2(path, &callback)
		if err != nil {
			if strings.Compare(err.Error(), common.ZK_NODE_DOSE_NOT_EXIST) == 0 {
				log.Log.Infof("config file not exist。filename: %v", name)
				return nil, errors.NewFinderError(errors.ConfigFileNotExist)
			}
			onUseConfigErrorWithCache(configFiles, n, f.config.CachePath, err)

		} else {
			_, fData, err := common.DecodeValue(data)
			if err != nil {
				onUseConfigErrorWithCache(configFiles, n, f.config.CachePath, err)
			} else {
				//
				confMap := make(map[string]interface{})
				if fileutil.IsTomlFile(n) {
					confMap = fileutil.ParseTomlFile(fData)
				}
				config := &common.Config{Name: n, File: fData, ConfigMap: confMap}
				configFiles[n] = config
				f.usedConfig.Store(n, config)
				//放到文件中
				err = CacheConfig(f.config.CachePath, config)
				if err != nil {
					log.Log.Errorf("CacheConfig: %s", err)
				}
			}
		}

	}
	return configFiles, nil
}

func (f *ConfigFinder) useAndSubscribeWithPrefix(prefix string, handler common.ConfigChangedHandler) (map[string]*common.Config, error) {
	f.locker.Lock()
	defer f.locker.Unlock()
	f.handler = handler
	configFiles := make(map[string]*common.Config)
	if f.storageMgr == nil {
		if f.config.CacheConfig {
			log.Log.Infof("init zk err,use cache")
			configFiles = getAllCachedConfig(f.config.CachePath, prefix)
			if configFiles != nil {
				for k, _ := range configFiles {
					f.fileSubscribe = append(f.fileSubscribe, k)
					// TODO: 去重？
				}
			}

			return configFiles, nil
		}

		log.Log.Infof("init zk err ,not use cache ,exit")
		return nil, nil
	}
	log.Log.Debugf("call UseAndSubscribeAll")

	if ok := f.checkDirExist(f.rootPath); !ok {
		log.Log.Infof("file not exist,path: %v", f.rootPath)
		return nil, errors.NewFinderError(errors.ConfigFileNotExist)
	}

	consumerPath := f.rootPath + "/consumer"
	consumerPath += "/normal/" + f.config.MeteData.Address
	f.storageMgr.SetTempPath(consumerPath)

	path := ""
	// watch dir
	dirCallback := NewConfigChangedCallback(prefix, CONFIG_DIR_CHANGED, f.rootPath, handler, f.config, f.storageMgr, f)
	names, err := f.storageMgr.GetChildrenWithWatch(f.rootPath, &dirCallback)
	if err != nil {
		if strings.Compare(err.Error(), common.ZK_NODE_DOSE_NOT_EXIST) == 0 {
			log.Log.Infof("config dir not exist,path: %v", f.rootPath)
			return nil, errors.NewFinderError(errors.ConfigDirNotExist)
		}

		return nil, err
	}

	for _, n := range names {
		if !strings.HasPrefix(n, prefix) {
			continue
		}
		f.fileSubscribe = append(f.fileSubscribe, n)
		basePath := f.rootPath
		path = basePath + "/" + n
		callback := NewConfigChangedCallback(n, CONFIG_CHANGED, path, handler, f.config, f.storageMgr, f)
		data, err := f.storageMgr.GetDataWithWatchV2(path, &callback)
		if err != nil {
			if strings.Compare(err.Error(), common.ZK_NODE_DOSE_NOT_EXIST) == 0 {
				log.Log.Infof("config file not exist。filename: %v", n)
				return nil, errors.NewFinderError(errors.ConfigFileNotExist)
			}
			onUseConfigErrorWithCache(configFiles, n, f.config.CachePath, err)
		} else {
			_, fData, err := common.DecodeValue(data)
			if err != nil {
				onUseConfigErrorWithCache(configFiles, n, f.config.CachePath, err)
			} else {
				//
				confMap := make(map[string]interface{})
				if fileutil.IsTomlFile(n) {
					confMap = fileutil.ParseTomlFile(fData)
				}
				config := &common.Config{Name: n, File: fData, ConfigMap: confMap}
				configFiles[n] = config
				f.usedConfig.Store(n, config)
				//放到文件中
				err = CacheConfig(f.config.CachePath, config)
				if err != nil {
					log.Log.Errorf("CacheConfig: %s", err)
				}
			}
		}

	}

	return configFiles, nil
}

func (f *ConfigFinder) checkFileExist(basePath string, names []string) bool {
	//TODO 判断文件是否存在，不存在则直接报错，
	log.Log.Debugf("basePath %s", basePath)
	files, err := f.storageMgr.GetChildren(basePath)
	if err != nil {
		log.Log.Errorf("query config file err : %v", err)
		return false
	}
	if len(names) > len(files) {
		log.Log.Infof("current file is %v, subscribe file is %v", files, names)
		return false
	}
	for _, subFileName := range names {
		var isExist = false
		for _, existFile := range files {
			if existFile == subFileName {
				isExist = true
			}
		}
		if !isExist {
			log.Log.Infof("current file is %v, subscribe file is %v", files, subFileName)
			return false
		}
	}
	return true

}

func (f *ConfigFinder) checkDirExist(basePath string) bool {
	//TODO 判断文件是否存在，不存在则直接报错，
	log.Log.Debugf("basePath %s", basePath)
	dirInfo, err := f.storageMgr.GetData(basePath)
	if err != nil {
		log.Log.Errorf("query config dir err : %v", err)
		return false
	}

	if dirInfo == nil {
		return false
	}

	return true
}

func (f *ConfigFinder) UnSubscribeConfig(name string) error {
	var err error
	if len(name) == 0 {
		err = errors.NewFinderError(errors.ConfigMissName)
		return err
	}
	for index, value := range f.fileSubscribe {
		if strings.Compare(name, value) == 0 {
			f.fileSubscribe = append(f.fileSubscribe[:index], f.fileSubscribe[index+1:]...)
		}
	}
	if len(f.fileSubscribe) == 0 {
		f.removeConfigConsumer()
	}

	return nil
}

func (f *ConfigFinder) removeConfigConsumer() {
	//如果订阅文件的个数为0，则取消注册者
	consumerPath := f.rootPath + "/consumer"
	if groupId, ok := f.grayConfig.Load(f.config.MeteData.Address); ok {
		//如果在灰度组。则进行注册到灰度组中
		consumerPath += "/gray/" + groupId.(string) + "/" + f.config.MeteData.Address
		f.storageMgr.Remove(consumerPath)
	} else {
		consumerPath += "/normal/" + f.config.MeteData.Address
		f.storageMgr.Remove(consumerPath)
	}
}
func (f *ConfigFinder) BatchUnSubscribeConfig(names []string) error {
	if len(names) == 0 {
		err := errors.NewFinderError(errors.ConfigMissName)
		return err
	}
	for _, name := range names {
		for index, value := range f.fileSubscribe {
			if strings.Compare(name, value) == 0 {
				f.fileSubscribe = append(f.fileSubscribe[:index], f.fileSubscribe[index+1:]...)
			}
		}
	}
	if len(f.fileSubscribe) == 0 {
		f.removeConfigConsumer()
	}
	return nil
}

// onUseConfigError with cache
func onUseConfigErrorWithCache(configFiles map[string]*common.Config, name string, cachePath string, err error) {
	log.Log.Errorf("onUseConfigError: %v, name: %v", err, name)
	configFiles[name] = getCachedConfig(name, cachePath)
}

func getCachedConfig(name string, cachePath string) *common.Config {
	config, err := GetConfigFromCache(cachePath, name)
	if err != nil {
		log.Log.Errorf("GetConfigFromCache: %s", err)
		return nil
	}

	return config
}

func getAllCachedConfig(cachePath string, prefix string) map[string]*common.Config {
	config, err := GetAllConfigFromCache(cachePath, prefix)
	if err != nil {
		log.Log.Errorf("GetConfigFromCache: %s", err)
		return nil
	}

	return config
}
