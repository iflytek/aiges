package ratelimit

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type RateConfig struct {
	ConcurrencyLimit int `json:"limit"`  //限制总数
	Appid         string `json:"appid"`  // appid
	Policy        string `json:"policy"` // 限制策略，数据存redis 还是 内存
	FaultTolerant bool `json:"tolerant"` // 调用限流出错是否容忍
	Ip string `json:"ip"`  // Ip地址
	LimitBy string `json:"limit_by"` // 根据appid还是ip限制
	Enabled bool `json:"enabled"`
}

func (r *RateConfig)AppipLimit()*AppidLimit  {
	return &AppidLimit{
		ConcurrencyLimit:r.ConcurrencyLimit,
		Appid:r.Appid,
		Policy:r.Policy,
		FaultTolerant:r.FaultTolerant,
	}
}

func (r *RateConfig)IpLimit()*IpLimit  {
	return &IpLimit{
		ConcurrencyLimit:r.ConcurrencyLimit,
		Ip:r.Ip,
		Policy:r.Policy,
		FaultTolerant:r.FaultTolerant,
	}
}
func (r *RateConfig)AppidConnLimit()*AppidConnLimit  {
	return &AppidConnLimit{
		ConcurrencyLimit:r.ConcurrencyLimit,
		Appid:r.Appid,
		Policy:r.Policy,
		FaultTolerant:r.FaultTolerant,
	}
}

func (r *RateConfig)MaxLimit()*MaxConnLimit  {
	return &MaxConnLimit{
		ConcurrencyLimit:r.ConcurrencyLimit,
		Policy:r.Policy,
		FaultTolerant:r.FaultTolerant,
	}
}


func LoadRateConfig(b []byte)([]*RateConfig  ,error){
	rcfg:=make([]*RateConfig,0)
	err:=json.Unmarshal(b,&rcfg)
	if err !=nil{
		return nil,err
	}
	return rcfg,nil
}

func LoadConfigToManager(rcfg []*RateConfig)  {
	for _,cfg:=range rcfg{
		switch cfg.LimitBy {
		case LimitByAppid:
			rateLimitManager.LoadRateLimit(cfg.AppipLimit())
			fmt.Println(cfg.AppipLimit())
		case LimitByIp:
			rateLimitManager.LoadRateLimit(cfg.IpLimit())
			fmt.Println(cfg.AppipLimit())
		case LimitByAppidConn:
			rateLimitManager.LoadRateLimit(cfg.AppidConnLimit())
		case LimitMax:
			rateLimitManager.LoadRateLimit(cfg.MaxLimit())
		}

	}
}

func LoadManager(b []byte)error  {
	cfg,err:=LoadRateConfig(b)
	if err !=nil{
		return err
	}
	LoadConfigToManager(cfg)
	return nil
}

func CheckLimit(appid,ip string) bool  {
	rl:=rateLimitManager.GetRateLimit(appid,ip)
	if rl==nil{
		return true
	}
	b,_:=metricData.Add(rl.GetKey())

	if rl.Climit()<b{
		return false
	}
	return true
}

func ReleaseLimit(appid,ip string)  {
	rl:=rateLimitManager.GetRateLimit(appid,ip)
	if rl==nil{
		return
	}
	metricData.Release(rl.GetKey())
}

func CheckConnLimit(appid,ip string) bool {
	rl:=rateLimitManager.GetConnLimit(appid,ip)
	if rl==nil{
		return true
	}
	b,_:=metricData.Add(rl.GetKey())

	if rl.Climit()<b{
		return false
	}
	return true
}

func ReleaseConnLimit(appid,ip string)  {
	rl:=rateLimitManager.GetConnLimit(appid,ip)
	if rl==nil{
		return
	}
	metricData.Release(rl.GetKey())
}


var(
	currentConn int64 = 0
	Max  int64 = 10000
	EnableLimit = false
)

func CheckMaxConnLimit() bool {
	//rl:=rateLimitManager.GetMaxLimit()
	//if rl==nil{
	//	return true
	//}
	//b,_:=metricData.Add(rl.GetKey())
	//if rl.Climit()<b{
	//	return false
	//}
	if !EnableLimit{
		return true
	}
	c := atomic.AddInt64(&currentConn,1)
	if c<=Max{
		return true
	}
	return false
}

func ReleaseMaxConnLimit()  {
	if !EnableLimit{
		return
	}
	atomic.AddInt64(&currentConn,-1)
	//rl:=rateLimitManager.GetMaxLimit()
	//if rl==nil{
	//	return
	//}
	//metricData.Release(rl.GetKey())

}