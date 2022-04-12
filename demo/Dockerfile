FROM loader:latest

WORKDIR /home/aiges

ADD requirements.txt /home/loader/requirements.txt

RUN  pip config set global.index-url https://pypi.mirrors.ustc.edu.cn/simple/   && pip install -r /home/loader/requirements.txt


ADD ./wrapper /home/aiges/wrapper

ENV PYTHONPATH=$PYTHONPATH:/home/aiges/wrapper 
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/wrapper/wrapper_lib
COPY yolov5 /root/.cache/torch/hub/ultralytics_yolov5_master
COPY Arial.ttf /root/.config/Ultralytics/Arial.ttf

COPY yolov5s.pt /home/aiges/yolov5s.pt
COPY xtest.toml /home/aiges
COPY xtest /home/aiges/xtest
COPY aiges.toml /home/aiges
COPY zidane.jpg /home/aiges

