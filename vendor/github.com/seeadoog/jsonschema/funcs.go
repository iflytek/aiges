package jsonschema

import (
	"fmt"
	"strings"
)

func init() {
	SetFunc("split", funcSplit)
	SetFunc("join", funcJoin)
	SetFunc("add", funcAdd)
	SetFunc("append", funcAppend)
	SetFunc("sub", funcSub)
	SetFunc("mod", funcMod)
	SetFunc("div", funcDiv)
	SetFunc("mul", funcMul)
	SetFunc("trimPrefix", funcTrimPrefix)
	SetFunc("trimSuffix", funcTrimSuffix)
	SetFunc("trim", funcTrim)
	SetFunc("replace", funcReplace)
	SetFunc("sprintf", funcSprintf)
	SetFunc("or", funcOr)
	SetFunc("delete", funcDelete)
}

func funcAppend(ctx Context, args ...Value) interface{} {
	bf := strings.Builder{}
	for _, arg := range args {
		v := arg.Get(ctx)
		bf.WriteString(StringOf(v))
	}
	return bf.String()
}

func funcAdd(ctx Context, args ...Value) interface{} {
	var sum float64 = 0
	for _, arg := range args {
		sum += NumberOf(arg.Get(ctx))
	}
	return sum
}
func funcMul(ctx Context, args ...Value) interface{} {
	var sum float64 = 0
	for _, arg := range args {
		sum *= NumberOf(arg.Get(ctx))
	}
	return sum
}

func funcSub(ctx Context, args ...Value) interface{} {
	if len(args) <= 2 {
		return 0
	}

	return NumberOf(args[0].Get(ctx)) - NumberOf(args[1].Get(ctx))
}

func funcDiv(ctx Context, args ...Value) interface{} {
	if len(args) <= 2 {
		return 0
	}

	return NumberOf(args[0].Get(ctx)) / NumberOf(args[1].Get(ctx))
}

func funcMod(ctx Context, args ...Value) interface{} {
	if len(args) <= 2 {
		return 0
	}

	return int(NumberOf(args[0].Get(ctx))) % int(NumberOf(args[1].Get(ctx)))
}

func funcSplit(ctx Context, args ...Value) interface{} {
	if len(args) < 2 {
		return nil
	}
	str := StringOf(args[0].Get(ctx))
	sep := StringOf(args[1].Get(ctx))
	num := -1
	if len(args) >= 3 {
		num = int(NumberOf(args[2].Get(ctx)))
	}
	return strings.SplitN(str, sep, num)
}

func funcJoin(ctx Context, args ...Value) interface{} {
	if len(args) < 2 {
		return ""
	}
	arri, ok := args[0].Get(ctx).([]string)
	sep := StringOf(args[1].Get(ctx))
	if ok {
		return strings.Join(arri, sep)
	}
	arr, ok := args[0].Get(ctx).([]interface{})
	if !ok {
		return ""
	}
	arrs := make([]string, len(arr))
	for i := range arr {
		arrs[i] = StringOf(arr[i])
	}
	return strings.Join(arrs, sep)
}

func funcTrimPrefix(ctx Context, args ...Value) interface{} {
	if len(args) < 2 {
		return ""
	}

	return strings.TrimPrefix(StringOf(args[0].Get(ctx)), StringOf(args[1].Get(ctx)))
}

func funcTrimSuffix(ctx Context, args ...Value) interface{} {
	if len(args) < 2 {
		return ""
	}

	return strings.TrimSuffix(StringOf(args[0].Get(ctx)), StringOf(args[1].Get(ctx)))
}

func funcTrim(ctx Context, args ...Value) interface{} {
	if len(args) < 2 {
		return ""
	}

	return strings.Trim(StringOf(args[0].Get(ctx)), StringOf(args[1].Get(ctx)))
}

func funcReplace(ctx Context, args ...Value) interface{} {
	if len(args) < 3 {
		return ""
	}

	return strings.Replace(StringOf(args[0].Get(ctx)),StringOf(args[1].Get(ctx)),StringOf(args[2].Get(ctx)),-1)
}

func funcSprintf(ctx Context, args ...Value) interface{}{
	if len(args) < 1{
		return nil
	}
	ags :=make([]interface{},0,len(args)-1)
	for _, value := range args[1:] {
		ags = append(ags,value.Get(ctx))
	}

	return fmt.Sprintf(StringOf(args[0].Get(ctx)),ags...)
}




func funcOr(ctx Context, args ...Value) interface{}{
	for _, arg := range args {
		val := arg.Get(ctx)
		if notNil(val){
			return val
		}
	}
	return nil
}

func funcDelete(ctx Context, args ...Value) interface{}{
	for _, arg := range args {
		delete(ctx,StringOf(arg.Get(ctx)))
	}
	return nil
}
