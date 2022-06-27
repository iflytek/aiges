# AI Service Engine

<!-- markdownlint-disable MD033 -->

<span class="badge-placeholder">[![Forks](https://img.shields.io/github/forks/xfyun/aiges)](https://img.shields.io/github/forks/xfyun/aiges)</span>
<span class="badge-placeholder">[![Stars](https://img.shields.io/github/stars/xfyun/aiges)](https://img.shields.io/github/stars/xfyun/aiges)</span>
<span class="badge-placeholder">[![Build Status](https://github.com/xfyun/aiges/actions/workflows/build.yaml/badge.svg)](https://github.com/xfyun/aiges/actions/workflows/build.yaml)</span>
<span class="badge-placeholder">[![GitHub release](https://img.shields.io/github/v/release/xfyun/aiges)](https://github.com/xfyun/aiges/releases/latest)</span>
<span class="badge-placeholder">[![GitHub contributors](https://img.shields.io/github/contributors/xfyun/AthenaServing)](https://github.com/xfyun/AthenaServing/graphs/contributors)</span>
<span class="badge-placeholder">[![License: Apache2.0](https://img.shields.io/github/license/xfyun/aiges)](https://img.shields.io/github/license/xfyun/aiges)</span>



<!-- markdownlint-restore -->

## 背景
>引擎托管平台（AI Service Engine）是一个专为AI能力开发者打造的AI算法引擎的无服务全托管式平台，您可以通过集成ASE中提供的插件，快速的部署AI算法引擎，并使用网络、分发策略、数据处理等配套辅助系统。引擎托管平台致力于加速AI算法引擎云服务化，并借助云原生架构，为云服务的稳定提供多重保障，您无需关注底层基础设施及服务化相关的开发、治理和运维，即可高效、安全对引擎进行部署、升级、扩缩、运营和监控。

## 整体架构
![img](https://raw.githubusercontent.com/xfyun/proposals/main/athenaloader/athena.png)


## 使用工作流

![img](https://github.com/xfyun/proposals/blob/main/athenaloader/usage.png?raw=true)

## 开源计划

| 任务项 |目标 |时间 |
|-----|-----|-----|
|&#9745; [加载器](#AIGes) 通用引擎/模型加载器|独立部署可运行，支持python快速推理服务化|2022/Q2|
|&#9745; [负载均衡器](#LoadBalance) 负载聚合组件|独立部署可运行|2022/Q3|
|&#9744; [WebGate](#Webgate) Web网关组件|可运行|2022/Q3|
|&#9745; [Polaris](#Polaris) 配置中心与服务发现|独立部署可运行|2022/Q2|
|&#9744; [Atom](#Atom) 协议转换组件|可运行|2022/Q3|
|&#9744; Serving on Kubernetes Helm Chart一键部署 (进行中)|支持在k8s集群上一键部署推理服务框架|2022/Q3|
|&#9744; Serving on Docker with docker-compose 一键部署|支持使用docker-compose部署推理服务框架|2022/Q4|
|&#9744; 各组件Documentation建设 (进行中)|各组件文档详设，门户建设|2022/Q4|
|&#9744; 多领域模型Demo演示示例、GIF (进行中)|部分领域模型推理示例，如mmdetection,mmocr,yolo等|2022/Q3|
|&#9744; AIServing [API](#API协议) AI能力协议规范|完善开源协议说明，schema自动生成、校验工具|2022/Q3|
|&#9744; AseCTl命令行工具 [API](#Asectl) 命令行工具|支持能力一键生成，运行，配置管理等|2022/Q4|

## 使用说明

### AIGes (AI General Engine Service)

通用引擎加载器(部分文档中loader， loader engine均为别名)

### Documentation

View Doc on [Documentation](https://xfyun.github.io/inferservice/architechture/architechture/)

### AIGes 项目编译
```
make
```
### AIGes 项目打包
```
make build-pack
```
### AIGes 清理项目文件
```
make clean
```

## Polaris

Go编写的服务发现、配置中心

### Documentation

View Doc on [Documentation](https://xfyun.github.io/inferservice/architechture/architechture/)

### Features

&#9745; 支持模型推理成RPC服务(Serving框架会转成HTTP服务)

&#9745; 支持C代码推理 support c++/c code infer

&#9745; 支持Python代码推理 Support python code infer

&#9745; 支持配置中心，服务发现

&#9745; 支持负载均衡配置

&#9744; 支持Java代码推理或者其它

&#9744; 支持计量授权

### Protocol

参见: [ase-proto](https://github.com/xfyun/ase_protocol)


### 开源版docker镜像

#### 基础镜像

基础镜像中提供

- 基础的编译好的 Python加载器AIService(包含支持python的libwrapper.so), 目录结构如下

    加载器默认放置于 容器/home/aiges 目录
    ```bash
    root@e38a9aacc355:/home/aiges# pwd
    /home/aiges
    root@e38a9aacc355:/home/aiges# ls -l /home/aiges/
    total 18760
    -rwxr-xr-x 1 root root 19181688 Jun 10 15:30 AIservice
    -rw-r--r-- 1 root root     2004 Jun 10 18:15 aiges.toml
    drwxr-xr-x 3 root root     4096 Jun 10 15:30 include
    drwxrwxrwx 1 root root     4096 Jun 10 15:31 library
    drwxr--r-- 2 root root     4096 Jun 10 18:16 log
    -rw-r--r-- 1 root root       96 Jun 10 18:15 start_test.sh
    drwxr-xr-x 2 root root     4096 Jun 10 18:16 xsf_status
    drwxr-xr-x 2 root root     17711057 Jun 10 18:16 xtest
    -rw-r--r-- 1 root root     4232 Jun 10 17:54 xtest.toml
    ```
    其中 aiges.toml 用于本地启动测试使用

- Python环境: 不推荐用户后续镜像构建修改Python版本


#### 业务镜像

业务镜像一般需要用户自己编写Dockerfile构建，业务镜像中用户可以根据场景需要定制安装

- 推理运行时，如onnxruntime、torch vision等

- gpu驱动，cuda，cudnn等驱动


示例Dockerfile地址为

* [openmmlab_mmocr](https://github.com/xfyun/aiges/blob/master/demo/mmocr/Dockerfile_cpu)
* [yolov5](https://github.com/xfyun/aiges/blob/master/demo/yolov5/Dockerfile)


### 注意事项

~~* 示例尚未提供gpu runtime安装方法~~

* 用户的python插件需要自行定义一个目录放置， 如上述的 /home/mmocr , /home/yolov5, 且该位置需要像

https://github.com/xfyun/aiges/blob/master/demo/mmocr/Dockerfile_cpu#L23 

声明PYTHONPATH指向该位置

* 插件wrapper.py 当前只支持非流式

* 插件wrapper.py 当前只支持返回,需要用户显式在插件中定义

```json
{"key": "boxes", "data": rlt, "len": len(rlt), "status": 3, "type": 0}
```
该结构数据

* 上条示例中, rlt 当前只支持单层数据，即不支持object嵌套，如果用户返回json，此处需要先把rlt dumps成文本再返回 

* 理论上用户除了上传 wrapper.py 以及相关依赖之外，还需要提供一些模型文件，这些文件比较大，一般不在Dockerfile中构建入镜像，会导致git代码库十分庞大,当前示例的的yolov5和 mmocr均在 wrapper init的时候下载模型

建议统一方式模型资源使用方式， 根据不同场景支持

* 用户下载模型，比如用户在代码wrapper_init中下载模型，需要用户提前准备好模型下载链接

* 平台提前挂载资源到指定目录如 /models，需要用户提前上传

#### 基础镜像构建(GPU)

***基础镜像仅在特殊需求时(如对cuda，python版本有要求时才需要重新构建,一般用户仅需关注构建业务镜像)***

1. cuda-go-python基础镜像,用于编译aiges项目的基础镜像，参见 [官方仓库](https://gitlab.com/nvidia/container-images/cuda/)，本仓 库引用了部分版本，存放于 [docker/gpu/cuda](docker/gpu/cuda)中
基础镜像当前基于 nvidia/cuda 官方的基础镜像作为base镜像 如 [cuda-10.1](docker/gpu/base/cuda-10.1)中所示: aiges基础镜像基于 ***形如 nvidia/cuda:10.1-devel-ubuntu18.04*** 构建

2. 基于 [cuda-10.1](docker/gpu/base/cuda-10.1) 已构建出 ***artifacts.iflytek.com/docker-private/atp/cuda-go-python-base:10.1-1.17-3.9.13-ubuntu1804***

3. aiges: 基于 [aiges-dockerifle](docker/gpu/aiges/ubuntu1804/Dockerfile) 构建出gpu的  aiges基础镜像

构建命令:

``docker buildx build -f docker/gpu/base/cuda-10.2/Dockerfile -t artifacts.iflytek.com/docker-private/atp/cuda-go-python-base:10.2-1.17-3.9.13-ubuntu1804  . --push``

***当前支持的cuda-go-python基础镜像列表(包含cuda go python编译环境)***

| repo                                                                                       | tag                           | python | cuda | os           |
|--------------------------------------------------------------------------------------------|-------------------------------|--------|------|--------------|
| artifacts.iflytek.com/docker-private/atp/cuda-go-python-base:10.1-1.17-3.9.13-ubuntu1804   | 10.1-1.17-3.9.13-ubuntu1804   | 3.9.13 | 10.1 | ubuntu 18.04 |
| artifacts.iflytek.com/docker-private/atp/cuda-go-python-base:10.2-1.17-3.9.13-ubuntu1804   | 10.2-1.17-3.9.13-ubuntu1804   | 3.9.13 | 10.2 | ubuntu 18.04 |
| artifacts.iflytek.com/docker-private/atp/cuda-go-python-base:11.6.2-1.17-3.9.13-ubuntu1804 | 11.6.2-1.17-3.9.13-ubuntu1804 | 3.9.13 | 11.6 | ubuntu 18.04 |

***当前支持的aiges基础镜像列表***

| repo                                                                      | tag                    | python | cuda | os           |
|---------------------------------------------------------------------------|------------------------|--------|------|--------------|
| artifacts.iflytek.com/docker-private/atp/aiges-gpu:10.1-3.9.13-ubuntu1804 | 10.1-3.9.13-ubuntu1804 | 3.9.13 | 10.1 | ubuntu 18.04 |
| artifacts.iflytek.com/docker-private/atp/aiges-gpu:10.2-3.9.13-ubuntu1804 | 10.2-3.9.13-ubuntu1804 | 3.9.13 | 10.2 | ubuntu 18.04 |

***构建命令***:

1. 使用buildx:
   - `docker buildx build  -f docker/gpu/aiges/ubuntu1804/Dockerfile . -t artifacts.iflytek.com/docker-private/atp/aiges-gpu:10.1-3.9.13-ubuntu1804`

2. 使用docker build
   - `docker buildx build  -f docker/gpu/aiges/ubuntu1804/Dockerfile . -t artifacts.iflytek.com/docker-private/atp/aiges-gpu:10.1-3.9.13-ubuntu1804`

3. 使用buildah
   - `buildah build  -f docker/gpu/aiges/ubuntu1804/Dockerfile . -t artifacts.iflytek.com/docker-private/atp/aiges-gpu:10.1-3.9.13-ubuntu1804`

#### 业务镜像构建方法

业务镜像需要基于 aiges基础镜像进行构建，用户可在此过程定制 python的依赖项目以及用户自研项目

参考示例:

[mmocr](demo/mmocr/Dockerfile_gpu)


### build.py使用

TODO 计划开发build.py 支持多cuda版本多python版本的基础镜像自动构建pipeline









