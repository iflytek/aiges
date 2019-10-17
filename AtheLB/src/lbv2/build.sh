#!/usr/bin/env bash

#export GOROOT=/usr/local/go
export PATH=${GOROOT}/bin:$PATH
export GOPATH=`pwd`/../../

go clean
go build -v -o AtheLB #编译
./AtheLB -v

#./lbv2 -m 0 -c lbv2.toml  -s lbv2
