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

```schema
{
  "type": "array",
  "items": {
    "type": "object",
    "required": [],
    "properties": {
      "filename": {
        "type": "string"
      },
      "result": {
        "type": "array",
        "items": {
          "type": "object",
          "required": [],
          "properties": {
            "box": {
              "type": "array",
              "items": {
                "type": "number"
              }
            },
            "box_score": {
              "type": "number"
            },
            "text": {
              "type": "string"
            },
            "text_score": {
              "type": "number"
            }
          }
        }
      }
    }
  }
}
```

真实数据
```json
[{
	"filename": "0",
	"result": [{
		"box": [190, 37, 253, 31, 254, 46, 191, 52],
		"box_score": 0.9566415548324585,
		"text": "nboroughofs",
		"text_score": 1.0
	}, {
		"box": [253, 47, 257, 36, 287, 47, 282, 58],
		"box_score": 0.9649642705917358,
		"text": "fsouthw",
		"text_score": 1.0
	}, {
		"box": [157, 59, 188, 41, 194, 52, 163, 70],
		"box_score": 0.9521175622940063,
		"text": "londond",
		"text_score": 0.9897959183673469
	}, {
		"box": [280, 58, 286, 50, 306, 67, 300, 74],
		"box_score": 0.9397556781768799,
		"text": "thwark",
		"text_score": 1.0
	}, {
		"box": [252, 78, 295, 78, 295, 98, 252, 98],
		"box_score": 0.9694718718528748,
		"text": "hill",
		"text_score": 1.0
	}, {
		"box": [165, 78, 247, 78, 247, 99, 165, 99],
		"box_score": 0.9548642039299011,
		"text": "octavia",
		"text_score": 1.0
	}, {
		"box": [164, 105, 215, 103, 216, 121, 165, 123],
		"box_score": 0.9806956052780151,
		"text": "social",
		"text_score": 1.0
	}, {
		"box": [219, 104, 294, 104, 294, 122, 219, 122],
		"box_score": 0.9688025116920471,
		"text": "reformer",
		"text_score": 1.0
	}, {
		"box": [150, 124, 226, 124, 226, 141, 150, 141],
		"box_score": 0.9752052426338196,
		"text": "established",
		"text_score": 1.0
	}, {
		"box": [229, 124, 255, 124, 255, 140, 229, 140],
		"box_score": 0.94972825050354,
		"text": "this",
		"text_score": 1.0
	}, {
		"box": [259, 125, 305, 123, 306, 139, 260, 142],
		"box_score": 0.9752089977264404,
		"text": "garden",
		"text_score": 1.1666666666666667
	}, {
		"box": [166, 142, 193, 141, 194, 156, 167, 157],
		"box_score": 0.9731062650680542,
		"text": "hall",
		"text_score": 1.0
	}, {
		"box": [198, 142, 223, 142, 223, 156, 198, 156],
		"box_score": 0.954893946647644,
		"text": "and",
		"text_score": 1.0
	}, {
		"box": [228, 144, 286, 144, 286, 159, 228, 159],
		"box_score": 0.977089524269104,
		"text": "cottages",
		"text_score": 1.25
	}, {
		"box": [180, 158, 205, 158, 205, 172, 180, 172],
		"box_score": 0.9400061964988708,
		"text": "and",
		"text_score": 1.0
	}, {
		"box": [210, 160, 279, 158, 279, 172, 210, 174],
		"box_score": 0.9543584585189819,
		"text": "pioneered",
		"text_score": 1.0
	}, {
		"box": [226, 176, 277, 176, 277, 188, 226, 188],
		"box_score": 0.9748533964157104,
		"text": "cadets",
		"text_score": 1.0
	}, {
		"box": [183, 177, 223, 177, 223, 189, 183, 189],
		"box_score": 0.9633154273033142,
		"text": "army",
		"text_score": 1.0
	}, {
		"box": [201, 190, 235, 190, 235, 204, 201, 204],
		"box_score": 0.971415102481842,
		"text": "1887",
		"text_score": 1.25
	}, {
		"box": [175, 213, 180, 201, 211, 212, 206, 225],
		"box_score": 0.9704344868659973,
		"text": "vted",
		"text_score": 0.9191176470588236
	}, {
		"box": [241, 213, 278, 200, 283, 213, 246, 227],
		"box_score": 0.9607459902763367,
		"text": "epeople",
		"text_score": 1.0
	}, {
		"box": [208, 224, 210, 212, 223, 214, 220, 227],
		"box_score": 0.9337806701660156,
		"text": "by",
		"text_score": 1.0
	}, {
		"box": [223, 214, 240, 214, 240, 226, 223, 226],
		"box_score": 0.969144344329834,
		"text": "the",
		"text_score": 1.0
	}]
}]
```


```webgate schema

```