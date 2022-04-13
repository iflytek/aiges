import torch
import numpy as np
from PIL import Image
import io
import flags

'''
初始化
config的值是由aiges.toml中[wrapper]各段设置的
'''

model = None

logger = flags.logger
model = torch.hub.load('.', 'yolov5s', source="local", pretrained=True)

def wrapperInit(config: {}) -> int:
    logger.info("model initializing...")
    logger.info("engine config %s" % str(config))
    global model

    logger.info("init success")
    return 0


'''
逆初始化
'''


def wrapperFini() -> int:
    logger.info("fini success", flush=True)
    return 0


'''
once接口执行函数
'''


def wrapperOnceExec(usrTag: str, params: {}, reqData: [], respData: [], psrIds: [], psrCnt: int) -> int:
    img = np.array(Image.open(io.BytesIO(reqData[0]["data"])).convert('RGB'))
    global model
    rlt = model(img)
    value = rlt.pandas().xyxy[0].to_json(orient="records")
    length = len(str(value).encode())
    respData.append({"key": "boxes", "data": value, "len": length, "status": 3, "type": 0})
    print(respData, flush=True)
    return 0


'''
根据不同错误码返回不同的错误描述
'''


def wrapperError(ret: int) -> str:
    if ret == 10013:
        return "reqData is empty"
    elif ret == 10001:
        return "load onnx model failed"
    else:
        return "other error code"
