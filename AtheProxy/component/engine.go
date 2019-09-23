package component

import (
	"errors"
	"strconv"
	"time"

	"client"
	"config"
	"consts"
	"protocol/biz"
	"protocol/engine"
	"util"

	"git.xfyun.cn/AIaaS/xsf-external/client"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"github.com/golang/protobuf/proto"
	"fmt"
)

var Engine = &engine{}

type engine struct {
}

//创建client对象
func createClient(configValue *config.AtmosConfig) *xsf.Caller {
	//创建client
	var c = client.GetRpcClient().CreateCaller()
	c.WithLBParams(configValue.Lb, configValue.Sub, nil)
	c.WithRetry(configValue.EngineRetry)
	return c
}

//初始化参数
func initParams(serverBiz *serverbiz.ServerBiz) (map[string]string, string) {
	//获取session
	newSession := make(map[string]string)
	if serverBiz.UpCall.Session != nil {
		for k, v := range serverBiz.UpCall.Session {
			newSession[k] = v
		}
	}

	//获取路由参数
	route := serverBiz.UpCall.BusinessArgs[consts.ATMOS_ROUTE]
	newSession[consts.SESSIONID] = serverBiz.GlobalRoute.SessionId
	return newSession, route
}

func (e *engine) GetStatus(serverBiz *serverbiz.ServerBiz) (serverbiz.GeneralData_DataStatus, error) {
	//遍历数据
	for _, v := range serverBiz.UpCall.DataList {
		return v.Status, nil
	}
	return serverbiz.GeneralData_BEGIN, consts.ERR_PARAM_INVALID
}

func (e *engine) OnException(configValue *config.AtmosConfig, serverBiz *serverbiz.ServerBiz, span *xsf.Span) ([]*enginebiz.MetaData, map[string]string, int32, error) {
	//初始化参数
	session, route := initParams(serverBiz)
	var r []*enginebiz.MetaData
	span.WithName("OnException")
	//获取handle
	handle := getEntHandle(session)
	//准备请求数据
	req := preRequestFullData(configValue, serverBiz.UpCall.BusinessArgs, serverBiz.UpCall.DataList, handle, serverBiz.GlobalRoute.SessionId, strconv.FormatInt(int64(serverBiz.UpCall.SeqNo), 10), route)
	//发送请求
	req.SetTraceID(span.Meta())

	c := createClient(configValue)
	res, code, err := c.SessionCall(xsf.CONTINUE, "continue", consts.ENT_OP_EXP, req, time.Duration(configValue.EngineTimeout)*time.Millisecond)

	if err != nil {
		recordError(span, serverBiz, err, code)
		return nil, session, code, err
	}

	//解析参数
	outParam := resolveParams(res, serverBiz.GlobalRoute.SessionId, "engine OnException", serverBiz.UpCall.SeqNo, span)
	session[consts.SESSION_HANDLE] = res.Session()

	isEnd := isFinished(outParam, span)
	if util.IsNeedShowInfoLog() {
		util.SugarLog.Infow("engine OnException", "SeqNo", serverBiz.UpCall.SeqNo, "session_id", serverBiz.GlobalRoute.SessionId, "isEnd", isEnd)
	}

	if len(outParam.DataList) > 0 {
		for _, v := range outParam.DataList {
			r = append(r, v)
		}
	} else {
		r = addDefaultData(isEnd, outParam.DataList, r)
	}
	code, err = rewriteError(code, outParam, err)
	return r, session, code, err

}

