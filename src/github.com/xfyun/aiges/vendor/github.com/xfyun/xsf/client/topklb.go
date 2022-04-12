package xsf

import (
	"github.com/xfyun/xsf/utils"
	finder "github.com/xfyun/finder-go/common"
	"math/rand"
	"sync"
	"time"
)

type topKLB struct {
	opt               *conOption
	sd                *serviceDiscovery
	probabilityMatrix []int //轮盘赌所用的概率矩阵
	//loadCollector     LoadCollectorInterface //负责调度队列
	loadCollectorMapRwMu sync.RWMutex
	loadCollectorMap     map[string]LoadCollectorInterface //k:svc,v:LoadCollectorInterface
	loadCalculator       *LoadCalculator                   //负责计算load
}

func newTopKLB(o *conOption) *topKLB {
	if !func(o *conOption) bool {
		if 0 == o.timePerSlice || 0 == o.winSize {
			return false
		}
		if 0 == len(o.probabilityMatrix) {
			return false
		}
		return true
	}(o) {
		return nil
	}
	dbgLoggerStd.Printf("newTopKLB ......")
	topKLbInst := new(topKLB)
	//鉴于业务节点数量有限且维护堆的资源消耗过大，故而采用普通轮训的方式
	topKLbInst.opt = o
	topKLbInst.loadCollectorMap = make(map[string]LoadCollectorInterface)
	topKLbInst.loadCalculator = newLoadCalculator(o.timePerSlice, o.winSize)
	topKLbInst.sd = newServiceDiscovery(o.fm)
	topKLbInst.sd.registerCallBackFunc(topKLbInst.updateAddr)
	topKLbInst.probabilityMatrix = o.probabilityMatrix

	go topKLbInst.pingAll()

	return topKLbInst
}
func (t *topKLB) getLoadCollector(svc string) (LoadCollectorInterface, bool) {
	t.loadCollectorMapRwMu.RLock()
	defer t.loadCollectorMapRwMu.RUnlock()
	loadCollector, loadCollectorOk := t.loadCollectorMap[svc]
	return loadCollector, loadCollectorOk
}
func (t *topKLB) setLoadCollector(svc string, item *Item) bool {
	t.loadCollectorMapRwMu.Lock()
	defer t.loadCollectorMapRwMu.Unlock()
	loadCollector, loadCollectorOk := t.loadCollectorMap[svc]
	if !loadCollectorOk {
		t.loadCollectorMap[svc] = newQueue(nil)
		loadCollector = t.loadCollectorMap[svc]
	}
	return loadCollector.store(item)
}
func (t *topKLB) delLoadCollector(svc string, addr string) bool {
	t.loadCollectorMapRwMu.RLock()
	defer t.loadCollectorMapRwMu.RUnlock()
	loadCollector, loadCollectorOk := t.loadCollectorMap[svc]
	if !loadCollectorOk {
		return false
	}
	return loadCollector.delete(loadCollector.getItem(addr))
}
func (t *topKLB) pingAll() {

	if utils.IsNil(t.opt.client) {
		//todo 用日志记录
		panic("init error")
	}

	//心跳
	dbgLoggerStd.Printf("fn:pingAll pingInterval:%v\n", t.opt.pingInterval)
	ticker := time.NewTicker(t.opt.pingInterval)
	for {
		select {
		case <-ticker.C:
			{
				dbgLoggerStd.Printf("fn:pingAll,allSvcKV:%+v\n", t.findAll())
				for svc, svcAddr := range t.findAll() {
					dbgLoggerStd.Printf("fn:pingAll,ping svc:%v,svcAddr:%+v\n", svc, svcAddr)
					t.ping(svc, svcAddr)
				}
			}
		}
	}
}
func (t *topKLB) ping(svc string, targets []string) {
	caller := NewCaller(t.opt.client)
	var baseTime time.Time
	for _, target := range targets {
		baseTime = time.Now()
		s, errcode, e := caller.CallWithAddr(
			"",
			PING,
			target,
			NewReq(),
			PINGTM,
		)
		dur := time.Since(baseTime).Nanoseconds()
		t.opt.client.updateLb(svc, target, s, errcode, e, dur)
	}

}
func (t *topKLB) findAll() map[string][]string {
	var rst = make(map[string][]string)
	for svc, svcAddrs := range t.sd.findAllService() {
		_, addrs := svcAddrs.addrs.NextInList(0)
		rst[svc] = addrs
	}
	return rst
}
func (t *topKLB) Find(param *LBParams) ([]string, []string, error) {
	//会话收敛
	if 0 != len(param.peerIp) {
		fuck.Println("peerIp not equal to zero")
		load, _ := t.getData(param.svc, param.peerIp)
		dbgLoggerStd.Printf("fn:Find,peerIp:%v(threshold:%v,load:%v) is not empty\n", param.peerIp, t.opt.threshold, load)
		if load < int64(t.opt.threshold) && t.opt.threshold != 0 {
			dbgLoggerStd.Printf("fn:Find,peerIp:%v(threshold:%v,load:%v) is ready\n", param.peerIp, t.opt.threshold, load)
			return []string{param.peerIp}, []string{param.peerIp}, nil
		} else {
			dbgLoggerStd.Printf("fn:Find,peerIp:%v(threshold:%v,load:%v) overflow or ignore threshold\n", param.peerIp, t.opt.threshold, load)
		}
	}

	_, e := t.sd.findAll(param.version, param.svc, param.logId, param.log)
	//s, e := t.sd.findAll(param.version, param.svc, param.logId, param.log)
	if e != nil {
		return nil, nil, e
	}

	{
		//nbesAddr, allAddr := s.addrs.NextInList(0)
		//dbgLoggerStd.recF("fn:Find,test 0,nbestAddr:%+v,allAddr:%+v\n", nbesAddr, allAddr)
	}

	addrs := make([]string, 0)
	loadCollectorInst, loadCollectorInstOk := t.getLoadCollector(param.svc)
	if !loadCollectorInstOk {
		return nil, nil, EINVALIDLADDR
	}
	dbgLoggerStd.Printf("fn:Find,nbest:%v\n", param.nbest)
	addrs = append(addrs, loadCollectorInst.load(param.nbest)...)
	if 0 == len(addrs) {
		param.log.Warnw("addrs is empty")
	}
	dbgLoggerStd.recLn("raw addrs:", addrs)

	//基于概率对最佳节点进行波动
	ix := t.roulette()
	dbgLoggerStd.recLn("roulette ix:", ix)
	if len(addrs) > ix {
		addrs[0], addrs[ix] = addrs[ix], addrs[0]
	}
	dbgLoggerStd.recLn("roulette addrs:", addrs)
	return addrs, nil, e
}

