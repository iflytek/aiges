FROM continuumio/miniconda3:4.12.0 as go-builder
# bullseye

ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

MAINTAINER ybyang7@iflytek.com
COPY docker/repos/bullseye/sources.list /etc/apt/sources.list


# Setup Go
# https://github.com/docker-library/golang/blob/master/1.17/bullseye/Dockerfile
# install cgo-related dependencies
# install cgo-related dependencies
RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
	; \
	rm -rf /var/lib/apt/lists/*

ENV PATH /usr/local/go/bin:$PATH

ENV GOLANG_VERSION 1.17.11

ENV PATH /usr/local/go/bin:$PATH
RUN set -eux; \
	arch="$(dpkg --print-architecture)"; arch="${arch##*-}"; \
	url=; \
	case "$arch" in \
		'amd64') \
			url='https://dl.google.com/go/go1.17.11.linux-amd64.tar.gz'; \
			sha256='d69a4fe2694f795d8e525c72b497ededc209cb7185f4c3b62d7a98dd6227b3fe'; \
			;; \
		'armel') \
			export GOARCH='arm' GOARM='5' GOOS='linux'; \
			;; \
		'armhf') \
			url='https://dl.google.com/go/go1.17.11.linux-armv6l.tar.gz'; \
			sha256='56fb4a344bfe01427398083b3e82018fcde1f168a94e9df9820010b29e558f27'; \
			;; \
		'arm64') \
			url='https://dl.google.com/go/go1.17.11.linux-arm64.tar.gz'; \
			sha256='adefa7412c6798f9cad02d1e8336fc2242f5bade30c5b32781759181e01961b7'; \
			;; \
		'i386') \
			url='https://dl.google.com/go/go1.17.11.linux-386.tar.gz'; \
			sha256='8e43b4c63dda0612df0e1ed3cf8bed566fe8ae88082a46343c40be76a61371e4'; \
			;; \
		'mips64el') \
			export GOARCH='mips64le' GOOS='linux'; \
			;; \
		'ppc64el') \
			url='https://dl.google.com/go/go1.17.11.linux-ppc64le.tar.gz'; \
			sha256='013624fc1964d576f314117f2f3b6a6a8f3d32f0ec8f55fa5962664537a34180'; \
			;; \
		's390x') \
			url='https://dl.google.com/go/go1.17.11.linux-s390x.tar.gz'; \
			sha256='0202a0b2ffaa9bdf1f2480789c2009ba6fde597e1f1c6d4cfb9559bb1964a39c'; \
			;; \
		*) echo >&2 "error: unsupported architecture '$arch' (likely packaging update needed)"; exit 1 ;; \
	esac; \
	build=; \
	if [ -z "$url" ]; then \
# https://github.com/golang/go/issues/38536#issuecomment-616897960
		build=1; \
		url='https://dl.google.com/go/go1.17.11.src.tar.gz'; \
		sha256='ac2649a65944c6a5abe55054000eee3d77196880da36a3555f62e06540e8eb54'; \
		echo >&2; \
		echo >&2 "warning: current architecture ($arch) does not have a compatible Go binary release; will be building from source"; \
		echo >&2; \
	fi; \
	\
	wget -O go.tgz.asc "$url.asc"; \
	wget -O go.tgz "$url" --progress=dot:giga; \
	echo "$sha256 *go.tgz" | sha256sum -c -; \
	\
# https://github.com/golang/go/issues/14739#issuecomment-324767697
	GNUPGHOME="$(mktemp -d)"; export GNUPGHOME; \
# https://www.google.com/linuxrepositories/
#	gpg --batch --keyserver keyserver.ubuntu.com --recv-keys 'EB4C 1BFD 4F04 2F6D DDCC  EC91 7721 F63B D38B 4796'; \
# let's also fetch the specific subkey of that key explicitly that we expect "go.tgz.asc" to be signed by, just to make sure we definitely have it
#	gpg --batch --keyserver keyserver.ubuntu.com --recv-keys '2F52 8D36 D67B 69ED F998  D857 78BD 6547 3CB3 BD13'; \
#	gpg --batch --verify go.tgz.asc go.tgz; \
#	gpgconf --kill all; \
	rm -rf "$GNUPGHOME" go.tgz.asc; \
	\
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	\
	if [ -n "$build" ]; then \
		savedAptMark="$(apt-mark showmanual)"; \
		apt-get update; \
		apt-get install -y --no-install-recommends golang-go; \
		\
		export GOCACHE='/tmp/gocache'; \
		\
		( \
			cd /usr/local/go/src; \
# set GOROOT_BOOTSTRAP + GOHOST* such that we can build Go successfully
			export GOROOT_BOOTSTRAP="$(go env GOROOT)" GOHOSTOS="$GOOS" GOHOSTARCH="$GOARCH"; \
			./make.bash; \
		); \
		\
		apt-mark auto '.*' > /dev/null; \
		apt-mark manual $savedAptMark > /dev/null; \
		apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false; \
		rm -rf /var/lib/apt/lists/*; \
		\
# remove a few intermediate / bootstrapping files the official binary release tarballs do not contain
		rm -rf \
			/usr/local/go/pkg/*/cmd \
			/usr/local/go/pkg/bootstrap \
			/usr/local/go/pkg/obj \
			/usr/local/go/pkg/tool/*/api \
			/usr/local/go/pkg/tool/*/go_bootstrap \
			/usr/local/go/src/cmd/dist/dist \
			"$GOCACHE" \
		; \
	fi; \
	\
	go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

#RUN mkdir -p output/include && go mod vendor && go build -v -o ./output/AIservice -gcflags "-N -l -c 10" main/main.go && \
#    cp ./cgo/header/widget/* ./output/include/ && \
#    cp  -r ./cgo/library/* ./output/

ENV GOPROXY=https://goproxy.cn,direct
RUN apt-get update && apt-get install -y libnuma-dev build-essential git sudo
COPY . /home/aiges
WORKDIR /home/aiges
RUN make build-pack

# 构建ai_cpython_wrapper需要依赖python 并匹配

# conda python setup python version
# Create the environment:
# Make RUN commands use the new environment:


FROM continuumio/miniconda3:4.12.0   as final
MAINTAINER ybyang7@iflytek.com
ENV TZ Asia/Shanghai

COPY docker/repos/bullseye/sources.list /etc/apt/sources.list
RUN apt update && apt install -y build-essential
RUN DEBIAN_FRONTEND=noninteractive apt update &&apt install -y libnuma-dev git vim
WORKDIR /home/aiges

# Make RUN commands use the new environment:
RUN conda install -c conda-forge ncurses

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/aiges:/home/aiges/library:
#/opt/conda/envs/aiges/lib
#RUN rm -f /opt/conda/envs/aiges/lib/libtinfo.so.6
#

# need xtest
# need aiges.toml
# xtest.toml
COPY --from=go-builder /home/aiges/dist/aiges_linux_amd64 .
COPY --from=go-builder /home/aiges/ai_cpython_wrapper  /home/aiges/ai_cpython_wrapper

RUN cd /home/aiges/ai_cpython_wrapper && make && cp -r /home/aiges/ai_cpython_wrapper/wrapper_lib/* /home/aiges/library && rm -rf /home/aiges/ai_cpython_wrapper


CMD ["sh", "-c", "./AIservice -m=0 -c=aiges.toml -s=svcName -u=http://companion.xfyun.iflytek:6868 -p=AIaaS -g=dx"]