//OnOnce
func (e *engine) OnOnce(configValue *config.AtmosConfig, serverBiz *serverbiz.ServerBiz, span *xsf.Span) ([]*enginebiz.MetaData, map[string]string, int32, error) {
	//初始化参数
	session, route := initParams(serverBiz)
	//遍历数据
	var r []*enginebiz.MetaData
	for _, v := range serverBiz.UpCall.DataList {
		if v.Status == serverbiz.GeneralData_ONCE {
			span.WithName("OnOnce")
			//准备请求数据
			req := preRequestFullData(configValue, serverBiz.UpCall.BusinessArgs, serverBiz.UpCall.DataList, "", serverBiz.GlobalRoute.SessionId, strconv.Itoa(int(serverBiz.UpCall.SeqNo)), route)
			if util.IsNeedShowInfoLog() {
				util.SugarLog.Infow("OnOnce req:", "session_id", serverBiz.GlobalRoute.SessionId, "BusinessArgs", serverBiz.UpCall.BusinessArgs, "serverBiz", serverBiz.String())
			}

			//发送请求
			req.SetTraceID(span.Meta())
			res, code, err := createClient(configValue).SessionCall(xsf.ONESHORT, route, consts.ENT_OP_IN, req, time.Duration(configValue.OnceTimeout)*time.Millisecond)
			span.WithRetTag(strconv.Itoa(int(code)))

			if err != nil {
				span.WithTag("SeqNo", strconv.FormatInt(int64(serverBiz.UpCall.SeqNo), 10)).WithRetTag(strconv.FormatInt(int64(code), 10)).
					WithTag("err", err.Error())
				util.SugarLog.Errorw("engine OnOnce", "session_id", serverBiz.GlobalRoute.SessionId, "code", code, "err", err)
				return r, session, code, err
			}
			session[consts.SESSION_HANDLE] = res.Session()
			//解析参数
			outParam := resolveParams(res, serverBiz.GlobalRoute.SessionId, "engine OnOnce", serverBiz.UpCall.SeqNo, span)
			if len(outParam.DataList) > 0 {
				for _, v := range outParam.DataList {
					r = append(r, v)
				}
			}
			code, err = rewriteError(code, outParam, err)
			return r, session, code, err
		}
	}
	return r, session, consts.MSP_SUCCESS, consts.ERR
}

//记录错误信息
func recordError(span *xsf.Span, serverBiz *serverbiz.ServerBiz, err error, code int32) {
	span.WithTag("SeqNo", strconv.FormatInt(int64(serverBiz.UpCall.SeqNo), 10)).WithRetTag(strconv.FormatInt(int64(code), 10))
	if err != nil {
		span.WithTag("err", err.Error())
	}
}

//错误码重写
func rewriteError(code int32, outParam *enginebiz.EngOutputData, err error) (int32, error) {
	if code == consts.MSP_SUCCESS && outParam.Ret != consts.MSP_SUCCESS {
		code = outParam.Ret
		if outParam.Err != "" {
			err = errors.New(outParam.Err)
		}
	}
	return code, err
}

//begin
func (e *engine) OnBegin(configValue *config.AtmosConfig, serverBiz *serverbiz.ServerBiz, span *xsf.Span) ([]*enginebiz.MetaData, map[string]string, int32, error) {
	//初始化参数
	session, route := initParams(serverBiz)
	//遍历数据
	var r []*enginebiz.MetaData
	if len(serverBiz.UpCall.DataList) == 0 {
		if serverBiz.UpCall.SeqNo == 1 {
			return r, session, consts.MSP_ERROR_INVALID_DATA, errors.New("DataList is nil")
		}
	}

	for _, v := range serverBiz.UpCall.DataList {
		if v.Status == serverbiz.GeneralData_BEGIN {
			span.WithName("OnBegin")
			//准备请求数据
			req := preRequestFullData(configValue, serverBiz.UpCall.BusinessArgs, serverBiz.UpCall.DataList, "", serverBiz.GlobalRoute.SessionId, strconv.Itoa(int(serverBiz.UpCall.SeqNo)), route)

			if util.IsNeedShowInfoLog() {
				util.SugarLog.Infow("OnBegin req:", "session_id", serverBiz.GlobalRoute.SessionId, "Sync", serverBiz.UpCall.Sync, "BusinessArgs", serverBiz.UpCall.BusinessArgs, "serverBiz", serverBiz.String())
			}

			//发送请求
			req.SetTraceID(span.Meta())
			c := createClient(configValue)
			res, code, err := c.SessionCall(xsf.CREATE, route, consts.ENT_OP_IN, req, time.Duration(configValue.EngineTimeout)*time.Millisecond)
			span.WithRetTag(strconv.Itoa(int(code)))
			if err != nil {
				recordError(span, serverBiz, err, code)
				util.SugarLog.Errorw("engine begin", "session_id", serverBiz.GlobalRoute.SessionId, "code", code, "err", err)
				return r, session, code, err
			}

			session[consts.SESSION_HANDLE] = res.Session()
			//解析参数
			outParam := resolveParams(res, serverBiz.GlobalRoute.SessionId, "engine begin", serverBiz.UpCall.SeqNo, span)
			code, err = rewriteError(code, outParam, err)

			if len(outParam.DataList) > 0 {
				for _, v := range outParam.DataList {
					r = append(r, v)
				}
			}
			return r, session, code, err
		}
	}
	return r, session, consts.MSP_SUCCESS, errors.New("")
}

