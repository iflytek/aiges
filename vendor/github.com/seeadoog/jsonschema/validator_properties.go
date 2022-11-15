package jsonschema

import "fmt"


func init(){
	RegisterValidator("minProperties",NewMinProperties)
	RegisterValidator("oneOf",NewOneOf)
	AddIgnoreKeys("description")
	//AddIgnoreKeys("additionalProperties")
	AddIgnoreKeys("$schema")
	AddIgnoreKeys("$comment")
	AddIgnoreKeys("examples")
}

type MinProperties struct {
	Path string
	Value int
}

func (m *MinProperties) Validate(c *ValidateCtx, value interface{}) {
	switch value.(type) {
	case map[string]interface{}:
		if len(value.(map[string]interface{})) < m.Value{
			c.AddError(Error{
				Path: m.Path,
				Info: fmt.Sprintf("sub items number at least has %d",m.Value),
			})
		}
	case []interface{}:
		if len(value.([]interface{})) < m.Value{
			c.AddError(Error{
				Path: m.Path,
				Info: fmt.Sprintf("sub items number at least has %d",m.Value),
			})
		}
	}
}


func NewMinProperties(i interface{},path string,parent Validator)(Validator,error){
	fi,ok:=i.(float64)
	if !ok{
		return nil, fmt.Errorf("value of minProperties must be number:%v,path:%s",desc(i),path)
	}
	if fi <0 {
		return nil, fmt.Errorf("value of minProperties must be >0 :%v,path:%s",fi,path)
	}
	return &MinProperties{
		Path:  path,
		Value: int(fi),
	}, nil
}


type OneOf []Validator

func (a OneOf) Validate(c *ValidateCtx, value interface{}) {
	allErrs := []Error{}
	for _, validator := range a {
		cb := c.Clone()
		validator.Validate(cb, value)
		if len(cb.errors) == 0 {
			return
		}
		allErrs = append(allErrs, cb.errors...)
	}
	// todo 区分errors

	c.AddErrors(allErrs...)
}

func NewOneOf(i interface{}, path string, parent Validator) (Validator, error) {
	m, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value of oneOf must be array:%v,path:%s", desc(i), path)
	}
	any := OneOf{}
	for idx, v := range m {
		ip, err := NewProp(v, path)
		if err != nil {
			return nil, fmt.Errorf("oneOf index:%d is invalid:%w %v,path:%s", idx, err, v, path)
		}
		any = append(any, ip)
	}
	if len(any) ==0{
		return nil, fmt.Errorf("oneof length must be > 0,path:%s",path)
	}
	return any, nil
}