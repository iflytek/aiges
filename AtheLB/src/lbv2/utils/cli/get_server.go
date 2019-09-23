package main

import (
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/client"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"github.com/cihub/seelog"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func getServer(cli *xsf.Client) {

	{
		//日志初始化
		if logger, err := seelog.LoggerFromConfigAsString(`<seelog type="sync">
			<outputs formatid="main">
				<filter levels="trace,debug,info,warn,error,critical">
					<console/>
				</filter>
				<filter levels="trace,debug,info,warn,error,critical">
					<file path="cli.log"/>
				</filter>
			</outputs>
			<formats>
				<format id="main" format="%Msg"/>
			</formats>
		</seelog>`); err != nil {
				log.Fatal(err)
			} else {
				seelog.ReplaceLogger(logger)
			}
	}

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGKILL)
		s := <-c
		switch s {
		case syscall.SIGINT, syscall.SIGKILL:
			{
				fmt.Printf("the program had received %v signal, will exit immediately -_-|||", s.String())
				seelog.Flush()
				os.Exit(1)
			}
		case syscall.SIGPIPE:
			{
				fmt.Printf("get broken pipe")
			}
		}
	}()


	_ = seelog.Criticalf("\n====================I am the dividing line========================\n")
	caller := xsf.NewCaller(cli)
	//创建请求
	req := utils.NewReq()
	req.SetParam("nbest", *nbest)
	req.SetParam("svc", *svc)
	req.SetParam("subsvc", *subsvc)
	req.SetParam("all", *all)

	if (*uid) != "-1" {
		req.SetParam("uid", *uid)
		req.SetParam("rp", "1")
	}

	allTime := int64(0)
	cnt := int64(0)
	total := int64(*n)
	fail := int64(0)
	var min, max int64 = math.MaxInt64, math.MinInt64
	qpsBase := time.Now()
	wg := sync.WaitGroup{}
	var firstAddrGet string
	var firstAddrGetOnce sync.Once
	var miss int64 //记录获取不一致地址的个数
	for i := 0; i < *c; i++ {
		wg.Add(1)
		go func() {
			defer func() { wg.Done() }()
			for {
				atomic.AddInt64(&cnt, 1)
				//执行请求
				ts := time.Now()
				res, code, e := caller.Call(*lbname, "getServer", req, time.Second)
				dur := time.Now().Sub(ts).Nanoseconds()
				if atomic.LoadInt64(&min) > dur {
					atomic.StoreInt64(&min, dur)
				}
				if atomic.LoadInt64(&max) < dur {
					atomic.StoreInt64(&max, dur)
				}
				atomic.AddInt64(&allTime, time.Now().Sub(ts).Nanoseconds())
				if e != nil {
					atomic.AddInt64(&fail, 1)
					panic(fmt.Sprintf("dur:%dms,code:%v,e:%v\n", dur/1e6, code, e))
				}
				dataArr := res.GetData()
				var bestNodes []string
				for _, data := range dataArr {
					bestNode := string(data.GetData())
					firstAddrGetOnce.Do(func() {
						firstAddrGet = bestNode
					})
					if bestNode != firstAddrGet {
						atomic.AddInt64(&miss, 1)
					}
					bestNodes = append(bestNodes, bestNode)
				}
				if *perf == 0 {
					seelog.Infof("bestNode:%+v,ts:%v\n", bestNodes, time.Now())
				}
				if atomic.LoadInt64(&cnt) >= atomic.LoadInt64(&total) {
					break
				}
			}
		}()
	}
	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Printf("currently completed: %v%%\n", float64(atomic.LoadInt64(&cnt))/float64(atomic.LoadInt64(&total)))
		}
	}()
	wg.Wait()
	qpsDur := time.Now().Sub(qpsBase).Nanoseconds()
	seelog.Criticalf("\n====================stats========================\n")
	seelog.Criticalf("nbest:%v\n", *nbest)
	seelog.Criticalf("svc:%v\n", *svc)
	seelog.Criticalf("subsvc:%v\n", *subsvc)
	seelog.Criticalf("all:%v\n", *all)
	seelog.Criticalf("lbname:%v\n", *lbname)
	seelog.Criticalf("----\n")
	seelog.Criticalf("concurrent:%v\n", *c)
	seelog.Criticalf("totalCount:%v\n", cnt)
	if (*uid) != "-1" {
		seelog.Criticalf("uid:%v\n", *uid)
		seelog.Criticalf("miss:%v\n", miss)
		seelog.Criticalf("hitRate:%v\n", float64(cnt-miss)/float64(cnt))
	}
	seelog.Criticalf("fail:%v\n", fail)
	seelog.Criticalf("rate:%v\n", float64(cnt-fail)/float64(cnt))
	seelog.Criticalf("qps:%v\n", float64(cnt)/(float64(qpsDur)/1e9))
	seelog.Criticalf("max:%vms\n", float64(max)/1e6)
	seelog.Criticalf("min:%vms\n", float64(min)/1e6)
	seelog.Criticalf("allTime:%vms\n", float64(allTime)/1e6)
	seelog.Criticalf("avg:%vms\n", (float64(allTime)/float64(cnt))/1e6)
}
