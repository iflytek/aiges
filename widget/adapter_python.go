package widget

import "C"
import (
	"fmt"
	"github.com/hashicorp/go-plugin"
	"github.com/xfyun/aiges/grpc/proto"
	"github.com/xfyun/aiges/grpc/shared"
	"github.com/xfyun/aiges/instance"
	"io"
	"log"
	"os"
	"os/exec"
)

///	wrapper适配器,提供golang至c/c++ wrapper层的数据适配及接口转换;
type enginePython struct {
	client    *plugin.Client
	rpcClient plugin.ClientProtocol
	wrapper   shared.PyWrapper
	stream    proto.WrapperService_CommunicateClient
}

func (ep *enginePython) open(cmd string) (errInfo error) {
	// We're a host. Start by launching the plugin process.
	ep.client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		SyncStdout:      os.Stdout,
		Cmd:             exec.Command("sh", "-c", cmd),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC},
	})
	// Connect via RPC
	var err error
	ep.rpcClient, err = ep.client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return err
	}
	wrapper, err := ep.rpcClient.Dispense("wrapper_grpc")
	if err != nil {
		fmt.Println("Error:", err.Error())
		return err

	}
	ep.wrapper = wrapper.(shared.PyWrapper)
	ep.stream, err = ep.wrapper.Communicate()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return err
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := ep.stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}

			engineCreateCallBackPy(in)
			// query handle
			if in.Tag != "" {

			}
			log.Printf("Got data %s at key , len(%s, %d)", in.List[0].Data, in.List[0].Key, in.List[0].Len)
		}
	}()
	notes := []*proto.Request{
		{Params: map[string]string{"Latitude": "0"}},
	}
	for _, note := range notes {
		if err := ep.stream.Send(note); err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", note, err)
		}
	}

	return
}

func (ep *enginePython) close() {
	ep.client.Kill()
	return
}

func (ep *enginePython) enginePythonInit(cfg map[string]string) (errNum int, errInfo error) {
	fmt.Println("call python wrapper Init ##")
	// Request the plugin
	ep.wrapper.WrapperInit(cfg)
	return
}

func (ep *enginePython) enginePythonOnceExec(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	log.Println("！ handle is :", handle)
	log.Println("req aysncmode is :", req.AsyncCode)
	log.Println("req wrapper handle :", req.WrapperHdl)
	log.Println("req params:", req.Params)
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
	go ep.wrapper.WrapperOnceExec(handle, req.Params, datas)

	return
}

func (ep *enginePython) enginePythonFini() (errNum int, errInfo error) {
	//ret := C.adapterFini()
	//if ret != 0 {
	//	errInfo = errors.New(enginePythonError(int(ret)))
	//	errNum = int(ret)
	//}
	return
}

func (ep *enginePython) enginePythonVersion() (ver string) {
	return "Devel"
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
	//uuid := catch.GenerateUuid()
	//catch.CallCgo(uuid, catch.Begin)
	//defer catch.CallCgo(uuid, catch.End)
	//// 参数对;
	//
	//paramList := C.paramListCreate()
	//defer C.paramListfree(paramList)
	//for k, v := range req.Params {
	//	key := C.CString(k)
	//	defer C.free(unsafe.Pointer(key))
	//	val := C.CString(v)
	//	defer C.free(unsafe.Pointer(val))
	//	valLen := C.uint(len(v))
	//	paramList = C.paramListAppend(paramList, key, val, valLen)
	//}
	//
	//// 个性化;
	//var psrPtr *C.uint
	//var psrIds []C.uint
	//psrCnt := len(req.PersonIds)
	//if psrCnt > 0 {
	//	psrIds = make([]C.uint, psrCnt)
	//	for k, v := range req.PersonIds {
	//		psrIds[k] = C.uint(v)
	//	}
	//	psrPtr = &psrIds[0]
	//}
	//
	//var errC C.int
	//var callback C.wrapperCallback = nil
	//if conf.WrapperAsync {
	//	callback = C.wrapperCallback(C.adapterCallback)
	//}
	//
	//usrTag := C.CString(handle)
	////defer C.free(unsafe.Pointer(usrTag)) callBack free
	//engHdl := C.adapterCreate(usrTag, paramList, callback, psrPtr, C.int(psrCnt), &errC)
	//if errC != 0 || engHdl == nil {
	//	errNum = int(errC)
	//	errInfo = errors.New(enginePythonError(int(errC)))
	//} else {
	//	resp.WrapperHdl = engHdl
	//}
	//
	////	C.adapterFreeParaList(pParaHead)
	//
	//if !conf.WrapperAsync {
	//	C.free(unsafe.Pointer(usrTag))
	//}
	return
}

