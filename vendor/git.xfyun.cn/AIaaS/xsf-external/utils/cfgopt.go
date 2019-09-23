package utils

import (
	"fmt"
	"time"
)

// CfgMode 表示配置加载的方式
type CfgMode int

// String CfgMode串化表示
func (s CfgMode) String() string {
	switch s {
	case Native:
		return "Native"
	case Centre:
		return "Centre"
	case Custom:
		return "Custom"
	default:
		return "Invalid-CfgMode"
	}
}

const (
	//  本地加载方式
	Native CfgMode = iota
	// 使用配置中心加载方式
	Centre
	// 自定义实现，暂不支持
	Custom
)

/*CfgOption:
配置选项
*/
type CfgOption struct {
	uuid string
	name string                  //配置名称（文件）
	cb   func(c *Configure) bool //配置更新类
	def  string                  //默认配置

	fm *FindManger // 配置中心操作句柄

	// 配置中心相关设置项
	prj   string // 项目名称
	group string // 集群名称
	srv   string // 服务名
	ver   string // 版本号，字符串
	url   string // 配置中心地址，可以支持域名、ip:port

	cachePath    string
	cacheConfig  bool
	cacheService bool

	tick    time.Duration // zk的tick超时
	stmout  time.Duration // 服务发现zk的session timeout
	zkcontm time.Duration // zk 连接超时
	zksleep time.Duration //tod0
	zkretry int           //zk 重试次数

	mode CfgMode // 配置的加载方式

	log *Logger // 日志句柄

	SvcIp   string //服务端监听ip，trace用
	SvcPort int32  //服务端监听端口，trace用

	//自动提取
	localIp string
}

func (co *CfgOption) String() string {
	return fmt.Sprintf("name:%v,def:%v,prj:%v,group:%v,srv:%v,ver:%v,url:%v,mode:%v",
		co.name, co.def, co.prj, co.group, co.srv, co.ver, co.url, co.mode)
}

/*

	设置cache默认值
*/
func (co *CfgOption) SetDef(
	defaultCacheService bool,
	defaultCacheConfig bool,
	defaultCachePath string,
	uuid string) {
	co.cacheService = defaultCacheService
	co.cacheConfig = defaultCacheConfig
	co.cachePath = defaultCachePath
	co.uuid = uuid
}
func (co *CfgOption) SetLocalIp(localIp string) {
	co.localIp = localIp
}

// FindManger 返回配置中心操作句柄
func (co *CfgOption) FindManger() *FindManger {
	return co.fm
}

// CfgOpt 设置配置选项的函数类型
type CfgOpt func(*CfgOption)

func WithCfgCacheService(cacheService bool) CfgOpt {
	return func(c *CfgOption) {
		c.cacheService = cacheService
	}
}

func WithCfgCacheConfig(cacheConfig bool) CfgOpt {
	return func(c *CfgOption) {
		c.cacheConfig = cacheConfig
	}
}

func WithCfgCachePath(cachePath string) CfgOpt {
	return func(c *CfgOption) {
		c.cachePath = cachePath
	}
}

// WithCfgName 指定配置文件名称
func WithCfgName(name string) CfgOpt {
	return func(c *CfgOption) {
		c.name = name
	}
}

// WithCfgDefault 指定配置默认加载项
func WithCfgDefault(cfg string) CfgOpt {
	return func(c *CfgOption) {
		c.def = cfg
	}

}

// WithCfgPrj 指定使用配置中心时的项目名
func WithCfgPrj(cfg string) CfgOpt {
	return func(c *CfgOption) {
		c.prj = cfg
	}
}

// WithCfgGroup 指定使用配置中心时的集群名
func WithCfgGroup(cfg string) CfgOpt {
	return func(c *CfgOption) {
		c.group = cfg
	}
}

// WithCfgService 指定使用配置中心时的服务名
func WithCfgService(cfg string) CfgOpt {
	return func(c *CfgOption) {
		c.srv = cfg
	}
}

func WithCfgSvcIp(svcIp string) CfgOpt {
	return func(c *CfgOption) {
		c.SvcIp = svcIp
	}
}

func WithCfgSvcPort(svcPort int32) CfgOpt {
	return func(c *CfgOption) {
		c.SvcPort = svcPort
	}
}

// WithCfgVersion 指定使用配置中心时的版本号
func WithCfgVersion(cfg string) CfgOpt {
	return func(c *CfgOption) {
		c.ver = cfg
	}
}

// WithCfgURL 指定使用配置中心时的URL，支持域名或者ip:port
func WithCfgURL(url string) CfgOpt {
	return func(c *CfgOption) {
		c.url = url
	}
}

// withCfgMode 指定使用配置的加载模式
func withCfgMode(mode CfgMode) CfgOpt {
	return func(c *CfgOption) {
		c.mode = mode
	}
}

// WithCfgCB 指定配置中心改动时触发的集成方的回调方式
func WithCfgCB(cb func(c *Configure) bool) CfgOpt {
	return func(c *CfgOption) {
		c.cb = cb
	}
}

// WithCfgReader 指定配置中心时的操作句柄
func WithCfgReader(fm *FindManger) CfgOpt {
	return func(c *CfgOption) {
		c.fm = fm
	}
}

// WithCfgLog 指定配置操作的日志句柄
func WithCfgLog(l *Logger) CfgOpt {
	return func(c *CfgOption) {
		c.log = l
	}
}

// WithCfgTick 指定配置zookeeper的心跳时间
func WithCfgTick(t time.Duration) CfgOpt {
	return func(c *CfgOption) {
		c.tick = t
	}
}

// WithCfgSessionTimeOut 指定配置zookeeper session 超时时间
func WithCfgSessionTimeOut(t time.Duration) CfgOpt {
	return func(c *CfgOption) {
		c.stmout = t
	}
}
