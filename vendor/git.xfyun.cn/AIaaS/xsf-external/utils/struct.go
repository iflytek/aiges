/*
* @file	struct.go
* @brief  rpc 请求与相应结构体定义以及若干操作方法
*
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */

package utils

import (
	"github.com/golang/protobuf/proto"
)

// todo: 使用对象池进行性能优化，补测一组直接new和对象池的性能对比测试
// 请求&响应消息体

// Req 请求响应结构
type Req struct {
	rd *ReqData
}

type Res struct {
	rd *ResData
}

// Data 结构
type Data struct {
	dm *DataMeta
}

func NewReq() *Req {
	req := new(Req)
	req.rd = new(ReqData)
	req.rd.S = new(Session)
	req.rd.Param = make(map[string]string)
	return req
}

func NewData() *Data {
	d := new(Data)
	d.dm = new(DataMeta)
	d.dm.Desc = make(map[string]string)
	return d
}

func NewRes() *Res {
	res := new(Res)
	res.rd = new(ResData)
	res.rd.Param = make(map[string]string)
	res.rd.S = new(Session)
	res.rd.S.Sess = make(map[string]string)

	return res
}

// 追加数据
func (d *Data) Append(b []byte) {
	d.dm.Data = append(d.dm.Data, b...)
}

// 获取数据
func (d *Data) Data() []byte {
	return d.dm.Data
}

//设置数据描述参数
func (d *Data) SetParam(k string, v string) {
	d.dm.Desc[k] = v
}

//设置数据描述参数
func (d *Data) GetParam(k string) string {
	return d.dm.Desc[k]
}

// 设置动作名
func (r *Req) SetOp(op string) {
	r.rd.Op = op
}
func (r *Req) Size() int {
	if nil != r.rd {
		return r.rd.Size()
	}
	return -1
}
func (r *Req) getSess() map[string]string {
	return r.rd.S.GetSess()
}
func (r *Req) GetPeerIp() (string, bool) {
	sessMap := r.getSess()
	if 0 == len(sessMap) {
		return "", false
	}
	peerIp, peerIpOk := sessMap[PEERADDR]
	return peerIp, peerIpOk
}
// 设置Req
func (r *Req) SetReq(req *ReqData) {
	r.rd = req
}

// 设置请求描述参数
func (r *Req) SetParam(k string, v string) {
	if nil == r.rd.Param {
		r.rd.Param = make(map[string]string)
	}
	r.rd.Param[k] = v
}

// 获取请求参数
func (r *Req) GetParam(k string) (string, bool) {
	v, ok := r.rd.Param[k]
	return v, ok
}

// 获取所有参数对
func (r *Req) GetAllParam() map[string]string {
	return r.rd.Param
}

// 追加数据
func (r *Req) Append(b []byte, desc map[string]string) {
	r.rd.Data = append(r.rd.Data, &DataMeta{Data: b, Desc: desc})
}

//
func (r *Req) AppendData(data *Data) {
	r.rd.Data = append(r.rd.Data, data.dm)
}

// 构造session，传进去的必须是从response中获取的
func (r *Req) Session(s string) error {
	sess := Session{}
	e := proto.UnmarshalText(s, &sess)

	//如果消息中已经携带trace id  以消息中trace id为准
	if nil != r.rd.S && len(r.rd.S.T) > 0 {
		sess.T = r.rd.S.T
		if nil != r.rd.S.GetSess() {
			sess.Sess = r.rd.S.GetSess()
		}
	}
	r.rd.S = &sess
	return e
}

// 追加session描述
func (r *Req) AppendSession(k string, v string) {
	if nil == r.rd.S.Sess {
		r.rd.S.Sess = make(map[string]string)
	}
	r.rd.S.Sess[k] = v
}

//  设置服务句柄参数描述，一般情况下请勿调用
func (r *Req) SetHandle(h string) {
	r.rd.S.H = h
}

func (r *Req) Handle() string {
	if nil == r.rd.S {
		return ""
	}
	return r.rd.S.H
}

// 设置trace日志的id
func (r *Req) SetTraceID(t string) {
	r.rd.S.T = t
}

// 获取trace日志的id
func (r *Req) TraceID() string {
	if nil == r.rd.S {
		return ""
	}
	return r.rd.S.T
}
func (r *Req) Req() *ReqData {
	return r.rd
}

func (r *Req) Op() string {
	return r.rd.Op
}

//获取data中的若干数据
func (r *Req) Data() []*Data {
	data := make([]*Data, 0, len(r.rd.Data))
	empty := true
	for _, d := range r.rd.Data {
		dm := new(Data)
		dm.dm = d
		data = append(data, dm)
		empty = false
	}
	if empty {
		return nil
	}
	return data
}

// 响应消息句柄
func (r *Res) GetParam(k string) (string, bool) {
	//todo 异常判断
	v, ok := r.rd.Param[k]
	return v, ok
}
func (r *Res) GetPeerIp() (string, bool) {
	sessMap := r.GetSess()
	if 0 == len(sessMap) {
		return "", false
	}
	peerIp, peerIpOk := sessMap[PEERADDR]
	return peerIp, peerIpOk
}
func (r *Res) SetParam(k string, v string) {
	r.rd.Param[k] = v
}

func (r *Res) GetAllParam() map[string]string {
	return r.rd.Param
}
func (r *Res) Size() int {
	if nil != r.rd {
		return r.rd.Size()
	}
	return -1
}
func (r *Res) GetData() []*DataMeta {
	return r.rd.Data
}

func (r *Res) GetRes() *ResData {
	return r.rd
}

func (r *Res) Session() string {
	if nil == r.rd {
		return ""
	}
	if nil == r.rd.S {
		return ""
	}
	return r.rd.S.String()
}

func (r *Res) Handle() string {
	if nil == r.rd || nil == r.rd.S {
		return ""
	}
	return r.rd.S.H
}

//  设置服务句柄参数描述，一般情况下请勿调用
func (r *Res) SetHandle(h string) {
	r.rd.S.H = h
}

// 设置trace日志的id
func (r *Res) SetTraceID(t string) {

	r.rd.S.T = t
}

// 获取trace日志的id
func (r *Res) TraceID() string {

	if nil == r.rd.S {
		return ""
	}

	return r.rd.S.T
}

// 追加session描述
func (r *Res) AppendSession(k string, v string) {
	r.rd.S.Sess[k] = v
}
func (r *Res) GetSess() map[string]string {
	return r.rd.S.GetSess()
}
func (r *Res) SetRes(rd *ResData) {
	r.rd = rd
}
func (r *Res) Res() *ResData {
	return r.rd
}
func (r *Res) Error() (int32, string) {
	return r.rd.Code, r.rd.ErrorInfo
}

//设置错误信息
func (r *Res) SetError(code int32, e string) {
	r.rd.Code = code
	r.rd.ErrorInfo = e
}

//拼接响应包中的数据
func (r *Res) AppendData(data *Data) {
	r.rd.Data = append(r.rd.Data, data.dm)
}
