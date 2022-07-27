package _var

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"syscall"
)

type Flag struct {
	/*	CmdMode		= xsf.Mode				// -m
		CmdCfg		= xsf.Cfg				// -c
		CmdProject	= xsf.Project			// -p
		CmdGroup	= xsf.Group				// -g
		CmdService	= xsf.Service			// -s
		CmdCompanionUrl = xsf.CompanionUrl	// -u
	*/
	// default 缺省配置模式为native
	CmdCfg       *string
	XTestVersion *bool
}

func NewFlag() Flag {
	return Flag{
		CmdCfg:       flag.String("f", "xtest.toml", "client cfg name"),
		XTestVersion: flag.Bool("v", false, "xtest version"),
	}
}

func (f *Flag) Parse()  {
	flag.Parse()
	if *f.XTestVersion {
		fmt.Println("2.5.2")
		//os.Exit(0)
		syscall.Exit(0)
		return
	}
}
//var (
//	/*	CmdMode		= xsf.Mode				// -m
//		CmdCfg		= xsf.Cfg				// -c
//		CmdProject	= xsf.Project			// -p
//		CmdGroup	= xsf.Group				// -g
//		CmdService	= xsf.Service			// -s
//		CmdCompanionUrl = xsf.CompanionUrl	// -u
//	*/
//	// default 缺省配置模式为native
//	CmdCfg       = flag.String("f", "xtest.toml", "client cfg name")
//	XTestVersion = flag.Bool("v", false, "xtest version")
//)

func Usage() {
	fmt.Println("usage of common test tool")
	fmt.Println("-f		specify config file")
	os.Exit(0)
}

// Input jbzhou5 Input data
func Input(data string) (int, error) {
	in := bufio.NewReader(os.Stdin)
	fmt.Print("Please input data: ")
	n, err := fmt.Fscanln(in, &data)
	if err != nil {
		return 0, err
	}
	return n, nil
}
