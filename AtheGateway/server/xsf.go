package server

import (
	"fmt"
	"common"
	"conf"
	"git.xfyun.cn/AIaaS/xsf-external/client"
	xsfs "git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"github.com/golang/protobuf/proto"
	"sync"
	"time"
	"os"
	"strconv"
	"errors"
	"sync/atomic"
)

const (
	XSF_CLIENT_NAME = "webgate-ws-c"
	LAST_MSG_STATUS = 2
)

var xsfClient *xsf.Client

//初始化Xsf客户端
func InitXsfClient(cfgName string) error {
	cli, err := xsf.InitClient(
		XSF_CLIENT_NAME,
		getCfgMode(),
		utils.WithCfgCacheService(conf.Conf.Xsf.CacheService),
		utils.WithCfgCacheConfig(conf.Conf.Xsf.CacheConfig),
		utils.WithCfgCachePath(conf.Conf.Xsf.CachePath),
		utils.WithCfgName(cfgName),
		utils.WithCfgURL(conf.Centra.CompanionUrl),
		utils.WithCfgPrj(conf.Centra.Project),
		utils.WithCfgGroup(conf.Centra.Group),
		utils.WithCfgService(conf.Centra.Service),
		utils.WithCfgVersion(conf.Centra.Version),
		utils.WithCfgSvcIp(conf.Conf.Server.Host),
		utils.WithCfgSvcPort(func() int32{
			a,_:=strconv.Atoi(conf.Conf.Server.Port)
			return int32(a)
		}()),
	)
	if err != nil {
		return err
	}
	xsfClient = cli
	return nil
}

//发送请求
func SendRequest(s *Session,sid string, biz *ServerBiz, frameId int) (*UpResult, *Error) {
	data, err := proto.Marshal(biz)
	if err != nil {
		common.Logger.Errorf("%s:ServerBiz(%+v) proto.Marshal is error.error:%s", sid, biz, err)
		return nil, NewErrorByError(ErrorCodeJSONParsing, errors.New("parse proto data error"), sid, frameId)
	}
	//初始化回调者
	xsfCALLER := xsf.NewCaller(xsfClient)

	xsfCALLER.WithRetry(conf.Conf.Xsf.CallRetry)
	//初始化发送参数
	req := xsf.NewReq()
	req.Session(sid)
	if s.reqParam !=nil{
		for k, v := range s.reqParam {
			req.SetParam(k,common.ConvertToString(v))
		}
	}
	req.Append(data, nil)
	common.Logger.Infof("sid=%s,frameid=%d datalen=%d reqdatalen=%d", biz.GetGlobalRoute().GetTraceId(), biz.GetUpCall().GetSeqNo(), len(data), len(req.Data()))
	//common.Logger.Infof("sid=%s,frameid=%d datalen=%d reqdatalen=%d,busi=%v", biz.GetGlobalRoute().GetTraceId(), biz.GetUpCall().GetSeqNo(), len(data), len(req.Data()),biz.UpCall.BusinessArgs)
	//发送请求
	res, code, err := xsfCALLER.Call(s.CallService, "req", req, time.Duration(5)*time.Second)
	if res!=nil{
		s.session =res.Session()
	}
	if err != nil {
		common.Logger.Errorf(":send request error %v %v ", err, code)
		return nil, NewError(int(code), err.Error(), sid, frameId)
	}

	//解析响应结果
	respMsg := &ServerBiz{}
	err = proto.Unmarshal(res.GetData()[0].Data, respMsg)
	if err != nil {
		common.Logger.Errorf("%s:proto.Unmarshal is error.error:%s", sid, biz, err)
		return nil, NewErrorByError(ErrorCodeJSONParsing, err, sid, frameId)
	}

	if respMsg.UpResult.GetRet() != 0 {
		common.Logger.Errorf("%v:send request error %v", respMsg.UpResult.Ret, respMsg.UpResult.ErrInfo)
		return nil, NewError(int(respMsg.UpResult.Ret), respMsg.UpResult.ErrInfo, sid, frameId)
	}
	common.Logger.Infof("sid=%s msgid=%d recieve %d %d", sid, respMsg.GetUpCall().GetSeqNo(), code, respMsg.GetUpResult().GetRet())

	return respMsg.UpResult, nil
}

type ServerHandler struct {
}

