package server

import (
	"encoding/base64"
	"common"
	"conf"
	"schemas"
	"errors"
	"github.com/json-iterator/go"
	"strconv"
)

type FrameReq struct {
	Business map[string]interface{} `json:"business"`
	Data     interface{}            `json:"data"`
	Common   CommonReq             `json:"common"`
	object   interface{}
}

const(
	CmdCLose =  "close"
	CmdKeepAlive = "kpal"
)

type CommonReq struct {
	AppId          string `json:"app_id"`
	Uid            string `json:"uid"`
	DeviceId       string `json:"device_id"`
	DeviceImei     string `json:"device.imei"`
	DeviceImsi     string `json:"device.imsi"`
	DeviceMac      string `json:"device.mac"`
	DeviceOther    string `json:"device.other"`
	AuthAuthId     string `json:"auth.auth_id"`
	AuthAuthSource string `json:"auth.auth_source"`
	NetType        string `json:"net.type"`
	NetIsp         string `json:"net.isp"`
	AppVer         string `json:"app.ver"`
	Wscid          string `json:"cid"`
	Cmd            string `json:"cmd"`
	Connection     string `json:"connection"`
	Sub            string `json:"sub"`
}



func (f *FrameReq)CommonString()string  {
	b:=common.NewStringBuilder()
	c:=f.Common
	b.Append("appid=",c.AppId," imei=",c.DeviceImei," imsi=",c.DeviceImei)
	b.Append(" other=",c.DeviceOther," mac=",c.DeviceMac," id=",c.DeviceId)
	b.Append(" uid=",c.Uid).Append(" netisp=",c.NetIsp)
	b.Append(" net_type=",c.NetType).Append(" ver=",c.AppVer)
	return b.ToString()
}

func (fr *FrameReq) GetStatus() (int, error) {
	status, err := schemas.GetByJPath(fr.Data, "$.status")
	if err != nil {
		return 0, err
	}
	if status == nil {
		return 0, nil
	}
	return int(status.(float64)), nil
}

func (fr *FrameReq) GetStatusBybiz(biz *ServerBiz) (int, error) {
	gdata:=biz.GetUpCall().DataList
	if gdata==nil || len(gdata)==0{
		return -1,errors.New("cannot find status because datalist is nil")
	}
	return int(gdata[0].Status),nil
}

func (fr *FrameReq)SetBizStatus(s *Session,biz *ServerBiz)  {
	for _,v:= range  biz.UpCall.DataList{
		v.Status=GeneralData_DataStatus(s.Status)
	}
}


func NewFrameReqWithConn(t *Trans,msg *[]byte,sid string, frameId int) (*FrameReq, *Error) {
	var req interface{}
	if err := jsoniter.Unmarshal(*msg, &req);err !=nil{
		return nil, NewError(int(ErrorCodeJSONParsing), "parse request json error", sid, frameId)
	}

	//schema 校验
	if err:=CheckDataBySchema(t.Mapping,req,sid,frameId);err!=nil{
		fr:= getFrameReq(convertMap(req))

		return fr,err
	}
	//把common参数设置到business中
	fr:= getFrameReq(convertMap(req))
	fr.SetBusiFromCommon()
	return fr,nil
}

func CheckDataBySchema(mp *schemas.RouteMapping, doc interface{}, sid string, frameId int) *Error {
	if !conf.Conf.Schema.Enable {
		return nil
	}
	msg, e := schemas.ValidateByMapping(mp,doc)
	if e != nil {
		return NewError(int(ErrorCodeGetUpCall),"param validate error:"+ msg, sid, frameId)
	}
	return nil
}

func getFrameReq(req map[string]interface{}) *FrameReq {
	//common.Logger.Infof("common:%v",req)
	frameReq := &FrameReq{
		Business: convertMap(req["business"]),
		Data:     req["data"],
		Common:   getCommonReq(convertMap(req["common"])),
		object:   req,
	}
	return frameReq
}

func getCommonReq(common map[string]interface{}) CommonReq {
	if common == nil || len(common) == 0 {
		return CommonReq{}
	}
	commonReq := CommonReq{
		AppId:          convertStr(common["app_id"]),
		Wscid:          convertStr(common["cid"]),
		Uid:            convertStr(common["uid"]),
		DeviceId:       convertStr(common["device_id"]),
		DeviceImei:     convertStr(common["device.imei"]),
		DeviceImsi:     convertStr(common["device.imsi"]),
		DeviceMac:      convertStr(common["device.mac"]),
		DeviceOther:    convertStr(common["device.other"]),
		NetType:        convertStr(common["net.type"]),
		NetIsp:         convertStr(common["net.isp"]),
		AppVer:         convertStr(common["app.ver"]),
		Sub:			convertStr(common["sub"]),
		Cmd:            convertStr(common["cmd"]),
	}
	return commonReq
}

func convertMap(req interface{}) map[string]interface{} {
	if m,ok:=req.(map[string]interface{});ok {
		return m
	}
	return map[string]interface{}{};
}
func convertStr(source interface{}) string {
	if s ,ok:=source.(string);ok {
		return s
	}
	return ""
}

