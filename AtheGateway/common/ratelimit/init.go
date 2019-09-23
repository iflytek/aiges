package ratelimit

import (
	"io/ioutil"
	"fmt"
)

func Init()  {
	//初始化 limitManager
	rateLimitManager = &RateLimtManager{
		globalCache:&RateLimitCache{},
		personalCache:&RateLimitCache{},
	}

	//初始化内存MetricData
	metricData = &LocalMetricDataService{data:map[string]int{}}

	//加载metricData
	f,err:=ioutil.ReadFile("cfg.json")
	if err !=nil{
		panic(err)
	}
	err =LoadManager(f)
	if err !=nil{
		panic(err)
	}

	r:=rateLimitManager.GetRateLimit("100IME","10.1.87.69")
	fmt.Println(r)
}

func InitByConf(max int)  {
	rateLimitManager = &RateLimtManager{
		globalCache:&RateLimitCache{},
		personalCache:&RateLimitCache{},
	}

	//初始化内存MetricData
	metricData = &LocalMetricDataService{data:map[string]int{}}

	LoadConfigToManager([]*RateConfig{
		{
			ConcurrencyLimit:max,
			LimitBy:LimitMax,
		},
	})
}