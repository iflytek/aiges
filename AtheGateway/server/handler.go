package server

import (
	"common"
	"conf"
	"git.xfyun.cn/AIaaS/xsf-external/utils"

	"git.xfyun.cn/AIaaS/json_script"
	"schemas"
	"github.com/golang/protobuf/proto"
	"sync/atomic"
)

const (
	FRAME_BEGIN    =  0
	FRAME_CONTINUE = 1
	FRAME_END =  2
	FRAME_CLOSE = 4
)

const (
	RESULT_STATUS_RECEIVING int32 = 1
	RESULT_STATUS_END       int32 = 2
)

const (
	BCLUSTER = "5s"
	ReqParamKey = "$req"
)

var Sidgenerator *utils.SidGenerator2


func handleFrame(s *Session, frameReq *FrameReq) {
	if s.checkTimeOut() {
		common.Logger.Errorf("%s:client(%s)send frame(%s) ,session time expired", s.Sid, s.ClientIp, s.SeqNo)
		err := NewError(int(ErrorCodeTimeOut), "session timeout", s.Sid, s.SeqNo)
		s.WriteError(err.GetErrorResp(s.Uid))
		s.Close()
		return
	}
	//closeStatus >=1 说明session 正在关闭中
	if atomic.LoadInt32(&s.closeStatus) >= closeStatusClosing{
		//s.WriteError(NewError(int(ErrorCodeSessionEnd),"call on closed session",s.Sid,0).GetErrorResp(s.Wscid))
		return
	}
	common.Logger.Infof("receive a frame,wscid=%s,sid=%s,status=%d", frameReq.Common.Wscid, s.Sid, s.Status)

	//每一帧都会执行的脚本
	if conf.Conf.Server.ScriptEnable{
		sp:=s.Mapping.RequestData.Script
		if sp !=nil{
			ctx,err:=schemas.ExecuteScript(sp.GetEvery(),frameReq.object)
			if err !=nil && ctx !=nil{
				if req,ok:=ctx.Get(ReqParamKey).(map[string]interface{});ok{
					s.reqParam =  req
				}
			}
		}
	}


	//校验请求参数状态和AppId
	if s.Status == FRAME_BEGIN { //后续的逻辑保证了此处只会被执行一次且必定执行
		//回话开始统计一次，
		if frameReq.Common.AppId != "" {
			s.Appid = frameReq.Common.AppId
		}
		if frameReq.Common.Uid != "" {
			s.Uid = frameReq.Common.Uid
		}
		//执行参数处理脚本脚本
		if conf.Conf.Server.ScriptEnable{
			sp:=s.Mapping.RequestData.Script
			if sp !=nil{
				ctx,err:=schemas.ExecuteScript(sp.GetFirst(),frameReq.object)
				if err,ok:=jsonscpt.IsExitError(err);ok{
					if err.Code!=0{
						s.WriteError(NewError(err.Code,err.Message,s.Sid,1).GetErrorResp(s.Wscid))
						s.Close()
						return
					}
				}
				if err !=nil && ctx !=nil{
					if req,ok:=ctx.Get(ReqParamKey).(map[string]interface{});ok{
						s.reqParam =  req
					}
				}
			}
		}
		//第一帧数据
		if err := CheckAppid(s, s.Appid); err != nil {
			common.Logger.Errorf("invalid appid:%s,sid=%s", s.Appid, s.Sid)
			s.writeJson(err.GetErrorResp(s.Wscid))
			s.setError(err.Code)
			s.Close()
			return
		}
	}

	var err *Error
	//是否mock
	if conf.Conf.Server.Mock {
		err = sendMockData(s, frameReq)
	} else {
		//根据状态向atmos发送数据
		err = sendFrameByStatus(frameReq, s)
	}

	if err != nil {
		common.Logger.Errorf("%s:send data to atmos by xsf failed:%s", err.Sid, err.Msg)
		s.setError(err.Code)
		s.WriteError(err.GetErrorResp(s.Wscid))
		if s.Status == FRAME_BEGIN{
			s.Close()
		}


	}

}

