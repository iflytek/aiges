# Web API AI能力接入说明文档

## 简介

   &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;全新的Web API服务webgate提供了基于配置的热插拔AI能力特性,只需要通过新增或编辑配置、推送配置两部操作即可完成AI能力的上线或者下线，免去了以往的定制开发工作，业务人员可以把更多的精力投入到AI能力协议的设计和实现上。


## Web API 2.0协议框架

- ##### 请求参数
```json
{
	"common": {},
	"business": {},
	"data": {}
}
```

- ##### 请求参数说明
| 参数名   | 类型   | 必传 | 描述  |
| ---| ---| ---| ---|
| common   | object | 否   | 公共参数|
| business | object | 否   | 业务参数 |
| data     | object | 是   | 业务数据流参数 |

- ##### 返回参数
```json
{
	"code": 0,
	"message": "",
	"data": {},
	"sid": ""
}
```

- ##### 返回参数说明
| 参数名  | 类型   | 描述                                    |
| -------|-------| --------------------------------------|
| code    | int    | 返回码                                |
| message | string | 描述信息                               |
| data    | object | 参考各AI能力协议中的详细定义内容        |
| sid     | string | 会话id，可用于排查分析问题 |

## 设计AI能力协议

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;为了更快的完成AI能力的接入，在进入后续的步骤前请先根据Web API 2.0协议框架来设计待接入AI能力的业务协议，这里拿听写能力举例：
  
  <strong>请求参数：</strong>
  
  ```json
  {  
    "common":{
        "app_id":"9a0e1c52",
        "uid": "203215623"
    },
	"business":{
	   "param1": "v1",
	   "param2": "v2",
	   "param3": "v3"
	},
	"data":{
        "status":0,
        "data":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."    
	 }
}
  
  ```
  
   <strong>返回参数：</strong>
   
   ```json
{
	"code": 0,
	"message": "",
	"data": [{
	   "data":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf",
	   "status":0
	}],
	"sid":"svc@49c7dc69@nc16763de65fc00153e0"
    
}

   ```


## 如何接入AI能力

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;webgate支持了两种通信协议，分别是websocket和http，用于满足流式和非流式AI能力需求。
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;在AI能力接入流程上，我们针对两套通信协议提供了统一的接入配置标准，通过一种类似DSL的配置语言来完成一个或多个AI能力的定义。
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;定义一个AI能力需要包含以下配置信息：协议版本(version)、服务名(service)、服务路由(route)、请求参数定义(schema)、请求参数映射(request.data.mapping)以及返回参数映射(response.data.mapping)。
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<strong>示例如下：</strong>

```json
[{
    "version": "1.0",
    "service": "svc",
    "route": "/v2/svc",
    "request.data.mapping": {},
    "response.data.mapping": {},
    "schema": {}
}]
```

下面以接入听写能力为例，详细介绍如何接入AI能力。

- ##### 第一步，配置基础信息
   
    - 配置业务协议版本号
    
      &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;通过version字段来配置，用来区分不同的业务协议版本，默认值1.0。拿听写能力举例，随着需求变化，听写能力所需的请求或者返回参数会有变化，我们可以通过这个版本号来升级业务协议，而老版本的协议配置只要不下线，用户就能继续使用，达到了向下兼容效果。
      
    - 配置服务名
    
      &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;通过service字段来配置，这里配置的服务名是用来标志具体的AI能力，例如svc。
    
    - 配置路由名称

      &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;通过route字段来配置，当前webgate属于二代Web API服务，因此route前缀是/v2/，完整的route示例：/v2/{service} ，这里的{service}不强制要求和service字段保持一致。
    
    
- ##### 第二步 配置schema 

  <strong>参数配置说明：</strong>
  
 | 参数名     |  类型   | 描述                     |
| ---------- | -------- | ------------------------------- |
| properties | object | 用来定义参数，可以嵌套|
| type       | string | 参数类型可支持的类型有object、integer、number(整型或浮点型),string、boolean、array |
| required   | array  | 用来指定必传的参数名称 |

 <strong>json schemma校验规则配置说明：</strong>
 
 不同类型的参数，可配置的校验规则有所不同，常见的校验规则配置如下:

- <strong>string类型</strong>

| 参数名    | 类型   | 描述                     |
| --------- | ------ | ------------------------ |
| maxLength | int    | 定义字符串的最大长度,>=0 |
| minLength | int    | 定义字符串的最小长度,>=0 |
| enum | array    | 用来限定string类型参数的取值范围 |
| pattern   | string | 用正则表达式约束字符串   |

- <strong>integer或number类型</strong>

| 参数名  | 类型 | 描述   |
| ------- | ---- | ------ |
| minimum | int  | 最小值 |
| maximum | int  | 最大值 |
  
 ```
{
      "type":"object",
      "properties":{
        "common":{
          "type":"object",
          "properties":{
            "app_id":{
              "type":"string",
              "maxLength":50
            },
            "uid":{
              "type":"string",
              "maxLength":50
            }
          },
          "required":["app_id"]
        },
        "business":{
          "type":"object",
          "properties":{
            "param1":{
              "type": "string",
              "maxLength": 30
            },
            "param2":{
              "type": "string",
              "maxLength": 30
            },
            "param3":{
              "type": "number"
            }
          }
        },
        "data":{
          "type":"object",
          "properties":{
            "status":{
              "type":"integer"
            },
            "data":{
              "type": "string",
              "maxLength": 30000000
            }
           
          }
        }
      }
    }
```
  
  <strong>备注：</strong> 更多校验规则的配置项可参考 http://json-schema.org/latest/json-schema-validation.html#rfc.section.6.7 。
    
    
