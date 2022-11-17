package main

import (
	"fmt"
	_ "github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/env"
	"github.com/xfyun/aiges/service"
	"github.com/xfyun/aiges/utils"
	"github.com/xfyun/aiges/widget"
	xsfUtil "github.com/xfyun/xsf/utils"
	"os"
)

func main() {
	flg := utils.NewFlag()
	env.Parse()
	flg.Parse()

	//profiler.Start(profiler.Config{
	//	ApplicationName: "AISERVICE",
	//	ServerAddress:   "http://172.31.98.182:44040",
	//})
	if len(os.Args) < 2 {
		usage()
	}
	if *conf.CmdVer {
		fmt.Println(service.VERSION)
		return
	}

	var err error
	if env.SYSArch == "linux" {
		// 设置cpu亲和性
		if err = utils.NumaBind(env.AIGES_ENV_NUAME); err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	//
	var ch = &utils.Coordinator{
		ConfChan: make(chan *xsfUtil.Configure),
		Ch2:      make(chan int),
	}

	var aisrv service.EngService
	widgetInst := widget.NewWidget(env.AIGES_PLUGIN_MODE, ch)
	// 控件初始化&逆初始化
	if err = widgetInst.Open(); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer widgetInst.Close()

	// 框架初始化&逆初始化
	if err = aisrv.Init(env.AIGES_ENV_VERSION); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer aisrv.Fini()

	// 注册行为
	if err = widgetInst.Register(&aisrv); err != nil {
		fmt.Println(err.Error())
		return
	}

	// 框架运行
	if err = aisrv.Run(ch); err != nil {
		fmt.Println(err.Error())
		return
	}
	return

}

func usage() {
	fmt.Printf("加载器运行方法:\n" +
		"- 本地模式运行\n" +
		"1: ./AIservice -init  , 初始化配置文件 aiges.toml (若存在，则不会替换)\n" +
		"2: ./AIservice -m=0 , 仅用于本地模式运行\n" +
		"3: ./AIservice -mnist , 下载mnistdemo\n" +
		"- 配置中心模式 (开源计划删除)\n" +
		"- 更多参数选项: 请执行 ./AIservice -h \n") // TODO usage() 完善
	os.Exit(0)
}
