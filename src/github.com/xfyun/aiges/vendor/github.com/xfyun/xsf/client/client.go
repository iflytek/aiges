/*
* @file	client.go
* @brief  客户端对象操作接口
*
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */

package xsf

import (
	"fmt"
	"github.com/xfyun/xsf/utils"
	"github.com/xfyun/uuid"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Client 客户端对象结构
type Client struct {
	// Log 日志句柄
	Log *Logger

	// lb 负载均衡
	lb *loadBalance

	// cfg 配置句柄
	cfg *utils.Configure

	// fm  配置&服务中心操作句柄
	fm *utils.FindManger

	// name 客户端名
	name string

	// timeout 处理超时，暂时不生效
	timeout int64

	// conTimeout 链接超时
	conTimeout int64

	circuit *Circuit
}

// NewClient 创建客户端请求对象
//
// 参数说明：
//
// @cname: 客户端名称，这个参数在读取配置时会用作定位配置的selection。
//
// @mode: 配置文件支持的模式，详见CfgMode
//
// @o...: 配置相关的其他属性,详见 CfgOpt
func NewClient(cname string, mode utils.CfgMode, o ...utils.CfgOpt) (*Client, error) {
	c := new(Client)
	if len(cname) == 0 {
		c.name = CFGSECTION
	}
	c.name = cname

	var co utils.CfgOption
	co.SetDef(
		defaultCacheService,
		defaultCacheConfig,
		defaultCachePath,
		fmt.Sprintf("%s", uuid.Must(uuid.NewV4())))
	for _, opt := range o {
		opt(&co)
	}

	// 初始化服务发现
	e := c.newFinder(mode, &co)
	if e != nil {
		return nil, e
	}

	// 初始化配置
	e = c.SetCfgWithOption(mode, &co)
	if e != nil {
		return nil, e
	}

	/*
		获取ip
	*/
	{
		host, _ := c.cfg.GetString(cname, CLIHOST)
		netcard, _ := c.cfg.GetString(cname, CLINETCARD)
		ipStr, e := utils.Host2Ip(host, netcard)
		if e != nil {
			return nil, e
		}
		co.SetLocalIp(ipStr)
	}

	e = c.newLog()
	if e != nil {
		//c.Log.Errorf("NewClient newLog failed: %v", e)
		return nil, e
	}

	c.Log.Errorw("print opt", "opt", co.String())
	e = c.newLB()
	if e != nil {
		c.Log.Errorf("NewClient newLB failed: %v", e)
		return nil, e
	}
	e = c.newTracer(co.SvcIp, co.SvcPort)
	if e != nil {
		c.Log.Errorf("NewClient newTracer failed: %v", e)
		return nil, e
	}

	c.Log.Infof("NewClient success")

	c.circuit = newCircuit()

	return c, nil
}

// NewClientWithCfg 创建客户端请求对象
//
// 参数说明：
//
// @cname: 客户端名称，这个参数在读取配置时会用作定位配置的selection。
//
// @cfg:外部已经初始化完成的配置操作句柄
//
// @o...: 配置相关的其他属性,详见 CfgOpt
func NewClientWithCfg(cname string, cfg *utils.Configure, o ...utils.CfgOpt) (*Client, error) {
	if cfg == nil {
		return nil, ECFGISNIL
	}
	c := new(Client)
	if len(cname) == 0 {
		c.name = CFGSECTION
	}
	c.name = cname

	var co utils.CfgOption
	for _, opt := range o {
		opt(&co)
	}

	if cfg.Option() != nil && cfg.Option().FindManger() != nil {
		c.fm = cfg.Option().FindManger()
	}

	e := c.SetCfg(cfg, &co)
	if e != nil {
		return nil, e
	}

	e = c.newLog()
	if e != nil {
		c.Log.Errorf("NewClientWithCfg newLog failed: %v", e)
		return nil, e
	}

	e = c.newLB()
	if e != nil {
		c.Log.Errorf("NewClientWithCfg newLB failed: %v", e)
		return nil, e
	}
	e = c.newTracer(co.SvcIp, co.SvcPort)
	if e != nil {
		c.Log.Errorf("NewClientWithCfg newTracer failed: %v", e)
		return nil, e
	}
	return c, nil
}

func DestroyClient(c *Client) {

	utils.DestroyFinder(c.fm)
	// 停止刷日志
	utils.StopLocalLog(c.Log)
	return
}

// GetConn 根据LB相关参数获取一个可用用于发起连接的对象
func (c *Client) GetConn(lbp *LBParams) (bool, *SFConn, error) {
	return c.lb.find(lbp)
}

// GetConnWithHandle 根据服务句柄（handle）相关参数获取一个可用用于发起连接的对象
func (c *Client) GetConnWithHandle(hanlde string) (*SFConn, error) {
	return c.lb.findWithHandle(hanlde)
}

// GetConnWithAddr 根据指定地址（ip:port）相关参数获取一个可用用于发起连接的对象
func (c *Client) GetConnWithAddr(addr string, logId string) (*SFConn, error) {
	return c.lb.findConn(addr, logId)
}

// SetCfgWithOption 根据相关参数设置相关的配置解析器
func (c *Client) SetCfgWithOption(mode utils.CfgMode, co *utils.CfgOption) error {
	cfg, e := utils.NewCfg(mode, co)
	if e != nil {
		return e
	}
	c.cfg = cfg
	return nil
}

// SetCfg 设置相关的配置解析器
func (c *Client) SetCfg(cfg *utils.Configure, co *utils.CfgOption) error {
	var e error
	//cfg.r = c.fm
	c.cfg, e = cfg.GenNewCfg(co)
	return e
}

// Cfg 获取配置操作句柄
func (c *Client) Cfg() *utils.Configure {
	return c.cfg
}

func (c *Client) updateLb(svc string, target string, s *Res, errcode int32, e error, dur int64) {
	dbgLoggerStd.Printf("target:%v,errcode:%v,e:%v,dur:%v\n", target, errcode, e, dur)
	var vCpu int64
	if !c.IsNil(s) {
		vCpuStr, _ := s.GetParam("vcpu")
		if 0 == len(vCpuStr) {
			vCpu = 0
		} else {
			v, _ := strconv.Atoi(vCpuStr)
			vCpu = int64(v)
		}
	}
	c.lb.update(svc, target, s, errcode, e, dur, vCpu)
}
func (c *Client) IsNil(obj interface{}) bool {
	type eface struct {
		rtype unsafe.Pointer
		data  unsafe.Pointer
	}
	if obj == nil {
		return true
	}
	return (*eface)(unsafe.Pointer(&obj)).data == nil
}

// newLB 创建负载均衡操作对象
func (c *Client) newLB() error {
	if c != nil && c.cfg != nil {
		// 读取配置中连接超时
		tm, e := c.cfg.GetInt(c.name, CFGCONTM)
		if e != nil {
			tm = CFGDEFAULTTIMEOUT
		}
		// 读取配置中连接超时
		lc, e := c.cfg.GetInt(c.name, CFGCONLC)
		if e != nil {
			lc = CFGCONDEFLC
		}
		// 读取重试次数
		r, e := c.cfg.GetInt(c.name, CFGRTY)
		if e != nil {
			r = CFGDEFRTY
		}
		// 读取单地址连接池大小
		p, e := c.cfg.GetInt(c.name, CFGCONPOOL)
		if e != nil {
			p = CFGDEFCONPOLL
		}
		rbuf, e := c.cfg.GetInt(c.name, CFGCONRBUF)
		if e != nil {
			rbuf = CFGDEFRBUF
		}
		maxReveive, e := c.cfg.GetInt(c.name, CFGCONMAXRECEIVE)
		if e != nil {
			maxReveive = DEFCFGCONMAXRECEIVE
		} else {
			maxReveive = maxReveive * 1024 * 1024
		}
		wbuf, e := c.cfg.GetInt(c.name, CFGCONWBUF)
		if e != nil {
			wbuf = CFGDEFWBUF
		}

		// 负载均衡的模式
		lbm, e := c.cfg.GetInt(c.name, CFGLBMODE)
		if e != nil {
			lbm = CFGDEFLBMODE
		}

		lbrty, e := c.cfg.GetInt(c.name, CFGLBRTY)
		if e != nil {
			lbrty = CFGDEFLBRTY
		}

		lbtm, e := c.cfg.GetInt(c.name, CFGLBTIMEOUT)
		if e != nil {
			lbtm = CFGDEFLBTIMEOUT
		}

		keepaliveTime := CFGDEFKEEPALIVE
		keepaliveTimeInt64, e := c.cfg.GetInt64(c.name, CFGKEEPALIVE)
		if e == nil {
			keepaliveTime = time.Millisecond * time.Duration(keepaliveTimeInt64)
		}
		keepaliveTimeout := CFGDEFKEEPALIVETIMEOUT
		keepaliveTimeoutInt64, e := c.cfg.GetInt64(c.name, CFGKEEPALIVETIMEOUT)
		if e == nil {
			keepaliveTimeout = time.Millisecond * time.Duration(keepaliveTimeoutInt64)
		}

		timePerSlice, e := c.cfg.GetInt(c.name, CFGWINTIMEPERSLICE)
		if e != nil {
			timePerSlice = CFGDEFWINTIMEPERSLICE
		}

		winsize, e := c.cfg.GetInt(c.name, CFGWINSIZE)
		if e != nil {
			winsize = CFGDEFWINSIZE
		}

		var probabilityMatrix []int
		probabilityStr, e := c.cfg.GetString(c.name, CFGPROBABILITY)
		if e != nil {
			probabilityMatrix = CFGDEFPROBABILITY
		} else {
			probabilityMatrixStr := strings.Split(probabilityStr, CFGPROBABILITYSEPARATOR)
			for _, probability := range probabilityMatrixStr {
				v, e := strconv.Atoi(probability)
				if e != nil {
					return fmt.Errorf("can't convert %v to number", probability)
				}
				probabilityMatrix = append(probabilityMatrix, v)
			}
		}
		threshold, e := c.cfg.GetInt(c.name, CFGTHRESHOLD)
		if e != nil {
			threshold = CFGDEFTHRESHOLD
		}
		pingInterval := CFGDEFPINGINTERVAL
		pingIntervalInt, e := c.cfg.GetInt(c.name, CFGPINGINTERVAL)
		if e == nil {
			pingInterval = time.Duration(pingIntervalInt) * time.Millisecond
		}
		c.lb = newLB(
			c.Log,
			uint(r),
			LBMode(lbm),
			WithKeepaliveTime(keepaliveTime),
			WithKeepaliveTimeout(keepaliveTimeout),
			WithConLifeCycle(lc),
			WithConTimeOut(tm),
			WithConMax(p),
			WithConFindManger(c.fm),
			WithLBRetry(lbrty),
			WithLBTimeOut(lbtm),
			WithReadBufSize(rbuf),
			WithWriteBufSize(wbuf),
			WithConMaxMsg(maxReveive),
			WithProbabilityMatrix(probabilityMatrix),
			WithWindowMeta(time.Duration(timePerSlice), int64(winsize)),
			WithThreshold(threshold),
			WithPingInterval(pingInterval),
			WithClient(c),
		)

		c.Log.Infof("read cfg: %d, %d, %d, %d, %d ", tm, lc, r, p, lbm)
		if c.lb == nil {
			return ELBMODE
		}
		// 负载均衡的模式
		addr, e := c.cfg.GetString(c.name, CFGADDRS)
		if e == nil {
			return setTeatAddr(addr)
		}
		c.Log.Infof("read cfg:tm:%d, lc:%d, r:%d, p:%d, lbm:%d, addr:%s, e:%v", tm, lc, r, p, lbm, addr, e)
		return nil
	}
	return ECFGISNIL
}

// newLog 创建日志操作对象
func (c *Client) newLog() error {
	lev, e := c.cfg.GetString(CFGLOGSECTON, CFGLOGLEVEL)
	if e != nil {
		lev = CFGLOGFILE
	}

	fn, e := c.cfg.GetString(CFGLOGSECTON, CFGLOGFILE)
	if e != nil {
		fn = CFGLODEFGFILE
	}
	size, e := c.cfg.GetInt(CFGLOGSECTON, CFGLOGSIZE)
	if e != nil {
		size = CFGLOGDEFSIZSE
	}
	count, e := c.cfg.GetInt(CFGLOGSECTON, CFGLOGCOUNT)
	if e != nil {
		count = CFGLOGDEFCOUNT
	}

	caller := CFGLOGDEFCALLER
	callerInt, e := c.cfg.GetInt(CFGLOGSECTON, CFGLOGCALLER)
	if e == nil && callerInt == 1 {
		caller = true
	}

	die, e := c.cfg.GetInt(CFGLOGSECTON, CFGLOGDIED)
	if e != nil {
		die = CFGLOGDEFDIED
	}

	as := true

	async, e := c.cfg.GetInt(CFGLOGSECTON, CFGLOGASYNC)
	if e == nil && async == 0 {
		as = false
	}

	cache, e := c.cfg.GetInt(CFGLOGSECTON, CFGLOGCACHE)
	if e != nil {
		die = CFGLOGDEFCACHE
	}

	batch, e := c.cfg.GetInt(CFGLOGSECTON, CFGLOGBATCH)
	if e != nil {
		die = CFGLOGDEFBATCH
	}

	wash, e := c.cfg.GetInt(CFGLOGSECTON, CFGLOGWASH)
	if e != nil {
		die = CFGLOGDEFWASH
	}

	//初始化日志
	//	c.Log, e = utils.NewLocalLog(lev,fn,size,count, die)
	c.Log, e = utils.NewLocalLog(utils.SetCaller(caller), utils.SetLevel(lev), utils.SetFileName(fn),
		utils.SetMaxSize(size), utils.SetMaxBackups(count), utils.SetMaxAge(die),
		utils.SetAsync(as), utils.SetCacheMaxCount(cache), utils.SetBatchSize(batch), utils.SetWash(wash))

	return e
}

// newFinder 创建配置中心操作实例
func (c *Client) newFinder(mode CfgMode, co *CfgOption) error {
	// 配置中心
	switch mode {
	case utils.Custom:
		fallthrough
	case utils.Centre:
		{
			var e error
			utils.WithCfgLog(c.Log)(co)
			c.fm, e = utils.NewFinder(co)
			if e == nil {
				utils.WithCfgReader(c.fm)(co)

			}
			return e
		}
	}
	return nil
}

// newTracer 创建Tracer对象
func (c *Client) newTracer(svcIp string, svcPort int32) error {

	addr, e := c.cfg.GetString(CFGTRACESECTION, CFGTRACEADDR)
	if e != nil {
		addr = CFGDEFTRACEADDR
	}

	port, e := c.cfg.GetString(CFGTRACESECTION, CFGTRACEPORT)
	if e != nil {
		port = CFGDEFTRACEPORT
	}
	enbale := true
	able, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACEENABLE)
	if e == nil && able == 0 {
		enbale = false
	}

	enableDump := false
	dump, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACEDUMP)
	//fmt.Println(e, dump)
	if e == nil && dump != 0 {
		enableDump = true
	}

	bakend, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACEBAKEND)
	if e != nil {
		bakend = CFGDEFTRACEBAKEND
	}

	enableDeliver := true
	deliver, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACEDELIVER)
	if e == nil && deliver == 0 {
		enableDeliver = false
	}

	enableSpill := true
	espill, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACESPILLABLE)
	if e == nil && espill == 0 {
		enableSpill = false
	}

	spill, e := c.cfg.GetString(CFGTRACESECTION, CFGTRACESPILL)
	if e != nil {
		spill = CFGDEFTRACESPILL
	}

	buffer, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACEBUFFER)
	if e != nil {
		buffer = CFGDEFTRACEBUFFER
	}

	batch, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACEBATCH)
	if e != nil {
		batch = CFGDEFTRACEBATCH
	}

	linger, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACELINGER)
	if e != nil {
		linger = CFGDEFTRACELINGER
	}

	watch := CFGLOGDEFWATCH
	watchInt, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACEWATCH)
	if e == nil {
		if watchInt == 1 {
			watch = true
		}
	}

	watchport, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACEWATCHPORT)
	if e != nil {
		watchport = CFGTRACEDEFWATCHPORT
	}

	spillsize, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACESPILLSIZE)
	if e != nil {
		spillsize = CFGTRACEDEFSPILLSIZE
	}

	loadts, e := c.cfg.GetInt(CFGTRACESECTION, CFGTRACELOADTS)
	if e != nil {
		loadts = CFGTRACEDEFLOADTS
	}

	bcluster, e := c.cfg.GetString(CFGTRACESECTION, CFGTRACEBCLUSTER)
	if e != nil {
		bcluster = CFGTRACEDEFCLUSTER
	}
	idc, e := c.cfg.GetString(CFGTRACESECTION, CFGTRACEIDC)
	if e != nil {
		idc = CFGTRACEDEFIDC
	}

	if able != CFGTRACEUNCHANGE {
		utils.AbleTrace(enbale)
		if enbale {
			return utils.InitTracer(addr, port,
				utils.WithLowLoadSleepTs(loadts),
				utils.WithMaxSpillContentSize(spillsize),
				utils.WithWatchPort(watchport),
				utils.WithWatch(watch),
				utils.WithDump(enableDump),
				utils.WithDeliver(enableDeliver),
				utils.WithBackend(bakend),
				utils.WithTraceLogger(c.Log),
				utils.WithSvcName(c.name),
				utils.WithTraceSpill(spill),
				utils.WithBufferSize(buffer),
				utils.WithBatchSize(batch),
				utils.WithLinger(linger),
				utils.WithSpillAble(enableSpill),
				utils.WithSvcBCluster(bcluster),
				utils.WithSvcIDC(idc),
				utils.WithSvcIp(svcIp),
				utils.WithSvcPort(svcPort))
		}
	}

	//fmt.Println(loadts)
	//fmt.Println(spillsize)
	//fmt.Println(watchport)
	//fmt.Println(watch)

	return nil
}
