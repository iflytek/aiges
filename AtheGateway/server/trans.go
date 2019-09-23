package server

import (
	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"sync"
	"common"
	"time"
	"conf"
	"schemas"
)

type Trans struct {
	ParentAppid string
	ConnId   int64
	Sessions sync.Map
	Conn     *websocket.Conn
	Call     string
	ClientIp string
	Ctx *gin.Context
	CallService string
	Mapping *schemas.RouteMapping
	sync.Mutex
	KeepAlive bool
}
//transaction
func NewTrans()*Trans  {
	return &Trans{}
}

func (t *Trans) Do() {
	for {
		_, msg, err := t.Conn.ReadMessage()
		if err != nil {
			//todo Write error
			if websocket.IsCloseError(err,websocket.CloseNoStatusReceived,websocket.CloseNormalClosure,websocket.CloseAbnormalClosure){
				common.Logger.Infof("websocket unnormal closed:",t.Ctx.Request.RemoteAddr)
			}
			//如果没有拿到结果，并且没有报出其他错误，发送exception
			t.Sessions.Range(func(key, value interface{}) bool {
				if ses,ok:=value.(*Session);ok{
					if ses.ResultStatus!=2 && ses.errorCode==0{
						SendException(ses)
					}
				}
				return true
			})

			return
		}
		t.resetTime()
		req, err2 := NewFrameReqWithConn(t,&msg, "", 0)
		if err2 != nil {
			//Write error
			if req !=nil{
				s,ok:=t.Sessions.Load(req.Common.Wscid)
				sub:=""
				if (ok){
					err2.Sid = s.(*Session).Sid
					sub = s.(*Session).Call
				}else{
					//获得sub
					if req.Common.Sub!=""{
						sub = req.Common.Sub
					}else{
						sub = t.Call
					}
					sid,_:=common.NewSid(sub)  //因为此时还没有sid，所有先生成一个sid
					err2.Sid = sid
				}
				//打一个trace
				t.Write(err2.GetErrorResp(req.Common.Wscid))
			}else{
				t.Write(err2.GetErrorResp(""))
			}

			common.Logger.Errorf("parse req failed:code=%d,msg=%s",err2.Code,err2.Msg)
			continue
		}
		//handle close frame
		if req.Common.Cmd == CmdCLose{
			sess,ok:=t.Sessions.Load(req.Common.Wscid)
			if ok{
				ses:=sess.(*Session)
				ses.closeAtOnce()
				if sess.(*Session).ResultStatus!=2 && ses.errorCode==0{
					SendException(ses)
				}
			}
			continue
		}

		var sess *Session
		s,ok := t.Sessions.Load(req.Common.Wscid)
		if !ok {
			sess = NewSession()
			sess.InitByTrans(t,req)
			sess.Wscid = req.Common.Wscid
			t.Sessions.Store(req.Common.Wscid,sess)
		}else{
			sess = s.(*Session)
		}

		//发送请求
		handleFrame(sess, req)

	}
}

func (t *Trans)Close()  {
	ConnTransManager.Delete(t.ConnId)
	t.Sessions.Range(func(key, value interface{}) bool {
		se:=value.(*Session)
		se.closeAtOnce() //连接关闭时需要立即关闭连接上的session
		return true
	})
	t.Conn.Close()

}



func (t *Trans) Write(v interface{}) error {
	t.Lock()
	err:=t.Conn.WriteJSON(v)
	t.Unlock()
	if err !=nil{
		return err
	}
	return nil
}

func (t *Trans)resetTime()  {
	t.Conn.SetReadDeadline(time.Now().Add(time.Second*time.Duration(conf.Conf.Session.TimeoutInterver)))
}

//连接管理
var ConnTransManager = &TransManager{counterMap:map[string]int{}}
//key:appid
//value:*Trans
type TransManager struct {
	cache sync.Map
	counterMap map[string]int
	cmu sync.Mutex
}
//增加一个
func (t *TransManager)AddCount(appid string)  {
	t.cmu.Lock()
	t.counterMap[appid]++
	t.cmu.Unlock()
}

//减少-个
func (t *TransManager)DelCount(appid string)  {
	t.cmu.Lock()
	t.counterMap[appid]--
	t.cmu.Unlock()
}

func (t *TransManager)GetTrans(key int64)*Trans  {
	if trans,ok:=t.cache.Load(key);ok{
		return trans.(*Trans)
	}
	return nil
}

func (t *TransManager)SetTrans(key int64,tr *Trans)  {
	t.cache.Store(key,tr)
	t.AddCount(tr.ParentAppid)
}
//连接断开，删除
func (t *TransManager)Delete(key int64)  {
	tr:=t.GetTrans(key)
	if tr !=nil{
		t.DelCount(tr.ParentAppid)
	}
	t.cache.Delete(key)
}
//获取当前appid连接数
func (t *TransManager)GetCount(appid string) int {
	t.cmu.Lock()
	c:= t.counterMap[appid]
	t.cmu.Unlock()
	return c
}

func (t *TransManager)GetConnectionInfo() map[string]int {

	cm:=map[string]int{}
	t.cmu.Lock()
	t.cache.Range(func(key, value interface{}) bool {
		cm[value.(*Trans).ParentAppid]++
		return true
	})
	t.cmu.Unlock()
	return cm
}
//杀掉该appid上超过的的连接
func (t *TransManager)Kill(appid string,num int)  {
	count:=0
	t.cache.Range(func(key, value interface{}) bool {
		if value.(*Trans).ParentAppid!=appid{
			return true
		}

		count++
		if count>num{
			trans:=value.(*Trans)
			if trans!=nil{
				trans.Write(NewError(10504,"max connections exceeded ","",0).GetErrorResp(""))
				trans.Close()
			}
		}
		return true
	})
}

