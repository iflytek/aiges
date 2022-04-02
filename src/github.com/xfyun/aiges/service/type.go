package service

import (
	"github.com/xfyun/aiges/instance"
)

// AI General Engine Service Version Information
const (
	VERSION string = "alpha"
)

// AI General Engine Service Operation
const (
	opAIEngIn   string = "AIIn"
	opAIEngOut  string = "AIOut"
	opException string = "AIExcp"
	opResUpdate string = "AIHotword"
	opLBType    string = "AILBType"
	opDetect    string = "AIDetect"
	opAILoad    string = "AILoad"
	opAIUnLoad  string = "AIUnLoad"
)

// 交互协议数据定义
const (
	reqVersion      = "version" // 请求pb协议版本 v1:AIaaS协议 v2:AIPaaS协议
	reqNumber       = "SeqNo"
	reqWaitTime     = "waitTime"
	reqFirst        = "1"
	syncRespTimeout = 500
)

// 注册事件类型(参照数据协议扩展)
type usrEvent instance.UserEvent

const (
	// 初始化&逆初始化
	EventUsrInit = usrEvent(instance.EventInit) //	业务服务初始化
	EventUsrFini = usrEvent(instance.EventFini) //	业务服务逆初始化

	EventUsrNew   = usrEvent(instance.EventNew) // 资源申请事件
	EventUsrDel   = usrEvent(instance.EventDel) // 资源释放事件
	EventUsrExcp  = usrEvent(instance.EventExcp)
	EventUsrDebug = usrEvent(instance.EventDebug)

	EventUsrRead     = usrEvent(instance.EventRead)     // 数据读事件
	EventUsrWrite    = usrEvent(instance.EventWrite)    // 数据写事件
	EventUsrOnceExec = usrEvent(instance.EventOnceExec) // 非会话处理事件

	EventUsrResLoad   = usrEvent(instance.EventResLoad)   // 资源加载事件
	EventUsrResUnload = usrEvent(instance.EventResUnload) // 资源卸载事件
)

// usrEvent trans to string;
func eventToString(ue usrEvent) (ueStr string) {
	switch ue {
	case EventUsrInit:
		return "EventUsrInit"
	case EventUsrFini:
		return "EventUsrFini"
	case EventUsrNew:
		return "EventUsrNew"
	case EventUsrDel:
		return "EventUsrDel"
	case EventUsrExcp:
		return "EventUsrExcp"
	case EventUsrDebug:
		return "EventUsrDebug"
	case EventUsrRead:
		return "EventUsrRead"
	case EventUsrWrite:
		return "EventUsrWrite"
	case EventUsrResLoad:
		return "EventUsrResLoad"
	case EventUsrResUnload:
		return "EventUsrResUnload"
	}
	return
}

// 业务初始化行为,对应事件EventInit
type actionInit func(map[string]string) (errNum int, errInfo error)

// 业务逆初始化行为,对应事件EventFini
type actionFini func() (errNum int, errInfo error)

// 用户自定义行为,对应事件EventUsrDefine
type actionUser instance.UsrAct
