# coding:utf-8
import sys

from aiges.sdk import WrapperBase, \
    StringParamField, \
    ImageBodyField, \
    StringBodyField, \
    AudioBodyField

# 

try:
    from aiges_embed import ResponseData, Response, DataListNode, DataListCls
except:
    from aiges.dto import Response, ResponseData, DataListNode, DataListCls


import base64
import json
import hashlib
import hmac
import os
import time
import requests
import re


'''
定义请求类:
 params:  params 开头的属性代表最终HTTP协议中的功能参数parameters部分， 
          params Field支持 StringParamField，
          NumberParamField，BooleanParamField,IntegerParamField，每个字段均支持枚举
          params 属性多用于协议中的控制字段，请求body字段不属于params范畴

 input:    input字段多用与请求数据段，即body部分，当前支持 ImageBodyField, StringBodyField, 和AudioBodyField
'''


class UserRequest(object):
    params1 = StringParamField(key="mode", enums=["music", "humming"], value='humming')

    input1 = AudioBodyField(key="data", path="/home/wrapper/test.wav")

'''
定义响应类:
 accepts:  accepts代表响应中包含哪些字段, 以及数据类型

 input:    input字段多用与请求数据段，即body部分，当前支持 ImageBodyField, StringBodyField, 和AudioBodyField
'''


class UserResponse(object):
    accept1 = StringBodyField(key="ouput_text")

class Wrapper(WrapperBase):
    serviceId = "music_third_api"
    version = "backup.0"
    requestCls = UserRequest()
    responseCls = UserResponse()

    requrl, http_method, http_uri = None, None, None
    # music
    access_key_music, access_secret_music = None, None
    # humming
    access_key_humming, access_secret_humming = None, None


    #config = {}
    #config = {
    #    "requrl" : "...",
    #    "http_method" : "...",
    #    "http_uri" : "...",
    #    "access_key_music" : "...",
    #    "access_secret_music" : "...",
    #    "access_key_humming" : "...",
    #    "access_secret_humming" : "..."
    #    }
   
   
   
    '''
    服务初始化
    @param config:
        插件初始化需要的一些配置，字典类型
        key: 配置名
        value: 配置的值
    @return
        ret: 错误码。无错误时返回0FYV2
    '''

    def wrapperInit(cls, config: {}) -> int:
        print("Initializing ..")
        config = config

        Wrapper.requrl, Wrapper.http_method, Wrapper.http_uri = config['requrl'], config['http_method'], config['http_uri']
        Wrapper.access_key_music, Wrapper.access_secret_music = config['access_key_music'], config['access_secret_music']
        Wrapper.access_key_humming, Wrapper.access_secret_humming = config['access_key_humming'], config['access_secret_humming']

        print('----------Finish Init--------------')
        
        return 0

    '''
    服务逆初始化

    @return
        ret:错误码。无错误码时返回0
    '''

    def wrapperFini(cls) -> int:
        print('------------------Finished-------------------')
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

    def wrapperOnceExec(self, params: {}, reqData: DataListCls) -> Response:
        print(" --------------Start Exec------------------")

        data_mode = params['mode']
        print(f'data_mode: {data_mode}')
        
        access_key = Wrapper.access_key_music if data_mode == 'music' else Wrapper.access_key_humming
        access_secret = Wrapper.access_secret_music if data_mode == 'music' else Wrapper.access_secret_humming
       
        src = reqData.list[0].data# binary files
        sample_bytes = reqData.list[0].len
        signature_version, data_type = '1', 'audio'
        print(type(src)) 

        timestamp = time.time()
        res = Response()
        
        string_to_sign = Wrapper.http_method + '\n' \
                    + Wrapper.http_uri + '\n' \
                    + access_key + '\n' \
                    + data_type + '\n' \
                    + signature_version + '\n' \
                    + str(timestamp)
        sign = base64.b64encode(hmac.new(access_secret.encode('ascii'), string_to_sign.encode('ascii'),
                                   digestmod=hashlib.sha1).digest()).decode('ascii')
        
        if sign is None:
            return res.response_err(5014)
        
        files = {'sample': src}
        data = {
            'access_key': access_key,
            'sample_bytes': sample_bytes,
            'timestamp': str(timestamp),
            'signature': sign,
            'data_type': data_type,
            'signature_version': signature_version
        }

        try:
            r = requests.post(Wrapper.requrl, files=files, data=data, timeout=5)
            print("--------")
            print(r.text)
            print("----------")
        except requests.exceptions.ConnectTimeout:
            return res.response_err(4408)
        if r is None:
            return res.response_err(4003)

        if r.status_code != 200:
            return res.response_err(4000 + r.status_code)

        pattern = re.compile('"code":\d+')
        error_code = re.findall(pattern, r.text)
        error_code = error_code[0].split(':')[-1]
        
        if int(error_code):
            return self.response_err(int(error_code))
        else:
            r.encoding = 'utf-8'
            print('-------------------------------------------')
            print(r.content)
            print('-------------------------------------------')
            
            l = ResponseData()
        
            l.key = "output_text"
            l.type = 0
            l.status = 3
            l.data = r.text
            l.len = len(r.text.encode())
            res.list = [l]
        return res 

    def wrapperError(cls, ret: int) -> str:
        if ret == 1001:
            return "识别无结果"
        elif ret == 2000:
            return "录音失败，可能是设备权限问题"
        elif ret == 2001:
            return "初始化错误或者初始化超时"
        elif ret == 2002:
            return "处理metadata错误"
        elif ret == 2004:
            return "无法生成指纹（有可能是静音）"
        elif ret == 2005:
            return "超时"
        elif ret == 3000:
            return "服务端错误"
        elif ret == 3001:
            return "Access Key不存在或错误"
        elif ret == 3002:
            return "HTTP内容非法"
        elif ret == 3003:
            return "请求数超出限制（需要升级账号）"
        elif ret == 3006:
            return "参数非法"
        elif ret == 3014:
            return "签名非法"
        elif ret == 3015:
            return "QPS超出限制（需要升级账号）"
        else:
            return f"Defined Error: {ret}"


    def wrapperTestFunc(cls, data: [], respData: []):
        r = Response()
        l = ResponseData()
        l.key = "ccc"
        l.status = 1
        d = open("pybind11/docs/pybind11-logo.png", "rb").read()
        l.len = len(d)
        l.data = d
        r.list = [l, l, l]

        print(r.list)
        print(444)
        return r

if __name__ == '__main__':
    m = Wrapper()
    m.schema()
    m.run()
