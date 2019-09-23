package widget

import "C"
import (
	"frame"
	"instance"
)

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

// plugin symbol list
const (
	wrapperInitSym    = "WrapperInit"
	wrapperFiniSym    = "WrapperFini"
	wrapperVerSym     = "WrapperVersion"
	wrapperLoadSym    = "WrapperLoadRes"
	wrapperUnloadSym  = "WrapperUnloadRes"
	wrapperCreateSym  = "WrapperCreate"
	wrapperWriteSym   = "WrapperWrite"
	wrapperReadSym    = "WrapperRead"
	wrapperDestroySym = "WrapperDestroy"
	wrapperExecSym    = "WrapperExec"
	wrapperDebugSym   = "WrapperDebugInfo"
)

type (
	wrapperInitPtr      func(cfg map[string]string) (err error)
	wrapperFiniPtr      func() (err error)
	wrapperVersionPtr   func() (version string)
	wrapperLoadResPtr   func(res wrapperData, resId int) (err error)
	wrapperUnloadResPtr func(resId int) (err error)
	wrapperCreatePtr    func(param map[string]string, prsIds []int, cb CallBackPtr) (hdl interface{}, err error)
	wrapperDestroyPtr   func(hdl interface{}) (err error)
	wrapperWritePtr     func(hdl interface{}, req []wrapperData) (err error)
	wrapperReadPtr      func(hdl interface{}) (resp []wrapperData, err error)
	wrapperOnceExecPtr  func(param map[string]string, req []wrapperData) (resp []wrapperData, err error)
	wrapperDebugInfoPtr func(hdl interface{}) (debug string)
)

func pluginInit(cfg map[string]string) (errNum int, errInfo error) {
	errInfo = pluginHdls[wrapperInitSym].(func(cfg map[string]string) (err error))(cfg)
	if errInfo != nil {
		errNum = frame.WrapperInitErr
	}
	return
}

func pluginFini() (errNum int, errInfo error) {
	errInfo = pluginHdls[wrapperFiniSym].(func() (err error))()
	if errInfo != nil {
		errNum = frame.WrapperFiniErr
	}
	return
}

func pluginVersion() (ver string) {
	return pluginHdls[wrapperVerSym].(func() (version string))()
}

// 资源加载卸载管理适配接口;
func pluginLoadRes(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	var res wrapperData
	res.data = req.PsrData
	res.desc = req.PsrDesc
	res.typ = DataPer
	errInfo = pluginHdls[wrapperLoadSym].(func(res wrapperData, resId int) (err error))(res, req.PsrId)
	if errInfo != nil {
		errNum = frame.WrapperLoadErr
	}
	return
}

func pluginUnloadRes(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	errInfo = pluginHdls[wrapperUnloadSym].(func(resId int) (err error))(req.PsrId)
	if errInfo != nil {
		errNum = frame.WrapperUnloadErr
	}
	return
}

// 资源申请行为
func pluginCreate(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	resp.WrapperHdl, errInfo = pluginHdls[wrapperCreateSym].(func(param map[string]string, prsIds []int, cb CallBackPtr) (hdl interface{}, err error))(req.Params, req.PersonIds, nil)
	if errInfo != nil {
		errNum = frame.WrapperCreateErr
	}
	return
}

// 资源释放行为
func pluginDestroy(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	errInfo = pluginHdls[wrapperDestroySym].(func(hdl interface{}) (err error))(req.WrapperHdl)
	if errInfo != nil {
		errNum = frame.WrapperDestroyErr
	}
	return
}

// 交互异常行为
func pluginExcp(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	resp, errNum, errInfo = pluginDestroy(handle, req)
	return
}

// 数据写行为
func pluginWrite(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	// 支持多数据写入
	inputs := make([]wrapperData, 0, len(req.DeliverData))
	for _, ele := range req.DeliverData {
		var data wrapperData
		data.typ = DataType(ele.DataType)
		data.status = DataStatus(ele.DataStatus)
		data.data = ele.Data
		data.desc = ele.DataDesc
		data.key = ele.DataId
		data.encoding = ele.DataFmt
		inputs = append(inputs, data)
	}

	errInfo = pluginHdls[wrapperWriteSym].(func(hdl interface{}, req []wrapperData) (err error))(req.WrapperHdl, inputs)
	if errInfo != nil {
		errNum = frame.WrapperWriteErr
	}
	return
}

// 数据读行为
func pluginRead(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	outputs, errInfo := pluginHdls[wrapperReadSym].(func(hdl interface{}) (resp []wrapperData, err error))(req.WrapperHdl)
	if errInfo != nil {
		errNum = frame.WrapperReadErr
	} else if len(outputs) > 0 {
		resp.DeliverData = make([]instance.DataMeta, 0, len(outputs))
		for k := range outputs {
			var ele instance.DataMeta
			ele.DataType = int(outputs[k].typ)
			ele.DataStatus = int(outputs[k].status)
			ele.DataDesc = outputs[k].desc
			ele.Data = outputs[k].data
			ele.DataId = outputs[k].key
			ele.DataFmt = outputs[k].encoding
			resp.DeliverData = append(resp.DeliverData, ele)
		}
	}

	return
}

func pluginOnceExec(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	// 支持多数据写入&输出
	inputs := make([]wrapperData, 0, len(req.DeliverData))
	for _, ele := range req.DeliverData {
		var data wrapperData
		data.typ = DataType(ele.DataType)
		data.status = DataStatus(ele.DataStatus)
		data.data = ele.Data
		data.desc = ele.DataDesc
		data.key = ele.DataId
		data.encoding = ele.DataFmt
		inputs = append(inputs, data)
	}

	outputs, errInfo := pluginHdls[wrapperExecSym].(func(param map[string]string, req []wrapperData) (resp []wrapperData, err error))(req.Params, inputs)
	if errInfo != nil {
		errNum = frame.WrapperExecErr
	} else if len(outputs) > 0 {
		resp.DeliverData = make([]instance.DataMeta, 0, len(outputs))
		for k := range outputs {
			var ele instance.DataMeta
			ele.DataType = int(outputs[k].typ)
			ele.DataStatus = int(outputs[k].status)
			ele.DataDesc = outputs[k].desc
			ele.Data = outputs[k].data
			ele.DataId = outputs[k].key
			ele.DataFmt = outputs[k].encoding
			resp.DeliverData = append(resp.DeliverData, ele)
		}
	}
	return
}

// 计算debug数据
func pluginDebug(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	resp.Debug = pluginHdls[wrapperDebugSym].(func(hdl interface{}) (debug string))(req.WrapperHdl)
	return
}
