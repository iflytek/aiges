package daemon

/**
*　　　　　　　　┏┓　　　┏┓+ +
*　　　　　　　┏┛┻━━━┛┻┓ + +
*　　　　　　　┃　　　　　　　┃
*　　　　　　　┃　　　━　　　┃ ++ + + +
*　　　　　　 ████━████ ┃+
*　　　　　　　┃　　　　　　　┃ +
*　　　　　　　┃　　　┻　　　┃
*　　　　　　　┃　　　　　　　┃ + +
*　　　　　　　┗━┓　　　┏━┛
*　　　　　　　　　┃　　　┃
*　　　　　　　　　┃　　　┃ + + + +
*　　　　　　　　　┃　　　┃　　　　Code is far away from bug with the animal protecting
*　　　　　　　　　┃　　　┃ + 　　　　神兽保佑,代码无bug
*　　　　　　　　　┃　　　┃
*　　　　　　　　　┃　　　┃　　+
*　　　　　　　　　┃　 　　┗━━━┓ + +
*　　　　　　　　　┃ 　　　　　　　┣┓
*　　　　　　　　　┃ 　　　　　　　┏┛
*　　　　　　　　　┗┓┓┏━┳┓┏┛ + + + +
*　　　　　　　　　　┃┫┫　┃┫┫
*　　　　　　　　　　┗┻┛　┗┻┛+ + + +
 */

