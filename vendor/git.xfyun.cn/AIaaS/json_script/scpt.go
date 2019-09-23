package jsonscpt

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)
var (
	systemId = map[string]int{
		"append":1,
		"len":1,
		"split":1,
		"printf":1,
		"print":1,
		"sprintf":1,
		"add":1,
		"json_m":1,
		"isnil":1,
		"delete":1,
		"and":1,
		"eq":1,
		"or":1,
		"true":1,
		"false":1,
		"gt":1,
		"ge":1,
		"lt":1,
		"le":1,
		"not":1,
		"return":1,
		"in":1,
		"contains":1,
		"join":1,
		"get":1,
		"set":1,
		"exit":1,
		"trim":1,
	}
)

func isSystemId(s string)bool  {
	if systemId[s]==1{
		return true
	}
	return false
}

type Context struct {
	table map[string]interface{}  // save all variables
}

func NewVm() *Context {
	c:=&Context{
		table:map[string]interface{}{},
	}
	c.init()
	return c
}


func (ctx *Context)init()  {
	ctx.SetFunc("append",apd)
	ctx.SetFunc("len",lens)
	ctx.SetFunc("split",split)
	ctx.SetFunc("printf",printf)
	ctx.SetFunc("print",printlnn)
	ctx.SetFunc("sprintf",sprintf)
	ctx.SetFunc("add",add)
	ctx.SetFunc("json_m", jsonMarshal)
	ctx.SetFunc("delete", deleteFun)
	ctx.SetFunc("isnil", isNil)
	ctx.SetFunc("and", and)
	ctx.SetFunc("or", or)
	ctx.SetFunc("eq", eq)
	ctx.SetFunc("gt", gt)
	ctx.SetFunc("ge", ge)
	ctx.SetFunc("le", le)
	ctx.SetFunc("lt", lt)
	ctx.SetFunc("not", not)
	ctx.SetFunc("return", ret)
	ctx.SetFunc("exit", exit)
	ctx.SetFunc("in", in)
	ctx.SetFunc("contains", contains)
	ctx.SetFunc("join", join)
	ctx.SetFunc("set", set)
	ctx.SetFunc("get", get)
	ctx.SetFunc("input", input)
	ctx.SetFunc("trim", trim)


}

func (ctx *Context)Func(name string,params ...interface{})  {

}

func (ctx *Context)SetFunc(name string,value Func)  {
	ctx.table[name] = value
}

func (ctx *Context)Set(k string,v interface{})  {
	MarshalInterface(k,ctx.table,v)
}

func (ctx *Context)Get(k string)interface{}  {
	v,_:=CachedJsonpathLookUp(ctx.table,k)
	return v
}

func (ctx *Context) ExecJsonObject(v interface{}) error {
	//if stv,ok:=v.(string);ok{
	//	json.Unmarshal([]byte(stv),&v)
	//}
	cpd,err:= CompileExpFromJsonObject(v)
	if err !=nil{
		return err
	}
	return ctx.Execute(cpd)
}
func (ctx *Context)ExecJson(s []byte) error {
	var i interface{}
	err:=json.Unmarshal(s,&i)
	if err !=nil{
		return err
	}
	return ctx.ExecJsonObject(i)
}

func (c *Context)Execute(exp Exp) error {
	return exp.Exec(c)
}

func (c *Context)SafeExecute(exp Exp, fatalHandler func(err interface{})) error {
	defer func() {
		if err := recover();err !=nil{
			if fatalHandler !=nil{
				fatalHandler(err)
			}
		}
	}()
	return exp.Exec(c)
}

func CompileExpFromJson(b []byte)(Exp,error){
	var i interface{}
	if err:=json.Unmarshal(b,&i);err !=nil{
		return nil,err
	}
	return CompileExpFromJsonObject(i)
}

func boolValid(s string) bool {
	return true
}





