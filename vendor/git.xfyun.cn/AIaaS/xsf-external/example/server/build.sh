#!/usr/bin/env bash

export GOROOT=/usr/local/go
export PATH=${GOROOT}/bin:$PATH

export BASEDIR=`pwd`

export GOPATH=${BASEDIR}/../../../../../../:${BASEDIR}
time go build -v