//根据状态向atmos发送数据
func sendFrameByStatus(frameReq *FrameReq, s *Session) *Error {
	//组装serverbiz
	serverBiz, err := frameReq.GetServerBiz(s)
	if err != nil {
		common.Logger.Errorf("getServerBiz error sid=%s,error=%s", s.Sid, err.Error())
		return err
	}
	//获取status
	status, error := frameReq.GetStatusBybiz(serverBiz)
	if error != nil {
		common.Logger.Errorf("getStatusByBiz error sid=%s,error=%v", s.Sid, error)
		return NewErrorByError(ErrorCodeGetUpCall, error, s.Sid, s.SeqNo)
	}

	//收到最后一帧数据
	if status == FRAME_END {
		s.Status = FRAME_END
	}

	//设置serverbiz中得status
	frameReq.SetBizStatus(s, serverBiz)
	//获取biz参数
	//serverBiz, err := frameReq.GetServerBiz(s)
	common.Logger.Infof("send frame:sid= %s seqNo=%d status=%d,busi=%s", s.Sid, s.SeqNo, s.Status, common.MapToString(frameReq.Business))

	var upr *UpResult

	switch s.Status {
	case FRAME_BEGIN: //处理第一帧
		//向atmos发送请求
		upr, err = SendRequest(s, s.Sid, serverBiz, s.SeqNo)
		if err != nil {
			return err
		}else{// 发送第一帧，带sid
			if conf.Conf.Server.WriteFirst{
				s.writeSuccess(&FrameResponse{Sid: s.Sid,Code:0,Wscid:s.Wscid,Message:"success"})
			}
		}
		s.Status = FRAME_CONTINUE

		//writeFirstFrameResp(s, upr)
		HandlerUpResult(s, upr)
		s.sessionMap = upr.Session

	case FRAME_CONTINUE: //处理中间的帧
		upr, err = SendRequest(s, s.Sid, serverBiz, s.SeqNo)
		if err != nil {
			return err
		}
		s.sessionMap = upr.Session
		HandlerUpResult(s, upr)

	case FRAME_END: //处理最后一帧
		upr, err = SendRequest(s, s.Sid, serverBiz, s.SeqNo)
		if err != nil {
			return err
		}
		s.sessionMap = upr.Session
		HandlerUpResult(s, upr)

		//回话结束，释放
	}
	s.SeqNo++
	return nil
}

//处理上行请求得数据，如果有则返回给用户，同步调用时有用
func HandlerUpResult(s *Session, upr *UpResult) {
	resp, err := NewRespByUpResult(s, upr)
	if err != nil {
		s.writeSuccess(err.GetErrorResp(s.Wscid))
		return
	}
	if resp != nil {
		s.writeSuccess(resp)
	}
}



//mock使用
func sendMockData(s *Session, frameReq *FrameReq) *Error {
	serverBiz, err := frameReq.GetServerBiz(s)
	if err != nil {
		return err
	}
	status, error := frameReq.GetStatusBybiz(serverBiz)
	if error != nil {
		return NewErrorByError(ErrorCodeGetUpCall, error, s.Sid, s.SeqNo)
	}
	if status == FRAME_END {
		s.Status = FRAME_END
	}
	frameReq.SetBizStatus(s, serverBiz)
	if s.Status == FRAME_BEGIN {
		//writeFirstFrameResp(s, &UpResult{})
		//resp, err := NewFirstFrameRespByUpResult(s, &UpResult{})
		//if err != nil {
		//	return err
		//}
		//s.writeJson(resp)
		s.Status = FRAME_CONTINUE
	}

	if status == FRAME_CLOSE {
		s.closeAtOnce()
		return nil
	}


	if s.closeStatus == 1{
		s.WriteError(NewError(int(ErrorCodeSessionEnd),"call on closed session",s.Sid,0).GetErrorResp(s.Wscid))
		return nil
	}

	data,_:=proto.Marshal(serverBiz)

	proto.Unmarshal(data,serverBiz)


	if s.Status == FRAME_CONTINUE {
		if s.SeqNo%10 == 0 {
			msg := &myMessage{
				msgNo: uint64(s.SeqNo / 10),
				body: &FrameResponse{
					Sid:s.Sid,
					Code:    0,
					Message: "success",
					Wscid:s.Wscid,
					Data: map[string]interface{}{
						"status": 1,
						"result": map[string]interface{}{
							"bg":  0,
							"ed":  0,
							"ls":  false,
							"pgs": "apd",
							"sn":  s.SeqNo / 10,
							"ws": []map[string]interface{}{
								{
									"bg": 0,
									"cw": []map[string]interface{}{
										{
											"sc": 0,
											"w":  "哈喽",
										},
									},
								},
								{
									"bg": 0,
									"cw": []map[string]interface{}{
										{
											"sc": 0,
											"w":  "哈喽",
										},
									},
								},
							},
						},
					},
				},
				Session: s,
			}
			proto.Unmarshal(data,serverBiz)
			writeMock(s, msg, false)
		}
	}

	if s.Status == FRAME_END {
		msg := &myMessage{
			msgNo: uint64(s.SeqNo),
			body: &FrameResponse{
				Code:    0,
				Message: "success",
				Wscid:s.Wscid,
				Data: map[string]interface{}{
					"status": 2,
					"result": map[string]interface{}{
						"bg":  0,
						"ed":  0,
						"ls":  false,
						"pgs": "apd",
						"sn":  s.SeqNo,
						"ws": []map[string]interface{}{
							{
								"bg": 0,
								"cw": []map[string]interface{}{
									{
										"sc": 0,
										"w":  "哈喽1234。",
									},
								},
							},
						},
					},
				},
			},
			Session: s,
		}
		proto.Unmarshal(data,serverBiz)
		writeMock(s, msg, true)
		s.Close()
	}
	s.SeqNo++
	return nil
}

func writeMock(s *Session, msg Message, islast bool) {
	msg.Send()
}
