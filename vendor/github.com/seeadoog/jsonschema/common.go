package jsonschema

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	sprintf = fmt.Sprintf
)

type Error struct {
	Path string
	Info string
}

type ValidateCtx struct {
	errors []Error
}

func (v *ValidateCtx) AddError(e Error) {
	v.errors = append(v.errors, e)
}

func (v *ValidateCtx) AddErrorInfo(path string, info string) {
	v.errors = append(v.errors, Error{Path: path, Info: info})
}

func (v *ValidateCtx) AddErrors(e ...Error) {
	for i, _ := range e {
		v.AddError(e[i])
	}
}

func (v *ValidateCtx) Clone() *ValidateCtx {
	return &ValidateCtx{}
}

type Validator interface {
	Validate(c *ValidateCtx, value interface{})
}

type NewValidatorFunc func(i interface{}, path string, parent Validator) (Validator, error)

func appendString(s ...string) string {
	sb := strings.Builder{}
	for _, str := range s {
		sb.WriteString(str)
	}
	return sb.String()
}

func panicf(f string, args ...interface{}) {
	panic(fmt.Sprintf(f, args...))
}

func StringOf(v interface{}) string {
	switch vv := v.(type) {
	case string:
		return vv
	case bool:
		if vv {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(vv, 'f', -1, 64)
	case int:
		return strconv.Itoa(vv)
	case nil:
		return ""

	}
	return fmt.Sprintf("%v", v)
}

func NumberOf(v interface{}) float64 {
	switch vv := v.(type) {
	case float64:
		return vv
	case bool:
		if vv {
			return 1
		}
		return 0
	case string:
		i, err := strconv.ParseFloat(vv, 64)
		if err != nil {
			return i
		}
		if vv == "true" {
			return 1
		}
		return 0
	}
	return 0
}

func BoolOf(v interface{}) bool {
	switch vv := v.(type) {
	case float64:
		return vv > 0
	case string:
		return vv == "true"
	case bool:
		return vv
	default:
		if NumberOf(v) > 0 {
			return true
		}
	}
	return false
}

func notNil(v interface{}) bool {
	switch v := v.(type) {
	case string:
		return v != ""
	case nil:
		return false

	}
	return true
}

func Equal(a, b interface{}) bool {
	return StringOf(a) == StringOf(b)
}

func desc(i interface{}) string {
	ty := reflect.TypeOf(i)
	return fmt.Sprintf("value:%v,type:%s", i, ty.String())
}
