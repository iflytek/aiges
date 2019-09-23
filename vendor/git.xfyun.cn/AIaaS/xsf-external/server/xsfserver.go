package xsf

import (
	"context"
	"errors"
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/client"
	"golang.org/x/time/rate"
	"time"

	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"google.golang.org/grpc/peer"
	"net"
	"os"
	"strconv"
	"sync/atomic"
)

type Net struct {
	ipStr   string
	portInt int
}

func (n *Net) GetIp() string {
	return n.ipStr
}
func (n *Net) GetPort() int {
	return n.portInt
}
func (n *Net) getLocalIp() string {
	return n.GetIp() + ":" + strconv.Itoa(n.GetPort())
}

/*
LookupHost:
desc:
根据host到指定的DNS上拉取对应的ip信息。windows 下默认使用系统DNS配置
*/
func (n *Net) GetHostByName(host string, dns string) ([]string, error) {
	if host == "" {
		hostname, hostnameErr := os.Hostname()
		if hostnameErr != nil {
			return nil, hostnameErr
		}
		return utils.LookupHost(hostname, dns)
	}
	return utils.LookupHost(host, dns)
}

type ToolBox struct {
	Cache      *SessionManager
	Qps        *QpsLimiter
	Cfg        *utils.Configure
	Log        *utils.Logger
	sid        *XrpcSidGenerator
	//Monitor    *SonarAdapter
	NetManager *Net
	Bc         *BootConfig

	lis net.Listener

	errWin   *errCodeWin
	delayWin *delayWin

	rateLimiter *rate.Limiter

	vCpuManager *VCpuManager
}
type TraceMeta struct {
	ip          string
	port        int
	serviceName string
}

func init() {
	go SignalHandle()
}
func (t *ToolBox) Init(bc BootConfig) error {
	t.Bc = &bc
	svcsection = bc.CfgData.Service

	//初始化configurator
	cfg, cfgErr := t.initCfg(bc, &utils.CfgOption{})
	if nil != cfgErr {
		return fmt.Errorf("CreateConfiguratorErr:%v,logCfgOpt:%+v", cfgErr, &utils.CfgOption{})
	}
	t.Cfg = cfg

	//读取日志相关配置，初始化日志
	if logErr := t.newLocalLog(cfg, bc); nil != logErr {
		return fmt.Errorf("loggerErr:%v", logErr)
	}
	addKillerCheck(killerLowPriority, "logger", t.Log)

	//读取bootConfig中service字段作为配置中ip、port、sub的入口
	//初始化lis、NetManager
	ip, port, netManagerErr := t.InitNetManager(cfg)
	if nil != netManagerErr {
		return netManagerErr
	}

	//=================================================================
	t.initGrpcVars(cfg)

	//=================================================================
	//初始化loadReporter
	//读取loadReporter的配置
	lbClient, lbClientErr := t.InitLbClient(cfg, ip, port)
	if nil != lbClientErr {
		return lbClientErr
	}

	//=================================================================
	//初始化hermesAdapter
	loggerStd.Println("about to deal with hermes.")
	lbv2Client, lbv2ClientErr := t.InitLbv2Client(cfg, bc, ip, port)
	if nil != lbv2ClientErr {
		return fmt.Errorf("lbv2Client init fail. err:%v\n", lbv2ClientErr)
	}

	//=================================================================
	//读取flowControl的配置
	fcInit := t.InitFc(cfg, bc, lbClient, lbv2Client)
	if nil != fcInit {
		return fcInit
	}
	//=================================================================
	//初始化rpcSidGenerator.go
	var sidErr error
	t.sid, sidErr = NewSidGenerator(sidVer, ip, int64(port))
	if nil != sidErr {
		return sidErr
	}
	loggerStd.Printf("NewSidGenerator success.\n")
	t.Log.Infof("NewSidGenerator success.\n")

	//=================================================================
	//初始化trace
	//traceErr := t.InitTrace(cfg, bc, port, ip)
	//if nil != traceErr {
	//	return traceErr
	//}

	//=================================================================
	//初始化sonar
	//sonarErr := t.InitSonar(cfg, ip, bc, port)
	//if nil != sonarErr {
	//	return sonarErr
	//}
	//=================================================================
	//初始化finder
	finderErr := t.InitFinder(cfg, bc, ip, port)
	if nil != finderErr {
		return finderErr
	}
	//=================================================================

	//=================================================================
	//初始化rateLimiter
	t.InitRateLimiter(cfg)

	//=================================================================
	//初始化vCpuManager
	vCpuManagerErr := t.InitVCpuManager(cfg, bc)
	if !(nil == vCpuManagerErr) {
		return vCpuManagerErr
	}

	return nil
}

