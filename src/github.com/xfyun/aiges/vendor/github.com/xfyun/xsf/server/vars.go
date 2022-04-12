package xsf

import (
	"fmt"
	"sync"
	"time"
)

//----------------------------------------------------------------
var rateLimiterErr = fmt.Errorf("request traffic exceeds limit,and rateFallback not defined")

const PEERADDR = "PeerAddr"
const remainVCpus = "vcpu"

//----------------------------------------------------------------
const (
	VCPUFILENAME        = "xsf.toml"
	VCPUGROUP           = "vcpu-g"    //虚拟cpu映射存储位置的group
	defaultVCPUGROUP    = "common"    //缺省common组
	VCPUSERVICE         = "vcpu-s"    //虚拟cpu映射存储位置的服务
	defaultVCPUSERVICE  = "xsf"       //缺省xsf
	VCPUVERSION         = "vcpu-v"    //虚拟cpu映射配置版本
	defaultVCPUVERSION  = "1.0.0"     //缺省1.0.0
	VCPUINTERVAL        = "vcpu-i"    //cpu信息采集的时间间隔
	defaultVCPUINTERVAL = time.Second //缺省一秒
	VCPUABLE            = "vcpu"
	defaultVCPUABLE     = false

	VCPUSEC = "cpu"

	VCPUFITTEDVALUE = 5
)

//----------------------------------------------------------------

const (
	//sth about keepalive
	//2019-07-30 16:13:31 add sth about keepalive
	KEEPALIVE           = "keepalive"
	DEFKEEPALIVE        = time.Duration(0) //启用keepalive check，值为探测的时间间隔，单位毫秒，0表示不启用，缺省不启用keeplive
	KEEPALIVETIMEOUT    = "keepalive-timeout"
	DEFKEEPALIVETIMEOUT = 1000 * time.Millisecond
)
const (
	//日志的读取字段，logsection为读取的Key，其余的为读取的val
	LOGSECTION       = "log"
	LOGLEVEL         = "level"
	FILENAME         = "file"
	MAXSIZE          = "size"
	MAXBACKUPS       = "count"
	MAXAGE           = "die"
	LOGASYNC         = "async"
	LOGCACHEMAXCOUNT = "cache"
	LOGBATCHSIZE     = "batch"
	LOGCALLER        = "caller"
	LOGWASH          = "wash"
	//日志的默认值
	defaultLOGLEVEL         = "warn"
	defaultFILENAME         = "xrpcs.log"
	defaultMAXSIZE          = 10
	defaultMAXBACKUPS       = 10
	defaultMAXAGE           = 10
	defaultLOGWASH          = 60
	defaultLOGASYNC         = true
	defaultLOGCACHEMAXCOUNT = -1
	defaultLOGBATCHSIZE     = 16 * 1024
	defaultCALLER           = false
)

//----------------------------------------------------------------
const (
	LOADREPORTER     = "lb"
	LBLBSTRATEGY     = "lbStrategy"
	LBZKLIST         = "zkList"
	LBROOT           = "root"
	LBROUTERTYPE     = "routerType"
	LBSUBROUTERTYPES = "subRouterTypes"
	LBREDIEHOST      = "redisHost"
	LBREDISPASSWD    = "redisPasswd"
	LBMAXACTIVE      = "maxActive"
	LBMAXIDLE        = "maxIdle"
	LBDB             = "db"
	LBIDLETIMEOUT    = "idleTimeOut"
	LBABLE           = "able"
	//lb默认值able
	defaultLBABLE = 0
)

//----------------------------------------------------------------
//此处仅配置在xsf里的部分
const (
	HERMES                 = "lbv2"
	HERMESABLE             = "able"
	HERMESSVC              = "sub"
	HERMESSUBSVC           = "subsvc"
	HERMESCFGNAME          = "sfc"
	HERMESAPIVERSION       = "apiversion"
	HERMESTASK             = "task"
	HERMESLBNAME           = "lbname"
	HERMESLBPROJECT        = "lbproject"
	HERMESLBGROUP          = "lbgroup"
	HERMESLBCLOUD          = "cloud"
	HERMESFINDERTTL        = "finderttl"
	HERMESBACKEND          = "backend"
	HERMESTIMEOUT          = "tm"
	HERMESFINDERMODE       = "findermode"
	defaultHERMESABLE      = false
	defaultHERMESFINDERTTL = time.Minute
	defaultHERMESBACKEND   = 4
	defaultHERMESTIMEOUT   = time.Second
	defaultHERMESTASK      = 10
	defaultHERMESCLOUD     = "0"

	HERMESTYPE = "type"
)

type lbReportExt struct {
	m     sync.Map
	empty bool
}

func (l *lbReportExt) get(key string) string {
	k, kOk := l.m.Load(key)
	if !kOk {
		return ""
	}
	return k.(string)
}
func (l *lbReportExt) set(key, val string) {
	l.empty = false
	l.m.Store(key, val)
}
func (l *lbReportExt) getAll() map[string]string {
	if l.empty {
		return nil
	}
	rst := make(map[string]string)
	l.m.Range(func(key, value interface{}) bool {
		rst[key.(string)] = value.(string)
		return true
	})
	return rst
}
func (l *lbReportExt) setType(val string) {
	l.set(HERMESTYPE, val)
}
func (l *lbReportExt) getType() string {
	return l.get(HERMESTYPE)
}

var lbReportExtInst = lbReportExt{empty: false}

func SetLbType(val string) {
	lbReportExtInst.setType(val)
}

//----------------------------------------------------------------
//local的读取字段，svcsection为读取的Key，其余的为读取的val

