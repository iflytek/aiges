package daemon

import (
	"encoding/json"
	"flag"
	"fmt"
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func Usage() {
	_, _ = fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	_, _ = fmt.Fprint(os.Stderr, "\n")
}

func RunServer() (err error) {
	flag.Usage = Usage
	flag.Parse()
	var serverInst xsf.XsfServer

	if err = serverInst.Run(xsf.BootConfig{
		CfgMode: utils.CfgMode(-1),
		CfgData: xsf.CfgMeta{
			CfgName:      "",
			Project:      "",
			Group:        "",
			Service:      "",
			Version:      lbVersion,
			ApiVersion:   lbApiVersion,
			CompanionUrl: ""}}, &Server{}); nil != err {
		log.Panic(err)
	}
	return
}

type monitor struct {
}

func (m *monitor) Query(in map[string]string, out io.Writer) {
	var rst []byte
	if "1" == in["down"] {
		rst = func() []byte {
			rstByte, _ := json.Marshal(getAbnormalNodeStats())
			return rstByte
		}()
	} else if forceOffline, ok := in[FORCEOFFLINE]; ok {
		rst = func() []byte {
			var status string
			engAddr := strings.Replace(in["engAddr"], `"`, "", -1)
			switch forceOffline {
			case ForceOffline:
				{
					blacklist.Store(engAddr, struct{}{})
					status = fmt.Sprintf(
						"successfully adding %v to blacklist(%v)",
						engAddr,
						blackListContent())
				}
			case CleanForceOffline:
				{
					blacklist.Delete(engAddr)
					status = fmt.Sprintf(
						"successfully deleting %v from blacklist(%v)",
						engAddr,
						blackListContent())
				}
			default:
				{
					status = ErrCmdServerIsIncorrect.Error()
				}
			}
			return []byte(status)
		}()
	} else {
		rst = func() []byte {
			svc := in["svc"]
			subsvc := in["subsvc"]
			return []byte(monitorWareHouseInst.query(svc, subsvc))
		}()
	}
	_, _ = out.Write(rst)
}

type killed struct {
}

func (k *killed) Closeout() {
	std.Println("lb is be killed")
}

type healthChecker struct {
}

func (h *healthChecker) Check() (err error) {
	return nil
}

type Server struct {
	toolbox *xsf.ToolBox
	LbHandle
}

func (s *Server) Init(toolbox *xsf.ToolBox) (err error) {
	std.Println("about to Server.init")

	s.toolbox = toolbox
	mssSidGenerator.Init("lb2", s.toolbox.NetManager.GetIp(), "xx")

	pprofInt, pprofErr := s.toolbox.Cfg.GetInt(BO, PPROF)
	if pprofErr == nil {
		pprofSrv(pprofInt)
	}

	xsf.AddKillerCheck(ComponentName, &killed{})
	xsf.AddHealthCheck(ComponentName, &healthChecker{})
	xsf.StoreMonitor(&monitor{})

	err = s.LbHandle.Init(s.toolbox)
	std.Println("Server.init Ok.")
	if nil != err {
		s.toolbox.Log.Errorw(
			"lb handle init",
			"lbInitError", err.Error())
		return err
	}
	return nil
}

func (s *Server) FInit() error {
	time.Sleep(time.Second * 15)
	return nil
}

func (s *Server) Call(in *xsf.Req, span *xsf.Span) (res *utils.Res, err error) {
	res = xsf.NewRes()
	switch s.strategy {
	case load:
		{
			s.toolbox.Log.Debugw("fn:Call=load")
			res, err = s.worker.serve(in, span, s.toolbox)
		}
	case poll:
		{
			s.toolbox.Log.Debugw("fn:Call=poll")
			res, err = s.worker.serve(in, span, s.toolbox)
		}
	case loadMini:
		{
			s.toolbox.Log.Debugw("fn:Call=ise")
			res, err = s.worker.serve(in, span, s.toolbox)
		}
	default:
		{
			s.toolbox.Log.Errorf("LbErrInputStrategy")
			res.SetError(ErrLbStrategyIsNotSupport.errCode, ErrLbStrategyIsNotSupport.errInfo)
		}
	}
	return res, err
}
