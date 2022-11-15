package jsonschema

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type _type byte

const (
	typeString _type = iota + 1
	typeInteger
	typeNumber
	typeArray
	typeBool
	typeObject
)

var types = map[string]_type{
	"string":  typeString,
	"integer": typeInteger,
	"number":  typeNumber,
	"bool":    typeBool,
	"object":  typeObject,
	"boolean": typeBool,
	"array":   typeArray,
}

type typeValidateFunc func(path string, c *ValidateCtx, value interface{})

var typeFuncs = [...]typeValidateFunc{
	0: func(path string, c *ValidateCtx, value interface{}) {

	},
	typeString: func(path string, c *ValidateCtx, value interface{}) {
		switch value.(type) {
		case string, *string:
			return
		}
		if isKind(reflect.TypeOf(value), reflect.String) {
			return
		}
		c.AddError(Error{
			Path: path,
			Info: "Invalid type, expected: string , given: " + reflect.TypeOf(value).String(),
		})
	},
	typeObject: func(path string, c *ValidateCtx, value interface{}) {
		switch value.(type) {
		case map[string]interface{}, map[string]string:
			return
		default:
			ty := reflect.TypeOf(value)
			if isKind(ty, reflect.Struct, reflect.Map) {
				return
			}
			//if ty.Kind() == reflect.Ptr || ty.Kind() == reflect.Struct || ty.Kind() == reflect.Map {
			//	return
			//}
		}

		c.AddError(Error{
			Path: path,
			Info: "Invalid type, expected: object , given: " + reflect.TypeOf(value).String(),
		})
	},
	typeInteger: func(path string, c *ValidateCtx, value interface{}) {
		if _, ok := value.(float64); !ok {
			rt := reflect.TypeOf(value)
			if isKind(rt, reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint8,
				reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint) {
				return
			}
			//switch rt.Kind() {
			//case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint8,
			//reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			//	return
			//default:
			//}
			c.AddError(Error{
				Path: path,
				Info: "Invalid type, expected: integer , given: " + reflect.TypeOf(value).String(),
			})
		} else {
			v := value.(float64)
			if v != float64(int(v)) {
				c.AddError(Error{
					Path: path,
					Info: sprintf("type should be integer, but float:%v", v),
				})
			}
		}
	},

	typeNumber: func(path string, c *ValidateCtx, value interface{}) {
		if _, ok := value.(float64); !ok {
			rt := reflect.TypeOf(value)
			if isKind(rt, reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64,
				reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Float32, reflect.Float64) {
				return
			}
			//switch rt.Kind() {
			//case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64,
			//reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Float32, reflect.Float64:
			//	return
			//}
			c.AddError(Error{
				Path: path,
				Info: "Invalid type, expected: number , given: " + reflect.TypeOf(value).String(),
			})
		}
	},

	typeBool: func(path string, c *ValidateCtx, value interface{}) {
		switch value.(type) {
		case bool:
			return
		}
		if isKind(reflect.TypeOf(value), reflect.Bool) {
			return
		}
		c.AddError(Error{
			Path: path,
			Info: "Invalid type, expected: boolean , given: " + reflect.TypeOf(value).String(),
		})

	},

	typeArray: func(path string, c *ValidateCtx, value interface{}) {
		switch value.(type) {
		case []interface{}:
			return
		default:
			if isKind(reflect.TypeOf(value), reflect.Slice, reflect.Array) {
				return
			}
		}
		c.AddError(Error{
			Path: path,
			Info: "Invalid type, expected: array , given: " + reflect.TypeOf(value).String(),
		})

	},
}

func isKind(t reflect.Type, wants ...reflect.Kind) bool {
	k := t.Kind()
	if k == reflect.Ptr {
		return isKind(t.Elem(), wants...)
	}
	for _, want := range wants {
		if k == want {
			return true
		}
	}
	return false
}

type Type struct {
	Path         string
	ValidateFunc typeValidateFunc
}

func (t *Type) Validate(c *ValidateCtx, value interface{}) {

	//t.ValidateFunc(t.Path,c,value)
	if value == nil {
		return
	}
	t.ValidateFunc(t.Path, c, value)
}

