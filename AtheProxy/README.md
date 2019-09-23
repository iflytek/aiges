# **atmos**
## **atmos简介**

​	atmos接收来自webgate接入层的请求，并根据请求的路由参数，请求负载均衡模块lb，获取引擎地址，将请求发送到引擎，并从引擎获取结果。

​	atmos 跟webgate交互的协议称为业务层协议。相对于webgate,atmos是服务端。

​	atmos 跟引擎交互的协议称为引擎协议。相对于引擎,atmos是客户端。

​	atmos交互模式分为**同步和异步**两种。**同步**是指webgate请求atmos，atmos同步请求引擎，并将引擎的处理结果，同步返回给webgate。**异步**是指webgate请求atmos，atmos立即返回一个响应给webgate，并同时调用引擎，如果引擎有处理结果或者错误返回，则回调webgate，把结果返回给webgate。如果引擎没有处理结果，则不会回调webgate。

​	而每种模式下，还支持**会话模式和once模式**。会话模式是指一次完整的业务功能，需要多次交互。比如语音识别，可以流式的将语音送到引擎。而引擎也可以将识别结果，随时返回给调用者。会话模式下，每次会话的所有请求，都会把请求发送同一个引擎实例。once模式是指一次请求把所有请求数据发送到引擎，并且引擎会一次性返回完整的处理结果。



## **atmos核心配置文件简介**

atmos需要的配置文件主要包含：**atmos.toml、xsfc.toml、xsfs.toml**。

- ### **atmos.toml是atmos组件的核心配置文件**


    ```
    [atmos-common]
    debugSwitch=0 #是否开启pprof，默认使用端口号8089 0表示关，1表示开
    subTypes="svc" #支持的sub类型，多个sub类型之间用逗号分隔，比如"svc,ocr",然后每个sub，要单独配置一组基础配置，规则是:atmos-sub，例如sub是svc，则为atmos-svc
    
    [atmos-svc]
    sub="svc"  #业务类型
    lb="lbv2" #lb名称，负载均衡的
    defaultTimeout=3000 #默认超时时间
    engineTimeout=3000  #调用引擎接口的超时时间 
    getEngineResultTimeout=15 #单位秒，获取引擎结果的时候，最大等待时间(可能包含多次引擎接口调用 )
    engineRetry = 2 #调用引擎失败的重试次数，默认值为3
    mockSwitch=1 #0表示mock，1表示正常
    
    [log]
    level = "warn"  #日志等级。可设置参数debug、info、warn、error。缺省warn。
    file = "/log/server/atmos.log"    #业务日志路径。
    size = 100      #日志大小。单位MB。缺省100
    count = 20      #日志文件留存数量。缺省20。
    die = 10        #日志文件的有效期。单位天。缺省10。
    async = 1       #是否启用异步日志。1是0否。缺省1。
    
    ```
    [atmos-common] 标签配置的是一些公共属性。主要包括是否开启pprof，支持的sub 业务类型。一个实例支持多种业务类型。支持的业务类型，需要配置，具体业务类型的相关属性。下面以sub svc为例子。
    
    配置具体的业务类型属性的格式为：[atmos-svc] ，然后可以配置svc这个业务的sub，负载均衡器服务名称。默认超时时间、单次调用引擎的接口超时时间、获取引擎结果的超时时间、重试次数等。
    
    [log] 日志信息配置。

- ### **xsfc.toml   客户端配置**

  主要包括客户端启动参数、重试、日志等的配置。

    ```
      [atmos-svc]#需要跟服务名保持一致
      conn-timeout = 2000 #连接超时时间
      lb-mode= 2  #lb模式
      lb-retry = 2  #重试次数
      [trace]
      ip = "127.0.0.1"
      able = 0
      dump = 1
      [log]
      level = "info"
      file = "/log/server/xsfc.log"
      size = 100
      count = 20
      die = 30
      async = 0
    ```

- ### **xsfs.tom 服务端配置**

    ```
    #服务自身的配置
    #注意此section名需对应bootConfig中的service
    [atmos-svc]#需要跟服务名保持一致
    
    #host = "172.16.154.172"#若host为空，则取netcard对应的ip，若二者均为空，则取hostname对应的ip
    #netcard = "eth0"
    
    #port = 50061  端口号如果不指定，则使用随机端口
    finder = 1 #缺省0
    
    #trace日志所用
    [trace]#已做缺省处理
    #trace收集服务的地址
    host = "127.0.0.1" #缺省127.0.0.1
    #trace收集服务的端口
    port = 4546 #缺省4545
    #trace服务的协程数
    backend = 1 #缺省4
    #是否将日志写入到远端
    deliver = 0 #缺省1
    #是否将日志 落入磁盘
    dump = 1 #缺省0
    #是否禁用trace
    able = 0 #缺省1
    
    [log]#已做缺省处理
    level = "info" #缺省warn
    file = "/log/xsfs.log" #缺省xsfs.log
    #日志文件的大小，单位MB
    size = 300 #缺省10
    #日志文件的备份数量
    count = 3 #缺省10
    #异步日志
    async = 0 #缺省异步
    
    ```

## **源码构建**

- golang版本要求1.9以上

- 设置好相应的GOPATH

- 下载 源码。

- 进入项目根目录，执行：go build -v -o atmos main.go

  

## **部署**

部署可以参考文档:doc/INSTALL.md

## **常见问题**

1、配置文件未推送到配置中心zookeeper，导致启动失败