ARG IMAGE_NAME
FROM ${IMAGE_NAME}:10.1-runtime-ubuntu18.04 as base

FROM base as base-amd64

ENV NV_CUDA_LIB_VERSION 10.1.243-1
ENV NV_CUDA_CUDART_DEV_VERSION 10.1.243-1
ENV NV_NVML_DEV_VERSION 10.1.243-1
ENV NV_LIBCUSPARSE_DEV_VERSION 10.1.243-1
ENV NV_LIBNPP_DEV_VERSION 10.1.243-1
ENV NV_LIBCUBLAS_DEV_PACKAGE_NAME libcublas-dev

ENV NV_LIBCUBLAS_DEV_VERSION 10.2.1.243-1
ENV NV_LIBCUBLAS_DEV_PACKAGE ${NV_LIBCUBLAS_DEV_PACKAGE_NAME}=${NV_LIBCUBLAS_DEV_VERSION}

ENV NV_LIBNCCL_DEV_PACKAGE_NAME libnccl-dev
ENV NV_LIBNCCL_DEV_VERSION 2.8.3-1
ENV NCCL_VERSION ${NV_LIBNCCL_DEV_VERSION}
ENV NV_LIBNCCL_DEV_PACKAGE ${NV_LIBNCCL_DEV_PACKAGE_NAME}=${NV_LIBNCCL_DEV_VERSION}+cuda10.1
ENV NV_LIBNCCL_PACKAGE_SHA256SUM 2e2218653517288004b25cafbf511f523c42a3fa7af21e7edf32f145a4deda94
ENV NV_LIBNCCL_PACKAGE_SOURCE https://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu1804/x86_64/libnccl2_2.8.3-1+cuda10.1_amd64.deb
ENV NV_LIBNCCL_PACKAGE_SOURCE_NAME libnccl2_2.8.3-1+cuda10.1_amd64.deb

RUN apt-get update && apt-get install -y --no-install-recommends wget

