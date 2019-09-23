package serviceutil

import (
	"encoding/json"
	"git.xfyun.cn/AIaaS/finder-go/common"
	common "git.xfyun.cn/AIaaS/finder-go/common"
	"log"
)

func ParseServiceConfigData(data []byte) *common.ServiceInstanceConfig {
	var configJson map[string]interface{}
	serviceInstanceConfig := &common.ServiceInstanceConfig{}
	err := json.Unmarshal(data, &configJson)
	if err != nil {
		log.Println("【ParseServiceConfigData】出错 ", err,"   ",string(data)  )
		return nil
	}
	log.Println("[ ParseServiceConfigData ] configJson: ",configJson,"  data: ",string(data))

	if sdkConfig ,ok:= configJson["sdk"].(map[string]interface{});ok {
		if isValid, ok := sdkConfig["is_valid"].(bool); ok {
			serviceInstanceConfig.IsValid = isValid
		} else {
			serviceInstanceConfig.IsValid = true
		}
	}


	delete(configJson, "sdk")
	userData, _ := json.Marshal(configJson)
	serviceInstanceConfig.UserConfig = string(userData)

	return serviceInstanceConfig
}



//采用内存换时间的策略
func CompareServiceInstanceList(prevProviderList []*common.ServiceInstance, currentProviderList []*common.ServiceInstance) []*common.ServiceInstanceChangedEvent {

	var providerMap = make(map[string]*common.ServiceInstance)
	var countMap = make(map[string]int8)
	eventList := []*common.ServiceInstanceChangedEvent{}
	addServiceInstance := common.ServiceInstanceChangedEvent{EventType: finder.INSTANCEADDED, ServerList: []*common.ServiceInstance{}}
	removeServiceInstance := common.ServiceInstanceChangedEvent{EventType: finder.INSTANCEREMOVE, ServerList: []*common.ServiceInstance{}}
	for _, provider := range prevProviderList {
		providerMap[provider.Addr] = provider
		countMap[provider.Addr] += 1
	}

	//看看是否有新增的
	var incrFlag = false
	for _, provider := range currentProviderList {
		if _, ok := providerMap[provider.Addr]; !ok {
			//找到新增的 就是在原来的不存在的
			var instance common.ServiceInstance
			instance.Addr = provider.Addr
			instance.Config = &common.ServiceInstanceConfig{IsValid: provider.Config.IsValid, UserConfig: provider.Config.UserConfig}
			addServiceInstance.ServerList = append(addServiceInstance.ServerList, &instance)
			incrFlag = true
		} else {
			countMap[provider.Addr] += 1
		}
	}
	log.Println(" countMap:",countMap)
	if incrFlag {
		eventList = append(eventList, &addServiceInstance)
	}
	//看是否有服务提供者减少了
	var decrFlag = false
	for key, value := range countMap {
		if value == 1 {
			provider := providerMap[key]
			var instance common.ServiceInstance
			instance.Addr = provider.Addr
			instance.Config = &common.ServiceInstanceConfig{IsValid: provider.Config.IsValid, UserConfig: provider.Config.UserConfig}
			removeServiceInstance.ServerList = append(removeServiceInstance.ServerList, &instance)
			decrFlag = true
		}
	}
	if decrFlag {
		eventList = append(eventList, &removeServiceInstance)
	}
	return eventList
}
