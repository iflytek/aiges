# Xtest测试框架

## 一、 目录总览

[TOC]

## 二、代码目录

### 2.1、根目录

#### Ⅰ.  xtest.go

> xtest为项目的入口，拥有main函数与linearCtl函数，
>
> main函数负责解析参数、初始化客户端，并实现多协程异步记录测试日志
>
> ```go
> // Xtest 结构体
> type Xtest struct {
>	r   request.Request // 与请求有关的对象
>	cli *xsfcli.Client // rpc客户端
>}
> // NewXtest 传入rpc客户端和全局配置，实例化Xtest对象
> func NewXtest(cli *xsfcli.Client, conf _var.Conf) Xtest
> // Run 启动测试流程
> func (x *Xtest) Run()
> // linearCtl 函数负责并发线性增长控制,防止瞬时并发请求冲击
> func (x *Xtest) linearCtl()
> ```

### 2.2、analy文件夹

#### Ⅰ. errdist.go： 运行错误信息有关的结构体与函数定义

> ```go
> type ErrInfo struct {
> 	errCode int
> 	errStr  error
> }
> 
> // errDistAnalyser 记录错误数据相关的信息
> type errDistAnalyser struct {
> 	errCnt   map[int]int64 // map[error]count 错误计数
> 	errDsc   map[int]error // 错误描述
> 	errTmp   []ErrInfo     // 临时存储区,用于channel满阻塞的极端场景;
> 	errMutex sync.Mutex	   // 互斥锁
> 	errChan  chan ErrInfo // errorInfo管道
> 	swg      sync.WaitGroup // 并发同步原语
> 	log      *utils.Logger // 日志写入
>     ErrAnaDst string // 错误日志存储地址
> }
> 
> // 初始化 errDistAnalyser 结构体成员，并启动一个count计数协程计算errCnt
> func (eda *errDistAnalyser) Start(clen int, logger *utils.Logger) 
> // 将错误信息发送到errChan中，如果errchan满了，则发送到临时存储区errTmp
> func (eda *errDistAnalyser) PushErr(code int, err error)
> // 关闭errChan并等待swg 并发协程执行完成
> func (eda *errDistAnalyser) Stop()
> // 首先读取errChan错误信息，再读取errTmp临时存储区错误信息，分别统计不同错误的次数，并落盘
> func (eda *errDistAnalyser) count()
> // 错误分布数据落盘
> func (eda *errDistAnalyser) dumpLog()
> ```

#### Ⅱ. perf.go： 性能指标统计有关的函数定义

> ```go
> // 计时类型,用于控制计时开关
> const (
>  FIFOPERF = 1 << iota // 首结果耗时
>  LILOPERF             // 尾结果耗时
>  SESSPERF             // 会话耗时
>  INTFPERF             // 接口耗时
> )
> 
> // 计时定点,用于标记计时位置
> const (
>  pointCreate int = 1 << iota
>  pointUpBegin
>  pointDownBegin
>  pointUpEnd
>  pointDownEnd
>  pointDestroy
> )
> 
> type PerfDetail struct {
>  cTime     time.Time // create 时间
>  dTime     time.Time // destroy 时间
>  firstUp   time.Time
>  lastUp    time.Time
>  firstDown time.Time
>  lastDown  time.Time
>  upCost    []time.Time // 上行接口耗时
>  downCost  []time.Time // 下行接口耗时
> }
> 
> type perfDist struct {
>  level   int  // 性能统计等级, 最高：FIFOPERF | LILOPERF | SESSPERF | INTFPERF
>  details map[string] /*sid*/ PerfDetail // 分布数据需要保存全量会话数据
> 
> }
> 
> func (pc *perfDist) Start(perfLevel int)
> // TODO check type and point, 根据性能等级判定当前point是否需要获取时间
> // TODO write to channel
> func (pc *perfDist) TickPoint(point int)
> func (pc *perfDist) Stop()
> // TODO read from channel
> // lock map
> func (pc *perfDist) analysis()
> // 性能指标落盘
> func (pc *perfDist) perfDump()
> ```

#### Ⅲ performance.go

