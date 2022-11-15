package jsonschema

import (
	"fmt"
)

const (
	keyCase    = "case"
	keyDefault = "default"
)

type AnyOf []Validator

func (a AnyOf) Validate(c *ValidateCtx, value interface{}) {
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

func NewAnyOf(i interface{}, path string, parent Validator) (Validator, error) {
	m, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value of anyOf must be array:%v,path:%s", desc(i), path)
	}
	any := AnyOf{}
	for idx, v := range m {
		ip, err := NewProp(v, path)
		if err != nil {
			return nil, fmt.Errorf("anyOf index:%d is invalid:%w %v,path:%s", idx, err, v, path)
		}
		any = append(any, ip)
	}
	return any, nil
}

type If struct {
	Then *Then
	Else *Else
	v    Validator
}

func (i *If) Validate(c *ValidateCtx, value interface{}) {
	cif := c.Clone()
	i.v.Validate(cif, value)
	if len(cif.errors) == 0 {
		if i.Then != nil {
			i.Then.v.Validate(c, value)
		}
	} else {
		if i.Else != nil {
			i.Else.v.Validate(c, value)
		}
	}
}

func NewIf(i interface{}, path string, parent Validator) (Validator, error) {
	ifp, err := NewProp(i, path)
	if err != nil {
		return nil, err
	}
	iff := &If{
		v: ifp,
	}
	pp, ok := parent.(*ArrProp)
	if ok {
		then, ok := pp.Get("then").(*Then)
		if ok {
			iff.Then = then
		}
		elsef, ok := pp.Get("else").(*Else)
		if ok {
			iff.Else = elsef
		}
	}
	return iff, nil
}

type Then struct {
	v Validator
}

func (t *Then) Validate(c *ValidateCtx, value interface{}) {
	// then 不能主动调用
}

type Else struct {
	v Validator
}

func (e *Else) Validate(c *ValidateCtx, value interface{}) {
	//panic("implement me")
}

func NewThen(i interface{}, path string, parent Validator) (Validator, error) {
	v, err := NewProp(i, path)
	if err != nil {
		return nil, err
	}
	return &Then{
		v: v,
	}, nil
}

func NewElse(i interface{}, path string, parent Validator) (Validator, error) {
	v, err := NewProp(i, path)
	if err != nil {
		return nil, err
	}
	return &Else{
		v: v,
	}, nil
}

type Not struct {
	v    Validator
	Path string
}

func (n Not) Validate(c *ValidateCtx, value interface{}) {
	cn := c.Clone()
	n.v.Validate(cn, value)
	//fmt.Println(ners,value)
	if len(cn.errors) == 0 {
		c.AddErrors(Error{
			Path: n.Path,
			Info: "is not valid",
		})
	}
}

func NewNot(i interface{}, path string, parent Validator) (Validator, error) {
	p, err := NewProp(i, path)
	if err != nil {
		return nil, err
	}
	return Not{v: p}, nil
}

type AllOf []Validator

func (a AllOf) Validate(c *ValidateCtx, value interface{}) {
	for _, validator := range a {
		validator.Validate(c, value)
	}
}

func NewAllOf(i interface{}, path string, parent Validator) (Validator, error) {
	arr, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value of 'allOf' must be array: %v", desc(i))
	}
	all := AllOf{}
	for _, ai := range arr {
		iv, err := NewProp(ai, path)
		if err != nil {
			return nil, err
		}
		all = append(all, iv)
	}
	return all, nil
}

type Dependencies struct {
	Val  map[string][]string
	Path string
}

func (d *Dependencies) Validate(c *ValidateCtx, value interface{}) {
	m, ok := value.(map[string]interface{})
	if !ok {
		return
	}
	// 如果存在key，那么必须存在某些key
	for key, vals := range d.Val {
		_, ok := m[key]
		if ok {
			for _, val := range vals {
				_, ok = m[val]
				if !ok {
					c.AddErrors(Error{
						Path: appendString(d.Path, ".", val),
						Info: "is required",
					})
				}
			}
		}
	}
}

func NewDependencies(i interface{}, path string, parent Validator) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value of dependencies must be map[string][]string :%v", desc(i))
	}
	vad := &Dependencies{
		Val:  map[string][]string{},
		Path: path,
	}
	for key, arris := range m {
		arrs, ok := arris.([]interface{})
		if !ok {
			return nil, fmt.Errorf("value of dependencies must be map[string][]string :%v,path:%s", desc(i), path)
		}
		strs := make([]string, len(arrs))
		for idx, item := range arrs {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("value of dependencies must be map[string][]string :%v,path:%s", desc(i), path)

			}
			strs[idx] = str
		}
		vad.Val[key] = strs

	}
	return vad, nil
}

/*
{
	"keyMatch":{
		"key1":"biaoge"
	}
}
*/

type KeyMatch struct {
	Val  map[string]interface{}
	Path string
}

func (k *KeyMatch) Validate(c *ValidateCtx, value interface{}) {
	m, ok := value.(map[string]interface{})
	if !ok {
		//c.AddError(Error{
		//	Path: k.Path,
		//	Info: "value is not object",
		//})
		return
	}
	for key, want := range k.Val {
		target := m[key]
		if target != want {
			c.AddError(Error{
				Path: appendString(k.Path, ".", key),
				Info: fmt.Sprintf("value must be %v", want),
			})
		}
	}
}

func NewKeyMatch(i interface{}, path string, parent Validator) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value of keyMatch must be map[string]interface{} :%v", desc(i))
	}
	return &KeyMatch{
		Val:  m,
		Path: path,
	}, nil
}

/*
{
	"switch":"tsy",
	"cases":{
		"key1":{},
		"key2":{}
	},
	"default":{}
}
*/
type Switch struct {
	Switch  string
	Case    map[string]Validator
	Default *Default
}

func (s *Switch) Validate(c *ValidateCtx, value interface{}) {
	m, ok := value.(map[string]interface{})
	if !ok {
		if s.Default != nil {
			s.Default.p.Validate(c, value)
		}
		return
	}
	for cas, validator := range s.Case {
		if cas == StringOf(m[s.Switch]) {
			validator.Validate(c, value)
			return
		}
	}
	if s.Default != nil {
		s.Default.p.Validate(c, value)
	}
}

func NewSwitch(i interface{}, path string, parent Validator) (Validator, error) {
	key, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("value of switch must be string path:%s", path)
	}

	s := &Switch{
		Switch: key,
		Case:   map[string]Validator{},
	}
	ap, ok := parent.(*ArrProp)
	if !ok {
		return s, nil
	}
	cases, ok := ap.Get(keyCase).(Cases)
	if ok {
		s.Case = cases
	}
	def, ok := ap.Get(keyDefault).(*Default)
	if ok {
		s.Default = def
	}
	return s, nil
}

type Default struct {
	p *ArrProp
}

func (d Default) Validate(c *ValidateCtx, value interface{}) {
}

func NewDefault(i interface{}, path string, parent Validator) (Validator, error) {
	da, err := NewProp(i, path)
	if err != nil {
		return nil, err
	}
	return &Default{p: da.(*ArrProp)}, nil
}

type Cases map[string]Validator

func (c2 Cases) Validate(c *ValidateCtx, value interface{}) {

}

func NewCases(i interface{}, path string, parent Validator) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value of case must be map,path: %s", path)
	}
	cases := make(Cases)
	for key, val := range m {
		vad, err := NewProp(val, path)
		if err != nil {
			return nil, err
		}
		cases[key] = vad
	}
	return cases, nil
}
