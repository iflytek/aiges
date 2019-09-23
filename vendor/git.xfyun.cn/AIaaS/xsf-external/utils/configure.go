/*
* @file	configure.go
* @brief	configure parse
* @author	kunliu2
* @version	1.0
* @date		2017.11.14
 */
package utils

import (
	"errors"
	"fmt"
	"git.xfyun.cn/AIaaS/finder-go/common"
	"github.com/BurntSushi/toml"
	"reflect"
	"strings"
	//"fmt"
	"strconv"
)

const (
	ENOTFIND  = "can't find secet&key"
	EINVTYPE  = "invalied type"
	EINVCFG   = "invalied Configure"
	EINVCFGCT = "configure read invailed file context"
)
const (
	PEERADDR = "PeerAddr"
)

/*
CfgReader:
配置数据源读取数据
*/
type CfgReader interface {
	Read(name string) (string, error) //读取配置文件内容
}

/*
Configure:
配置操作实现
*/
type Configure struct {
	opt  *CfgOption //配置选项
	file string     //文件内容

	/*
		配置中心错误回调信息
	*/
	fileName string
	errCode  int
	errMsg   string

	r CfgReader
	c CfgMode

	v  map[string]interface{} //结构化配置类型
	dv map[string]interface{} //结构化配置类型
}

func NewCfg(c CfgMode, co *CfgOption) (*Configure, error) {
	/*var co CfgOption
	for _,opt := range o {
		opt(&co)
	}*/
	co.mode = c
	switch c {
	case Native:
		return NewNative(co)
	case Centre:
		if nil == co.fm {
			var e error
			co.fm, e = NewFinder(co)
			if nil != e {
				return nil, e
			}
		}
		return NewCentreWithFinder(co)
		//return nil,errors.New("NewCfg | Centre Not complete")
	case Custom:
		return nil, errors.New("NewCfg | Custom Not complete")
	default:
		return nil, errors.New("NewCfg | CfgMode Not Support")
	}
}
func NewCfgWithBytes(cfgData string) (*Configure, error) {
	return newBytesReader(cfgData)
}
func (c *Configure) GenNewCfg(co *CfgOption) (*Configure, error) {

	//todo： 如果其他不为空，且和原Configure不等，则报错
	if co.prj != c.opt.prj ||
		co.group != c.opt.group ||
		co.ver != c.opt.ver ||
		co.mode != c.opt.mode {
		return nil, errors.New("GenNewCfg | CfgMode Not Support")
	}

	if c.opt.mode == Centre {
		return NewCentreWithFinder(co)
	} else if c.opt.mode == Native {
		return nil, errors.New("NewCfg | CfgMode Not Support")
	}

	return nil, errors.New("NewCfg | CfgMode Not Support")
}

func (c *Configure) Option() *CfgOption {
	return c.opt
}

/*
Init:
初始化一个配置对象
*/
func (c *Configure) init(r CfgReader, co *CfgOption) error {
	var e error
	c.opt = co
	c.r = r
	c.file, e = r.Read(c.opt.name)
	if nil != e {
		return e
	}
	// todo: 合并default
	var v1 interface{}
	if len(c.opt.def) > 0 {
		_, e = toml.Decode(c.opt.def, &v1)
		switch v1.(type) {
		case map[string]interface{}:
			c.dv = v1.(map[string]interface{})
		default:
			e = nil
		}
	}

	var v interface{}
	_, e = toml.Decode(c.file, &v)
	if nil != e {
		return e
	}
	switch v.(type) {
	case map[string]interface{}:
		{
			c.v = v.(map[string]interface{})
		}
	default:
		{
			e = errors.New(EINVCFGCT)
			return e
		}
	}
	//合并默认默认配置
	return e
}

