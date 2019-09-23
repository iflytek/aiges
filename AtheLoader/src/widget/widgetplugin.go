package widget

import (
	"errors"
	"fmt"
	"plugin"
	"reflect"
	"service"
)

// 定义初始plugin.Symbol用于plugin类型校验
var (
	wrapperInit      wrapperInitPtr
	wrapperFini      wrapperFiniPtr
	wrapperVersion   wrapperVersionPtr
	wrapperLoadRes   wrapperLoadResPtr
	wrapperUnloadRes wrapperUnloadResPtr
	wrapperCreate    wrapperCreatePtr
	wrapperDestroy   wrapperDestroyPtr
	wrapperWrite     wrapperWritePtr
	wrapperRead      wrapperReadPtr
	wrapperOnceExec  wrapperOnceExecPtr
	wrapperDebugInfo wrapperDebugInfoPtr
)

var pluginHdls = map[string]plugin.Symbol{
	wrapperInitSym:    wrapperInit,
	wrapperFiniSym:    wrapperFini,
	wrapperVerSym:     wrapperVersion,
	wrapperLoadSym:    wrapperLoadRes,
	wrapperUnloadSym:  wrapperUnloadRes,
	wrapperCreateSym:  wrapperCreate,
	wrapperWriteSym:   wrapperWrite,
	wrapperReadSym:    wrapperRead,
	wrapperDestroySym: wrapperDestroy,
	wrapperExecSym:    wrapperOnceExec,
	wrapperDebugSym:   wrapperDebugInfo,
}

var pluginTypes = map[string]reflect.Type{
	wrapperInitSym:    reflect.TypeOf(func(cfg map[string]string) (err error) { return }),
	wrapperFiniSym:    reflect.TypeOf(func() (err error) { return }),
	wrapperVerSym:     reflect.TypeOf(func() (version string) { return }),
	wrapperLoadSym:    reflect.TypeOf(func(res wrapperData, resId int) (err error) { return }),
	wrapperUnloadSym:  reflect.TypeOf(func(resId int) (err error) { return }),
	wrapperCreateSym:  reflect.TypeOf(func(param map[string]string, prsIds []int, cb CallBackPtr) (hdl interface{}, err error) { return }),
	wrapperWriteSym:   reflect.TypeOf(func(hdl interface{}, req []wrapperData) (err error) { return }),
	wrapperReadSym:    reflect.TypeOf(func(hdl interface{}) (resp []wrapperData, err error) { return }),
	wrapperDestroySym: reflect.TypeOf(func(hdl interface{}) (debug string) { return }),
	wrapperExecSym:    reflect.TypeOf(func(param map[string]string, req []wrapperData) (resp []wrapperData, err error) { return }),
	wrapperDebugSym:   reflect.TypeOf(func(hdl interface{}) (debug string) { return }),
}

// 加载器 go plugin集成方案实现
type WidgetGo struct {
	eng *plugin.Plugin
}

func (inst *WidgetGo) Open() (errInfo error) {
	inst.eng, errInfo = plugin.Open(wrapperGo)
	if errInfo == nil {
		for sym := range pluginHdls {
			hdl, err := inst.eng.Lookup(sym)
			if err != nil {
				return errors.New(fmt.Sprintf("plugin lookup sym %s fail", sym))
			}
			// 类型校验
			if reflect.TypeOf(hdl) != (pluginTypes[sym]) {
				return errors.New(fmt.Sprintf("typeof plugin symbol:%s invalid", sym))
			}
			pluginHdls[sym] = hdl
		}
	}
	return
}

func (inst *WidgetGo) Close() {
	// nothing to do
	return
}

func (inst *WidgetGo) Register(srv *service.EngService) (errInfo error) {
	// 注册"事件-行为"对
	errInfo = srv.Register(service.EventUsrInit, pluginInit)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrFini, pluginFini)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrNew, pluginCreate)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrDel, pluginDestroy)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrRead, pluginRead)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrWrite, pluginWrite)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrOnceExec, pluginOnceExec)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrExcp, pluginExcp)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrDebug, pluginDebug)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrResLoad, pluginLoadRes)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrResUnload, pluginUnloadRes)
	if errInfo != nil {
		return
	}
	return
}

func (inst *WidgetGo) Version() (ver string) {
	return pluginVersion()
}
