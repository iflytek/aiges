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

### CPU推理

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


