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
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type PollSubSvc struct {
	subSvc string

	ix        uint32                 //索引
	addrRwMu  sync.RWMutex           //保护subSvcMap、subSvcSlice
	addrMap   map[string]*SubSvcItem //K:addr type:*SubSvcItem
	addrSlice []*SubSvcItem

	metaData *PollMeta
}

func (p *PollSubSvc) String() string {
	rst := make(map[string]string)
	p.addrRwMu.RLock()
	defer p.addrRwMu.RUnlock()
	for k, v := range p.addrMap {
		rst[k] = v.String()
	}
	rstT, _ := json.Marshal(rst)
	return string(rstT)
}
func (p *PollSubSvc) init(subSvc string, meta *PollMeta) error {
	p.subSvc = subSvc
	p.addrMap = make(map[string]*SubSvcItem)
	return nil
}

type PollSvc struct {
	svc           string
	subSvcMap     map[string]*PollSubSvc
	subSvcMapRwMu sync.RWMutex

	metaData *PollMeta
}

func (p *PollSvc) init(svc string, meta *PollMeta) error {
	p.svc = svc
	p.subSvcMap = make(map[string]*PollSubSvc)
	return nil
}

type PollMeta struct {
	ttl int64 //纳秒,定时清除无效节点的时间间隔

	ticker *time.Ticker //清除无用节点的扫描周期

	monitorInterval int64 //监控数据的更新时间,单位毫秒，缺省一分钟

	toolbox *xsf.ToolBox
}

type Poll struct {
	metaData   *PollMeta
	svcMap     map[string]*PollSvc
	svcMapRwMu sync.RWMutex
}

func newPoll(toolbox *xsf.ToolBox) *Poll {
	std.Println("strategy:poll => about to init poll")

	PollTmp := &Poll{metaData: &PollMeta{}}
	PollTmp.init(toolbox)

	return PollTmp
}
func (p *Poll) purge() {
	/*
		1、循环取出svc、subsvc
		2、subsvc判断
	*/
	for {
		select {
		case <-p.metaData.ticker.C:
			{
				var svcInstSet []*PollSvc
				p.svcMapRwMu.RLock()
				for _, svcInst := range p.svcMap {
					svcInstSet = append(svcInstSet, svcInst)
				}
				p.svcMapRwMu.RUnlock()

				var subSvcInstSet []*PollSubSvc
				for _, svcInst := range svcInstSet {
					svcInst.subSvcMapRwMu.RLock()
					for _, subSvcInst := range svcInst.subSvcMap {
						subSvcInstSet = append(subSvcInstSet, subSvcInst)
					}
					svcInst.subSvcMapRwMu.RUnlock()
				}

				for _, subSvcInst := range subSvcInstSet {
					p.doPurge(subSvcInst)
				}
			}
		}
	}
}
func (p *Poll) doPurge(subSvcInst *PollSubSvc) {
	var addrSet []string
	subSvcInst.addrRwMu.RLock()

	for addr, addrInst := range subSvcInst.addrMap {
		nowTs := time.Now().UnixNano()
		nodeTtl := nowTs - addrInst.timestamp
		if nodeTtl > p.metaData.ttl {
			addrSet = append(addrSet, addr)
		}
	}
	subSvcInst.addrRwMu.RUnlock()

	subSvcInst.addrRwMu.Lock()

	for _, addr := range addrSet {
		delete(subSvcInst.addrMap, addr)
		p.rmAddrItem(&subSvcInst.addrSlice, addr)
	}

	subSvcInst.addrRwMu.Unlock()

}
func (p *Poll) toMonitorWareHouse(*montorWareHouse) {
	//todo to complete
}

/*
	更新全量数据至监控中
*/
func (p *Poll) traversal() (rst *monitor) {
	//todo to complete
	return
}

func (p *Poll) rmAddrItem(addrInstSet *[]*SubSvcItem, addr string) {
	ix := -1
	for k, v := range *addrInstSet {
		if addr == v.addr {
			ix = k
			break
		}
	}
	if ix == -1 {
		return
	}
	for ; ix < len(*addrInstSet)-1; ix++ {
		(*addrInstSet)[ix] = (*addrInstSet)[ix+1]
	}
	*addrInstSet = (*addrInstSet)[:len(*addrInstSet)-1]
}

