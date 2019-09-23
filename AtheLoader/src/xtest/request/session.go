package request

import (
	"errors"
	"frame"
	xsfcli "git.xfyun.cn/AIaaS/xsf-external/client"
	"github.com/golang/protobuf/proto"
	"protocol"
	"strconv"
	"sync"
	"time"
	"xtest/util"
	_var "xtest/var"
)

func SessionCall(cli *xsfcli.Client, index int64) (code int, err error) {
	// 下行结果缓存
	var indexs []int = make([]int, 0, len(_var.UpStreams))
	for _, v := range _var.UpStreams {
		streamIndex := index % int64(len(v.DataList))
		indexs = append(indexs, int(streamIndex))
	}
	// go routine 区分不同frame slice数据流
	var thrRslt []protocol.EngOutputData = make([]protocol.EngOutputData, 0, 1)
	var thrLock sync.Mutex
	reqSid := util.NewSid(_var.TestSub)
	hdl, status, code, err := sessAIIn(cli, indexs, &thrRslt, &thrLock, reqSid)
	if err != nil {
		if len(hdl) != 0 {
			_ = sessAIExcp(cli, hdl, reqSid)
			return
		}
	} else if status != protocol.EngOutputData_END {
		code, err = sessAIOut(cli, hdl, reqSid, &thrRslt)
		if err != nil {
			_ = sessAIExcp(cli, hdl, reqSid)
			return
		}
	}
	_ = sessAIExcp(cli, hdl, reqSid)

	// 结果落盘
	tmpMerge := make(map[string] /*streamId*/ *protocol.MetaData)
	for k, _ := range thrRslt {
		for _, d := range thrRslt[k].DataList {
			meta, exist := tmpMerge[d.DataId]
			if exist {
				tmpMerge[d.DataId].Data = append(meta.Data, d.Data...)
			} else {
				tmpMerge[d.DataId] = d
			}
		}
	}

	for _, v := range tmpMerge {
		var outType string = "invalidType"
		switch v.DataType {
		case protocol.MetaData_TEXT:
			outType = "text"
		case protocol.MetaData_AUDIO:
			outType = "audio"
		case protocol.MetaData_IMAGE:
			outType = "image"
		case protocol.MetaData_VIDEO:
			outType = "video"
		}

		select {
		case _var.AsyncDrop <- _var.OutputMeta{reqSid, outType, v.Format, v.Encoding, v.Data}:
		default:
			// 异步channel满, 同步写;	key: sid-type-format-encoding, value: data
			key := reqSid + "-" + outType + "-" + v.Format + "-" + v.Encoding
			downOutput(key, v.Data, cli.Log)
		}
	}
	return
}

func sessAIIn(cli *xsfcli.Client, indexs []int, thrRslt *[]protocol.EngOutputData, thrLock *sync.Mutex, reqSid string) (hdl string, status protocol.EngOutputData_DataStatus, code int, err error) {
	// request构包；构造首包SeqNo=1,同加载器建立会话上下文信息; 故首帧不携带具体数据
	req := xsfcli.NewReq()
	req.SetParam("SeqNo", "1")
	req.SetParam("baseId", "0")
	req.SetParam("waitTime", strconv.Itoa(_var.TimeOut))
	dataIn := protocol.EngInputData{}
	dataIn.EngParam = make(map[string]string)
	for k, v := range _var.UpParams {
		dataIn.EngParam[k] = v
	}
	dataIn.EngParam["sid"] = reqSid

	input, err := proto.Marshal(&dataIn)
	if err != nil {
		cli.Log.Errorw("sessAIIn marshal create request fail", "err", err.Error(), "params", dataIn.EngParam)
		return hdl, status, -1, err
	}

	rd := xsfcli.NewData()
	rd.Append(input)
	req.AppendData(rd)

	caller := xsfcli.NewCaller(cli)
	resp, ecode, err := caller.SessionCall(xsfcli.CREATE, _var.SvcName, "AIIn", req, time.Duration(_var.TimeOut+_var.LossDeviation)*time.Millisecond)
	if err != nil {
		cli.Log.Errorw("sessAIIn Create request fail", "err", err.Error(), "code", ecode, "params", dataIn.EngParam)
		return hdl, status, int(ecode), err
	}
	hdl = resp.Session()

	// data stream: 相同UpInterval合并发送;
	merge := make(map[int] /*UpInterval*/ []int /*indexs's index*/)
	for k, _ := range indexs {
		_, exist := merge[_var.UpStreams[k].UpInterval]
		if !exist {
			merge[_var.UpStreams[k].UpInterval] = make([]int, 0, 1)
		}
		merge[_var.UpStreams[k].UpInterval] = append(merge[_var.UpStreams[k].UpInterval], k)
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
		go multiUpStream(cli, &rwg, hdl, k, v, reqSid, thrRslt, thrLock, errChan)
	}
	rwg.Wait() // 异步协程上行数据交互结束
	select {
	case einfo := <-errChan:
		return hdl, status, einfo.code, einfo.err
	default:
		// unblock; check status
		for k, _ := range *thrRslt {
			if (*thrRslt)[k].Status == protocol.EngOutputData_END {
				status = (*thrRslt)[k].Status
			}
		}
	}
	return
}