> ```go
> type direction int
> type DataStatus int
> type SessStatus int
> 
> const (
>     UP   direction = 1
>     DOWN direction = 2
> 
>     DataBegin    DataStatus = 0
>     DataContinue DataStatus = 1
>     DataEnd      DataStatus = 2
>     DataTotal    DataStatus = 3
> 
>     SessBegin    SessStatus = 0
>     SessContinue SessStatus = 1
>     SessEnd      SessStatus = 2
>     SessOnce     SessStatus = 3
> 
>     outputPerfFile   = "perf.txt"
>     outputRecordFile = "perfReqRecord.csv"
>     outputPerfImg    = "perf.jpg"
> )
> 
> /*
> xtest 性能检测工具
> */
> type callDetail struct {
>     ID       string     //uuid
>     Handle   string     //会话模式时的hdl
>     Tm       time.Time  //时间戳
>     dataStat DataStatus //数据状态 ，0,1,2,3
>     sessStat SessStatus //会话状态,0,1,2,3
>     Dire     direction  //输入 还是输出
>     ErrCode  int
>     ErrInfo  string
> }
> 
> type performance struct {
>     Max         float32 `json:"max"`
>     Min         float32 `json:"min"`
>     FailRate    float32 `json:"failRate"`
>     SuccessRate float32 `json:"successRate"`
>     //平均值95 99线
>     Delay95      float32 `json:"delay95"`
>     Delay99      float32 `json:"delay99"`
>     DelayAverage float32 `json:"delayAverage"`
>     //首结果95 99线
>     DelayFirstMin     float32 `json:"delayFirstMin"`
>     DelayFirstMax     float32 `json:"delayFirstMax"`
>     DelayFirst95      float32 `json:"delayFirst95"`
>     DelayFirst99      float32 `json:"delayFirst99"`
>     DelayFirstAverage float32 `json:"delayFirstAverage"`
>     //尾结果95 99线
>     DelayLastMin     float32 `json:"delayLatMin"`
>     DelayLastMax     float32 `json:"delayLatMax"`
>     DelayLast95      float32 `json:"delayLast95"`
>     DelayLast99      float32 `json:"delayLast99"`
>     DelayLastAverage float32 `json:"delayLastAverage"`
> }
> 
> type singlePerfCost struct {
>     id        string
>     cost      float32
>     firstCost float32 //首个结果耗时
>     lastCost  float32 //最后一个结果耗时
> }
> 
> type errMsg struct {
>     ErrInfo string `json:"errInfo"`
>     Handle  string `json:"handle"`
> }
> 
> type PerfModule struct {
>     idx            int
>     collectChan    chan callDetail
>     mtx            sync.Mutex
>     control        chan bool
>     correctReqPath map[string][]callDetail //正确的请求路径图
> 
>     errReqRecord map[int][]errMsg //错误的请求记录
> 
>     correctReqCost []singlePerfCost //正确的请求花费的时间记录
> 
>     perf performance //性能结果
> 
>     reqRecordFile *os.File
> 
>     Log *utils.Logger
> }
> 
> var Perf *PerfModule
> 
> // 初始化Performance实例，并启动一个collect协程收集性能日志
> func (pf *PerfModule) Start() (err error) 
> // 关闭请求记录文件，calcDelay计算请求的性能指标，dump将性能指标数据落盘并关闭collectChan收集管道
> func (pf *PerfModule) Stop() 
> // 将采集详细数据写入collectChan收集管道
> func (pf *PerfModule) Record(id, handle string, stat DataStatus, stat2 SessStatus, dire direction, errCode int, errInfo string) 
> // 读取collectChan收集管道，将信息分类为正确与错误信息并记录correctReqPath和errReqRecord
> func (pf *PerfModule) collect() 
> // 计算correctReqPath[id]请求响应的时间开销
> func (pf *PerfModule) pretreatment(id string) 
> // 从性能日志文件中解析数据到实例
> func (pf *PerfModule) loadRecord() error 
> // 计算最后才能知道的性能指标，例如正确率、失败率、95、99指标
> func (pf *PerfModule) calcDelay()
> // 写入性能指标到outputPerfFile日志文件
> func (pf *PerfModule) dump()
> // 从data数据中计算出性能指标
> func (pf *PerfModule) anallyArray(data []float32) (min, max, average, aver95, aver99 float32) 
> ```

### 2.3、 inclue文件夹

#### Ⅰ.h264_nalu_spilt.h

#### Ⅱ. type.h

### 2.4、 lib文件夹

#### Ⅰ. libh264bitstream.so.0

### 2.5 log 文件夹
> 记录Xtest运行日志，常见的有如下：
> - .png 运行资源变化折线图
> - errDist: 错误分析日志
> - perf.txt perf.cvs: 性能分析日志
> - result: 运行结果
> - xtest*.log: Xtest运行日志

