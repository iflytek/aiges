package server

import "github.com/golang/protobuf/proto"
import "fmt"
import "math"

import "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// 用于消息类型
type ServerBiz_MsgType int32

const (
	ServerBiz_UP_CALL     ServerBiz_MsgType = 0
	ServerBiz_UP_RESULT   ServerBiz_MsgType = 1
	ServerBiz_DOWN_CALL   ServerBiz_MsgType = 2
	ServerBiz_DOWN_RESULT ServerBiz_MsgType = 3
)

var ServerBiz_MsgType_name = map[int32]string{
	0: "UP_CALL",
	1: "UP_RESULT",
	2: "DOWN_CALL",
	3: "DOWN_RESULT",
}

var ServerBiz_MsgType_value = map[string]int32{
	"UP_CALL":     0,
	"UP_RESULT":   1,
	"DOWN_CALL":   2,
	"DOWN_RESULT": 3,
}

func (x ServerBiz_MsgType) String() string {
	return proto.EnumName(ServerBiz_MsgType_name, int32(x))
}
func (ServerBiz_MsgType) EnumDescriptor() ([]byte, []int) { return fileDescriptorServerBiz, []int{0, 0} }

// 区分数据类型
type GeneralData_DataType int32

const (
	GeneralData_TEXT  GeneralData_DataType = 0
	GeneralData_AUDIO GeneralData_DataType = 1
	GeneralData_IMAGE GeneralData_DataType = 2
	GeneralData_VIDEO GeneralData_DataType = 3
)

var GeneralData_DataType_name = map[int32]string{
	0: "TEXT",
	1: "AUDIO",
	2: "IMAGE",
	3: "VIDEO",
}
var GeneralData_DataType_value = map[string]int32{
	"TEXT":  0,
	"AUDIO": 1,
	"IMAGE": 2,
	"VIDEO": 3,
}

func (x GeneralData_DataType) String() string {
	return proto.EnumName(GeneralData_DataType_name, int32(x))
}
func (GeneralData_DataType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorServerBiz, []int{6, 0}
}

// 区分数据状态
type GeneralData_DataStatus int32

const (
	GeneralData_BEGIN    GeneralData_DataStatus = 0
	GeneralData_CONTINUE GeneralData_DataStatus = 1
	GeneralData_END      GeneralData_DataStatus = 2
	GeneralData_ONCE     GeneralData_DataStatus = 3
)

var GeneralData_DataStatus_name = map[int32]string{
	0: "BEGIN",
	1: "CONTINUE",
	2: "END",
	3: "ONCE",
}
var GeneralData_DataStatus_value = map[string]int32{
	"BEGIN":    0,
	"CONTINUE": 1,
	"END":      2,
	"ONCE":     3,
}

func (x GeneralData_DataStatus) String() string {
	return proto.EnumName(GeneralData_DataStatus_name, int32(x))
}
func (GeneralData_DataStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorServerBiz, []int{6, 1}
}

// 服务业务协议
type ServerBiz struct {
	MsgType     ServerBiz_MsgType `protobuf:"varint,1,opt,name=msg_type,json=msgType,proto3,enum=serverbiz.ServerBiz_MsgType" json:"msg_type,omitempty"`
	Version     string            `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	GlobalRoute *GlobalRoute      `protobuf:"bytes,3,opt,name=global_route,json=globalRoute" json:"global_route,omitempty"`
	UpCall      *UpCall           `protobuf:"bytes,4,opt,name=up_call,json=upCall" json:"up_call,omitempty"`
	UpResult    *UpResult         `protobuf:"bytes,5,opt,name=up_result,json=upResult" json:"up_result,omitempty"`
	DownCall    *DownCall         `protobuf:"bytes,6,opt,name=down_call,json=downCall" json:"down_call,omitempty"`
	DownResult  *DownResult       `protobuf:"bytes,7,opt,name=down_result,json=downResult" json:"down_result,omitempty"`
}

func (m *ServerBiz) Reset()                    { *m = ServerBiz{} }
func (m *ServerBiz) String() string            { return proto.CompactTextString(m) }
func (*ServerBiz) ProtoMessage()               {}
func (*ServerBiz) Descriptor() ([]byte, []int) { return fileDescriptorServerBiz, []int{0} }

func (m *ServerBiz) GetMsgType() ServerBiz_MsgType {
	if m != nil {
		return m.MsgType
	}
	return ServerBiz_UP_CALL
}

func (m *ServerBiz) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ServerBiz) GetGlobalRoute() *GlobalRoute {
	if m != nil {
		return m.GlobalRoute
	}
	return nil
}

func (m *ServerBiz) GetUpCall() *UpCall {
	if m != nil {
		return m.UpCall
	}
	return nil
}

func (m *ServerBiz) GetUpResult() *UpResult {
	if m != nil {
		return m.UpResult
	}
	return nil
}

func (m *ServerBiz) GetDownCall() *DownCall {
	if m != nil {
		return m.DownCall
	}
	return nil
}

func (m *ServerBiz) GetDownResult() *DownResult {
	if m != nil {
		return m.DownResult
	}
	return nil
}

// 路由信息
type GlobalRoute struct {
	SessionId    string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	TraceId      string `protobuf:"bytes,2,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	UpRouterId   string `protobuf:"bytes,3,opt,name=up_router_id,json=upRouterId,proto3" json:"up_router_id,omitempty"`
	GuiderId     string `protobuf:"bytes,4,opt,name=guider_id,json=guiderId,proto3" json:"guider_id,omitempty"`
	DownRouterId string `protobuf:"bytes,5,opt,name=down_router_id,json=downRouterId,proto3" json:"down_router_id,omitempty"`
	Appid        string `protobuf:"bytes,6,opt,name=appid,proto3" json:"appid,omitempty"`
	Uid          string `protobuf:"bytes,7,opt,name=uid,proto3" json:"uid,omitempty"`
	Did          string `protobuf:"bytes,8,opt,name=did,proto3" json:"did,omitempty"`
	ClientIp     string `protobuf:"bytes,9,opt,name=client_ip,json=clientIp,proto3" json:"client_ip,omitempty"`
}

