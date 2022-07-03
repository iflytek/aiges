import logging

log = logging.getLogger()

SUPPORTED_DISTRO_LIST = ["ubuntu1804"]
SUPPORTED_PYVERSION_LIST = ["3.9.13"]
SUPPORTED_GOLANG_LIST = ["1.17"]
SUPPORTED_CUDA_LIST = ["10.1", "10.2", "11.2", "11.6"]

ECR_REPO = "public.ecr.aws/iflytek-open"
INNER_REPO = "artifacts.iflytek.com/docker-private/atp"
TEMP_GEN_DIR = "./dist/aiges"
Dockerfile = "Dockerfile"