FROM public.ecr.aws/iflytek-open/aiges-gpu:10.1-1.17-3.9.13-ubuntu1804-v1.3.2 

RUN  pip3 config set global.index-url https://pypi.mirrors.ustc.edu.cn/simple/ 

copy wrapper /home/wrapper

RUN pip3 install requests


# 拷贝mmocr配置文件到  /home/aiges, 因为工作目录位于 /home/aiges
#RUN cp -ra /home/mmocr/configs /home/aiges



ENV PYTHONPATH=$PYTHONPATH:/home/wrapper:/home:/home/aiges

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/wrapper




#ENV PYTHONPATH=$PYTHONPATH:/home/mmocr

#ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/wrapper/wrapper_lib


COPY xtest.toml /home/aiges
COPY xtest /home/aiges/xtest
COPY aiges.toml /home/aiges
#CMD ["sh", "-c", "./AIservice -m=0 -c=aiges.toml -s=svcName -u=http://companion.xfyun.iflytek:6868 -p=AIaaS -g=dx"]