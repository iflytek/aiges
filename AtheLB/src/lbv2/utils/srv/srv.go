package main

import (
	"flag"
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/client"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"github.com/cihub/seelog"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	svc    = flag.String("svc", "iat", "svc")
	subsvc = flag.String("subsvc", "sms", "subsvc")
	gaddr  = flag.String("addr", "", "addr,if empty,take rand port")
	gtotal = flag.String("total", "10", "total")
	gbest  = flag.String("best", "10", "best")
	gidle  = flag.String("idle", "10", "idle")

	licMin = flag.Int("min", 0, "idle min")
	licMax = flag.Int("max", 0, "idle max")

	live   = flag.String("live", "1", "live")
	dur    = flag.Int("dur", 1000, "dur ms")
	lbname = flag.String("lbname", "lbv2", "lbname")
)

func init() {
	flag.Parse()
	if logger, err := seelog.LoggerFromConfigAsString(`<seelog type="sync">
    <outputs formatid="main">
        <filter levels="trace,debug,info,warn,error,critical">
        <console/>
        </filter>
        <filter levels="trace,debug,info,warn,error,critical">
            <file path="srv.log"/>
        </filter>
    </outputs>
    <formats>
        <format id="main" format="%LEV=> %Msg"/>
    </formats>
</seelog>`); err != nil {
		log.Fatal(err)
	} else {
		seelog.ReplaceLogger(logger)
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
}
func srv(cli *xsf.Client) {
	var addrStr = *gaddr
	if (*gaddr) == "" {
		addrStr = func() string {
			ipStr := func() string {
				var ipStr string
				addrs, err := net.InterfaceAddrs()
				if err != nil {
					panic("get ip err")
				}
				for _, a := range addrs {
					if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
						if ipNet.IP.To4() != nil {
							ipStr = ipNet.IP.String()
							return ipStr
						}
					}
				}
				return ipStr
			}()
			l, e := net.Listen("tcp", ipStr+":0")
			if e != nil {
				panic(e)
			}
			return l.Addr().String()
		}()
	}
	var totalStr, bestStr string
	dynamic := true
	if (*licMin) == (*licMax) {
		dynamic = false
	} else {
		totalStr = strconv.Itoa(*licMax)
		bestStr = strconv.Itoa((*licMax) / 2)
	}

	c := xsf.NewCaller(cli)
	var cnt int64
	for {
		//创建请求
		req := utils.NewReq()
		req.SetParam("svc", *svc)
		req.SetParam("subsvc", *subsvc)
		req.SetParam("addr", addrStr)

		if dynamic {
			req.SetParam("total", totalStr)
			req.SetParam("best", bestStr)
			req.SetParam("idle", fmt.Sprintf("%v", RandInt64(*licMin, *licMax)))
		} else {
			req.SetParam("total", *gtotal)
			req.SetParam("best", *gbest)
			req.SetParam("idle", *gidle)
		}

		req.SetParam("live", *live)
		//执行请求
		_, code, e := c.Call(*lbname, "setServer", req, 1000*time.Millisecond)

		totalTmp, _ := req.GetParam("total")
		bestTmp, _ := req.GetParam("best")
		idleTmp, _ := req.GetParam("idle")

		if e != nil {
			_ = seelog.Errorf("NO.%d,code:%+v,e:%+v => total:%v,best:%v,idle:%v\n", atomic.AddInt64(&cnt, 1), code, e, totalTmp, bestTmp, idleTmp)
		} else {
			seelog.Infof("NO.%d,code:%+v,e:%+v => total:%v,best:%v,idle:%v\n", atomic.AddInt64(&cnt, 1), code, e, totalTmp, bestTmp, idleTmp)
		}
		time.Sleep(time.Millisecond * time.Duration(*dur))
	}
}
func RandInt64(min, max int) (rst int) {
	for {
		rst = rand.Intn(max)
		if rst < min {
			continue
		} else {
			break
		}
	}
	return
}
