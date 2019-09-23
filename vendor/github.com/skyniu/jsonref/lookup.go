package jsonref

import (
	"errors"
	"fmt"
)

func JsonPathLookup(i interface{},path string) (interface{},error )  {
	return Lookup(path,i)
}

func Lookup(query string,val interface{}) (interface{},error )  {
	tokens,err:=tokenize2(query)
	if err !=nil{
		return nil,err
	}
	//log.Println(tokens)
	var cp = val
	for k,v:=range tokens{
		field,idx,err:=yyp(v)
		//log.Println(yyp(v))
		if err !=nil{
			return nil,err
		}
		if k< len(tokens)-1{
			if idx == TYPE_KEY{
				cpm,ok:=cp.(map[string]interface{})
				if !ok{
					return nil,errors.New(fmt.Sprintf("read field failed ,%s cannot convert1 to map",field))
				}
				cp = cpm[field]
			}else{
				if field==""{
					cps,ok:=cp.([]interface{})
					if !ok{
						cpss,ok:=cp.([]map[string]interface{})
						if !ok{
							return nil,errors.New(fmt.Sprintf("read field failed ,%s cannot convert1 to array",field))
						}
						cp = cpss[idx]
						continue
					}
					cp = cps[idx]
					continue
				}
				cpm,ok:=cp.(map[string]interface{})
				if !ok{
					return nil,errors.New(fmt.Sprintf("read field failed ,%s cannot convert2 to map",field))
				}
				//log.Println(cpm[field],reflect.TypeOf(cpm[field]))
				cps,ok:=cpm[field].([]interface{})
				if !ok{
					return nil,errors.New(fmt.Sprintf("read field failed ,%s cannot convert2 to array",field))
				}
				cp =cps[idx]
			}
		}else{
			if idx == TYPE_KEY{

				cpm,ok:=cp.(map[string]interface{})
				if !ok{
					return nil,errors.New(fmt.Sprintf("read field failed ,%s cannot convert3 to map",field))
				}
				return cpm[field],nil
			}else{
				if field==""{
					cps,ok:=cp.([]interface{})
					if !ok{
						cpss,ok:=cp.([]map[string]interface{})
						if !ok{
							return nil,errors.New(fmt.Sprintf("read field failed ,%s cannot convert3 to array",field))
						}
						return cpss[idx],nil
					}
					return cps[idx],nil
				}
				cpm,ok:=cp.(map[string]interface{})
				if !ok{
					return nil,errors.New(fmt.Sprintf("create field failed ,%s cannot convert4 to map",field))
				}
				cps,ok:=cpm[field].([]interface{})
				if !ok{
					return nil,errors.New(fmt.Sprintf("create field failed ,%s cannot convert4 to array",field))
				}
				return cps[idx],nil
				return nil,nil
			}
		}

	}
	return nil,nil
}

