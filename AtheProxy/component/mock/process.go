package mock

import (
	"protocol/biz"
	"git.xfyun.cn/AIaaS/xsf-external/client"
	"consts"
	"config"
	"util"
	"client"
	"protocol/engine"
	"time"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

//引擎
func MockEngine(serverBiz *serverbiz.ServerBiz, span *xsf.Span, param map[string]string) []byte {

	return AsynchronousSessionMode(serverBiz, span, param);

	//sync := serverBiz.UpCall.Sync
	//
	//util.SugarLog.Debugw("MockEngine", "sync", sync)
	////同步
	//if sync {
	//	//判断是否是会话模式
	//	if Engine.IsOnce(serverBiz) {
	//		util.SugarLog.Debugw("MockEngine", "IsOnce", true)
	//		outParam, session, code, err := Engine.OnOnce(serverBiz, span)
	//		return SetUpResult(serverBiz, outParam, session, span, code, err);
	//	}
	//	//同步会话模式
	//	return syncSessionMode(serverBiz, span);
	//
	//	//异步
	//} else {
	//	//判断是否是会话模式:非会话模式
	//	if Engine.IsOnce(serverBiz) {
	//		go func() {
	//			outParam, _, code, err := Engine.OnOnce(serverBiz, span)
	//			if (err != nil) {
	//				o, _ := json.Marshal(err)
	//				param[consts.ERROR_INFO] = string(o)
	//			}
	//			if code != consts.MSP_SUCCESS {
	//				param[consts.TASK_STATUS] = consts.STATUS_FAILED
	//			} else {
	//				param[consts.TASK_STATUS] = consts.STATUS_SUCCEED
	//			}
	//			Callback(param, serverBiz, outParam, span, code)
	//		}()
	//		return SetUpResult(serverBiz, nil, nil, span, consts.MSP_SUCCESS, nil)
	//	} else {
	//		//会话模式
	//		return asynchronousSessionMode(serverBiz, span, param);
	//	}
	//}
}

/**
 同步会话模式
 */
func syncSessionMode(serverBiz *serverbiz.ServerBiz, span *xsf.Span) []byte {
	//begin
	outParam, session, isContinue, code, err := Engine.OnBegin(serverBiz, span)
	if code != consts.MSP_SUCCESS {
		return SetUpResult(serverBiz, nil, nil, span, code, err)
	}
	if !isContinue {
		return SetUpResult(serverBiz, outParam, session, span, consts.MSP_SUCCESS, nil)
	}

	//continue
	outParam, session, isContinue, code, err = Engine.OnContinue(serverBiz, span)
	if code != consts.MSP_SUCCESS {
		return SetUpResult(serverBiz, nil, nil, span, code, err)
	}
	if !isContinue {
		return SetUpResult(serverBiz, outParam, session, span, consts.MSP_SUCCESS, nil)
	}

	//end
	outParam, session, code, err = Engine.OnEnd(serverBiz, span)
	if code != consts.MSP_SUCCESS {
		return SetUpResult(serverBiz, nil, nil, span, code, err)
	}
	return SetUpResult(serverBiz, outParam, session, span, consts.MSP_SUCCESS, nil)
}

//异步会话模式
func AsynchronousSessionMode(serverBiz *serverbiz.ServerBiz, span *xsf.Span, param map[string]string) []byte {
	if param == nil {
		param = make(map[string]string)
	}
	//begin
	outParam, session, isContinue, code, err := Engine.OnBegin(serverBiz, span)
	if code != consts.MSP_SUCCESS {
		return SetUpResult(serverBiz, nil, nil, span, code, err)
	}
	if !isContinue {
		//go Callback(param, serverBiz, outParam, span, consts.MSP_SUCCESS)
		return SetUpResult(serverBiz, nil, session, span, consts.MSP_SUCCESS, nil)
	}

	newSpan := utils.NewSpanFromMeta(span.Meta(), utils.CliSpan)

	go func() {
		defer newSpan.Flush()
		//continue
		outParam, session, isContinue, code, err = Engine.OnContinue(serverBiz, newSpan)
		if code != consts.MSP_SUCCESS {
			return
		}
		if !isContinue {
			//引擎没有返回数据直接return
			if outParam == nil || len(outParam) == 0 {
				return
			}
			Callback(param, serverBiz, outParam, newSpan, code)
			return
		}

		//end
		outParam, session, code, err = Engine.OnEnd(serverBiz, newSpan)
		if code != consts.MSP_SUCCESS {
			return
		}
		Callback(param, serverBiz, outParam, newSpan, code)
	}()
	return SetUpResult(serverBiz, nil, session, span, consts.MSP_SUCCESS, nil)
}

//回调
func Callback(param map[string]string, serverBiz *serverbiz.ServerBiz, outParam []*enginebiz.MetaData, span *xsf.Span, code int32) {
	from := serverBiz.UpCall.From
	var addr string
	switch from {
	//如果来自guider则取GuiderId
	case consts.SOURCE_GUIDER:
		addr = serverBiz.GlobalRoute.GetGuiderId()
	default:
		addr = serverBiz.GlobalRoute.GetUpRouterId()
	}

	//准备请求数据
	req := xsf.NewReq()
	req.SetTraceID(span.Meta())
	data := setDownCall(serverBiz, outParam, span, code)
	req.Append(data, nil)

	//创建client
	c := client.GetRpcClient().CreateCaller()

	//编排需要把请求参数回设
	if param != nil {
		for k, v := range param {
			req.SetParam(k, v)
		}
	}
	util.SugarLog.Infow("callback param", "param", param, "addr", addr)
	begin := time.Now();
	_, errcode, err := c.CallWithAddr(from, "resp", addr, req, time.Duration(config.TIMEOUT)*time.Millisecond)
	if err != nil {
		util.SugarLog.Errorw("callback", "begin", begin, "end", time.Now(), "addr", addr, "errcode:", errcode, "err", err, "config.TIMEOUT", config.TIMEOUT, "SessionId", serverBiz.GlobalRoute.SessionId)
	}
}
