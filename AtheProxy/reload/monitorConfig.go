package reload

import (
	"time"
	common "git.xfyun.cn/AIaaS/finder-go/common"
	"config"
	"git.xfyun.cn/AIaaS/finder-go"

	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"fmt"
	"strings"
	"util"
)

type ConfigChangedHandle struct {
}

// OnConfigFileChanged OnConfigFileChanged
func (s *ConfigChangedHandle) OnConfigFileChanged(fconfig *common.Config) bool {

	cfg, err := utils.NewCfgWithBytes(string(fconfig.File))
	if err == nil {
		subs := strings.Split(config.SUB_TYPES, ",")
		subMap := make(map[string]*config.AtmosConfig)
		for _, sub := range subs {
			if sub == "" {
				continue
			}
			subMap[sub] = util.GetParamsBySub(sub, cfg)
		}
		config.SUB_MAP = subMap

		fmt.Println(time.Now())
		for k, v := range config.SUB_MAP {
			fmt.Println(k, v)
		}

	}
	return true
}

func (s *ConfigChangedHandle) OnError(errInfo common.ConfigErrInfo) {

}

var co *utils.CfgOption

func MonitorConfig(cfgOption *utils.CfgOption) {
	co = cfgOption
	config := common.BootConfig{
		//companion地址
		CompanionUrl: config.CompanionUrl,
		//缓存路径
		CachePath: "",
		//是否缓存服务信息
		CacheService: false,
		//是否缓存配置信息
		CacheConfig:   true,
		ExpireTimeout: 5 * time.Second,
		MeteData: &common.ServiceMeteData{
			Project: config.Project,
			Group:   config.Group,
			Service: config.Service,
			Version: config.Version,
			Address: "monitorconfig",
		},
	}
	f, err := finder.NewFinderWithLogger(config, util.SugarLog)
	if err == nil {
		f.ConfigFinder.UseAndSubscribeConfig([]string{"atmos.toml"}, &ConfigChangedHandle{})
	} else {
		fmt.Println("MonitorConfig err:", err)
	}
}
