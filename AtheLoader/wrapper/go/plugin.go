package _go

type DataType int
type DataStatus int
type CallBackPtr func(hdl interface{}, resp []wrapperData)

const (
	DataText  DataType = 0 // 文本数据
	DataAudio DataType = 1 // 音频数据
	DataImage DataType = 2 // 图像数据
	DataVideo DataType = 3 // 视频数据
	DataPer   DataType = 4 // 个性化数据

	DataBegin    DataStatus = 0 // 首数据
	DataContinue DataStatus = 1 // 中间数据
	DataEnd      DataStatus = 2 // 尾数据
	DataOnce     DataStatus = 3 // 非会话单次输入
)

type wrapperData struct {
	key      string     // 数据标识
	data     []byte     // 数据实体
	desc     string     // 数据描述
	encoding string     // 数据编码
	typ      DataType   // 数据类型
	status   DataStatus // 数据状态
}

// plugin interface
type (
	WrapperInit      func(cfg map[string]string) (err error)
	WrapperFini      func() (err error)
	WrapperVersion   func() (version string)
	WrapperLoadRes   func(res wrapperData, resId int) (err error)
	WrapperUnloadRes func(resId int) (err error)
	WrapperCreate    func(param map[string]string, prsIds []int, cb CallBackPtr) (hdl interface{}, err error)
	WrapperDestroy   func(hdl interface{}) (err error)
	WrapperWrite     func(hdl interface{}, req []wrapperData) (err error)
	WrapperRead      func(hdl interface{}) (resp []wrapperData, err error)
	WrapperExec      func(param map[string]string, req []wrapperData) (resp []wrapperData, err error)
	WrapperDebugInfo func(hdl interface{}) (debug string)
)
