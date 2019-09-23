package widget

/*
#cgo linux CFLAGS: -I../cgo/header -Wno-attributes
#cgo linux LDFLAGS: -ldl
#include <stdlib.h>
#include <dlfcn.h>
#include "./widget/wrapper.h"

// @return library handle
void* cLibOpen(const char* libName, char** err){
	void* hdl = dlopen(libName, RTLD_NOW);
	if (hdl == NULL){
		*err = (char*)dlerror();
	}
	return hdl;
}

// @return symbol address
void* cLibLoad(void* hdl, const char* sym, char** err){
	void* addr = dlsym(hdl, sym);
	if (addr == NULL){
		*err = (char*)dlerror();
	}
	return addr;
}

int  cLibClose(void* hdl){
	int ret = dlclose(hdl);
	if (ret != 0)
		return -1;
	return 0;
}

// 符号列表
const char* InitSym = "wrapperInit";
const char* FiniSym = "wrapperFini";
const char* ErrorSym = "wrapperError";
const char* VerSym = "wrapperVersion";
const char* LoadSym = "wrapperLoadRes";
const char* UnloadSym = "wrapperUnloadRes";
const char* CreateSym = "wrapperCreate";
const char* WriteSym = "wrapperWrite";
const char* ReadSym = "wrapperRead";
const char* DestroySym = "wrapperDestroy";
const char* ExecSym = "wrapperExec";
const char* ExecFreeSym = "wrapperExecFree";
const char* ExecAsynSym = "wrapperExecAsync";
const char* DebugSym = "wrapperDebugInfo";

// 接口定义
wrapperInitPtr cPtrInit;
wrapperFiniPtr cPtrFini;
wrapperErrorPtr cPtrError;
wrapperVersionPtr cPtrVersion;
wrapperLoadResPtr cPtrLoadRes;
wrapperUnloadResPtr cPtrUnloadResPtr;
wrapperCreatePtr cPtrCreate;
wrapperWritePtr cPtrWrite;
wrapperReadPtr cPtrRead;
wrapperDestroyPtr cPtrDestroy;
wrapperExecPtr cPtrExec;
wrapperExecFreePtr cPtrExecFree;
wrapperExecAsyncPtr cPtrExecAsync;
wrapperDebugInfoPtr cPtrDebugInfo;

// 接口寻址
// @return excp symbol if load fail, else NULL;
const char* loadWrapperSym(void* hdl, char** loadErr){
	// load all interface
	if ((cPtrInit = cLibLoad(hdl, InitSym, loadErr)) == NULL)
		return InitSym;
	if ((cPtrFini = cLibLoad(hdl, FiniSym, loadErr)) == NULL)
		return FiniSym;
	if ((cPtrError = cLibLoad(hdl, ErrorSym, loadErr)) == NULL)
		return ErrorSym;
	if ((cPtrVersion = cLibLoad(hdl, VerSym, loadErr)) == NULL)
		return VerSym;
	if ((cPtrLoadRes = cLibLoad(hdl, LoadSym, loadErr)) == NULL)
		return LoadSym;
	if ((cPtrUnloadResPtr = cLibLoad(hdl, UnloadSym, loadErr)) == NULL)
		return UnloadSym;
	if ((cPtrCreate = cLibLoad(hdl, CreateSym, loadErr)) == NULL)
		return CreateSym;
	if ((cPtrWrite = cLibLoad(hdl, WriteSym, loadErr)) == NULL)
		return WriteSym;
	if ((cPtrRead = cLibLoad(hdl, ReadSym, loadErr)) == NULL)
		return ReadSym;
	if ((cPtrDestroy = cLibLoad(hdl, DestroySym, loadErr)) == NULL)
		return DestroySym;
	if ((cPtrExec = cLibLoad(hdl, ExecSym, loadErr)) == NULL)
		return ExecSym;
	if ((cPtrExecFree = cLibLoad(hdl, ExecFreeSym, loadErr)) == NULL)
		return ExecFreeSym;
	if ((cPtrExecAsync = cLibLoad(hdl, ExecAsynSym, loadErr)) == NULL)
		return ExecAsynSym;
	if ((cPtrDebugInfo = cLibLoad(hdl, DebugSym, loadErr)) == NULL)
		return DebugSym;
	return NULL;
}

// 接口适配封装
int adapterInit(pConfig cfg){
	return (*cPtrInit)(cfg);
}

int adapterFini(){
	return (*cPtrFini)();
}

const char* adapterVersion(){
	return (*cPtrVersion)();
}

const char* adapterError(int errNum){
	return (*cPtrError)(errNum);
}

int adapterLoadRes(pDataList perData, unsigned int resId){
	return (*cPtrLoadRes)(perData, resId);
}

int adapterUnloadRes(unsigned int resId){
	return (*cPtrUnloadResPtr)(resId);
}

const char* adapterCreate(const void* usrTag, pParamList params, wrapperCallback cb, unsigned int psrIds[], int psrCnt, int* errNum){
	return (*cPtrCreate)(usrTag, params, cb, psrIds, psrCnt, errNum);
}

int adapterWrite(const char* handle, pDataList reqData){
	return (*cPtrWrite)(handle, reqData);
}

int adapterRead(const char* handle, pDataList* respData){
	return (*cPtrRead)(handle, respData);
}

int adapterDestroy(const char* handle){
	return (*cPtrDestroy)(handle);
}

int adapterExec(pParamList params, pDataList reqData, pDataList* respData){
	return (*cPtrExec)(params, reqData, respData);
}

int adapterExecFree(pDataList* respData){
	return (*cPtrExecFree)(respData);
}

int adapterExecAsync(const void* handle, pParamList params, pDataList reqData, wrapperCallback callback, int timeout){
	return (*cPtrExecAsync)(handle, params, reqData, callback, timeout);
}

const char* adapterDebugInfo(const char* handle){
	return (*cPtrDebugInfo)(handle);
}

extern int engineCreateCallBack(void* handle, pDataList respData, int ret);
int adapterCallback(const void* handle, pDataList respData, int ret){
	return engineCreateCallBack((void*)handle, respData, ret);
}
*/
import "C"
import (
	"conf"
	"errors"
	"instance"
	"unsafe"
)

