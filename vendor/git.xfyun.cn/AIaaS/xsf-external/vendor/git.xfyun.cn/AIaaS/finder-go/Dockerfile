FROM 172.16.59.153/aiaas/ubuntugo_gcc:1.9.2
MAINTAINER yangzhou10@iflytek.com
COPY ./finder-go/bin /root/go/src/git.xfyun.cn/AIaaS/finder-go
ENV GOPATH /root/go
WORKDIR /root/go/src/git.xfyun.cn/AIaaS/finder-go
CMD ["bash", "demo"]