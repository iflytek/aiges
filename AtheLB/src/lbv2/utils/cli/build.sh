#!/usr/bin/env bash

export GOROOT=/usr/local/go
export PATH=${GOROOT}/bin:$PATH
export GOPATH=`pwd`/../../../../

go clean
go build -v -o cli #编译

#./lbv2 -m 0 -c lbv2.toml  -s lbv2