package xsf

import (
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/client"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	CNAME = "lbv2"
)

type hermesAdapter struct {
	able   bool
	svc    string
	subsvc string
	addr   string
	uid    string
	total  int
	idle   int
	best   int

	lbname     string
	apiversion string
	finderTtl  time.Duration //更新本地地址的时间，通过访问服务发现实现
	backend    int           //上报的的协程数，缺省4
	timeout    time.Duration //上报的超时时间，缺省一秒

	bc BootConfig

	cli    *xsf.Client
	caller *xsf.Caller

	taskInChan chan func() //任务通道，用来传送上报任务
	//taskOutChan chan callRstItem

	hermesTask int
	svcIp      string //服务端监听ip，trace用
	svcPort    int32  //服务端监听端口，trace用
}
type HermesAdapterCfgOpt func(*hermesAdapter)

func WithHermesAdapterSvcIp(svcIp string) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.svcIp = svcIp
	}
}
func WithHermesAdapterTask(task int) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.hermesTask = task
	}
}
func WithHermesAdapterSvcPort(svcPort int32) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.svcPort = svcPort
	}
}
func WithHermesAdapterTimeout(timeout time.Duration) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.timeout = timeout
	}
}
func WithHermesAdapterLbName(lbname string) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.lbname = lbname
	}
}
func WithHermesAdapterLbApiVersion(apiversion string) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.apiversion = apiversion
	}
}
func WithHermesAdapterBackEnd(backend int) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.backend = backend
	}
}
func WithHermesAdapterAddr(addr string) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.addr = addr
	}
}
func WithHermesAdapterSubsvc(subsvc string) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.subsvc = subsvc
	}
}
func WithHermesAdapterSvc(svc string) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.svc = svc
	}
}

func WithHermesAdapterFinderTtl(finderTtl time.Duration) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.finderTtl = finderTtl
	}
}

func WithHermesAdapterBootConfig(bc BootConfig) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.bc = bc
	}
}

func (h *hermesAdapter) Init(opts ...HermesAdapterCfgOpt) (err error) {
	for _, o := range opts {
		o(h)
	}
	if !h.able {
		loggerStd.Printf("hermes not enable\n")
		return
	}
	h.cli, err = xsf.InitClient(
		CNAME,
		h.bc.CfgMode,
		utils.WithCfgCacheService(true),
		utils.WithCfgCacheConfig(true),
		utils.WithCfgCachePath("."),
		utils.WithCfgName(h.bc.CfgData.CfgName),
		utils.WithCfgURL(h.bc.CfgData.CompanionUrl),
		utils.WithCfgPrj(h.bc.CfgData.Project),
		utils.WithCfgGroup(h.bc.CfgData.Group),
		utils.WithCfgService(h.bc.CfgData.Service),
		utils.WithCfgVersion(h.bc.CfgData.Version),
		utils.WithCfgSvcIp(h.svcIp),
		utils.WithCfgSvcPort(h.svcPort))
	if nil != err {
		return fmt.Errorf("InitClient fail err:%v", err)
	}
	h.caller = xsf.NewCaller(h.cli)
	h.caller.WithApiVersion(h.apiversion)
	h.taskInChan = make(chan func(), h.hermesTask)

	/*
		消费task
	*/
	go h.writer()
	return
}

type callRstItem struct {
	s       *Res
	errcode int32
	e       error
	addr    string
}

func (h *hermesAdapter) writer() {
	if !h.able {
		return
	}
	h.cli.Log.Infow("about to start writer", "backend", h.backend)
	wg := sync.WaitGroup{}
	for ix := 0; ix < h.backend; ix++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range h.taskInChan {

				h.cli.Log.Debugw("receive task")

				task()

			}
		}()
	}
	wg.Wait()
	h.cli.Log.Infow("writer exiting")
}

var setServerOnce sync.Once

const reportDef = -1
const reportFailure = 0
const reportSuccess = 1

var reportFlag = reportDef //这部分后续优化吧,-1:not set,0:failure,1:success

func (h *hermesAdapter) setServer(svc, subsvc, addr string, getAuthInfo func() (int32, int32, int32), live string) error {
	if !h.able {
		return nil
	}
	//req := utils.NewReq()
	//req.SetParam("svc", svc)
	//req.SetParam("subsvc", subsvc)
	//req.SetParam("addr", addr)
	//req.SetParam("live", live)

	reportAddr, reportAddrErr := h.caller.GetSrv(h.lbname)
	if reportAddrErr != nil {
		h.cli.Log.Errorw("can't get valid report addrs", "reportAddr", strings.Join(reportAddr, ";"))
		return reportAddrErr
	} else {
		h.cli.Log.Infow("report addrs", "reportAddr", strings.Join(reportAddr, ";"))
	}

	setServerOnce.Do(func() {
		loggerStd.Printf("start reporting->target:%v\n", reportAddr)
	})

	for _, lbAddr := range reportAddr {
		lbAddrTmp := lbAddr
		task := func() {

			req := utils.NewReq()
			req.SetParam("svc", svc)
			req.SetParam("subsvc", subsvc)
			req.SetParam("addr", addr)
			req.SetParam("live", live)
			for k, v := range lbReportExtInst.getAll() {
				req.SetParam(k, v)
			}

			if getAuthInfo != nil {

				maxLic, idle, bestLic := getAuthInfo()

				req.SetParam("total", strconv.Itoa(int(maxLic)))
				req.SetParam("best", strconv.Itoa(int(bestLic)))
				req.SetParam("idle", strconv.Itoa(int(idle)))
			} else {

				req.SetParam("total", "0")
				req.SetParam("best", "0")
				req.SetParam("idle", "0")
			}

			_, errcode, e := h.caller.CallWithAddr(h.lbname, xsf.LBOPSET, lbAddrTmp, req, time.Second)
			if errcode != 0 || e != nil {
				reportFlag = reportFailure
				h.cli.Log.Errorw("fn:setServer h.caller.CallWithAddr", "errcode", errcode, "err", e, "addr", lbAddrTmp)
			} else {
				reportFlag = reportSuccess
			}
		}

		select {
		case h.taskInChan <- task:
		default:
			{
				h.cli.Log.Warnw("taskInChan overflow")
			}
		}
	}

	h.cli.Log.Debugw("create report ctx", "timeout(ns)", int(h.timeout))

	return nil
}

func (h *hermesAdapter) report(getAuthInfo func() (int32, int32, int32)) error {
	if !h.able {
		return nil
	}
	return h.setServer(h.svc, h.subsvc, h.addr, getAuthInfo, "1")
}
func (h *hermesAdapter) offline() error {
	if !h.able {
		return nil
	}
	return h.setServer(h.svc, h.subsvc, h.addr, nil, "0")
}

func (h *hermesAdapter) getServer(uid, svc, subsvc, nbest, all string) (s *Res, errcode int32, e error) {
	if !h.able {
		return
	}
	req := utils.NewReq()
	req.SetParam("uid", uid)
	req.SetParam("svc", svc)
	req.SetParam("subsvc", subsvc)
	req.SetParam("nbest", nbest)
	req.SetParam("all", all)
	return h.caller.Call(h.lbname, "getServer", req, time.Second)
}
func (h *hermesAdapter) request(uid, svc, subsvc, nbest, all string) (s *Res, errcode int32, e error) {
	if !h.able {
		return
	}
	return h.getServer(uid, svc, subsvc, nbest, "0")
}
func (h *hermesAdapter) requestAll(uid, svc, subsvc, nbest, all string) (s *Res, errcode int32, e error) {
	if !h.able {
		return
	}
	return h.getServer(uid, svc, subsvc, nbest, "1")
}
