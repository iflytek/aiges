package widget

/*
#cgo linux CFLAGS: -I../cgo/header -I/usr/include/python2.7 -Wno-attributes -std=c99
#cgo linux LDFLAGS: -ldl -lpython2.7
#include <stdlib.h>
#include <stdio.h>
#include <Python.h>
#include <string.h>
#include "./widget/type.h"

// python module object
PyObject *pWrapperModule = NULL;
PyObject *pEngineInst = NULL;	// 类实例
PyObject *pEngineClass = NULL;	// "EngineBase"
PyObject *pFuncInit = NULL;		// "wrapperInit"
PyObject *pFuncFini = NULL;		// "wrapperFini"
PyObject *pFuncErr = NULL;		// "wrapperError"
PyObject *pFuncExec = NULL;		// "wrapperExec"
PyObject *pFuncDebug = NULL;	// "wrapperDebugInfo"
PyObject *pStrVer = NULL;		// "version"

// python adapter error code
const int PyErrLoadModule = 10001;
const int PyErrLoadEngine = 10002;
const int PyErrNewEngine = 10003;
const int PyErrLoadFunc = 10004;
const int PyErrFuncCallable = 10005;
const int PyErrValConvert = 10006;
const int PyErrFuncExec = 10007;

char* version = NULL;
pDataList tempResult = NULL;

// python初始化
int pyInit(){
	// Py_SetProgramName(argv[0]);	// optional but recommended
	Py_Initialize();
	// check
	if (Py_IsInitialized() == 0){
		return -1;
	}

	// GIL check
	PyEval_InitThreads();
	if (PyEval_ThreadsInitialized() == 0) {
		return -1;
	}

	// 切换python工作目录
	PyRun_SimpleString("import sys");
	PyRun_SimpleString("sys.path.append('./')");
	return 0;
}

// python逆初始化
int pyFini(){
	Py_XDECREF(pWrapperModule);
	Py_XDECREF(pEngineClass);
	Py_XDECREF(pEngineInst);
	Py_XDECREF(pFuncInit);
	Py_XDECREF(pFuncFini);
	Py_XDECREF(pFuncErr);
	Py_XDECREF(pFuncExec);
	Py_XDECREF(pFuncDebug);
	Py_XDECREF(pStrVer);
	Py_Finalize();
	return 0;
}

int pyModuleLoad(const char* module){
	// load python模块&函数
	pWrapperModule = PyImport_ImportModule(module);
	if (!pWrapperModule){
		PyErr_Print();
		return PyErrLoadModule;
	}

	PyObject* pDict = PyModule_GetDict(pWrapperModule);
	pEngineClass = PyDict_GetItemString(pDict, "EngineBase");
	if (!pEngineClass){
		printf("pyModuleLoad get EngineBase fail:");
		PyErr_Print();
		return PyErrLoadEngine;
	}
	pEngineInst = PyInstance_New(pEngineClass, NULL, NULL);
	if (!pEngineInst){
		printf("pyModuleLoad Instance_New fail:");
		PyErr_Print();
		return PyErrNewEngine;
	}
	pStrVer = PyObject_GetAttrString(pEngineClass, "version");	// TODO check 属性
	if (!pStrVer){
		printf("pyModuleLoad get version fail:");
		PyErr_Print();
		return PyErrNewEngine;
	}
	return 0;
}

const char* pyVersion(){
	if (version == NULL && pStrVer != NULL){
		version = PyString_AsString(pStrVer);
		// PyArg_Parse(pStrVer, "s", &version);
	}
	return version;
}

int pyEngineInit(pConfig cfg){
	// build python args
	PyObject* pcfg = PyDict_New();
	while (cfg != NULL){
		PyDict_SetItemString(pcfg, cfg->key, Py_BuildValue("s#", cfg->value, cfg->vlen));
		cfg = cfg->next;
	}
	//
	//PyObject* pArgsT = PyTuple_New(1);
	//PyTuple_SetItem(pArgsT, 0, pcfg);
	PyObject* pRet = PyObject_CallMethod(pEngineInst, "wrapperInit", "O", pcfg);
	if (!pRet){
		printf("pyEngineInit fail:");
		PyErr_Print();
		return PyErrFuncExec;
	}
	int ret = (int)PyInt_AsLong(pRet);
	Py_XDECREF(pRet);
	Py_XDECREF(pcfg);
	return ret;
}

int pyEngineFini(){
	PyObject* pRet = PyObject_CallMethod(pEngineInst,  "wrapperFini", NULL, NULL);
	if (!pRet){
		printf("pyEngineFini fail:");
		PyErr_Print();
		return PyErrFuncExec;
	}
	int ret = PyInt_AsLong(pRet);
	Py_XDECREF(pRet);
	return ret;
}

int pyEngineExec(pParamList params, pDataList reqData, pDataList* respData){
	// build python params
	PyObject* pPara = PyDict_New();
	while (params != NULL){
		PyDict_SetItemString(pPara, params->key, Py_BuildValue("s", params->value));
		params = params->next;
	}

	// build python inputData
	PyObject* pData = NULL;
	while (reqData != NULL){
		PyObject* pTemp = PyDict_New();
		PyDict_SetItemString(pTemp, "Key", Py_BuildValue("s", reqData->key));
		PyDict_SetItemString(pTemp, "Data", Py_BuildValue("s#", reqData->data, reqData->len));
		PyDict_SetItemString(pTemp, "Type", Py_BuildValue("i", reqData->type));
		PyDict_SetItemString(pTemp, "Status", Py_BuildValue("i", reqData->status));
		PyDict_SetItemString(pTemp, "Desc", Py_BuildValue("s", reqData->desc));
		PyDict_SetItemString(pTemp, "Encoding", Py_BuildValue("s", reqData->encoding));
		if (pData == NULL){
			pData = PyList_New(1);
			PyList_SetItem(pData, 0, pTemp);
		}else {
			PyList_Append(pData, pTemp);
		}
		// Py_XDECREF(pTemp);  TODO 释放临时资源
		reqData = reqData->next;
	}

	PyObject* pTuple = PyObject_CallMethod(pEngineInst, "wrapperExec", "OO", pPara, pData);
	if (!pTuple){
		printf("pyEngineExec fail:");
		PyErr_Print();
		return PyErrFuncExec;
	}
	// TODO PyTuple_Check() & PyTuple_Size()
	PyObject* pRet = PyTuple_GetItem(pTuple, 0); // error code;
	int ret = PyInt_AsLong(pRet);
	if (ret == 0){
		PyObject* pResp = PyTuple_GetItem(pTuple, 1); // calc result;
		int respNum = PyList_Size(pResp);
		// free python->c result tmp cache;
		while (tempResult != NULL){
			pDataList temp = tempResult;
			tempResult = tempResult->next;
			free(temp);
		}
		pDataList* cur = &tempResult;
		for (int i = 0; i < respNum; i++){
			*cur = (struct DataList*)malloc(sizeof(struct DataList));
			PyObject* pDict = PyList_GetItem(pResp, i);
			PyObject* pKeys = PyDict_Keys(pDict);
			for(int i = 0; i < PyList_Size(pKeys); ++i) {
				PyObject *pKey = PyList_GetItem(pKeys, i);	// TODO fix
				PyObject *pValue = PyDict_GetItem(pDict, pKey);
				char* key = PyString_AsString(pKey);
				if (strcmp(key, "KEY") == 0){
					(*cur)->key = PyString_AsString(pValue);
				}else if (strcmp(key, "Type") == 0){
					(*cur)->type = PyInt_AsLong(pValue);
				}else if (strcmp(key, "Status") == 0){
					(*cur)->status = PyInt_AsLong(pValue);
				}else if (strcmp(key, "Desc") == 0){
					(*cur)->desc = PyString_AsString(pValue);
				}else if (strcmp(key, "Encoding") == 0){
					(*cur)->encoding = PyString_AsString(pValue);
				}else if (strcmp(key, "Data") == 0){
					(*cur)->data = PyString_AsString(pValue);
				}
			}
			(*cur)->len = strlen((*cur)->data);
			cur = &((*cur)->next);
		}
		respData = &tempResult;

	}
	return ret;
}

char* pyEngineErr(int err){
	PyObject* pInt = PyInt_FromLong((long)err);
	PyObject* pErr = PyObject_CallMethod(pEngineInst, "wrapperError", "i", pInt);
	if (!pErr){
		printf("pyEngineErr fail:");
		PyErr_Print();
		return "pyEngineErr exception";
	}
	char* desc = PyString_AsString(pErr);
	Py_XDECREF(pInt);
	Py_XDECREF(pErr);
	return desc;
}

*/
import "C"
import (
	"errors"
	"instance"
	"strconv"
	"strings"
	"unsafe"
)

