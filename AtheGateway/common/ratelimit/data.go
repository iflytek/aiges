package ratelimit

import (
	"sync"
	"github.com/go-redis/redis"
)

var(
	metricData MetricDataService
)
type MetricDataService interface {
	Add(key string)(int,error)
	Release(key string)(error)
	AddDelta(key string,deta int)error
}


// 内存中的metric数据
type LocalMetricDataService struct {
	data map[string]int
	lock sync.Mutex
}

func (l *LocalMetricDataService)Add(key string) (int,error) {
	l.lock.Lock()
	c:=l.data[key]
	l.data[key] = c+1
	l.lock.Unlock()
	return c+1,nil
}

func (l *LocalMetricDataService)Release(key string)error  {
	l.lock.Lock()
	c:=l.data[key]
	if c-1>=0{
		l.data[key] = c-1
	}else{
		l.data[key] = 0
	}
	l.lock.Unlock()
	return nil
}

func (l *LocalMetricDataService)AddDelta(key string,deta int)error  {
	l.lock.Lock()
	c:=l.data[key]
	if c-deta>=0{
		l.data[key] = c-deta
	}else{
		l.data[key] = 0
	}
	l.lock.Unlock()
	return nil
}

func (l *LocalMetricDataService)Update(key string,value int)error  {
	l.lock.Lock()
	l.data[key] = value
	l.lock.Unlock()
	return nil
}


//redis 中的metric data 数据
type RedisMetricDataService struct {
	redis redis.Cmdable
}
// lua 脚本实现原子操作
func (this *RedisMetricDataService)Release(key string)(error){
	scrpit:=`
local ns = redis.call("get",KEYS[1])
if (not ns) then
ns = 0
end
local nn = tonumber(ns) -1
if (nn < 0) then
	return 0
end
redis.call("set",KEYS[1],tostring(nn))

return nn

`
	r:=this.redis.Eval(scrpit,[]string{key})


	return r.Err()
}

func (this *RedisMetricDataService)Add(key string)(int,error){
	scrpit:=`
local ns = redis.call("get",KEYS[1])
if (not ns) then
ns = 0
end

local nn = tonumber(ns) +1
redis.call("set",KEYS[1],tostring(nn))
return nn

`
	r:=this.redis.Eval(scrpit,[]string{key})

	return r.Int()
}

func (this *RedisMetricDataService)AddDelta(key string,deta int)error  {

	scrpit:=`
local ns = redis.call("get",KEYS[1])
if (not ns) then
ns = 0
end
local nn = tonumber(ns) + ARGV[1]
if (nn < 0) then
    redis.call("set",KEYS[1],0)
	return 0
end
redis.call("set",KEYS[1],tostring(nn))

return nn

`
	r:=this.redis.Eval(scrpit,[]string{key},deta)


	return r.Err()
}

