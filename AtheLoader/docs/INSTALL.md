##编译
	本组件集成CI构建，如果使用非CI构建方式请参照以下步骤构建。
	1.请区分编译机器已经安装golang程序并且版本号 >1.9.2
	2.下载对应tag代码并执行解压操作
	3.进入AIGES目录，执行对应shell脚本即可
	├── build.sh		// 框架构建脚本;
	└── demobuild.sh	// iat demo构建脚本;
	4.构建完成于本地生成组件AIservice || xats(其中AIservice为框架服务组件, xats为iat demo服务组件)


##demo部署
	引擎服务部署依赖引擎内核相关组件及其他三方组件,相关组件已拷贝至depends目录;若使用CI构建则直接拉取docker镜像即可，若手动构建可按以下方案部署
	1. 将当前目录的xats拷贝至demo/depends/iat-gpu01/bin/目录
	2. 将demo目录的iat.toml拷贝至demo/depends/iat-gpu01/bin/目录
	3. 将src/cLibrary/*拷贝至demo/depends/iat-gpu01/bin/目录
	3. 将demo/depends/iat-gpu01/*拷贝至部署目录xxx
	4. 设置环境变量export GODEBUG=cgocheck=0
	5. 进入xxx/bin/按demo/start.sh方式启动服务
	├── src
	└── demo


	使用CI构建拉取docker镜像使用方式如下:
	1. 挂载资源目录/msp/resource至docker容器; docker -v挂载
	2. 挂载日志目录/log/server至docker容器; docker -v挂载
	3. GPU版本docker部署时需注意挂载相关nvidia设备; docker --device挂载
	
	eg: docker run -itd --name devTest --net="host" -v /msp/resource/sms:/msp/resource/sms -v /log/server:/log/server --device=/dev/nvidia0:/dev/nvidia0 --device=/dev/nvidia1:/dev/nvidia1 --device=/dev/nvidiactl:/dev/nvidiactl --device=/dev/nvidia-uvm:/dev/nvidia-uvm --device=/dev/nvidia-uvm-tools:/dev/nvidia-uvm-tools --privileged=true 172.16.59.153/aiaas/aiges:1.0.1


##命令行
	服务运行需以命令行方式启动(见iatstart.sh),对命令行参数作相关解释如下
	-v					打印服务版本信息
	-m      	int		配置模式|cfgmode(0:本地配置, 1:配置中心)
	-c	        string	配置文件名|cfgname
	-u		    string	配置中心|url
	-p      	string	项目名|project
	-g  		string	集群名|group
	-s      	string	服务名|service

##配置项
	对AIGES服务框架相关配置项作相关解释如下,其他框架配置项解释见xrpc相关文档
	[aisrv] # 该配置section需与启动命令中service保持一致
	finder = 0					# 是否开启服务发现
	#host = "172.16.154.100"	# 服务地址(该项缺省不配置,取首个非回环地址)
	port = 5090					# 服务端口
	report = 1					# 负载上报间隔时间(单位:s)

	[aiges]
	sessMode = 1				# 是否会话模式,缺省1
	numaNode = -1				# 是否设置numa绑定(-1:不设置, 其他:cpu Node)
	elogRemote = 0				# 是否开启eventlog远端上传
	elogLocal = 0				# 是否开启eventlog本地dump
	elogHost = "127.0.0.1"		# eventlog flume地址
	elogPort = "4545"			# eventlog flume端口
	elogXml = "seelog.xml"		# eventlog 本地日志配置
	elogSpill = "/log/server/iatspill"  						# eventlog spill路径,建议挂载至docker外
	libCodec = "libamr.so;libamr_wb.so;libspeex.so;libico.so" 	# 相关编解码库

	[aires] # 引擎所依赖的基础资源配置; "资源类型" = "资源路径"
	HMM_16K = "/msp/resource/sms/acmod_16KRnn_sms.bin"
	WFST =  "/msp/resource/sms/wfst.bin"
	LM = "/msp/resource/sms/gram.bin"
	RLM = "/msp/resource/sms/nextg.rnnlmwords.bin"

##建议
	由于GPU服务器单台双进程授权过高,为防止docker崩溃导致授权急剧下降,建议每个docker一个服务进程;