func (c *Configure) OnConfigFileChanged(con *finder.Config) bool {

	// 解析配置
	var v interface{}
	_, e := toml.Decode(string(con.File), &v)
	if nil != e {
		return false
	}
	c.file = string(con.File)
	switch v.(type) {
	case map[string]interface{}:
		c.v = v.(map[string]interface{})
	default:
		return false
	}
	// 调用用户自定义的cb
	if nil != c.opt.cb {
		return c.opt.cb(c)

	}
	return true
}
func (c *Configure) OnError(errInfo finder.ConfigErrInfo) {
	c.fileName = errInfo.FileName
	c.errCode = errInfo.ErrCode
	c.errMsg = errInfo.ErrMsg
}
func (c *Configure) GetSection(s string) interface{} {
	v, o := c.v[s]
	if o {
		return v

	}
	v, _ = c.dv[s]
	return v
	//return nil
}

func (c *Configure) SetAsObject(s string, k string, v interface{}) error {
	if nil != c.v {
		sv, e := c.v[s]
		//找到对应段落下的map
		switch tv := sv.(type) {
		case map[string]interface{}:
			_, e = tv[k]
			if true != e {
				return errors.New(ENOTFIND)
			}
			tv[k] = v
			return nil
		default:
			return errors.New(EINVTYPE)
		}
	}
	return fmt.Errorf("%v in read %v=>%v", EINVCFG, s, k)
}
func (c *Configure) GetString(s string, k string) (string, error) {
	sv, e := c.getString(c.v, s, k)
	if nil != e { //配置文件取不到，则去默认配置
		return c.getString(c.dv, s, k)
	}
	return sv, e
}

func (c *Configure) GetInt(s string, k string) (int, error) {
	sv, e := c.getInt(c.v, s, k)
	if nil != e { //配置文件取不到，则去默认配置
		return c.getInt(c.dv, s, k)
	}
	return sv, e
}

func (c *Configure) GetInt64(s string, k string) (int64, error) {
	sv, e := c.getInt64(c.v, s, k)
	if nil != e { //配置文件取不到，则去默认配置
		return c.getInt64(c.dv, s, k)
	}
	return sv, e
}
func (c *Configure) GetLocalIp() string {
	return c.opt.localIp
}
func (c *Configure) GetSvcIp() string {
	return c.opt.SvcIp
}
func (c *Configure) GetSvcPort() string {
	//如果没有传port，则用uuid替代
	var resPort string
	if 0 == c.opt.SvcPort {
		resPort = c.opt.uuid
	} else {
		resPort = strconv.Itoa(int(c.opt.SvcPort))
	}

	return resPort
}
func (c *Configure) GetBool(s string, k string) (bool, error) {
	sv, e := c.getBool(c.v, s, k)
	if nil != e { //配置文件取不到，则去默认配置
		return c.getBool(c.dv, s, k)
	}
	return sv, e
}
func (c *Configure) GetRawCfg() string {
	return c.file
}
func (c *Configure) GetAsObject(s string, k string) interface{} {
	o := c.getObject(c.v, s, k)
	if nil == o {
		return c.getObject(c.dv, s, k)
	}
	return o
}

func (c *Configure) GetSecs() []string {
	return c.getSecs(c.v)
}

func (c *Configure) GetTopSecs() map[string]interface{} {
	return c.getTopSecs(c.v)
}

func (c *Configure) GetInterface(s string, k string, v interface{}) error {
	e := c.getInterface(c.v, s, k, v)
	if nil != e {
		return c.getInterface(c.dv, s, k, v)
	}
	return e
}

/*
GetAsObject:
获取对应段以及关键字对应的object
*/

func (c *Configure) getInterface(mp map[string]interface{}, s string, k string, val interface{}) error {
	if nil != mp {

		v, o := mp[s]
		if !o {
			return fmt.Errorf("%v in read %v=>%v", EINVCFG, s, k)
		}
		val = v
	}
	return fmt.Errorf("%v in read %v=>%v", EINVCFG, s, k)
}

func (c *Configure) getObject(mp map[string]interface{}, s string, k string) interface{} {
	if nil != mp {
		sv, e := mp[s]
		if true == e {

			//找到对应段落下的map
			switch tv := sv.(type) {
			case map[string]interface{}:
				v, e := tv[k]
				if true != e {
					return nil
				}
				return v
			default:
				return nil
			}

		}
	}
	return nil
}

