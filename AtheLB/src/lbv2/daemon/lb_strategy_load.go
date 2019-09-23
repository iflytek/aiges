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
	"fmt"
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

type LoadSubSvc struct {
	//dbFlag bool //如果true，则表示需要进行db等操作
	subSvc string
	//ttl         int64                  //纳秒,定时清除无效节点的时间间隔,暂时没用到
	subSvcRwMu  sync.RWMutex           //保护subSvcMap、subSvcSlice
	subSvcMap   map[string]*SubSvcItem //K:addr
	subSvcSlice SubSvcItemSlice

	/*
		1、通过此字典查询seg_id对应的addr
		2、通过addr查询具体的节点信息
		3、对节点信息进行处理，然后决定是否返回此节点
		4、如此循环，一直到获得nBest熟练的节点后返回或报错
	*/
	//seed          int64             //uid分段的种子，如10000
	//segIdAddr     map[string]string //k:seg_id,v:server_ip  用于维护seg_id与server_ip的对应关系
	//segIdAddrRwMu sync.RWMutex

	//threshold int64 //阈值（百分数，如20代表20%），一旦服务节点超过这个值则不再为该节点分配请求
}

func (L *LoadSubSvc) String() string {
	subSvcStr := func() string {
		subSvcMapStr := func() string {
			var rst = make(map[string]string)

			L.subSvcRwMu.RLock()
			for k, v := range L.subSvcMap {
				rst[k] = fmt.Sprintf("%+v", v)
			}
			L.subSvcRwMu.RUnlock()

			return fmt.Sprintf("SubSvcMap:%+v", rst)
		}()
		subSvcSliceStr := func() string {
			var rst []string

			L.subSvcRwMu.RLock()
			for _, v := range L.subSvcSlice {
				rst = append(rst, func() string {
					return fmt.Sprintf("%+v", v)
				}())
			}
			L.subSvcRwMu.RUnlock()

			return fmt.Sprintf("subSvcSlice:%+v", rst)
		}()
		return strings.Join([]string{subSvcMapStr, subSvcSliceStr}, ",")
	}()

	//L.segIdAddrRwMu.RLock()
	//segIdAddrStr := func() string {
	//	return fmt.Sprintf("segIdAddr:%+v", L.segIdAddr)
	//}()
	//L.segIdAddrRwMu.RUnlock()

	//return strings.Join([]string{subSvcStr, segIdAddrStr}, ",")
	return strings.Join([]string{subSvcStr}, ",")
}

type LoadMeta struct {
	//svc    string //大业务名，如svc
	//defSub string //兜底路由
	ttl int64 //纳秒,定时清除无效节点的时间间隔

	ticker *time.Ticker //清除无用节点的扫描周期

	//dbTime    time.Duration
	deadNodes sync.Map //存储无效节点

	//rcTime time.Duration //第一次访问数据库失败后，后续重新访问的时间间隔

	threshold int64
	//rmqTopic  string
	//rmqGroup    string
	//rmqInterval time.Duration //存储rmq时间间隔
	//rmqAble     int64
	//consumer    int64

	monitorInterval int64 //监控数据的更新时间,单位毫秒，缺省一分钟
	preAuth         int64

	nodeDur time.Duration //异常节点的保存时间，缺省5min

	toolbox *xsf.ToolBox
}

type LoadSvc struct {
	svc        string                 //大业务名，如svc
	defSub     string                 //兜底路由
	ttl        int64                  //纳秒,定时清除无效节点的时间间隔
	svcMap     map[string]*LoadSubSvc //K:type（如sms、sms-5s）
	svcMapRwMu sync.RWMutex

	ticker *time.Ticker //清除无用节点的扫描周期

	threshold   int64
	rmqTopic    string
	rmqGroup    string
	rmqInterval time.Duration //消费rmq的时间间隔
	rmqAble     int64

	monitorInterval int64 //监控数据的更新时间,单位毫秒，缺省一分钟
	preAuth         int64

	toolbox *xsf.ToolBox
}

func (l *LoadSvc) Init(in *LoadMeta) {
	//l.svc = in.svc
	//l.defSub = in.defSub
	l.ttl = in.ttl
	l.svcMap = make(map[string]*LoadSubSvc)

	l.ticker = in.ticker

	l.threshold = in.threshold
	//l.rmqTopic = in.rmqTopic
	//l.rmqGroup = in.rmqGroup
	//l.rmqInterval = in.rmqInterval
	//l.rmqAble = in.rmqAble

	l.monitorInterval = in.monitorInterval
	l.preAuth = in.preAuth

	l.toolbox = in.toolbox
}

type nodeSnap struct {
	addr       string
	addrDetail *SubSvcItem //addrDetail
	addrSegIds []string    //存储节点对应的segIds {seg1,seg2...}
}

func (n *nodeSnap) setAddr(in string) {
	n.addr = in
}
func (n *nodeSnap) setAddrDetail(in *SubSvcItem) {
	n.addrDetail = in
}
func (n *nodeSnap) setAddrSegIds(segId string) {
	n.addrSegIds = append(n.addrSegIds, segId)
}

type Load struct {
	metaData       *LoadMeta
	LoadSvcMap     map[string]*LoadSvc
	LoadSvcMapRWMu sync.RWMutex
}

func (l *Load) map2string(m map[string]string) (rst string) {
	rst = fmt.Sprintf("%+v", m)
	return
}

/*
	1、分析rmq消息包
	2、rmq消息为逗号分隔的长字符串
*/
func (l *Load) parseRmqMsg(sid atomic.Value, in string) (rst map[string]string) {
	rst = make(map[string]string)

	in = strings.Replace(in, "common.", "", -1)
	l.metaData.toolbox.Log.Debugw(
		"remove string common.",
		"sid", sid.Load(), "rmqMsg", in)
	cutByComma := strings.Split(in, ",")
	for _, v := range cutByComma {
		tmp := strings.Split(v, "=")
		if len(tmp) > 1 {
			rst[tmp[0]] = tmp[1]
		}
	}
	l.metaData.toolbox.Log.Debugw(
		"parseRmqMsg.",
		"sid", sid.Load(), "rmqMsg", in, "rst", rst)
	return
}

/*
ats通知函数，通知ats缓存个性化数据
*/


/*
	dbFlag标志判断是否进行db等操作
*/
func (l *Load) doPurge(in *LoadSubSvc, logId string) bool {

	/*
		1、临时取出来，减少锁的作用时间
		2、读取时
			a、先处理subSvcMapEx，取出无效节点的addr
		3、删除时
			a、先处理segIdAddr，否则会读取到脏数据，删除
	*/

	l.metaData.toolbox.Log.Debugw("GarbageCollection=> enter doPurge",
		"logid", logId)
	var addrList []string

	/*
		标记过期的addr节点
	*/
	l.metaData.toolbox.Log.Debugw("GarbageCollection=> mark expired data",
		"logid", logId)
	emptySubSvc := true
	in.subSvcRwMu.RLock()
	for k, v := range in.subSvcMap {
		emptySubSvc = false
		nowTs := time.Now().UnixNano()
		nodeTtl := nowTs - v.timestamp
		if nodeTtl > l.metaData.ttl {
			/*
				1、选择超时的地址
				2、取出所有映射到该地址的段号
			*/
			addrList = append(addrList, k)
			l.metaData.toolbox.Log.Debugw(
				"the node deadline has been reached",
				"nodeTtl", nodeTtl, "nowTs", nowTs, "v.timestamp", v.timestamp, "l.ttl", l.metaData.ttl, "logid", logId)
		}
	}
	in.subSvcRwMu.RUnlock()

	l.metaData.toolbox.Log.Debugw("GarbageCollection=> invalid addr",
		"addrList", addrList, "logid", logId)

	if len(addrList) == 0 {
		l.metaData.toolbox.Log.Debugw("GarbageCollection=> invalid addr len==0, return",
			"addrList", addrList, "logid", logId)
		return emptySubSvc
	}

	/*
		此处需要保证segIdAddr与subSvcMap同时删除，故将segIdAddrRwMu放入内部
	*/
	in.subSvcRwMu.Lock()
	//in.segIdAddrRwMu.Lock()

	var segIds []string

	//删除addr列表
	for _, v := range addrList {
		nodeSnapTmp := &nodeSnap{}
		//if in.dbFlag {
		//	nodeSnapTmp.setAddr(v)
		//	//删除segId列表
		//	for segId, segIdAddr := range in.segIdAddr {
		//		if v == segIdAddr {
		//			nodeSnapTmp.setAddrSegIds(segId)
		//			delete(in.segIdAddr, segId)
		//			segIds = append(segIds, segId)
		//		}
		//	}
		//} else {
		l.metaData.toolbox.Log.Debugw("GarbageCollection=> ignore segId operation",
			"subSvc", in.subSvc, "toDbAddr", v, "logid", logId)
		//}
		if addrDetail, addrDetailOk := in.subSvcMap[v]; addrDetailOk {
			nodeSnapTmp.setAddrDetail(addrDetail)
		} else {
			l.metaData.toolbox.Log.Errorw("fn:doPurge can't take addrDetail from subSvcMap",
				"addr", v, "logid", logId)
		}
		delete(in.subSvcMap, v)
		l.rmSubSvcItemSliceItem(&in.subSvcSlice, v)
		setAbnormalNode(
			time.Now().Format("2006-01-02-15:04:05"),
			in.subSvc,
			v,
		)

		l.metaData.toolbox.Log.Debugw("GarbageCollection=> DelServerAsync",
			"toDbAddr", v, "logid", logId)
		//if in.dbFlag {
		//	l.metaData.toolbox.Log.Debugw("GarbageCollection=> deadNodes.Store",
		//		"delAddr", v, "logid", logId)
		//	l.metaData.deadNodes.Store(v, nodeSnapTmp)
		//}
	}
	//in.segIdAddrRwMu.Unlock()
	in.subSvcRwMu.Unlock()

	l.metaData.toolbox.Log.Debugw("GarbageCollection=> invalid segIds",
		"segIds", segIds, "logid", logId)

	return emptySubSvc
}

