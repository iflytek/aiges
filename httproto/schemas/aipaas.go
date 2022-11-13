package schemas

import (
	"encoding/json"
	"fmt"
	"github.com/seeadoog/jsonschema"
	"github.com/xfyun/aiges/httproto/common"
	"github.com/xfyun/aiges/httproto/pb"
	"reflect"
	//"git.xfyun.cn/AIaaS/webgate-ws/common"
)

type Accept struct {
	DataType string `json:"dataType"`
}

type AIService struct {
	Input  map[string]interface{} `json:"input"`
	Accept map[string]Accept      `json:"accept"`
}

func (s *AIService) GetInput() map[string]interface{} {
	return s.Input
}

func (s *AIService) GetAccept() map[string]Accept {
	return s.Accept
}

const (
	From         = "webgate-ws"
	Sync         = false
	ProtoVersion = "2.0"
	//CallBackKey  = "call_back_addr"
)

var (
	CallBackAddr = ""
)

type AIpaas struct {
	Meta         map[string]interface{} `json:"meta"`
	SchemaInput  *jsonschema.Schema     `json:"schema_input"`
	SchemaOutput *JsonElement           `json:"schema_output"`
}

type Context struct {
	SeqNo    int32
	Session  map[string]string
	Header   map[string]string
	Sync     bool
	IsStream bool
}

func (c *Context) getSeqNo() int32 {
	if c == nil {
		return 0
	}
	return c.SeqNo
}

func (c *Context) GetSession() map[string]string {
	if c == nil {
		return map[string]string{}
	}
	return c.Session
}

func (c *Context) GetHeader() map[string]string {
	if c == nil {
		return map[string]string{}
	}
	return c.Header
}

// serviceId
// version
// hosts
// service
// route
// sub
// call
// call_type

type Meta map[string]interface{}

func (m Meta) GetServiceId() string {
	return String(m["serviceId"])
}

func (m Meta) GetCloudId() string {
	return String(m["cloudId"])
}

func (m Meta) GetVersion() string {
	return String(m["version"])
}

func (m Meta) GetHost() []string {
	h := m["hosts"]
	switch h.(type) {
	case string:
		host := h.(string)
		if host == "" {
			return []string{}
		}
		return []string{host}
	case []interface{}:
		hs := h.([]interface{})
		hss := make([]string, len(hs))
		for _, v := range hs {
			hss = append(hss, String(v))
		}
		return hss
	}
	return []string{}
}

func (m Meta) GetService() []string {
	if seri, ok := m["service"].([]interface{}); ok {
		services := make([]string, len(seri))
		for idx, si := range seri {
			services[idx] = String(si)
		}
		return services
	}
	return []string{}
}

func (m Meta) GetRoute() []string {
	route := m["route"]
	switch route.(type) {
	case string:
		return []string{route.(string)}
	case []interface{}:
		rs := []string{}
		for _, route := range route.([]interface{}) {
			rs = append(rs, String(route))
		}
		return rs
	}
	return nil
}

func (m Meta) GetSesssionStat(ctx *Context) pb.UpCall_SessState {
	if ctx.IsStream {
		return pb.UpCall_STREAM
	}
	return pb.UpCall_ONCE
}

func (m Meta) GetSub() string {
	return String(m["sub"])
}

func parseAccept(in interface{}) Accept {
	a := Accept{}
	if m, ok := in.(map[string]interface{}); ok {
		a.DataType = String(m["dataType"])
	}
	return a
}

// 获取service
func (m Meta) GetServiceByNickName(nickName string) *AIService {
	if sm, ok := m[nickName].(map[string]interface{}); ok {
		ais := &AIService{}
		if input, ok := sm["input"].(map[string]interface{}); ok {
			ais.Input = input
		}
		if accept, ok := sm["accept"].(map[string]interface{}); ok {
			acpt := make(map[string]Accept)
			for k, v := range accept {
				acpt[k] = parseAccept(v)
			}
			ais.Accept = acpt
		} else {
			//return nil
			return ais
		}

		return ais
	}
	return &AIService{}
}

func (m Meta) Sync() bool {
	return false
}

func (m Meta) GetCallService() string {
	return String(m["call"])
}

func (m Meta) GetCallType() int {
	if t, ok := m["call_type"]; ok {
		return int(Number(t))
	}
	return 0
}

