package xsf

import (
	"fmt"
	"github.com/xfyun/lb_client"
	"time"
)

type LbAdapter struct {
	svc            string   //负载上报的地址 eg.192.168.86.60:2181
	lbStrategy     int      //负载策略(必传)
	zkList         []string //zk列表(必传)
	root           string   //根目录
	routerType     string   //路由类型(如：iat)(必传)
	subRouterTypes []string //子路由类型列表(如:["iat_gray","iat_hefei"])
	redisHost      string   //redis主机(必传)
	redisPasswd    string   //redis密码
	maxActive      int      //redis最大连接数
	maxIdle        int      //redis最大空闲连接数
	db             int      //redis数据库
	idleTimeOut    time.Duration

	able int //是否启用LB
	lc   lb_client.LbClienter
}
type LbAdapterCfgOpt func(*LbAdapter)

//配置负载的业务名
func WithLbAdapterSvc(svc string) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.svc = svc
	}
}

//配置负载的策略
func WithLbAdapterStrategy(lbStrategy int) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.lbStrategy = lbStrategy
	}
}

//配置负载的zk列表
func WithLbAdapterZkList(zkList []string) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.zkList = zkList
	}
}

//配置负载的根目录
func WithLbAdapterRoot(root string) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.root = root
	}
}

//配置负载的路由类型
func WithLbAdapterRouterType(routerType string) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.routerType = routerType
	}
}

//配置负载的子路由类型列表(如:["iat_gray","iat_hefei"])
func WithLbAdapterSubRouterTypes(subRouterTypes []string) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.subRouterTypes = subRouterTypes
	}
}

//配置负载的redis主机
func WithLbAdapterSRedisHost(redisHost string) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.redisHost = redisHost
	}
}

//配置负载的redis密码
func WithLbAdapterSRedisPasswd(redisPasswd string) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.redisPasswd = redisPasswd
	}
}

//配置负载的redis最大连接数
func WithLbAdapterMaxActive(maxActive int) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.maxActive = maxActive
	}
}

//配置负载的redis最大空闲连接数
func WithLbAdapterMaxIdle(maxIdle int) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.maxIdle = maxIdle
	}
}

//配置负载的redis数据库
func WithLbAdapterDb(db int) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.db = db
	}
}

//配置负载的redis空闲连接数超时时间，单位秒
func WithLbAdapterIdleTimeOut(idleTimeOut time.Duration) LbAdapterCfgOpt {
	return func(lb *LbAdapter) {
		lb.idleTimeOut = idleTimeOut
	}
}

/*
Init(svc string, lbStrategy int, zkList []string, root string, routerType string, subRouterTypes []string, redisHost string, redisPasswd string, maxActive int, maxIdle int, db int, idleTimeOut time.Duration)
}*/

func (l *LbAdapter) Init(opts ...LbAdapterCfgOpt) error {
	if l.able == defaultLBABLE {
		return nil
	}
	for _, opt := range opts {
		opt(l)
	}

	l.lc = &lb_client.LbClient{}
	InitErr := l.lc.Init(
		lb_client.WithLbStrategy(l.lbStrategy),
		lb_client.WithZkList(l.zkList),
		lb_client.WithRoot(l.root),
		lb_client.WithRouterType(l.routerType),
		lb_client.WithSubRouterTypes(l.subRouterTypes),
		lb_client.WithRedisHost(l.redisHost),
		lb_client.WithRedisPassword(l.redisPasswd),
		lb_client.WithRedisDb(l.db),
		lb_client.WithRedisMaxActive(l.maxActive),
		lb_client.WithRedisMaxIdle(l.maxIdle),
		lb_client.WithRedisIdleTimeOut(l.idleTimeOut))
	if InitErr != nil {
		return fmt.Errorf("LbAdapter.Init fail InitErr:%v", InitErr)
	}
	return nil
}

//此svc暂时无用，暂不作修改
func (l *LbAdapter) Login(svc string, totalInst, idleInst, bestInst int32, param map[string]string) error {
	if l.able == defaultLBABLE {
		return nil
	}
	//启用初始化时提供的ip+port，暂时忽略此处的svc
	if LoginErr := l.lc.Login(l.svc, totalInst, idleInst, bestInst, param); LoginErr != nil {
		return fmt.Errorf("LbAdapter.Login fail LoginErr:%v", LoginErr)
	}
	return nil
}
func (l *LbAdapter) Update(totalInst, idleInst, bestInst int32) error {

	if l.able == defaultLBABLE {
		return nil
	}
	UpdateErr := l.lc.Upadate(totalInst, idleInst, bestInst)
	if UpdateErr != nil {
		return fmt.Errorf("LbAdapter.update fail UpdateErr:%v", UpdateErr)
	}
	return nil
}
func (l LbAdapter) LoginOut() error {
	if l.able == defaultLBABLE {
		return nil
	}
	LoginOutErr := l.lc.LoginOut()
	if LoginOutErr != nil {
		return fmt.Errorf("LbAdapter.LoginOut fail LoginOutErr:%v", LoginOutErr)
	}
	return nil
}
