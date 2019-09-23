FROM 172.16.59.153/aiaas/ubuntu:14.04.gcc.gdb.golang.probes
MAINTAINER yjhuang3@iflytek.com
COPY ./atmos/dest/resource /root/go/src/resource
COPY ./atmos/exe /root/go/src/git.xfyun.cn/AIaaS/atmos
ENV GOPATH /root/go
ENV LD_LIBRARY_PATH /root/go/src/resource
WORKDIR /root/go/src/resource
CMD ["bash", "start.sh"]