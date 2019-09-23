/*
* @file	caller.go
* @brief  框架客户端的接口定义
*         封装了rpc以及框架行为的接口，让调用者以native的方式进行rpc调用，使用者也无需关注底层rpc的使用方式，降低服务编码门槛。
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */

package xsf

import (
	"errors"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"golang.org/x/net/context"
	"strconv"
	"sync/atomic"
	"time"
)

/*
1、连接超时
2、连接池的大小
3、负载地址的维护
4、配置的嵌入
*/

type Req = utils.Req
type Res = utils.Res
type Data = utils.Data

// 一次调用的对象
type Caller struct {
	//prx *utils.XsfCallClient  // 调用代理
	cli *Client
	//span trace
	//sp *utils.Span 已和大坤沟通，无需保持状态

	retry     int
	timeLimit time.Duration //重试时间上限制

	tm      int64 //超时时间
	lbname  string
	catgory string

	apiVersion string

	lbExt map[string]string

	hashKey atomic.Value
}

type callCfg struct {
	retry bool
	logId string
	req   *Req
}
type callOpt func(opt *callCfg)

func withCallReq(req *Req) callOpt {
	return callOpt(func(opt *callCfg) {
		opt.req = req
	})
}
func withCallLogId(logId string) callOpt {
	return callOpt(func(opt *callCfg) {
		opt.logId = logId
	})
}
func withCallRetry(retry bool) callOpt {
	return callOpt(func(opt *callCfg) {
		opt.retry = retry
	})
}

// SessStat 用于表示会话的请求的状态
type SessStat int32

const (
	// 创建会话
	CREATE SessStat = iota
	// 后续会话
	CONTINUE
	//  非会话模式
	ONESHORT
	// 无效会话
	INVALID
)

// InitClient 创建一个Client的全局操作句柄，需要设置到具体的caller对象中。
//
// 参数说明：
//
// @cname: 客户端名称，这个参数在读取配置时会用作定位配置的selection。
//
// @mode: 配置文件支持的模式，详见CfgMode
//
// @o...: 配置相关的其他属性,详见 CfgOpt
func InitClient(cname string, mode CfgMode, o ...CfgOpt) (*Client, error) {
	return NewClient(cname, mode, o...)
}

// InitClientWithCfg 创建一个Client的全局操作句柄，需要设置到具体的caller对象中。
//
// 参数说明：
//
// @cname: 客户端名称，这个参数在读取配置时会用作定位配置的selection。
//
// @cfg:外部已经初始化完成的配置操作句柄
//
// @o...: 配置相关的其他属性,详见 CfgOpt
func InitClientWithCfg(cname string, cfg *utils.Configure, o ...CfgOpt) (*Client, error) {
	var e error
	cli, e := NewClientWithCfg(cname, cfg, o...)
	return cli, e
}

// NewReq 生成请求对应
func NewReq() *Req {
	return utils.NewReq()
}

// NewRes 生成响应对应
func NewRes() *Res {
	return utils.NewRes()
}

// NewData 生成元数据的对象
func NewData() *Data {
	return utils.NewData()
}

// NewCaller 生成一次调用的调用对像，建议作为局部变量使用
//
// 参数说明：
//
// @cli : 调用对象所依赖的Client 实体，是InitClientWithCfg或者InitClient生成的。
func NewCaller(cli *Client) *Caller {
	if nil == cli {
		return nil
	}
	c := new(Caller)
	c.cli = cli
	c.retry = 0
	c.apiVersion = defaultApiVersion
	return c
}

// TimeOut 当次调用的超时时间设置，该接口目前不生效
func (c *Caller) TimeOut(mesc int64) *Caller {
	c.tm = mesc
	return c
}

// WithLBParams 设置当次调用，LB相关的参数。当前，仅在LB在Remote下生效
//
// 参数说明:
//
// @lbaname: Remote LoadBalance 的服务名
//
// @busin: 请求目标服务的业务类型
//
// @ext: 扩展参数，根据LoadBalance实际需要进行设置
func (c *Caller) WithLBParams(lbaname string, busin string, ext map[string]string) {
	c.lbname = lbaname
	c.catgory = busin
	c.lbExt = ext
}

