package schemas

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"git.xfyun.cn/AIaaS/finder-go"
	common "git.xfyun.cn/AIaaS/finder-go/common"
	js "git.xfyun.cn/AIaaS/json_script"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"github.com/oliveagle/jsonpath"
	"github.com/qri-io/jsonschema"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

/**
	This module is used for loading AI ability and define the shcema of a AI ability interface.
The part of schema define the schema of input of request parameters.The request parameters will be limit
by the rules defined in schema.
	The part of mapping defines how to output the upstream responce to client.
We can easily convert the format of upstream response to any format we want to output to the client by
define some mapping rules ahead.
	The part of script will be used to handle request parameters.It's a light and effective rule engine which depende on
json to generate ast.
**/
var RouteMappingCache = &MappingCache{}

var Logger *utils.Logger

type RouteMapping struct {
	Version      string                 `json:"version"`
	Service      string                 `json:"service"`
	Route        string                 `json:"route"`
	CallService  string                 `json:"call"`
	RequestData  *RequestData           `json:"request.data.mapping"`
	ResponseData *ResponseData          `json:"response.data.mapping"`
	Schema       *jsonschema.RootSchema `json:"schema"`    // 的map
	AtmosMap     map[string]string      `json:"atmos_map"` //调用的atmos映射map
}

type Rule struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type RequestData struct {
	DataType []int32 `json:"data_type"`
	Rule     []Rule  `json:"rule"`
	Script   *Script `json:"script"`
}

func GetMapping(route, version, sub string) *RouteMapping {
	return GetMappingByKey(sub + route + version)
}

//响应的数据类型格式
const (
	//json类型
	DATA_TYPE_JSON uint32 = 0
	//普通字符串类型
	DATA_TYPE_STRING uint32 = 1
	//byte类型
	DATA_TYPE_BYTE uint32 = 2
)

// 表示响应映射关系的结构
type ResponseData struct {
	DataType []uint32 `json:"data_type"`
	Rule     []Rule   `json:"rule"`
	Script   *Script  `json:"script"`
}

//根据json path获取value
func GetByJPath(jsonData interface{}, jPath string) (interface{}, error) {
	return jsonpath.JsonPathLookup(jsonData, jPath)
}

func GetMappingByKey(key string) *RouteMapping {
	return RouteMappingCache.Get(key)
}

//func GetSchemaByKey(service string)*jsonschema.RootSchema  {
//	return SchemaCaches.Get(service)
//}
func (mp *RouteMapping) GetAtmos(sub string) string {
	//mp:=RouteMappingCache.Get(key)
	if mp == nil {
		return ""
	}
	if mp.AtmosMap == nil {
		return mp.CallService
	}

	atmos := mp.AtmosMap[sub]
	if atmos != "" {
		return atmos
	}
	return mp.CallService
}

//获取响应数据类型
func GetRespDataType(key string) []uint32 {
	respData := GetMappingByKey(key).ResponseData
	return respData.DataType
}

//获取回调请求结果
func GetUpCallReqByCall(key string, jsonData interface{}) ([]map[string]interface{}, error) {
	mapping := GetMappingByKey(key)
	if mapping == nil {
		return make([]map[string]interface{}, 0), errors.New("ability unsupported")
	}
	result, err := mapping.ResolveUpCallReq(jsonData)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//获取回调响应结果
func GetRespByCall(key string, jsonData interface{}) (interface{}, error) {
	mapping := GetMappingByKey(key)
	if mapping == nil {
		return make([]interface{}, 0), nil
	}

	result, err := mapping.ResolveResp(jsonData)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//解析上传请求参数
func (routeMap *RouteMapping) ResolveUpCallReq(jsonData interface{}) ([]map[string]interface{}, error) {
	if jsonData == nil {
		return nil, nil
	}
	resolveResult := []map[string]interface{}{}
	var i interface{}
	//reqDate["data_type"] = elem.DataType
	for _, value := range routeMap.RequestData.Rule {
		if value.Src == "" {
			return nil, fmt.Errorf("server config error %v", routeMap.Service)
		}
		res, err := CachedJsonpathLookUp(jsonData, value.Src)
		if err != nil {
			Logger.Warnf("resolve request error,retrive %s err:%v", value.Src, err)
			continue
			//return nil, errors.New("cannot parse request:"+err.Error())
		}
		MarshalInterface(value.Dst, &i, res)
		//if err != nil{
		//	return nil,errors.New("server error,cannot assemble request : invalid expression")
		//}

	}

	ii, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("server config error: request mapping error,request data must be array")
	}
	if len(ii) > len(routeMap.RequestData.DataType) {
		return nil, fmt.Errorf("server config error:length of request data_type array  is not equal to that of data array")
	}
	for idx, v := range ii {
		mi, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("server config error:request datas item must be map")
		}
		dtp := routeMap.RequestData.DataType[idx]
		if dtp != -1 {
			mi["data_type"] = dtp
		} else {
			//dataType,err:=CachedJsonpathLookUp(jsonData,"$[0].data_type")
			//if err !=nil{
			//	Logger.Warnf("cannot find dataType in request data error:%v",err)
			//}
			//mi["data_type"] = dataType
		}
		resolveResult = append(resolveResult, mi)
	}
	return resolveResult, nil
}

