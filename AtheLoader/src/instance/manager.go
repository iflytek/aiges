/*
	实例管理器,用于管理服务框架实例
	1. 并非用户引擎计算实例;
	2. 支持会话&非会话模式,非会话模式不做授权控制;
*/
package instance

import (
	"buffer"
	"catch"
	"conf"
	"frame"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"strconv"
	"sync"
)

var svcAddr string
var svcPort string

type usrData struct {
	inst  *ServiceInst
	mngr  *Manager
	appid string
}

type Manager struct {
	sessMode   bool
	instCaches map[*ServiceInst]bool /*idle*/
	instMutex  sync.Mutex
	// session cache
	tool *xsf.ToolBox

	// session history info
	hisSess  [hstSessionSize]string
	hisIndex int
	hisMutex sync.Mutex

	delChan  chan string
	actMap   map[UserEvent]UsrAct

	nonSessCache map[string] /*handle*/ *ServiceInst
	cacheMutex   sync.Mutex
}

// TODO 非会话模式的请求并发管理,待讨论,计划引入性能指标及lb改造
func (mngr *Manager) Init(sm bool, lic int, delThr int, act map[UserEvent]UsrAct, tool *xsf.ToolBox) (errInfo error) {
	// 会话模式
	if sm {
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
	} else {
		mngr.nonSessCache = make(map[string]*ServiceInst)
	}

	// 异步delete协程; // 考虑外部重试可能多次调用导致写满channel阻塞, channel长度设置为授权5倍;
	mngr.delChan = make(chan string, lic*5)
	for j := 0; j < delThr; j++ {
		go mngr.asyncRelease()
	}

	// manager cache
	mngr.tool = tool
	mngr.actMap = act
	mngr.sessMode = sm

	svcAddr = tool.NetManager.GetIp()
	svcPort = strconv.Itoa(tool.NetManager.GetPort())
	return
}

func (mngr *Manager) Fini() {
	for inst := range mngr.instCaches {
		inst.Fini()
		delete(mngr.instCaches, inst)
	}
}

/*
	Acquire&Query接口拆分原因：
	0. 外部进行"事件-行为"判定,不同事件对应不同接口行为,单一职责接口;
	1. 若使用同一接口,根据handle无法判定行为,若查询不到是否需要申请实例;
	2. 若dataWrite请求先与sessionBegin到达,manager无法判定;
*/
// 获取会话实例(会话&非会话);
func (mngr *Manager) Acquire(handle string, param map[string]string) (inst *ServiceInst, errNum int, errInfo error) {
	if mngr.sessMode {
		mngr.instMutex.Lock()
		appid, _ := param["appid"]
		for srvInst := range mngr.instCaches {
			if mngr.instCaches[srvInst] {
				mngr.instCaches[srvInst] = false
				mngr.instMutex.Unlock()
				errInfo = mngr.tool.Cache.SetSessionData(handle, &usrData{srvInst, mngr, appid}, CCReportCallBack)
				if errInfo != nil { // reset if setSessionData fail
					mngr.instMutex.Lock()
					mngr.instCaches[srvInst] = true
					mngr.instMutex.Unlock()
					errNum = frame.AigesErrorLicNotEnough
					mngr.tool.Log.Errorw("sessMngr SetSessionData fail", "hdl", handle, "errNum", errNum, "errInfo", errInfo.Error())
				} else {
					inst = srvInst
					inst.context = handle
					mngr.tool.Cache.UpdateDelay()
					appidCCInc(appid)
				}
				return
			}
		}
		errNum = frame.AigesErrorLicNotEnough
		errInfo = frame.ErrorInstNotEnouth
		mngr.instMutex.Unlock()
		mngr.tool.Log.Errorw("instManager acquire fail, service license not enough", "errNum", errNum)
	} else {
		// 非会话模式无需控制授权及流量,由xsf:qpsLimiter功能支持,直接返回可用实例;
		inst = &ServiceInst{context: handle, mngr: mngr}
		errInfo = inst.Init(mngr.actMap, mngr.tool)
		if errInfo != nil {
			errNum = frame.AigesErrorLicNotEnough
			mngr.tool.Log.Errorw("instManager nonSession initialize fail", "errInfo", errInfo.Error())
			return
		}
		mngr.cacheMutex.Lock()
		defer mngr.cacheMutex.Unlock()
		mngr.nonSessCache[handle] = inst
	}
	return
}

// 查询会话实例(会话模式, Session Continue)
func (mngr *Manager) Query(handle string) (inst *ServiceInst, errNum int, errInfo error) {
	// 非会话模式,不支持会话查询;
	if !mngr.sessMode {
		errNum = frame.AigesErrorInvalidSessMode
		errInfo = frame.ErrorSessNotSupport
		return
	}

	srvData, errCache := mngr.tool.Cache.GetSessionData(handle)
	if errCache != nil {
		if mngr.isHistory(handle) {
			errNum = frame.AigesErrorEngInactive
			errInfo = frame.ErrorInstNotActive
		} else {
			errNum = frame.AigesErrorInvalidHdl
			errInfo = frame.ErrorInvalidInstHdl
		}
		mngr.tool.Log.Errorw("instManager query inst fail", "errNum", errNum)
	} else {
		inst = srvData.(*usrData).inst
	}
	return
}