func CheckAppid(s *Session,appid string) *Error {
	if !conf.Conf.Auth.EnableAppidCheck{
		return nil
	}

	if appid==""{
		return NewError(int(ErrorInvalidAppid),"app_id cannot be empty",s.Sid,0)
	}
	//if realAppid=="" 那么可能不是走kong，或者kong没有开启鉴权，也放过
	if realAppid :=s.Ctx.GetHeader("X-Consumer-Username");(realAppid !=appid)&& realAppid !=""{
		common.Logger.Errorf("invalid app_id ,expected:%s sid=%s",appid,s.Sid)
		return NewError(int(ErrorInvalidAppid),"invalid app_id："+appid,s.Sid,0)
	}
	return nil
}

func (fr *FrameReq) GetServerBiz(s *Session) (*ServerBiz, *Error) {
	dataList, err := getDataList(s, fr.Data)
	if err != nil {
		return nil, err
	}
	gr := &GlobalRoute{
		SessionId:  s.Sid,
		TraceId:    s.Sid,
		UpRouterId: XsfCallBackAddr,
		GuiderId:   conf.Centra.Service,
		Appid:      s.Appid,
		Uid:        s.Uid,
		ClientIp:   s.ClientIp,
	}

	upCall := &UpCall{
		Call:         s.Call,
		SeqNo:        int32(s.SeqNo),
		From:         s.From,
		Sync:         false,
		BusinessArgs: getBusinessArgs(fr.Business),
		Session:      s.sessionMap,
		DataList:     dataList,
	}
	serverBiz := &ServerBiz{
		Version:     conf.Centra.Version,
		MsgType:     ServerBiz_UP_CALL,
		GlobalRoute: gr,
		UpCall:      upCall,
	}
	return serverBiz, nil
}

//把commo参数设置到business中
func (fr *FrameReq)SetBusiFromCommon()  {
	if fr.Business==nil{
		return
	}
	fr.setbusi("imsi",fr.Common.DeviceImsi)
	fr.setbusi("imei",fr.Common.DeviceImei)
	fr.setbusi("net_type",fr.Common.NetType)
	fr.setbusi("cver",fr.Common.AppVer)
	fr.setbusi("device_mac",fr.Common.DeviceMac)
	fr.setbusi("device_oth",fr.Common.DeviceOther)
	fr.setbusi("net_isp",fr.Common.NetIsp)
	fr.setbusi("cver",fr.Common.AppVer)
	fr.setbusi("uid",fr.Common.Uid)

}

func (fr *FrameReq)setbusi(k string,v string)  {
	if v==""{
		return
	}
	fr.Business[k]=v
}

//获取business参数
func getBusinessArgs(business map[string]interface{}) map[string]string {
	businessArgs := make(map[string]string)
	for k, v := range business {
		businessArgs[k] = common.ConvertToString(v)
	}
	return businessArgs
}

//获取dataList
func getDataList(s *Session, data interface{}) ([]*GeneralData, *Error) {
//	result, err := schemas.GetUpCallReqByCall(s.Key, data)
	result,err:=s.Mapping.ResolveUpCallReq(data)
	if err != nil {
		common.Logger.Errorf("sid:%s,seqNo:%d=>根据请参数获取上行请求参数失败,失败的原因是:%s", s.Sid, s.SeqNo, err.Error())
		return nil, NewErrorByError(ErrorCodeGetUpCall, err, s.Sid, s.SeqNo)
	}

	dataList := make([]*GeneralData, len(result))
	for index, elem := range result {
		data, err := base64.StdEncoding.DecodeString(getData(elem))
		if err != nil {
			common.Logger.Errorf("%s:对客户端发送帧%v数据解码失败,失败的原因是:%s", s.Sid, s.SeqNo, err.Error())
			return nil, NewError(int(ErrorCodeDecoding), "parse base64 string error", s.Sid, s.SeqNo)
		}
		//
		gd := &GeneralData{
			DataId:   getString("id",elem),
			FrameId:  uint32(s.SeqNo)-1,
			DataType: GeneralData_DataType(getInt("data_type",elem)),
			Status:   GeneralData_DataStatus(getInt("status",elem)),
			Format:   getString("format",elem),
			Encoding: getString("encoding",elem),
			DescArgs:getMap("desc_args",elem),
			Data:     data,
		}
		dataList[index] = gd
		s.dataargs["format["+strconv.Itoa(index)+"]"] = getString("format",elem)
		s.dataargs["encoding["+strconv.Itoa(index)+"]"] = getString("encoding",elem)
		s.dataargs["data_type["+strconv.Itoa(index)+"]"] = elem["data_type"]
		s.dataargs["len["+strconv.Itoa(index)+"]"] = len(data)
	}
	return dataList, nil
}

func getData(elem map[string]interface{}) string {
	if e,ok:=elem["data"].(string) ;ok {
		return e
	}

	return ""
}


func getString(key string,elem map[string]interface{})string  {
	if e,ok:=elem[key].(string);ok{
		return e
	}
	return ""
}

func getInt(key string,elem map[string]interface{})int  {
	if e,ok:=elem[key].(float64);ok{
		return int(e)
	}
	if e,ok:=elem[key].(int32);ok{
		return int(e)
	}
	return 0
}


func getMap(key string,elem map[string]interface{})map[string][]byte  {
	if m,ok:=elem[key].(map[string]interface{});ok{
		var tm = make(map[string][]byte)
		for k, v := range m {
			vs:=common.ConvertToString(v)
			tm[k] = common.ToBytes(&vs)
		}
		return tm
	}
	return nil

}