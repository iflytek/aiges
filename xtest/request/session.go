package request

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/xfyun/aiges/frame"
	"github.com/xfyun/aiges/protocol"
	"github.com/xfyun/aiges/xtest/analy"
	"github.com/xfyun/aiges/xtest/util"
	_var "github.com/xfyun/aiges/xtest/var"
	xsfcli "github.com/xfyun/xsf/client"
	"strconv"
	"sync"
	"time"
)

func SessionCall(cli *xsfcli.Client, index int64) (code int, err error) {
	// 下行结果缓存
	var indexs []int = make([]int, 0, len(_var.UpStreams))
	for _, v := range _var.UpStreams {
		streamIndex := index % int64(len(v.DataList))
		indexs = append(indexs, int(streamIndex))
	}
	// go routine 区分不同frame slice数据流
	var thrRslt []protocol.LoaderOutput = make([]protocol.LoaderOutput, 0, 1)
	var thrLock sync.Mutex
	reqSid := util.NewSid(_var.TestSub)
	hdl, status, code, err := sessAIIn(cli, indexs, &thrRslt, &thrLock, reqSid)
	if err != nil {
		if len(hdl) != 0 {
			_ = sessAIExcp(cli, hdl, reqSid)
			return
		}
	} else if status != protocol.LoaderOutput_END {
		code, err = sessAIOut(cli, hdl, reqSid, &thrRslt)
		if err != nil {
			_ = sessAIExcp(cli, hdl, reqSid)
			return
		}
	}
	_ = sessAIExcp(cli, hdl, reqSid)

	// 结果落盘
	tmpMerge := make(map[string] /*streamId*/ *protocol.Payload)
	for k, _ := range thrRslt {
		for _, d := range thrRslt[k].Pl {
			meta, exist := tmpMerge[d.Meta.Name]
			if exist {
				tmpMerge[d.Meta.Name].Data = append(meta.Data, d.Data...)
			} else {
				tmpMerge[d.Meta.Name] = d
			}
		}
	}

	for _, v := range tmpMerge {
		var outType string = "invalidType"
		switch v.Meta.DataType {
		case protocol.MetaDesc_TEXT:
			outType = "text"
		case protocol.MetaDesc_AUDIO:
			outType = "audio"
		case protocol.MetaDesc_IMAGE:
			outType = "image"
		case protocol.MetaDesc_VIDEO:
			outType = "video"
		}

		select {
		case _var.AsyncDrop <- _var.OutputMeta{reqSid, outType, v.Meta.Name, v.Meta.Attribute, v.Data}:
		default:
			// 异步channel满, 同步写;	key: sid-type-format-encoding, value: data
			key := reqSid + "-" + outType + "-" + v.Meta.Name
			downOutput(key, v.Data, cli.Log)
		}
	}
	return
}

