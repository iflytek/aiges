package conf

// 框架配置项
const (
	// section AIGES
	sectionAiges             = "aiges"
	engSub                   = "sub"
	wrapperMode              = "asyncMode"
	sessGort                 = "sessDelGt"
	usrCfgName               = "usrCfg"
	realTimeRlt              = "realTimeRlt"
	httpRetry                = "httpRetry"
	grayMark                 = "gray"
	wrapperTrace             = "wrapperTrace"
	headerPass               = "headerPass"
	localPassFile            = "headerPass.wrapper"       // 平台header参数透传列表,透传至wrapper插件(按行读取)
	wrapperDelayDetectPeriod = "wrapperDelayDetectPeriod" //平台判断是否发生卡死的检测周期
	storageData              = "storageData"
	asyncRelease             = "asyncRelease"

	// grpc python bin path
	pythonPluginCmd = "pythonPluginCmd"

	// section flowControl
	sectionFc = "fc"
	maxLic    = "max"

	// section pprof
	sectionPProf   = "pprof"
	pprofAble      = "able"
	pporfHost      = "host"
	pprofPort      = "port"
	pprofSvcName   = "svcName"
	pprofProxyAddr = "proxyAddr"

	// wrapper cfg
	sectionWrapper = "wrapper"

	// section downAsync
	sectionAsyncDown = "downAsync"
	rabHost          = "rabHost"
	rabUser          = "rabUser"
	rabPass          = "rabPass"
	rabQueue         = "rabQueue"
	rabRetry         = "retry"
	nrtDBUrl         = "nrtDBUrl"

	// section dataProcess
	sectionDp = "dataProcess"
	dpRsAble  = "rsable" // reSample able

)

const (
	defaultEngSud string = "ase"
)
