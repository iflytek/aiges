package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	aiges_conf "github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/httproto/common"
	"github.com/xfyun/aiges/httproto/conf"
	dto "github.com/xfyun/aiges/httproto/http"
	"github.com/xfyun/aiges/httproto/pb"
	"github.com/xfyun/aiges/httproto/schemas"
	"github.com/xfyun/aiges/instance"
	"github.com/xfyun/aiges/protocol"
	xsf "github.com/xfyun/xsf/client"
	xsfserver "github.com/xfyun/xsf/server"
	"log"

	"github.com/xfyun/xsf/utils"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type WsSession struct {
	Handle         string
	Sid            string
	AppId          string
	Uid            string
	ServiceId      string
	RawServiceId   string
	MagicServiceId string
	Conn           *websocket.Conn
	buffer         []byte
	Sub            string
	lock           *sync.Mutex
	ReadTimeout    int
	SessonTimeout  int
	CallService    string
	CallType       int
	Ctx            *gin.Context
	schema         *schemas.AISchema
	//session        string
	sessionContext *schemas.Context
	span           *utils.Span
	startTime      time.Time
	Status         int32
	errorCode      int
	logger         *xsf.Logger
	conf           *conf.Config
	firstout       int
	lastout        int
	firstin        int
	lastin         int
	//CloudId        string
	index           int
	handlers        HandlerChain
	targetSub       string
	cid             string
	lastActiveTime  time.Time
	sessionAddress  string
	sessionCall     bool
	routerInfo      string
	XsfCallBackAddr string
	Si              xsfserver.UserInterface
	Mngr            func() *instance.Manager
	CloseChan       chan bool
}

func (s *WsSession) Run() {
	n := len(s.handlers)
	for s.index < n {
		s.handlers[s.index](s)
		s.index++
	}
}

func (s *WsSession) Next() {
	s.index++
	s.Run()
}

func (s *WsSession) Abort() {
	s.index = len(s.handlers)
}

const (
	CtxKeyClientIp = "client_ip"
	CtxKeyKongIp   = "kong_ip"
	CtxKeyHost     = "host"
)

func NewWsSession(ctx *gin.Context, schema *schemas.AISchema, cfg *conf.Config, logger *xsf.Logger, conn *websocket.Conn, lock *sync.Mutex) *WsSession {
	meta := schema.Meta
	sid := common.NewSid(meta.GetSub())
	//s, _ := uuid.NewV4()
	//#sid := s.String()
	//sess := sessionPool.GetSession()
	sess := &WsSession{}
	sess.Sid = sid
	sess.Sub = meta.GetSub()
	sess.Conn = conn
	sess.buffer = bytePool.Get()
	sess.SessonTimeout = cfg.Session.SessionTimeout
	sess.ReadTimeout = cfg.Session.TimeoutInterver
	//calls := cfg.Server.MockService
	sess.lock = lock
	// next service
	sess.CallService = meta.GetCallService()
	sess.CallType = meta.GetCallType()
	sess.Ctx = ctx
	sess.Status = 0
	sess.startTime = time.Now()
	sess.schema = schema
	sess.cid = ""
	if sess.sessionContext != nil {
		sess.sessionContext.SeqNo = 1
		sess.sessionContext.Header = nil
		sess.sessionContext.Session = nil
		sess.sessionContext.Sync = false
	} else {
		sess.sessionContext = &schemas.Context{SeqNo: 1, Sync: false, IsStream: true, InputSyncId: 0, OutPutSyncId: 0}
	}

	sess.logger = logger
	sess.conf = cfg
	//sess.sonar = nil
	sess.errorCode = 0
	sess.ServiceId = meta.GetServiceId()
	sess.RawServiceId = sess.ServiceId
	sess.MagicServiceId = ""
	if t := meta.GetSessonTimeout(); t > 0 {
		sess.SessonTimeout = t
	}
	if t := meta.GetReadTimeout(); t > 0 {
		sess.ReadTimeout = t
	}
	sess.ReadTimeout = aiges_conf.ReadTimeout
	sess.SessonTimeout = aiges_conf.SessiontTimeout
	sess.firstout = 0
	sess.lastout = 0
	sess.firstin = 0
	sess.lastin = 0
	sess.sessionAddress = ""
	sess.routerInfo = ""
	aiSessGroup.Set(sid, sess)

	return sess
}

