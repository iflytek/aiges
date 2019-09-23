package _var

import (
	"flag"
	"fmt"
	"os"
)

var (
	/*	CmdMode		= xsf.Mode				// -m
		CmdCfg		= xsf.Cfg				// -c
		CmdProject	= xsf.Project			// -p
		CmdGroup	= xsf.Group				// -g
		CmdService	= xsf.Service			// -s
		CmdCompanionUrl = xsf.CompanionUrl	// -u
	*/
	// default 缺省配置模式为native
	CmdCfg = flag.String("f", "xtest.toml", "client cfg name")
	//CmdVer	= flag.Bool("v", false, "print cli version")	// -v
	//CmdHelp = flag.Bool("h", false, "print cli usage")	// -h
)

func Usage() {
	fmt.Println("usage of common test tool")
	//fmt.Println("-h		print help information")
	//fmt.Println("-v		print tool version")
	fmt.Println("-f		specify config file")
	os.Exit(0)
}
