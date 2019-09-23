

# 部署说明

## Prerequisites

部署好配置中心环境：包括zookeeper集群、companion、配置中心网站等组件。



## 部署步骤

1、在配置中心创建项目、集群、服务、服务版本，比如AIaaS、aitest、atmos-svc、2.0.0

2、在配置中心上传配置文件：
上传的路径 参考：AIaaS / aitest / atmos-svc / 2.0.0
需要上传的配置文件列表为：

1）atmos.toml  atmos组件相关配置

2）xsfc.toml  xsf客户端配置配置

3）xsfs.toml  xsf服务端配置

示例配置文件参考配置文件示例目录下的文件
上传后，根据实际环境，来修改相关配置

3、基于docker部署

- 获取最新docker镜像版本

  比如：172.16.59.153/aiaas/atmos:2.0.0

- docker启动参考命令

  ```
  docker run -d  --net="host" --name atmos -v /data/atmos/logs:/log/server  172.16.59.153/aiaas/atmos:2.0.0 ../atmos -m 1 -p AIaaS -g aitest -s atmos-svc -u http://10.1.86.223:9080
  
  启动命令 -m 1指定是以配置中心模式启动
  -p 指定项目名称
  -g 指定集群名称
  -s 指定服务名称
  -u 指定配置中心companion地址
  ```

  

4、基于物理部署

```
    ../atmos -m 1 -p AIaaS -g aitest -s atmos-svc -u http://10.1.86.223:9080
```





