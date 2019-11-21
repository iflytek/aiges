#!/bin/bash

CURDIR=`pwd`
OUTPUT="`pwd`/build"

function prebuild() {
  mkdir -p ${OUTPUT}/output
}

function getback() {
  cd ${CURDIR}
}

function build_AtheLB() {
  echo "build AtheLB"
  prebuild
  rm -rf /tmp/AtheLB
  export GOPATH=/tmp/AtheLB/
  mkdir -p  /tmp/AtheLB/ 
  cp AtheLB/* /tmp/AtheLB/ -R
  cp vendor /tmp/AtheLB/src/lbv2/ -R
  cd /tmp/AtheLB/src/lbv2 && sh build.sh || exit 1
  install AtheLB ${OUTPUT}
  # cp /tmp/AtheLB/src/lbv2/lbv2 output/
  getback
}

function build_AtheGateway() {
  echo "build AtheGateway"
  prebuild
  rm -rf /tmp/AtheGateway
  mkdir -p /tmp/AtheGateway/src
  cp AtheGateway/* /tmp/AtheGateway/src -R
  cp vendor /tmp/AtheGateway/src/ -R
  export GOPATH=/tmp/AtheGateway/
  cd /tmp/AtheGateway/ && go build -v -o ${OUTPUT}/AtheGateway src/main/main.go || exit 1
  getback
}

function build_AtheProxy() {
  echo "build AtheProxy"
  rm -rf /tmp/AtheProxy
  mkdir -p /tmp/AtheProxy/src
  cp AtheProxy/* /tmp/AtheProxy/src -R
  cp vendor /tmp/AtheProxy/src/ -R
  export GOPATH=/tmp/AtheProxy/
  cd /tmp/AtheProxy/ && go build -o ${OUTPUT}/AtheProxy src/main/main.go || exit 1
  getback
}

function build_AtheLoader() {
  echo "build AtheLoader"
  cp AtheLoader/ /tmp/ -R
  cp vendor /tmp/AtheLoader/src/ -R
  export GOPATH=/tmp/AtheLoader/
  cd /tmp/AtheLoader/
  sh ./build.sh || exit 1
  install output/* ${OUTPUT}/*
  getback
}

function cleanup() {
	rm -r /tmp/AtheLoader
	rm -r /tmp/AtheProxy
	rm -r /tmp/AtheGateway
	rm -r /tmp/AtheLB
}

TARGET=$1
if [[ -z "${TARGET}" ]]; then
	TARGET="all"
fi

if [[ "${TARGET}" == "all" ]]; then
	build_AtheLB
	build_AtheGateway
	build_AtheProxy
	build_AtheLoader
elif [[ "${TARGET}" == "AtheLB" ]]; then
	build_AtheLB
elif [[ "${TARGET}" == "AtheGateway" ]]; then
	build_AtheGateway
elif [[ "${TARGET}" == "AtheProxy" ]]; then
	build_AtheProxy
elif [[ "${TARGET}" == "AtheLoader" ]]; then
	build_AtheLoader
else
  	echo "Invalid build target"
  	exit 1
fi

cleanup
