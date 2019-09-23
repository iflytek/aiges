package jsonscpt

import (
	"errors"
	"strconv"
	"strings"
	"fmt"
	"github.com/oliveagle/jsonpath"
	"sync"
	"regexp"
)

type QueryProp struct {
	Query string
	Value interface{}
}

type SwitchExp struct {
	SrcExp  string
	DataExp string
}

const (
	TYPE_KEY = -1
)

// src must be map[string]interface{}
// Marshal() do not support expression start with  $,$[0] .
//if you need this function ,please use MarshalInterface()


//resolve the token is a field or array
//return -1 and field name if token is field
//return the idx of token and array name if token is array
func yyp(token string) (string, int, error) {
	numidx_start := 0
	numidx_end := 0
	for k, v := range token {
		t := string(v)
		if t == "[" {
			numidx_start = k
		}
		if t == "]" {
			numidx_end = k
		}
	}
	if numidx_end > 0 && numidx_start >= 0 {
		num, err := strconv.Atoi(token[numidx_start+1 : numidx_end])
		if err != nil {
			return "", TYPE_KEY, err
		}
		return token[:numidx_start], num, nil
	}
	return token, TYPE_KEY, nil
}

//cahce tokens
var tokensCache sync.Map

func tokenize2(query string) ([]string, error) {
	tkn,ok:=tokensCache.Load(query)
	if !ok{
		tkns :=strings.Split(query, ".")
		tokensCache.Store(query,tkns)
		return tkns,nil
	}
	return tkn.([]string),nil
}

//marshal and set the value to interface{}
//MarshalInterface() is power than Marshal().
// MarshalInterface() support expression such as $ ,$[0] which  Marshal() doesn't support
//attention that $dst must be *interface{}
func MarshalInterface(query string, dst interface{}, value interface{}) error {
	return marshalInterface(query, dst, value)
}

func marshalInterface(query string, dst interface{}, value interface{}) error {
	//fmt.Println(reflect.TypeOf(dst))
	tks, err := compileTokens(query)

	if err != nil {
		return err
	}
	var cp = dst
	if cpi, ok := dst.(*interface{}); ok {
		if query == "" {
			*cpi = value
			return nil
		}
		if _, ok := (*cpi).(map[string]interface{}); ok {
			// do not handle
		} else if _, ok := (*cpi).([]interface{}); ok {
			cp = cpi
			goto done
		} else {
			//yp, idx, err := yyp(tks[0])
			yp:=tks[0].key
			idx:=tks[0].index
			//if err != nil {
			//	return err
			//}
			if idx == TYPE_KEY {
				*cpi = map[string]interface{}{}
			} else {
				if yp != "" {
					*cpi = map[string]interface{}{}
				} else {
					*cpi = make([]interface{}, idx+1)
					cp = cpi
					goto done
				}

			}
		}
		cp = *cpi

	}
done:
	return parserToken(tks, cp, value)

}

