#!/bin/bash

function build_AtheLB() {
  rm -rf /tmp/AtheLB
  export GOPATH=/tmp/AtheLB/
  cd /tmp/AtheLB/src/lbv2 && sh build.sh
}

function build_AtheGateway() {
  rm -rf /tmp/AtheGateway
  mkdir -p /tmp/AtheGateway/src
  cp AtheGateway/* /tmp/AtheGateway/src -R
  export GOPATH=/tmp/AtheGateway/
  cd /tmp/AtheGateway/
  go build -o AtheGateway src/main/main.go
}

function build_AtheProxy() {
  rm -rf /tmp/AtheProxy
  mkdir -p /tmp/AtheProxy/src
  export GOPATH=/tmp/AtheProxy/
  cd /tmp/AtheProxy/
  go build -o AtheProxy src/main/main.go
}

function build_AtheLoader() {
  rm -rf /tmp/AtheProxy
  mkdir -p /tmp/AtheProxy/
  export GOPATH=/tmp/AtheProxy/
  cd /tmp/AtheProxy/
  go build -o AtheProxy src/main/main.go
}
