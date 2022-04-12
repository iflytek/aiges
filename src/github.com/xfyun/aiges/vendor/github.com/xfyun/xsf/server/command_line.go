package xsf

import (
	"flag"
)

var (
	Mode         = flag.Int("m", -1, "0、1、2 refer to Native、Centre、Custom")
	Cfg          = flag.String("c", "", "cfgName")
	DefaultCfg   = flag.String("dc", "", "default cfgName")
	Project      = flag.String("p", "", "Project")
	Group        = flag.String("g", "", "Group")
	Service      = flag.String("s", "", "Service")
	CompanionUrl = flag.String("u", "", "CompanionUrl")
)
