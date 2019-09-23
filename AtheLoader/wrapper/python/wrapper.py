#!/usr/bin/env python

# Note:GIL问题(PyGILState_Ensure/PyGILState_Release);

class EngineBase:   # clase EngineBase(object):  c call exception

    '''
    aiges engine wrapper
    '''
    version = '1.0.1'                   # 类属性-版本号,子类需实现,上层通过版本号获取服务配置

    def wrapperInit(self, **kwargs):
        '''
        @初始化接口,上层逻辑限制,替代__init__接口
        输入：配置对kwargs(dict)
        输出：return int, 返回错误码
        '''
        pass

    def wrapperFini(self):
        '''
        @逆初始化接口,替代__del__接口
        输出：return int, 返回错误码
        '''
        pass

    def wrapperCreate(self, *args, **kwargs):
        '''
        @实例创建接口
        输入：个性化资源id列表args, 会话参数对kwargs
        输出：成功则返回实例object inst 用于关联会话相关接口, 失败None
        '''
        pass

    def wrapperDestroy(self, inst):
        '''
        @实例销毁接口
        输入：由create接口申请所得实例object inst
        输出：return int, 返回错误码
        '''
        pass

    def wrapperError(self, code):
        '''
        @错误码描述转换接口
        输入：error code(int)
        输出：return str, 返回错误描述信息
        '''
        pass

    def wrapperDebugInfo(self, inst):
        '''
        @调试日志接口
        输入：由create接口申请所得实例object inst
        输出：return str, 返回会话计算过程中的debug信息
        '''
        pass

    def wrapperWrite(self, inst, dataInput):
        '''
        @数据写接口,由上层框架调用写数据
        输入：引擎实例object inst, 用户请求数据[]
            dataInput = [{
                "Key":"",
                "Data":"",
                "Type":"",
                "Status":"",
                "Desc":"",
                "Encoding":"",
            }]
        输出：return int, 返回错误码
        '''
        pass

    def wrapperRead(self, inst):
        '''
        @数据读接口,上层框架调用读取引擎结果
        输入：引擎实例object inst
        输出：return [], int ;返回引擎计算结果dataOutput,错误信息err; 即tuple
            dataOutput = [{
                "Key":"",
                "Data":"",
                "Type":"",
                "Status":"",
                "Desc":"",
                "Encoding":"",
            }]
        '''
        pass

    def wrapperExec(self, **kwargs, msg):
        '''
        @非会话oneShot接口
        输入：用户请求参数对kwargs, 用户请求数据msg []
            dataInput = [{
                "Key":"",
                "Data":"",
                "Type":"",
                "Status":"",
                "Desc":"",
                "Encoding":"",
            }]
        输出：return int, [] ; 返回错误信息err, 引擎计算结果dataOutput; 即tuple
        '''
        pass

    def wrapperLoadRes(self, id, res):
            '''
            @资源加载接口
            输入：资源数据res, 资源id(int)
                res = {
                            "Key":"",
                            "Data":"",
                            "Type":"",
                            "Status":"",
                            "Desc":"",
                            "Encoding":"",
                }
            输出：return int, 返回错误码
            '''
            pass

    def wrapperUnloadRes(self, id):
        '''
        @资源卸载接口
        输入：待卸载资源id(int)
        输出：return int, 返回错误码
        '''
        pass

if __name__ == "__main__":
    engine = EngineBase()
    print(engine.version)