func CompileExpFromJsonObject(v interface{}) (Exp,error) {
	if exp,ok:=v.(string);ok{
		if exp=="break"{   // parse break
			return &BreakExp{},nil
		}

		e,err:=parseSetExp(exp)
		if err !=nil{
			return nil, err
		}
		return e,nil
	}
	if m,ok:=v.(map[string]interface{});ok{

		if ifexp,ok:=m["if"].(string);ok{
			exp:=&IfExp{}
			efe,err:=parseBoolExp(ifexp);
			if err==nil{
				exp.If = efe
			}else{
				return nil,err
			}
			if then,ok:=m["then"];ok && then !=nil {
				var parsedExp,err = CompileExpFromJsonObject(then)
				if err !=nil{
					return nil,err
				}
				exp.Then = parsedExp
			}else{
				//return  nil,errors.New("line:"+ifexp+" has no then block")
			}
			if el,ok:=m["else"];ok && el !=nil{
				var parsedExp,err = CompileExpFromJsonObject(el)
				if err !=nil{
					return nil,err
				}
				exp.Else = parsedExp
			}
			return exp,nil
		}else if forexp,ok:=m["for"].(string);ok{  //parse for

			if forrangeReg.MatchString(forexp){ // for range
				r:=forrangeReg.FindAllStringSubmatch(forexp,-1)
				if len(r)>0{
					if len(r[0])>3{
						key:=r[0][1]
						val:=r[0][2]
						if isSystemId(key){
							return nil,errors.New("system id cannot be variable:"+key)
						}
						if isSystemId(val){
							return nil,errors.New("system id cannot be variable:"+val)
						}
						forVal:=r[0][3]
						parsedForValue,err:=parseValue(forVal)
						if err !=nil{
							return nil,err
						}
						do,err:=CompileExpFromJsonObject(m["do"])
						if err !=nil{
							return nil,err
						}
						r:=&ForRangeExp{
							Value:parsedForValue,
							Do:do,
							SubIdx:key,
							SubValue:val,
						}
						return r,nil
					}
				}
				return nil,errors.New("invalid for range exp"+ConvertToString(v))

			}else{  // for bool
				exp:=&ForExp{}
				efe,err:=parseBoolExp(forexp);
				if err !=nil{
					return nil,err
				}
				exp.Addtion = efe
				if do,ok:=m["do"];ok && do !=nil{
					blexp,err:=CompileExpFromJsonObject(do)
					if err !=nil{
						return nil,err
					}
					exp.Do = blexp
				}else{
					return nil,errors.New("for do is nil")
				}
				return exp,nil
			}


		}else if data,ok:=m["data"];ok{
			return &DataExp{
				Key:ConvertToString(m["key"]),
				Data:data,
			},nil
		}else if gofun,ok:=m["go"];ok{

			exp,err:=CompileExpFromJsonObject(gofun)
			if err !=nil{
				return nil,err
			}
			return &GoFunc{Exp:exp},nil
		}else if fun,ok:=m["func"].(string);ok{
			if params,ok:=m["do"];ok{
				return parseFunc(fun,params)
			}
		}

		return nil,errors.New("invalid object:"+fmt.Sprintf("%v",v))
	}

	if eps,ok:=v.([]interface{});ok {
		var parsedExp = Exps{}
		for i := 0; i < len(eps); i++ {
			e,err:= CompileExpFromJsonObject(eps[i])
			if err !=nil{
				return nil,err
			}else{
				parsedExp = append(parsedExp,e)
			}
		}
		return parsedExp,nil
	}

	return nil,errors.New("invalid exp:"+ConvertToString(v))
}
var funcReg = regexp.MustCompile(`(\w+)\((.*)\)`)
func parseFunc(s string,body interface{})(Exp,error){
	if !funcReg.MatchString(s){
		return nil, errors.New("invalid func define:" + s)
	}
	r:=funcReg.FindAllStringSubmatch(s,-1)
	if len(r)>0{
		v:=r[0]
		if len(v)>=2{
			fun:=v[1]
			if isSystemId(fun){
				return nil, errors.New("variable cannot be system Id"+fun)
			}
			do,err:=CompileExpFromJsonObject(body)
			if err !=nil{
				return nil, err
			}
			var params []string
			if len(v[2])>0{
				params = strings.Split(v[2],",")
				for i:=0;i< len(params);i++{
					params[i] = strings.Trim(params[i]," ")
					if isSystemId(params[i]){
						return nil, errors.New("variable cannot be system Id"+params[i])
					}
				}
			}
			fd:=&FuncDefine{
				FuncName:fun,
				Params:params,
				Body:do,
			}
			return fd, nil

		}
	}
	return nil, nil
}

