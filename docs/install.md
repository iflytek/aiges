# 部署说明
 * [服务发现和配置中心部署部署服务发现](#服务发现和配置中心部署)
   * [依赖](#依赖)
   * [部署页面](#部署页面)
   * [添加区域](#添加区域)
   * [推送配置](#推送配置)
 * [部署服务](#部署服务)
   * [部署](#部署)
   * [验证](#验证)
 * [清理](#清理)

服务发现和配置中心部署
========

依赖
-----
  - mysql
  - zookeeper
  - docker

部署页面
--------

```
./run_cynosure.sh listen_port mysql_host mysql_addr mysql_username mysql_password
#./tools/run_cynosure.sh 8011 127.0.0.1 3306 test password
./run_companion.sh -h local_ip -p listen_port -z zk_ip -w cynosure_ip:cynosure_port 
#./run_companion.sh -h 127.0.0.1 -p 11223 -z 127.0.0.1:2181 -w 127.0.0.1:8011
```
通过访问http://127.0.0.1:8011 即可访问配置中心页面, 默认密码admin 123456

添加区域
-------

![添加区域](../pics/add_area.png)

在弹出的框中填入
```
区域名称: local
companion: http://${companion_url}:${companion_port}

# 其中companion:是之前启动componion监听ip 和host
```

推送配置
------

![配置推送](../pics/push_config.png)

将每一个服务的配置推送到对应componion

部署服务
=======

部署
----
```
./run.sh http://cynosure_ip:cynosure_port
# ./run.sh http://127.0.0.1:11223
```

验证
-----
利用websocket 发送一下请求

url
```
ws://10.1.87.66:9021/v2/svc
```

data
```
{
  "common": {
    "app_id": "123456",
    "uid": "123243443"
  },
  "business": {
    "ent": "svc"
  },
  "data": [
    {
      "id": "",
      "encoding": "",
      "format": "",
      "data": "aGVsbG8lMjB3b3JsZA==",
      "status": 2
    }
  ]
}
```

拿到结果
```
{
  "code": 0,
  "message": "success",
  "sid": "svc00010001@dx16eb5cd52375742902",
  "data": [
    {
      "data": "cmVzcG9uc2UgcmVzdWx0IGZyb20gQXRoZUxvYWRlciB3cmFwcGVy",
      "data_id": "",
      "data_type": 0,
      "desc_args": {
        "sub": "svc"
      },
      "encoding": "",
      "format": "utf8",
      "frame_id": 0,
      "status": 2
    }
  ]
}
```

清理
====
```
./cleanup.sh
```