///	wrapper适配器,提供golang至c/c++ wrapper层的数据适配及接口转换;
type engineC struct {
	libHdl unsafe.Pointer // 引擎句柄;
}

func (ec *engineC) open(libName string) (errInfo error) {
	var errC *C.char
	var libNameC *C.char = C.CString(libName)
	defer C.free(unsafe.Pointer(libNameC))
	ec.libHdl = C.cLibOpen(libNameC, &errC)
	if ec.libHdl == nil {
		errInfo = errors.New("wrapper.open: load library " + libName + "failed, " + C.GoString(errC))
		return
	}

	errSym := C.loadWrapperSym(ec.libHdl, &errC)
	if errSym != nil {
		C.cLibClose(ec.libHdl)
		errInfo = errors.New("wrapper.open: load symbol " + C.GoString(errSym) + " failed, " + C.GoString(errC))
		return
	}
	return
}

func (ec *engineC) close() {
	C.cLibClose(ec.libHdl)
	ec.libHdl = nil
	return
}

func engineInit(cfg map[string]string) (errNum int, errInfo error) {
	// 配置参数的语言栈转换;
	var pCfgHead, pCfgTail C.pConfig
	for k, v := range cfg {
		// config k-v
		var pair C.struct_ParamList
		pair.key = C.CString(k)
		pair.value = C.CString(v)
		pair.vlen = C.uint(len(v))
		pair.next = nil
		defer C.free(unsafe.Pointer(pair.key))
		defer C.free(unsafe.Pointer(pair.value))

		if pCfgHead == nil {
			pCfgHead = &pair
			pCfgTail = &pair
		} else {
			(*pCfgTail).next = &pair
			pCfgTail = (*pCfgTail).next
		}
	}

	ret := C.adapterInit(pCfgHead)
	if ret != 0 {
		errInfo = errors.New(engineError(int(ret)))
		errNum = int(ret)
	}
	return
}

