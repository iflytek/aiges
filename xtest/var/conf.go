package _var

import (
	"errors"
	"fmt"
	"git.iflytek.com/AIaaS/xsf/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/atomic"
	"io/ioutil"
	"os"
	"protocol"
	"reflect"
	"strconv"
	"strings"
	"xtest/util"
)

const (
	CliName = "xtest"
	CliVer  = "2.0.1"
)

type InputMeta struct {
	Name       string                     // 上行数据流key值
	DataSrc    string                     // 上行实体数据来源;数据集则配置对应目录
	SliceOn    int                        // 上行数据切片开关, !0:切片. 0:不切片
	UpSlice    int                        // 上行数据切片大小: byte
	UpInterval int                        // slice发包间隔: ms
	DataType   protocol.MetaDesc_DataType // audio/text/image/video
	DataDesc   map[string]string

	// DataList map[string/*file*/] []byte /*data*/
	DataList [][]byte /*data*/
}

type OutputMeta struct {
	Name string            // 下行数据流key
	Sid  string            // sid
	Type string            // 下行数据类型
	Desc map[string]string // 数据描述
	Data []byte            // 下行数据实体
}

//var (
//	// [svcMode]
//	SvcId            string        = "s12345678"
//	SvcName          string        = "AIservice"            // dst service name
//	TimeOut          int           = 1000                   // 超时时间: ms, 对应加载器waitTime
//	LossDeviation    int           = 50                     // 自身性能损耗误差, ms.
//	MultiThr         int           = 100                    // 请求并发数
//	DropThr                        = 100                    // 下行数据异步输出线程数
//	LoopCnt          *atomic.Int64 = atomic.NewInt64(10000) // 请求总次数
//	ReqMode          int           = 0                      // 0: 非会话模式, 1: 常规会话模式 2.文本按行会话模式 3.文件会话模式
//	LinearNs         int           = 0                      // 并发模型线性增长时间,用于计算并发增长斜率(单位：ns). default:0,瞬时并发压测.
//	TestSub          string        = "ase"                  // 测试业务sub, 缺省test
//	InputCmd         bool          = false                  // jbzhou5 非会话模式切换为命令行输入
//	PrometheusSwitch bool          = false                  // jbzhou5 Prometheus写入开关
//	Plot             bool          = true                   // jbzhou5 绘制图形开关
//	PlotFile         string        = "./log/line.png"       // jbzhou5 绘制图像保存路径
//	FileSorted       int           = 0                      // jbzhou5 文件排序方式
//	FileNameSeq      string        = "/"                    // 文件名分割方式
//	PerfConfigOn     bool          = false                  //true: 开启性能检测 false: 不开启性能检测
//	PerfLevel        int           = 0                      //非会话模式默认0
//	//会话模式0: 从发第一帧到最后一帧的性能
//	//会话模式1:首结果(发送第一帧到最后一帧的性能)
//	//会话模式2:尾结果(发送最后一帧到收到最后一帧的性能)
//	// 请求参数对
//	Header map[string]string = make(map[string]string)
//	Params map[string]string = make(map[string]string)
//
//	Payload []string // 上行数据流
//	Expect  []string // 下行数据流
//
//	// 上行数据流配置, 多数据流通过section [data]中payload进行配置
//	UpStreams []InputMeta = make([]InputMeta, 0, 1)
//
//	DownExpect []protocol.MetaDesc
//
//	// [downstream]
//	Output = 0 // 0：输出至公共文件outputDst(sid+***:data)
//	// 1：以独立文件形式输出至目录outputDst(文件名:sid+***)
//	// 2：输出至终端
//	//-1：不输出
//	OutputDst = "./log/result" // output=0时,该项配置输出文件名; output=1时,该项配置输出目录名
//	ErrAnaDst = "./log/errDist"
//	AsyncDrop chan OutputMeta // 下行数据异步落盘同步通道
//
//	// jbzhou5 性能资源日志保存目录
//	// ResourcesDst = "./"
//	// jbzhou5 Prometheus并发协程计数器
//	ServicePid     int = 0 // jbzhou5 Aiservice的PID号
//	ConcurrencyCnt     = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "Xtest_Concurrency_Go_Routine",
//		Help: "The total number of processed events",
//	})
//	// jbzhou5 Prometheus监听参数
//	CpuPer = promauto.NewGauge(prometheus.GaugeOpts{ // CPU 利用率
//		Name: "Xtest_CPU_Percent",
//		Help: "Xtest cpu percent",
//	})
//	MemPer = promauto.NewGauge(prometheus.GaugeOpts{ // 内存利用率
//		Name: "Xtest_MEM_Percent",
//		Help: "Xtest mem percent",
//	})
//)

