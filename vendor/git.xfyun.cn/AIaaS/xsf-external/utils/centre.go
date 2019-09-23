package utils

import (
	"errors"
	"git.xfyun.cn/AIaaS/finder-go/common"
	"sync"
)

var (
	ECFGNOTFOUND    = errors.New("File not found")
	ECFGINVALREADER = errors.New("Invaled cfg reader")
)

var (
	centerManger = make(map[string]*Configure)
	cmMutex      = new(sync.Mutex)
)

type centerReader struct {
	fm *FindManger
	f  map[string]*finder.Config
}

/*
func newCentre(co *CfgOption) (*Configure, error){

	so := common.BootConfig{
		CompanionUrl:co.url,
		TickerDuration:   5000,
		MeteData: common.ServiceMeteData{
				Project: co.prj,
				Group:   co.group,
				Service: co.srv,
				Version: co.ver,
			},
	}

	var e error
	cc := new(centerReader)
	cfg := new(Configure)

	cc.fm, e = finder.NewFinder(so)
	if e != nil {
		return nil, e
	}

	// 订阅
	//var handler common.ConfigChangedHandler
	cc.f, e = cc.fm.ConfigFinder.UseAndSubscribeConfig([]string{co.name},cfg)
	fmt.Println("cc.fm.ConfigFinder.UseAndSubscribeConfig",e,cc.f)
	if e != nil {
		return nil, e
	}


	e = cfg.init( cc, co)
	if e != nil {
		return nil, e
	}

	cmMutex.Lock()
	centerManger[co.name] =cfg
	cmMutex.Unlock()

	return cfg, nil
}
*/

func NewCentreWithFinder(co *CfgOption) (*Configure, error) {

	if nil == co.fm {
		return nil, ECFGINVALREADER
	}

	cmMutex.Lock()
	v, ok := centerManger[co.name]
	cmMutex.Unlock()
	if ok {
		return v, nil
	}

	// 重新申请配置对象
	var e error
	var cc *centerReader
	cfg := new(Configure)
	if nil == cfg.r {
		cfg.r = new(centerReader)
	}

	cc = cfg.r.(*centerReader)

	if nil == cc.fm {
		// 设置findmanger
		cc.fm = co.fm
	}

	// 订阅
	cc.f, e = cc.fm.UseCfgAndSub(co.name, cfg)
	if nil != e {
		return nil, e
	}
	e = cfg.init(cc, co)
	if nil != e {
		return nil, e
	}

	return cfg, nil
}

func (cc *centerReader) Read(name string) (string, error) {
	v, o := cc.f[name]
	if o {
		return string(v.File), nil
	}

	return "", ECFGNOTFOUND
}

/*
func (cc *centerReader)FindAll(name string)(string, error){

	cc.f, e = cc.fm.ServiceFinder.
}*/
