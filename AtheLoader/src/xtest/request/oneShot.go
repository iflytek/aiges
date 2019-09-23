package request

import (
	xsfcli "git.xfyun.cn/AIaaS/xsf-external/client"
	"github.com/golang/protobuf/proto"
	"protocol"
	"strconv"
	"time"
	"xtest/util"
	_var "xtest/var"
)

func OneShotCall(cli *xsfcli.Client, index int64) (code int, err error) {
	// request构包, 通过oneShot方式请求AIIn方法.
	req := xsfcli.NewReq()
	req.SetParam("SeqNo", "1") // 相关协议约定;
	req.SetParam("baseId", "0")
	req.SetParam("waitTime", strconv.Itoa(_var.TimeOut))
	dataIn := protocol.EngInputData{}
	dataIn.EngParam = make(map[string]string)
	for k, v := range _var.UpParams {
		dataIn.EngParam[k] = v
	}
	reqSid := util.NewSid(_var.TestSub)
	dataIn.EngParam["sid"] = reqSid

	for k, v := range _var.UpStreams {
		streamIndex := index % int64(len(v.DataList))
		desc := make(map[string][]byte)
		for dk, dv := range v.DataDesc {
			desc[dk] = []byte(dv)
		}
		inputmeta := protocol.MetaData{strconv.Itoa(k), 0, v.DataType,
			protocol.MetaData_ONCE, v.DataFmt, v.DataEnc, v.DataList[streamIndex], desc}

		dataIn.DataList = append(dataIn.DataList, &inputmeta)
	}

	// input data marshal
	input, errMsl := proto.Marshal(&dataIn)
	if errMsl != nil {
		cli.Log.Errorw("OneShotCall marshal fail", "err", errMsl.Error(), "params", dataIn.EngParam)
		return -1, errMsl
	}

	rd := xsfcli.NewData()
	rd.Append(input)
	req.AppendData(rd)

	caller := xsfcli.NewCaller(cli)
	resp, ecode, err := caller.SessionCall(xsfcli.ONESHORT, _var.SvcName, "AIIn", req, time.Duration(_var.TimeOut+_var.LossDeviation)*time.Millisecond)
	if err != nil {
		cli.Log.Errorw("OneShotCall request fail", "err", err.Error(), "code", ecode, "params", dataIn.EngParam)
		return int(ecode), err
	}

	// 解析结果、输出落盘
	dataOut := protocol.EngOutputData{}
	err = proto.Unmarshal(resp.GetData()[0].Data, &dataOut)
	if err != nil {
		cli.Log.Errorw("OneShotCall Resp Unmarshal fail", "err", err.Error(), "respData", resp.GetData()[0].Data)
		return -1, err
	}

	// got rec result
	if len(dataOut.DataList) > 0 {
		recResult := dataOut.GetDataList()
		for _, v := range recResult {
			// 结果输出 & 异步写channel失败则同步写入;
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
	}

	return
}
