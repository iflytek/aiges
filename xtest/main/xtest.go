package main

import (
	"flag"
	"fmt"
	"github.com/xfyun/aiges/xtest/analy"
	"github.com/xfyun/aiges/xtest/request"
	"github.com/xfyun/aiges/xtest/var"
	xsfcli "github.com/xfyun/xsf/client"
	"github.com/xfyun/xsf/utils"
	"sync"
	"time"
)

func main() {
	flag.Parse()

	// xrpc框架初始化;
	cli, e := xsfcli.InitClient(_var.CliName, utils.CfgMode(0), utils.WithCfgName(*_var.CmdCfg),
		utils.WithCfgURL(""), utils.WithCfgPrj(""), utils.WithCfgGroup(""),
		utils.WithCfgService(""), utils.WithCfgVersion(""))
	if e != nil {
		fmt.Println("cli xsf init fail with ", e.Error())
		return
	}

	// cli配置初始化;
	e = _var.ConfInit(cli.Cfg())
	if e != nil {
		fmt.Println("cli conf init fail with ", e.Error())
		return
	}
	// 数据分析初始化
	analy.ErrAnalyser.Start(_var.MultiThr, cli.Log)
	if _var.PerfConfigOn {
		analy.Perf = new(analy.PerfModule)
		analy.Perf.Log = cli.Log
		startErr := analy.Perf.Start()
		if startErr != nil {
			fmt.Println("failed to open req record file.", startErr.Error())
			return
		}
		defer analy.Perf.Stop()
	}
	// 启动异步输出打印&落盘
	var rwg sync.WaitGroup
	for i := 0; i < _var.DropThr; i++ {
		rwg.Add(1)
		go request.DownStreamWrite(&rwg, cli.Log)
	}

	var wg sync.WaitGroup
	for i := 0; i < _var.MultiThr; i++ {
		wg.Add(1)
		go func() {
			for {
				// loopCnt测试请求次数控制;
				loopIndex := _var.LoopCnt.Dec()
				if loopIndex < 0 {
					break
				}

				switch _var.ReqMode {
				case true:
					code, err := request.SessionCall(cli, loopIndex) // loopIndex % len(stream.dataList)
					analy.ErrAnalyser.PushErr(code, err)
				default:
					code, err := request.OneShotCall(cli, loopIndex)
					analy.ErrAnalyser.PushErr(int(code), err)
				}
			}
			wg.Done()
		}()
		linearCtl() // 并发线性增长控制,防止瞬时并发请求冲击
	}
	wg.Wait()
	// 关闭异步落盘协程&wait
	close(_var.AsyncDrop)
	analy.ErrAnalyser.Stop()
	rwg.Wait()
	xsfcli.DestroyClient(cli)
	fmt.Println("cli finish")
	return
}

func linearCtl() {
	if _var.LinearNs > 0 {
		time.Sleep(time.Duration(time.Nanosecond) * time.Duration(_var.LinearNs))
	}
}
