package server

type ErrorCode int

const (
	ErrorCodeJSONParsing     ErrorCode = 10160  //json解析异常
	ErrorCodeGetUpCall       ErrorCode = 10163  //根据映射获取上行参数失败
	ErrorCodeDecoding        ErrorCode = 10161  //解码错误
	ErrorCodeGetRespCall     ErrorCode = 10164  //根据映射响应结果失败
	ErrorCodeTimeOut         ErrorCode = 10114  //超时错误码
	ErrorCodeSessionEnd      ErrorCode = 10202  //会话已经结束
	ErrorInvalidAppid        ErrorCode = 10313 //appId非法
	ErrorServerError         ErrorCode = 11502 //服务端配置错误
	ErrorCodeDownCall        ErrorCode = 11503  //atmos回调DownCall为空
	ErrorCodePASEDATA        ErrorCode = 10118  //解析响应的Data失败
	//
)


type Error struct {
	Code    int
	Msg     string
	Sid     string
	FrameId int
}

func NewErrorByError(code ErrorCode, err error, sid string, frameId int) *Error {
	return NewError(int(code), err.Error(), sid, frameId)
}

func NewError(code int, msg string, sid string, frameId int) *Error {
	return &Error{
		Code:    code,
		Msg:     msg,
		Sid:     sid,
		FrameId: frameId,
	}
}

//获取错误的响应
func (err *Error) GetErrorResp(wscid string) *FrameResponse {
	return NewFrameRespByError(err, wscid)
}

func (err *Error) string() string {
	return string(err.Code) + ":" + err.Msg
}

func (err *Error)Error()string  {
	return err.string()
}

