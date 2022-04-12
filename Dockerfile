FROM artifacts.iflytek.com/docker-private/jianjiang/ubuntu_go:20.04_1.16 as builder
MAINTAINER ybyang7@iflytek.com

ENV GOPROXY=https://goproxy.cn,direct
RUN apt-get update && apt-get install -y libnuma-dev build-essential
COPY src/github.com/xfyun/aiges /home/aiges
WORKDIR /home/aiges

RUN mkdir -p output/include && go mod vendor && go build -v -o ./output/AIservice -gcflags "-N -l -c 10" main/main.go && \
    cp ./cgo/header/widget/* ./output/include/ && \
    cp  -r ./cgo/library/* ./output/
  

#FROM  artifacts.iflytek.com/docker-private/jianjiang/ubuntu_go:20.04_1.16 as prod 
FROM artifacts.iflytek.com/docker-private/atp/miniconda3:latest
MAINTAINER ybyang7
RUN echo '''deb https://mirrors.aliyun.com/debian/ bullseye main non-free contrib \
deb-src https://mirrors.aliyun.com/debian/ bullseye main non-free contrib \
deb https://mirrors.aliyun.com/debian-security/ bullseye-security main \
deb-src https://mirrors.aliyun.com/debian-security/ bullseye-security main \
deb https://mirrors.aliyun.com/debian/ bullseye-updates main non-free contrib \
deb-src https://mirrors.aliyun.com/debian/ bullseye-updates main non-free contrib \
deb https://mirrors.aliyun.com/debian/ bullseye-backports main non-free contrib \
deb-src https://mirrors.aliyun.com/debian/ bullseye-backports main non-free contrib''' >/etc/apt/sources.list

WORKDIR /home/loader
RUN apt update && apt install -y build-essential vim


# Create the environment:
RUN echo '''name: loader \n\
channels: \n \
- defaults \n\
dependencies: \n \
- python=3.9.12 \n\
prefix: /opt/loader ''' > /home/loader/environment.yml &&  conda env create -f environment.yml




# Make RUN commands use the new environment:
SHELL ["conda", "run", "-n", "loader", "/bin/bash", "-c"]

RUN echo "conda activate loader " >> ~/.bashrc
ENV TZ Asia/Shanghai

RUN DEBIAN_FRONTEND=noninteractive apt update &&apt install -y libnuma-dev libboost-all-dev
WORKDIR /home/aiges
COPY --from=builder /home/aiges/output .


COPY ai_cpython_wrapper/ /home/wrapper
 
RUN cd /home/wrapper && make

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/aiges:/home/wrapper/wrappere_lib