import (
	"encoding/json"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type LoadMiniSubSvc struct {
	subSvc string

	addrRwMu sync.Mutex
	addrMap  map[string]*SubSvcItem //K:addr type:*SubSvcItem

	metaData *LoadMiniMeta
}

func (i *LoadMiniSubSvc) init(subSvc string, meta *LoadMiniMeta) {
	i.subSvc = subSvc
	i.metaData = meta
	i.addrMap = make(map[string]*SubSvcItem)
}

type LoadMiniMeta struct {
	ttl int64 //纳秒,定时清除无效节点的时间间隔

	ticker *time.Ticker //清除无用节点的扫描周期

	monitorInterval int64 //监控数据的更新时间,单位毫秒，缺省一分钟

	svc string

	toolbox *xsf.ToolBox
}

type LoadMiniSvc struct {
	svc           string
	metaData      *LoadMiniMeta
	subSvcMap     map[string]*LoadMiniSubSvc
	subSvcMapRwMu sync.RWMutex
}

func (i *LoadMiniSvc) init(svc string, meta *LoadMiniMeta) {
	i.svc = svc
	i.metaData = meta
	i.subSvcMap = make(map[string]*LoadMiniSubSvc)
	i.subSvcMap = make(map[string]*LoadMiniSubSvc)
}

type LoadMini struct {
	metaData           *LoadMiniMeta
	LoadMiniSvcMap     map[string]*LoadMiniSvc
	LoadMiniSvcMapRWMu sync.RWMutex
}

func newLoadMini(toolbox *xsf.ToolBox) *LoadMini {
	loadMiniTmp := &LoadMini{}
	loadMiniTmp.init(toolbox)
	return loadMiniTmp
}
func (i *LoadMini) toMonitorWareHouse(in *montorWareHouse) {
	in.Svc = i.metaData.svc
	in.Ttl = i.metaData.ttl
	in.MonitorInterval = i.metaData.monitorInterval
	in.toolbox = i.metaData.toolbox

	in.handle = i

	go in.loop()
}

/*
	更新全量数据至监控中
*/
func (i *LoadMini) traversal() *monitorSvc {
	//拷贝svc映射
	var loadMiniSvcMap = make(map[string]*LoadMiniSvc)
	i.LoadMiniSvcMapRWMu.RLock()
	for svc, loadSvc := range i.LoadMiniSvcMap {
		loadMiniSvcMap[svc] = loadSvc
	}
	i.LoadMiniSvcMapRWMu.RUnlock()

	var monitorLoadMiniTmp monitorSvc

	//使用svc映射
	for svc, loadMiniSvc := range loadMiniSvcMap {
		var monitorLoadMiniSubSvcTmp monitorSubSvc

		//拷贝subsvc映射
		var loadMiniSubSvcMap = make(map[string]*LoadMiniSubSvc)
		loadMiniSvc.subSvcMapRwMu.RLock()
		for subSvc, loadSubSvc := range loadMiniSvc.subSvcMap {
			loadMiniSubSvcMap[subSvc] = loadSubSvc
		}
		loadMiniSvc.subSvcMapRwMu.RUnlock()

		//使用subsvc映射
		for subSvc, loadMiniSubSvc := range loadMiniSubSvcMap {
			i.metaData.toolbox.Log.Infow("fn:traversal", "subsvc", subSvc)

			var monitorLoadMiniAddrTmp monitorAddr

			//拷贝subsvcfunc映射
			var subSvcItemMap = make(map[string]*SubSvcItem)
			loadMiniSubSvc.addrRwMu.Lock()
			for loadMiniSubSvcAddr, loadMiniSubSvcAuth := range loadMiniSubSvc.addrMap {
				subSvcItemMap[loadMiniSubSvcAddr] = loadMiniSubSvcAuth
			}
			loadMiniSubSvc.addrRwMu.Unlock()

			//使用subsvcfunc映射
			for loadMiniSubSvcAddr, loadMiniSubSvcAuth := range subSvcItemMap {
				var monitorSubSvcItemTmp monitorSubSvcItem

				monitorSubSvcItemTmp.Addr = loadMiniSubSvcAuth.addr
				monitorSubSvcItemTmp.Timestamp = loadMiniSubSvcAuth.timestamp
				monitorSubSvcItemTmp.BestInst = loadMiniSubSvcAuth.bestInst
				monitorSubSvcItemTmp.IdleInst = func() int64 {
					if loadMiniSubSvcAuth.idleInst < 0 {
						return 0
					}
					return loadMiniSubSvcAuth.idleInst
				}()
				monitorSubSvcItemTmp.TotalInst = loadMiniSubSvcAuth.totalInst

				if nil == monitorLoadMiniAddrTmp.AddrMap {
					monitorLoadMiniAddrTmp.AddrMap = make(map[string]*monitorSubSvcItem)
				}
				monitorLoadMiniAddrTmp.AddrMap[loadMiniSubSvcAddr] = &monitorSubSvcItemTmp
			}

			if nil == monitorLoadMiniSubSvcTmp.SubSvcMap {
				monitorLoadMiniSubSvcTmp.SubSvcMap = make(map[string]*monitorAddr)
			}
			monitorLoadMiniSubSvcTmp.SubSvcMap[subSvc] = &monitorLoadMiniAddrTmp

		}

		if nil == monitorLoadMiniTmp.SvcMap {
			monitorLoadMiniTmp.SvcMap = make(map[string]*monitorSubSvc)
		}
		monitorLoadMiniTmp.SvcMap[svc] = &monitorLoadMiniSubSvcTmp
	}
	return &monitorLoadMiniTmp
}

func (i *LoadMini) purge() {

	var purgeCnt int64
	for {
		select {
		case <-i.metaData.ticker.C:
			{
				logId := "purge@" +
					strconv.Itoa(int(atomic.AddInt64(&purgeCnt, 1))) +
					strconv.Itoa(time.Now().Nanosecond())

				i.metaData.toolbox.Log.Debugw("GarbageCollection=> begin to start collecting garbage",
					"logid", logId)

				/*
					※ 垃圾回收
						1、取出所有svc,减少锁粒度
						2、对于每个svc取出所有subsvc
						3、据svc+subsvc取出节点详细信息，剔除无效节点
				*/

				var subSvcDetaSet []*LoadMiniSubSvc

				i.metaData.toolbox.Log.Debugw("GarbageCollection=> about to collect svcDataSet",
					"logid", logId)
				var svcDataSet []*LoadMiniSvc
				var svcs []string
				i.LoadMiniSvcMapRWMu.RLock()
				for _, svcData := range i.LoadMiniSvcMap {
					svcs = append(svcs, svcData.svc)
					svcDataSet = append(svcDataSet, svcData)
				}
				i.LoadMiniSvcMapRWMu.RUnlock()
				i.metaData.toolbox.Log.Debugw("GarbageCollection=> svcSet",
					"svcs", svcs, "logid", logId)

				i.metaData.toolbox.Log.Debugw("GarbageCollection=> about to collect subSvcDetaSet",
					"logid", logId)
				var svcSubSvcs = make(map[string][]string)
				for _, svcData := range svcDataSet {
					svcData.subSvcMapRwMu.RLock()
					for subSvc, subSvcData := range svcData.subSvcMap {
						svcSubSvcs[svcData.svc] = append(svcSubSvcs[svcData.svc], subSvc)
						subSvcDetaSet = append(subSvcDetaSet, subSvcData)
					}
					svcData.subSvcMapRwMu.RUnlock()
				}

				/*
					此处以k
					判断是否对后续的移除操作添加db操作
				*/
				svcSubSvcsByte, _ := json.Marshal(svcSubSvcs)
				i.metaData.toolbox.Log.Debugw("GarbageCollection=> svcSet",
					"svcSubSvcs", string(svcSubSvcsByte), "logid", logId)

				i.metaData.toolbox.Log.Debugw("GarbageCollection=> about to deal subSvcDetaSet",
					"ts", time.Now(), "logid", logId)
				for _, subSvcData := range subSvcDetaSet {
					i.metaData.toolbox.Log.Debugw("GarbageCollection=> about to doPurge subSvcData",
						"logid", logId)

					i.doPurge(subSvcData, logId)
				}
			}
		}
	}
}

func (i *LoadMini) doPurge(subSvcInst *LoadMiniSubSvc, logId string) {
	subSvcInst.addrRwMu.Lock()
	i.metaData.toolbox.Log.Infow("about to clean dead nodes",
		"fn", "doPurge")
	for addr, addrInst := range subSvcInst.addrMap {
		nowTs := time.Now().UnixNano()
		nodeTtl := nowTs - addrInst.timestamp
		i.metaData.toolbox.Log.Infow(
			"in range subSvcInst.addrMap",
			"addr", addr, "addrInst", addrInst.String(),
			"fn", "doPurge", "nodeTtl", nodeTtl, "i.metaData.ttl", i.metaData.ttl)
		if nodeTtl > i.metaData.ttl {
			i.metaData.toolbox.Log.Warnw(
				"del addr from subSvcInst.addrMap",
				"addr", addr, "fn", "doPurge", "nodeTtl", nodeTtl, "i.metaData.ttl", i.metaData.ttl)
			delete(subSvcInst.addrMap, addr)
		}
	}
	subSvcInst.addrRwMu.Unlock()

}
func (i *LoadMini) init(toolbox *xsf.ToolBox) {
	if i.metaData == nil {
		i.metaData = &LoadMiniMeta{}
	}
	i.metaData.toolbox = toolbox

	i.LoadMiniSvcMap = make(map[string]*LoadMiniSvc)

	tickerInt64, tickerInt64Err := i.metaData.toolbox.Cfg.GetInt64(BO, TICKER)
	if nil != tickerInt64Err || tickerInt64 <= 0 {
		log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, TICKER)
	}

	ttlInt64, ttlInt64Err := i.metaData.toolbox.Cfg.GetInt64(BO, TTL)
	if nil != ttlInt64Err || ttlInt64 <= 0 {
		log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, TTL)
	}

	svcString, svcStringErr := i.metaData.toolbox.Cfg.GetString(BO, SVC)
	if nil != svcStringErr {
		log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, SVC)
	}

	monitorInt64, _ := i.metaData.toolbox.Cfg.GetInt64(BO, MONITOR)

	i.metaData.toolbox.Log.Infow(
		"init",
		"ticker", tickerInt64, "ttl", ttlInt64, "svc", svcString)

	i.metaData.svc = svcString
	i.metaData.ttl = ttlInt64 * 1e6 //ms to ns
	i.metaData.ticker = time.NewTicker(time.Millisecond * time.Duration(tickerInt64))
	i.metaData.monitorInterval = monitorInt64

	go i.purge()
	go i.toMonitorWareHouse(&monitorWareHouseInst)

}

