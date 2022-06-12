ARG IMAGE_NAME
FROM ${IMAGE_NAME}:10.1-runtime-ubuntu18.04 as base

FROM base as base-amd64

ENV NV_CUDNN_PACKAGE_VERSION 8.0.5.39-1
ENV NV_CUDNN_VERSION 8.0.5.39

ENV NV_CUDNN_PACKAGE_NAME libcudnn8
ENV NV_CUDNN_PACKAGE ${NV_CUDNN_PACKAGE_NAME}=${NV_CUDNN_PACKAGE_VERSION}+cuda10.1

FROM base as base-ppc64le

ENV NV_CUDNN_PACKAGE_VERSION 8.0.4.30-1
ENV NV_CUDNN_VERSION 8.0.4.30

ENV NV_CUDNN_PACKAGE_NAME libcudnn8
ENV NV_CUDNN_PACKAGE ${NV_CUDNN_PACKAGE_NAME}=${NV_CUDNN_PACKAGE_VERSION}+cuda10.1
ENV NV_CUDNN_DL_HASHCMD sha256sum
ENV NV_CUDNN_DL_SUM da448059bdbd4585c8855f93438654503fa75bf75dc8c6de39eceabd7c9dc76a
ENV NV_CUDNN_DL_BASENAME libcudnn8_8.0.4.30-1+cuda10.1_ppc64el.deb
ENV NV_CUDNN_DL_URL https://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu1804/ppc64el/libcudnn8_8.0.4.30-1+cuda10.1_ppc64el.deb

RUN apt-get update && apt-get install -y --no-install-recommends wget

RUN wget -q ${NV_CUDNN_DL_URL} \
    && echo "${NV_CUDNN_DL_SUM}  ${NV_CUDNN_DL_BASENAME}" | ${NV_CUDNN_DL_HASHCMD} -c - \
    && dpkg -i ${NV_CUDNN_DL_BASENAME} \
    && rm -f ${NV_CUDNN_DL_BASENAME} \
    && apt-get purge --autoremove -y wget


FROM base-${TARGETARCH}

ARG TARGETARCH

LABEL maintainer "NVIDIA CORPORATION <cudatools@nvidia.com>"

ENV CUDNN_VERSION ${NV_CUDNN_VERSION}

LABEL com.nvidia.cudnn.version="${CUDNN_VERSION}"

RUN apt-get update && apt-get install -y --no-install-recommends \
    ${NV_CUDNN_PACKAGE} \
    ${NV_CUDNN_PACKAGE_DEV} \
    && apt-mark hold ${NV_CUDNN_PACKAGE_NAME} && \
    rm -rf /var/lib/apt/lists/*