# Athena Serving

<!-- markdownlint-capture -->
<!-- markdownlint-disable MD033 -->

<span class="badge-placeholder">[![Build Status](https://img.shields.io/drone/build/thegeeklab/hugo-geekdoc?logo=drone&server=https%3A%2F%2Fdrone.thegeeklab.de)](https://drone.thegeeklab.de/thegeeklab/hugo-geekdoc)</span>
<span class="badge-placeholder">[![GitHub release](https://img.shields.io/github/v/release/xfyun/aiges)](https://github.com/xfyun/aiges/releases/latest)</span>
<span class="badge-placeholder">[![License: Apache2.0](https://img.shields.io/github/license/xfyun/aiges)](https://github.com/xfyun/aiges/blob/master/LICENSE)</span>

<!-- markdownlint-restore -->
## 整体架构

![img](https://raw.githubusercontent.com/xfyun/proposals/main/athenaloader/athena.png)

## 使用工作流

![img](https://github.com/xfyun/proposals/blob/main/athenaloader/usage.png?raw=true)

## 开源计划

- [x] [加载器](#AI General Engine Service (AIges)) 通用引擎/模型加载器
- [x] [负载均衡器](#LoadBalance) 负载聚合组件
- [ ] [WebGate](#Webgate) Web网关组件
- [x] [Polaris](#Polaris) 配置中心与服务发现
- [ ] [Atom](#Atom) 协议转换组件
- [ ] Serving on Kubernetes Helm Chart一键部署 (进行中)
- [ ] Serving on Docker with docker-compose 一键部署
- [ ] 各组件Documentation建设 (进行中)
- [ ] 多领域模型Demo演示示例、GIF (进行中)



## AI General Engine Service (AIges)

通用引擎加载器(部分文档中loader， loader engine均为别名)

### Documentation

View Doc on [Documentation](https://xfyun.github.io/inferservice/architechture/architechture/)


## Features


- [x] 支持模型推理成RPC服务(Serving框架会转成HTTP服务)
- [x] 支持C代码推理 support c++/c code infer
- [x] 支持Python代码推理 Support python code infer
- [x] 支持配置中心，服务发现
- [x] 支持负载均衡配置
- [ ] 支持Java代码推理或者其它
- [ ] 支持计量授权

## Polaris

Go编写的服务发现、配置中心

### Documentation

View Doc on [Documentation](https://xfyun.github.io/inferservice/architechture/architechture/)

## Features
- Full category for managing service
- Support multi service versions
- Support roll back config
- Support feedback for pushing config
- Support management for provider and consumer online
- High available for some not expected cases
- Easy Integration
- Support Delivery by docker 



