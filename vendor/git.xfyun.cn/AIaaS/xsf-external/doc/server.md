## 功能特性 ##

```
DNS寻址
配置中心
服务注册&发现
会话管理器
qps限流器
metrics
连接池
协程池
优雅上下线
管理后台
健康检测
监控&日志数据
```

## 基本接口 ##

```
//查询接口，可通过 xrpc-proxy 执行此查询
type monitor interface {
	Query(map[string]string, io.Writer)
}

//优雅的关闭，当接收到syscall.SIGINT, syscall.SIGKILL时，会回调这接口
type Killer interface {
	Closeout()
}

//服务自检接口，可通过 xrpc-proxy 执行
type HealthChecker interface {
	Check() error
}

//用户业务逻辑接口
type UserInterface interface {
	Init(*ToolBox) error
	FInit() error
	Call(*Req, *Span) (*Res, error)
}
```
- monitor为数据查询接口，用于特殊的数据的导出，如服务实现了Query接口后，可通过xrpc-proxy远程提取
- Killer为信号监听接口，如服务实现了Killer接口后，那么当框架收到Interrupt、SIGINT、SIGKILL、SIGTERM信号时，会主动调用此接口，以做相关的资源回收等收尾工作
- HealthChecker为健康检测接口，用户需实现此接口，框架是根据Check() error的返回值判断服务是否正常运行
- UserInterface为用户实际的服务接口，分别有Init(*ToolBox) error、Finit() error、Call(*Req, *Span) (*Res, error)三个方法，代表服务的初始化、逆初始化、服务调用接口

## 集成步骤 ##

```
※ 实现相关接口
※ 入口 main 包内执行 server.run 运行即可
```

## 集成示例 ##

- 接口实现

```
package main

import (
	"errors"
	"fmt"
	"io"

	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"strconv"
	"time"
)

type monitor struct {
}

func (m *monitor) Query(in map[string]string, out io.Writer) {
	_, _ = out.Write([]byte( fmt.Sprintf("%+v", in)))
}

//当接收到syscall.SIGINT, syscall.SIGKILL时，会回调这接口
type killed struct {
}

func (k *killed) Closeout() {
	fmt.Println("server be killed.")
}

type healthChecker struct {
}

//服务自检接口，cmdserver用
func (h *healthChecker) Check() error {
	return errors.New("this is check function from health check")
}

//用户业务逻辑接口
type server struct {
	tool *xsf.ToolBox
}

//业务初始化接口
func (c *server) Init(toolbox *xsf.ToolBox) error {
	fmt.Println("begin init")
	c.tool = toolbox
	{
		xsf.AddKillerCheck("server", &killed{})
		xsf.AddHealthCheck("server", &healthChecker{})
		xsf.StoreMonitor(&monitor{})
	}
	fmt.Println("server init success.")
	return nil
}

//业务逆初始化接口
func (c *server) FInit() error {
	fmt.Println("user logic FInit success.")
	return nil
}

//业务服务接口
func (c *server) Call(in *xsf.Req, span *xsf.Span) (response *utils.Res, err error) {
	switch in.Op() {
	case "ssb":
		return c.ssbRouter(in)
	case "auw":
		return c.auwRouter(in)
	case "sse":
		return c.sseRouter(in)
	case "req":
		return c.reqRouter(in)
	default:
		break
	}
	return c.unknown(in)
}

func (c *server) unknown(in *xsf.Req) (*utils.Res, error) {
	fmt.Printf("the op -> %v is not supported.\n", op)
	res := xsf.NewRes()
	res.SetHandle(in.Handle())
	res.SetError(1, fmt.Sprintf("the op -> %v is not supported.", op))
	res.SetParam("intro", "received data")
	res.SetParam("op", "illegal")
	return res, nil
}

func (c *server) reqRouter(in *xsf.Req) (*utils.Res, error) {
	c.tool.Log.Debugw("in process", "op", "req")
	res := xsf.NewRes()
	res.SetHandle(in.Handle())
	res.SetParam("intro", "received data")
	res.SetParam("op", "req")
	res.SetParam("ip", c.tool.NetManager.GetIp())
	res.SetParam("port", strconv.Itoa(c.tool.NetManager.GetPort()))
	data := xsf.NewData()
	data.SetParam("intro", "for test")
	data.Append([]byte("test data"))
	res.AppendData(data)
	return res, nil
}

func (c *server) sseRouter(in *xsf.Req) (*utils.Res, error) {
	c.tool.Log.Debugw("in process", "op", "sse")
	res := xsf.NewRes()
	res.SetHandle(in.Handle())
	c.tool.Cache.DelSessionData(in.Handle())
	_ = c.tool.Cache.Update()
	{
		res.SetParam("intro", "received data")
		res.SetParam("op", "sse")
		res.SetParam("ip", c.tool.NetManager.GetIp())
		res.SetParam("port", strconv.Itoa(c.tool.NetManager.GetPort()))
	}
	return res, nil
}

func (c *server) auwRouter(in *xsf.Req) (*utils.Res, error) {
	c.tool.Log.Debugw("in process", "op", "auw")
	res := xsf.NewRes()
	res.SetHandle(in.Handle())
	if _, GetSessionDataErr := c.tool.Cache.GetSessionData(in.Handle()); GetSessionDataErr != nil {
		res.SetError(1, fmt.Sprintf("GetSessionData failed. ->GetSessionDataErr:%v", GetSessionDataErr))
	}
	{
		res.SetParam("intro", "received data")
		res.SetParam("op", "auw")
		res.SetParam("ip", c.tool.NetManager.GetIp())
		res.SetParam("port", strconv.Itoa(c.tool.NetManager.GetPort()))
	}
	return res, nil
}

func (c *server) ssbRouter(in *xsf.Req) (*utils.Res, error) {
	c.tool.Log.Debugw("in process", "op", "ssb")
	res := xsf.NewRes()
	res.SetHandle(in.Handle())
	sessionCb := func(sessionTag interface{}, svcData interface{}, exception ...xsf.CallBackException) {
		c.tool.Log.Infow("this is callback function", "timestamp", time.Now(), sessionTag, in.Handle())
	}
	SetSessionDataErr := c.tool.Cache.SetSessionData(in.Handle(), "svcData", sessionCb)
	if nil != SetSessionDataErr {
		res.SetError(1, fmt.Sprintf("Set %s failed. ->SetErr:%v ->addr:%v",
			in.Handle(), SetSessionDataErr, fmt.Sprintf("%v:%v", c.tool.NetManager.GetIp(), c.tool.NetManager.GetPort())))
	} else {
		_ = c.tool.Cache.Update()
	}
	return res, nil
}

```

