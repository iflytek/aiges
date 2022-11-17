package utils

import (
	"github.com/xfyun/xsf/utils"
)

// 用于协作 python初始化协程和xsf 框架初始化 同步顺序

type Coordinator struct {
	ConfChan chan *utils.Configure
	Ch2      chan int
}
