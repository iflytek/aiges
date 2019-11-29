#!/bin/bash

#help="docker run -itd --net host -v /log/athena/:/log/ --rm littlescw00/cynosure:latest sh watchdog.sh --server.port=8011 --mysql.addr=10.1.87.68:3306 --mysql.database=ifly_cynosure_3 --spring.datasource.username=test --spring.datasource.password=password"
help="usage: run_cynosure.sh server_port mysql_host mysql_port mysql_user mysql_password"

server_port=$1
mysql_host=$2
mysql_port=$3
mysql_user=$4
mysql_password=$5

cleanup() {
	docker rm -f cynosure
	echo "drop database ifly_cynosure_3;" | mysql -h${mysql_host} -P${mysql_port} -u${mysql_user} -p${mysql_password}
}


if [ $# != "5" ]; then
	echo ${help}
	exit 1
fi


#mysql -h10.1.87.68 -P3306 -utest -ppassword
mysql -h${mysql_host} -P${mysql_port} -u${mysql_user} -p${mysql_password} < sqls/init.sql || cleanup
for i in `ls sqls |grep -v init`; do
	mysql -h${mysql_host} -P${mysql_port} -u${mysql_user} -p${mysql_password} < sqls/${i} || cleanup;
done

docker run -itd --name cynosure --net host -v /log/athena/:/log/ --rm littlescw00/cynosure:latest sh watchdog.sh --server.port=${server_port} --mysql.addr=${mysql_host}:${mysql_port} --mysql.database=ifly_cynosure_3 --spring.datasource.username=${mysql_user} --spring.datasource.password=${mysql_password} || cleanup
