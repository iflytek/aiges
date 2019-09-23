#ifndef __AIGES_WRAPPER_H__
#define __AIGES_WRAPPER_H__

#include "type.h"

#ifdef __cplusplus
extern "C" {
#endif


/*
    wrapper服务层初始化
    @param  cfg         服务层配置对
*/
int WrapperAPI wrapperInit(pConfig cfg);

/*
    wrapper服务层逆初始化
*/
int WrapperAPI wrapperFini();

/*
    获取服务错误信息
    @param  errNum      服务层异常错误码
    @return             错误码对应的错误描述信息
*/
const char* WrapperAPI wrapperError(int errNum);

/*
    获取服务版本信息
    @return             服务版本信息
*/
const char* WrapperAPI wrapperVersion();


/// 以下接口为会话模式请求调用接口;
/// 1. 包含个性化资源加载/卸载接口;
/// 2. 包含会话模式上下文相关接口：create/write/read/destroy;
/// 3. 包含同步模式/异步模式接口;

/*
    个性化数据加载
    @param  perData     个性化加载数据
    @param  resId       个性化数据标记,由框架层生成传入
    @return             接口错误码
*/

int WrapperAPI wrapperLoadRes(pDataList perData, unsigned int resId);

/*
    个性化数据卸载
    @param  resId       个性化数据标记
    @return             接口错误码
*/
int WrapperAPI wrapperUnloadRes(unsigned int resId);

/*
    回调接口定义
    @param  handle      实例句柄,用于关联异步请求
    @param  respData    异步计算结果,通过回调返回框架层
    @note   无需集成方实现(框架实现),由集成方于请求数据计算完毕后调用;
*/
typedef int(*wrapperCallback)(const char* handle, pDataList respData);

/*
    创建计算资源
    @param  params      会话参数对
    @param  cb          异步回调:若同步响应则cb为null,通过wrapperRead获取结果
                                若异步响应则传入cb,通过回调cb返回结果至框架;
    @param  psrIds      会话所需个性化资源id
    @param  psrCnt      会话个性化资源Count
    @param  errNum      接口错误码[in/out]
    @return             引擎服务实例句柄,用于关联上下文;
*/
const char* WrapperAPI wrapperCreate(pParamList params, wrapperCallback cb, unsigned int[] psrIds, int psrCnt, int* errNum);

/*
    写入计算数据
    @param  handle      实例句柄,用于关联上下文;
    @param  reqData     写入数据实体
    @return             接口错误码
*/
int WrapperAPI wrapperWrite(const char* handle, pDataList reqData);

/*
    读取计算结果
    @param  handle      实例句柄,用于关联上下文;
    @param  respData    同步读取结果实体
    @return             接口错误码
    @note               respData内存由底层自行维护,在destroy阶段销毁
*/
int WrapperAPI wrapperRead(const char* handle, pDataList* respData);

/*
    释放计算资源
    @param  handle      会话句柄,用于关联上下文;
*/
int WrapperAPI wrapperDestroy(const char* handle);


/// 以下接口为非会话模式请求调用接口,对应引擎框架oneShot协议消息;
/// 1. 其中wrapperExec()为同步阻塞接口,要求引擎服务阻塞带计算完成返回计算结果;
/// 2. wrapperExecFree()为同步临时资源释放接口,用于释放wrapperExec产生的临时结果数据;
/// 3. wrapperExec()为异步非阻塞接口,要求引擎服务即时返回,异步计算结果通过回调callback返回;

/*
    非会话模式计算接口,对应oneShot请求
    @param  reqData     写入数据实体
    @param  respData    返回结果实体,内存由底层服务层申请维护,通过execFree()接口释放
    @return 接口错误码
    @note   同步操作接口
*/
int WrapperAPI wrapperExec(pParamList params, pDataList reqData, pDataList* respData);


/*
    同步接口响应数据缓存释放接口
    @param  respData    由同步接口exec获取的响应结果数据
*/
int WrapperAPI wrapperExecFree(pDataList* respData);


/*
    非会话模式计算接口,对应oneShot请求
    @param  handle      请求句柄,关联异步连接
    @param  reqData     写入数据实体
    @param  callback    异步回调接口,用于异步返回计算结果(框架实现)
    @param  timeout     异步超时时间,集成方实现该超时控制,ms;
    @note   异步操作接口
*/
int WrapperAPI wrapperExecAsync(const char* handle, pParamList params, pDataList reqData, wrapperCallback callback, int timeout);


/*
    调试信息输出接口
    @return 会话调试信息;
    @note   单次会话destroy前调用一次;
*/
const char* WrapperAPI wrapperDebugInfo();


#ifdef __cplusplus
}
#endif

#endif