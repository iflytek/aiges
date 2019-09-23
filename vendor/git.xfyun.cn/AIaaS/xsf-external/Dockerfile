FROM 172.16.59.153/aiaas/ubuntu:14.04.probes
MAINTAINER sqjian@iflytek.com

COPY example/client/*              /xsf/example/client/
COPY example/server/*              /xsf/example/server/
###############################
COPY probes/live/*                 /xsf/probes/live/
###############################
COPY lab/rate/client/*             /xsf/rate/client/
COPY lab/rate/server/*             /xsf/rate/server/
###############################
COPY lab/perf/client/*        /xsf/perf/client/
COPY lab/perf/server/*        /xsf/perf/server/
###############################
COPY lab/grpc/client/*        /xsf/grpc/client/
COPY lab/grpc/server/*        /xsf/grpc/server/
###############################
COPY lab/basic/client/*       /xsf/basic/client/
COPY lab/basic/cmdClient/*      /xsf/basic/cmdClient/
COPY lab/basic/server/*       /xsf/basic/server/
COPY lab/basic/other/hostAdapter/*       /xsf/basic/other/hostAdapter/
###############################
COPY lab/cgroup/*       /xsf/cgroup/

WORKDIR /xsf