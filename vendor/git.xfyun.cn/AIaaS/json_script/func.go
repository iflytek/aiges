package jsonscpt

import (
	"bytes"
	"strconv"
	"fmt"
	"strings"
	"encoding/json"
)

type Func func(...interface{})interface{}
//get len of string
func lens(i ...interface{}) interface{} {
	if len(i)==0{
		return 0
	}
	v:=i[0]
	switch v.(type) {
	case string:
		return len(v.(string))
	case map[string]interface{}:
		return len(v.(map[string]interface{}))
	case []interface{}:
		return len(v.([]interface{}))
	case []string:
		return len(v.([]string))

	default:
		return 0
	}
}
// append string and number
func apd(i ...interface{}) interface{} {
	var bf bytes.Buffer
	for _, v := range i {
		bf.WriteString(ConvertToString(v))
	}
	return bf.String()
}
// split str
var split Func = func(i ...interface{}) interface{} {
	if len(i)>=2{
		if s,ok:=i[0].(string);ok{
			if sp,ok:=i[1].(string);ok{
				if len(i)>=3{
					if n,ok:=i[2].(float64);ok{
						return strings.SplitN(s,sp,int(n))
					}
				}else{
					return strings.Split(s,sp)
				}
			}
		}
	}
	return nil
}

var add Func = func(i ...interface{}) interface{} {
	var sum  float64 = 0
	for _, v := range i {
		sum+=number(v)
	}
	return sum
}

var printf Func = func(i ...interface{}) interface{} {
	if len(i)>0{
		if format,ok:=i[0].(string);ok{
			fmt.Printf(format+"\n",i[1:]...)
		}
	}
	return nil
}

var printlnn Func = func(i ...interface{}) interface{} {
		fmt.Println(i...)
	return nil
}

var sprintf Func = func(i ...interface{}) interface{} {
	if len(i)>0{
		if format,ok:=i[0].(string);ok{
			return fmt.Sprintf(format,i[1:]...)
		}
	}
	return ""
}
//to json string
var jsonMarshal Func = func(i ...interface{}) interface{} {
	if len(i)>0{
		b,err:=json.Marshal(i[0])
		if err !=nil{
			return "{}"
		}
		return toString(b)

	}
	return ""
}
//test if value is nil
var isNil Func = func(i ...interface{}) interface{} {
	if len(i)>0{
		if i[0]==nil{
			return true
		}else{
			return false
		}
	}
	return true
}
//delete key of a object
var deleteFun Func = func(i ...interface{}) interface{} {
	if len(i)>0{
		if m,ok:=i[0].(map[string]interface{});ok{
			for j:=1;j< len(i);j++{
				if k,ok:=i[j].(string);ok{
					delete(m,k)
				}
			}
		}
	}
	return nil
}
// &&
var and Func = func(i ...interface{}) interface{} {
	for _, v := range i {
		if !convertToBool(v){
				return false
		}
	}
	return true
}
//||
var or Func = func(i ...interface{}) interface{} {
	for _, v := range i {
		if convertToBool(v){
			return true
		}
	}
	return false
}
//==
var eq Func = func(i ...interface{}) interface{} {
	if len(i)<2{
		return false
	}
	return fmt.Sprintf("%v",i[0])==fmt.Sprintf("%v",i[1])
}

var not Func = func(i ...interface{}) interface{} {
	if len(i)<1{
		return false
	}
	return !convertToBool(i[0])
}
// >
var gt Func = func(i ...interface{}) interface{} {
	if len(i)>=2{
		//fmt.Println(number(i[0]),number(i[1]),number(i[0])>number(i[1]))
		return number(i[0])>number(i[1])
	}
return false
}
// >=
var ge Func = func(i ...interface{}) interface{} {
	if len(i)>=2{
		return number(i[0])>=number(i[1])
	}
	return false
}
// <
var lt Func = func(i ...interface{}) interface{} {
	if len(i)>=2{
		return number(i[0])<number(i[1])
	}
	return false
}
//<=
var le Func = func(i ...interface{}) interface{} {
	if len(i)>=2{
		return number(i[0])<=number(i[1])
	}
	return false
}