func sessAIIn(cli *xsfcli.Client, indexs []int, thrRslt *[]protocol.LoaderOutput, thrLock *sync.Mutex, reqSid string) (hdl string, status protocol.LoaderOutput_RespStatus, code int, err error) {
	// request构包；构造首包SeqNo=1,同加载器建立会话上下文信息; 故首帧不携带具体数据
	req := xsfcli.NewReq()
	req.SetParam("SeqNo", "1")
	req.SetParam("baseId", "0")
	req.SetParam("version", "v2")
	req.SetParam("waitTime", strconv.Itoa(_var.TimeOut))
	dataIn := protocol.LoaderInput{}
	dataIn.State = protocol.LoaderInput_STREAM
	dataIn.ServiceId = _var.SvcId
	dataIn.ServiceName = _var.SvcName
	// 平台参数header
	dataIn.Headers = make(map[string]string)
	dataIn.Headers["sid"] = reqSid
	dataIn.Headers["status"] = "0"
	for k, v := range _var.Header {
		dataIn.Headers[k] = v
	}
	// 能力参数params
	dataIn.Params = make(map[string]string)
	for k, v := range _var.Params {
		dataIn.Params[k] = v
	}
	// 期望输出expect
	for k, _ := range _var.DownExpect {
		dataIn.Expect = append(dataIn.Expect, &_var.DownExpect[k])
	}

	input, err := proto.Marshal(&dataIn)
	if err != nil {
		cli.Log.Errorw("sessAIIn marshal create request fail", "err", err.Error(), "params", dataIn.Params)
		return hdl, status, -1, err
	}

	rd := xsfcli.NewData()
	rd.Append(input)
	req.AppendData(rd)

	caller := xsfcli.NewCaller(cli)
	analy.Perf.Record(reqSid, "", analy.DataBegin, analy.SessBegin, analy.UP, 0, "")
	resp, ecode, err := caller.SessionCall(xsfcli.CREATE, _var.SvcName, "AIIn", req, time.Duration(_var.TimeOut+_var.LossDeviation)*time.Millisecond)
	if err != nil {
		cli.Log.Errorw("sessAIIn Create request fail", "err", err.Error(), "code", ecode, "params", dataIn.Params)
		analy.Perf.Record(reqSid, resp.Handle(), analy.DataBegin, analy.SessBegin, analy.DOWN, int(ecode), err.Error())
		return hdl, status, int(ecode), err
	}
	hdl = resp.Session()
	analy.Perf.Record(reqSid, resp.Handle(), analy.DataBegin, analy.SessBegin, analy.DOWN, 0, "")

	// data stream: 相同UpInterval合并发送;
	merge := make(map[int] /*UpInterval*/ map[int]int /*stream's index*/)
	for k, v := range indexs {
		_, exist := merge[_var.UpStreams[k].UpInterval]
		if !exist {
			merge[_var.UpStreams[k].UpInterval] = make(map[int]int)
		}
		merge[_var.UpStreams[k].UpInterval][k] = v
		// streamIndex: index of indexs ; fileIndex: value of indexs[index]
	}

	errChan := make(chan struct {
		code int
		err  error
	}, 1) // 仅保存首个错误码;
	defer close(errChan)
	var rwg sync.WaitGroup
	for k, v := range merge {
		rwg.Add(1)
		go multiUpStream(cli, &rwg, hdl, k, v, thrRslt, thrLock, errChan)
	}
	rwg.Wait() // 异步协程上行数据交互结束
	select {
	case einfo := <-errChan:
		return hdl, status, einfo.code, einfo.err
	default:
		// unblock; check status
		for k, _ := range *thrRslt {
			if (*thrRslt)[k].Status == protocol.LoaderOutput_END {
				status = (*thrRslt)[k].Status
			}
		}
	}
	return
}

