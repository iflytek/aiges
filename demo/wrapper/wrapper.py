import sys
if not hasattr(sys, 'argv'):
    sys.argv  = ['']


'''
服务初始化
@param config:
    插件初始化需要的一些配置，字典类型
    key: 配置名
    value: 配置的值
@return
    ret: 错误码。无错误时返回0
'''
def wrapperInit(config: {}) -> int:
    return 0


'''
服务逆初始化

@return
    ret:错误码。无错误码时返回0
'''
def wrapperFini() -> int:
    return 0

'''
非会话模式计算接口,对应oneShot请求,可能存在并发调用

@param usrTag 句柄
#param params 功能参数
@param  reqData     写入数据实体
@param  respData    返回结果实体,内存由底层服务层申请维护,通过execFree()接口释放
@param psrIds 需要使用的个性化资源标识列表
@param psrCnt 需要使用的个性化资源个数

@return 接口错误码
    reqDat
    ret:错误码。无错误码时返回0
'''
def wrapperOnceExec(usrTag:str,params:{},reqData:[],respData:[],psrIds:[],psrCnt:int) -> int:
    print("hello world")
    print(usrTag)
    print(params)
    print(reqData)
    print(psrIds)
    print(psrCnt)
    return 100


def wrapperCreate(usrTag: str, params: [], psrIds: [], psrCnt: int) -> str:
    return ""


def wrapperWrite(handle: str, datas: []) -> int:
    return 0


def wrapperRead(handle: str) -> []:
    return []


def wrapperDestroy(handle: str) -> int:
    return 0


def wrapperError(ret:int)->str:
    if ret==100:
        return "this is a tese error return"
    return ""