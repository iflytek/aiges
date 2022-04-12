FROM loader:latest

WORKDIR /home/aiges

ADD requirements.txt /home/loader/requirements.txt

RUN  pip config set global.index-url https://pypi.mirrors.ustc.edu.cn/simple/   && pip install -r /home/loader/requirements.txt



ENV PYTHONPATH=$PYTHONPATH:/home/yolov5

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/wrapper/wrapper_lib

COPY yolov5  /home/yolov5

COPY Arial.ttf /root/.config/Ultralytics/Arial.ttf
COPY yolov5s.pt /home/aiges/yolov5s.pt
COPY xtest.toml /home/aiges
COPY xtest /home/aiges/xtest
COPY aiges.toml /home/aiges
COPY zidane.jpg /home/aiges
CMD ["sh", "-c", "./AIservice -m=0 -c=aiges.toml -s=svcName -u=http://companion.xfyun.iflytek:6868 -p=AIaaS -g=dx" ]

