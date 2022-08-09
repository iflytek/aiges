FROM debian:bullseye-slim

LABEL maintainer="ybyang7@iflytek.com"

ENV LANG=C.UTF-8 LC_ALL=C.UTF-8

# bullseye

ENV TZ=Asia/Shanghai \
DEBIAN_FRONTEND=noninteractive

MAINTAINER ybyang7@iflytek.com
COPY docker/repos/bullseye/sources.list /etc/apt/sources.list


# Setup Go
# https://github.com/docker-library/golang/blob/master/1.17/bullseye/Dockerfile
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

ENV GOLANG_VERSION 1.18.5

RUN set -eux; \
arch="$(dpkg --print-architecture)"; arch="${arch##*-}"; \
url=; \
case "$arch" in \
'amd64') \
url='https://dl.google.com/go/go1.18.5.linux-amd64.tar.gz'; \
sha256='9e5de37f9c49942c601b191ac5fba404b868bfc21d446d6960acc12283d6e5f2'; \
;; \
'armel') \
export GOARCH='arm' GOARM='5' GOOS='linux'; \
;; \
'armhf') \
url='https://dl.google.com/go/go1.18.5.linux-armv6l.tar.gz'; \
sha256='d5ac34ac5f060a5274319aa04b7b11e41b123bd7887d64efb5f44ead236957af'; \
;; \
'arm64') \
url='https://dl.google.com/go/go1.18.5.linux-arm64.tar.gz'; \
sha256='006f6622718212363fa1ff004a6ab4d87bbbe772ec5631bab7cac10be346e4f1'; \
;; \
'i386') \
url='https://dl.google.com/go/go1.18.5.linux-386.tar.gz'; \
sha256='0c44f85d146c6f98c34e8ff436a42af22e90e36fe232d3d9d3101f23fd61362b'; \
;; \
'mips64el') \
export GOARCH='mips64le' GOOS='linux'; \
;; \
'ppc64el') \
url='https://dl.google.com/go/go1.18.5.linux-ppc64le.tar.gz'; \
sha256='2e37fb9c7cbaedd4e729492d658aa4cde821fc94117391a8105c13b25ca1c84b'; \
;; \
's390x') \
url='https://dl.google.com/go/go1.18.5.linux-s390x.tar.gz'; \
sha256='e3d536e7873639f85353e892444f83b14cb6670603961f215986ae8e28e8e07a'; \
;; \
*) echo >&2 "error: unsupported architecture '$arch' (likely packaging update needed)"; exit 1 ;; \
esac; \
build=; \
if [ -z "$url" ]; then \
# https://github.com/golang/go/issues/38536#issuecomment-616897960
		build=1; \
		url='https://dl.google.com/go/go1.18.5.src.tar.gz'; \
		sha256='9920d3306a1ac536cdd2c796d6cb3c54bc559c226fc3cc39c32f1e0bd7f50d2a'; \
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


FROM debian:bullseye-slim

LABEL maintainer="ybyang7@iflytek.com"

ENV LANG=C.UTF-8 LC_ALL=C.UTF-8
MAINTAINER ybyang7@iflytek.com
ENV TZ Asia/Shanghai

COPY docker/repos/bullseye/sources.list /etc/apt/sources.list
RUN DEBIAN_FRONTEND=noninteractive apt update &&apt install -y libnuma-dev git vim
WORKDIR /home/aiges

# setup native python

