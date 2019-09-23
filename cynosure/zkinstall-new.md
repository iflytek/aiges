
### 基础环境
docker版本1.12.0以上，请自行安装。

#### 拉取镜像 
````
sudo docker pull 172.16.59.153/develop/zookeeper:3.4.12_private
//该镜像是基于官方的镜像基础构建。增加了通过命令行的添加环境变量的方式传入配置文件参数，从而不用挂载配置文件，减少了部署难度
//镜像详细的信息可以参考Dockerfile 和docker-entrypoint.sh
````
#### 机器上创建根目录，用于挂载到docker容器中
````
mkdir -p /opt/polaris/zookeeper
````
#### 在zookeeper目录下创建logs、data、datalog,conf目录,并在conf目录下创建log4j.properties和zoo.cfg
````
mkdir mkdir -p /opt/polaris/zookeeper/logs
mkdir mkdir -p /opt/polaris/zookeeper/data
mkdir mkdir -p /opt/polaris/zookeeper/datalog
chmod 644 /opt/polaris/zookeeper/logs
````
zookeeper/conf/zoo.cfg:
````
clientPort=2183
dataDir=/data
dataLogDir=/datalog
tickTime=2000
initLimit=5
syncLimit=2
server.1=10.1.86.211:2888:3888
server.2=10.1.86.211:2788:3788
server.3=10.1.86.211:2688:3688
autopurge.snapRetainCount=3
autopurge.purgeInterval=1
````
conf/log4j.properties:
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
#### 创建启动脚本start.sh，
  脚本中的severIP请自行修改，环境变量的值也可以根据需要修改，环境变量ZOO_MY_ID指定的是每台机器的唯一标识，每个脚本都要不一样。ZOO_AUTOPURGE_SNAPRETAINCOUNT指定autopurge.snapRetainCount，这个参数指定了需要保留的文件数目，默认保留3个；ZOO_AUTOPURGE_PURGEINTERVAL指定　　autopurge.purgeInterval这个参数指定了清理频率，单位是小时，需要填写一个1或者更大的数据，默认0表示不开启自动清理功能。
ZOO_SERVERS指定了集群服务器的IP地址
                                                                                                                         
                                                                                                                         　　autopurge.purgeInterval这个参数指定了清理频率，单位是小时，需要填写一个1或者更大的数据，默认0表示不开启自动清理功能。
zookeeper/start.sh
````
sudo docker run -d --name zookeeper --net="host"  -v /opt/polaris/zookeeper/conf:/conf -v /opt/polaris/zookeeper/logs:/zookeeper-3.4.12/logs -v /opt/polaris/zookeeper/data:/data -v /opt/polaris/zookeeper/datalog:/datalog -e "ZOO_MY_ID=1" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" -e "ZOO_AUTOPURGE_SNAPRETAINCOUNT=3" -e "ZOO_AUTOPURGE_PURGEINTERVAL=1" 172.16.59.153/develop/zookeeper:3.4.12_private
````

#### server2 的启动脚本
````
sudo docker run -d --name zookeeper --net="host"  -v /opt/polaris/zookeeper/conf:/conf -v /opt/polaris/zookeeper/logs:/zookeeper-3.4.12/logs -v /opt/polaris/zookeeper/data:/data -v /opt/polaris/zookeeper/datalog:/datalog -e "ZOO_MY_ID=2" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" -e "ZOO_AUTOPURGE_SNAPRETAINCOUNT=3" -e "ZOO_AUTOPURGE_PURGEINTERVAL=1"  172.16.59.153/develop/zookeeper:3.4.12_private
````

#### server3 的启动脚本
````
sudo docker run -d --name zookeeper --net="host"  -v /opt/polaris/zookeeper/conf:/conf -v /opt/polaris/zookeeper/logs:/zookeeper-3.4.12/logs -v /opt/polaris/zookeeper/data:/data -v /opt/polaris/zookeeper/datalog:/datalog -e "ZOO_MY_ID=3" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" -e "ZOO_AUTOPURGE_SNAPRETAINCOUNT=3" -e "ZOO_AUTOPURGE_PURGEINTERVAL=1"  172.16.59.153/develop/zookeeper:3.4.12_private
````

