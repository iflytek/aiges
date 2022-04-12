/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package flange

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	// inner-span split byte
	innerSplit = "\u001D"
	// outer-span split byte
	outerSplit = "\u001E"
)

var (
	// cluster and cluster length
	onceBCluster    = ""
	onceBClusterLen = 0

	// idc and idc length
	onceIDC    = ""
	onceIDCLen = 0

	// ip, port, svc, may move to init_global_config()
	onceIPRune      = ""
	oncePort        string
	oncePortLen     = 0
	onceServiceName = ""

	// atomic value for new trace id
	atomicTraceId uint64 = 0

	// tag zero for new trace id
	tagZero = []byte("00000000")

	// byte buffer pool for `ToString()`
	byteBufferPool = &sync.Pool{
		New: func() interface{} {
			var buffer bytes.Buffer
			buffer.Grow(16384)
			return &buffer
		},
	}
)

// Span define span as a fixed size `map[string]string`
type Span struct {
	// identifier for a trace, set on all items within it.
	traceId []byte

	// span name , rpc method for example.
	name string

	// identifier of this span within a trace.
	// Id []byte
	// span id for address & ts
	spanIdTs []byte

	// short span id for compute
	spanIdHierarchy []byte

	// epoch microseconds of the start of this span.
	timestamp int64

	// measurement in microseconds of the critical path.
	duration int64

	// annotations
	// annotations map[string]int64
	annos string

	// tags
	// tags map[string]string
	tags string

	// other fields for build span
	// span type.
	spanType int32

	// current child id
	currentChildId int32

	// temp meta
	tmpMeta []byte
}

// NewSpan Creates root span (default span type is SERVER).
// @param ip:port:serverName is your server deploy info,
// @param spanType span type, [server|client]
// @param abandon sample switch, true for sample, false for non-sample
func NewSpan(spanType int32, abandon bool) *Span {
	span := getSpan()

	if abandon {
		span.traceId[25] = 'a'
	} else {
		span.traceId[25] = 'n'
	}
	span.timestamp = CurrentTimeMicros()
	span.spanType = spanType
	span.currentChildId = 0
	span.spanIdHierarchy = append(span.spanIdHierarchy, '0')

	return span
}

// Next creates child span.
func (span *Span) Next(spanType int32) *Span {
	nextSpan := getSpan()
	if span == nil || nextSpan == nil {
		return nil
	}

	nextSpan.traceId = span.traceId
	nextSpan.timestamp = CurrentTimeMicros()
	nextSpan.spanType = spanType
	nextSpan.spanIdHierarchy = append(nextSpan.spanIdHierarchy, span.spanIdHierarchy...)
	nextSpan.spanIdHierarchy = append(nextSpan.spanIdHierarchy, '.')
	// TODO if < 10, may do 1+40 = '1'
	nextSpan.spanIdHierarchy = strconv.AppendInt(nextSpan.spanIdHierarchy, int64(atomic.AddInt32(&span.currentChildId, 1)), 10)

	return nextSpan
}

// FromMeta creates span from tmpMeta.
// @param tmpMeta tmpMeta info retrieve with rpc
// @param ip:port:serverName deploy server info
// @spanType span type, [server|client]
func FromMeta(meta string, spanType int32) *Span {
	// tmpMeta should like `c0a8380115174896423630105n008088#xxx`, at least > 34
	mLen := len(meta)
	if mLen < 34 {
		return nil
	}

	span := getSpan()
	if span == nil {
		return nil
	}

	span.tmpMeta = []byte(meta)
	// check traceId
	if mLen >= 32 {
		span.traceId = span.tmpMeta[0:32]
	} else {
		errorf("invalid tmpMeta traceId in %s", span.tmpMeta)
		return nil
	}
	// check spanIdTs
	if mLen > 33 && mLen >= 33+int(span.tmpMeta[32]) {
		span.spanIdTs = span.tmpMeta[33 : 33+int(span.tmpMeta[32])]
	} else {
		errorf("invalid spanIdTs in %s", span.tmpMeta)
		return nil
	}
	// check spanHierarchy
	if mLen >= 33+int(span.tmpMeta[32]) {
		span.spanIdHierarchy = span.tmpMeta[33+int(span.tmpMeta[32]):]
	} else {
		errorf("invalid spanIdHierarchy in %s", span.tmpMeta)
		return nil
	}

	span.timestamp = CurrentTimeMicros()
	span.spanType = spanType
	span.currentChildId = 0

	return span
}