// 从websocket 连接中读取数据，重用buffer ，提升性能
func (t *WsSession) readMessage() (int, []byte, error) {
	var r io.Reader
	messageType, r, err := t.Conn.NextReader()
	if err != nil {
		return messageType, nil, err
	}
	ed := 0 // buffer 尾
	for {
		n, err := r.Read(t.buffer[ed:])
		ed += n
		if err != nil {
			if err == io.EOF {
				return messageType, t.buffer[:ed], nil
			}
			return messageType, nil, err
		}
		if ed == len(t.buffer) { // reader 的 长度== 扩容
			old := t.buffer
			t.buffer = make([]byte, len(old)*2) // buffer cap double
			copy(t.buffer, old)
		}
	}
	//b,err:=ioutil.ReadAll(r)
	return messageType, t.buffer[:ed], err
}

// 向websocket 写消息
func (s *WsSession) WriteMessage(message interface{}) {
	data, _ := jsoniter.Marshal(message)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Conn.WriteMessage(websocket.TextMessage, data)
}

func (s *WsSession) WriteError(code common.ErrorCode, msg string) {
	// 忽略10101 等错误码，不返回给用户
	for _, c := range s.conf.Server.IgnoreRespCodes {
		if c == int(code) {
			return
		}
	}
	s.WriteMessage(&common.ErrorResp{
		Header: common.Header{
			Code:    code,
			Message: msg,
			Sid:     s.Sid,
			Cid:     s.cid,
		},
	})
}

func (s *WsSession) WriteSuccess(payload interface{}, status int, headers map[string]string) {
	rss := common.NewSuccessResp(s.Sid, payload, status, s.cid)
	for key, val := range s.schema.BuildResponseHeader(headers) {
		rss.SetHeader(key, val)
	}
	s.WriteMessage(rss)
}

func (s *WsSession) WriteSuccessWithGuiderStatus(payload interface{}, status int, guiderStatus string, headers map[string]string) {
	wfStatus := 0
	if guiderStatus != "" {
		wfStatus, _ = strconv.Atoi(guiderStatus)
	}
	rss := common.NewSuccessResp(s.Sid, payload, status, s.cid)
	rss.WfStatus = wfStatus
	for key, val := range s.schema.BuildResponseHeader(headers) {
		rss.SetHeader(key, val)
	}
	s.WriteMessage(rss)
}

// 写close 帧
func (s *WsSession) WriteClose(reason string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason))

}

func (s *WsSession) formatArgs(kvs []interface{}) []interface{} {

	common := []interface{}{"sid", s.Sid, "app_id", s.AppId, "uid", s.Uid, "serviceId", s.ServiceId, "couldId", "local"}
	args := make([]interface{}, len(kvs)+len(common))
	copy(args, common)
	copy(args[len(common):], kvs)
	return args
}

func (s *WsSession) Errorw(msg string, kvs ...interface{}) {
	//s.logger.Errorw("","")

	log.Printf(msg, s.formatArgs(kvs)...)
}

func (s *WsSession) Infow(msg string, kvs ...interface{}) {
	if s.conf.Log.Level == "error" { //日志级别为error 时就不用复制args了
		return
	}
	log.Printf(msg, s.formatArgs(kvs)...)
}

func (s *WsSession) Debugw(msg string, kvs ...interface{}) {
	if s.conf.Log.Level == "error" {
		return
	}
	log.Printf(msg, s.formatArgs(kvs)...)
}

func (s *WsSession) StartSpan() {
	s.span = utils.NewSpan(utils.SrvSpan).Start()
	s.span.WithName("aiges-ws").WithTag("sub", s.Sub).WithRetTag("0").WithTag("sid", s.Sid)
	s.span.WithTag("goroutines", strconv.Itoa(runtime.NumGoroutine()))
}

func (s *WsSession) SpanTagString(k, v string) {
	s.span.WithTag(k, v)
}

func (s *WsSession) SpanMeta() string {
	return s.span.Meta()
}

func (s *WsSession) SpanTagErr(err string) {
	s.span.WithErrorTag(err)
}