func (p *Poll) init(toolbox *xsf.ToolBox) {

	p.metaData.toolbox = toolbox

	p.svcMap = make(map[string]*PollSvc)

	tickerInt64, tickerInt64Err := p.metaData.toolbox.Cfg.GetInt64(BO, TICKER)
	if nil != tickerInt64Err || tickerInt64 <= 0 {
		log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, TICKER)
	}

	ttlInt64, ttlInt64Err := p.metaData.toolbox.Cfg.GetInt64(BO, TTL)
	if nil != ttlInt64Err || ttlInt64 <= 0 {
		log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, TTL)
	}

	monitorInt64, _ := p.metaData.toolbox.Cfg.GetInt64(BO, MONITOR)

	p.metaData.toolbox.Log.Infow("init", "ticker", tickerInt64, "ttl", ttlInt64)

	p.metaData.ttl = ttlInt64
	p.metaData.ticker = time.NewTicker(time.Millisecond * time.Duration(tickerInt64))
	p.metaData.monitorInterval = monitorInt64

	go p.purge()

}

func (p *Poll) serve(in *xsf.Req, span *xsf.Span, toolbox *xsf.ToolBox) (res *utils.Res, err error) {
	sid := mssSidGenerator.GenerateSid("serve")
	p.metaData.toolbox.Log.Debugw("serve:begin",
		"sid", sid, "op", in.Op())
	defer p.metaData.toolbox.Log.Debugw("serve:end",
		"sid", sid, "op", in.Op())

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
			setErr := p.setServer(
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
			exParamString, _ := in.GetParam(EXPARAM)

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
						p.metaData.toolbox.Log.Errorw(
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

			nBestNodes, nBestNodesErr := p.getServer(
				withGetSid(sid),
				withGetExParam(exParamString),
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
func (p *Poll) setServer(opt ...SetInPutOpt) (err LbErr) {
	optInst := &SetInPut{}
	for _, optFunc := range opt {
		optFunc(optInst)
	}

	subSvcs := strings.Split(optInst.subSvc, ",")

	p.metaData.toolbox.Log.Infow("setServer",
		"sid", optInst.sid, "optInst", optInst.String())

	var errs []LbErr

	for _, subSvc := range subSvcs {
		optInst.subSvc = subSvc
		errs = append(errs, p.setServerUnit(optInst))

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

func (p *Poll) setServerUnit(optInst *SetInPut) (err LbErr) {
	sid := optInst.sid
	p.metaData.toolbox.Log.Debugw("setServer:begin", "sid", sid)
	defer func() {
		p.metaData.toolbox.Log.Debugw("setServer:end", "sid", sid)
		if nil != err {
			p.metaData.toolbox.Log.Debugw(
				"setServer:attention!!!",
				"sid", sid, "err", err)
		}
	}()

	p.metaData.toolbox.Log.Debugw(
		"setServer",
		"sid", sid, "strategy", "poll", "optInst", optInst)

	/*
		多svc支持
	*/
	p.svcMapRwMu.RLock()
	svcInst, svcInstOk := p.svcMap[optInst.svc]
	p.svcMapRwMu.RUnlock()

	if !svcInstOk {
		p.metaData.toolbox.Log.Infow(
			"setServerUnit svc not found, new item",
			"sid", optInst.sid, "svc", optInst.svc)
		//不存在svc，新建
		svcInst = &PollSvc{}
		_ = svcInst.init(optInst.svc, p.metaData)

		p.svcMapRwMu.Lock()
		if _, svcInstOk = p.svcMap[optInst.svc]; !svcInstOk {
			p.svcMap[optInst.svc] = svcInst
		}
		p.svcMapRwMu.Unlock()
	}

	svcInst.subSvcMapRwMu.RLock()
	subSvcInst, subSvcInstOk := svcInst.subSvcMap[optInst.subSvc]
	svcInst.subSvcMapRwMu.RUnlock()

	if !subSvcInstOk {
		p.metaData.toolbox.Log.Infow(
			"setServerUnit subSvc not found, new item",
			"sid", optInst.sid, "svc", optInst.svc, "subSvc", optInst.subSvc)
		//不存在subSvc，新建
		subSvcInst = &PollSubSvc{}
		subSvcInst.init(optInst.subSvc, p.metaData)

		svcInst.subSvcMapRwMu.Lock()
		if _, subSvcInstOk = svcInst.subSvcMap[optInst.subSvc]; !subSvcInstOk {
			svcInst.subSvcMap[optInst.subSvc] = subSvcInst
		}
		svcInst.subSvcMapRwMu.Unlock()
	}

	//更新授权
	return p.updateStats(sid, optInst, subSvcInst)
}
func (p *Poll) updateStats(sid string, optInst *SetInPut, subSvcInst *PollSubSvc) (err LbErr) {
	if nil == subSvcInst {
		p.metaData.toolbox.Log.Errorw("subSvcInst equal nil")
		return ErrInternalIncorrect
	}

	subSvcInst.addrRwMu.RLock()
	addrInst, addrInstOk := subSvcInst.addrMap[optInst.addr]
	subSvcInst.addrRwMu.RUnlock()

	if !addrInstOk {
		p.metaData.toolbox.Log.Infow(
			"updateStats addr not found, new item",
			"sid", optInst.sid, "svc", optInst.svc, "subSvc", optInst.subSvc, "addr", optInst.addr)
		//节点不存在，新建

		addrInst = &SubSvcItem{addr: optInst.addr}

		subSvcInst.addrRwMu.Lock()
		if _, addrInstOk := subSvcInst.addrMap[optInst.addr]; !addrInstOk {
			subSvcInst.addrMap[optInst.addr] = addrInst
			subSvcInst.addrSlice = append(subSvcInst.addrSlice, addrInst)
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
	p.metaData.toolbox.Log.Infow("updateStats addr update",
		"sid", optInst.sid, "svc", optInst.svc, "subSvc", optInst.subSvc, "addr", optInst.addr)
	addrInst.set(optInst.total, optInst.best, optInst.idle, live2dead())
	return nil
}
func (p *Poll) getServer(opt ...GetInPutOpt) (nBestNodes []string, nBestNodesErr LbErr) {
	optInst := &GetInPut{}
	for _, optFunc := range opt {
		optFunc(optInst)
	}

	sid := optInst.sid
	p.metaData.toolbox.Log.Debugw(
		"getServer:begin",
		"sid", sid)

	defer func() {
		p.metaData.toolbox.Log.Debugw(
			"getServer:end",
			"sid", sid)
		if nil != nBestNodesErr {
			p.metaData.toolbox.Log.Debugw(
				"getServer:attention!!!",
				"sid", sid, "nBestNodesErr", nBestNodesErr)
		}
	}()

	p.metaData.toolbox.Log.Infow(
		"getServer optInst data",
		"sid", sid, "all", optInst.all, "uid", optInst.uid, "svc", optInst.svc, "subsvc",
		optInst.subSvc, "exParam", optInst.exParam, "nBest", optInst.nBest)

	p.svcMapRwMu.RLock()
	svcInst, svcInstOk := p.svcMap[optInst.svc]
	p.svcMapRwMu.RUnlock()

	if !svcInstOk {
		nBestNodesErr = ErrLbSvcIncorrect
		return
	}

	svcInst.subSvcMapRwMu.RLock()
	subSvcInst, subSvcInstOk := svcInst.subSvcMap[optInst.subSvc]
	svcInst.subSvcMapRwMu.RUnlock()

	if !subSvcInstOk {
		nBestNodesErr = ErrLbSubSvcIncorrect
		return
	}

	var addrInstRst []*SubSvcItem

	subSvcInst.addrRwMu.RLock()
	addrSliceLen := len(subSvcInst.addrSlice)
	subSvcInst.addrRwMu.RUnlock()
	if addrSliceLen <= 0 {
		nBestNodesErr = ErrLbNoSurvivingNode
		return
	}

	subSvcInst.addrRwMu.RLock()
	ix := atomic.AddUint32(&subSvcInst.ix, 1) % uint32(addrSliceLen)
	addrInstRst = p.getNBest(subSvcInst.addrSlice, int(ix), int(optInst.nBest))
	subSvcInst.addrRwMu.RUnlock()

	p.metaData.toolbox.Log.Debugw("take nbest",
		"sid", sid, "subSvcInst", subSvcInst.String())

	for _, addrInst := range addrInstRst {
		nBestNodes = append(nBestNodes, addrInst.addr)
	}

	return
}
func (p *Poll) getNBest(addrInstSet []*SubSvcItem, ix int, nBest int) []*SubSvcItem {
	addrInstSetLen := len(addrInstSet)

	if 0 == addrInstSetLen && nBest <= 0 && ix < 0 {
		return nil
	}

	if addrInstSetLen < nBest {
		return addrInstSet
	}

	var rst []*SubSvcItem
	if ix+nBest > addrInstSetLen {
		rst = addrInstSet[ix:]
		if t := nBest - len(rst); t > 0 {
			rst = append(rst, addrInstSet[:t]...)
		}
	} else {
		rst = addrInstSet[ix : ix+nBest]
	}

	return rst
}
