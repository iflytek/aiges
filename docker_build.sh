#!/bin/bash

docker build -f Dockerfile.AtheProxy -t littlescw00/atheproxy:latest  .
docker build -f Dockerfile.AtheGateway -t littlescw00/athegateway:latest  .
docker build -f Dockerfile.AtheLB -t littlescw00/athelb:latest  .