func (c *Caller) WithApiVersion(apiVersion string) {
	c.apiVersion = apiVersion
}

// WithRetry 设置调用的重试次数。当前未生效
func (c *Caller) WithRetry(retry int) {
	c.retry = retry
}

func (c *Caller) WithHashKey(key string) {
	c.hashKey.Store(key)
}
func (c *Caller) GetHashKey() (key string) {
	key, _ = c.hashKey.Load().(string)
	return
}

// Call 发起无状态的调用。
//
// 参数说明：
//
// @service: 目标服务的服务名
//
// @op: 请求的操作名，建议和目标服务的实际处理的函数同名
//
// @r: 请求消息对象
//
// @tm: 请求超时时间，根据具体情况进行设置
//
// 返回值说明：
//
// @s: 响应消息对象
//
// @errcode: 处理错误码，沿用msp_errors
//
// @e: 接口的错误描述
func (c *Caller) Call(service string, op string, r *Req, tm time.Duration) (s *Res, errcode int32, e error) {
	//参数异常
	if 0 == len(op) {
		return nil, INVAILDPARAM, EINAILDOP
	}
	r.SetOp(op)
	return c.oneShortCall(c.apiVersion, service, r, tm)
}

//todo 这里待完成，返回结果的数据结构待设计
func (c *Caller) CallAll(service string, op string, r *Req, tm time.Duration) (s *Res, errcode int32, e error) {
	//参数异常
	if 0 == len(op) {
		return nil, INVAILDPARAM, EINAILDOP
	}
	r.SetOp(op)
	return c.oneShortCall(c.apiVersion, service, r, tm)
}

// SessionCall 发起有状态的调用。
//
// 参数说明：
//
// @ss: 会话状态
//
// @service: 目标服务的服务名
//
// @op: 请求的操作名，建议和目标服务的实际处理的函数同名
//
// @r: 请求消息对象
//
// @tm: 请求超时时间，根据具体情况进行设置
//
// 返回值说明：
//
// @s: 响应消息对象
//
// @errcode: 处理错误码，沿用msp_errors
//
// @e: 接口的错误描述
func (c *Caller) SessionCall(ss SessStat, service string, op string, r *Req, tm time.Duration) (s *Res, errcode int32, e error) {
	// 输入异常
	if nil == r {
		return nil, INVAILDDATA, EINAILIDDATA
	}
	if CONTINUE != ss {
		if 0 == len(op) || 0 == len(service) {
			return nil, INVAILDPARAM, EINAILDOP
		}
	}

	// 进入函数调用
	r.SetOp(op)
	switch ss {
	case CREATE:
		return c.createCall(c.apiVersion, service, r, tm)
	case CONTINUE:
		return c.continueCall(c.apiVersion, service, r, tm)
	case ONESHORT:
		return c.oneShortCall(c.apiVersion, service, r, tm)
	default:
		return nil, INVAILDPARAM, EINAILDOP
	}
}

