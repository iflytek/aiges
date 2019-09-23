#! /bin/sh

echo "begin to start cynosure ..."
BASEDIR='/opt/server/cynosure'
#pid=`ps ax | grep -i 'com.iflytek.ccr.polaris.cynosure.Application' |grep java| grep $BASEDIR/bin | grep -v grep | awk '{print $1}'`
#pid=`ps ax | grep -i 'cynosure' |grep java`
pid=`ps ax | grep -i 'cynosure' |grep java |grep -v grep | awk '{print $1}'`
if [ -n "$pid" ] ; then
        echo "cynosure is already running ..."
        exit 0
fi
LogPath='/log/server'

if [ ! -d "$LogPath" ]; then
        mkdir "$LogPath"
fi

JAVA=$JAVA_HOME/bin/java
JAVA_OPT_1="-server -Xms4g -Xmx16g -Xmn4g -XX:PermSize=1g -XX:MaxPermSize=2g -XX:+PrintGCDetails -XX:+PrintGCTimeStamps"
JAVA_OPT_2="-XX:+UseConcMarkSweepGC -XX:+UseCMSCompactAtFullCollection -XX:CMSInitiatingOccupancyFraction=70 -XX:+CMSParallelRemarkEnabled -XX:SoftRefLRUPolicyMSPerMB=0 -XX:+CMSClassUnloadingEnabled -XX:SurvivorRatio=8 -XX:+DisableExplicitGC"
JAVA_OPT_3="-verbose:gc -Xloggc:$LogPath/cynosure_gc.log -XX:+PrintGCDetails -XX:+PrintGCDateStamps"
JAVA_OPT_4="-XX:-OmitStackTraceInFastThrow"
JAVA_OPT_5="-XX:+HeapDumpOnOutOfMemoryError"
JAVA_OPT_6="-XX:HeapDumpPath=$LOG_PATH"
JAVA_OPT_7="-XX:ErrorFile=$LogPath/java_error_%p.log"

JAVA_OPTS="${JAVA_OPT_1} ${JAVA_OPT_2} ${JAVA_OPT_3} ${JAVA_OPT_4} ${JAVA_OPT_5} ${JAVA_OPT_6} ${JAVA_OPT_7}"
for i in $BASEDIR/lib/*.jar
do
    CLASSPATH=$CLASSPATH:"$i"
done


export CLASSPATH=.:"$BASEDIR":$CLASSPATH
#$JAVA $JAVA_OPTS com.iflytek.ccr.polaris.cynosure.Application $@ -d $BASEDIR/bin  >> $LogPath/run.log 2>&1 &
$JAVA $JAVA_OPTS -jar $BASEDIR/cynosure-1.0.0-SNAPSHOT.jar $@ >> $LogPath/run.log 2>&1 &