var exit Func = func(i ...interface{}) interface{} {
	if len(i)>=3{
		return &ErrorExit{Code: int(number(i[0])),Message:ConvertToString(i[1]),Value:i[2]}
	}
	if len(i)>=2{
		return &ErrorExit{Code: int(number(i[0])),Message:ConvertToString(i[1])}
	}
	if len(i)>=1{
		return &ErrorExit{Code: int(number(i[0]))}
	}

	return &ErrorExit{}
}
//return
var ret Func = func(i ...interface{}) interface{} {

	if len(i)>=1{
		return &ErrorReturn{Value:i[0]}
	}

	return &ErrorReturn{}
}

var join Func = func(i ...interface{}) interface{} {
	if len(i)>1{
		var sp = ConvertToString(i[len(i)-1])
		var joined = make([]string,0, len(i))
		for k:=0;k< len(i)-1;k++{
			switch i[k].(type) {
			case []string:
				joined = append(joined,i[k].([]string)...)
			default:
				joined = append(joined,ConvertToString(i[k]))
			}
		}
		return strings.Join(joined,sp)
	}
	return ""
}
var contains Func = func(i ...interface{}) interface{} {
	if len(i)<2{
		return false
	}
	return strings.Contains(ConvertToString(i[0]),ConvertToString(i[1]))
}

var in Func = func(i ...interface{}) interface{} {
	if len(i)>=1{

		t:=ConvertToString(i[0])
		for k:=1;k< len(i);k++{
			switch i[k].(type) {
			case string:
				if t == i[k].(string){
					return true
				}
			case []interface{}:
				for _,vv:=range i[k].([]interface{}){
					if t == ConvertToString(vv){
						return true
					}
				}
			case []string:
				for _,vv:=range i[k].([]string){
					if t == vv{
						return true
					}
				}
			default:
				if t==ConvertToString(i[k]){
					return true
				}
			}


		}
	}
	return false
}

var index Func = func(i ...interface{}) interface{} {
	if len(i)<2{
		return nil
	}

	switch i[0].(type) {
	case []string:
		return i[0].([]string)[int(number(i[1]))]
	case []int:
		return i[0].([]int)[int(number(i[1]))]
	case []interface{}:
		return i[0].([]interface{})[int(number(i[1]))]
	}
	return nil
}

var set Func = func(i ...interface{}) interface{} {
	if len(i)>2{
		switch i[0].(type) {
		case map[string]interface{}:
			i[0].(map[string]interface{})[ConvertToString(i[1])]=i[2]
		case map[string]string:
			i[0].(map[string]string)[ConvertToString(i[1])]=ConvertToString(i[2])
		}
	}
	return nil
}
//get the value of map or array
var get Func  = func(i ...interface{}) interface{} {
	if len(i)>1{
		key :=ConvertToString(i[1])
		switch i[0].(type) {
		case map[string]interface{}:
			return i[0].(map[string]interface{})[key]
		case map[string]string:
			return i[0].(map[string]string)[key]
		case []interface{}:
			return i[0].([]interface{})[int(number(key))]
		case []string:
			return i[0].([]string)[int(number(key))]
		case []int:
			return i[0].([]int)[int(number(i[1]))]

		}
	}
	return nil
}

var input Func = func(i ...interface{}) interface{} {
	var s string
	fmt.Scanf("%s",&s)
	return s
}


var trim Func = func(i ...interface{}) interface{} {
	if len(i)>=3{
		if number(i[2])>=0{
			return strings.TrimSuffix(ConvertToString(i[0]),ConvertToString(i[1]))
		}

		if number(i[2])<0{
			return strings.TrimPrefix(ConvertToString(i[0]),ConvertToString(i[1]))
		}
	}
	if len(i)>=2{
		return strings.TrimPrefix(ConvertToString(i[0]),ConvertToString(i[1]))
	}
	if len(i)==1{
		return i[0]
	}
	return ""
}

func number(i interface{})  float64{
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
		 n,_:=strconv.Atoi(i.(string))
		 return float64(n)
	}
	return 0
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
