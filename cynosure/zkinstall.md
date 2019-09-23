## 基础环境
docker版本1.12.0，请自行安装。

#### 拉取镜像 
````
sudo docker pull 172.16.59.153/develop/zookeeper.3.4.12:1.1
````
#### 机器上创建根目录，用于挂载到docker容器中
````
mkdir -p /opt/polaris/zookeeper
````
#### 在zookeeper目录下创建logs、conf、data、datalog目录。conf 目录用于存放zk配置文件，logs目录用于存放日志,data目录存放myid
````
mkdir mkdir -p /opt/polaris/zookeeper/logs
mkdir mkdir -p /opt/polaris/zookeeper/data
mkdir mkdir -p /opt/polaris/zookeeper/datalog
mkdir mkdir -p /opt/polaris/zookeeper/conf
chmod +777 /opt/polaris/zookeeper/logs
````
#### 在第一台机器上的conf 目录下创建zoo.cfg作为zk的启动配置如下server1,server2,server3的地址为集群的地址，请自行修改
zookeeper/conf/zoo.cfg:
````
clientPort=2181
dataDir=/data
dataLogDir=/datalog
tickTime=2000
initLimit=5
syncLimit=2
server.1=10.1.86.211:2888:3888
server.2=10.1.86.70:2888:3888
server.3=10.1.86.212:2888:3888
autopurge.snapRetainCount=3
autopurge.purgeInterval=1
````
#### 在conf目录下再创建log4j.properties,内容如下
````
# Define some default values that can be overridden by system properties
zookeeper.root.logger=INFO, CONSOLE
zookeeper.console.threshold=INFO
zookeeper.log.dir=.
zookeeper.log.file=zookeeper.log
zookeeper.log.threshold=DEBUG
zookeeper.tracelog.dir=.
zookeeper.tracelog.file=zookeeper_trace.log

#
# ZooKeeper Logging Configuration
#

# Format is "<default threshold> (, <appender>)+

# DEFAULT: console appender only
log4j.rootLogger=${zookeeper.root.logger}

# Example with rolling log file
#log4j.rootLogger=DEBUG, CONSOLE, ROLLINGFILE

# Example with rolling log file and tracing
#log4j.rootLogger=TRACE, CONSOLE, ROLLINGFILE, TRACEFILE

#
# Log INFO level and above messages to the console
#
log4j.appender.CONSOLE=org.apache.log4j.ConsoleAppender
log4j.appender.CONSOLE.Threshold=${zookeeper.console.threshold}
log4j.appender.CONSOLE.layout=org.apache.log4j.PatternLayout
log4j.appender.CONSOLE.layout.ConversionPattern=%d{ISO8601} [myid:%X{myid}] - %-5p [%t:%C{1}@%L] - %m%n

#
# Add ROLLINGFILE to rootLogger to get log file output
#    Log DEBUG level and above messages to a log file
log4j.appender.ROLLINGFILE=org.apache.log4j.RollingFileAppender
log4j.appender.ROLLINGFILE.Threshold=${zookeeper.log.threshold}
log4j.appender.ROLLINGFILE.File=${zookeeper.log.dir}/${zookeeper.log.file}

````
#### data 目录下创建myid文件内容如下：myid 是机器的标识，每一台机器的myid值一定要不一样。
zookeeper/data/myid:
````
1
````

#### 创建启动脚本start.sh
zookeeper/start.sh
````
sudo docker run -d --name zookeeper1 --net="host"  -v /opt/polaris/zookeeper/logs:/zookeeper-3.4.12/logs  -v /opt/polaris/zookeeper/conf:/conf -v /opt/polaris/zookeeper/data:/data -v /opt/polaris/zookeeper/datalog:/datalog -e "ZOO_MY_ID=1" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" 172.16.59.153/develop/zookeeper.3.4.12:1.1
````

#### 另外的两台机器的配置文件和第一台机器一样，注意修改data/myid 的值

#### server2的 data/myid 
````
2
```` 

#### server2 的启动脚本
````
sudo docker run -d --name zookeeper2 --net="host"  -v /opt/polaris/zookeeper/logs:/zookeeper-3.4.12/logs  -v /opt/polaris/zookeeper/conf:/conf -v /opt/polaris/zookeeper/data:/data -v /opt/polaris/zookeeper/datalog:/datalog -e "ZOO_MY_ID=2" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" 172.16.59.153/develop/zookeeper.3.4.12:1.1
````

#### server3的 data/myid 
````
3
```` 

