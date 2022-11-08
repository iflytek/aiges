package httproto

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/xfyun/aiges/protocol"
	"reflect"
	"strconv"
)

type Request struct {
	Header    map[string]interface{}            `json:"header"`
	Parameter map[string]map[string]interface{} `json:"parameter"`
	Payload   map[string]map[string]interface{} `json:"payload"`
}

func (r *Request) ConvertToPb(serviceName string, stat protocol.LoaderInput_SessState) (*protocol.LoaderInput, error) {
	in := &protocol.LoaderInput{
		ServiceId:   serviceName,
		ServiceName: serviceName,
		State:       stat,
		Headers:     toStringMap(r.Header),
		Params:      map[string]string{},
		Expect:      nil,
		Pl:          nil,
		SyncId:      0,
	}
	r.readParameter(in)
	return in, r.readPayload(in)
}

func (r *Request) readParameter(pb *protocol.LoaderInput) {

	for _, val := range r.Parameter {
		for sk, sv := range val {
			switch sv.(type) {
			case map[string]interface{}:
				desc := &protocol.MetaDesc{
					Name:      sk,
					DataType:  getDataType(toString(val["data_type"])),
					Attribute: map[string]string{},
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

func outputToJson(pb *protocol.LoaderOutput, sid string) *CommonResponse {
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

		pd[getDataTypeName(payload.GetMeta().GetDataType())] = base64.StdEncoding.EncodeToString(payload.GetData())

		res.Payload[payload.GetMeta().GetName()] = pd
	}
	return res
}
