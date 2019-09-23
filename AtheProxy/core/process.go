package core

import (
	"encoding/json"
	"time"
	"git.xfyun.cn/AIaaS/xsf-external/client"

	"client"
	"component"
	"config"
	"consts"
	"protocol/biz"
	"protocol/engine"
	"util"
	"runtime/debug"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"errors"
	"github.com/golang/protobuf/proto"
	"component/mock"
)

var Core = &core{}

type core struct {
}

//执行过程
func (c *core) Process(op string, data []byte, span *xsf.Span, param map[string]string) []byte {
	span.Start()
	serverBiz := &serverbiz.ServerBiz{}
	err := proto.Unmarshal(data, serverBiz)
	if err != nil {
		defer func() {
			span.End().Flush()
		}()
		util.SugarLog.Errorw("exception", "err", err)
		return SetUpResult(serverBiz, nil, nil, int32(consts.MSP_ERROR_MSG_PARAM_ERROR), err)
	}

	if util.IsNeedShowDebugLog() {
		util.SugarLog.Debugw("Process req:", "serverBiz", serverBiz.String(), "BusinessArgs", serverBiz.UpCall.BusinessArgs)
	}

	//获取全局异常
	defer func() {
		if v := recover(); v != nil {
			util.SugarLog.Errorw("exception", "err", v, "detail", string(debug.Stack()))
		}
	}()

	if config.SUB_MAP[serverBiz.UpCall.Call] == nil {
		defer func() {
			span.End().Flush()
		}()
		content := SetUpResult(serverBiz, nil, nil, consts.MSP_ERROR_MSG_INVALID_SUBJECT, errors.New("sub:"+serverBiz.UpCall.Call+" is invalid"))
		util.SugarLog.Errorw("sub error", "err", "sub:"+serverBiz.UpCall.Call+" is invalid")
		return content
	}

	if serverBiz.UpCall.SeqNo == 1 {
		//如果预处理返回的是处理完成的状态，则不进行后续的处理
		isFinish, rst := preprocessing(config.SUB_MAP[serverBiz.UpCall.Call], serverBiz, span)
		if isFinish {
			return rst
		}
	}

	//1 为mock，0为正常模式
	if config.SUB_MAP[serverBiz.UpCall.Call].EnableMock == consts.OPEN {
		return mock.MockEngine(serverBiz, span, param)
	} else {
		//引擎
		return engine(op, config.SUB_MAP[serverBiz.UpCall.Call], serverBiz, span, param)
	}
}

/**
调用引擎前的预处理
*/
func preprocessing(configValue *config.AtmosConfig, serverBiz *serverbiz.ServerBiz, span *xsf.Span) (isFinish bool, result []byte) {
	var businessArgs map[string]string
	if serverBiz != nil && serverBiz.UpCall != nil {
		businessArgs = serverBiz.UpCall.BusinessArgs
	} else {
		businessArgs = make(map[string]string)
		businessArgs["upcall"] = "upcall is nil"
	}

	//是否需要结束本次请求
	isFinish = false

	//检查参数
	ret := component.Verify.CheckParams(serverBiz)
	if ret != int32(consts.MSP_SUCCESS) {
		defer func() {
			span.End().Flush()
		}()
		isFinish = true
		return isFinish, SetUpResult(serverBiz, nil, nil, int32(ret), nil)
	}

	//设置兜底ent
	if serverBiz.UpCall.BusinessArgs[consts.ATMOS_ROUTE] == "" {
		util.SugarLog.Errorw("route error,set default route", "SessionId", serverBiz.GlobalRoute.SessionId, "DefaultRoute", configValue.DefaultRoute)
		serverBiz.UpCall.BusinessArgs[consts.ATMOS_ROUTE] = configValue.DefaultRoute
	}

	return isFinish, nil
}

