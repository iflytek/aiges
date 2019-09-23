package daemon

import (
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const ComponentName = "Server"

const (
	NBESTTAG  = "nbest"
	ALL       = "all"
	REPORTER  = "setServer"
	CLIENT    = "getServer"
	CMDSERVER = "cmdServer"
	UID       = "uid"
	RP        = "rp"
	EXPARAM   = "exparam"
)

const (
	/*-------指定组件的版本号-------*/
	lbVersion    = "0.0.0"
	lbApiVersion = "1.0.0"
	date         = "2019-09-04 16:36:01"
	name         = "hermes"
	author       = "sqjian@iflytek.com"
	repo         = "https://git.xfyun.cn/AIaaS/hermes.git"
	comments     = "external"
)

const (
	SEED      = 10000 //节点的分段基数
	SVC       = "svc"
	DEFSUB    = "defSub"
	THRESHOLD = "threshold"
	PREAUTH   = "preauth"
	TICKER    = "ticker"
	TTL       = "ttl"
	STRATEGY  = "strategy"
	BO        = "bo"
	//RMQABLE   = "rmqable"
	//RMQADDRS  = "rmqaddrs"
	//RMQTOPIC  = "rmqtopic"
	//RMQGROUP  = "rmqgroup"
	//RMQTICKER = "rmqticker"
	//CONSUMER  = "consumer"
	PPROF   = "pprof"
	MONITOR = "monitor"
	NODEDUR = "nodedur"

	defNODEDUR = time.Minute * 5
	//defaultRMQABLE  = 1
	defaultPREAUTH = -1
	//defaultCONSUMER = 1
)

//const (
//	DC = "dc"
//)
//const (
//	DB             = "db"
//	DBABLE         = "able"
//	DBTIME         = "dbtime"
//	RCTIME         = "rctime"
//	DBBASEURL      = "baseurl"
//	DBCALLER       = "caller"
//	DBCALLERKEY    = "callerkey"
//	DBTIMEOUT      = "timeout"
//	DBTOKEN        = "token"
//	DBVERSION      = "version"
//	DBIDC          = "idc"
//	DBSCHEMA       = "schema"
//	DBTABLE        = "table"
//	DBBATCH        = "batch"
//	defaultDBABLE  = 1
//	defaultDBBATCH = 10000
//	defaultDBTIME  = time.Hour * 48
//	defaultRCTIME  = time.Minute * 10
//)

//worker
const (
	LICLIVE   = "live"
	LICADDR   = "addr"
	LICSVC    = "svc"
	LICSUBSVC = "subsvc"
	LICTOTAL  = "total"
	LICIDLE   = "idle"
	LICBEST   = "best"
)

const (
	//some vars about rmq
	RmqUid        = "uid"
	rmqSvcName    = "svc_name"
	rmqSubSvcName = "mc.arm.ent"
	notifyOp      = "AIHotword"
	tmpNode       = "_TmpNode" //标识临时节点
)

const (
	//lbv2c
	cacheService = true
	cacheConfig  = true
	cachePath    = "./findercache"
)

/*
	移除字母，保留数字
*/
var regex = regexp.MustCompile(`[A-Za-z]`)

/*
	单独存储通知引擎失败用
*/
var personalizedLogger *utils.Logger

func init() {
	var err error
	personalizedLogger, err = utils.NewLocalLog(
		utils.SetCaller(false),
		utils.SetLevel("debug"),
		utils.SetFileName("log/personalized.log"),
		utils.SetMaxSize(100),
		utils.SetMaxBackups(10),
		utils.SetMaxAge(10),
		utils.SetAsync(false),
		utils.SetCacheMaxCount(-1),
		utils.SetBatchSize(1024))
	if nil != err {
		panic(err)
	}
}

var mssSidGenerator MssSidGenerator

type StrategyClassify int

const (
	load StrategyClassify = iota
	poll
	loadMini
)
const Unknown = "unknown"

func (s StrategyClassify) String() string {
	switch s {
	case load:
		return "load"
	case poll:
		return "poll"
	case loadMini:
		return "loadMini"
	default:
		return Unknown
	}
}

//const (
//	rcCnt = 10 //重试数据库的次数
//)

/////////////////////////////////////////////////////
var (
	nodeDurWin *abnormalNodeWindow
)

func setAbnormalNode(ts, subSvc, node string) {
	if nil != nodeDurWin {
		nodeDurWin.setAbnormalNode(
			ts + nodeWindowsBoundary + subSvc + nodeWindowsBoundary + node)
	}
}

type abnormalNode struct {
	Ts     string `json:"ts"`
	SubSvc string `json:"sub_svc"`
	Node   string `json:"node"`
}

func getAbnormalNodeStats() []abnormalNode {
	if nil == nodeDurWin {
		return nil
	}
	var abnormalNodes []abnormalNode
	nodes := nodeDurWin.getStats()
	for _, node := range nodes {
		tmp := strings.Split(node, nodeWindowsBoundary)
		if len(tmp) < 2 {
			return nil
		}
		abnormalNodes = append(abnormalNodes, abnormalNode{Ts: tmp[0], SubSvc: tmp[1], Node: tmp[2]})
	}
	return abnormalNodes
}

var (
	blacklist sync.Map
)

const (
	/*
		1、黑名单，用于外部主动下线某个节点
		2、forceOffline参数为0表示强制引擎下线，为1表示解除强制下线
		3、下线后添加至此名单，并不再接受上报
	*/

	FORCEOFFLINE      = "forceOffline"
	ForceOffline      = "1"
	CleanForceOffline = "0"
)

func blackListContent() []string {
	var rst []string
	blacklist.Range(func(key, value interface{}) bool {
		rst = append(rst, key.(string))
		return true
	})
	return rst
}

var Std = std
var std = log.New(os.Stderr, "", log.LstdFlags)