//清除函数，定时清除无效节点
func (l *Load) purge() {
	var purgeCnt int64
	for {
		select {
		case <-l.metaData.ticker.C:
			{
				//维护subsvc与svc之间的映射关系
				logId := "purge@" + strconv.Itoa(int(atomic.AddInt64(&purgeCnt, 1))) + strconv.Itoa(time.Now().Nanosecond())
				l.metaData.toolbox.Log.Debugw("GarbageCollection=> begin to start collecting garbage", "logid", logId)

				/*
					※ 垃圾回收
						1、取出所有svc,减少锁粒度
						2、对于每个svc取出所有subsvc
						3、据svc+subsvc取出节点详细信息，剔除无效节点
				*/

				subSvcDataSet, subSvc2Svc := l.extractSubSvc(logId)


				l.dealSubSvc(logId, subSvcDataSet, subSvc2Svc)
			}
		}
	}
}

func (l *Load) dealSubSvc(logId string, subSvcDataSet []*LoadSubSvc, subSvc2Svc map[string]string) {
	l.metaData.toolbox.Log.Debugw("GarbageCollection=> about to deal subSvcDataSet",
		"ts", time.Now(), "logid", logId)
	for _, subSvcData := range subSvcDataSet {
		l.metaData.toolbox.Log.Debugw("GarbageCollection=> about to doPurge subSvcData",
			"logid", logId)

		/*
			emptySubSvc用做移除操作，确保空的subsvc能清理到
		*/
		emptySubSvc := l.doPurge(subSvcData, logId)
		if emptySubSvc {
			l.metaData.toolbox.Log.Errorw("GarbageCollection=> subSvc is empty,about to delete",
				"subSvc", subSvcData.subSvc, "emptySubSvc", emptySubSvc)

			l.LoadSvcMapRWMu.RLock()
			svcDataMark, svcDataMarkOk := l.LoadSvcMap[subSvc2Svc[subSvcData.subSvc]]
			l.LoadSvcMapRWMu.RUnlock()

			if !svcDataMarkOk {
				l.metaData.toolbox.Log.Errorw("GarbageCollection=> can't find svcData from loadSvcMap",
					"subSvc", subSvcData.subSvc, "emptySubSvc", emptySubSvc, "svcDataMarkOk", svcDataMarkOk)
				continue
			}

			svcDataMark.svcMapRwMu.Lock()
			delete(svcDataMark.svcMap, subSvcData.subSvc)
			svcDataMark.svcMapRwMu.Unlock()

		}
	}
}

func (l *Load) extractSubSvc(logId string) ([]*LoadSubSvc, map[string]string) {
	l.metaData.toolbox.Log.Debugw("GarbageCollection=> about to collect svcDataSet", "logid", logId)

	var subSvc2Svc = make(map[string]string)
	var subSvcDataSet []*LoadSubSvc
	var svcDataSet []*LoadSvc
	var svcs []string

	l.LoadSvcMapRWMu.RLock()
	for _, svcData := range l.LoadSvcMap {
		svcs = append(svcs, svcData.svc)
		svcDataSet = append(svcDataSet, svcData)
	}
	l.LoadSvcMapRWMu.RUnlock()

	l.metaData.toolbox.Log.Debugw("GarbageCollection=> svcSet", "svcs", svcs, "logid", logId)
	var svcSubSvcs = make(map[string][]string)
	for _, svcData := range svcDataSet {
		svcData.svcMapRwMu.RLock()
		for subSvc, subSvcData := range svcData.svcMap {
			svcSubSvcs[svcData.svc] = append(svcSubSvcs[svcData.svc], subSvc)
			subSvcDataSet = append(subSvcDataSet, subSvcData)

			//供后续subsvc清除所用
			subSvc2Svc[subSvc] = svcData.svc
		}
		svcData.svcMapRwMu.RUnlock()
	}

	svcSubSvcsByte, _ := json.Marshal(svcSubSvcs)
	l.metaData.toolbox.Log.Debugw("GarbageCollection=> svcSet",
		"svcSubSvcs", string(svcSubSvcsByte), "logid", logId)
	return subSvcDataSet, subSvc2Svc
}
func (l *Load) recoveryNode(sid string, optInst *SetInPut, SubSvcTmp *LoadSubSvc, node *nodeSnap) (err LbErr) {
	l.metaData.toolbox.Log.Debugw(
		"node recovery",
		"sid", sid, "strategy", "load", "fn", "setServer")

	tmp := &SubSvcItem{
		timestamp: time.Now().UnixNano(),
		addr:      node.addrDetail.addr,
		bestInst:  node.addrDetail.bestInst,
		idleInst:  node.addrDetail.idleInst,
		totalInst: node.addrDetail.totalInst}

	l.metaData.toolbox.Log.Debugw("tmp node detail", "sid", sid, "tmp", tmp.String())

	SubSvcTmp.subSvcRwMu.Lock()
	/*
		保证subSvcSlice不会重复append数据
	*/
	if _, ok := SubSvcTmp.subSvcMap[optInst.addr]; !ok {
		l.metaData.toolbox.Log.Debugw("can't take addrDetail,this is new node", "sid", sid, "addr", optInst.addr)
		SubSvcTmp.subSvcMap[optInst.addr] = tmp
		SubSvcTmp.subSvcSlice = append(SubSvcTmp.subSvcSlice, tmp)
	} else {
		l.metaData.toolbox.Log.Debugw("this is not a new node", "sid", sid)
	}
	SubSvcTmp.subSvcRwMu.Unlock()

	/*
		恢复segId映射关系
	*/
	//SubSvcTmp.segIdAddrRwMu.Lock()
	//for _, segId := range node.addrSegIds {
	//	SubSvcTmp.segIdAddr[segId] = node.addr
	//}
	//
	//SubSvcTmp.segIdAddrRwMu.Unlock()
	return
}
func (l *Load) newNode(sid string, optInst *SetInPut, SubSvcTmp *LoadSubSvc) (err LbErr) {

	l.metaData.toolbox.Log.Debugw(
		"node first report",
		"sid", sid, "strategy", "load", "fn", "setServer")
	/*
		1、节点第一次上报
		2、在subSvcMap中新增加addr key，并为此节点选取合适的segId
		3、segId为从零递增的自然数序列（从SegIdManager实例中获取）
	*/

	tmp := &SubSvcItem{
		timestamp: time.Now().UnixNano(),
		addr:      optInst.addr,
		bestInst:  optInst.best,
		idleInst:  optInst.idle,
		totalInst: optInst.total}

	l.metaData.toolbox.Log.Debugw("tmp node detail", "sid", sid, "tmp", tmp.String())

	SubSvcTmp.subSvcRwMu.Lock()
	/*
		保证subSvcSlice不会重复append数据
	*/
	if _, ok := SubSvcTmp.subSvcMap[optInst.addr]; !ok {
		l.metaData.toolbox.Log.Debugw("can't take addrDetail,this is new node",
			"sid", sid, "addr", optInst.addr)
		SubSvcTmp.subSvcMap[optInst.addr] = tmp
		SubSvcTmp.subSvcSlice = append(SubSvcTmp.subSvcSlice, tmp)
	} else {
		l.metaData.toolbox.Log.Debugw("this is not a new node", "sid", sid)
	}
	SubSvcTmp.subSvcRwMu.Unlock()

	/*
		目前仅svc支持数据库分段等操作
	*/
	//if l.metaData.svc == optInst.svc {
	//	segIdStr := strconv.Itoa(int(segIdManagerInst.getMin()))
	//
	//	/*
	//		1、判断segId是否已存在
	//		2、如不存在则添加segId并写入数据库
	//	*/
	//	SubSvcTmp.segIdAddrRwMu.Lock()
	//	if _, ok := SubSvcTmp.segIdAddr[segIdStr]; !ok {
	//
	//		l.metaData.toolbox.Log.Debugw("ready add data to mysql", "sid", sid, "segIdStr", segIdStr)
	//		SubSvcTmp.segIdAddr[segIdStr] = optInst.addr
	//		row := RowData{segIdDb: segIdStr, typeDb: optInst.subSvc, serverIpDb: optInst.addr}
	//		MysqlManagerInst.AddNewSegIdDataAsync(row)
	//
	//	} else {
	//		l.metaData.toolbox.Log.Errorw("didn't add data to mysql", "sid", sid, "segIdStr", segIdStr)
	//	}
	//	SubSvcTmp.segIdAddrRwMu.Unlock()
	//}

	return
}

