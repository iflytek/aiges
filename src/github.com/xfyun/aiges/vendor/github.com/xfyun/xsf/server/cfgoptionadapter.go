package xsf

import (
	"github.com/xfyun/xsf/utils"
)

const (
	//defaultCacheService = true
	//defaultCacheConfig  = true
	defaultCachePath = "."
)

type CfgMeta struct {
	CfgName      string
	CfgDefault   string
	Project      string
	Group        string
	Service      string
	Version      string
	ApiVersion   string
	CompanionUrl string

	CachePath string

	CallBack func(c *utils.Configure) bool
}

// native stand for local cfg
// online stand for online cfg
type BootConfig struct {
	CfgMode utils.CfgMode
	CfgData CfgMeta
}
