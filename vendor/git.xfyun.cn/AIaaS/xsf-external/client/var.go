/*
* @file	var.go
* @brief  定义包内使用的变量
*
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */

package xsf

import (
	"errors"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"time"
)

// 配置字段
var (
	CFGSECTION        = "sfc"
	CFGCONTM          = "conn-timeout" //连接超时，单位ms
	CFGDEFAULTTIMEOUT = 1000

	CFGCONLC        = "conn-life-cycle" //连接生命周期，单位毫秒
	CFGCONDEFLC int = 120 * 1e3

	CFGRTY    = "conn-retry" //连接重试次数
	CFGDEFRTY = 3

	CFGCONPOOL    = "conn-pool-size" //连接池最大的个数
	CFGDEFCONPOLL = 3

	CFGCONRBUF          = "conn-rbuf" // 连接读缓冲区
	CFGDEFRBUF          = GRCMRS
	CFGCONWBUF          = "conn-wbuf" // 连接写缓冲区
	CFGDEFWBUF          = GRCMWS
	CFGCONMAXRECEIVE    = "maxreceive"
	DEFCFGCONMAXRECEIVE = 4 * 1024 * 1024

	CFGLBMODE    = "lb-mode" //0禁用lb,2使用lb。缺省0
	CFGDEFLBMODE = 0

	//lb相关配置
	CFGLBRTY    = "lb-retry" //当负载均衡失败时，重试几次
	CFGDEFLBRTY = 1

	CFGLBTIMEOUT    = "lb-timeout"
	CFGDEFLBTIMEOUT = 500

	//用于测试
	CFGADDRS = "taddrs" //连接池最大的个数

	// tracer  配置
	//CFGTRACESECTION = "trace"
	//CFGTRACEADDR    = "host"
	//CFGDEFTRACEADDR = "127.0.0.1"
	//
	//CFGTRACEPORT    = "port"
	//CFGDEFTRACEPORT = "4545"
	//
	//CFGTRACEENABLE    = "able"
	//CFGDEFTRACEENABLE = 1
	//CFGTRACEUNCHANGE  = -1
	//
	//CFGTRACEDUMP    = "dump"
	//CFGDEFTRACEDUMP = 0
	//
	//CFGTRACEBAKEND    = "backend"
	//CFGDEFTRACEBAKEND = 12
	//
	//CFGTRACESPILL     = "spill"
	//CFGDEFTRACESPILL  = "/log/spill"
	//CFGTRACESPILLABLE = "spill-able"
	//
	//CFGTRACEDELIVER    = "deliver"
	//CFGDEFTRACEDELIVER = 1
	//
	//CFGTRACEBUFFER    = "buffer"
	//CFGDEFTRACEBUFFER = 100000
	//
	//CFGTRACEBATCH    = "batch"
	//CFGDEFTRACEBATCH = 100
	//
	//CFGTRACELINGER    = "linger"
	//CFGDEFTRACELINGER = 5

	CFGLOGSECTON    = "log"
	CFGLOGLEVEL     = "level"
	CFGLOGDEFLEVEL  = "warn"
	CFGLOGCALLER    = "caller"
	CFGLOGDEFCALLER = false

	CFGLOGFILE    = "file"
	CFGLODEFGFILE = "./xsfc.log"

	CFGLOGSIZE     = "size"
	CFGLOGDEFSIZSE = 10

	CFGLOGCOUNT    = "count"
	CFGLOGDEFCOUNT = 20

	CFGLOGDIED    = "die"
	CFGLOGDEFDIED = 30

	CFGLOGASYNC    = "async"
	CFGLOGDEFASYNC = 1

	CFGLOGCACHE    = "cache"
	CFGLOGDEFCACHE = -1

	CFGLOGBATCH    = "batch"
	CFGLOGDEFBATCH = 16 * 1024

	CFGLOGWASH    = "wash"
	CFGLOGDEFWASH = 60

	//////////////////////////////
	//CFGTRACEWATCH  = "watch"
	//CFGLOGDEFWATCH = false
	//
	//CFGTRACEWATCHPORT    = "watchport"
	//CFGTRACEDEFWATCHPORT = 12331
	//
	//CFGTRACESPILLSIZE    = "spillsize"
	//CFGTRACEDEFSPILLSIZE = 1
	//
	//CFGTRACELOADTS    = "loadts"
	//CFGTRACEDEFLOADTS = 1
	//
	//CFGTRACEBCLUSTER   = "bcluster"
	//CFGTRACEDEFCLUSTER = "3s"
	//
	//CFGTRACEIDC    = "idc"
	//CFGTRACEDEFIDC = "dz"

	CLIHOST    = "host"
	CLINETCARD = "netcard"

	defaultCacheService = true
	defaultCacheConfig  = true
	defaultCachePath    = "."
	defaultApiVersion   = "1.0.0"

	//2019-07-30 16:13:31 add sth about keepalive
	CFGKEEPALIVE           = "keepalive"
	CFGDEFKEEPALIVE        = time.Duration(0) //启用keepalive check，值为探测的时间间隔，单位毫秒，0表示不启用，缺省不启用keeplive
	CFGKEEPALIVETIMEOUT    = "keepalive-timeout"
	CFGDEFKEEPALIVETIMEOUT = 1000 * time.Millisecond

	//2019-07-23 15:14:53
	//topKLb新增配置
	CFGWINTIMEPERSLICE    = "time-per-slice"
	CFGDEFWINTIMEPERSLICE = 1000 //ms
	CFGWINSIZE            = "win-size"
	CFGDEFWINSIZE         = 100

	CFGPROBABILITYSEPARATOR = ","
	CFGPROBABILITY          = "probability" //逗号分隔
	CFGDEFPROBABILITY       = []int{80, 16, 4}

	CFGTHRESHOLD    = "threshold"
	CFGDEFTHRESHOLD = 0 //表示忽略阈值判断

	CFGPINGINTERVAL    = "ping"
	CFGDEFPINGINTERVAL = time.Second
)

