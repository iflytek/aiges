'''
服务初始化
@param config:
	type:	{}
	desc:	根据配置（来源于用户启动进程时指定配置文件aiges.toml中【wrapper】下属的键值对）初始化插件
        key: 配置名
        value: 配置的值
@return 接口错误码
	type:	int
	desc:	错误码,无错误时返回0
'''
def wrapperInit(config: {}) -> int:
    print("init success",flush=True)
    return 0

'''
服务逆初始化,结束进程时进行一些回收操作
@return 接口错误码
	type: 	int
    desc:	ret:错误码。无错误码时返回0
'''
def wrapperFini() -> int:
    print("fini success",flush=True)
    return 0

'''
非会话模式计算接口,对应oneShot请求,可能存在并发调用
@param usrTag
	type:	string 
	desc:	句柄
@param params 
	type:	{}
    desc:	功能参数(对应ASE平台上的功能参数)
@param  reqData
	type:	[{}] 
	desc:	写入数据实体(对应ASE平台上的请求数据)
        key: 请求数据段的名称
        type: 请求的数据类型 0: 文本 1:音频 2:图像 3:视频 4:个性化数据
        stauts: 数据状态 非流式请固定为3
        data: 请求的数据实体
        len: 请求的数据实体的长度
@param  respData
	type:	[{}]
	desc:	返回结果实体(对应ASE平台上的响应数据)
        key: 返回数据段的名称
        type: 返回的数据类型 0: 文本 1:音频 2:图像 3:视频 4:个性化数据
        stauts: 数据状态 非流式请固定为3
        data: 返回的数据实体
        len: 返回的数据实体的长度
@param psrIds
	type:	[int]
	desc:	需要使用的个性化资源标识列表(如无个性化需要可忽略)
@param psrCnt 
	type:	int
	desc:	需要使用的个性化资源个数

@return 接口错误码
	type:	int
	desc:	错误码,无错误时返回0
'''
def wrapperOnceExec(usrTag:str,params:{},reqData:[],respData:[],psrIds:[],psrCnt:int) -> int:
    rlt="this is a result from python"
    singleData={}
    singleData["key"]="result"
    singleData["type"]=0
    singleData["status"]=3
    singleData["data"]=rlt
    singleData["len"]=len(rlt)
    respData.append(singleData)    
    return 0
'''
不同处理码返回不同的错误描述信息
@param
@return:
	type: string
	desc: 错误码对应的错误描述
'''
def wrapperError(ret:int)->str:
    print("call wrapperError"+str(ret),flush=True)
    return "custom error str with code:"+str(ret)