// SessionCall 像指定服务地址发起的调用。
//
// 参数说明：
//
// @service: 目标服务的服务名
//
// @op: 请求的操作名，建议和目标服务的实际处理的函数同名
//
// @addr: 服务地址，格式为ip:port
//
// @r: 请求消息对象
//
// @tm: 请求超时时间，根据具体情况进行设置
//
// 返回值说明：
//
// @s: 响应消息对象
//
// @errcode: 处理错误码，沿用msp_errors
//
// @e: 接口的错误描述
func (c *Caller) CallWithAddr(service string, op string, addr string, r *Req, tm time.Duration) (s *Res, errcode int32, e error) {
	logId := logSidGeneratorInst.GenerateSid("CallWithAddr")
	c.cli.Log.Infow("CallWithAddr enter", "logId", logId)
	//// trace 日志句柄
	//var span *utils.Span
	//if len(r.TraceID()) > 0 {
	//	span = utils.FromMeta(r.TraceID(), "", 0, service, utils.CliSpan)
	//} else {
	//	span = utils.NewSpan(utils.CliSpan)
	//}
	//span = span.Next(utils.CliSpan)
	//span = span.WithRpcCallType()
	//span.WithName(op).Start()
	r.SetOp(op)
	//span.WithTag("call", "CallWithAddr")
	//if op != LBOPSET && op != PING {
	//	defer span.Flush()
	//}
	//defer span.End()

	// 根据地址获取本地调用代理（stub）
	prx, e := c.getPrxyWithAddr(addr, logId)
	if nil != e {
		c.cli.Log.Errorw("CallWithAddr:getPrxyWithAddr",
			"error", e.Error(), "addr", addr, "logId", logId)

		//span.WithFuncTag("getPrxyWithAddr")
		//span.WithErrorTag(e.Error())
		//span.WithRetTag(strconv.Itoa(int(NOUSEFULCONN)))

		return nil, NOUSEFULCONN, e
	}

	//r.SetTraceID(span.Meta())
	r.SetOp(op)

	callBase := time.Now()
	s, errcode, e = c.doCall(prx, r, tm)
	c.cli.Log.Infow("record doCall perf",
		"fn", "CallWithAddr", "logId", logId, "addr", addr, "dur", time.Since(callBase).Nanoseconds())

	if nil != e {
		c.cli.Log.Errorw("CallWithAddr:doCall",
			"logId", logId, "error", e.Error(), "addr", addr)

		//span.WithFuncTag("doCall")
		//span.WithErrorTag(e.Error())
		//span.WithTag("addr", addr)
		//span.WithRetTag(string(errcode))
	}
	//span.WithRetTag(strconv.Itoa(int(errcode)))
	//span.WithTag("handle", s.Handle())

	return s, errcode, e
}

// prxCall 基于grpc 本地stub的调用实现
func (c *Caller) prxCall(prx utils.XsfCallClient, r *Req, tm time.Duration) (s *Res, errcode int32, e error) {
	c.cli.Log.Infow("prxCall enter")

	s = utils.NewRes()
	ctx, cancel := context.WithTimeout(context.Background(), tm)
	defer cancel()

	rd, e := prx.Call(ctx, r.Req())
	if nil != e {
		return nil, NETWORKEXCEPT, e
	}
	if len(rd.ErrorInfo) > 0 {
		e = errors.New(rd.ErrorInfo)
	}
	s.SetRes(rd)
	return s, rd.Code, e
}

// createCall 会话状态下的首次会话的调用模式
func (c *Caller) createCall(apiVersion, service string, r *Req, tm time.Duration) (s *Res, errcode int32, e error) {
	logId := logSidGeneratorInst.GenerateSid("createCall")
	c.cli.Log.Infow("createCall enter", "logId", logId)
	// trace 日志句柄
	//var span *utils.Span
	//if len(r.TraceID()) > 0 {
	//	span = utils.FromMeta(r.TraceID(), "", 0, service, utils.CliSpan)
	//} else {
	//	span = utils.NewSpan(utils.CliSpan)
	//}
	//
	//spnxt := span.Next(utils.CliSpan)
	//spnxt = spnxt.WithRpcCallType()
	//spnxt.WithName(r.Op()).Start()
	//spnxt.WithTag("call", "createCall")
	//defer spnxt.Flush()
	//defer spnxt.End()

	tmpFlag, prx, conn, e := c.getPrxy(apiVersion, service, nil, withCallLogId(logId), withCallReq(r))

	if tmpFlag {
		r.SetParam(tmpKV.tmpK, tmpKV.tmpV)
	}
	if nil != e {
		c.cli.Log.Errorw("createCall:getPrxy", "logId", logId, "error", e.Error(), "service", service)

		errcode = c.err2errCode(e)

		//spnxt.WithFuncTag("getPrxy")
		//spnxt.WithErrorTag(e.Error())
		//span.WithRetTag(strconv.Itoa(int(errcode)))

		return nil, errcode, e
	}

	// todo: 失败的addr需要计数，进行降级等
	//r.SetTraceID(spnxt.Meta())

	callBase := time.Now()
	s, errcode, e = c.doCall(prx, r, tm)
	c.cli.Log.Infow("record doCall perf", "fn", "createCall", "dur", time.Since(callBase).Nanoseconds())

	if nil != e {
		c.cli.Log.Errorw("createCall:doCall", "logId", logId, "error", e.Error(), "service", service, "addr", conn.Addr)

		//spnxt.WithFuncTag("doCall")
		//spnxt.WithErrorTag(e.Error())
		//spnxt.WithTag("addr", conn.Addr)
		//spnxt.WithRetTag(string(errcode))
	}
	//spnxt.WithRetTag(strconv.Itoa(int(errcode)))

	// 进入重试,换机器
	if nil != e && c.retry > 0 {
		return c.retryCall(apiVersion, service, r, tm, withCallReq(r))
	}

	return s, errcode, e
}

