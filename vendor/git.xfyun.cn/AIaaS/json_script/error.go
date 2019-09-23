package jsonscpt

import "fmt"
const(
	CodeFuncDoesNotExists = -1
)
type ErrorExit struct {
	Message string
	Code int
	Value interface{}
}

func (e *ErrorExit)Error()string  {
	return fmt.Sprintf("code=%d,message=%s",e.Code,e.Message)
}

var breakError = &BreakError{}
//break error
type BreakError struct {

}

func (e *BreakError)Error()string  {
	return "break"
}

func IsReturnError(err error)(e *ErrorReturn,ok bool ) {
	e,ok =err.(*ErrorReturn)
	return
}

func IsExitErrorI(err interface{})(e *ErrorExit,ok bool ) {
	e,ok =err.(*ErrorExit)
	return
}

type ErrorReturn struct {
	Value interface{}
}

func (e *ErrorReturn)Error()string  {
	return fmt.Sprintf("%v",e.Value)
}

func IsExitError(err error)(e *ErrorExit,ok bool ) {
	e,ok =err.(*ErrorExit)
	return
}