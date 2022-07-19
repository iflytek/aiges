FROM public.ecr.aws/iflytek-open/aiges-gpu:10.1-1.17-3.9.13-ubuntu1804-v1.3.2 

RUN  pip3 config set global.index-url https://pypi.mirrors.ustc.edu.cn/simple/ 

copy wrapper /home/wrapper

RUN pip3 install requests

ENV PYTHONPATH=$PYTHONPATH:/home/wrapper:/home:/home/aiges

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/wrapper

COPY xtest.toml /home/aiges
COPY xtest /home/aiges/xtest
COPY aiges.toml /home/aiges
#CMD ["sh", "-c", "./AIservice -m=0 -c=aiges.toml -s=svcName -u=http://companion.xfyun.iflytek:6868 -p=AIaaS -g=dx"]