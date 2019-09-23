## BootConfig对象简介
```
private String companionUrl;// companion组件的url，http://127.0.0.1:7091
private String cachePath; // 缓存路径
private boolean configCache = true;//配置是否使用本地缓存
private boolean serviceCache = true;//配置是否使用本地缓存

```
## ServiceMeteData对象简介

    /**
     * 项目名称
     */
    private String project;

    /**
     * 集群名称
     */
    private String group;

    /**
     * 服务名称
     */
    private String service;

    /**
     *  组件版本
     */
    private String version;

    /**
     *  组件唯一标识（一般场景推荐ip：port）
     */
    private String address;

## 主要接口列表
```
//sdk管理对象
public FinderManager()

//sdk初始化
public CommonResult init(final BootConfig bootConfig)

//服务注册接口
public CommonResult registerService(String apiVersion)

//服务注册接口
public CommonResult registerService(String addr, String apiVersion)

//取消服务注册接口
public CommonResult unRegisterService(String apiVersion)

//取消服务注册接口
public CommonResult unRegisterService(String addr, String apiVersion)

//获取并订阅服务变更
public CommonResult<List<Service>> useAndSubscribeService(List<SubscribeRequestValue> requestValueList, ServiceHandle serviceHandle)

//取消服务变更的订阅
public CommonResult unSubscribeService(SubscribeRequestValue requestValue) 

//获取并订阅配置变更
public CommonResult<List<Config>> useAndSubscribeConfig(List<String> configNameList, ConfigChangedHandler configChangedHandler) 

//取消文件配置变更的订阅
public CommonResult unSubscribeConfig(String configName)
```
## maven地址
```
   <dependency>
         <groupId>com.iflytek.ccr.polaris</groupId>
         <artifactId>finder-java</artifactId>
         <version>2.0.0</version>
   </dependency>
   
    <mirror>
         <id>nexus_public</id>
         <mirrorOf>*,!lib</mirrorOf>
         <name>nexus_public_repo</name>
         <url>http://172.16.59.13:20221/nexus/content/groups/public/</url>
       </mirror>
```        

## 集成说明
首先需要在配置中心的网站，创建基础数据。这个可以联系李玉功。
之后就是代码的集成了。
无论是集成服务发现还是集成配置管理，第一步都是相同的。

### 1、初始化FinderManager对象

#### 1）初始化服务元数据ServiceMeteData对象

project：项目名字
group：集群名称

service：服务名字

version：服务版本号

address：服务地址（ip:port）,ip和端口号的字符串，能够唯一定位到该服务。

public ServiceMeteData(String project, String group, String service, String version, String address)

#### 2）初始化BootConfig对象，可以设置如下属性
    private String companionUrl;//companion 地址
    private String cachePath;//缓存的路径，该缓存sdk内部使用
    private long tickerDuration; //可以不设置
    private boolean configCache = true;// 配置订阅的时候，如果无法从服务端获取到配置信息，是否使用缓存的配置信息，默认为true，使用缓存。
    private boolean serviceCache = true;//服务订阅的时候，如果无法从服务端获取到服务信息，是否使用缓存的服务信息，默认为true，使用缓存。
    private int zkSessionTimeout;//zk session超时时间，单位毫秒 可以不设置
    private int zkConnectTimeout;//zk 连接超时时间，单位毫秒 可以不设置
    private int zkMaxSleepTime; //zk 连接最大休眠时间 可以不设置
    private int zkMaxRetryNum; //zk连接最大重试次数 可以不设置
    private ServiceMeteData meteData;//服务元数据对象
    
#### 3)初始化FinderManager对象：FinderManager()，调用init()方法进行初始化
由于FinderManager init()，需要调用外部接口，存在失败的可能性。如果抛出异常，需要进行处理。可以尝试再次进行初始化。

### 2、注册服务
有两个方法：一个带地址参数，一个不带地址参数。地址是指该服务的地址加端口号（ip：port）。如果调用不带参数的注册，则使用初始化时ServiceMeteData中的地址。apiVersion是两个组件交互的api的版本号，跟组件自身的版本没有关系。
```
//服务注册接口
public CommonResult registerService(String apiVersion)
//服务注册接口
public CommonResult registerService(String addr, String apiVersion)
```
需要对返回的CommonResult 进行解析。判断是否注册成功了。如果没有注册成功，可以尝试重新注册。<br>
CommonResult ret 属性为0，则表明注册成功。

### 3、订阅服务
```
//获取并订阅服务变更
public CommonResult<List<Service>> useAndSubscribeService(List<SubscribeRequestValue> requestValueList, ServiceHandle serviceHandle)
```
参数说明：List<SubscribeRequestValue> 参数是服务对象的集合。可以一次订阅多个服务。<br>
ServiceHandle serviceHandle 参数，这个是回调函数接口。需要实现该接口。<br>
CommonResult<List<Service>> 返回值。这个是当前服务的集合。<br>
获取该服务后，如果服务相关信息再有变更的话，serviceHandle 回调会被触发。<br>

### 4、订阅配置变更
 ```
public CommonResult<List<Config>> useAndSubscribeConfig(List<String> configNameList, ConfigChangedHandler configChangedHandler)
```
参数说明：List<String> configNameList 该参数是配置文件名称的集合。可以一次订阅多个配置文件。<br>
ConfigChangedHandler configChangedHandler 是配置变更的回调接口。<br>
CommonResult<List<Config>> 返回值，返回当前订阅配置文件的当前信息。

### 5、取消服务的订阅
```
//取消服务变更的订阅
public CommonResult unSubscribeService(SubscribeRequestValue requestValue) 
```
参数说明：
serviceName 服务名称

### 6、取消订阅配置
```
public CommonResult unSubscribeConfig(String configName)
```
configName 该参数是配置文件名称

### 7、注销服务
```
  //取消服务注册接口
  public CommonResult unRegisterService(String apiVersion)
  
  //取消服务注册接口
  public CommonResult unRegisterService(String addr, String apiVersion)
  ```
  如果调用不带地址参数的注销服务，则使用初始化时ServiceMeteData中的地址。
