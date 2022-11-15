package jsonschema

var valueFuncs = map[string]Func{
	"append": funcAppend,
	"add":    funcAdd,
}

func SetFunc(name string, fun Func) {
	valueFuncs[name] = fun
}

type Context map[string]interface{}
type Value interface {
	Get(ctx Context) interface{}
}

type Const struct {
	Val interface{}
}

func (c Const) Get(ctx Context) interface{} {
	return c.Val
}

type Var struct {
	Key *JsonPathCompiled
}

func (v Var) Get(ctx Context) interface{} {
	val, err := v.Key.Get(map[string]interface{}(ctx))
	if err != nil {
		return nil
	}
	return val
}

type VarFunc struct {
	funName string
	args    []Value
}

func (v VarFunc) Get(ctx Context) interface{} {
	fun := valueFuncs[v.funName]
	if fun == nil {
		return nil
	}
	return fun(ctx, v.args...)

}

type Func func(ctx Context, args ...Value) interface{}

func parseFuncValue(name string,args []interface{})(Value,error){
	argsv := make([]Value, len(args))
	for idx, arg := range args {
		argv, err := parseValue(arg)
		if err != nil {
			return nil, err
		}
		argsv[idx] = argv
	}

	return &VarFunc{
		funName: name,
		args: argsv,
	},nil
}

func parseValue(i interface{}) (Value, error) {
	switch i.(type) {
	case map[string]interface{}:
		m := i.(map[string]interface{})
		funName := StringOf(m["func"])
		if valueFuncs[funName]  != nil{
			args, ok := m["args"].([]interface{})
			if !ok {
				return &Const{
					Val: i,
				}, nil
			}
			return parseFuncValue(funName,args)
		}
		return &Const{
			Val: i,
		}, nil

	case string:
		str := i.(string)
		if len(str) > 3 && str[0] == '$' && str[1] == '{' && str[len(str)-1] == '}' {
			jp, err := parseJpathCompiled(str[2 : len(str)-1])
			if err != nil {
				return nil, err
			}
			return &Var{Key: jp}, nil
		}
		return &Const{Val: i}, nil
	case []interface{}:
		vv := i.([]interface{})
		if len(vv) > 1{
			str  := StringOf(vv[0])
			if len(str) > 0 && str[0]=='$'{
				funcName := str[1:]
				if valueFuncs[funcName] != nil{
					args := vv[1:]
					return parseFuncValue(funcName,args)
				}

			}
		}
		return &Const{Val: i}, nil
	default:
		return &Const{Val: i}, nil
	}
}
