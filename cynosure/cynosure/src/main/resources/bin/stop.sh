#!/bin/sh

BASEDIR='/opt/server/cynosure'
pid=`ps ax | grep -i 'com.iflytek.ccr.polaris.cynosure.Application' |grep java | grep $BASEDIR/bin | grep -v grep | awk '{print $1}'`
if [ -z "$pid" ] ; then
        echo "No cynosure running."
        exit 0
fi

echo "The cynosure(${pid}) is running..."

kill ${pid}

echo "Send shutdown request to cynosure(${pid}) OK"
