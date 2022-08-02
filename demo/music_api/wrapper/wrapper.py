# encoding: utf-8
import base64
import json
import hashlib
import hmac
import os
import sys
import time
import requests
import logging
import re

log_path = 'usr.log'
logging.basicConfig(level=logging.INFO, filename=log_path, 
                     format='{"ts": %(asctime)s, "msg": %(message)s}', #日志输出的格式
                     datefmt="%Y-%m-%d %H:%M:%S") #时间输出的格式

requrl, http_method, http_uri = None, None, None
# music
access_key_music, access_secret_music = None, None
# humming
access_key_humming, access_secret_humming = None, None

# 保存响应数据
def mkdir(file_path):
    if not os.path.exists(file_path):
        os.mknod(file_path)
    return 0

def wrapperInit(config) -> int:
    '''
        config 是一个字典，在aiges.toml的wrapper字段配置
    '''
    global requrl, http_method, http_uri
    global access_key_music, access_secret_music, access_key_humming, access_secret_humming
    
    requrl, http_method, http_uri = config['requrl'], config['http_method'], config['http_uri']
    access_key_music, access_secret_music = config['access_key_music'], config['access_secret_music']
    access_key_humming, access_secret_humming = config['access_key_humming'], config['access_secret_humming']

    global log_path
    mkdir(log_path)

    logging.info('Init successfully.')
    return 0

def wrapperFini() -> int:
    logging.info('Wrapper finished.')
    return 0

def wrapperOnceExec(usrTag: str, params: {}, reqData: [], respData: [], psrIds: [], psrCnt: int) -> int:
    global requrl, http_method, http_uri
    global access_key_music, access_secret_music, access_key_humming, access_secret_humming

    data_mode = params['mode']

    if data_mode not in ('music', 'humming'):
        logging.error('非法的识别模式（仅限于music或者humming）')
        return 5006

    access_key = access_key_music if data_mode is 'music' else access_key_humming
    access_secret = access_secret_music if data_mode is 'music' else access_secret_humming
    
    src = reqData[0]['data']# binary files
    sample_bytes = reqData[0]['len']
    signature_version, data_type = '1', 'audio'

    timestamp = time.time()

    string_to_sign = http_method + '\n' \
                    + http_uri + '\n' \
                    + access_key + '\n' \
                    + data_type + '\n' \
                    + signature_version + '\n' \
                    + str(timestamp)
    sign = base64.b64encode(hmac.new(access_secret.encode('ascii'), string_to_sign.encode('ascii'),
                                   digestmod=hashlib.sha1).digest()).decode('ascii')
 
    if sign is None:
        logging.error('HMAC failure.')
        return 5014
     
    files = {'sample': src}
    data = {
        'access_key': access_key,
        'sample_bytes': sample_bytes,
        'timestamp': str(timestamp),
        'signature': sign,
        'data_type': data_type,
        'signature_version': signature_version
    }
    logging.info('Post http request.')
    try:
        r = requests.post(requrl, files=files, data=data, timeout=5)
    except requests.exceptions.ConnectTimeout:
        logging.error('Http post timeout.')
        return 4408# http timeout
    if r is None:
        logging.error("HTTP内容非法")
        return 4003

    if r.status_code != 200:
        return 4000 + r.status_code

    pattern = re.compile('"code":\d+')
    error_code = re.findall(pattern, r.text)
    error_code = error_code[0].split(':')[-1]
    if int(error_code):
        return int(error_code)
    else:
        r.encoding = 'utf-8'
        logging.warning(r.text)
        #print('-------------------------------------------')
        #print(r.text)
        #print('-------------------------------------------')
    
        respData.append({
            'key': 'output_text',
            'type': 0,
            'status': 3,
            "data": r.text,
            "len": len(r.text.encode())
        })
        return 0

def wrapperError(ret: int) -> str:
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
        return f"User Defined Error: {ret}"