// continueCall 会话模式下的持续请求
func (c *Caller) continueCall(apiVersion, service string, r *Req, tm time.Duration) (s *Res, errcode int32, e error) {

	logId := logSidGeneratorInst.GenerateSid("continueCall")

	c.cli.Log.Infow("continueCall enter", "logId", logId)

	//var span *utils.Span
	//if len(r.TraceID()) > 0 {
	//	span = utils.FromMeta(r.TraceID(), "", 0, service, utils.CliSpan)
	//} else {
	//	span = utils.NewSpan(utils.CliSpan)
	//}
	//
	//span = span.Next(utils.CliSpan)
	//span = span.WithRpcCallType()
	//
	//span.WithName(r.Op()).Start()
	//span.WithTag("call", "continueCall")
	//
	//defer span.Flush()
	//defer span.End()

	var h string
	h = r.Handle()
	if len(h) < 26 {
		c.cli.Log.Errorw("continueCall:Handle", "logId", logId, "error", EINAILDHDL.Error(), "handle", h)
		//span.WithFuncTag("getHandle")
		//span.WithTag("xsf-handle", h)
		//span.WithErrorTag(EINAILDHDL.Error())
		//span.WithRetTag(strconv.Itoa(int(INVAILDHANDLE)))
		return nil, INVAILDHANDLE, EINAILDHDL
	}

	// 根据服务句柄获取stub，handle为26字节的实现
	prx, e := c.getPrxyWithHandle(h)
	if nil != e {

		c.cli.Log.Errorw("continueCall:getPrxyWithHandle", "logId", logId, "error", e.Error(), "service", service, "handle", h)
		//span.WithFuncTag("getPrxyWithHandle")
		//span.WithTag("xsf-handle", h)
		//span.WithErrorTag(e.Error())
		//span.WithRetTag(strconv.Itoa(int(NOUSEFULCONN)))
		return nil, NOUSEFULCONN, e
	}

	//r.SetTraceID(span.Meta())

	callBase := time.Now()
	s, errcode, e = c.doCall(prx, r, tm)
	c.cli.Log.Infow("record doCall perf", "handle", h, "fn", "continueCall", "dur", time.Since(callBase).Nanoseconds())

	if nil != e {
		c.cli.Log.Errorw("continueCall:doCall", "logId", logId, "error", e.Error(), "service", service, "handle", h)
		//span.WithFuncTag("doCall")
		//span.WithErrorTag(e.Error())
		//span.WithTag("xsf-handle", h)

	}
	//span.WithRetTag(strconv.Itoa(int(errcode)))

	return s, errcode, e
}

func (c *Caller) err2errCode(err error) (errcode int32) {

	//错误码区分
	switch err {
	case EINVALIDADDR:
		{
			errcode = INVAILDLB
		}
	case INVALIDLB:
		{
			errcode = INVAILDLBSRV
		}
	case INVALIDSRV:
		{
			errcode = INVAILDSRV
		}
	case INVALIDRMLB:
		{
			errcode = INVAILDRMLB
		}
	default:
		{
			errcode = NETWORKEXCEPT
		}
	}
	return errcode
}