func (s *WsSession) setError(code int) {
	if s.errorCode == 0 {
		s.errorCode = code
	}
}
func (s *WsSession) SetError(code int) {
	if s.errorCode == 0 {
		s.errorCode = code
	}
}

// 检查session是否超时
func (s *WsSession) Alive() bool {
	if time.Since(s.startTime) > time.Duration(s.SessonTimeout)*time.Second {
		return false
	}
	return true
}

func (s *WsSession) SchemaCheck(o interface{}) error {
	if s.schema != nil {
		if err := s.schema.Validate(o); err != nil {
			return err
		}
	}
	return nil
}

func Handle(s *WsSession) {
	s.CloseChan = make(chan bool)
	go func() {
		for s.Status != 2 {
			select {
			case <-s.CloseChan:
				return
			default:
				if s.Handle != "" {
					inst, _, errInfo := s.Mngr().Query(s.Handle)
					if errInfo != nil {
						goto END
					}
					engOutput, _, err := inst.DataRespByChan(uint(0), nil)
					if err != nil {
						goto END
					}
					if len(engOutput.GetPl()) > 0 {
						s.ResolveLoadoutput(&engOutput, nil, s.Ctx)
					}

				}
			}
		END:
			time.Sleep(time.Duration(50) * time.Millisecond)

		}
	}()
	for {
		s.ResetReadDeadline()
		// 上行数据for循环读取 请求
		_, msg, err := s.readMessage()
		if err != nil {
			//if !websocket.IsCloseError(err,websocket.CloseNormalClosure,websocket.CloseNoStatusReceived,websocket.CloseAbnormalClosure){
			s.StartSpan()
			s.SpanTagString("sessionCloseErr", fmt.Sprintf("connection close: err:%s ; cost:%d", err.Error(), time.Since(s.startTime)/time.Second))
			s.SpanTagString("appid", s.AppId)
			s.Errorw("read message error:", "error", err.Error(), "cost", time.Since(s.startTime))
			s.FlushSpan()
			s.WriteClose("time out")
			//}
			return
		}

		if !s.Alive() {
			s.StartSpan()
			s.SpanTagString("sessionCloseErr", fmt.Sprintf("session timeout! used:%d   allow:%d", time.Since(s.startTime)/time.Second, s.SessonTimeout))
			s.SpanTagString("appid", s.AppId)
			s.FlushSpan()
			s.Errorw("session timeout")
			s.WriteError(common.ErrorCodeSetReadDeadline, "session timeout")
			s.WriteClose("session timeout")
			return
		}
		s.sessionContext.InputSyncId++
		code, info := s.handleAIUpMessage(msg)
		if code != 0 {
			s.WriteError(code, info)
			time.Sleep(5 * time.Millisecond)
			s.WriteClose(info)
			s.setError(int(code))

			return
		}
	}
}

// close 只能被执行一次
func (s *WsSession) CloseSession() {
	aiSessGroup.Delete(s.Sid)
	time.Sleep(1 * time.Second)
	s.Conn.Close()
	bytePool.Put(s.buffer)
	if s.errorCode == 0 && s.Status != 2 {
		s.SendException()
	}
	s.CloseChan <- true
	//sessionPool.PutSession(s)
}

func (s *WsSession) CloseMulitplex() {
	aiSessGroup.Delete(s.Sid)
	//time.Sleep(1 * time.Second)
	//s.Conn.Close()
	bytePool.Put(s.buffer)
	if s.errorCode == 0 && s.Status != 2 {
		s.SendException()
	}
	//sessionPool.PutSession(s)
}

// close 只能被执行一次
func (s *WsSession) CloseHttpSession() {
	bytePool.Put(s.buffer)
	aiSessGroup.Delete(s.Sid)
	//sessionPool.PutSession(s)
}

