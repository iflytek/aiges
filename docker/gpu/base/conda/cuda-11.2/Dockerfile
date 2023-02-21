FROM nvidia/cuda:11.2.0-devel-ubuntu18.04
ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

LABEL maintainer="ybyang7, iflytek"

ENV LANG=C.UTF-8 LC_ALL=C.UTF-8
#ADD sources.list /etc/apt/sources.list
ENV GOLANG_VERSION 1.17.11
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
    	wget \
        pkg-config \
    ; \
	rm -rf /var/lib/apt/lists/*

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
	gpg --batch --keyserver keyserver.ubuntu.com --recv-keys 'EB4C 1BFD 4F04 2F6D DDCC  EC91 7721 F63B D38B 4796'; \
# let's also fetch the specific subkey of that key explicitly that we expect "go.tgz.asc" to be signed by, just to make sure we definitely have it
	gpg --batch --keyserver keyserver.ubuntu.com --recv-keys '2F52 8D36 D67B 69ED F998  D857 78BD 6547 3CB3 BD13'; \
	gpg --batch --verify go.tgz.asc go.tgz; \
	gpgconf --kill all; \
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
# hadolint ignore=DL3008
RUN apt-get update -q && \
    apt-get install -q -y --no-install-recommends \
        bzip2 \
        ca-certificates \
        git \
        libglib2.0-0 \
        libsm6 \
        libxext6 \
        libxrender1 \
        mercurial \
        openssh-client \
        procps \
        subversion \
        wget \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ENV PATH /opt/conda/bin:$PATH

CMD [ "/bin/bash" ]

# Leave these args here to better use the Docker build cache
ARG CONDA_VERSION=py39_4.12.0

ARG PYTHONVERSION=3.6
ARG TS_CONDA_ADDR=https://mirrors.tuna.tsinghua.edu.cn/anaconda/miniconda
ARG CONDA_ADDR=https://repo.anaconda.com/miniconda
RUN set -x && \
    UNAME_M="$(uname -m)" && \
    if [ "${UNAME_M}" = "x86_64" ]; then \
        MINICONDA_URL="${TS_CONDA_ADDR}/Miniconda3-${CONDA_VERSION}-Linux-x86_64.sh"; \
        SHA256SUM="78f39f9bae971ec1ae7969f0516017f2413f17796670f7040725dd83fcff5689"; \
    elif [ "${UNAME_M}" = "s390x" ]; then \
        MINICONDA_URL="${TS_CONDA_ADDR}/Miniconda3-${CONDA_VERSION}-Linux-s390x.sh"; \
        SHA256SUM="ff6fdad3068ab5b15939c6f422ac329fa005d56ee0876c985e22e622d930e424"; \
    elif [ "${UNAME_M}" = "aarch64" ]; then \
        MINICONDA_URL="${TS_CONDA_ADDR}/Miniconda3-${CONDA_VERSION}-Linux-aarch64.sh"; \
        SHA256SUM="5f4f865812101fdc747cea5b820806f678bb50fe0a61f19dc8aa369c52c4e513"; \
    elif [ "${UNAME_M}" = "ppc64le" ]; then \
        MINICONDA_URL="${TS_CONDA_ADDR}/Miniconda3-${CONDA_VERSION}-Linux-ppc64le.sh"; \
        SHA256SUM="1fe3305d0ccc9e55b336b051ae12d82f33af408af4b560625674fa7ad915102b"; \
    fi && \
    wget "${MINICONDA_URL}" -O miniconda.sh -q && \
    echo "${SHA256SUM} miniconda.sh" > shasum && \
    if [ "${CONDA_VERSION}" != "latest" ]; then sha256sum --check --status shasum; fi && \
    mkdir -p /opt && \
    sh miniconda.sh -b -p /opt/conda && \
    rm miniconda.sh shasum && \
    ln -s /opt/conda/etc/profile.d/conda.sh /etc/profile.d/conda.sh && \
    echo ". /opt/conda/etc/profile.d/conda.sh" >> ~/.bashrc && \
    echo "conda activate base" >> ~/.bashrc && \
    find /opt/conda/ -follow -type f -name '*.a' -delete && \
    find /opt/conda/ -follow -type f -name '*.js.map' -delete && \
    /opt/conda/bin/conda clean -afy


WORKDIR /home/aiges
## Create the environment:
#RUN echo "name: base \n\
#channels: \n \
#- defaults \n\
#dependencies: \n \
#- python=${PYTHONVERSION} \n\
#prefix: /opt/aiges " > /home/aiges/environment.yml &&  conda env create -f environment.yml

# Make RUN commands use the new environment:
SHELL ["conda", "run", "-n", "base", "/bin/bash", "-c"]

RUN echo "conda activate base " >> ~/.bashrc
RUN conda install -c "conda-forge/label/cf202003" cxx-compiler && conda install sysroot_linux-64=2.17

