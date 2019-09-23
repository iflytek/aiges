package catch

import (
	"conf"
	xsf "git.xfyun.cn/AIaaS/xsf-external/server"
	"os"
	"runtime/debug"
	"strconv"
	"utils"
)

// go协程recover, for each goroutine;
func RecoverHandle(tag string) {
	if switchOn {
		if err := recover(); err != nil {
			meta := rpMeta{
				*conf.CmdService,
				conf.SvcVersion,
				strconv.Itoa(os.Getpid()),
				strconv.Itoa(utils.GetGoroutineID()),
				xsf.GetNetaddr(),
				tag,         // inst catch hook:
				err.(error), // panic error
			}
			catchFlush(meta, nil, debug.Stack())
		}
	}
}