// WithName set span name(rpc method name).
func (span *Span) WithName(name string) *Span {
	if span == nil {
		return nil
	}

	span.name = name
	return span
}

// Start records the start timestamp of rpc span.
func (span *Span) Start() *Span {
	if span == nil {
		return nil
	}

	span.timestamp = CurrentTimeMicros()

	span.annos = span.annos +
			outerSplit + getStartAnnatationType(span.spanType) + innerSplit + strconv.Itoa(int(span.timestamp))
	return span
}

// End records the duration of rpc span.
func (span *Span) End() *Span {
	if span == nil {
		return nil
	}

	ts := CurrentTimeMicros()
	span.duration = ts - span.timestamp

	span.annos = span.annos +
			outerSplit + getEndAnnatationType(span.spanType) + innerSplit + strconv.Itoa(int(ts))
	return span
}

// Send records the message send timestamp of mq span.
func (span *Span) Send() *Span {
	if span == nil {
		return nil
	}

	span.timestamp = CurrentTimeMicros()
	span.duration = 0

	span.annos = span.annos +
			outerSplit + getStartAnnatationType(span.spanType) + innerSplit + strconv.Itoa(int(span.timestamp))
	return span
}

// Recv records the message receive timestamp of mq span.
func (span *Span) Recv() *Span {
	if span == nil {
		return nil
	}

	span.timestamp = CurrentTimeMicros()
	span.duration = 0

	span.annos = span.annos +
			outerSplit + getEndAnnatationType(span.spanType) + innerSplit + strconv.Itoa(int(span.timestamp))
	return span
}

// WithTag set custom tag.
func (span *Span) WithTag(key string, value string) *Span {
	if span == nil {
		return nil
	}

	span.tags = span.tags +
			outerSplit + key + innerSplit + value
	return span
}

// WithRetTag set ret
func (span *Span) WithRetTag(value string) *Span {
	if span == nil {
		return nil
	}

	span.tags = span.tags +
			outerSplit + "ret" + innerSplit + value
	return span
}

// WithErrorTag set error
func (span *Span) WithErrorTag(value string) *Span {
	if span == nil {
		return nil
	}

	span.tags = span.tags +
			outerSplit + "error" + innerSplit + value
	return span
}

// WithLocalComponent set local component.
func (span *Span) WithLocalComponent() *Span {
	if span == nil {
		return nil
	}

	span.tags = span.tags +
			outerSplit + "lc" + innerSplit + "true"
	return span
}

// WithRpcComponent set rpc tag.
func (span *Span) WithRpcComponent() *Span {
	if span == nil {
		return nil
	}

	span.tags = span.tags +
			outerSplit + "rpc" + innerSplit + "true"
	return span
}

// WithRpcCallType set 'call.type' to 'rpc'
func (span *Span) WithRpcCallType() *Span {
	if span == nil {
		return nil
	}

	span.tags = span.tags +
			outerSplit + "call.type" + innerSplit + "rpc"
	return span
}

// WithClientAddr set client address.
func (span *Span) WithClientAddr() *Span {
	if span == nil {
		return nil
	}

	span.tags = span.tags +
			outerSplit + "ca" + innerSplit + "true"
	return span
}

// WithServerAddr set server address.
func (span *Span) WithServerAddr() *Span {
	if span == nil {
		return nil
	}

	span.tags = span.tags +
			outerSplit + "sa" + innerSplit + "true"
	return span
}

// WithMessageAddr set message address.
func (span *Span) WithMessageAddr() *Span {
	if span == nil {
		return nil
	}

	span.tags = span.tags +
			outerSplit + "ma" + innerSplit + "true"
	return span
}

// WithDescf set desc with format, `fmt.Sprintf` will cause performance
// Deprecated: this function will cause performance duo to `fmt.Printf()`
func (span *Span) WithDescf(format string, values ...interface{}) *Span {
	if span == nil {
		return nil
	}

	return span
}