type Conf struct {
	// [xtest]
	Taddrs			string
	// [svcMode]
	SvcId            string
	SvcName          string        // dst service name
	TimeOut          int           // 超时时间: ms, 对应加载器waitTime
	LossDeviation    int           // 自身性能损耗误差, ms.
	MultiThr         int           // 请求并发数
	DropThr          int           // 下行数据异步输出线程数
	LoopCnt          *atomic.Int64 // 请求总次数
	ReqMode          int           // 0: 非会话模式, 1: 常规会话模式 2.文本按行会话模式 3.文件会话模式
	LinearNs         int           // 并发模型线性增长时间,用于计算并发增长斜率(单位：ns). default:0,瞬时并发压测.
	TestSub          string        // 测试业务sub, 缺省test
	InputCmd         bool          // jbzhou5 非会话模式切换为命令行输入
	PrometheusSwitch bool          // jbzhou5 Prometheus写入开关
	PrometheusPort   int           // jbzhou5 Prometheus指标服务端口
	Plot             bool          // jbzhou5 绘制图形开关
	PlotFile         string        // jbzhou5 绘制图像保存路径
	FileSorted       int           // jbzhou5 文件排序方式
	FileNameSeq      string        // 文件名分割方式
	PerfConfigOn     bool          //true: 开启性能检测 false: 不开启性能检测
	PerfLevel        int           //非会话模式默认0
	//会话模式0: 从发第一帧到最后一帧的性能
	//会话模式1:首结果(发送第一帧到最后一帧的性能)
	//会话模式2:尾结果(发送最后一帧到收到最后一帧的性能)
	// 请求参数对
	Header map[string]string
	Params map[string]string

	Payload []string // 上行数据流
	Expect  []string // 下行数据流

	// 上行数据流配置, 多数据流通过section [data]中payload进行配置
	UpStreams []InputMeta

	DownExpect []protocol.MetaDesc

	// [downstream]
	Output int // 0：输出至公共文件outputDst(sid+***:data)
	// 1：以独立文件形式输出至目录outputDst(文件名:sid+***)
	// 2：输出至终端
	//-1：不输出
	OutputDst string // output=0时,该项配置输出文件名; output=1时,该项配置输出目录名
	ErrAnaDst string
	AsyncDrop chan OutputMeta // 下行数据异步落盘同步通道

	// jbzhou5 性能资源日志保存目录
	// ResourcesDst = "./"
	// jbzhou5 Prometheus并发协程计数器
	ConcurrencyCnt prometheus.Gauge
	// jbzhou5 Prometheus监听参数
	CpuPer prometheus.Gauge
	MemPer prometheus.Gauge
}

