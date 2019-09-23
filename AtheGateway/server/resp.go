package server

import (
	"common"
	"schemas"
	"github.com/json-iterator/go"
	"errors"
)

const sucessee = 0


type FrameResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Sid     string      `json:"sid,omitempty"`
	Uid     string      `json:"uid,omitempty"`
	Wscid     string      `json:"cid,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

//type DataResp struct {
//	Result interface{} `json:"result,omitempty"`
//
//	Status int `json:"status"`
//}

func NewFrameRespByServerBiz(generalData []*GeneralData, s *Session, sid string, wscid string, frameId int) (*FrameResponse, *Error) {
	result, err := getResultByCall(s, generalData, sid, frameId)
	if err != nil {
		return nil, err
	}
	return &FrameResponse{
		Code:    sucessee,
		Message: "success",
		Data:    result,
		Sid:sid,
		Wscid:wscid,
	}, nil
}

//创建错误得响应结果
func NewFrameRespByError(err *Error, cid string) *FrameResponse {
	return &FrameResponse{
		Code:    err.Code,
		Message: err.Msg,
		Sid:     err.Sid,
		Wscid:cid,
		//	Uid:     uid,
	}
}

//根据回调结果获取响应结果
func getResultByCall(s *Session, dataList []*GeneralData, sid string, frameId int) (interface{}, *Error) {
	parseResult, error := parseDataList(dataList,s, sid, frameId)
	if error != nil {
		common.Logger.Errorf("sid=%s,seqNo=%d:解析异步回调DataList失败,失败得原因:%s", sid, frameId, error.Msg)
		return nil, error
	}
//	resp, err := schemas.GetRespByCall(key, parseResult)
	resp,err:=s.Mapping.ResolveResp(parseResult)
	if err != nil {
		common.Logger.Errorf("%s:根据配置映射获取响应结果失败 call:%s,失败的原因是:%s", sid, frameId, err.Error())
		return nil, NewError(int(ErrorCodeGetRespCall), "cannot parse server data", sid, frameId)
	}
	return resp, nil
}

//解析DataList
func parseDataList(dataList []*GeneralData,s *Session, sid string, frameId int) (interface{}, *Error) {
	result := make([]map[string]interface{}, len(dataList))
	//dataTypes := schemas.GetRespDataType(key)
	dataTypes:=s.Mapping.ResponseData.DataType
	for index, elem := range dataList {
		temp := make(map[string]interface{})
		temp["data_id"] = elem.DataId
		temp["frame_id"] = elem.FrameId
		temp["desc_args"] = getDescArgs(elem.DescArgs)
		temp["status"] = elem.Status
		temp["format"] = elem.Format
		temp["encoding"] = elem.Encoding
		temp["data_type"] = elem.DataType
		if dataTypes != nil && (index< len(dataTypes)) {
			data, err := parseData(elem.Data, sid, frameId, dataTypes[index])
			if err != nil {
				common.Logger.Errorf("%s:%d解析Data失败,失败的原因是:%s", sid, frameId, err.Error())
				return nil, NewErrorByError(ErrorCodePASEDATA, errors.New("server cannot parse response data"), sid, frameId)
			}
			temp["data"] = data
		}else{
			return nil,NewError(int(ErrorServerError),"server config error: too many datas in resp",sid,frameId)
		}
		result[index] = temp
	}
	return result, nil
}

//解析Data
func parseData(data []byte, sid string, frameId int, dataType uint32) (interface{}, error) {
	if data != nil && len(data) != 0 {
		switch dataType {
		case schemas.DATA_TYPE_STRING:
			return common.ToString(data[:]), nil
		case schemas.DATA_TYPE_BYTE:
			return common.EncodingTobase64String(data), nil
		case schemas.DATA_TYPE_JSON:
			var parseResult interface{}
			err := jsoniter.Unmarshal(data, &parseResult)
			if err != nil {
				common.Logger.Errorf("%s:%d解析Data失败,失败的原因是:%s", sid, frameId, err.Error())
				return common.EncodingTobase64String(data), nil
			}
			return parseResult, nil
		default: //默认
			return common.EncodingTobase64String(data), nil
		}
	}
	return nil, nil
}

//获取描述参数
func getDescArgs(args map[string][]byte) map[string]string {
	descArgs := make(map[string]string)
	for key, value := range args {
		descArgs[key] = string(value)
	}
	return descArgs
}

//创建第一帧的响应结果
func NewFirstFrameRespByUpResult(s *Session, upr *UpResult) (*FrameResponse, *Error) {
	result, err := getResultByUpResult(s,upr)
	if err != nil {
		return nil, NewErrorByError(ErrorCodeJSONParsing, err, s.Sid, 1)
	}
	//if result == nil{
	//	result = map[string]interface{}{
	//		"status":1,
	//	}
	//}
	return &FrameResponse{
		Code:    sucessee,
		Message: "success",
		//	Uid:     uid,
		Sid:  s.Sid,
		Data: result,
	}, nil
}
//创建非第一帧的响应结果
func NewRespByUpResult(s *Session, upr *UpResult) (*FrameResponse, *Error) {
	result, err := getResultByUpResult(s,upr)
	if err != nil {
		return nil, NewErrorByError(ErrorCodeJSONParsing, err, s.Sid, 1)
	}
	//data := make(map[string]interface{}, 2)
	//data["status"] = 1
	if result == nil{
		return nil,nil
	}
	return &FrameResponse{
		Code:    sucessee,
		Message: "success",
		//	Uid:     uid,
		Sid:  s.Sid,
		Data: result,
	}, nil
}
//根据上行结果获取result
func getResultByUpResult(s *Session,upr *UpResult) (interface{}, error) {
	if upr.DataList == nil || len(upr.DataList) == 0 {
		return nil, nil
	}
	result,err:=getResultByCall(s,upr.DataList,s.Sid,s.SeqNo-1)

	if err != nil {
		return nil, err
	}
	return result, nil
}