// public property functions
// Meta gets tmpMeta string, format: <traceId>#<id>.
func (span *Span) Meta() string {
	if span == nil {
		return ""
	}

	span.tmpMeta = span.tmpMeta[:0]
	span.tmpMeta = append(span.tmpMeta, span.traceId...)
	span.tmpMeta = append(span.tmpMeta, byte(len(span.spanIdTs)))
	span.tmpMeta = append(span.tmpMeta, span.spanIdTs...)
	span.tmpMeta = append(span.tmpMeta, span.spanIdHierarchy...)
	return string(span.tmpMeta)
}

// ToString convert to string in json.
func (span *Span) ToString() string {
	if span == nil {
		return ""
	}

	buffer := byteBufferPool.Get().(*bytes.Buffer)

	// basic fields
	buffer.WriteString("{")
	buffer.WriteString("\"traceId\":\"" + string(span.traceId) + "\",")
	buffer.WriteString("\"name\":\"" + span.name + "\",")
	buffer.WriteString("\"id\":\"" + string(span.spanIdTs) + string(span.spanIdHierarchy) + "\",")
	buffer.WriteString("\"timestamp\":" + strconv.FormatInt(span.timestamp, 10) + ",")
	buffer.WriteString("\"duration\":" + strconv.FormatInt(span.duration, 10) + ",")

	// annotations
	buffer.WriteString("\"annotations\":[")
	annoBuffer := byteBufferPool.Get().(*bytes.Buffer)
	for _, anno := range strings.Split(span.annos[1:], outerSplit) {
		annoSplit := strings.Split(anno, innerSplit)

		annoBuffer.WriteString("{\"timestamp\":")
		annoBuffer.WriteString(annoSplit[1])
		annoBuffer.WriteString(", \"value\":\"")
		annoBuffer.WriteString(annoSplit[0])
		annoBuffer.WriteString("\", \"endpoint\":{\"serviceName\":\"")
		annoBuffer.WriteString(onceServiceName)
		annoBuffer.WriteString("\", \"ip\":")
		annoBuffer.WriteString(onceIPRune)
		annoBuffer.WriteString(", \"port\":")
		annoBuffer.WriteString(oncePort)
		annoBuffer.WriteString("}},")
	}
	annotationStr := annoBuffer.String()
	annoBuffer.Reset()
	byteBufferPool.Put(annoBuffer)
	if len(annotationStr) > 0 {
		buffer.WriteString(annotationStr[:len(annotationStr)-1])
	}
	buffer.WriteString("],")

	// tags
	buffer.WriteString("\"tags\":[")
	tagBuffer := byteBufferPool.Get().(*bytes.Buffer)
	for _, tag := range strings.Split(span.tags[1:], outerSplit) {
		tagSplit := strings.Split(tag, string(innerSplit))

		tagBuffer.WriteString("{\"key\":\"")
		tagBuffer.WriteString(strings.Replace(tagSplit[0], "\"", "\\\"", -1))
		tagBuffer.WriteString("\", \"value\":\"")
		tagBuffer.WriteString(strings.Replace(tagSplit[1], "\"", "\\\"", -1))
		tagBuffer.WriteString("\", \"endpoint\":{\"serviceName\":\"")
		tagBuffer.WriteString(onceServiceName)
		tagBuffer.WriteString("\", \"ip\":")
		tagBuffer.WriteString(onceIPRune)
		tagBuffer.WriteString(", \"port\":")
		tagBuffer.WriteString(oncePort)
		tagBuffer.WriteString("}},")
	}
	tagStr := tagBuffer.String()
	tagBuffer.Reset()
	byteBufferPool.Put(tagBuffer)
	if len(tagStr) > 0 {
		buffer.WriteString(tagStr[:len(tagStr)-1])
	}
	buffer.WriteString("]}")

	result := buffer.String()
	buffer.Reset()
	byteBufferPool.Put(buffer)
	return result
}