// oneShortCall 无需维护状态的调用的实现
func (c *Caller) oneShortCall(apiVersion, service string, r *Req, tm time.Duration) (s *Res, errcode int32, e error) {
	logId := logSidGeneratorInst.GenerateSid("oneShortCall")
	c.cli.Log.Infow("oneShortCall enter", "logId", logId)

	//var span *utils.Span
	//if len(r.TraceID()) > 0 {
	//	span = utils.FromMeta(r.TraceID(), "", 0, service, utils.CliSpan)
	//} else {
	//	span = utils.NewSpan(utils.CliSpan)
	//}
	//
	//spnxt := span.Next(utils.CliSpan)
	//spnxt = spnxt.WithRpcCallType()
	//spnxt.WithName(r.Op()).Start()
	//spnxt.WithTag("call", "oneShortCall")
	//defer spnxt.Flush()
	//defer spnxt.End()

	tmpFlag, prx, conn, e := c.getPrxy(apiVersion, service, nil, withCallLogId(logId), withCallReq(r))
	if nil != e {
		c.cli.Log.Errorw("oneShortCall:getPrxy", "logId", logId, "error", e.Error(), "service", service)

		errcode = c.err2errCode(e)

		//spnxt.WithRetTag(strconv.Itoa(int(errcode)))
		//spnxt.WithFuncTag("getPrxy")
		//spnxt.WithErrorTag(e.Error())

		return nil, errcode, e
	}
	//r.SetTraceID(spnxt.Meta())
	if tmpFlag {
		r.SetParam(tmpKV.tmpK, tmpKV.tmpV)
	}
	callBase := time.Now()
	s, errcode, e = c.doCall(prx, r, tm)
	dur := time.Since(callBase).Nanoseconds()
	c.cli.Log.Infow("record doCall perf", "fn", "oneShort:qCall", "dur", dur)

	c.cli.updateLb(service, conn.Addr, s, errcode, e, dur)

	if nil != e {
		c.cli.Log.Errorw("oneShortCall:doCall",
			"logId", logId, "error", e.Error(), "service", service, "addr", conn.Addr)

		//spnxt.WithFuncTag("doCall")
		//spnxt.WithErrorTag(e.Error())
		//spnxt.WithTag("addr", conn.Addr)

	}
	//spnxt.WithRetTag(strconv.Itoa(int(errcode)))

	// 进入重试,换机器
	if nil != e && c.retry > 0 {
		return c.retryCall(apiVersion, service, r, tm, nil, withCallLogId(logId), withCallRetry(true))
	}

	return s, errcode, e
}

// doCall 正式发起请求的调用
func (c *Caller) doCall(prx utils.XsfCallClient, r *Req, tm time.Duration) (s *Res, errcode int32, e error) {
	ctx, cancel := context.WithTimeout(context.Background(), tm)
	defer cancel()

	s = utils.NewRes()
	res, e := prx.Call(ctx, r.Req())
	if nil != e {
		return nil, NETWORKEXCEPT, e
	}
	if len(res.ErrorInfo) > 0 {
		e = errors.New(res.ErrorInfo)
	}
	s.SetRes(res)
	return s, res.Code, e
}

// getPrxy 根据服务名获取本地调用代理
func (c *Caller) getPrxy(
	apiVersion string,
	service string,
	sp *Span,
	opt ...callOpt) (bool, utils.XsfCallClient, *SFConn, error) {
	callCfgInst := &callCfg{}
	for _, optItem := range opt {
		optItem(callCfgInst)
	}
	hashKey := c.isHash(callCfgInst)

	lbp := new(LBParams)
	if len(c.lbname) != 0 { //启用了远端LB模式
		lbp.WithLogId(callCfgInst.logId)
		lbp.WithHashKey(hashKey)
		lbp.WithVersion(apiVersion)
		lbp.WithName(c.lbname)
		lbp.WithSvc(service)
		lbp.WithCatgory(c.catgory)
		lbp.WithExtend(c.lbExt)
		lbp.WithLog(c.cli.Log)
		lbp.WithTracer(sp)
		lbp.WithNBest(int(c.retry))
		lbp.WithDirectEngIp(func() string {
			if callCfgInst.req != nil {
				directEngIp, _ := callCfgInst.req.GetParam(DIRECTENGIP)
				return directEngIp
			}
			return ""
		}())
	} else {
		lbp.WithLogId(callCfgInst.logId)
		lbp.WithHashKey(hashKey)
		lbp.WithVersion(apiVersion)
		lbp.WithName(service)
		lbp.WithSvc(service)
		lbp.WithLog(c.cli.Log)
		lbp.WithLocalIp(c.cli.cfg.GetLocalIp())
		lbp.WithRetry(callCfgInst.retry)
		lbp.WithNBest(int(c.retry + 1))
		lbp.WithPeerIp(func() string {
			if callCfgInst.req != nil {
				peerIp, _ := callCfgInst.req.GetPeerIp()
				return peerIp
			}
			return ""
		}())
	}

	tmpFlag, conn, e := c.cli.GetConn(lbp)
	//fmt.Println("getPrxy conn:", conn)
	if e != nil {
		c.cli.Log.Infow("c.cli.GetConn(lbp) failed",
			"logId", callCfgInst.logId)
		return tmpFlag, nil, nil, e
	}
	return tmpFlag, utils.NewXsfCallClient(conn.Conn), conn, nil
}

