//go:build linux
// +build linux

package widget

/*
#cgo linux CFLAGS: -I../cgo/header -Wno-attributes
#cgo linux LDFLAGS: -ldl
#include <stdlib.h>
#include <dlfcn.h>
#include <stdio.h>
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
const char* SetCtrlSym = "wrapperSetCtrl";

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
wrapperSetCtrlPtr	cPtrSetCtrl;

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
	if ((cPtrSetCtrl = cLibLoad(hdl, SetCtrlSym, loadErr)) == NULL)
		return SetCtrlSym;
	return NULL;
}

// 接口适配封装
int adapterSetCtrl(CtrlType type, void* func){
	return (*cPtrSetCtrl)(type, func);
}

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


const void* adapterCreate(const char* usrTag, pParamList params, wrapperCallback cb, unsigned int psrIds[], int psrCnt, int* errNum){
	return (*cPtrCreate)(usrTag, params, cb, psrIds, psrCnt, errNum);
}


struct ParamList* paramListCreate(){
	struct ParamList* pParamPtr= (struct ParamList*)malloc(sizeof(struct ParamList));
	pParamPtr->key=NULL;
	pParamPtr->value=NULL;
	pParamPtr->vlen=0;
	pParamPtr->next=NULL;
	return pParamPtr;
}

struct ParamList* paramListAppend(struct ParamList* ptr,char* key,char* value,unsigned int vlen){
	struct ParamList* head=ptr;
	if(ptr!=NULL && ptr->key==NULL){
		ptr->key=key;
		ptr->value=value;
		ptr->vlen=vlen;
		return ptr;
	}
	while(ptr->next!=NULL){
		ptr=ptr->next;
	}
	struct ParamList* pParamPtr= (struct ParamList*)malloc(sizeof(struct ParamList));
	pParamPtr->key=key;
	pParamPtr->value=value;
	pParamPtr->vlen=vlen;
	pParamPtr->next=NULL;
	ptr->next=pParamPtr;
	return head;
}

void paramListfree(struct ParamList* ptr){
	struct ParamList* current=ptr;
	while(current!=NULL){
		ptr=ptr->next;
		free(current);
		current=ptr;
	}
	return ;
}

struct DataList* dataListCreate(){
	struct DataList* pDataPtr= (struct DataList*)malloc(sizeof(struct DataList));
	pDataPtr->key=NULL;
	pDataPtr->data=NULL;
	pDataPtr->len=0;
	pDataPtr->type=DataText;
	pDataPtr->status=DataBegin;
	pDataPtr->desc=NULL;
	pDataPtr->next=NULL;
	return pDataPtr;
}

struct DataList* dataListAppend(struct DataList* ptr,char* key,void* data,unsigned int len,int type,int status,struct ParamList* descPtr){
	struct DataList* head=ptr;
	if(ptr!=NULL && ptr->key==NULL){
		ptr->key=key;
		ptr->data=data;
		ptr->len=len;
		ptr->type=(DataType)(type);
		ptr->status=(DataStatus)(status);
		ptr->desc=descPtr;
		ptr->next=NULL;
		return ptr;
	}
	while(ptr->next!=NULL){
		ptr=ptr->next;
	}
	struct DataList* pDataPtr= (struct DataList*)malloc(sizeof(struct DataList));
	pDataPtr->key=key;
	pDataPtr->data=data;
	pDataPtr->len=len;
	pDataPtr->type=(DataType)(type);
	pDataPtr->status=(DataStatus)(status);
	pDataPtr->desc=descPtr;
	pDataPtr->next=NULL;
	ptr->next=pDataPtr;

	return head;
}

void DataListfree(struct DataList* ptr){
	struct DataList* currentData=ptr;
	while(currentData!=NULL){
		ptr=ptr->next;
		free(currentData);
		currentData=ptr;
	}
	return ;
}


int adapterWrite(const void* handle, pDataList reqData){
	return (*cPtrWrite)(handle, reqData);
}

int adapterRead(const void* handle, pDataList* respData){
	return (*cPtrRead)(handle, respData);
}

pDataList* getDataList(){

	pDataList* ptr;
	ptr=(pDataList*)malloc(sizeof(pDataList));
	memset(ptr,0,sizeof(pDataList));
	return ptr;
}
void releaseDataList(pDataList*  pDataListPtr){
	free((pDataList*)(pDataListPtr));
}

int adapterDestroy(const void* handle){
	return (*cPtrDestroy)(handle);
}

int adapterExec(const char* usrTag, pParamList params, pDataList reqData, pDataList* respData, unsigned int psrIds[], int psrCnt){
	return (*cPtrExec)(usrTag, params, reqData, respData, psrIds, psrCnt);
}

int adapterExecFree(const char* usrTag, pDataList* respData){
	return (*cPtrExecFree)(usrTag, respData);
}

int adapterExecAsync(const char* usrTag, pParamList params, pDataList reqData, wrapperCallback callback, int timeout, unsigned int psrIds[], int psrCnt){
	return (*cPtrExecAsync)(usrTag, params, reqData, callback, timeout, psrIds, psrCnt);
}

const char* adapterDebugInfo(const void* handle){
	return (*cPtrDebugInfo)(handle);
}

extern int engineCreateCallBack(void* handle, pDataList respData, int ret);
int adapterCallback(const void* handle, pDataList respData, int ret){
	return engineCreateCallBack((void*)handle, respData, ret);
}


int adapterFreeParaList(pParamList pl){
	for (;pl != NULL;){
		pParamList tmp = pl;
		pl = pl->next;
		free(tmp->key);
		tmp->key = NULL;
		free(tmp->value);
		tmp->value = NULL;
		free(tmp);
		tmp = NULL;
	}
	return 0;
}

int adapterFreeDataList(pDataList dl){
	for (;dl != NULL;) {
		pDataList tmp = dl;
		dl = dl->next;
		free(tmp->key);
		tmp->key = NULL;
		free(tmp->data);
		tmp->data = NULL;
		pDescList desc = tmp->desc;
		for (;desc != NULL;) {
			pDescList tmp1 = desc;
			desc = desc->next;
			free(tmp1->key);
			tmp1->key = NULL;
			free(tmp1->value);
			tmp1->value = NULL;
			free(tmp1);
			tmp1 = NULL;
		}
		tmp->desc = NULL;
		free(tmp);
		tmp = NULL;
	}
	return 0;
}

*/
import "C"
import (
	"errors"
	"github.com/xfyun/aiges/catch"
	"github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/instance"
	"unsafe"
)

