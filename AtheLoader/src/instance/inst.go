package instance

import "C"
import (
	"buffer"
	"catch"
	"codec"
	"conf"
	"dp"
	"encoding/json"
	"errors"
	"frame"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"github.com/golang/protobuf/proto"
	"protocol"
	"storage"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
	服务实例,集成服务框架内部各功能模块,进行加载器层服务处理,包括但不限于：
	1. 进行数据缓存排序
	2. 进行数据编解码
	3. 进行数据预处理
	4. 进行个性化判定&查询
	5. 进行参数校验&控制过滤
	6. 进行计量授权管理
	7. 进行数据收集、链路跟踪记录及性能监控
	8. 进行错误码拆分管理及告警
	9. 事件-行为 判定
*/

type jsonDown struct {
	ReqParam map[string]string `json:"reqParam"`
	Data  []byte `json:"data"`
	SpanMeta string `json:"spanMeta"`
	Sid   string `json:"sid"`
	Ret int		`json:"ret"`
	ErrDesc string  `json:"errDesc"`
}

// 异步协程(单一)读取有序解码重采样数据,写入引擎服务对应判定行为接口;
type ServiceInst struct {
	instHdl    string      // 服务实例唯一标识,sid;
	wrapperHdl interface{} // wrapper层实例句柄
	alive      bool        // 服务实例状态flag
	aliveLock  sync.RWMutex

	tool     *xsf.ToolBox
	dpar     map[string]*dp.AudioResampler  // 数据处理(重采样)
	deCodec  map[string]*codec.AucodecInst // 输入数据处理
	enCodec  *codec.AucodecInst // 输出数据处理
	downCtrl string             // 数据下行方式控制：

	mngr    *Manager // 归属实例管理器指针;
	context string   // 请求上下文句柄,request.handle;

	// 事件判定; 适配层进行wrapper.so c接口层调用;
	usrActions map[UserEvent]UsrAct // 用户注册事件;

	readTvCnt    int32           // 读取缓冲区失败次数;
	inputData    buffer.MultiBuf // 输入数据缓冲区;
	outPutDatas  buffer.MultiBuf // 输出数据缓冲区;
	outPutLastId uint            // 输出数据id;

	instWg sync.WaitGroup

	// 协程异常状态同步
	excpNum  int
	excpErr  error
	excpLock sync.RWMutex
	respChan chan ActMsg

	// APM log
	spanMeta string		// 非实时服务离线span
	nrtId	string		// 非实时服务任务id
	nrtTime int64		// 非实时服务任务开始时间

	// 服务属性参数
	uid   string
	appid string
	sub   string

	upParams  map[string]string // 上行请求参数对
	upDatas   []eventStorage    // 上行请求数据流
	downDatas []eventStorage    // 请求下行数据流
}

func (si *ServiceInst) Init(acts map[UserEvent]UsrAct, tool *xsf.ToolBox) (errInfo error) {
	// 数据缓冲区
	si.inputData.Init(seqBufTimeout, seqAudSize, 0)
	si.outPutDatas.Init(0, seqRltSize, 0)
	si.tool = tool
	si.usrActions = acts
	si.alive = true
	return
}

func (si *ServiceInst) Fini() {
	si.setAlive(false)
	si.instWg.Wait()
	si.release()
	si.outPutDatas.Fini()
	si.inputData.Fini()
	si.tool = nil
	si.usrActions = nil
	return
}

// 申请服务计算资源;
func (si *ServiceInst) ResLoad(base int, params map[string]string, span* xsf.Span) (errNum int, errInfo error) {
	si.release()
	si.upParams = params
	si.spanMeta = span.Meta()
	si.nrtId = params[nrtTask]
	if len(si.nrtId) > 0 {
		status, err := nrtQuery(si.nrtId)
		if err != nil{
			si.tool.Log.Errorw("service query nrt task fail", "task", si.nrtId, "err", err.Error())
		}else if status != nrtCreate {
			si.tool.Log.Errorw("service query nrt task invalid", "task", si.nrtId, "status", status)
			return frame.AigesErrorInvalidParaValue, frame.ErrorInvalidParaValue
		}
		si.nrtTime = time.Now().UnixNano()
		if err = nrtUpdate(si.nrtId, 0, nil, nrtActive); err != nil {
			si.tool.Log.Errorw("service update nrt task active fail", "task", si.nrtId)
			// 数据库状态失败降级
		}
	}
	// 音频重采样
	si.dpar = make(map[string]*dp.AudioResampler)
	// 音频编解码, 上行数据解码操作lazy init
	si.deCodec = make(map[string]*codec.AucodecInst)
	if si.enCodec, errNum, errInfo = codec.NewAucodec(&params); errInfo != nil {
		return
	}

	// 服务属性参数
	si.downCtrl = params[downMethod]
	if si.downCtrl == downAsync && len(conf.RabbitHost) == 0 {
		return frame.AigesErrorInvalidPara, errors.New("invalid downMethod")
	}
	si.uid = params[userId]
	si.appid = params[appId]
	si.inputData.SetBase(uint(base))
	si.instHdl = params[sessionId]
	si.setAlive(true)
	si.instWg.Add(1)
	go si.asyncCalc(params)
	return
}

// 释放会话临时资源;
func (si *ServiceInst) release() {
	if si.enCodec != nil {
		if code, err := codec.CloseAucodec(si.enCodec); err != nil {
			si.tool.Log.Errorw("service inst encoder end fail", "code", code, "err", err.Error(), "inst", si.instHdl)
		}
		si.enCodec = nil
	}

	for _,v := range si.deCodec {
		if code, err := codec.CloseAucodec(v); err != nil {
			si.tool.Log.Errorw("service inst decoder end fail", "code", code, "err", err.Error(), "inst", si.instHdl)
		}
	}
	si.deCodec = nil

	for _, v := range si.dpar {
		if err := v.Destroy(); err != nil {
			si.tool.Log.Errorw("service inst resample destroy fail", "err", err.Error(), "inst", si.instHdl)
		}
	}
	si.dpar = nil
	si.outPutDatas.Release()
	si.inputData.Release()
	si.instHdl = ""
	si.wrapperHdl = ""
	si.spanMeta = ""
	si.nrtId = ""
	si.readTvCnt = 0
	si.outPutLastId = 0
	si.nrtTime = 0
	si.upParams = nil
	si.upDatas = nil
	si.downDatas = nil
	si.respChan = nil
	si.setAlive(false)
}

/*
同步写入：写入输入数据缓冲区;
异步协程：数据读取 -> 重采样 -> 编解码(解码) -> 个性化 -> AI计算 -> 编解码(编码) -> 数据处理 -> 输出
*/
func (si *ServiceInst) DataReq(dataIn *protocol.EngInputData, span *xsf.Span) (errNum int, errInfo error) {

	// 会话异常检查
	errNum, errInfo = si.checkExcp()
	if errInfo != nil {
		si.tool.Log.Errorw("DataReq Check Goroutine Error", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
		span.WithTag("errNum", strconv.Itoa(errNum)).WithTag("errInfo", errInfo.Error()).WithTag("location", "DataReq::checkExcp")
		return
	}

	if si.sessAlive() {
		// 状态判定, 写数据排序缓冲区;
		inputs := make([]buffer.DataMeta, 0, len(dataIn.DataList))
		for _, v := range dataIn.DataList {
			var input interface{} = nil // Note: 直接传递v.Data时, interface not nil
			dataLen := len(v.Data)
			if dataLen > 1 { // 协议限定: 适配听写尾音频大小为1场景请求,若尾音频大小为1,写入空数据
				input = v.Data
			}

			inputs = append(inputs, buffer.DataMeta{input, v.DataId, uint(v.FrameId),
				buffer.DataStatus(v.Status), buffer.DataType(v.DataType),
				v.Format, v.Encoding, v.Desc})
		}
		errNum, errInfo = si.inputData.WriteData(inputs)
	}

	return
}

/*
查询获取计算结果：结果数据
1. 支持超时阻塞查询 & 非阻塞查询;
timeout: 0 非阻塞； other 超时ms;
*/
func (si *ServiceInst) DataResp(timeout uint, span *xsf.Span) (dataOut protocol.EngOutputData, errNum int, errInfo error) {
	// 会话异常检查;
	errNum, errInfo = si.checkExcp()
	if errInfo != nil {
		si.tool.Log.Errorw("DataResp Check Goroutine Error", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
		span.WithTag("errNum", strconv.Itoa(errNum)).WithTag("errInfo", errInfo.Error()).WithTag("location", "DataResp::checkExcp")
		return
	}

	// 查询输出结果;
	if si.sessAlive() && si.downCtrl != downAsync {
		// 读取响应结果;
		var respRlt []byte
		var status protocol.EngOutputData_DataStatus
		engResult, _, err := si.outPutDatas.ReadDataWithTime(timeout)
		if err != nil && err != frame.ErrorSeqBufferEmpty {
			errNum, errInfo = frame.AigesErrorEngInactive, err // 适配错误码： EngInactive error code
		} else if len(engResult) > 0 {
			for _, rslt := range engResult {
				if rslt.Data != nil {
					switch rslt.DataType {
					case buffer.DataText:
						respRlt = []byte(rslt.Data.(string))
					case buffer.DataAudio, buffer.DataImage, buffer.DataVideo:
						respRlt = rslt.Data.([]byte)
					}
				}
				result := protocol.MetaData{rslt.DataId, uint32(rslt.FrameId),
					protocol.MetaData_DataType(rslt.DataType), protocol.MetaData_DataStatus(rslt.Status),
					rslt.Format, rslt.Encoding, respRlt, rslt.Desc}
				dataOut.DataList = append(dataOut.DataList, &result)
				status = protocol.EngOutputData_DataStatus(rslt.Status)
				// TODO status 需考虑多数据流的不同stream数据状态,所有stream end即会话end;
			}
		} else {
			// 超时状态查询错误码,判定异步是否异常
			errNum, errInfo = si.checkExcp()
		}

		// 构造响应包;
		dataOut.Status = status
		dataOut.Ret = int32(errNum)
		if errInfo != nil {
			dataOut.Err = errInfo.Error()
		} else if dataOut.DataList == nil && !conf.SessMode {
			// 非会话模式读取超时场景
			errInfo, errNum = frame.ErrorOnceExecTimeout, frame.AigesErrorOnceTimeout
		}
	}
	return
}

/*
异常中断
*/
func (si *ServiceInst) DataException(span *xsf.Span) (dataOut protocol.EngOutputData, errNum int, errInfo error) {
	// 重置服务实例
	if si.setAlive(false) {
		si.setExcp(frame.AigesErrorFinRoutine, frame.ErrorFinishRoutine)
	}
	return
}

func (si *ServiceInst) dataDpLazy(meta *buffer.DataMeta) (errNum int, errInfo error) {
	switch meta.DataType {
	case buffer.DataAudio:
		// 上行多数据流音频重采样
		if si.dpar[meta.DataId] == nil {
			si.dpar[meta.DataId] = dp.NewResampler()
			var rateStr string
			if len(meta.Format) > 0 {
				fmts := strings.Split(meta.Format, ";")
				if len(fmts) == 2 {
					pairs := strings.Split(fmts[1], "=")
					if len(pairs) == 2 {
						rateStr = pairs[1]
					}
				}
			} else {
				rateStr = si.upParams[audioRate]
			}
			if len(rateStr) > 0 {
				rate, _ := strconv.Atoi(rateStr)
				errInfo = checkAudioRate(rate)
				if errInfo != nil {
					errNum = frame.AigesErrorInvalidPara
					si.tool.Log.Errorw("service inst get invalid audio format", "format", meta.Format, "rate", rateStr, "inst", si.instHdl)
					return
				}
				if rate == AudioRate8k /*&& reSampler able*/ {
					errInfo = si.dpar[meta.DataId].Init(1, AudioRate8k, AudioRate16k, ResampleQuality)
					if errInfo != nil {
						errNum = frame.AigesErrorDpInit
						return
					}
				}
				si.upParams[audioRate] = rateStr // 同步刷新rate参数至upParams, codec初始化所需
			}
		}
		// 上行多数据流音频编解码
		if si.deCodec[meta.DataId] == nil {
			bak := si.upParams[frame.ParaAuCodec]
			if len(meta.Encoding) > 0 {
				si.upParams[frame.ParaAuCodec] = meta.Encoding
			}
			if auInst, errNum, errInfo := codec.NewAucodec(&si.upParams); errInfo != nil {
				return errNum, errInfo
			}else {
				si.deCodec[meta.DataId] = auInst
			}
			si.upParams[frame.ParaAuCodec] = bak
		}
	}
	return
}

//	异步协程: 单线程消费计算;
func (si *ServiceInst) asyncCalc(params map[string]string) {
	defer catch.RecoverHandle(si.instHdl)
	var errNum int
	var errInfo error
	// TODO 参数过滤控制;
	if conf.SessMode {
		errNum, errInfo = si.sessionCalc(params)
	} else {
		// 非会话模式(eg:translate场景)
		errNum, errInfo = si.noneSessCalc(params)
	}

	if errInfo != nil {
		si.setExcp(errNum, errInfo)
		si.outPutDatas.Signal() // 中断结果查询超时阻塞操作
	}
	if si.downCtrl == downAsync {
		if errInfo != nil {
			si.dataPostBack(nil, errNum, errInfo)
			si.tool.Log.Errorw("request fail with downMethod async", "code", errNum, "err", errInfo.Error())
		}
		// 实例释放：
		si.DataException(nil)
		si.mngr.Release(si.context)
		// TODO FreeChan((si.wrapperHdl)) TODO !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	}
	si.instWg.Done()
	return
}

// 会话模式异步计算
func (si *ServiceInst) sessionCalc(params map[string]string) (errNum int, errInfo error) {
	// 用户事件判定：资源申请事件(eg:ssb场景)
	actNew, ok := si.usrActions[EventNew] // 上层接口保障非nil
	if ok {
		var resp ActMsg
		var ids []int

		resp, errNum, errInfo = actNew(si.instHdl, &ActMsg{Params: params, PersonIds: ids})
		if errInfo != nil {
			si.tool.Log.Errorw("EventNew action fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
			errNum = frame.WrapperCreateErr
			return
		}
		si.tool.Log.Debugw("EventNew action detail", "inst", si.instHdl, "params", params)
		si.wrapperHdl = resp.WrapperHdl
		if conf.WrapperAsync {
			si.respChan, errInfo = AllocChan(si.instHdl)
			if errInfo != nil {
				errNum = frame.WrapperCreateErr
				return
			}
		}
	}

	fin := false
	si.upDatas = make([]eventStorage, 0, 1)
	si.downDatas = make([]eventStorage, 0, 1)
	for si.sessAlive() && !fin {
		var merge []buffer.DataMeta
		merge, errNum, errInfo = si.inputData.ReadMergeData()
		if errInfo != nil && errInfo != frame.ErrorSeqBufferEmpty {
			si.tool.Log.Errorw("read audio buffer fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
			break
		} else if merge == nil || len(merge) == 0 {
			if si.readTvCnt >= sessTimeoutCnt {
				// TODO 写入last输入数据;
				errNum, errInfo = frame.AigesErrorSessTimeout, frame.ErrorInstRwTimeout
				si.tool.Log.Errorw("read input buf timeout then exit", "sid", si.instHdl)
				break
			}
			atomic.AddInt32(&si.readTvCnt, 1)
			continue
		}

		// storage event log & media data
		for _, v := range merge {
			si.upDatas = append(si.upDatas, eventStorage{v.Data, v.DataId,
				v.DataType, v.Format, v.Encoding})
			si.tool.Log.Infow("calc read buffer", "stream", v.DataId,
				"frame", v.FrameId, "status", v.Status, "format", v.Format, "encoding", v.Encoding, "sid", si.instHdl)
		}

		errNum, errInfo = si.dataCompletion(&merge)
		if errInfo != nil {
			break
		}
		errNum, errInfo = si.checkExcp()
		if errInfo != nil {
			si.tool.Log.Errorw("sessionCalc checkExcp then exit", "sid", si.instHdl)
			break
		}
		fin, errNum, errInfo = si.stask(merge, &si.downDatas)
		if errInfo != nil {
			break
		}
		si.readTvCnt = 0 // 重置读取超时次数;
	}

	// 用户事件判定：资源释放事件(eg:sse场景)
	actDel, ok := si.usrActions[EventDel]
	if ok {
		_, _, _ = actDel(si.instHdl, &ActMsg{WrapperHdl: si.wrapperHdl})
		si.tool.Log.Debugw("EventDel action detail", "inst", si.instHdl, "wrapperHdl", si.wrapperHdl)
	}

	return
}

// 非会话模式异步计算
func (si *ServiceInst) noneSessCalc(params map[string]string) (errNum int, errInfo error) {
	// 读取数据 & 校验状态once
	output, errNum, errInfo := si.inputData.ReadMergeData()
	if errInfo != nil && errInfo != frame.ErrorSeqBufferEmpty {
		si.tool.Log.Errorw("noneSessCalc read data fail", "errNum", errNum, "errInfo", errInfo.Error())
		return
	} else if output == nil || len(output) == 0 {
		errNum, errInfo = frame.AigesErrorInvalidData, frame.ErrorInvalidData
		si.tool.Log.Errorw("noneSessCalc read invalid data, nil/empty data")
		return
	}
	// eventLog storage
	si.upDatas = make([]eventStorage, 0, len(output))
	for _, stream := range output {
		si.upDatas = append(si.upDatas, eventStorage{stream.Data, stream.DataId,
			stream.DataType, stream.Format, stream.Encoding})
	}

	// 离线类引擎不做数据存储,于eventLog后数据补齐
	errNum, errInfo = si.dataCompletion(&output)
	if errInfo != nil {
		return
	}
	engData := make([]DataMeta, 0, len(output))
	for _, stream := range output {
		// 输入数据处理
		if errNum,errInfo = si.dataDpLazy(&stream); errInfo != nil {
			return
		}
		var finData []byte
		switch stream.DataType {
		case buffer.DataAudio:
			codecTmp, errNum, errInfo := si.deCodec[stream.DataId].Decode(stream.Data.([]byte), true) // TODO 判定interface类型
			if errInfo != nil {
				si.tool.Log.Errorw("codec::Decode fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
				return errNum, errInfo
			}
			finData, errInfo = si.dpar[stream.DataId].ProcessInt(0, codecTmp)
			if errInfo != nil {
				si.tool.Log.Errorw("dpar::Resample fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
				return errNum, errInfo
			}
		default:
			finData = stream.Data.([]byte)
		}

		engData = append(engData, DataMeta{stream.DataId, finData, int(stream.DataType),
			int(stream.Status), stream.FrameId, stream.Format, "raw", "TODO "})
		// TODO 传递数据编解码依据功能设置value, 数据描述传递
	}

	// 事件判定&once处理
	if conf.WrapperAsync {
		si.respChan, errInfo = AllocChan(si.instHdl)
		if errInfo != nil {
			errNum = frame.WrapperCreateErr
			return
		}
		// TODO defer FreeChan(si.instHdl)！！！！！！！！！！
	}
	actExec, ok := si.usrActions[EventOnceExec]
	if ok {
		// TODO actExec: exec & release
		var resp ActMsg
		resp, errNum, errInfo = actExec(si.instHdl, &ActMsg{Params: params, DeliverData: engData})
		if errInfo != nil {
			si.tool.Log.Errorw("OnceExec action fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
			errNum = frame.WrapperExecErr
			return
		}

		if conf.WrapperAsync {
			resp, _ = <-si.respChan
			if resp.AsyncErr != nil {
				si.tool.Log.Errorw("OnceExec async callback error", "errInfo", resp.AsyncErr.Error(), "sid", si.instHdl)
			}
		}
		//  输出数据处理
		si.wrapperHdl = resp.WrapperHdl
		si.downDatas = make([]eventStorage, 0, len(resp.DeliverData))
		respData := make([]buffer.DataMeta, 0, len(resp.DeliverData))
		for _, data := range resp.DeliverData {
			var tmpData interface{}
			var respEncode []byte = data.Data
			switch buffer.DataType(data.DataType) {
			case buffer.DataText:
				tmpData = string(data.Data)
			case buffer.DataAudio:
				respEncode, errNum, errInfo = si.enCodec.Encode(data.Data, true) // TODO check data.Data if nil
				if errInfo != nil {
					si.tool.Log.Errorw("codec::Encode fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
					return
				}
				tmpData = respEncode
			default:
				tmpData = data.Data
			}
			respData = append(respData, buffer.DataMeta{tmpData, data.DataId, data.DataFrame,
				buffer.DataStatus(data.DataStatus), buffer.DataType(data.DataType),
				data.DataFmt, data.DataEnc, nil /*data desc*/})
			si.downDatas = append(si.downDatas, eventStorage{respEncode, data.DataId,
				buffer.DataType(data.DataType), data.DataFmt, data.DataEnc})
		}
		if si.downCtrl == downAsync {
			si.dataPostBack(&respData, 0, nil)
		} else {
			si.outPutDatas.WriteData(respData)
		}
	} else {
		si.tool.Log.Errorw("noneSessCalc can't find OnceExec Event to call")
		errNum, errInfo = frame.AigesErrorNilEvent, frame.ErrorNilRegEvent
	}
	return
}

func (si *ServiceInst) dataCompletion(bufData *[]buffer.DataMeta) (errNum int, errInfo error) {
	// check if data need to download from http/s3.
	for k, _ := range *bufData {
		ds, exist := (*bufData)[k].Desc[dataSrc]
		if exist {
			switch string(ds) {
			case dataHttp:
				url, _ := (*bufData)[k].Desc[dataHttpUrl]
				if len(url) == 0 {
					return frame.AigesErrorInvalidData, errors.New("input invalid http url")
				}

				// download from http, return err if download fail
				data, errInfo := storage.HttpDownload(string(url))
				if errInfo != nil {
					return frame.AigesErrorInvalidData, errInfo
				}
				(*bufData)[k].Data = data
			case dataS3:
				access, _ := (*bufData)[k].Desc[dataS3Access]
				secret, _ := (*bufData)[k].Desc[dataS3Secret]
				endpoint, _ := (*bufData)[k].Desc[dataS3Ep]
				bucket, _ := (*bufData)[k].Desc[dataS3Bucket]
				key, _ := (*bufData)[k].Desc[dataS3Key]
				if len(access) == 0 || len(secret) == 0 || len(endpoint) == 0 ||
					len(bucket) == 0 || len(key) == 0 {
					return frame.AigesErrorInvalidData, errors.New("input invalid s3 tags")
				}

				// TODO download from s3, return err if download fail
			case dataClient:
				// nothing to do
			default:
				return frame.AigesErrorInvalidData, errors.New("input invalid data source")
			}
		}
	}

	return
}

// push downstream data to s3, produce rmq message
func (si *ServiceInst) dataPostBack(data *[]buffer.DataMeta, code int, err error) {
	var dataDown protocol.EngOutputData
	var status protocol.EngOutputData_DataStatus
	if data != nil {
		for _, rslt := range *data {
			var respRlt []byte
			if rslt.Data != nil {
				switch rslt.DataType {
				case buffer.DataText:
					respRlt = []byte(rslt.Data.(string))
				case buffer.DataAudio, buffer.DataImage, buffer.DataVideo:
					respRlt = rslt.Data.([]byte)
					// TODO 异步回传富媒体分离方案
				}
			}
			result := protocol.MetaData{rslt.DataId, uint32(rslt.FrameId),
				protocol.MetaData_DataType(rslt.DataType), protocol.MetaData_DataStatus(rslt.Status),
				rslt.Format, rslt.Encoding, respRlt, rslt.Desc}
			dataDown.DataList = append(dataDown.DataList, &result)
			status = protocol.EngOutputData_DataStatus(rslt.Status)

		}
	}
	dataDown.Ret = int32(code)
	if err != nil {
		dataDown.Err = err.Error()
	}
	dataDown.Status = status

	var jd jsonDown
	output, errMsl := proto.Marshal(&dataDown)
	jd.Sid = si.instHdl
	jd.ReqParam = si.upParams
	jd.SpanMeta = si.spanMeta
	jd.Data = output
	jd.Ret = 0
	if errMsl != nil {
		si.tool.Log.Errorw("dataPostBack proto marshal fail", "sid", si.instHdl, "data", string(output), "err", errMsl.Error())
		jd.Ret = frame.AigesErrorPbMarshal
		jd.ErrDesc = frame.ErrorPbMarshal.Error()
	}
	jdata, merr := json.Marshal(jd)
	if merr != nil {
		si.tool.Log.Errorw("dataPostBack json marshal fail", "sid", si.instHdl, "jsondata", jd, "err", merr.Error())
	}

	// 更新数据库状态
	nrtCost := (time.Now().UnixNano() - si.nrtTime) / (1000 * 1000)
	nrtStatus := nrtSuccess
	if err != nil {
		nrtStatus = nrtExecFail
	}
	if err := nrtUpdate(si.nrtId, nrtCost, output, nrtStatus); err != nil {
		si.tool.Log.Errorw("nrt status update fail", "hdl", si.instHdl, "err", err.Error())
	}

	// rabbitMQ 队列存储
	perr := storage.RabPuslish(jdata)
	if perr != nil {
		si.tool.Log.Errorw("dataPostBack produce rmq message fail", "sid", si.instHdl, "err", perr.Error())
	}
	si.tool.Log.Debugw("dataPostBack produce rmq message", "sid", si.instHdl, "body", jdata)
	return
}

func (si *ServiceInst) stask(inputs []buffer.DataMeta, downstream *[]eventStorage) (fin bool, errNum int, errInfo error) {
	lastInput := false
	fin = false
	deliver := make([]DataMeta, 0, len(inputs))
	for _, v := range inputs {
		// 非空数据进行编解码&重采样操作;
		meta := DataMeta{v.DataId, nil, int(v.DataType), int(v.Status),
			v.FrameId, v.Format, v.Encoding, "" /*TODO temp nil for now*/}
		if v.Data != nil {
			if errNum,errInfo = si.dataDpLazy(&v); errInfo != nil {
				return
			}
			switch v.DataType {
			case buffer.DataAudio:
				codecTmp, errNum, errInfo := si.deCodec[v.DataId].Decode(v.Data.([]byte), buffer.DataStatus(v.Status) == buffer.DataStatusLast) // TODO check typeof v.Data
				if errInfo != nil {
					si.tool.Log.Errorw("codec::Decode fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
					return true, errNum, errInfo
				}

				finData, errInfo := si.dpar[v.DataId].ProcessInt(0, codecTmp)
				if errInfo != nil {
					si.tool.Log.Errorw("dpar::Resample fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
					return true, errNum, errInfo
				}
				meta.Data = finData
			default:
				meta.Data = v.Data.([]byte)
			}
		}
		deliver = append(deliver, meta)
		lastInput = (v.Status == buffer.DataStatusLast || v.Status == buffer.DataStatusOnce)
		// TODO lastInput 判定条件需考虑多数据流的不同stream数据状态,所有stream end即会话end;
	}

	// 数据写事件;
	actWrite, ok := si.usrActions[EventWrite]
	if !ok {
		return true, frame.WrapperReadErr, errors.New("there is no register EventWrite by Register()")
	}

	si.tool.Log.Debugw("EventWrite action detail", "inst", si.instHdl, "wrapperHdl", si.wrapperHdl, "input", deliver)
	_, errNum, errInfo = actWrite(si.instHdl, &ActMsg{WrapperHdl: si.wrapperHdl, DeliverData: deliver})
	if errInfo != nil {
		si.tool.Log.Errorw("usrAction EventWrite fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
		errNum = frame.WrapperWriteErr
		return true, errNum, errInfo
	}

	// 服务实时读 || 音频last
	if conf.RealTimeRead || lastInput {
		waitCnt := 0
		for {
			// 数据读事件; 若写入last数据则循环读取;
			var resp ActMsg
			if conf.WrapperAsync {
				// 读取异步会话回调模式回写channel数据;
				var alive bool
				if !lastInput { // 中间结果读取,非阻塞读;
					select {
					case resp, alive = <-si.respChan:
						if !alive {
							fin = true // channel关闭,finish
						}
					default:
					}
				} else {
					resp, alive = <-si.respChan
					if !alive {
						fin = true
					}
				}
				if resp.AsyncErr != nil {
					si.tool.Log.Errorw("receive async callback error", "errInfo", resp.AsyncErr.Error(), "sid", si.instHdl)
					return true, frame.WrapperAsyncErr, resp.AsyncErr
				}
			} else {
				actRead, ok := si.usrActions[EventRead]
				if !ok {
					return true, frame.WrapperReadErr, errors.New("there is no register EventRead by Register()")
				}
				resp, errNum, errInfo = actRead(si.instHdl, &ActMsg{WrapperHdl: si.wrapperHdl})
				if errInfo != nil {
					si.tool.Log.Errorw("usrAction EventRead fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
					errNum = frame.WrapperReadErr
					return true, errNum, errInfo
				}
				si.tool.Log.Debugw("EventRead action detail", "inst", si.instHdl, "wrapperHdl", si.wrapperHdl, "output", resp)
			}

			outputs := make([]buffer.DataMeta, 0, len(resp.DeliverData))
			for _, v := range resp.DeliverData {
				// TODO 考虑多个输出数据流的数据id;
				meta := buffer.DataMeta{nil, v.DataId, si.outPutLastId, buffer.DataStatus(v.DataStatus),
					buffer.DataType(v.DataType), v.DataFmt, v.DataEnc, nil /*TODO nil for now*/}
				si.outPutLastId++
				// 输出数据处理(编码处理)
				if v.Data != nil {
					waitCnt = 0 // 重置超时cnt;
					var respData interface{}
					var respEncode []byte = v.Data
					switch buffer.DataType(v.DataType) {
					case buffer.DataText:
						respData = string(respEncode)
					case buffer.DataAudio:
						respEncode, errNum, errInfo = si.enCodec.Encode(v.Data, meta.Status == buffer.DataStatusLast)
						if errInfo != nil {
							si.tool.Log.Errorw("codec::Encode fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
							return true, errNum, errInfo
						}
					default:
						respData = respEncode
					}
					meta.Data = respData
					// event downstream
					*downstream = append(*downstream, eventStorage{respEncode, v.DataId,
						buffer.DataType(v.DataType), v.DataFmt, v.DataEnc})
				}
				outputs = append(outputs, meta)
				// 读取结果完毕,中断循环&置结束位
				if buffer.DataStatus(v.DataStatus) == buffer.DataStatusOnce || buffer.DataStatus(v.DataStatus) == buffer.DataStatusLast {
					fin = true // TODO fin 判定条件需考虑多数据流的不同stream数据状态,所有stream end即会话end;
				}
			}

			if si.downCtrl == downAsync {
				si.dataPostBack(&outputs, 0, nil)
			} else {
				si.outPutDatas.WriteData(outputs) // TODO 合成场景确认：fmt & enc
			}

			if fin || !lastInput || waitCnt > syncRltWaitCnt {
				break
			}
			// 若同步读未取得计算结果,循环等待timeWait(50ms);
			if len(resp.DeliverData) == 0 { // 等价于 || resp.DeliverData == nil
				time.Sleep(time.Duration(syncRltWaitTime) * time.Millisecond)
				waitCnt++
			}
		}
	}
	return
}

// 查询当前实例状态;
func (si *ServiceInst) sessAlive() (alive bool) {
	si.aliveLock.RLock()
	defer si.aliveLock.RUnlock()
	return si.alive
}

// 设置实例状态,返回设置前状态;
func (si *ServiceInst) setAlive(alive bool) (last bool) {
	si.aliveLock.Lock()
	defer si.aliveLock.Unlock()
	last = si.alive
	si.alive = alive
	return
}

// 会话异常同步
func (si *ServiceInst) setExcp(errNum int, errInfo error) {
	si.excpLock.Lock()
	defer si.excpLock.Unlock()
	// 保留首个错误信息;
	if si.excpErr == nil {
		si.excpNum = errNum
		si.excpErr = errInfo
	}
	return
}

func (si *ServiceInst) checkExcp() (errNum int, errInfo error) {
	si.excpLock.RLock()
	defer si.excpLock.RUnlock()

	errNum = si.excpNum
	errInfo = si.excpErr
	return
}

func (si *ServiceInst) resetExcp() {
	si.excpLock.Lock()
	defer si.excpLock.Unlock()

	si.excpNum = frame.AigesSuccess
	si.excpErr = nil
	return
}

func (si *ServiceInst) CatchTag() (params map[string]string, datas []eventStorage) {
	return si.upParams, si.upDatas
}