func (l *Load) newSubSvc(sid string, optInst *SetInPut, svcInst *LoadSvc, SubSvcTmp *LoadSubSvc) (err LbErr) {

	/*
		1、生成节点实例
		2、创建节点map
		3、创建subsvc map
	*/

	l.metaData.toolbox.Log.Debugw(
		"subsvc first report",
		"sid", sid, "strategy", "load", "fn", "setServer")
	/*
		1、节点第一次上报
		2、在subSvcMap中新增加addr key，并为此节点选取合适的segId
		3、segId为从零递增的自然数序列（从SegIdManager实例中获取）
	*/

	tmp := &SubSvcItem{
		timestamp: time.Now().UnixNano(),
		addr:      optInst.addr,
		bestInst:  optInst.best,
		idleInst:  optInst.idle,
		totalInst: optInst.total}

	subSvcMap := map[string]*SubSvcItem{optInst.addr: tmp}
	subSvcSlice := append(make(SubSvcItemSlice, 0, 1000), tmp)

	//dbFlag := false
	//var segIdStr string
	//if l.metaData.svc == optInst.svc {
	//	dbFlag = true
	//	segIdStr = strconv.Itoa(int(segIdManagerInst.getMin()))
	//}
	//segIdAddr := map[string]string{segIdStr: optInst.addr}

	svcInst.svcMapRwMu.Lock()
	svcInst.svcMap[optInst.subSvc] = &LoadSubSvc{
		//dbFlag:      dbFlag,
		subSvc: optInst.subSvc,
		//segIdAddr:   segIdAddr,
		//seed:        SEED,
		//threshold:   svcInst.threshold,
		subSvcSlice: subSvcSlice,
		subSvcMap:   subSvcMap}
	svcInst.svcMapRwMu.Unlock()

	//if l.metaData.svc == optInst.svc {
	//	row := RowData{segIdDb: segIdStr, typeDb: optInst.subSvc, serverIpDb: optInst.addr}
	//	MysqlManagerInst.AddNewSegIdDataAsync(row)
	//}

	return
}
func (l *Load) updateSubSvc(sid string, optInst *SetInPut, svcInst *LoadSvc, SubSvcTmp *LoadSubSvc) (err LbErr) {
	/*
		能取到表示节点曾经上传过
	*/
	SubSvcTmp.subSvcRwMu.RLock()
	item, itemOk := SubSvcTmp.subSvcMap[optInst.addr]
	SubSvcTmp.subSvcRwMu.RUnlock()
	/*
		live
		0 stand for offline
		1 stand for online
	*/
	live2dead := func() (dead int64) {
		if 1 == optInst.live {
			dead = 0
		} else {
			dead = 1
		}
		return
	}
	if !itemOk {
		deadNode, deadNodeOk := l.metaData.deadNodes.Load(optInst.addr)
		if deadNodeOk {
			l.metaData.deadNodes.Delete(optInst.addr)

			deadNodeInst, deadNodeInstOk := deadNode.(*nodeSnap)

			if deadNodeInstOk {
				_ = l.recoveryNode(sid, optInst, SubSvcTmp, deadNodeInst)
			} else {
				l.metaData.toolbox.Log.Errorw("fn:updateSubSvc can't convert deadNode to *nodeSnap",
					"sid", sid, "strategy", "addr", optInst.addr)
			}
		} else {
			err = l.newNode(sid, optInst, SubSvcTmp)
		}

	} else {
		l.metaData.toolbox.Log.Debugw(
			"node update",
			"sid", sid, "strategy", "load", "fn", "setServer")
		/*
			1、节点授权更新
			2、通过key addr索引到相关的节点数据，进行递增或递减操作
		*/

		//节点主动下线
		if 0 == optInst.live {
			l.metaData.toolbox.Log.Debugw(
				"node active offline",
				"sid", sid, "strategy", "load", "fn", "setServer")

			SubSvcTmp.subSvcRwMu.Lock()
			//SubSvcTmp.segIdAddrRwMu.Lock()

			//if l.metaData.svc == optInst.svc {
			//	//删除segId列表
			//	for segId, segIdAddr := range SubSvcTmp.segIdAddr {
			//		if optInst.addr == segIdAddr {
			//			delete(SubSvcTmp.segIdAddr, segId)
			//
			//			/*
			//				释放sedId
			//			*/
			//			i, e := strconv.ParseInt(segId, 10, 64)
			//			if e == nil {
			//				segIdManagerInst.free(i)
			//			}
			//		}
			//	}
			//}

			delete(SubSvcTmp.subSvcMap, optInst.addr)
			l.rmSubSvcItemSliceItem(&SubSvcTmp.subSvcSlice, optInst.addr)
			//SubSvcTmp.segIdAddrRwMu.Unlock()
			SubSvcTmp.subSvcRwMu.Unlock()

			//if l.metaData.svc == optInst.svc {
			//MysqlManagerInst.DelServer(optInst.addr)
			//}
		}

		item.set(optInst.total, optInst.best, optInst.idle, live2dead())
	}
	return
}

