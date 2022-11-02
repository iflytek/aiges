package main

import (
	"flag"
	"fmt"
	_ "github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/env"
	"github.com/xfyun/aiges/service"
	"github.com/xfyun/aiges/utils"
	"github.com/xfyun/aiges/widget"
	"os"
)

func main() {
	flag.Parse()
	env.Parse()
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

	var aisrv service.EngService
	widgetInst := widget.NewWidget(env.AIGES_PLUGIN_MODE)
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
	if err = aisrv.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}
	return

}

func usage() {
	fmt.Printf("TODO:加载器参数说明\n") // TODO usage() 完善
	os.Exit(0)
}
