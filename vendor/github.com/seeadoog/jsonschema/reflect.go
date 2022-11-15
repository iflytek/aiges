package jsonschema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	_String     = "string"
	_Int        = "integer"
	_Bool       = "boolean"
	_Number     = "number"
	_Object     = "object"
	_Type       = "type"
	_Properties = "properties"
	_Array      = "array"
	_Items      = "items"
	_Enum       = "enum"
	_Maximum    = "maximum"
	_Minimum    = "minimum"
	_MaxLength  = "maxLength"
	_MinLength  = "minLength"
	_Required   = "required"
)

//generate jsonschema from giving template
func GenerateSchema(i interface{}) (*Schema, error) {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	schema := map[string]interface{}{}
	err := parseSchema(schema, t, nil)
	if err != nil {
		return nil, err
	}

	sc, err := NewSchema(schema)
	return sc, err
}

func GenerateSchemaAsString(i interface{}) (string, error) {
	schema, err := GenerateSchema(i)
	if err != nil {
		return "", err
	}
	bs, _ := json.Marshal(schema)
	return string(bs), nil
}

var (
	stringFieldsParses = []string{}
)

func AddRefString(validates ...string) {
	stringFieldsParses = append(stringFieldsParses, validates...)
}

func parseSchema(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) (err error) {
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		return parseSchema(sc, t, field)
	case reflect.Struct:
		properties := map[string]interface{}{}
		sc[_Properties] = properties
		sc[_Type] = _Object
		requires := make([]interface{}, 0)
		for i := 0; i < t.NumField(); i++ {
			fi := t.Field(i)
			if fi.Anonymous {
				fisc := map[string]interface{}{}
				err = parseSchema(fisc, fi.Type, &fi)
				if err != nil {
					return err
				}
				fp, _ := fisc[_Properties].(map[string]interface{})
				for key, val := range fp {
					properties[key] = val
				}
				continue
			}
			tag := fi.Tag.Get("json")
			if tag == "" {
				tag = fi.Name
			}
			fiv := map[string]interface{}{}
			properties[tag] = fiv
			required := fi.Tag.Get("required")
			if isTure(required) {
				requires = append(requires, tag)
			}
			if err := parseSchema(fiv, fi.Type, &fi); err != nil {
				return err
			}
		}
		if len(requires) > 0 {
			sc[_Required] = requires
		}
	case reflect.String:
		sc[_Type] = _String
		if field != nil {
			funs := []parseFunc{
				parseEnumString,
				parseMaxlength,
				parseMinlength,
				parseDefaultValue,
				parsePattern,
				parseFormat,
			}
			funs = append(funs, newParseFuncs(stringFieldsParses)...)
			err = doParses(funs, sc, t, field)
			if err != nil {
				return err
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sc[_Type] = _Int
		if field != nil {
			err = doParses([]parseFunc{
				parseEnumInt,
				parseMaximum,
				parseMinimum,
				parseDefaultValue,
				parseMultipleOf,
			}, sc, t, field)
			if err != nil {
				return err
			}
		}
	case reflect.Float32, reflect.Float64:
		sc[_Type] = _Number
		if field != nil {
			err = doParses([]parseFunc{
				parseEnumNumber,
				parseMaximum,
				parseMinimum,
				parseDefaultValue,
				parseMultipleOf,
			}, sc, t, field)
			if err != nil {
				return err
			}
		}
	case reflect.Bool:
		sc[_Type] = _Bool
		if field != nil {
			err = doParses([]parseFunc{
				parseDefaultValue,
			}, sc, t, field)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		sc[_Type] = _Array
		items := map[string]interface{}{}
		sc[_Items] = items
		err = parseSchema(items, t.Elem(), nil)
		if err != nil {
			return err
		}
		err = doParses([]parseFunc{
			parseMaxItems,
			parseMinItems,
			parseUniqueItems,
			parseDefaultValue,
		}, sc, t, field)
		if err != nil {
			return err
		}
	case reflect.Map:
		sc[_Type] = _Object
		if t.Elem().Kind() == reflect.Interface {
			sc["additionalProperties"] = true
		} else {
			addi := map[string]interface{}{}
			err = parseSchema(addi, t.Elem(), nil)
			if err != nil {
				return err
			}
			sc["additionalProperties"] = addi
			sc["properties"] = map[string]interface{}{}

		}

	default:
		return fmt.Errorf("unvalid type while parse schema:" + t.Name())
	}
	return nil
}

type parseFunc = func(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error

func parseMaximum(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	maximum := field.Tag.Get(_Maximum)
	if maximum != "" {
		num, err := strconv.ParseFloat(maximum, 64)
		if err != nil {
			return fmt.Errorf("parse int maximum tag error,value is not integer:%s:%s", field.Name, maximum)
		}
		sc[_Maximum] = num
	}
	return nil
}

func parseMinimum(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	minimum := field.Tag.Get(_Minimum)
	if minimum != "" {
		num, err := strconv.ParseFloat(minimum, 64)
		if err != nil {
			return fmt.Errorf("parse int minimum tag error,value is not integer:%s:%s", field.Name, minimum)
		}
		sc[_Minimum] = num
	}
	return nil
}

func doParses(funs []parseFunc, sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	for _, fun := range funs {
		if err := fun(sc, t, field); err != nil {
			return err
		}
	}
	return nil
}

func parseEnumString(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {

	enums := field.Tag.Get(_Enum)
	if len(enums) > 0 {
		eus := strings.Split(enums, ",")
		eusi := make([]interface{}, len(eus))
		for i, s := range eus {
			eusi[i] = s
		}
		sc[_Enum] = eusi
	}
	return nil
}

func parseEnumNumber(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	enums := field.Tag.Get(_Enum)
	if len(enums) > 0 {
		eus := strings.Split(enums, ",")
		eusi := make([]interface{}, len(eus))
		for i, s := range eus {
			num, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("parse int eumus tag error,tag value is not int:%s:%s", field.Name, enums)
			}
			eusi[i] = float64(num) // 主要是用做生成schema，做校验使用
		}
		sc[_Enum] = eusi
	}
	return nil
}

func parseEnumInt(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	enums := field.Tag.Get(_Enum)
	if len(enums) > 0 {
		eus := strings.Split(enums, ",")
		eusi := make([]interface{}, len(eus))
		for i, s := range eus {
			num, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return fmt.Errorf("parse int eumus tag error,tag value is not int:%s:%s", field.Name, enums)
			}
			eusi[i] = num // 主要是用做生成schema，做校验使用
		}
		sc[_Enum] = eusi
	}
	return nil
}

func parseMaxlength(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	maxLen := field.Tag.Get(_MaxLength)
	if len(maxLen) > 0 {
		num, err := strconv.Atoi(maxLen)
		if err != nil {
			return fmt.Errorf("parse maxLength tag error ,val is not int:%s:%s", field.Name, maxLen)
		}
		sc[_MaxLength] = float64(num)
	}
	return nil
}

func parseMinlength(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	minLength := field.Tag.Get(_MinLength)
	if len(minLength) > 0 {
		num, err := strconv.Atoi(minLength)
		if err != nil {
			return fmt.Errorf("parse minLength tag error ,val is not int:%s:%s", field.Name, minLength)
		}
		sc[_MinLength] = float64(num)
	}
	return nil
}

func parseDefaultValue(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	def := field.Tag.Get("default")
	if def != "" {
		val, err := formatValue(def, t)
		if err != nil {
			return fmt.Errorf("field '%v' default value should be type:%v but now value is :%v", field.Name, t, def)
		}
		sc["defaultVal"] = val
	}
	return nil
}

func formatValue(val string, t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.String:
		return val, nil
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		data, err := strconv.Atoi(val)
		return float64(data), err
	case reflect.Bool:
		return strconv.ParseBool(val)
	case reflect.Float32, reflect.Float64:
		return strconv.ParseFloat(val, 64)
	case reflect.Slice:
		return val, nil
	default:
		return nil, fmt.Errorf("%v type cannot have default value", t.String())
	}
}

func parsePattern(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	def := field.Tag.Get("pattern")
	if def != "" {
		sc["pattern"] = def
	}
	return nil
}

func parseFormat(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	def := field.Tag.Get("format")
	if def != "" {
		sc["format"] = def
	}
	return nil
}

func parseMultipleOf(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	def := field.Tag.Get("multipleOf")
	if def != "" {
		data, err := strconv.ParseFloat(def, 64)
		if err != nil {
			return fmt.Errorf("mulitpleOf val is not number : got:%v", def)
		}
		sc["multipleOf"] = data
	}
	return nil
}

func parseMaxItems(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	def := field.Tag.Get("maxItems")
	if def != "" {
		data, err := strconv.Atoi(def)
		if err != nil {
			return fmt.Errorf("maxItems val is not int : got:%v", def)
		}
		sc["maxItems"] = float64(data)
	}
	return nil
}

func parseMinItems(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	def := field.Tag.Get("minItems")
	if def != "" {
		data, err := strconv.Atoi(def)
		if err != nil {
			return fmt.Errorf("minItems val is not int : got:%v", def)
		}
		sc["minItems"] = float64(data)
	}
	return nil
}
func parseUniqueItems(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
	def := field.Tag.Get("uniqueItems")
	if def != "" {
		data, err := strconv.ParseBool(def)
		if err != nil {
			return fmt.Errorf("uniqueItems val is not bool : got:%v", def)
		}
		sc["uniqueItems"] = data
	}
	return nil
}

func newParseFuncs(ss []string) []parseFunc {
	funs := make([]parseFunc, 0, len(ss))

	for _, s := range ss {
		funs = append(funs, newParseFunc(s))
	}
	return funs
}

func newParseFunc(f string) parseFunc {
	return func(sc map[string]interface{}, t reflect.Type, field *reflect.StructField) error {
		def := field.Tag.Get(f)
		if def != "" {
			sc[f] = def
		}
		return nil
	}
}

func isTure(b string) bool {
	return b == "true" || b == "1" || b == "True" || b == "TRUE"
}