### 2.6 prometheus文件夹
#### Ⅰ. resource.go: 资源记录相关函数定义
> ```go
> // Resource 资源条目
> type Resource struct {
>	Mem  float64 // 内存
>	Cpu  float64 // cpu
>	Time float64 // 时间
>}
>
>type Resources struct {
>	resourceChan chan Resource // 发送资源条目的管道
>	resources    []Resource // 记录条目的数组
>	stopChan     chan bool // 通知关闭记录的管道
>	wg           sync.WaitGroup
>}
>
> // NewResources 实例化资源结构体
>func NewResources() Resources
>
>// Serve 启动Prometheus监听
>func (rs *Resources) Serve() 
>// ReadMem 获取内存使用, 传入AiService的地址
>func (rs *Resources) ReadMem(taddrs string) error 
> // 从管道中读取数据存储到数组中
>func (rs *Resources) GenerateData()
>
>// MetricValue 获取metric的Value值
>func (rs *Resources) MetricValue(m prometheus.Gauge) >(float64, error) 
>// Draw 绘制图片
>func (rs *Resources) Draw(dst string) error 
> 关闭管道，通知停止采集数据
>func (rs *Resources) Stop()
>
>// bToMb bit转Mb
>func bToMb(b uint64) uint64 

### 2.7、request文件夹

#### Ⅰ. fileSession.go

> ```go
> // 文件session请求
> func FileSessionCall(cli *xsfcli.Client, index int64) (info analy.ErrInfo)
> // 文件AI上行请求
> func FilesessAIIn(cli *xsfcli.Client, indexs int64, thrRslt *[]protocol.LoaderOutput, thrLock *sync.Mutex, reqSid string) (hdl string, status protocol.LoaderOutput_RespStatus, info analy.ErrInfo) 
> // 多线程文件上传流请求
> func FilemultiUpStream(cli *xsfcli.Client, swg *sync.WaitGroup, session string, pm *[]protocol.LoaderOutput, sm *sync.Mutex, errchan chan struct {
> code int
> err  error
> }) 
> 
> // 实时性校准,用于校准发包大小及发包时间间隔之间的实时性.
> func FilertCalibration(curReq int, interval int, sTime time.Time)
> 
> // downStream 下行调用单线程;
> func FilesessAIOut(cli *xsfcli.Client, hdl string, sid string, rslt *[]protocol.LoaderOutput) (info analy.ErrInfo) 
> // 文件session报错
> func FilesessAIExcp(cli *xsfcli.Client, hdl string, sid string) (err error)
> 
> // upStream first error ，将上传流错误写入ch管道
> func FileunBlockChanWrite(ch chan analy.ErrInfo, err analy.ErrInfo) 
> ```

#### Ⅱ. oneShot.go：与RPC通信有关的函数定义

> ```go
> // 使用xsf框架发起RPC通信，设置协议参数、上行数据键值对，使用ONESHORT方式发起SessionCall，然后
> // 下行数据解析到AsyncDrop下行数据异步落盘同步通道，如果通道满了，使用downOutput函数写入本地文件。
> func OneShotCall(cli *xsfcli.Client, index int64) (info analy.ErrInfo)
> ```

#### Ⅲ. output.go：下行数据输出有关的函数定义

> ```go
> // 读取AsyncDrop通道中的下行数据，调用downOutput函数写入本地文件
> func DownStreamWrite(wg *sync.WaitGroup, log *utils.Logger) 
> // 写入数据到本地文件（_var.OutputDst和_var.Output配置相关文件路径）
> func downOutput(key string, data []byte, log *utils.Logger)
> ```
#### Ⅳ. request.go 
>```go
> // 请求结构体，只有一个配置成员
> type Request struct {
>	C _var.Conf // 请求所需的配置
>}


#### Ⅴ. session.go

> ```go
> // Session调用
> func SessionCall(cli *xsfcli.Client, index int64) (info analy.ErrInfo) 
> 
> // AI调用输入
> func sessAIIn(cli *xsfcli.Client, indexs []int, thrRslt *[]protocol.EngOutputData, thrLock *sync.Mutex, reqSid string) (hdl string, status protocol.EngOutputData_DataStatus, info analy.ErrInfo)
> 
> // 多线程上传流
> func multiUpStream(cli *xsfcli.Client, swg *sync.WaitGroup, session string, interval int, indexs map[int]int, sid string, pm *[]protocol.EngOutputData, sm *sync.Mutex, errchan chan analy.ErrInfo)
> // 实时性校准,用于校准发包大小及发包时间间隔之间的实时性.
> func rtCalibration(curReq int, interval int, sTime time.Time)
> // downStream 下行调用单线程;
> func sessAIOut(cli *xsfcli.Client, hdl string, sid string, rslt *[]protocol.EngOutputData) (info analy.ErrInfo) 
> // 
> func sessAIExcp(cli *xsfcli.Client, hdl string, sid string) (err error)
> // upStream first error 将错误信息发送至ch管道
> func unBlockChanWrite(ch chan analy.ErrInfo, err analy.ErrInfo)
> ```

#### Ⅵ. signal.go：xtest退出有关的函数定义

> ```go
> // 通过signal.Notify转发信号，优雅退出程序
> func SigRegister() 
> ```

#### Ⅶ. splitFrame.go

> ```go
> func GetH264Frames(video []byte) (frameSizes []int)
> ```

#### Ⅷ. textLine.go

