package protocol

import "strconv"

// common headers
const (
	Uid       = "uid"
	AppId     = "appid"
	CloudId   = "cloud_id"
	ComposeId = "compose_id"
	Sub       = "sub"
	MeterPara = "meter_param"
	DeviceId  = "did"
	Imei      = "device.imei"
	SessionId = "sid"

	DataCleanTTL = "data_cleaner_ttl"
)

// payload attribute
const (
	Encoding   = "encoding"
	Compress   = "compress"
	SampleRate = "sample_rate"
	FrameRate  = "frame_rate"
	FrameSize  = "frame_size"
	Channels   = "channels"
	BitDepth   = "bit_depth"
	Height     = "height"
	Width      = "width"
	Sequence   = "seq"
	Status     = "status"
	AudioUrl   = "audio_url"
)

// TODO 关键参数校验;
type BaseAttr struct {
	Name   string // 数据名
	Type   int    // 数据类型
	Seq    int    // 数据序号
	Status int    // 数据状态
}

// 文本数据属性
type TextAttr struct {
	Encoding string // 文本编码
	Compress string // 压缩格式
	Seq      int    // 序号
	Status   int    // 状态
}

// 音频数据属性
type AudioAttr struct {
	Encoding   string // 编解码
	SampleRate string // 采样率
	FrameSize  string // 帧大小(开源speex)
	Channels   int    // 通道数
	BitDepth   int    // 位深
	Seq        int    // 序号
	Status     int    // 状态
}

type VideoAttr struct {
	Encoding  string // 编解码
	FrameRate string // 帧率
	Width     string // 分辨率-宽
	Height    string // 分辨率-高
	Seq       int    // 序号
	Status    int    // 状态
}

type ImageAttr struct {
	Encoding string // 编解码
	Seq      int    // 序号
	Status   int    // 状态
}

func GetAllAttr(desc *MetaDesc) (attr interface{}) {
	seq, _ := strconv.Atoi(desc.Attribute[Sequence])
	status, _ := strconv.Atoi(desc.Attribute[Status])
	switch desc.DataType {
	case MetaDesc_TEXT:
		return TextAttr{desc.Attribute[Encoding],
			desc.Attribute[Compress], seq, status}
	case MetaDesc_AUDIO:
		depth, _ := strconv.Atoi(desc.Attribute[BitDepth])
		channel, _ := strconv.Atoi(desc.Attribute[Channels])
		return AudioAttr{desc.Attribute[Encoding],
			desc.Attribute[SampleRate], desc.Attribute[FrameSize],
			channel, depth, seq, status}
	case MetaDesc_IMAGE:
		return ImageAttr{desc.Attribute[Encoding], seq, status}
	case MetaDesc_VIDEO:
		//rate,_ := strconv.Atoi(desc.Attribute[FrameRate])
		return VideoAttr{desc.Attribute[Encoding], desc.Attribute[FrameRate],
			desc.Attribute[Width], desc.Attribute[Height], seq, status}
	default:
		return nil
	}
	return
}

func GetBaseAttr(desc *MetaDesc) (attr BaseAttr) {
	// TODO desc为空, 崩溃问题；调整报错;
	seq, _ := strconv.Atoi(desc.Attribute[Sequence])
	status, _ := strconv.Atoi(desc.Attribute[Status])
	return BaseAttr{desc.Name, int(desc.DataType), seq, status}
}