func (m *GlobalRoute) Reset()                    { *m = GlobalRoute{} }
func (m *GlobalRoute) String() string            { return proto.CompactTextString(m) }
func (*GlobalRoute) ProtoMessage()               {}
func (*GlobalRoute) Descriptor() ([]byte, []int) { return fileDescriptorServerBiz, []int{1} }

func (m *GlobalRoute) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *GlobalRoute) GetTraceId() string {
	if m != nil {
		return m.TraceId
	}
	return ""
}

func (m *GlobalRoute) GetUpRouterId() string {
	if m != nil {
		return m.UpRouterId
	}
	return ""
}

func (m *GlobalRoute) GetGuiderId() string {
	if m != nil {
		return m.GuiderId
	}
	return ""
}

func (m *GlobalRoute) GetDownRouterId() string {
	if m != nil {
		return m.DownRouterId
	}
	return ""
}

func (m *GlobalRoute) GetAppid() string {
	if m != nil {
		return m.Appid
	}
	return ""
}

func (m *GlobalRoute) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *GlobalRoute) GetDid() string {
	if m != nil {
		return m.Did
	}
	return ""
}

func (m *GlobalRoute) GetClientIp() string {
	if m != nil {
		return m.ClientIp
	}
	return ""
}

// 上行数据请求
type UpCall struct {
	Call         string            `protobuf:"bytes,1,opt,name=call,proto3" json:"call,omitempty"`
	SeqNo        int32             `protobuf:"zigzag32,2,opt,name=seq_no,json=seqNo,proto3" json:"seq_no,omitempty"`
	From         string            `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	Sync         bool              `protobuf:"varint,4,opt,name=sync,proto3" json:"sync,omitempty"`
	BusinessArgs map[string]string `protobuf:"bytes,5,rep,name=business_args,json=businessArgs" json:"business_args,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TempArgs     map[string]string `protobuf:"bytes,6,rep,name=temp_args,json=tempArgs" json:"temp_args,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	UserArgs     map[string][]byte `protobuf:"bytes,7,rep,name=user_args,json=userArgs" json:"user_args,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Session      map[string]string `protobuf:"bytes,8,rep,name=session" json:"session,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	DataList     []*GeneralData    `protobuf:"bytes,9,rep,name=data_list,json=dataList" json:"data_list,omitempty"`
}

func (m *UpCall) Reset()                    { *m = UpCall{} }
func (m *UpCall) String() string            { return proto.CompactTextString(m) }
func (*UpCall) ProtoMessage()               {}
func (*UpCall) Descriptor() ([]byte, []int) { return fileDescriptorServerBiz, []int{2} }

func (m *UpCall) GetCall() string {
	if m != nil {
		return m.Call
	}
	return ""
}

func (m *UpCall) GetSeqNo() int32 {
	if m != nil {
		return m.SeqNo
	}
	return 0
}

func (m *UpCall) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *UpCall) GetSync() bool {
	if m != nil {
		return m.Sync
	}
	return false
}

func (m *UpCall) GetBusinessArgs() map[string]string {
	if m != nil {
		return m.BusinessArgs
	}
	return nil
}

func (m *UpCall) GetTempArgs() map[string]string {
	if m != nil {
		return m.TempArgs
	}
	return nil
}

func (m *UpCall) GetUserArgs() map[string][]byte {
	if m != nil {
		return m.UserArgs
	}
	return nil
}

func (m *UpCall) GetSession() map[string]string {
	if m != nil {
		return m.Session
	}
	return nil
}

func (m *UpCall) GetDataList() []*GeneralData {
	if m != nil {
		return m.DataList
	}
	return nil
}

