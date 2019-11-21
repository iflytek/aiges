export GOPATH=`pwd`
go build -v -o AIservice -gcflags "-N -l -c 10" ./src/main/main.go
go build -v -o ./src/xtest/xtest -gcflags "-N -l -c 10" ./src/xtest/xtest.go
gcc -fPIC -shared -Wno-attributes -g -O0 -rdynamic -o ./output/libwrapper.so ./test/wrapper/wrapper.cpp
mv ./AIservice ./output/
cp -r ./src/cgo/library/* ./output/
cp -r ./test/script/* ./output/
