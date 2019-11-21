#include "wrapper.h"
#include "stdio.h"
#include "string.h"
#include "stdlib.h"

/*
 * 用于实现wrapper demo的mock版本
 * 目前仅支持同步插件接口,非实时结果返回模式的mock
 * */

const void* wrapperInnerHdl = "wrapperTestHandle";
pDataList wrapperInnerRslt = NULL;
pDataList wrapperOnceRslt = NULL;

int WrapperAPI wrapperInit(pConfig cfg){
	printf("wrapperInit\n");
    
    // 打印输入配置项
    while (cfg != NULL) {
        printf("key=%s, value=%s\n", cfg->key, cfg->value);
        cfg = cfg->next;
    }

    // 构建内部测试read值
    wrapperInnerRslt = (struct DataList*)malloc(sizeof(struct DataList));
    wrapperInnerRslt->key = (char*)"result";
    wrapperInnerRslt->data = (void*)"response result from AtheLoader wrapper";
    wrapperInnerRslt->len = strlen("response result from AtheLoader wrapper");
    wrapperInnerRslt->desc = NULL;
    wrapperInnerRslt->encoding = (char*)"utf8";
    wrapperInnerRslt->type = DataText;
    wrapperInnerRslt->status = DataEnd;
    wrapperInnerRslt->next = NULL;

    wrapperOnceRslt = (struct DataList*)malloc(sizeof(struct DataList));
    wrapperOnceRslt->key = (char*)"result";
    wrapperOnceRslt->data = (void*)"response result from AtheLoader wrapper";
    wrapperOnceRslt->len = strlen("response result from AtheLoader wrapper");
    wrapperOnceRslt->desc = NULL;
    wrapperOnceRslt->encoding = (char*)"utf8";
    wrapperOnceRslt->type = DataText;
    wrapperOnceRslt->status = DataOnce;
    wrapperOnceRslt->next = NULL;
    return 0;
}

int WrapperAPI wrapperFini(){
	printf("wrapperFini\n");
	return 0;
}

const char* WrapperAPI wrapperError(int errNum){
    printf("wrapperError\n");	
    return "inner error";
}

const char* WrapperAPI wrapperVersion(){
	printf("wrapperVersion\n");
    return "1.0.0";
}

int WrapperAPI wrapperLoadRes(pDataList perData, unsigned int resId){
	printf("wrapperLoadRes\n");
	return 0;
}

int WrapperAPI wrapperUnloadRes(unsigned int resId){
	printf("wrapperUnloadRes\n");
	return 0;
}

const void* WrapperAPI wrapperCreate(const void* usrTag, pParamList params, wrapperCallback cb, unsigned int psrIds[], int psrCnt, int* errNum){
	printf("wrapperCreate\n");
	return wrapperInnerHdl;
}

int WrapperAPI wrapperWrite(const void* handle, pDataList reqData){
	printf("wrapperWrite\n");
	return 0;
}

int WrapperAPI wrapperRead(const void* handle, pDataList* respData){
	printf("wrapperRead\n");
    /*
Note: demo 关闭实时读功能,wrapperRead仅需返回最终last结果,
    若需实时返回结果,可通过修改配置[aiges].realTimeRlt = true
     */
    respData = &wrapperInnerRslt; 
	return 0;
}

int WrapperAPI wrapperDestroy(const void* handle){
	printf("wrapperDestroy\n");
	return 0;
}

int WrapperAPI wrapperExec(pParamList params, pDataList reqData, pDataList* respData){
	printf("wrapperExec\n");
	respData = &wrapperOnceRslt;
    return 0;
}

int WrapperAPI wrapperExecFree(pDataList* respData){
	printf("wrapperExecFree\n");
	return 0;
}

int WrapperAPI wrapperExecAsync(const void* usrTag, pParamList params, pDataList reqData, wrapperCallback callback, int timeout){
	printf("wrapperExecAsync\n");
	return 0;
}

const char* WrapperAPI wrapperDebugInfo(const void* handle){
	printf("wrapperDebugInfo\n");
	return "DebugInfo";
}
