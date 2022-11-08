package utils

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	xsf "github.com/xfyun/xsf/server"
	"os"
)

const defaultCfgFile = "aiges.toml"

var (
	//go:embed default.cfg
	defaultCfg []byte
)

type Flag struct {
	/*	CmdMode		= xsf.Mode				// -m
		CmdCfg		= xsf.Cfg				// -c
		CmdProject	= xsf.Project			// -p
		CmdGroup	= xsf.Group				// -g
		CmdService	= xsf.Service			// -s
		CmdCompanionUrl = xsf.CompanionUrl	// -u

	*/
	InitNativeCfg *bool
}

func NewFlag() Flag {
	return Flag{
		InitNativeCfg: flag.Bool("init", false, "init default config"),
	}
}

func InitDefaultCfg() error {
	_, err := os.Stat(defaultCfgFile)
	if err == nil {
		return nil
	}
	file, err := os.OpenFile(defaultCfgFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.Write(defaultCfg)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	return nil
}

func (f *Flag) Parse() {
	flag.Parse()
	if *f.InitNativeCfg {
		// 默认生成一个aiges.toml文件
		fmt.Println("Generating default cfg...")
		InitDefaultCfg()
	}
	if *xsf.Mode == 0 {
		if err := InitDefaultCfg(); err == nil {
			if *xsf.Cfg == "" {
				*xsf.Cfg = defaultCfgFile
			}
			if *xsf.Group == "" {
				*xsf.Group = "hu"
			}
			if *xsf.Project == "" {
				*xsf.Project = "AIpaas"
			}
			if *xsf.CompanionUrl == "" {
				*xsf.CompanionUrl = "http://companion.xfyun.iflytek:6868"
			}
			if *xsf.Service == "" {
				*xsf.Service = "svcName"
			}
		}
	}
}
