package main

import (
	"flag"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"log"
	"sync"
)

const (
	cfgName      = "server.toml"
	project      = "3s"
	group        = "3s"
	service      = "xsf-server"
	version      = "x.x.x"
	apiVersion   = "1.0.0"
	cachePath    = "xxx"
	companionUrl = "http://10.1.87.70:6868"
)

func init() {
	flag.Parse()
}
func main() {

	//定义一个服务实例
	var serverInst xsf.XsfServer

	//定义相关的启动参数
	/*
		1、CfgMode可选值Native、Centre，native为本地配置读取模式，Centre为配置中心模式，当此值为-1时，表示有命令行传入
		2、CfgName 配置文件名
		3、Project 配置中心用 项目名
		4、Group 配置中心用 组名
		5、Service 配置中心用 服务名
		6、Version 配置中心用 配置版本名
		7、CompanionUrl 配置中心用 配置中心地址
	*/
	bc := xsf.BootConfig{
		CfgMode: utils.Native,
		CfgData: xsf.CfgMeta{
			CfgName:      cfgName,
			Project:      project,
			Group:        group,
			Service:      service,
			Version:      version,
			ApiVersion:   apiVersion,
			CachePath:    cachePath,
			CompanionUrl: companionUrl}}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		/*
			1、启动服务
			2、若有异常直接报错，注意需用户自己实现协程等待
		*/
		if err := serverInst.Run(
			bc,
			&server{},
			xsf.SetOpRouter(generateOpRouter())); err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()
}
