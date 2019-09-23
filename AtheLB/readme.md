## 功能 ##
> intro:全局负载，基于收集的信息做负载，并按特定策略返回相应的节点

- 当前实现的策略有load、loadMini、poll
	- load：基于授权数的高性能版本，利用率偏低（单节点大于95%，小于100%）
	- loadMini：基于授权数的加锁版本（利用率100%）
	- poll：轮训

- 部分功能列举如下

- 数据同步
	- 同步DB数据并至本地
	- 同步新节点数据至DB

- 节点上报
	- 引擎段的分配
	- 节点生命周期的维护

- 节点请求
	- 根据负载返回响应节点
	- 根据uid返回响应节点

- 个性化索引
	- 同步rmq至本地
	- 根据策略路由至后端节点

- 监控相关
	- 监控全量节点
	- 监控局部节点
	

## 接口 ##

### 上报接口 ###

```
req.SetParam("svc", "svc)				//大业务名
req.SetParam("subsvc", "sms")			//子业务名
req.SetParam("addr", "x.x.x.x:xxxx")	//上报节点的服务地址
req.SetParam("total", "xxx")			//节点的总授权量
req.SetParam("best", "xxx)				//节点的最佳授权量
req.SetParam("idle", "xxx")				//节点的空闲授权量
req.SetParam("live", "xxx")				//节点的存活状态(若此值为0，则代表节点主动下线)
```

### 请求接口 ###

```
req.SetParam("nbest", "xxx")			//期望获取的节点数
req.SetParam("svc", "svc")				//大业务名
req.SetParam("subsvc", "svc")			//子业务名
req.SetParam("all", "1")				//是否返回所有的节点
req.SetParam("uid", "xxxx")				//universal id 全局id（个性化所需）
```

### 基本流程 ###

- 引擎上报相关信息至lb实例，通过xsf内置支持完成
- lb负责管理与调度引擎实例，包括接受业务节点请求，回收失效引擎实例等
- 业务节点请求lb实例获取相应的引擎节点

**参考配置**

```
[AtheLB] #已做缺省处理,此section如不传缺省不启用
#缺省1000，批量上报的超时时间，单位毫秒
tm = 1000 
#上报的的协程数，缺省4
backend = 100
#更新本地地址的时间，通过访问服务发现实现，缺省一分钟
finderttl = 100 
lbname = "lbv2"
apiversion = "1.0"
able = 1
sub = "svc"
subsvc = "sms,sms-16k"
#任务队列长度
task = 10 #任务队列长度
#链接超时
conn-timeout = 2000
#连接池大小
conn-pool-size = 100         #rpc连接池数量。缺省4
#负载均衡模式
lb-mode= 0  #0：无权重轮询,1：普通哈希，2：hermes，3、一致性哈希，4：业务负载，缺省0
#负载超时
lb-timeout=500
#当负载均衡失败时，重试几次
lb-retry = 2
#连接生命周期，单位毫秒
conn-life-cycle = 120 * 1000
# 连接重试次数
conn-retry=3
#连接读缓冲区
conn-rbuf=4*1-24
#连接写缓冲区
conn-wbuf=32*1024*1024
#启用keepalive check，值为探测的时间间隔，单位毫秒，缺省不启用keeplive
keepalive = 10 
#在keeplive启用的条件下，表示heatbeat的超时时间，单位毫秒，缺省1000ms
keepalive-timeout = 3 
```

## 构建&部署 ##

### 构建 ###

- golang版本要求1.9以上

- 设置好相应的GOPATH(如:export GOPATH=lb-external)

- 下载 tag版本并进入lb-external\src\lbv2 

- 在lb目录下,执行cmd：go build -v -o lbv2


### 部署 ###

- 本地模式

```
# -m 制定运行方式(0-本地配置,1-配置中心)    
# -c 配置文件名(Native模式时，用本地的配置文件，Center模式时使用配置中心的配置文件)    
# -s 启动的服务名(注意:配置文件，配置中心的服务名需一致) 
#!/usr/bin/env bash
./lbv2 -v
./lbv2 -m 0 -c lbv2.toml  -s lbv2
```


- 配置中心
	