RUN wget -q ${NV_LIBNCCL_PACKAGE_SOURCE} \
    && echo "$NV_LIBNCCL_PACKAGE_SHA256SUM  ${NV_LIBNCCL_PACKAGE_SOURCE_NAME}" | sha256sum -c --strict - \
    && dpkg -i ${NV_LIBNCCL_PACKAGE_SOURCE_NAME} \
    && rm -f ${NV_LIBNCCL_PACKAGE_SOURCE_NAME} \
    && apt-get purge --autoremove -y wget \
    && rm -rf /var/lib/apt/lists/*

ENV NV_LIBNCCL_DEV_PACKAGE_SHA256SUM fb3f5f11ad8ee6e35f24ab9ed2e601a6684b5524f47e0a362db11041644696b3
ENV NV_LIBNCCL_DEV_PACKAGE_SOURCE https://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu1804/x86_64/libnccl-dev_2.8.3-1+cuda10.1_amd64.deb
ENV NV_LIBNCCL_DEV_PACKAGE_SOURCE_NAME libnccl-dev_2.8.3-1+cuda10.1_amd64.deb
RUN apt-get update && apt-get install -y --no-install-recommends wget

RUN wget -q ${NV_LIBNCCL_DEV_PACKAGE_SOURCE} \
    && echo "$NV_LIBNCCL_DEV_PACKAGE_SHA256SUM  ${NV_LIBNCCL_DEV_PACKAGE_SOURCE_NAME}" | sha256sum -c --strict - \
    && dpkg -i ${NV_LIBNCCL_DEV_PACKAGE_SOURCE_NAME} \
    && rm -f ${NV_LIBNCCL_DEV_PACKAGE_SOURCE_NAME} \
    && apt-get purge --autoremove -y wget \
    && rm -rf /var/lib/apt/lists/*

FROM base as base-ppc64le

ENV NV_CUDA_LIB_VERSION 10.1.243-1
ENV NV_CUDA_CUDART_DEV_VERSION 10.1.243-1
ENV NV_NVML_DEV_VERSION 10.1.243-1
ENV NV_LIBCUSPARSE_DEV_VERSION 10.1.243-1
ENV NV_LIBNPP_DEV_VERSION 10.1.243-1
ENV NV_LIBCUBLAS_DEV_PACKAGE_NAME libcublas-dev
ENV NV_LIBCUBLAS_DEV_VERSION 10.2.1.243-1
ENV NV_LIBCUBLAS_DEV_PACKAGE ${NV_LIBCUBLAS_DEV_PACKAGE_NAME}=${NV_LIBCUBLAS_DEV_VERSION}

ENV NV_LIBNCCL_DEV_PACKAGE_NAME libnccl-dev
ENV NV_LIBNCCL_DEV_VERSION 2.8.3-1
ENV NCCL_VERSION ${NV_LIBNCCL_DEV_VERSION}
ENV NV_LIBNCCL_DEV_PACKAGE ${NV_LIBNCCL_DEV_PACKAGE_NAME}=${NV_LIBNCCL_DEV_VERSION}+cuda10.1
ENV NV_LIBNCCL_PACKAGE_SHA256SUM e5f73701b0af959de36db8fc6549d698e452bb8bc3c64da5b6e9d5c40d8bab01
ENV NV_LIBNCCL_PACKAGE_SOURCE https://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu1804/ppc64el/libnccl2_2.8.3-1+cuda10.1_ppc64el.deb
ENV NV_LIBNCCL_PACKAGE_SOURCE_NAME libnccl2_2.8.3-1+cuda10.1_ppc64el.deb

RUN apt-get update && apt-get install -y --no-install-recommends wget

RUN wget -q ${NV_LIBNCCL_PACKAGE_SOURCE} \
    && echo "$NV_LIBNCCL_PACKAGE_SHA256SUM  ${NV_LIBNCCL_PACKAGE_SOURCE_NAME}" | sha256sum -c --strict - \
    && dpkg -i ${NV_LIBNCCL_PACKAGE_SOURCE_NAME} \
    && rm -f ${NV_LIBNCCL_PACKAGE_SOURCE_NAME} \
    && apt-get purge --autoremove -y wget \
    && rm -rf /var/lib/apt/lists/*

ENV NV_LIBNCCL_DEV_PACKAGE_SHA256SUM 90461ea41c2053a886257f16b3c7f76b69efa47eabb8956510260f8e6468f873
ENV NV_LIBNCCL_DEV_PACKAGE_SOURCE https://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu1804/ppc64el/libnccl-dev_2.8.3-1+cuda10.1_ppc64el.deb
ENV NV_LIBNCCL_DEV_PACKAGE_SOURCE_NAME libnccl-dev_2.8.3-1+cuda10.1_ppc64el.deb
RUN apt-get update && apt-get install -y --no-install-recommends wget

RUN wget -q ${NV_LIBNCCL_DEV_PACKAGE_SOURCE} \
    && echo "$NV_LIBNCCL_DEV_PACKAGE_SHA256SUM  ${NV_LIBNCCL_DEV_PACKAGE_SOURCE_NAME}" | sha256sum -c --strict - \
    && dpkg -i ${NV_LIBNCCL_DEV_PACKAGE_SOURCE_NAME} \
    && rm -f ${NV_LIBNCCL_DEV_PACKAGE_SOURCE_NAME} \
    && apt-get purge --autoremove -y wget \
    && rm -rf /var/lib/apt/lists/*

FROM base-${TARGETARCH}

ARG TARGETARCH
LABEL maintainer "NVIDIA CORPORATION <cudatools@nvidia.com>"

RUN apt-get update && apt-get install -y --no-install-recommends \
    cuda-nvml-dev-10-1=${NV_NVML_DEV_VERSION} \
    cuda-command-line-tools-10-1=${NV_CUDA_LIB_VERSION} \
    cuda-nvprof-10-1=${NV_CUDA_LIB_VERSION} \
    cuda-npp-dev-10-1=${NV_LIBNPP_DEV_VERSION} \
    cuda-libraries-dev-10-1=${NV_CUDA_LIB_VERSION} \
    cuda-minimal-build-10-1=${NV_CUDA_LIB_VERSION} \
    ${NV_LIBCUBLAS_DEV_PACKAGE} \
    ${NV_LIBNCCL_DEV_PACKAGE} \
    && rm -rf /var/lib/apt/lists/*

# apt from auto upgrading the cublas package. See https://gitlab.com/nvidia/container-images/cuda/-/issues/88
RUN apt-mark hold ${NV_LIBCUBLAS_DEV_PACKAGE_NAME} ${NV_LIBNCCL_DEV_PACKAGE_NAME}

ENV LIBRARY_PATH /usr/local/cuda/lib64/stubs