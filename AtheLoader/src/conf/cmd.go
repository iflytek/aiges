package conf

import (
	"flag"
	xsf "git.xfyun.cn/AIaaS/xsf-external/server"
)

// 外部命令行参数及环境变量
var (
	CmdProject      = xsf.Project
	CmdGroup        = xsf.Group
	CmdService      = xsf.Service
	CmdCompanionUrl = xsf.CompanionUrl
	CmdVer          = flag.Bool("v", false, "aiges version info")
	CmdPrf          = flag.Bool("pprof", false, "pprof switch")
	CmdPrfAddr      = flag.String("prfAddr", "127.0.0.1:1234", "pprof ip: port")
)