# runtime dependencies
RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		libbluetooth-dev \
		tk-dev \
		uuid-dev \
    	libffi-dev \
    	libssl-dev \
        libbz2-dev \
        liblzma-dev \
    ; \
	rm -rf /var/lib/apt/lists/*

ENV GPG_KEY E3FF2839C048B25C084DEBE9B26995E310250568
ENV PYTHON_VERSION 3.9.13

RUN set -eux; \
	\
	wget -O python.tar.xz "https://registry.npmmirror.com/-/binary/python/${PYTHON_VERSION%%[a-z]*}/Python-$PYTHON_VERSION.tar.xz"; \
	wget -O python.tar.xz.asc "https://registry.npmmirror.com/-/binary/python/${PYTHON_VERSION%%[a-z]*}/Python-$PYTHON_VERSION.tar.xz.asc"; \
	GNUPGHOME="$(mktemp -d)"; export GNUPGHOME; \
	gpg --batch --keyserver hkps://keys.openpgp.org --recv-keys "$GPG_KEY"; \
	gpg --batch --verify python.tar.xz.asc python.tar.xz; \
	command -v gpgconf > /dev/null && gpgconf --kill all || :; \
	rm -rf "$GNUPGHOME" python.tar.xz.asc; \
	mkdir -p /usr/src/python; \
	tar --extract --directory /usr/src/python --strip-components=1 --file python.tar.xz; \
	rm python.tar.xz; \
	\
	cd /usr/src/python; \
	gnuArch="$(dpkg-architecture --query DEB_BUILD_GNU_TYPE)"; \
	./configure \
		--build="$gnuArch" \
		--enable-loadable-sqlite-extensions \
		--enable-optimizations \
		--enable-option-checking=fatal \
		--enable-shared \
		--with-system-expat \
		--without-ensurepip \
	; \
	nproc="$(nproc)"; \
	make -j "$nproc" \
	; \
	make install; \
	\
# enable GDB to load debugging data: https://github.com/docker-library/python/pull/701
	bin="$(readlink -ve /usr/local/bin/python3)"; \
	dir="$(dirname "$bin")"; \
	mkdir -p "/usr/share/gdb/auto-load/$dir"; \
	cp -vL Tools/gdb/libpython.py "/usr/share/gdb/auto-load/$bin-gdb.py"; \
	\
	cd /; \
	rm -rf /usr/src/python; \
	\
	find /usr/local -depth \
		\( \
			\( -type d -a \( -name test -o -name tests -o -name idle_test \) \) \
			-o \( -type f -a \( -name '*.pyc' -o -name '*.pyo' -o -name 'libpython*.a' \) \) \
		\) -exec rm -rf '{}' + \
	; \
	\
	ldconfig; \
	\
	python3 --version

# make some useful symlinks that are expected to exist ("/usr/local/bin/python" and friends)
RUN set -eux; \
	for src in idle3 pydoc3 python3 python3-config; do \
		dst="$(echo "$src" | tr -d 3)"; \
		[ -s "/usr/local/bin/$src" ]; \
		[ ! -e "/usr/local/bin/$dst" ]; \
		ln -svT "$src" "/usr/local/bin/$dst"; \
	done

# if this is called "PIP_VERSION", pip explodes with "ValueError: invalid truth value '<VERSION>'"
ENV PYTHON_PIP_VERSION 22.0.4
# https://github.com/docker-library/python/issues/365
ENV PYTHON_SETUPTOOLS_VERSION 58.1.0
# https://github.com/pypa/get-pip
ENV PYTHON_GET_PIP_URL https://bootstrap.pypa.io/get-pip.py
ENV PYTHON_GET_PIP_SHA256 ba3ab8267d91fd41c58dbce08f76db99f747f716d85ce1865813842bb035524d

# will cause error here
#COPY ./docker/get-pip.py /

RUN set -eux; \
	\
	wget -O get-pip.py "$PYTHON_GET_PIP_URL"; \
	echo "$PYTHON_GET_PIP_SHA256 *get-pip.py" | sha256sum -c -; \
	\
	export PYTHONDONTWRITEBYTECODE=1; \
	\
	python get-pip.py \
		--disable-pip-version-check \
		--no-cache-dir \
		--no-compile \
		"pip==$PYTHON_PIP_VERSION" \
		"setuptools==$PYTHON_SETUPTOOLS_VERSION" \
	; \
	rm -f get-pip.py; \
	\
	pip --version


ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/aiges:/home/aiges/library


# need xtest
# need aiges.toml
# xtest.toml
COPY --from=go-builder /home/aiges/dist/aiges_linux_amd64 .
COPY --from=go-builder /home/aiges/ai_cpython_wrapper  /home/aiges/ai_cpython_wrapper

RUN cd /home/aiges/ai_cpython_wrapper && make && cp -r /home/aiges/ai_cpython_wrapper/wrapper_lib/* /home/aiges/library && rm -rf /home/aiges/ai_cpython_wrapper

CMD ["sh", "-c", "./AIservice -m=0 -c=aiges.toml -s=svcName -u=http://companion.xfyun.iflytek:6868 -p=AIaaS -g=dx"]