///	wrapper适配器,提供golang至c/c++ wrapper层的数据适配及接口转换;
type engineC struct {
	libHdl unsafe.Pointer // 引擎句柄;
}

func (ec *engineC) open(libName string) (errInfo error) {
	var errC *C.char
	var libNameC *C.char = C.CString(libName)
	ec.libHdl = C.cLibOpen(libNameC, &errC)
	if ec.libHdl == nil {
		C.free(unsafe.Pointer(libNameC))
		errInfo = errors.New("wrapper.open: load library " + libName + "failed, " + C.GoString(errC))
		return
	}

	errSym := C.loadWrapperSym(ec.libHdl, &errC)
	if errSym != nil {
		C.cLibClose(ec.libHdl)
		C.free(unsafe.Pointer(libNameC))
		errInfo = errors.New("wrapper.open: load symbol " + C.GoString(errSym) + " failed, " + C.GoString(errC))
		return
	}
	C.free(unsafe.Pointer(libNameC))
	return
}

func (ec *engineC) close() {
	C.cLibClose(ec.libHdl)
	ec.libHdl = nil
	return
}

func engineInit(cfg map[string]string) (errNum int, errInfo error) {
	// 配置参数的语言栈转换;
	configList := C.paramListCreate()
	defer C.paramListfree(configList)
	for k, v := range cfg {
		key := C.CString(k)
		defer C.free(unsafe.Pointer(key))
		val := C.CString(v)
		defer C.free(unsafe.Pointer(val))
		valLen := C.uint(len(v))
		configList = C.paramListAppend(configList, key, val, valLen)
	}

	ret := C.adapterInit(configList)
	if ret != 0 {
		errInfo = errors.New(engineError(int(ret)))
		errNum = int(ret)
		handle := new(catch.StartFailedHandle)
		handle.Occur()
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
	uuid := catch.GenerateUuid()
	catch.CallCgo(uuid, catch.Begin)
	defer catch.CallCgo(uuid, catch.End)

	dataList := C.dataListCreate()
	defer C.DataListfree(dataList)

	descList := C.paramListCreate()
	defer C.paramListfree(descList)
	for k, v := range req.PsrDesc {
		key := C.CString(k)
		defer C.free(unsafe.Pointer(key))
		val := C.CString(v)
		defer C.free(unsafe.Pointer(val))
		valLen := C.uint(len(v))
		descList = C.paramListAppend(descList, key, val, valLen)
	}
	tmpKey := C.CString(req.PsrKey)
	defer C.free(unsafe.Pointer(tmpKey))
	tmpData := C.CBytes(req.PsrData)
	defer C.free(unsafe.Pointer(tmpData))
	//4 个性化数据 2传完
	dataList = C.dataListAppend(dataList, tmpKey, tmpData, C.uint(len(req.PsrData)), C.int(4), C.int(2), descList)

	errC := C.adapterLoadRes(dataList, C.uint(req.PsrId))
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
	}
	return
}

