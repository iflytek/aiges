// 错误码分类,各模块分段;
package frame

import "errors"

// errNum list
const (
	AigesSuccess = 0

	// 服务实例错误码
	AigesErrorInvalidOp        = 10003
	AigesErrorInvalidSessMode  = 10004
	AigesErrorInvalidPara      = 10106
	AigesErrorInvalidParaValue = 10107
	AigesErrorInvalidHdl       = 10008
	AigesErrorInvalidData      = 10009
	AigesErrorLicNotEnough     = 10110
	AigesErrorDpInit           = 10011
	AigesErrorSessTimeout      = 10019
	AigesErrorEngInactive      = 10101 // TODO 兼容保留
	AigesErrorOnceTimeout      = 10102

	// 排序缓冲区错误
	AigesErrorBufferEmpty   = 10300
	AigesErrorSeqChanClosed = 10301

	// 协议错误
	AigesErrorPbMarshal   = 10400
	AigesErrorPbUnmarshal = 10401

	// 内部同步错误
	AigesErrorFinRoutine = 10500

	// 事件异常错误
	AigesErrorNilEvent = 10600

	// wrapper接口异常
	WrapperInitErr    = 100001
	WrapperFiniErr    = 100002
	WrapperLoadErr    = 100003
	WrapperUnloadErr  = 100004
	WrapperCreateErr  = 100005
	WrapperDestroyErr = 100006
	WrapperWriteErr   = 100007
	WrapperReadErr    = 100008
	WrapperExecErr    = 100009
	WrapperAsyncErr   = 100010
)

// errInfo list
var (
	// Error information of EngService
	ErrorInvalidSampleRate = errors.New("invalid audio sampleRate")
	ErrorSessNotSupport    = errors.New("not support session mode")

	// Error of Aiges instance & manager module
	ErrorInvalidParaValue = errors.New("service para value invalid")
	ErrorInvalidInstHdl  = errors.New("service instance invalid")
	ErrorInstNotEnouth   = errors.New("service license not enough")
	ErrorInstNotActive   = errors.New("service instance inactive, already release")
	ErrorInstRwTimeout   = errors.New("service read buffer timeout, session timeout")
	ErrorOnceExecTimeout = errors.New("service exec once request timeout")

	// Error of Aiges protocol module
	ErrorInvalidOp   = errors.New("invalid service operation")
	ErrorInvalidData = errors.New("input invalid data")
	ErrorPbMarshal   = errors.New("marshal pb message fail")
	ErrorPbUnmarshal = errors.New("unmarshal pb message fail")

	// Error of Aiges configure module
	ErrorGetUsrConfigure = errors.New("can't get user configure")
	ErrorInvalidUsrCfg   = errors.New("invalid user configure data")

	// Error of Aiges buffer module
	ErrorSeqBufferEmpty = errors.New("seqBuffer empty right now")
	ErrorSeqChanClosed  = errors.New("seq channel already closed")

	// Error of Aiges Codec module
	ErrorAudioCodingStart  = errors.New("audioCoding start fail")
	ErrorAudioCodingStop   = errors.New("audioCoding stop fail")
	ErrorAudioCodingEncode = errors.New("audioCoding encode fail")
	ErrorAudioCodingDecode = errors.New("audioCoding decode fail")

	// Error of Aiges DataProcess module
	ErrorAudioResampleInit = errors.New("audio reSample init fail")

	// Error of Aiges inner sync
	ErrorFinishRoutine = errors.New("finish instance go routine")

	ErrorNilRegEvent = errors.New("can't find reg event to exec request")

	// Error personal resource system
	ErrorAsyncResTimeout = errors.New("async query resource timeout")
)
