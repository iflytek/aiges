package finder

type ReturnCode int

func (retCode ReturnCode) String() string {
	if errString, ok := retCodeToString[retCode]; ok {
		return errString
	}
	return "unknow error for " + string(retCode)
}

var retCodeToString = map[ReturnCode]string{
	Success:                          "成功",
	InvalidParam:                     "无效的参数",
	MissCompanionUrl:                 "缺少companionUrl",
	ConfigMissName:                   "[config] 配置信息中缺少配置的名字",
	ConfigMissCacheFile:              "[config] 丢失缓存文件",
	ZkGetInfoError:                   "获取zk信息错误",
	ZkInfoMissConfigRootPath:         "获取的zkInfo中不存在config path",
	ZkInfoMissServiceRootPath:        "获取的zkInfo中不存在server path",
	ZkInfoMissAddr:                   "获取的zkInfo中不存在zk的地址信息",
	ZkInfoAddrConvertError:           "转换获取的zkInfo中的zk地址信息出错",
	ZkParamsMissServers:              "zk参数中缺少zk的服务地址信息",
	ZkParamsMissSessionTimeout:       "zk参数中缺少sessionTimeout的配置信息",
	ZkDataCanotNil:                   "zk中不能设置nil的数据",
	CompanionRegisterServiceErr:      "向companion注册服务失败",
	FeedbackServiceErr:               "feedback service数据到companion失败",
	FeedbackConfigErr:                "feedback config数据到companion失败",
	FeedbackPostErr:                  "向companion feedback post请求时出错",
	DecodeVauleDataEmptyErr:          "解码数据的时候，数据为空",
	DecodeVauleDataNotFullErr:        "解码数据的时候，数据不完整",
	DecodeVauleDataFormatErr:         "解码数据的时候，数据格式出错",
	ServiceMissItem:                  "[service] 没有service的信息",
	ServiceMissAddr:                  "[service] 缺失service对应的地址信息",
	JsonUnmarshalErr:                 "Json在反序列化的时候出错",
	JsonMarshalErr:                   "json在序列化的时候出错",
	ZkPathCannotNil:                  "zk中的path不能为空",
	ZkPathPrefixIllegal:              "zk中的path开始符号必须是 / ",
	ZkPathSuffixIllegal:              "zk中的path不能以 / 结束",
	ZkPathNullCharacterNotAllowed:    "zk中的path字符不能有null",
	ZkPathEmptyNodeNameNotAllowed:    "zk中的path不允许存在空的节点名字",
	ZkPathRelativePathNotAllowed:     "zk中的path不允许相对路径",
	ZkPathInvalidCharacterNotAllowed: "zk中的path不允许非法字符",
	ZkGetDataErr:                     "从zk中获取数据出错",
	ServiceMissApiVersion:            "[service] 缺失版本号",
	ZkGetNilData:                     "zk中节点上的数据为空",
	ZkConnectionLoss:                 "zk连接不存在",
	ZkInfoMissZkNodePath:             "zk的地址所在的节点信息不存在",
	ConfigFileNotExist:               "配置文件不存在",
	ConfigDirNotExist:                "配置目录不存在",
}

const (
	ConfigSuccess     = iota // 0 获取配置成功
	ConfigReadFailure        // 1 读数据失败
	ConfigLoadFailure        // 2 加载配置失败
)

const (
	Success          ReturnCode = 0
	InvalidParam     ReturnCode = 10000
	MissCompanionUrl ReturnCode = 10001
)

//config相关错误
const (
	ConfigMissName ReturnCode = 10100 + iota
	ConfigFileNotExist
	ConfigMissCacheFile
	ConfigDirNotExist
)

//zk相关错误
const (
	ZkGetInfoError ReturnCode = 10200 + iota
	ZkGetNilData
	ZkConnectionLoss
	ZkInfoMissRootPath
	ZkInfoMissConfigRootPath
	ZkInfoMissServiceRootPath
	ZkInfoMissAddr
	ZkInfoMissZkNodePath
	ZkInfoAddrConvertError
	ZkGetDataErr
	ZkParamsMissServers
	ZkParamsMissSessionTimeout
	ZkDataCanotNil
	ZkPathCannotNil
	ZkPathPrefixIllegal
	ZkPathSuffixIllegal
	ZkPathNullCharacterNotAllowed
	ZkPathEmptyNodeNameNotAllowed
	ZkPathRelativePathNotAllowed
	ZkPathInvalidCharacterNotAllowed
)

//service相关错误
const (
	ServiceMissAddr ReturnCode = 10300 + iota
	ServiceMissItem
	ServiceMissApiVersion
)

//feedback相关错误
const (
	FeedbackConfigErr ReturnCode = 10400 + iota
	FeedbackServiceErr
	FeedbackPostErr
)

//Companion相关错误
const (
	CompanionRegisterServiceErr ReturnCode = 10500 + iota
)

const (
	DecodeVauleDataEmptyErr ReturnCode = 10600 + iota
	DecodeVauleDataNotFullErr
	DecodeVauleDataFormatErr
)

//json处理数据的时候出错
const (
	JsonUnmarshalErr ReturnCode = 10700 + iota
	JsonMarshalErr
)