func engineFini() (errNum int, errInfo error) {
	ret := C.adapterFini()
	if ret != 0 {
		errInfo = errors.New(engineError(int(ret)))
		errNum = int(ret)
	}
	return
}

func engineVersion() (ver string) {
	verC := C.adapterVersion()
	return C.GoString(verC)
}

// 资源加载卸载管理适配接口;
func engineLoadRes(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	var perDataC C.struct_DataList
	perDataC.data = C.CBytes(req.PsrData)
	defer C.free(unsafe.Pointer(perDataC.data))
	perDataC.desc = C.CString(req.PsrDesc)
	defer C.free(unsafe.Pointer(perDataC.desc))
	perDataC._type = C.DataPer
	perDataC.len = C.uint(len(req.PsrData))

	errC := C.adapterLoadRes(&perDataC, C.uint(req.PsrId))
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
	}
	return
}

func engineUnloadRes(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	errC := C.adapterUnloadRes(C.uint(req.PsrId))
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
	}
	return
}

// 资源申请行为
func engineCreate(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	// 参数对;
	var pParaHead, pParaTail C.pConfig
	for k, v := range req.Params {
		var pair C.struct_ParamList
		pair.key = C.CString(k)
		pair.value = C.CString(v)
		pair.vlen = C.uint(len(v))
		pair.next = nil
		defer C.free(unsafe.Pointer(pair.key))
		defer C.free(unsafe.Pointer(pair.value))

		if pParaHead == nil {
			pParaHead = &pair
			pParaTail = &pair
		} else {
			(*pParaTail).next = &pair
			pParaTail = (*pParaTail).next
		}
	}

	// 个性化;
	var psrPtr *C.uint
	var psrIds []C.uint
	psrCnt := len(req.PersonIds)
	if psrCnt > 0 {
		psrIds = make([]C.uint, psrCnt)
		for k, v := range req.PersonIds {
			psrIds[k] = C.uint(v)
		}
		psrPtr = &psrIds[0]
	}

	var errC C.int
	var callback C.wrapperCallback = nil
	if conf.WrapperAsync {
		callback = C.wrapperCallback(C.adapterCallback)
	}

	usrTag := C.CString(handle)
	//defer C.free(unsafe.Pointer(usrTag)) callBack free
	engHdl := C.adapterCreate(unsafe.Pointer(usrTag), pParaHead, callback, psrPtr, C.int(psrCnt), &errC)
	if errC != 0 || engHdl == nil {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
		return
	}
	resp.WrapperHdl = engHdl
	return
}

// 资源释放行为
func engineDestroy(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	errC := C.adapterDestroy(req.WrapperHdl.(*C.char))
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
	}
	return
}

// 交互异常行为
func engineExcp(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	resp, errNum, errInfo = engineDestroy(handle, req)
	return
}

// 数据写行为
func engineWrite(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	// 写数据转换;
	var pDataHead, pDataTail C.pDataList
	for _, ele := range req.DeliverData {
		var reqDataC C.struct_DataList
		reqDataC._type = C.DataType(ele.DataType)
		reqDataC.status = C.DataStatus(ele.DataStatus)
		reqDataC.encoding = C.CString(ele.DataFmt)
		defer C.free(unsafe.Pointer(reqDataC.encoding))
		reqDataC.desc = C.CString(ele.DataDesc)
		defer C.free(unsafe.Pointer(reqDataC.desc))
		reqDataC.data = C.CBytes(ele.Data)
		defer C.free(unsafe.Pointer(reqDataC.data))
		reqDataC.len = C.uint(len(ele.Data))
		reqDataC.next = nil

		if pDataHead == nil {
			pDataHead = &reqDataC
			pDataTail = &reqDataC
		} else {
			(*pDataTail).next = &reqDataC
			pDataTail = (*pDataTail).next
		}
	}

	errC := C.adapterWrite(req.WrapperHdl.(*C.char), pDataHead)
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
	}
	return
}

