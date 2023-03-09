package dto

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/xfyun/aiges/httproto/schemas"
	"github.com/xfyun/aiges/protocol"
	"reflect"
	"strconv"
)

type Request struct {
	Header    map[string]interface{}            `json:"header"`
	Parameter map[string]map[string]interface{} `json:"parameter"`
	Payload   map[string]map[string]interface{} `json:"payload"`
}

func (r *Request) ConvertToPb(serviceName string, stat protocol.LoaderInput_SessState, ctx *context.Context, input_syncId int32) (*protocol.LoaderInput, error) {
	sch := schemas.GetSvcSchemaFromPython()
	sch.Meta.GetService()
	in := &protocol.LoaderInput{
		ServiceId:   serviceName,
		ServiceName: serviceName,
		State:       stat,
		Headers:     toStringMap(r.Header),
		Params:      map[string]string{},
		Expect:      nil,
		Pl:          nil,
		SyncId:      input_syncId,
	}
	r.readParameter(in)
	err := r.readPayload(in)
	//st, _ := sch.InputSchema.MarshalJSON()
	// todo
	//err = sch.Validate(r)

	return in, err
}

func (r *Request) readParameter(pb *protocol.LoaderInput) {
	// ** 理论上这里应该要根据schema 做返回数据结果的限制 ,前面要加schema校验 数据合法性

	// todo
	for _, val := range r.Parameter {
		for sk, sv := range val {
			switch sv.(type) {
			case map[string]interface{}:
				var tmp = sv.(map[string]interface{})
				var attr = map[string]string{}
				encoding, ok := tmp["encoding"]
				if ok {
					attr["encoding"] = encoding.(string)
				}
				format, ok := tmp["format"]
				if ok {
					attr["format"] = format.(string)
				}
				compress, ok := tmp["compress"]
				if ok {
					attr["compress"] = compress.(string)
				}
				desc := &protocol.MetaDesc{
					Name:      sk,
					DataType:  getDataType(toString(tmp["data_type"])),
					Attribute: attr,
				}
				pb.Expect = append(pb.Expect, desc)
			case string, int, bool, float64:
				pb.Params[sk] = toString(sv)
			default:

			}
		}
	}

}

func (r *Request) readPayload(pb *protocol.LoaderInput) error {
	for name, mm := range r.Payload {
		pe := &protocol.Payload{
			Meta: &protocol.MetaDesc{
				Name:      name,
				Attribute: map[string]string{},
			},
			Data: nil,
		}
		var err error
		for key, val := range mm {
			switch key {
			case "text":
				pe.Meta.DataType = protocol.MetaDesc_TEXT
				pe.Data, err = parsePayload(val)
			case "image":
				pe.Meta.DataType = protocol.MetaDesc_IMAGE
				pe.Data, err = parsePayload(val)
			case "video":
				pe.Meta.DataType = protocol.MetaDesc_VIDEO
				pe.Data, err = parsePayload(val)
			case "audio":
				pe.Meta.DataType = protocol.MetaDesc_AUDIO
				pe.Data, err = parsePayload(val)
			case "other":
				pe.Meta.DataType = protocol.MetaDesc_OTHER
				pe.Data, err = parsePayload(val)
			case "status":
				pe.Meta.Attribute[key] = convertStatus(val)
			default:
				pe.Meta.Attribute[key] = toString(val)
			}
			if err != nil {
				return err
			}
		}
		pb.Pl = append(pb.Pl, pe)
	}
	return nil
}

func getString(in map[string]interface{}, key string) string {
	if in == nil {
		return ""
	}
	return toString(in[key])
}

func convertStatus(in interface{}) string {
	switch i := in.(type) {
	case float64:
		return strconv.FormatFloat(i, 'f', 0, 64)
	default:
		return fmt.Sprintf("%v", i)
	}
}
func toString(in interface{}) string {
	switch i := in.(type) {
	case nil:
		return ""
	case string:
		return i
	case float64:
		return strconv.FormatFloat(i, 'b', -1, 64)
	case bool:
		return strconv.FormatBool(i)
	case int:
		return strconv.Itoa(i)
	default:
		return fmt.Sprintf("%v", i)
	}
}

func toStringMap(in map[string]interface{}) map[string]string {
	res := make(map[string]string, len(in))
	for key, val := range in {
		res[key] = toString(val)
	}
	return res
}

func parsePayload(in interface{}) ([]byte, error) {
	switch i := in.(type) {
	case string:
		return base64.StdEncoding.DecodeString(i)
	case []interface{}, map[string]interface{}:
		bf := bytes.NewBuffer(nil)
		err := json.NewEncoder(bf).Encode(i)
		return bf.Bytes(), err
	default:
		return nil, fmt.Errorf("not support type of payload:%v", reflect.TypeOf(in))
	}
}

func getDataType(s string) protocol.MetaDesc_DataType {
	switch s {
	case "text":
		return protocol.MetaDesc_TEXT
	case "image":
		return protocol.MetaDesc_IMAGE
	case "audio":
		return protocol.MetaDesc_AUDIO
	case "video":
		return protocol.MetaDesc_VIDEO
	default:
		return protocol.MetaDesc_OTHER
	}
}

func getDataTypeName(p protocol.MetaDesc_DataType) string {
	switch p {
	case protocol.MetaDesc_VIDEO:
		return "video"
	case protocol.MetaDesc_TEXT:
		return "text"
	case protocol.MetaDesc_AUDIO:
		return "audio"
	case protocol.MetaDesc_IMAGE:
		return "image"
	case protocol.MetaDesc_OTHER:
		return "other"
	default:
		return "other"

	}
}

type CommonResponse struct {
	Header  map[string]interface{} `json:"header"`
	Payload map[string]interface{} `json:"payload"`
}

func OutputToJson(pb *protocol.LoaderOutput, sid string, expect []*protocol.MetaDesc) *CommonResponse {
	res := &CommonResponse{
		Header: map[string]interface{}{
			"code":   0,
			"status": pb.Status,
			"sid":    sid,
		},
		Payload: make(map[string]interface{}),
	}

	for _, payload := range pb.Pl {
		pd := make(map[string]interface{})
		for key, val := range payload.GetMeta().GetAttribute() {
			pd[key] = val
		}
		dt := payload.GetMeta().GetDataType()
		if dt == protocol.MetaDesc_TEXT {

			pd[getDataTypeName(payload.GetMeta().GetDataType())] = string(payload.GetData())

		} else {
			pd[getDataTypeName(payload.GetMeta().GetDataType())] = base64.StdEncoding.EncodeToString(payload.GetData())

		}

		res.Payload[payload.GetMeta().GetName()] = pd
	}
	return res
}
