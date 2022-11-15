package jsonschema

import "fmt"

type ConstVal struct {
	Val interface{}
}

func (cc ConstVal) Validate(c *ValidateCtx,value interface{}) {

}

type DefaultVal struct {
	Val interface{}
}

func (d DefaultVal) Validate(c *ValidateCtx,value interface{}) {

}

type ReplaceKey string

func (r ReplaceKey) Validate(c *ValidateCtx,value interface{}) {

}

func NewConstVal(i interface{},path string ,parent Validator) (Validator, error) {
	return &ConstVal{
		Val: i,
	}, nil
}

func NewDefaultVal(i interface{},path string,parent Validator) (Validator, error) {
	return &DefaultVal{i}, nil
}

func NewReplaceKey(i interface{},path string,parent Validator) (Validator, error) {
	s, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("value of 'replaceKey' must be string :%v", i)
	}
	return ReplaceKey(s), nil

}

type FormatVal _type

func (f FormatVal) Validate(c *ValidateCtx, value interface{}) {

}

func (f FormatVal)Convert(value interface{})interface{}{
	switch _type(f) {
	case typeString:
		return StringOf(value)
	case typeBool:
		return BoolOf(value)
	case typeInteger,typeNumber:
		return NumberOf(value)
	}
	return value
}

func NewFormatVal(i interface{},path string,parent Validator)(Validator,error){
	str,ok:=i.(string)
	if !ok{
		return nil,fmt.Errorf("value of format must be string:%s",str)
	}
	return FormatVal(types[str]),nil
}

/*
{
	"setVal":{
		"key1":1,
		"key2":"val2",
		"key3":"${key1}",
		"key4":{
			"func":"append",
			"args":["${key1}","${key2}",{"func":"add","args":[1,2]}]
		},
	}
}
{
	"if":{
		"op":"eq",
		"l":"",
		"r":""
	}
	"then":{

	},

	"else":{

	},
	"and":[
		{
			"if":{}
		}
	],
	"set":{
		"k1":"",


	}
}

 */

type SetVal map[*JsonPathCompiled]Value

func (s SetVal) Validate(c *ValidateCtx,value interface{}) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		return
	}
	ctx:=Context(m)
	for key, val := range s {
		v:=val.Get(ctx)
		key.Set(m,v)
	}
}

func NewSetVal(i interface{},path string ,parent Validator)(Validator,error){
	m,ok:=i.(map[string]interface{})
	if !ok{
		return nil, fmt.Errorf("value of setVal must be map[string]interface{} :%v", i)
	}
	setVal:=SetVal{}
	for key, val := range m {
		v,err:=parseValue(val)
		if err != nil{
			return nil, err
		}
		jp,err:=parseJpathCompiled(key)
		if err != nil{
			return nil,err
		}
		setVal[jp] = v
	}
	return setVal, nil
}