func (t *ToolBox) InitVCpuManager(cfg *utils.Configure, bc BootConfig) error {
	vCpuAble := defaultVCPUABLE
	vCpuAbleInt, vCpuAbleErr := cfg.GetInt(svcsection, VCPUABLE)
	if nil == vCpuAbleErr {
		if 1 == vCpuAbleInt {
			vCpuAble = true
		}
	}
	if vCpuAble {
		vCpuGroup, e := cfg.GetString(svcsection, VCPUGROUP)
		if e != nil {
			vCpuGroup = defaultVCPUGROUP
		}
		vCpuService, e := cfg.GetString(svcsection, VCPUSERVICE)
		if e != nil {
			vCpuService = defaultVCPUSERVICE
		}
		vCpuVersion, e := cfg.GetString(svcsection, VCPUVERSION)
		if e != nil {
			vCpuVersion = defaultVCPUVERSION
		}

		vCpuInterval := defaultVCPUINTERVAL
		vCpuIntervalInt64, e := cfg.GetInt64(svcsection, VCPUINTERVAL)
		if e == nil {
			vCpuInterval = time.Millisecond * time.Duration(vCpuIntervalInt64)
		}

		vCpuMap, vCpuMapErr := getVCpuCfg(
			vCpuVersion,
			bc.CfgData.Project,
			vCpuGroup,
			vCpuService,
			VCPUFILENAME,
			bc.CfgData.CompanionUrl,
			bc.CfgMode,
		)
		if nil != vCpuMapErr {
			t.Log.Errorf("fetch vManager cfg failed:%v,ver:%v,prj:%v,group:%v,service:%v,url:%v,file:%v,mode:%v\n",
				vCpuMapErr, vCpuVersion, bc.CfgData.Project, vCpuGroup, vCpuService, bc.CfgData.CompanionUrl, VCPUFILENAME, bc.CfgMode)
		}
		loggerStd.Printf("vCpuManager:%+v\n", vCpuMap)
		t.vCpuManager, e = NewVCpuManager(vCpuInterval, vCpuMap)
		if nil != e {
			return fmt.Errorf("init vCpuManager failed,err:%v", e)
		} else {
			loggerStd.Println("vCpuManager is disable")
		}
	}
	return nil
}

func (t *ToolBox) InitRateLimiter(cfg *utils.Configure) {
	rateInt, _ := cfg.GetInt(svcsection, rateLimiterRate)
	burstInt, _ := cfg.GetInt(svcsection, rateLimiterBurst)
	if 0 != rateInt || 0 != burstInt {
		loggerStd.Printf("rateLimiter,rate:%v,burst:%v\n", rateInt, burstInt)
		t.rateLimiter = rate.NewLimiter(rate.Every(time.Duration(rateInt)*time.Millisecond), burstInt)
	}
}

func (t *ToolBox) InitFinder(cfg *utils.Configure, bc BootConfig, ip string, port int) error {
	finderSwitch, finderSwitchErr := cfg.GetInt(svcsection, FINDERSWITCH_)
	loggerStd.Printf("finderSwitch:%v,finderSwitchErr:%v\n", finderSwitch, finderSwitchErr)
	t.Log.Infof("finderSwitch:%v,finderSwitchErr:%v\n", finderSwitch, finderSwitchErr)
	if nil != finderSwitchErr {
		finderSwitch = defaultFINDER
	}
	if 0 != finderSwitch {
		finderCfgOpt := &utils.CfgOption{}
		utils.WithCfgDefault(bc.CfgData.CfgName)(finderCfgOpt)
		utils.WithCfgVersion(bc.CfgData.Version)(finderCfgOpt)
		utils.WithCfgPrj(bc.CfgData.Project)(finderCfgOpt)
		utils.WithCfgGroup(bc.CfgData.Group)(finderCfgOpt)
		utils.WithCfgService(bc.CfgData.Service)(finderCfgOpt)
		utils.WithCfgName(bc.CfgData.CfgName)(finderCfgOpt)
		utils.WithCfgURL(bc.CfgData.CompanionUrl)(finderCfgOpt)
		utils.WithCfgCB(bc.CfgData.CallBack)(finderCfgOpt)
		utils.WithCfgLog(t.Log)(finderCfgOpt)
		finder, finderErr := utils.NewFinder(finderCfgOpt)

		if nil != finderErr {
			return fmt.Errorf("CreateFinder fail -> bc:%+v, finderErr:%v", bc, finderErr)
		}
		loggerStd.Printf("CreateFinder success.\n")
		t.Log.Infof("CreateFinder success.\n")

		loggerStd.Printf("about to call finderadapter.AddRegister. addr:%s\n", fmt.Sprintf("%v:%v", ip, port))
		t.Log.Infof("about to call finderAdapter.AddRegister. addr:%s\n", fmt.Sprintf("%v:%v", ip, port))

		finderadapter.AddRegister(bc.CfgData.ApiVersion, fmt.Sprintf("%v:%v", ip, port), finder)
	}
	return nil
}

//func (t *ToolBox) InitSonar(cfg *utils.Configure, ip string, bc BootConfig, port int) error {
//	sonarHost, sonarHostErr := cfg.GetString(SONAR, SONARHOST)
//	if sonarHostErr != nil {
//		sonarHost = defaultSONARHOST
//	}
//	sonarPort, sonarPortErr := cfg.GetInt(SONAR, SONARPORT)
//	if sonarPortErr != nil {
//		sonarPort = defaultSONARPORT
//	}
//	sonarBackend, sonarBackendErr := cfg.GetInt(SONAR, SONARBACKEND)
//	if sonarBackendErr != nil {
//		sonarBackend = defaultSONARBACKEND
//	}
//	sonarDeliver, sonarDeliverErr := cfg.GetInt(SONAR, SONARDELIVER)
//	if sonarDeliverErr != nil {
//		sonarDeliver = defaultSONARDELIVER
//	}
//	sonarDump, sonarDumpErr := cfg.GetInt(SONAR, SONARDUMP)
//	if sonarDumpErr != nil {
//		sonarDump = defaultSONARDUMP
//	}
//	sonarAble, sonarAbleErr := cfg.GetInt(SONAR, SONARABLE)
//	if sonarAbleErr != nil {
//		sonarAble = defaultSONARABLE
//	}
//	sonarDS, sonarDSErr := cfg.GetString(SONAR, SONARDS)
//	if sonarDSErr != nil {
//		sonarDS = defaultSONARDS
//	}
//	sonardeliver := false
//	sonardump := false
//	sonarable := false
//	if sonarDeliver != 0 {
//		sonardeliver = true
//	}
//	if sonarDump != 0 {
//		sonardump = true
//	}
//	if sonarAble != 0 {
//		sonarable = true
//	}
//	t.Monitor = &SonarAdapter{}
//	sonarErr := t.Monitor.initSonar(
//		WithSonarAdapterAble(sonarable),
//		WithSonarAdapterDs(sonarDS),
//		WithSonarAdapterMetricEndpoint(ip),
//		WithSonarAdapterMetricServiceName(bc.CfgData.Service),
//		WithSonarAdapterMetricPort(strconv.Itoa(port)),
//		WithSonarAdapterLogger(nil),
//		WithSonarAdapterSonarDumpEnable(sonardump),
//		WithSonarAdapterSonarDeliverEnable(sonardeliver),
//		WithSonarAdapterSonarHost(sonarHost),
//		WithSonarAdapterSonarPort(strconv.Itoa(sonarPort)),
//		WithSonarAdapterSonarBackend(sonarBackend))
//	if sonarable {
//		if sonarErr != nil {
//			return fmt.Errorf("sonarErr:%v", sonarErr)
//		}
//		loggerStd.Printf("sonar init success.\n")
//		t.Log.Infof("sonar init success.\n")
//	}
//	return nil
//}