//continue
func (e *engine) OnContinue(configValue *config.AtmosConfig, serverBiz *serverbiz.ServerBiz, span *xsf.Span) ([]*enginebiz.MetaData, map[string]string, bool, int32, error) {
	//初始化参数
	session, route := initParams(serverBiz)

	//遍历数据
	var r []*enginebiz.MetaData
	for _, v := range serverBiz.UpCall.DataList {
		if v.Status == serverbiz.GeneralData_CONTINUE {
			span.WithName("OnContinue")
			//准备请求数据
			req := preRequestFullData(configValue, serverBiz.UpCall.BusinessArgs, serverBiz.UpCall.DataList, getEntHandle(session), serverBiz.GlobalRoute.SessionId, strconv.Itoa(int(serverBiz.UpCall.SeqNo)), route)

			//发送请求
			req.SetTraceID(span.Meta())
			res, code, err := createClient(configValue).SessionCall(xsf.CONTINUE, route, consts.ENT_OP_IN, req, time.Duration(configValue.EngineTimeout)*time.Millisecond)
			span.WithRetTag(strconv.Itoa(int(code)))
			if err != nil {
				span.WithTag("SeqNo", strconv.FormatInt(int64(serverBiz.UpCall.SeqNo), 10)).WithRetTag(strconv.FormatInt(int64(code), 10))
				if err != nil {
					span.WithTag("err", err.Error())
				}
				if needNotifyEngine(code) {
					util.SugarLog.Errorw("engine continue", "SeqNo", serverBiz.UpCall.SeqNo, "session_id", serverBiz.GlobalRoute.SessionId, "handle", session[consts.SESSION_HANDLE], "code", code, "err", err)
					newSpan := span.Next(utils.SrvSpan).WithRetTag(strconv.FormatInt(int64(code), 10)).WithName("OnContinue exception");
					newSpan.End().Flush()
					//通知引擎，会话出错
					req.SetTraceID(newSpan.Meta())
					createClient(configValue).SessionCall(xsf.CONTINUE, route, consts.ENT_OP_EXP, req, time.Duration(configValue.EngineTimeout)*time.Millisecond)
				}
				return r, session, false, code, err
			}

			//解析参数
			outParam := resolveParams(res, serverBiz.GlobalRoute.SessionId, "engine continue", serverBiz.UpCall.SeqNo, span)
			session[consts.SESSION_HANDLE] = res.Session()
			code, err = rewriteError(code, outParam, err)
			isEnd := isFinished(outParam, span)
			if util.IsNeedShowInfoLog() {
				util.SugarLog.Infow("engine continue", "SeqNo", serverBiz.UpCall.SeqNo, "session_id", serverBiz.GlobalRoute.SessionId, "isEnd", isEnd)
			}
			if len(outParam.DataList) > 0 {
				for _, v := range outParam.DataList {
					r = append(r, v)
				}
			} else {
				r = addDefaultData(isEnd, outParam.DataList, r)
			}
			return r, session, isEnd, code, err
		}
	}
	return r, session, false, consts.MSP_SUCCESS, errors.New("")
}

