# Athena Serving

<!-- markdownlint-capture -->
<!-- markdownlint-disable MD033 -->

<span class="badge-placeholder">[![Build Status](https://img.shields.io/drone/build/thegeeklab/hugo-geekdoc?logo=drone&server=https%3A%2F%2Fdrone.thegeeklab.de)](https://drone.thegeeklab.de/thegeeklab/hugo-geekdoc)</span>
<span class="badge-placeholder">[![GitHub release](https://img.shields.io/github/v/release/xfyun/AthenaServing)](https://github.com/xfyun/AthenaServing/releases/latest)</span>
<span class="badge-placeholder">[![GitHub contributors](https://img.shields.io/github/contributors/xfyun/AthenaServing)](https://github.com/xfyun/AthenaServing/graphs/contributors)</span>
<span class="badge-placeholder">[![License: Apache2.0](https://img.shields.io/github/license/xfyun/AthenaServing)](https://github.com/xfyun/AthenaServing/blob/master/LICENSE)</span>

<!-- markdownlint-restore -->
## 整体架构

![img](https://raw.githubusercontent.com/xfyun/proposals/main/athenaloader/athena.png)

## AI General Engine Service (AIges)

通用引擎加载器(部分文档中loader， loader engine均为别名)

### Documents

View Doc on [Documentation](https://xfyun.github.io/inferservice/architechture/architechture/)


## Features

-[X] 支持模型推理成RPC服务(Serving框架会转成HTTP服务)
-[X] 支持C代码推理 support c++/c code infer
-[X] 支持Python代码推理 Support python code infer
-[X] 支持配置中心，服务发现
-[X] 支持负载均衡配置
-[] 支持Java代码推理或者其它
-[] 支持计量授权
 




