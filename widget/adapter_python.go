package widget

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/grpc/proto"
	"github.com/xfyun/aiges/grpc/shared"
	"github.com/xfyun/aiges/httproto/schemas"
	"github.com/xfyun/aiges/instance"
	"github.com/xfyun/aiges/utils"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
	"unsafe"
)

// /	wrapper适配器,提供golang至c/c++ wrapper层的数据适配及接口转换;
type enginePython struct {
	client    *plugin.Client
	rpcClient plugin.ClientProtocol
	wrapper   shared.PyWrapper
	stream    proto.WrapperService_CommunicateClient
	Schema    string
}

func (ep *enginePython) open(ch *utils.Coordinator) (errInfo error) {
	// open 似乎没必要
	var a = <-ch.ConfChan
	//logLevelStr, _ := cfg["log.level"]

	logLevelStr, ok := conf.UsrCfgData["log.level"]
	if !ok {
		logLevelStr = "info"
	}
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "python-plugin",
		Output: os.Stdout,
		Level:  enginePythonLogLvl(logLevelStr),
	})
	// We're a host. Start by launching the plugin process.
	ep.client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		SyncStdout:      os.Stdout,
		Cmd:             exec.Command("bash", "-c", conf.PythonCmd),
		Logger:          logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC},
	})
	// Connect via RPC
	var err error
	for i := 0; i < 3; i++ {
		ep.rpcClient, err = ep.client.Client()
		if err != nil {
			time.Sleep(time.Second * 5)
			fmt.Printf("Retrying connect ...")
			continue
		} else {
			break
		}
	}
	if err != nil {
		ep.client.Kill()
		log.Fatalln("Error:", err.Error())
		return err

	}

	wrapper, err := ep.rpcClient.Dispense("wrapper_grpc")
	if err != nil {
		ep.client.Kill()
		log.Fatalln("Error:", err.Error())
		return err

	}
	ep.wrapper = wrapper.(shared.PyWrapper)
	ep.stream, err = ep.wrapper.Communicate()
	if err != nil {
		ep.client.Kill()
		log.Fatalln("Error:", err.Error())
		return err
	}
	//waitc := make(chan struct{})
	countErr := 0
	go func() {
		for {
			in, err := ep.stream.Recv()
			if err == io.EOF {
				// read done.
				//close(waitc)
				continue
			}
			if err != nil {
				countErr += 1
				log.Printf("Client Recv the response failed: %v, retrying... %d\n", err, countErr)
				if countErr >= 4 {
					ep.client.Kill()
					return
				}
				time.Sleep(time.Second * 1)
				continue
			}
			// query handle
			if in.Tag != "" {
				engineCreateCallBackPy(in)
			}
		}
	}()
	ch.ConfChan <- a
	return
}

func (ep *enginePython) close() {
	ep.client.Kill()
	return
}

func (ep *enginePython) enginePythonInit(cfg map[string]string) (errNum int, errInfo error) {

	// Init the plugin
	err := ep.wrapper.WrapperInit(cfg)
	if err != nil {
		ep.client.Kill()
		log.Fatalf("err: %v\n", err)
		return -1, err
	}
	// Get schema 这里传入的参数目前无用，，实际python测那边没有用到
	schema, err := ep.wrapper.WrapperSchema("svcName")
	if err != nil {
		ep.client.Kill()
		log.Fatalln("Error:", err.Error())
		return -1, err
	}

	// 设置schema
	schemas.SetSchemaFromPython(schema.GetData())

	ep.Schema = schema.GetData()

	return
}

func (ep *enginePython) enginePythonOnceExec(userTag string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	var datas []*proto.RequestData
	for _, dd := range req.DeliverData {
		datas = append(
			datas,
			&proto.RequestData{
				Data:   dd.Data,
				Key:    dd.DataId,
				Len:    uint64(len(dd.Data)),
				Type:   uint32(dd.DataType),
				Status: uint32(dd.DataStatus),
			},
		)
	}
	// 这里只需要把handle、tag带过去， grpc 那边通过双工流返回回来即可。
	ep.wrapper.WrapperOnceExec(userTag, req.Params, datas)

	return
}

func (ep *enginePython) enginePythonFini() (errNum int, errInfo error) {
	log.Println("Calling Python Fini in Aiges...")
	return
}

func (ep *enginePython) enginePythonVersion() (ver string) {
	return "Devel-3.0"
}

// 资源加载卸载管理适配接口;
func (ep *enginePython) enginePythonLoadRes(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {

	return
}

func (ep *enginePython) enginePythonUnloadRes(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	//uuid := catch.GenerateUuid()
	//catch.CallCgo(uuid, catch.Begin)
	//defer catch.CallCgo(uuid, catch.End)
	//errC := C.adapterUnloadRes(C.uint(req.PsrId))
	//if errC != 0 {
	//	errNum = int(errC)
	//	errInfo = errors.New(enginePythonError(int(errC)))
	//}
	return
}

// 资源申请行为
func (ep *enginePython) enginePythonCreate(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {

	// Init the plugin
	hdl, err := ep.wrapper.WrapperCreate(handle, req.Params)
	if err != nil {
		errNum = int(30002)
		errInfo = err
		return
	}
	resp.WrapperHdl = unsafe.Pointer(&hdl.Handle)
	b := (*[]byte)(resp.WrapperHdl)
	fmt.Printf(handle, string(*b))
	return
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// 资源释放行为
func (ep *enginePython) enginePythonDestroy(userTag string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	//uuid := catch.GenerateUuid()
	hdl := (*[]byte)(req.WrapperHdl)
	h := string(*hdl)
	ret, errInfo := ep.wrapper.WrapperDestroy(h)
	errNum = int(ret.Ret)
	return
}

// 交互异常行为
func (ep *enginePython) enginePythonExcp(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	//uuid := catch.GenerateUuid()
	//catch.CallCgo(uuid, catch.Begin)
	//defer catch.CallCgo(uuid, catch.End)
	//resp, errNum, errInfo = enginePythonDestroy(handle, req)

	return
}

// 数据写行为
func (ep *enginePython) enginePythonWrite(userTag string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	var datas []*proto.RequestData
	for _, dd := range req.DeliverData {
		datas = append(
			datas,
			&proto.RequestData{
				Data:   dd.Data,
				Key:    dd.DataId,
				Len:    uint64(len(dd.Data)),
				Type:   uint32(dd.DataType),
				Status: uint32(dd.DataStatus),
			},
		)
	}
	hdl := string(*(*[]byte)(req.WrapperHdl))
	ret, err := ep.wrapper.WrapperWrite(hdl, userTag, req.Params, datas)
	if err != nil {
		errNum = int(ret.Ret)
		errInfo = err
	}
	return
}

// 数据读行为
func (ep *enginePython) enginePythonRead(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	fmt.Println("cccc")
	return
}

// 计算debug数据
func (ep *enginePython) enginePythonDebug(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {

	return
}

func enginePythonError(errNum int) (errInfo string) {
	//err := C.adapterError(C.int(errNum))
	return "errr"
}

func enginePythonLogLvl(logstr string) hclog.Level {
	switch logstr {
	case "debug":
		return hclog.Debug
	case "info":
		return hclog.Info
	case "warn":
		return hclog.Warn
	default:
		return hclog.Error
	}
}