func (m Meta) GetSessonTimeout() int {
	return int(Number(m["session_timeout"]))
}

func (m Meta) GetReadTimeout() int {
	return int(Number(m["read_timeout"]))
}

// convert map[string]interface{} ot map[string]string
func OMapToStringMap(in map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for k, v := range in {
		out[k] = String(v)
	}
	return out
}

//resolve global route
func (m Meta) resolveGlobalRoute(headerMap map[string]interface{}, ctx *Context) *pb.GlobalRoute {
	header := OMapToStringMap(headerMap)
	for k, v := range ctx.GetHeader() {
		if _, ok := header[k]; ok {
			continue
		}
		header[k] = v
	}

	gr := &pb.GlobalRoute{
		Headers: header,
	}
	return gr
}

func getDataKeyByDataType(dataType string) string {
	switch dataType {
	case "audio":
		return "audio"
	case "image":
		return "image"
	case "text":
		return "text"
	case "video":
		return "video"
	case "other":
		return "other"

	}
	return "audio"
}

func getDataPbType(dataType string) pb.MetaDesc_DataType {
	switch dataType {
	case "audio":
		return pb.MetaDesc_AUDIO
	case "image":
		return pb.MetaDesc_IMAGE
	case "video":
		return pb.MetaDesc_VIDEO
	case "text":
		return pb.MetaDesc_TEXT
	case "other":
		return pb.MetaDesc_OTHER

	}
	return pb.MetaDesc_AUDIO
}

func getAttribute(inputData map[string]interface{}, dataKey string) map[string]string {
	attr := make(map[string]string)
	for k, v := range inputData {
		if k == dataKey {
			continue
		}
		attr[k] = String(v)
	}
	return attr
}

func (m Meta) resolveDataList(payload map[string]interface{}, services []string) ([]*pb.GeneralData, error) {
	dataList := make([]*pb.GeneralData, 0, 1)
	if payload == nil {
		return nil, nil
	}
	for _, srvName := range services {

		service := m.GetServiceByNickName(srvName)
		if service == nil {
			continue
		}
		for input, dataType := range service.GetInput() {
			dataKey := ""
			switch dataType.(type) {
			case string:
				dataKey = dataType.(string)
			case map[string]interface{}:
				dtm := dataType.(map[string]interface{})
				dataKey = String(dtm["dataType"])
			default:
				return nil, fmt.Errorf("invalid data type:%s,val:%v", reflect.TypeOf(dataType).String(), dataType)
			}

			dataKey = getDataKeyByDataType(dataKey)
			inputData, ok := payload[input].(map[string]interface{}) // 当前service 的 数据流
			if !ok {
				//return nil, fmt.Errorf("resolve payload error: input stream %s is not object", input)
				continue
			}
			meta := &pb.MetaDesc{
				Name:      input,
				DataType:  getDataPbType(dataKey),
				Attribute: getAttribute(inputData, dataKey),
			}

			data, err := common.DecodeBase64string(String(inputData[dataKey]))
			if err != nil {
				return nil, fmt.Errorf("base64 decode error,input stream,field : '%s' must be encode to base64 string. error:%s", input, err)
			}
			gd := &pb.GeneralData{
				Meta: meta,
				Data: data,
			}
			dataList = append(dataList, gd)
		}
	}
	return dataList, nil
}

func (m Meta) resolveBusinessArgsAndPle(parameter map[string]interface{}, services []string) (map[string]*pb.ArgsData, []*pb.MetaDesc, error) {
	if parameter == nil {
		return nil, nil, nil
	}
	businessArgs := make(map[string]*pb.ArgsData)

	for _, srvName := range services {
		acceptPbs := make(map[string]*pb.MetaDesc)
		srvMap, ok := parameter[srvName].(map[string]interface{})
		if !ok {
			//return nil, nil, fmt.Errorf("%s parameter name is not object", srvName)
			continue
		}
		service := m.GetServiceByNickName(srvName)
		if service == nil {
			continue
		}
		busi := make(map[string]string)
		for accept, dataType := range service.Accept {

			acceptMap, ok := srvMap[accept].(map[string]interface{})
			if !ok {
				//return nil, nil, fmt.Errorf("resolve business error:%s.%s in not object", srvName, accept)
				continue
			}
			acceptPb := &pb.MetaDesc{
				//todo add service_name
				Name:      String(accept),
				DataType:  getDataPbType(dataType.DataType),
				Attribute: OMapToStringMap(acceptMap),
			}
			acceptPbs[String(accept)] = acceptPb
		}
		for key, val := range srvMap {
			if service.Accept != nil {
				if _, ok := service.Accept[key]; ok {
					continue

				}
			}
			busi[key] = String(val)
		}
		busiPb := &pb.ArgsData{
			BusinessArgs: busi,
			Ple:          acceptPbs,
		}
		businessArgs[srvName] = busiPb
	}

	return businessArgs, nil, nil
}

