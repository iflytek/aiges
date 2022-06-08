package utils

import (
	"github.com/xfyun/finder-go"
	"github.com/xfyun/finder-go/common"
	"time"
)

type FindManger struct {
	//cfg *utils.Configure
	fm *finder.FinderManager
}

//var findManagerInst *FindManger

func NewFinder(co *CfgOption) (*FindManger, error) {

	//if findManagerInst != nil {
	//	return findManagerInst, nil
	//}

	if 0 == co.tick {
		co.tick = 50 * time.Second
	}
	if 0 == co.stmout {
		co.stmout = 5 * time.Second
	}
	so := common.BootConfig{
		CompanionUrl:  co.url,
		CachePath:     co.cachePath,
		CacheConfig:   co.cacheConfig,
		CacheService:  co.cacheService,
		ExpireTimeout: co.stmout,
		MeteData: &common.ServiceMeteData{
			Project: co.prj,
			Group:   co.group,
			Service: co.srv,
			Version: co.ver,
		},
	}
	/* var log *customLogImpl
	 if co.log != nil{
		 log = NewLogImpl(co.log)
	 }*/
	fm, e := finder.NewFinderWithLogger(so, co.log)
	if e != nil {
		return nil, e
	}

	f := new(FindManger)
	f.fm = fm
	//findManagerInst = f
	return f, nil
}

func DestroyFinder(fm *FindManger) {
	if fm != nil {
		finder.DestroyFinder(fm.fm)
	}
}

func (fm *FindManger) UseCfgAndSub(name string, cb common.ConfigChangedHandler) (map[string]*common.Config, error) {
	return fm.fm.ConfigFinder.UseAndSubscribeConfig([]string{name}, cb)
}

func (fm *FindManger) UseCfg(name string) (map[string]*common.Config, error) {
	return fm.fm.ConfigFinder.UseConfig([]string{name})
}

func (fm *FindManger) RegisterSrv(version string) error {
	return fm.fm.ServiceFinder.RegisterService(version)
}

func (fm *FindManger) RegisterSrvWithAddr(addr string, version string) error {
	return fm.fm.ServiceFinder.RegisterServiceWithAddr(addr, version)
}

func (fm *FindManger) UnRegisterSrv(version string) error {
	return fm.fm.ServiceFinder.UnRegisterService(version)
}
func (fm *FindManger) UnRegisterSrvWithAddr(version string, addr string) error {
	return fm.fm.ServiceFinder.UnRegisterServiceWithAddr(version, addr)
}

func (fm *FindManger) UseSrvAndSub(apiVersion, name string, handler common.ServiceChangedHandler) (map[string]common.Service, error) {
	return fm.fm.ServiceFinder.UseAndSubscribeService([]common.ServiceSubscribeItem{{ApiVersion: apiVersion, ServiceName: name}}, handler)

}

func (fm *FindManger) UnSubSrv(name string) error {
	return fm.fm.ServiceFinder.UnSubscribeService(name)
}
func (fm *FindManger) QueryService(project, group string) (map[string][]common.ServiceInfo, error) {
	return fm.fm.ServiceFinder.QueryService(project, group)
}
func (fm *FindManger) QueryServiceWatch(project, group string, handler common.ServiceChangedHandler) (map[string][]common.ServiceInfo, error) {
	return fm.fm.ServiceFinder.QueryServiceWatch(project, group, handler)
}
