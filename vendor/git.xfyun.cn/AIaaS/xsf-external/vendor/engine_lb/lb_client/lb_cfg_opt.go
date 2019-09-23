package lb_client

import (
	"time"
)

//配置选项
type LbCfgOption struct {
	LbStrategy     int           //负载策略(必传)
	ZkServerList   []string      //zk列表(必传)
	ZkConnTimeOut  time.Duration //zk连接超时
	Root           string        //根目录
	RouterType     string        //路由类型(如:svc)(必传)
	SubRouterTypes []string      //子路由类型列表(如：[svc_gray,svc_hefei])(必传)
	RedisHost      string        //redis主机(必传)
	RedisPasswd    string        //redis密码
	MaxActive      int           //redis最大连接数(maxActive=0代表没连接限制)
	MaxIdle        int           //redis最大空闲实例数
	Db             int           //redis数据库
	IdleTimeOut    time.Duration //redis空闲实例超时设置(单位:s)
}

type LbCfgOpt func(*LbCfgOption)

//配置负载的策略
func WithLbStrategy(strategy int) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.LbStrategy = strategy
	}
}

//配置zk列表
func WithZkList(zkList []string) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.ZkServerList = zkList
	}
}

//配置zk连接超时时间
func WithZkConnTimeOut(timeout time.Duration) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.ZkConnTimeOut = timeout
	}
}

//配置根目录
func WithRoot(root string) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.Root = root
	}
}

//配置路由类型(业务类型)
func WithRouterType(routerType string) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.RouterType = routerType
	}
}

//配置子路由类型
func WithSubRouterTypes(subRouterType []string) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.SubRouterTypes = subRouterType
	}
}

//配置redis主机
func WithRedisHost(redisHost string) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.RedisHost = redisHost
	}
}

//配置redis主机密码
func WithRedisPassword(redisPwd string) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.RedisPasswd = redisPwd
	}
}

//配置redis数据库
func WithRedisDb(db int) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.Db = db
	}
}

//配置redis最大连接数
func WithRedisMaxActive(maxActive int) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.MaxActive = maxActive
	}
}

//配置redis最大空闲连接数
func WithRedisMaxIdle(maxIdle int) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.MaxIdle = maxIdle
	}
}

//配置redis空闲连接超时
func WithRedisIdleTimeOut(idleTimeOut time.Duration) LbCfgOpt {
	return func(lc *LbCfgOption) {
		lc.IdleTimeOut = idleTimeOut
	}
}
