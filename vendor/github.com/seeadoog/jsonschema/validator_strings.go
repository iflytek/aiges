package jsonschema

import (
	"fmt"
	"regexp"
)

type Pattern struct {
	regexp *regexp.Regexp
	Path   string
	pattern string
}

func (p *Pattern) Validate(c *ValidateCtx, value interface{}) {
	str, ok := value.(string)
	if !ok {
		return
	}
	if !p.regexp.MatchString(str) {
		c.AddError(Error{
			Path: p.Path,
			Info: appendString( str," ,value does not match pattern: ",p.pattern),
		})
	}
}

func NewPattern(i interface{}, path string, parent Validator) (Validator, error) {
	str, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("%s is not a string when assign regexp,path:%s", str, path)
	}
	reg, err := regexp.Compile(str)
	if err != nil {
		return nil, fmt.Errorf("regexp compile error:%w", err)
	}
	return &Pattern{regexp: reg, Path: path,pattern: str}, nil
}

type FormatValidateFunc func(c *ValidateCtx, path string, value string)

var formats = map[string]FormatValidateFunc{
	"date-time":             wrapValidateFunc(isValidDateTime),
	"date":                  wrapValidateFunc(isValidDate),
	"email":                 wrapValidateFunc(isValidEmail),
	"hostname":              wrapValidateFunc(isValidHostname),
	"idn-email":             wrapValidateFunc(isValidIDNEmail),
	"idn-hostname":          wrapValidateFunc(isValidIDNHostname),
	"ipv6":                  wrapValidateFunc(isValidIPv6),
	"iri-reference":         wrapValidateFunc(isValidIriRef),
	"iri":                   wrapValidateFunc(isValidIri),
	"json-pointer":          wrapValidateFunc(isValidJSONPointer),
	"ipv4":                  wrapValidateFunc(isValidIPv4),
	"regex":                 wrapValidateFunc(isValidRegex),
	"relative-json-pointer": wrapValidateFunc(isValidRelJSONPointer),
	"time":                  wrapValidateFunc(isValidTime),
	"uri":                   wrapValidateFunc(isValidURI),
	"uri-reference":         wrapValidateFunc(isValidURIRef),
	"uri-template":          wrapValidateFunc(isValidURITemplate),
	"phone":                 wrapValidateFunc(isValidPhone),
}

func AddFormatValidateFunc(name string,f FormatValidateFunc){
	formats[name] = f
}

func wrapValidateFunc(fun func(value string) error) FormatValidateFunc {
	return func(c *ValidateCtx, path string, value string) {
		if err := fun(value); err != nil {
			c.AddError(Error{
				Path: path,
				Info: err.Error(),
			})
		}
	}
}

type Format struct {
	Path         string
	validateFunc FormatValidateFunc
}

func (f *Format) Validate(c *ValidateCtx, value interface{}) {
	str, ok := value.(string)
	if !ok {
		return
	}
	f.validateFunc(c, f.Path, str)
}

func NewFormat(i interface{}, path string, parent Validator) (Validator, error) {
	s, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("value of format must be string:%v,path:%s", desc(i), path)
	}
	vf, ok := formats[s]
	if !ok {
		return nil, fmt.Errorf("invalid format value:%v,path:%s", i, path)
	}
	return &Format{
		Path:         path,
		validateFunc: vf,
	}, nil
}
