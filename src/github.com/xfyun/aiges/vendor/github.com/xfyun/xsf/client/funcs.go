package xsf

import (
	"context"
	"sync"
)

type HystrixFallback func(
	ctx context.Context,
	err error,
	service string,
	op string,
	r *Req) (
	s *Res,
	errcode int32,
	e error)

type Func struct {
	fallbackFuncC HystrixFallback
}
type CircuitFuncs struct {
	circuit           *Circuit
	circuitFuncs      map[string]*Func
	circuitFuncsMutex *sync.RWMutex
}

func (c *CircuitFuncs) init(parent *Circuit) {
	c.circuitFuncs = make(map[string]*Func)
	c.circuitFuncsMutex = &sync.RWMutex{}
	c.circuit = parent
}
func (c *CircuitFuncs) getFunc(name string) *Func {
	c.circuitFuncsMutex.RLock()
	_func, exists := c.circuitFuncs[name]
	c.circuitFuncsMutex.RUnlock()

	if !exists {
		c.configureFunc(name, CommandConfig{})
		_func = c.getFunc(name)
	}

	return _func
}
func (c *CircuitFuncs) configureFunc(name string, config CommandConfig) {
	c.circuitFuncsMutex.Lock()
	c.circuitFuncs[name] = &Func{config.ErrorFallback}
	c.circuitFuncsMutex.Unlock()
}
