package finder

import (
	"fmt"
	"runtime"
)

type FinderError struct {
	Ret  ReturnCode
	Func string
	File string
	Line int
	Desc string
}

func (fe *FinderError) Error() string {
	format := `An error caught in %s：%d, errCode : %d [%s].`

	return fmt.Sprintf(format, fe.File, fe.Line, fe.Ret, fe.Desc)
}

func NewFinderError(ret ReturnCode) *FinderError {
	fe := FinderError{
		Ret: ret,
	}
	fe.Desc = ret.String()
	if _, file, line, ok := runtime.Caller(1); ok {
		fe.File = file
		fe.Line = line
	}
	return &fe
}
