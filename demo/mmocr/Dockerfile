FROM artifacts.iflytek.com/docker-private/atp/py_loader:py39

RUN  pip3 config set global.index-url https://pypi.mirrors.ustc.edu.cn/simple/ 

## install cuda tookit
RUN apt -y install nvidia-cuda-toolkit nvidia-cuda-dev


RUN pip3 install torch==1.10  torchvision 
RUN pip3 install openmim && \
    mim install mmcv-full && \
    mim install mmdet

 copy mmocr /home/mmocr
