package instance

import (
	"sync"
)

// <serviceId,appId, reqCCR> 用于上报应用级的并发信息;
var appReqCCR = map[string]*map[string]uint{}
var lock sync.Mutex

const (
	sampleApp       = "appid"
	sampleServiceID = "service_id"
	sampleSub       = "sub"
)

// 运营管控数据采样入口
func sampleEnter(req *map[string]string) {
	if sub, ok := (*req)[sampleSub]; ok {
		if sub == "ase" {
			appCcrInc((*req)[sampleServiceID], (*req)[sampleApp])
		} else {
			appCcrInc(sub, (*req)[sampleApp])
		}
	} else {
		appCcrInc((*req)[sampleServiceID], (*req)[sampleApp])
	}
	// 新增其他运营管控数据
}

// 运营管控数据采样出口
func sampleExit(req *map[string]string) {
	if sub, ok := (*req)[sampleSub]; ok {
		if sub == "ase" {
			appCcrDec((*req)[sampleServiceID], (*req)[sampleApp])
		} else {
			appCcrDec(sub, (*req)[sampleApp])
		}
	} else {
		appCcrDec((*req)[sampleServiceID], (*req)[sampleApp])
	}
}

// 应用请求并发数据统计
func appCcrDec(svcId, appid string) {
	lock.Lock()
	defer lock.Unlock()

	appidMap, ok := appReqCCR[svcId]
	if ok {
		cc, exit := (*appidMap)[appid]
		if exit {
			(*appidMap)[appid] = cc - 1
			if (*appidMap)[appid] == 0 {
				delete(*appidMap, appid)
				if len(*appidMap) == 0 {
					delete(appReqCCR, svcId)
				}
			}
		}
	}
}

func appCcrInc(svcId, appid string) {
	lock.Lock()
	defer lock.Unlock()
	appidMap, exist := appReqCCR[svcId]
	if !exist {
		tempMap := new(map[string]uint)
		*tempMap = make(map[string]uint)
		(*tempMap)[appid] = 1
		appReqCCR[svcId] = tempMap
		return
	}
	cc, _ := (*appidMap)[appid] // if not exist, cc=0
	(*appidMap)[appid] = cc + 1
}

// 避免appCC读写并发异常,入参深拷贝;
// 减小锁粒度, 外部clear rpCC
func CCQuery(rpCC map[string]map[string]uint) (rlt map[string]map[string]uint) {
	lock.Lock()
	defer lock.Unlock()
	for k, v := range appReqCCR {
		tempMap := make(map[string]uint)
		for a, c := range *v {
			tempMap[a] = c
		}
		rpCC[k] = tempMap
	}
	rlt = rpCC
	return rlt
}