> ```go
> // 文本session请求
> func TextCall(cli *xsfcli.Client, index int64) (info analy.ErrInfo)
> // 文本AI上行
> func TextAIIn(cli *xsfcli.Client, indexs int64, thrRslt *[]protocol.LoaderOutput, thrLock *sync.Mutex, reqSid string) (hdl string, status protocol.LoaderOutput_RespStatus, info analy.ErrInfo)
> 
> // 多线程文本上行数据流
> func TextmultiUpStream(cli *xsfcli.Client, swg *sync.WaitGroup, session string, pm *[]protocol.LoaderOutput, sm *sync.Mutex, errchan chan analy.ErrInfo) 
> 
> // 实时性校准,用于校准发包大小及发包时间间隔之间的实时性.
> func TextrtCalibration(curReq int, interval int, sTime time.Time)
> 
> // downStream 下行调用单线程;
> func TextsessAIOut(cli *xsfcli.Client, hdl string, sid string, rslt *[]protocol.LoaderOutput) (info analy.ErrInfo) 
> 
> func TextsessAIExcp(cli *xsfcli.Client, hdl string, sid string) (err error)
> 
> // upStream first error，错误日志写入管道
> func TextunBlockChanWrite(ch chan analy.ErrInfo, err analy.ErrInfo) 
> ```

### 2.8、script文件夹

#### Ⅰ. test.sh： 运行脚本

> 启动xtest脚本

#### Ⅱ. xtest.toml：配置文件

> ```toml
> [xtest]
> #测试目标服务配置，配置格式如下,注意分割符的差异. 业务1@ip1:port1;ipn:portn,业务2@ip2:port2;ipn:portn
> taddrs="AIservice@127.0.0.1:5090"
> trace-ip = "172.16.51.13"
> 
> [svcMode]
> service = "AIservice"           # 请求目标服务名称, eg:sms
> timeout = 1000                  # 服务超时时间, 对应服务端waitTime
> multiThr = 10                  # 请求并发数
> loopCnt = 100                 # 请求总次数
> reqMode = 0                     # 服务请求模式. 0:非会话模式 1:会话模式
> reqPara = "k1=v1,k2=v2,k3=v3"   # 服务请求参数对, 多个参数对以","分隔
> linearMs = 5000                 # 并发压测线性增长区间,单位:ms
> 
> 
> [upstream]
> inputSrc = "./test"               # 上行数据流数据源, 配置文件路径(配置为目录则循环读取目录中文件)
> sliceOn = 1                     # 切片开关, 0:关闭切换, !0:开启切片
> sliceSize = 1280                # 上行数据切片大小,用于会话模式: byte
> interval = 40                   # 上行数据发包间隔,用于会话模式: ms. 注：动态校准,每个包间隔非严格interval
> type = "audio"                  # 数据类型
> format = "audio/L16;rate=16000" # 数据格式
> encoding = "raw"                # 数据编码
> describe = "k=v;k=v"            # 数据描述信息
> 
> 
> #[upstream-N]                    # 用于实现多数据流上行配置, 对于多个数据流可按照upstream-N规则叠加配置
> #inputSrc = "path-N"             # 同upstream
> #sliceSize = 1280                # 同upstream
> #interval = 40                   # 同upstream, 注:多个upstream发包间隔相同时,数据流发送包合并
> #type = "audio"                  # 同upstream
> #format = "audio/L16;rate=16000" # 同upstream
> #encoding = "raw"                # 同upstream
> #describe = "k=v;k=v"            # 同upstream
> 
> 
> [downstream]                    # 下行数据流存储输出
> output = 0                      # 输出方式： 0:输出至公共文件outputDst 1:以独立文件形式输出至目录outputDst(文件名:sid+**) 2：输出至终端
> outputDst = "./log/result"          # 响应数据输出目标, 若output=0则配置目录, output=1则配置文件
> 
> 
> [log]
> file = "./log/xtest.log"              # 日志文件名
> level = "debug"                 # 日志打印等级
> size = 100                      # 日志文件大小
> count = 20                      # 日志备份数量
> async = 0                       # 异步开关
> die = 30
> 
> [trace]
> able = 0
> ```

#### Ⅲ. xtest_example.toml文件

### 2.9、testdata文件夹

### 2.10、util文件夹

#### Ⅰ charts.go： 绘制图表有关的函数定义

```go

const (
   lineChartXAxisName = "Time"
   lineChartYAxisName = "Percentage"
   lineChartHeight    = 700
   lineChartWidth     = 1280
   colorMultiplier    = 256
)

var (
   lineChartStyle = chart.Style{
      Padding: chart.Box{
         Top:  30,
         Left: 150,
      },
   }
   timeFormat = GetHMS
)

type Charts struct {
   Vals    LinesData
   Dst     string    // 保存文件
   XValues []float64 // X轴时间戳
}

type LineYValue struct {
   Name   string
   Values []float64
}

type LinesData struct {
   Title     string
   BarValues []LineYValue
}

// createLineChart 创建线性图
func (c *Charts) createLineChart(title string, xValues []float64, values []LineYValue) error

// Draw 传入绘制数据，绘制条形图
func (c *Charts) Draw() error 
// GetHMS 格式化时间获取时分秒
func GetHMS(v interface{}) string

// getNsec 获取纳秒数
func getNsec(cur time.Time) float64 
```

