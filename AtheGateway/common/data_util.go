package common

import (
	"bytes"
	"fmt"
	"strconv"
	"encoding/base64"
	"unsafe"
	"runtime"
)

func MapToString(m map[string]interface{})string  {
	if m==nil{
		return ""
	}
	b:=bytes.NewBuffer(make([]byte,0,64))
	for k,v:=range m{
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(ConvertToString(v))
		b.WriteString(",")
	}
	return b.String()
}


func ConvertToString(v interface{}) string {
	switch v.(type) {
	case string:return v.(string)
	case float64:return strconv.FormatFloat(v.(float64),'f',-1,64)
	case bool:return strconv.FormatBool(v.(bool))
	case int: return strconv.Itoa(v.(int))
	default:
		return fmt.Sprintf("%v",v)
	}
}


var enc = base64.StdEncoding
func EncodingTobase64String(src []byte)(string)  {
	buf := make([]byte,enc.EncodedLen(len(src)))
	enc.Encode(buf, src)
	return ToString(buf)
}

func ToString(buf []byte)string  {

	return *(*string)(unsafe.Pointer(&buf))
}

func ToBytes(s *string)[]byte  {
	return *(*[]byte)(unsafe.Pointer(s))
}


type StringBuildler struct {
	bf *bytes.Buffer
}

func NewStringBuilder() *StringBuildler  {
	return &StringBuildler{bf:bytes.NewBuffer(make([]byte,0,64))}
}

func (s *StringBuildler)Append(strs ...string) *StringBuildler {
	for _,v:=range strs{
		s.bf.WriteString(v)
	}
	return s
}

func (s *StringBuildler)ToString()string  {
	 return s.bf.String()
}

func (s *StringBuildler)GetBytes()[]byte  {
	return s.bf.Bytes()
}

func (s *StringBuildler)Len()int  {
	return s.bf.Len()
}

func RunWithRecovery( f func(),recoveryfun func())  {
	defer func() {
		if err :=recover();err !=nil{
			recoveryfun()
		}
	}()
	f()
}

func GetMaxConn(every int) int {
	if every <=0{
		every = 100
	}
	cpu:=runtime.NumCPU()
	return every *cpu
}

