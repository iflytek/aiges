## ※ 限制

- 客户端提供的接口只有几类，接口形式固定。这种方式对语音云业务以及需要业务编排的场景很有必要，但可能引起部分不适。
- 目前的负载策略有RoundRobin、hash、ConHash、topK、remote LB，其中remote LB依赖于独立的旁路负载实例
- 框架提供了裁剪模块的配置，尚不支持自定义模块替换框架依赖，如链路跟踪、配置中心、本地日志等
- 目前只有go语言版本的实现

## ※ 调用流程

- 初始化客户端：xsf.InitClient
- 准备请求数据：xsf.NewReq()
- 创建调用执行者：xsf.NewCaller 
- 发起调用：Caller.Call
- 处理结果：Res,error

## ※ 流程说明

### 1、初始化客户端

#### 接口

```
func InitClient(cname string,mode CfgMode, o ...CfgOpt) (*Client ,error )
```

#### 说明：

- 本接口是创建客户端，此实例生命周期可以和进程一致，error需要被处理，如果出错，则退出程序

#### 参数：

- cname: 客户端名称，这个参数在读取配置时会用作定位配置的selection。
- mode: 配置文件支持的模式，详见CfgMode定义时的注释，如果为Centre时，需要依赖配置中心服务的部署
- o...: 配置相关的其他属性,详见 CfgOpt定义时的注释

#### 备注

- 如果配置文件中启用了trace，还需要依赖日志跟踪服务。非必要。

### 2、准备请求数据

#### 主要接口以及结构

```
func NewReq()*Req
func (r *Req) Append(b []byte, desc map[string]string) 
func (r *Req) AppendData(data *Data)
func (r *Req) SetParam(k string, v string)
func (r *Req) GetParam(k string) (string, bool) 
func (r *Req) Session(s string) error
func (r *Req) AppendSession(k string, v string) 
func (r *Req) SetHandle(h string)
func (r *Req) Handle() (string)
func (r *Req) SetTraceID(t string)
func (r *Req) TraceID() string
func (r *Req) Data() []*Data

type DataMeta struct {
	Data []byte            
	Desc map[string]string 
}

type ReqData struct {
	Op    string            
	S     *Session          
	Param map[string]string 
	Data  []*DataMeta       
}

```

#### 结构说明：

- 结构设计基于：
	- 一次请求可能会携带多类数据，每一类数据可能都有自己的描述
	- 对于一些请求，操作需要明确，也会有相关的session记录，当次请求也需要全局性的描述。
- DataMeta：
	- data成员是携带用户请求数据的
	- Desc成员是用于描述改数据的参数
- ReqData：
	- Op当次请求的动作名，一般为对端函数名，构造消息体中可以不设置
	- S请求的session，Param用于描述请求的全局描述，Data携带请求数据
- Req结构基本上是基于以上接口的封装器

#### 接口说明：

- NewReq：申请请求消息对象
- Append/AppendData：向ReqData中追加Data，两种形式的接口
- SetParam/GetParam： 设置&获取参数
- Session/AppendSession：获取&追加自定session
- SetHandle/Handle：设置&获取服务句柄，这个句柄可用于会话模式下定位下游服务器，无状态的请求可以不关注。为框架生成，session默认字段
- SetTraceID/TraceID：设置&获取trace id，链路跟踪的索引，低频接口。为框架生成，session默认字段

### 3、创建调用执行者

#### 接口

```
func NewCaller(cli *Client)( *Caller)
```

#### 说明

- 创建调用的执行者，建议为局部变量。如果是无状态的服务接口，也可以作为全局对象使用

#### 参数

- cli: 客户端实例， InitClient返回。

### 4、发起调用

#### 接口

```
func (c *Caller)Call(service string, op string, r *Req, tm time.Duration)(s *Res,errcode int32, e error)
func (c *Caller)SessionCall(ss SessStat,service string, op string, r *Req, tm time.Duration )(s *Res,errcode int32, e error)
func (c *Caller)CallWithAddr(service string, op string,addr string, r *Req, tm time.Duration )(s *Res,errcode int32, e error)
```

### 说明

- Call：无状态服务的接口调用
- SessionCall：需要关注session上文的调用，该方法的CONTINUE请求，会根据Req session中的Handle定位服务器
- CallWithAddr：像制定的addr发起请求

### 参数

- service：服务名称
- op： 请求的对端动作名
- r: 请求的消息体
- tm：当次会话的超时时间
- ss：会话状态，SessStat详见SessStat
- addr:服务地址，格式：ip:port

### 返回值

- Res:响应而回的对象，携带响应数据
- errcode：携带错误码
- e：error信息

## 5、remote LB

### 接口

```
func (c *Caller)WithLBParams(lbaname string, busin string, ext map[string]string)
```

### 说明

- 用于设置remote LB参数，lb-mode为2时生效。需要依赖独立的LB服务

### 参数

- lbname：remote LB的服务名
- busin：需要接入负载的业务名
- ext：扩展参数，暂时未使用

## 3.6 处理返回结果以及报错

- todo:: 响应的处理，可以参照请求打包的逻辑类推


## 3.7 配置说明

- 配置为toml格式，也可以精简。配置样例如下

```
[xsf] #InitClient传入的cname
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
# 窗口bucket大小 ms
timeperslice = 10
# 窗口bucket数量
winsize = 100
# 概率
probability="80,16,4"
# 阈值
threshold= 1000
# 心跳间隔,单位ms
ping = 1000 

#测试目标服务配置，配置格式如下,注意分割符的差异
#业务1@ip1:port1;ipn:portn,业务2@ip2:port2;ipn:portn
taddrs="xsf-server@127.0.0.1:1997"

#trace相关配置
[trace]
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

#日志相关配置，默认为异步日志
[log]
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
```


