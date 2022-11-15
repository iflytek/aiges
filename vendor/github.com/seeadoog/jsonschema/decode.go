package jsonschema

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"unsafe"
)

//UnmarshalFromMap 将map 中的值序列化到 struct 中
func UnmarshalFromMap(in interface{}, template interface{}) error {
	v := reflect.ValueOf(template)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		panic("template value is nil or not pointer")
	}
	return unmarshalObject2Struct("", in, v)
}

var (
	bytesType = reflect.TypeOf([]byte(nil))
)

func unmarshalObject2Struct(path string, in interface{}, v reflect.Value) error {
	if in == nil {
		return nil
	}
	// 是非导出的变量
	if v.Kind() != reflect.Ptr && !v.CanSet() {

		return nil
	}

	switch {
	// 目标是字节数组
	case bytesType == v.Type():
		switch inv := in.(type) {
		case []byte:
			v.Set(reflect.ValueOf(in))
			return nil
		case string:
			bytes, err := base64.StdEncoding.DecodeString(inv)
			if err != nil {
				return fmt.Errorf("%s  type is not []byte , cannot decode as base64 string :%v", path, err)
			}
			v.Set(reflect.ValueOf(bytes))
			return nil
		default:
			return fmt.Errorf("%s type is not []byte", path)
		}

	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			vt := v.Type()
			elemType := vt.Elem()
			var nv reflect.Value
			switch elemType.Kind() {
			default:
				nv = reflect.New(elemType)
			}
			err := unmarshalObject2Struct(path, in, nv.Elem())
			if err != nil {
				return err
			}
			v.Set(nv)
			return nil
		}
		return unmarshalObject2Struct(path, in, v.Elem())
	case reflect.Slice:
		arr, ok := in.([]interface{})
		t := v.Type()
		if !ok {
			return fmt.Errorf("type of %s should be slice", path)
		}

		elemType := t.Elem()
		slice := reflect.MakeSlice(t, 0, len(arr))
		for _, v := range arr {
			elemVal := reflect.New(elemType)
			err := unmarshalObject2Struct(path, v, elemVal)
			if err != nil {
				return err
			}
			slice = reflect.Append(slice, elemVal.Elem())
		}
		v.Set(slice)
		return nil
	case reflect.String:
		vv, ok := in.(string)
		if !ok {
			return fmt.Errorf("type of %s should be string", path)
		}
		v.SetString(vv)
	case reflect.Map:
		vmap, ok := in.(map[string]interface{})
		if !ok {
			return fmt.Errorf("type of %s should be object", path)
		}
		t := v.Type()
		elemT := t.Elem()
		newV := v
		if v.IsNil() {
			newV = reflect.MakeMap(v.Type())
		}
		for key, val := range vmap {
			elemV := reflect.New(elemT)
			err := unmarshalObject2Struct(key, val, elemV)
			if err != nil {
				return err
			}
			newV.SetMapIndex(reflect.ValueOf(key), elemV.Elem())
		}
		v.Set(newV)
		return nil
	case reflect.Struct:
		t := v.Type()

		vmap, ok := in.(map[string]interface{})
		if !ok {
			return fmt.Errorf("type of %s should be object", path)
		}
		for i := 0; i < t.NumField(); i++ {
			fieldT := t.Field(i)
			name := fieldT.Tag.Get("json")
			if name == "" {
				name = fieldT.Name
			}
			if fieldT.Anonymous {
				err := unmarshalObject2Struct(name, in, v.Field(i))
				if err != nil {
					return err
				}
				continue
			}

			elemV := vmap[name]
			if elemV == nil {
				continue
			}
			// 是包进

			err := unmarshalObject2Struct(name, elemV, v.Field(i))
			if err != nil {
				return err
			}

		}
		return nil
	case reflect.Interface:
		inVal := reflect.ValueOf(in)
		if inVal.Type().Implements(v.Type()) {
			v.Set(inVal)
		}
		return nil
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		intV, err := intValueOf(in)
		if err != nil {
			return err
		}
		v.SetInt(intV)
		return nil
	case reflect.Bool:
		boolV, err := boolValueOf(in)
		if err != nil {
			return fmt.Errorf("%s error:%w", path, err)
		}
		v.SetBool(boolV)
		return nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		intV, err := intValueOf(in)
		if err != nil {
			return err
		}
		v.SetUint(uint64(intV))
		return nil
	case reflect.Float64, reflect.Float32:
		floatV, err := floatValueOf(in)
		if err != nil {
			return err
		}
		v.SetFloat(floatV)
		return nil
	case reflect.Array:
		arr, ok := in.([]interface{})
		//t := v.Type()
		if !ok {
			return fmt.Errorf("type of %s should be slice", path)
		}

		arType := reflect.ArrayOf(v.Len(), v.Type().Elem())
		arrv := reflect.New(arType)
		pointer := arrv.Pointer()
		eleSize := v.Type().Elem().Size()
		if v.Len() < len(arr) {
			return fmt.Errorf("length of %s is %d . but target value length is %d", path, v.Len(), len(arr))
		}
		for i, vv := range arr {
			elemV := reflect.New(v.Type().Elem())
			err := unmarshalObject2Struct(path, vv, elemV)
			if err != nil {
				return err
			}
			memCopy(pointer+uintptr(i)*eleSize, elemV.Pointer(), eleSize)
		}
		v.Set(arrv.Elem())
	default:
		panic("not support :" + v.Kind().String())
	}
	return nil
}

func intValueOf(v interface{}) (int64, error) {
	switch t := v.(type) {
	case float64:
		return int64(t), nil
	case float32:
		return int64(t), nil
	case int:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case int64:
		return int64(t), nil
	default:
		return 0, fmt.Errorf("type is %v ,not int ", reflect.TypeOf(v))
	}
}

func boolValueOf(v interface{}) (bool, error) {
	switch v := v.(type) {
	case bool:
		return v, nil
	case int:
		return v > 0, nil
	case float64:
		return v > 0, nil
	default:
		return false, fmt.Errorf("invalid bool value:%v", v)
	}
}

func floatValueOf(v interface{}) (float64, error) {
	switch v := v.(type) {
	case int:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("invalid float value:%v", v)
	}
}

func bytesOf(p uintptr, len uintptr) []byte {
	h := &reflect.SliceHeader{
		Data: p,
		Len:  int(len),
		Cap:  int(len),
	}
	return *(*[]byte)(unsafe.Pointer(h))
}

func memCopy(dst, src uintptr, len uintptr) {
	db := bytesOf(dst, len)
	sb := bytesOf(src, len)
	copy(db, sb)
}
