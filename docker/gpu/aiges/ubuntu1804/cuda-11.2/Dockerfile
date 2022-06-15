FROM artifacts.iflytek.com/docker-private/atp/cuda-go-python-base:11.2.0-1.17-3.9.13-ubuntu1804 as go-builder

MAINTAINER ybyang7@iflytek.com

ENV GOPROXY=https://goproxy.cn,direct
COPY docker/repos/ubuntu1804/sources.list /etc/apt/sources.list
RUN apt-get update && apt-get install -y libnuma-dev build-essential git sudo
COPY . /home/aiges
WORKDIR /home/aiges

#RUN mkdir -p output/include && go mod vendor && go build -v -o ./output/AIservice -gcflags "-N -l -c 10" main/main.go && \
#    cp ./cgo/header/widget/* ./output/include/ && \
#    cp  -r ./cgo/library/* ./output/

RUN make build-pack

FROM artifacts.iflytek.com/docker-private/atp/cuda-go-python-base:10.1-1.17-3.9.13-ubuntu1804 as cpython-builder
MAINTAINER ybyang7@iflytek.com

COPY docker/repos/ubuntu1804/sources.list /etc/apt/sources.list

WORKDIR /home/loader
RUN apt update && apt install -y build-essential vim

RUN DEBIAN_FRONTEND=noninteractive apt update &&apt install -y libnuma-dev libboost-all-dev
COPY ai_cpython_wrapper/ /home/wrapper

RUN cd /home/wrapper && make

#ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/aiges:/home/wrapper/wrappere_lib

FROM artifacts.iflytek.com/docker-private/atp/cuda-go-python-base:10.1-1.17-3.9.13-ubuntu1804 as final
MAINTAINER ybyang7@iflytek.com
ENV TZ Asia/Shanghai

COPY docker/repos/ubuntu1804/sources.list /etc/apt/sources.list
RUN apt update && apt install -y build-essential vim
RUN DEBIAN_FRONTEND=noninteractive apt update &&apt install -y libnuma-dev libboost-all-dev git vim
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/aiges:/home/aiges/library
WORKDIR /home/aiges
COPY --from=go-builder /home/aiges/dist/aiges_linux_amd64_v1 .
COPY --from=cpython-builder /home/wrapper/wrapper_lib /home/aiges/library
