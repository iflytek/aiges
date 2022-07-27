package request

import (
	xsfcli "git.iflytek.com/AIaaS/xsf/client"
	"github.com/golang/protobuf/proto"
	"protocol"
	"strconv"
	"time"
	"xtest/analy"
	"xtest/util"
	_var "xtest/var"
)

func (r *Request) OneShotCall(cli *xsfcli.Client, index int64) (info analy.ErrInfo) {
	// request构包, 通过oneShot方式请求AIIn方法.
	sessId := util.NewSid(r.C.TestSub)
	req := xsfcli.NewReq()
	req.SetParam("SeqNo", "1") // 相关协议约定;
	req.SetParam("baseId", "0")
	req.SetParam("version", "v2")
	req.SetParam("waitTime", strconv.Itoa(r.C.TimeOut))
	dataIn := protocol.LoaderInput{}
	dataIn.State = protocol.LoaderInput_ONCE
	dataIn.ServiceId = r.C.SvcId
	dataIn.ServiceName = r.C.SvcName
	// 平台参数header
	dataIn.Headers = make(map[string]string)
	dataIn.Headers["sid"] = sessId
	dataIn.Headers["status"] = "3"
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
	// 上行数据payload
	for _, v := range r.C.UpStreams {
		streamIndex := index % int64(len(v.DataList))
		desc := protocol.MetaDesc{Name: v.Name, DataType: v.DataType}
		desc.Attribute = make(map[string]string)
		for k, v := range v.DataDesc {
			desc.Attribute[k] = v
		}
		desc.Attribute["status"] = "3"
		desc.Attribute["seq"] = "0"
		payload := protocol.Payload{Data: v.DataList[streamIndex], Meta: &desc}
		dataIn.Pl = append(dataIn.Pl, &payload)
	}

	// input data marshal
	input, err := proto.Marshal(&dataIn)
	if err != nil {
		cli.Log.Errorw("OneShotCall marshal fail", "err", err.Error(),
			"header", dataIn.Headers, "params", dataIn.Params)
		return analy.ErrInfo{ErrCode: -1, ErrStr: err}
	}

	rd := xsfcli.NewData()
	rd.Append(input)
	req.AppendData(rd)

	caller := xsfcli.NewCaller(cli)

	analy.Perf.Record(sessId, "", analy.DataTotal, analy.SessOnce, analy.UP, 0, "")
	r.C.ConcurrencyCnt.Add(1) // jbzhou5 启动协程时+1
	resp, code, err := caller.SessionCall(xsfcli.ONESHORT, r.C.SvcName, "AIIn", req, time.Duration(r.C.TimeOut+r.C.LossDeviation)*time.Millisecond)
	r.C.ConcurrencyCnt.Dec() // jbzhou5 任务完成时-1
	if err != nil {
		cli.Log.Errorw("OneShotCall request fail", "err", err.Error(), "code", code,
			"header", dataIn.Headers, "params", dataIn.Params)
		analy.Perf.Record(sessId, "", analy.DataTotal, analy.SessOnce, analy.DOWN, int(code), err.Error())
		return analy.ErrInfo{ErrCode: int(code), ErrStr: err}
	}

	// 解析结果、输出落盘
	dataOut := protocol.LoaderOutput{}
	err = proto.Unmarshal(resp.GetData()[0].Data, &dataOut)
	if err != nil {
		cli.Log.Errorw("OneShotCall Resp Unmarshal fail", "err", err.Error(), "respData", resp.GetData()[0].Data)
		analy.Perf.Record(sessId, "", analy.DataTotal, analy.SessOnce, analy.DOWN, -1, err.Error())
		return analy.ErrInfo{ErrCode: -1, ErrStr: err}
	}
	analy.Perf.Record(sessId, "", analy.DataTotal, analy.SessOnce, analy.DOWN, int(dataOut.Code), dataOut.Err)
	// get result
	for _, v := range dataOut.Pl {
		// 结果输出 & 异步写channel失败则同步写入;
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
		case r.C.AsyncDrop <- _var.OutputMeta{v.Meta.Name, sessId,
			outType, v.Meta.Attribute, v.Data}:
		default:
			// 异步channel满, 同步写;	key: sid-type-name, value: data
			key := sessId + "-" + outType + "-" + v.Meta.Name
			r.downOutput(key, v.Data, cli.Log)
		}
	}

	return
}