//probability采用百分制，如80%记为80
func (t *topKLB) roulette() int {
	randIx := rand.Intn(100)
	dbgLoggerStd.Printf("fn:roulette,randIx:%v\n", randIx)
	var pointer int
	for i := 0; i < len(t.probabilityMatrix); i++ {
		pointer += t.probabilityMatrix[i]
		if randIx <= pointer {
			return i
		}
	}
	return -1
}
func (t *topKLB) updateAddr(svc string, notifyType string, instance []*finder.ServiceInstance) {
	dbgLoggerStd.recF("notifyType:%+v,instance:%+v\n", notifyType, instance)
	if notifyType == string(finder.INSTANCEADDED) {
		for _, addrUnit := range instance {
			t.setLoadCollector(svc, newItem(addrUnit.Addr, 0))
		}
	} else if notifyType == string(finder.INSTANCEREMOVE) {
		for _, addrUnit := range instance {
			t.delLoadCollector(svc, addrUnit.Addr)
		}
	} else {
	}
}

func (t *topKLB) updateData(svc string, target string, s *Res, errcode int32, e error, dur int64, vCpu int64) {
	load := t.loadCalculator.syncWithLoad(
		target,
		cellCalculatorCell{
			errCode: int64(errcode),
			dur:     dur / int64(time.Millisecond),
			vCpu:    vCpu,
		},
	)
	dbgLoggerStd.Printf("fn:updateData,load:%v\n", load)
	loadCollector, loadCollectorOk := t.getLoadCollector(svc)
	if !loadCollectorOk {
		return
	}

	loadCollector.update(nil, target, int(load))
}
func (t *topKLB) getData(svc string, target string) (int64, bool) {
	dbgLoggerStd.Printf("fn:updateData,svc:%v,target:%v\n", svc, target)
	loadCollector, loadCollectorOk := t.getLoadCollector(svc)
	if !loadCollectorOk {
		return 0, false
	}
	return loadCollector.getItem(target).priority, true
}
