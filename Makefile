.PHONY: build-linux

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) unpack
BINARY_NAME=main
BINARY_LINUX=$(BINARY_NAME)-linux
GORELEASER_BIN = $(shell pwd)/bin/goreleaser


SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

install-goreleaser: ## check license if not exist install go-lint tools
	#goimports -l -w cmd
	#goimports -l -w pkg
	$(call go-get-tool,$(GORELEASER_BIN),github.com/goreleaser/goreleaser@v1.6.3)

build:
	mkdir bin
	$(GOBUILD) -v -o ./bin/AIservice -gcflags "-N -l -c 10" ./main/main.go
	cp -r ./cgo/library/* ./bin/
	mkdir -p bin/include
	cp -ra ./cgo/header/widget/* ./bin/include

clean:
	rm -rf bin dist

pack:
	tar -acvf aiservice.tar.gz ./bin
	mkdir -p dist
	mv aiservice.tar.gz dist

dist: build pack


build-pack: SHELL:=/bin/bash
build-pack: install-goreleaser  ## build binaries by default
	@echo "build aiges bin"
	$(GORELEASER_BIN) build --snapshot --rm-dist  --timeout=1h

build-release: SHELL:=/bin/bash
build-release: install-goreleaser ## build binaries by default
	@echo "build sealos bin"
	$(GORELEASER_BIN) release --timeout=1h  --release-notes=hack/release/Note.md --debug  --rm-dist

