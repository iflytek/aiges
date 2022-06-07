/*
* @file	loadbalance.go
* @brief  负载均衡封装层
*         所有的负载策略均封装在这个对象的操作
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */

package xsf

import (
	"github.com/xfyun/xsf/utils"
	"google.golang.org/grpc"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// conOption 连接相关属性设置
type conOption struct {
	// retry 重试次数
	retry int

	// lbtimeout 超时时间,单位ms
	lbtimeout int

	// lbretry 重试次数
	lbretry int

	// timeout 连接超时，单位ms
	timeout int

	// 连接读缓冲区大小，默认GRCMRS
	rbuf int

	// 连接写缓冲区大小，默认GRCMWS
	wbuf int

	// maxmsgsize 消息最大尺寸
	maxmsgsize int

	// max 最大连接数
	max int

	// lc 地址活跃有效期，超过这个时间则认为地址过期。单位ms
	lc int

	// fm 服务发现的操作句柄
	fm *FindManger

	//keepalive检查时间间隔
	keepaliveTime time.Duration

	//keepalive检查超时时间
	keepaliveTimeout time.Duration

	//窗口每个时间片的时长
	timePerSlice time.Duration
	//窗口长度
	winSize int64
	//轮盘赌概率矩阵
	probabilityMatrix []int

	//topK负载时请求特定的机器前的阈值
	threshold int

	//心跳间隔期
	pingInterval time.Duration

	//topK心跳用
	client *Client
}

// connOpt  属性设置函数返回值
type connOpt func(*conOption)

func WithKeepaliveTime(to time.Duration) connOpt {
	return func(o *conOption) {
		o.keepaliveTime = to
	}
}
func WithKeepaliveTimeout(to time.Duration) connOpt {
	return func(o *conOption) {
		o.keepaliveTimeout = to
	}
}

func WithProbabilityMatrix(probabilityMatrix []int) connOpt {
	return func(o *conOption) {
		o.probabilityMatrix = probabilityMatrix
	}
}

func WithWindowMeta(timePerSlice time.Duration, winSize int64) connOpt {
	return func(o *conOption) {
		o.timePerSlice = timePerSlice
		o.winSize = winSize
	}
}
func WithThreshold(threshold int) connOpt {
	return func(o *conOption) {
		o.threshold = threshold
	}
}
func WithPingInterval(pingInterval time.Duration) connOpt {
	return func(o *conOption) {
		o.pingInterval = pingInterval
	}
}
func WithClient(client *Client) connOpt {
	return func(o *conOption) {
		o.client = client
	}
}

// WithConTimeOut 设置连接超时
func WithConTimeOut(to int) connOpt {
	return func(o *conOption) {
		o.timeout = to
	}
}

// WithConMaxMsg 设置最大消息大小
func WithConMaxMsg(s int) connOpt {
	return func(o *conOption) {
		o.maxmsgsize = s
	}
}

// WithConMax 设置最大连接数
func WithConMax(s int) connOpt {
	return func(o *conOption) {
		o.max = s
	}
}

// WithConLifeCycle 设置连接的生命周期
func WithConLifeCycle(l int) connOpt {
	return func(o *conOption) {
		o.lc = l
	}
}

// WithConFindManger 设置配置中心操作句柄
func WithConFindManger(fm *FindManger) connOpt {
	return func(o *conOption) {
		o.fm = fm
	}
}

// WithLBTimeOut 请求LB的超时时间
func WithLBTimeOut(tm int) connOpt {
	return func(o *conOption) {
		o.lbtimeout = tm
	}
}

// WithLBRetry 设置LB重试次数
func WithLBRetry(r int) connOpt {
	return func(o *conOption) {
		if r > 0 {
			o.lbretry = r
		} else {
			o.lbretry = CFGDEFLBRTY
		}
	}
}

// WithReadBufSize 设置链接读缓冲区
func WithReadBufSize(r int) connOpt {
	return func(o *conOption) {
		o.rbuf = r
	}
}

// WithWriteBufSize 设置链接读缓冲区
func WithWriteBufSize(r int) connOpt {
	return func(o *conOption) {
		o.wbuf = r
	}
}

// loadBalance 负载均衡对外的操作结构
type loadBalance struct {
	// conns 连接池
	conns *connPool

	// lbi LB实现策略
	lbi LBI

	log  *utils.Logger
	r    uint
	mode LBMode
	c    *cleaner
}

// SFConn 连接对象，lazy方式
type SFConn struct {
	// Conn 可用连接，为Addr的第一个地址的连接
	Conn *grpc.ClientConn

	//Addr 可用地址池
	Addr string
}

// newLB 创建loadBalance 操作对象
func newLB(log *utils.Logger, r uint, m LBMode, o ...connOpt) *loadBalance {

	lb := new(loadBalance)
	lb.conns = newConnPool(o...) // 连接属性
	lb.log = log

	lb.r = r + 1
	lb.mode = m
	// 设置选项
	var co conOption
	co.retry = 2
	co.lbtimeout = 500

	for _, opt := range o {
		opt(&co)
	}
	lb.log.Infof("NewLB mode: %d", m)
	dbgLoggerStd.Printf("lb mode:%v\n", m)
	switch m {
	case RoundRobin:
		lb.lbi = newRRLB(&co)
	case Hash:
		lb.lbi = newHashLB(&co)
	case Remote:
		lb.lbi = newRemoteLB(&co, newRRLB(&co), lb.conns)
	case ConHash:
		lb.lbi = newConsistencyHashLB(&co)
	case TopKHash:
		lb.lbi = newTopKLB(&co)
	default:
		return nil
	}
	if func(obj interface{}) bool {
		type eface struct {
			rtype unsafe.Pointer
			data  unsafe.Pointer
		}
		if obj == nil {
			return true
		}
		return (*eface)(unsafe.Pointer(&obj)).data == nil
	}(lb.lbi) {
		return nil
	}
	lb.c = newCleaner(time.Duration(lb.conns.o.lc)*time.Millisecond, lb)
	return lb
}
func (lb *loadBalance) isHash() bool {
	if lb.mode == Hash || lb.mode == ConHash {
		return true
	}
	return false
}
func (lb *loadBalance) update(svc string, target string, s *Res, errcode int32, e error, dur int64, vCpu int64) {
	topKLBInst, ok := lb.lbi.(*topKLB)
	if !ok {
		return
	}
	topKLBInst.updateData(svc, target, s, errcode, e, dur, vCpu)
}

//是否需要本地优先
func (lb *loadBalance) isLocalPriority() bool {
	_, ok := lb.lbi.(*topKLB)
	if ok {
		return false
	}
	return true
}

// findWithID 找到一个合适的可用的连接。id暂时未实现
func (lb *loadBalance) findWithID(lbp *LBParams, id string) (tmp bool, conn *SFConn, e error) {
	// todo: 增加地址的重试
	var addr []string
	//var i uint
	//for  ;i < lb.r ; i ++  {
	var allAddr []string
	addr, allAddr, e = lb.lbi.Find(lbp)
	if len(addr) > 0 {

		lb.log.Infow(
			"loadBalance.FindWithID",
			"logId", lbp.logId, "name", lbp.name, "addr", addr)
		{
			//2019-05-27 16:38:50
			if !lbp.retry && lb.isLocalPriority() && len(lbp.localIp) != 0 {
				lb.log.Infow(
					"loadBalance.FindWithID,about to use local address first",
					"logId", lbp.logId, "name", lbp.name, "nbest", lbp.nbest, "localIp", lbp.localIp, "retry", lbp.retry)
				for _, val := range allAddr {
					if strings.Contains(val, lbp.localIp) {
						lb.log.Infow(
							"loadBalance.FindWithID,found an peerAddr that is the same as the local address",
							"logId", lbp.logId, "name", lbp.name, "nbest", lbp.nbest, "localIp", val)
						var c *grpc.ClientConn
						c, e = lb.conns.get(val, id)
						if e == nil {
							conn = new(SFConn)
							conn.Conn = c
							conn.Addr = val
							lb.c.update(val)
							return false, conn, e
						} else {
							lb.log.Warnw("lb.conns.Get failed", "addr", val, "logId", lbp.logId, "err", e.Error())
						}
					}
				}
			} else {
				lb.log.Infow("loadBalance.FindWithID,ignore local address first")
			}
		}
		/////////////////////////
		//2018/10/11 11:36
		addrIx := 0
		var i uint
		for i = 0; i < lb.r; i++ {
			lb.log.Infof("loadBalance.findWithID id: %s, addr:%s,lb.r:%v", id, addr[addrIx%len(addr)], lb.r)
			var c *grpc.ClientConn
			addrTmp, tmpFlag := lb.tmpNode(addr[addrIx%len(addr)])
			c, e = lb.conns.get(addrTmp, id)
			if e == nil {
				conn = new(SFConn)
				conn.Conn = c
				conn.Addr = addrTmp
				lb.c.update(addrTmp)
				return tmpFlag, conn, e
			} else {
				lb.log.Warnw("lb.conns.Get failed", "addr", addr[addrIx%len(addr)], "logId", lbp.logId, "err", e.Error())
			}
			addrIx++
		}
	} else {
		lb.log.Warnw("loadBalance.FindWithID name failed", "logId", lbp.logId, "error", lbp.name, "error", e)
	}

	//经过重试后，仍未找到
	lb.log.Errorw("loadBalance.FindWithID can't find addr",
		"logId", lbp.logId, "name", lbp.name, "addr", addr, "allAddr", allAddr)
	if e != nil {
		return false, nil, e
	}
	return false, nil, EINVALIDADDR
}

// tmpNode标记处理
func (lb *loadBalance) tmpNode(addr string) (string, bool) {
	lb.log.Infow("loadBalance.tmpNode", "addr", addr)
	if strings.Contains(addr, tmpNode) {
		return strings.Replace(addr, tmpNode, "", -1), true
	} else {
		return addr, false
	}
}

// find 找到一个可用的连接
func (lb *loadBalance) find(lbp *LBParams) (bool, *SFConn, error) {
	return lb.findWithID(lbp, "")
}

// findWithHandle 根据handle找到一个可用的连接
// todo: 增加handle规则的注释
func (lb *loadBalance) findWithHandle(handle string) (*SFConn, error) {

	// 根据handle解析出ip、port
	ip, e := lb.handleToIP(handle)
	if e != nil {
		return nil, EBADHANDLE
	}
	port, e := strconv.ParseInt(handle[9:13], 16, 32)
	if e != nil {
		return nil, EBADHANDLE
	}
	addr := ip + ":" + strconv.Itoa(int(port))

	// 根据地址获取可用连接
	var i uint
	for i = 0; i < lb.r; i++ {
		lb.log.Infof("loadBalance.FindWithHandle name: %s, addr:%s", handle, addr)
		var c *grpc.ClientConn
		c, e = lb.conns.get(addr, "")
		if e == nil {
			conn := new(SFConn)
			conn.Conn = c
			conn.Addr = addr
			lb.c.update(addr)
			return conn, e
		}
	}

	return nil, EBADADDR
}

// findConn 获取一个可用连接
func (lb *loadBalance) findConn(addr string, logId string) (*SFConn, error) {
	baseTime := time.Now()
	c, e := lb.conns.get(addr, "")
	dur := time.Since(baseTime).Nanoseconds()
	if e == nil {
		conn := new(SFConn)
		conn.Conn = c
		conn.Addr = addr
		lb.c.update(addr)
		return conn, e
	}
	lb.log.Errorw("CallWithAddr:findConn",
		"error", e.Error(), "addr", addr, "logId", logId, "dur", dur)
	return nil, e
}

// remove 从LB中清除一个地址
func (lb *loadBalance) remove(addr string) {
	//todo 需要完善,清理LB中的addr地址
	//	lb.lbi.Remove(addr)
	lb.conns.remove(&addr)
}

// handleToIP 从handle中解析出ip
func (lb *loadBalance) handleToIP(h string) (string, error) {
	ip1, e := strconv.ParseInt(string(h[1:3]), 16, 32)
	if e != nil {
		return "", e
	}
	ip2, e := strconv.ParseInt(string(h[3:5]), 16, 32)
	if e != nil {
		return "", e
	}
	ip3, e := strconv.ParseInt(string(h[5:7]), 16, 32)
	if e != nil {
		return "", e
	}
	ip4, e := strconv.ParseInt(string(h[7:9]), 16, 32)
	if e != nil {
		return "", e
	}

	return strconv.Itoa(int(ip1)) + "." + strconv.Itoa(int(ip2)) + "." + strconv.Itoa(int(ip3)) + "." + strconv.Itoa(int(ip4)), nil
}
