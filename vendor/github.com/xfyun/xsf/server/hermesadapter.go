package xsf

import (
	"fmt"
	"github.com/xfyun/xsf/client"
	"github.com/xfyun/xsf/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	CNAME        = "lbv2"
	MoveToSvc    = "moveToSvc"
	MoveToSubSvc = "moveToSubSvc"
	RawSvc       = "rawSvc"
	RawSubSvc    = "rawSubSvc"
)

type SvcUnit struct {
	Svc    string `json:"svc"`
	SubSvc string `json:"sub_svc"`
}
type MoveTo struct {
	Original SvcUnit `json:"original"`
	From     SvcUnit `json:"from"`
	Now      SvcUnit `json:"now"`
}
type hermesAdapter struct {
	mode           utils.CfgMode
	able           bool
	originalSvc    string
	originalSubSvc string
	fromSvc        string
	fromSubSvc     string
	extra          string //保存负载以外的一些信息
	svc            string
	subsvc         string
	addr           string
	uid            string
	total          int
	idle           int
	best           int

	lbname     string
	lbprj      string
	lbgro      string
	apiversion string
	finderTtl  time.Duration //更新本地地址的时间，通过访问服务发现实现
	backend    int           //上报的的协程数，缺省4
	timeout    time.Duration //上报的超时时间，缺省一秒

	bc BootConfig

	cli    *xsf.Client
	caller *xsf.Caller

	detectors []*detector

	taskInChan chan func() //任务通道，用来传送上报任务
	//taskOutChan chan callRstItem

	hermesTask int
	svcIp      string //服务端监听ip，trace用
	svcPort    int32  //服务端监听端口，trace用

	cloud string //cloud_id
}
type HermesAdapterCfgOpt func(*hermesAdapter)

