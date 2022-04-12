FROM 172.16.59.153/aiaas/ubuntugo_gcc:1.9.2
MAINTAINER yangzhou10@iflytek.com
COPY ./finder-go/bin /root/go/src/github.com/xfyun/finder-go/v3
ENV GOPATH /root/go
WORKDIR /root/go/src/github.com/xfyun/finder-go/v3
CMD ["bash", "demo"]
