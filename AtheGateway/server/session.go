package server

import (
	"common"
	"conf"
	ws "github.com/gorilla/websocket"
	"time"
	"github.com/gin-gonic/gin"
	"strings"
	"schemas"
	"sync/atomic"
)


const (
	closeStatusWorking int32 = 0
	closeStatusClosing int32 = 1
	closeStatusClosed int32 = 2
)

type Session struct {
	Pipeline 	Pipeline  //排序管道
	Sid         string
	Appid       string
	CallService string   //调用的atmos
	From        string  //from
	EnableSort   bool
	Ctx          *gin.Context
	Conn         *ws.Conn
	Call         string //本次会话使用的ai能力
	Status       int  //发送数据的status
	ResultStatus int32 // 接收数据的status
	SeqNo        int
	Uid          string
	Wscid          string
	ClientIp     string
	readTimeOut  time.Duration
	startTime    time.Time
	connTime     time.Time
	sessionMap   map[string]string
	session      string    //xsf调用的session
	endpoint     string
	errorCode    int
	dataargs     map[string]interface{}  //记录data中的xsflog
	trans *Trans
	closeStatus int32
	Mapping *schemas.RouteMapping
	reqParam map[string]interface{}
}

func NewSession() *Session {
	s := &Session{}
	//Insert(sid, s)
	return s
}

func (s *Session)Init(conn *ws.Conn, sid string, clientIp string)  {
	s.Sid = sid
	//s.Call = call
	s.Conn = conn
	s.startTime = time.Now()
	s.connTime = time.Now()
	s.Status = 0
	s.SeqNo = 1
	s.ResultStatus = RESULT_STATUS_RECEIVING
	s.ClientIp = clientIp
	s.readTimeOut = time.Duration(conf.Conf.Session.TimeoutInterver) * time.Second
	s.sessionMap = nil
	s.errorCode = 0
	s.session=""
	s.endpoint = ""
	s.dataargs =map[string]interface{}{};
	//结束
	Insert(sid, s)

}

func (s *Session) InitByTrans(t *Trans ,req *FrameReq)  {
	s.trans = t
	s.Ctx = t.Ctx
	s.Call = t.Call
	s.CallService = t.CallService
	s.From = conf.Conf.Xsf.From
	s.EnableSort = conf.Conf.Xsf.EnableRespsort
	s.Mapping = t.Mapping
	//aikit
	if len(req.Common.Sub) > 0{
		s.Call = req.Common.Sub
	}
	s.CallService = t.Mapping.GetAtmos(s.Call)
	sid,_:=common.NewSid(s.Call)

	if s.EnableSort{
		s.Pipeline =NewPipeline(uint64(conf.Conf.Server.PipeDepth),time.Duration(conf.Conf.Server.PipeTimeout)*time.Millisecond)
		s.Pipeline.Open()
	}

	s.Init(t.Conn, sid, t.ClientIp)

}

func (s *Session) checkSessionTimeOut() bool {
	return time.Now().Sub(s.startTime) > time.Duration(conf.Conf.Session.SessionTimeout) * time.Second
}


func (s *Session) checkTimeOut() bool {
	if s.checkSessionTimeOut(){
		Remove(s.Sid)
		s.Close()
		return true
	}
	return false
}

func (s *Session) GetEndpoint() string {
	if s.endpoint==""{
		hp:=strings.SplitN(conf.Conf.Xsf.XsfLocalIp,":",2)
		s.endpoint = hp[0]
		return s.endpoint
	}
	return s.endpoint
}
//立即关闭session
func (s *Session) closeAtOnce()  {
	s.clear()
}

//有错误
func (s *Session) Close() error {
	//err := s.Conn.Close()
	if !atomic.CompareAndSwapInt32(&s.closeStatus,closeStatusWorking,closeStatusClosing) {
		return nil
	}

	time.AfterFunc(time.Second*time.Duration(conf.Conf.Session.SessionCloseWait),s.clear)
	atomic.StoreInt32(&s.closeStatus ,closeStatusClosing)//
	//SendException(s)
	return nil
}


func (s *Session)clear()  {
	if atomic.SwapInt32(&s.closeStatus,closeStatusClosed) == closeStatusClosed{  //如果closeStatus 已经是2 则直接return
		return
	}

	s.trans.Sessions.Delete(s.Wscid) //删除当前连接中的会话
	if !s.trans.KeepAlive{
		s.trans.Conn.Close()  //关闭连接。非长连接模式下
	}

	Remove(s.Sid) // 删除sid
	if s.Pipeline!=nil{
		s.Pipeline.Close()
	}
}

func (s *Session) writeJson(data interface{}) {
	error := s.trans.Write(data)
	if error != nil {
		common.Logger.Errorf("%s:向客户端:%s写数据时失败,失败的原因是:%s", s.Sid, s.ClientIp, error.Error())
	}
}

func (s *Session) WriteError(data *FrameResponse) {
	s.writeJson(data)
	//time.Sleep(100 * time.Millisecond)
}

func (s *Session) writeSuccess(resp interface{}) {
	s.writeJson(resp)
}

func (s *Session)setError(code int)  {

	if s.errorCode==0{
		s.errorCode=code
	}

}