- ##### 第三步，配置请求参数映射规则
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;webgate的上游服务是AIaaS 2.0架构中的Atmos服务，未了保证通用，使用了比较抽象的的业务协议。为了降低理解和使用Web API的门槛，我们对后端请求协议中的data参数做了一层映射，暴露给用户简洁明了的请求参数。

 <strong>服务配置示例：</strong>

 ```json
{
	"request.data.mapping": {
		"data_type":[1],
		"rule": [
			{"dst":"$[0].data","src":"$.data"},
			{"dst":"$[0].status","src":"$.status"}
		]
	}
}
```

参数名     | 类型    | 描述       
--- | --- | ---
request.data_type | int array   | 请求数据类型，0:文本、1:音频、2:图像、3:视频，数组的大小和映射到后端的data数组大小保持一致 
rule               | array | 前后端映射规则，dst为后端协议data参数中定义的字段，，src为Web API暴露给开发者参数data中定义的字段 


 <strong>说明：</strong> rule字段中用到的符号“$”表示请求协议中暴露给用户的data参数。


 <strong>参考信息：</strong>

 后端协议中的data参数示例：
    
```json
{
	"data": [{
		"data_id": "",         
		"frame_index": 1,
		"data_type": 1,
		"status": 1,
		"desc_args": {},
		"format": "",
		"encoding": "",
		"data": ""
	}]
}
    
```
    
后端协议参数说明：
    
参数名|类型|说明
-----|-----|-----
data_id |string | 数据编号  
frame_id |uint32 | 数据序号  
data_type |DataType |枚举类型，用于区分数据类型<br> TEXT     = 0;    // 文本<br>AUDIO    = 1;    // 音频<br>IMAGE    = 2;    // 图像<br> VIDEO    = 3;    // 视频 
status |DataStatus | 枚举类型，用于区分数据状态<br>BEGIN	= 0;//开始<br>CONTINUE	= 1;//跟流<br>END		= 2;//结束<br>ONCE		= 3;//一次调用结束  
desc_args |map(string,bytes) | 数据描述参数  
format |string | 数据的编码格式 
encoding |string | 数据的压缩格式
data |bytes |数据内容  
    
- ##### 第四步，配置返回参数映射规则
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;与上一步中的配置请求参数映射的背景相同，上游服务返回给webgate的协议中，data也是一个抽象的协议，为了降低理解和使用Web API的门槛，我们对后端返回协议中的data参数也做了一层映射，最终暴露给用户简洁明了的返回参数。

 <strong>服务示例：</strong>

 ```json
{
	"response.data.mapping": {
		"data_type": 0,
		"rule": [
             {"dst":"$.result","src":"$[0].data"},
             {"dst":"$.status","src":"$[0].status"}
		]
	}
}
```

 参数名 | 类型  | 描述                                                         
-----| ------ | --------
 data_type | int  array  | 返回数据的类型，0:json、1:普通字符串、2:字节流  其中json处理失败会降级为字节流处理         
 rule  | array | 前后端映射规则，src为后端协议data参数中定义的字段，dst为Web API暴露给开发者的参数data中定义的字段 


 <strong>说明：</strong>
    
- ##### 第五步，配置其他AI能力
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;如果希望在一个webgate中接入多个AI能力，只需要定义多组配置即可，示例如下：
   
 ```json
[{
    "version": "1.0",
    "service": "svc0",
    "route": "/v2/svc0",
    "request.data.mapping": [],
    "response.data.mapping": [],
    "schema": {}
}，
{
    "version": "1.0",
    "service": "svc01",
    "route": "/v2/svc01",
    "request.data.mapping": [],
    "response.data.mapping": [],
    "schema": {}
}]
```


## 服务完整配置示例

```json
[
  {
    "service": "svc",
    "version": "1.0",
    "call":"atmos-svc",
    "route": "/v2/svc",
    "request.data.mapping": {
      "data_type": [2],
      "rule":[
        {"dst":"$","src":"$"}
      ]
    },
    "response.data.mapping": {
      "data_type": [2],
      "rule":[
        {"dst":"$","src":"$"}
      ]
    },
    "schema": {
      "type":"object",
      "properties":{
        "common":{
          "type":"object",
          "properties":{
            "app_id":{
              "type":"string",
              "maxLength":50
            },
            "uid":{
              "type":"string",
              "maxLength":50
            }
          },
          "required":["app_id"]
        },
        "business":{
          "type":"object",
          "properties":{
            "param1":{
              "type": "string",
              "maxLength": 30
            },
            "param2":{
              "type": "string",
              "maxLength": 30
            },
            "param3":{
              "type": "number"
            }
          }
        },
        "data":{
          "type":"object",
          "properties":{
            "status":{
              "type":"integer"
            },
            "data":{
              "type": "string",
              "maxLength": 30000000
            }
           
          }
        }
      }
    }
  }
]
```

#### schema 参数扩展属性

key|value|作用
---|---|---
replaceKey|string|将参数的替换为replaceKey的值
constVal|string|无论用户是否传了该参数，改参数都会是该值
defaultVal|string|用户没有传该参数时，默认会是该值

三个属性的作用顺序为replaceKey,constVal,defaultVal
