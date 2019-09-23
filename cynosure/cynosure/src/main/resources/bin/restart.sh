#! /bin/sh

dir=`pwd`
pid=`ps ax | grep -i 'com.iflytek.ccr.polaris.cynosure.Application' |grep java | grep $dir | grep -v grep | awk '{print $1}'`
params=`ps ax | grep -i 'com.iflytek.ccr.polaris.cynosure.Application' |grep java | grep $dir | grep -v grep`
params=${params##*com.iflytek.sis.uup.main.Program}

if [ -z "$pid" ] ; then
     echo "No cynosure running."
fi

while [ -n "$pid" ];
    do
        echo "The cynosure(${pid}) is running..."
        kill ${pid}
	    echo "Send shutdown request to cynosure(${pid})"
        sleep 1
        pid=`ps ax | grep -i 'com.iflytek.ccr.polaris.cynosure.Application' |grep java | grep $dir | grep -v grep | awk '{print $1}'`
    done

if [ $# -ne 0 ];then
    ./start.sh $@
else
    optConf=
    optPort=
    optHost=
    GETOPT_ARGS=`getopt -o c:h:p:d: -al configDir:,host:,port:,programdir: -- ${params}`
    eval set -- "$GETOPT_ARGS"
    while [ -n "$1" ]
    do
            case "$1" in
                    -c|--configDir) optConf=$2;shift 2;;
                    -h|--host) optHost=$2;shift 2;;
                    -p|--port) optPort=$2;shift 2;;
                    -d|--programdir) shift 2;;
                    --) break;;
                    *) echo $1,$2;shift 2;;
            esac
    done
    startParams=""
    if [ -n "$optConf" ];then
        startParams="$startParams""   -c  ""$optConf"
    fi
    if [ -n "$optPort" ];then
        startParams="$startParams""   -p  ""$optPort"
    fi
    if [ -n "$optHost" ];then
        startParams="$startParams""   -h  ""$optHost"
    fi
    echo "$startParams"
    ./start.sh "${startParams}"
fi





