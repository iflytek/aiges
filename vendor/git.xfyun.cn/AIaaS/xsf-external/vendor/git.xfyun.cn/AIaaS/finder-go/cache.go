package finder

import (
	"encoding/json"
	"fmt"

	common "git.xfyun.cn/AIaaS/finder-go/common"
	errors "git.xfyun.cn/AIaaS/finder-go/errors"
	"git.xfyun.cn/AIaaS/finder-go/utils/fileutil"
	"git.xfyun.cn/AIaaS/finder-go/log"
)

func CacheStorageInfo(cachePath string, zkInfo *common.StorageInfo) error {
	cachePath = fmt.Sprintf("%s/storage_%s.findercache", cachePath, "info")
	data, err := json.Marshal(zkInfo)
	if err != nil {
		log.Log.Errorf("CacheStorageInfo err: %s",err)
		return err
	}
	err = fileutil.WriteFile(cachePath, data)
	if err != nil {
		log.Log.Errorf("CacheStorageInfo err: %s",err)

		return err
	}

	return nil
}

func GetStorageInfoFromCache(cachePath string) (*common.StorageInfo, error) {
	cachePath = fmt.Sprintf("%s/storage_%s.findercache", cachePath, "info")

	data, err := fileutil.ReadFile(cachePath)
	if err != nil {
		log.Log.Errorf("GetStorageInfoFromCache err: %s",err)

		return nil, err
	}
	zkInfo := &common.StorageInfo{}
	err = json.Unmarshal(data, zkInfo)
	if err != nil {
		log.Log.Errorf("GetStorageInfoFromCache err: %s",err)
		return nil, err
	}

	return zkInfo, nil
}

func CacheConfig(cachePath string, config *common.Config) error {
	cachePath = fmt.Sprintf("%s/config_%s.findercache", cachePath, config.Name)
	var err error
	err = fileutil.WriteFile(cachePath, config.File)
	if err != nil {
		log.Log.Errorf("CacheConfig err: %s",err)
		return err
	}
	return nil
}

func GetConfigFromCache(cachePath string, name string) (*common.Config, error) {
	cachePath = fmt.Sprintf("%s/config_%s.findercache", cachePath, name)

	exist, err := fileutil.ExistPath(cachePath)
	if err != nil {
		log.Log.Errorf("GetConfigFromCache err: %s",err)
	}
	if !exist {
		err = errors.NewFinderError(errors.ConfigMissCacheFile)
		log.Log.Errorf("GetConfigFromCache err  %s:",err)
		return nil, err
	}
	data, err := fileutil.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}
	return &common.Config{Name: name, File: data}, nil
}

func CacheService(cachePath string, service *common.Service) error {
	if service==nil {
		return nil
	}
	cachePath = fmt.Sprintf("%s/service_%s_%s.findercache", cachePath, service.ServiceName, service.ApiVersion)
	data, err := json.Marshal(service)
	if err != nil {
		log.Log.Errorf("CacheService err: %s",err)
		return err
	}
	err = fileutil.WriteFile(cachePath, data)
	if err != nil {
		log.Log.Errorf("CacheService err: %s",err)
		return err
	}

	return nil
}

func GetServiceFromCache(cachePath string, item common.ServiceSubscribeItem) (*common.Service, error) {
	cachePath = fmt.Sprintf("%s/service_%s_%s.findercache", cachePath, item.ServiceName, item.ApiVersion)
	data, err := fileutil.ReadFile(cachePath)
	if err != nil {
		log.Log.Errorf("GetServiceFromCache err: %s",err)
		return nil, err
	}
	service := &common.Service{}
	err = json.Unmarshal(data, service)
	if err != nil {
		log.Log.Errorf("GetServiceFromCache err: %s",err)
		return nil, err
	}

	return service, nil
}
