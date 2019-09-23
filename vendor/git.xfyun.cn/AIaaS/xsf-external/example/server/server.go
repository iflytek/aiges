package main

import (
	"errors"
	"fmt"
	"io"

	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"strconv"
	"time"
)

type monitor struct {
}

func (m *monitor) Query(in map[string]string, out io.Writer) {
	_, _ = out.Write([]byte( fmt.Sprintf("%+v", in)))
}

//当接收到syscall.SIGINT, syscall.SIGKILL时，会回调这接口
type killed struct {
}

func (k *killed) Closeout() {
	fmt.Println("server be killed.")
}

type healthChecker struct {
}

//服务自检接口，cmdserver用
func (h *healthChecker) Check() error {
	return errors.New("this is check function from health check")
}

//用户业务逻辑接口
type server struct {
	tool *xsf.ToolBox
}

//业务初始化接口
func (c *server) Init(toolbox *xsf.ToolBox) error {
	fmt.Println("begin init")
	c.tool = toolbox
	{
		xsf.AddKillerCheck("server", &killed{})
		xsf.AddHealthCheck("server", &healthChecker{})
		xsf.StoreMonitor(&monitor{})
	}
	fmt.Println("server init success.")
	return nil
}

//业务逆初始化接口
func (c *server) FInit() error {
	fmt.Println("user logic FInit success.")
	return nil
}

//业务服务接口
func (c *server) Call(in *xsf.Req, span *xsf.Span) (response *utils.Res, err error) {
	switch in.Op() {
	case "ssb":
		return c.ssbRouter(in)
	case "auw":
		return c.auwRouter(in)
	case "sse":
		return c.sseRouter(in)
	case "req":
		return c.reqRouter(in)
	default:
		break
	}
	return c.unknown(in)
}

func (c *server) unknown(in *xsf.Req) (*utils.Res, error) {
	fmt.Printf("the op -> %v is not supported.\n", in.Op())
	res := xsf.NewRes()
	res.SetHandle(in.Handle())
	res.SetError(1, fmt.Sprintf("the op -> %v is not supported.", in.Op()))
	res.SetParam("intro", "received data")
	res.SetParam("op", "illegal")
	return res, nil
}

func (c *server) reqRouter(in *xsf.Req) (*utils.Res, error) {
	c.tool.Log.Debugw("in process", "op", "req")
	res := xsf.NewRes()
	res.SetHandle(in.Handle())
	res.SetParam("intro", "received data")
	res.SetParam("op", "req")
	res.SetParam("ip", c.tool.NetManager.GetIp())
	res.SetParam("port", strconv.Itoa(c.tool.NetManager.GetPort()))
	data := xsf.NewData()
	data.SetParam("intro", "for test")
	data.Append([]byte("test data"))
	res.AppendData(data)
	return res, nil
}

func (c *server) sseRouter(in *xsf.Req) (*utils.Res, error) {
	c.tool.Log.Debugw("in process", "op", "sse")
	res := xsf.NewRes()
	res.SetHandle(in.Handle())
	c.tool.Cache.DelSessionData(in.Handle())
	_ = c.tool.Cache.Update()
	{
		res.SetParam("intro", "received data")
		res.SetParam("op", "sse")
		res.SetParam("ip", c.tool.NetManager.GetIp())
		res.SetParam("port", strconv.Itoa(c.tool.NetManager.GetPort()))
	}
	return res, nil
}

func (c *server) auwRouter(in *xsf.Req) (*utils.Res, error) {
	c.tool.Log.Debugw("in process", "op", "auw")
	res := xsf.NewRes()
	res.SetHandle(in.Handle())
	if _, GetSessionDataErr := c.tool.Cache.GetSessionData(in.Handle()); GetSessionDataErr != nil {
		res.SetError(1, fmt.Sprintf("GetSessionData failed. ->GetSessionDataErr:%v", GetSessionDataErr))
	}
	{
		res.SetParam("intro", "received data")
		res.SetParam("op", "auw")
		res.SetParam("ip", c.tool.NetManager.GetIp())
		res.SetParam("port", strconv.Itoa(c.tool.NetManager.GetPort()))
	}
	return res, nil
}

func (c *server) ssbRouter(in *xsf.Req) (*utils.Res, error) {
	c.tool.Log.Debugw("in process", "op", "ssb")
	res := xsf.NewRes()
	res.SetHandle(in.Handle())
	sessionCb := func(sessionTag interface{}, svcData interface{}, exception ...xsf.CallBackException) {
		c.tool.Log.Infow("this is callback function", "timestamp", time.Now(), sessionTag, in.Handle())
	}
	SetSessionDataErr := c.tool.Cache.SetSessionData(in.Handle(), "svcData", sessionCb)
	if nil != SetSessionDataErr {
		res.SetError(1, fmt.Sprintf("Set %s failed. ->SetErr:%v ->addr:%v",
			in.Handle(), SetSessionDataErr, fmt.Sprintf("%v:%v", c.tool.NetManager.GetIp(), c.tool.NetManager.GetPort())))
	} else {
		_ = c.tool.Cache.Update()
	}
	return res, nil
}