//func (t *ToolBox) InitTrace(cfg *utils.Configure, bc BootConfig, port int, ip string) error {
//	traceHost, traceHostErr := cfg.GetString(TRACE, TRACEHOST)
//	if traceHostErr != nil {
//		traceHost = defaultTRACEHOST
//	}
//	tracePort, tracePortErr := cfg.GetInt(TRACE, TRACEPORT)
//	if tracePortErr != nil {
//		tracePort = defaultTRACEPORT
//	}
//	backend, backendErr := cfg.GetInt(TRACE, BACKEND)
//	if backendErr != nil {
//		backend = defaultBACKEND
//	}
//	deliver_, deliverErr := cfg.GetInt(TRACE, DELIVER)
//	if deliverErr != nil {
//		deliver_ = defaultDELIVER
//	}
//	dump_, dumpErr := cfg.GetInt(TRACE, DUMP)
//	if dumpErr != nil {
//		dump_ = defaultDUMP
//	}
//	able_, ableErr := cfg.GetInt(TRACE, ABLE)
//	if ableErr != nil {
//		able_ = defaultABLE
//	}
//	spill, spillErr := cfg.GetString(TRACE, SPILL)
//	if spillErr != nil {
//		spill = defaultSPILL
//	}
//	buffer, bufferErr := cfg.GetInt(TRACE, BUFFER)
//	if bufferErr != nil {
//		buffer = defaultBUFFER
//	}
//	batch, batchErr := cfg.GetInt(TRACE, BATCH)
//	if batchErr != nil {
//		batch = defaultBATCH
//	}
//	linger, lingerErr := cfg.GetInt(TRACE, LINGER)
//	if lingerErr != nil {
//		linger = defaultLINGER
//	}
//	watchBool := defaultWATCH
//	watchInt, watchErr := cfg.GetInt(TRACE, WATCH)
//	if watchErr == nil {
//		if watchInt == 1 {
//			watchBool = true
//		}
//	}
//	bcluster := defaultTRACEBCLUSTER
//	bclusterStr, bclusterStrErr := cfg.GetString(TRACE, TRACEBCLUSTER)
//	if bclusterStrErr == nil {
//		bcluster = bclusterStr
//	}
//	idc := defaultTRACEIDC
//	idcStr, idcStrErr := cfg.GetString(TRACE, TRACEIDC)
//	if idcStrErr == nil {
//		idc = idcStr
//	}
//	watchPort := defaultWatchPort
//	watchPortInt, watchPortErr := cfg.GetInt(TRACE, WATCHPORT)
//	if watchPortErr == nil {
//		watchPort = watchPortInt
//	}
//	spillSize := defaultSpillSize
//	spillSizeInt, spillSizeErr := cfg.GetInt(TRACE, SPILLSIZE)
//	if spillSizeErr == nil {
//		spillSize = spillSizeInt
//	}
//	loadTs := defaultLoadTs
//	loadTsInt, loadTsErr := cfg.GetInt(TRACE, LOADTS)
//	if loadTsErr == nil {
//		loadTs = loadTsInt
//	}
//	deliver := false
//	dump := false
//	able := false
//	if deliver_ == 1 {
//		deliver = true
//	}
//	if dump_ == 1 {
//		dump = true
//	}
//	if able_ == 1 {
//		able = true
//	}
//	if able_ != defaultUNCHANGE {
//		utils.AbleTrace(able)
//		if able {
//			loggerStd.Printf("traceHost:%v, tracePort:%v, loadTs:%v, spillSize:%v, watchPort:%v, watchBool:%v,bc.CfgData.Service:%v, port:%v, ip:%v,buffer:%v, batch:%v, linger:%v, spill:%v, backend:%v, deliver:%v, dump:%v\n",
//				traceHost, tracePort, loadTs, spillSize, watchPort, watchBool, bc.CfgData.Service, port, ip, buffer, batch, linger, spill, backend, deliver, dump)
//			if traceErr := utils.InitTracer(
//				traceHost,
//				strconv.Itoa(tracePort),
//				utils.WithLowLoadSleepTs(loadTs),
//				utils.WithMaxSpillContentSize(spillSize),
//				utils.WithWatchPort(watchPort),
//				utils.WithWatch(watchBool),
//				utils.WithSvcName(bc.CfgData.Service),
//				utils.WithSvcPort(int32(port)),
//				utils.WithSvcIp(ip),
//				utils.WithBufferSize(buffer),
//				utils.WithBatchSize(batch),
//				utils.WithLinger(linger),
//				utils.WithTraceSpill(spill),
//				utils.WithBackend(backend),
//				utils.WithDeliver(deliver),
//				utils.WithDump(dump),
//				utils.WithSvcBCluster(bcluster),
//				utils.WithSvcIDC(idc),
//				utils.WithTraceLogger(t.Log)); traceErr != nil {
//				return fmt.Errorf("InitTracer failed -> able:%v,ip:%v,port:%v,backend:%v,deliver:%v,dump:%v -> traceErr:%v", able, traceHost, tracePort, backend, deliver, dump, traceErr)
//			}
//		}
//	}
//	return nil
//}

