FROM loader:latest

WORKDIR /home/aiges

ADD requirements.txt /home/loader/requirements.txt

RUN  pip config set global.index-url https://pypi.mirrors.ustc.edu.cn/simple/   && pip install -r /home/loader/requirements.txt


ADD ./wrapper /home/aiges/wrapper

ENV PYTHONPATH=$PYTHONPATH:/home/aiges/wrapper 
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/wrapper/wrapper_lib

COPY yolov5s.pt /home/aiges/wrapper