#### server3 的启动脚本
````
sudo docker run -d --name zookeeper3 --net="host"  -v /opt/polaris/zookeeper/logs:/zookeeper-3.4.12/logs  -v /opt/polaris/zookeeper/conf:/conf -v /opt/polaris/zookeeper/data:/data -v /opt/polaris/zookeeper/datalog:/datalog -e "ZOO_MY_ID=3" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" 172.16.59.153/develop/zookeeper.3.4.12:1.1
````
#### 全部启动完成后，进入容器里面，检查一下zookeeper集群的状态。具体步骤为：
```
docker exec -ti zookeeper1 bash

cd bin

./zkServer.sh status
````
## 单机部署集群
#### 单机部署集群时端口号要不一样，在polaris 目录下创建以下目录
````
mkdir /opt/polaris/zookeeper1
mkdir /opt/polaris/zookeeper2
mkdir /opt/polaris/zookeeper3
````
#### 在每个zookeeper目录下创建 logs、conf、data、datalog目录
````
mkdir mkdir -p /opt/polaris/zookeeper1/logs
mkdir mkdir -p /opt/polaris/zookeeper1/data
mkdir mkdir -p /opt/polaris/zookeeper1/datalog
mkdir mkdir -p /opt/polaris/zookeeper1/conf
chmod +777 /opt/polaris/zookeeper1/logs

mkdir mkdir -p /opt/polaris/zookeeper2/logs
mkdir mkdir -p /opt/polaris/zookeeper2/data
mkdir mkdir -p /opt/polaris/zookeeper2/datalog
mkdir mkdir -p /opt/polaris/zookeeper2/conf
chmod +777 /opt/polaris/zookeeper2/logs

mkdir mkdir -p /opt/polaris/zookeeper3/logs
mkdir mkdir -p /opt/polaris/zookeeper3/data
mkdir mkdir -p /opt/polaris/zookeeper3/datalog
mkdir mkdir -p /opt/polaris/zookeeper3/conf
chmod +777 /opt/polaris/zookeeper3/logs
````

#### 创建conf目录下的配置文件zoo.cfg,log4j.properties
zookeeper1/conf/zoo.cfg:
````
clientPort=2181
dataDir=/data
dataLogDir=/datalog
tickTime=2000
initLimit=5
syncLimit=2
server.1=10.1.86.211:2888:3888
server.2=10.1.86.70:2788:3788
server.3=10.1.86.212:2688:3688
autopurge.snapRetainCount=5
autopurge.purgeInterval=2
````
zookeeper2/conf/zoo.cfg:
````
clientPort=2182
dataDir=/data
dataLogDir=/datalog
tickTime=2000
initLimit=5
syncLimit=2
server.1=10.1.86.211:2888:3888
server.2=10.1.86.70:2788:3788
server.3=10.1.86.212:2688:3688
autopurge.snapRetainCount=5
autopurge.purgeInterval=2
````
zookeeper3/conf/zoo.cfg:
````
clientPort=2183
dataDir=/data
dataLogDir=/datalog
tickTime=2000
initLimit=5
syncLimit=2
server.1=10.1.86.211:2888:3888
server.2=10.1.86.70:2788:3788
server.3=10.1.86.212:2688:3688
autopurge.snapRetainCount=5
autopurge.purgeInterval=2
````
#### data目录下的myid文件分别修改为 1，2，3
zookeeper1/data/myid:
````
1
````
zookeeper2/data/myid:
````
2
````
zookeeper3/data/myid:
````
3
````

#### 启动脚本
zookeeper1/start.sh:
````
sudo docker run -d --name zookeeper1 --net="host"  -v /opt/polaris/zookeeper1/logs:/zookeeper-3.4.12/logs  -v /opt/polaris/zookeeper/conf:/conf -v /opt/polaris/zookeeper/data:/data -v /opt/polaris/zookeeper/datalog:/datalog -e "ZOO_MY_ID=1" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" 172.16.59.153/develop/zookeeper.3.4.12:1.1
````
zookeeper2/start.sh:
````
sudo docker run -d --name zookeeper2 --net="host"  -v /opt/polaris/zookeeper1/logs:/zookeeper-3.4.12/logs  -v /opt/polaris/zookeeper/conf:/conf -v /opt/polaris/zookeeper/data:/data -v /opt/polaris/zookeeper/datalog:/datalog -e "ZOO_MY_ID=2" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" 172.16.59.153/develop/zookeeper.3.4.12:1.1
````
zookeeper3/start.sh:
````
sudo docker run -d --name zookeeper3 --net="host"  -v /opt/polaris/zookeeper1/logs:/zookeeper-3.4.12/logs  -v /opt/polaris/zookeeper/conf:/conf -v /opt/polaris/zookeeper/data:/data -v /opt/polaris/zookeeper/datalog:/datalog -e "ZOO_MY_ID=3" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" 172.16.59.153/develop/zookeeper.3.4.12:1.1
````

