package xsf

import (
	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

type Req = utils.Req
type Res = utils.Res
type Data = utils.Data
type Span = utils.Span

func NewReq() *Req {
	return utils.NewReq()
}
func NewReqEx(req_ *utils.ReqData) *Req {
	req := NewReq()
	if req_.S == nil {
		req_.S = new(utils.Session)
	}
	req.SetReq(req_)
	return req
}
func NewData() *Data {
	return utils.NewData()
}
func NewRes() *Res {
	return utils.NewRes()
}
