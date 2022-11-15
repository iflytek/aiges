package jsonschema

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Schema struct {
	prop Validator // root validator
	i    interface{}
}

func NewSchema(i map[string]interface{}) (*Schema, error) {
	s := &Schema{}
	s.i = i
	p, err := NewProp(i, "$")
	if err != nil {
		return nil, err
	}
	s.prop = p
	return s, nil
}

func NewSchemaFromJSON(j []byte) (*Schema, error) {
	var i map[string]interface{}
	err := json.Unmarshal(j, &i)
	if err != nil {
		return nil, err
	}
	return NewSchema(i)
}
func (s *Schema) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	s.i = i
	p, err := NewProp(i, "$")
	if err != nil {
		return err
	}
	s.prop = p
	return nil
}

func (s *Schema) MarshalJSON() (b []byte, err error) {
	data, err := json.Marshal(s.i)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Schema) Validate(i interface{}) error {
	c := &ValidateCtx{}
	ii, err := scaleObject(i)
	if err != nil {
		return err
	}
	s.prop.Validate(c, ii)
	if len(c.errors) == 0 {
		return nil
	}
	return errors.New(errsToString(c.errors))
}

func (s *Schema) ValidateAndUnmarshalJSON(data []byte, template interface{}) (err error) {
	var i interface{}
	err = json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	err = s.Validate(i)
	if err != nil {
		return err
	}
	return UnmarshalFromMap(i, template)
}

func scaleObject(i interface{}) (o interface{}, err error) {
	switch d := i.(type) {
	case []byte:
		err = json.Unmarshal(d, &o)
		if err != nil {
			return o, err
		}
		return o, nil
	case string:
		err = json.Unmarshal([]byte(d), &o)
		if err != nil {
			return o, err
		}
		return o, nil
	default:
		return i, nil
	}
}

func (s *Schema) ValidateError(i interface{}) []Error {
	c := &ValidateCtx{}
	s.prop.Validate(c, i)
	return c.errors
}

func (s *Schema) Bytes() []byte {
	bs, _ := json.Marshal(s.i)
	return bs
}

func (s *Schema) FormatBytes() []byte {
	bf := bytes.NewBuffer(nil)
	bs := s.Bytes()
	err := json.Indent(bf, bs, "", "   ")
	if err != nil {
		return bs
	}
	return bf.Bytes()
}

func errsToString(errs []Error) string {
	sb := strings.Builder{}
	n := 0
	for _, err := range errs {
		n += len(err.Path) + len(err.Info) + 5
	}
	sb.Grow(n)
	for _, err := range errs {
		sb.WriteString(appendString("'", err.Path, "' ", err.Info, "; "))
	}
	return sb.String()
}

var (
	globalSchemas = map[reflect.Type]*Schema{}
)

//RegisterSchema  will generate schema by giving type and register it  to global map.
//use Validate() to validate the giving value
func RegisterSchema(typ interface{}) error {
	sc, err := GenerateSchema(typ)
	if err != nil {
		return err
	}
	globalSchemas[reflect.TypeOf(typ)] = sc
	return nil
}

func MustRegisterSchema(typ interface{}) {
	if err := RegisterSchema(typ); err != nil {
		panic("register schema error" + err.Error())
	}
}

func Validate(i interface{}) error {
	t := reflect.TypeOf(i)
	sc := globalSchemas[t]
	if sc == nil {
		return fmt.Errorf("no schema found for:%v", t.String())
	}
	return sc.Validate(i)
}