var svcsection = "local"

var (
	MetricsCluster string
	MetricsIdc     string
	MetricsSub     string
)

const (
	METRICS            = "metrics"
	METRICSIDC         = "idc"
	defaultMETRICSIDC  = "dx"
	METRICSSUB         = "sub"
	defaultMETRICSSUB  = "xxx"
	METRICSCS          = "cs"
	defaultMETRICSCS   = "3s"
	METRICSABLE        = "able"
	defaultMETRICSABLE = false

	//2019-04-11 17:11:32
	METRICSTIMEPERSLICE        = "timePerSlice"
	defaultMETRICSTIMEPERSLICE = time.Second
	METRICSWINSIZE             = "winSize"
	defaultMETRICSWINSIZE      = 60
)
const (
	GRPCTIMEOUT_  = "grpctimeout"
	IP_           = "host"
	NETCARD_      = "netcard"
	PORT_         = "port"
	REUSEPORT_    = "reuseport"
	FINDERSWITCH_ = "finder"
	MAXRECEIVE    = "maxreceive"
	MAXSEND       = "maxsend"
	CONRBUF       = "conn-rbuf" // 连接读缓冲区
	CONWBUF       = "conn-wbuf" // 连接写缓冲区

	rateLimiterRate  = "rate"  //令牌的填充速率
	rateLimiterBurst = "burst" //令牌缓存数量

	//默认值
	defaultREUSEPORT   = 0
	defaultFINDER      = 0
	defaultPORT        = 0
	defaultGRPCTIMEOUT = 120
	defaultMAXRECEIVE  = 4 * 1024 * 1024 //能收取的最大消息包大小，单位MB，缺省16MB
	defaultMAXSEND     = 4 * 1024 * 1024 //能发送的最大消息包大小，单位MB，缺省16MB
	defaultCONRBUF     = 0
	defaultCONWBUF     = 2 * 1024 * 1024
)

var GRPCTIMEOUT int

//----------------------------------------------------------------

const (
	FLOWCONTROL           = "fc"
	FCABLE                = "able"
	ROUTER                = "router"
	ROUTER2SESSIONMANAGER = "sessionManager"
	ROUTER2QPSLIMITER     = "qpsLimiter"
	STRATEGY              = "strategy"
	MAX                   = "max"
	BEST                  = "best"
	TTL                   = "ttl"
	WAVE                  = "wave"
	ROLLTIMEOUT           = "roll"
	REPORT                = "report"

	TASKSIZE        = "tasksize"
	TASKCHANNELSIZE = "taskchsize"

	AUTHFILTERWIN        = "filterwin"
	defaultAUTHFILTERWIN = time.Second * 1e1

	defaultROLLTIMEOUT     time.Duration = 5000
	defaultREPORT          time.Duration = 1000
	defaultTTL                           = 15000
	defaultWAVE                          = 10
	defaultFcAble                        = 0
	defaultSTRATEGY                      = 0
	defaultTASKSIZE                      = 10
	defaultTASKCHANNELSIZE               = 10
)

//----------------------------------------------------------------
const (
	//读取trace相关信息
	TRACE     = "trace"
	TRACEHOST = "host"
	TRACEPORT = "port"
	DUMP      = "dump"
	ABLE      = "able"
	DELIVER   = "deliver"
	BACKEND   = "backend"
	SPILL     = "spill"
	BUFFER    = "buffer"
	BATCH     = "batch"
	LINGER    = "linger"
	WATCH     = "watch"
	WATCHPORT = "watchport"
	SPILLSIZE = "spillsize"
	LOADTS    = "loadts"

	TRACEBCLUSTER = "bcluster" //业务集群标识
	TRACEIDC      = "idc"      //IDC标识位

	//默认值
	defaultLoadTs    = 100
	defaultSpillSize = 1
	defaultWATCH     = false
	defaultWatchPort = 12331
	defaultTRACEHOST = "127.0.0.1"
	defaultTRACEPORT = 4545
	defaultDUMP      = 0
	defaultABLE      = 1
	defaultDELIVER   = 1
	defaultBACKEND   = 4
	defaultSPILL     = "/log/spill"
	defaultBUFFER    = 100000
	defaultBATCH     = 100
	defaultLINGER    = 5

	defaultTRACEBCLUSTER = "3s"
	defaultTRACEIDC      = "dz"
	defaultUNCHANGE      = -1
)
const (
	//读取sonar相关信息
	SONAR        = "sonar"
	SONARHOST    = "host"
	SONARPORT    = "port"
	SONARDUMP    = "dump"
	SONARABLE    = "able"
	SONARDS      = "ds"
	SONARDELIVER = "deliver"
	SONARBACKEND = "backend"

	//默认值
	defaultSONARHOST    = "127.0.0.1"
	defaultSONARPORT    = 4545
	defaultSONARDUMP    = 0
	defaultSONARABLE    = 1
	defaultSONARDELIVER = 1
	defaultSONARBACKEND = 4
	defaultSONARDS      = "vagus"
)

//----------------------------------------------------------------
const (
	BVT               = "bvt"
	BVTABLE           = "able"
	defBVTABLE        = false
	BVTSERVICE        = "service"
	BVTVERSION        = "version"
	BVTCFGFILE        = "cfgFile"
	BVTTIMEOUT        = "timeout"
	BVTPLATFORM       = "platform"
	BVTID             = "id"
	BVTSERVICEADDRESS = "serviceAddress"
	BVTNAMESPACE      = "namespace"
	BVTCALLBACK       = "callback"
	BVTASYNC          = "async"
)

//----------------------------------------------------------------
const (
	PING = "ping"
)
