# AIges 多版本Docker镜像构建


## Cuda Cudnn版本

当前 宿主机默认统一提供 10.1


## OS版本

10.1 当前支持的os 推荐ubuntu18.04


## Python版本

Python参考官方Dockerfile build


## CPU or GPU？

建议统一使用 GPU镜像，通过环境变量切换是否使用GPU卡，牺牲点空间

## Base镜像

考虑到使用Go语言和 CGO C++语言，以及Python语言，初期统一提供 ubuntu18.04-go作为最基础镜像







