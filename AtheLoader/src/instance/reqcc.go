package instance

import "sync"

// <appid, reqCC> 用于上报应用级的并发信息;
var appidCC = map[string]uint{}
var lock sync.Mutex

func appidCCDec(appid string) {
	lock.Lock()
	defer lock.Unlock()
	cc, exit := appidCC[appid]
	if exit {
		appidCC[appid] = cc - 1
		if appidCC[appid] == 0 {
			delete(appidCC, appid)
		}
	}
}

func appidCCInc(appid string) {
	lock.Lock()
	defer lock.Unlock()
	cc, _ := appidCC[appid] // if not exist, cc=0
	appidCC[appid] = cc + 1
}

// 避免appidCC读写并发异常,入参深拷贝;
// 减小锁粒度, 外部clear rpCC
func CCQuery(rpCC map[string]uint) {
	lock.Lock()
	defer lock.Unlock()
	for k, v := range appidCC {
		rpCC[k] = v
	}
}
