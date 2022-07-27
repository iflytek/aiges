package main

import (
	"fmt"
	xsfcli "git.iflytek.com/AIaaS/xsf/client"
	"github.com/pterm/pterm"
	"sync"
	"time"
	"xtest/analy"
	"xtest/prometheus"
	"xtest/request"
	"xtest/util"
	"xtest/var"
)

type Xtest struct {
	r   request.Request
	cli *xsfcli.Client
}

func NewXtest(cli *xsfcli.Client, conf _var.Conf) Xtest {
	return Xtest{r: request.Request{C: conf}, cli: cli}
}

func (x *Xtest) Run() {
	// 数据分析初始化、性能数据
	analy.ErrAnalyser.Start(x.r.C.MultiThr, x.cli.Log, x.r.C.ErrAnaDst)
	if x.r.C.PerfConfigOn {
		analy.Perf = new(analy.PerfModule)
		analy.Perf.Log = x.cli.Log
		startErr := analy.Perf.Start()
		if startErr != nil {
			fmt.Println("failed to open req record file.", startErr.Error())
			return
		}
		defer analy.Perf.Stop()
	}
	// 启动异步输出打印&落盘
	var rwg sync.WaitGroup
	for i := 0; i < x.r.C.DropThr; i++ {
		rwg.Add(1)
		go x.r.DownStreamWrite(&rwg, x.cli.Log)
	}

	var wg sync.WaitGroup

	// jbzhou5
	r := prometheus.NewResources()     // 开启资源监听实例
	stp := util.NewScheduledTaskPool() // 开启一个定时任务池
	if x.r.C.PrometheusSwitch {
		go r.Serve(x.r.C.PrometheusPort) // jbzhou5 启动一个协程写入Prometheus
	}

	if x.r.C.Plot {
		r.GenerateData()
	}

	// 启动一个系统资源定时任务
	stp.Start(time.Microsecond*100, func() {
		err := r.ReadMem(&x.r.C)
		if err != nil {
			return
		}
	})

	go util.ProgressShow(x.r.C.LoopCnt, x.r.C.LoopCnt.Load())

	for i := 0; i < x.r.C.MultiThr; i++ {
		wg.Add(1)
		go func() {
			for {
				loopIndex := x.r.C.LoopCnt.Load()
				x.r.C.LoopCnt.Dec()
				if x.r.C.LoopCnt.Load() < 0 {
					break
				}

				switch x.r.C.ReqMode {
				case 0:
					info := x.r.OneShotCall(x.cli, loopIndex)
					analy.ErrAnalyser.PushErr(info)
				case 1:
					info := x.r.SessionCall(x.cli, loopIndex) // loopIndex % len(stream.dataList)
					analy.ErrAnalyser.PushErr(info)
				case 2:
					info := x.r.TextCall(x.cli, loopIndex) // loopIndex % len(stream.dataList)
					analy.ErrAnalyser.PushErr(info)
				case 3:
					info := x.r.FileSessionCall(x.cli, loopIndex) // loopIndex % len(stream.dataList)
					analy.ErrAnalyser.PushErr(info)
				default:
					println("Unsupported Mode!")
				}
			}
			wg.Done()
		}()
		x.linearCtl() // 并发线性增长控制,防止瞬时并发请求冲击
	}
	wg.Wait()
	// 关闭异步落盘协程&wait
	close(x.r.C.AsyncDrop)
	analy.ErrAnalyser.Stop()
	rwg.Wait()
	xsfcli.DestroyClient(x.cli)
	stp.Stop() // 关闭定时任务
	r.Stop()   // 关闭资源收集
	r.Dump() // 持久化资源日志
	if x.r.C.Plot {
		r.Draw(x.r.C.PlotFile)
	}
	pterm.DefaultBasicText.Println(pterm.LightGreen("\ncli finish "))
}

func (x *Xtest) linearCtl() {
	if x.r.C.LinearNs > 0 {
		time.Sleep(time.Duration(time.Nanosecond) * time.Duration(x.r.C.LinearNs))
	}
}