func (i *LoadMini) serve(in *xsf.Req, span *xsf.Span, toolbox *xsf.ToolBox) (res *utils.Res, err error) {
	sid := mssSidGenerator.GenerateSid("serve")
	i.metaData.toolbox.Log.Debugw("serve:begin", "sid", sid, "op", in.Op())
	defer i.metaData.toolbox.Log.Debugw("serve:end", "sid", sid, "op", in.Op())

	res = xsf.NewRes()
	switch in.Op() {
	case REPORTER:
		{
			//获取addr
			addrString, addrOk := in.GetParam(LICADDR)
			if !addrOk {
				res.SetError(ErrLbAddrIsIncorrect.errCode, ErrLbAddrIsIncorrect.errInfo)
				return res, nil
			}
			//获取svc
			svcString, svcOk := in.GetParam(LICSVC)
			if !svcOk {
				res.SetError(ErrLbAddrIsIncorrect.errCode, ErrLbAddrIsIncorrect.errInfo)
				return res, nil
			}

			//获取subSvc
			subSvcString, subSvcOk := in.GetParam(LICSUBSVC)
			if !subSvcOk {
				res.SetError(ErrLbAddrIsIncorrect.errCode, ErrLbAddrIsIncorrect.errInfo)
				return res, nil
			}
			//获取total
			totalString, totalOk := in.GetParam(LICTOTAL)
			if !totalOk {
				res.SetError(ErrLbTotalIsIncorrect.errCode, ErrLbTotalIsIncorrect.errInfo)
				return res, nil
			}
			totalInt, totalErr := strconv.Atoi(totalString)
			if nil != totalErr {
				toolbox.Log.Errorf("totalErr:%v", totalErr)
				res.SetError(ErrLbTotalIsIncorrect.errCode, ErrLbTotalIsIncorrect.errInfo)
				return res, nil
			}

			//获取idle
			idleString, idleOk := in.GetParam(LICIDLE)
			if !idleOk {
				res.SetError(ErrLbIdleIsIncorrect.errCode, ErrLbIdleIsIncorrect.errInfo)
				return res, nil
			}
			idleInt, idleErr := strconv.Atoi(idleString)
			if nil != idleErr {
				toolbox.Log.Errorf("idleErr:%v", idleErr)
				res.SetError(ErrLbIdleIsIncorrect.errCode, ErrLbIdleIsIncorrect.errInfo)
				return res, nil
			}

			//获取best
			bestString, bestOk := in.GetParam(LICBEST)
			if !bestOk {
				res.SetError(ErrBestIsIncorrect.errCode, ErrBestIsIncorrect.errInfo)
				return res, nil
			}
			bestInt, bestErr := strconv.Atoi(bestString)
			if nil != bestErr {
				toolbox.Log.Errorf("bestErr:%v", bestErr)
				res.SetError(ErrBestIsIncorrect.errCode, ErrBestIsIncorrect.errInfo)
				return res, nil
			}

			//获取live
			liveString, liveOk := in.GetParam(LICLIVE)
			if !liveOk {
				res.SetError(ErrLiveIsIncorrect.errCode, ErrLiveIsIncorrect.errInfo)
				return res, nil
			}
			liveInt, liveErr := strconv.Atoi(liveString)
			if nil != liveErr {
				toolbox.Log.Errorf("liveErr:%v", liveErr)
				res.SetError(ErrLiveIsIncorrect.errCode, ErrLiveIsIncorrect.errInfo)
				return res, nil
			}
			totalInt64, idleInt64, bestInt64 := int64(totalInt), int64(idleInt), int64(bestInt)
			setErr := i.setServer(
				withSetSid(sid),
				withSetLive(liveInt),
				withSetAddr(addrString),
				withSetSvc(svcString),
				withSetSubSvc(subSvcString),
				withSetTotal(totalInt64),
				withSetIdle(idleInt64),
				withSetBest(bestInt64),
			)
			if nil != setErr {
				res.SetError(setErr.ErrorCode(), setErr.ErrInfo())
			}
		}
	case CLIENT:
		{
			//获取svc
			svcString, svcOk := in.GetParam(LICSVC)
			if !svcOk {
				res.SetError(ErrLbAddrIsIncorrect.errCode, ErrLbAddrIsIncorrect.errInfo)
				return res, nil
			}
			//获取subSvc
			subSvcString, subSvcOk := in.GetParam(LICSUBSVC)
			if !subSvcOk {
				res.SetError(ErrLbAddrIsIncorrect.errCode, ErrLbAddrIsIncorrect.errInfo)
				return res, nil
			}
			//获取nbest
			nBestString, nBestOk := in.GetParam(NBESTTAG)
			if !nBestOk {
				res.SetError(ErrLbNbestIsIncorrect.errCode, ErrLbNbestIsIncorrect.errInfo)
				return res, nil
			}

			//获取exparam，选传
			exparamString, _ := in.GetParam(EXPARAM)

			//获取all
			allString, allOk := in.GetParam(ALL)
			all := false
			if allOk {
				if allString == "1" {
					all = true
				}
			}

			//获取uid
			uidString, uidOk := in.GetParam(UID)
			rpString, rpOk := in.GetParam(RP)
			/*
				-1表示uid值没有传
			*/
			var uidInt64 int64 = -1
			if uidOk && rpOk {
				if len(rpString) != 0 {
					uidInt, uidErr := strconv.Atoi(regex.ReplaceAllString(uidString, ""))
					if nil != uidErr {
						i.metaData.toolbox.Log.Errorw(
							"uid is incorrect",
							"uid", uidString)

						/*
							res.SetError(ErrLbUidIsIncorrect.errCode, ErrLbUidIsIncorrect.errInfo)
							return res, nil
						*/
					} else {
						uidInt64 = int64(uidInt)
					}
				}

			}

			nBestInt, nBestErr := strconv.Atoi(nBestString)
			if nil != nBestErr || nBestInt <= 0 {
				toolbox.Log.Errorw(
					"nBestString illegal",
					"nBestInt", nBestInt, "nBestErr:%v", nBestErr)
				res.SetError(ErrLbNbestIsIncorrect.errCode, ErrLbNbestIsIncorrect.errInfo)
				return res, nil
			}

			nBestNodes, nBestNodesErr := i.getServer(
				withGetSid(sid),
				withGetExParam(exparamString),
				withGetUid(uidInt64),
				withGetAll(all),
				withGetNBest(int64(nBestInt)),
				withGetSubSvc(subSvcString),
				withGetSvc(svcString),
			)
			if nil != nBestNodesErr {
				res.SetError(nBestNodesErr.ErrorCode(), nBestNodesErr.ErrInfo())
			}
			for _, node := range nBestNodes {
				data := utils.NewData()
				data.Append([]byte(node))
				res.AppendData(data)
			}
		}
	default:
		{
			toolbox.Log.Errorf("op:%v -> errCode:%v,errInfo:%v",
				in.Op(), ErrLbInputOperation.errCode, ErrLbInputOperation.errInfo)
			res.SetError(ErrLbInputOperation.errCode, ErrLbInputOperation.errInfo)
		}
	}
	return res, nil
}
func (i *LoadMini) setServer(opt ...SetInPutOpt) (err LbErr) {
	optInst := &SetInPut{}
	for _, optFunc := range opt {
		optFunc(optInst)
	}

	subSvcs := strings.Split(optInst.subSvc, ",")

	i.metaData.toolbox.Log.Infow("setServer",
		"sid", optInst.sid, "optInst", optInst.String())

	var errs []LbErr

	for _, subSvc := range subSvcs {
		optInst.subSvc = subSvc
		errs = append(errs, i.setServerUnit(optInst))

	}

	successFlag := true
	errInfo := func() string {
		var rst []string
		for _, setErr := range errs {
			if nil != setErr {
				successFlag = false
				rst = append(rst, setErr.Error())
			}
		}
		return strings.Join(rst, ",")
	}()

	if !successFlag {
		err = NewLbErrImpl(-1, errInfo)
	}
	return err

}

