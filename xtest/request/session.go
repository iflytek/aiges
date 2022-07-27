package request

import (
	"encoding/json"
	"errors"
	"frame"
	xsfcli "git.iflytek.com/AIaaS/xsf/client"
	"github.com/golang/protobuf/proto"
	"protocol"
	"strconv"
	"sync"
	"time"
	"xtest/analy"
	"xtest/util"
	_var "xtest/var"
)

func (r *Request) SessionCall(cli *xsfcli.Client, index int64) (info analy.ErrInfo) {
	// 下行结果缓存

	data, _ := json.Marshal(r.C.UpStreams)
	println(string(data))
	var indexs []int = make([]int, 0, len(r.C.UpStreams))
	for _, v := range r.C.UpStreams {
		streamIndex := index % int64(len(v.DataList))
		indexs = append(indexs, int(streamIndex))
	}
	// go routine 区分不同frame slice数据流
	var thrRslt []protocol.LoaderOutput = make([]protocol.LoaderOutput, 0, 1)
	var thrLock sync.Mutex
	reqSid := util.NewSid(r.C.TestSub)
	hdl, status, info := r.sessAIIn(cli, indexs, &thrRslt, &thrLock, reqSid)
	if info.ErrStr != nil {
		if len(hdl) != 0 {
			_ = r.sessAIExcp(cli, hdl, reqSid)
			return info
		}
	} else if status != protocol.LoaderOutput_END {
		info = r.sessAIOut(cli, hdl, reqSid, &thrRslt)
		if info.ErrStr != nil {
			_ = r.sessAIExcp(cli, hdl, reqSid)
			return info
		}
	}
	_ = r.sessAIExcp(cli, hdl, reqSid)
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
		case r.C.AsyncDrop <- _var.OutputMeta{reqSid, outType, v.Meta.Name, v.Meta.Attribute, v.Data}:
		default:
			// 异步channel满, 同步写;	key: sid-type-format-encoding, value: data
			key := reqSid + "-" + outType + "-" + v.Meta.Name
			r.downOutput(key, v.Data, cli.Log)
		}
	}
	return
}