func (c *Caller) isHash(callCfgInst *callCfg) string {
	var hashKey string
	var hashKeyOk bool
	if !c.cli.lb.isHash() {
		return hashKey
	}

	if callCfgInst.retry {
		c.stringAdd()
		hashKey, hashKeyOk = c.hashKey.Load().(string)
	} else {
		hashKey, hashKeyOk = c.hashKey.Load().(string)
	}

	if !hashKeyOk {
		//缺省填充ip+port
		ipStr := c.cli.cfg.GetSvcIp()
		if ipStr == "" {
			ipStr = c.cli.cfg.GetLocalIp()
		}
		portStr := c.cli.cfg.GetSvcPort()
		c.cli.Log.Infow("set local ip+port for default",
			"ipStr", ipStr, "portStr", portStr, "logId", callCfgInst.logId)
		c.hashKey.Store(ipStr + portStr)
	}
	return hashKey
}

func (c *Caller) getSrv(service string) (addr []string, addrErr error) {

	lbp := new(LBParams)
	if 0 != len(c.lbname) { //启用了远端LB模式
		lbp.WithName(c.lbname)
		lbp.WithSvc(service)
		lbp.WithVersion(c.apiVersion)
		lbp.WithCatgory(c.catgory)
		lbp.WithExtend(c.lbExt)
		lbp.WithLog(c.cli.Log)
		lbp.WithNBest(int(c.retry + 1))
	} else {
		lbp.WithName(service)
		lbp.WithVersion(c.apiVersion)
		lbp.WithSvc(service)
		lbp.WithLog(c.cli.Log)
		lbp.WithNBest(int(c.retry + 1))
	}

	addr, _, addrErr = c.cli.lb.lbi.Find(lbp)
	return
}

// getPrxy 根据服务名获取服务地址列表
func (c *Caller) GetSrv(service string) (addr []string, addrErr error) {
	startIx, addrErr := c.getSrv(service)
	if nil != addrErr {
		return
	}
	addr = append(addr, startIx[0])

	/*
		默认线上单个服务最多10个lb
	*/
	for i := 0; i < 10; i++ {
		addrTmp, addrErr := c.getSrv(service)
		if nil == addrErr && addrTmp[0] == startIx[0] {
			break
		}
		if nil != addrErr || len(addrTmp) > 0 {
			addr = append(addr, addrTmp[0])
		}
	}
	return
}

// getConn 按照服务名获取连接相关属性
func (c *Caller) getConn(service string, sp *Span) (bool, *SFConn, error) {
	// todo: lb 模式需要根据lbmode来判断
	lbp := new(LBParams)
	if 0 != len(c.lbname) { //启用了远端LB模式
		lbp.WithName(c.lbname)
		lbp.WithSvc(service)
		lbp.WithCatgory(c.catgory)
		lbp.WithExtend(c.lbExt)
		lbp.WithLog(c.cli.Log)
		lbp.WithTracer(sp)
	} else {
		lbp.WithName(service)
		lbp.WithSvc(service)
		lbp.WithLog(c.cli.Log)
		lbp.WithTracer(sp)

	}

	tmpFlag, conn, e := c.cli.GetConn(lbp)
	//fmt.Println("getPrxy conn:", conn)
	if nil != e {
		return tmpFlag, nil, e
	}
	return tmpFlag, conn, nil
}