func (i *LoadMini) setServerUnit(optInst *SetInPut) (err LbErr) {
	sid := optInst.sid
	i.metaData.toolbox.Log.Debugw("setServer:begin", "sid", sid)
	defer func() {
		i.metaData.toolbox.Log.Debugw("setServer:end", "sid", sid)
		if nil != err {
			i.metaData.toolbox.Log.Debugw(
				"setServer:attention!!!",
				"sid", sid, "err", err)
		}
	}()

	i.metaData.toolbox.Log.Debugw(
		"setServer",
		"sid", sid, "strategy", "loadMini", "optInst", optInst)

	i.LoadMiniSvcMapRWMu.RLock()
	svcInst, svcInstOk := i.LoadMiniSvcMap[optInst.svc]
	i.LoadMiniSvcMapRWMu.RUnlock()

	if !svcInstOk {
		//todo new svcInst and sync to LoadMiniSvcMap
		i.metaData.toolbox.Log.Infow(
			"new svcInst and sync to LoadMiniSvcMap",
			"sid", optInst.sid, "svc", optInst.svc, "subSvc", optInst.subSvc, "addr", optInst.addr)

		svcInst = &LoadMiniSvc{}
		svcInst.init(optInst.svc, i.metaData)

		i.LoadMiniSvcMapRWMu.Lock()
		if _, svcInstOk := i.LoadMiniSvcMap[optInst.svc]; !svcInstOk {
			i.LoadMiniSvcMap[optInst.svc] = svcInst
		}
		i.LoadMiniSvcMapRWMu.Unlock()
	}

	svcInst.subSvcMapRwMu.RLock()
	subSvcInst, subSvcInstOk := svcInst.subSvcMap[optInst.subSvc]
	svcInst.subSvcMapRwMu.RUnlock()

	if !subSvcInstOk {
		//new subSvcInst and sync to subSvcMap
		i.metaData.toolbox.Log.Infow(
			"new subSvcInst and sync to subSvcMap",
			"sid", optInst.sid, "svc", optInst.svc, "subSvc", optInst.subSvc, "addr", optInst.addr)

		subSvcInst = &LoadMiniSubSvc{}
		subSvcInst.init(optInst.subSvc, i.metaData)

		svcInst.subSvcMapRwMu.Lock()
		if _, subSvcInstOk := svcInst.subSvcMap[optInst.subSvc]; !subSvcInstOk {
			svcInst.subSvcMap[optInst.subSvc] = subSvcInst
		}
		svcInst.subSvcMapRwMu.Unlock()
	}

	//更新授权
	return i.updateStats(sid, optInst, subSvcInst)
}
func (i *LoadMini) updateStats(sid string, optInst *SetInPut, subSvcInst *LoadMiniSubSvc) (err LbErr) {
	if subSvcInst == nil {
		i.metaData.toolbox.Log.Errorw("subSvcInst equal nil",
			"sid", optInst.sid, "svc", optInst.svc, "subSvc", optInst.subSvc, "addr", optInst.addr)
		return ErrInternalIncorrect
	}
	//主动下线
	if optInst.live == 0 {
		i.metaData.toolbox.Log.Warnw("engine will be offline",
			"sid", optInst.sid, "svc", optInst.svc, "subSvc", optInst.subSvc, "addr", optInst.addr)
		subSvcInst.addrRwMu.Lock()
		delete(subSvcInst.addrMap, optInst.addr)
		subSvcInst.addrRwMu.Unlock()
	}

	subSvcInst.addrRwMu.Lock()
	addrInst, addrInstOk := subSvcInst.addrMap[optInst.addr]
	subSvcInst.addrRwMu.Unlock()

	if !addrInstOk {
		i.metaData.toolbox.Log.Infow(
			"updateStats addr not found, new item",
			"sid", optInst.sid, "svc", optInst.svc, "subSvc", optInst.subSvc, "addr", optInst.addr)
		//节点不存在，新建

		addrInst = &SubSvcItem{addr: optInst.addr}

		subSvcInst.addrRwMu.Lock()
		if _, addrInstOk := subSvcInst.addrMap[optInst.addr]; !addrInstOk {
			subSvcInst.addrMap[optInst.addr] = addrInst
		}
		subSvcInst.addrRwMu.Unlock()
	}

	/*
		live
		0 stand for offline
		1 stand for online
	*/
	live2dead := func() (dead int64) {
		if optInst.live == 1 {
			dead = 0
		} else {
			dead = 1
		}
		return
	}
	i.metaData.toolbox.Log.Infow("updateStats addr update",
		"sid", optInst.sid, "svc", optInst.svc, "subSvc", optInst.subSvc, "addr", optInst.addr)
	addrInst.set(optInst.total, optInst.best, optInst.idle, live2dead())
	return nil
}