func (c *Configure) getTopSecs(mp map[string]interface{}) map[string]interface{} {
	rst := make(map[string]interface{})

	for k, v := range mp {
		if strings.Contains(reflect.TypeOf(v).String(), "map") {
			continue
		}
		rst[k] = v
	}

	return rst
}

func (c *Configure) getSecs(mp map[string]interface{}) []string {
	var rst []string

	for k, v := range mp {
		if !strings.Contains(reflect.TypeOf(v).String(), "map") {
			continue
		}
		rst = append(rst, k)
	}

	return rst
}

func (c *Configure) getBool(mp map[string]interface{}, s string, k string) (bool, error) {
	if nil != mp {
		//找到对应的段落
		sv, e := mp[s]
		if true != e {
			return false, errors.New(ENOTFIND)
		}

		//找到对应段落下的map
		var v interface{}
		switch tv := sv.(type) {
		case map[string]interface{}:
			v, e = tv[k]
			if true != e {
				return false, errors.New(ENOTFIND)
			}
		default:
			return false, errors.New(EINVTYPE)
		}

		switch tv := v.(type) {
		case bool:
			return tv, nil
		default:
			return false, errors.New(EINVTYPE)
		}
	}
	return false, fmt.Errorf("%v in read %v=>%v", EINVCFG, s, k)
}
func (c *Configure) getString(mp map[string]interface{}, s string, k string) (string, error) {
	if nil != mp {
		//找到对应的段落
		sv, e := mp[s]
		if true != e {
			return "", errors.New(ENOTFIND)
		}

		//找到对应段落下的map
		var v interface{}
		switch tv := sv.(type) {
		case map[string]interface{}:
			v, e = tv[k]
			if true != e {
				return "", errors.New(ENOTFIND)
			}
		default:
			return "", errors.New(EINVTYPE)
		}

		switch tv := v.(type) {
		case string:
			return tv, nil
		case int:
			return strconv.Itoa(tv), nil
		case int64:
			return strconv.FormatInt(tv, 10), nil
		case float64:
			return strconv.FormatFloat(v.(float64), 'E', -1, 64), nil
		case float32:
			return strconv.FormatFloat(v.(float64), 'E', -1, 32), nil
		default:
			return "", errors.New(EINVTYPE)
		}
	}
	return "", fmt.Errorf("%v in read %v=>%v", EINVCFG, s, k)
}

func (c *Configure) getInt(mp map[string]interface{}, s string, k string) (int, error) {

	if nil != mp {
		sv, e := mp[s]
		if true != e {
			return 0, errors.New(ENOTFIND)
		}

		//找到对应段落下的map
		var v interface{}
		switch tv := sv.(type) {
		case map[string]interface{}:
			v, e = tv[k]
			if true != e {
				return 0, errors.New(ENOTFIND)
			}
		default:
			return 0, errors.New(EINVTYPE)
		}

		switch tv := v.(type) {
		case int:
			return tv, nil
		case string:
			return strconv.Atoi(tv)
		case int64:
			return int(tv), nil
		case int32:
			return int(tv), nil
		default:
			return 0, errors.New(EINVTYPE)
		}
	}
	return 0, fmt.Errorf("%v in read %v=>%v", EINVCFG, s, k)
}

func (c *Configure) getInt64(mp map[string]interface{}, s string, k string) (int64, error) {
	i, e := c.getInt(mp, s, k)
	if nil != e {
		return 0, fmt.Errorf("%v in read %v=>%v", EINVCFG, s, k)
	}
	return int64(i), nil
}

// formatAtom formats a value without inspecting its internal structure.
/*
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
		// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path,
				formatAtom(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem())
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

func Dispaly(name string,x interface{}){
	fmt.Println("Display:%s (%T)", name, x)
	display(name, reflect.ValueOf(x))
}*/