//
func (m Meta) resolveUpCall(root map[string]interface{}, ctx *Context) (*pb.UpCall, error) {

	payload, ok := root["payload"].(map[string]interface{})
	if !ok {
		//return nil, fmt.Errorf("payload is not object")
		payload = nil
	}

	parameter, ok := root["parameter"].(map[string]interface{})
	if !ok {
		//return nil, fmt.Errorf("parameter is not object")
		parameter = nil
	}

	dataList, err := m.resolveDataList(payload, m.GetService())
	if err != nil {
		return nil, err
	}

	businessArgs, _, err := m.resolveBusinessArgsAndPle(parameter, m.GetService())
	if err != nil {
		return nil, err
	}

	up := &pb.UpCall{
		Call:         m.GetSub(),
		SeqNo:        ctx.getSeqNo(),
		From:         From,
		Sync:         ctx.Sync,
		BusinessArgs: businessArgs,
		Session:      ctx.GetSession(),
		SessionState: m.GetSesssionStat(ctx),
		//Ple:          ple,
		DataList: dataList,
	}
	return up, nil
}

//var pbPool = sync.Pool{}

func (m Meta) ResolveServerBiz(root map[string]interface{}, ctx *Context) (*pb.ServerBiz, error) {
	headerMap, ok := root["header"].(map[string]interface{})
	if !ok {
		// Header 不一定每一帧都会有
		headerMap = map[string]interface{}{}
	}
	up, err := m.resolveUpCall(root, ctx)
	if err != nil {
		return nil, err
	}

	biz := &pb.ServerBiz{
		MsgType:     pb.ServerBiz_UP_CALL,
		Version:     ProtoVersion,
		GlobalRoute: m.resolveGlobalRoute(headerMap, ctx),
		UpCall:      up,
		UpResult:    nil,
		DownCall:    nil,
		DownResult:  nil,
	}
	return biz, nil
}

//
var dataTypeEnums = [...]string{
	pb.MetaDesc_TEXT:  "text",
	pb.MetaDesc_AUDIO: "audio",
	pb.MetaDesc_IMAGE: "image",
	pb.MetaDesc_VIDEO: "video",
	pb.MetaDesc_OTHER: "other",
}

func getRespDataKeyByDataType(t pb.MetaDesc_DataType) string {
	//switch t {
	//case pb.MetaDesc_TEXT:
	//	return "text"
	//case pb.MetaDesc_AUDIO:
	//	return "audio"
	//case pb.MetaDesc_IMAGE:
	//	return "image"
	//case pb.MetaDesc_VIDEO:
	//	return "video"
	//case pb.MetaDesc_OTHER:
	//	return "other"
	//}
	if int(t) >= len(dataTypeEnums) || t < 0 {
		return "other"
	}
	return dataTypeEnums[t]
}

func format(v interface{}, typ string) interface{} {
	switch typ {
	case "string":
		return String(v)
	case "integer", "number":
		return Number(v)
	case "bool":
		return Bool(v)

	}
	return v
}

func (m Meta) resolveDataListResp(datalist []*pb.GeneralData, outPutSchema *JsonElement) map[string]interface{} {
	if len(datalist) == 0 {
		return nil
	}
	payload := make(map[string]interface{})
	for _, generalData := range datalist {
		name := generalData.GetMeta().GetName()
		dataType := generalData.GetMeta().GetDataType()
		//payload[name] =
		data := make(map[string]interface{})
		data[getRespDataKeyByDataType(dataType)] = m.formatData(generalData.GetMeta().GetAttribute(), generalData.GetData())
		attrs := outPutSchema.Get("properties").Get("payload").Get("properties").Get(name).Get("properties")
		for key, val := range generalData.GetMeta().GetAttribute() {
			typ, ok := attrs.Get(key).Get("type").GetAsString()
			if ok {
				data[key] = format(val, typ)
			} else {
				data[key] = val
			}
		}
		payload[name] = data
	}
	return payload
}

