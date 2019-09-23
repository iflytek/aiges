export GOPATH=`pwd`
#cp -rf ./src/cgo/library/* ./
go test -v ./src/buffer
go test -v ./src/dp
#go test -v ./src/apm