- oprouter

```
package main

import (
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"time"
)

func generateOpRouter() *xsf.OpRouter {
	router := &xsf.OpRouter{}
	router.Store("op", func(in *xsf.Req, span *xsf.Span, tool *xsf.ToolBox) (*utils.Res, error) {
		res := xsf.NewRes()
		res.SetHandle(in.Handle())
		fmt.Printf("info:this is op operator. -> timestamp:%v,Handle:%v\n", time.Now(), in.Handle())
		return res, nil
	})
	return router
}
```

- 接口调度

```
package main

import (
	"flag"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"log"
	"sync"
)

const (
	cfgName      = "server.toml"
	project      = "3s"
	group        = "3s"
	service      = "xsf-server"
	version      = "x.x.x"
	apiVersion   = "1.0.0"
	cachePath    = "xxx"
	companionUrl = "http://10.1.87.70:6868"
)

func init() {
	flag.Parse()
}
func main() {

	//定义一个服务实例
	var serverInst xsf.XsfServer

	//定义相关的启动参数
	/*
		1、CfgMode可选值Native、Centre，native为本地配置读取模式，Centre为配置中心模式，当此值为-1时，表示有命令行传入
		2、CfgName 配置文件名
		3、Project 配置中心用 项目名
		4、Group 配置中心用 组名
		5、Service 配置中心用 服务名
		6、Version 配置中心用 配置版本名
		7、CompanionUrl 配置中心用 配置中心地址
	*/
	bc := xsf.BootConfig{
		CfgMode: utils.Native,
		CfgData: xsf.CfgMeta{
			CfgName:      cfgName,
			Project:      project,
			Group:        group,
			Service:      service,
			Version:      version,
			ApiVersion:   apiVersion,
			CachePath:    cachePath,
			CompanionUrl: companionUrl}}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		/*
			1、启动服务
			2、若有异常直接报错，注意需用户自己实现协程等待
		*/
		if err := serverInst.Run(
			bc,
			&server{},
			xsf.SetOpRouter(generateOpRouter())); err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()
}
```

- 代码里有较为详尽的注释，如仍有疑问请联系sqjian@iflytek.com

## 最小配置 ##
```
#服务自身的配置
#注意此section名需对应bootConfig中的service
[test]#已做缺省处理
host = "127.0.0.1"#若host为空，则取netcard对应的ip，若二者均为空，则取hostname对应的ip
port = 1997 #不填则采用随机端口


#trace日志所用
[trace]#已做缺省处理
able = 0 #缺省1

[log]#已做缺省处理
level = "info" #缺省warn
file = "serverNative.log" #缺省xsfs.log  

#sonar日志所用
[sonar]#已做缺省处理,此section如不传缺省启用
able = 0 #缺省1
```