#### 在data目录下创建myid，作为每一台机器的唯一标识，每台机器的myid都要不一样
````
1
````

#### 全部启动完成后，进入容器里面，检查一下zookeeper集群的状态。具体步骤为：
````
docker exec -ti zookeeper1 bash

cd bin

./zkServer.sh status
````
### 单机部署集群,需要在zoo.cfg中指定client 的端口
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
chmod +644 /opt/polaris/zookeeper1/logs

mkdir mkdir -p /opt/polaris/zookeeper2/logs
mkdir mkdir -p /opt/polaris/zookeeper2/data
mkdir mkdir -p /opt/polaris/zookeeper2/datalog
mkdir mkdir -p /opt/polaris/zookeeper2/conf
chmod +644 /opt/polaris/zookeeper2/logs

mkdir mkdir -p /opt/polaris/zookeeper3/logs
mkdir mkdir -p /opt/polaris/zookeeper3/data
mkdir mkdir -p /opt/polaris/zookeeper3/datalog
mkdir mkdir -p /opt/polaris/zookeeper3/conf
chmod +644 /opt/polaris/zookeeper3/logs
````

#### 创建conf目录下的配置文件zoo.cfg,log4j.properties,ip地址为本机
zookeeper1/conf/zoo.cfg:
````
clientPort=2181
dataDir=/data
dataLogDir=/datalog
tickTime=2000
initLimit=5
syncLimit=2
server.1=10.1.86.211:2888:3888
server.2=10.1.86.211:2788:3788
server.3=10.1.86.211:2688:3688
autopurge.snapRetainCount=3
autopurge.purgeInterval=1
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
server.2=10.1.86.211:2788:3788
server.3=10.1.86.211:2688:3688
autopurge.snapRetainCount=3
autopurge.purgeInterval=1
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
server.2=10.1.86.211:2788:3788
server.3=10.1.86.211:2688:3688
autopurge.snapRetainCount=3
autopurge.purgeInterval=1
````

#### 每个目录下的conf 目录下都要创建log4j.properties

#### 启动脚本
zookeeper1/start.sh:
````
sudo docker run -d --name zookeeper1 --net="host"  -v /opt/polaris/zookeeper1/logs:/zookeeper-3.4.12/logs  -v /opt/polaris/zookeeper/conf:/conf -v /opt/polaris/zookeeper/data:/data -v /opt/polaris/zookeeper/datalog:/datalog -e "ZOO_MY_ID=1" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" -e "ZOO_AUTOPURGE_SNAPRETAINCOUNT=3" -e "ZOO_AUTOPURGE_PURGEINTERVAL=1"  172.16.59.153/develop/zookeeper:3.4.12_private
````
zookeeper2/start.sh:
````
sudo docker run -d --name zookeeper2 --net="host"  -v /opt/polaris/zookeeper2/logs:/zookeeper-3.4.12/logs  -v /opt/polaris/zookeeper2/conf:/conf -v /opt/polaris/zookeeper2/data:/data -v /opt/polaris/zookeeper2/datalog:/datalog -e "ZOO_MY_ID=2" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" -e "ZOO_AUTOPURGE_SNAPRETAINCOUNT=3" -e "ZOO_AUTOPURGE_PURGEINTERVAL=1"  172.16.59.153/develop/zookeeper:3.4.12_private
````
zookeeper3/start.sh:
````
sudo docker run -d --name zookeeper3 --net="host"  -v /opt/polaris/zookeeper3/logs:/zookeeper-3.4.12/logs  -v /opt/polaris/zookeeper3/conf:/conf -v /opt/polaris/zookeeper3/data:/data -v /opt/polaris/zookeeper3/datalog:/datalog -e "ZOO_MY_ID=3" -e "ZOO_LOG4J_PROP=INFO,ROLLINGFILE" -e "ZOO_LOG_DIR=/zookeeper-3.4.12/logs"  -e "ZOO_SERVERS:server.1=10.1.86.211:2888:3888 server.2=10.1.86.70:2888:3888 server.3=10.1.86.212:2888:3888" -e "ZOO_AUTOPURGE_SNAPRETAINCOUNT=3" -e "ZOO_AUTOPURGE_PURGEINTERVAL=1"  172.16.59.153/develop/zookeeper:3.4.12_private
````
#### 