#### Ⅱ file.go: 与文件读取有关的函数定义

```go
func ReadDir(fi os.FileInfo, src string, sep string, flag int) ([][]byte, error) 
func CompFunc(flag int, i, j string) bool
```

#### Ⅲ ptermShow.go： 与进度条有关的函数定义
>```go
> // 输出一个任务进度条
> func ProgressShow(cnt *atomic.Int64, cnt1 int64)

#### Ⅳ sid.go：与sid生成有关的函数定义

> ```go
> var (
>  index        int64  = 0		// 生成的SID索引
>  Location     string = "dx"
>  LocalIP      string 			// 本地IP
>  ShortLocalIP string			// 本地短IP
>  Port         string			// 端口
> )
> // 获取本地ip地址与短地址ip
> func init() 
> // 生成sid
> func NewSid(sub string) string 
> ```

#### Ⅴ task.go： 与定时任务有关的函数定义
>```go
> // ScheduledTaskPool 定时任务池
>type ScheduledTaskPool struct {
>	Size         int            // 任务个数
>	TimeStopChan chan bool      // // 通知定时任务结束协程
>	wg           sync.WaitGroup // 同步原语
>}
> // NewScheduledTaskPool 实例化一个任务池对象
>func NewScheduledTaskPool() ScheduledTaskPool 
>
>// Start 启动一个定时任务 jbzhou5
>func (stp *ScheduledTaskPool) Start(d time.Duration, f func())
>
>// Stop 结束定时任务
>func (stp *ScheduledTaskPool) Stop() 

### 2.11、var文件夹

#### Ⅰ. cmd.go：命令行输入相关

> ```go
>type Flag struct {
>	/*	CmdMode		= xsf.Mode				// -m
>		CmdCfg		= xsf.Cfg				// -c
>		CmdProject	= xsf.Project			// -p
>		CmdGroup	= xsf.Group				// -g
>		CmdService	= xsf.Service			// -s
>		CmdCompanionUrl = xsf.CompanionUrl	// -u
>	*/
>	// default 缺省配置模式为native
>	CmdCfg       *string // 指定配置文件
>	XTestVersion *bool // xtest版本
>}
>
> // 使用Flag包将Flag结构体中的变量和命令行参数绑定
>func NewFlag() Flag 
>
> // 解析命令行参数填充到flag结构体
>func (f *Flag) Parse() 
> // 打印xtest用法配置选项
> func Usage() 
> ```

#### Ⅱ. conf.go：xtest配置相关定义和函数