## 完整配置详解 ##

```
#服务自身的配置
#注意此section名需对应bootConfig中的service
[xsf-server]#已做缺省处理
host = "127.0.0.1"#若host为空，则取netcard对应的ip，若二者均为空，则取hostname对应的ip
#host = "127.0.0.1"#若host为空，则取netcard对应的ip，若二者均为空，则取hostname对应的ip
#netcard = "eth0"
port = 1997 #不填则采用随机端口
#reuseport = 1 #缺省0
finder = 0 #缺省0
maxreceive = 4 #能收取的最大消息包大小，单位MB，缺省16MB
maxsend = 4  #能发送的最大消息包大小，单位MB，缺省16MB
conn-rbuf = 4 #连接读缓冲区
conn-wbuf  = 4 #连接写缓冲区
keepalive = 20 #启用keepalive check，值为探测的时间间隔，单位毫秒，缺省不启用keeplive
keepalive-timeout = 3 #在keeplive启用的条件下，表示heatbeat的超时时间，单位毫秒，缺省1000ms
vcpu = 1 #vcpu采集开关

[metrics]
#参数齐则开启metrics
idc = "dx"
sub = "xsf"
cs = "3s"
timePerSlice = 1000 #滑动窗口bucket大小，单位毫秒
winSize = 10 #窗口大小

#trace日志所用
[trace]#已做缺省处理,此section如不传缺省启用
#trace收集服务的地址
host = "172.16.51.3" #缺省127.0.0.1
#trace收集服务的端口
port = 4546 #缺省4545
#trace服务的协程数
backend = 1 #缺省4
#是否将日志写入到远端
deliver = 0 #缺省1
#是否将日志 落入磁盘
dump = 1 #缺省0
#是否禁用trace
able = 1 #缺省1
watch = 0 #缺省1
watchport = 1234
spill = "./log/spill"  #缺省/log/spill
buffer =    10 #缓冲,缺省1000000
batch =     100     #大小,缺省100
linger =    5       #每几秒检查一次buffer,缺省5
bcluster ="xxx" #业务集群标识，缺省3s
idc      ="yyy" #IDC标识位，缺省dz

[log]#已做缺省处理
level = "info" #缺省warn
file = "serverNative.log" #缺省xsfs.log
#日志文件的大小，单位MB
size = 300 #缺省10
#日志文件的备份数量
count = 3 #缺省10
#日志文件的有效期，单位Day
die = 3 #缺省10
#缓存大小，单位条数,超过会丢弃
cache = 100000 #缺省-1，代表不丢数据，堆积到内存中
#批处理大小，单位条数，一次写入条数（触发写事件的条数）
batch = 160#缺省16*1024
#异步日志
async = 0 #缺省异步
#是否添加调用行信息
caller = 1 #缺省0
wash = 60 #写入磁盘的缺省时间

[lb] #已做缺省处理,此section如不传缺省不启用
able           = 0 #缺省0
lbStrategy     = 0                                                                          #负载策略(必传)
zkList         = ["192.168.86.60:2191", "192.168.86.60:2192", "192.168.86.60:2190"]         #zk列表(必传)
root           = "/"                                                                         #根目录
routerType     = "xsf"                                                                      #路由类型(如：svc)(必传)
subRouterTypes = ["xsf_gray", "xsf_hefei"]                                                  #子路由类型列表(如:["svc_gray","svc_hefei"])
redisHost      = "192.168.86.60:6379"                                                       #redis主机(必传)
redisPasswd    = ""                                                                         #redis密码
maxActive      = 100                                                                        #redis最大连接数
maxIdle        = 50                                                                         #redis最大空闲连接数
db             = 0                                                                          #redis数据库
idleTimeOut    = 10                                                                         #redis空闲连接数超时时间，单位秒

[lbv2] #已做缺省处理,此section如不传缺省不启用
tm = 1000 #缺省1000，批量上报的超时时间，单位毫秒
backend = 100#上报的的协程数，缺省4
finderttl = 100 #更新本地地址的时间，通过访问服务发现实现，缺省一分钟
lbname = "lbv2"
apiversion = "1.0"
able = 1
sub = "svc"
subsvc = "sms,sms-16k"
task = 10 #任务队列长度

#测试目标服务配置，配置格式如下,注意分割符的差异
#业务1@ip1:port1;ipn:portn,业务2@ip2:port2;ipn:portn
conn-timeout = 100
conn-pool-size = 4         #rpc连接池数量。缺省4
lb-mode= 0  #0禁用lb,2使用lb。缺省0
lb-retry = 0
taddrs="lbv2@10.1.87.68:1999"


[fc]#flowControl 包括sessionManager和qpsLimiter
#限流器的类型，若所填值非sessionManager和qpsLimiter或者没填，那么限流器不会初始化
able = 1 #缺省为0
router = "sessionManager"   #路由字段，可选项为sessionManager和qpsLimiter
max = 100                   #会话模式时代表最大的授权量，非会话模式代表间隔时间里的最大请求数
ttl = 30000                     #会话模式代表会话的超时时间，非会话模式代表有效期（间隔时间），缺省15000ms
best = 10                   #最佳授权数
roll = 1000                    #sessionManager内部遍历超时session的时间间隔  缺省1000ms
report =1000                   #当策略为0即定时上报时，此为上报时间间隔 缺省1000ms，当策略为1即根据授权波动变化是，此值代表检查波动值的时间间隔
strategy = 2                #可选值为0、1、2（缺省为0），0.代表定时上报(v1)；1.根据授权范围上报(v1)；2.基于hermes（v2）
wave = 2                  #波动值，当授权数变化值大于等于该值时，出发触发上报行为,缺省10

#sonar日志所用
[sonar]#已做缺省处理,此section如不传缺省启用
#trace收集服务的地址
host = "172.16.51.3" #缺省127.0.0.1
#trace收集服务的端口
port = 4546 #缺省4545
#trace服务的协程数
backend = 1 #缺省4
#是否将日志写入到远端
deliver = 1 #缺省1
#是否将日志 落入磁盘
dump = 1 #缺省0
#是否禁用trace
able = 0 #缺省1
ds = "vagus" #缺省vagus
```