const (
	PING   = "ping"
	PINGTM = time.Second
)

// 错误信息
var (
	ECFGISNIL     = errors.New("cfg is nil")
	ELBMODE       = errors.New("new lb failed. check lb mode")
	NOUSECONN     = errors.New("No useful connetion!")
	EBADADDR      = errors.New("ContinueCall: invalid address")
	EBADHANDLE    = errors.New("ContinueCall: invalid handle")
	EINVALIDADDR  = errors.New("RemoteLB: can't find valued addr")
	INVALIDLB     = errors.New("RemoteLB: can't find valued lb")
	INVALIDRMLB   = errors.New("request remoteLb failed")
	EINVALIDLADDR = errors.New("LocalLB: can't find valued addr")
	INVALIDSRV    = errors.New("Finder: can't find busin service")
	INVALIDFINDER = errors.New("Finder: finder is nil")

	EINAILDOP    = errors.New("Invaild Op Params")
	EINAILDHDL   = errors.New("Invaild Session Handle")
	EINAILIDDATA = errors.New("Invaild Data")

	//INVAILDPARAM  int32 = 10106
	//INVAILDHANDLE int32 = 10109
	//INVAILDDATA   int32 = 10108
	//NOUSEFULCONN  int32 = 10200
	//NETWORKEXCEPT int32 = 10201 //网络异常
	//INVAILDLB     int32 = 10202 //lb找不到有效节点
	//INVAILDLBSRV  int32 = 10203 //找不到lb节点
	//INVAILDSRV    int32 = 10204 //找不到业务节点

	INVAILDPARAM  int32 = 10139
	INVAILDHANDLE int32 = 10140
	INVAILDDATA   int32 = 10141

	NOUSEFULCONN  int32 = 10221
	NETWORKEXCEPT int32 = 10222 //网络异常
	INVAILDLB     int32 = 10223 //lb找不到有效节点
	INVAILDLBSRV  int32 = 10224 //找不到lb节点
	INVAILDSRV    int32 = 10225 //找不到业务节点
	INVAILDRMLB   int32 = 10226 //请求lb失败
)

// 常量
var (
	GRCMRS = 4 * 1024
	GRCMWS = 32 * 1024 * 1024

	ZKSESSION = 5 //单位秒
)

var (
	tmpNode = "_TmpNode" //标识临时节点
	tmpKV   = struct {
		tmpK string
		tmpV string
	}{
		tmpK: "tmpNode",
		tmpV: "1",
	}                    //发给引擎临时节点键值对
)

const (
	DIRECTENGIP = "directEngIp"
)

var loggerStd = (&utils.LoggerStderr{}).Init("")
