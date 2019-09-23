package main

import (
	"flag"
	"config"
	"server"
	"util"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"fmt"
	"time"
	"client"
	"reload"
)

//初始化
func init() {
	//解析命令行参数
	flag.Parse()

	//命令行获取参数
	getArgs()

	//加载应用程序配置文件
	util.LoadConfig()

	//初始化日志
	util.InitLogger()
}

//程序入口
func main() {
	defer utils.StopLocalLog(util.SugarLog)

	fmt.Println("server start ................")
	fmt.Println("server start at ", time.Now(), "version:", config.Version)
	fmt.Println("server start ................")

	//监听配置文件的变化，有变动则重新加载
	reload.MonitorConfig(util.CfgOption)

	//初始化rpc server
	client.InitRpcClient()

	//初始化rpc server
	server.InitRpcServer()
}

//命令行获取参数
func getArgs() {
	if *xsf.Mode != -1 {
		config.UseCfgCentre = *xsf.Mode
	}
	if *xsf.CompanionUrl != "" {
		config.CompanionUrl = *xsf.CompanionUrl
	}
	if *xsf.Project != "" {
		config.Project = *xsf.Project
	}
	if *xsf.Group != "" {
		config.Group = *xsf.Group
	}
	if *xsf.Service != "" {
		config.Service = *xsf.Service
	}
}
