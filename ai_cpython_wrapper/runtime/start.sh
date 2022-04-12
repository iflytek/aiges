export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./
#nohup ./AIservice -m=0 -c=aiges.toml -s=svcName -u=http://companion.xfyun.iflytek:6868 -p=AIaaS -g=dx &
nohup ./AIservice -m=1 -s=mocksvc -c=aiges.toml -p=guiderAllService -g=gas -u=http://10.1.87.70:6868 &
