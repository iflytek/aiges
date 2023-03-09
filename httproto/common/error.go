package common

import "fmt"

type ErrorCode int

const (
	ErrorCodeUpgrade         ErrorCode = 10031 //http协议升级为ws
	ErrorCodeGenerateSID     ErrorCode = 10324 //Sid生成失败
	ErrorCodeSetReadDeadline ErrorCode = 10114 //设置读取超时时间
	ErrorCodeCommonArgsNil   ErrorCode = 10106 //解码错误
	ErrorDataNil             ErrorCode = 10109 //解码错误
	ErrorSchemaValidFailed   ErrorCode = 10107 //schema 校验失败
	//ErrorCodeAppId         ErrorCode = 10313 //appId为空
	ErrorCodeJSONParsing          ErrorCode = 10160 //json解析异常
	ErrorCodeGetUpCall            ErrorCode = 10163 //根据映射获取上行参数失败
	ErrorCodeInvalidSessionHandle ErrorCode = 10165 //根据映射获取上行参数失败
	ErrorCodeDecoding             ErrorCode = 10161 //解码错误
	ErrorCodeGetRespCall          ErrorCode = 10164 //根据映射响应结果失败
	ErrorCodeTimeOut              ErrorCode = 10114 //超时错误码
	ErrorCodeConnRead             ErrorCode = 10200 //网络读取失败
	ErrorCodeSessionEnd           ErrorCode = 10202 //会话已经结束
	ErrorInvalidAppid             ErrorCode = 10313 //appId非法
	ErrorServerError              ErrorCode = 11502 //服务端配置错误
	ErrorCodeDownCall             ErrorCode = 11503 //atmos回调DownCall为空
	ErrorCodePASEDATA             ErrorCode = 10118 //解析响应的Data失败
	ErrorCodeNoResult             ErrorCode = 12000 //客户端数据发送完毕没有收到最后的响应结果看
	ErrorSubServiceNotFound       ErrorCode = 10404 //客户端数据发送完毕没有收到最后的响应结果看
	ErrorNotFound                 ErrorCode = 10404 //客户端数据发送完毕没有收到最后的响应结果看
	ErrorNotFoundSvc              ErrorCode = 10225
	ErrorReqError                 ErrorCode = 10240
	ErrorDoReq                    ErrorCode = 10250
	ErrCodeInvalidSession         ErrorCode = 11345
	//
)

type HttpError struct {
	Code     ErrorCode
	HttpCode int
	Message  string
}

func (e *HttpError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%v|%s", e.Code, e.Message)
}

func NewHttpError(code ErrorCode, httpCode int, msg string) error {
	return &HttpError{Code: code, HttpCode: httpCode, Message: msg}
}