// ==== other functions ====
// setGlobalConfig for global value pre-set
func setGlobalConfig(bcluster string, idc string, serviceIP string, servicePort string, serviceName string) {
	// set the global value
	onceBClusterLen = len(bcluster)
	if onceBClusterLen > 4 {
		infof("len(%s) > 4, will be cut", bcluster)
		onceBCluster = bcluster[0:4]
		onceBClusterLen = 4
	} else {
		onceBCluster = bcluster
	}

	onceIDCLen = len(idc)
	if onceIDCLen > 4 {
		infof("len(%s) > 4, will be cut", idc)
		onceIDC = idc[0:4]
		onceIDCLen = 4
	} else {
		onceIDC = idc
	}

	onceServiceName = serviceName
	if len(serviceIP) > 0 {
		onceIPRune = strconv.FormatInt(int64(IPv4toRune(serviceIP)), 10)
	}
	if len(servicePort) > 0 {
		oncePort = servicePort
		oncePortLen = len(oncePort)
	}

	// prepare fixed endpoint buf
	endpointString = onceServiceName + innerSplit + onceIPRune + innerSplit + oncePort
}

// create a new span with default value
func getSpan() *Span {
	// init a raw span
	span := &Span{
		traceId:         make([]byte, 0, 32),
		spanIdTs:        make([]byte, 0, 32),
		spanIdHierarchy: make([]byte, 0, 8),
		tmpMeta:         make([]byte, 0, 40),
	}

	microsTs := CurrentTimeMicros()
	millisTs := microsTs / 1000

	// third, set & reset span
	span.traceId = span.traceId[:0]
	{
		// new trace id logical
		span.traceId = append(span.traceId, tagZero[0:4-onceBClusterLen]...)
		span.traceId = append(span.traceId, onceBCluster...)

		span.traceId = append(span.traceId, tagZero[0:4-onceIDCLen]...)
		span.traceId = append(span.traceId, onceIDC...)

		span.traceId = strconv.AppendInt(span.traceId, millisTs, 10)
		// span.TraceId = appendInt(span.TraceId, genGlobalMillisTs)

		// replace Rand with auto-increased value
		rand := atomic.AddUint64(&atomicTraceId, 1) % 10000
		var tempRandStrLen = 0
		if rand > 999 {
			tempRandStrLen = 4
		} else if rand > 99 {
			tempRandStrLen = 3
		} else if rand > 9 {
			tempRandStrLen = 2
		} else {
			tempRandStrLen = 1
		}
		span.traceId = append(span.traceId, tagZero[0:4-tempRandStrLen]...)
		span.traceId = strconv.AppendInt(span.traceId, int64(rand), 10)
		// span.TraceId = appendInt(span.TraceId, int64(rand))

		span.traceId = append(span.traceId, 'a')

		span.traceId = append(span.traceId, tagZero[0:6-oncePortLen]...)
		span.traceId = append(span.traceId, oncePort...)
	}
	span.name = ""
	span.timestamp = microsTs
	// id with format [&address|ts|0.1.1]
	span.spanIdTs = span.spanIdTs[:0]
	{
		// span id logical
		span.spanIdTs = strconv.AppendInt(span.spanIdTs, int64(uintptr(unsafe.Pointer(&span.traceId))), 16)
		// span.spanIdTs = appendInt(span.spanIdTs, int64(uintptr(unsafe.Pointer(&span.TraceId))))
		span.spanIdTs = strconv.AppendInt(span.spanIdTs, span.timestamp, 10)
		// span.spanIdTs = appendInt(span.spanIdTs, span.Timestamp)
	}
	span.duration = 0
	span.spanType = SERVER
	span.currentChildId = 0
	span.spanIdHierarchy = span.spanIdHierarchy[:0]

	// add span version tag
	span.WithTag("span_version", FLANGE_VERSION)

	return span
}

// gets annotation type while rpc start or message send.
func getStartAnnatationType(spanType int32) string {
	if spanType == CLIENT {
		return CLIENT_SEND
	} else if spanType == SERVER {
		return SERVER_RECV
	} else {
		return MESSAGE_SEND
	}
}

// gets annotation type while rpc end or message receive.
func getEndAnnatationType(spanType int32) string {
	if spanType == CLIENT {
		return CLIENT_RECV
	} else if spanType == SERVER {
		return SERVER_SEND
	} else {
		return MESSAGE_RECV
	}
}

// dump span to file.
// ${dir}/${traceId}_${spanId}
func (span *Span) dump(dir string) {
	filename := dir + string(os.PathSeparator) + string(span.traceId) + "_" + string(span.spanIdTs) + string(span.spanIdHierarchy)
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		errorf("open %s failed: %v\n", filename, err)
		return
	}

	// Dump span
	fmt.Fprintln(fp, span.ToString())
	fp.Close()
}