func parserToken(tks []tokens, cp, value interface{}) error {
	for k:=0;k< len(tks);k++ {
		field:=tks[k].key
		idx:=tks[k].index
		if idx == TYPE_KEY {
			cpm, ok := cp.(map[string]interface{})
			if !ok {
				cpms,ok:=cp.(map[string]string)
				if !ok{
					return errors.New(fmt.Sprintf("create field failed ,%s cannot convert_ to map", field))
				}
				//return errors.New(fmt.Sprintf("create field failed ,%s->parent cannot convert_ to map", field))
				cpms[field] = ConvertToString(value)
				return nil
			}

			if k < len(tks)-1 {
				if cpm[field] == nil {
					cpm[field] = map[string]interface{}{}
				}
				cpm2, ok := cpm[field].(map[string]interface{})
				if !ok {
					cpms,ok:=cpm[field].(map[string]string)
					if !ok{
						return errors.New(fmt.Sprintf("create field failed ,%s cannot convert_ to map", field))
					}
					cp = cpms
				}else{
					cp = cpm2
				}
			} else {
				//filed is last token ,set value to interface
				cpm[field] = value
			}
		} else { //array
			if field == "" && k == 0 {  //root array
				cpi, ok := cp.(*interface{})
				//	fmt.Println(reflect.TypeOf(cp))
				if !ok {
					return errors.New("root is not pointer")
				}
				if _, ok := (*cpi).([]interface{}); !ok {
					return errors.New("root is not array")
				}
				if len((*cpi).([]interface{})) < idx+1 {
					for i := len((*cpi).([]interface{})); i < idx+1; i++ {
						*cpi = append((*cpi).([]interface{}), nil)
					}
				}

				if k < len(tks)-1 {
					for i := 0; i < idx+1; i++ {
						if (*cpi).([]interface{})[i] == nil {
							(*cpi).([]interface{})[i] = map[string]interface{}{}
						}
					}
					cp = (*cpi).([]interface{})[idx]
				} else {
					//filed is last token ,set value to interface
					(*cpi).([]interface{})[idx] = value
				}
				//fmt.Println((*cpi).([]interface{}))
				continue
			}

			cpm, ok := cp.(map[string]interface{})
			if !ok {
				return errors.New("nil array child")
			}
			if cpm[field] == nil {
				cpm[field] = make([]interface{}, idx+1)
			}
			cps, ok := cpm[field].([]interface{})
			if !ok {
				return errors.New(fmt.Sprintf("create array failed ,%s cannot convert2 to array", field))
			}
			lenmap := len(cps)
			if lenmap < idx+1 {
				for i := lenmap; i < idx+1; i++ {
					cpm[field] = append(cpm[field].([]interface{}), nil)
				}
			}
			if k < len(tks)-1 {
				for i := 0; i < idx+1; i++ {
					if cpm[field].([]interface{})[i] == nil {
						cpm[field].([]interface{})[i] = map[string]interface{}{}
					}
				}
				cp = cpm[field].([]interface{})[idx]
			} else {
				//filed is last token ,set value to interface
				cpm[field].([]interface{})[idx] = value
			}
		}
	}

	return nil
}
/*
SwitchJson() can switch format of json from $data to $dst by specific expression strings.
attention that $dst must be type of *interface{}
*/
func SwitchJson(exps []SwitchExp, dst interface{}, data interface{}) error {
	for _, v := range exps {
		val, err := jsonpath.JsonPathLookup( data,v.DataExp)
		if err != nil {
			return err
		}
		err = marshalInterface(v.SrcExp, dst, val)
		if err != nil {
			return err
		}
	}
	return nil
}

var ruleRegexp= regexp.MustCompile(`(\w+|\$)(\[\d+\])?(\.\w+(\[\d+\])?)*$`)
func checkRule(rule string) bool {
	return ruleRegexp.MatchString(rule)
}


type tokens struct {
	index int
	key string
}

var cachedTokens = sync.Map{}

func compileTokens(s string)([]tokens,error){

	cachedTkns,ok:=cachedTokens.Load(s)
	if ok{
		return cachedTkns.([]tokens),nil
	}
	tkns,err:=tokeniz3(s)
	if err !=nil{
		return nil,err
	}

	var ctkns= make([]tokens,0, len(tkns))
	//fmt.Println("--------------------")
	for i:=0;i< len(tkns);i++{
		key,idx,err:=yyp(tkns[i])
		if err !=nil{
			return nil,err
		}
		ctkns = append(ctkns,tokens{index:idx,key:key})
	}
	cachedTokens.Store(s,ctkns)
	return ctkns,nil
}
//转义
func tokeniz3(s string)  ([]string, error) {
	tokens:=make([]string,0,5)
	token:=make([]byte,0, len(s))
	//fmt.Println("====================")
	for i:=0;i< len(s);i++{
		v:=s[i]
		token = append(token,v)

		if v=='.'{
			if len(token)>1{
				if token[len(token)-2]=='\\'{
					continue
				}
			}
			tokens = append(tokens,strings.Replace(string(token[:len(token)-1]),"\\","",-1))
			token=token[:0]
		}
	}
	if len(token)>0{
		tokens = append(tokens,strings.Replace(string(token),"\\","",-1))
	}
	return tokens,nil
}