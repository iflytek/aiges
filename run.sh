#!/bin/bash

if [[ -z "$1" ]]; then
	echo "usage ./run.sh componion_url"
	exit 1
fi


if [ "$1" == "clean" ]; then
	docker rm -f athelb
	docker rm -f athegateway
	docker rm -f atheproxy
	docker rm -f atheloader
	exit 0
fi

url=$1
docker run -itd --net host --name athelb  littlescw00/athelb -m 1 -p athena -g athena -s athelb -u ${url} -c xsfs.toml || exit 1
docker run -itd --net host  --name athegateway littlescw00/athegateway  -project athena -group athena -service athegateway -version 1.0.0 -url ${url} || exit 1
docker run -itd --net host  --name atheproxy littlescw00/atheproxy  -m 1 -p athena -g athena -s atheproxy -u ${url} || exit 1
docker run  --net host -itd  -e "GODEBUG=cgocheck=0" -e "LD_LIBRARY_PATH=/AtheLoader" --name atheloader littlescw00/atheloader:latest -p athena -g athena -s atheloader -u ${url} -c  aiges.toml -m 1 || exit 1

