## 简介 ##

### 赫尔墨斯（Hermes） ###

![hermes](res/hermes.jpg)

- hermes，取名于赫尔墨斯，为众神使者，意在为各服务指引方向
- 赫尔墨斯（Hermes），希腊神话人物，十二主神之一，宙斯与迈亚的儿子,是众神使者

## 功能 ##

### 一句话 ###

- 全局负载，基于收集的信息做负载，并按特定策略返回相应的节点

### 详细 ###

> 当前仅实现了基于授权数的个性化负载策略，相关功能列举如下

- 数据同步
	- 同步DB数据并至本地
	- 同步新节点数据至DB

- 节点上报
	- 引擎段的分配
	- 节点生命周期的维护

- 节点请求
	- 根据负载返回响应节点
	- 根据uid返回响应节点

- 个性化索引
	- 同步rmq至本地
	- 根据策略路由至后端节点

- 监控相关
	- 监控全量节点
	- 监控局部节点
	
## 模块 ##

- hermes采用模块化设计，以此支持不同的策略
- 当前仅实现了基于授权实例数的个性化负载

## 接口 ##

### 上报接口 ###

```
req.SetParam("svc", "svc)				//大业务名
req.SetParam("subsvc", "sms")			//子业务名
req.SetParam("addr", "x.x.x.x:xxxx")	//上报节点的服务地址
req.SetParam("total", "xxx")			//节点的总授权量
req.SetParam("best", "xxx)				//节点的最佳授权量
req.SetParam("idle", "xxx")				//节点的空闲授权量
req.SetParam("live", "xxx")				//节点的存活状态(若此值为0，则代表节点主动下线)
```

### 请求接口 ###

```
req.SetParam("nbest", "xxx")			//期望获取的节点数
req.SetParam("svc", "svc")				//大业务名
req.SetParam("subsvc", "svc")			//子业务名
req.SetParam("all", "1")				//是否返回所有的节点
req.SetParam("uid", "xxxx")				//universal id 全局id（个性化所需）
req.SetParam("sg", "xxxx")				//strategy（负载策略）
```

## 流程 ##

> 补充中。。。

## 编译 ##

### 本地编译 ###

- golang版本要求1.9以上

- 设置好相应的GOPATH(如:export GOPATH=hermes)

- 下载 tag版本并进入hermes\src\lbv2 

- 在lb目录下,执行cmd：go build -v -o lbv2

### docker ###

> 构建中。。。

## 部署 ##

- 配置中心

	- ./lbv2 -m 1 -c lbv2.toml -p 3s -s lbv2 -u http://10.1.86.223:9080 -g 3s
	
```
注释:
    -m 制定运行方式(0-本地配置,1-配置中心)    
    -c 配置文件名(Native模式时，用本地的配置文件，Center模式时使用配置中心的配置文件)    
    -p 项目名    
    -s 启动的服务名(注意:配置文件，配置中心的服务名需一致)  
    -u 配置中心地址    
    -g 配置项目组
```

## 自测 ##

[hermes_test_data.xlsx](hermes_test_data.xlsx)

## FAQ ##

> 补充中。。。
