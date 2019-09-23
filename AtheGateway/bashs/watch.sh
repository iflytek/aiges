if [ -z "$KONG_TARGET_WEIGHT" ]
then
    KONG_TARGET_WEIGHT="100"
fi

if [ -z "$KONG_ADMIN_API" ]
then
    KONG_ADMIN_API="http://10.1.87.68"
fi

if [ -z "$KONG_APIKEY" ]
then
    KONG_APIKEY="password"
fi

if [ -z "$KONG_APISECRET" ]
then
    KONG_APISECRET="secret"
fi

if [ -n "$KONG_UPSTREAM" ]
then
    KONG_SERVICE_NAME=$KONG_UPSTREAM
fi

function generateHeader() {

    echo "upstream server address=> $targetAddr"
    echo "request_Url=> $2"
    method=$1
    requestUrl=$2
    datenow=`date -u "+%a, %d %b %G %H:%M:%S MST"` # 获取时间戳
    url=${requestUrl#*//}
    uri=/${url#*/}
    requestline="$method $uri HTTP/1.1"
    host=${KONG_ADMIN_API#*//}
    host=${host%%/*}
    signStr="host: $host\ndate: $datenow\n$requestline"
    #计算hmac signature
    signature=`printf "$signStr" |openssl dgst -sha256 -hmac  "$KONG_SECRET" -binary | base64`
    authHeader="hmac username=\"$KONG_APIKEY\", algorithm=\"hmac-sha256\", headers=\"host date request-line\", signature=\"$signature\""

}

function register() {
    echo "start register------------------------------"
    requestUrl=$KONG_ADMIN_API/upstreams/$KONG_SERVICE_NAME/targets
    method=POST
    generateHeader $method $requestUrl
    curl -X POST -i $requestUrl -H "Authorization:$authHeader" -H "Host:$host" -H "Date:$datenow" -d "target=$targetAddr" -d "weight=$KONG_TARGET_WEIGHT"
}


function unregister() {
    echo "start unregister-----------------------------"
    requestUrl=$KONG_ADMIN_API/upstreams/$KONG_SERVICE_NAME/targets/$targetAddr
    echo "request_url=> $requestUrl"
    method=DELETE
    generateHeader $method $requestUrl
    curl -X DELETE -i $requestUrl -H "Authorization:$authHeader" -H "Host:$host" -H "Date:$datenow"
}



#### 启动服务
check_time=10
#注册服务
#收到退出信号，取消注册服务
startcmd="./webgate-ws-app $*"

trap 'unregister' SIGQUIT
echo "service was started ..."
#启动命令
#启动服务
$startcmd &
echo "service is running"
sleep 10
#初始化 ip，port
source /etc/webgate-env
if [ -z "$APP_HOST" ]
then
    ip=`hostname -i`
    ip=${ip%% *}
else
    ip=$APP_HOST
fi

targetAddr=$ip:$APP_PORT


pid=`ps -ax |grep "$startcmd" |grep -v grep |awk '{print $1}'`
#看看程序是否启动成功 没有则直接退出
if [ -z "$pid" ];then
    echo "start service failed"
    unregister
    exit 1
fi
# 更新环境变量

register
#服务不存在时退出
while true
do
pid=`ps -ax |grep "$startcmd" |grep -v grep |awk '{print $1}'`
if [ -z "$pid" ]; then
    unregister  #程序不在时取消注册服务
    sleep 10
    exit 1
fi
sleep $check_time
done
