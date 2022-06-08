package xsf

import (
	"context"
	"errors"
	"fmt"
	"github.com/xfyun/xsf/client"
	"github.com/xfyun/xsf/server/internal/bvt"
	"golang.org/x/time/rate"
	"time"

	"github.com/xfyun/xsf/utils"
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
	Monitor    *SonarAdapter
	NetManager *Net
	Bc         *BootConfig

	lis net.Listener

	errWin   *errCodeWindow
	delayWin *delayWindow

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
	//=================================================================
	//初始化configurator
	logCfgOpt := &utils.CfgOption{}
	utils.WithCfgDefault(bc.CfgData.CfgDefault)(logCfgOpt)
	utils.WithCfgVersion(bc.CfgData.Version)(logCfgOpt)
	utils.WithCfgPrj(bc.CfgData.Project)(logCfgOpt)
	utils.WithCfgGroup(bc.CfgData.Group)(logCfgOpt)
	utils.WithCfgService(bc.CfgData.Service)(logCfgOpt)
	utils.WithCfgName(bc.CfgData.CfgName)(logCfgOpt)
	utils.WithCfgURL(bc.CfgData.CompanionUrl)(logCfgOpt)
	utils.WithCfgCB(bc.CfgData.CallBack)(logCfgOpt)
	utils.WithCfgCachePath(bc.CfgData.CachePath)(logCfgOpt)
	cfg, err := utils.NewCfg(utils.CfgMode(bc.CfgMode), logCfgOpt)

	if err != nil {
		return fmt.Errorf("CreateConfiguratorErr:%v,logCfgOpt:%+v", err, logCfgOpt)
	}
	t.Cfg = cfg
	//=================================================================
	//读取日志相关配置，初始化日志
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

	t.Log, loggerErr = utils.NewLocalLog(utils.SetWash(logWash),
		utils.SetCaller(logCaller), utils.SetBatchSize(logBatchSize),
		utils.SetCacheMaxCount(logCacheMaxCount), utils.SetAsync(logasync),
		utils.SetLevel(logLevel), utils.SetFileName(fileName),
		utils.SetMaxSize(maxsize), utils.SetMaxBackups(maxBackups), utils.SetMaxAge(maxAge))
	if loggerErr != nil {
		return fmt.Errorf("loggerErr:%v", loggerErr)
	}
	addKillerCheck(killerLowPriority, "logger", t.Log)
	loggerStd.Printf("utils.NewLocalLog success. -> LOGLEVEL:%v, FILENAME:%v, MAXSIZE:%v, MAXBACKUPS:%v, MAXAGE:%v\n",
		logLevel, fileName, maxsize, maxBackups, maxAge)
	t.Log.Errorf("xsfVer:%v,service:%v", utils.GetVer(), bc.CfgData.Service)
	t.Log.Infof("utils.NewLocalLog success. -> LOGLEVEL:%v, FILENAME:%v, MAXSIZE:%v, MAXBACKUPS:%v, MAXAGE:%v\n",
		logLevel, fileName, maxsize, maxBackups, maxAge)
	//=================================================================
	//读取bootConfig中service字段作为配置中ip、port、sub的入口
	ip, _ := cfg.GetString(svcsection, IP_)
	netCard, _ := cfg.GetString(svcsection, NETCARD_)
	ip, ipCov := utils.Host2Ip(ip, netCard)
	if ipCov != nil {
		return fmt.Errorf("host2Ip:%v,ip:%v,netCard:%v", ipCov, ip, netCard)
	}
	port, portErr := cfg.GetInt(svcsection, PORT_)
	if portErr != nil {
		port = defaultPORT
	}

	reusePort, reusePortErr := t.Cfg.GetInt(svcsection, REUSEPORT_)
	if reusePortErr != nil {
		reusePort = defaultREUSEPORT
	}
	t.lis, err = NewListener(reusePort, net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return errors.New(fmt.Sprintf("can't listen %v", ip) + func() string {
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
		return fmt.Errorf("can't get ip and port from %v", t.lis.Addr().String())
	}
	port, err = strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("can't convert %v to int,err->%v", portStr, err)
	}
	loggerStd.Printf("host2ip->ip:%v,port:%v\n", ip, port)
	t.Log.Infof("host2ip->ip:%v,port:%v\n", ip, port)
	t.NetManager = &Net{portInt: port, ipStr: ip}

	GRPCTIMEOUT = defaultGRPCTIMEOUT
	grpctimeout, grpctimeoutErr := cfg.GetInt(svcsection, GRPCTIMEOUT_)
	if grpctimeoutErr == nil {
		GRPCTIMEOUT = grpctimeout

	}

	finderSwitch, finderSwitchErr := cfg.GetInt(svcsection, FINDERSWITCH_)
	loggerStd.Printf("finderSwitch:%v,finderSwitchErr:%v\n", finderSwitch, finderSwitchErr)
	t.Log.Infof("finderSwitch:%v,finderSwitchErr:%v\n", finderSwitch, finderSwitchErr)
	if finderSwitchErr != nil {
		finderSwitch = defaultFINDER
	}

	//=================================================================
	maxreceiveInt, maxreceiveIntErr := cfg.GetInt(svcsection, MAXRECEIVE)
	if maxreceiveIntErr != nil {
		maxreceiveInt = defaultMAXRECEIVE
	} else {
		maxreceiveInt = maxreceiveInt * 1024 * 1024
	}
	grpcOptInst.setMaxReceiveMessageSize(maxreceiveInt)

	maxsendInt, maxsendIntErr := cfg.GetInt(svcsection, MAXSEND)
	if maxsendIntErr != nil {
		maxsendInt = defaultMAXSEND
	} else {
		maxsendInt = maxsendInt * 1024 * 1024
	}
	grpcOptInst.setMaxSendMessageSize(maxsendInt)

	conrbufInt, conrbufIntErr := cfg.GetInt(svcsection, CONRBUF)
	if conrbufIntErr != nil {
		conrbufInt = defaultCONRBUF
	} else {
		conrbufInt = conrbufInt * 1024 * 1024
	}
	grpcOptInst.setReadBufferSize(conrbufInt)

	conwbufInt, conwbufIntErr := cfg.GetInt(svcsection, CONWBUF)
	if conwbufIntErr != nil {
		conwbufInt = defaultCONWBUF
	} else {
		conwbufInt = conwbufInt * 1024 * 1024
	}
	grpcOptInst.setWriteBufferSize(conwbufInt)

	//=================================================================
	//初始化loadReporter
	//读取loadReporter的配置
	lbClient := LbAdapter{}
	if cfg.GetSection(LOADREPORTER) != nil {
		lbAble, lbAbleErr := cfg.GetInt(LOADREPORTER, LBABLE)
		if lbAbleErr != nil {
			lbAble = defaultABLE
		}
		lbLbStrategy, lbLbStrategyErr := cfg.GetInt(LOADREPORTER, LBLBSTRATEGY)
		lbZkList, lbZkListErr := interface2stringslice(cfg.GetAsObject(LOADREPORTER, LBZKLIST))
		lbRoot, lbRootErr := cfg.GetString(LOADREPORTER, LBROOT)
		if lbRootErr != nil {
			lbRoot = ""
		}
		lbRouterType, lbRouterTypeErr := cfg.GetString(LOADREPORTER, LBROUTERTYPE)
		lbSubRouterTypes, lbSubRouterTypesErr := interface2stringslice(cfg.GetAsObject(LOADREPORTER, LBSUBROUTERTYPES))
		if lbSubRouterTypesErr != nil {
			lbSubRouterTypes = nil
		}
		lbRedisHost, lbRedisHostErr := cfg.GetString(LOADREPORTER, LBREDIEHOST)
		lbRedisPasswd, lbRedisPasswdErr := cfg.GetString(LOADREPORTER, LBREDISPASSWD)
		if lbRedisPasswdErr != nil {
			lbRedisPasswd = ""
		}
		lbMaxActive, lbMaxActiveErr := cfg.GetInt(LOADREPORTER, LBMAXACTIVE)
		if lbMaxActiveErr != nil {
			lbMaxActive = 0
		}
		lbMaxIdle, lbMaxIdleErr := cfg.GetInt(LOADREPORTER, LBMAXIDLE)
		if lbMaxIdleErr != nil {
			lbMaxIdle = 0
		}
		lbDb, lbDbErr := cfg.GetInt(LOADREPORTER, LBDB)
		if lbDbErr != nil {
			lbDb = 0
		}
		lbIdleTimeout, lbIdleTimeoutErr := cfg.GetInt(LOADREPORTER, LBIDLETIMEOUT)
		if lbIdleTimeoutErr != nil {
			lbIdleTimeout = 0
		}
		if lbAble != defaultLBABLE {
			if lbLbStrategyErr != nil || lbZkListErr != nil || lbRouterTypeErr != nil || lbRedisHostErr != nil {
				return fmt.Errorf("lbLbStrategyErr:%v,lbZkListErr:%v,lbRouterTypeErr:%v,lbredishostErr:%v", lbLbStrategyErr, lbZkListErr, lbRouterTypeErr, lbRedisHostErr)
			}
		}
		lbClient = LbAdapter{able: lbAble}
		if InitErr := lbClient.Init(
			WithLbAdapterSvc(fmt.Sprintf("%v:%v", ip, port)),
			WithLbAdapterStrategy(lbLbStrategy),
			WithLbAdapterZkList(lbZkList),
			WithLbAdapterRoot(lbRoot),
			WithLbAdapterRouterType(lbRouterType),
			WithLbAdapterSubRouterTypes(lbSubRouterTypes),
			WithLbAdapterSRedisHost(lbRedisHost),
			WithLbAdapterSRedisPasswd(lbRedisPasswd),
			WithLbAdapterMaxActive(lbMaxActive),
			WithLbAdapterMaxIdle(lbMaxIdle),
			WithLbAdapterDb(lbDb),
			WithLbAdapterIdleTimeOut(time.Second*time.Duration(lbIdleTimeout))); InitErr != nil {
			return fmt.Errorf("InitErr:%v\n", InitErr)
		}
		if lbAble == defaultLBABLE {
			loggerStd.Printf("lb is not enabled\n")
			t.Log.Infof("lb is not enabled\n")
		} else {
			loggerStd.Printf("lbClient.Init success. -> lbClient:%+v\n", lbClient)
			t.Log.Infof("lbClient.Init success. -> lbClient:%+v\n", lbClient)
		}

	} else {
		lbClient = LbAdapter{able: defaultLBABLE}
	}
	//=================================================================
	//初始化hermesAdapter
	loggerStd.Println("about to deal with hermes.")
	lbv2Client, lbv2ClientErr := func() (hermes hermesAdapter, hermesErr error) {
		if cfg.GetSection(HERMES) != nil {
			lbAbleInt, lbAbleErr := cfg.GetInt(HERMES, HERMESABLE)
			lbAble := defaultHERMESABLE
			if lbAbleErr == nil {
				if lbAbleInt != 0 {
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
			if lbv2SvcErr != nil {
				hermesErr = fmt.Errorf("can't take %v from %v,err:%v", HERMESSVC, HERMES, lbv2SvcErr)
				return
			}
			lbv2SubSvc, lbv2SubSvcErr := cfg.GetString(HERMES, HERMESSUBSVC)
			if lbv2SubSvcErr != nil {
				hermesErr = lbv2SubSvcErr
			}

			lbname, lbnameErr := cfg.GetString(HERMES, HERMESLBNAME)
			if lbnameErr != nil {
				hermesErr = lbnameErr
			}
			lbprj, _ := cfg.GetString(HERMES, HERMESLBPROJECT)
			lbgro, _ := cfg.GetString(HERMES, HERMESLBGROUP)

			apiversion := bc.CfgData.ApiVersion
			apiVersionStr, apiVersionStrErr := cfg.GetString(HERMES, HERMESAPIVERSION)
			if apiVersionStrErr == nil {
				apiversion = apiVersionStr
			}

			hermesTask := defaultHERMESTASK
			hermesTaskInt64, hermesTaskInt64Err := cfg.GetInt64(HERMES, HERMESTASK)
			if hermesTaskInt64Err == nil {
				hermesTask = int(hermesTaskInt64)
			}

			finderTtl := defaultHERMESFINDERTTL
			finderTtlInt, finderTtlErr := cfg.GetInt(HERMES, HERMESFINDERTTL)
			if finderTtlErr == nil {
				finderTtl = time.Millisecond * time.Duration(finderTtlInt)
			}

			backendInt, backendIntErr := cfg.GetInt(HERMES, HERMESBACKEND)
			if backendIntErr != nil {
				backendInt = defaultHERMESBACKEND
			}

			HermesTimeout := defaultHERMESTIMEOUT
			HermesTimeoutInt, HermesTimeoutIntErr := cfg.GetInt(HERMES, HERMESTIMEOUT)
			if HermesTimeoutIntErr == nil {
				HermesTimeout = time.Duration(HermesTimeoutInt) * time.Millisecond
			}

			HermesMode := bc.CfgMode
			HermesModeInt, HermesModeIntErr := cfg.GetInt(HERMES, HERMESFINDERMODE)
			if HermesModeIntErr == nil {
				HermesTimeout = time.Duration(HermesTimeoutInt) * time.Millisecond
				HermesMode = utils.CfgMode(HermesModeInt)
			}

			cloud := defaultHERMESCLOUD
			cloudString, cloudStringErr := cfg.GetString(HERMES, HERMESLBCLOUD)
			if cloudStringErr == nil {
				cloud = cloudString
			}

			if InitErr := hermes.Init(
				WithHermesAdapterMode(HermesMode),
				WithHermesAdapterLbApiVersion(apiversion),
				WithHermesAdapterLbName(lbname),
				WithHermesAdapterLbPrj(lbprj),
				WithHermesAdapterLbGro(lbgro),
				WithHermesAdapterAddr(fmt.Sprintf("%v:%v", t.NetManager.ipStr, t.NetManager.portInt)),
				WithHermesAdapterSvcAndSubSvc(lbv2Svc, lbv2SubSvc),
				WithHermesAdapterSubsvc(lbv2SubSvc),
				WithHermesAdapterFinderTtl(finderTtl),
				WithHermesAdapterBackEnd(backendInt),
				WithHermesAdapterTimeout(HermesTimeout),
				WithHermesAdapterBootConfig(bc),
				WithHermesAdapterSvcIp(ip),
				WithHermesAdapterCloudId(cloud),
				WithHermesAdapterTask(hermesTask),
				WithHermesAdapterSvcPort(int32(port))); InitErr != nil {
				hermesErr = fmt.Errorf("InitErr:%v\n", InitErr)
				return
			}
			loggerStd.Printf("hermes init success with mode:%v.", HermesMode.String())
		} else {
			hermes.able = false
			loggerStd.Println("hermes is disable.")
		}
		return
	}()
	time.Sleep(time.Second)
	if lbv2ClientErr != nil {
		return fmt.Errorf("lbv2Client init fail. err:%v\n", lbv2ClientErr)
	}

	//=================================================================
	//读取flowControl的配置
	fcAble, fcAbleErr := cfg.GetInt(FLOWCONTROL, FCABLE)
	if fcAbleErr != nil {
		fcAble = defaultFcAble
	}

	var licMax int

	if fcAble != 0 {
		router, routerErr := cfg.GetString(FLOWCONTROL, ROUTER)
		max, maxErr := cfg.GetInt(FLOWCONTROL, MAX)
		licMax = max
		ttl, ttlErr := cfg.GetInt(FLOWCONTROL, TTL)
		if ttlErr != nil {
			ttl = defaultTTL
		}
		best, bestErr := cfg.GetInt(FLOWCONTROL, BEST)
		if bestErr != nil {
			best = max
		}
		wave, waveErr := cfg.GetInt(FLOWCONTROL, WAVE)
		if waveErr != nil {
			wave = defaultWAVE
		}
		strategy, strategyErr := cfg.GetInt(FLOWCONTROL, STRATEGY)
		if strategyErr != nil {
			strategy = defaultSTRATEGY
		}
		roll, rollErr := cfg.GetInt(FLOWCONTROL, ROLLTIMEOUT)
		rollTime := time.Duration(roll)
		if rollErr != nil {
			rollTime = defaultROLLTIMEOUT
		}
		report, reportErr := cfg.GetInt(FLOWCONTROL, REPORT)
		reportInterval := time.Duration(report)
		if reportErr != nil {
			reportInterval = defaultREPORT
		}
		taskSize, taskSizeErr := cfg.GetInt(FLOWCONTROL, TASKSIZE)
		if taskSizeErr != nil {
			taskSize = defaultTASKSIZE
		}
		taskChannelSize, taskChannelSizeErr := cfg.GetInt(FLOWCONTROL, TASKCHANNELSIZE)
		if taskChannelSizeErr != nil {
			taskChannelSize = defaultTASKCHANNELSIZE
		}
		if routerErr != nil || maxErr != nil {
			return errors.New(fmt.Sprintf("can't get the router、max from configurator -> routerErr:%v, maxErr:%v", routerErr, maxErr))
		}

		authFilterWin := defaultAUTHFILTERWIN
		authFilterWinInt, authFilterWinIntErr := cfg.GetInt(FLOWCONTROL, AUTHFILTERWIN)
		if authFilterWinIntErr == nil {
			authFilterWin = time.Duration(authFilterWinInt) * time.Millisecond
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
				WithSessionManagerWave(wave),
				WithSessionManagerAuthFilterWin(authFilterWin))
			if SessionManagerErr != nil {
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
				WithQpsLimiterReporterv2(lbv2Client),
				WithQpsLimiterLogger(t.Log))
			if QpsLimiterErr != nil {
				return fmt.Errorf("QpsLimiterErr:%v", QpsLimiterErr)
			}
			loggerStd.Printf("NewQpsLimiter success.\n")
			t.Log.Infof("NewQpsLimiter success.\n")
		}
	}
	//=================================================================
	//初始化rpcsidgenerator.go
	var sidErr error
	t.sid, sidErr = NewSidGenerator(sidVer, ip, int64(port))
	if sidErr != nil {
		return sidErr
	}
	loggerStd.Printf("NewSidGenerator success.\n")
	t.Log.Infof("NewSidGenerator success.\n")

	//=================================================================
	//初始化trace
	traceHost, traceHostErr := cfg.GetString(TRACE, TRACEHOST)
	if traceHostErr != nil {
		traceHost = defaultTRACEHOST
	}
	tracePort, tracePortErr := cfg.GetInt(TRACE, TRACEPORT)
	if tracePortErr != nil {
		tracePort = defaultTRACEPORT
	}
	backend, backendErr := cfg.GetInt(TRACE, BACKEND)
	if backendErr != nil {
		backend = defaultBACKEND
	}
	deliver_, deliverErr := cfg.GetInt(TRACE, DELIVER)
	if deliverErr != nil {
		deliver_ = defaultDELIVER
	}
	dump_, dumpErr := cfg.GetInt(TRACE, DUMP)
	if dumpErr != nil {
		dump_ = defaultDUMP
	}
	able_, ableErr := cfg.GetInt(TRACE, ABLE)
	if ableErr != nil {
		able_ = defaultABLE
	}
	spill, spillErr := cfg.GetString(TRACE, SPILL)
	if spillErr != nil {
		spill = defaultSPILL
	}
	buffer, bufferErr := cfg.GetInt(TRACE, BUFFER)
	if bufferErr != nil {
		buffer = defaultBUFFER
	}
	batch, batchErr := cfg.GetInt(TRACE, BATCH)
	if batchErr != nil {
		batch = defaultBATCH
	}
	linger, lingerErr := cfg.GetInt(TRACE, LINGER)
	if lingerErr != nil {
		linger = defaultLINGER
	}
	watchBool := defaultWATCH
	watchInt, watchErr := cfg.GetInt(TRACE, WATCH)
	if watchErr == nil {
		if watchInt == 1 {
			watchBool = true
		}
	}

	bcluster := defaultTRACEBCLUSTER
	bclusterStr, bclusterStrErr := cfg.GetString(TRACE, TRACEBCLUSTER)
	if bclusterStrErr == nil {
		bcluster = bclusterStr
	}

	idc := defaultTRACEIDC
	idcStr, idcStrErr := cfg.GetString(TRACE, TRACEIDC)
	if idcStrErr == nil {
		idc = idcStr
	}
	watchPort := defaultWatchPort
	watchPortInt, watchPortErr := cfg.GetInt(TRACE, WATCHPORT)
	if watchPortErr == nil {
		watchPort = watchPortInt
	}

	spillSize := defaultSpillSize
	spillSizeInt, spillSizeErr := cfg.GetInt(TRACE, SPILLSIZE)
	if spillSizeErr == nil {
		spillSize = spillSizeInt
	}

	loadTs := defaultLoadTs
	loadTsInt, loadTsErr := cfg.GetInt(TRACE, LOADTS)
	if loadTsErr == nil {
		loadTs = loadTsInt
	}

	deliver := false
	dump := false
	able := false
	if deliver_ == 1 {
		deliver = true
	}
	if dump_ == 1 {
		dump = true
	}
	if able_ == 1 {
		able = true
	}
	if able_ != defaultUNCHANGE {
		utils.AbleTrace(able)
		if able {
			if traceErr := utils.InitTracer(
				traceHost,
				strconv.Itoa(tracePort),
				utils.WithLowLoadSleepTs(loadTs),
				utils.WithMaxSpillContentSize(spillSize),
				utils.WithWatchPort(watchPort),
				utils.WithWatch(watchBool),
				utils.WithSvcName(bc.CfgData.Service),
				utils.WithSvcPort(int32(port)),
				utils.WithSvcIp(ip),
				utils.WithBufferSize(buffer),
				utils.WithBatchSize(batch),
				utils.WithLinger(linger),
				utils.WithTraceSpill(spill),
				utils.WithBackend(backend),
				utils.WithDeliver(deliver),
				utils.WithDump(dump),
				utils.WithSvcBCluster(bcluster),
				utils.WithSvcIDC(idc),
				utils.WithTraceLogger(t.Log)); traceErr != nil {
				return fmt.Errorf("InitTracer failed -> able:%v,ip:%v,port:%v,backend:%v,deliver:%v,dump:%v -> traceErr:%v", able, traceHost, tracePort, backend, deliver, dump, traceErr)
			}
		}
	}

	//=================================================================
	//初始化sonar
	sonarHost, sonarHostErr := cfg.GetString(SONAR, SONARHOST)
	if sonarHostErr != nil {
		sonarHost = defaultSONARHOST
	}
	sonarPort, sonarPortErr := cfg.GetInt(SONAR, SONARPORT)
	if sonarPortErr != nil {
		sonarPort = defaultSONARPORT
	}
	sonarBackend, sonarBackendErr := cfg.GetInt(SONAR, SONARBACKEND)
	if sonarBackendErr != nil {
		sonarBackend = defaultSONARBACKEND
	}
	sonarDeliver, sonarDeliverErr := cfg.GetInt(SONAR, SONARDELIVER)
	if sonarDeliverErr != nil {
		sonarDeliver = defaultSONARDELIVER
	}
	sonarDump, sonarDumpErr := cfg.GetInt(SONAR, SONARDUMP)
	if sonarDumpErr != nil {
		sonarDump = defaultSONARDUMP
	}
	sonarAble, sonarAbleErr := cfg.GetInt(SONAR, SONARABLE)
	if sonarAbleErr != nil {
		sonarAble = defaultSONARABLE
	}
	sonarDS, sonarDSErr := cfg.GetString(SONAR, SONARDS)
	if sonarDSErr != nil {
		sonarDS = defaultSONARDS
	}
	sonardeliver := false
	sonardump := false
	sonarable := false
	if sonarDeliver != 0 {
		sonardeliver = true
	}
	if sonarDump != 0 {
		sonardump = true
	}
	if sonarAble != 0 {
		sonarable = true
	}
	t.Monitor = &SonarAdapter{}
	sonarErr := t.Monitor.initSonar(
		WithSonarAdapterAble(sonarable),
		WithSonarAdapterDs(sonarDS),
		WithSonarAdapterMetricEndpoint(ip),
		WithSonarAdapterMetricServiceName(bc.CfgData.Service),
		WithSonarAdapterMetricPort(strconv.Itoa(port)),
		WithSonarAdapterLogger(nil),
		WithSonarAdapterSonarDumpEnable(sonardump),
		WithSonarAdapterSonarDeliverEnable(sonardeliver),
		WithSonarAdapterSonarHost(sonarHost),
		WithSonarAdapterSonarPort(strconv.Itoa(sonarPort)),
		WithSonarAdapterSonarBackend(sonarBackend))
	if sonarable {
		if sonarErr != nil {
			return fmt.Errorf("sonarErr:%v", sonarErr)
		}
		loggerStd.Printf("sonar init success.\n")
		t.Log.Infof("sonar init success.\n")
	}
	//=================================================================
	//初始化finder
	loggerStd.Printf("about to deal finder.\n")
	if finderSwitch != 0 {
		loggerStd.Printf("about to new finder.\n")
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

		if finderErr != nil {
			loggerStd.Printf("CreateFinder fail -> bc:%+v, finderErr:%v", bc, finderErr)
			return fmt.Errorf("CreateFinder fail -> bc:%+v, finderErr:%v", bc, finderErr)
		}
		loggerStd.Printf("CreateFinder success.\n")
		t.Log.Infof("CreateFinder success.\n")

		loggerStd.Printf("about to call finderadapter.AddRegister. addr:%s\n", fmt.Sprintf("%v:%v", ip, port))
		t.Log.Infof("about to call finderadapter.AddRegister. addr:%s\n", fmt.Sprintf("%v:%v", ip, port))

		finderadapter.AddRegister(bc.CfgData.ApiVersion, fmt.Sprintf("%v:%v", ip, port), finder)
	}
	//=================================================================
	loggerStd.Printf("about to deal metrics.\n")
	//初始化metrics
	if func() bool {
		metricsAble := defaultMETRICSABLE
		metricsAbleInt, metricsAbleIntErr := cfg.GetInt(METRICS, METRICSABLE)
		if nil == metricsAbleIntErr {
			if 1 == metricsAbleInt {
				metricsAble = true
			}
		}
		return metricsAble
	}() {
		metricsIdc, metricsIdcErr := cfg.GetString(METRICS, METRICSIDC)
		if nil != metricsIdcErr || 0 == len(metricsIdc) {
			metricsIdc = defaultMETRICSIDC
		}
		MetricsIdc = metricsIdc

		metricsSub, metricsSubErr := cfg.GetString(METRICS, METRICSSUB)
		if nil != metricsSubErr || 0 == len(metricsSub) {
			metricsSub = defaultMETRICSSUB
		}
		MetricsSub = metricsSub

		metricsCs, metricsCsErr := cfg.GetString(METRICS, METRICSCS)
		if nil != metricsCsErr || 0 == len(metricsCs) {
			metricsCs = defaultMETRICSCS
		}
		MetricsCluster = metricsCs
		loggerStd.Println("begin to init registryInst")
		metricsErr := registryInst.initEx(metricsSub, svcsection, metricsIdc, metricsCs)
		if metricsErr != nil {
			return metricsErr
		}

		//初始化slidingWindow 2019-04-11 17:19:49
		loggerStd.Println("begin to init sliding window")
		var timePerSlice = defaultMETRICSTIMEPERSLICE
		timePerSliceInt, timePerSliceIntErr := cfg.GetInt64(METRICS, METRICSTIMEPERSLICE)
		if timePerSliceIntErr == nil {
			timePerSlice = time.Duration(timePerSliceInt) * time.Millisecond
		}

		winSize, winSizeErr := cfg.GetInt64(METRICS, METRICSWINSIZE)
		if nil != winSizeErr {
			winSize = defaultMETRICSWINSIZE
		}

		loggerStd.Printf("begin to init slidingWindow,timePerSlice:%v,winSize:%v\n", timePerSlice, winSize)

		t.delayWin = newDelayWindow(timePerSlice, winSize)
		t.errWin = newErrCodeWindow(timePerSlice, winSize)

		AddSlidingDelayWindows(t.delayWin)
		AddSlidingErrCodeWindows(t.errWin)
	} else {
		loggerStd.Println("metrics is disable")
	}

	//=================================================================
	//初始化rateLimiter
	loggerStd.Printf("about to deal rateLimiter.\n")
	rateInt, _ := cfg.GetInt(svcsection, rateLimiterRate)
	burstInt, _ := cfg.GetInt(svcsection, rateLimiterBurst)
	if 0 != rateInt || 0 != burstInt {
		loggerStd.Printf("rateLimiter,rate:%v,burst:%v\n", rateInt, burstInt)
		t.rateLimiter = rate.NewLimiter(rate.Every(time.Duration(rateInt)*time.Millisecond), burstInt)
	}

	//=================================================================
	//初始化vCpuManager
	loggerStd.Printf("about to deal vCpuManager.\n")
	vCpuAble := defaultVCPUABLE
	vCpuAbleInt, e := cfg.GetInt(svcsection, VCPUABLE)
	if e == nil {
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
			utils.CfgMode(bc.CfgMode),
		)
		if vCpuMapErr != nil {
			t.Log.Errorf("fetch vManager cfg failed:%v,ver:%v,prj:%v,group:%v,service:%v,url:%v,file:%v,mode:%v\n",
				vCpuMapErr, vCpuVersion, bc.CfgData.Project, vCpuGroup, vCpuService, bc.CfgData.CompanionUrl, VCPUFILENAME, bc.CfgMode)
		}
		loggerStd.Printf("vCpuManager:%+v\n", vCpuMap)
		t.vCpuManager, e = NewVCpuManager(vCpuInterval, vCpuMap)
		if err != nil {
			return fmt.Errorf("init vCpuManager failed,err:%v", e)
		} else {
			loggerStd.Println("vCpuManager is disable")
		}
	}

	//=================================================================
	//初始化 bvtVerifier
	loggerStd.Printf("about to deal bvtVerifier.\n")
	bvtAbleInt, bvtAbleIntErr := cfg.GetInt(BVT, BVTABLE)
	bvtService, bvtServiceErr := cfg.GetString(BVT, BVTSERVICE)
	bvtVersion, bvtVersionErr := cfg.GetString(BVT, BVTVERSION)
	bvtCfgFile, bvtCfgFileErr := cfg.GetString(BVT, BVTCFGFILE)
	bvtTimeoutInt, bvtTimeoutIntErr := cfg.GetInt(BVT, BVTTIMEOUT)
	bvtPlatform, bvtPlatformErr := cfg.GetString(BVT, BVTPLATFORM)
	bvtId, bvtIdErr := cfg.GetString(BVT, BVTID)
	bvtServiceAddress, bvtServiceAddressErr := cfg.GetString(BVT, BVTSERVICEADDRESS)
	bvtCallback, bvtCallbackErr := cfg.GetString(BVT, BVTCALLBACK)
	bvtAsyncInt, bvtAsyncIntErr := cfg.GetInt(BVT, BVTASYNC)

	bvtNamespace, bvtNamespaceErr := cfg.GetString(BVT, BVTNAMESPACE)
	if bvtNamespaceErr != nil {
		loggerStd.Println("namespace not set, use default")
		bvtNamespace = "default"
	}

	bvtAble := defBVTABLE
	if bvtAbleIntErr == nil {
		if bvtAbleInt == 0 {
			bvtAble = false
		}
		if bvtAbleInt == 1 {
			bvtAble = true
		}
	}

	if !bvtAble {
		loggerStd.Println("bvt is disable")
		return nil
	}

	//检查相关配置读取是否ok
	if bvtServiceErr != nil ||
		bvtVersionErr != nil ||
		bvtCfgFileErr != nil ||
		bvtTimeoutIntErr != nil ||
		bvtPlatformErr != nil ||
		bvtIdErr != nil ||
		bvtCallbackErr != nil ||
		bvtAsyncIntErr != nil ||
		bvtServiceAddressErr != nil {
		loggerStd.Printf("bvt configuration error,"+
			"bvtServiceErr:%v,"+
			"bvtVersionErr:%v,"+
			"bvtCfgFileErr:%v,"+
			"bvtTimeoutIntErr:%v,"+
			"bvtPlatformErr:%v,"+
			"bvtIdErr:%v,"+
			"bvtCallbackErr:%v,"+
			"bvtAsyncIntErr:%v,"+
			"bvtServiceAddressErr:%v\n",
			bvtServiceErr,
			bvtVersionErr,
			bvtCfgFileErr,
			bvtTimeoutIntErr,
			bvtPlatformErr,
			bvtIdErr,
			bvtCallbackErr,
			bvtAsyncIntErr,
			bvtServiceAddressErr)
		return nil
	}

	bvtAsync := false
	if bvtAsyncInt == 1 {
		bvtAsync = true
	}
	bvtTimeout := time.Millisecond * time.Duration(bvtTimeoutInt)

	loggerStd.Println("begin to init bvtVerifier")
	loggerStd.Printf(
		"bvtGProject:%v,"+
			"bvtGGroup:%v,"+
			"bvtGService:%v,"+
			"bvtGVersion:%v,"+
			"bvtGCfgFile:%v,"+
			"bvtGCompanionUrl:%v,"+
			"bvtTimeout:%v,"+
			"bvtPlatform:%v,"+
			"bvtId:%v,"+
			"bvtEngAddr:%v,"+
			"bvtCallback:%v,"+
			"bvtServiceAddress:%v\n",
		bc.CfgData.Project,
		bc.CfgData.Group,
		bvtService,
		bvtVersion,
		bvtCfgFile,
		bc.CfgData.CompanionUrl,
		bvtTimeout,
		bvtPlatform,
		bvtId,
		t.lis.Addr().String(),
		bvtCallback,
		bvtServiceAddress)

	bvt.Init(
		bc.CfgData.Project,
		bc.CfgData.Group,
		bvtService,
		bvtVersion,
		bvtCfgFile,
		bc.CfgData.CompanionUrl,

		bvtTimeout,
		bvtPlatform,
		bvtId,
		t.lis.Addr().String(),
		bvtCallback,
		bvtAsync,
		bvtServiceAddress,

		licMax,
		bc.CfgData.Service,
		bvtNamespace,
	)

	return nil
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
	incrConcurrentStatistics()
	defer decrConcurrentStatistics()

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

	utils.SyncUserInitStatus()

	if err := srv.Init(toolbox); err != nil {
		return err
	}

	addKillerCheck(killerHighPriority, "srv.Finit", &killerWrapper{callback: func() {
		finitErr := srv.Finit()
		if finitErr != nil {
			toolbox.Log.Errorw("srv.Finit failed")
		}
	}})

	loggerStd.Println("about to x.run")
	toolbox.Log.Infof("about to x.run\n")
	var x xsfServer
	if err := x.run(bc, toolbox.lis, &callserver{ui: srv, tool: toolbox, tm: TraceMeta{ip: toolbox.NetManager.GetIp(), port: toolbox.NetManager.GetPort(), serviceName: svcsection}, opts: opts}); err != nil {
		return err
	}
	return nil
}
func XrpcsRun(bc BootConfig, toolbox *ToolBox, srv UserInterface, opts *options) error {
	return xrpcsRun(bc, toolbox, srv, opts)
}
