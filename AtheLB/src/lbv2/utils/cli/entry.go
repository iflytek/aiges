package main

import (
	"flag"
	xsf "git.xfyun.cn/AIaaS/xsf-external/client"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"log"
)

const (
	clientCfg  = "cli.toml"
	cname      = "lb_cli"
	cfgService = "lb_cli"
	cfgVersion = "0.0.0"
)

var (
	cfgUrl   = flag.String("u", "http://10.1.87.69:6868", "cfgUrl")
	cfgPrj   = flag.String("p", "xsf", "cfgPrj")
	cfgGroup = flag.String("g", "xsf", "cfgGroup")
	mode     = flag.Int("mode", 0, "0:native,1:centre")
)
var (
	nbest  = flag.String("nbest", "1", "nbest")
	svc    = flag.String("svc", "iat", "svc")
	subsvc = flag.String("subsvc", "sms", "subsvc")
	uid    = flag.String("uid", "-1", "uid")
	all    = flag.String("all", "0", "all")
	lbname = flag.String("lbname", "xsf-lbv2", "lbname")
	c      = flag.Int("c", 1, "concurrent")
	n      = flag.Int("n", 10, "total count")
	perf   = flag.Int("perf", 0, "perf")
	action = flag.Int("action", 0, "0:get_server,1:force_offline")
)

func main() {
	flag.Parse()

	var c *xsf.Client
	var e error

	switch *mode {
	case 0:
		{
			c, e = xsf.InitClient(
				cname,
				utils.Native,
				utils.WithCfgName(clientCfg),
				utils.WithCfgURL(*cfgUrl),
				utils.WithCfgPrj(*cfgPrj),
				utils.WithCfgGroup(*cfgGroup),
				utils.WithCfgService(cfgService),
				utils.WithCfgVersion(cfgVersion),
			)
			if e != nil {
				log.Fatal("main | InitCient error:", e)
			}
		}
	case 1:
		{
			c, e = xsf.InitClient(
				cname,
				utils.Centre,
				utils.WithCfgName(clientCfg),
				utils.WithCfgURL(*cfgUrl),
				utils.WithCfgPrj(*cfgPrj),
				utils.WithCfgGroup(*cfgGroup),
				utils.WithCfgService(cfgService),
				utils.WithCfgVersion(cfgVersion),
			)
			if e != nil {
				log.Panic("main | InitCient error:", e)
			}
		}
	default:
		{
			panic("Oops!!!")
		}
	}

	switch *action {
	case 0:
		{
			getServer(c)
		}
	case 1:
		{
			forceOffline(c)
		}

	}
}