//引擎
func engine(op string, configValue *config.AtmosConfig, serverBiz *serverbiz.ServerBiz, span *xsf.Span, param map[string]string) []byte {

	status, err := component.Engine.GetStatus(serverBiz)
	if err != nil {
		return SetUpResult(serverBiz, nil, nil, consts.MSP_ERROR_MSG_PARAM_ERROR, err)
	}

	//同步
	if serverBiz.UpCall.Sync {
		defer func() {
			span.End().Flush()
		}()
		span.WithTag("Sync", "true")

		//判断是否是异常模式
		if consts.OP_EXCEPTION == op {
			outParam, session, code, err := component.Engine.OnException(configValue, serverBiz, span)
			return SetUpResult(serverBiz, outParam, session, code, err)
		}

		switch status {
		case serverbiz.GeneralData_ONCE:
			outParam, session, code, err := component.Engine.OnOnce(configValue, serverBiz, span)
			return SetUpResult(serverBiz, outParam, session, code, err)
		case serverbiz.GeneralData_BEGIN:
			outParam, session, code, err := component.Engine.OnBegin(configValue, serverBiz, span)
			return SetUpResult(serverBiz, outParam, session, code, err)
		case serverbiz.GeneralData_CONTINUE:
			outParam, session, _, code, err := component.Engine.OnContinue(configValue, serverBiz, span)
			return SetUpResult(serverBiz, outParam, session, code, err)
		case serverbiz.GeneralData_END:
			outParam, session, _, code, err := component.Engine.OnEnd(configValue, serverBiz, span, consts.ENT_OP_IN, "")
			if code != consts.MSP_SUCCESS {
				return SetUpResult(serverBiz, nil, session, code, err)
			}
			return SetUpResult(serverBiz, outParam, session, consts.MSP_SUCCESS, nil)
		default:
			return SetUpResult(serverBiz, nil, nil, consts.MSP_ERROR_MSG_PARAM_ERROR, nil)
		}
		//异步
	} else {
		span.WithTag("Sync", "false")

		if param == nil {
			param = make(map[string]string)
		}

		if status == serverbiz.GeneralData_BEGIN {
			defer func() {
				span.End().Flush()
			}()
			//begin
			outParam, session, code, err := component.Engine.OnBegin(configValue, serverBiz, span)
			//ssb同步返回结果，不需要异步
			return SetUpResult(serverBiz, outParam, session, code, err)
		} else {
			newSpan := span
			go func() {
				defer func() {
					span.End().Flush()
					//获取全局异常
					if v := recover(); v != nil {
						util.SugarLog.Errorw("panic", "err", v, "detail", string(debug.Stack()))
					}
				}()

				//判断是否是异常模式
				if consts.OP_EXCEPTION == op {
					component.Engine.OnException(configValue, serverBiz, span)
					return
				}

				switch status {
				//once模式
				case serverbiz.GeneralData_ONCE:
					outParam, _, code, err := component.Engine.OnOnce(configValue, serverBiz, span)
					if err != nil {
						o, _ := json.Marshal(err)
						param[consts.ERROR_INFO] = string(o)
					}
					if code != consts.MSP_SUCCESS {
						param[consts.TASK_STATUS] = consts.STATUS_FAILED
					} else {
						param[consts.TASK_STATUS] = consts.STATUS_SUCCEED
					}
					Callback(param, serverBiz, outParam, span, code)
					//发送数据
				case serverbiz.GeneralData_CONTINUE:
					//continue
					outParam, _, isEngineEnd, code, err := component.Engine.OnContinue(configValue, serverBiz, newSpan)
					if code != consts.MSP_SUCCESS {
						if err != nil {
							param[consts.ERROR_INFO] = err.Error()
						}
						param[consts.TASK_STATUS] = consts.STATUS_FAILED
						Callback(param, serverBiz, outParam, newSpan, code)
						return
					}

					//引擎没有返回数据直接return
					if outParam == nil || len(outParam) == 0 {
						return
					}
					if isEngineEnd {
						param[consts.TASK_STATUS] = consts.STATUS_SUCCEED
					} else {
						param[consts.TASK_STATUS] = consts.STATUS_IN_PROCESS
					}
					Callback(param, serverBiz, outParam, newSpan, code)
					return
					//获取结果
				case serverbiz.GeneralData_END:
					//end
					op := consts.ENT_OP_IN
					engineHandle := ""
					for {
						outParam, session, isEnd, code, err := component.Engine.OnEnd(configValue, serverBiz, newSpan, op, engineHandle)
						op = consts.ENT_OP_OUT
						if err != nil {
							param[consts.ERROR_INFO] = err.Error()
						}
						if code != consts.MSP_SUCCESS {
							param[consts.TASK_STATUS] = consts.STATUS_FAILED
							Callback(param, serverBiz, outParam, newSpan, code)
							return
						}

						if isEnd {
							param[consts.TASK_STATUS] = consts.STATUS_SUCCEED
							Callback(param, serverBiz, outParam, newSpan, code)
							return
						}
						if outParam != nil && len(outParam) > 0 {
							param[consts.TASK_STATUS] = consts.STATUS_IN_PROCESS
							Callback(param, serverBiz, outParam, newSpan, code)
						}
						engineHandle = getHandle(session)
						newSpan = span.Next(utils.SrvSpan);
						newSpan.End().Flush()
					}
				default:
				}
			}()
		}
		return SetUpResult(serverBiz, nil, serverBiz.UpCall.Session, consts.MSP_SUCCESS, nil)
	}
}

//回调
func Callback(param map[string]string, serverBiz *serverbiz.ServerBiz, outParam []*enginebiz.MetaData, span *xsf.Span, code int32) {
	var addr string
	switch serverBiz.UpCall.From {
	//如果来自guider则取GuiderId
	case consts.SOURCE_GUIDER:
		addr = serverBiz.GlobalRoute.GetGuiderId()
	default:
		addr = serverBiz.GlobalRoute.GetUpRouterId()
	}

	//准备请求数据
	req := xsf.NewReq()
	req.SetTraceID(span.Meta())
	req.Append(setDownCall(serverBiz, outParam, span, code), nil)

	//编排需要把请求参数回设
	if param != nil {
		for k, v := range param {
			req.SetParam(k, v)
		}
	}
	if util.IsNeedShowInfoLog() {
		util.SugarLog.Infow("callback param", "param", param, "addr", addr, "SessionId", serverBiz.GlobalRoute.SessionId)
	}

	_, errcode, err := client.GetRpcClient().CreateCaller().CallWithAddr(serverBiz.UpCall.From, "resp", addr, req, time.Duration(config.TIMEOUT)*time.Millisecond)
	if err != nil {
		util.SugarLog.Errorw("callback", "addr", addr, "errcode:", errcode, "err", err, "SessionId", serverBiz.GlobalRoute.SessionId)
	}
}

//获取 handle
func getHandle(m map[string]string) string {
	if v, ok := m[consts.SESSION_HANDLE]; ok {
		return v
	}
	return ""
}
