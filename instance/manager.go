/*
	实例管理器,用于管理服务框架实例
	1. 并非用户引擎计算实例;
	2. 支持会话&非会话模式,非会话模式不做授权控制;
*/
package instance

import (
	"context"
	"fmt"
	"github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/frame"
	"github.com/xfyun/xsf/server"
	"strconv"
	"sync"
	"time"
)

var svcAddr string
var svcPort string

type usrData struct {
	inst *ServiceInst
}

type Manager struct {
	instCaches map[*ServiceInst]bool /*idle*/
	instMutex  sync.Mutex
	// session cache
	tool    *xsf.ToolBox
	action  map[UserEvent]UsrAct
	delChan chan string

	// session history info
	hisSess  [hstSessionSize]string
	hisIndex int
	hisMutex sync.Mutex
}

// TODO 非会话模式的请求并发管理,待讨论,计划引入性能指标及lb改造
func (mngr *Manager) Init(lic int, delThr int, act map[UserEvent]UsrAct, tool *xsf.ToolBox) (errInfo error) {
	// 预申请实例管理
	mngr.tool = tool
	mngr.instCaches = make(map[*ServiceInst]bool)
	for i := 0; i < lic; i++ {
		inst := ServiceInst{mngr: mngr}
		errInfo = inst.Init(act, tool)
		if errInfo != nil {
			for i := range mngr.instCaches {
				i.Fini()
			}
			tool.Log.Errorw("instManager initialize fail", "errInfo", errInfo.Error())
			return
		}
		mngr.instCaches[&inst] = true
	}

	// 异步delete协程; // 考虑外部重试可能多次调用导致写满channel阻塞, channel长度设置为授权5倍;
	mngr.delChan = make(chan string, lic*5)
	for j := 0; j < delThr; j++ {
		go mngr.asyncRelease()
	}

	mngr.action = make(map[UserEvent]UsrAct)
	for e, a := range act {
		mngr.action[e] = a
	}

	svcAddr = tool.NetManager.GetIp()
	svcPort = strconv.Itoa(tool.NetManager.GetPort())
	return
}

func (mngr *Manager) Fini() {
	for inst := range mngr.instCaches {
		inst.Fini()
		delete(mngr.instCaches, inst)
	}
	fmt.Println("aiService.Finit: fini instMngr success!")
}

/*
	Acquire&Query接口拆分原因：
	0. 外部进行"事件-行为"判定,不同事件对应不同接口行为,单一职责接口;
	1. 若使用同一接口,根据handle无法判定行为,若查询不到是否需要申请实例;
	2. 若dataWrite请求先与sessionBegin到达,manager无法判定;
*/
// 获取会话实例(会话&非会话);
func (mngr *Manager) Acquire(handle string, param map[string]string) (inst *ServiceInst, errNum int, errInfo error) {
	mngr.instMutex.Lock()
	defer mngr.instMutex.Unlock()
	for srvInst := range mngr.instCaches {
		if mngr.instCaches[srvInst] {
			mngr.instCaches[srvInst] = false
			errInfo = mngr.tool.Cache.SetSessionData(handle, &usrData{srvInst}, CCReportCallBack)
			if errInfo != nil { // reset if setSessionData fail
				mngr.instCaches[srvInst] = true
				errNum = frame.AigesErrorLicNotEnough
				mngr.tool.Log.Errorw("sessMngr SetSessionData fail", "hdl", handle, "errNum", errNum, "errInfo", errInfo.Error())
			} else {
				inst = srvInst
				inst.context = handle
				mngr.tool.Cache.UpdateDelay()
				sampleEnter(&param)
			}
			return
		}
	}
	errNum = frame.AigesErrorLicNotEnough
	errInfo = frame.ErrorInstNotEnouth
	mngr.tool.Log.Errorw("instManager acquire fail, service license not enough", "errNum", errNum)
	return
}

// 查询会话实例(会话模式, Session Continue)
func (mngr *Manager) Query(handle string) (inst *ServiceInst, errNum int, errInfo error) {
	srvData, errCache := mngr.tool.Cache.GetSessionData(handle)
	if errCache != nil {
		errNum, errInfo = frame.AigesErrorInvalidHdl, frame.ErrorInvalidInstHdl
		if mngr.isHistory(handle) {
			errNum, errInfo = frame.AigesErrorEngInactive, frame.ErrorInstNotActive
		}
		mngr.tool.Log.Errorw("instManager query inst fail", "hdl", handle, "errNum", errNum)
	} else {
		inst = srvData.(*usrData).inst
	}
	return
}