// 资源释放行为
func (ep *enginePython) enginePythonDestroy(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	//uuid := catch.GenerateUuid()
	//catch.CallCgo(uuid, catch.Begin)
	//defer catch.CallCgo(uuid, catch.End)
	//errC := C.adapterDestroy(req.WrapperHdl)
	//if errC != 0 {
	//	errNum = int(errC)
	//	errInfo = errors.New(enginePythonError(int(errC)))
	//}
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
func (ep *enginePython) enginePythonWrite(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	//uuid := catch.GenerateUuid()
	//catch.CallCgo(uuid, catch.Begin)
	//defer catch.CallCgo(uuid, catch.End)
	//// 写数据转换;
	//dataList := C.dataListCreate()
	//defer C.DataListfree(dataList)
	//for _, ele := range req.DeliverData {
	//	tmpKey := C.CString(ele.DataId)
	//	defer C.free(unsafe.Pointer(tmpKey))
	//	tmpData := C.CBytes(ele.Data)
	//	defer C.free(unsafe.Pointer(tmpData))
	//
	//	descList := C.paramListCreate()
	//	defer C.paramListfree(descList)
	//	for k, v := range ele.DataDesc {
	//		key := C.CString(k)
	//		defer C.free(unsafe.Pointer(key))
	//		val := C.CString(v)
	//		defer C.free(unsafe.Pointer(val))
	//		valLen := C.uint(len(v))
	//		descList = C.paramListAppend(descList, key, val, valLen)
	//	}
	//	dataList = C.dataListAppend(dataList, tmpKey, tmpData, C.uint(len(ele.Data)), C.int(ele.DataType), C.int(ele.DataStatus), descList)
	//}
	//
	//errC := C.adapterWrite(req.WrapperHdl, dataList)
	//if errC != 0 {
	//	errNum = int(errC)
	//	errInfo = errors.New(enginePythonError(int(errC)))
	//}
	//
	////	C.adapterFreeDataList(pDataHead)
	return
}

// 数据读行为
func (ep *enginePython) enginePythonRead(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	//uuid := catch.GenerateUuid()
	//catch.CallCgo(uuid, catch.Begin)
	//defer catch.CallCgo(uuid, catch.End)
	//respDataC := C.getDataList()
	//defer C.releaseDataList(respDataC)
	//errC := C.adapterRead(req.WrapperHdl, respDataC)
	//if errC != 0 {
	//	errNum = int(errC)
	//	errInfo = errors.New(enginePythonError(int(errC)))
	//	return
	//}
	//
	//// 读数据转换;
	////	resp.DeliverData = make([]instance.DataMeta, 0, 1)
	//for *respDataC != nil {
	//	var ele instance.DataMeta
	//	ele.DataId = C.GoString((*(*respDataC)).key)
	//	ele.DataType = int((*(*respDataC))._type)
	//	ele.DataStatus = int((*(*respDataC)).status)
	//	ele.DataDesc = make(map[string]string)
	//	pDesc := (*(*respDataC)).desc
	//	for pDesc != nil {
	//		ele.DataDesc[C.GoString((*pDesc).key)] = C.GoStringN((*pDesc).value, C.int((*pDesc).vlen))
	//		pDesc = (*pDesc).next
	//	}
	//	if int((*(*respDataC)).len) != 0 {
	//		ele.Data = C.GoBytes((*(*respDataC)).data, C.int((*(*respDataC)).len))
	//	}
	//	resp.DeliverData = append(resp.DeliverData, ele)
	//	*respDataC = (*respDataC).next
	//}

	return
}

// 计算debug数据
func (ep *enginePython) enginePythonDebug(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	//uuid := catch.GenerateUuid()
	//catch.CallCgo(uuid, catch.Begin)
	//defer catch.CallCgo(uuid, catch.End)
	//debug := C.adapterDebugInfo(req.WrapperHdl)
	//resp.Debug = C.GoString(debug)
	return
}

func enginePythonError(errNum int) (errInfo string) {
	//err := C.adapterError(C.int(errNum))
	return "errr"
}