func NewType(i interface{}, path string, parent Validator) (Validator, error) {
	iv, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("value of 'type' must be string! v:%v,path:%s", desc(i), path)
	}
	ivs := strings.Split(iv, "|")
	if len(ivs) > 1 {
		return NewTypes(iv, path, parent)
	}

	t, ok := types[iv]
	if !ok {
		return nil, fmt.Errorf("invalie type:%s,path:%s", iv, path)
	}

	return &Type{
		ValidateFunc: typeFuncs[t],
		Path:         path,
	}, nil
}

type Types struct {
	Vals []Validator
	Path string
	Type string
}

func (t *Types) Validate(c *ValidateCtx, value interface{}) {

	for _, v := range t.Vals {
		cc := c.Clone()
		v.Validate(cc, value)
		if len(cc.errors) == 0 {
			return
		}
	}
	c.AddErrors(Error{
		Path: t.Path,
		Info: appendString("type should be one of ", t.Type),
	})
}

func NewTypes(i interface{}, path string, parent Validator) (Validator, error) {
	str, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("value of types must be string !like 'string|number'")
	}
	arr := strings.Split(str, "|")
	tys := &Types{
		Vals: nil,
		Path: path,
		Type: str,
	}
	for _, s := range arr {
		//fmt.Println(s)
		ts, err := NewType(s, path, parent)
		if err != nil {
			return nil, fmt.Errorf("parse type items error!%w", err)
		}
		tys.Vals = append(tys.Vals, ts)
	}
	return tys, nil
}

type MaxLength struct {
	Val  int
	Path string
}

func (l *MaxLength) Validate(c *ValidateCtx, value interface{}) {

	switch value.(type) {
	case string:
		if len(value.(string)) > int(l.Val) {
			c.AddError(Error{
				Path: l.Path,
				Info: "length must be less or equal than " + strconv.Itoa(int(l.Val)),
			})
		}
	case []interface{}:
		if len(value.([]interface{})) > int(l.Val) {
			c.AddError(Error{
				Path: l.Path,
				Info: "length must be less or equal than " + strconv.Itoa(int(l.Val)),
			})
		}
	}

}

func NewMaxLen(i interface{}, path string, parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'maxLength' must be int: %v,path:%s", desc(i), path)
	}
	if v < 0 {
		return nil, fmt.Errorf("value of 'maxLength' must be >=0,%v path:%s", i, path)
	}
	return &MaxLength{
		Path: path,
		Val:  int(v),
	}, nil
}

func NewMinLen(i interface{}, path string, parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'minLengtg' must be int: %v,path:%s", desc(i), path)
	}
	if v < 0 {
		return nil, fmt.Errorf("value of 'minLength' must be >=0,%v path:%s", i, path)
	}
	return &MinLength{
		Val:  int(v),
		Path: path,
	}, nil
}

func NewMaximum(i interface{}, path string, parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'maximum' must be int")
	}
	return &Maximum{
		Val:  v,
		Path: path,
	}, nil
}

func NewMinimum(i interface{}, path string, parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'minimum' must be int:%v,path:%s", desc(i), path)
	}
	return &Minimum{
		Path: path,
		Val:  v,
	}, nil
}

type MinLength struct {
	Val  int
	Path string
}

func (l *MinLength) Validate(c *ValidateCtx, value interface{}) {
	switch value.(type) {
	case string:
		if len(value.(string)) < int(l.Val) {
			c.AddError(Error{
				Info: "length must be larger or equal than " + strconv.Itoa(int(l.Val)),
				Path: l.Path,
			})
		}
	case []interface{}:
		if len(value.([]interface{})) < int(l.Val) {
			c.AddError(Error{
				Info: "length must be larger or equal than " + strconv.Itoa(int(l.Val)),
				Path: l.Path,
			})
		}
	}
}

type Maximum struct {
	Val  float64
	Path string
}

func (m *Maximum) Validate(c *ValidateCtx, value interface{}) {
	val, ok := valueOfFloat(value)
	if !ok {
		return
	}
	if val > m.Val {
		c.AddError(Error{
			Info: appendString("value must be less or equal than ", strconv.FormatFloat(float64(m.Val), 'f', -1, 64)),
			Path: m.Path,
		})
	}
}