func (m Meta) ResolveDownResponseFromServerBiz(biz *pb.ServerBiz, e *JsonElement) map[string]interface{} {
	return m.resolveDataListResp(biz.GetDownCall().GetDataList(), e)
}

func (m Meta) ResolveUpResult(biz *pb.UpResult, e *JsonElement) map[string]interface{} {
	return m.resolveDataListResp(biz.GetDataList(), e)
}

type AISchema struct {
	Meta         Meta                   `json:"meta"`
	InputSchema  *jsonschema.Schema     `json:"schemainput"`
	SchemaOutput *JsonElement           `json:"schemaoutput"`
	subRouteMap  *routeMap              `json:"sub_route_map"`
	outPutFormat map[string]string      `json:"out_put_format"`
	headerSchema map[string]interface{} `json:"header_schema"`
}

func (s *AISchema) GetSubServiceId(req interface{}) (subServiceid string, routerInfo string) {
	return s.subRouteMap.getSubServiceId(req)
}

func (s *AISchema) GetSource() string {
	return String(s.Meta["serviceSource"])
}

func (s *AISchema) Validate(in interface{}) error {
	if s.InputSchema == nil {
		return nil
	}
	return s.InputSchema.Validate(in)
}

func (s *AISchema) ResolveServerBiz(root map[string]interface{}, ctx *Context) (*pb.ServerBiz, error) {
	return s.Meta.ResolveServerBiz(root, ctx)
}

func (s *AISchema) ResolveDownResponseByBiz(biz *pb.ServerBiz) map[string]interface{} {
	return s.Meta.ResolveDownResponseFromServerBiz(biz, s.SchemaOutput)
}

func (s *AISchema) ResolveUpResult(biz *pb.UpResult) interface{} {
	return s.Meta.ResolveUpResult(biz, s.SchemaOutput)
}

func (m Meta) GetCompanion() string {
	return String(m["companion_route"])
}
func (m Meta) GetCodeMap() map[string]interface{} {
	mp, ok := m["codeMap"].(map[string]interface{})
	if ok {
		return mp
	}
	return nil
}
func (m Meta) IsCategory() bool {
	return m["type"] == "category" || String(m["companion_route"]) != ""
}

func (m Meta) BuildHeader() bool {
	return Bool(m["build_header"])
}
func (m Meta) EnableClientSession() bool {
	return Bool(m["enable_client_session"])
}
func (s *AISchema) Init() error {
	//if len(s.Meta.GetRoute()) == 0 {
	//	return fmt.Errorf("init schema error,route is empty,serviceId=%s", s.Meta.GetServiceId())
	//}

	if s.Meta.GetServiceId() == "" {
		return fmt.Errorf("init schema error, serviceId is empty")
	}
	services := s.Meta.GetService()

	for _, host := range services {
		if host == "" {
			return fmt.Errorf("host item cannot be empty")
		}
	}
	s.outPutFormat = map[string]string{}
	return nil
}

func (m Meta) formatData(desc map[string]string, data []byte) interface{} {
	if desc == nil {
		return common.EncodingTobase64String(data)
	}

	format := desc["format"]
	if format == "json" {
		if json.Valid(data) {
			return JsonRawString(data)
		}
	}
	return common.EncodingTobase64String(data)
}

func (s *AISchema) BuildResponseHeader(headers map[string]string) map[string]interface{} {
	if headers == nil || !s.Meta.BuildHeader() {
		return nil
	}
	rhd := map[string]interface{}{}
	hsc := s.headerSchema
	for key, val := range hsc {
		typ, ok := NewJsonElem(val).Get("type").GetAsString()
		if ok {
			vv, has := headers[key]
			if has {
				rhd[key] = format(vv, typ)
			}
		}
	}
	return rhd
}

type Request struct {
}

type JsonRawString []byte

func (j JsonRawString) MarshalJSON() ([]byte, error) {
	return []byte(j), nil
}
