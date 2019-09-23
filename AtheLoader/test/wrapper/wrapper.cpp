#include "wrapper.h"
#include "stdio.h"

int WrapperAPI wrapperInit(pConfig cfg){
	printf("wrapperInit::debug\n");
	return 0;
}

int WrapperAPI wrapperFini(){
	printf("wrapperFini::debug\n");
	return 0;
}

const char* WrapperAPI wrapperError(int errNum){
	printf("wrapperError::debug, errNum %d\n", errNum);
	return "wrapperError::debug";
}

const char* WrapperAPI wrapperVersion(){
	printf("wrapperVersion::debug\n");
	return "1.0.1";
}

int WrapperAPI wrapperLoadRes(pDataList perData, unsigned int resId){
	printf("wrapperLoadRes::debug\n");
	return 0;
}

int WrapperAPI wrapperUnloadRes(unsigned int resId){
	printf("wrapperUnloadRes::debug\n");
	return 0;
}

const char* WrapperAPI wrapperCreate(pParamList params, wrapperCallback cb, unsigned int psrIds[], int psrCnt, int* errNum){
	printf("wrapperCreate::debug\n");
	return "wrapperCreate::debug";
}

int WrapperAPI wrapperWrite(const char* handle, pDataList reqData){
	printf("wrapperWrite::debug\n");
	return 0;
}

int WrapperAPI wrapperRead(const char* handle, pDataList* respData){
	printf("wrapperRead::debug\n");
	return 0;
}

int WrapperAPI wrapperDestroy(const char* handle){
	printf("wrapperDestroy::debug\n");
	return 0;
}

int WrapperAPI wrapperExec(pParamList params, pDataList reqData, pDataList* respData){
	printf("wrapperExec::debug\n");
	return 0;
}

int WrapperAPI wrapperExecFree(pDataList* respData){
	printf("wrapperExecFree::debug\n");
	return 0;
}

int WrapperAPI wrapperExecAsync(const char* handle, pParamList params, pDataList reqData, wrapperCallback callback, int timeout){
	printf("wrapperExecAsync::debug\n");
	return 0;
}

const char* WrapperAPI wrapperDebugInfo(const char* handle){
	printf("wrapperDebugInfo\n");
	return "wrapperDebugInfo";
}
