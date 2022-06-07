.PHONY: build-linux

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) unpack
BINARY_NAME=main
BINARY_LINUX=$(BINARY_NAME)-linux

bulid:
	mkdir bin
	$(GOBUILD) -v -o ./bin/AIservice -gcflags "-N -l -c 10" ./main/main.go
	cp -r ./cgo/library/* ./bin/

clean:
	rm -rf bin

unpack:
	tar -acvf aiservice.tar.gz ./bin