func (s *WsSession) SendException() {
	s.Errorw("session close unexpected")
	s.StartSpan()
	defer s.FlushSpan()
	biz := &protocol.LoaderOutput{
		ServiceId: s.ServiceId,
		Code:      int32(common.ErrorCodeTimeOut),
		Status:    StatusEnd,
		Err:       "session close unexpected",
		Pl:        nil,
	}

	data, err := proto.Marshal(biz)
	if err != nil {
		s.Errorw("pbMarshal error while send exception", "error", err)
		return
	}

	//xsfCALLER := xsf.NewCaller(xsfClient)
	//xsfCALLER.WithRetry(1)
	req := xsf.NewReq()
	req.SetTraceID(s.SpanMeta())
	req.Append(data, nil)
	span := utils.NewSpan(utils.SrvSpan)
	span.WithName(s.ServiceId)
	span.WithTag("sid", s.Sid)
	_, err = s.Si.Call(req, span)
	if err != nil {
		s.Errorw("send exception error", "error", err.Error())
	}

}

func (s *WsSession) FlushSpan() {
	if s.span == nil {
		return
	}
	s.SpanTagString("ret", strconv.Itoa(s.errorCode))
	s.span.End().Flush()
}

func bizToString(biz *pb.ServerBiz) string {
	busiStr := map[string]string{}
	for k, v := range biz.GetUpCall().GetBusinessArgs() {
		busiStr[k+"args"] = common.MapstrToString(v.GetBusinessArgs())
		ples := make(map[string]string)
		for acpt, desc := range v.GetPle() {
			ple := common.NewStringBuilder()
			ple.AppendIfNotEmpty("attr", common.MapstrToString(desc.GetAttribute()))
			ple.AppendIfNotEmpty("accept", desc.GetName())
			ple.AppendIfNotEmpty("data_type", desc.GetDataType().String())
			ples[acpt] = ple.ToString()
		}
		busiStr[k+"ple"] = common.MapstrToString(ples)
	}
	m := map[string]interface{}{
		"header":   common.MapstrToString(biz.GetGlobalRoute().GetHeaders()),
		"business": common.MapstrToString(busiStr),
		//"payload":  common.MapToString(payload),
	}
	return common.MapToString(m)
}

func bizPayloadString(biz *pb.ServerBiz) string {
	//payload := map[string]string{}
	sb := strings.Builder{}
	for _, v := range biz.GetUpCall().GetDataList() {
		sb.WriteString(v.GetMeta().GetName())
		sb.WriteString(" dataLen:")
		sb.WriteString(strconv.Itoa(len(v.GetData())))
		//sb.WriteString("service:")
		//sb.WriteString(v.GetMeta().GetServiceName())
		sb.WriteString(" dataType:")
		sb.WriteString(pb.MetaDesc_DataType_name[int32(v.GetMeta().GetDataType())])
		sb.WriteString(" attr:")
		sb.WriteString(common.MapstrToString(v.GetMeta().GetAttribute()))
		//payload[v.GetMeta().GetName()] = fmt.Sprintf("service=%s,dataLen=%d,dataType=%v,attr=%s", v.GetMeta().GetServiceName(), len(v.GetData()), v.GetMeta().GetDataType(), common.String(v.GetMeta().GetAttribute()))
		sb.WriteString(" | ")
	}
	return sb.String()
}

// s.Ctx.GetHeader("X-Consumer-Username")
func CheckAppIdMatching(readAppid, consumerAppid string, enabled bool) bool {
	if !enabled {
		return true
	}
	//if realAppid=="" 那么可能不是走kong，或者kong没有开启鉴权，也放过
	if readAppid == "" {
		return true
	}
	if readAppid != consumerAppid {
		return false
	}
	return true
}

