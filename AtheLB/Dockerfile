FROM 172.16.59.153/aiaas/ubuntu:14.04
MAINTAINER sqjian@iflytek.com

#cli
COPY src/lbv2/utils/cli/cli /utils/cli/
COPY src/lbv2/utils/cli/*.sh /utils/cli/
COPY src/lbv2/utils/cli/*.toml /utils/cli/
#srv
COPY src/lbv2/utils/srv/srv /utils/srv/
COPY src/lbv2/utils/srv/*.sh /utils/srv/
COPY src/lbv2/utils/srv/*.toml /utils/srv/

RUN chmod -R +x /utils/*


COPY src/lbv2/lbv2 /lbv2/
COPY src/lbv2/run.sh /lbv2/
COPY src/lbv2/lbv2.toml /lbv2/
COPY src/lbv2/live /lbv2/
COPY docker-compose.yml /lbv2/
RUN chmod +x /lbv2/*



ENV PATH=/lbv2:$PATH
ENV LIVE_PORT=1995

WORKDIR /lbv2
RUN ls -lh
RUN ./lbv2 -v
CMD ./run.sh