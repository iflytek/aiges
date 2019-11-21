export LD_LIBRARY_PATH=./:$LD_LIBRARY_PATH
./AIservice -m=0 -c=aiges.toml -s=svcname -u=http://companion.xfyun.iflytek:6868 -p=AIaaS -g=dx
#./AIservice -m=1 -c=aiges.toml -s=AIservice -u=http://10.1.87.70:6868 -p=guiderAllService -g=gas -pprof=true -prfAddr=10.1.87.61:5555
