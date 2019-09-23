package _var

import (
	"errors"
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"go.uber.org/atomic"
	"io/ioutil"
	"os"
	"protocol"
	"strings"
)

const (
	CliName = "xtest"
	CliVer  = "2.0.1"
)

type InputMeta struct {
	DataSrc    string                     // 上行实体数据来源;数据集则配置对应目录
	SliceOn    int                        // 上行数据切片开关, !0:切片. 0:不切片
	UpSlice    int                        // 上行数据切片大小: byte
	UpInterval int                        // slice发包间隔: ms
	DataType   protocol.MetaData_DataType // audio/text/image/video
	DataFmt    string
	DataEnc    string
	DataDesc   map[string]string

	// DataList map[string/*file*/] []byte /*data*/
	DataList [][]byte /*data*/
}

type OutputMeta struct {
	Sid  string // sid
	Type string // 下行数据类型
	Fmt  string // 下行数据格式
	Enc  string // 下行数据编码
	//Desc string			// map对以";"分隔
	Data []byte // 下行数据实体
}

var (
	// [svcMode]
	SvcName       string        = "AIservice"            // dst service name
	TimeOut       int           = 1000                   // 超时时间: ms, 对应加载器waitTime
	LossDeviation int           = 50                     // 自身性能损耗误差, ms.
	MultiThr      int           = 100                    // 请求并发数
	DropThr                     = 100                    // 下行数据异步输出线程数
	LoopCnt       *atomic.Int64 = atomic.NewInt64(10000) // 请求总次数
	ReqMode       int           = 1                      // 1: sessMode, 0: unSessMode
	// unSessMode: 对应一次上下行交互, 存在一次Exec接口调用(内部调用AIInput+AIExcp)
	// sessMode: 对应多次上下行交互, 存在多次Write/Read接口调用(内部调用AIInput/AIOutput/AIExcp)
	LinearNs int    = 0     // 并发模型线性增长时间,用于计算并发增长斜率(单位：ns). default:0,瞬时并发压测.
	TestSub  string = "tst" // 测试业务sub, 缺省test

	// 请求参数对
	UpParams map[string]string = make(map[string]string)

	// [upstream]
	// 上行数据流配置, 多数据流通过section [upstream-N]进行配置
	UpStreams []InputMeta = make([]InputMeta, 0, 1)

	// [downstream]
	Output = 0 // 0：输出至公共文件outputDst(sid+***:data)
	// 1：以独立文件形式输出至目录outputDst(文件名:sid+***)
	// 2：输出至终端
	//-1：不输出
	OutputDst                 = "./log/result" // output=0时,该项配置输出文件名; output=1时,该项配置输出目录名
	ErrAnaDst                 = "./log/errDist"
	AsyncDrop chan OutputMeta // 下行数据异步落盘同步通道
)

func ConfInit(conf *utils.Configure) error {

	// 获取[svcMode]
	secTmp := "svcMode"
	if sn, err := conf.GetString(secTmp, "service"); err == nil {
		SvcName = sn
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

	if rm, err := conf.GetInt(secTmp, "reqMode"); err == nil {
		ReqMode = rm
	}
	linearms, _ := conf.GetInt(secTmp, "linearMs")
	if linearms > 0 && MultiThr > 0 {
		LinearNs = (linearms * 1000 * 1000) / MultiThr
	}
	paraList, _ := conf.GetString(secTmp, "reqPara")
	arrTmp := strings.Split(paraList, ";")
	for _, para := range arrTmp {
		kv := strings.Split(para, "=")
		if len(kv) != 2 {
			return errors.New("input invalid reqPara")
		}
		UpParams[kv[0]] = kv[1]
	}
	AsyncDrop = make(chan OutputMeta, MultiThr*10) // channel长度取并发数*10, channel满则同步写.

	// 获取[downstream]
	secTmp = "downstream"
	if op, err := conf.GetInt(secTmp, "output"); err == nil {
		Output = op
	}
	if opd, err := conf.GetString(secTmp, "outputDst"); err == nil {
		OutputDst = opd
	}

	// 获取配置section, 判定upstream*
	secs := conf.GetSecs()
	for _, sec := range secs {
		if strings.Contains(sec, "upstream") {
			meta := InputMeta{}
			meta.DataSrc, _ = conf.GetString(sec, "inputSrc")
			meta.SliceOn, _ = conf.GetInt(sec, "sliceOn") // default slice off
			meta.UpSlice, _ = conf.GetInt(sec, "sliceSize")
			meta.UpInterval, _ = conf.GetInt(sec, "interval")
			typ, _ := conf.GetString(sec, "type")
			switch typ {
			case "audio":
				meta.DataType = protocol.MetaData_AUDIO
			case "text":
				meta.DataType = protocol.MetaData_TEXT
			case "image":
				meta.DataType = protocol.MetaData_IMAGE
			case "video":
				meta.DataType = protocol.MetaData_VIDEO
			default:
				return errors.New("input invalid upstream type, section: " + sec)
			}
			meta.DataFmt, _ = conf.GetString(sec, "format")
			meta.DataEnc, _ = conf.GetString(sec, "encoding")
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
				return errors.New("input invalid upstream, section: " + sec)
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
				return errors.New("can't read valid file from upstream " + sec)
			}

			// 上行数据流
			UpStreams = append(UpStreams, meta)
		}
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