//启动Xsf的服务
func StartXsfServer(cfgName string) {
	bc := xsfs.BootConfig{
		CfgMode: getCfgMode(),
		CfgData: xsfs.CfgMeta{
			CfgName:      cfgName,
			Project:      conf.Centra.Project,
			Group:        conf.Centra.Group,
			Service:      conf.Centra.Service,
			Version:      conf.Centra.Version,
			CompanionUrl: conf.Centra.CompanionUrl,
		},
	}
	var wg sync.WaitGroup
	var server = &xsfs.XsfServer{}
	//set spill enable
	go func() {
		defer wg.Done()
		err := (server).Run(bc, &ServerHandler{})
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()

}

func SendException(s *Session)  {
	if conf.Conf.Server.Mock{
		return
	}
	xsfCALLER := xsf.NewCaller(xsfClient)
	xsfCALLER.WithRetry(conf.Conf.Xsf.CallRetry)
	//初始化发送参数
	req := xsf.NewReq()
	req.Session(s.session)
	s.SeqNo++
	biz:=&ServerBiz{
		GlobalRoute:&GlobalRoute{
			SessionId:s.Sid,
			UpRouterId:XsfCallBackAddr,
			GuiderId:conf.Centra.Service,
			Appid:s.Appid,
			Uid:s.Uid,
			ClientIp:s.ClientIp,

		},
		UpCall:&UpCall{
			Call:s.Call,
			SeqNo:int32(s.SeqNo),
			From:s.From,
			Sync:false,
			Session:s.sessionMap,

		},
		MsgType:ServerBiz_UP_CALL,
		Version:conf.Centra.Version,
	}
	data,err:=proto.Marshal(biz)
	if err !=nil{
		common.Logger.Errorf("send exception error")
		return
	}
	req.Append(data,nil)
	xsfCALLER.Call(s.CallService,"exception",req,time.Duration(5)*time.Second)

}

var XsfCallBackAddr string
//业务逆初始化接口
func (serHandler *ServerHandler) FInit() error {
	time.Sleep(15 * 1000 * time.Millisecond)
	return nil
}

func (serHandler *ServerHandler) Init(toolbox *xsfs.ToolBox) error {

	XsfCallBackAddr = fmt.Sprintf("%s:%d",conf.Conf.Server.Host,toolbox.NetManager.GetPort())
	fmt.Println("xsf callback addr:",XsfCallBackAddr)
	xsfs.AddKillerCheck("server",&killed{})
	return nil
}

//回调处理
func (c *ServerHandler) Call(in *xsf.Req, span1 *xsf.Span) (*utils.Res, error) {
//	span = span.Next(utils.SrvSpan)

	serverBiz := getServerBiz(in, span1)
	if serverBiz == nil {
		return xsf.NewRes(), nil
	}

	if serverBiz.GlobalRoute == nil {
		common.Logger.Errorf("downcall globalRoute nil")
		return xsf.NewRes(), fmt.Errorf("downcall globalRoute nil")
	}

	sid := serverBiz.GlobalRoute.SessionId
	s := Get(sid)
	if s == nil {
		common.Logger.Errorf("%s:session time out is nil", sid)
		return xsf.NewRes(), nil
	}


	//重置session时间
	s.trans.resetTime()
	if serverBiz.DownCall == nil {
		common.Logger.Errorf("%s:downcall DownCall nil", sid)
		s.writeJson(NewError(int(ErrorCodeDownCall), "server error :atmos return an error data", s.Sid, s.SeqNo).GetErrorResp(s.Wscid))
		return xsf.NewRes(), nil
	}

	ret := int(serverBiz.DownCall.Ret)
	msg := in.Req().Param["error_info"]

	if ret != 0 {
		common.Logger.Errorf("downcall result: ,sid:%s,ret:%d,msg:%s", sid, ret, msg)
		if !conf.IsIgnoreSonarCode(ret){
			//sonarFail(s,s.metPtr,ret)
			s.setError(ret)
		}
		if !conf.IsIgnoreRespCode(ret){
			s.WriteError(NewError(ret, msg, sid, s.SeqNo).GetErrorResp(s.Wscid))
		}
		s.Close() //引擎给了错误码，close session
		return xsf.NewRes(), nil
	}


	if serverBiz.DownCall.DataList == nil || len(serverBiz.DownCall.DataList) == 0 {

		common.Logger.Errorf("%s:downcall DownCall is %s", sid, msg)
		s.WriteError(NewError(ret, msg, s.Sid, s.SeqNo).GetErrorResp(s.Wscid))
		return xsf.NewRes(), nil
	}
	//构建响应结果
	resp, error := NewFrameRespByServerBiz(serverBiz.DownCall.DataList, s, sid, s.Wscid, s.SeqNo)
	if error != nil {
		common.Logger.Errorf("%s:create response data failed, call:%s,reason:%s ", s.Sid, s.Call, error.Msg)
		s.writeJson(error.GetErrorResp(s.Wscid))
		s.setError(error.Code)
		return xsf.NewRes(), nil
	}

	//获取排序所需要的参数
	frameId:=serverBiz.DownCall.DataList[0].FrameId
	status:=serverBiz.DownCall.DataList[0].Status
	//收到最后一帧，改变结果状态
	if int32(status) == RESULT_STATUS_END{
		atomic.StoreInt32(&s.ResultStatus,RESULT_STATUS_END)
		s.Close()
		//common.Logger.Infof("downcall last result:sid=%s,status=%d",s.Sid,status)
	}
	common.Logger.Infof("downcall datalist: ret=%d, msg=%s,len=%d,status=%d,sid=%s", ret,msg,len(serverBiz.DownCall.DataList),status,s.Sid)

	if s.EnableSort{
		m := NewMessage(s,uint64(frameId),resp)
		common.Logger.Infof("downcall result ,frame id:%d,status:%d",frameId,status)
		if status==LAST_MSG_STATUS{
			s.Pipeline.Push(m,true)
		}else{
			s.Pipeline.Push(m,false)
		}
	}else {
		s.writeSuccess(resp)
	}

	return xsf.NewRes(), nil

}

//解析回调的请求参数
func getServerBiz(in *xsf.Req, span *xsf.Span) *ServerBiz {
	inData := in.Data()
	if inData == nil || len(inData) == 0 {
		common.Logger.Errorf("downcall from atmos err ,data is nil")

		return nil
	}
	msg := &ServerBiz{}
	if err := proto.Unmarshal(inData[0].Data(), msg); err != nil {
		common.Logger.Errorf("reveive downcall  err %v", err)
		span.WithErrorTag("downcall:" + err.Error())
		return nil
	}
	//获取数据
	return msg
}

//获取xfs初始化客户端加载配置文件的方式
func getCfgMode() utils.CfgMode {
	if *conf.BootMode {
		return utils.Native
	}
	return utils.Centre
}


type killed struct {

}


func (k *killed) Closeout() {
	fmt.Println("server be killed.")
	for true{
		num:=getCurrentSessionNum()
		if num==0{
			break
		}
		time.Sleep(200*time.Millisecond)
	}
	os.Exit(0)
}


