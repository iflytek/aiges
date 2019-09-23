package xsf

import (
	"sync"
)

type Op func(*Req, *Span, *ToolBox) (*Res, error)
type OpRouter struct {
	router sync.Map
}

func NewOpRouter() *OpRouter {
	sm := sync.Map{}
	return &OpRouter{router: sm}
}
func (o *OpRouter) Store(k string, op Op) {
	o.router.Store(k, op)
}
func (o *OpRouter) load(k string) (op Op, ok bool) {
	opInterface, opInterfaceOk := o.router.Load(k)
	if !opInterfaceOk {
		return nil, opInterfaceOk
	}
	return opInterface.(Op), opInterfaceOk
}