func (t *ToolBox) InitFc(cfg *utils.Configure, bc BootConfig, lbClient *LbAdapter, lbv2Client *hermesAdapter) error {
	fcAble, fcAbleErr := cfg.GetInt(FLOWCONTROL, FCABLE)
	if nil != fcAbleErr {
		fcAble = defaultFcAble
	}

	if 0 != fcAble {
		router, routerErr := cfg.GetString(FLOWCONTROL, ROUTER)
		max, maxErr := cfg.GetInt(FLOWCONTROL, MAX)
		ttl, ttlErr := cfg.GetInt(FLOWCONTROL, TTL)
		if nil != ttlErr {
			ttl = defaultTTL
		}
		best, bestErr := cfg.GetInt(FLOWCONTROL, BEST)
		if nil != bestErr {
			best = max
		}
		wave, waveErr := cfg.GetInt(FLOWCONTROL, WAVE)
		if nil != waveErr {
			wave = defaultWAVE
		}
		strategy, strategyErr := cfg.GetInt(FLOWCONTROL, STRATEGY)
		if nil != strategyErr {
			strategy = defaultSTRATEGY
		}
		roll, rollErr := cfg.GetInt(FLOWCONTROL, ROLLTIMEOUT)
		rollTime := time.Duration(roll)
		if nil != rollErr {
			rollTime = defaultROLLTIMEOUT
		}
		report, reportErr := cfg.GetInt(FLOWCONTROL, REPORT)
		reportInterval := time.Duration(report)
		if nil != reportErr {
			reportInterval = defaultREPORT
		}
		taskSize, taskSizeErr := cfg.GetInt(FLOWCONTROL, TASKSIZE)
		if nil != taskSizeErr {
			taskSize = defaultTASKSIZE
		}
		taskChannelSize, taskChannelSizeErr := cfg.GetInt(FLOWCONTROL, TASKCHANNELSIZE)
		if nil != taskChannelSizeErr {
			taskChannelSize = defaultTASKCHANNELSIZE
		}
		if nil != routerErr || nil != maxErr {
			return errors.New(fmt.Sprintf("can't get the router、max from configurator -> routerErr:%v, maxErr:%v",
				routerErr, maxErr))
		}
		if router == ROUTER2SESSIONMANAGER {
			var SessionManagerErr error
			t.Cache, SessionManagerErr = NewSessionManager(
				WithSessionManagerTaskSize(taskSize),
				WithSessionManagerTaskChannelSize(taskChannelSize),
				WithSessionManagerBc(bc),
				WithSessionManagerMaxLic(int32(max)),
				WithSessionManagerBestLic(int32(best)),
				WithSessionManagerTimeout(time.Duration(ttl)*time.Millisecond),
				WithSessionManagerRollTime(time.Duration(rollTime)*time.Millisecond),
				WithSessionManagerReportInterval(int32(reportInterval)),
				WithSessionManagerReporter(lbClient),
				WithSessionManagerReporterv2(lbv2Client),
				WithSessionManagerLogger(t.Log),
				WithSessionManagerStrategy(strategy),
				WithSessionManagerWave(wave))
			if nil != SessionManagerErr {
				return SessionManagerErr
			}
			loggerStd.Printf("NewSessionManager success.\n")
			t.Log.Infof("NewSessionManager success.\n")
		} else if router == ROUTER2QPSLIMITER {
			var QpsLimiterErr error
			t.Qps, QpsLimiterErr = NewQpsLimiter(
				WithQpsLimiterBc(bc),
				WithQpsLimiterMaxReqCount(int32(max)),
				WithQpsLimiterBestReqCount(int32(best)),
				WithQpsLimiterInterval(int32(ttl)),
				WithQpsLimiterReportInterval(reportInterval),
				WithQpsLimiterReporter(lbClient),
				WithQpsLimiterReporterV2(lbv2Client),
				WithQpsLimiterLogger(t.Log))
			if nil != QpsLimiterErr {
				return fmt.Errorf("QpsLimiterErr:%v", QpsLimiterErr)
			}
			loggerStd.Printf("NewQpsLimiter success.\n")
			t.Log.Infof("NewQpsLimiter success.\n")
		}
	}
	return nil
}

