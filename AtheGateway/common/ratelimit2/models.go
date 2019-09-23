package ratelimit2

import (
	"sync/atomic"
	"encoding/json"
	"fmt"
	"sync"
)

type Config struct {
	Appid string `json:"appid"`
	ConnLimit int `json:"conn_limit"`
	EffGlobal  bool `json:"eff_global"`
}

type limitConfigMap struct {
	Configs []*Config `json:"configs"`
}

func (c *Config)init()  {
	if c.ConnLimit<=0{
		c.ConnLimit = 100000
	}
}

func init() {
	updateLimitConfigCache(&limitConfigMap{})
}

var configCache atomic.Value
//获取配置缓存示例
func ConfigCacheInstance() *LimitConfigCache {
	return configCache.Load().(*LimitConfigCache)
}
//更新配置缓存
func updateLimitConfigCache(cfgs *limitConfigMap)  {
	cache:=&LimitConfigCache{cache:make(map[string]*Config)}
	for _, v := range cfgs.Configs {
		v.init()
		fmt.Println("initConf:",v.Appid,v.ConnLimit)
		cache.StoreConfig(v)
	}
	configCache.Store(cache)
}
//加载json配置
func LoadConfigCache(b []byte) (err error) {
	var cfgs = &limitConfigMap{}

	err =json.Unmarshal(b,&cfgs)
	if err !=nil{
		return
	}
	updateLimitConfigCache(cfgs)
	return
}

//配置缓存
type LimitConfigCache struct {
	cache map[string]*Config
	 // 白名单
	globalConfig *GlobalConfig    // 配置是否对所有的白名单生效

}

func (c *LimitConfigCache)GetConfig(key string)*Config  {
	return c.cache[key]
}

func (c *LimitConfigCache)StoreConfig(config *Config)  {
	c.cache[config.Appid] = config

}

func (c *LimitConfigCache)Range(f func(string,*Config)bool)  {
	for k,v:=range c.cache{
		if !f(k,v){
			return
		}
	}
}

type GlobalConfig struct {
	Enabled bool
	ConnLimit int
	WhiteList []string
}


type Map struct {
	read atomic.Value
	write map[string]interface{}
	mu sync.Mutex
}

type read struct {
	data map[string]interface{}
}

func NewMap()*Map  {
	rea:=&read{data:make(map[string]interface{})}
	m:=&Map{
		write:make(map[string]interface{}),
	}
	m.read.Store(rea)
	return m
}

func (m *Map)Get(key string)interface{}{
	read:=m.read.Load().(*read)
	return read.data[key]
}

func (m *Map)Set(key string,value interface{})  {
	m.mu.Lock()
	m.write[key] = value
	newR:=&read{data:m.write}
	old:=m.read.Load().(*read)
	m.read.Store(newR)
	old.data[key] = value
	m.write = old.data
	m.mu.Unlock()
}

func (m *Map)Range(f func(string,interface{}) bool)  {
	read:=m.read.Load().(map[string]interface{})
	for k,v:=range read{
		if !f(k,v){
			return
		}
	}
}