func (i *LoadMini) getServer(opt ...GetInPutOpt) (nBestNodes []string, nBestNodesErr LbErr) {
	optInst := &GetInPut{}
	for _, optFunc := range opt {
		optFunc(optInst)
	}

	sid := optInst.sid
	defer func() {
		if 0 == len(nBestNodes) || nil != nBestNodesErr {
			i.metaData.toolbox.Log.Errorw(
				"getServer warning",
				"nBestNodes", nBestNodes, "nBestNodesErr", nBestNodesErr, "sid", sid)
		} else {
			i.metaData.toolbox.Log.Infow(
				"getServer info",
				"nBestNodes", nBestNodes, "nBestNodesErr", nBestNodesErr, "sid", sid)
		}
	}()
	i.metaData.toolbox.Log.Infow(
		"getServer optInst data",
		"sid", sid, "all", optInst.all, "uid", optInst.uid, "svc", optInst.svc,
		"subsvc", optInst.subSvc, "exParam", optInst.exParam, "nBest", optInst.nBest)

	i.LoadMiniSvcMapRWMu.RLock()
	svcInst, svcInstOk := i.LoadMiniSvcMap[optInst.svc]
	i.LoadMiniSvcMapRWMu.RUnlock()

	if !svcInstOk {
		i.metaData.toolbox.Log.Warnw(
			"can't take optInst.svc from loadMiniSvcMap",
			"optInst.svc", optInst.svc, "sid", sid)
		nBestNodesErr = ErrLbSvcIncorrect
		return
	}
	svcInst.subSvcMapRwMu.RLock()
	subSvcInst, subSvcInstOk := svcInst.subSvcMap[optInst.subSvc]
	svcInst.subSvcMapRwMu.RUnlock()

	if !subSvcInstOk {
		i.metaData.toolbox.Log.Warnw("can't take subsvc from suSvcMap",
			"optInst.svc", optInst.svc, "optInst.subSvc", optInst.subSvc, "sid", sid)
		nBestNodesErr = ErrLbSubSvcIncorrect
		return
	}

	return i.dealSubSvcInst(subSvcInst, optInst)
}
func (i *LoadMini) dealSubSvcInst(subSvcInst *LoadMiniSubSvc, optInst *GetInPut) (nBestNodes []string, nBestNodesErr LbErr) {

	subSvcInst.addrRwMu.Lock()
	addrMapLen := len(subSvcInst.addrMap)
	subSvcInst.addrRwMu.Unlock()
	if addrMapLen <= 0 {
		i.metaData.toolbox.Log.Warnw(
			"addrMapLen <= 0",
			"addrMapLen", addrMapLen, "optInst.svc", optInst.svc, "optInst.subSvc", optInst.subSvc, "sid", optInst.sid)
		nBestNodesErr = ErrLbNoSurvivingNode
		return
	}

	var maxAddr string
	var maxIdle int64 = math.MinInt64
	//todo 年后详细设计
	subSvcInst.addrRwMu.Lock()

	for addr, addrInst := range subSvcInst.addrMap {
		i.metaData.toolbox.Log.Infow(
			"take addrInst from addrMap",
			"addr", addr, "addrInst", addrInst.String(), "sid", optInst.sid, "maxIdle", maxIdle, "maxAddr", maxAddr)
		if addrInst.idleInst > maxIdle {
			maxIdle = addrInst.idleInst
			maxAddr = addr
		}
	}

	//预授
	if addrInst, addrInstOk := subSvcInst.addrMap[maxAddr]; addrInstOk {
		i.metaData.toolbox.Log.Infow(
			"pre auth",
			"addrInst", addrInst.String(), "sid", optInst.sid, "maxAddr", maxAddr)
		addrInst.idleInst--
	}

	subSvcInst.addrRwMu.Unlock()

	nBestNodes = append(nBestNodes, maxAddr)
	return
}
