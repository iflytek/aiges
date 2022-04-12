package xsf

import (
	"log"
	"github.com/xfyun/xsf/utils"
)

//框架初始化会首先执行此函数获取本机地址
//默认获取本机第一个非回环地址，更多的实现函数可参考utils中的localaddr文件
var (
	netaddrNotLoopback string
)

func init() {
	var neterr error
	netaddrNotLoopback, neterr = utils.GetAddr()
	if neterr != nil {
		log.Panic(neterr)
	}
}

//return network addr
func GetNetaddr() (string) {
	return netaddrNotLoopback
}