// getPrxyWithHandle 根据服务handle 获取调用代理
func (c *Caller) getPrxyWithHandle(handle string) (utils.XsfCallClient, error) {
	conn, e := c.cli.GetConnWithHandle(handle)
	if nil != e {
		return nil, e
	}
	return utils.NewXsfCallClient(conn.Conn), nil
}

// getPrxyWithAddr 根据指定地址获取调用代理
func (c *Caller) getPrxyWithAddr(addr string, logId string) (utils.XsfCallClient, error) {
	conn, e := c.cli.GetConnWithAddr(addr, logId)
	if nil != e {
		c.cli.Log.Errorw("CallWithAddr:GetConnWithAddr",
			"error", e.Error(), "addr", addr, "logId", logId)
		return nil, e
	}
	return utils.NewXsfCallClient(conn.Conn), nil
}
func (c *Caller) stringAdd() {

	var out string

	inStr, inStrOk := c.hashKey.Load().(string)
	if !inStrOk {
		return
	}
	inTmp := []byte(inStr)
	cf := false
	for k, v := range inTmp {
		vt := uint8(v)
		if (vt >= 48 && vt < 57) ||
			(vt >= 65 && vt < 90) ||
			(vt >= 97 && vt < 122) {
			inTmp[k] = v + 1
			cf = true
			break
		}
	}
	if !cf {
		out = inStr + strconv.Itoa(1)
	} else {
		out = string(inTmp)
	}

	c.hashKey.Store(out)
}

// retryCall
func (c *Caller) retryCall(
	apiVersion string,
	service string,
	r *Req,
	tm time.Duration,
//sp *Span,
	opt ...callOpt) (s *Res, errcode int32, e error) {

	callCfgInst := &callCfg{}
	for _, optItem := range opt {
		optItem(callCfgInst)
	}

	// 重试调用
	var prx utils.XsfCallClient
	var conn *SFConn

	for i := 0; i < c.retry; i++ {

		//spnxt := sp.Next(utils.CliSpan)
		//spnxt.WithName(r.Op()).Start()
		//spnxt = spnxt.WithRpcCallType()
		var tmpFlag bool
		tmpFlag, prx, conn, e = c.getPrxy(apiVersion, service, nil, opt...)
		if e != nil {
			//if conn != nil {
			//	spnxt.WithTag("addr", conn.Addr)
			//}

			errcode = c.err2errCode(e)

			//spnxt.WithFuncTag("getPrxy")
			//spnxt.WithErrorTag(e.Error())
			//spnxt.WithRetTag(strconv.Itoa(int(errcode)))
			//spnxt.End()
			//spnxt.Flush()
			continue
		}
		//spnxt.WithTag("Call", "retryCall")
		//spnxt.WithFuncTag("doCall")
		//
		//r.SetTraceID(spnxt.Meta())
		if tmpFlag {
			r.SetParam(tmpKV.tmpK, tmpKV.tmpV)
		}
		callBase := time.Now()
		s, errcode, e = c.doCall(prx, r, tm)
		c.cli.Log.Infow("record doCall perf",
			"fn", "retryCall", "dur", time.Since(callBase).Nanoseconds())

		if nil != e {
			c.cli.Log.Errorw("retryCall",
				"logId", callCfgInst.logId, "error", e.Error(), "service", service, "addr:", conn.Addr)

			//spnxt.WithErrorTag(e.Error())
			//spnxt.WithTag("addr", conn.Addr)
			//spnxt.WithRetTag(strconv.Itoa(int(errcode)))
			//
			//spnxt.End()
			//spnxt.Flush()

		} else {
			//spnxt.WithRetTag("0")
			//spnxt.End()
			//spnxt.Flush()
			return s, errcode, e
			//break  //成功退出
		}

	}
	return s, errcode, e
}

func (c *Caller) GetHashAddr(hashKey string, svc string) (string, error) {
	lbp := new(LBParams)
	lbp.WithHashKey(hashKey)
	lbp.WithVersion(c.apiVersion)
	lbp.WithSvc(svc)
	lbp.WithLog(c.cli.Log)

	addr, _, err := c.cli.lb.lbi.Find(lbp)
	if nil == err && len(addr) >= 1 {
		return addr[0], err
	}
	return "", err
}
