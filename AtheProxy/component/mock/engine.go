package mock

import (
	"errors"
	"consts"
	"util"
	"protocol/engine"
	"protocol/biz"
	"git.xfyun.cn/AIaaS/xsf-external/client"

	"github.com/golang/protobuf/proto"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"encoding/json"
)

var Engine = &engine{}

type engine struct {
}

//初始化参数
func initParams(serverBiz *serverbiz.ServerBiz) (*xsf.Caller, map[string]string, string, string) {
	//获取session
	newSession := make(map[string]string)
	newSession[consts.SESSIONID] = serverBiz.GlobalRoute.SessionId
	return nil, newSession, serverBiz.GlobalRoute.SessionId, ""
}

/**
判断是否是会话模式:once表示非会话模式
 */
func (e *engine) IsOnce(serverBiz *serverbiz.ServerBiz) (bool) {
	//遍历数据
	dt := serverBiz.UpCall.DataList
	for _, v := range dt {
		status := v.Status
		if status == serverbiz.GeneralData_ONCE {
			return true;
		}
	}
	return false;
}

//OnOnce
func (e *engine) OnOnce(serverBiz *serverbiz.ServerBiz, span *xsf.Span) ([]*enginebiz.MetaData, map[string]string, int32, error) {
	//初始化参数
	_, session, sessionId, _ := initParams(serverBiz)

	//遍历数据
	dt := serverBiz.UpCall.DataList
	var r []*enginebiz.MetaData
	for _, v := range dt {
		status := v.Status
		if status == serverbiz.GeneralData_ONCE {
			//准备请求数据
			req := preRequestData(serverBiz.UpCall.BusinessArgs, v, "", sessionId)
			if util.IsNeedShowInfoLog() {
				o, _ := json.Marshal(req)
				util.SugarLog.Infow("OnOnce req:", "session_id", sessionId, "content", string(o), "BusinessArgs", serverBiz.UpCall.BusinessArgs, "GeneralData", v)
			}
			//发送请求
			var tempSession = utils.Session{"h", sessionId, make(map[string]string)}
			dataMeta := &utils.DataMeta{[]byte("mock data"), make(map[string]string)}
			var temp = utils.ResData{consts.MSP_SUCCESS, "", &tempSession, make(map[string]string), []*utils.DataMeta{dataMeta}}
			var res = utils.NewRes();
			res.SetRes(&temp)
			util.SugarLog.Infow("engine mock OnContinue", "session_id", sessionId, "mock OnContinue", "mock OnContinue")
			session[consts.SESSION_HANDLE] = res.Session()

			//解析参数
			var metaData = enginebiz.MetaData{"1", uint32(serverBiz.UpCall.SeqNo), enginebiz.MetaData_TEXT, enginebiz.MetaData_ONCE, "", "", []byte("mock once data"), nil}
			var mm = []*enginebiz.MetaData{&metaData}
			if len(mm) > 0 {
				for _, v := range mm {
					r = append(r, v)
				}
			}
			return r, session, consts.MSP_SUCCESS, errors.New("")
		}
	}
	return r, session, consts.MSP_SUCCESS, errors.New("")
}

//begin
func (e *engine) OnBegin(serverBiz *serverbiz.ServerBiz, span *xsf.Span) ([]*enginebiz.MetaData, map[string]string, bool, int32, error) {
	session := make(map[string]string)
	session[consts.SESSIONID] = serverBiz.GlobalRoute.SessionId

	//遍历数据
	var r []*enginebiz.MetaData
	for _, v := range serverBiz.UpCall.DataList {
		if v.Status == serverbiz.GeneralData_BEGIN {
			session[consts.SERVICE_NAME] = serverBiz.GlobalRoute.SessionId
			return r, session, false, consts.MSP_SUCCESS, errors.New("")
		}
	}
	return r, session, true, consts.MSP_SUCCESS, errors.New("")
}

//continue
func (e *engine) OnContinue(serverBiz *serverbiz.ServerBiz, span *xsf.Span) ([]*enginebiz.MetaData, map[string]string, bool, int32, error) {
	//遍历数据
	var r []*enginebiz.MetaData
	for _, v := range serverBiz.UpCall.DataList {
		if v.Status == serverbiz.GeneralData_CONTINUE {
			session := make(map[string]string)
			session[consts.SESSIONID] = serverBiz.GlobalRoute.SessionId
			return r, session, false, consts.MSP_SUCCESS, errors.New("")
		}
	}
	return r, nil, true, consts.MSP_SUCCESS, errors.New("")
}

//end
func (e *engine) OnEnd(serverBiz *serverbiz.ServerBiz, span *xsf.Span) ([]*enginebiz.MetaData, map[string]string, int32, error) {
	//遍历数据
	var r []*enginebiz.MetaData
	for _, v := range serverBiz.UpCall.DataList {
		if v.Status == serverbiz.GeneralData_END {
			session := make(map[string]string)
			session[consts.SESSIONID] = serverBiz.GlobalRoute.SessionId
			var metaData = enginebiz.MetaData{"9999", uint32(serverBiz.UpCall.SeqNo), enginebiz.MetaData_AUDIO, enginebiz.MetaData_END, "", "", []byte("mock end data"), nil}
			dd := []*enginebiz.MetaData{&metaData}
			var r []*enginebiz.MetaData
			for _, v := range dd {
				r = append(r, v)
			}
			return r, session, consts.MSP_SUCCESS, errors.New("")
		}
	}
	return r, nil, consts.MSP_SUCCESS, errors.New("")
}

//准备请求数据
func preRequestData(args map[string]string, gd *serverbiz.GeneralData, handle string, sessionId string) (*xsf.Req) {
	//创建输入参数
	in := &enginebiz.EngInputData{
		EngParam: args,
		DataList: []*enginebiz.MetaData{
			{
				DataId:   gd.DataId,
				FrameId:  gd.FrameId,
				Format:   gd.Format,
				Encoding: gd.Encoding,
				DataType: enginebiz.MetaData_DataType(gd.DataType),
				Status:   enginebiz.MetaData_DataStatus(gd.Status),
				Data:     gd.Data,
				Desc:     gd.DescArgs,
			},
		},
	}

	//准备请求数据
	data, _ := proto.Marshal(in)
	req := xsf.NewReq()
	req.Append(data, nil)
	if handle != "" {
		req.Session(handle)
	}
	return req
}