// 数据读行为
func engineRead(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	var respDataC C.pDataList
	errC := C.adapterRead(req.WrapperHdl.(*C.char), &respDataC)
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
		return
	}

	// 读数据转换;
	resp.DeliverData = make([]instance.DataMeta, 0, 1)
	for respDataC != nil {
		var ele instance.DataMeta
		ele.DataType = int((*respDataC)._type)
		ele.DataStatus = int((*respDataC).status)
		ele.DataDesc = C.GoString((*respDataC).desc)
		ele.DataFmt = C.GoString((*respDataC).encoding)
		ele.Data = C.GoBytes(unsafe.Pointer((*respDataC).data), C.int((*respDataC).len))
		resp.DeliverData = append(resp.DeliverData, ele)
		respDataC = (*respDataC).next
	}

	return
}

func engineOnceExec(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	// 非会话参数;
	var pParaHead, pParaTail C.pConfig
	for k, v := range req.Params {
		var pair C.struct_ParamList
		pair.key = C.CString(k)
		pair.value = C.CString(v)
		pair.vlen = C.uint(len(v))
		pair.next = nil
		defer C.free(unsafe.Pointer(pair.key))
		defer C.free(unsafe.Pointer(pair.value))

		if pParaHead == nil {
			pParaHead = &pair
			pParaTail = &pair
		} else {
			(*pParaTail).next = &pair
			pParaTail = (*pParaTail).next
		}
	}

	// 非会话数据;
	var pDataHead, pDataTail C.pDataList
	for _, ele := range req.DeliverData {
		var reqDataC C.struct_DataList
		reqDataC._type = C.DataType(ele.DataType)
		reqDataC.status = C.DataStatus(ele.DataStatus)
		reqDataC.encoding = C.CString(ele.DataFmt)
		defer C.free(unsafe.Pointer(reqDataC.encoding))
		reqDataC.desc = C.CString(ele.DataDesc)
		defer C.free(unsafe.Pointer(reqDataC.desc))
		reqDataC.data = C.CBytes(ele.Data)
		defer C.free(unsafe.Pointer(reqDataC.data))
		reqDataC.len = C.uint(len(ele.Data))
		reqDataC.next = nil

		if pDataHead == nil {
			pDataHead = &reqDataC
			pDataTail = &reqDataC
		} else {
			(*pDataTail).next = &reqDataC
			pDataTail = (*pDataTail).next
		}
	}

	// 处理函数：exec() & execAsync()
	if conf.WrapperAsync {
		usrTag := C.CString(handle)
		//defer C.free(unsafe.Pointer(usrTag)) callBack free
		callback := C.wrapperCallback(C.adapterCallback)
		errC := C.adapterExecAsync(unsafe.Pointer(usrTag), pParaHead, pDataHead, callback, C.int(0))
		if errC != 0 {
			errNum = int(errC)
			errInfo = errors.New(engineError(int(errC)))
			return
		}
	} else {
		var respDataC C.pDataList
		errC := C.adapterExec(pParaHead, pDataHead, &respDataC)
		if errC != 0 {
			errNum = int(errC)
			errInfo = errors.New(engineError(int(errC)))
			return
		}

		// 输出拷贝&转换
		tmpDataPtr := respDataC
		resp.DeliverData = make([]instance.DataMeta, 0, 1)
		for tmpDataPtr != nil {
			var ele instance.DataMeta
			ele.DataType = int((*tmpDataPtr)._type)
			ele.DataStatus = int((*tmpDataPtr).status)
			ele.DataDesc = C.GoString((*tmpDataPtr).desc)
			ele.DataFmt = C.GoString((*tmpDataPtr).encoding)
			ele.Data = C.GoBytes(unsafe.Pointer((*tmpDataPtr).data), C.int((*tmpDataPtr).len))
			resp.DeliverData = append(resp.DeliverData, ele)
			tmpDataPtr = (*tmpDataPtr).next
		}

		// tmp数据释放：execFree()
		errC = C.adapterExecFree(&respDataC)
	}
	return
}

// 计算debug数据
func engineDebug(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	debug := C.adapterDebugInfo(req.WrapperHdl.(*C.char))
	resp.Debug = C.GoString(debug)
	return
}

func engineError(errNum int) (errInfo string) {
	err := C.adapterError(C.int(errNum))
	return C.GoString(err)
}