func NewConf() Conf {
	return Conf{
		Taddrs: "",
		SvcId:            "s12345678",
		SvcName:          "AIservice",            // dst service name
		TimeOut:          1000,                   // 超时时间: ms, 对应加载器waitTime
		LossDeviation:    50,                     // 自身性能损耗误差, ms.
		MultiThr:         100,                    // 请求并发数
		DropThr:          100,                    // 下行数据异步输出线程数
		LoopCnt:          atomic.NewInt64(10000), // 请求总次数
		ReqMode:          0,                      // 0: 非会话模式, 1: 常规会话模式 2.文本按行会话模式 3.文件会话模式
		LinearNs:         0,                      // 并发模型线性增长时间,用于计算并发增长斜率(单位：ns). default:0,瞬时并发压测.
		TestSub:          "ase",                  // 测试业务sub, 缺省test
		InputCmd:         false,                  // jbzhou5 非会话模式切换为命令行输入
		PrometheusSwitch: false,                  // jbzhou5 Prometheus写入开关
		Plot:             true,                   // jbzhou5 绘制图形开关
		PlotFile:         "./log/line.png",       // jbzhou5 绘制图像保存路径
		FileSorted:       0,                      // jbzhou5 文件排序方式
		FileNameSeq:      "/",                    // 文件名分割方式
		PerfConfigOn:     false,                  //true: 开启性能检测 false: 不开启性能检测
		PerfLevel:        0,                      //非会话模式默认0
		//会话模式0: 从发第一帧到最后一帧的性能
		//会话模式1:首结果(发送第一帧到最后一帧的性能)
		//会话模式2:尾结果(发送最后一帧到收到最后一帧的性能)
		// 请求参数对
		Header: make(map[string]string),
		Params: make(map[string]string),

		Payload: []string{}, // 上行数据流
		Expect:  []string{}, // 下行数据流

		// 上行数据流配置, 多数据流通过section [data]中payload进行配置
		UpStreams: make([]InputMeta, 0, 1),

		DownExpect: []protocol.MetaDesc{},

		// [downstream]
		Output: 0, // 0：输出至公共文件outputDst(sid+***:data)
		// 1：以独立文件形式输出至目录outputDst(文件名:sid+***)
		// 2：输出至终端
		//-1：不输出
		OutputDst: "./log/result", // output=0时,该项配置输出文件名; output=1时,该项配置输出目录名
		ErrAnaDst: "./log/errDist",
		AsyncDrop: make(chan OutputMeta), // 下行数据异步落盘同步通道

		// jbzhou5 性能资源日志保存目录
		// ResourcesDst = "./"
		// jbzhou5 Prometheus并发协程计数器
		ConcurrencyCnt: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "Xtest_Concurrency_Go_Routine",
			Help: "The total number of processed events",
		}),
		// jbzhou5 Prometheus监听参数
		CpuPer: promauto.NewGauge(prometheus.GaugeOpts{ // CPU 利用率
			Name: "Xtest_CPU_Percent",
			Help: "Xtest cpu percent",
		}),
		MemPer: promauto.NewGauge(prometheus.GaugeOpts{ // 内存利用率
			Name: "Xtest_MEM_Percent",
			Help: "Xtest mem percent",
		}),
	}
}

func (c *Conf) ConfInit(conf *utils.Configure) error {
	if err := c.secParseXtest(conf); err != nil {
		return err
	}

	if err := c.secParseSvc(conf); err != nil {
		return err
	}

	if err := c.secParseHeader(conf); err != nil {
		return err
	}

	if err := c.secParseParams(conf); err != nil {
		return err
	}

	if err := c.secParseData(conf); err != nil {
		return err
	}

	//jbzhou5 输入流
	if c.InputCmd { // 开启命令行输入
		c.Payload = []string{} // 清空payload输入
		if err := c.secParseCmd(conf); err != nil {
			return err
		}
	} else {
		if err := c.secParsePl(conf); err != nil {
			return err
		}
	}

	//输出流
	if err := c.secParseEp(conf); err != nil {
		return err
	}

	if err := c.secParseDStream(conf); err != nil {
		return err
	}
	return nil
}

func (c *Conf) secParseXtest(conf *utils.Configure) error {
	secTmp := "xtest"
	if taddrs, err := conf.GetString(secTmp, "taddrs"); err == nil {
		c.Taddrs = taddrs
	}
	return nil
}