/**
判断是否需要通知引擎，会话出错
 */
func needNotifyEngine(code int32) bool {

	var flag = true
	for _, v := range consts.IGNORE_ERROR_ARRAY {
		if v == code {
			flag = false
		}
	}
	return flag
}

//end
func (e *engine) OnEnd(configValue *config.AtmosConfig, serverBiz *serverbiz.ServerBiz, span *xsf.Span, op string, engineHandle string) ([]*enginebiz.MetaData, map[string]string, bool, int32, error) {
	//初始化参数
	session, finalEnt := initParams(serverBiz)
	span.WithName("OnEnd")
	//获取handle
	handle := getEntHandle(session)

	//适用于第一次请求就是end状态，而且返回结果后，还有后续结果的场景，原始请求中没有handler字段
	if handle == "" {
		handle = engineHandle
	}

	var req *xsf.Req
	//创建请求对象
	op, req = createEndReq(serverBiz.UpCall.DataList, op, req, configValue, serverBiz, handle, serverBiz.GlobalRoute.SessionId, finalEnt)

	var r []*enginebiz.MetaData
	begin := time.Now()

	req.SetTraceID(span.Meta())
	newSpan := span;
	var i int64 = 1
	for {
		if i > 1 {
			newSpan = span.Next(utils.SrvSpan).WithName("OnEnd" + strconv.FormatInt(i-1, 10));
			req.SetTraceID(newSpan.Meta())
		}
		//发送请求
		sessStat := xsf.CONTINUE
		if serverBiz.UpCall.SeqNo == 1 && i == 1 && op == consts.ENT_OP_IN {
			sessStat = xsf.CREATE
		}
		c := createClient(configValue)
		res, code, err := c.SessionCall(sessStat, finalEnt, op, req, time.Duration(configValue.EngineTimeout)*time.Millisecond)
		newSpan.WithRetTag(strconv.Itoa(int(code)))

		if err != nil {
			recordError(newSpan, serverBiz, err, code)
			session[consts.SESSION_HANDLE] = handle
			//通知引擎异常
			util.SugarLog.Errorw("engine end", "SeqNo", serverBiz.UpCall.SeqNo, "session_id", serverBiz.GlobalRoute.SessionId, "handle", handle, "code", code, "err", err)
			tempSpan := newSpan.Next(utils.SrvSpan).WithRetTag(fmt.Sprintf("%v", code)).WithName("exception")
			tempSpan.End().Flush()
			req.SetTraceID(tempSpan.Meta())
			c.SessionCall(xsf.CONTINUE, finalEnt, consts.ENT_OP_EXP, req, time.Duration(configValue.EngineTimeout)*time.Millisecond)
			if i > 1 {
				newSpan.End().Flush()
			}
			return r, session, true, code, err
		}

		//解析参数
		outParam := resolveParams(res, serverBiz.GlobalRoute.SessionId, "engine end"+strconv.FormatInt(i-1, 10), serverBiz.UpCall.SeqNo, newSpan)
		session[consts.SESSION_HANDLE] = res.Session()
		code, err = rewriteError(code, outParam, err)

		if serverBiz.UpCall.SeqNo == 1 && i == 1 && op == consts.ENT_OP_IN {
			op = consts.ENT_OP_OUT
			_, req = createEndReq(serverBiz.UpCall.DataList, op, req, configValue, serverBiz, res.Session(), serverBiz.GlobalRoute.SessionId, finalEnt)
			handle = res.Session()
		}

		//是否为end
		isEnd := isFinished(outParam, newSpan)
		if util.IsNeedShowInfoLog() {
			util.SugarLog.Infow("engine OnEnd", "SeqNo", serverBiz.UpCall.SeqNo, "session_id", serverBiz.GlobalRoute.SessionId, "isEnd", isEnd)
		}

		if len(outParam.DataList) > 0 || isEnd {
			for _, v := range outParam.DataList {
				r = append(r, v)
			}
			r = addDefaultData(isEnd, outParam.DataList, r)
			if i > 1 {
				newSpan.End().Flush()
			}
			return r, session, isEnd, code, err
		}

		if code != consts.MSP_SUCCESS {
			recordError(newSpan, serverBiz, err, code)
			if i > 1 {
				newSpan.End().Flush()
			}
			return r, session, isEnd, code, err
		}

		now := time.Now()
		//获取结果时间大于超时时间，则直接返回超时错误码
		if int(now.Sub(begin).Seconds()) > configValue.GetEngineResultTimeout {
			//通知引擎异常
			util.SugarLog.Errorw("engine end", "SeqNo", serverBiz.UpCall.SeqNo, "session_id", serverBiz.GlobalRoute.SessionId, "handle", handle, "code", consts.SERVER_ERROR_TIME_OUT)
			tempSpan := newSpan.Next(utils.SrvSpan).WithRetTag(fmt.Sprintf("%v", code)).WithName("exception")
			tempSpan.End().Flush()
			req.SetTraceID(tempSpan.Meta())
			c.SessionCall(xsf.CONTINUE, finalEnt, consts.ENT_OP_EXP, req, time.Duration(configValue.EngineTimeout)*time.Millisecond)
			if i > 1 {
				newSpan.End().Flush()
			}
			return r, session, true, consts.SERVER_ERROR_TIME_OUT, consts.TIME_OUT_ERR
		}
		i++
		op = consts.ENT_OP_OUT
	}

	return r, session, false, consts.MSP_SUCCESS, errors.New("")
}

