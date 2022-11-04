package instance

import "C"
import (
	"errors"
	"github.com/xfyun/aiges/buffer"
	"github.com/xfyun/aiges/catch"
	"github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/frame"
	"github.com/xfyun/aiges/protocol"
	aigesUtils "github.com/xfyun/aiges/utils"
	"github.com/xfyun/xsf/server"
	"github.com/xfyun/xsf/utils"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
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

type ServiceInst struct {
	instHdl string // 服务实例唯一标识,sid;
	//wrapperHdl interface{} // wrapper层实例句柄
	wrapperHdl unsafe.Pointer
	alive      bool // 服务实例状态flag
	aliveLock  sync.RWMutex

	tool *xsf.ToolBox
	//dpar     map[string] /*dataId*/ *dp.AudioResampler // 数据处理(重采样)
	upStatus map[string]buffer.DataStatus // 上行数据流状态,用于判定最终请求状态
	downCtrl string                       // 数据下行方式控制：

	instWg  sync.WaitGroup
	mngr    *Manager // 归属实例管理器指针;
	context string   // 请求上下文句柄,request.handle;

	// 事件判定; 适配层进行wrapper.so c接口层调用;
	usrActions map[UserEvent]UsrAct // 用户注册事件;

	// 输入输出缓冲区
	inputData    buffer.DataBuf  // 输入数据缓冲区;
	outPutData   buffer.DataBuf  // 输出数据缓冲区;
	outPutId     uint            // 输出业务包id;
	outPutLastId map[string]uint // 输出数据流id;

	// 协程异常状态同步
	excpNum  int
	excpErr  error
	excpLock sync.RWMutex
	respChan chan ActMsg
	onceTrig chan bool

	// apm log
	spanMeta     string    // 非实时服务离线span
	debug        *xsf.Span // 用于调试的span日志
	debugWrapper *xsf.Span // 用于插件的span日志
	nrtTime      int64     // 非实时服务任务开始时间
	nrtDegrade   bool      // 非实时任务状态更新降级

	// 框架通用属性参数
	uid       string
	appid     string
	cloudId   string
	composeId string
	sub       string
	svcId     string

	// 排障及数据分析
	upDatas   []eventStorage
	downDatas []eventStorage

	// 计量数据
	meterParam string // header: "meterPara"
	audioLen   int    // header："audioLen.xxx"
	bytes      int    // header: "bytes.xxx"
	mpCtrl     bool
	// =============================================
	// v2 protocol
	headers  map[string]string
	params   map[string]string
	expect   map[string]*protocol.MetaDesc
	sessType protocol.LoaderInput_SessState
	//encoder  map[string] /*dataId*/ codec.Instance
	//decoder  map[string] /*dataId*/ codec.Instance
	//outdpar  map[string] /*dataId*/ *dp.AudioResampler // 数据处理(重采样)

}

func (si *ServiceInst) Init(acts map[UserEvent]UsrAct, tool *xsf.ToolBox) (errInfo error) {
	// 数据缓冲区
	si.inputData.Init(seqBufTimeout, seqAudSize, 0)
	si.outPutData.Init(0, seqRltSize, 0)
	si.onceTrig = make(chan bool, 1)
	si.tool = tool
	si.usrActions = acts
	si.alive = true
	return
}

func (si *ServiceInst) Fini() {
	si.setAlive(false)
	si.instWg.Wait()
	si.release()
	si.outPutData.Fini()
	si.inputData.Fini()
	si.tool = nil
	si.usrActions = nil
	return
}

// 申请服务计算资源;
func (si *ServiceInst) ResLoad(proto *map[string]string, input *protocol.LoaderInput, span *xsf.Span) (errNum int, errInfo error) {
	si.release()
	si.spanMeta = span.Meta()
	si.debug = span.Next(utils.SrvSpan).Start().WithName("AILoader").WithServerAddr()
	si.debugWrapper = si.debug.Next(utils.SrvSpan).Start().WithName("AIWrapper").WithServerAddr()
	si.headers = input.GetHeaders()
	si.params = input.GetParams()
	for _, pv := range conf.HeaderPass {
		if si.params == nil {
			si.params = make(map[string]string)
		}
		if _, ep := si.params[pv]; ep == false {
			si.params[pv] = si.headers[pv]
		}
	}
	si.sessType = input.GetState()
	si.expect = make(map[string]*protocol.MetaDesc)
	for _, v := range input.GetExpect() {
		si.expect[v.Name] = v
	}
	// 上行数据缓冲区设置
	si.inputData.SetBase(uint(input.SyncId))

	si.outPutLastId = make(map[string]uint)
	si.outPutId = 0

	// 平台控制参数
	si.uid = si.headers[protocol.Uid]
	si.appid = si.headers[protocol.AppId]
	si.cloudId = si.headers[protocol.CloudId]
	si.composeId = si.headers[protocol.ComposeId]
	si.sub = func() string {
		if val, ok := si.headers[protocol.Sub]; ok {
			return val
		} else {
			return "ase"
		}
	}()
	si.svcId = input.ServiceId
	si.meterParam = si.headers[protocol.MeterPara]
	si.instHdl = si.headers[protocol.SessionId]
	si.downCtrl = si.headers[downMethod]
	if si.downCtrl == downAsync && len(conf.RabbitHost) == 0 {
		return frame.AigesErrorInvalidPara, errors.New("invalid downMethod")
	}
	si.upStatus = make(map[string]buffer.DataStatus)
	si.setAlive(true)
	si.instWg.Add(1)
	go si.asyncCalc()
	return
}

// 释放会话临时资源;
func (si *ServiceInst) release() {
	// clear onceTrig
	select {
	case <-si.onceTrig:
	default:
	}
	si.outPutData.Release()
	si.inputData.Release()
	si.instHdl = ""
	si.wrapperHdl = nil
	si.spanMeta = ""
	si.outPutId = 0
	si.outPutLastId = nil
	si.nrtTime = 0
	si.audioLen = 0
	si.bytes = 0
	si.meterParam = ""
	si.nrtDegrade = false
	si.upDatas = nil
	si.downDatas = nil
	si.respChan = nil
	si.upStatus = nil
	si.expect = nil
	si.setAlive(false)
}

/*
数据写入：同步写入,异步处理
*/
func (si *ServiceInst) DataReq(input *protocol.LoaderInput, span *xsf.Span) (errNum int, errInfo error) {
	// 会话异常检查
	errNum, errInfo = si.checkExcp()
	if errInfo != nil {
		si.tool.Log.Errorw("DataReq Check Goroutine Error", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
		span.WithTag("errNum", strconv.Itoa(errNum)).WithTag("errInfo", errInfo.Error()).WithTag("location", "DataReq::checkExcp")
		return
	}

	if si.sessAlive() {
		// 状态判定, 写数据排序缓冲区;
		inputs := make([]buffer.DataMeta, 0, len(input.Pl))
		for _, v := range input.Pl {
			attr := protocol.GetBaseAttr(v.Meta)
			meta := buffer.DataMeta{v.Data, attr.Name, uint(attr.Seq),
				buffer.DataStatus(attr.Status), buffer.DataType(attr.Type), v.Meta}
			inputs = append(inputs, meta)
			si.tool.Log.Debugw("input data segment", "data key", attr.Name,
				"data status", attr.Status,
				"data type", attr.Type,
				"data length", len(v.Data), "sid", si.instHdl)
		}
		errNum, errInfo = si.inputData.WriteData(uint(input.SyncId), inputs)
		if input.State == protocol.LoaderInput_ONCE {
			select {
			case si.onceTrig <- true:
			default:
			}
		}
	}

	return
}

/*
查询获取计算结果：结果数据
1. 支持超时阻塞查询 & 非阻塞查询;
timeout: 0 非阻塞； other 超时ms;
*/
func (si *ServiceInst) DataResp(timeout uint, span *xsf.Span) (dataOut protocol.LoaderOutput, errNum int, errInfo error) {
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
		var status protocol.LoaderOutput_RespStatus
		engResult, _, err := si.outPutData.ReadDataWithTime(&timeout)
		if err != nil && err != frame.ErrorSeqBufferEmpty {
			if errNum, errInfo = si.checkExcp(); errInfo == nil { // buffer关闭,先检查是否存在错误,无错误则返回10101;
				errNum, errInfo = frame.AigesErrorEngInactive, err // 适配错误码： EngInactive error code
			}
		} else if len(engResult) > 0 {
			for _, rslt := range engResult {
				var respRlt []byte
				if rslt.Data != nil {
					respRlt = rslt.Data.([]byte)
				}
				result := protocol.Payload{Meta: rslt.Desc, Data: respRlt}
				si.tool.Log.Debugw("read from data cache :", "dataLen", len(respRlt),
					"dataId", rslt.DataId, "dataType", rslt.DataType, "status", rslt.Status, "desc", rslt.Desc,
					"frameId", rslt.FrameId, "sid", si.instHdl)
				dataOut.Pl = append(dataOut.Pl, &result)
				status = protocol.LoaderOutput_RespStatus(rslt.Status)
				// TODO status 需考虑多数据流的不同stream数据状态,所有stream end即会话end;
			}
		} else {
			// 超时状态查询错误码,判定异步是否异常
			errNum, errInfo = si.checkExcp()
		}

		// 构造响应包;
		dataOut.Status = status
		dataOut.Code = int32(errNum)
		if errInfo != nil {
			dataOut.Err = errInfo.Error()
		} else if dataOut.Pl == nil && si.sessType == protocol.LoaderInput_ONCE {
			// 非会话模式读取超时场景
			errInfo, errNum = frame.ErrorOnceExecTimeout, frame.AigesErrorOnceTimeout
		}
	}
	return
}

/*
异常中断
*/
func (si *ServiceInst) DataException(span *xsf.Span) (errNum int, errInfo error) {
	// 重置服务实例
	if si.setAlive(false) {
		si.setExcp(frame.AigesErrorEngInactive, frame.ErrorInstNotActive)
	}
	return
}

// 上行数据处理
func (si *ServiceInst) inputProc(input *[]buffer.DataMeta) (output []DataMeta, errNum int, errInfo error) {
	for _, stream := range *input {
		// 输入数据处理
		meta := DataMeta{stream.DataId, stream.Data.([]byte), int(stream.DataType),
			int(stream.Status), stream.FrameId, stream.Desc.Attribute}

		output = append(output, meta)
	}
	return
}

// 下行数据处理
func (si *ServiceInst) outputProc(input []DataMeta) (output []buffer.DataMeta, errNum int, errInfo error) {
	for _, v := range input {
		/*		if errNum, errInfo = si.enCodecV1(v); errInfo != nil {
				return
			}*/
		if attr := si.expect[v.DataId]; attr == nil {
			si.tool.Log.Errorw("invalid expect . attr is nil", "key", v.DataId, "sid", si.instHdl)
			return nil, frame.AigesErrorInvalidOut, frame.ErrorInvalidOutput
		}
		if si.expect[v.DataId].DataType != protocol.MetaDesc_DataType(v.DataType) {
			si.tool.Log.Errorw("invalid expect. datatype is not match", "key", v.DataId, "expect", si.expect,
				"output", v.DataType, "sid", si.instHdl)
			return nil, frame.AigesErrorInvalidOut, frame.ErrorInvalidOutput
		}

		meta := buffer.DataMeta{v.Data, v.DataId, si.outPutLastId[v.DataId], buffer.DataStatus(v.DataStatus),
			buffer.DataType(v.DataType), nil}
		si.outPutLastId[v.DataId]++
		meta.Desc = &protocol.MetaDesc{}
		meta.Desc.Attribute = make(map[string]string)
		for dk, dv := range v.DataDesc {
			meta.Desc.Attribute[dk] = dv
		}
		if pmd := si.expect[v.DataId]; pmd != nil {
			meta.Desc.Name = pmd.Name
			meta.Desc.DataType = pmd.DataType
			for tk, tv := range pmd.Attribute {
				meta.Desc.Attribute[tk] = tv
			}
		}
		meta.Desc.Attribute[protocol.Sequence] = strconv.Itoa(int(si.outPutLastId[v.DataId]))
		meta.Desc.Attribute[protocol.Status] = strconv.Itoa(v.DataStatus)

		if v.Data != nil {
			si.downDatas = append(si.downDatas, eventStorage{meta.Data, meta.DataId, meta.Desc})
		}

		output = append(output, meta)
	}
	return
}

//	异步协程: 单线程消费计算;
func (si *ServiceInst) asyncCalc() {
	defer catch.RecoverHandle()
	errNum, errInfo := nrtCheck(si)
	if errInfo == nil {
		switch si.sessType {
		case protocol.LoaderInput_STREAM:
			errNum, errInfo = si.sessionCalc()
		default:
			<-si.onceTrig // once多数据流写状态同步;
			errNum, errInfo = si.noneSessCalc()
		}
	}

	if errInfo != nil {
		si.setExcp(errNum, errInfo)
		si.outPutData.Signal() // 中断结果查询超时阻塞操作
		si.debug.WithErrorTag(errInfo.Error()).WithRetTag(strconv.Itoa(errNum))
		if si.downCtrl == downAsync {
			dataPostBack(si, nil, errNum, errInfo)
			si.tool.Log.Errorw("request fail with downMethod async", "code", errNum, "err", errInfo.Error())
		}
	}
	if si.downCtrl == downAsync {
		// 实例释放：
		si.DataException(nil)
		si.mngr.Release(si.context)
		// TODO FreeChan((si.wrapperHdl))
	}

	si.tool.Log.Debugw("asyncCalc flush eventlog start.",
		"sid", si.instHdl)

	si.tool.Log.Debugw("asyncCalc flush eventlog end.",
		"sid", si.instHdl)
	wrapperRet := si.eventUsrDetailCode()
	si.tool.Log.Debugw("asyncCalc flush metric info start.",
		"wrapperRet", wrapperRet,
		"sid", si.instHdl)
	si.excpLock.RLock()
	excpNum := si.excpNum
	si.excpLock.RUnlock()
	si.tool.Monitor.NewMetric(
		"label", 1,
		xsf.KV{"sid", si.instHdl},
		xsf.KV{"gray", conf.GrayLabel},
		xsf.KV{"ret", excpNum},
		xsf.KV{"wrapperRet", wrapperRet})
	si.tool.Log.Debugw("asyncCalc flush metric info end.",
		"sid", si.instHdl)
	si.tool.Log.Debugw("asyncCalc flush span info start.", "sid", si.instHdl)
	si.debugWrapper.End().Flush()
	si.debug.End().Flush()

	si.tool.Log.Debugw("asyncCalc flush span info end.", "sid", si.instHdl)

	si.instWg.Done()

	si.tool.Log.Debugw("asyncCalc end .release inst", "sid", si.instHdl)

	return
}

// 会话模式异步计算
func (si *ServiceInst) sessionCalc() (errNum int, errInfo error) {
	// 用户事件判定：资源申请事件(eg:ssb场景)
	actNew, ok := si.usrActions[EventNew] // 上层接口保障非nil
	if ok {
		var resp ActMsg
		if conf.WrapperAsync {
			si.respChan, errInfo = AllocChan(si.context)
			if errInfo != nil {
				errNum = frame.AigesErrorSessRepeat
				return
			}
		}
		resp, errNum, errInfo = actNew(si.context, &ActMsg{Params: si.params})
		if errInfo != nil {
			si.tool.Log.Errorw("EventNew action fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
			errInfo = eventUsrErr(errNum, errInfo)
			errNum = frame.WrapperCreateErr
			return
		}
		si.tool.Log.Debugw("EventNew action detail", "inst", si.instHdl, "params", si.params, "wrapperHdl", resp.WrapperHdl)
		si.wrapperHdl = resp.WrapperHdl
	}

	fin := false
	si.upDatas = make([]eventStorage, 0, 1)
	si.downDatas = make([]eventStorage, 0, 1)
	for si.sessAlive() && !fin {
		var merge []buffer.DataMeta
		merge, errNum, errInfo = si.inputData.ReadDataWithTime(nil)
		if errInfo != nil && errInfo != frame.ErrorSeqBufferEmpty {
			if errInfo == frame.ErrorSeqChanClosed {
				si.tool.Log.Infow("session dead. read audio buffer fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
				break
			}
			si.tool.Log.Errorw("read audio buffer fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
			break
		} else if merge == nil || len(merge) == 0 {
			// 内部重试读取,直至session timeout或AIExcp;
			continue
		}

		// storage event log & media data
		for _, v := range merge {
			si.upStatus[v.DataId] = v.Status
			si.upDatas = append(si.upDatas, eventStorage{v.Data, v.DataId, v.Desc})
			si.tool.Log.Infow("calc read buffer", "stream", v.DataId,
				"frame", v.FrameId, "status", v.Status, "sid", si.instHdl)
		}

		errNum, errInfo = nrtDataFill(&merge)
		if errInfo != nil {
			break
		}
		errNum, errInfo = si.checkExcp()
		if errInfo != nil {
			si.tool.Log.Errorw("sessionCalc checkExcp then exit", "sid", si.instHdl)
			break
		}
		fin, errNum, errInfo = si.metaTask(merge)
		if errInfo != nil {
			break
		}
	}

	if actDel, ok := si.usrActions[EventDel]; ok {
		_, _, _ = actDel(si.context, &ActMsg{WrapperHdl: si.wrapperHdl})
		si.tool.Log.Debugw("EventDel action detail", "inst", si.instHdl, "wrapperHdl", si.wrapperHdl)
	}

	return
}

// 非会话模式异步计算
func (si *ServiceInst) noneSessCalc() (errNum int, errInfo error) {
	// 读取数据 & 校验状态once
	output, errNum, errInfo := si.inputData.ReadDataWithTime(nil)
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
		si.upDatas = append(si.upDatas, eventStorage{stream.Data, stream.DataId, stream.Desc})
	}

	// 数据补齐操作(离线类引擎不做数据存储)
	errNum, errInfo = nrtDataFill(&output)
	if errInfo != nil {
		return
	}
	// 数据处理,包含编解码/重采样等
	engData, errNum, errInfo := si.inputProc(&output)
	if errInfo != nil {
		return
	}

	// 事件判定&once处理
	// todo # 如果python grpc也走此模式
	if conf.WrapperAsync {
		si.respChan, errInfo = AllocChan(si.context)
		if errInfo != nil {
			errNum = frame.AigesErrorSessRepeat
			return
		}
		// TODO defer FreeChan(si.instHdl)！！！！！！！！！！
	}
	actExec, ok := si.usrActions[EventOnceExec]
	if ok {
		var resp ActMsg
		resp, errNum, errInfo = actExec(si.context, &ActMsg{Params: si.params, DeliverData: engData})
		if errInfo != nil {
			si.tool.Log.Errorw("OnceExec action fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
			errInfo = eventUsrErr(errNum, errInfo)
			errNum = frame.WrapperExecErr
			return
		}
		if conf.StorageData {
			for _, tmpData := range engData {
				sErr := aigesUtils.StorageDecodeData(tmpData.Data, tmpData.DataId)
				if sErr != nil {
					si.tool.Log.Errorw("OnceExec storage data failed", "errInfo", sErr.Error(), "sid", si.instHdl)
				}
			}
		}

		// 如果是异步取数据方式
		if conf.WrapperAsync {
			resp, _ = <-si.respChan
			if resp.AsyncErr != nil {
				si.tool.Log.Errorw("OnceExec async callback error", "errInfo", resp.AsyncErr.Error(), "sid", si.instHdl)
				errInfo = eventUsrErr(resp.AsyncCode, resp.AsyncErr)
				errNum = frame.WrapperAsyncErr
				return
			}
		}
		//  输出数据处理
		si.wrapperHdl = resp.WrapperHdl
		si.downDatas = make([]eventStorage, 0, len(resp.DeliverData))
		respData := make([]buffer.DataMeta, 0, len(resp.DeliverData))
		for _, data := range resp.DeliverData {
			/*			if errNum, errInfo = si.enCodecV1(data); errInfo != nil {
						return
					}*/
			if attr := si.expect[data.DataId]; attr == nil {
				si.tool.Log.Errorw("invalid expect", "key", data.DataId, "sid", si.instHdl)
				return frame.AigesErrorInvalidOut, frame.ErrorInvalidOutput
			}
			if si.expect[data.DataId].DataType != protocol.MetaDesc_DataType(data.DataType) {
				si.tool.Log.Errorw("invalid expect type", "key", data.DataId, "type", si.expect[data.DataId].DataType, "sid", si.instHdl)
				return frame.AigesErrorInvalidOut, frame.ErrorInvalidOutput
			}
			//记录文本的数据返回
			if data.DataType == 0 {
				si.tool.Log.Debugw("OnceExec get result .", "key", data.DataId,
					"status", data.DataStatus, "type", data.DataType, "data", string(data.Data), "data length", len(data.Data), "sid", si.instHdl)
			} else {
				si.tool.Log.Debugw("OnceExec get result .", "key", data.DataId,
					"status", data.DataStatus, "type", data.DataType, "data length", len(data.Data), "sid", si.instHdl)
			}
			meta := buffer.DataMeta{data.Data, data.DataId, data.DataFrame,
				buffer.DataStatus(data.DataStatus), buffer.DataType(data.DataType), nil}
			meta.Desc = &protocol.MetaDesc{}
			meta.Desc.Attribute = make(map[string]string)
			if pmd := si.expect[data.DataId]; pmd != nil {
				meta.Desc.Name = pmd.Name
				meta.Desc.DataType = pmd.DataType
				for k, v := range pmd.Attribute {
					meta.Desc.Attribute[k] = v
				}
			}
			meta.Desc.Attribute[protocol.Sequence] = strconv.Itoa(int(si.outPutLastId[data.DataId]))
			meta.Desc.Attribute[protocol.Status] = strconv.Itoa(data.DataStatus)

			respData = append(respData, meta)
			si.downDatas = append(si.downDatas, eventStorage{meta.Data, meta.DataId, meta.Desc})
		}
		if si.downCtrl == downAsync {
			dataPostBack(si, &respData, 0, nil)
		} else {
			si.outPutData.WriteData(si.outPutId, respData)
			si.outPutId++
		}
	} else {
		si.tool.Log.Errorw("noneSessCalc can't find OnceExec Event to call")
		errNum, errInfo = frame.AigesErrorNilEvent, frame.ErrorNilRegEvent
	}
	return
}

func (si *ServiceInst) metaTask(inputs []buffer.DataMeta) (fin bool, errNum int, errInfo error) {
	fin = false
	deliver, errNum, errInfo := si.inputProc(&inputs)
	if errInfo != nil {
		return true, errNum, errInfo
	}

	// 数据写事件;
	actWrite, ok := si.usrActions[EventWrite]
	if !ok {
		return true, frame.WrapperWriteErr, errors.New("there is no register EventWrite by Register()")
	}

	si.tool.Log.Debugw("EventWrite action detail", "inst", si.instHdl, "wrapperHdl", si.wrapperHdl, "input", deliver)
	_, errNum, errInfo = actWrite(si.context, &ActMsg{WrapperHdl: si.wrapperHdl, DeliverData: deliver})
	if errInfo != nil {
		si.tool.Log.Errorw("usrAction EventWrite fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
		errInfo = eventUsrErr(errNum, errInfo)
		errNum = frame.WrapperWriteErr
		return true, errNum, errInfo
	}

	// check input upstream status
	lastInput := false
	if len(deliver) > 0 {
		lastInput = true
		for _, v := range si.upStatus {
			if v == buffer.DataStatusLast || v == buffer.DataStatusOnce {
				continue
			}
			lastInput = false
			break
		}
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
					errInfo = eventUsrErr(resp.AsyncCode, resp.AsyncErr)
					return true, frame.WrapperAsyncErr, errInfo
				}
			} else {
				actRead, ok := si.usrActions[EventRead]
				if !ok {
					return true, frame.WrapperReadErr, errors.New("there is no register EventRead by Register()")
				}
				resp, errNum, errInfo = actRead(si.context, &ActMsg{WrapperHdl: si.wrapperHdl})
				if errInfo != nil {
					si.tool.Log.Errorw("usrAction EventRead fail", "errNum", errNum, "errInfo", errInfo.Error(), "sid", si.instHdl)
					errInfo = eventUsrErr(errNum, errInfo)
					errNum = frame.WrapperReadErr
					return true, errNum, errInfo
				}
				si.tool.Log.Debugw("EventRead action detail", "inst", si.instHdl, "wrapperHdl", si.wrapperHdl, "output",
					resp)
			}

			output, errNum, errInfo := si.outputProc(resp.DeliverData)
			if errInfo != nil {
				return true, errNum, errInfo
			}

			if si.downCtrl == downAsync {
				dataPostBack(si, &output, 0, nil)
			} else {
				si.outPutData.WriteData(si.outPutId, output) // TODO 合成场景确认：fmt & enc
				si.outPutId++
			}
			// 冗余 todo del
			for k, _ := range output {
				if output[k].Status == buffer.DataStatusOnce || output[k].Status == buffer.DataStatusLast {
					fin = true
				}
			}

			if fin || !lastInput || waitCnt > syncRltWaitCnt {
				break
			}
			// 若同步读未取得计算结果,循环等待timeWait(50ms);
			if len(resp.DeliverData) == 0 { // 等价于 || resp.DeliverData == nil
				time.Sleep(time.Duration(syncRltWaitTime) * time.Millisecond)
				waitCnt++
				continue
			}
			waitCnt = 0
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

func (si *ServiceInst) CatchTag() (header map[string]string, params map[string]string, datas []eventStorage) {
	return si.headers, si.params, si.upDatas
}

func (si *ServiceInst) GetSessState() protocol.LoaderInput_SessState {
	return si.sessType
}

func (si *ServiceInst) eventUsrDetailCode() (code int) {
	si.excpLock.RLock()
	defer si.excpLock.RUnlock()
	if si.excpNum == 0 {
		return 0
	}
	if si.excpErr != nil && len(si.excpErr.Error()) > 0 {
		parts := strings.Split(si.excpErr.Error(), ";code=")
		if len(parts) == 2 {
			code, _ = strconv.Atoi(parts[1])
		}
	}
	return code
}

func eventUsrErr(code int, err error) error {
	return errors.New(err.Error() + ";code=" + strconv.Itoa(code))
}

func (si *ServiceInst) TaskAttr() (sub, svc string, appid string, uid string, cloudId string, composeId string) {
	return si.sub, si.svcId, si.appid, si.uid, si.cloudId, si.composeId
}

func (si *ServiceInst) WrapperTraceTag(key string, value string) {
	if conf.WrapperTrace {
		si.debugWrapper.WithTag(key, value)
	}
}
