// 错误码分类,各模块分段;
package frame

import "errors"

// errNum list
const (
	AigesSuccess = 0

	// 服务实例错误码
	AigesErrorInvalidOp        = 10003
	AigesErrorInvalidSessMode  = 10004
	AigesErrorInvalidHdl       = 10008
	AigesErrorInvalidData      = 10009
	AigesErrorHttpReq          = 10010
	AigesErrorHttpFail         = 10011
	AigesErrorHttpInvalidData  = 10012
	AigesErrorHttpTimeout      = 10013
	AigesErrorCodecStart       = 10041
	AigesErrorCodecEncode      = 10042
	AigesErrorCodecDecode      = 10043
	AigesErrorCodecStop        = 10044
	AigesErrorInvalidPara      = 10106
	AigesErrorInvalidParaValue = 10107
	AigesErrorLicNotEnough     = 10110
	AigesErrorDpInit           = 10011
	AigesErrorSessTimeout      = 10019
	AigesErrorEngInactive      = 10101 // TODO 兼容保留
	AigesErrorOnceTimeout      = 10102
	AigesErrorNrtUpdate        = 10103
	AigesErrorElasticLic       = 10104
	AigesErrorSessRepeat       = 10105
	AigesErrorInvalidOut       = 10106
	AigesErrorCrash            = 10109

	// 排序缓冲区错误
	AigesErrorBufferEmpty   = 10300
	AigesErrorSeqChanClosed = 10301

	// 协议错误
	AigesErrorPbMarshal                = 10400
	AigesErrorPbUnmarshal              = 10401
	AigesErrorPbVersion                = 10402
	AigesErrorJsonMarshal              = 10403
	AigesErrorJsonUnmarshal            = 10404
	AigesErrorPbAdapter                = 10405
	AigesErrorPbAdapterIatParamInvalid = 10406
	// 事件异常错误
	AigesErrorNilEvent = 10600

	// wrapper接口异常
	WrapperInitErr    = 10901
	WrapperFiniErr    = 10902
	WrapperLoadErr    = 10903
	WrapperUnloadErr  = 10904
	WrapperCreateErr  = 10905
	WrapperDestroyErr = 10906
	WrapperWriteErr   = 10907
	WrapperReadErr    = 10908
	WrapperExecErr    = 10909
	WrapperAsyncErr   = 10910

	// 插件回调接口异常
	InnerInvalidCustom = 20001
	InnerInvalidHandle = 20002
)

// errInfo list
var (
	// Error information of EngService
	ErrorInvalidSampleRate = errors.New("invalid audio sampleRate")

	// Error of Aiges instance & manager module
	ErrorInvalidParaValue = errors.New("service para value invalid")
	ErrorInvalidInstHdl   = errors.New("service instance invalid")
	ErrorInstNotEnouth    = errors.New("service license not enough")
	ErrorElasticLicInst   = errors.New("service license elastic fail")
	ErrorInvalidStatus    = errors.New("service request status invalid")
	ErrorInstNotActive    = errors.New("service instance inactive, already release")
	ErrorInstRwTimeout    = errors.New("service read buffer timeout, session timeout")
	ErrorOnceExecTimeout  = errors.New("service exec once request timeout")
	ErrorInvalidOutput    = errors.New("wrapper output data invalid(key or type)")
	ErrorNullOutput       = errors.New("wrapper output data is nil ")

	ErrorPersonalIndex = errors.New("personal index service response fail")

	// Error of Aiges protocol module
	ErrorInvalidOp                = errors.New("invalid service operation")
	ErrorInvalidData              = errors.New("input invalid data")
	ErrorPbMarshal                = errors.New("marshal pb message fail")
	ErrorPbUnmarshal              = errors.New("unmarshal pb message fail")
	ErrorPbVersion                = errors.New("pb message version invalid")
	ErrorPbAdapter                = errors.New("pb message adapter fail")
	ErrorPbAdapterIatParamInvalid = errors.New("pb message adapter param invalid")
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

	ErrorNilRegEvent = errors.New("can't find reg event to exec request")

	// Error personal resource system
	ErrorAsyncResTimeout = errors.New("async query resource timeout")
)
