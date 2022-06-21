package _var

import (
	"errors"
	"fmt"
	"github.com/xfyun/aiges/protocol"
	"github.com/xfyun/xsf/utils"
	"go.uber.org/atomic"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
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

var (
	// [svcMode]
	SvcId         string        = "s12345678"
	SvcName       string        = "AIservice"            // dst service name
	TimeOut       int           = 1000                   // 超时时间: ms, 对应加载器waitTime
	LossDeviation int           = 50                     // 自身性能损耗误差, ms.
	MultiThr      int           = 100                    // 请求并发数
	DropThr                     = 100                    // 下行数据异步输出线程数
	LoopCnt       *atomic.Int64 = atomic.NewInt64(10000) // 请求总次数
	ReqMode       bool          = true                   // true: sessMode, false: unSessMode
	LinearNs      int           = 0                      // 并发模型线性增长时间,用于计算并发增长斜率(单位：ns). default:0,瞬时并发压测.
	TestSub       string        = "ase"                  // 测试业务sub, 缺省test
	PerfConfigOn  bool          = false                  //true: 开启性能检测 false: 不开启性能检测
	PerfLevel     int           = 0                      //非会话模式默认0
	//会话模式0: 从发第一帧到最后一帧的性能
	//会话模式1:首结果(发送第一帧到最后一帧的性能)
	//会话模式2:尾结果(发送最后一帧到收到最后一帧的性能)
	// 请求参数对
	Header map[string]string = make(map[string]string)
	Params map[string]string = make(map[string]string)

	Payload []string // 上行数据流
	Expect  []string // 下行数据流

	// 上行数据流配置, 多数据流通过section [data]中payload进行配置
	UpStreams []InputMeta = make([]InputMeta, 0, 1)

	DownExpect []protocol.MetaDesc

	// [downstream]
	Output = 0 // 0：输出至公共文件outputDst(sid+***:data)
	// 1：以独立文件形式输出至目录outputDst(文件名:sid+***)
	// 2：输出至终端
	//-1：不输出
	OutputDst = "./log/result" // output=0时,该项配置输出文件名; output=1时,该项配置输出目录名
	ErrAnaDst = "./log/errDist"
	AsyncDrop chan OutputMeta // 下行数据异步落盘同步通道
)

func ConfInit(conf *utils.Configure) error {

	if err := secParseSvc(conf); err != nil {
		return err
	}

	if err := secParseHeader(conf); err != nil {
		return err
	}

	if err := secParseParams(conf); err != nil {
		return err
	}

	if err := secParseData(conf); err != nil {
		return err
	}

	if err := secParsePl(conf); err != nil {
		return err
	}

	if err := secParseEp(conf); err != nil {
		return err
	}

	if err := secParseDStream(conf); err != nil {
		return err
	}

	return nil
}

func secParseEp(conf *utils.Configure) error {
	for _, sec := range Expect {
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
		DownExpect = append(DownExpect, meta)
	}
	return nil
}

func secParsePl(conf *utils.Configure) error {
	for _, sec := range Payload {
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
			// 遍历目录文件
			files, err := ioutil.ReadDir(meta.DataSrc)
			if err != nil {
				return err
			}
			// 过滤空文件及子目录
			for _, file := range files {
				if !file.IsDir() && file.Size() != 0 {
					data, err := ioutil.ReadFile(meta.DataSrc + "/" + file.Name())
					if err != nil {
						fmt.Printf("read file %s fail, %s ", meta.DataSrc+"/"+file.Name(), err.Error())
						return err
					}
					meta.DataList = append(meta.DataList, data)
				}
			}
		} else if fi.Size() != 0 {
			data, err := ioutil.ReadFile(meta.DataSrc)
			if err != nil {
				fmt.Printf("read file %s fail, %s", meta.DataSrc, err.Error())
				return err
			}
			meta.DataList = append(meta.DataList, data)
		}

		// 判定len(meta.FileList)
		if len(meta.DataList) == 0 {
			return errors.New("can't read valid file from payload source path " + sec)
		}

		// 上行数据流
		UpStreams = append(UpStreams, meta)
	}
	return nil
}

