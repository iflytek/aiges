#!/bin/bash
#Eric 2017/1/13
#ywip=`ip addr | grep 'inet '$1 |awk -F '[ /]+' ' NR==1{print $3}'`
#ip addr|grep 'inet 172.16.52' |awk -F '[ /]+' ' NR==1{print $3}'
name="cynosure" #etcd key use this
businesspath="/opt/server/"$name"/"
classpath="com.iflytek.ccr.polaris.cynosure.Application"
startcmd=""
startcmd="/opt/server/"$name"/bin/start.sh"
check_time=10 #seconds
#=======Start program and daemon it======
echo "service was started..."
$startcmd $@
sleep 10
echo $name" is running"
while true
do
pid=`ps ax | grep -i $classpath |grep java| grep $businesspath | grep -v grep | awk '{print $1}'`
if [ -z "$pid" ]; then
        $startcmd $@
        echo "service was started..."
        sleep 10
        pid=`ps ax | grep -i $classpath |grep java| grep $businesspath | grep -v grep | awk '{print $1}'`
        if [ -z "$pid" ]; then
             echo $startcmd" can not run,you can edit here and send message"
             #exit 1
        fi
        echo $name" is running"
fi
sleep $check_time
done