> ```go
> const (
> 	CliName = "xtest"
> 	CliVer  = "2.0.1"
> )
> 
> type InputMeta struct {
> 	Name       string                     // 上行数据流key值
> 	DataSrc    string                     // 上行实体数据来源;数据集则配置对应目录
> 	SliceOn    int                        // 上行数据切片开关, !0:切片. 0:不切片
> 	UpSlice    int                        // 上行数据切片大小: byte
> 	UpInterval int                        // slice发包间隔: ms
> 	DataType   protocol.MetaDesc_DataType // audio/text/image/video
> 	DataDesc   map[string]string
> 
> 	// DataList map[string/*file*/] []byte /*data*/
> 	DataList [][]byte /*data*/
> }
> 
> type OutputMeta struct {
> 	Name string            // 下行数据流key
> 	Sid  string            // sid
> 	Type string            // 下行数据类型
> 	Desc map[string]string // 数据描述
> 	Data []byte            // 下行数据实体
> }
> // 配置
>type Conf struct {
>  // [xtest]
>  Taddrs            string
>	// [svcMode]
>	SvcId            string
>	SvcName          string        // dst service name
>	TimeOut          int           // 超时时间: ms, 对应加载器>waitTime
>	LossDeviation    int           // 自身性能损耗误差, ms.
>	MultiThr         int           // 请求并发数
>	DropThr          int           // 下行数据异步输出线程数
>	LoopCnt          *atomic.Int64 // 请求总次数
>	ReqMode          int           // 0: 非会话模式, 1: 常规会话>模式 2.文本按行会话模式 3.文件会话模式
>	LinearNs         int           // 并发模型线性增长时间,用于计>算并发增长斜率(单位：ns). default:0,瞬时并发压测.
>	TestSub          string        // 测试业务sub, 缺省test
>	InputCmd         bool          // jbzhou5 非会话模式切换为命>令行输入
>	PrometheusSwitch bool          // jbzhou5 Prometheus写入开>关
> PrometheusPort   int           // jbzhou5 Prometheus指标服务端口
> Plot             bool          // jbzhou5 绘制图形开关
>	PlotFile         string        // jbzhou5 绘制图像保存路径
>	FileSorted       int           // jbzhou5 文件排序方式
>	FileNameSeq      string        // 文件名分割方式
>	PerfConfigOn     bool          //true: 开启性能检测 false: >不开启性能检测
>	PerfLevel        int           //非会话模式默认0
>	//会话模式0: 从发第一帧到最后一帧的性能
>	//会话模式1:首结果(发送第一帧到最后一帧的性能)
>	//会话模式2:尾结果(发送最后一帧到收到最后一帧的性能)
>	// 请求参数对
>	Header map[string]string
>	Params map[string]string
>
>	Payload []string // 上行数据流
>	Expect  []string // 下行数据流
>
>	// 上行数据流配置, 多数据流通过section [data]中payload进行配置
>	UpStreams []InputMeta
>
>	DownExpect []protocol.MetaDesc
>
>	// [downstream]
>	Output int // 0：输出至公共文件outputDst(sid+***:data)
>	// 1：以独立文件形式输出至目录outputDst(文件名:sid+***)
>	// 2：输出至终端
>	//-1：不输出
>	OutputDst string // output=0时,该项配置输出文件名; output=1>时,该项配置输出目录名
>	ErrAnaDst string
>	AsyncDrop chan OutputMeta // 下行数据异步落盘同步通道
>
>	// jbzhou5 性能资源日志保存目录
>	// ResourcesDst = "./"
>	// jbzhou5 Prometheus并发协程计数器
>	ServicePid     int // jbzhou5 Aiservice的PID号
>	ConcurrencyCnt prometheus.Gauge
>	// jbzhou5 Prometheus监听参数
>	CpuPer prometheus.Gauge
>	MemPer prometheus.Gauge
>}
>
>func NewConf() Conf
> // 调用解析函数，将配置文件数据填充到conf结构体
>func (c *Conf) ConfInit(conf *utils.Configure) error
>
> // 解析Expect
>func (c *Conf) secParseEp(conf *utils.Configure) error 
>
>//解析payload
>func (c *Conf) secParsePl(conf *utils.Configure) error
> // 解析service
>func (c *Conf) secParseSvc(conf *utils.Configure) error
> // 解析Header
>func (c *Conf) secParseHeader(conf *utils.Configure) error
> // 解析params
>func (c *Conf) secParseParams(conf *utils.Configure) error
> // 解析数据
>func (c *Conf) secParseData(conf *utils.Configure) error 
> // 解析下行流
>func (c *Conf) secParseDStream(conf *utils.Configure) error
>
>//jbzhou5 解析命令行输入的数据
>func (c *Conf) secParseCmd(conf *utils.Configure) error
>
> ```

## 三、需求分析

### 3.1 功能需求

- 模式-支持流式、非流式、异步回调三种模式
- 非流式模式下
    - [x] 读取文件输入
    - [x] 配置中心手动输入数据 √ 使用scan输入文本数据   success
- 流式模式下
    - [x] 单一文件一次输入
    - [x] 单一文件按照固定长度输入
    - [ ] 文本文件按行读取 √  优化代码
    - [ ] 多个文件，每个文件一帧输入
        - [x] 文件有序 √ 优化代码
        - [x] 文件无序：Shuffle √ 实现功能

### 3.2 性能需求

- [x] 并发，显示当前路数 √ 实现功能 success
- [ ] 成功率，性能数据及性能分布，输出本地相关数据 √ 优化代码
- [x] 内存、显存定时统计： 内存：syscall 包调用。显存：cmd执行 nvidia-smi ？node_exporter？

### 3.3 其他需求

- 文档说明
- demo样例

## 四、功能规划

开启性能测试：日志大小预估

100次： 207KB

1000次：2115KB

10000次：21639KB

100000次：219MB

todolist:

## 五、配置说明

### 5.1 样例配置

