package mock

import (
	"errors"

	"github.com/golang/protobuf/proto"

	"util"
	"protocol/biz"
	"protocol/engine"
	"git.xfyun.cn/AIaaS/xsf-external/server"
)

//设置up result
func SetUpResult(serverBiz *serverbiz.ServerBiz, outParam []*enginebiz.MetaData, session map[string]string, span *xsf.Span, code int32, msg error) ([]byte) {
	newServerBiz := &serverbiz.ServerBiz{}
	//设置msg_type
	newServerBiz.MsgType = serverbiz.ServerBiz_UP_RESULT

	if serverBiz != nil {
		//version
		newServerBiz.Version = serverBiz.Version

		//设置globalRoute
		newServerBiz.GlobalRoute = serverBiz.GlobalRoute
	}

	//设置数据
	dataList := createData(outParam)

	//设置upresult
	if msg == nil {
		msg = errors.New("")
	}

	//BusinessArgs数据追加到session中
	if session == nil {
		session = make(map[string]string)
	}

	//设置upresult
	newServerBiz.UpResult = &serverbiz.UpResult{
		Ret:      code,
		AckNo:    serverBiz.UpCall.SeqNo,
		ErrInfo:  msg.Error(),
		From:      serverBiz.UpCall.From,
		Session:  session,
		DataList: dataList,
	}

	//输出响应参数
	content, err := proto.Marshal(newServerBiz)
	if err != nil {
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
	}

	//设置数据
	dataList := createData(outParam)

	//设置session map
	newServerBiz.DownCall = &serverbiz.DownCall{
		Ret:      code,
		SeqNo:    seqNo,
		From:     from,
		DataList: dataList,
	}

	//输出响应参数
	content, err := proto.Marshal(newServerBiz)
	if err != nil {
		util.SugarLog.Errorw("down call response", "sessionId", newServerBiz.GlobalRoute.SessionId, "err", err)
	}
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
			Format:   v.Format,
			Encoding: v.Encoding,
			Data:     v.Data,
		}
		dataList = append(dataList, gd)
	}
	return dataList
}