func (t *ToolBox) InitLbv2Client(cfg *utils.Configure, bc BootConfig, ip string, port int) (hermes *hermesAdapter, hermesErr error) {
	hermes = &hermesAdapter{}
	if nil != cfg.GetSection(HERMES) {
		lbAbleInt, lbAbleErr := cfg.GetInt(HERMES, HERMESABLE)
		lbAble := defaultHERMESABLE
		if nil == lbAbleErr {
			if 0 != lbAbleInt {
				lbAble = true
			}
		}
		if !lbAble {
			hermes.able = false
			return
		} else {
			hermes.able = true
		}
		lbv2Svc, lbv2SvcErr := cfg.GetString(HERMES, HERMESSVC)
		if nil != lbv2SvcErr {
			hermesErr = fmt.Errorf("can't take %v from %v,err:%v", HERMESSVC, HERMES, lbv2SvcErr)
			return
		}
		lbv2SubSvc, lbv2SubSvcErr := cfg.GetString(HERMES, HERMESSUBSVC)
		if nil != lbv2SubSvcErr {
			hermesErr = lbv2SubSvcErr
		}

		lbName, lbNameErr := cfg.GetString(HERMES, HERMESLBNAME)
		if nil != lbNameErr {
			hermesErr = lbNameErr
		}

		apiVersion := bc.CfgData.ApiVersion
		apiVersionStr, apiVersionStrErr := cfg.GetString(HERMES, HERMESAPIVERSION)
		if nil == apiVersionStrErr {
			apiVersion = apiVersionStr
		}

		hermesTask := defaultHERMESTASK
		hermesTaskInt64, hermesTaskInt64Err := cfg.GetInt64(HERMES, HERMESTASK)
		if nil == hermesTaskInt64Err {
			hermesTask = int(hermesTaskInt64)
		}

		finderTtl := defaultHERMESFINDERTTL
		finderTtlInt, finderTtlErr := cfg.GetInt(HERMES, HERMESFINDERTTL)
		if nil == finderTtlErr {
			finderTtl = time.Millisecond * time.Duration(finderTtlInt)
		}

		backendInt, backendIntErr := cfg.GetInt(HERMES, HERMESBACKEND)
		if nil != backendIntErr {
			backendInt = defaultHERMESBACKEND
		}

		HermesTimeout := defaultHERMESTIMEOUT
		HermesTimeoutInt, HermesTimeoutIntErr := cfg.GetInt(HERMES, HERMESTIMEOUT)
		if nil == HermesTimeoutIntErr {
			HermesTimeout = time.Duration(HermesTimeoutInt) * time.Millisecond
		}

		if InitErr := hermes.Init(
			WithHermesAdapterLbApiVersion(apiVersion),
			WithHermesAdapterLbName(lbName),
			WithHermesAdapterAddr(fmt.Sprintf("%v:%v", t.NetManager.ipStr, t.NetManager.portInt)),
			WithHermesAdapterSvc(lbv2Svc),
			WithHermesAdapterSubsvc(lbv2SubSvc),
			WithHermesAdapterFinderTtl(finderTtl),
			WithHermesAdapterBackEnd(backendInt),
			WithHermesAdapterTimeout(HermesTimeout),
			WithHermesAdapterBootConfig(bc),
			WithHermesAdapterSvcIp(ip),
			WithHermesAdapterTask(hermesTask),
			WithHermesAdapterSvcPort(int32(port)), ); InitErr != nil {
			hermesErr = fmt.Errorf("InitErr:%v\n", InitErr)
			return
		}
		loggerStd.Println("hermes init success.")
	} else {
		hermes.able = false
		loggerStd.Println("hermes is disable.")
	}
	return
}

func (t *ToolBox) InitLbClient(cfg *utils.Configure, ip string, port int) (*LbAdapter, error) {
	var lbClient *LbAdapter
	if nil == cfg.GetSection(LOADREPORTER) {
		lbClient = &LbAdapter{able: defaultLBABLE}
		return lbClient, nil
	}

	lbAble, lbAbleErr := cfg.GetInt(LOADREPORTER, LBABLE)
	if nil != lbAbleErr {
		lbAble = defaultLBABLE
	}
	lbLbStrategy, lbLbStrategyErr := cfg.GetInt(LOADREPORTER, LBLBSTRATEGY)
	lbZkList, lbZkListErr := interface2stringslice(cfg.GetAsObject(LOADREPORTER, LBZKLIST))
	lbRoot, lbRootErr := cfg.GetString(LOADREPORTER, LBROOT)
	if nil != lbRootErr {
		lbRoot = ""
	}
	lbRouterType, lbRouterTypeErr := cfg.GetString(LOADREPORTER, LBROUTERTYPE)
	lbSubRouterTypes, lbSubRouterTypesErr := interface2stringslice(cfg.GetAsObject(LOADREPORTER, LBSUBROUTERTYPES))
	if nil != lbSubRouterTypesErr {
		lbSubRouterTypes = nil
	}
	lbRedisHost, lbRedisHostErr := cfg.GetString(LOADREPORTER, LBREDIEHOST)
	lbRedisPassWd, lbRedisPassWdErr := cfg.GetString(LOADREPORTER, LBREDISPASSWD)
	if nil != lbRedisPassWdErr {
		lbRedisPassWd = ""
	}
	lbMaxActive, lbMaxActiveErr := cfg.GetInt(LOADREPORTER, LBMAXACTIVE)
	if nil != lbMaxActiveErr {
		lbMaxActive = 0
	}
	lbMaxIdle, lbMaxIdleErr := cfg.GetInt(LOADREPORTER, LBMAXIDLE)
	if nil != lbMaxIdleErr {
		lbMaxIdle = 0
	}
	lbDb, lbDbErr := cfg.GetInt(LOADREPORTER, LBDB)
	if nil != lbDbErr {
		lbDb = 0
	}
	lbIdleTimeout, lbIdleTimeoutErr := cfg.GetInt(LOADREPORTER, LBIDLETIMEOUT)
	if nil != lbIdleTimeoutErr {
		lbIdleTimeout = 0
	}
	if lbAble != defaultLBABLE {
		if nil != lbLbStrategyErr || nil != lbZkListErr || nil != lbRouterTypeErr || nil != lbRedisHostErr {
			return nil, fmt.Errorf("lbLbStrategyErr:%v,lbZkListErr:%v,lbRouterTypeErr:%v,lbredishostErr:%v",
				lbLbStrategyErr, lbZkListErr, lbRouterTypeErr, lbRedisHostErr)
		}
	}
	lbClient = &LbAdapter{able: lbAble}
	if lbAble == defaultLBABLE {
		loggerStd.Printf("lb is not enabled\n")
		t.Log.Infof("lb is not enabled\n")
		return nil, nil
	}
	InitErr := lbClient.Init(
		WithLbAdapterSvc(fmt.Sprintf("%v:%v", ip, port)),
		WithLbAdapterStrategy(lbLbStrategy),
		WithLbAdapterZkList(lbZkList),
		WithLbAdapterRoot(lbRoot),
		WithLbAdapterRouterType(lbRouterType),
		WithLbAdapterSubRouterTypes(lbSubRouterTypes),
		WithLbAdapterSRedisHost(lbRedisHost),
		WithLbAdapterSRedisPasswd(lbRedisPassWd),
		WithLbAdapterMaxActive(lbMaxActive),
		WithLbAdapterMaxIdle(lbMaxIdle),
		WithLbAdapterDb(lbDb),
		WithLbAdapterIdleTimeOut(time.Second*time.Duration(lbIdleTimeout)))
	if nil != InitErr {
		return nil, fmt.Errorf("InitErr:%v\n", InitErr)
	}
	loggerStd.Printf("lbClient.Init success. -> lbClient:%+v\n", lbClient)
	t.Log.Infof("lbClient.Init success. -> lbClient:%+v\n", lbClient)

	return lbClient, nil
}

