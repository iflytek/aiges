package catch

import (
	"fmt"
	"github.com/xfyun/xsf/utils"
	"os"
	"os/signal"
)

var (
	switchOn bool = false
	catchLog *utils.Logger
	sigChan  chan os.Signal
	dumpDir  string = "/log/loaderEvent"
)

type rpMeta struct {
	svc  string // 服务名
	ver  string // 加载器版本号
	pid  string // 进程号
	gid  string // 协程号
	addr string // 服务地址
	err  error  // 异常信息, panic error
}

type TagData struct {
	Typ  string            // 数据类型
	Desc map[string]string // 数据描述
	Data []byte            // 数据实体
}
type TagParam struct {
	Sid    string
	Header map[string]string
	Param  map[string]string
}
type TagRequest struct {
	TagParam
	DataList []TagData
}

func Open(flagOn bool, threshold int, log *utils.Logger, cb func() (req []TagRequest)) {
	switchOn = flagOn
	if switchOn {
		catchLog = log
		instMgrCallBack = cb
		sigChan = make(chan os.Signal, 1)
		err := os.MkdirAll(dumpDir, os.ModeDir)
		if err != nil {
			catchLog.Errorw("catchOpen make dump dir fail", "error", err.Error())
		}
		signalHandle()
		DeadLockDetectInit(threshold)
	}
	return
}

func Close() {
	if switchOn {
		DeadLockDetectFini()
		switchOn = false
		instMgrCallBack = nil
		signal.Stop(sigChan)
		close(sigChan)
		fmt.Println("aiService.Fini: fini catch success!")
	}
}