## 服务状态监控 ##

> 需配合 xrpc-proxy，见 xrpc-proxy 仓库

**※ 常规查询**

```
查修服务状态信息：
curl "http://xxx.xxx.xxx:xxx/index.html?svc=xxx&cmd=status"

查修服务健康信息：
curl "http://xxx.xxx.xxx:xxx/index.html?svc=xxx&cmd=health"

查修服务协程信息：
curl "http://xxx.xxx.xxx:xxx/index.html?svc=xxx&cmd=goroutine"

查修服务堆栈信息：
curl "http://xxx.xxx.xxx:xxx/index.html?svc=xxx&cmd=heap"

查修服务线程信息：
curl "http://xxx.xxx.xxx:xxx/index.html?svc=xxx&cmd=threadcreate"

查修服务内存信息：
curl "http://xxx.xxx.xxx:xxx/index.html?svc=xxx&cmd=block"

查修服务垃圾收集信息：
curl "http://xxx.xxx.xxx:xxx/index.html?svc=xxx&cmd=gcsummary"

```

**※ 自定义查询**

```
查询 svc 对应节点的授权信息
curl "http://xxx.xxx.xxx:xxx/index.html?svc=svc&cmd=query&svc=svc"

查询 svc & subsvc 的授权信息
curl "http://xxx.xxx.xxx:xxx/index.html?svc=svc&cmd=query&svc=svc&subsvc=sms"

```


**※ 查询指定节点**

```
查询 addr 节点的信息
curl "http://xxx.xxx.xxx:xxx/index.html?addr=xxx.xxx.xxx:xxx&cmd=xxx"

```

## log格式注解 ##

- 日志示例
- 
```
{"level":"info","ts":"2018-01-12T14:18:38.675+0800","caller":"server/sessionmanager.go:277","msg":"enter timer 1 ->interval:3s,mapLen:0","pid":12644}
```
- level可选值为如下四种
	- debug
	- info
	- warn
	- error
- ts
	- 日志的时间格式采用国际标准化组织的国际标准ISO 8601，称为《数据存储和交换形式·信息交换·日期和时间的表示方法》
- caller
	- 日志被调用的位置
- msg
	- 用户真实的日志数据
- pid
	- 进程的pid，部署docker的pid

## 服务地址 ##

```
host = "0.0.0.0"	//可指定域名或者ip
netcard = "eth0"	//读取的网卡
port = 50061		//监听的端口
```

1. 如果host存在，则取host对应的ip（dns寻址），若host为ip，则去host  
1. 如果host不存在，netcard存在，则取netcard对应的ip  
1. 如果host、和netcard都没传，则  
	- 去本机hostname  
	- 调用LookupHost查找hostname对应的ip  
