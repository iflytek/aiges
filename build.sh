export GOPATH=`pwd`
mkdir output
go build -v -o ./output/AIservice -gcflags "-N -l -c 10" ./src/main/main.go
cp ./src/cgo/header/widget/* ./output/include/
cp -r ./src/cgo/library/* ./output/
#cp -r ./src/xtest/script/* ./output/
#cp -r ./test/script/* ./output/