func createEndReq(dt []*serverbiz.GeneralData, op string, req *xsf.Req, configValue *config.AtmosConfig, serverBiz *serverbiz.ServerBiz, handle string, sessionId string, finalEnt string) (string, *xsf.Req) {
	if dt == nil || len(dt) == 0 {
		op = consts.ENT_OP_OUT
		//准备请求数据
		req = preRequestData(configValue, serverBiz.UpCall.BusinessArgs, nil, handle, sessionId, strconv.Itoa(int(serverBiz.UpCall.SeqNo)))
	} else {
		//遍历数据
		for _, v := range dt {
			if v.Status == serverbiz.GeneralData_END {
				//准备请求数据
				req = preRequestFullData(configValue, serverBiz.UpCall.BusinessArgs, dt, handle, sessionId, strconv.Itoa(int(serverBiz.UpCall.SeqNo)), finalEnt)
				break
			}
		}
	}

	if util.IsNeedShowInfoLog() {
		util.SugarLog.Infow("OnEnd req:", "op", op, "handle", req.Handle(), "session_id", sessionId, "content", serverBiz.String(), "BusinessArgs", serverBiz.UpCall.BusinessArgs)
	}

	return op, req
}
func addDefaultData(isEnd bool, data []*enginebiz.MetaData, r []*enginebiz.MetaData) []*enginebiz.MetaData {
	// 如果 end状态，并且 DataList 是空，则构造一个只包括状态字段的datalist
	if isEnd && (data == nil || len(data) == 0) {
		metaData := &enginebiz.MetaData{}
		metaData.Status = enginebiz.MetaData_END
		r = append(r, metaData)
	}
	return r
}

