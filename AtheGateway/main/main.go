package main

import (
	"fmt"
	"common"
	"common/ratelimit2"
	"conf"
	"schemas"
	"server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var gd  ws.Upgrader

//sid生成器
var (
	routerMap    = sync.Map{}
	sidGenerator = &utils.SidGenerator2{}
	sessionPool = NewSessionPool()
)
const (
	DEFAULT_VERSION = "1.0"
	KEY_SESSION = "webgate-session"
)

func main() {
	//初始化配置

	conf.InitConf()

	gd = ws.Upgrader{
		HandshakeTimeout: time.Duration(conf.Conf.Session.HandshakeTimeout) * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//初始化XsfLog
	if err:=common.InitXsfLog(conf.Conf.Log.File, conf.Conf.Log.Level, conf.Conf.Log.Size, conf.Conf.Log.Count, conf.Conf.Log.Caller, conf.Conf.Log.Batch,conf.Conf.Log.Asyn);err !=nil{
		panic(err)
	}

	//初始化auth Mysql
	//auth.InitDB(conf.Conf.Mysql.Mysql)
	//初始化session缓存
	server.InitCache(time.Duration(conf.Conf.Session.ScanInterver) * time.Second)

	common.InitSidGenerator(conf.Conf.Server.Host,conf.Conf.Server.Port,conf.Conf.Xsf.Location)
	server.Sidgenerator = sidGenerator
	//初始化sid生成器
	//sidGenerator.Init(conf.Conf.Xsf.Location, conf.Conf.Server.Host, conf.Conf.Server.Port)
	//启动XSF服务端
	server.StartXsfServer("xsf.toml")
	// 初始化XSF客户端
	err := server.InitXsfClient("xsf.toml")

	if err != nil {
		fmt.Printf("InitXsfClient is error:%s\n", err)
		common.Logger.Errorf("InitXsfClient is error:%s", err)
		return
	}

	initConnLimit() //初始化限流配置


	//启动路由
	router := initGin()
	//router.Use(common.Loggers())
	//监控路由
	router.GET("/monitor/:option",server.MonitorHandler)
	router.GET("/conns/:appid",server.HandlerGetConncection)
	router.POST("/conns/:appid/kill/:remain",server.HandlerKill)

	//服务路由
	authorized := router.Group("")
	authorized.Use(authRequired)
	LoadRouteMap(authorized)
	conf.AddConfigChangerHander(func(s string, bytes []byte)bool {
		if schemas.IsSchemaFile(s){
			err:=schemas.LoadRoteMapping(bytes)
			if err !=nil{
				conf.ConsoleError("reload schema error:"+s+err.Error())
				return false
			}
			LoadRouteMap(authorized)
		}

		if s==conf.APP_CONFIG{
			if err:=common.InitXsfLog(conf.Conf.Log.File, conf.Conf.Log.Level, conf.Conf.Log.Size, conf.Conf.Log.Count, conf.Conf.Log.Caller, conf.Conf.Log.Batch,conf.Conf.Log.Asyn);err!=nil{
				return false
			}
			LoadRouteMap(authorized)
		}
		return true
	})
	// 开启 pprof
	ginpprof.Wrapper(router)

	err = router.Run(conf.Conf.Server.Host + ":" + conf.Conf.Server.Port)
	if err != nil {
		fmt.Printf("start gin server error:%s\n", err)
		panic(err)
		common.Logger.Errorf("Router.Run  is error:%s", err)
		return
	}


}

//请求认证
func authRequired(context *gin.Context) {
	context.Next()
}
var connId int64 = 0
//请求处理
func mainHandler(context *gin.Context) {
	// 获取客户端的版本
	var version,ok = context.GetQuery("version")
	if !ok{
		version  = DEFAULT_VERSION
	}
	//
	route:=context.Request.URL.Path
	cal,_:=routerMap.Load(route)
	call:=cal.(string)
//	key:=call+route+version
	key:=schemas.GetMappingKey(call,route,version)
	mp:=schemas.GetMappingByKey(key)
	if mp==nil{
		context.AbortWithStatusJSON(http.StatusNotFound,gin.H{"message":"cannot not found route:"+route+" version:"+version})
		return
	}
	//升级http协议为websocket
	conn, err := gd.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		common.Logger.Errorf("client(%s) request, upgrade protocol from http to websocket failed :%s", context.ClientIP(), err.Error())
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	common.Logger.Infof("client(%s)request, ability:%s", context.ClientIP(), call)
	t:=server.NewTrans()
	t.Mapping = mp
	t.ConnId = atomic.AddInt64(&connId,1)
	t.Conn = conn
	t.Ctx = context
	t.Call=call
	t.ClientIp = context.ClientIP()
	appid:=context.GetHeader("X-Consumer-Username")
	t.ParentAppid = appid
	if appid!=""{
		server.ConnTransManager.SetTrans(t.ConnId,t)
	}
	if kep,ok:=context.GetQuery("stream_mode");ok && kep == "multiplex"{
		t.KeepAlive = true
	}
	context.Set(KEY_SESSION,t)
	t.Do()
	t.Close()
	//当返回结果需要排序时打开排序的管道

}

func LoadRouteMap(g *gin.RouterGroup) {
	//g.GET("/",mainHandler)
	routers := schemas.RouteMappingCache
	routers.Range(func(key, value interface{}) bool {
		v:=value.(*schemas.RouteMapping)
		if r,ok:=routerMap.Load(v.Route);ok && r.(string)!=""{
			routerMap.Store(v.Route,v.Service)
			return true
		}
		routerMap.Store(v.Route,v.Service)
		g.GET(v.Route, mainHandler)
		return true
	})
	fmt.Println("routers:", routerMap)
}

func initGin()*gin.Engine  {
	g:=gin.New()
	if conf.Conf.Server.Mode == "debug"{
		gin.SetMode(gin.DebugMode)
	}else{
		gin.SetMode(gin.ReleaseMode)
		//g.Use(recoveryHandler) // 生产环境启用recovery()
	}
	g.Use(rateLimitHandler) // 开启限流
	return g
}
//限流
func rateLimitHandler(ctx *gin.Context)  {
	appid:=ctx.GetHeader("X-Consumer-Username")
	//限制连接数
	limitConfig:=ratelimit2.ConfigCacheInstance().GetConfig(appid)
	//当appid 为空时不做流控
	if limitConfig !=nil && appid !=""{
		currentConn:=server.ConnTransManager.GetCount(appid)
		//超过限制
		if currentConn>=limitConfig.ConnLimit{
			ctx.AbortWithStatusJSON(http.StatusForbidden,gin.H{"message":"max connections exceeded"})
			return
		}
	}
	ctx.Next()
}

func recoveryHandler(ctx *gin.Context)  {
	defer func() {
		if err:=recover();err !=nil{
			common.Logger.Errorf("server come across unexpected error:%v",err)
			if s,ok:=ctx.Get(KEY_SESSION);ok{
				if t,ok:=s.(*server.Trans);ok{
					t.Write(server.FrameResponse{Code:int(server.ErrorServerError),Message:"unexpected server error"})
					t.Close()
				}
			}
			ctx.AbortWithStatus(http.StatusBadGateway) //
		}
	}()
	ctx.Next()  //保证defer 函数体在所有的逻辑执行完成后再执行
}


//session pool

type SessionPoll struct {
	sync.Pool
}

func (p *SessionPoll)GetSession()*server.Session  {
	return p.Get().(*server.Session)
}

func (p *SessionPoll)PutSession(session *server.Session)  {
	p.Put(session)
}

func NewSessionPool()*SessionPoll  {
	p:=&SessionPoll{}
	p.New = func() interface{} {
		return server.NewSession()
	}
	return p
}


//初始化限流配置
func initConnLimit()  {
	conf.AddConfigChangerHander(func(s string, bytes []byte) bool {
		if s == conf.LimitConf{
			// 收到配置推送，重新加载限流配置，并根据配置杀掉对应appid多余的连接
			ratelimit2.LoadConfigCache(bytes)
			ratelimit2.ConfigCacheInstance().Range(func(apppid string, config *ratelimit2.Config) bool {
				server.ConnTransManager.Kill(config.Appid,config.ConnLimit)
				return true
			})

		}
		return true
	})
}

func InitAdmin()  {
	go initAdminApi()
}

func initAdminApi()  {
	g:=gin.New()
	//reconvery while panic occered
	g.Use(func(context *gin.Context) {
		if err:=recover();err !=nil{
			context.AbortWithStatusJSON(502,gin.H{"message":"unexpected error"})
			return
		}
		context.Next()
	})
	// authorization
	g.Use(func(context *gin.Context) {
		//todo
		context.Next()
	})

	g.GET("/monitor/:option",server.MonitorHandler)
	g.GET("/conns/:appid",server.HandlerGetConncection)
	g.POST("/conns/:appid/kill/:remain",server.HandlerKill)
	err:=g.Run(conf.Conf.Server.AdminListen)
	if err !=nil{
		fmt.Println("admin listen error:",err.Error())

	}


}