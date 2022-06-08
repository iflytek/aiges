/*
* @file	remote.go
* @brief  远端负载均衡实现
*         从远端负载均衡获取可用的机器列表
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */

package xsf

import (
	"context"
	"errors"
	"github.com/xfyun/xsf/utils"
	"strconv"
	"time"
)

var (
	CLOUD      = "cloud"
	NBSET      = "nbest"
	ROUTER     = "router"
	UID        = "uid"
	SVC        = "svc"
	SUBSVC     = "subsvc"
	SUBROUTER  = "sub_router"
	LBOPGET    = "getServer"
	LBOPSET    = "setServer"
	FAILEDNODE = "broken"
)

type remoteLB struct {
	lbi     LBI
	conns   *connPool
	o       *conOption
	retry   int
	timeout int
}

func newRemoteLB(o *conOption, lbi LBI, conns *connPool) *remoteLB {
	rrl := new(remoteLB)
	rrl.retry = o.lbretry
	rrl.timeout = o.lbtimeout
	rrl.lbi = lbi
	rrl.conns = conns

	return rrl
}

func (rlb *remoteLB) Find(param *LBParams) ([]string, []string, error) {
	param.log.Infow("remoteLB:Find", "logId", param.logId)
	param.log.Infow("print param", "param", param, "logId", param.logId)

	if param.directEngIp != "" {
		return []string{param.directEngIp}, nil, nil
	}
	svc := param.svc
	//nbest := param.nbest
	param.svc = param.name
	param.nbest = rlb.retry
	addr, _, e := rlb.lbi.Find(param)
	param.log.Infow("fn:remoteLb-Find", "addr", addr, "logId", param.logId)
	var raddrs []string
	var re error
	var grpcErr error
	for i := 0; i < rlb.retry+1; i++ {

		if e == nil && len(addr) > 0 {
			//	for _,v :=range addr {
			param.log.Infow("remoteLB:Find", "i", i, "retry", rlb.retry, "lbAddr", addr, "callLbAddr", addr[i%len(addr)], "logId", param.logId)

			//发起remote的LB请求，需要包含LB的超时
			param.svc = svc
			//param.nbest = nbest
			param.nbest = 1 //每次重新交互取数据，硬编码的问题后续改吧
			raddrs, re, grpcErr = rlb.callLB(addr[i%len(addr)], param)
			if re == nil {
				return raddrs, nil, re
			} else {
				param.log.Infow("retry callLB", "logId", param.logId)
				continue
				//return nil, EINVALIDADDR
			}
			//	}

		} else {
			param.log.Errorw("remoteLB:rlb.lbi.Find", "error", e, "logId", param.logId, "addr", addr)
			return nil, nil, INVALIDLB
		}
	}
	if grpcErr != nil {
		return nil, nil, INVALIDRMLB
	}
	return nil, nil, EINVALIDADDR
}

//调用远端服务
func (rlb *remoteLB) callLB(addr string, param *LBParams) (addrs []string, businErr error, grpcErr error) {

	var sp *Span
	if param.span != nil {
		sp = param.span.Next(utils.CliSpan)
		if sp != nil {
			sp.WithName(param.name).WithTag("saddr:", addr).Start()

			//defer sp.flush()
			defer sp.End()
		}
	}
	// 构建请求包
	req := new(utils.ReqData)
	if param.ext == nil {
		req.Param = make(map[string]string)
	} else {
		req.Param = param.ext
	}

	req.Param[CLOUD] = func() string {
		if len(param.ext) == 0 {
			return ""
		}
		return param.ext[CLOUD]
	}()
	req.Param[NBSET] = strconv.Itoa(param.nbest)
	req.Param[ROUTER] = param.catgory
	req.Param[SUBROUTER] = param.svc
	req.Param[SVC] = param.catgory
	req.Param[SUBSVC] = param.svc
	req.Param[FAILEDNODE] = param.failed
	if uidStr, uidStrOk := param.ext[UID]; uidStrOk {
		req.Param[UID] = uidStr
	}
	for extK, extV := range param.ext {
		req.Param[extK] = extV
	}
	req.S = new(utils.Session)
	req.S.T = param.span.Meta()
	req.Op = LBOPGET
	c, e := rlb.conns.get(addr, "")
	if e != nil {
		if sp != nil {
			sp.WithErrorTag(e.Error())
		}
		param.log.Errorw("remoteLB:callLB ",
			"error", e.Error(), "addr", addr, "logId", param.logId)
		return nil, e, e
	}
	prx := utils.NewXsfCallClient(c)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rlb.timeout)*time.Millisecond)
	defer cancel()
	//ctx:=context.Background()
	res, e := prx.Call(ctx, req)
	if e != nil {
		param.log.Errorw("remoteLB:callLB ",
			"error", e.Error(), addr, rlb.timeout, "router", param.catgory, "subRouter", param.svc, "logId", param.logId)

		if sp != nil {
			sp.WithErrorTag(e.Error())
		}
		return nil, e, e
	}

	if res != nil && len(res.ErrorInfo) != 0 {
		param.log.Errorw("remoteLB:callLB ",
			"res_error", res.ErrorInfo, addr, rlb.timeout, "router", param.catgory, "subRouter", param.svc, "logId", param.logId)
		if sp != nil {
			sp.WithErrorTag(res.ErrorInfo)
		}
		return nil, errors.New(res.ErrorInfo), nil
	}
	addrs = make([]string, 0, rlb.retry)
	for _, v := range res.Data {
		addrs = append(addrs, string(v.Data))
	}
	return addrs, nil, nil
}
