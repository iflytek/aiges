package conf

// 框架配置项
const (
	// section AIGES
	sectionAiges = "aiges"
	engSub       = "sub"
	sessMode     = "sessMode"
	wrapperMode  = "asyncMode"
	numaNode     = "numaNode"
	sessGort     = "sessDelGt"
	gesMock      = "gesMock"
	usrCfgName   = "usrCfg"
	realTimeRlt  = "realTimeRlt"
	catchSwitch  = "catchOn"
	catchDump    = "catchDump"

	// section flowControl
	sectionFc = "fc"
	maxLic    = "max"

	// wrapper cfg
	sectionWrapper = "wrapper"

	// section downAsync
	sectionAsyncDown = "downAsync"
	rabHost			 = "rabHost"
	rabUser			 = "rabUser"
	rabPass			 = "rabPass"
	rabQueue		 = "rabQueue"
	nrtDBUrl		 = "nrtDBUrl"
)

const (
	defaultNumaNode int  = -1
	defaultEngSud string = "sub"
)