func (r *Request) sessAIIn(cli *xsfcli.Client, indexs []int, thrRslt *[]protocol.LoaderOutput, thrLock *sync.Mutex, reqSid string) (hdl string, status protocol.LoaderOutput_RespStatus, info analy.ErrInfo) {
	// jbzhou5 并行网络协程监听
	r.C.ConcurrencyCnt.Add(1)
	defer r.C.ConcurrencyCnt.Dec() // jbzhou5 任务完成时-1

	// request构包；构造首包SeqNo=1,同加载器建立会话上下文信息; 故首帧不携带具体数据
	println(indexs)
	req := xsfcli.NewReq()
	req.SetParam("SeqNo", "1")
	req.SetParam("baseId", "0")
	req.SetParam("version", "v2")
	req.SetParam("waitTime", strconv.Itoa(r.C.TimeOut))
	dataIn := protocol.LoaderInput{}
	dataIn.State = protocol.LoaderInput_STREAM
	dataIn.ServiceId = r.C.SvcId
	dataIn.ServiceName = r.C.SvcName
	// 平台参数header
	dataIn.Headers = make(map[string]string)
	dataIn.Headers["sid"] = reqSid
	dataIn.Headers["status"] = "0"
	for k, v := range r.C.Header {
		dataIn.Headers[k] = v
	}
	// 能力参数params
	dataIn.Params = make(map[string]string)
	for k, v := range r.C.Params {
		dataIn.Params[k] = v
	}
	// 期望输出expect
	for k, _ := range r.C.DownExpect {
		dataIn.Expect = append(dataIn.Expect, &r.C.DownExpect[k])
	}

	input, err := proto.Marshal(&dataIn)
	if err != nil {
		cli.Log.Errorw("sessAIIn marshal create request fail", "err", err.Error(), "params", dataIn.Params)
		return hdl, status, analy.ErrInfo{ErrCode: -1, ErrStr: err}
	}

	rd := xsfcli.NewData()
	rd.Append(input)
	req.AppendData(rd)

	caller := xsfcli.NewCaller(cli)
	analy.Perf.Record(reqSid, "", analy.DataBegin, analy.SessBegin, analy.UP, 0, "")
	resp, ecode, err := caller.SessionCall(xsfcli.CREATE, r.C.SvcName, "AIIn", req, time.Duration(r.C.TimeOut+r.C.LossDeviation)*time.Millisecond)
	if err != nil {
		cli.Log.Errorw("sessAIIn Create request fail", "err", err.Error(), "code", ecode, "params", dataIn.Params)
		analy.Perf.Record(reqSid, resp.Handle(), analy.DataBegin, analy.SessBegin, analy.DOWN, int(ecode), err.Error())
		return hdl, status, analy.ErrInfo{ErrCode: int(ecode), ErrStr: err}
	}
	hdl = resp.Session()
	analy.Perf.Record(reqSid, resp.Handle(), analy.DataBegin, analy.SessBegin, analy.DOWN, 0, "")

	// data stream: 相同UpInterval合并发送;
	merge := make(map[int] /*UpInterval*/ map[int]int /*stream's index*/)
	for k, v := range indexs {
		_, exist := merge[r.C.UpStreams[k].UpInterval]
		if !exist {
			merge[r.C.UpStreams[k].UpInterval] = make(map[int]int)
		}
		merge[r.C.UpStreams[k].UpInterval][k] = v
		// streamIndex: index of indexs ; fileIndex: value of indexs[index]
	}

	errChan := make(chan analy.ErrInfo, 1) // 仅保存首个错误码;
	defer close(errChan)
	var rwg sync.WaitGroup
	for k, v := range merge {
		println("merge ", v, k)
		rwg.Add(1)
		go r.multiUpStream(cli, &rwg, hdl, k, v, thrRslt, thrLock, errChan)
	}
	rwg.Wait() // 异步协程上行数据交互结束
	select {
	case einfo := <-errChan:
		return hdl, status, einfo
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

func (r *Request) multiUpStream(cli *xsfcli.Client, swg *sync.WaitGroup, session string, interval int, indexs map[int]int, pm *[]protocol.LoaderOutput, sm *sync.Mutex, errchan chan analy.ErrInfo) {
	// jbzhou5 并行网络协程监听
	r.C.ConcurrencyCnt.Add(1)
	defer r.C.ConcurrencyCnt.Dec() // jbzhou5 任务完成时-1

	defer swg.Done()

	sTime := time.Now()
	dataSendLen := make([]int, len(indexs)) // send current size
	dataSeq := make([]int, len(indexs))     // send current index
	endLen := 0

	h264FrameSizes := make(map[int][]int)
	for streamId, fileId := range indexs {
		if r.C.UpStreams[streamId].DataType == protocol.MetaDesc_DataType(protocol.MetaData_VIDEO) {
			h264FrameSizes[streamId] = GetH264Frames(r.C.UpStreams[streamId].DataList[fileId])
			cli.Log.Debugw("upstream get h264 frames size. ",
				"frames", h264FrameSizes, "length", len(h264FrameSizes))
		}
	}
	for dataId := 1; len(indexs) > 0; dataId++ {
		req := xsfcli.NewReq()
		req.SetParam("baseId", "0")
		req.SetParam("version", "v2")
		req.SetParam("waitTime", strconv.Itoa(r.C.TimeOut))
		_ = req.Session(session)
		dataIn := protocol.LoaderInput{}
		dataIn.SyncId = int32(dataId)
		// 上行数据流实体数据
		for streamId, fileId := range indexs {
			var size int
			if r.C.UpStreams[streamId].DataType == protocol.MetaDesc_DataType(protocol.MetaData_VIDEO) {
				if (dataId - 1) > len(h264FrameSizes[streamId])-1 {
					size = 0
				} else {
					size = h264FrameSizes[streamId][dataId-1]
				}
			} else {
				size = r.C.UpStreams[streamId].UpSlice
			}
			upStatus := protocol.EngInputData_CONTINUE
			if dataSendLen[streamId] >= len(r.C.UpStreams[streamId].DataList[fileId]) {
				continue // 该上行数据流已发送完毕
			}
			//if dataSendLen[streamId] == 0 {
			//	upStatus = protocol.EngInputData_BEGIN
			//}
			if dataSendLen[streamId]+size >= len(r.C.UpStreams[streamId].DataList[fileId]) {
				size = len(r.C.UpStreams[streamId].DataList[fileId]) - dataSendLen[streamId]
				upStatus = protocol.EngInputData_END
				endLen += 1
			}
			upData := r.C.UpStreams[streamId].DataList[fileId][dataSendLen[streamId] : dataSendLen[streamId]+size]
			desc := make(map[string]string)
			for dk, dv := range r.C.UpStreams[streamId].DataDesc {
				desc[dk] = dv
			}
			cli.Log.Debugw("send data", "status", upStatus, "streamId", streamId, "fileId", fileId,
				"sendLen", size, "totalSendLen", dataSendLen[streamId], "totalLen", len(r.C.UpStreams[streamId].DataList[fileId]))
			md := protocol.MetaDesc{
				Name:      r.C.UpStreams[streamId].Name,
				DataType:  r.C.UpStreams[streamId].DataType,
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
			r.unBlockChanWrite(errchan, analy.ErrInfo{-1, err})
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

		resp, ecode, err := caller.SessionCall(xsfcli.CONTINUE, r.C.SvcName, "AIIn", req, time.Duration(r.C.TimeOut+r.C.LossDeviation)*time.Millisecond)
		if err != nil && ecode != frame.AigesErrorEngInactive {
			cli.Log.Errorw("multiUpStream Create request fail", "err", err.Error(), "code", ecode, "params", dataIn.Params)
			r.unBlockChanWrite(errchan, analy.ErrInfo{ErrCode: int(ecode), ErrStr: err})
			analy.Perf.Record("", req.Handle(), analy.DataContinue, analy.SessContinue, analy.DOWN, int(ecode), err.Error())
			return
		}
		// 下行结果输出
		dataOut := protocol.LoaderOutput{}
		err = proto.Unmarshal(resp.GetData()[0].Data, &dataOut)
		if err != nil {
			cli.Log.Errorw("multiUpStream Resp Unmarshal fail", "err", err.Error(), "respData", resp.GetData()[0].Data)
			r.unBlockChanWrite(errchan, analy.ErrInfo{ErrCode: -1, ErrStr: err})
			return
		}

		switch dataOut.Code {
		case 0: // nothing to do
		case frame.AigesErrorEngInactive:
			return
		default:
			cli.Log.Errorw("multiUpStream get engine err", "err", dataOut.Err, "code", dataOut.Code, "params", dataIn.Params)
			r.unBlockChanWrite(errchan, analy.ErrInfo{ErrCode: int(dataOut.Code), ErrStr: errors.New(dataOut.Err)})
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
		r.rtCalibration(dataId, interval, sTime)
	}
}

// 实时性校准,用于校准发包大小及发包时间间隔之间的实时性.
func (r *Request) rtCalibration(curReq int, interval int, sTime time.Time) {
	cTime := int(time.Now().Sub(sTime).Nanoseconds() / (1000 * 1000)) // ssb至今绝对时长.ms
	expect := interval * (curReq + 1)                                 // 期望发包时间
	if expect > cTime {
		time.Sleep(time.Millisecond * time.Duration(expect-cTime))
	}
}

// downStream 下行调用单线程;
func (r *Request) sessAIOut(cli *xsfcli.Client, hdl string, sid string, rslt *[]protocol.LoaderOutput) (info analy.ErrInfo) {
	// jbzhou5 并行网络协程监听
	r.C.ConcurrencyCnt.Add(1)
	defer r.C.ConcurrencyCnt.Dec() // jbzhou5 任务完成时-1

	// loop read downstream result
	for {
		req := xsfcli.NewReq()
		req.SetParam("baseId", "0")
		req.SetParam("version", "v2")
		req.SetParam("waitTime", strconv.Itoa(r.C.TimeOut))
		_ = req.Session(sid)
		dataIn := protocol.LoaderInput{}

		input, err := proto.Marshal(&dataIn)
		if err != nil {
			cli.Log.Errorw("sessAIOut marshal create request fail", "err", err.Error(), "params", dataIn.Params)
			return analy.ErrInfo{ErrCode: -1, ErrStr: err}
		}

		rd := xsfcli.NewData()
		rd.Append(input)
		req.AppendData(rd)
		_ = req.Session(hdl)

		caller := xsfcli.NewCaller(cli)
		resp, ecode, err := caller.SessionCall(xsfcli.CONTINUE, r.C.SvcName, "AIOut", req, time.Duration(r.C.TimeOut+r.C.LossDeviation)*time.Millisecond)
		if err != nil {
			cli.Log.Errorw("sessAIOut request fail", "err", err.Error(), "code", ecode, "params", dataIn.Params)
			if ecode == frame.AigesErrorEngInactive { // reset 10101 inactive
				err = nil
			}
			analy.Perf.Record("", req.Handle(), analy.DataContinue, analy.SessContinue, analy.DOWN, int(ecode), err.Error())

			return analy.ErrInfo{ErrCode: int(ecode), ErrStr: err}
		}

		// 解析结果、输出落盘
		dataOut := protocol.LoaderOutput{}
		err = proto.Unmarshal(resp.GetData()[0].Data, &dataOut)
		if err != nil {
			cli.Log.Errorw("sessAIOut Resp Unmarshal fail", "err", err.Error(), "respData", resp.GetData()[0].Data)
			return analy.ErrInfo{ErrCode: -1, ErrStr: err}
		}

		*rslt = append(*rslt, dataOut)
		analy.Perf.Record("", req.Handle(), analy.DataStatus(int(dataOut.Status)), analy.SessContinue, analy.DOWN, int(dataOut.Code), dataOut.Err)
		cli.Log.Debugw("sessAIOut get resp result", "hdl", sid, "result", dataOut)
		if dataOut.Status == protocol.LoaderOutput_END {
			return analy.ErrInfo{ErrStr: err} // last result
		}
	}

	return
}

func (r *Request) sessAIExcp(cli *xsfcli.Client, hdl string, sid string) (err error) {
	// jbzhou5 并行网络协程监听
	r.C.ConcurrencyCnt.Add(1)
	defer r.C.ConcurrencyCnt.Dec() // jbzhou5 任务完成时-1

	req := xsfcli.NewReq()
	req.SetParam("baseId", "0")
	req.SetParam("waitTime", strconv.Itoa(r.C.TimeOut))
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
	_, _, err = caller.SessionCall(xsfcli.CONTINUE, r.C.SvcName, "AIExcp", req, time.Duration(r.C.TimeOut+r.C.LossDeviation)*time.Millisecond)
	return
}

// upStream first error
func (r *Request) unBlockChanWrite(ch chan analy.ErrInfo, err analy.ErrInfo) {
	select {
	case ch <- err:
	default:
		// ch full, return. save first err code
	}
}
