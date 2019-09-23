package xsf

import (
	"fmt"
	"sync"

	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

var finderadapter finderAdapter

const defaultApiVersion = "1.0.0"

type finderAdapter struct {
	versionMap map[string]string
	finderMap  map[string]*utils.FindManger
	finderRWMU sync.RWMutex
}

func init() {
	finderadapter.finderMap = make(map[string]*utils.FindManger)
	finderadapter.versionMap = make(map[string]string)
	addKillerCheck(killerFirstPriority, "finder", &finderadapter)
}
func (f *finderAdapter) AddRegister(version, addr string, findManager *utils.FindManger) {
	f.finderRWMU.Lock()
	defer f.finderRWMU.Unlock()
	f.finderMap[addr] = findManager
	f.versionMap[addr] = version
}

//bool返回值表示是否有finder配置
func (f *finderAdapter) Register(addr string, version string) (bool, error) {
	f.finderRWMU.RLock()
	defer f.finderRWMU.RUnlock()
	if finder, finderOk := f.finderMap[addr]; finderOk {
		loggerStd.Printf("about to call finder.RegisterSrvWithAddr -> addr:%v,version:%s\n", addr, version)
		registerErr := finder.RegisterSrvWithAddr(addr, version)
		if registerErr != nil {
			return true, fmt.Errorf("finder.RegisterSrvWithAddr fail -> addr:%v,version:%s", addr, version)
		}
		loggerStd.Printf("success call finder.RegisterSrvWithAddr -> addr:%v,version:%s\n", addr, version)
		return true, nil
	}
	return false, nil
}
func (f *finderAdapter) Closeout() {
	loggerStd.Println("about to UnRegisterSrv")
	f.finderRWMU.RLock()
	defer f.finderRWMU.RUnlock()
	for k := range f.finderMap {
		_ = f.finderMap[k].UnRegisterSrvWithAddr(f.versionMap[k], k)
	}
}