func (mngr *Manager) Release(handle string) (errNum int, errInfo error) {
	// 写入异步delete channel; TODO 暂不做过滤或超时处理;故授权并发上报迁移至SessCallBack
	mngr.setHistory(handle)
	mngr.delChan <- handle
	return
}

func (mngr *Manager) cacheRelease(handle string, inst *ServiceInst) (errNum int, errInfo error) {
	mngr.instMutex.Lock()
	defer mngr.instMutex.Unlock()
	if mngr.sessMode {
		mngr.instCaches[inst] = true
	} else {
		// 非会话模式无需处理,runtime gc回收资源;
	}
	return
}

// 会话模式场景,历史会话数据
func (mngr *Manager) setHistory(handle string) {
	mngr.hisMutex.Lock()
	mngr.hisSess[mngr.hisIndex] = handle
	mngr.hisIndex++
	mngr.hisIndex %= hstSessionSize
	mngr.hisMutex.Unlock()
	return
}

// 历史会话查询
func (mngr *Manager) isHistory(handle string) bool {
	mngr.hisMutex.Lock()
	for i := range mngr.hisSess {
		if mngr.hisSess[i] == handle {
			mngr.hisMutex.Unlock()
			return true
		}
	}
	mngr.hisMutex.Unlock()
	return false
}

func (mngr *Manager) asyncRelease() {
	for {
		select {
		case sessionTag := <-mngr.delChan:
			if mngr.sessMode {
				mngr.tool.Cache.DelSessionData(sessionTag)
				// 该操作触发xsf cache调用SessCallBack操作;
			} else {
				// 非会话模式释放,查找cache释放资源 TODO 非会话模式如何与xsf同步压力状态,是否需要同步
				mngr.cacheMutex.Lock()
				inst, ok := mngr.nonSessCache[sessionTag]
				delete(mngr.nonSessCache, sessionTag)
				mngr.cacheMutex.Unlock()
				if ok {
					inst.Fini()
				}
			}
		}
	}
}

// @caller	回调方及回调原因;
func CCReportCallBack(sessionTag interface{}, svcData interface{}, caller ...xsf.CallBackException) {
	usr := svcData.(*usrData)
	// service instance release wait
	usr.inst.context = "" // reset context
	usr.inst.instWg.Wait()
	usr.inst.resetExcp()
	// service manager release
	usr.mngr.cacheRelease(sessionTag.(string), usr.inst)
	usr.mngr.tool.Cache.UpdateDelay()
	appidCCDec(usr.appid)
	usr.mngr.tool.Log.Infow("Call SessCallBack", "hdl", sessionTag)
}

func (mngr *Manager) CatchCallBack(tag string) (reqDoubt []catch.TagRequest) {
	switch conf.SessMode {
	case true:
		mngr.instMutex.Lock()
		defer mngr.instMutex.Unlock()
		for inst, idle := range mngr.instCaches {
			if !idle {
				td := make([]catch.TagData, 0, 1)
				params, datas := inst.CatchTag()
				for _, v := range datas {
					meta := catch.TagData{Fmt: v.dataFmt, Enc: v.dataEnc, Typ: buffer.DataTypeToString(v.dataType)}
					switch v.dataType {
					case buffer.DataText:
						meta.Data = []byte(v.data.(string))
					default:
						meta.Data = v.data.([]byte)
					}
					td = append(td, meta)
				}
				tr := catch.TagRequest{}
				tr.Sid, _ = params[sessionId]
				tr.Param = params
				tr.DataList = td
				reqDoubt = append(reqDoubt, tr)
			}
		}

	default:
		mngr.cacheMutex.Lock()
		defer mngr.cacheMutex.Unlock()
		for _, inst := range mngr.nonSessCache {
			td := make([]catch.TagData, 0, 1)
			params, datas := inst.CatchTag()
			for _, v := range datas {
				meta := catch.TagData{Fmt: v.dataFmt, Enc: v.dataEnc, Typ: buffer.DataTypeToString(v.dataType)}
				switch v.dataType {
				case buffer.DataText:
					meta.Data = []byte(v.data.(string))
				default:
					meta.Data = v.data.([]byte)
				}
				td = append(td, meta)
			}
			tr := catch.TagRequest{}
			tr.Sid, _ = params[sessionId]
			tr.Param = params
			tr.DataList = td
			reqDoubt = append(reqDoubt, tr)
		}
	}
	return
}
