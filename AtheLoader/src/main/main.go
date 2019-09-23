package main

import (
	"conf"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"service"
	"widget"
)

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		usage()
	}
	if *conf.CmdVer {
		fmt.Println(service.VERSION)
		return
	}
	if *conf.CmdPrf {
		go func() {
			lsAddr := *conf.CmdPrfAddr
			if len(lsAddr) == 0 {
				lsAddr = ":1234" // default pprof listen port
			}
			log.Println(http.ListenAndServe(lsAddr, nil))
		}()
	}

	var aisrv service.EngService
	widgetInst := widget.NewWidget()
	// 控件初始化&逆初始化
	errInfo := widgetInst.Open()
	if errInfo != nil {
		fmt.Println(errInfo.Error())
		return
	}
	defer widgetInst.Close()

	// 框架初始化&逆初始化
	errInfo = aisrv.Init(widgetInst.Version())
	if errInfo != nil {
		fmt.Println(errInfo.Error())
		return
	}
	defer aisrv.Fini()

	// 注册行为
	errInfo = widgetInst.Register(&aisrv)
	if errInfo != nil {
		fmt.Println(errInfo.Error())
		return
	}

	// 框架运行
	errInfo = aisrv.Run()
	if errInfo != nil {
		fmt.Println(errInfo.Error())
		return
	}
	return
}

func usage() {
	fmt.Printf("TODO:加载器参数说明\n") // TODO usage() 完善
	os.Exit(0)
}
