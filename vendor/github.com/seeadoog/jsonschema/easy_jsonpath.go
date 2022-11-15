package jsonschema

import (
	"fmt"
	"regexp"
	"strconv"
)

func parseTokens(exp string) ([]string, error) {
	token := make([]byte, 0, len(exp))
	tokens := make([]string, 0, 1)
	skip := false
	for i := 0; i < len(exp); i++ {
		v := exp[i]
		if v == '\\' && !skip {
			skip = true
			continue
		}
		token = append(token, v)
		if v == '.' && !skip {
			tokens = append(tokens, string(token[:len(token)-1]))
			token = token[:0]
		}
		skip = false
	}
	if len(token) > 0 {
		tokens = append(tokens, string(token))
	}

	return tokens, nil
}

var reg = regexp.MustCompile(`(.*)(\[(\d+)\])?$`)

type token struct {
	key string
	index int
}

func parseJpathToken(tkn string) (*token, error) {
	if !reg.MatchString(tkn) {
		return nil, fmt.Errorf("invalid token:%s", tkn)
	}
	result := reg.FindAllStringSubmatch(tkn, -1)
	if len(result) == 0 || len(result[0]) < 4 {
		return nil, fmt.Errorf("invalid token:%s", tkn)
	}
	key := result[0][1]
	idxs := result[0][3]
	idx:=-1
	if idxs !=""{
		idx,_ =strconv.Atoi(idxs)
	}
	return &token{key: key,index: idx},nil
}

func parseAsTokens(exp string)([]*token,error){
	tokens,err:=parseTokens(exp)
	if err != nil{
		return nil, err
	}
	tkns:=make([]*token, len(tokens))
	for idx, item := range tokens {
		tkn,err:=parseJpathToken(item)
		if err != nil{
			return nil,err
		}
		tkns[idx] = tkn
	}
	return tkns,nil
}

type JsonPathCompiled struct {
	tokens []*token
}
func (c *JsonPathCompiled)Get(i interface{})(interface{},error){
	vi:=i
	for _, token := range c.tokens {
		if token.key != ""{
			m,ok:=vi.(map[string]interface{})
			if !ok{
				return nil,fmt.Errorf("try to get '%s' at not object value",token.key)
			}
			vi = m[token.key]
		}
		if token.index >=0{
			arr,ok:=vi.([]interface{})
			if !ok{
				return nil,fmt.Errorf("try to index '%d' at not array value",token.index)
			}
			if len(arr)<=token.index{
				return nil, fmt.Errorf("index out of range :%d",token.index)
			}
			vi = arr[token.index]
		}
	}
	return vi, nil
}
//key1.busi
func (c *JsonPathCompiled)Set(in interface{},val interface{})error{
	vi:=in
	vip:=in
	for i, token := range c.tokens {
		if i< len(c.tokens)-1{
			if token.key != ""{
				m,ok:= vi.(map[string]interface{})
				if !ok{
					return fmt.Errorf("try to set at not map val:in=%v",in)
				}
				vi = m[token.key]
				vip=m
			}
			if token.index >=0{
				arr,ok:=vi.([]interface{})
				if !ok{
					return fmt.Errorf("try to index '%d' at not array value",token.index)
				}
				if len(arr)<=token.index{
					return fmt.Errorf("index out of range :%d",token.index)
				}
				vi = arr[token.index]
				vip = arr
			}
		}else{
			if token.key != ""{
				m,ok:= vi.(map[string]interface{})
				if !ok{
					return fmt.Errorf("try to set at not map val:in=%v",in)
				}
				if token.index<0{
					m[token.key] = val
				}else{
					vi = m[token.key]
					vip = m
				}
			}
			if token.index>=0{
				arr,ok:=vi.([]interface{})
				if !ok{
					return fmt.Errorf("try to set index '%d' at not array value",token.index)
				}
				if len(arr) <= token.index{
					arr = append(arr,make([]interface{},token.index-len(arr)+1)...)
					m,ok:=vip.(map[string]interface{})
					if ok{
						m[token.key] = arr
					}
				}
				arr[token.index] = val
			}
		}
	}
	return nil
}

func parseJpathCompiled(exp string)(*JsonPathCompiled,error){
	tokens,err:=parseAsTokens(exp)
	if err != nil{
		return nil, err
	}
	return &JsonPathCompiled{tokens: tokens}, nil
}
//func CompileJpath(exp string)(JPath,error){
//
//}