```toml
[xtest]
taddrs="AIservice@127.0.0.1:5090"
trace-ip = "172.16.51.13"

[svcMode]
service = "AIservice"           # 请求目标服务名称, eg:sms
svcId = "s12345678"             # 服务id
timeout = 1000                  # 服务超时时间, 对应服务端waitTime
multiThr = 100                  # 请求并发数
loopCnt = 100000                # 请求总次数
sessMode = 0                    # 0: 非会话模式, 1: 常规会话模式 2.文本按行会话模式 3.文件会话模式
linearMs = 5000                 # 并发压测线性增长区间,单位:ms
perfOn=true                     # 是否开启性能测试
perfLevel=0                     # 非会话模式默认0
                                # 会话模式0: 从发第一帧到最后一帧的性能
                                # 会话模式1:首结果(发送第一帧到最后一帧的性能)
                                # 会话模式2:尾结果(发送最后一帧到收到最后一帧的性能)
inputCmd = false 				# 切换为命令行输入，仅在非会话模式生效
prometheus_switch = true  		# Prometheus开关， 开启后开启双写，同时写入prometheus与本地日志
prometheus_port = 2117    # jbzhou5 Prometheus指标暴露端口
plot = true  					# 绘制资源图， 默认开启
plot_file = "./log/line.png"    # 绘制图形保存路径
file_sorted = 0  				# 传入文件是否排序， 0： 随机， 1： 升序， 2： 降序
file_name_seq = "_" 			# 传入文件名分割方式 例如传入'_', 则1_2.txt -> 1，2_2.txt -> 2, 为空或者传入非法则不处理
[header]
"appid" = "100IME"
"uid" = "1234567890"

[parameter]
"key" = 2
"x" = 1

[data]
payload = "dataKey2"   			# 输入数据流配置段,多个数据流以";"分割， 如果开启了inputCmd， 该值会被清空
expect = "dataKey3"             # 输出数据流配置段,多个数据流以";"分割

[dataKey1]  # 输入数据流dataKey1描述
inputSrc = "path"               # 上行数据流数据源, 配置文件路径(配置为目录则循环读取目录中文件)
sliceOn = false                 # 切片开关, false:关闭切换, true:开启切片
sliceSize = 1280                # 上行数据切片大小,用于会话模式: byte
interval = 40                   # 上行数据发包间隔,用于会话模式: ms. 注：动态校准,每个包间隔非严格interval
name = "input1"                 # 输入数据流key值
type = "image"                  # 数据类型，支持"audio","text","image","video"
describe = "k1=v1;k2=v2"        # 数据描述信息,多个描述信息以";"分割
                                # 图像支持如下属性："encoding", 如"encoding=jpg"
[dataKey2]
inputSrc = "./testdata/text2"   # 上行数据流数据源, 配置文件路径(配置为目录则循环读取目录中文件)
sliceOn = false                 # 切片开关, false:关闭切换, true:开启切片
sliceSize = 1280                # 上行数据切片大小,用于会话模式: byte
interval = 40                   # 上行数据发包间隔,用于会话模式: ms. 注：动态校准,每个包间隔非严格interval
name = "input2"                 # 输入数据流key值
type = "text"                  	# 数据类型，支持"audio","text","image","video"
describe = "encoding=utf8"      # 数据描述信息,多个描述信息以";"分割
                                # 图像支持如下属性："encoding", 如"encoding=jpg"
[dataKey3]
name = "result"                 # 输入数据流key值
type = "text"                   # 输出数据类型，支持"audio","text","image","video"
describe = "k1=v1;k2=v2"        # 数据描述信息,多个描述信息以";"分割
                                # 文本支持如下属性："encoding","compress", 如"encoding=utf8;compress=gzip"


[downstream]                    # 下行数据流存储输出
output = 0                      # 输出方式： 0:输出至公共文件outputDst 1:以独立文件形式输出至目录outputDst(文件名:sid+**) 2：输出至终端
outputDst = "./log/result"      # 响应数据输出目标, 若output=0则配置目录, output=1则配置文件


[log]
file = "./log/xtest.log"        # 日志文件名
level = "debug"                 # 日志打印等级
size = 100                      # 日志文件大小
count = 20                      # 日志备份数量
async = 0                       # 异步开关
die = 30

[trace]
able = 0
```
## 六、字段说明
> **xtest.toml 中大部分字段一般保持不变，下面仅对常用字段进行说明解释。**
- ```[xtest]```
  - ```taddrs="AIservice@127.0.0.1:5090"```： 与Aiservice的通信地址，与AIservice的启动端口对应，其中端口会被解析用于获取Aiservice的进程，监听其使用资源信息。
  - ``` trace-ip = "172.16.51.13```