func engineUnloadRes(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	uuid := catch.GenerateUuid()
	catch.CallCgo(uuid, catch.Begin)
	defer catch.CallCgo(uuid, catch.End)
	errC := C.adapterUnloadRes(C.uint(req.PsrId))
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
	}
	return
}

// 资源申请行为
func engineCreate(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	uuid := catch.GenerateUuid()
	catch.CallCgo(uuid, catch.Begin)
	defer catch.CallCgo(uuid, catch.End)
	// 参数对;

	paramList := C.paramListCreate()
	defer C.paramListfree(paramList)
	for k, v := range req.Params {
		key := C.CString(k)
		defer C.free(unsafe.Pointer(key))
		val := C.CString(v)
		defer C.free(unsafe.Pointer(val))
		valLen := C.uint(len(v))
		paramList = C.paramListAppend(paramList, key, val, valLen)
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
	engHdl := C.adapterCreate(usrTag, paramList, callback, psrPtr, C.int(psrCnt), &errC)
	if errC != 0 || engHdl == nil {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
	} else {
		resp.WrapperHdl = engHdl
	}

	//	C.adapterFreeParaList(pParaHead)

	if !conf.WrapperAsync {
		C.free(unsafe.Pointer(usrTag))
	}
	return
}

// 资源释放行为
func engineDestroy(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	uuid := catch.GenerateUuid()
	catch.CallCgo(uuid, catch.Begin)
	defer catch.CallCgo(uuid, catch.End)
	errC := C.adapterDestroy(req.WrapperHdl)
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
	}
	return
}

// 交互异常行为
func engineExcp(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	uuid := catch.GenerateUuid()
	catch.CallCgo(uuid, catch.Begin)
	defer catch.CallCgo(uuid, catch.End)
	resp, errNum, errInfo = engineDestroy(handle, req)
	return
}

// 数据写行为
func engineWrite(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	uuid := catch.GenerateUuid()
	catch.CallCgo(uuid, catch.Begin)
	defer catch.CallCgo(uuid, catch.End)
	// 写数据转换;
	dataList := C.dataListCreate()
	defer C.DataListfree(dataList)
	for _, ele := range req.DeliverData {
		tmpKey := C.CString(ele.DataId)
		defer C.free(unsafe.Pointer(tmpKey))
		tmpData := C.CBytes(ele.Data)
		defer C.free(unsafe.Pointer(tmpData))

		descList := C.paramListCreate()
		defer C.paramListfree(descList)
		for k, v := range ele.DataDesc {
			key := C.CString(k)
			defer C.free(unsafe.Pointer(key))
			val := C.CString(v)
			defer C.free(unsafe.Pointer(val))
			valLen := C.uint(len(v))
			descList = C.paramListAppend(descList, key, val, valLen)
		}
		dataList = C.dataListAppend(dataList, tmpKey, tmpData, C.uint(len(ele.Data)), C.int(ele.DataType), C.int(ele.DataStatus), descList)
	}

	errC := C.adapterWrite(req.WrapperHdl, dataList)
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
	}

	//	C.adapterFreeDataList(pDataHead)
	return
}

// 数据读行为
func engineRead(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	uuid := catch.GenerateUuid()
	catch.CallCgo(uuid, catch.Begin)
	defer catch.CallCgo(uuid, catch.End)
	respDataC := C.getDataList()
	defer C.releaseDataList(respDataC)
	errC := C.adapterRead(req.WrapperHdl, respDataC)
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(engineError(int(errC)))
		return
	}

	// 读数据转换;
	//	resp.DeliverData = make([]instance.DataMeta, 0, 1)
	for *respDataC != nil {
		var ele instance.DataMeta
		ele.DataId = C.GoString((*(*respDataC)).key)
		ele.DataType = int((*(*respDataC))._type)
		ele.DataStatus = int((*(*respDataC)).status)
		ele.DataDesc = make(map[string]string)
		pDesc := (*(*respDataC)).desc
		for pDesc != nil {
			ele.DataDesc[C.GoString((*pDesc).key)] = C.GoStringN((*pDesc).value, C.int((*pDesc).vlen))
			pDesc = (*pDesc).next
		}
		if int((*(*respDataC)).len) != 0 {
			ele.Data = C.GoBytes((*(*respDataC)).data, C.int((*(*respDataC)).len))
		}
		resp.DeliverData = append(resp.DeliverData, ele)
		*respDataC = (*respDataC).next
	}

	return
}

