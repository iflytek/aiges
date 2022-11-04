/*
	控件层通过注册"事件-行为"方法对,实现对控件接口的注册及上层接口调用的统一调度适配；
	1. 控件层定义各业务的业务逻辑接口,使之更适配业务属性及接口调用逻辑；
	2. 控件接口实现了上层golang至底层c/c++的不同语言栈数据类型转换适配；
	3. 控件层需实现如下接口：Run()用于注册"事件-行为"; WrapperVersion()用于获取服务版本;
*/
package widget

import (
	"github.com/xfyun/aiges/service"
	"log"
)

type WidgetPython struct {
	eng enginePython
}

/*
	控件初始化&逆初始化
	@param clib	适配引擎库;
*/
func (inst *WidgetPython) Open() (errInfo error) {
	log.Println("Starting Using Python : ")
	return inst.eng.open()
}

func (inst *WidgetPython) Close() {
	inst.eng.close()
	return
}

/*
	控件入口,框架实现通用AI能力适配层,适配当前各类AI引擎;
	@param srv	服务框架实例,用于控件向框架注册"事件-行为"
*/
func (inst *WidgetPython) Register(srv *service.EngService) (errInfo error) {
	// 注册"事件-行为"对
	errInfo = srv.Register(service.EventUsrInit, inst.eng.enginePythonInit)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrFini, inst.eng.enginePythonFini)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrNew, inst.eng.enginePythonCreate)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrDel, inst.eng.enginePythonDestroy)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrRead, inst.eng.enginePythonRead)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrWrite, inst.eng.enginePythonWrite)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrOnceExec, inst.eng.enginePythonOnceExec)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrExcp, inst.eng.enginePythonExcp)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrDebug, inst.eng.enginePythonDebug)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrResLoad, inst.eng.enginePythonLoadRes)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrResUnload, inst.eng.enginePythonUnloadRes)
	if errInfo != nil {
		return
	}
	return
}

func (inst *WidgetPython) Version() (ver string) {
	return inst.eng.enginePythonVersion()
}