func (t *ToolBox) initGrpcVars(cfg *utils.Configure) {
	GRPCTIMEOUT = defaultGRPCTIMEOUT
	grpcTimeout, grpcTimeoutErr := cfg.GetInt(svcsection, GRPCTIMEOUT_)
	if nil == grpcTimeoutErr {
		GRPCTIMEOUT = grpcTimeout
	}
	maxReceiveInt, maxReceiveIntErr := cfg.GetInt(svcsection, MAXRECEIVE)
	if nil != maxReceiveIntErr {
		maxReceiveInt = defaultMAXRECEIVE
	} else {
		maxReceiveInt = maxReceiveInt * 1024 * 1024
	}
	grpcOptInst.setMaxReceiveMessageSize(maxReceiveInt)
	maxSendInt, maxSendIntErr := cfg.GetInt(svcsection, MAXSEND)
	if nil != maxSendIntErr {
		maxSendInt = defaultMAXSEND
	} else {
		maxSendInt = maxSendInt * 1024 * 1024
	}
	grpcOptInst.setMaxSendMessageSize(maxSendInt)
	conRBufInt, conRBufIntErr := cfg.GetInt(svcsection, CONRBUF)
	if nil != conRBufIntErr {
		conRBufInt = defaultCONRBUF
	} else {
		conRBufInt = conRBufInt * 1024 * 1024
	}
	grpcOptInst.setReadBufferSize(conRBufInt)
	conWBufInt, conWBufIntErr := cfg.GetInt(svcsection, CONWBUF)
	if nil != conWBufIntErr {
		conWBufInt = defaultCONWBUF
	} else {
		conWBufInt = conWBufInt * 1024 * 1024
	}
	grpcOptInst.setWriteBufferSize(conWBufInt)
}

func (t *ToolBox) InitNetManager(cfg *utils.Configure) (ip string, port int, err error) {
	ip, _ = cfg.GetString(svcsection, IP_)
	netCard, _ := cfg.GetString(svcsection, NETCARD_)
	ip, ipErr := utils.Host2Ip(ip, netCard)
	if nil != ipErr {
		return ip, port, fmt.Errorf("host2Ip:%v,ip:%v,netCard:%v", ipErr, ip, netCard)
	}
	port, portErr := cfg.GetInt(svcsection, PORT_)
	if nil != portErr {
		port = defaultPORT
	}

	reusePort, reusePortErr := t.Cfg.GetInt(svcsection, REUSEPORT_)
	if reusePortErr != nil {
		reusePort = defaultREUSEPORT
	}
	t.lis, err = NewListener(reusePort, net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return ip, port, errors.New(fmt.Sprintf("can't listen %v", ip) + func() string {
			if 0 == port {
				return ""
			} else {
				return ":" + strconv.Itoa(port)
			}
		}())
	}
	var portStr string
	_, portStr, err = net.SplitHostPort(t.lis.Addr().String())
	if err != nil {
		return ip, port, fmt.Errorf("can't get ip and port from %v", t.lis.Addr().String())
	}
	port, err = strconv.Atoi(portStr)
	if err != nil {
		return ip, port, fmt.Errorf("can't convert %v to int,err->%v", portStr, err)
	}
	loggerStd.Printf("host2ip->ip:%v,port:%v\n", ip, port)
	t.Log.Infof("host2ip->ip:%v,port:%v\n", ip, port)
	t.NetManager = &Net{portInt: port, ipStr: ip}
	return
}