func (s *WsSession) sendMessage(in *dto.Request) (common.ErrorCode, string) {
	//sonar

	if s.Status == StatusBegin {
		if s.schema.Meta.IsCategory() {
			subServiceId, routerInfo := s.schema.GetSubServiceId(in)
			s.routerInfo = routerInfo
			// 子serviceId 不为空，使用子serviceID ，并在 sub=ase 时获取子schema
			if subServiceId != "" {
				s.ServiceId = subServiceId
				s.SpanTagString("useMappedServiceId", "true")
				s.SpanTagString("routerInfo", routerInfo)
				// tood schema sc := schemas.GetSchemaByServiceId(subServiceId, cloudId)
				s.Debugw("companion route", "subSrv", subServiceId, "routeInfo", routerInfo)
				if s.schema == nil {
					return common.ErrorSubServiceNotFound, fmt.Sprintf("sub service not found:%s", subServiceId)
				}
				//s.schema = s
			} else {
				return common.ErrorSubServiceNotFound, fmt.Sprintf("no category route find")
			}
		}
	}
	//clientIp := s.Ctx.GetString(CtxKeyClientIp)
	// 组装加载器输入 PB
	ctx := s.Ctx.Request.Context()
	lin, err := in.ConvertToPb(s.ServiceId, protocol.LoaderInput_STREAM, &ctx, s.sessionContext.InputSyncId)
	if err != nil {
		return 10001, err.Error()
	}

	upResult, respheader, bizE := s.SendAIBizLocal(lin, s.CallType)
	///
	s.SpanTagString("servieId2", s.MagicServiceId)
	if bizE != nil {
		s.Errorw("send biz error", "error", bizE.Message, "code", bizE.Code, "call", s.CallService)
		s.SpanTagErr(bizE.Error())
		return bizE.Code, bizE.Message
	}
	s.ResolveLoadoutput(upResult, respheader, s.Ctx)
	s.sessionContext.SeqNo++
	return 0, ""
}

func (s *WsSession) handleAIUpMessage(msg []byte) (common.ErrorCode, string) {
	s.StartSpan()
	defer s.FlushSpan()
	req := new(dto.Request)
	err := json.Unmarshal(msg, &req)
	if err != nil {
		s.Errorw("json unmarshal error", "error", err.Error(), "json_data", common.ToString(msg))
		return common.ErrorCodeGetUpCall, "parse request json error"
	}
	// schema 校验，校验请求参数的合法性
	//if err := s.SchemaCheck(in); err != nil {
	//	s.Errorw("schema validate error", "error", err.Error(), "data", common.ToString(msg))
	//	return ErrorCodeGetUpCall, err.Error()
	//}
	return s.sendMessage(req)
}

//func bizToLogString(biz *pb.ServerBiz)string{
//	for srv, args := range biz.GetUpCall().GetBusinessArgs() {
//
//	}
//}

func (s *WsSession) ResolveUpResult(upr *pb.UpResult, header map[string]string) {
	if s.targetSub == "" {
		if upr.GetSession() != nil {
			s.targetSub = upr.GetSession()[KeyAipaaSSUb]
		}
	}
	if s.sessionContext.Session == nil {
		s.sessionContext.Session = upr.GetSession()
	}
	s.Debugw("up result", "session", upr.GetSession())
	if len(upr.GetDataList()) == 0 {
		if s.Status == StatusBegin {
			s.WriteSuccess(nil, 0, header)
			s.Status = StatusContinue
		}
		return
	}

	if s.Status == StatusBegin {
		s.Status = StatusContinue
	}

	result := s.schema.ResolveUpResult(upr)
	s.WriteSuccess(result, int(upr.GetStatus()), header)

}

func (s *WsSession) ResolveLoadoutput(upr *protocol.LoaderOutput, header map[string]string, ctx *gin.Context) {
	if len(upr.GetPl()) == 0 {
		if s.Status == StatusBegin {
			s.WriteSuccess(nil, 0, header)
			s.Status = StatusContinue
		}
		return
	}
	if s.Status == StatusBegin {
		s.Status = StatusContinue
	}

	p := s.schema.ResolveLoadOutput(upr)
	s.WriteSuccess(p, int(upr.GetStatus()), header)
	if upr.GetStatus() == protocol.LoaderOutput_ONCE {
		s.Status = StatusEnd
		s.CloseSession()
	}
}

func (s *WsSession) ResetReadDeadline() {
	s.Conn.SetReadDeadline(time.Now().Add(time.Duration(s.ReadTimeout) * time.Second))
}

var mockData = []byte("this is mocked data")

func (s *WsSession) MockCall(header map[string]string) {
	status := header["status"]
	if s.sessionContext.SeqNo == 0 || status == "2" || s.sessionContext.SeqNo%10 == 0 {
		data := map[string]interface{}{
			"result": map[string]interface{}{
				"text": "1234",
			},
		}
		s.WriteSuccess(data, common.IntFromString(status), nil)
		s.Status = StatusEnd
	}
}

type SendBizError struct {
	Code    common.ErrorCode
	Message string
}

func (e *SendBizError) Error() string {
	return fmt.Sprintf("%d|%s", e.Code, e.Message)
}

