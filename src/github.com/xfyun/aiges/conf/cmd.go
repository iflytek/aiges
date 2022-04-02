package conf

import (
	"flag"
	xsf "github.com/xfyun/xsf/server"
)

// 外部命令行参数及环境变量
var (
	CmdMode         = xsf.Mode
	CmdCfg          = xsf.Cfg
	CmdProject      = xsf.Project
	CmdGroup        = xsf.Group
	CmdService      = xsf.Service
	CmdCompanionUrl = xsf.CompanionUrl
	CmdVer          = flag.Bool("v", false, "aiges version info")
)
