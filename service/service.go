package service

import (
	"fmt"
	"github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/env"
	"github.com/xfyun/aiges/frame"
	"github.com/xfyun/aiges/instance"
	"github.com/xfyun/aiges/protocol"
	"github.com/xfyun/aiges/storage"
	aigesUtils "github.com/xfyun/aiges/utils"
	"github.com/xfyun/xsf/server"
	"strconv"
)

type reqMode int

const (
	reqAlloc  reqMode = 0
	reqCalcu  reqMode = 1
	reqAssist reqMode = 2
)

var InstMngr instance.Manager

type aiService struct {
	tool         *xsf.ToolBox
	instMngr     *instance.Manager
	callbackInit actionInit
	callbackFini actionFini
	callbackUser map[usrEvent]actionUser

	Coordinator *aigesUtils.Coordinator
}

func (srv *aiService) Init(box *xsf.ToolBox) (err error) {
	srv.tool = box
	srv.instMngr = &InstMngr
	// 框架配置初始化
	if err = conf.Construct(srv.tool.Cfg); err != nil {
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
		conf.Licence,
		conf.DelSessRt,
		usrDef,
		srv.tool); err != nil {
		srv.tool.Log.Errorf(err.Error())
		return
	}

	// 如果是python grpc模式需要阻塞下
	if env.AIGES_PLUGIN_MODE == "python" {
		srv.Coordinator.ConfChan <- box.Cfg
		<-srv.Coordinator.Ch2
	}

	// wrapper引擎初始化;
	if srv.callbackInit != nil {
		code, err := srv.callbackInit(conf.UsrCfgData)
		if err != nil {
			srv.tool.Log.Errorw("call wrapper init func fail", "code", code, "error", err.Error())
			fmt.Println("call wrapper init func fail", "code", code, "error", err.Error())
			return err
		}
		srv.tool.Log.Debugw("aiService.Init", "wrapper config", conf.UsrCfgData)
	}

	aigesUtils.SetCommonLogger(srv.tool.Log)
	srv.tool.Log.Infof("aiService.Init: init success!")
	fmt.Println("aiService.Init: init success!")
	return nil
}

// 无缝退出下线/;

func (srv *aiService) Finit() error {
	fmt.Println("aiService.Finit: fini begin!")
	// 框架逆初始化操作
	srv.instMngr.Fini()

	storage.RabFini()
	// wrapper引擎逆初始化
	if srv.callbackFini != nil {
		srv.callbackFini()
	}
	srv.tool.Log.Infof("aiService.Finit: fini ALL success!")
	fmt.Println("aiService.Finit: fini ALL success!")
	return nil
}

