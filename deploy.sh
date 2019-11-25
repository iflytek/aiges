#!/bin/bash
echo "docker login"
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
echo $?
docker push littlescw00/atheproxy:latest
docker push littlescw00/athegateway:latest
docker push littlescw00/athelb:latest
docker push littlescw00/atheloader:latest