func multiUpStream(cli *xsfcli.Client, swg *sync.WaitGroup, session string, interval int, indexs map[int]int, pm *[]protocol.LoaderOutput, sm *sync.Mutex, errchan chan struct {
	code int
	err  error
}) {
	defer swg.Done()

	sTime := time.Now()
	dataSendLen := make([]int, len(indexs)) // send current size
	dataSeq := make([]int, len(indexs))     // send current index
	endLen := 0

	h264FrameSizes := make(map[int][]int)
	for streamId, fileId := range indexs {
		if _var.UpStreams[streamId].DataType == protocol.MetaDesc_DataType(protocol.MetaData_VIDEO) {
			h264FrameSizes[streamId] = GetH264Frames(_var.UpStreams[streamId].DataList[fileId])
			cli.Log.Debugw("upstream get h264 frames size. ",
				"frames", h264FrameSizes, "length", len(h264FrameSizes))
		}
	}
	for dataId := 1; len(indexs) > 0; dataId++ {
		req := xsfcli.NewReq()
		req.SetParam("baseId", "0")
		req.SetParam("version", "v2")
		req.SetParam("waitTime", strconv.Itoa(_var.TimeOut))
		_ = req.Session(session)
		dataIn := protocol.LoaderInput{}
		dataIn.SyncId = int32(dataId)
		// 上行数据流实体数据
		for streamId, fileId := range indexs {
			var size int
			if _var.UpStreams[streamId].DataType == protocol.MetaDesc_DataType(protocol.MetaData_VIDEO) {
				if (dataId - 1) > len(h264FrameSizes[streamId])-1 {
					size = 0
				} else {
					size = h264FrameSizes[streamId][dataId-1]
				}
			} else {
				size = _var.UpStreams[streamId].UpSlice
			}
			upStatus := protocol.EngInputData_CONTINUE
			if dataSendLen[streamId] >= len(_var.UpStreams[streamId].DataList[fileId]) {
				continue // 该上行数据流已发送完毕
			}
			//if dataSendLen[streamId] == 0 {
			//	upStatus = protocol.EngInputData_BEGIN
			//}
			if dataSendLen[streamId]+size >= len(_var.UpStreams[streamId].DataList[fileId]) {
				size = len(_var.UpStreams[streamId].DataList[fileId]) - dataSendLen[streamId]
				upStatus = protocol.EngInputData_END
				endLen += 1
			}
			upData := _var.UpStreams[streamId].DataList[fileId][dataSendLen[streamId] : dataSendLen[streamId]+size]
			desc := make(map[string]string)
			for dk, dv := range _var.UpStreams[streamId].DataDesc {
				desc[dk] = dv
			}
			cli.Log.Debugw("send data", "status", upStatus, "streamId", streamId, "fileId", fileId,
				"sendLen", size, "totalSendLen", dataSendLen[streamId], "totalLen", len(_var.UpStreams[streamId].DataList[fileId]))
			md := protocol.MetaDesc{
				Name:      _var.UpStreams[streamId].Name,
				DataType:  _var.UpStreams[streamId].DataType,
				Attribute: desc}
			md.Attribute["seq"] = strconv.Itoa(dataSeq[streamId])
			md.Attribute["status"] = strconv.Itoa(int(upStatus))
			inputmeta := protocol.Payload{Meta: &md, Data: upData}
			dataIn.Pl = append(dataIn.Pl, &inputmeta)
			dataSendLen[streamId] += size
			dataSeq[streamId]++
		}

		if len(dataIn.Pl) == 0 {
			break // 所有上行数据流发送完毕;
		}

		input, err := proto.Marshal(&dataIn)
		if err != nil {
			cli.Log.Errorw("multiUpStream marshal create request fail", "err", err.Error(), "params", dataIn.Params)
			unBlockChanWrite(errchan, struct {
				code int
				err  error
			}{-1, err})
			return
		}

		rd := xsfcli.NewData()
		rd.Append(input)
		req.AppendData(rd)
		caller := xsfcli.NewCaller(cli)

		//TODO 会话模式当前只支持一个数据流或者多个上传间隔相同的数据流做性能测试
		if endLen < len(indexs) {
			analy.Perf.Record("", req.Handle(), analy.DataContinue, analy.SessContinue, analy.UP, 0, "")
		} else {
			analy.Perf.Record("", req.Handle(), analy.DataEnd, analy.SessContinue, analy.UP, 0, "")
		}

		resp, ecode, err := caller.SessionCall(xsfcli.CONTINUE, _var.SvcName, "AIIn", req, time.Duration(_var.TimeOut+_var.LossDeviation)*time.Millisecond)
		if err != nil && ecode != frame.AigesErrorEngInactive {
			cli.Log.Errorw("multiUpStream Create request fail", "err", err.Error(), "code", ecode, "params", dataIn.Params)
			unBlockChanWrite(errchan, struct {
				code int
				err  error
			}{int(ecode), err})
			analy.Perf.Record("", req.Handle(), analy.DataContinue, analy.SessContinue, analy.DOWN, int(ecode), err.Error())
			return
		}
		// 下行结果输出
		dataOut := protocol.LoaderOutput{}
		err = proto.Unmarshal(resp.GetData()[0].Data, &dataOut)
		if err != nil {
			cli.Log.Errorw("multiUpStream Resp Unmarshal fail", "err", err.Error(), "respData", resp.GetData()[0].Data)
			unBlockChanWrite(errchan, struct {
				code int
				err  error
			}{-1, err})
			return
		}

		switch dataOut.Code {
		case 0: // nothing to do
		case frame.AigesErrorEngInactive:
			return
		default:
			cli.Log.Errorw("multiUpStream get engine err", "err", dataOut.Err, "code", dataOut.Code, "params", dataIn.Params)
			unBlockChanWrite(errchan, struct {
				code int
				err  error
			}{int(dataOut.Code), errors.New(dataOut.Err)})
			analy.Perf.Record("", req.Handle(), analy.DataContinue, analy.SessContinue, analy.DOWN, int(dataOut.Code), dataOut.Err)
			return // engine err but not 10101
		}

		// 同步下行数据
		if len(dataOut.Pl) > 0 {
			(*sm).Lock()
			*pm = append(*pm, dataOut)
			cli.Log.Debugw("multiUpStream get resp result", "hdl", session, "result", dataOut)
			(*sm).Unlock()
			analy.Perf.Record("", req.Handle(), analy.DataStatus(int(dataOut.Status)), analy.SessContinue, analy.DOWN, 0, "")
		}
		if dataOut.Status == protocol.LoaderOutput_END {
			return // last result
		}

		// wait ms.动态调整校准上行数据实时率, 考虑其他接口耗时.
		rtCalibration(dataId, interval, sTime)
	}
}

