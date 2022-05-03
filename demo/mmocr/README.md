# MMOCR ASE Python推理



## 快速开始

* clone仓库

```bash
git clone https://github.com/open-mmlab/mmocr.git
```

若要使用ase python加载器进行mmocr推理， 根据mmocr官方介绍依赖于:

*  Linux | Windows | macOS

*  Python 3.7

* PyTorch 1.6 or higher

* torchvision 0.7.0

* CUDA 10.1

* NCCL 2

* GCC 5.4.0 or higher

* MMCV

* MMDetection

根据官方引导， 我们基于 ase py39 docker镜像进行推理镜像构建

### MMOCR CPU推理镜像构建

```dockerfile
# 基于py加载器镜像
FROM artifacts.iflytek.com/docker-private/atp/py_loader:py39

# 设置国内pip源，加速pip安装
RUN  pip3 config set global.index-url https://pypi.mirrors.ustc.edu.cn/simple/ 

# 安装mmocr依赖
RUN pip3 install torch==1.10  torchvision 
RUN  apt install -y  libgl1-mesa-glx && pip3 install openmim && \
    mim install mmcv-full && \
    mim install mmdet

# 拷贝mmocr 项目 ，根据需要调整，可用git clone
copy mmocr /home/mmocr

# 安装mmocr自身依赖
RUN cd /home/mmocr/ && pip3 install -e .
# 安装wrapper demo脚本依赖
RUN pip3 install iflags


# 拷贝提前写好的wrapper脚本
COPY wrapper.py /home/mmocr

# 拷贝mmocr配置文件到  /home/aiges, 因为工作目录位于 /home/aiges
RUN cp -ra /home/mmocr/configs /home/aiges

# 设置wrapper.py搜索路径
ENV PYTHONPATH=$PYTHONPATH:/home/mmocr

# 设置加载器搜索路径,后期会屏蔽到基础镜像中去
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/wrapper/wrapper_lib


# 测试工具文件
COPY xtest.toml /home/aiges
COPY xtest /home/aiges/xtest
COPY aiges.toml /home/aiges
CMD ["sh", "-c", "./AIservice -m=0 -c=aiges.toml -s=svcName -u=http://companion.xfyun.iflytek:6868 -p=AIaaS -g=dx"]

```

### wrapper.py 编写

```python
import numpy as np
from PIL import Image
import io
import flags
from mmocr.utils.ocr import MMOCR
# 导入 推理所需要的一些库以及一些工具工具
'''
初始化
config的值是由aiges.toml中[wrapper]各段设置的
'''

model = None

logger = flags.logger

def wrapperInit(config: {}) -> int:
    # init中实现 全局model/引擎初始化
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
当为once类型接口调用时会调用此接口
'''
def wrapperOnceExec(usrTag: str, params: {}, reqData: [], respData: [], psrIds: [], psrCnt: int) -> int:
    # 转换请求过来的图片为 ndarray
    img = np.array(Image.open(io.BytesIO(reqData[0]["data"])).convert('RGB'))
    global model
    # 调用mmocr的 model推理模块进行推理
    rlt = model.readtext(img)
    respData.append(rlt)
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

```


