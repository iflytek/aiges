package catch

import (
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"os"
	"os/signal"
)

var (
	switchOn  bool = false
	dumpCatch bool = false
	catchLog  *utils.Logger
	sigChan   chan os.Signal
	dumpDir   string = "./catchDump"
)

type rpMeta struct {
	svc  string // 服务名
	ver  string // 加载器版本号
	pid  string // 进程号
	gid  string // 协程号
	addr string // 服务地址
	sid  string // 异常实例, 为空则flush全量会话信息
	err  error  // 异常信息, panic error
}

type TagData struct {
	Typ  string // 数据类型
	Fmt  string // 数据编码
	Enc  string // 数据压缩
	Data []byte // 数据实体
}
type TagParam struct {
	Sid   string
	Param map[string]string
}
type TagRequest struct {
	TagParam
	DataList []TagData
}

func Open(flagOn bool, flagDump bool, log *utils.Logger, cb func(tag string) (reqDoubt []TagRequest)) {
	switchOn = flagOn
	dumpCatch = flagDump
	catchLog = log
	instMgrCallBack = cb
	sigChan = make(chan os.Signal, 1)
	if dumpCatch {
		err := os.MkdirAll(dumpDir, os.ModeDir)
		if err != nil {
			catchLog.Errorw("catchOpen make dump dir fail", "error", err.Error())
		}
	}
	return
}

func Close() {
	if switchOn {
		switchOn = false
		dumpCatch = false
		instMgrCallBack = nil
		catchLog = nil
		signal.Stop(sigChan)
		close(sigChan)
	}
}
