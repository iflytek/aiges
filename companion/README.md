# **companion**
## **companion简介**
​	companion是java语言开发的，是配置中心的服务代理层。提供网站操作zookeeper的接口、SDK获取zookeeper集群地址的接口，以及服务路径、配置变更的回调接口。

## **源码构建**

- 安装好maven环境、java环境

- 进入项目目录，执行命令：mvn install -U  -Dmaven.test.skip=true

- 也可以基于docker进行构建，项目中提供了docker构建的脚本，执行即可

## **部署**

获取最新的镜像，然后基于镜像进行部署，启动命令可以如下命令：

```
docker run -d --name companion --net="host" -v /opt/polaris/companion/logs:/log/server -v /etc/localtime:/etc/localtime XXXXX/develop/companion:2.0.3 sh watchdog.sh -h10.1.86.211 -p6868 -z10.X.86.211:2181,10.X.86.70:2181,10.X.86.212:2181 -whttps://10.1.87.69:8095

注：启动companion参数分别代表companion的ip地址-h10.1.86.211，端口-p6868
zk集群地址-z10.X.86.211:2181,10.X.86.70:2181,10.X.86.212:2181
cynosure地址-whttps://10.1.87.69:8095
```



## **常见问题**

