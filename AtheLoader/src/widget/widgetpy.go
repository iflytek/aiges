package widget

import (
	"service"
)

type WidgetPy struct {
	//	eng string // py script
}

func (inst *WidgetPy) Open() (err error) {
	return pythonOpen(wrapperPy)
}

func (inst *WidgetPy) Close() {
	pythonClose()
}

func (inst *WidgetPy) Register(srv *service.EngService) (errInfo error) {
	// 注册"事件-行为"对
	errInfo = srv.Register(service.EventUsrInit, pythonInit)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrFini, pythonFini)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrNew, pythonCreate)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrDel, pythonDestroy)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrRead, pythonRead)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrWrite, pythonWrite)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrOnceExec, pythonOnceExec)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrExcp, pythonExcp)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrDebug, pythonDebug)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrResLoad, pythonLoadRes)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrResUnload, pythonUnloadRes)
	if errInfo != nil {
		return
	}
	return
}

func (inst *WidgetPy) Version() (ver string) {
	return pythonVersion()
}
