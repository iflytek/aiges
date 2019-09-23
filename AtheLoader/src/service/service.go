package service

import (
	"catch"
	"codec"
	"conf"
	"fmt"
	"frame"
	xsf "git.xfyun.cn/AIaaS/xsf-external/server"
	"github.com/golang/protobuf/proto"
	"instance"
	"protocol"
	"storage"
	"strconv"
	"utils"
)

type aiService struct {
	tool         *xsf.ToolBox
	instMngr     *instance.Manager
	callbackInit actionInit
	callbackFini actionFini
	callbackUser map[usrEvent]actionUser
}

func (srv *aiService) Init(box *xsf.ToolBox) (err error) {
	srv.tool = box
	srv.instMngr = &instance.Manager{}
	// 框架配置初始化
	if err = conf.Construct(srv.tool.Cfg); err != nil {
		srv.tool.Log.Errorf(err.Error())
		return
	}
	// 设置cpu亲和性
	if err = utils.NumaBind(conf.NumaNode); err != nil {
		srv.tool.Log.Errorf(err.Error())
		return
	}
	// 编解码初始化
	if err = codec.AudioCodecInit(); err != nil {
		srv.tool.Log.Errorf(err.Error())
		return
	}
	// 异步下发模式消息队列
	if err = storage.RabInit(
		conf.RabbitHost,
		conf.RabbitUser,
		conf.RabbitPass,
		conf.RabbitQueue,
		srv.tool.Log); err != nil {
			srv.tool.Log.Errorf(err.Error())
			return
	}
	// 服务实例初始化
	usrDef := make(map[instance.UserEvent]instance.UsrAct)
	for e, a := range srv.callbackUser {
		usrDef[instance.UserEvent(e)] = instance.UsrAct(a)
	}
	if err = srv.instMngr.Init(
		true /*conf.SessMode*/,
		conf.Licence,
		conf.DelSessRt,
		usrDef,
		srv.tool); err != nil {
		srv.tool.Log.Errorf(err.Error())
		return
	}

	// 异常捕获模块
	catch.Open(conf.Catch, conf.CatchDump, srv.tool.Log, srv.instMngr.CatchCallBack)
	catch.SignalHook()

	// wrapper引擎初始化;
	if srv.callbackInit != nil {
		code, err := srv.callbackInit(conf.UsrCfgData)
		if err != nil {
			srv.tool.Log.Errorw("call wrapper init func fail", "code", code, "error", err.Error())
			return err
		}
		srv.tool.Log.Debugw("aiService.Init", "wrapper.cfg", conf.UsrCfgData)
		srv.tool.Log.Infow("aiService.Init: initialize engine wrapper succ!")
	}

	srv.tool.Log.Infof("aiService.Init: init success!")
	fmt.Println("aiService.Init: init success!")
	return nil
}

// 无缝退出下线/;
func (srv *aiService) FInit() error {
	// 框架逆初始化操作
	srv.instMngr.Fini()
	codec.AudioCodecFini()
	catch.Close()
	storage.RabFini()

	// wrapper引擎逆初始化
	if srv.callbackFini != nil {
		srv.callbackFini()
	}
	srv.tool.Log.Infof("aiService.Finit: fini success!")
	return nil
}