func secParseSvc(conf *utils.Configure) error {
	secTmp := "svcMode"
	if sn, err := conf.GetString(secTmp, "service"); err == nil {
		SvcName = sn
	}

	if si, err := conf.GetString(secTmp, "svcId"); err == nil {
		SvcId = si
	}

	if ts, err := conf.GetString(secTmp, "sub"); err == nil {
		TestSub = ts
	}

	if to, err := conf.GetInt(secTmp, "timeout"); err == nil {
		TimeOut = to
	}

	if ld, err := conf.GetInt(secTmp, "timeLoss"); err == nil {
		LossDeviation = ld
	}

	if mt, err := conf.GetInt(secTmp, "multiThr"); err == nil {
		MultiThr = mt
	}
	DropThr = MultiThr
	if cnt, err := conf.GetInt64(secTmp, "loopCnt"); err == nil {
		LoopCnt.Store(int64(cnt))
	}

	if perfOn, err := conf.GetBool(secTmp, "perfOn"); err == nil {
		PerfConfigOn = perfOn
	}

	if perfLevel, err := conf.GetInt(secTmp, "perfLevel"); err == nil {
		PerfLevel = perfLevel
		if perfLevel != 0 && perfLevel != 1 && perfLevel != 2 {
			fmt.Println("perfLevel set invalid. use default: 0")
			PerfLevel = 0
		}
	}
	if !ReqMode {
		PerfLevel = 0
	}

	if rm, err := conf.GetBool(secTmp, "sessMode"); err == nil {
		ReqMode = rm
	}
	linearms, _ := conf.GetInt(secTmp, "linearMs")
	if linearms > 0 && MultiThr > 0 {
		LinearNs = (linearms * 1000 * 1000) / MultiThr
	}

	AsyncDrop = make(chan OutputMeta, MultiThr*10) // channel长度取并发数*10, channel满则同步写.
	return nil
}

func secParseHeader(conf *utils.Configure) error {
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
				Header[key] = valStr
			}
		} else {
			return errors.New("invalid section header key, type of key must be string")
		}
	}
	return nil
}

func secParseParams(conf *utils.Configure) error {
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
				Params[key] = valStr
			}
		} else {
			return errors.New("invalid section parameter key, type of key must be string")
		}
	}
	return nil
}

func secParseData(conf *utils.Configure) error {
	inputList, _ := conf.GetString("data", "payload")
	Payload = strings.Split(inputList, ";")
	if len(Payload) == 0 {
		return errors.New("invalid configure: data.payload")
	}

	outputList, _ := conf.GetString("data", "expect")
	Expect = strings.Split(outputList, ";")
	if len(outputList) == 0 {
		return errors.New("invalid configure: data.expect")
	}
	return nil
}

func secParseDStream(conf *utils.Configure) error {
	if op, err := conf.GetInt("downstream", "output"); err == nil {
		Output = op
	}
	if opd, err := conf.GetString("downstream", "outputDst"); err == nil {
		OutputDst = opd
	}
	// 下行数据输出目录/文件预处理
	switch Output {
	case 0: // 输出至统一公共文件 outputDst(!IsDir())
		fi, err := os.Stat(OutputDst)
		if err == nil && fi.IsDir() {
			err = os.RemoveAll(OutputDst)
			if err != nil {
				return err
			}
		}
		fp, err := os.Create(OutputDst)
		if err != nil {
			return err
		}
		_ = fp.Close()

	case 1: // 输出至目录outputDst(IsDir()), 一次请求的输出以单独的文件存储
		fi, err := os.Stat(OutputDst)
		if err == nil && !fi.IsDir() {
			err = os.Remove(OutputDst)
			if err != nil {
				return err
			}
		}
		err = os.MkdirAll(OutputDst, os.ModeDir)
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