/* TODO c/c++ 调用python GIL问题
PyGILState_STATE gstate;
gstate = PyGILState_Ensure();

// Perform Python actions here.
result = CallSomeFunction();
// evaluate result or handle exception

// Release the thread. No Python API allowed beyond this point.
PyGILState_Release(gstate);
*/

func pythonOpen(file string) (err error) {
	ret := C.pyInit()
	if ret != 0 {
		return errors.New("python init fail")
	}

	if !strings.HasSuffix(file, ".py") {
		return errors.New("python init, open invalid python file " + file)
	}
	mName := strings.TrimSuffix(file, ".py")
	cName := C.CString(mName)
	defer C.free(unsafe.Pointer(cName))

	ret = C.pyModuleLoad(cName)
	if ret != 0 {
		return errors.New("python init load fail," + ", name: " + mName + ", err:" + strconv.Itoa(int(ret)))
	}

	return nil
}

func pythonClose() {
	C.pyFini()
}

func pythonInit(cfg map[string]string) (errNum int, errInfo error) {
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
	ret := C.pyEngineInit(pCfgHead)
	if ret != 0 {
		errInfo = errors.New(pythonError(int(ret)))
		errNum = int(ret)
	}
	return
}

func pythonFini() (errNum int, errInfo error) {
	ret := C.pyEngineFini()
	if ret != 0 {
		errInfo = errors.New(pythonError(int(ret)))
		errNum = int(ret)
	}
	return
}

func pythonOnceExec(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
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

	// 处理函数：exec() // TODO lock
	var respDataC C.pDataList
	errC := C.pyEngineExec(pParaHead, pDataHead, &respDataC)
	if errC != 0 {
		errNum = int(errC)
		errInfo = errors.New(pythonError(int(errC)))
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
	// TODO unlock (data copy)
	return
}

func pythonVersion() (ver string) {
	verC := C.pyVersion()
	return C.GoString(verC)
}

func pythonError(code int) (err string) {
	errC := C.pyEngineErr(C.int(code))
	return C.GoString(errC)
}

// TODO 后续补充实现
// 资源加载卸载管理适配接口;
func pythonLoadRes(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {

	return
}

func pythonUnloadRes(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {

	return
}

// 资源申请行为
func pythonCreate(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {

	return
}

// 资源释放行为
func pythonDestroy(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {

	return
}

// 交互异常行为
func pythonExcp(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {

	return
}

// 数据写行为
func pythonWrite(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {

	return
}

// 数据读行为
func pythonRead(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {

	return
}

// 计算debug数据
func pythonDebug(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	return
}