func (c *Conf) secParseEp(conf *utils.Configure) error {
	for _, sec := range c.Expect {
		meta := protocol.MetaDesc{}
		meta.Name, _ = conf.GetString(sec, "name")
		typ, _ := conf.GetString(sec, "type")
		switch typ {
		case "audio":
			meta.DataType = protocol.MetaDesc_AUDIO
		case "text":
			meta.DataType = protocol.MetaDesc_TEXT
		case "image":
			meta.DataType = protocol.MetaDesc_IMAGE
		case "video":
			meta.DataType = protocol.MetaDesc_VIDEO
		default:
			return errors.New("input invalid payload type, section: " + sec)
		}
		meta.Attribute = make(map[string]string)
		descstr, _ := conf.GetString(sec, "describe")
		descarr := strings.Split(descstr, ";")
		for _, desc := range descarr {
			tmp := strings.Split(desc, "=")
			if len(tmp) == 2 {
				meta.Attribute[tmp[0]] = tmp[1]
			}
		}
		// 期望输出流
		c.DownExpect = append(c.DownExpect, meta)
	}
	return nil
}

//解析每一个数据字段
func (c *Conf) secParsePl(conf *utils.Configure) error {
	for _, sec := range c.Payload {
		meta := InputMeta{}
		meta.Name, _ = conf.GetString(sec, "name")
		meta.DataSrc, _ = conf.GetString(sec, "inputSrc")
		meta.SliceOn, _ = conf.GetInt(sec, "sliceOn") // default slice off
		meta.UpSlice, _ = conf.GetInt(sec, "sliceSize")
		meta.UpInterval, _ = conf.GetInt(sec, "interval")
		typ, _ := conf.GetString(sec, "type")
		switch typ {
		case "audio":
			meta.DataType = protocol.MetaDesc_AUDIO
		case "text":
			meta.DataType = protocol.MetaDesc_TEXT
		case "image":
			meta.DataType = protocol.MetaDesc_IMAGE
		case "video":
			meta.DataType = protocol.MetaDesc_VIDEO
		default:
			return errors.New("input invalid payload type, section: " + sec)
		}
		meta.DataDesc = make(map[string]string)
		descstr, _ := conf.GetString(sec, "describe")
		descarr := strings.Split(descstr, ";")
		for _, desc := range descarr {
			tmp := strings.Split(desc, "=")
			if len(tmp) == 2 {
				meta.DataDesc[tmp[0]] = tmp[1]
			}
		}

		// check DataSrc , check sliceSize & interval
		if len(meta.DataSrc) == 0 || (meta.SliceOn != 0 && meta.UpSlice <= 0) || meta.UpInterval <= 0 {
			return errors.New("invalid payload configure, section: " + sec)
		}
		// read upstream valid file list
		meta.DataList = make([][]byte, 0, 1)
		fi, err := os.Stat(meta.DataSrc)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			data, err := util.ReadDir(meta.DataSrc, c.FileNameSeq, c.FileSorted)
			if err != nil {
				return err
			}
			meta.DataList = append(meta.DataList, data...)
		} else if fi.Size() != 0 {
			data, err := ioutil.ReadFile(meta.DataSrc)
			if err != nil {
				fmt.Printf("read file %s fail, %s", meta.DataSrc, err.Error())
				return err
			}
			meta.DataList = append(meta.DataList, data)
		}

		println(len(meta.DataList))

		// 判定len(meta.FileList)
		if len(meta.DataList) == 0 {
			return errors.New("can't read valid file from payload source path " + sec)
		}

		// 上行数据流
		c.UpStreams = append(c.UpStreams, meta)
	}
	println(len(c.UpStreams))
	return nil
}

