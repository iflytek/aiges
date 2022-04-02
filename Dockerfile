FROM hub.iflytek.com/aiaas/ubuntu:14.04.gcc.golang1.14
RUN apt-get update && apt-get install -y libnuma-dev
COPY ./output /home/