//解析响应参数
func (routeMap *RouteMapping) ResolveResp(jsonData interface{}) (interface{}, error) {
	if jsonData == nil {
		return nil, nil
	}
	var respData interface{}
	if routeMap.ResponseData.Script != nil {
		ExecuteScript(routeMap.ResponseData.Script.GetEvery(), jsonData)
	}
	//create response by specific rules
	for _, value := range routeMap.ResponseData.Rule {
		if value.Src == "" {
			return nil, fmt.Errorf("server config error %v", routeMap.Service)
		}
		res, err := CachedJsonpathLookUp(jsonData, value.Src)
		if err != nil {
			Logger.Warnf("resolve response error:retrive field %s error:%v", value.Src, err)
			continue
		}
		MarshalInterface(value.Dst, &respData, res)
	}

	return respData, nil
}

func LoadRoteMapping(mappingBuffer []byte) error {
	mappings := &[]RouteMapping{}
	err := json.Unmarshal(mappingBuffer, mappings)
	if err != nil {
		return err
	}

	var errs = make([]error, 0)
	for _, elem := range *mappings {
		routeMapping := elem
		ok := checkMappingRules(&routeMapping, &errs)
		if routeMapping.RequestData.Script != nil {
			err := ParseScript(routeMapping.RequestData.Script)
			if err != nil {
				fmt.Println("ERROR:parse script:" + err.Error())
			}

		}
		if routeMapping.ResponseData.Script != nil {
			err := ParseScript(routeMapping.ResponseData.Script)
			if err != nil {
				fmt.Println("ERROR:parse script:" + err.Error())
			}
		}
		if ok {
			RouteMappingCache.Set(GetMappingKey(elem.Service, elem.Route, elem.Version), &routeMapping)
		}
	}
	if err := checkMapping(); err != nil {
		return err
	}

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println("ERROR:", e)
		}
		return errors.New("check mapping rule error!!!")
	}
	fmt.Println(".......................  success load mapping ...........................")
	return nil
}

func GetMappingKey(service, route, version string) string {
	return service + route + version
}

func LoadRouteMappingFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return LoadRoteMapping(data)
}

//check if schema is loaded expected
func checkMapping() error {

	RouteMappingCache.Range(func(k, val interface{}) bool {
		v := val.(*RouteMapping)
		if v.Service == "" || v.Route == "" || v.CallService == "" {
			//log.("--------------mapping has error at " + k.(string))
			fmt.Println("load mapping error--------------mapping has error at " + k.(string))
		}
		fmt.Println("name:", k, "mapping:", v.ResponseData.DataType, )

		return true
	})

	return nil
}

//check mapping

func checkMappingRules(mapping *RouteMapping, errs *[]error) bool {
	var cers []error
	//check request
	for _, m := range mapping.RequestData.Rule {
		if !checkRule(m.Src) {
			cers = append(cers, fmt.Errorf("mapping rule has error at %s %s request.data.mapping rule.src:'%s'", mapping.Service, mapping.Version, m.Src))
		}
		if !checkRule(m.Dst) {
			cers = append(cers, fmt.Errorf("mapping rule has error at %s %s request.data.mapping rule.dst:'%s'", mapping.Service, mapping.Version, m.Dst))
		}
	}

	for _, m := range mapping.ResponseData.Rule {
		if !checkRule(m.Src) {
			cers = append(cers, fmt.Errorf("mapping rule has error at %s %s response.data.mapping rule.src:'%s'", mapping.Service, mapping.Version, m.Src))
		}
		if !checkRule(m.Dst) {
			cers = append(cers, fmt.Errorf("mapping rule has error at %s %s response.data.mapping rule.dst:'%s'", mapping.Service, mapping.Version, m.Dst))
		}
	}
	*errs = append(*errs, cers...)

	if len(cers) > 0 {
		return false
	}
	return true

}

type MappingCache struct {
	sync.Map
}

