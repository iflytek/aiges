#include "../cgo/header/widget/wrapper.h"
#include "stdio.h"
#include "string.h"
#include "stdlib.h"
#include<iostream>
#include<fstream>
#include <pthread.h>
#include <map>
#include <thread>
#include <mutex>
#include<string>
#include<unistd.h>
/*
 * 用于实现wrapper demo的mock版本
 * 目前仅支持同步插件接口,非实时结果返回模式的mock
 * */
std::mutex my_mutex;
const void* wrapperInnerHdl = "wrapperTestHandle";
const char* defResult = "response result from AtheLoader wrapper";
pDataList wrapperInnerRslt = NULL;
pDataList wrapperOnceRslt = NULL;
std::map< int* /*handle*/, int /*status*/> RsltStatus;
pthread_mutex_t  idle_mutex_;
//std::ofstream g_log("debug.log");
int KJINDEX=0;
int reqCount=0;

int WrapperAPI wrapperSetCtrl(CtrlType type, void* func){
	return 0;
}

int WrapperAPI wrapperInit(pConfig cfg){
	printf("wrapperInit\n");
    //std::cout.rdbuf(g_log.rdbuf());
    // 打印输入配置项
    while (cfg != NULL) {
        printf("key=%s, value=%s\n", cfg->key, cfg->value);
        if(std::string(cfg->key)=="initFail" && std::string(cfg->value)=="1"){
            std::cout<<"init failed"<<std::endl;
            return -1;
        }
        cfg = cfg->next;
    }

    // 构建内部测试read值
    wrapperInnerRslt = (struct DataList*)malloc(sizeof(struct DataList));
    wrapperInnerRslt->key = (char*)"result";
    wrapperInnerRslt->data = (void*)defResult;
    wrapperInnerRslt->len = strlen(defResult);
    wrapperInnerRslt->desc = NULL;
    wrapperInnerRslt->type = DataText;
    wrapperInnerRslt->status = DataEnd;
    wrapperInnerRslt->next = NULL;

    wrapperOnceRslt = (struct DataList*)malloc(sizeof(struct DataList));
    wrapperOnceRslt->key = (char*)"result";
    wrapperOnceRslt->data = (void*)defResult;
    wrapperOnceRslt->len = strlen(defResult);
    wrapperOnceRslt->desc = NULL;
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
//  printf("wrapperError\n");
    return "inner error";
}

const char* WrapperAPI wrapperVersion(){
//	printf("wrapperVersion\n");
    return "1.0.0";
}

int WrapperAPI wrapperLoadRes(pDataList perData, unsigned int resId){
    std::cout<<"load res "<<resId<<std::endl;
    DataList *p = perData;
    while (p != NULL)
    {
        std::cout<<"key:"<<p->key<<std::endl;
        std::cout<<"len:"<<p->len<<std::endl;
        std::cout<<"status:"<<p->status<<std::endl;
        std::cout<<"type:"<<p->type<<std::endl;
        ParamList *de = p->desc;
        while(de!=NULL)
        {
            std::cout<<"这里是描述信息"<<std::endl;
            std::cout<<"key:"<<de->key<<std::endl;
            std::cout<<"value:"<<de->value<<std::endl;
            std::cout<<"vlen:"<<de->vlen<<std::endl;
            de=de->next;
        }
        delete de;
        p=p->next;
    }
    delete p;
//	printf("wrapperLoadRes\n");
	return 0;
}

int WrapperAPI wrapperUnloadRes(unsigned int resId){
//	printf("wrapperUnloadRes\n");
    std::cout<<"unload res "<<resId<<std::endl;
	return 0;
}

const void* WrapperAPI wrapperCreate(const char* usrTag, pParamList params, wrapperCallback cb, unsigned int psrIds[], int psrCnt, int* errNum){
//	printf("wrapperCreate\n");
     pthread_mutex_lock(&idle_mutex_);

    KJINDEX++;
    int* a=new int();
    *a=KJINDEX;
    RsltStatus.insert(std::pair<int*, int>(a, 0));
    pthread_mutex_unlock(&idle_mutex_);

    return a;
	//return wrapperInnerHdl;
}

int WrapperAPI wrapperWrite(const void* handle, pDataList reqData){
//	printf("wrapperWrite\n");
    pthread_mutex_lock(&idle_mutex_);

    void* pChar = const_cast<void*>(handle);
    int* sp = static_cast<int*>(pChar);
    std::map<int*, int>::iterator ite = RsltStatus.find(sp);
    if (ite != RsltStatus.end() && reqData != NULL){
        ite->second = reqData->status;
    }
    pthread_mutex_unlock(&idle_mutex_);
	return 0;
}

