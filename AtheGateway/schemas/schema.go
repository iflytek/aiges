package schemas

import (
	"errors"
	"github.com/qri-io/jsonschema"
	"github.com/qri-io/jsonpointer"
	"fmt"
	"regexp"
)

const (

	MAGIC_KEY = "magic"
	CONST_VALUE = "constVal"
	REPLACE_KEY =  "replaceKey"
	DEFAULT_KEY =  "defaultVal"

)

func init() {
	jsonschema.RegisterValidator("properties",NewProperties)
	jsonschema.RegisterValidator(MAGIC_KEY,NewMagic)
	jsonschema.RegisterValidator(CONST_VALUE, func() jsonschema.Validator {
		return new(SetVal)
	})
	jsonschema.RegisterValidator(REPLACE_KEY, func() jsonschema.Validator {
		return new(Replace)
	})
	jsonschema.RegisterValidator(DEFAULT_KEY, func() jsonschema.Validator {
		return new(Defalut)
	})
}


func Validate(name string, doc interface{}) (vales []jsonschema.ValError, err error) {
	mp := GetMappingByKey(name)
	if mp == nil {
		err = errors.New("cannot find mapping:" + name)
		return
	}
	var schema  = mp.Schema


	if schema == nil{
		return nil,nil
	}
	errs := []jsonschema.ValError{}
	schema.Validate("/",doc,&errs)
	//vales, err = sc.Schema.ValidateBytes(data)
	return errs, nil
}

//验证并返回message
func ValidateOfMsg(name string, doc interface{}) (string, error) {
	es, e := Validate(name, doc)
	if e != nil {
		return e.Error(), e
	}

	if len(es) == 0 {
		return "", nil
	}
	msg := ""
	for _, v := range es {
		msg += getLimitString(v.PropertyPath+" "+v.Message, 100)+","
	}

	return msg, errors.New("validate err ")
}

func ValidateByMapping(mp *RouteMapping,doc interface{}) (string, error)  {
	es, e := ValidateMp(mp, doc)
	if e != nil {
		return e.Error(), e
	}

	if len(es) == 0 {
		return "", nil
	}
	msg := ""
	for _, v := range es {
		msg += getLimitString(v.PropertyPath+" "+v.Message, 100)+","
	}

	return msg, errors.New("validate err ")
}

func ValidateMp(mp *RouteMapping, doc interface{}) (vales []jsonschema.ValError, err error) {
	//mp := GetMappingByKey(name)
	//if mp == nil {
	//	err = errors.New("cannot find mapping:" )
	//	return
	//}
	var schema  = mp.Schema
	if schema == nil{
		return nil,nil
	}
	errs := []jsonschema.ValError{}
	schema.Validate("/",doc,&errs)
	//vales, err = sc.Schema.ValidateBytes(data)
	return errs, nil
}

func getLimitString(s string, limit int) string {
	if len(s) > limit {
		return s[:limit]
	}
	return s
}


type Properties jsonschema.Properties

func NewProperties() jsonschema.Validator {

	return &Properties{}
//	return (*jsonschema.Properties)(unsafe.Pointer(p))
}

func (p Properties) Validate(propPath string, data interface{}, errs *[]jsonschema.ValError) {

	jp, err := jsonpointer.Parse(propPath)
	if err != nil {
		jsonschema.AddError(errs, propPath, nil, "invalid property path")
		return
	}

	if obj, ok := data.(map[string]interface{}); ok {
		for key, val := range obj {
			//限定参数集合
			if  p[key]==nil {
				jsonschema.AddError(errs, propPath, data, fmt.Sprintf(`invalid param:'%s'`,key))
				return
			}else{
				d, _ := jp.Descendant(key)
				p[key].Validate(d.String(), val, errs)
			}
		}
		//参数转换----------
		for k,v:=range p{
			magic:=v.JSONProp(MAGIC_KEY)
			if mac,ok:=magic.(*Magic);ok && mac!=nil{
				if mac.Key!="" && mac.Enable{
					GetOp(mac.Op).Do(obj,mac,k)
				}
			}
			//设置定值
			constVal:=v.JSONProp(CONST_VALUE)
			if val,ok:=constVal.(*SetVal);ok&&val !=nil{
					obj[k] = (*val)
			}
			//替换key
			replaceKey:=v.JSONProp(REPLACE_KEY)
			if replaceKey,ok:=replaceKey.(*Replace);ok{
				value:=obj[k]
				if value==nil || replaceKey==nil || obj[string(*replaceKey)]!=nil{

				}else{
					delete(obj,k)
					obj[string(*replaceKey)]=value;
				}

			}
			//设置默认值
			defaultKey:=v.JSONProp(DEFAULT_KEY)
			if def,ok:=defaultKey.(*Defalut);ok && def!=nil{
				if obj[k]==nil{
					obj[k] = *def
				}
			}
		}
	}

}

//参数转换的验证器

type Magic struct {
	Key string `json:"key"`
	Enable bool `json:"enable"`
	Op string `json:"op"`
	Val interface{} `json:"val"`
}

func NewMagic() jsonschema.Validator {
	return &Magic{}
}

func (m Magic) Validate(propPath string, data interface{}, errs *[]jsonschema.ValError) {

}

var opMap = map[string]Op{
	"set":&Set{},
	"rpl":&Rpl{},

}

var defaultOp  = &Rpl{}

func GetOp(k string)Op  {
	op:=opMap[k]
	if op ==nil{
		return defaultOp
	}
	return op
}

type Op interface {
	Do(map[string]interface{},*Magic,string)
}

type Rpl struct {}

func (r *Rpl)Do(obj map[string]interface{},mac *Magic,k string)  {
	value:=obj[k]
	if value==nil || obj[mac.Key]!=nil{
		return
	}
	delete(obj,k)
	obj[mac.Key]=value;
}

type Set struct {}

func (s *Set)Do(obj map[string]interface{},mac *Magic,k string)  {
	obj[mac.Key]=mac.Val;
}

type SetFun struct {}
var exprg = regexp.MustCompile(`\w+\(\w+\)`)
func (s *SetFun)Do(obj map[string]interface{},mac *Magic,k string)  {
	exp,ok:=mac.Val.(string)
	if !ok || !exprg.MatchString(exp){
		return
	}
	if obj["data_args"]==nil{
		obj["data_args"]= make(map[string]interface{})
	}
	obj[mac.Key]=mac.Val;

}

type Replace string


func (p *Replace) Validate(propPath string, data interface{}, errs *[]jsonschema.ValError){

}

type SetVal string

func (p *SetVal) Validate(propPath string, data interface{}, errs *[]jsonschema.ValError){

}

type SetArgs struct {
	Key string
	Val interface{}
}

func (p *SetArgs) Validate(propPath string, data interface{}, errs *[]jsonschema.ValError){

}

type Defalut string

func (p *Defalut) Validate(propPath string, data interface{}, errs *[]jsonschema.ValError){

}