func (t *ToolBox) newLocalLog(cfg *utils.Configure, bc BootConfig) error {
	logLevel, logLevelErr := cfg.GetString(LOGSECTION, LOGLEVEL)
	fileName, fileNameErr := cfg.GetString(LOGSECTION, FILENAME)
	maxsize, maxSizeErr := cfg.GetInt(LOGSECTION, MAXSIZE)
	maxBackups, maxBackupsErr := cfg.GetInt(LOGSECTION, MAXBACKUPS)
	maxAge, maxAgeErr := cfg.GetInt(LOGSECTION, MAXAGE)
	logAsyncInt, logAsyncErr := cfg.GetInt(LOGSECTION, LOGASYNC)
	logCacheMaxCount, logCacheMaxCountErr := cfg.GetInt(LOGSECTION, LOGCACHEMAXCOUNT)
	logBatchSize, logBatchSizeErr := cfg.GetInt(LOGSECTION, LOGBATCHSIZE)
	logCallerInt, logCallerErr := cfg.GetInt(LOGSECTION, LOGCALLER)
	logWash, logWashErr := cfg.GetInt(LOGSECTION, LOGWASH)
	if logWashErr != nil {
		logWash = defaultLOGWASH
	}
	//日志的默认值
	if logLevelErr != nil {
		logLevel = defaultLOGLEVEL
	}
	if fileNameErr != nil {
		fileName = defaultFILENAME
	}
	if maxSizeErr != nil {
		maxsize = defaultMAXSIZE
	}
	if maxBackupsErr != nil {
		maxBackups = defaultMAXBACKUPS
	}
	if maxAgeErr != nil {
		maxAge = defaultMAXAGE
	}
	logAsync := false
	if logAsyncErr != nil {
		logAsync = defaultLOGASYNC
	} else if logAsyncInt != 0 {
		logAsync = true
	}
	if logCacheMaxCountErr != nil {
		logCacheMaxCount = defaultLOGCACHEMAXCOUNT
	}
	if logBatchSizeErr != nil {
		logBatchSize = defaultLOGBATCHSIZE
	}
	logCaller := false
	if logCallerErr != nil {
		logCaller = defaultCALLER
	} else if logCallerInt != 0 {
		logCaller = true
	}
	var loggerErr error
	t.Log, loggerErr = utils.NewLocalLog(
		utils.SetWash(logWash),
		utils.SetCaller(logCaller),
		utils.SetBatchSize(logBatchSize),
		utils.SetCacheMaxCount(logCacheMaxCount),
		utils.SetAsync(logAsync),
		utils.SetLevel(logLevel),
		utils.SetFileName(fileName),
		utils.SetMaxSize(maxsize),
		utils.SetMaxBackups(maxBackups),
		utils.SetMaxAge(maxAge),
	)
	if nil != loggerErr {
		return fmt.Errorf("loggerErr:%v", loggerErr)
	}
	loggerStd.Printf("utils.NewLocalLog success. -> LOGLEVEL:%v, FILENAME:%v, MAXSIZE:%v, MAXBACKUPS:%v, MAXAGE:%v\n",
		logLevel, fileName, maxsize, maxBackups, maxAge)
	t.Log.Errorf("xsfVer:%v,service:%v",
		utils.GetVer(), bc.CfgData.Service)
	t.Log.Infof("utils.NewLocalLog success. -> LOGLEVEL:%v, FILENAME:%v, MAXSIZE:%v, MAXBACKUPS:%v, MAXAGE:%v\n",
		logLevel, fileName, maxsize, maxBackups, maxAge)
	return nil
}

func (t *ToolBox) getLogCfg(cfg *utils.Configure) (string, string, int, int, int, int, int, int, bool, bool, error) {
	logLevel, logLevelErr := cfg.GetString(LOGSECTION, LOGLEVEL)
	fileName, fileNameErr := cfg.GetString(LOGSECTION, FILENAME)
	maxsize, maxSizeErr := cfg.GetInt(LOGSECTION, MAXSIZE)
	maxBackups, maxBackupsErr := cfg.GetInt(LOGSECTION, MAXBACKUPS)
	maxAge, maxAgeErr := cfg.GetInt(LOGSECTION, MAXAGE)
	logAsyncInt, logAsyncErr := cfg.GetInt(LOGSECTION, LOGASYNC)
	logCacheMaxCount, logCacheMaxCountErr := cfg.GetInt(LOGSECTION, LOGCACHEMAXCOUNT)
	logBatchSize, logBatchSizeErr := cfg.GetInt(LOGSECTION, LOGBATCHSIZE)
	logCallerInt, logCallerErr := cfg.GetInt(LOGSECTION, LOGCALLER)
	logWash, logWashErr := cfg.GetInt(LOGSECTION, LOGWASH)
	if logWashErr != nil {
		logWash = defaultLOGWASH
	}
	//日志的默认值
	if logLevelErr != nil {
		logLevel = defaultLOGLEVEL
	}
	if fileNameErr != nil {
		fileName = defaultFILENAME
	}
	if maxSizeErr != nil {
		maxsize = defaultMAXSIZE
	}
	if maxBackupsErr != nil {
		maxBackups = defaultMAXBACKUPS
	}
	if maxAgeErr != nil {
		maxAge = defaultMAXAGE
	}
	logasync := false
	if logAsyncErr != nil {
		logasync = defaultLOGASYNC
	} else if logAsyncInt != 0 {
		logasync = true
	}
	if logCacheMaxCountErr != nil {
		logCacheMaxCount = defaultLOGCACHEMAXCOUNT
	}
	if logBatchSizeErr != nil {
		logBatchSize = defaultLOGBATCHSIZE
	}
	logCaller := false
	if logCallerErr != nil {
		logCaller = defaultCALLER
	} else if logCallerInt != 0 {
		logCaller = true
	}
	var loggerErr error
	return logLevel, fileName, maxsize, maxBackups, maxAge, logCacheMaxCount, logBatchSize, logWash, logasync, logCaller, loggerErr
}

