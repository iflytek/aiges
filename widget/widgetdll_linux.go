//go:build linux || cgo
// +build linux cgo

/*
	控件层通过注册"事件-行为"方法对,实现对控件接口的注册及上层接口调用的统一调度适配；
	1. 控件层定义各业务的业务逻辑接口,使之更适配业务属性及接口调用逻辑；
	2. 控件接口实现了上层golang至底层c/c++的不同语言栈数据类型转换适配；
	3. 控件层需实现如下接口：Run()用于注册"事件-行为"; WrapperVersion()用于获取服务版本;
*/
package widget

import (
	"github.com/xfyun/aiges/service"
)

type WidgetC struct {
	eng engineC
}

/*
	控件初始化&逆初始化
	@param clib	适配引擎库;
*/
func (inst *WidgetC) Open() (errInfo error) {
	return inst.eng.open(wrapperC)
}

func (inst *WidgetC) Close() {
	inst.eng.close()
	return
}

/*
	控件入口,框架实现通用AI能力适配层,适配当前各类AI引擎;
	@param srv	服务框架实例,用于控件向框架注册"事件-行为"
*/
func (inst *WidgetC) Register(srv *service.EngService) (errInfo error) {
	// 注册"事件-行为"对
	errInfo = srv.Register(service.EventUsrInit, engineInit)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrFini, engineFini)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrNew, engineCreate)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrDel, engineDestroy)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrRead, engineRead)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrWrite, engineWrite)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrOnceExec, engineOnceExec)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrExcp, engineExcp)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrDebug, engineDebug)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrResLoad, engineLoadRes)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrResUnload, engineUnloadRes)
	if errInfo != nil {
		return
	}
	return
}

func (inst *WidgetC) Version() (ver string) {
	return engineVersion()
}
