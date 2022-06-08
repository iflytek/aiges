import numpy as np
from PIL import Image
import io
import flags
from mmocr.utils.ocr import MMOCR
import json

'''
初始化
config的值是由aiges.toml中[wrapper]各段设置的
'''

model = None

logger = flags.logger

def wrapperInit(config: {}) -> int:
    logger.info("model initializing...")
    logger.info("engine config %s" % str(config))
    global model
    # Load models into memory
    model = MMOCR(det='TextSnake', recog=None)
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
    rlt = model.readtext(img,details=True)
    rlt = json.dumps(rlt)
    respData.append({"key": "boxes", "data": rlt, "len": len(rlt), "status": 3, "type": 0})
    #respData.append(rlt)
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
