package jsonref

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type QueryProp struct {
	Query string
	Value interface{}
}

const (
	TYPE_KEY=-1
)

// src must be map[string]interface{}
func Marshal(query string,src interface{}, value interface{})error{
	return marshal(query,src,value,0)
}

func marshal(query string,src interface{}, value interface{},start int) error {
	tks,err:=tokenize2(query)
	if err !=nil{
		return err
	}
	var cp = src
	for k,v:=range tks{
		field,idx,err:=yyp(v)
		if err!=nil{
			return err
		}
		if k < len(tks)-1{
			//map
			if idx ==TYPE_KEY{
				cpm ,ok:= cp.(map[string]interface{})
				if !ok{
					return errors.New(fmt.Sprintf("create field failed ,%s.parent cannot convert_ to map",field))
				}
				if cpm[field]==nil{
					cpm[field]= map[string]interface{}{}
				}
				cpm,ok = cpm[field].(map[string]interface{})
				if !ok{
					return errors.New(fmt.Sprintf("create field failed ,%s cannot convert_ to map",field))
				}
				cp = cpm
			}else{   //array
				cpm,ok := cp.(map[string]interface{})
				if !ok{
					return errors.New("nil array child")
				}
				if cpm[field]==nil{
					cpm[field]=make([]interface{},idx+1)
				}
				//arrmp:=cp[field].([]map[string]interface{})
				//log.Println(reflect.TypeOf(cpm[field]))
				cps,ok:=cpm[field].([]interface{})
				if !ok{
					return errors.New(fmt.Sprintf("create array failed ,%s cannot convert2 to array",field))
				}
				lenmap:=len(cps)
				if lenmap<idx+1{
					for i:=lenmap;i<idx+1;i++{
						cpm[field]=append(cpm[field].([]interface{}),nil)
					}
				}
				for i:=0;i<idx+1;i++{
					if cpm[field].([]interface{})[i]==nil{
						cpm[field].([]interface{})[i]=map[string]interface{}{}
					}
				}
				cp = cpm[field].([]interface{})[idx]
			}

		}else{
			if idx ==TYPE_KEY{
				cpm,ok:=cp.(map[string]interface{})
				if !ok{
					return errors.New(fmt.Sprintf("set val failed,%s is not map",field))
				}
				cpm[field]= value
			}else{
				if field==""{
					return errors.New("cannot set val for array for empty name")
				}
				cpm,ok:=cp.(map[string]interface{})
				if !ok{
					return errors.New(fmt.Sprintf("set val failed ,%s cannot convert1 to map -",field))
				}
				if cpm[field]==nil{
					cpm[field]=make([]interface{},idx+1)
				}
				cps,ok:=cpm[field].([]interface{})
				if !ok{
					return errors.New(fmt.Sprintf("set val failed,%s cannot convert3 to array",field))
				}
				if len(cps)<idx+1{
					for i:=len(cpm[field].([]interface{}));i<idx+1;i++{
						cpm[field]= append(cpm[field].([]interface{}),nil)
					}
				}
				cpm[field].([]interface{})[idx]= value

			}
		}


	}
	return nil

}

func Marshals(querys []QueryProp)(tmp map[string]interface{} ,err error){
	m:=map[string]interface{}{}
	for _,v:=range querys{
		if err:=Marshal(v.Query,m,v.Value);err!=nil{
			return nil,err
		}
	}
	return m,nil

}


func yyp(token string)(string,int,error){
	numidx_start:=0
	numidx_end:=0
	for k,v:=range token{
		t:=string(v)
		if t=="["{
			numidx_start=k
		}
		if t=="]"{
			numidx_end=k
		}
	}
	if numidx_end>0 && numidx_start>=0{
		num,err:=strconv.Atoi(token[numidx_start+1:numidx_end])
		if err !=nil{
			return "",TYPE_KEY,err
		}
		return token[:numidx_start],num,nil
	}
	return token,TYPE_KEY,nil
}


func tokenize2(query string) ([]string, error){
	return strings.Split(strings.Trim(query,"$."),"."),nil
}

func tokenize(query string) ([]string, error) {
	tokens := []string{}
	//	token_start := false
	//	token_end := false
	token := ""

	// fmt.Println("-------------------------------------------------- start")
	for idx, x := range query {
		token += string(x)
		// //fmt.Printf("idx: %d, x: %s, token: %s, tokens: %v\n", idx, string(x), token, tokens)
		if idx == 0 {
			if token == "$" || token == "@" {
				tokens = append(tokens, token[:])
				token = ""
				continue
			} else {
				return nil, fmt.Errorf("should start with '$'")
			}
		}
		if token == "." {
			continue
		} else if token == ".." {
			if tokens[len(tokens)-1] != "*" {
				tokens = append(tokens, "*")
			}
			token = "."
			continue
		} else {
			// fmt.Println("else: ", string(x), token)
			if strings.Contains(token, "[") {
				// fmt.Println(" contains [ ")
				if x == ']' && !strings.HasSuffix(token, "\\]") {
					if token[0] == '.' {
						tokens = append(tokens, token[1:])
					} else {
						tokens = append(tokens, token[:])
					}
					token = ""
					continue
				}
			} else {
				// fmt.Println(" doesn't contains [ ")
				if x == '.' {
					if token[0] == '.' {
						tokens = append(tokens, token[1:len(token)-1])
					} else {
						tokens = append(tokens, token[:len(token)-1])
					}
					token = "."
					continue
				}
			}
		}
	}
	if len(token) > 0 {
		if token[0] == '.' {
			token = token[1:]
			if token != "*" {
				tokens = append(tokens, token[:])
			} else if tokens[len(tokens)-1] != "*" {
				tokens = append(tokens, token[:])
			}
		} else {
			if token != "*" {
				tokens = append(tokens, token[:])
			} else if tokens[len(tokens)-1] != "*" {
				tokens = append(tokens, token[:])
			}
		}
	}
	// fmt.Println("finished tokens: ", tokens)
	// fmt.Println("================================================= done ")
	return tokens, nil
}