func (c *Conf) secParseSvc(conf *utils.Configure) error {
	secTmp := "svcMode"
	if sn, err := conf.GetString(secTmp, "service"); err == nil {
		c.SvcName = sn
	}

	if si, err := conf.GetString(secTmp, "svcId"); err == nil {
		c.SvcId = si
	}

	if ts, err := conf.GetString(secTmp, "sub"); err == nil {
		c.TestSub = ts
	}

	if to, err := conf.GetInt(secTmp, "timeout"); err == nil {
		c.TimeOut = to
	}

	if ld, err := conf.GetInt(secTmp, "timeLoss"); err == nil {
		c.LossDeviation = ld
	}

	if mt, err := conf.GetInt(secTmp, "multiThr"); err == nil {
		c.MultiThr = mt
	}
	c.DropThr = c.MultiThr
	if cnt, err := conf.GetInt64(secTmp, "loopCnt"); err == nil {
		c.LoopCnt.Store(int64(cnt))
	}

	if perfOn, err := conf.GetBool(secTmp, "perfOn"); err == nil {
		c.PerfConfigOn = perfOn
	}

	if perfLevel, err := conf.GetInt(secTmp, "perfLevel"); err == nil {
		c.PerfLevel = perfLevel
		if perfLevel != 0 && perfLevel != 1 && perfLevel != 2 {
			fmt.Println("perfLevel set invalid. use default: 0")
			c.PerfLevel = 0
		}
	}
	if c.ReqMode == 0 {
		c.PerfLevel = 0
	}

	if rm, err := conf.GetInt(secTmp, "sessMode"); err == nil {
		c.ReqMode = rm
	}
	linearms, _ := conf.GetInt(secTmp, "linearMs")
	if linearms > 0 && c.MultiThr > 0 {
		c.LinearNs = (linearms * 1000 * 1000) / c.MultiThr
	}

	// jbzhou5 当模式为非会话且配置了cmd输入，才开启手动输入
	if inputCmd, err := conf.GetBool(secTmp, "inputCmd"); err == nil && c.ReqMode == 0 {
		c.InputCmd = inputCmd
	}

	// jbzhou5 Prometheus写入开关
	if prometheusSwitch, err := conf.GetBool(secTmp, "prometheus_switch"); err == nil {
		c.PrometheusSwitch = prometheusSwitch
	}

	// jbzhou5 Prometheus指标端口
	if prometheusPort, err := conf.GetInt(secTmp, "prometheus_port"); err == nil {
		c.PrometheusPort = prometheusPort
	}

	// jbzhou5 资源监控绘图开关
	if plot, err := conf.GetBool(secTmp, "plot"); err == nil {
		c.Plot = plot
	}

	// jbzhou5 绘图保存路径
	if plotFile, err := conf.GetString(secTmp, "plot_file"); err == nil {
		c.PlotFile = plotFile
	}

	// jbzhou5 设置读入文件顺序
	if fileSorted, err := conf.GetInt(secTmp, "file_sorted"); err == nil {
		c.FileSorted = fileSorted
	}

	// jbzhou5 文件名分割方式
	if filenameSeq, err := conf.GetString(secTmp, "file_name_seq"); err == nil && filenameSeq != "" {
		c.FileNameSeq = filenameSeq
	}

	c.AsyncDrop = make(chan OutputMeta, c.MultiThr*10) // channel长度取并发数*10, channel满则同步写.
	return nil
}

func (c *Conf) secParseHeader(conf *utils.Configure) error {
	secData := conf.GetSection("header")
	if secData != nil {
		kv, ok := secData.(map[string]interface{})
		if ok {
			for key, value := range kv {
				var valStr string
				switch value.(type) {
				case string:
					valStr = value.(string)
				case int:
					valStr = strconv.Itoa(value.(int))
				case int64:
					valStr = strconv.FormatInt(value.(int64), 10)
				case uint:
					valStr = strconv.FormatUint(uint64(value.(uint)), 10)
				case uint64:
					valStr = strconv.FormatUint(value.(uint64), 10)
				case bool:
					valStr = strconv.FormatBool(value.(bool))
				case float64:
					valStr = strconv.FormatFloat(value.(float64), 'f', -1, 64)
				case float32:
					valStr = strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)
				default:
					return errors.New("invalid header configure, type/key " + reflect.TypeOf(value).String() + key)
				}
				c.Header[key] = valStr
			}
		} else {
			return errors.New("invalid section header key, type of key must be string")
		}
	}
	return nil
}

