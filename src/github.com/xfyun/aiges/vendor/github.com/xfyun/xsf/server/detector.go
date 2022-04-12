package xsf

import (
	"fmt"
	finder "github.com/xfyun/finder-go/common"
	"github.com/xfyun/xsf/utils"
	"sync"
	"time"
)

type detector struct {
	finderInst *utils.FindManger
	addrCache  *sync.Map
}

func (d *detector) OnServiceInstanceConfigChanged(name string, version string, instance string, config *finder.ServiceInstanceConfig) bool {
	fmt.Println("OnServiceInstanceConfigChanged", name, version, instance, config)
	if config.IsValid {
		loggerStd.Println("OnServiceInstanceConfigChanged insert adddr: ", instance)
		d.addrCache.Store(instance, true)
	} else {
		loggerStd.Println("OnServiceInstanceConfigChanged delete adddr: ", instance)
		d.addrCache.Delete(instance)
	}
	return true
}
func (d *detector) OnServiceConfigChanged(name, version string, config *finder.ServiceConfig) bool {
	fmt.Println("OnServiceConfigChanged", name, version, config)
	return true
}
func (d *detector) OnServiceInstanceChanged(name, version string, instances []*finder.ServiceInstanceChangedEvent) bool {
	for _, v := range instances {
		if v.EventType == finder.INSTANCEADDED {
			for _, inst := range v.ServerList {
				loggerStd.Println("OnServiceInstanceChanged insert addrs:", inst.Addr)
				d.addrCache.Store(inst.Addr, true)
			}
		} else if v.EventType == finder.INSTANCEREMOVE {
			for _, inst := range v.ServerList {
				loggerStd.Println("OnServiceInstanceChanged delete addrs:", inst.Addr)
				d.addrCache.Delete(inst.Addr)
			}
		}
	}
	return true
}

func newDetector(
	cfgUrl string,
	cfgPrj string,
	cfgGroup string,
	cfgName string,
	log *utils.Logger,
) (*detector, error) {
	co := &utils.CfgOption{}
	utils.WithCfgTick(time.Second)(co)
	utils.WithCfgSessionTimeOut(time.Second)(co)
	utils.WithCfgURL(cfgUrl)(co)
	utils.WithCfgCachePath(".")(co)
	utils.WithCfgCacheConfig(true)(co)
	utils.WithCfgCacheService(true)(co)
	utils.WithCfgPrj(cfgPrj)(co)
	utils.WithCfgGroup(cfgGroup)(co)
	utils.WithCfgLog(log)(co)
	finderInst, err := utils.NewFinder(co)
	if err != nil {
		return nil, err
	}

	detectorInst := &detector{finderInst: finderInst}

	addrCache, addrCacheErr := func() (*sync.Map, error) {
		m := &sync.Map{}
		srvs, srvsErr := finderInst.UseSrvAndSub("1.0.0", cfgName, detectorInst)
		if srvsErr != nil {
			return m, srvsErr
		}
		for _, v := range srvs {
			for _, provider := range v.ProviderList {
				m.Store(provider.Addr, true)
			}
		}
		return m, nil
	}()
	if addrCacheErr != nil {
		return nil, addrCacheErr
	}
	detectorInst.addrCache = addrCache
	return detectorInst, nil
}

func (d *detector) getAll() []string {
	var addrSet []string
	d.addrCache.Range(func(key, value interface{}) bool {
		if value.(bool) {
			addrSet = append(addrSet, key.(string))
		}
		return true
	})
	return addrSet
}