/*
	兼容多个subsvc上报,subsvc以逗号切割
*/
func (l *Load) setServer(opt ...SetInPutOpt) (err LbErr) {
	optInst := &SetInPut{}
	for _, optFunc := range opt {
		optFunc(optInst)
	}
	//不黑名单受限，不再接受上报
	if _, ok := blacklist.Load(optInst.addr); ok {
		return ErrBlacklistIsIncorrect
	}
	subSvcs := strings.Split(optInst.subSvc, ",")

	var errs []LbErr

	for _, subSvc := range subSvcs {
		optInst.subSvc = subSvc
		errs = append(errs, l.setServerUnit(optInst))

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
func (l *Load) setServerUnit(optInst *SetInPut) (err LbErr) {

	sid := optInst.sid
	l.metaData.toolbox.Log.Debugw("setServer:begin", "sid", sid)
	defer func() {
		l.metaData.toolbox.Log.Debugw("setServer:end", "sid", sid)
		if nil != err {
			l.metaData.toolbox.Log.Debugw(
				"setServer:attention!!!",
				"sid", sid, "err", err)
		}
	}()

	l.metaData.toolbox.Log.Debugw(
		"setServer",
		"sid", sid, "strategy", "load", "optInst", optInst)

	/*
		多svc支持
	*/
	l.LoadSvcMapRWMu.RLock()
	svcInst, svcInstOk := l.LoadSvcMap[optInst.svc]
	l.LoadSvcMapRWMu.RUnlock()

	if !svcInstOk {
		svcInst = &LoadSvc{}
		svcInst.Init(l.metaData)
		//svcInst = &LoadSvc{}.Init(l.metaData)
	}

	svcInst.svcMapRwMu.RLock()
	SubSvcTmp, SubSvcTmpOk := svcInst.svcMap[optInst.subSvc]
	svcInst.svcMapRwMu.RUnlock()

	if svcInstOk && SubSvcTmpOk {

		err = l.updateSubSvc(sid, optInst, svcInst, SubSvcTmp)

	} else {

		err = l.newSubSvc(sid, optInst, svcInst, SubSvcTmp)

		l.LoadSvcMapRWMu.Lock()
		if nil == l.LoadSvcMap {
			l.LoadSvcMap = make(map[string]*LoadSvc)
		}
		l.LoadSvcMap[optInst.svc] = svcInst
		l.LoadSvcMapRWMu.Unlock()
	}

	return
}

/*
	1、从段表里取出地址处理
*/
//func (l *Load) oldSegment(sid string, addr string, optInst *GetInPut, svcInst *LoadSvc, SubSvcTmp *LoadSubSvc) (nBestNodes []string, nBestNodesErr LbErr) {
//
//	/*
//		1、从段里索引出地址
//		2、通过addr取出节点的详细数据
//	*/
//
//	SubSvcTmp.subSvcRwMu.RLock()
//	addrDetails, addrDetailsOk := SubSvcTmp.subSvcMap[addr]
//	SubSvcTmp.subSvcRwMu.RUnlock()
//	if !addrDetailsOk {
//		//当分段表第一次从数据库中拉取时，没有相关的负载信息
//
//		nBestNodesErr = ErrLbNoSurvivingNode
//		l.metaData.toolbox.Log.Errorw(
//			"getServer can't take addr from subSvcMapEx,maybe logic error",
//			"sid", sid, "addrDetailsOk", addrDetailsOk, "SubSvcTmp", SubSvcTmp, "addr", addr)
//		return
//	}
//	//如果是个性化通知的请求，不用判断负载
//	if optInst.personalized {
//		l.metaData.toolbox.Log.Infow(
//			"getServer",
//			"sid", sid, "addrDetailsOk", addrDetailsOk, "SubSvcTmp", SubSvcTmp, "personalized", "personal request")
//		nBestNodes = append(nBestNodes, addrDetails.addr)
//		return
//	}
//	l.metaData.toolbox.Log.Infow(
//		"getServer",
//		"sid", sid, "addrDetailsOk", addrDetailsOk, "SubSvcTmp", SubSvcTmp)
//
//	if addrDetails.totalInst == 0 {
//		l.metaData.toolbox.Log.Errorw(
//			"getServer warnning",
//			"sid", sid, "addrDetails.totalInst", addrDetails.totalInst, "SubSvcTmp", SubSvcTmp)
//
//		nBestNodes = append(nBestNodes, addrDetails.addr)
//		return
//	}
//	l.metaData.toolbox.Log.Debugw(
//		"judge whether or not be overload",
//		"sid", sid, )
//
//	if ((addrDetails.totalInst - addrDetails.idleInst) * 100 / addrDetails.totalInst) < SubSvcTmp.threshold {
//		l.metaData.toolbox.Log.Debugw(
//			"not overload",
//			"sid", sid, "idleInst", addrDetails.idleInst, "totalInst", addrDetails.totalInst, "threshold", SubSvcTmp.threshold)
//
//		//此处仅判断当前节点是否过载，没过载则返回当前节点
//		nBestNodes = append(nBestNodes, addrDetails.addr)
//	} else {
//		l.metaData.toolbox.Log.Errorw(
//			"already overload",
//			"sid", sid, "idleInst", addrDetails.idleInst, "totalInst", addrDetails.totalInst, "threshold", SubSvcTmp.threshold)
//
//		addrMin, addrMinErr := l.takeAddrMin(svcInst, SubSvcTmp)
//		if nil != addrMinErr {
//			nBestNodesErr = ErrLbNoSurvivingNode
//			l.metaData.toolbox.Log.Errorw(
//				"takeAddrMin error,maybe totalInst error",
//				"sid", sid, )
//			return
//		}
//		svcInst.toolbox.Log.Debugw(
//			"takeAddrMin",
//			"addrMin", addrMin, "addrMinErr", addrMinErr, "SubSvcTmp", SubSvcTmp)
//		nBestNodes = append(nBestNodes, addrMin+tmpNode)
//	}
//	return
//}
//func (l *Load) segIs2int(in string) (rst int) {
//	rst, _ = strconv.Atoi(in)
//	return
//}
//func (l *Load) newSegment(sid string, segId string, optInst *GetInPut, svcInst *LoadSvc, SubSvcTmp *LoadSubSvc) (
//nBestNodes []string, nBestNodesErr LbErr) {
//
///*
//	1、段里取不到地址，表明这是一个新段
//	2、分配新段，并同步到数据库
//*/
//
//addrMin, addrMinErr := l.takeAddrMin(svcInst, SubSvcTmp)
//l.metaData.toolbox.Log.Debugw(
//	"takeAddrMin",
//	"sid", sid, "addrMin", addrMin, "addrMinErr", addrMinErr, "SubSvcTmp", SubSvcTmp)
//
//if nil != addrMinErr {
//	nBestNodesErr = ErrLbNoSurvivingNode
//	l.metaData.toolbox.Log.Errorw(
//		"takeAddrMin error,maybe totalInst error",
//		"sid", sid, "SubSvcTmp", SubSvcTmp)
//	return
//}
//
//nBestNodes = append(nBestNodes, addrMin)
//
//segIdManagerInst.freeze(int64(l.segIs2int(segId)))
//
////添加segId并写入数据库
//SubSvcTmp.segIdAddrRwMu.Lock()
//SubSvcTmp.segIdAddr[segId] = addrMin
//SubSvcTmp.segIdAddrRwMu.Unlock()
//
//l.metaData.toolbox.Log.Debugw("newSegment", "segId", segId, "addrMin", addrMin, "sid", sid)
//
//row := RowData{segIdDb: segId, typeDb: optInst.subSvc, serverIpDb: addrMin}
//MysqlManagerInst.AddNewSegIdDataAsync(row)
//return
//}

/*
	1、个性化处理
*/
//func (l *Load) personalized(sid string, optInst *GetInPut, svcInst *LoadSvc, SubSvcTmp *LoadSubSvc) (nBestNodes []string, nBestNodesErr LbErr) {
//
//	l.metaData.toolbox.Log.Infow(
//		"success take SubSvcTmp from l.svcMapEx",
//		"sid", sid, "optInst.subSvc", optInst.subSvc, "personalized", true)
//
//	//将uid转换为segId
//	segIdStr := strconv.FormatInt(optInst.uid%SubSvcTmp.seed, 10)
//
//	//从segId字典中取出addr
//	SubSvcTmp.segIdAddrRwMu.RLock()
//
//	addrTmp, addrTmpOk := SubSvcTmp.segIdAddr[segIdStr]
//	l.metaData.toolbox.Log.Debugw(
//		"getServer:sth about segId",
//		"sid", sid, "optInst.uid", optInst.uid, "SubSvcTmp.seed", SubSvcTmp.seed, "segIdStr", segIdStr, "SubSvcTmp.segIdAddr",
//		SubSvcTmp.segIdAddr, "addrTmpOk", addrTmpOk, "addrTmp", addrTmp)
//
//	SubSvcTmp.segIdAddrRwMu.RUnlock()
//
//	if !addrTmpOk {
//		l.metaData.toolbox.Log.Debugw("about to call newSegment", "sid", sid, "segIdStr", segIdStr)
//		nBestNodes, nBestNodesErr = l.newSegment(sid, segIdStr, optInst, svcInst, SubSvcTmp)
//
//	} else {
//		l.metaData.toolbox.Log.Debugw("about to call oldSegment", "sid", sid, "segIdStr", segIdStr, "addrTmp", addrTmp)
//		nBestNodes, nBestNodesErr = l.oldSegment(sid, addrTmp, optInst, svcInst, SubSvcTmp)
//
//	}
//
//	/*
//		填充节点数量，第一个为按uid规则选取的节点，后面的nbest-1个节点按load规则选取，可能有重复
//	*/
//	l.metaData.toolbox.Log.Debugw(
//		"personalized complete",
//		"sid", sid, "nBestNodesLen", len(nBestNodes), "nBestNodes", nBestNodes, "nBestNodesErr", nBestNodesErr)
//
//	if nBestNodesErr == nil && (optInst.nBest-1) > 0 {
//
//		l.metaData.toolbox.Log.Debugw(
//			"ready to reAllocate additional nodes",
//			"sid", sid, "nbest", optInst.nBest)
//
//		/*
//			填充nbest-1
//		*/
//		optInst.nBest--
//		nBestNodesAdditional, nBestNodesAdditionalErr := l.notPersonalized(sid, optInst, svcInst, SubSvcTmp, true)
//		if nil != nBestNodesAdditionalErr {
//
//			l.metaData.toolbox.Log.Debugw(
//				"reAllocate additional nodes failed",
//				"sid", sid, "nBestNodesAdditionalErr", nBestNodesAdditionalErr)
//			nBestNodesErr = nBestNodesAdditionalErr
//		} else {
//
//			l.metaData.toolbox.Log.Debugw(
//				"reAllocate additional nodes success",
//				"sid", sid, "nBestNodesAdditional", nBestNodesAdditional, "nBestNodesAdditionalErr", nBestNodesAdditionalErr)
//			for _, node := range nBestNodesAdditional {
//				nBestNodes = append(nBestNodes, node)
//			}
//		}
//	}
//	return
//}

/*
	1、2018.9.28 新增preAuthorized
	1.1、如果preAuthorized为true，则表示已经授权，否则需对第一个addr授权

*/
func (l *Load) notPersonalized(sid string, optInst *GetInPut, svcInst *LoadSvc, SubSvcTmp *LoadSubSvc, preAuthorized bool) (nBestNodes []string, nBestNodesErr LbErr) {

	l.metaData.toolbox.Log.Infow(
		"success take SubSvcTmp from l.svcMapEx",
		"sid", sid, "optInst.subSvc", optInst.subSvc, "personalized", false)
	/*
		如果传了all参数，则其它参数会被忽略
	*/
	if optInst.all {
		all := func() string {
			var rst []string
			SubSvcTmp.subSvcRwMu.RLock()
			for _, v := range SubSvcTmp.subSvcSlice {
				m := make(map[string]interface{}, 10)
				rst = append(rst, func() string {
					m["addr"] = v.addr
					m["idleInst"] = v.idleInst
					m["bestInst"] = v.bestInst
					m["bestInst"] = v.bestInst
					r, _ := json.Marshal(m)
					return string(r)
				}())
			}
			SubSvcTmp.subSvcRwMu.RUnlock()
			return strings.Join(rst, ";")
		}()
		nBestNodes = append(nBestNodes, all)
	} else {
		SubSvcTmp.subSvcRwMu.RLock()
		nBestNodes = l.takeNBest(sid, SubSvcTmp.subSvcSlice, optInst.nBest, preAuthorized)
		SubSvcTmp.subSvcRwMu.RUnlock()

		if 0 == len(nBestNodes) {
			nBestNodesErr = ErrLbNoSurvivingNode
		}
		if nil != nBestNodesErr {
			//单纯的日志代码
			var tmpStrSlice []string

			SubSvcTmp.subSvcRwMu.RLock()
			for _, v := range SubSvcTmp.subSvcSlice {
				tmpStrSlice = append(tmpStrSlice, v.String())
			}
			SubSvcTmp.subSvcRwMu.RUnlock()

			l.metaData.toolbox.Log.Debugw("sth about load", "snapshot", strings.Join(tmpStrSlice, ";"))
		}
	}
	return
}
func (l *Load) cmdServer(addr, forceOffline string) (status string, err LbErr) {
	switch forceOffline {
	case ForceOffline:
		{
			blacklist.Store(addr, struct{}{})
			status = fmt.Sprintf("successfully adding %v to blacklist(%v)", addr, blackListContent())
		}
	case CleanForceOffline:
		{
			blacklist.Delete(addr)
			status = fmt.Sprintf("successfully deleting %v from blacklist(%v)", addr, blackListContent())
		}
	default:
		{
			err = ErrCmdServerIsIncorrect
		}
	}
	return
}
func (l *Load) getServer(opt ...GetInPutOpt) (nBestNodes []string, nBestNodesErr LbErr) {
	optInst := &GetInPut{}
	for _, optFunc := range opt {
		optFunc(optInst)
	}
	sid := optInst.sid
	l.metaData.toolbox.Log.Debugw(
		"getServer:begin",
		"sid", sid)

	defer func() {
		l.metaData.toolbox.Log.Debugw(
			"getServer:end",
			"sid", sid)
		if nil != nBestNodesErr {
			l.metaData.toolbox.Log.Debugw(
				"getServer:attention!!!",
				"sid", sid, "nBestNodesErr", nBestNodesErr)
		}
	}()

	l.metaData.toolbox.Log.Infow(
		"getServer optInst data",
		"sid", sid, "all", optInst.all, "uid", optInst.uid, "svc", optInst.svc, "subsvc", optInst.subSvc, "exParam", optInst.exParam, "nBest", optInst.nBest)

	/*
		多svc支持
	*/
	l.LoadSvcMapRWMu.RLock()
	svcInst, svcInstOk := l.LoadSvcMap[optInst.svc]
	l.LoadSvcMapRWMu.RUnlock()

	if !svcInstOk {
		nBestNodesErr = ErrLbSvcIncorrect
		return
	}

	l.metaData.toolbox.Log.Debugw(
		"take subsvc from svcMap",
		"sid", sid, "subSvc", optInst.subSvc)

	svcInst.svcMapRwMu.RLock()
	SubSvcTmp, SubSvcTmpOk := svcInst.svcMap[optInst.subSvc]
	svcInst.svcMapRwMu.RUnlock()

	if !SubSvcTmpOk {
		l.metaData.toolbox.Log.Debugw(
			"can't subsvc from svcMap,about to take defSub",
			"sid", sid, "subSvc", optInst.subSvc, "defSub", svcInst.defSub)

		svcInst.svcMapRwMu.RLock()
		SubSvcTmp, SubSvcTmpOk = svcInst.svcMap[svcInst.defSub]
		svcInst.svcMapRwMu.RUnlock()

		if !SubSvcTmpOk {
			l.metaData.toolbox.Log.Debugw(
				"no node reported in defSub",
				"sid", sid, "defSub", svcInst.defSub)
		}
	}

	if SubSvcTmpOk {
		nBestNodes, nBestNodesErr = l.notPersonalized(sid, optInst, svcInst, SubSvcTmp, false)
		/*
			1、判断是否走个性化
			2、目前判断的唯一标准为是否传uid，uid默认-1
		*/
		//switch func() bool {
		//	if optInst.uid == -1 {
		//		return false
		//	}
		//	if l.metaData.svc != optInst.svc {
		//		return false
		//	}
		//	return true
		//}() {
		//case true:
		//	{
		//		nBestNodes, nBestNodesErr = l.personalized(sid, optInst, svcInst, SubSvcTmp)
		//
		//		l.metaData.toolbox.Log.Debugw(
		//			"personalized complete",
		//			"sid", sid, "nBestNodes", nBestNodes, "nBestNodesErr", nBestNodesErr)
		//
		//	}
		//case false:
		//	{
		//		nBestNodes, nBestNodesErr = l.notPersonalized(sid, optInst, svcInst, SubSvcTmp, false)
		//
		//		l.metaData.toolbox.Log.Debugw(
		//			"notPersonalized complete",
		//			"sid", sid, "nBestNodes", nBestNodes, "nBestNodesErr", nBestNodesErr)
		//
		//	}
	} else {
		nBestNodesErr = ErrLbSubSvcIncorrect
	}
	return
}

/*
	更新全量数据至监控中
*/
func (l *Load) traversal() *monitorSvc {

	//拷贝svc映射
	var loadSvcMap = make(map[string]*LoadSvc)
	l.LoadSvcMapRWMu.RLock()
	for svc, loadSvc := range l.LoadSvcMap {
		loadSvcMap[svc] = loadSvc
	}
	l.LoadSvcMapRWMu.RUnlock()

	var monitorLoadTmp monitorSvc

	//使用svc映射
	for svc, loadSvc := range loadSvcMap {
		var monitorLoadSubSvcTmp monitorSubSvc

		//拷贝subsvc映射
		var loadSubSvcMap = make(map[string]*LoadSubSvc)
		loadSvc.svcMapRwMu.RLock()
		for subSvc, loadSubSvc := range loadSvc.svcMap {
			loadSubSvcMap[subSvc] = loadSubSvc
		}
		loadSvc.svcMapRwMu.RUnlock()

		//使用subsvc映射
		for subSvc, loadSubSvc := range loadSubSvcMap {
			l.metaData.toolbox.Log.Infow("fn:traversal", "subsvc", subSvc)

			var monitorLoadAddrTmp monitorAddr

			//拷贝subsvcfunc映射
			var subSvcItemMap = make(map[string]*SubSvcItem)
			loadSubSvc.subSvcRwMu.RLock()
			for loadSubSvcAddr, loadSubSvcAuth := range loadSubSvc.subSvcMap {
				subSvcItemMap[loadSubSvcAddr] = loadSubSvcAuth
			}
			loadSubSvc.subSvcRwMu.RUnlock()

			//使用subsvcfunc映射
			for loadSubSvcAddr, loadSubSvcAuth := range subSvcItemMap {
				var monitorSubSvcItemTmp monitorSubSvcItem

				monitorSubSvcItemTmp.Addr = loadSubSvcAuth.addr
				monitorSubSvcItemTmp.Timestamp = loadSubSvcAuth.timestamp
				monitorSubSvcItemTmp.BestInst = loadSubSvcAuth.bestInst
				monitorSubSvcItemTmp.IdleInst = func() int64 {
					if loadSubSvcAuth.idleInst < 0 {
						return loadSubSvcAuth.idleInst
					}
					return loadSubSvcAuth.idleInst
				}()
				monitorSubSvcItemTmp.TotalInst = loadSubSvcAuth.totalInst

				if nil == monitorLoadAddrTmp.AddrMap {
					monitorLoadAddrTmp.AddrMap = make(map[string]*monitorSubSvcItem)
				}
				monitorLoadAddrTmp.AddrMap[loadSubSvcAddr] = &monitorSubSvcItemTmp
			}

			if nil == monitorLoadSubSvcTmp.SubSvcMap {
				monitorLoadSubSvcTmp.SubSvcMap = make(map[string]*monitorAddr)
			}
			monitorLoadSubSvcTmp.SubSvcMap[subSvc] = &monitorLoadAddrTmp

		}

		if nil == monitorLoadTmp.SvcMap {
			monitorLoadTmp.SvcMap = make(map[string]*monitorSubSvc)
		}
		monitorLoadTmp.SvcMap[svc] = &monitorLoadSubSvcTmp
	}

	return &monitorLoadTmp
}
func (l *Load) init(toolbox *xsf.ToolBox) {

	l.metaData.toolbox = toolbox

	std.Println("strategy:load => about to read config")
	////读取配置文件
	//dbAble, dbAbleErr := l.metaData.toolbox.Cfg.GetInt64(DB, DBABLE)
	//if nil != dbAbleErr {
	//	dbAble = defaultDBABLE
	//}
	//dbTime := defaultDBTIME
	//dbTimeInt64, dbTimeInt64Err := l.metaData.toolbox.Cfg.GetInt64(DB, DBTIME)
	//if nil == dbTimeInt64Err {
	//	dbTime = time.Second * time.Duration(dbTimeInt64)
	//}
	//rcTime := defaultRCTIME
	//reTimeInt, reTimeIntErr := l.metaData.toolbox.Cfg.GetInt64(DB, RCTIME)
	//if nil == reTimeIntErr {
	//	rcTime = time.Second * time.Duration(reTimeInt)
	//}
	//baseUrlString, baseUrlErr := l.metaData.toolbox.Cfg.GetString(DB, DBBASEURL)
	//if nil != baseUrlErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", DB, DBBASEURL)
	//}
	//callerString, callerErr := l.metaData.toolbox.Cfg.GetString(DB, DBCALLER)
	//if nil != callerErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", DB, DBCALLER)
	//}
	//callerKeyString, callerKeyErr := l.metaData.toolbox.Cfg.GetString(DB, DBCALLERKEY)
	//if nil != callerKeyErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", DB, DBCALLERKEY)
	//}
	//timeoutInt, timeoutErr := l.metaData.toolbox.Cfg.GetInt(DB, DBTIMEOUT)
	//if nil != timeoutErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", DB, DBTIMEOUT)
	//}
	//tokenString, tokenErr := l.metaData.toolbox.Cfg.GetString(DB, DBTOKEN)
	//if nil != tokenErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", DB, DBTOKEN)
	//}
	//versionString, versionErr := l.metaData.toolbox.Cfg.GetString(DB, DBVERSION)
	//if nil != versionErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", DB, DBVERSION)
	//}
	//idcString, idcErr := l.metaData.toolbox.Cfg.GetString(DB, DBIDC)
	//if nil != idcErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", DB, DBIDC)
	//}
	//schemaString, schemaErr := l.metaData.toolbox.Cfg.GetString(DB, DBSCHEMA)
	//if nil != schemaErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", DB, DBSCHEMA)
	//}
	//batchInt, batchIntErr := l.metaData.toolbox.Cfg.GetInt(DB, DBBATCH)
	//if nil != batchIntErr {
	//	batchInt = defaultDBBATCH
	//}
	//tableString, tableErr := l.metaData.toolbox.Cfg.GetString(DB, DBTABLE)
	//if nil != tableErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", DB, DBTABLE)
	//}
	//svcString, svcStringErr := l.metaData.toolbox.Cfg.GetString(BO, SVC)
	//if nil != svcStringErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, SVC)
	//}
	//defSubString, defSubStringErr := l.metaData.toolbox.Cfg.GetString(BO, DEFSUB)
	//if nil != defSubStringErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, DEFSUB)
	//}
	preAuthInt64, preAuthInt64Err := l.metaData.toolbox.Cfg.GetInt64(BO, PREAUTH)
	if nil != preAuthInt64Err {
		preAuthInt64 = defaultPREAUTH
	}
	if preAuthInt64 > 0 {
		log.Fatalf("preauth:%v > 0", preAuthInt64)
	}
	if preAuthInt64 == 0 {
		std.Printf("warnning:preauth:%v\n", preAuthInt64)
	}
	ttlInt64, ttlInt64Err := l.metaData.toolbox.Cfg.GetInt64(BO, TTL)
	if nil != ttlInt64Err {
		log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, TTL)
	}
	tickerInt64, tickerInt64Err := l.metaData.toolbox.Cfg.GetInt64(BO, TICKER)
	if nil != tickerInt64Err {
		log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, TICKER)
	}

	thresholdInt64, thresholdInt64Err := l.metaData.toolbox.Cfg.GetInt64(BO, THRESHOLD)
	if nil != thresholdInt64Err {
		log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, THRESHOLD)
	}

	//rmqTopicString, rmqTopicStringErr := l.metaData.toolbox.Cfg.GetString(BO, RMQTOPIC)
	//if nil != rmqTopicStringErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, RMQTOPIC)
	//}
	//
	//rmqAbleInt64, rmqAbleInt64Err := toolbox.Cfg.GetInt64(BO, RMQABLE)
	//if nil != rmqAbleInt64Err {
	//	rmqAbleInt64 = defaultRMQABLE
	//}
	//consumerInt64, consumerInt64Err := toolbox.Cfg.GetInt64(BO, CONSUMER)
	//if nil != consumerInt64Err {
	//	consumerInt64 = defaultCONSUMER
	//}
	//
	//rmqGroupString, rmqGroupStringErr := l.metaData.toolbox.Cfg.GetString(BO, RMQGROUP)
	//if nil != rmqGroupStringErr {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, RMQGROUP)
	//}
	//
	//rmqTickerInt64, rmqTickerInt64Err := l.metaData.toolbox.Cfg.GetInt64(BO, RMQTICKER)
	//if nil != rmqTickerInt64Err {
	//	log.Fatalf("l.toolbox.Cfg.GetString(%v, %v)", BO, RMQTICKER)
	//}

	monitorInt64, _ := l.metaData.toolbox.Cfg.GetInt64(BO, MONITOR)

	nodeDur := defNODEDUR
	nodeDurInt64, _ := l.metaData.toolbox.Cfg.GetInt64(BO, NODEDUR)
	if 0 != nodeDurInt64 {
		nodeDur = time.Duration(nodeDurInt64) * time.Millisecond
	}

	std.Println("strategy:load => about to init mysql")

	//_, _ = MysqlManagerInst.Init(
	//	dbAble,
	//	batchInt,
	//	l.metaData.toolbox.Log,
	//	baseUrlString,
	//	callerString,
	//	callerKeyString,
	//	time.Duration(timeoutInt)*time.Millisecond,
	//	tokenString,
	//	versionString,
	//	idcString,
	//	schemaString,
	//	tableString)
	//
	//l.metaData.dbTime = dbTime
	//l.metaData.rcTime = rcTime
	//l.metaData.svc = svcString
	//l.metaData.defSub = defSubString
	l.metaData.ttl = ttlInt64 * 1e6 //配置文件里单位是毫秒，此处转换为纳秒
	l.metaData.threshold = thresholdInt64
	l.metaData.monitorInterval = monitorInt64
	l.metaData.preAuth = preAuthInt64
	l.metaData.nodeDur = nodeDur
	nodeDurWin = newAbnormalNodeWindow(l.metaData.nodeDur)
	std.Printf("strategy:load => preAuth:%v\n", l.metaData.preAuth)

	l.metaData.ticker = time.NewTicker(time.Millisecond * time.Duration(tickerInt64))

	//if err := dcInst.Init(withDcBc(l.metaData.toolbox.Bc)); nil != err {
	//	log.Fatalf("dcInst.Init fail -> err:%v", err)
	//}

	//l.metaData.rmqTopic = rmqTopicString
	//l.metaData.rmqGroup = rmqGroupString
	//
	//l.metaData.rmqInterval = time.Millisecond * time.Duration(rmqTickerInt64)
	//l.metaData.rmqAble = rmqAbleInt64
	//l.metaData.consumer = consumerInt64

	//LoadSvcInst := LoadSvc{}
	//LoadSvcInst.Init(l.metaData)
	//l.metaData.toolbox.Log.Debugw("fn:init",
	//	"svc", LoadSvcInst.svc)

	std.Println("strategy:load => about to init Internal data structure")

	/*
		data:subsvc--segid--srvip
	*/
	//data, dataErr := MysqlManagerInst.GetSubSvcSegIdSrvipEx()
	//if nil != dataErr {
	//	l.metaData.toolbox.Log.Errorw("pos:MysqlManagerInst.GetSubSvcSegIdSrvip",
	//		"data", data, "dataErr", dataErr)
	//	go l.recoveryDbData()
	//} else {
	//	dataJson, _ := json.Marshal(data)
	//	l.metaData.toolbox.Log.Debugw("mysqlRst", "data", string(dataJson))
	//}
	//
	//for subsvc, segIdAddr := range data {
	//	l.metaData.toolbox.Log.Debugw("load.init", "subsvc", subsvc, "segIdAddrLen", len(segIdAddr), "segIdAddr", segIdAddr)
	//	LoadSvcInst.svcMapRwMu.Lock()
	//
	//	subSvcMap := func() (rst map[string]*SubSvcItem) {
	//		rst = make(map[string]*SubSvcItem)
	//		for segId, serverIp := range segIdAddr {
	//			segIdInt, _ := strconv.Atoi(segId)
	//			segIdManagerInst.freeze(int64(segIdInt))
	//
	//			rst[serverIp] = &SubSvcItem{
	//				timestamp: time.Now().UnixNano(),
	//				addr:      serverIp}
	//		}
	//		return
	//	}()
	//	subSvcSlice := func() (rst SubSvcItemSlice) {
	//		rst = make(SubSvcItemSlice, 0, 1000)
	//		for _, v := range subSvcMap {
	//			rst = append(rst, v)
	//		}
	//		return
	//	}()
	//
	//	LoadSvcInst.svcMap[subsvc] = &LoadSubSvc{
	//		dbFlag:      true,
	//		subSvc:      subsvc,
	//		segIdAddr:   segIdAddr,
	//		seed:        SEED,
	//		threshold:   l.metaData.threshold,
	//		subSvcSlice: subSvcSlice,
	//		subSvcMap:   subSvcMap}
	//
	//	l.metaData.toolbox.Log.Debugw("fn:init",
	//		"subsvc", subsvc, "segIdAddr", segIdAddr, "subSvcSlice", subSvcSlice, "subSvcMap", subSvcMap)
	//
	//	LoadSvcInst.svcMapRwMu.Unlock()
	//}

	l.LoadSvcMapRWMu.Lock()
	if nil == l.LoadSvcMap {
		l.LoadSvcMap = make(map[string]*LoadSvc)
	}
	//l.LoadSvcMap[l.metaData.svc] = &LoadSvcInst
	l.LoadSvcMapRWMu.Unlock()

	l.toMonitorWareHouse(&monitorWareHouseInst)

	go l.purge() //定时清除无效节点
	//go l.DealDeadNodes()
	/*
		是否开启rmq消费
	*/
	/*	std.Printf("strategy:load => rmqAble:%v,fn:init\n", l.metaData.rmqAble)
		if l.metaData.rmqAble != 0 {
			rmqAddrs, rmqAddrsErr := toolbox.Cfg.GetString(BO, RMQADDRS)
			if nil != rmqAddrsErr {
				panic(fmt.Sprintf("can't GetString %v from %v", RMQADDRS, BO))
			}

			for ix := int64(0); ix < l.metaData.consumer; ix++ {
				std.Printf("NO.%d rmq consumer init.\n", ix)

				var RmqManagerInst RmqManager
				if rmqInitErr := RmqManagerInst.Init(strings.Split(rmqAddrs, ",")); nil != rmqInitErr {
					std.Printf("RmqManagerInst init failed,err:%v\n", rmqInitErr)
				}
				go l.notify(RmqManagerInst)
			}
		}*/

	std.Println("strategy:load => init complete")
}

