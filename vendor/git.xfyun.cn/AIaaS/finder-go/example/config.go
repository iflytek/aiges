package main

import (
	"fmt"

	common "git.xfyun.cn/AIaaS/finder-go/common"
	"strings"
	"log"
)

// ConfigChangedHandle ConfigChangedHandle
type ConfigChangedHandle struct {
}

// OnConfigFileChanged OnConfigFileChanged
func (s *ConfigChangedHandle) OnConfigFileChanged(config *common.Config) bool {
	if strings.HasSuffix(config.Name, ".toml") {
		fmt.Println(config.Name, " has changed:\r\n", string(config.File), " \r\n 解析后的map为 ：", config.ConfigMap)
	} else {
		fmt.Println(config.Name, " has changed:\r\n", string(config.File))
	}
	config.File=nil
	config.Name=""
	config.ConfigMap=nil
	config=nil
	return true
}

func (s *ConfigChangedHandle) OnError(errInfo common.ConfigErrInfo){
	log.Println("配置文件出错：",errInfo)
}