func (mngr *Manager) Release(handle string) (errNum int, errInfo error) {
	// 写入异步delete channel; TODO 暂不做过滤或超时处理;故授权并发上报迁移至SessCallBack
	mngr.setHistory(handle)
	if conf.AsyncRelease {
		mngr.delChan <- handle
	} else {
		mngr.tool.Cache.DelSessionData(handle)
	}
	return
}

func (mngr *Manager) UpdateLic(param map[string]string) (errNum int, errInfo error) {
	xsf.SetLbType(param[reqLbType])
	lic, err := strconv.Atoi(param[reqDynLic])
	if err != nil {
		return frame.AigesErrorInvalidParaValue, frame.ErrorInvalidParaValue
	}
	tm, err := strconv.Atoi(param[reqDynTime])
	if err != nil {
		return frame.AigesErrorInvalidParaValue, frame.ErrorInvalidParaValue
	}

	// 动态授权调整;
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tm)*time.Millisecond)
	defer cancel()
	err = mngr.tool.Cache.UpdateOTF(ctx, xsf.WithSessionManagerMaxLicOTF(int32(lic)))
	if err != nil {
		return frame.AigesErrorElasticLic, err // frame.ErrorElasticLicInst
	}

	// TODO 服务缓存实例调整;
	mngr.instMutex.Lock()
	defer mngr.instMutex.Unlock()
	for diff := lic - len(mngr.instCaches); diff > 0; diff-- {
		inst := ServiceInst{mngr: mngr}
		errInfo = inst.Init(mngr.action, mngr.tool)
		if errInfo != nil {
			mngr.tool.Log.Errorw("UpdateLic initialize inst fail", "errInfo", errInfo.Error())
			return frame.AigesErrorInvalidParaValue, errInfo
		}
		mngr.instCaches[&inst] = true
	}

	return
}

func (mngr *Manager) cacheRelease(handle string, inst *ServiceInst) (errNum int, errInfo error) {
	mngr.instMutex.Lock()
	defer mngr.instMutex.Unlock()
	mngr.instCaches[inst] = true
	return
}

// 会话模式场景,历史会话数据
func (mngr *Manager) setHistory(handle string) {
	mngr.hisMutex.Lock()
	defer mngr.hisMutex.Unlock()
	mngr.hisSess[mngr.hisIndex] = handle
	mngr.hisIndex++
	mngr.hisIndex %= hstSessionSize
	return
}

// 历史会话查询
func (mngr *Manager) isHistory(handle string) bool {
	mngr.hisMutex.Lock()
	defer mngr.hisMutex.Unlock()
	for i := range mngr.hisSess {
		if mngr.hisSess[i] == handle {
			return true
		}
	}
	return false
}

func (mngr *Manager) asyncRelease() {
	if conf.AsyncRelease {
		for {
			select {
			case sessionTag := <-mngr.delChan:
				mngr.tool.Cache.DelSessionData(sessionTag)
			}
		}
	} else {
		fmt.Println("sync release handle")
	}
}

// @caller	回调方及回调原因;
func CCReportCallBack(sessionTag interface{}, svcData interface{}, caller ...xsf.CallBackException) {
	usr := svcData.(*usrData)
	// service instance release wait
	select {
	case usr.inst.onceTrig <- true:
	default:
	}
	usr.inst.context = "" // reset context
	usr.inst.setAlive(false)
	func() {
		defer func() {
			if recover() != nil {
				//有可能调用Inputdata的时候，signal已经被关闭了
			}
		}()
		usr.inst.inputData.Signal()
	}()
	usr.inst.instWg.Wait()
	usr.inst.resetExcp()
	// service manager release
	usr.inst.mngr.cacheRelease(sessionTag.(string), usr.inst)
	usr.inst.mngr.tool.Cache.UpdateDelay()
	sampleExit(&usr.inst.headers)
	usr.inst.mngr.tool.Log.Infow("Call SessCallBack", "hdl", sessionTag)
}