/*
	第一次拉取数据失败时，循环从数据库拉取数据
*/
//func (l *Load) recoveryDbData() {
//	logId := "recoveryDbData@" + strconv.Itoa(time.Now().Nanosecond())
//	//todo to complete
//	l.metaData.toolbox.Log.Infow("enter recoveryDbData",
//		"logId", logId, "rcTime", l.metaData.rcTime)
//	timer := time.NewTimer(l.metaData.rcTime)
//
//	cnt := 0
//end:
//	for {
//		select {
//		case <-timer.C:
//			{
//				if cnt++; cnt <= rcCnt {
//
//					//data:subsvc--segid--srvip
//					data, dataErr := MysqlManagerInst.GetSubSvcSegIdSrvipEx()
//					if nil != dataErr {
//						l.metaData.toolbox.Log.Errorw(
//							"pos:MysqlManagerInst.GetSubSvcSegIdSrvip fail",
//							"logId", logId, "data", data, "dataErr", dataErr, "cnt", cnt)
//					} else {
//						l.metaData.toolbox.Log.Infow(
//							"pos:MysqlManagerInst.GetSubSvcSegIdSrvip success",
//							"logId", logId, "data", data, "dataErr", dataErr, "cnt", cnt)
//						//恢复数据
//						l.db2Internal(data, logId)
//
//						break end
//					}
//				} else {
//					break end
//				}
//
//				timer.Reset(l.metaData.rcTime)
//			}
//		}
//	}
//
//}
//
////subsvc--segid--srvip
//func (l *Load) db2Internal(data map[string]map[string]string, logId string) {
//	l.metaData.toolbox.Log.Infow("enter db2Internal")
//	for subSvc, segIdAddrs := range data {
//
//		l.LoadSvcMapRWMu.RLock()
//		svcInst, svcInstOk := l.LoadSvcMap[l.metaData.svc]
//		l.LoadSvcMapRWMu.RUnlock()
//
//		if svcInstOk {
//			svcInst.svcMapRwMu.RLock()
//			subSvcInst, subSvcInstOk := svcInst.svcMap[subSvc]
//			svcInst.svcMapRwMu.RUnlock()
//
//			if subSvcInstOk {
//				//取到subsvc
//				l.metaData.toolbox.Log.Infow(
//					"fn:db2Internal,take subsvc",
//					"logId", logId, "subsvc", subSvc)
//				for segId, srvIp := range segIdAddrs {
//
//					subSvcInst.segIdAddrRwMu.RLock()
//					_, subSvcSegIdOk := subSvcInst.segIdAddr[segId]
//					subSvcInst.segIdAddrRwMu.RUnlock()
//
//					if subSvcSegIdOk {
//						l.metaData.toolbox.Log.Infow(
//							"fn:db2Internal,segId already used,ignore",
//							"segId", segId, "srvIp", srvIp, "logId", logId, "subsvc", subSvc)
//						//segId已经被使用，忽略数据库中的数据
//						continue
//					} else {
//						//恢复数据
//						l.metaData.toolbox.Log.Infow(
//							"fn:db2Internal,recovery data",
//							"segId", segId, "srvIp", srvIp, "logId", logId, "subsvc", subSvc)
//						subSvcInst.segIdAddrRwMu.Lock()
//						subSvcInst.segIdAddr[segId] = srvIp
//						subSvcInst.segIdAddrRwMu.Unlock()
//					}
//
//				}
//
//			} else {
//				//没有取到subsvc，新建subsvc数据
//				l.metaData.toolbox.Log.Infow(
//					"fn:db2Internal,new subsvc data",
//					"logId", logId, "subsvc", subSvc)
//				subSvcMap := func() (rst map[string]*SubSvcItem) {
//					rst = make(map[string]*SubSvcItem)
//					for segId, serverIp := range segIdAddrs {
//						segIdInt, _ := strconv.Atoi(segId)
//						segIdManagerInst.freeze(int64(segIdInt))
//
//						rst[serverIp] = &SubSvcItem{
//							timestamp: time.Now().UnixNano(),
//							addr:      serverIp}
//					}
//					return
//				}()
//				subSvcSlice := func() (rst SubSvcItemSlice) {
//					rst = make(SubSvcItemSlice, 0, 1000)
//					for _, v := range subSvcMap {
//						rst = append(rst, v)
//					}
//					return
//				}()
//
//				subSvcTmp := &LoadSubSvc{
//					dbFlag:      true,
//					subSvc:      subSvc,
//					segIdAddr:   segIdAddrs,
//					seed:        SEED,
//					threshold:   l.metaData.threshold,
//					subSvcSlice: subSvcSlice,
//					subSvcMap:   subSvcMap}
//
//				svcInst.svcMapRwMu.Lock()
//				svcInst.svcMap[subSvc] = subSvcTmp
//				svcInst.svcMapRwMu.Unlock()
//			}
//
//		} else {
//			//没有取到svc，新建svc数据
//			l.metaData.toolbox.Log.Infow(
//				"fn:db2Internal,new svc data",
//				"logId", logId, "subsvc", subSvc)
//			LoadSvcInst := LoadSvc{}
//			LoadSvcInst.Init(l.metaData)
//			for subsvc, segIdAddr := range data {
//				l.metaData.toolbox.Log.Debugw(
//					"load.init",
//					"logId", logId, "subsvc", subsvc, "segIdAddrLen", len(segIdAddr), "segIdAddr", segIdAddr)
//				LoadSvcInst.svcMapRwMu.Lock()
//
//				subSvcMap := func() (rst map[string]*SubSvcItem) {
//					rst = make(map[string]*SubSvcItem)
//					for segId, serverIp := range segIdAddr {
//						segIdInt, _ := strconv.Atoi(segId)
//						segIdManagerInst.freeze(int64(segIdInt))
//
//						rst[serverIp] = &SubSvcItem{
//							timestamp: time.Now().UnixNano(),
//							addr:      serverIp}
//					}
//					return
//				}()
//				subSvcSlice := func() (rst SubSvcItemSlice) {
//					rst = make(SubSvcItemSlice, 0, 1000)
//					for _, v := range subSvcMap {
//						rst = append(rst, v)
//					}
//					return
//				}()
//
//				LoadSvcInst.svcMap[subsvc] = &LoadSubSvc{
//					dbFlag:      true,
//					subSvc:      subsvc,
//					segIdAddr:   segIdAddr,
//					seed:        SEED,
//					threshold:   l.metaData.threshold,
//					subSvcSlice: subSvcSlice,
//					subSvcMap:   subSvcMap}
//
//				l.metaData.toolbox.Log.Debugw("fn:init",
//					"logId", logId, "subsvc", subsvc, "segIdAddr", segIdAddr, "subSvcSlice", subSvcSlice, "subSvcMap", subSvcMap)
//
//				LoadSvcInst.svcMapRwMu.Unlock()
//			}
//
//			l.LoadSvcMapRWMu.Lock()
//			if nil == l.LoadSvcMap {
//				l.LoadSvcMap = make(map[string]*LoadSvc)
//			}
//			l.LoadSvcMap[l.metaData.svc] = &LoadSvcInst
//			l.LoadSvcMapRWMu.Unlock()
//		}
//
//	}
//}

