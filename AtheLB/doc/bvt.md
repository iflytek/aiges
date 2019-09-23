## 上报

```
单svc、多svc上报 测试通过
```


## 请求

```
单svc、多svc读取 测试通过
```


## 数据库分段

```
segId对应关系 测试通过
```


## rmq消费

```
rmq消费是否成功，以及是否通知了ats,暂未验证!!!
```


## 监控
```
svc、subsvc两个个维度监控查询 测试通过
addr 维度后续看需求增加
```

## 
```
docker run -itd --net=host --name dev-lbv2 172.16.59.153/aiaas/hermes:2.2.5 ./lbv2 -m 1 -c lbv2.toml -p AIaaS -s lbv2 -u http://10.1.87.70:6868 -g dev
docker run -itd --net=host --name dev-lb-niche 172.16.59.153/aiaas/hermes:2.2.5 ./lbv2 -m 1 -c lbv2.toml -p AIaaS -s lb-niche -u http://10.1.87.70:6868 -g dev
```