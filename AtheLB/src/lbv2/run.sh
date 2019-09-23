#!/usr/bin/env bash
#./lbv2 -m 0 -c lbv2.toml  -s lbv2
./lbv2 -v
./lbv2 -m 0 -c lbv2.toml  -s lbv2
#./lbv2 -m 1 -c lbv2.toml -p 3s -s lbv2 -u http://10.1.87.70:6868 -g 3s

#docker run -itd --rm --net=host --name hermes 172.16.59.153/aiaas/hermes:1.0.5 ./lbv2 -m 1 -c lbv2.toml -p 3s -s lbv2 -u http://10.1.87.70:6868 -g 3s
#docker run -itd --rm --net=host --name hermes --ulimit nofile=65535 --ulimit nproc=2048  172.16.59.153/aiaas/hermes:1.0.7 ./lbv2 -m 1 -c lbv2.toml -p 3s -s lbv2 -u http://10.1.87.70:6868 -g 3s
#docker run -itd --rm --net=host --name hermes --ulimit nofile=65535 --ulimit nproc=2048 -v /home/sqjian/hermes_new/lbv2.toml:/lbv2/lbv2.toml  -v /home/sqjian/hermes_new/run.sh:/lbv2/run.sh 172.16.59.153/aiaas/hermes:1.0.7 ./run.sh
#docker run -itd --net=host --name hermes --ulimit nofile=65535 --ulimit nproc=2048 -v /home/sqjian/hermes_new/lbv2.toml:/lbv2/lbv2.toml  -v /home/sqjian/hermes_new/run.sh:/lbv2/run.sh 172.16.59.153/aiaas/hermes:1.0.7 ./run.sh
#docker run -itd --rm --net=host --name hermes_x.x.x 172.16.59.153/aiaas/hermes:x.x.x ./lbv2 -m 1 -c lbv2.toml -p 3s -s lbv2 -u http://10.1.87.70:6868 -g 3s
#docker run -itd --rm --net=host --name lbv2.x.x.x 172.16.59.153/aiaas/hermes:2.0.15 ./lbv2 -m 1 -c lbv2.toml -p 3s -s lbv2.x.x.x -u http://10.1.87.70:6868 -g 3s
#docker run -itd --net=host --name lbv2.x.x.x 172.16.59.153/aiaas/hermes:2.0.15 ./lbv2 -m 1 -c lbv2.toml -p 3s -s lbv2.x.x.x -u http://10.1.87.70:6868 -g 3s
#docker run -itd --net=host --name lbv2-sqjian 172.16.59.153/aiaas/hermes:2.0.15 ./lbv2 -m 1 -c lbv2.toml -p sqjian -s lbv2 -u http://10.1.87.70:6868 -g sqjian
#docker run -itd --net=host --name lbv2 172.16.59.153/aiaas/hermes:2.2.3 ./lbv2 -m 1 -c lbv2.toml -p guiderAllService -s lbv2 -u http://10.1.87.70:6868 -g gas
#docker run -itd --net=host --name lbv2- 172.16.59.153/aiaas/hermes:2.3.12 ./lbv2 -m 1 -c lbv2.toml -p guiderAllService -s lbv2- -u http://10.1.87.70:6868 -g gas

#外部挂载配置文件，本地启动
#docker run -itd --rm --net=host --name=lbv2-hunter --ulimit nofile=65535 --ulimit nproc=65535 -v /home/sqjian/hermes/lbv2.toml:/lbv2/lbv2.toml 172.16.59.153/aiaas/hermes:2.0.6 ./lbv2 -m 0 -c lbv2.toml  -s lbv2


#docker run -it --rm --net=host --name=lbv1 172.16.59.153/aiaas/hermes:v0.0.0.2cbd66230d79ea069c5c8969918d2539d83468be ./lbv2 -m 1 -c lbv1.toml -p AIaaS -s lbv2 -u http://companion.xfyun.iflytek:6868 -g dx
#docker run -it --rm --net=host --name=lbv2 172.16.59.153/aiaas/hermes:v0.0.0.2cbd66230d79ea069c5c8969918d2539d83468be ./lbv2 -m 1 -c lbv2.toml -p AIaaS -s lbv2 -u http://companion.xfyun.iflytek:6868 -g dx
#docker run -it --rm --net=host --name=lbv3 172.16.59.153/aiaas/hermes:v0.0.0.2cbd66230d79ea069c5c8969918d2539d83468be ./lbv2 -m 1 -c lbv3.toml -p AIaaS -s lbv2 -u http://companion.xfyun.iflytek:6868 -g dx

#docker run -it --rm --net=host --name=xxx 172.16.59.153/aiaas/hermes:0.0.0 ./lbv2 -m 1 -c lbv2.toml -p xsf -s xsf-lbv2 -u http://companion.xfyun.iflytek:6868 -g xsf