func NewSendBizError(code common.ErrorCode, msg string) *SendBizError {
	return &SendBizError{
		Code:    code,
		Message: msg,
	}
}

func (s *WsSession) call(req *xsf.Req, op string) (res *xsf.Res, code int32, err error) {
	//xsfCALLER := xsf.NewCaller(xsfClient)
	//xsfCALLER.WithRetry(s.conf.Xsf.CallRetry)
	if !s.sessionCall || s.sessionAddress == "" {
		// todo do
		fmt.Errorf("print $$$$$$$$$", op)
		//res, code, err = xsfCaller.Call(s.CallService, op, req, time.Duration(5)*time.Second)
		if res != nil && s.sessionAddress == "" {
			s.sessionAddress, _ = res.GetPeerIp()
			s.SpanTagString("serviceAddr", s.sessionAddress)
		}
	} else {
		//res, code, err = xsfCALLER.CallWithAddr(s.CallService, op, s.sessionAddress, req, time.Duration(5)*time.Second)
		s.SpanTagString("callWithAddr", s.sessionAddress)
	}
	return
}

func (s *WsSession) SendAIBizLocal(biz *protocol.LoaderInput, callType int) (*protocol.LoaderOutput, map[string]string, *SendBizError) {
	s.Debugw("success send request to backend")
	cl2 := time.Now()
	biz.Headers["sid"] = s.Sid
	//in.Expect[0].DataType = protocol.MetaDesc_DataType(protocol.MetaDesc_TEXT)
	bytes, _ := proto.Marshal(biz)
	xsfReq := xsf.NewReq()
	xsfReq.Append(bytes, nil)
	xsfReq.SetOp("AIIn")
	xsfReq.SetParam("SeqNo", strconv.Itoa(int(s.sessionContext.SeqNo)))
	xsfReq.SetParam("version", "v2")
	xsfReq.SetParam("waitTime", "1000")
	xsfReq.SetParam("baseId", "0")
	xsfReq.SetHandle(s.Sid)
	s.Handle = s.Sid

	span := utils.NewSpan(utils.SrvSpan)
	span.WithName(s.ServiceId)
	span.WithTag("sid", s.Sid)

	res, err := s.Si.Call(xsfReq, span)

	if err != nil {
		s.Errorw(":send request error", "error", err.Error(), "code", common.ErrorReqError)
		return nil, nil, NewSendBizError(common.ErrorReqError, "send request to backend error:"+err.Error())
	}

	cl3 := time.Now()
	s.SpanTagString("xsfCallCost", strconv.Itoa(int(cl3.Sub(cl2).Nanoseconds())))

	//解析响应结果
	respMsg := &protocol.LoaderOutput{}
	err = proto.Unmarshal(res.GetData()[0].Data, respMsg)
	if err != nil {
		s.Errorw("proto.Unmarshal up result error", "error", err.Error())
		return nil, nil, NewSendBizError(common.ErrorCodeJSONParsing, "invalid up result message")
	}
	return respMsg, nil, nil
}

