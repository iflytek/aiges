#!/bin/bash
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push littlescw00/atheproxy:latest
docker push littlescw00/athegateway:latest
docker push littlescw00/athelb:latest
