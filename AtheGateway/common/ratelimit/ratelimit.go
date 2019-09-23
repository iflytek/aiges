package ratelimit

import (
	"sync"
)

const(
	PolicyLocal = "local"
	PolicyRedis = "redis"
	KeyGlobal = "global"
	KeyConnGlobal = "conn_global"
	KeyConnMax = "connmax"
)

const (
	LimitByAppid = "appid"
	LimitByIp    = "ip"
	LimitByAppidConn    = "conn"
	LimitMax    = "max"
)

const KeyPrefix = "webgate-ws-ratelimit-"
const KeyConnPrefix = "webgate-ws-conn-"

const(
	MaxData int = 0x7fffffffffffffff
)

var (
	limitDataBase LimitDataBase
)

type RateMeta struct {
	Ip string
	Appid string
}

type LimitDataBase interface {
	Limit(key string)(int64,error)
	Release(key string)(int64,error)
}

type LocalDataBase struct {
	data sync.Map
}


type RateLimitCache struct {
	cache sync.Map
}

func (c *RateLimitCache)Get(key string) RateLimit {
	r,ok:= c.cache.Load(key)
	if ok{
		return r.(RateLimit)
	}
	return nil
}

func (c *RateLimitCache)Set(val RateLimit)  {
	c.cache.Store(val.GetKey(),val)
}

var rateLimitManager *RateLimtManager

type RateLimtManager struct {
	globalCache *RateLimitCache  // 全局的ratelimit 配置
	personalCache *RateLimitCache // 私有的ratelimit配置
}

func generateKey(s string)string  {
	if s==""{
		return KeyGlobal
	}
	return KeyPrefix+s
}

func generateConnKey(s string)string  {
	if s==""{
		return KeyConnGlobal
	}
	return KeyConnPrefix+s
}

func (rm *RateLimtManager)GetRateLimit(appid string,ip string)RateLimit  {
	rl:=rm.personalCache.Get(generateKey(appid))
	if rl ==nil{
		rl = rm.personalCache.Get(generateKey(ip))
	}
	if rl == nil{
		rl = rm.globalCache.Get(KeyGlobal)
	}
	return rl
}

func (rm *RateLimtManager)GetConnLimit(appid string,ip string)RateLimit  {
	rl:=rm.personalCache.Get(generateConnKey(appid))
	if rl == nil{
		rl = rm.globalCache.Get(generateConnKey(ip))
	}
	if rl == nil{
		rl = rm.globalCache.Get(KeyConnGlobal)
	}

	return rl
}

func (rm *RateLimtManager)GetMaxLimit() RateLimit  {
	return rm.globalCache.Get(KeyConnMax)
}

func (rm *RateLimtManager)LoadRateLimit(rl RateLimit)  {
	if rl.GetKey()==KeyGlobal || rl.GetKey() == KeyConnGlobal || rl.GetKey()==KeyConnMax{
		rm.globalCache.Set(rl)
	}else{
		rm.personalCache.Set(rl)
	}
}

func LoadRateLimitConfig(cfgs []RateLimit)  {
	for _,rl:=range cfgs{
		rateLimitManager.LoadRateLimit(rl)
	}

}

func CheckLimitRate(appid string,ip string,rl *AppidLimit) bool {
	if rl == nil{
		return true
	}
	return false
}



