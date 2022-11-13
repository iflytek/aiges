package common

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"unsafe"
)

func MapToString(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	b := bytes.NewBuffer(make([]byte, 0, 64))
	for k, v := range m {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(String(v))
		b.WriteString(",")
	}
	return b.String()
}
func MapstrToString(m map[string]string) string {
	if m == nil {
		return ""
	}
	b := bytes.NewBuffer(make([]byte, 0, 64))
	for k, v := range m {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(v)
		b.WriteString(",")
	}
	return b.String()
}

func String(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v.(bool))
	case int:
		return strconv.Itoa(v.(int))
	case int32:
		return strconv.Itoa(int(v.(int32)))
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

func Number(i interface{}) float64 {
	switch i.(type) {
	case float64:
		return i.(float64)
	case int:
		return float64(i.(int))
	case int32:
		return float64(i.(int32))
	case int64:
		return float64(i.(int64))
	case string:
		n, _ := strconv.Atoi(i.(string))
		return float64(n)
	}
	return 0
}

func Int(i interface{}) int {
	switch i.(type) {
	case float64:
		return int(i.(float64))
	case int:
		return i.(int)
	case int32:
		return int(i.(int32))
	case int64:
		return int(i.(int64))
	case string:
		n, _ := strconv.Atoi(i.(string))
		return n
	}
	return 0
}

func IntFromString(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func Bool(v interface{}) bool {
	switch v.(type) {
	case bool:
		return v.(bool)
	case string:
		return len(v.(string)) > 0
	case float64:
		return int(v.(float64)) > 0
	}
	if v != nil {
		return true
	}
	return false
}

var enc = base64.StdEncoding

func EncodingTobase64String(src []byte) string {
	buf := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)
	return ToString(buf)
}

func DecodeBase64string(s string) ([]byte, error) {
	dbuf := make([]byte, enc.DecodedLen(len(s)))
	n, err := enc.Decode(dbuf, *(*[]byte)(unsafe.Pointer(&s)))
	return dbuf[:n], err
}

func ToString(buf []byte) string {

	return *(*string)(unsafe.Pointer(&buf))
}

func ToBytes(s *string) []byte {
	return *(*[]byte)(unsafe.Pointer(s))
}

type StringBuildler struct {
	bf *bytes.Buffer
}

func NewStringBuilder() *StringBuildler {
	return &StringBuildler{bf: bytes.NewBuffer(make([]byte, 0, 64))}
}

func (s *StringBuildler) Append(strs ...string) *StringBuildler {
	for _, v := range strs {
		s.bf.WriteString(v)
	}
	return s
}

func (s *StringBuildler) AppendIfNotEmpty(k, v string) *StringBuildler {
	if len(v) > 0 {
		s.Append(k, "=", v, ",")
	}
	return s
}

func (s *StringBuildler) AppendIfNotEmptyI(k string, v interface{}) *StringBuildler {
	if v != nil {
		s.Append(k, "=", String(v), ",")
	}
	return s
}
func (s *StringBuildler) ToString() string {
	return ToString(s.bf.Bytes())
}

func (s *StringBuildler) GetBytes() []byte {
	return s.bf.Bytes()
}

func (s *StringBuildler) Len() int {
	return s.bf.Len()
}

func RunWithRecovery(f func(), recoveryFun func(err interface{})) {
	defer func() {
		if err := recover(); err != nil {
			if recoveryFun != nil {
				recoveryFun(err)
			}
		}
	}()
	f()
}

func GetMaxConn(every int) int {
	if every <= 0 {
		every = 100
	}
	cpu := runtime.NumCPU()
	return every * cpu
}

//stored by big end
func Int32toBytes(a int32) []byte {
	b := &reflect.SliceHeader{
		Data: (uintptr)(unsafe.Pointer(&a)),
		Cap:  4,
		Len:  4,
	}
	return *(*[]byte)(unsafe.Pointer(b))
}

//stored by bing end
func Int64toBytes(a int64) []byte {
	b := &reflect.SliceHeader{
		Data: (uintptr)(unsafe.Pointer(&a)),
		Cap:  8,
		Len:  8,
	}
	return *(*[]byte)(unsafe.Pointer(b))
}

//stored by big end
func Bytes2int32(b []byte) int32 {
	if len(b) < 4 {
		panic("[]byte convert to int32 failed,slice too short!")
	}
	return *(*int32)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&b)).Data))
}

//stored by big end
func Bytes2int64(b []byte) int64 {
	if len(b) < 8 {
		panic("[]byte convert to int64 failed,slice too short!")
	}
	return *(*int64)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&b)).Data))
}

type StringSlice []string

func (s *StringSlice) Get(i int) string {
	if i >= len(*s) {
		return ""
	}
	return (*s)[i]
}

func (s *StringSlice) Set(i int, val string) {
	if i >= len(*s) {
		old := *s
		*s = make([]string, i+5)
		copy(*s, old)
	}
	(*s)[i] = val
}

type Code int

var message = StringSlice{}

func (c Code) GetMessage() string {
	return message.Get(int(c))
}

func (c Code) GetInt() int {
	return int(c)
}

func WrapCode(code int, msg string) Code {
	message.Set(code, msg)
	return Code(code)
}