int WrapperAPI wrapperRead(const void* handle, pDataList* respData){
//	printf("wrapperRead\n");
    /*
Note: demo 关闭实时读功能,wrapperRead仅需返回最终last结果,
    若需实时返回结果,可通过修改配置[aiges].realTimeRlt = true
     */
    *respData = wrapperInnerRslt;
    pthread_mutex_lock(&idle_mutex_);

     void* pChar = const_cast<void*>(handle);
     int* sp = static_cast<int*>(pChar);
     std::map<int*, int>::iterator ite = RsltStatus.find(sp);
    if (ite != RsltStatus.end()){
        (*respData)->status = DataStatus(ite->second);
    }

    pthread_mutex_unlock(&idle_mutex_);
	return 0;
}

int WrapperAPI wrapperDestroy(const void* handle){
//	printf("wrapperDestroy\n");
    pthread_mutex_lock(&idle_mutex_);
     void* pChar = const_cast<void*>(handle);
        int* sp = static_cast<int*>(pChar);
        std::map<int*, int>::iterator ite = RsltStatus.find(sp);
    if (ite != RsltStatus.end()){
        RsltStatus.erase(ite);
    }

    pthread_mutex_unlock(&idle_mutex_);
	return 0;
}

int WrapperAPI wrapperExec(const char* usrTag, pParamList params, pDataList reqData, pDataList* respData, unsigned int psrIds[], int psrCnt){
//	printf("wrapperExec\n");
    *respData = wrapperOnceRslt;
    while (params != NULL) {
        if (strcmp(params->key, "ctrl") == 0) {
            std::string val=std::string(params->value);
            if (val=="normal"){
                std::cout<<"call exec normal"<<std::endl;
                DataList* wrapperRslt = (struct DataList*)malloc(sizeof(struct DataList));
                wrapperRslt->key = (char*)"result";
                char* temp = (char*)malloc(params->vlen);
                strcpy(temp, params->value);
                wrapperRslt->data = (void*)temp;
                wrapperRslt->len = params->vlen;
                wrapperRslt->desc = NULL;
                wrapperRslt->type = DataText;
                wrapperRslt->status = DataOnce;
                wrapperRslt->next = NULL;
                *respData = wrapperRslt;
            }else if(val=="deadlock"){
               std::cout<<"call exec deadlock"<<std::endl;
               my_mutex.lock();
            }else if(val=="crash"){
                std::cout<<"call exec crash"<<std::endl;
                free((*respData)->data);
                free(*respData);
                free((*respData)->data);
            }else if(val=="random"){
                //根据调用次数的奇数 偶数 返回正确 失败
                reqCount+=1;
                if(reqCount%2==0){
                    return 21111;
                }
                std::cout<<"call exec random"<<std::endl;
                DataList* wrapperRslt = (struct DataList*)malloc(sizeof(struct DataList));
                wrapperRslt->key = (char*)"result";
                char* temp = (char*)malloc(params->vlen);
                strcpy(temp, params->value);
                wrapperRslt->data = (void*)temp;
                wrapperRslt->len = params->vlen;
                wrapperRslt->desc = NULL;
                wrapperRslt->type = DataText;
                wrapperRslt->status = DataOnce;
                wrapperRslt->next = NULL;
                *respData = wrapperRslt;
            }else if(val=="limit"){
                reqCount+=1;
                if(reqCount%2==0){
                    return 21111;
                }
            }else{
                return 22222;
            }
        }
        params = params->next;
    }
    return 0;
}

int WrapperAPI wrapperExecFree(const char* usrTag, pDataList* respData){
//	printf("wrapperExecFree\n");
	if (*respData != wrapperOnceRslt) {
	    free((*respData)->data);
	    free(*respData);
	}
	return 0;
}

int WrapperAPI wrapperExecAsync(const char* usrTag, pParamList params, pDataList reqData, wrapperCallback callback, int timeout, unsigned int psrIds[], int psrCnt){
//	printf("wrapperExecAsync\n");
    callback(usrTag,wrapperOnceRslt,0);
	return 0;
}

const char* WrapperAPI wrapperDebugInfo(const void* handle){
//	printf("wrapperDebugInfo\n");
	return "DebugInfo";
}
