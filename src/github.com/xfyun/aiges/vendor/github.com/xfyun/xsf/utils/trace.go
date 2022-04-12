package utils

import (
	"fmt"
	//flange "genitus/flange/mock"
	"github.com/xfyun/flange"
)

type SpanType int

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

	abandon bool
}

type TraceOpt func(*TracerOption)

func WithAbandon(a bool) TraceOpt {
	return func(c *TracerOption) {
		c.abandon = a
	}
}
func WithLowLoadSleepTs(t int) TraceOpt {
	return func(c *TracerOption) {
		c.LowLoadSleepTs = t
	}
}
func WithMaxSpillContentSize(s int) TraceOpt {
	return func(c *TracerOption) {
		c.MaxSpillContentSize = s
	}
}
func WithWatchPort(p int) TraceOpt {
	return func(c *TracerOption) {
		c.watchPort = p
	}
}
func WithDump(d bool) TraceOpt {
	return func(c *TracerOption) {
		c.dump = d
	}
}
func WithDeliver(d bool) TraceOpt {
	return func(c *TracerOption) {
		c.deliver = d
	}
}

func WithSpillAble(d bool) TraceOpt {
	return func(c *TracerOption) {
		c.spillAble = d
	}
}

func WithBufferSize(b int) TraceOpt {
	return func(c *TracerOption) {
		c.buffer = b
	}
}

func WithBatchSize(b int) TraceOpt {
	return func(c *TracerOption) {
		c.batch = b
	}
}

func WithWatch(b bool) TraceOpt {
	flange.WatchLogEnable = b
	return func(c *TracerOption) {
		c.watchAble = b
	}
}
func WithLinger(l int) TraceOpt {
	return func(c *TracerOption) {
		c.linger = l
	}
}
func WithBackend(b int) TraceOpt {
	return func(c *TracerOption) {
		c.b = b
	}
}

func WithTraceLogger(l *Logger) TraceOpt {
	return func(c *TracerOption) {
		c.log = l
	}
}

func WithTraceSpill(spill string) TraceOpt {
	return func(c *TracerOption) {
		c.spill = spill
	}
}
func WithSvcIp(ip string) TraceOpt {
	return func(c *TracerOption) {
		c.svcIp = ip
	}
}
func WithSvcPort(port int32) TraceOpt {
	return func(c *TracerOption) {
		c.svcPort = port
	}
}
func WithSvcName(name string) TraceOpt {
	return func(c *TracerOption) {
		c.svcName = name
	}
}
func WithSvcBCluster(bcluster string) TraceOpt {
	return func(c *TracerOption) {
		c.bcluster = bcluster
	}
}

func WithSvcIDC(idc string) TraceOpt {
	return func(c *TracerOption) {
		c.idc = idc
	}
}

func InitTracer(addr string, port string, o ...TraceOpt) error {

	to := &TracerOption{dump: false, deliver: true, b: 4, spill: "/log/spill", buffer: 1000000, batch: 100, linger: 5, spillAble: true}
	for _, opt := range o {
		opt(to)
	}

	if to.svcName == ignoreSvc {
		return nil
	}
	loggerStd.Printf("fn:InitTracer,TracerOption:%+v\n", to)
	if to.log != nil {
		flange.Logger = NewLogImpl(to.log)
	}
	flange.DumpEnable = to.dump
	flange.DeliverEnable = to.deliver
	flange.SpillDir = to.spill
	flange.BuffSize = int32(to.buffer)
	flange.BatchSize = to.batch
	flange.LingerSec = to.linger
	flange.SpillEnable = to.spillAble
	flange.WatchLogPort = to.watchPort
	flange.MaxSpillContentSize = int64(to.MaxSpillContentSize)
	flange.LowLoadSleepTs = to.LowLoadSleepTs
	//bug: 如下
	//todo: 这样的处理在客户端比较恶心，有可能得不到有效的ip和端口
	/*if len(to.svcIp) >0 && len(to.svcName) > 0 {
		flange.SetGlobalConfig(to.svcIp, to.svcPort, to.svcName)
	}*/
	//todo 0.19版本的此接口已不复存在，待确认

	//panic(fmt.Sprintf("%v:%v", to.bcluster, to.idc))

	toStr := func(in int32) string {
		if in == 0 {
			return ""
		}
		return fmt.Sprintf("%d", in)
	}

	return flange.Init(addr, port, to.b, to.bcluster, to.idc, to.svcIp, toStr(to.svcPort), to.svcName)
}