func (m *MappingCache) Get(k string) *RouteMapping {
	if v, ok := m.Load(k); ok {
		return v.(*RouteMapping)
	}
	return nil
}

func (m *MappingCache) Set(k string, v *RouteMapping) {
	m.Store(k, v)
}

func (m *MappingCache) RangeM(f func(k string, val *RouteMapping) bool) {
	m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*RouteMapping))
	})
}

//jsonpath cache map

var jsonPathComplied sync.Map

func CachedJsonpathLookUp(obj interface{}, jpath string) (interface{}, error) {
	c, ok := jsonPathComplied.Load(jpath)
	if !ok {
		co, err := jsonpath.Compile(jpath)
		if err != nil {
			return nil, err
		}
		jsonPathComplied.Store(jpath, co)
		c = co
	}

	return c.(*jsonpath.Compiled).Lookup(obj)
}

//schema caches
//var SchemaCaches = &SchemaCache{}
const (
	SchemaFilePrerix = "schema_"
	SchemaFileSuffix = ".json"
)

var (
	SchemaFilePrefixLen = len(SchemaFilePrerix)
	SchemaFileSuffixLen = len(SchemaFileSuffix)
)

func IsSchemaFile(key string) bool {
	return strings.HasPrefix(key, SchemaFilePrerix) && strings.HasSuffix(key, SchemaFileSuffix)
}

func assembleSchemaFileName(service string) string {
	return SchemaFilePrerix + service + SchemaFileSuffix
}

func loadMappingByConfigs(fm map[string]*common.Config) error {
	for _, f := range fm {
		err := LoadRoteMapping(f.File)
		if err != nil {
			return err
		}
	}
	return nil
}

// 加载mapping 文件

var loadedSchemas = sync.Map{}

func LoadMapping(finderManager *finder.FinderManager, services []string, configChangeHandler common.ConfigChangedHandler) error {
	for _, service := range services {
		loaded, ok := loadedSchemas.Load(service)
		if (!ok) || (!loaded.(bool)) {
			loadedSchemas.Store(service, true)
			fm, err := finderManager.ConfigFinder.UseAndSubscribeConfig([]string{assembleSchemaFileName(service)}, configChangeHandler)
			if err != nil {
				return fmt.Errorf("get file error:%s,%v", assembleSchemaFileName(service), err)
			}
			if err := loadMappingByConfigs(fm); err != nil {
				return errors.New(assembleSchemaFileName(service) + " " + err.Error())
			}
		}
	}
	return nil
}

//script

type Script struct {
	first  js.Exp // 只在请求第一帧执行
	every  js.Exp // 在请求的每一帧都会执行
	before js.Exp
	Before interface{} `json:"before"`
	First  interface{} `json:"first"`
	Every  interface{} `json:"every"`
}

func (s *Script) GetFirst() js.Exp {
	return s.first
}

func (s *Script) GetEvery() js.Exp {
	return s.every
}

func (s *Script) GetBefore() js.Exp {
	return s.before
}

func ParseScript(script *Script) (error) {
	if script.First != nil {
		scp, err := js.CompileExpFromJsonObject(script.First)
		if err != nil {
			return err
		}
		script.first = scp
	}
	if script.Every != nil {
		scp, err := js.CompileExpFromJsonObject(script.Every)
		if err != nil {
			return err
		}
		script.every = scp
	}
	if script.Before != nil {
		scp, err := js.CompileExpFromJsonObject(script.Before)
		if err != nil {
			return err
		}
		script.before = scp
	}
	return nil
}

func ExecuteScript(exp js.Exp, param interface{}) (*js.Context, error) {

	if exp == nil {
		return nil, nil
	}
	vm := js.NewVm()
	vm.Set("$", param)
	vm.Set("log_err", log_error)
	vm.Set("log_info", log_Info)
	vm.SetFunc("base64len", base64Len)
	return vm, vm.SafeExecute(exp, nil)
}

var base64Len js.Func = func(i ...interface{}) interface{} {
	if len(i) > 0 {
		base := js.ConvertToString(i[0])
		return base64.StdEncoding.DecodedLen(len(base))
	}
	return 0
}

var log_error js.Func = func(i ...interface{}) interface{} {
	if len(i) > 0 && Logger != nil {
		Logger.Errorf(js.ConvertToString(i[0]), i[1:]...)
	}
	return nil
}

var log_Info js.Func = func(i ...interface{}) interface{} {
	if len(i) > 0 && Logger != nil {
		Logger.Infof(js.ConvertToString(i[0]), i[1:]...)
	}
	return nil
}