func multiUpStream(cli *xsfcli.Client, swg *sync.WaitGroup, session string, interval int, indexs []int, sid string, pm *[]protocol.EngOutputData, sm *sync.Mutex, errchan chan struct {
	code int
	err  error
}) {
	defer swg.Done()

	sTime := time.Now()
	dataSendLen := make([]int, len(indexs))
	for dataId := 0; len(indexs) > 0; dataId++ {
		req := xsfcli.NewReq()
		req.SetParam("baseId", "0")
		req.SetParam("waitTime", strconv.Itoa(_var.TimeOut))
		_ = req.Session(session)
		dataIn := protocol.EngInputData{}
		dataIn.EngParam = make(map[string]string)
		for k, v := range _var.UpParams {
			dataIn.EngParam[k] = v
		}
		dataIn.EngParam["sid"] = sid

		// 上行数据流实体数据
		for streamId, fileId := range indexs {
			upStatus := protocol.MetaData_CONTINUE
			if dataSendLen[streamId] == 0 {
				upStatus = protocol.MetaData_BEGIN
			}
			var upData []byte
			sendLen := _var.UpStreams[streamId].UpSlice
			if _var.UpStreams[streamId].SliceOn == 0 { // 切片开关：关闭
				sendLen = len(_var.UpStreams[streamId].DataList[fileId])
			}
			if sendLen >= len(_var.UpStreams[streamId].DataList[fileId])-dataSendLen[streamId] {
				// 尾数据
				sendLen = len(_var.UpStreams[streamId].DataList[fileId]) - dataSendLen[streamId]
				upStatus = protocol.MetaData_END
				upData = _var.UpStreams[streamId].DataList[fileId][dataSendLen[streamId]:]
				indexs = append(indexs[:streamId], indexs[streamId+1:]...)
				dataSendLen = append(dataSendLen[:streamId], dataSendLen[streamId+1:]...)
			} else {
				// 正常数据
				upData = _var.UpStreams[streamId].DataList[fileId][dataSendLen[streamId] : dataSendLen[streamId]+sendLen]
				dataSendLen[streamId] += sendLen
			}

			desc := make(map[string][]byte)
			for dk, dv := range _var.UpStreams[streamId].DataDesc {
				desc[dk] = []byte(dv)
			}
			inputmeta := protocol.MetaData{strconv.Itoa(streamId), uint32(dataId),
				_var.UpStreams[streamId].DataType,
				upStatus, _var.UpStreams[streamId].DataFmt,
				_var.UpStreams[streamId].DataEnc, upData, desc}

			dataIn.DataList = append(dataIn.DataList, &inputmeta)
		}

		input, err := proto.Marshal(&dataIn)
		if err != nil {
			cli.Log.Errorw("multiUpStream marshal create request fail", "err", err.Error(), "params", dataIn.EngParam)
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
		resp, ecode, err := caller.SessionCall(xsfcli.CONTINUE, _var.SvcName, "AIIn", req, time.Duration(_var.TimeOut+_var.LossDeviation)*time.Millisecond)
		if err != nil && ecode != frame.AigesErrorEngInactive {
			cli.Log.Errorw("multiUpStream Create request fail", "err", err.Error(), "code", ecode, "params", dataIn.EngParam)
			unBlockChanWrite(errchan, struct {
				code int
				err  error
			}{int(ecode), err})
			return
		}
		// 下行结果输出
		dataOut := protocol.EngOutputData{}
		err = proto.Unmarshal(resp.GetData()[0].Data, &dataOut)
		if err != nil {
			cli.Log.Errorw("multiUpStream Resp Unmarshal fail", "err", err.Error(), "respData", resp.GetData()[0].Data)
			unBlockChanWrite(errchan, struct {
				code int
				err  error
			}{-1, err})
			return
		}

		switch dataOut.Ret {
		case 0: // nothing to do
		case frame.AigesErrorEngInactive:
			return
		default:
			cli.Log.Errorw("multiUpStream get engine err", "err", dataOut.Err, "code", dataOut.Ret, "params", dataIn.EngParam)
			unBlockChanWrite(errchan, struct {
				code int
				err  error
			}{int(dataOut.Ret), errors.New(dataOut.Err)})
			return // engine err but not 10101
		}

		// 同步下行数据
		if len(dataOut.DataList) > 0 {
			(*sm).Lock()
			*pm = append(*pm, dataOut)
			cli.Log.Debugw("multiUpStream get resp result", "hdl", session, "result", dataOut)
			(*sm).Unlock()
		}
		if dataOut.Status == protocol.EngOutputData_END {
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
func sessAIOut(cli *xsfcli.Client, hdl string, sid string, rslt *[]protocol.EngOutputData) (code int, err error) {
	// loop read downstream result
	for {
		req := xsfcli.NewReq()
		req.SetParam("baseId", "0")
		req.SetParam("waitTime", strconv.Itoa(_var.TimeOut))
		dataIn := protocol.EngInputData{}
		dataIn.EngParam = make(map[string]string)
		dataIn.EngParam["sid"] = sid

		input, err := proto.Marshal(&dataIn)
		if err != nil {
			cli.Log.Errorw("sessAIOut marshal create request fail", "err", err.Error(), "params", dataIn.EngParam)
			return -1, err
		}

		rd := xsfcli.NewData()
		rd.Append(input)
		req.AppendData(rd)
		_ = req.Session(hdl)

		caller := xsfcli.NewCaller(cli)
		resp, ecode, err := caller.SessionCall(xsfcli.CONTINUE, _var.SvcName, "AIOut", req, time.Duration(_var.TimeOut+_var.LossDeviation)*time.Millisecond)
		if err != nil {
			cli.Log.Errorw("sessAIOut request fail", "err", err.Error(), "code", ecode, "params", dataIn.EngParam)
			if ecode == frame.AigesErrorEngInactive { // reset 10101 inactive
				err = nil
			}
			return int(ecode), err
		}

		// 解析结果、输出落盘
		dataOut := protocol.EngOutputData{}
		err = proto.Unmarshal(resp.GetData()[0].Data, &dataOut)
		if err != nil {
			cli.Log.Errorw("sessAIOut Resp Unmarshal fail", "err", err.Error(), "respData", resp.GetData()[0].Data)
			return -1, err
		}

		*rslt = append(*rslt, dataOut)
		cli.Log.Debugw("sessAIOut get resp result", "hdl", sid, "result", dataOut)
		if dataOut.Status == protocol.EngOutputData_END {
			return code, err // last result
		}
	}

	return
}

func sessAIExcp(cli *xsfcli.Client, hdl string, sid string) (err error) {

	req := xsfcli.NewReq()
	req.SetParam("baseId", "0")
	req.SetParam("waitTime", strconv.Itoa(_var.TimeOut))
	dataIn := protocol.EngInputData{}
	dataIn.EngParam = make(map[string]string)
	dataIn.EngParam["sid"] = sid

	input, err := proto.Marshal(&dataIn)
	if err != nil {
		cli.Log.Errorw("sessAIExcp marshal create request fail", "err", err.Error(), "params", dataIn.EngParam)
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