- ```[svcMode]```
  - ```service = "AIservice"``` ： 请求目标服务名称, eg:sms
  - ```svcId = "s12345678"```   ：服务id
  - ```timeout = 1000``` ：服务超时时间, 对应服务端waitTime
  - ```multiThr = 100``` ：请求并发数，即同时开启多个协程发送请求测试
  - ```loopCnt = 100000``` ： multiThr个协程发送的请求总次数
  - ```sessMode = 0``` ： 0: 非会话模式, 1: 常规会话模式 2.文本按行会话模式 3.文件会话模式
  - ```linearMs = 5000``` ：并发压测线性增长区间,单位:ms
  - ```perfOn=true``` ： 是否开启性能测试，即是否在log文件夹底下记录perf.txt和PerfRecord.csv，用于记录成功率、失败率、发送延迟等性能指标。
  - ```perfLevel=0 ```：与sessMode字段对应，非会话模式默认0，会话模式0: 从发第一帧到最后一帧的性能，会话模式1:首结果(发送第一帧到最后一帧的性能)，会话模式2:尾结果(发送最后一帧到收到最后一帧的性能)
  - ```inputCmd = false ```：切换为命令行输入，仅在非会话模式生效，配置该字段时，所配置的[data] 字段将失效，仅读取用户命令行输入数据
  - ```prometheus_switch = true``` ：Prometheus开关，开启后会开放一个Prometheus监控端口，可使用grafana进行数据的展示。关闭/打开都会在Log目录生成一个Resource.csv 资源监听文件。
  - ```prometheus_port = 2117```：Prometheus指标暴露端口
  - ```plot = true``` ：绘制资源图，默认开启，将绘制Aiservice使用的资源变化折线图
  - ```plot_file = "./log/line.png"``` ：绘制图形保存路径
  - ```file_sorted = 0``` ：[data] 传入文件是否按名称排序， 0： 随机， 1： 升序， 2： 降序
  - ```file_name_seq = "_"``` ： 传入文件名分割方式 例如传入'_', 则1_2.txt -> 1，2_2.txt -> 2, 为空或者传入非法（即不能作为文件名的字符）则不处理， 注意此处分割为仅保留前半部分，若文件名为1_2_3.txt， 则得到的分割文件名为1。
- 
- ```[parameter]```：使用的AI模型需要传入的字段，根据自己需要填写
  - ```"key" = 2```
  - ```"x" = 1```

- ```[data]```
  - ```payload = "dataKey2"``` ：输入数据流配置段,多个数据流以";"分割， jbzhou5 如果开启了inputCmd， 该值会被清空
  - ```expect = "dataKey3"```：输出数据流配置段,多个数据流以";"分割

- ```[dataKey1] ``` ：名称可自定义，主要用于在[data] 字段的payload属性中方便标记加载数据，输入数据流dataKey1描述
  - ```inputSrc = "path" ```：上行数据流数据源, 配置文件路径(配置为目录则循环读取目录中文件)
  - ```sliceOn = false```：切片开关, false:关闭切换, true:开启切片
  - ```sliceSize = 1280```：上行数据切片大小,用于会话模式: byte
  - ```interval = 40```： 上行数据发包间隔,用于会话模式: ms. 注：动态校准,每个包间隔非严格interval
  - ```name = "input1"```： 输入数据流key值
  - ```type = "image"```： 数据类型，支持"audio","text","image","video"
  - ```describe = "k1=v1;k2=v2"``： 数据描述信息,多个描述信息以";"分割，图像支持如下属性："encoding", 如"encoding=jpg"

- ```[dataKey2]```
  - ```inputSrc = "./testdata/text2"```：上行数据流数据源, 配置文件路径(配置为目录则循环读取目录中文件)
  - ```sliceOn = false```： 切片开关, false:关闭切换, true:开启切片
  - ```sliceSize = 1280```：上行数据切片大小,用于会话模式: byte
  - ```interval = 40```： 上行数据发包间隔,用于会话模式: ms. 注：动态校准,每个包间隔非严格interval
  - ```name = "input2"```： 输入数据流key值
  - ```type = "text"``` ： 数据类型，支持"audio","text","image","video"
  - ```describe = "encoding=utf8"``` ： 数据描述信息,多个描述信息以";"分割图像支持如下属性："encoding", 如"encoding=jpg"

- ``` [dataKey3]```
  - ```name = "result"```： 输入数据流key值
  - ```type = "text"```：输出数据类型，支持"audio","text","image","video"
  - ```describe = "k1=v1;k2=v2"```： 数据描述信息,多个描述信息以";"分割，文本支持如下属性："encoding","compress", 如"encoding=utf8;compress=gzip"


- ```[downstream] ``` ：下行数据流存储输出
  - ```output = 0 ```： 输出方式： 0:输出至公共文件outputDst 1:以独立文件形式输出至目录outputDst(文件名:sid+**) 2：输出至终端
  - ```outputDst = "./log/result"```：响应数据输出目标, 若output=0则配置目录, output=1则配置文件


- ```[log]```
  - ```file = "./log/xtest.log"```：日志文件名
  - ```level = "debug" ```：日志打印等级
  - ```size = 100 ```：日志文件大小，单个超过size，将会写入新文件。
  - ```count = 20``` ： 日志备份数量
  - ```async = 0```： 异步开关
  - ```die = 30```

- ```[trace]```
  - ```able = 0```
## 七、使用说明
1. 启动Aiservice，注意监听的端口是否有改变。
2. 根据自己的AI模型修改xtest.toml文件，例如参数，测试轮数等配置，具体请参考说明六。
3. 由于监听资源需要管理员权限，请保证本地已经安装**netstat**命令，然后执行： ```sudo ./xtest``` 或```sudo ./xtest -f *.toml ```命令启动， 否则资源文件将为空，但并不影响其他任务。