func (c *Conf) secParseParams(conf *utils.Configure) error {
	secData := conf.GetSection("parameter")
	if secData != nil {
		kv, ok := secData.(map[string]interface{})

		if ok {
			for key, value := range kv {
				var valStr string
				switch value.(type) {
				case string:
					valStr = value.(string)
				case int:
					valStr = strconv.Itoa(value.(int))
				case int64:
					valStr = strconv.FormatInt(value.(int64), 10)
				case uint:
					valStr = strconv.FormatUint(uint64(value.(uint)), 10)
				case uint64:
					valStr = strconv.FormatUint(value.(uint64), 10)
				case bool:
					valStr = strconv.FormatBool(value.(bool))
				case float64:
					valStr = strconv.FormatFloat(value.(float64), 'f', -1, 64)
				case float32:
					valStr = strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)
				default:
					return errors.New("invalid parameter configure, type/key " + reflect.TypeOf(value).String() + key)
				}
				c.Params[key] = valStr
			}
		} else {
			return errors.New("invalid section parameter key, type of key must be string")
		}
	}
	return nil
}

func (c *Conf) secParseData(conf *utils.Configure) error {
	inputList, _ := conf.GetString("data", "payload")
	c.Payload = strings.Split(inputList, ";")
	if len(c.Payload) == 0 {
		return errors.New("invalid configure: data.payload")
	}

	outputList, _ := conf.GetString("data", "expect")
	c.Expect = strings.Split(outputList, ";")
	if len(outputList) == 0 {
		return errors.New("invalid configure: data.expect")
	}
	return nil
}

func (c *Conf) secParseDStream(conf *utils.Configure) error {
	if op, err := conf.GetInt("downstream", "output"); err == nil {
		c.Output = op
	}
	if opd, err := conf.GetString("downstream", "outputDst"); err == nil {
		c.OutputDst = opd
	}
	// 下行数据输出目录/文件预处理
	switch c.Output {
	case 0: // 输出至统一公共文件 outputDst(!IsDir())
		fi, err := os.Stat(c.OutputDst)
		if err == nil && fi.IsDir() {
			err = os.RemoveAll(c.OutputDst)
			if err != nil {
				return err
			}
		}
		fp, err := os.Create(c.OutputDst)
		if err != nil {
			return err
		}
		_ = fp.Close()

	case 1: // 输出至目录outputDst(IsDir()), 一次请求的输出以单独的文件存储
		fi, err := os.Stat(c.OutputDst)
		if err == nil && !fi.IsDir() {
			err = os.Remove(c.OutputDst)
			if err != nil {
				return err
			}
		}
		err = os.MkdirAll(c.OutputDst, os.ModeDir)
		if err != nil {
			return err
		}

	case 2: // 输出至终端
		// nothing to do, print to terminal
	default:
		return errors.New("downstream output invalid, output=0/1/2")
	}
	return nil
}

//jbzhou5 解析命令行输入的数据
func (c *Conf) secParseCmd(conf *utils.Configure) error {
	meta := InputMeta{}
	meta.Name = "CMD"
	meta.DataSrc = "CMD"
	meta.SliceOn = 0
	meta.UpSlice = 1280
	meta.UpInterval = 40

	meta.DataDesc = make(map[string]string)
	descstr := "encoding=utf8;compress=gzip"
	descarr := strings.Split(descstr, ";")
	for _, desc := range descarr {
		tmp := strings.Split(desc, "=")
		if len(tmp) == 2 {
			meta.DataDesc[tmp[0]] = tmp[1]
		}
	}
	var data string
	n, err := Input(data)
	if n == 0 || err != nil {
		return errors.New("cmd input error")
	}
	// read upstream valid file list
	meta.DataList = make([][]byte, 0, 1)
	meta.DataList = append(meta.DataList, []byte(data))

	println("DataList:", len(meta.DataList))

	// 上行数据流
	c.UpStreams = append(c.UpStreams, meta)
	return nil
}