// 上行结果ack
type UpResult struct {
	Ret      int32             `protobuf:"zigzag32,1,opt,name=ret,proto3" json:"ret,omitempty"`
	AckNo    int32             `protobuf:"zigzag32,2,opt,name=ack_no,json=ackNo,proto3" json:"ack_no,omitempty"`
	ErrInfo  string            `protobuf:"bytes,3,opt,name=err_info,json=errInfo,proto3" json:"err_info,omitempty"`
	From     string            `protobuf:"bytes,4,opt,name=from,proto3" json:"from,omitempty"`
	Session  map[string]string `protobuf:"bytes,5,rep,name=session" json:"session,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	DataList []*GeneralData    `protobuf:"bytes,6,rep,name=data_list,json=dataList" json:"data_list,omitempty"`
}

func (m *UpResult) Reset()                    { *m = UpResult{} }
func (m *UpResult) String() string            { return proto.CompactTextString(m) }
func (*UpResult) ProtoMessage()               {}
func (*UpResult) Descriptor() ([]byte, []int) { return fileDescriptorServerBiz, []int{3} }

func (m *UpResult) GetRet() int32 {
	if m != nil {
		return m.Ret
	}
	return 0
}

func (m *UpResult) GetAckNo() int32 {
	if m != nil {
		return m.AckNo
	}
	return 0
}

func (m *UpResult) GetErrInfo() string {
	if m != nil {
		return m.ErrInfo
	}
	return ""
}

func (m *UpResult) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *UpResult) GetSession() map[string]string {
	if m != nil {
		return m.Session
	}
	return nil
}

func (m *UpResult) GetDataList() []*GeneralData {
	if m != nil {
		return m.DataList
	}
	return nil
}

// 下行数据请求
type DownCall struct {
	Ret      int32             `protobuf:"zigzag32,1,opt,name=ret,proto3" json:"ret,omitempty"`
	SeqNo    int32             `protobuf:"zigzag32,2,opt,name=seq_no,json=seqNo,proto3" json:"seq_no,omitempty"`
	From     string            `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	Args     map[string]string `protobuf:"bytes,4,rep,name=args" json:"args,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	DataList []*GeneralData    `protobuf:"bytes,5,rep,name=data_list,json=dataList" json:"data_list,omitempty"`
}

func (m *DownCall) Reset()                    { *m = DownCall{} }
func (m *DownCall) String() string            { return proto.CompactTextString(m) }
func (*DownCall) ProtoMessage()               {}
func (*DownCall) Descriptor() ([]byte, []int) { return fileDescriptorServerBiz, []int{4} }

func (m *DownCall) GetRet() int32 {
	if m != nil {
		return m.Ret
	}
	return 0
}

func (m *DownCall) GetSeqNo() int32 {
	if m != nil {
		return m.SeqNo
	}
	return 0
}

func (m *DownCall) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *DownCall) GetArgs() map[string]string {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *DownCall) GetDataList() []*GeneralData {
	if m != nil {
		return m.DataList
	}
	return nil
}

// 下行结果ack
type DownResult struct {
	Ret     int32             `protobuf:"zigzag32,1,opt,name=ret,proto3" json:"ret,omitempty"`
	AckNo   int32             `protobuf:"zigzag32,2,opt,name=ack_no,json=ackNo,proto3" json:"ack_no,omitempty"`
	ErrInfo string            `protobuf:"bytes,3,opt,name=err_info,json=errInfo,proto3" json:"err_info,omitempty"`
	Args    map[string]string `protobuf:"bytes,4,rep,name=args" json:"args,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *DownResult) Reset()                    { *m = DownResult{} }
func (m *DownResult) String() string            { return proto.CompactTextString(m) }
func (*DownResult) ProtoMessage()               {}
func (*DownResult) Descriptor() ([]byte, []int) { return fileDescriptorServerBiz, []int{5} }

func (m *DownResult) GetRet() int32 {
	if m != nil {
		return m.Ret
	}
	return 0
}

func (m *DownResult) GetAckNo() int32 {
	if m != nil {
		return m.AckNo
	}
	return 0
}

func (m *DownResult) GetErrInfo() string {
	if m != nil {
		return m.ErrInfo
	}
	return ""
}

func (m *DownResult) GetArgs() map[string]string {
	if m != nil {
		return m.Args
	}
	return nil
}

// 数据
type GeneralData struct {
	DataId   string                 `protobuf:"bytes,1,opt,name=data_id,json=dataId,proto3" json:"data_id,omitempty"`
	FrameId  uint32                 `protobuf:"varint,2,opt,name=frame_id,json=frameId,proto3" json:"frame_id,omitempty"`
	DataType GeneralData_DataType   `protobuf:"varint,3,opt,name=data_type,json=dataType,proto3,enum=serverbiz.GeneralData_DataType" json:"data_type,omitempty"`
	Status   GeneralData_DataStatus `protobuf:"varint,4,opt,name=status,proto3,enum=serverbiz.GeneralData_DataStatus" json:"status,omitempty"`
	DescArgs map[string][]byte      `protobuf:"bytes,5,rep,name=desc_args,json=descArgs" json:"desc_args,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Format   string                 `protobuf:"bytes,6,opt,name=format,proto3" json:"format,omitempty"`
	Encoding string                 `protobuf:"bytes,7,opt,name=encoding,proto3" json:"encoding,omitempty"`
	Data     []byte                 `protobuf:"bytes,8,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *GeneralData) Reset()                    { *m = GeneralData{} }
func (m *GeneralData) String() string            { return proto.CompactTextString(m) }
func (*GeneralData) ProtoMessage()               {}
func (*GeneralData) Descriptor() ([]byte, []int) { return fileDescriptorServerBiz, []int{6} }

func (m *GeneralData) GetDataId() string {
	if m != nil {
		return m.DataId
	}
	return ""
}

func (m *GeneralData) GetFrameId() uint32 {
	if m != nil {
		return m.FrameId
	}
	return 0
}

func (m *GeneralData) GetDataType() GeneralData_DataType {
	if m != nil {
		return m.DataType
	}
	return GeneralData_TEXT
}

func (m *GeneralData) GetStatus() GeneralData_DataStatus {
	if m != nil {
		return m.Status
	}
	return GeneralData_BEGIN
}

func (m *GeneralData) GetDescArgs() map[string][]byte {
	if m != nil {
		return m.DescArgs
	}
	return nil
}

func (m *GeneralData) GetFormat() string {
	if m != nil {
		return m.Format
	}
	return ""
}

func (m *GeneralData) GetEncoding() string {
	if m != nil {
		return m.Encoding
	}
	return ""
}

func (m *GeneralData) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*ServerBiz)(nil), "serverbiz.ServerBiz")
	proto.RegisterType((*GlobalRoute)(nil), "serverbiz.GlobalRoute")
	proto.RegisterType((*UpCall)(nil), "serverbiz.UpCall")
	proto.RegisterType((*UpResult)(nil), "serverbiz.UpResult")
	proto.RegisterType((*DownCall)(nil), "serverbiz.DownCall")
	proto.RegisterType((*DownResult)(nil), "serverbiz.DownResult")
	proto.RegisterType((*GeneralData)(nil), "serverbiz.GeneralData")
	proto.RegisterEnum("serverbiz.ServerBiz_MsgType", ServerBiz_MsgType_name, ServerBiz_MsgType_value)
	proto.RegisterEnum("serverbiz.GeneralData_DataType", GeneralData_DataType_name, GeneralData_DataType_value)
	proto.RegisterEnum("serverbiz.GeneralData_DataStatus", GeneralData_DataStatus_name, GeneralData_DataStatus_value)
}
func (m *ServerBiz) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ServerBiz) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.MsgType != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(m.MsgType))
	}
	if len(m.Version) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.Version)))
		i += copy(dAtA[i:], m.Version)
	}
	if m.GlobalRoute != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(m.GlobalRoute.Size()))
		n1, err := m.GlobalRoute.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.UpCall != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(m.UpCall.Size()))
		n2, err := m.UpCall.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.UpResult != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(m.UpResult.Size()))
		n3, err := m.UpResult.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if m.DownCall != nil {
		dAtA[i] = 0x32
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(m.DownCall.Size()))
		n4, err := m.DownCall.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	if m.DownResult != nil {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(m.DownResult.Size()))
		n5, err := m.DownResult.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	return i, nil
}