func WithHermesAdapterCloudId(cloud string) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.cloud = cloud
	}
}

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
func WithHermesAdapterLbPrj(lbprj string) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.lbprj = lbprj
	}
}
func WithHermesAdapterLbGro(lbgro string) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.lbgro = lbgro
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
func WithHermesAdapterMode(mode utils.CfgMode) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.mode = mode
	}
}
func WithHermesAdapterSvcAndSubSvc(svc, subSvc string) HermesAdapterCfgOpt {
	return func(h *hermesAdapter) {
		h.svc = svc
		h.originalSvc = svc

		h.subsvc = subSvc
		h.originalSubSvc = subSvc

		extra, extraErr := utils.AddExtraTag("", map[string]string{RawSvc: h.originalSvc, RawSubSvc: h.originalSubSvc})
		if extraErr != nil {
			panic(fmt.Sprintf("addExtraTag failed,err:%v", extraErr))
		}
		h.extra = extra
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
	loggerStd.Printf("hermes init->cfgName:%v,companion:%v,prj:%v,grp:%v,srv:%v,ver:%v\n",
		h.bc.CfgData.CfgName, h.bc.CfgData.CompanionUrl, h.bc.CfgData.Project, h.bc.CfgData.Group, h.bc.CfgData.Service, h.bc.CfgData.Version)
	h.cli, err = xsf.InitClient(
		CNAME,
		h.mode,
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
	if err != nil {
		panic(fmt.Sprintf("InitClient fail err:%v", err))
		return
	}
	h.caller = xsf.NewCaller(h.cli)
	h.caller.WithApiVersion(h.apiversion)
	h.taskInChan = make(chan func(), h.hermesTask)

	detectorPrj, detectorGro, detectorNam := func() (string, string, []string) {
		var prj, gro string
		var nam []string
		if len(h.lbprj) != 0 {
			prj = h.lbprj
		} else {
			prj = h.bc.CfgData.Project
		}
		if len(h.lbgro) != 0 {
			gro = h.lbgro
		} else {
			gro = h.bc.CfgData.Group
		}
		nam = strings.Split(h.lbname, ",")
		return prj, gro, nam
	}()
	loggerStd.Printf("hermes detector params,prj:%v,gro:%v,nam:%v\n", detectorPrj, detectorGro, detectorNam)
	for _, nam := range detectorNam {
		detector, detectorErr := newDetector(h.bc.CfgData.CompanionUrl, detectorPrj, detectorGro, nam, h.cli.Log)
		if detectorErr != nil {
			panic(fmt.Sprintf("init detector fail err:%v", err))
		}
		h.detectors = append(h.detectors, detector)
	}

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

const MoveFromSvc = "moveFromSvc"
const MoveFromSubSvc = "moveFromSubSvc"
const Extra = "extra"

var reportFlag = reportDef //这部分后续优化吧,-1:not set,0:failure,1:success

func (h *hermesAdapter) setServer(from SvcUnit, svc, subsvc, addr string, getAuthInfo func() (int32, int32, int32), live string) error {
	if !h.able {
		return nil
	}

	h.cli.Log.Infow("fn:setServer",
		"from", from, "svc", svc, "subsvc", subsvc, "addr", addr)

	var reportAddrs []string
	for _, detector := range h.detectors {
		reportAddrs = append(reportAddrs, detector.getAll()...)
	}

	h.cli.Log.Infow("get reportAddrs", "lbname", h.lbname, "addrs", reportAddrs)

	setServerOnce.Do(func() {
		loggerStd.Printf("start reporting->target:%v\n", reportAddrs)
	})

	for _, lbAddr := range reportAddrs {
		h.cli.Log.Infow("report task", "lbname", h.lbname, "addr", lbAddr)
		lbAddrTmp := lbAddr
		task := func() {

			req := utils.NewReq()
			req.SetParam(HERMESLBCLOUD, h.cloud)
			req.SetParam(Extra, h.extra)
			req.SetParam(MoveFromSvc, from.Svc)
			req.SetParam(MoveFromSubSvc, from.SubSvc)
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

			res, errcode, e := h.caller.CallWithAddr(h.lbname, xsf.LBOPSET, lbAddrTmp, req, time.Second)
			if errcode != 0 || e != nil {
				reportFlag = reportFailure
				h.cli.Log.Errorw("fn:setServer h.caller.CallWithAddr", "errcode", errcode, "err", e, "addr", lbAddrTmp)
			} else {
				h.FilterMoveTo(res)
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

func (h *hermesAdapter) FilterMoveTo(res *xsf.Res) {
	moveToSvc, moveToSvcOk := res.GetParam(MoveToSvc)
	moveToSubSvc, moveToSubSvcOk := res.GetParam(MoveToSubSvc)
	if !moveToSvcOk || !moveToSubSvcOk || len(moveToSvc) == 0 || len(moveToSubSvc) == 0 {
		h.cli.Log.Infow("fn:FilterMoveTo",
			"moveToSvcOk", moveToSvcOk, "moveToSubSvcOk", moveToSubSvcOk,
			"moveToSvc", moveToSvc, "moveToSubSvc", moveToSubSvc,
		)
		return
	}
	h.fromSvc = h.svc
	h.fromSubSvc = h.subsvc
	h.svc = moveToSvc
	h.subsvc = moveToSubSvc

	if h.svc == h.originalSvc && h.subsvc == h.originalSubSvc {
		//如果归还则移除from标识
		h.fromSvc = ""
		h.fromSubSvc = ""
	}

	h.cli.Log.Infow("fn:FilterMoveTo",
		"fromSvc", h.fromSvc, "fromSubSvc", h.fromSubSvc,
		"svc", h.svc, "subsvc", h.subsvc,
		"originalSvc", h.originalSvc, "originalSubSvc", h.originalSubSvc,
	)
	return
}

func (h *hermesAdapter) report(getAuthInfo func() (int32, int32, int32)) error {
	if !h.able {
		return nil
	}
	return h.setServer(SvcUnit{h.fromSvc, h.fromSubSvc}, h.svc, h.subsvc, h.addr, getAuthInfo, "1")
}
func (h *hermesAdapter) offline() error {
	if !h.able {
		return nil
	}
	return h.setServer(SvcUnit{h.svc, h.subsvc}, h.svc, h.subsvc, h.addr, nil, "0")
}