func engineOnceExec(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	uuid := catch.GenerateUuid()
	catch.CallCgo(uuid, catch.Begin)
	defer catch.CallCgo(uuid, catch.End)
	// 非会话参数;
	paramList := C.paramListCreate()
	defer C.paramListfree(paramList)
	for k, v := range req.Params {
		key := C.CString(k)
		defer C.free(unsafe.Pointer(key))
		val := C.CString(v)
		defer C.free(unsafe.Pointer(val))
		valLen := C.uint(len(v))
		paramList = C.paramListAppend(paramList, key, val, valLen)
	}

	// 非会话数据;
	dataList := C.dataListCreate()
	defer C.DataListfree(dataList)
	for _, ele := range req.DeliverData {
		tmpKey := C.CString(ele.DataId)
		defer C.free(unsafe.Pointer(tmpKey))
		tmpData := C.CBytes(ele.Data)
		defer C.free(unsafe.Pointer(tmpData))

		descList := C.paramListCreate()
		defer C.paramListfree(descList)
		for k, v := range ele.DataDesc {
			key := C.CString(k)
			defer C.free(unsafe.Pointer(key))
			val := C.CString(v)
			defer C.free(unsafe.Pointer(val))
			valLen := C.uint(len(v))
			descList = C.paramListAppend(descList, key, val, valLen)
		}
		dataList = C.dataListAppend(dataList, tmpKey, tmpData, C.uint(len(ele.Data)), C.int(ele.DataType), C.int(ele.DataStatus), descList)
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

	// 处理函数：exec() & execAsync()
	if conf.WrapperAsync {
		usrTag := C.CString(handle)
		//defer C.free(unsafe.Pointer(usrTag)) callBack free
		callback := C.wrapperCallback(C.adapterCallback)
		errC := C.adapterExecAsync(usrTag, paramList, dataList, callback, C.int(0), psrPtr, C.int(psrCnt))
		if errC != 0 {
			errNum = int(errC)
			errInfo = errors.New(engineError(int(errC)))
		}
	} else {
		respDataC := C.getDataList()
		defer C.releaseDataList(respDataC)

		tag := C.CString(handle)
		errC := C.adapterExec(tag, paramList, dataList, respDataC, psrPtr, C.int(psrCnt))
		C.free(unsafe.Pointer(tag))
		if errC != 0 {
			errNum = int(errC)
			errInfo = errors.New(engineError(int(errC)))
		} else {
			// 输出拷贝&转换
			tmpDataPtr := *respDataC
			resp.DeliverData = make([]instance.DataMeta, 0, 1)
			for *respDataC != nil {
				var ele instance.DataMeta
				ele.DataId = C.GoString((*(*respDataC)).key)
				ele.DataType = int((*(*respDataC))._type)
				ele.DataStatus = int((*(*respDataC)).status)
				ele.DataDesc = make(map[string]string)
				pDesc := (*(*respDataC)).desc
				for pDesc != nil {
					ele.DataDesc[C.GoString((*pDesc).key)] = C.GoStringN((*pDesc).value, C.int((*pDesc).vlen))
					pDesc = (*pDesc).next
				}
				if int((*(*respDataC)).len) != 0 {
					ele.Data = C.GoBytes((*(*respDataC)).data, C.int((*(*respDataC)).len))
				}
				resp.DeliverData = append(resp.DeliverData, ele)
				*respDataC = (*respDataC).next
			}
			*respDataC = tmpDataPtr
			// tmp数据释放：execFree()
			errC = C.adapterExecFree(tag, respDataC)
		}
	}
	return
}

// 计算debug数据
func engineDebug(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	uuid := catch.GenerateUuid()
	catch.CallCgo(uuid, catch.Begin)
	defer catch.CallCgo(uuid, catch.End)
	debug := C.adapterDebugInfo(req.WrapperHdl)
	resp.Debug = C.GoString(debug)
	return
}

func engineError(errNum int) (errInfo string) {
	err := C.adapterError(C.int(errNum))
	return C.GoString(err)
}