```
# -m 制定运行方式(0-本地配置,1-配置中心)    
# -c 配置文件名(Native模式时，用本地的配置文件，Center模式时使用配置中心的配置文件)    
# -p 项目名    
# -s 启动的服务名(注意:配置文件，配置中心的服务名需一致)  
# -u 配置中心地址    
# -g 配置项目组

#!/usr/bin/env bash
./lbv2 -m 1 -c lbv2.toml -p 3s -s lbv2 -u http://10.1.86.223:9080 -g 3s
```

- 配置文件（按需调整）

```
[lbv2]
#host="0.0.0.0"                 #若host为空，则取netcard对应的ip，若二者均为空，则取hostname对应的ip
#netcard = "eth0"
port = 1995                       #指定端口
reuseport = 0                     #缺省0
cmdserver = 0                     #缺省0
finder = 0                        #使用服务发现,缺省0

[bo]#Business Object
#pprof = 9999 #缺省关闭pprof
strategy = 0                      #策略 0：load  1: poll 2: ise
threshold = 90                    #服务阈值，(totalInst - idleInst) * 100 / totalInst
ticker = 5000                    #清除无效节点的扫描周期，单位毫秒
svc = "svc"                       #后续可优化
defsub = "xxx"                    #兜底路由
ttl = 3000                      #节点的生存周期,单位毫秒
#preauth = 8  #预授开关，缺省预授-1
rmqable = 0 #缺省1，开启
rmqaddrs = "172.16.154.26:10700,172.16.154.26:10800,172.16.154.26:10900"
#rmqaddrs = "10.1.87.19:10800,10.1.87.68:10800"
rmqtopic = "mc_svc_test"
rmqgroup = "group"
rmqticker = 1000 #消费rmq的时间间隔
consumer = 3
monitor = 1000 #监控数据的更新周期，单位毫秒
nodedur = 10000 #异常节点的保存周期，单位毫秒

[db]
able = 0 #缺省1
batch = 2#每次最多从数据库拉取的数据条数，缺省10000
baseurl   = "http://172.16.154.235:808/ws"
caller    = "xfyun"
callerkey = "12345678"
timeout   = 3000 #毫秒
token     = "100IME"
version   = "db-service-v3-3.0.0.1001"
idc       = "bj"
schema    = "ifly_cp_msp_balance"
table     = "seg_list_lbv2"
dbtime    = 1 #清理db的时间间隔，缺省172800s，单位s
rctime    = 1 #第一次访问db失败后，后续重新访问数据的时间间隔，缺省600s

[dc]#通知ats拉取个性化资源所用
#测试目标服务配置，配置格式如下,注意分割符的差异
#业务1@ip1:port1;ipn:portn,业务2@ip2:port2;ipn:portn
conn-timeout = 100
conn-pool-size = 4         #rpc连接池数量。缺省4
lb-mode= 0  #0禁用lb,2使用lb。缺省0
lb-retry = 0
#taddrs="lbv2@10.1.87.61:9095"

[trace]
host = "172.16.51.3"              #trace收集服务的地址,缺省127.0.0.1
port = 4546                       #trance的端口号,缺省4545
backend = 1                       #trace服务的协程数,缺省4
deliver = 1                       #是否将日志写入到远端,缺省1
dump = 1                          #是否将日志 落入磁盘,缺省0
able = 0                          #是否禁用trace,缺省1

[log]
level = "debug"                  #日志文件类型,缺省warn
file = "log/lb.log"              #日志文件名,缺省xsfs.log
size = 100                         #日志文件的大小,单位MB,缺省10
count = 10                        #日志文件的备份数量,缺省10
die = 3                          #日志文件的有效期,单位Day,缺省10
cache = -1                       #缓存大小,单位条数,超过会丢弃,(缺省-1，代表不丢数据，堆积到内存)
batch = 1600                     #批处理大小,单位条数,一次写入条数（触发写事件的条数）
async = 0                        #异步日志,缺省异步
caller = 1                       #是否添加调用行信息,缺省0

[sonar]
host = "10.1.86.60"              #trace收集服务的地,缺省127.0.0.1
port = 4546                      #trace收集服务的端口,缺省4545
backend = 1                      #trace服务的协程数,缺省4
deliver = 1                      #是否将日志写入到远端,缺省1
dump = 1                         #是否将日志写入磁盘,缺省0
able = 0                         #是否禁用trace,缺省1
rate = 5000                      #上报频率,单位毫秒
ds = "vagus"                     #缺省vagus
```

## 自测 ##

[hermes_test_data.xlsx](hermes_test_data.xlsx)

## FAQ ##
> waiting...
