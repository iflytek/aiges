package buffer

/*
	seqBuf数据缓存排序模块
*/
// 有序数据状态;
type DataStatus int

const (
	DataStatusFirst    DataStatus = 0
	DataStatusContinue DataStatus = 1
	DataStatusLast     DataStatus = 2
	DataStatusOnce     DataStatus = 3
)

// 数据类型
type DataType int

const (
	DataText  DataType = 0 // 文本
	DataAudio DataType = 1 // 音频
	DataImage DataType = 2 // 图像
	DataVideo DataType = 3 // 视频
)

type DataMeta struct {
	Data     interface{}       // 数据;
	DataId   string            // 数据id;
	FrameId  uint              // 排序id;
	Status   DataStatus        // 数据状态;
	DataType DataType          // 数据类型;
	Format   string            // 数据编码格式
	Encoding string            // 数据压缩格式
	Desc     map[string][]byte // 数据描述
}

func DataTypeToString(dataType DataType) (typeStr string) {
	switch dataType {
	case DataText:
		typeStr = "text"
	case DataAudio:
		typeStr = "audio"
	case DataImage:
		typeStr = "image"
	case DataVideo:
		typeStr = "video"
	default:
		typeStr = "invalid data type"
	}
	return
}