func (m *GlobalRoute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GlobalRoute) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.SessionId) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.SessionId)))
		i += copy(dAtA[i:], m.SessionId)
	}
	if len(m.TraceId) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.TraceId)))
		i += copy(dAtA[i:], m.TraceId)
	}
	if len(m.UpRouterId) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.UpRouterId)))
		i += copy(dAtA[i:], m.UpRouterId)
	}
	if len(m.GuiderId) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.GuiderId)))
		i += copy(dAtA[i:], m.GuiderId)
	}
	if len(m.DownRouterId) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.DownRouterId)))
		i += copy(dAtA[i:], m.DownRouterId)
	}
	if len(m.Appid) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.Appid)))
		i += copy(dAtA[i:], m.Appid)
	}
	if len(m.Uid) > 0 {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.Uid)))
		i += copy(dAtA[i:], m.Uid)
	}
	if len(m.Did) > 0 {
		dAtA[i] = 0x42
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.Did)))
		i += copy(dAtA[i:], m.Did)
	}
	if len(m.ClientIp) > 0 {
		dAtA[i] = 0x4a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.ClientIp)))
		i += copy(dAtA[i:], m.ClientIp)
	}
	return i, nil
}

func (m *UpCall) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpCall) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Call) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.Call)))
		i += copy(dAtA[i:], m.Call)
	}
	if m.SeqNo != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64((uint32(m.SeqNo)<<1)^uint32((m.SeqNo >> 31))))
	}
	if len(m.From) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.From)))
		i += copy(dAtA[i:], m.From)
	}
	if m.Sync {
		dAtA[i] = 0x20
		i++
		if m.Sync {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.BusinessArgs) > 0 {
		for k, _ := range m.BusinessArgs {
			dAtA[i] = 0x2a
			i++
			v := m.BusinessArgs[k]
			mapSize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			i = encodeVarintServerBiz(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if len(m.TempArgs) > 0 {
		for k, _ := range m.TempArgs {
			dAtA[i] = 0x32
			i++
			v := m.TempArgs[k]
			mapSize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			i = encodeVarintServerBiz(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if len(m.UserArgs) > 0 {
		for k, _ := range m.UserArgs {
			dAtA[i] = 0x3a
			i++
			v := m.UserArgs[k]
			byteSize := 0
			if len(v) > 0 {
				byteSize = 1 + len(v) + sovServerBiz(uint64(len(v)))
			}
			mapSize := 1 + len(k) + sovServerBiz(uint64(len(k))) + byteSize
			i = encodeVarintServerBiz(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			if len(v) > 0 {
				dAtA[i] = 0x12
				i++
				i = encodeVarintServerBiz(dAtA, i, uint64(len(v)))
				i += copy(dAtA[i:], v)
			}
		}
	}
	if len(m.Session) > 0 {
		for k, _ := range m.Session {
			dAtA[i] = 0x42
			i++
			v := m.Session[k]
			mapSize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			i = encodeVarintServerBiz(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if len(m.DataList) > 0 {
		for _, msg := range m.DataList {
			dAtA[i] = 0x4a
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *UpResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpResult) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Ret != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64((uint32(m.Ret)<<1)^uint32((m.Ret >> 31))))
	}
	if m.AckNo != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64((uint32(m.AckNo)<<1)^uint32((m.AckNo >> 31))))
	}
	if len(m.ErrInfo) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.ErrInfo)))
		i += copy(dAtA[i:], m.ErrInfo)
	}
	if len(m.From) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.From)))
		i += copy(dAtA[i:], m.From)
	}
	if len(m.Session) > 0 {
		for k, _ := range m.Session {
			dAtA[i] = 0x2a
			i++
			v := m.Session[k]
			mapSize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			i = encodeVarintServerBiz(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if len(m.DataList) > 0 {
		for _, msg := range m.DataList {
			dAtA[i] = 0x32
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *DownCall) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DownCall) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Ret != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64((uint32(m.Ret)<<1)^uint32((m.Ret >> 31))))
	}
	if m.SeqNo != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64((uint32(m.SeqNo)<<1)^uint32((m.SeqNo >> 31))))
	}
	if len(m.From) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.From)))
		i += copy(dAtA[i:], m.From)
	}
	if len(m.Args) > 0 {
		for k, _ := range m.Args {
			dAtA[i] = 0x22
			i++
			v := m.Args[k]
			mapSize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			i = encodeVarintServerBiz(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if len(m.DataList) > 0 {
		for _, msg := range m.DataList {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *DownResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DownResult) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Ret != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64((uint32(m.Ret)<<1)^uint32((m.Ret >> 31))))
	}
	if m.AckNo != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64((uint32(m.AckNo)<<1)^uint32((m.AckNo >> 31))))
	}
	if len(m.ErrInfo) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.ErrInfo)))
		i += copy(dAtA[i:], m.ErrInfo)
	}
	if len(m.Args) > 0 {
		for k, _ := range m.Args {
			dAtA[i] = 0x22
			i++
			v := m.Args[k]
			mapSize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			i = encodeVarintServerBiz(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	return i, nil
}

func (m *GeneralData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GeneralData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.DataId) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.DataId)))
		i += copy(dAtA[i:], m.DataId)
	}
	if m.FrameId != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(m.FrameId))
	}
	if m.DataType != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(m.DataType))
	}
	if m.Status != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(m.Status))
	}
	if len(m.DescArgs) > 0 {
		for k, _ := range m.DescArgs {
			dAtA[i] = 0x2a
			i++
			v := m.DescArgs[k]
			byteSize := 0
			if len(v) > 0 {
				byteSize = 1 + len(v) + sovServerBiz(uint64(len(v)))
			}
			mapSize := 1 + len(k) + sovServerBiz(uint64(len(k))) + byteSize
			i = encodeVarintServerBiz(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintServerBiz(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			if len(v) > 0 {
				dAtA[i] = 0x12
				i++
				i = encodeVarintServerBiz(dAtA, i, uint64(len(v)))
				i += copy(dAtA[i:], v)
			}
		}
	}
	if len(m.Format) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.Format)))
		i += copy(dAtA[i:], m.Format)
	}
	if len(m.Encoding) > 0 {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.Encoding)))
		i += copy(dAtA[i:], m.Encoding)
	}
	if len(m.Data) > 0 {
		dAtA[i] = 0x42
		i++
		i = encodeVarintServerBiz(dAtA, i, uint64(len(m.Data)))
		i += copy(dAtA[i:], m.Data)
	}
	return i, nil
}