func (srv *aiService) Call(req *xsf.Req, span *xsf.Span) (res *xsf.Res, xsferr error) {
	defer span.Flush()
	defer span.End()
	span.WithName(req.Op())
	res = xsf.NewRes()
	res.SetHandle(req.Handle())
	srv.tool.Log.Debugw("call req detail before adapter", "params", req.GetAllParam(), "dataLen", len(req.Data()))
	// 数据实体解析&适配
	var code int
	var err error
	engInput := protocol.LoaderInput{}
	engOutput := protocol.LoaderOutput{}
	if len(req.Data()) > 0 {
		input := req.Data()[0].Data()
		code, err = protocol.InputAdapter(input, &engInput)
		if err != nil {
			errResp(code, err, res, span)
			srv.tool.Log.Errorw("protocol invalid", "op", req.Op(), "hdl", req.Handle())
			return
		}
	}
	srv.tool.Log.Debugw("call req after adapter", "params", engInput.GetParams(),
		"dataLen", len(engInput.GetPl()), "hdl", req.Handle())
	for k, pl := range engInput.GetPl() {
		srv.tool.Log.Debugw("call req after adapter data part ", "idx", k, "meta", pl.Meta, "hdl", req.Handle())
	}
	// 协议交互参数&服务框架参数
	protos := req.GetAllParam()
	headers := engInput.GetHeaders()
	params := engInput.GetParams()
	waitTime := syncRespTimeout
	if wt, exist := protos[reqWaitTime]; exist {
		waitTime, _ = strconv.Atoi(wt)
	}

	// 获取服务实例
	var inst *instance.ServiceInst
	if req.Op() == opAIEngIn || req.Op() == opAIEngOut || req.Op() == opException {
		switch protos[reqNumber] {
		case reqFirst:
			if inst, code, err = srv.instMngr.Acquire(req.Handle(), headers); err == nil {
				if code, err = inst.ResLoad(&protos, &engInput, span); err == nil {
					// _ = auth.MeterCount(headers) // ssb success, calc add TODO 平台通用计量规则待进一步确认
				}
				span.WithTag("usrTag", req.Handle())
			}
			srv.tool.Log.Debugw("request first", "header", headers,
				"param", params, "proto", protos, "state", engInput.State, "hdl", req.Handle())
		default:
			inst, code, err = srv.instMngr.Query(req.Handle())
		}
	}

	if err == nil {
		switch req.Op() {
		case opAIEngIn:
			{
				// 写入数据
				code, err = inst.DataReq(&engInput, span)
				if err == nil {
					if inst.GetSessState() == protocol.LoaderInput_STREAM {
						waitTime = 0 // 会话模式上行接口,非阻塞查询;
					}
					engOutput, code, err = inst.DataResp(uint(waitTime), span)
				}
			}
		case opAIEngOut:
			{
				engOutput, code, err = inst.DataResp(uint(waitTime), span)
			}
		case opException:
			{
				inst.DataException(span)
				srv.instMngr.Release(req.Handle())
			}
		case opResUpdate:
			{
			}
		case opLBType:
			{
				srv.instMngr.UpdateLic(protos)
			}
		case opDetect:
		case opAILoad:
		case opAIUnLoad:
		default:
			errResp(frame.AigesErrorInvalidOp, frame.ErrorInvalidOp, res, span)
			srv.tool.Log.Errorw("invalid operation", "op", req.Op(), "hdl", req.Handle())
			return
		}
	}

	// 服务主动中断条件：
	// 1. 非会话模式请求结束时中断
	// 2. 会话模式输出完整响应时中断
	// 3. 其他任意异常错误抛出时中断
	if (err != nil || code != 0 ||
		engOutput.Status == protocol.LoaderOutput_END ||
		engOutput.Status == protocol.LoaderOutput_ONCE) &&
		inst != nil {
		inst.DataException(span) // TODO 导致异常返回10500, 非10101??
		srv.instMngr.Release(req.Handle())
		srv.tool.Log.Debugw("call release handle",
			"op", req.Op(), "hdl", req.Handle(), "code", code)
	}
	if err != nil {
		errResp(code, err, res, span)
		srv.tool.Log.Errorw("call op fail", "op", req.Op(), "hdl", req.Handle(), "code", code, "err", err.Error())
		return
	}
	// 结果序列化
	engOutput.ServiceId = engInput.ServiceId
	output, code, err := protocol.OutputAdapter(&engOutput)
	if err != nil {
		errResp(code, err, res, span)
		srv.tool.Log.Errorw("response marshal fail", "hdl", req.Handle(), "err", err.Error())
	} else {
		rd := xsf.NewData()
		rd.Append(output)
		res.AppendData(rd)
	}

	srv.tool.Log.Infow("Call detail data", "op", req.Op(), "hdl", req.Handle(), "wait", waitTime, "input", engInput, "output", engOutput)
	return
}

func errResp(errNum int, errInfo error, res *xsf.Res, span *xsf.Span) {
	engOutput := protocol.LoaderOutput{}
	engOutput.Code = int32(errNum)
	engOutput.Err = errInfo.Error()
	engOutput.Status = protocol.LoaderOutput_END
	output, _, _ := protocol.OutputAdapter(&engOutput)
	rd := xsf.NewData()
	rd.Append([]byte(output))
	res.AppendData(rd)
	if errNum == frame.AigesErrorLicNotEnough {
		// 防止报错导致xsf重试,限制仅授权不足时重试
		res.SetError(int32(errNum), errInfo.Error())
	}
	span.WithErrorTag(errInfo.Error()).WithRetTag(strconv.Itoa(errNum))
}
