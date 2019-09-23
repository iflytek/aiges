package service

import (
	"conf"
	"errors"
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"instance"
	"log"
	"sync"
)

type EngService struct {
	xsfInst    xsf.XsfServer // xsf基础框架实例
	aiInst     aiService
	wg         sync.WaitGroup
	usrVersion string
	usrActions map[usrEvent]actionUser // 其他事件类型(支持复合事件:eventOpAIIn & eventDataFirst)
	initAction actionInit              // 引擎初始化事件
	finiAction actionFini              // 引擎逆初始化事件
}

/*
	该接口用于框架初始化,提供框架初始化功能,包括各类系统设置及环境变量设置;
*/
func (srv *EngService) Init(srvVer string) (errInfo error) {
	srv.usrVersion = srvVer
	conf.SvcVersion = srvVer
	srv.usrActions = make(map[usrEvent]actionUser)
	return
}

/*
	该接口用于集成方注册"事件-行为"对, 框架在注册的对应条件发生时,会调用注册的对应方法action;
	@param et			事件类型：业务初始化|业务逆初始化|用户自定义
	@param event		事件描述; (仅用户自定义事件需要进行描述,描述规则见"docs")
	@param action		事件对应的行为;
*/
func (srv *EngService) Register(event usrEvent, action interface{}) (errInfo error) {
	switch event {
	case EventUsrInit:
		{
			// 类型校验
			init, ok := action.(func(map[string]string) (errNum int, errInfo error))
			if ok {
				srv.initAction = actionInit(init)
				return
			}
		}
	case EventUsrFini:
		{
			// 类型校验
			fini, ok := action.(func() (errNum int, errInfo error))
			if ok {
				srv.finiAction = actionFini(fini)
				return
			}
		}
	default:
		{
			// 类型校验 && 事件解析校验
			usrAct, ok := action.(func(hdl string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error))
			if ok {
				srv.usrActions[event] = actionUser(usrAct)
				return
			}
		}

	}
	return errors.New("EngService.Register: invalid action, event " + eventToString(event))
}

/*
	该接口为框架服务运行接口,调用之后框架实际运行进行各类模块初始化/服务初始化/并监听端口接收请求;
	@note	该接口与RunWithWidget()区别,需要自行服务框架register注册操作;
*/
func (srv *EngService) Run() (errInfo error) {
	srv.wg.Add(1)
	go func() {
		defer func() {
			srv.wg.Done()
			fmt.Println("EngService exit")
		}()

		srv.aiInst = aiService{callbackInit: srv.initAction, callbackFini: srv.finiAction, callbackUser: srv.usrActions}
		xsf.AddKillerCheck(SERVICE, &sigClose{srv})

		if errInfo = srv.xsfInst.Run(xsf.BootConfig{CfgMode: utils.CfgMode(-1),
			CfgData: xsf.CfgMeta{CfgName: "", Project: "", Group: "", Service: "",
				Version: srv.usrVersion, CompanionUrl: ""}}, &srv.aiInst); errInfo != nil {
			log.Fatal(errInfo)
		}
	}()
	srv.wg.Wait()
	return
}

/*
	该接口用于框架逆初始化;
*/
func (srv *EngService) Fini() {
	// nothing to do, Fini by `kill aiges`;
	// 逆初始化接口对用户不可见, 防止用户误使用导致下线过程中服务异常;
	return
}

/*
	该接口用于获取框架版本号
*/
func (srv *EngService) Version() string {
	return VERSION
}