// 实时性校准,用于校准发包大小及发包时间间隔之间的实时性.
func rtCalibration(curReq int, interval int, sTime time.Time) {
	cTime := int(time.Now().Sub(sTime).Nanoseconds() / (1000 * 1000)) // ssb至今绝对时长.ms
	expect := interval * (curReq + 1)                                 // 期望发包时间
	if expect > cTime {
		time.Sleep(time.Millisecond * time.Duration(expect-cTime))
	}
}

// downStream 下行调用单线程;
func sessAIOut(cli *xsfcli.Client, hdl string, sid string, rslt *[]protocol.LoaderOutput) (code int, err error) {
	// loop read downstream result
	for {
		req := xsfcli.NewReq()
		req.SetParam("baseId", "0")
		req.SetParam("version", "v2")
		req.SetParam("waitTime", strconv.Itoa(_var.TimeOut))
		_ = req.Session(sid)
		dataIn := protocol.LoaderInput{}

		input, err := proto.Marshal(&dataIn)
		if err != nil {
			cli.Log.Errorw("sessAIOut marshal create request fail", "err", err.Error(), "params", dataIn.Params)
			return -1, err
		}

		rd := xsfcli.NewData()
		rd.Append(input)
		req.AppendData(rd)
		_ = req.Session(hdl)

		caller := xsfcli.NewCaller(cli)
		resp, ecode, err := caller.SessionCall(xsfcli.CONTINUE, _var.SvcName, "AIOut", req, time.Duration(_var.TimeOut+_var.LossDeviation)*time.Millisecond)
		if err != nil {
			cli.Log.Errorw("sessAIOut request fail", "err", err.Error(), "code", ecode, "params", dataIn.Params)
			if ecode == frame.AigesErrorEngInactive { // reset 10101 inactive
				err = nil
			}
			analy.Perf.Record("", req.Handle(), analy.DataContinue, analy.SessContinue, analy.DOWN, int(ecode), err.Error())

			return int(ecode), err
		}

		// 解析结果、输出落盘
		dataOut := protocol.LoaderOutput{}
		err = proto.Unmarshal(resp.GetData()[0].Data, &dataOut)
		if err != nil {
			cli.Log.Errorw("sessAIOut Resp Unmarshal fail", "err", err.Error(), "respData", resp.GetData()[0].Data)
			return -1, err
		}

		*rslt = append(*rslt, dataOut)
		analy.Perf.Record("", req.Handle(), analy.DataStatus(int(dataOut.Status)), analy.SessContinue, analy.DOWN, int(dataOut.Code), dataOut.Err)
		cli.Log.Debugw("sessAIOut get resp result", "hdl", sid, "result", dataOut)
		if dataOut.Status == protocol.LoaderOutput_END {
			return code, err // last result
		}
	}

	return
}

func sessAIExcp(cli *xsfcli.Client, hdl string, sid string) (err error) {

	req := xsfcli.NewReq()
	req.SetParam("baseId", "0")
	req.SetParam("waitTime", strconv.Itoa(_var.TimeOut))
	dataIn := protocol.LoaderInput{}
	input, err := proto.Marshal(&dataIn)
	if err != nil {
		cli.Log.Errorw("sessAIExcp marshal create request fail", "err", err.Error(), "params", dataIn.Params)
		return
	}

	rd := xsfcli.NewData()
	rd.Append(input)
	req.AppendData(rd)
	_ = req.Session(hdl)

	caller := xsfcli.NewCaller(cli)
	_, _, err = caller.SessionCall(xsfcli.CONTINUE, _var.SvcName, "AIExcp", req, time.Duration(_var.TimeOut+_var.LossDeviation)*time.Millisecond)
	return
}

// upStream first error
func unBlockChanWrite(ch chan struct {
	code int
	err  error
}, err struct {
	code int
	err  error
}) {
	select {
	case ch <- err:
	default:
		// ch full, return. save first err code
	}
}