func AbleTrace(able bool) {
	loggerStd.Printf("fn:AbleTrace,able:%v\n", able)
	enableTrace = able
}
func Abandon(able bool) {
	abandonTrace = able
}

type Span struct {
	sp *flange.Span
}

// creates span
// @param meta meta info retrieve with rpc
// @param ip:port:serverName deploy server info
// @spanType span type, [server|client]
func NewSpan(spanType SpanType) *Span {
	if !enableTrace {
		return nil
	}
	//trace自0.19后，接口变动，此处移除ip，port，serviceName
	sp := flange.NewSpan(int32(spanType), abandonTrace)
	if sp == nil {
		return nil
	}

	return &Span{sp: sp}

}
func NewSpanWithOpts(spanType SpanType, opts ...TraceOpt) *Span {
	to := &TracerOption{abandon: false}
	for _, opt := range opts {
		opt(to)
	}
	if !enableTrace {
		return nil
	}
	//trace自0.19后，接口变动，此处移除ip，port，serviceName
	sp := flange.NewSpan(int32(spanType), to.abandon)
	if sp == nil {
		return nil
	}

	return &Span{sp: sp}

}
func NewSpanFromMeta(meta string, spanType SpanType) *Span {
	if !enableTrace {
		return nil
	}
	//trace自0.19后，接口变动，此处移除ip，port，serviceName
	sp := flange.FromMeta(meta, int32(spanType))
	if sp == nil {
		return nil
	}
	return &Span{sp: sp}
}

// creates span from meta.
// @param meta meta info retrieve with rpc
// @param ip:port:serverName deploy server info
// @spanType span type, [server|client]
func FromMeta(meta string, ip string, port int32, serviceName string, spanType SpanType) *Span {
	if !enableTrace {
		return nil
	}
	//trace自0.19后，接口变动，此处移除ip，port，serviceName
	sp := flange.FromMeta(meta, int32(spanType))
	if sp == nil {
		return nil
	}
	return &Span{sp: sp}

}

// set span name(rpc method name).
func (span *Span) WithName(name string) *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.WithName(name)
	return span
}

// set tag.
func (span *Span) WithTag(key string, value string) *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.WithTag(key, value)
	return span
}

// set function
func (span *Span) WithFuncTag(value string) *Span {
	if !enableTrace || span == nil {
		return span
	}
	return span.WithTag("func", value)
}

func (span *Span) WithRpcComponent() *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.WithRpcComponent()
	return span
}

func (span *Span) WithRpcCallType() *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.WithRpcCallType()
	return span
}

// set ret
func (span *Span) WithRetTag(value string) *Span {
	if !enableTrace || span == nil {
		return span
	}
	return span.WithTag("ret", value)
}

// sert error
func (span *Span) WithErrorTag(value string) *Span {
	if !enableTrace || span == nil {
		return span
	}
	return span.WithTag("error", value)
}

// set local component.
func (span *Span) WithLocalComponent() *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.WithLocalComponent()
	return span
}

// set client address.
func (span *Span) WithClientAddr() *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.WithClientAddr()
	return span
}

// set server address.
func (span *Span) WithServerAddr() *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.WithServerAddr()
	return span
}

// set message address.
func (span *Span) WithMessageAddr() *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.WithServerAddr()
	return span
}

// records the start timestamp of rpc span.
func (span *Span) Start() *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.Start()
	return span
}

// records the duration of rpc span.
func (span *Span) End() *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.End()
	return span
}

// records the message send timestamp of mq span.
func (span *Span) Send() *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.Send()
	return span
}

// records the message receive timestamp of mq span.
func (span *Span) Recv() *Span {
	if !enableTrace || span == nil {
		return span
	}
	span.sp.Recv()
	return span
}

// creates child span.
func (span *Span) Next(st SpanType) *Span {
	if !enableTrace || span == nil {
		return span
	}
	return &Span{sp: span.sp.Next(int32(st))}
}

// gets meta string, format: <traceId>#<id>.
func (span *Span) Meta() string {
	if !enableTrace || span == nil {
		return ""
	}
	return span.sp.Meta()
}

// convert to string in json.
func (span *Span) ToString() string {
	if !enableTrace || span == nil {
		return ""
	}
	return span.sp.ToString()
}

// convert to string in json.
func (span *Span) Flush() {
	if !enableTrace || span == nil {
		return
	}
	flange.Flush(span.sp)

}