func valueOfFloat(value interface{}) (float64, bool) {
	val, ok := value.(float64)
	if ok {
		return val, true
	}
	return valueFloatByReflect(reflect.ValueOf(value))
}

func valueFloatByReflect(v reflect.Value) (float64, bool) {
	switch v.Kind() {
	case reflect.Ptr:
		return valueFloatByReflect(v.Elem())
	case reflect.Int, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return float64(v.Uint()), true
	case reflect.Float32:
		return float64(v.Float()), true
	}
	return 0, false
}

type Minimum struct {
	Val  float64
	Path string
}

func (m Minimum) Validate(c *ValidateCtx, value interface{}) {
	val, ok := valueOfFloat(value)
	if !ok {
		return
	}
	if val < (m.Val) {
		c.AddError(Error{
			Path: m.Path,
			Info: appendString("value must be larger or equal than ", strconv.FormatFloat(m.Val, 'f', -1, 64)),
		})
	}
}

type Enums struct {
	Val  []interface{}
	Path string
}

func (enums *Enums) Validate(c *ValidateCtx, value interface{}) {
	if value == nil {
		return
	}
	for _, e := range enums.Val {
		if e == value {
			return
		}
	}

	for _, e := range enums.Val {
		if Equal(e, value) {
			return
		}
	}
	c.AddError(Error{
		Path: enums.Path,
		Info: fmt.Sprintf("value is invalid , shoule be one of %v", enums.Val),
	})
}

func NewEnums(i interface{}, path string, parent Validator) (Validator, error) {
	arr, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value of 'enums' must be arr:%v,path:%s", desc(i), path)
	}
	return &Enums{
		Val:  arr,
		Path: path,
	}, nil
}

type Required struct {
	Val  []string
	Path string
	rMap map[string]bool
}

func (r *Required) Validate(c *ValidateCtx, value interface{}) {
	m, ok := value.(map[string]interface{})
	if !ok {
		r.validateStruct(c, reflect.ValueOf(value))
		return
	}
	for _, key := range r.Val {
		if _, ok := m[key]; !ok {
			c.AddError(Error{
				Path: appendString(r.Path, ".", key),
				Info: "field is required",
			})
		}
	}
}

func (r *Required) validateStruct(c *ValidateCtx, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		r.validateStruct(c, v.Elem())
		return
	case reflect.Struct:
		requiedMap := r.rMap
		for _, required := range r.Val {
			requiedMap[required] = true
		}
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			ft := t.Field(i)
			name := ft.Tag.Get("json")
			if name == "" {
				name = ft.Name
			}
			if !r.rMap[name] {
				continue
			}
			switch fv.Kind() {
			case reflect.Ptr:
				if fv.IsNil() {
					c.AddError(Error{
						Path: appendString(r.Path, ".", name),
						Info: "field is required",
					})
				}
			case reflect.String:
				if fv.String() == "" {
					c.AddError(Error{
						Path: appendString(r.Path, ".", name),
						Info: "field is required",
					})
				}
			}
		}
	}

}

func NewRequired(i interface{}, path string, parent Validator) (Validator, error) {
	arr, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value of 'required' must be array:%v", i)
	}
	var properties *Properties
	ap, ok := parent.(*ArrProp)
	if ok {
		pptis, ok := ap.Get("properties").(*Properties)
		if ok {
			properties = pptis
		}
	}
	req := make([]string, len(arr))
	for idx, item := range arr {
		itemStr, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("value of 'required item' must be string:%v of %v", item, i)
		}
		if properties != nil && !properties.EnableUnknownField {
			if _, ok := properties.properties[itemStr]; !ok {
				return nil, fmt.Errorf("required '%s' is not defined in propertis when additionalProperties is not enabled! path:%s", itemStr, path)
			}
		}

		req[idx] = itemStr

	}
	rm := make(map[string]bool)
	for _, re := range req {
		rm[re] = true
	}
	return &Required{
		Val:  req,
		Path: path,
		rMap: rm,
	}, nil
}

type Items struct {
	Val  *ArrProp
	Path string
}

