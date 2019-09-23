package xsf

import (
	"log"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

//框架初始化会首先执行此函数获取本机地址
//默认获取本机第一个非回环地址，更多的实现函数可参考utils中的localaddr文件
var (
	netAddrNotLoopback string
)

func init() {
	var netErr error
	netAddrNotLoopback, netErr = utils.GetAddr()
	if netErr != nil {
		log.Panic(netErr)
	}
}

//return network addr
func GetNetAddr() (string) {
	return netAddrNotLoopback
}