func encodeVarintServerBiz(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ServerBiz) Size() (n int) {
	var l int
	_ = l
	if m.MsgType != 0 {
		n += 1 + sovServerBiz(uint64(m.MsgType))
	}
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if m.GlobalRoute != nil {
		l = m.GlobalRoute.Size()
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if m.UpCall != nil {
		l = m.UpCall.Size()
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if m.UpResult != nil {
		l = m.UpResult.Size()
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if m.DownCall != nil {
		l = m.DownCall.Size()
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if m.DownResult != nil {
		l = m.DownResult.Size()
		n += 1 + l + sovServerBiz(uint64(l))
	}
	return n
}

func (m *GlobalRoute) Size() (n int) {
	var l int
	_ = l
	l = len(m.SessionId)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.TraceId)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.UpRouterId)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.GuiderId)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.DownRouterId)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.Appid)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.Uid)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.Did)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.ClientIp)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	return n
}

func (m *UpCall) Size() (n int) {
	var l int
	_ = l
	l = len(m.Call)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if m.SeqNo != 0 {
		n += 1 + sozServerBiz(uint64(m.SeqNo))
	}
	l = len(m.From)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if m.Sync {
		n += 2
	}
	if len(m.BusinessArgs) > 0 {
		for k, v := range m.BusinessArgs {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			n += mapEntrySize + 1 + sovServerBiz(uint64(mapEntrySize))
		}
	}
	if len(m.TempArgs) > 0 {
		for k, v := range m.TempArgs {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			n += mapEntrySize + 1 + sovServerBiz(uint64(mapEntrySize))
		}
	}
	if len(m.UserArgs) > 0 {
		for k, v := range m.UserArgs {
			_ = k
			_ = v
			l = 0
			if len(v) > 0 {
				l = 1 + len(v) + sovServerBiz(uint64(len(v)))
			}
			mapEntrySize := 1 + len(k) + sovServerBiz(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovServerBiz(uint64(mapEntrySize))
		}
	}
	if len(m.Session) > 0 {
		for k, v := range m.Session {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			n += mapEntrySize + 1 + sovServerBiz(uint64(mapEntrySize))
		}
	}
	if len(m.DataList) > 0 {
		for _, e := range m.DataList {
			l = e.Size()
			n += 1 + l + sovServerBiz(uint64(l))
		}
	}
	return n
}

func (m *UpResult) Size() (n int) {
	var l int
	_ = l
	if m.Ret != 0 {
		n += 1 + sozServerBiz(uint64(m.Ret))
	}
	if m.AckNo != 0 {
		n += 1 + sozServerBiz(uint64(m.AckNo))
	}
	l = len(m.ErrInfo)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.From)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if len(m.Session) > 0 {
		for k, v := range m.Session {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			n += mapEntrySize + 1 + sovServerBiz(uint64(mapEntrySize))
		}
	}
	if len(m.DataList) > 0 {
		for _, e := range m.DataList {
			l = e.Size()
			n += 1 + l + sovServerBiz(uint64(l))
		}
	}
	return n
}

func (m *DownCall) Size() (n int) {
	var l int
	_ = l
	if m.Ret != 0 {
		n += 1 + sozServerBiz(uint64(m.Ret))
	}
	if m.SeqNo != 0 {
		n += 1 + sozServerBiz(uint64(m.SeqNo))
	}
	l = len(m.From)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if len(m.Args) > 0 {
		for k, v := range m.Args {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			n += mapEntrySize + 1 + sovServerBiz(uint64(mapEntrySize))
		}
	}
	if len(m.DataList) > 0 {
		for _, e := range m.DataList {
			l = e.Size()
			n += 1 + l + sovServerBiz(uint64(l))
		}
	}
	return n
}

