package server

type Message interface {
	//获取消息编号
	MsgNo() uint64
	//发送消息
	Send() error
	//获取消息的字符串形式
	String() string
}

type myMessage struct {
	msgNo   uint64
	body    interface{}
	Session *Session
}

func NewMessage(s *Session,msgNo uint64, body interface{}) Message {
	return &myMessage{
		msgNo: msgNo,
		body:  body,
		Session:s,
	}
}

func (msg *myMessage) MsgNo() uint64 {
	return msg.msgNo
}

func (msg *myMessage) Send() error {
	msg.Session.writeSuccess(msg.body)
	return nil
}

func (msg *myMessage) String() string {

	return ""
}

type Writer interface {
	Write()
}
