export GOPATH=`pwd`
mkdir output
mkdir -p output/include
go build -v -o ./output/AIservice -gcflags "-N -l -c 10" ./src/github.com/xfyun/aiges/main/main.go
cp ./src/github.com/xfyun/aiges/cgo/header/widget/* ./output/include/
cp -r ./src/github.com/xfyun/aiges/cgo/library/* ./output/
