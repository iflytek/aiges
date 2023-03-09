package stream

const (
	KeySid          = "sid"
	KeyUid          = "uid"
	KeyAppId        = "appid"
	KeyClientIp     = "client_ip"
	KeyCallBackAddr = "callback_addr"
	KeyStatus       = "status"
	KeyTraceId      = "trace_id"
	KeyServiceId    = "service_id"
	KeySeqNo        = "seq_no"
	KeyCloudId      = "cloud_id"
	KeyRoute        = "route"
	KeyAipaaSSUb    = "aipaas_sub"
	keyConnId       = "cid"
	KeySub          = "sub"
)

const (
	StatusBegin = iota
	StatusContinue
	StatusEnd
	StatusOnce
)
