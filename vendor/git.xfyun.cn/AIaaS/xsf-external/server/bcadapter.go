package xsf

import (
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

//在配置文件里没有对应参数的情况下用命令行参数补齐
//发生警告时，打印警告日志
func bcCheck(bc BootConfig) (BootConfig, error, error) {
	if "" == bc.CfgData.ApiVersion {
		bc.CfgData.ApiVersion = defaultApiVersion
	}

	//check whether the params in bc is missed
	//检查bc参数是否有缺少
	var ext error
	if -1 == bc.CfgMode {
		if -1 == *Mode {
			return bc, fmt.Errorf("can't find cfgmode in preset and command params"), nil
		} else {
			ext = fmt.Errorf("warning ->: %v",
				"the cfgmode use the command params instead of preset.")
			bc.CfgMode = utils.CfgMode(*Mode)
		}
	}
	if "" == bc.CfgData.CfgName {
		if "" == *Cfg {
			return bc, fmt.Errorf("can't find cfgname in preset and command params"), nil
		} else {
			ext = fmt.Errorf("warning ->: %v",
				"the cfgname use the command params instead of preset.")
			bc.CfgData.CfgName = *Cfg
		}
	}
	if "" == bc.CfgData.CfgDefault {
		bc.CfgData.CfgDefault = *DefaultCfg
	}
	if "" == bc.CfgData.Service {
		if "" == *Service {
			return bc, fmt.Errorf("can't find service in preset and command params"), nil
		} else {
			ext = fmt.Errorf("warning ->: %v",
				"the service use the command params instead of preset.")
			bc.CfgData.Service = *Service
		}
	}
	/////////
	if utils.Centre == bc.CfgMode {
		if "" == bc.CfgData.Project {
			if "" == *Project {
				return bc, fmt.Errorf("can't find project in preset and command params"), nil
			} else {
				bc.CfgData.Project = *Project
				ext = fmt.Errorf("warning ->: %v",
					"the project use the command params instead of preset.")
			}
		}
		if "" == bc.CfgData.Group {
			if *Group == "" {
				return bc, fmt.Errorf("can't find group in preset and command params"), nil
			} else {
				ext = fmt.Errorf("warning ->: %v",
					"the group use the command params instead of preset.")
				bc.CfgData.Group = *Group
			}
		}
		if "" == bc.CfgData.CompanionUrl {
			if "" == *CompanionUrl {
				return bc, fmt.Errorf("can't find companionurl in preset and command params"), nil
			} else {
				ext = fmt.Errorf("warning ->: %v",
					"the companionurl use the command params instead of preset.")
				bc.CfgData.CompanionUrl = *CompanionUrl
			}
		}
	}
	/////////

	/*
		finder 2.0增加的参数
	*/
	if "" == bc.CfgData.CachePath {
		bc.CfgData.CachePath = defaultCachePath
	}

	return bc, nil, ext
}
