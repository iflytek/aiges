### 部署说明(以听写为例)

#### 配置中心部署

**依赖配置中心**

*1* 配置中心建立项目
```cgo
project=AIaaS
group=dx
service=webgate-ws-iat
version=1.1.0

```
*2* 配置中心上传以下配置

[xsf.toml](../conf/xsf.toml) 、[app.toml](../conf/app.toml) 和 [shcema_svc.json](../conf/schema_svc.json)


*3* 拉取镜像
```
docker pull 172.16.59.153/aiaas/webgate-ws
```

*4* 修改schema.json文件配置，配置要部署的ai能力服务。后面如果要新增ai能力，不需要修改代码，通过新增schema文件，命名为schema_$name.json,并在app.toml配置中 schema.services增加$name。即可实现新能力上线。（例子中已经配置了听写的schema）schema [配置见](config schema.md)

app.toml

````
[schema]
services=["svc"]
````

*5* 启动镜像

```
docker run --name webgate-ws --net=host -d \
-v /data/webgate-ws/log:/log/server \
-v /etc/localtime:/etc/localtime  \
-e "KONG_TARGET_WEIGHT=100" \
-e "KONG_ADMIN_API=http://10.1.87.70:8000" \
-e "KONG_UPSTREAM=webgates" \
-e "KONG_APIKEY=XXXXX" \
-e "KONG_SECRET=XXXXX" \
 172.16.59.153/aiaas/webgate-ws /bin/bash bashs/watch.sh \
 --nativeBoot=false --project=AIaaS --group=dx --service=webgate-ws-svc --version=1.1.0 --url=http://companion.xfyun.iflytek:6868
```

webgate的鉴权和公网代理依赖 kong，如果接口需要暴露到公网，需要部署kong，通过kong代理到webgate并配置鉴权插件，kong的使用可以参考官方文档。部署流程见[kong部署与配置](kong-deploy-doc.pdf)



