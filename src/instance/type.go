package instance

import (
	"github.com/xfyun/aiges/frame"
	"github.com/xfyun/aiges/protocol"
	"unsafe"
)

// 服务框架协议相关交互参数
const (
	nrtTask = frame.ParaNrtTask
)

// 离线动态授权相关协议
const (
	reqLbType  = "lbType"
	reqDynLic  = "authzNo"
	reqDynTime = "timeout"
)

func checkAudioRate(rate int) (errInfo error) {
	switch rate {
	case AudioRate16k:
	case AudioRate8k:
	default:
		errInfo = frame.ErrorInvalidSampleRate
	}
	return
}

const (
	hstSessionSize    int   = 1000 // 历史会话记录维护列表大小;
	seqBufTimeout     uint  = 1000 // 数据排序超时时间5s;
	seqAudSize        uint  = 500  // 输入缓冲区有序队列大小;
	seqRltSize        uint  = 100  // 输出缓冲区有序队列大小;
	sessTimeoutCnt    int32 = 10   // 会话超时10s;
	AudioRate16k      int   = 16000
	AudioRate8k       int   = 8000
	ResampleQuality   int   = 10
	syncRltWaitTime   int   = 50  // 同步结果读取waitTime 50ms
	syncRltWaitCnt    int   = 100 // 同步结果读取超时次数
	defaultEventCache int   = 32000

	eventUid     string = "uid"
	eventAppid   string = "appid"
	dataSrc      string = "dataSrc"
	dataHttp     string = "http"
	dataHttpUrl  string = "url"
	dataS3       string = "s3"
	dataS3Access string = "access"
	dataS3Secret string = "secret"
	dataS3Ep     string = "endpoint"
	dataS3Bucket string = "bucket"
	dataS3Key    string = "key"
	dataClient   string = "client"

	downMethod string = "downMethod"
	downAsync  string = "async"
	downSync   string = "sync"
)

type UserEvent int

const (
	EventInit UserEvent = 1 << iota // 用户初始化事件
	EventFini                       // 逆初始化事件

	EventNew   // 资源申请事件
	EventDel   // 资源释放事件
	EventExcp  // 异常销毁事件
	EventDebug // debug事件,获取debug信息

	EventRead     // 数据读事件
	EventWrite    // 数据写事件
	EventOnceExec // 非会话处理事件,对应once请求

	EventResLoad   // 资源load事件
	EventResUnload // 资源卸载事件
)

type eventStorage struct {
	data   interface{}
	dataId string             // 区分数据流
	Desc   *protocol.MetaDesc // 数据属性;
}

type DataMeta struct {
	DataId     string // 数据流标识：用于多输入||多输出场景;
	Data       []byte // 数据
	DataType   int    // 数据类型
	DataStatus int    // 数据状态
	DataFrame  uint   // 数据分片id
	//DataFmt    string // 数据格式,用于wrapper引擎
	//DataEnc    string // 数据编码,用于wrapper自解码引擎
	DataDesc map[string]string // 数据描述
}

type ActMsg struct {
	//WrapperHdl interface{} // 适配c/go引擎句柄
	WrapperHdl unsafe.Pointer `json:"-"`
	// 计算数据
	DeliverData []DataMeta

	// 引擎资源申请
	Params    map[string]string // 服务参数对
	PersonIds []int             // 个性化数据集

	// 调试&联调信息
	Debug string // debug信息;

	// 个性化load/unload
	PsrData []byte            // 个性化数据
	PsrDesc map[string]string // 个性化描述
	PsrId   int               // 个性化id
	PsrKey  string            // 个性化key

	AsyncErr  error // 异步回调错误
	AsyncCode int   // 异步回调错误码(c接口错误码)
}

// 用户自定义注册事件
type UsrAct func(hdl string, req *ActMsg) (resp ActMsg, errNum int, errInfo error)

func DataTypeToString(dataType protocol.MetaDesc_DataType) (typeStr string) {
	switch dataType {
	case protocol.MetaDesc_TEXT:
		typeStr = "text"
	case protocol.MetaDesc_AUDIO:
		typeStr = "audio"
	case protocol.MetaDesc_IMAGE:
		typeStr = "image"
	case protocol.MetaDesc_VIDEO:
		typeStr = "video"
	default:
		typeStr = "invalid data type"
	}
	return
}
