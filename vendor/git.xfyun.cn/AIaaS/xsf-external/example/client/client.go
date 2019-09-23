package main

import (
	"flag"
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/client"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"log"
	"time"
)

const (
	cname = "client" //配置文件的主段名

	clientCfg  = "client.toml"
	cfgUrl     = "http://10.1.87.70:6868"
	cfgPrj     = "3s"
	cfgGroup   = "3s"
	cfgService = "xsf-client" //服务发现的服务名
	cfgVersion = "x.x.x"      //配置文件的版本号
	apiVersion = "1.0.0"      //api版本号，一般不用修改

	cacheService = true
	cacheConfig  = true
	cachePath    = "./findercache" //配置缓存路径
)

var (
	tm   = flag.Int64("tm", 1000, "timeout")
	mode = flag.Int64("mode", 0, "0:native;1:center")
)

func ssb(c *xsf.Caller, tm time.Duration) (*xsf.Res, string, int32, error) {
	req := xsf.NewReq()
	req.SetParam("k1", "v1")
	req.SetParam("k2", "v2")
	req.SetParam("k3", "v3")
	res, code, e := c.SessionCall(xsf.CREATE, "xsf-server", "ssb", req, tm)
	if code != 0 || e != nil {
		log.Fatal("ssb err")
	}

	var sess string
	if e == nil {
		sess = res.Session()
	}
	return res, sess, code, e
}

func auw(c *xsf.Caller, sess string, tm time.Duration) (*xsf.Res, int32, error) {
	req := xsf.NewReq()
	req.SetParam("k1", "v1")
	req.SetParam("k2", "v2")
	req.SetParam("k3", "v3")

	req.Session(sess)

	res, code, e := c.SessionCall(xsf.CONTINUE, "xsf-server", "auw", req, tm)
	if code != 0 || e != nil {
		log.Fatal("auw err")
	}

	return res, code, e
}

func sse(c *xsf.Caller, sess string, tm time.Duration) (*xsf.Res, int32, error) {
	req := xsf.NewReq()
	req.SetParam("k1", "v1")
	req.SetParam("k2", "v2")
	req.SetParam("k3", "v3")

	req.Session(sess)
	res, code, e := c.SessionCall(xsf.CONTINUE, "xsf-server", "sse", req, tm)
	if code != 0 || e != nil {
		log.Fatal("sse err")
	}
	return res, code, e
}

func sessionCallExample(c *xsf.Caller, tm time.Duration) {

	c.WithApiVersion(apiVersion)

	_, sess, _, _ := ssb(c, tm)
	_, _, _ = auw(c, sess, tm)
	_, _, _ = sse(c, sess, tm)

}
func callExample(c *xsf.Caller, tm time.Duration) {
	span := utils.NewSpan(utils.CliSpan).Start()
	defer span.End().Flush()

	span = span.WithName("callExample")
	span = span.WithTag("customKey1", "customVal1")
	span = span.WithTag("customKey2", "customVal2")
	span = span.WithTag("customKey3", "customVal3")
	c.WithApiVersion(apiVersion)

	req := xsf.NewReq()
	req.SetParam("k1", "v1")
	req.SetParam("k2", "v2")
	req.SetParam("k3", "v3")

	req.SetTraceID(span.Meta()) //将span信息带到后端
	_, code, e := c.Call("xsf-server", "req", req, tm)
	if code != 0 || e != nil {
		log.Fatal("sse err", code, e)
	}

}
func callWithAddr(c *xsf.Caller, tm time.Duration) {

	req := xsf.NewReq()
	req.SetParam("k1", "v1")
	req.SetParam("k2", "v2")
	req.SetParam("k3", "v3")

	_, code, e := c.CallWithAddr("xsf-server", "req", "127.0.0.1:1997", req, tm)
	if code != 0 || e != nil {
		log.Fatal("sse err", code, e)
	}

}

func main() {
	flag.Parse()

	cli, cliErr := xsf.InitClient(
		cname,
		func() xsf.CfgMode {
			switch *mode {
			case 0:
				{
					fmt.Println("about to init native client")
					return utils.Native
				}
			default:
				{
					fmt.Println("about to init centre client")
					return utils.Centre
				}
			}
		}(),
		utils.WithCfgCacheService(cacheService),
		utils.WithCfgCacheConfig(cacheConfig),
		utils.WithCfgCachePath(cachePath),
		utils.WithCfgName(clientCfg),
		utils.WithCfgURL(cfgUrl),
		utils.WithCfgPrj(cfgPrj),
		utils.WithCfgGroup(cfgGroup),
		utils.WithCfgService(cfgService),
		utils.WithCfgVersion(cfgVersion),
	)
	if cliErr != nil {
		log.Fatal("main | InitCient error:", cliErr)
	}
	fmt.Println("ip:", cli.Cfg().GetLocalIp())

	callWithAddr(xsf.NewCaller(cli), time.Second)
	//callExample(xsf.NewCaller(cli), time.Millisecond*time.Duration(*tm))
	//sessionCallExample(xsf.NewCaller(cli), time.Millisecond*time.Duration(*tm))

}
