package main

import (
	"flag"
	"git.xfyun.cn/AIaaS/xsf-external/client"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"log"
)

const (
	clientCfg  = "srv.toml"
	cname      = "lb_srv"
	cfgService = "lb_srv"
	cfgVersion = "0.0.0"
)

var (
	cfgUrl   = flag.String("u", "http://10.1.87.69:6868", "cfgUrl")
	cfgPrj   = flag.String("p", "xsf", "cfgPrj")
	cfgGroup = flag.String("g", "xsf", "cfgGroup")
	mode     = flag.Int("mode", 0, "0:native,1:centre")
)

func main() {
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
				log.Panic("main | InitCient error:", e)
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
	srv(c)
}