func (m *DownResult) Size() (n int) {
	var l int
	_ = l
	if m.Ret != 0 {
		n += 1 + sozServerBiz(uint64(m.Ret))
	}
	if m.AckNo != 0 {
		n += 1 + sozServerBiz(uint64(m.AckNo))
	}
	l = len(m.ErrInfo)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if len(m.Args) > 0 {
		for k, v := range m.Args {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovServerBiz(uint64(len(k))) + 1 + len(v) + sovServerBiz(uint64(len(v)))
			n += mapEntrySize + 1 + sovServerBiz(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *GeneralData) Size() (n int) {
	var l int
	_ = l
	l = len(m.DataId)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	if m.FrameId != 0 {
		n += 1 + sovServerBiz(uint64(m.FrameId))
	}
	if m.DataType != 0 {
		n += 1 + sovServerBiz(uint64(m.DataType))
	}
	if m.Status != 0 {
		n += 1 + sovServerBiz(uint64(m.Status))
	}
	if len(m.DescArgs) > 0 {
		for k, v := range m.DescArgs {
			_ = k
			_ = v
			l = 0
			if len(v) > 0 {
				l = 1 + len(v) + sovServerBiz(uint64(len(v)))
			}
			mapEntrySize := 1 + len(k) + sovServerBiz(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovServerBiz(uint64(mapEntrySize))
		}
	}
	l = len(m.Format)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.Encoding)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovServerBiz(uint64(l))
	}
	return n
}

func sovServerBiz(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozServerBiz(x uint64) (n int) {
	return sovServerBiz(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ServerBiz) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerBiz
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ServerBiz: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ServerBiz: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MsgType", wireType)
			}
			m.MsgType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MsgType |= (ServerBiz_MsgType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GlobalRoute", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.GlobalRoute == nil {
				m.GlobalRoute = &GlobalRoute{}
			}
			if err := m.GlobalRoute.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpCall", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.UpCall == nil {
				m.UpCall = &UpCall{}
			}
			if err := m.UpCall.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpResult", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.UpResult == nil {
				m.UpResult = &UpResult{}
			}
			if err := m.UpResult.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DownCall", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DownCall == nil {
				m.DownCall = &DownCall{}
			}
			if err := m.DownCall.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DownResult", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DownResult == nil {
				m.DownResult = &DownResult{}
			}
			if err := m.DownResult.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GlobalRoute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerBiz
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GlobalRoute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GlobalRoute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SessionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SessionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TraceId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TraceId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpRouterId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UpRouterId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GuiderId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GuiderId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DownRouterId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DownRouterId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Appid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Appid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Did", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Did = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientIp", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientIp = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UpCall) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerBiz
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UpCall: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpCall: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Call", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Call = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SeqNo", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			v = int32((uint32(v) >> 1) ^ uint32(((v&1)<<31)>>31))
			m.SeqNo = v
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.From = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sync", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Sync = bool(v != 0)
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BusinessArgs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BusinessArgs == nil {
				m.BusinessArgs = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowServerBiz
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipServerBiz(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthServerBiz
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.BusinessArgs[mapkey] = mapvalue
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TempArgs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TempArgs == nil {
				m.TempArgs = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowServerBiz
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipServerBiz(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthServerBiz
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.TempArgs[mapkey] = mapvalue
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserArgs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.UserArgs == nil {
				m.UserArgs = make(map[string][]byte)
			}
			var mapkey string
			mapvalue := []byte{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowServerBiz
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapbyteLen uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapbyteLen |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intMapbyteLen := int(mapbyteLen)
					if intMapbyteLen < 0 {
						return ErrInvalidLengthServerBiz
					}
					postbytesIndex := iNdEx + intMapbyteLen
					if postbytesIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = make([]byte, mapbyteLen)
					copy(mapvalue, dAtA[iNdEx:postbytesIndex])
					iNdEx = postbytesIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipServerBiz(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthServerBiz
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.UserArgs[mapkey] = mapvalue
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Session", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Session == nil {
				m.Session = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowServerBiz
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipServerBiz(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthServerBiz
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Session[mapkey] = mapvalue
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataList = append(m.DataList, &GeneralData{})
			if err := m.DataList[len(m.DataList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UpResult) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerBiz
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UpResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ret", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			v = int32((uint32(v) >> 1) ^ uint32(((v&1)<<31)>>31))
			m.Ret = v
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AckNo", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			v = int32((uint32(v) >> 1) ^ uint32(((v&1)<<31)>>31))
			m.AckNo = v
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrInfo", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ErrInfo = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.From = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Session", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Session == nil {
				m.Session = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowServerBiz
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipServerBiz(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthServerBiz
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Session[mapkey] = mapvalue
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataList = append(m.DataList, &GeneralData{})
			if err := m.DataList[len(m.DataList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DownCall) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerBiz
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DownCall: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DownCall: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ret", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			v = int32((uint32(v) >> 1) ^ uint32(((v&1)<<31)>>31))
			m.Ret = v
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SeqNo", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			v = int32((uint32(v) >> 1) ^ uint32(((v&1)<<31)>>31))
			m.SeqNo = v
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.From = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Args", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Args == nil {
				m.Args = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowServerBiz
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipServerBiz(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthServerBiz
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Args[mapkey] = mapvalue
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataList = append(m.DataList, &GeneralData{})
			if err := m.DataList[len(m.DataList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DownResult) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerBiz
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DownResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DownResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ret", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			v = int32((uint32(v) >> 1) ^ uint32(((v&1)<<31)>>31))
			m.Ret = v
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AckNo", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			v = int32((uint32(v) >> 1) ^ uint32(((v&1)<<31)>>31))
			m.AckNo = v
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrInfo", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ErrInfo = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Args", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Args == nil {
				m.Args = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowServerBiz
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipServerBiz(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthServerBiz
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Args[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GeneralData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerBiz
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GeneralData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GeneralData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FrameId", wireType)
			}
			m.FrameId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FrameId |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataType", wireType)
			}
			m.DataType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DataType |= (GeneralData_DataType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= (GeneralData_DataStatus(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DescArgs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DescArgs == nil {
				m.DescArgs = make(map[string][]byte)
			}
			var mapkey string
			mapvalue := []byte{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowServerBiz
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthServerBiz
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapbyteLen uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowServerBiz
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapbyteLen |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intMapbyteLen := int(mapbyteLen)
					if intMapbyteLen < 0 {
						return ErrInvalidLengthServerBiz
					}
					postbytesIndex := iNdEx + intMapbyteLen
					if postbytesIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = make([]byte, mapbyteLen)
					copy(mapvalue, dAtA[iNdEx:postbytesIndex])
					iNdEx = postbytesIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipServerBiz(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthServerBiz
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.DescArgs[mapkey] = mapvalue
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Format", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Format = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Encoding", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Encoding = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthServerBiz
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipServerBiz(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowServerBiz
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowServerBiz
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthServerBiz
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowServerBiz
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipServerBiz(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthServerBiz = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowServerBiz   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("ServerBiz.proto", fileDescriptorServerBiz) }

var fileDescriptorServerBiz = []byte{
	// 1034 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0x4f, 0x73, 0xdb, 0x54,
	0x10, 0x8f, 0x2c, 0x5b, 0x7f, 0xd6, 0x4e, 0xaa, 0x3c, 0x68, 0x11, 0x81, 0x26, 0xc6, 0xf4, 0x90,
	0xe1, 0xe0, 0x81, 0x64, 0xa0, 0x69, 0x60, 0x86, 0x49, 0x62, 0x4d, 0xd0, 0x4c, 0xea, 0x74, 0x14,
	0x1b, 0xb8, 0x69, 0x14, 0xeb, 0xc5, 0xa3, 0x89, 0x2c, 0xa9, 0xef, 0x49, 0xe9, 0xa4, 0x57, 0x0e,
	0x1c, 0xb9, 0xf2, 0x11, 0xf8, 0x14, 0x9c, 0x39, 0xf2, 0x11, 0x98, 0x70, 0xe1, 0xce, 0x17, 0x60,
	0xde, 0xea, 0xd9, 0x96, 0x9b, 0x14, 0x92, 0xd2, 0x4b, 0xb2, 0xfb, 0x76, 0x7f, 0xfb, 0xbc, 0xbf,
	0x5d, 0xed, 0x3e, 0xb8, 0x77, 0x42, 0xd9, 0x05, 0x65, 0xfb, 0xd1, 0xcb, 0x6e, 0xc6, 0xd2, 0x3c,
	0x25, 0x26, 0xc7, 0x83, 0xd3, 0xe8, 0x65, 0xe7, 0x17, 0x15, 0xcc, 0x99, 0x99, 0x3c, 0x06, 0x63,
	0xc2, 0xc7, 0x7e, 0x7e, 0x99, 0x51, 0x5b, 0x69, 0x2b, 0x9b, 0x2b, 0x5b, 0x1f, 0x76, 0x67, 0xbe,
	0xdd, 0x79, 0x98, 0xa7, 0x7c, 0x3c, 0xb8, 0xcc, 0xa8, 0xa7, 0x4f, 0x4a, 0x81, 0xd8, 0xa0, 0x5f,
	0x50, 0xc6, 0xa3, 0x34, 0xb1, 0x6b, 0x6d, 0x65, 0xd3, 0xf4, 0xa6, 0x2a, 0x79, 0x02, 0xad, 0x71,
	0x9c, 0x9e, 0x06, 0xb1, 0xcf, 0xd2, 0x22, 0xa7, 0xb6, 0xda, 0x56, 0x36, 0x9b, 0x5b, 0x0f, 0x2a,
	0x61, 0x0f, 0xd1, 0xec, 0x09, 0xab, 0xd7, 0x1c, 0xcf, 0x15, 0xf2, 0x09, 0xe8, 0x45, 0xe6, 0x8f,
	0x82, 0x38, 0xb6, 0xeb, 0x88, 0x5a, 0xad, 0xa0, 0x86, 0xd9, 0x41, 0x10, 0xc7, 0x9e, 0x56, 0xe0,
	0x7f, 0xf2, 0x29, 0x98, 0x45, 0xe6, 0x33, 0xca, 0x8b, 0x38, 0xb7, 0x1b, 0xe8, 0xfd, 0xce, 0x82,
	0xb7, 0x87, 0x26, 0xcf, 0x28, 0xa4, 0x24, 0x10, 0x61, 0xfa, 0x22, 0x29, 0xe3, 0x6b, 0xd7, 0x10,
	0xbd, 0xf4, 0x45, 0x82, 0x37, 0x18, 0xa1, 0x94, 0xc8, 0x17, 0xd0, 0x44, 0x84, 0xbc, 0x45, 0x47,
	0xcc, 0xfd, 0x57, 0x30, 0xf2, 0x1e, 0x08, 0x67, 0x72, 0xc7, 0x01, 0x5d, 0x12, 0x46, 0x9a, 0xa0,
	0x0f, 0x9f, 0xf9, 0x07, 0x7b, 0x47, 0x47, 0xd6, 0x12, 0x59, 0x06, 0x73, 0xf8, 0xcc, 0xf7, 0x9c,
	0x93, 0xe1, 0xd1, 0xc0, 0x52, 0x84, 0xda, 0x3b, 0xfe, 0xae, 0x5f, 0x5a, 0x6b, 0xe4, 0x1e, 0x34,
	0x51, 0x95, 0x76, 0xb5, 0xf3, 0x43, 0x0d, 0x9a, 0x15, 0xae, 0xc8, 0x43, 0x00, 0x4e, 0xb9, 0x20,
	0xd9, 0x8f, 0x42, 0x2c, 0x97, 0xe9, 0x99, 0xf2, 0xc4, 0x0d, 0xc9, 0xfb, 0x60, 0xe4, 0x2c, 0x18,
	0x51, 0x61, 0x94, 0x35, 0x41, 0xdd, 0x0d, 0x49, 0x1b, 0x5a, 0x82, 0x2c, 0x11, 0x85, 0x09, 0xb3,
	0x8a, 0x66, 0x28, 0x32, 0x0c, 0xcc, 0xdc, 0x90, 0x7c, 0x00, 0xe6, 0xb8, 0x88, 0xc2, 0xd2, 0x5c,
	0x47, 0xb3, 0x51, 0x1e, 0xb8, 0x21, 0x79, 0x04, 0x2b, 0x25, 0x0f, 0xb3, 0x00, 0x0d, 0xf4, 0x68,
	0x61, 0xce, 0xd3, 0x10, 0xef, 0x42, 0x23, 0xc8, 0xb2, 0x28, 0x44, 0x6e, 0x4d, 0xaf, 0x54, 0x88,
	0x05, 0x6a, 0x11, 0x85, 0xc8, 0x9d, 0xe9, 0x09, 0x51, 0x9c, 0x84, 0x51, 0x68, 0x1b, 0xe5, 0x49,
	0x18, 0xe1, 0xe5, 0xa3, 0x38, 0xa2, 0x49, 0xee, 0x47, 0x99, 0x6d, 0x96, 0x97, 0x97, 0x07, 0x6e,
	0xd6, 0xf9, 0xb1, 0x01, 0x5a, 0x59, 0x7b, 0x42, 0xa0, 0x8e, 0xc5, 0x2b, 0x53, 0x47, 0x99, 0xdc,
	0x07, 0x8d, 0xd3, 0xe7, 0x7e, 0x92, 0x62, 0xce, 0xab, 0x5e, 0x83, 0xd3, 0xe7, 0xfd, 0x54, 0xb8,
	0x9e, 0xb1, 0x74, 0x22, 0x33, 0x45, 0x59, 0x9c, 0xf1, 0xcb, 0x64, 0x84, 0xe9, 0x19, 0x1e, 0xca,
	0xe4, 0x1b, 0x58, 0x3e, 0x2d, 0x78, 0x94, 0x50, 0xce, 0xfd, 0x80, 0x8d, 0xb9, 0xdd, 0x68, 0xab,
	0x9b, 0xcd, 0xad, 0x8f, 0xaf, 0x35, 0x5e, 0x77, 0x5f, 0xba, 0xed, 0xb1, 0x31, 0x77, 0x92, 0x9c,
	0x5d, 0x7a, 0xad, 0xd3, 0xca, 0x11, 0xf9, 0x0a, 0xcc, 0x9c, 0x4e, 0xb2, 0x32, 0x8a, 0x86, 0x51,
	0x36, 0xae, 0x47, 0x19, 0xd0, 0x49, 0x36, 0x8f, 0x60, 0xe4, 0x52, 0x15, 0xe8, 0x82, 0x53, 0x56,
	0xa2, 0xf5, 0xd7, 0xa1, 0x87, 0x9c, 0xb2, 0x0a, 0xba, 0x90, 0x2a, 0xd9, 0x01, 0x5d, 0xf6, 0x81,
	0x6d, 0x20, 0x76, 0xfd, 0x3a, 0xf6, 0xa4, 0x74, 0x28, 0xa1, 0x53, 0x77, 0xb2, 0x0d, 0x66, 0x18,
	0xe4, 0x81, 0x1f, 0x47, 0x3c, 0xb7, 0x4d, 0xc4, 0x2e, 0x7c, 0xaa, 0x34, 0xa1, 0x2c, 0x88, 0x7b,
	0x41, 0x1e, 0x78, 0x86, 0x70, 0x3c, 0x8a, 0x78, 0xbe, 0xf6, 0x35, 0xac, 0x5e, 0x63, 0x43, 0x94,
	0xf5, 0x9c, 0x5e, 0xca, 0xda, 0x08, 0x51, 0x34, 0xc4, 0x45, 0x10, 0x17, 0x54, 0x76, 0x63, 0xa9,
	0xec, 0xd6, 0x76, 0x94, 0xb5, 0x2f, 0x61, 0x79, 0x81, 0x88, 0xbb, 0x82, 0x17, 0x78, 0xf8, 0x2f,
	0x70, 0xab, 0x0a, 0xde, 0x85, 0x56, 0x95, 0x88, 0xbb, 0x5c, 0xdc, 0xf9, 0xa9, 0x06, 0xc6, 0x74,
	0xae, 0x08, 0x20, 0xa3, 0x39, 0x02, 0x57, 0x3d, 0x21, 0x8a, 0x4e, 0x0c, 0x46, 0xe7, 0x95, 0x4e,
	0x0c, 0x46, 0xe7, 0xfd, 0x54, 0x7c, 0x96, 0x94, 0x31, 0x3f, 0x4a, 0xce, 0x52, 0xd9, 0x8d, 0x3a,
	0x65, 0xcc, 0x4d, 0xce, 0xe6, 0x4d, 0x5a, 0xaf, 0x34, 0xe9, 0xee, 0xbc, 0x94, 0x65, 0x2b, 0xb6,
	0x6f, 0x98, 0x6a, 0xb7, 0x29, 0xa6, 0x76, 0xcb, 0x62, 0xfe, 0x1f, 0x46, 0xfe, 0x52, 0xc0, 0x98,
	0xce, 0xcd, 0x9b, 0x19, 0xb9, 0xed, 0xb7, 0xf9, 0x19, 0xd4, 0xb1, 0xf5, 0xeb, 0xf8, 0xab, 0x1f,
	0xde, 0x30, 0x97, 0xbb, 0xf3, 0xc6, 0x47, 0xd7, 0xc5, 0x6c, 0x1b, 0xb7, 0xcc, 0xf6, 0x31, 0x98,
	0x6f, 0xd4, 0x75, 0x9d, 0x5f, 0x15, 0x80, 0xf9, 0xb8, 0x7f, 0x2b, 0xe5, 0xdf, 0x5e, 0xc8, 0x79,
	0xe3, 0xc6, 0xbd, 0xf2, 0x6a, 0xd6, 0x6f, 0x9e, 0xc0, 0xdf, 0x2a, 0x34, 0x2b, 0x9c, 0x90, 0xf7,
	0x40, 0x47, 0xfa, 0x66, 0xab, 0x44, 0x13, 0x6a, 0xb9, 0x47, 0xce, 0x58, 0x30, 0x99, 0xed, 0x91,
	0x65, 0x4f, 0x47, 0xdd, 0x0d, 0xc5, 0x94, 0x42, 0x0c, 0xbe, 0x17, 0x54, 0x7c, 0x2f, 0x6c, 0xdc,
	0x4c, 0x79, 0x57, 0xfc, 0xc1, 0x27, 0x03, 0x72, 0x8f, 0xbb, 0xf0, 0x09, 0x68, 0x3c, 0x0f, 0xf2,
	0x82, 0x63, 0xc3, 0xaf, 0x6c, 0x7d, 0xf4, 0x2f, 0xd0, 0x13, 0x74, 0xf4, 0x24, 0x80, 0xec, 0x81,
	0x19, 0x52, 0x3e, 0xaa, 0x8e, 0xe8, 0x47, 0xaf, 0x43, 0x53, 0x3e, 0xaa, 0xcc, 0xc8, 0x50, 0xaa,
	0xe4, 0x01, 0x68, 0x67, 0x29, 0x9b, 0x04, 0xb9, 0xdc, 0x4f, 0x52, 0x23, 0x6b, 0x60, 0xd0, 0x64,
	0x94, 0x86, 0x51, 0x32, 0x96, 0x5b, 0x6a, 0xa6, 0x8b, 0x4e, 0x15, 0xbf, 0x1e, 0x77, 0x55, 0xcb,
	0x43, 0x59, 0x8c, 0x9f, 0x85, 0x2b, 0xee, 0x32, 0x7e, 0x3a, 0x9f, 0x83, 0x31, 0x25, 0x86, 0x18,
	0x50, 0x1f, 0x38, 0xdf, 0x0f, 0xac, 0x25, 0x62, 0x42, 0x63, 0x6f, 0xd8, 0x73, 0x8f, 0x2d, 0x45,
	0x88, 0xee, 0xd3, 0xbd, 0x43, 0xc7, 0xaa, 0x09, 0xf1, 0x5b, 0xb7, 0xe7, 0x1c, 0x5b, 0x6a, 0x67,
	0x07, 0x60, 0x4e, 0x8a, 0x30, 0xec, 0x3b, 0x87, 0x6e, 0xdf, 0x5a, 0x22, 0x2d, 0x30, 0x0e, 0x8e,
	0xfb, 0x03, 0xb7, 0x3f, 0x74, 0x2c, 0x85, 0xe8, 0xa0, 0x3a, 0xfd, 0x9e, 0x55, 0x13, 0xa1, 0x8f,
	0xfb, 0x07, 0x8e, 0xa5, 0xee, 0x5b, 0xbf, 0x5d, 0xad, 0x2b, 0xbf, 0x5f, 0xad, 0x2b, 0x7f, 0x5c,
	0xad, 0x2b, 0x3f, 0xff, 0xb9, 0xbe, 0x74, 0xaa, 0xe1, 0x93, 0x70, 0xfb, 0x9f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x60, 0xf7, 0xec, 0x95, 0x25, 0x0a, 0x00, 0x00,
}