func (s *WsSession) SendAIBizByXsf(biz *pb.ServerBiz, callType int) (*pb.UpResult, map[string]string, *SendBizError) {
	cl1 := time.Now()
	data, err := proto.Marshal(biz)
	if err != nil {
		s.Errorw("pbMarshal error", "error", err)
		return nil, nil, NewSendBizError(common.ErrorCodeGetUpCall, "pb marshal error"+err.Error())
	}
	cl2 := time.Now()
	s.SpanTagString("pbMarshalCost", strconv.Itoa(int(cl2.Sub(cl1).Nanoseconds())))
	//初始化回调者
	//xsfCALLER := xsf.NewCaller(xsfClient)
	//
	//xsfCALLER.WithRetry(s.conf.Xsf.CallRetry)

	//初始化发送参数
	req := xsf.NewReq()
	req.SetTraceID(s.SpanMeta())
	//req.Session(s.Sid)

	req.Append(data, nil)

	var res *xsf.Res
	var code int32
	//var err error
	req.SetHandle(s.Sid)

	span := utils.NewSpan(utils.SrvSpan)
	span.WithName(s.ServiceId)
	span.WithTag("sid", s.Sid)
	res, err = s.Si.Call(req, span)
	res, code, err = s.call(req, "req")
	//if !s.sessionCall || s.sessionAddress == "" {
	//	res, code, err = xsfCALLER.Call(s.CallService, "req", req, time.Duration(5)*time.Second)
	//	if res != nil && s.sessionAddress == "" {
	//		s.sessionAddress, _ = res.GetPeerIp()
	//		s.SpanTagString("serviceAddr", s.sessionAddress)
	//	}
	//} else {
	//	res, code, err = xsfCALLER.CallWithAddr(s.CallService, "req", s.sessionAddress, req, time.Duration(5)*time.Second)
	//	s.SpanTagString("callWithAddr", s.sessionAddress)
	//}
	//if res != nil {
	//	s.session = res.Session()
	//}
	if err != nil {
		s.Errorw(":send request error", "error", err.Error(), "code", code)
		return nil, nil, NewSendBizError(common.ErrorCode(code), "send request to backend error:"+err.Error())
	}
	cl3 := time.Now()
	s.SpanTagString("xsfCallCost", strconv.Itoa(int(cl3.Sub(cl2).Nanoseconds())))

	//解析响应结果
	respMsg := &pb.ServerBiz{}
	err = proto.Unmarshal(res.GetData()[0].Data, respMsg)
	if err != nil {
		fmt.Println(string(res.GetData()[0].Data))
		s.Errorw("proto.Unmarshal up result error", "error", err.Error())
		return nil, nil, NewSendBizError(common.ErrorCodeJSONParsing, "invalid up result message")
	}

	header := respMsg.GetGlobalRoute().GetHeaders()
	if header != nil {
		serviceId := header["service_id"]
		if serviceId != "" && s.MagicServiceId == "" {
			s.MagicServiceId = serviceId
			s.ServiceId = serviceId
		}
	}

	if respMsg.GetUpResult().GetRet() != 0 {
		s.Errorw("get up result error", "error", respMsg.GetUpResult().GetErrInfo(), "code", respMsg.GetUpResult().GetRet())
		return respMsg.GetUpResult(), header, NewSendBizError(common.ErrorCode(respMsg.GetUpResult().GetRet()), respMsg.GetUpResult().GetErrInfo())
	}
	//s.Debugw("success send request to backend")
	return respMsg.GetUpResult(), header, nil
}

func (s *WsSession) readBody(body io.Reader) ([]byte, error) {
	bf := bytes.NewBuffer(s.buffer[:0])
	_, err := bf.ReadFrom(body)
	if err != nil {
		if err == io.EOF {
			return bf.Bytes(), nil
		}
		return nil, err
	}
	return bf.Bytes(), nil
}

// once
var (
	aiSessGroup *AISessionGroup
)

type AISessionGroup struct {
	lock     sync.RWMutex
	sess     map[string]*WsSession
	interval time.Duration
}

func InitSessionGroup(interval int) {
	aiSessGroup = &AISessionGroup{sess: map[string]*WsSession{}, interval: time.Duration(interval) * time.Second}
}

func (g *AISessionGroup) Get(sid string) *WsSession {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.sess[sid]
}

func (g *AISessionGroup) Set(sid string, sess *WsSession) {
	g.lock.Lock()
	g.sess[sid] = sess
	g.lock.Unlock()
}

func (g *AISessionGroup) Delete(sid string) {
	g.lock.Lock()
	delete(g.sess, sid)
	g.lock.Unlock()
}

// attention read only ，if write in this function ，will occur dead lock
func (g *AISessionGroup) Range(f func(sid string, sess *WsSession) bool) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	for k, v := range g.sess {
		if !f(k, v) {
			return
		}
	}
}

func (g *AISessionGroup) CheckIdleInBackground() {
	if g.interval == 0 {
		g.interval = 30 * time.Second
	}

	go func() {
		for range time.Tick(g.interval) {
			g.lock.RLock()
			deleted := make([]string, 0, 10)
			for sid, s := range g.sess {
				if !s.Alive() {
					deleted = append(deleted, sid)
				}
			}
			g.lock.RUnlock()
			for _, sid := range deleted {
				g.Delete(sid)
			}
		}
	}()

}
