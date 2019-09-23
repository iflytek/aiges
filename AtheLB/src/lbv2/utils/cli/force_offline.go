package main

import (
	"fmt"
	xsf "git.xfyun.cn/AIaaS/xsf-external/client"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"time"
)

const (
	/*
		1、黑名单，用于外部主动下线某个节点
		2、forceOffline参数为0表示强制引擎下线，为1表示解除强制下线
		3、下线后添加至此名单，并不再接受上报
	*/

	ADDR              = "addr"
	FORCEOFFLINE      = "forceOffline"
	ForceOffline      = "0"
	CleanForceOffline = "1"
	CMDSERVER         = "cmdServer"
)

func forceOffline(cli *xsf.Client) {
	const lbAddr = "127.0.0.1:1995"
	const engAddr = "1.1.3.1:1111"
	//offline(cli, engAddr, lbAddr)
	online(cli, engAddr, lbAddr)
}

func offline(cli *xsf.Client, engAddr, lbAddr string) {
	req := utils.NewReq()
	req.SetParam(ADDR, engAddr)
	req.SetParam(FORCEOFFLINE, ForceOffline)
	res, _, err := xsf.NewCaller(cli).CallWithAddr("", CMDSERVER, lbAddr, req, time.Second)
	if nil != err {
		fmt.Println(err)
	} else {
		fmt.Println(res.GetParam("status"))
	}
}
func online(cli *xsf.Client, engAddr, lbAddr string) {
	req := utils.NewReq()
	req.SetParam(ADDR, engAddr)
	req.SetParam(FORCEOFFLINE, CleanForceOffline)
	res, _, err := xsf.NewCaller(cli).CallWithAddr("", CMDSERVER, lbAddr, req, time.Second)
	if nil != err {
		fmt.Println(err)
	} else {
		fmt.Println(res.GetParam("status"))
	}
}
