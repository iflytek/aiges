# Athena Serving

<!-- markdownlint-capture -->
<!-- markdownlint-disable MD033 -->

<span class="badge-placeholder">[![Build Status](https://img.shields.io/drone/build/thegeeklab/hugo-geekdoc?logo=drone&server=https%3A%2F%2Fdrone.thegeeklab.de)](https://drone.thegeeklab.de/thegeeklab/hugo-geekdoc)</span>
<span class="badge-placeholder">[![GitHub release](https://img.shields.io/github/v/release/xfyun/aiges)](https://github.com/xfyun/aiges/releases/latest)</span>
<span class="badge-placeholder">[![License: Apache2.0](https://img.shields.io/github/license/xfyun/aiges)](https://github.com/xfyun/aiges/blob/master/LICENSE)</span>

<!-- markdownlint-restore -->
## 整体架构

![img](https://raw.githubusercontent.com/xfyun/proposals/main/athenaloader/athena.png)

## 使用工作流

![img](https://github.com/xfyun/proposals/blob/main/athenaloader/usage.png?raw=true)

## 开源计划

| 任务项 |目标 |时间 |
|-----|-----|-----|
|&#9745; [加载器](#AIGes) 通用引擎/模型加载器|独立部署可运行，支持python快速推理服务化|2022/Q2|
|&#9745; [负载均衡器](#LoadBalance) 负载聚合组件|独立部署可运行|2022/Q3|
|&#9744; [WebGate](#Webgate) Web网关组件|可运行|2022/Q3|
|&#9745; [Polaris](#Polaris) 配置中心与服务发现|独立部署可运行|2022/Q2|
|&#9744; [Atom](#Atom) 协议转换组件|可运行|2022/Q3|
|&#9744; Serving on Kubernetes Helm Chart一键部署 (进行中)|支持在k8s集群上一键部署推理服务框架|2022/Q3|
|&#9744; Serving on Docker with docker-compose 一键部署|支持使用docker-compose部署推理服务框架|2022/Q4|
|&#9744; 各组件Documentation建设 (进行中)|各组件文档详设，门户建设|2022/Q4|
|&#9744; 多领域模型Demo演示示例、GIF (进行中)|部分领域模型推理示例，如mmdetection,mmocr,yolo等|2022/Q3|
|&#9744; AIServing [API](#API协议) AI能力协议规范|完善开源协议说明，schema自动生成、校验工具|2022/Q3|
|&#9744; AseCTl命令行工具 [API](#Asectl) 命令行工具|支持能力一键生成，运行，配置管理等|2022/Q4|



## AIGes (AI General Engine Service)

通用引擎加载器(部分文档中loader， loader engine均为别名)

### Documentation

View Doc on [Documentation](https://xfyun.github.io/inferservice/architechture/architechture/)


### Features


&#9745; 支持模型推理成RPC服务(Serving框架会转成HTTP服务)
&#9745; 支持C代码推理 support c++/c code infer
&#9745; 支持Python代码推理 Support python code infer
&#9745; 支持配置中心，服务发现
&#9745; 支持负载均衡配置
&#9744; 支持Java代码推理或者其它
&#9744; 支持计量授权

## Polaris

Go编写的服务发现、配置中心

### Documentation

View Doc on [Documentation](https://xfyun.github.io/inferservice/architechture/architechture/)

### Features
- Full category for managing service
- Support multi service versions
- Support roll back config
- Support feedback for pushing config
- Support management for provider and consumer online
- High available for some not expected cases
- Easy Integration
- Support Delivery by docker 

## API协议

### 1. 整体体说明
本协议面向AI PaaS项目进行了面向应用者开发、内部业务协议进行了约束和定义

### 2 面向应用开发者协议
协议采用json进行描述，

####  2.1 元数据
由于AI的服务请求一般是用于数据计算与处理，在交互的往来协议中，需要携带一些基础数据。本协议将数据进行抽象约束。为提高适配性，本协议对元数据进行了枚举定义。元数据结构中定义了对应数据类型基本的描述进行了定义，也允许增加K-V单层结构的的描述。用户在定义数据时，选择数据类型后，系统会主动提供结构字段。

#####  2.1.1 文本描述
结构举例：

    {
		"encoding":"utf8", 
		"status":0,
		"seq": 1，
		"compress":"gzip",
		"custom1":"zz",
		"text":"hello word"
    } 

字段 | 含义 |  类型| 说明 |
-|-|-|-
|encoding | 文本编码 | string | 取值范围可枚举|
|status | 数据状态 | int | 取值范围为0（开始）、1（继续）、2（结束）、3（一次传完） |
|text | 文本数据 | string | 为文本数据，base64|
|seq| 数据序号 | int | 标明数据为第几块，可选|
|custom1 | 用户自定义参数 | string或 int | 用户自定义，一级参数|

##### 2.1.2 音频描述

	{
		"status":0，
		"sample_rate":16000,
		"channels":1,
		"bit_depth":16,
		"encoding":"opus",
		"seq": 1,
		"custom1":"zz",
		"audio":"xxxxxxxxxxxxxxxxxxxxxxxxxx" # 音频数据
	} 

字段 | 含义 |  类型| 说明 |
-|-|-|-
|encoding | 音频编码 | string | 取值范围可枚举|
|status | 数据状态 | int | 取值范围为0（开始）、1（继续）、2（结束）、3（一次传完） |
|audio | 音频数据 | string | 为文本数据，base64|
|sample_rate | 采样率 | int | 音频采样率，可枚举|
|channels | 声道数 | int | 声道数，可枚举|
|bit_depth | 位深 | int | 单位bit，可枚举|
|seq| 数据序号 | int | 标明数据为第几块，可选|
|custom1 | 用户自定义参数 | string或 int | 用户自定义，一级参数|

##### 2.1.2 待完善
1 图像
2 音频
3 其他二进制数据等


####  2.2  请求协议说明
分为三部分：平台参数、能力参数、传输数据。请求协议举模板如下：

    {
	     "header":{},
	     "parameter":{
		     "service_iat":{
			      "accept_1":{
			    }
		     }
	     }, 
	     "payload":{
	     	"service_iat":{
	     	}
	     }
    }

字段说明：

字段 | 含义 |  类型| 说明 |
-|-|-|-
header |平台参数 | Object | 用于控制平台特性的参数，如appid等，此结构只有一个层级，开发者不可修改。|
parameter | 能力参数 |Object | 用于控制AI引擎特性开关的。该段将会被直接透传至引擎,多层级结构,结构确定,未枚举结构为一级结构|
payload | 输入数据|Object | 用于携带请求的数据，多层级结构,结构确定 |
accept_1 | 输入数据|Object | 用户自定义,用于描述返回结果的编码\数据格式的参数,此结构由元数据中定义的属性中定义的 |
service_iat | 输入数据|Object | 用户自定义, 可默认与accept_1一致。 用于携带请求的数据，多层级结构,结构确定 |

##### 2.2.1 单能力单数据举例如


    {
	    "header":{
		"appid":"1234",
		"uid":"d233"
	    },
	    "parameter":{
		    "service_iat":{
		       "language": "zh_cn",
		       "domain":"iat",
		       "accent": "mandarin"
		       }
	    }, 
	    "payload":{
	    	"service_iat":{
			"status":0,
			"format":"audio/L16;rate=16000",
			"encoding":"raw",
			"audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
			}
	     }
     } 


##### 2.2.2  单能力单多数据举例

		{
	    "header":{
			"appid":"1234",
			"uid":"1234"
	    },

	    "parameter":{
	      "service_iat":{
		       "language": "zh_cn",
		       "domain":"iat",
		       "accent": "mandarin"
		       }
	    },

	    "payload":{
			"data_1":{
				"status":0,
				"format":"audio/L16;rate=16000",
				"encoding":"raw",
			      "audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
		     },
		     "data_2":{
				"status":0,
				"format":"audio/L16;rate=16000",
				"encoding":"raw",
			       "audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
	     }
	}
	} 


结构说明：
1.  此时payload为多级结构
2.  多数据情况想，payload通过不同的数据ID作为key来描述，如示例中的data_1、data_2
3. 数据ID为用户自定义


#####   2.2.3 多数据范围特性时的描述
针对存在多个数据流返回的，需要在输入参数（parameter段）中指定多个数据特性的描述。

	"parameter":{
	      "service_iat":{
		       "language": "zh_cn",
		       "domain":"iat",
		       "accent": "mandarin",
		       "accept_1":{
				"format":"audio/L16;rate=16000",
				"encoding":"raw"
			    },
			   "accept_2":{
				"format":"audio/L16;rate=16000",
				"encoding":"raw"
		       }
	       ｝

	    } 

结构说明：
1.  此时parameter为多级结构
2.  多数据情况下，parameter通过不同的ID作为key来对不同的结果进行描述，如示例中的accept_1、accept_2
3.  Accept ID为用户自定义

##### 2.2.4 多能力多数据描述
		{
		    "header":{
				"appid":"1234",
				"uid":"1234"
		    },
		    "parameter":{
				"service_1":{

				       "language": "zh_cn",
				       "domain":"iat",
				       "accent": "mandarin"
				},
				"service_2":{

				       "language": "zh_cn",
				       "domain":"iat",
				       "accent": "mandarin"
				}
		    },

		    "payload":{
				"service_1_data_1":{
					"status":0,
					"format":"audio/L16;rate=16000",
					"encoding":"raw",
					 "audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
				      },

				 "service_2_data_1":{
					"status":0,
					"format":"audio/L16;rate=16000",
					"encoding":"raw",
					 "audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
				  },
				 "service_2_data_2":{
					"status":0,
					"format":"audio/L16;rate=16000",
					"encoding":"raw",
				 	"audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
				     }
			}
		} 

结构说明：
1.  此时parameter为多级结构，此时payload为多级结构
2.  多数据情况下，parameter通过不同的能力ID作为key来分别描述不同能力的特性描述，service_1、service_2
3. 能力ID为用户自定义
4.  多数据情况想，payload通过不同的数据ID作为key来描述，如示例中的service_1_data_1、service_2_data_1、service_2_data_2
5. 此场景下数据ID为编排平台的借口生成

##### 2.2.5  流式场景下，后续数据包描述
header、parameter为可选，一般情况不携带。以多数据为例：

	 {
	 	"header":{
	 	     "sid":"iat000704fa@dx16ade44e4d87a1c802", #可选	
	 		},
		
		"payload":{
			"data_1":{
				"status":1,
				"format":"audio/L16;rate=16000",
				"encoding":"raw",
		                "audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
		     },

		 	"data_2":{
				"status":1,
				"format":"audio/L16;rate=16000",
				"encoding":"raw",
		                "audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
		     }
	 } 

字段说明：

字段 | 含义 |  类型| 说明 |
-|-|-|-
sid |平台返回会话句柄 | Object | 返回消息中携带的会话句柄，可选|

####  2.3 返回协议说明
返回协议定义了会话的计算状态，以及数据段

##### 2.3.1 单输出描述



	{
	  "header":{
		 "code": 0, 
		 "message": "success", 
		 "sid": "iat000704fa@dx16ade44e4d87a1c802"

	},
	 # 此结构为元数据结构（data），描述返回结果。
	 "payload": {
		 “result_1”：{
			"status":0, #数据状态
			"format":"audio/L16;rate=16000",
			"encoding":"raw",
		       "audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
		    }
	   }
	}

字段 | 含义 |  类型| 说明 |
-|-|-|-
code |错误码 | int | 可枚举|
message | 错误描述 |string | 错误信息描述|
sid |平台返回会话句柄 | Object | 返回消息中携带的会话句柄，可选|
payload | 输入数据|Object | 用于携带返回的数据，元数据结构见 [2.1 元数据 的定义](”2.1 元数据 的定义“) |

##### 2.3.2 多输出描述

	{
		"header":{
			"code": 0,
			 "message": "success",
			 "sid": "iat000704fa@dx16ade44e4d87a1c802",
			 "status":0		
			},
	 
	"payload":{
		 "result_1":{
			"status":0,
			"format":"audio/L16;rate=16000",
			"encoding":"raw",
		    "audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
		     },

		 "result_2":{
			"status":0,
			"format":"audio/L16;rate=16000",
			"encoding":"raw",
			"audio":"exSI6ICJlbiIsCgkgICAgInBvc2l0aW9uIjogImZhbHNlIgoJf..."
		     }
	}
	} 

字段说明：
同 [2.3.1 单输出描述]("2.3.1 单输出描述")
结构说明：
1.  此时payload为多级结构
2.  多数据情况下，payload通过不同的结果ID作为key来分别描述不同的结果特性描述，如：result_1、result_2
3.  结果ID为用户自定义

#### 2.4 系统参数的约束

##### 2.4.1 header字段
字段 | 含义 |  类型| 说明 |
-|-|-|-
app_id |应用id | string |必选|
ath_id | 三方用户ID |string | 可选|


### 3.1 路由信息
新增session状态（SessState），用于描述流式与否

	message GlobalRoute {
	    string session_id = 1; //session id
	    string trace_id = 2; //trace id
	    string up_router_id = 3; //上行数据路由标识
	    string guider_id = 4; //调度中心标识
	    string down_router_id = 5; //下行数据路由标识
	    string appid = 6; //应用标识
	    string uid = 7; //用户标识
	    string did = 8; //设备标识
	    string client_ip = 9; //客户端ip
	    
	    SessState  session_state = 10; //新增：会话状态，流式、非流式,STATE为枚举

	} 

#### 3.1 元数据调整
去除原format、encoding字段，描述归入：desc_args字段，且将原来的byte类型调整为string类型

	message GeneralData {
	    string data_id = 1; //数据编号
	    uint32 frame_id = 2; //数据序号

	    //区分数据类型
	    enum DataType {
		TEXT = 0; // 文本
		AUDIO = 1; // 音频
		IMAGE = 2; // 图像
		VIDEO = 3; // 视频
	    }

	    DataType data_type = 3; //数据类型

	    //区分数据状态
	   enum DataStatus {
		BEGIN = 0; //开始
		CONTINUE = 1; //跟流
		END = 2; //结束
		ONCE = 3; //一次调用结束
	    }

	    DataStatus status = 4; //数据状态

	    map<string, string> desc_args = 5; //数据描述参数

	    bytes data = 8; //数据内容
	}