func (t *ToolBox) initCfg(bc BootConfig, logCfgOpt *utils.CfgOption) (*utils.Configure, error) {
	utils.WithCfgDefault(bc.CfgData.CfgDefault)(logCfgOpt)
	utils.WithCfgVersion(bc.CfgData.Version)(logCfgOpt)
	utils.WithCfgPrj(bc.CfgData.Project)(logCfgOpt)
	utils.WithCfgGroup(bc.CfgData.Group)(logCfgOpt)
	utils.WithCfgService(bc.CfgData.Service)(logCfgOpt)
	utils.WithCfgName(bc.CfgData.CfgName)(logCfgOpt)
	utils.WithCfgURL(bc.CfgData.CompanionUrl)(logCfgOpt)
	utils.WithCfgCB(bc.CfgData.CallBack)(logCfgOpt)
	utils.WithCfgCachePath(bc.CfgData.CachePath)(logCfgOpt)
	return utils.NewCfg(utils.CfgMode(bc.CfgMode), logCfgOpt)
}

type callserver struct {
	tool *ToolBox
	ui   UserInterface
	tm   TraceMeta
	opts *options
}

var getClientAddr = func(in context.Context) (addr string) {
	p, ok := peer.FromContext(in)
	if ok {
		return p.Addr.String()
	}
	return
}

func (c *callserver) Call(ctx context.Context, in *utils.ReqData) (*utils.ResData, error) {
	//meta := in.S.T
	meta := func() string {
		if in.S == nil {
			return ""
		}
		return in.S.T
	}()
	//当meta信息不合法时，不能生成合适的span，若为nil，则重新生成span
	span := utils.FromMeta(meta, c.tm.ip, int32(c.tm.port), c.tm.serviceName, utils.SrvSpan)
	if span == nil {
		span = utils.NewSpan(utils.SrvSpan)
	}
	//span = span.WithName("Call").Start()
	span = span.WithName(in.Op).Start()
	span = span.WithRpcCallType()
	if in.Op != xsf.LBOPGET && in.Op != xsf.LBOPSET && in.Op != PING {
		defer span.Flush()
	}
	defer span.End()

	peerAddr := getClientAddr(ctx)
	//将*utils.ReqData转换为*utils.Req
	inC := NewReqEx(in)
	inC.AppendSession(PEERADDR, peerAddr)
	//inC.SetParam(PEERADDR, peerAddr)
	sid := inC.Handle()
	if sid == "" {
		sid = c.tool.sid.generateSid()
		inC.SetHandle(sid)
	}

	var out = NewRes()
	var err error
	var start time.Time
	var dur int64
	abandon := false
	if c.tool.rateLimiter != nil {
		if !c.tool.rateLimiter.Allow() {
			abandon = true
			c.tool.Log.Errorw("request traffic exceeds limit")
			if c.opts.rateFallback != nil {
				start = time.Now()
				out, err = c.opts.rateFallback.Call(inC, span.Next(utils.SrvSpan))
				end := time.Now()
				dur = end.Sub(start).Nanoseconds()
			} else {
				err = rateLimiterErr
			}
		}
	}
	if (!abandon) && (PING != in.Op) {
		start = time.Now()
		var router1, router2, router3 int32 = 0, 0, 0
		if c.opts.router != nil {
			op, ok := c.opts.router.load(in.Op)
			if ok {
				atomic.AddInt32(&router1, 1)
				out, err = op(inC, span.Next(utils.SrvSpan), c.tool)
			} else {
				atomic.AddInt32(&router2, 1)
				out, err = c.ui.Call(inC, span.Next(utils.SrvSpan))
			}
		} else {
			atomic.AddInt32(&router3, 1)
			out, err = c.ui.Call(inC, span.Next(utils.SrvSpan))
		}
		end := time.Now()
		dur = end.Sub(start).Nanoseconds()
	}

	c.tool.Log.Infow("record call perf", "handle", sid, "cIp", peerAddr, "dur", dur)
	if err != nil {
		return nil, err
	}
	out.SetHandle(sid)
	out.SetTraceID(span.Meta())
	//将*utils.Res转换为*utils.ResData
	if !utils.IsNil(c.tool.vCpuManager) {
		out.SetParam(remainVCpus, strconv.FormatInt(c.tool.vCpuManager.RemainVCpuInt64, 10))
	}
	out.AppendSession(PEERADDR, c.tool.NetManager.getLocalIp())
	outC := out.Res()

	{
		//sync data to slidingWindow
		if !utils.IsNil(c.tool.errWin) {
			c.tool.errWin.setErrCode(int64(out.Res().GetCode()))
		}
		if !utils.IsNil(c.tool.delayWin) {
			c.tool.delayWin.setDur(dur)
		}
	}
	return outC, nil
}

func xrpcsRun(bc BootConfig, toolbox *ToolBox, srv UserInterface, opts *options) error {
	if err := srv.Init(toolbox); nil != err {
		return err
	}

	addKillerCheck(
		killerHighPriority,
		"srv.FInit",
		&killerWrapper{callback: func() {
			fInitErr := srv.FInit()
			if fInitErr != nil {
				toolbox.Log.Errorw("srv.FInit failed")
			}
		}},
	)

	loggerStd.Println("about to x.run")
	toolbox.Log.Infof("about to x.run\n")
	var x xsfServer
	if err := x.run(
		bc,
		toolbox.lis,
		&callserver{
			ui:   srv,
			tool: toolbox,
			tm: TraceMeta{
				ip:          toolbox.NetManager.GetIp(),
				port:        toolbox.NetManager.GetPort(),
				serviceName: svcsection,
			},
			opts: opts},
	); err != nil {
		return err
	}
	return nil
}
func xrpcsRunWrapper(
	bc BootConfig,
	toolbox *ToolBox,
	srv UserInterface,
	opts *options,
) error {
	return xrpcsRun(bc, toolbox, srv, opts)
}
