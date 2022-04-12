package route

import (
	"encoding/json"
	common "github.com/xfyun/finder-go/common"
	"strings"
)

func newEmptyServiceRoute() *common.ServiceRoute {
	return &common.ServiceRoute{RouteItem: []*common.RouteItem{}}
}
func ParseRouteData(data []byte) *common.ServiceRoute {
	var serviceRoute common.ServiceRoute
	err := json.Unmarshal([]byte(`{"RouteItem":`+string(data)+"}"), &serviceRoute)
	if err != nil {
		return newEmptyServiceRoute()
	}
	return &serviceRoute
}

func FilterServiceByRouteData(serviceRoute *common.ServiceRoute, consumerAddr string, providerList []*common.ServiceInstance) []*common.ServiceInstance {
	var removeProviderList = make([]string, 0)
	var resultProvider = make([]*common.ServiceInstance, 0)
	if serviceRoute == nil {
		return providerList
	}
	for _, item := range serviceRoute.RouteItem {
		routeConsumer := item.Consumer
		routeProvider := item.Provider
		only := item.Only
		//如果在消费者组中，则直接返回对应的数据
		for _, consumer := range routeConsumer {
			if strings.Compare(consumerAddr, consumer) == 0 {
				for _, value := range routeProvider {
					if provide, ok := containProvider(providerList, value); ok {
						resultProvider = append(resultProvider, provide)
					}
				}
				return resultProvider
			}
		}
		//不在当前轮次的消费者中
		if strings.Compare("Y", only) == 0 {
			//如果这些提供者需要只能当前消费者看到
			removeProviderList = append(removeProviderList, routeProvider...)
		}
	}

	if len(removeProviderList) == 0 {
		return providerList
	}
	for _, value := range removeProviderList {
		providerList = deleteProvider(providerList, value)
	}
	return providerList
}
func containProvider(providerList []*common.ServiceInstance, provider string) (*common.ServiceInstance, bool) {
	for _, value := range providerList {
		if strings.Compare(value.Addr, provider) == 0 {
			return value, true
		}
	}
	return nil, false
}

func deleteProvider(providerList []*common.ServiceInstance, provider string) []*common.ServiceInstance {
	for index, value := range providerList {
		if strings.Compare(value.Addr, provider) == 0 {
			providerList = append(providerList[:index], providerList[index+1:]...)
			break
		}
	}
	return providerList
}
