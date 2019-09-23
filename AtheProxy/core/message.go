package core

import (
	"github.com/golang/protobuf/proto"

	"util"
	"protocol/biz"
	"protocol/engine"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"strconv"
	"consts"
)

//设置up result
func SetUpResult(serverBiz *serverbiz.ServerBiz, outParam []*enginebiz.MetaData, session map[string]string, code int32, msg error) ([]byte) {
	newServerBiz := &serverbiz.ServerBiz{}
	//设置msg_type
	newServerBiz.MsgType = serverbiz.ServerBiz_UP_RESULT

	var seqNo int32
	var from string
	if serverBiz != nil {
		//version
		newServerBiz.Version = serverBiz.Version
		//设置globalRoute
		newServerBiz.GlobalRoute = serverBiz.GlobalRoute

		upCall := serverBiz.UpCall
		if upCall != nil {
			seqNo = upCall.SeqNo
			from = upCall.From
		}

		if session == nil {
			session = make(map[string]string)
		}

		//设置数据
		dataList := createData(outParam)

		//设置upresult
		if msg == nil {
			msg = consts.ERR
		}

		//设置upresult
		newServerBiz.UpResult = &serverbiz.UpResult{
			Ret:      code,
			AckNo:    seqNo,
			ErrInfo:  msg.Error(),
			From:     from,
			Session:  session,
			DataList: dataList,
		}

		if util.IsNeedShowDebugLog() {
			util.SugarLog.Debugw("SetUpResult", "session_id", serverBiz.GlobalRoute.SessionId, "result serverBiz", newServerBiz.String())
		}

	} else {
		util.SugarLog.Errorw("serverBiz is nil", "code", code, "session", session)

		//设置upresult
		newServerBiz.UpResult = &serverbiz.UpResult{
			Ret:      code,
			AckNo:    seqNo,
			ErrInfo:  msg.Error(),
			From:     from,
			Session:  session,
			DataList: nil,
		}
	}

	//输出响应参数
	content, err := proto.Marshal(newServerBiz)
	if err != nil && newServerBiz.GlobalRoute != nil {
		util.SugarLog.Errorw("up result response proto.Marshal fail", "sessionId", newServerBiz.GlobalRoute.SessionId)
	}
	return content
}

//设置down call
func setDownCall(serverBiz *serverbiz.ServerBiz, outParam []*enginebiz.MetaData, span *xsf.Span, code int32) ([]byte) {
	newServerBiz := &serverbiz.ServerBiz{}
	//设置msg_type和version
	newServerBiz.MsgType = serverbiz.ServerBiz_DOWN_CALL

	var seqNo int32
	var from string
	var sub string
	if serverBiz != nil {
		//version
		newServerBiz.Version = serverBiz.Version
		//设置globalRoute
		newServerBiz.GlobalRoute = serverBiz.GlobalRoute
		if serverBiz.UpCall != nil {
			seqNo = serverBiz.UpCall.SeqNo
			from = serverBiz.UpCall.From
			sub = serverBiz.UpCall.Call
		}
	}

	//设置数据
	dataList := createDownCallData(sub, outParam)

	//设置session map
	newServerBiz.DownCall = &serverbiz.DownCall{
		Ret:      code,
		SeqNo:    seqNo,
		From:     from,
		DataList: dataList,
		//Args: args,
	}

	if util.IsNeedShowDebugLog() {
		util.SugarLog.Debugw("down call response", "session_id", serverBiz.GlobalRoute.SessionId, "result serverBiz", newServerBiz.String())
	}
	span.WithTag("desc", "down call response").
		WithTag("seqNo", strconv.FormatInt(int64(newServerBiz.DownCall.GetSeqNo()), 10))

	//输出响应参数
	content, _ := proto.Marshal(newServerBiz)
	return content
}

//创建数据
func createData(outParam []*enginebiz.MetaData) []*serverbiz.GeneralData {
	dataList := make([]*serverbiz.GeneralData, 0, len(outParam))
	for _, v := range outParam {
		gd := &serverbiz.GeneralData{
			DataId:   v.DataId,
			FrameId:  v.FrameId,
			DataType: serverbiz.GeneralData_DataType(v.DataType),
			Status:   serverbiz.GeneralData_DataStatus(v.Status),
			DescArgs: v.Desc,
			Format:   v.Format,
			Encoding: v.Encoding,
			Data:     v.Data,
		}
		dataList = append(dataList, gd)
	}
	return dataList
}

func createDownCallData(sub string, outParam []*enginebiz.MetaData) []*serverbiz.GeneralData {
	dataList := make([]*serverbiz.GeneralData, 0, len(outParam))
	for _, v := range outParam {
		if v != nil {
			if v.Desc != nil {
				v.Desc["sub"] = []byte(sub)
			} else {
				v.Desc = make(map[string][]byte)
				v.Desc["sub"] = []byte(sub)
			}
		}
		gd := &serverbiz.GeneralData{
			DataId:   v.DataId,
			FrameId:  v.FrameId,
			DataType: serverbiz.GeneralData_DataType(v.DataType),
			Status:   serverbiz.GeneralData_DataStatus(v.Status),
			DescArgs: v.Desc,
			Format:   v.Format,
			Encoding: v.Encoding,
			Data:     v.Data,
		}
		dataList = append(dataList, gd)
	}
	return dataList
}