//准备请求数据serverBiz.UpCall.DataList
func preRequestFullData(configValue *config.AtmosConfig, args map[string]string, dt []*serverbiz.GeneralData, handle string, sessionId string, seqNo string, finalEnt string) *xsf.Req {

	newArgs := make(map[string]string)
	for key, value := range args {
		_, has := configValue.EngineSkipParmMap[key]
		if has {
			continue
		}
		newArgs[key] = value
	}
	newArgs[consts.SESSIONID] = sessionId
	newArgs["waitTime"] = strconv.Itoa(int(configValue.EngineTimeout / 2))
	if util.IsNeedShowInfoLog() {
		util.SugarLog.Infow("requestData", "EngParam", newArgs, "SessionId", sessionId)
	}
	req := xsf.NewReq()
	dataList := []*enginebiz.MetaData{}
	for _, gd := range dt {
		if gd != nil {
			dataList = append(dataList, &enginebiz.MetaData{
				gd.DataId,
				gd.FrameId,
				enginebiz.MetaData_DataType(gd.DataType),
				enginebiz.MetaData_DataStatus(gd.Status),
				gd.Format,
				gd.Encoding,
				gd.Data,
				gd.DescArgs})
		}
	}

	//创建输入参数
	in := &enginebiz.EngInputData{
		EngParam: newArgs,
		DataList: dataList,
	}
	//准备请求数据
	data, _ := proto.Marshal(in)
	req.Append(data, nil)

	if handle != "" {
		req.Session(handle)
	}
	req.SetParam("SeqNo", seqNo)
	req.SetParam("waitTime", newArgs["waitTime"])
	if seqNo == "1" {
		//baseId 如果ssb有数据，则设为0 传给引擎，否则传1给引擎
		if dt == nil || len(dt) == 0 {
			req.SetParam("baseId", "1")
		} else {
			req.SetParam("baseId", "0")
		}
	}
	return req
}

//准备请求数据
func preRequestData(configValue *config.AtmosConfig, args map[string]string, gd *serverbiz.GeneralData, handle string, sessionId string, seqNo string) *xsf.Req {

	newArgs := make(map[string]string)
	for key, value := range args {
		_, has := configValue.EngineSkipParmMap[key]
		if has {
			continue
		}
		newArgs[key] = value
	}
	newArgs[consts.SESSIONID] = sessionId
	newArgs["waitTime"] = strconv.Itoa(int(configValue.EngineTimeout / 2))
	if util.IsNeedShowInfoLog() {
		util.SugarLog.Infow("requestData", "EngParam", newArgs, "SessionId", sessionId)
	}
	req := xsf.NewReq()
	if gd != nil {
		//创建输入参数
		in := &enginebiz.EngInputData{
			EngParam: newArgs,
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
		req.Append(data, nil)
	}

	if handle != "" {
		req.Session(handle)
	}
	req.SetParam("SeqNo", seqNo)
	req.SetParam("waitTime", newArgs["waitTime"])
	if seqNo == "1" {
		//baseId 如果ssb有数据，则设为0 传给引擎，否则传1给引擎
		if gd == nil || gd.Data == nil {
			req.SetParam("baseId", "1")
		} else {
			req.SetParam("baseId", "0")
		}
	}
	return req
}

//获取handle
func getEntHandle(m map[string]string) string {
	if v, ok := m[consts.SESSION_HANDLE]; ok {
		return v
	}
	return ""
}

//引擎是否返回end状态
func isFinished(outPut *enginebiz.EngOutputData, span *xsf.Span) bool {
	if outPut.Status == enginebiz.EngOutputData_END {
		span.WithTag("isEnd", "true")
		return true
	}
	span.WithTag("isEnd", "false")
	return false
}

//解析参数
func resolveParams(res *xsf.Res, sessionId string, flow string, seqNo int32, span *xsf.Span) *enginebiz.EngOutputData {
	outPut := &enginebiz.EngOutputData{}
	if res.GetData() != nil && len(res.GetData()) > 0 {
		data := res.GetData()[0].Data
		proto.Unmarshal(data, outPut)
	}

	span.WithTag("desc", flow).
		WithTag("SeqNo", strconv.Itoa(int(seqNo)))

	if outPut.DataList != nil && len(outPut.DataList) > 0 {
		span.WithTag("hasResult", "true")
		if util.IsNeedShowInfoLog() {
			util.SugarLog.Infow(flow+" ,result:", "session_id", sessionId, "handle", res.Session(), "SeqNo", seqNo, "hasResult", true, "Ret", outPut.Ret)
		}
	} else {
		span.WithTag("hasResult", "false")
		if util.IsNeedShowInfoLog() {
			util.SugarLog.Infow(flow+" ,result:", "session_id", sessionId, "handle", res.Session(), "SeqNo", seqNo, "hasResult", false, "Ret", outPut.Ret)
		}
	}
	return outPut
}