func (srv *aiService) Call(req *xsf.Req, span *xsf.Span) (res *xsf.Res, errXsf error) {
	span.WithName(req.Op())
	res = xsf.NewRes()
	res.SetHandle(req.Handle())

	// 获取解析数据实体
	engInput := protocol.EngInputData{}
	engOutput := protocol.EngOutputData{}
	if len(req.Data()) > 0 {
		// 协议限定: 使用protocol多数据流替换xsf多数据流;
		input := req.Data()[0].Data()
		err := proto.Unmarshal(input, &engInput)
		if err != nil {
			res.SetError(frame.AigesErrorPbUnmarshal, frame.ErrorPbUnmarshal.Error())
			span.WithErrorTag(frame.ErrorPbUnmarshal.Error())
			srv.tool.Log.Errorw("protocol invalid", "op", req.Op(), "hdl", req.Handle(), "msg", input)
			return
		}
	}

	// 交互协议参数pro & 框架服务参数bus
	proParam := req.GetAllParam()
	busParam := engInput.GetEngParam()
	// 服务框架业务逻辑
	var code int
	var err error // xsf框架缺陷,非nil导致客户端报错10201
	var inst *instance.ServiceInst
	if proParam[reqNumber] == reqFirst && req.Op() == opAIEngIn { // AIIn首帧,获取句柄.
		inst, code, err = srv.instMngr.Acquire(req.Handle(), busParam)
		if err == nil {
			var baseId int
			if baseId, err = strconv.Atoi(proParam[reqBaseId]); err != nil {
				baseId = defaultBaseId
			}
			code, err = inst.ResLoad(baseId, busParam, span)
			srv.tool.Log.Debugw("engine request params", "params", busParam, "hdl", req.Handle())
		}
	} else if req.Op() != opLBType {
		inst, code, err = srv.instMngr.Query(req.Handle())
	}

	waitTime := syncRespTimeout
	if wt, exist := proParam[reqWaitTime]; exist {
		waitTime, _ = strconv.Atoi(wt)
	}

	if err == nil {
		switch req.Op() {
		case opAIEngIn:
			{
				// 异常检查:链路交互约定; opAIEngOut接口允许携带空数据;
				if len(req.Data()) < 0 {
					res.SetError(frame.AigesErrorInvalidData, frame.ErrorInvalidData.Error())
					span.WithErrorTag(frame.ErrorInvalidData.Error())
					srv.tool.Log.Errorw("receive empty request data body", "op", req.Op(), "hdl", req.Handle())
					return
				}
				// 写入数据
				code, err = inst.DataReq(&engInput, span)
				if err == nil {
					if conf.SessMode {
						waitTime = 0 // 会话模式上行接口,非阻塞查询;
					}
					engOutput, code, err = inst.DataResp(uint(waitTime), span)
				}

				if err != nil {
					engOutput.Ret = int32(code)
					engOutput.Err = err.Error()
					engOutput.Status = protocol.EngOutputData_END
				}
			}
		case opAIEngOut:
			{
				// 阻塞方式查询结果
				// default:1000
				engOutput, code, err = inst.DataResp(uint(waitTime), span)
				if err != nil {
					engOutput.Ret = int32(code)
					engOutput.Err = err.Error()
					engOutput.Status = protocol.EngOutputData_END
				}
			}
		case opException:
			{
				// 异常释放资源
				inst.DataException(span)
				srv.instMngr.Release(req.Handle())
			}
		case opLBType:
			{
				xsf.SetLbType(busParam[reqLbType])
			}
		default:
			res.SetError(frame.AigesErrorInvalidOp, frame.ErrorInvalidOp.Error())
			span.WithErrorTag(frame.ErrorInvalidOp.Error())
			srv.tool.Log.Errorw("receive invalid operation", "op", req.Op(), "hdl", req.Handle())
			return
		}
	}

	srv.tool.Log.Infow("Call detail data In", "op", req.Op(), "hdl", req.Handle(), "inputParam", busParam, "waitTime", waitTime)

	// 服务主动中断条件：
	// 1. 非会话模式请求结束时中断
	// 2. 会话模式输出完整响应时中断
	// 3. 其他任意异常错误抛出时中断
	if (err != nil || engOutput.Ret != 0 ||
		engOutput.Status == protocol.EngOutputData_END ||
		engOutput.Status == protocol.EngOutputData_ONCE) &&
		inst != nil {
		srv.tool.Log.Debugw("DelSessionData", "hdl", req.Handle())
		inst.DataException(span) // TODO 导致异常返回10500, 非10101??
		srv.instMngr.Release(req.Handle())
	}

	// 序列化响应输出
	output, errMsl := proto.Marshal(&engOutput)
	if errMsl != nil {
		res.SetError(frame.AigesErrorPbMarshal, frame.ErrorPbMarshal.Error())
		span.WithErrorTag(frame.ErrorPbMarshal.Error())
		srv.tool.Log.Errorw("response marshal frame", "hdl", req.Handle())
	} else if err != nil {
		errorResp(code, err, res)
		res.SetError(int32(code), err.Error())
	} else {
		rd := xsf.NewData()
		rd.Append(output)
		res.AppendData(rd)
	}

	srv.tool.Log.Infow("Call detail data", "op", req.Op(), "hdl", req.Handle(), "inputParam", busParam, "engout", engOutput)
	return
}

func errorResp(errNum int, errInfo error, res *xsf.Res) {
	engOutput := protocol.EngOutputData{}
	engOutput.Ret = int32(errNum)
	engOutput.Err = errInfo.Error()
	engOutput.Status = protocol.EngOutputData_END
	output, _ := proto.Marshal(&engOutput)
	rd := xsf.NewData()
	rd.Append([]byte(output))
	res.AppendData(rd)
}