func (item *Items) validateStruct(c *ValidateCtx, val interface{}) {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Slice:
		//t := v.Type()
		for i := 0; i < v.Len(); i++ {
			vi := v.Index(i)
			if vi.CanInterface() {
				item.Val.Validate(c, vi.Interface())
			}
		}
	}
}

func (i *Items) Validate(c *ValidateCtx, value interface{}) {
	if value == nil {
		return
	}
	arr, ok := value.([]interface{})
	if !ok {
		i.validateStruct(c, value)
		return
	}
	for _, item := range arr {
		for _, validator := range i.Val.Val {
			if validator.Val != nil {
				validator.Val.Validate(c, item)
			}
		}
	}
}

func NewItems(i interface{}, path string, parent Validator) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot create items with not object type: %v,path:%s", desc(i), path)
	}
	p, err := NewProp(m, path)
	if err != nil {
		return nil, err
	}
	p.(*ArrProp).Path = path + "[*]"
	return &Items{
		Val:  p.(*ArrProp),
		Path: path + "[*]",
	}, nil
}

type MultipleOf struct {
	Val  float64
	Path string
}

func (m MultipleOf) Validate(c *ValidateCtx, value interface{}) {
	v, ok := value.(float64)
	if !ok {
		return
	}
	a := (v / m.Val)

	if a != float64(int(a)) {
		c.AddError(Error{
			Path: m.Path,
			Info: sprintf("value must be multipleOf %v,but:%v, divide:%v", m.Val, v, v/m.Val),
		})
	}
}

func NewMultipleOf(i interface{}, path string, parent Validator) (Validator, error) {
	m, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf(" value of multipleOf must be an active number %v,path:%s", desc(i), path)
	}
	if m <= 0 {
		return nil, fmt.Errorf(" value of multipleOf must be an active number %v,path:%s", desc(i), path)
	}
	return &MultipleOf{Val: m, Path: path}, nil
}

// base64 解码后的长度校验器。以base64解码后的长度为准
type MaxB64DLength struct {
	Val  int
	Path string
}

func (l *MaxB64DLength) Validate(c *ValidateCtx, value interface{}) {

	switch value.(type) {
	case string:
		s := value.(string)
		n := base64.StdEncoding.DecodedLen(len(s))
		if n > int(l.Val) {
			c.AddError(Error{
				Path: l.Path,
				Info: "length is invalid, max length is  " + strconv.Itoa(int(l.Val)),
			})
		}
	}

}

func NewMaxB64DLen(i interface{}, path string, parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'maxB64DLen' must be int: %v,path:%s", desc(i), path)
	}
	if v < 0 {
		return nil, fmt.Errorf("value of 'maxB64DLen' must be >=0,%v path:%s", i, path)
	}
	return &MaxB64DLength{
		Path: path,
		Val:  int(v),
	}, nil
}

type MinB64DLength struct {
	Val  int
	Path string
}

func (l *MinB64DLength) Validate(c *ValidateCtx, value interface{}) {

	switch value.(type) {
	case string:
		s := value.(string)
		n := base64.StdEncoding.DecodedLen(len(s))
		if n < int(l.Val) {
			c.AddError(Error{
				Path: l.Path,
				Info: "length is invalid ,min length is  " + strconv.Itoa(int(l.Val)),
			})
		}
	}

}

func NewMinB64DLength(i interface{}, path string, parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'minB64DLen' must be int: %v,path:%s", desc(i), path)
	}
	if v < 0 {
		return nil, fmt.Errorf("value of 'minB64DLen' must be >=0,%v path:%s", i, path)
	}
	return &MinB64DLength{
		Path: path,
		Val:  int(v),
	}, nil
}

type constValidator struct {
	Path string
	V    string
}

func (c2 constValidator) Validate(c *ValidateCtx, value interface{}) {
	if StringOf(value) == c2.V {
		return
	}
	c.AddError(Error{
		Path: c2.Path,
		Info: "value is invalid , expected: " + c2.V,
	})
}

func NewConst(i interface{}, path string, parent Validator) (Validator, error) {
	return &constValidator{
		Path: path,
		V:    StringOf(i),
	}, nil
}