/*
	拷贝数据至montorWareHouse
*/
func (l *Load) toMonitorWareHouse(in *montorWareHouse) {
	//in.Svc = l.metaData.svc
	//in.Defsub = l.metaData.defSub
	in.Ttl = l.metaData.ttl
	in.Threshold = l.metaData.threshold
	//in.RmqTopic = l.metaData.rmqTopic
	//in.RmqGroup = l.metaData.rmqGroup

	in.toolbox = l.metaData.toolbox

	in.MonitorInterval = func() int64 {
		if 0 == l.metaData.monitorInterval {
			return defaultMonitorInterval.Nanoseconds() / 1e6
		}
		return l.metaData.monitorInterval
	}()

	in.handle = l

	go in.loop()
}
func (l *Load) serve(in *xsf.Req, span *xsf.Span, toolbox *xsf.ToolBox) (res *utils.Res, err error) {
	sid := mssSidGenerator.GenerateSid("serve")
	l.metaData.toolbox.Log.Debugw("serve:begin", "sid", sid, "op", in.Op())
	defer l.metaData.toolbox.Log.Debugw("serve:end", "sid", sid, "op", in.Op())

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
			setErr := l.setServer(
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
			l.metaData.toolbox.Log.Debugw(
				"rec uid & rp", "sid", sid,
				"uid", uidString, "uidOk", uidOk, "rpString", rpString, "rpOk", rpOk)
			/*
				-1表示uid值没有传
			*/
			var uidInt64 int64 = -1
			if uidOk && rpOk {
				if 0 != len(rpString) {
					uidInt, uidErr := strconv.Atoi(regex.ReplaceAllString(uidString, ""))
					if nil != uidErr {
						l.metaData.toolbox.Log.Errorw(
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

			nBestNodes, nBestNodesErr := l.getServer(
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
	case CMDSERVER:
		{
			//获取addr
			addrString, addrOk := in.GetParam(LICADDR)
			if !addrOk {
				res.SetError(ErrLbAddrIsIncorrect.errCode, ErrLbAddrIsIncorrect.errInfo)
				return res, nil
			}
			//获取live
			forceOfflineString, forceOfflineOk := in.GetParam(FORCEOFFLINE)
			if !forceOfflineOk {
				res.SetError(ErrForceOfflineIsIncorrect.errCode, ErrForceOfflineIsIncorrect.errInfo)
				return res, nil
			}
			cmdStatus, cmdErr := l.cmdServer(addrString, forceOfflineString)
			if nil != cmdErr {
				res.SetError(cmdErr.ErrorCode(), cmdErr.ErrInfo())
			}
			res.SetParam("status", cmdStatus)
		}
	default:
		{
			toolbox.Log.Errorf("op:%v -> errCode:%v,errInfo:%v", in.Op(), ErrLbInputOperation.errCode, ErrLbInputOperation.errInfo)
			res.SetError(ErrLbInputOperation.errCode, ErrLbInputOperation.errInfo)
		}
	}
	return res, nil
}
func newLoad(toolbox *xsf.ToolBox) *Load {
	std.Println("strategy:load => about to init load")

	loadTmp := &Load{metaData: &LoadMeta{}}
	loadTmp.init(toolbox)

	return loadTmp
}
func (l *Load) takeNBest(sid string, in SubSvcItemSlice, n int64, preAuthorized bool) (bestNodes []string) {
	/*
		signMap用来标记，那些节点是已经被选过了的
	*/

	var preAuthorizedFlag = false
	zeroCnt := 0
	l.metaData.toolbox.Log.Debugw("takeNBest:enter", "sid", sid, "n", n)

	signMap := make(map[int64]struct{})
	for i := int64(0); i < n; i++ {
		var max, maxIx int64 = math.MinInt64, -1
		for k, v := range in {
			/*
				如果v.totalInst == 0 || v.bestInst == 0表示这里的数据仅仅来自于mysql，节点还没有上报
			*/
			if v.totalInst == 0 || v.bestInst == 0 { //尽量保证能够返回节点，此处不应该return
				l.metaData.toolbox.Log.Debugw(
					"totalInst or bestInst equal zero",
					"sid", sid, "addr", v.addr, "inLen", len(in))
				zeroCnt++
				continue
			}
			if v.idleInst > max && !(func() bool { _, existed := signMap[int64(k)]; return existed })() {
				max = v.idleInst
				maxIx = int64(k)
			}
			l.metaData.toolbox.Log.Debugw("takeNBest", "sid", sid, "v.idleInst", v.idleInst, "max", max, "maxIx", maxIx)
		}
		if maxIx != -1 {
			signMap[maxIx] = struct{}{}
			if (!preAuthorized) && (!preAuthorizedFlag) {
				/*
					如果已经已经预授和本次已处理预授，不再处理
				*/
				in[maxIx].preAuthorization(l.metaData.preAuth) //提前将授权减1
				preAuthorizedFlag = true
			}
			bestNodes = append(bestNodes, in[maxIx].addr)
			l.metaData.toolbox.Log.Debugw("takeNBest",
				"sid", sid, "maxIx", maxIx, "bestNode", in[maxIx].addr)
		} else {
			l.metaData.toolbox.Log.Debugw("takeNBest:can't take enough nodes",
				"sid", sid)
		}
	}

	l.metaData.toolbox.Log.Debugw("takeNBest rstTmp",
		"bestNodes", bestNodes, "sid", sid)

	nTmp := 0
	if zeroCnt == len(in)*int(n) {
		l.metaData.toolbox.Log.Warnw("takeNBest warning, all nodes is zero",
			"sid", sid, "zeroCnt", zeroCnt)
		for _, v := range in {
			nTmp++
			bestNodes = append(bestNodes, v.addr)
		}
	}

	l.metaData.toolbox.Log.Debugw("takeNBest complete",
		"bestNodes", bestNodes, "sid", sid, "zeroCnt", zeroCnt)

	return
}
func (l *Load) rmSubSvcItemSliceItem(in *SubSvcItemSlice, item string) {
	ix := -1
	for k, v := range *in {
		if item == v.addr {
			ix = k
			break
		}
	}
	if ix == -1 {
		return
	}
	for ; ix < len(*in)-1; ix++ {
		(*in)[ix] = (*in)[ix+1]
	}
	*in = (*in)[:len(*in)-1]
}

func (l *Load) takeAddrMin(svcInst *LoadSvc, in *LoadSubSvc) (addrMin string, addrMinErr error) {

	/*
		如果超过阈值则临时转到其它节点中去
		如果segId查不到也临时转移到其它节点
		遍历所有节点寻找负载最小者
	*/
	zeroCnt := 0
	zeroAddr := ""
	loadMin := int64(math.MaxInt64)
	var loadNow int64
	in.subSvcRwMu.RLock()
	for _, v := range in.subSvcMap {
		if 0 == v.totalInst {
			//这种情况仅在第一次拉取分段表的时候存在
			zeroCnt++

			/*
				此处任意保存一个为0的地址，当所有节点都为0时，返回此地址
			*/
			zeroAddr = v.addr
			continue
		}

		loadNow = (v.totalInst - v.idleInst) * 100 / v.totalInst
		if loadNow < loadMin {
			loadMin = loadNow
			addrMin = v.addr
		}
	}

	/*

		if "" == addrMin && zeroCnt != 0 {
		if zeroCnt == len(in.subSvcMap) {
	*/
	in.subSvcRwMu.RUnlock()
	if "" == addrMin && 0 != zeroCnt {
		if zeroCnt == len(in.subSvcMap) {
			/*
				此处表明所有节点totalInst均为0,节点仍未上报
			*/
			addrMin = zeroAddr
		} else {
			addrMinErr = ErrLbNoSurvivingNode
		}
	} else {
		l.metaData.toolbox.Log.Infow(
			"about to call preAuthorization",
			"l.preAuth", l.metaData.preAuth, "addrMin", addrMin)
		if "" == addrMin {
			addrMinErr = ErrLbNoSurvivingNode
		} else {
			if addrDetail, addrDetailOk := in.subSvcMap[addrMin]; addrDetailOk {
				addrDetail.preAuthorization(l.metaData.preAuth)
			}
		}
	}

	return
}

//func (l *Load) missOrder(in []string) {
//	rand.Seed(int64(time.Now().Nanosecond()))
//	randIx := rand.Intn(len(in))
//	in[0], in[randIx] = in[randIx], in[0]
//}
