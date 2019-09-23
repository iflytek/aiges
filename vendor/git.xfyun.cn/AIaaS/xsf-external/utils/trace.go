package utils

//
type SpanType int

//
const (
	// unknown type.
	UnknowSpan SpanType = iota
	// for CLIENT_SEND or CLIENT_RECV.
	CliSpan
	// for SERVER_RECV or SERVER_SEND.
	SrvSpan
	// for MESSAGE_SEND.
	ProduceSpan
	// for MESSAGE_RECV.
	ConsumerSpan
)

const ignoreSvc = "lbv2"

var (
	enableTrace  = true
	abandonTrace = false
)

/*CfgOption:
配置选项
*/
type TracerOption struct {
	dump    bool
	deliver bool
	b       int
	log     *Logger
	spill   string

	buffer              int
	batch               int
	linger              int
	spillAble           bool
	watchAble           bool //trace的监控接口
	watchPort           int  //trace的监控端口
	MaxSpillContentSize int  //spill的最大大小
	LowLoadSleepTs      int  //span池的填充速度

	svcIp   string
	svcPort int32
	svcName string

	bcluster string
	idc      string
}

type TraceOpt func(*TracerOption)

type Span struct {
}

// creates span
// @param meta meta info retrieve with rpc
// @param ip:port:serverName deploy server info
// @spanType span type, [server|client]
func NewSpan(spanType SpanType) *Span {
	return &Span{}

}

func NewSpanFromMeta(meta string, spanType SpanType) *Span {
	return nil
}

// creates span from meta.
// @param meta meta info retrieve with rpc
// @param ip:port:serverName deploy server info
// @spanType span type, [server|client]
func FromMeta(meta string, ip string, port int32, serviceName string, spanType SpanType) *Span {
	return nil
}

// set span name(rpc method name).
func (span *Span) WithName(name string) *Span {
	return nil
}

// set tag.
func (span *Span) WithTag(key string, value string) *Span {
	return nil
}

// set function
func (span *Span) WithFuncTag(value string) *Span {
	return nil
}

func (span *Span) WithRpcComponent() *Span {
	return nil
}

func (span *Span) WithRpcCallType() *Span {
	return nil
}

// set ret
func (span *Span) WithRetTag(value string) *Span {
	return nil
}

// sert error
func (span *Span) WithErrorTag(value string) *Span {
	return nil
}

// set local component.
func (span *Span) WithLocalComponent() *Span {
	return nil
}

// set client address.
func (span *Span) WithClientAddr() *Span {
	return nil
}

// set server address.
func (span *Span) WithServerAddr() *Span {
	return nil
}

// set message address.
func (span *Span) WithMessageAddr() *Span {
	return nil
}

// records the start timestamp of rpc span.
func (span *Span) Start() *Span {
	return nil
}

// records the duration of rpc span.
func (span *Span) End() *Span {
	return nil
}

// records the message send timestamp of mq span.
func (span *Span) Send() *Span {
	return nil
}

// records the message receive timestamp of mq span.
func (span *Span) Recv() *Span {
	return nil
}

// creates child span.
func (span *Span) Next(st SpanType) *Span {
	return nil
}

// gets meta string, format: <traceId>#<id>.
func (span *Span) Meta() string {
	return ""
}

// convert to string in json.
func (span *Span) ToString() string {
	return ""
}

// convert to string in json.
func (span *Span) Flush() {
	return

}
