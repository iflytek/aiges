package xsf

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type CircuitError struct {
	Message string
	Code    int32
}

func (e CircuitError) Error() string {
	return e.Message
}
func (e CircuitError) ErrCode() int32 {
	return e.Code
}

const (
	UnknownCode int32 = -1
)
const (
	RecursiveThreshold int32  = 1
	currentHierarchy   string = "currentHierarchy"
)

var (
	ErrMaxConcurrency = CircuitError{Message: "max concurrency"}

	ErrCircuitOpen = CircuitError{Message: "circuit open"}

	ErrCircuitOpenLackOfTicket = CircuitError{Message: "circuit open:lack of ticket"}

	ErrTimeout = CircuitError{Message: "timeout", Code: -1}

	ErrRecursiveHierarchyExceed = CircuitError{Message: "recursive hierarchy exceed threshold", Code: -1}
)

type xsfOpt interface {
	apply(interface{})
}

type commandFunc func(*CommandConfig)

func (f commandFunc) apply(l interface{}) {
	f(l.(*CommandConfig))
}
func WithCommandTimeout(Timeout int) xsfOpt {
	return commandFunc(func(in *CommandConfig) { in.Timeout = Timeout })
}
func WithCommandMaxConcurrentRequests(MaxConcurrentRequests int) xsfOpt {
	return commandFunc(func(in *CommandConfig) { in.MaxConcurrentRequests = MaxConcurrentRequests })
}
func WithCommandRequestVolumeThreshold(RequestVolumeThreshold int) xsfOpt {
	return commandFunc(func(in *CommandConfig) { in.RequestVolumeThreshold = RequestVolumeThreshold })
}
func WithCommandSleepWindow(SleepWindow int) xsfOpt {
	return commandFunc(func(in *CommandConfig) { in.SleepWindow = SleepWindow })
}
func WithCommandErrorPercentThreshold(ErrorPercentThreshold int) xsfOpt {
	return commandFunc(func(in *CommandConfig) { in.ErrorPercentThreshold = ErrorPercentThreshold })
}
func WithCommandFallback(fallback HystrixFallback) xsfOpt {
	return commandFunc(func(in *CommandConfig) { in.ErrorFallback = fallback })
}
func (c *Client) ConfigureHystrix(service string, op string, opts ...xsfOpt) error {
	hystrixConfigInst := &CommandConfig{}
	for _, opt := range opts {
		opt.apply(hystrixConfigInst)
	}
	c.circuit.configure(service, op, *hystrixConfigInst)
	return nil
}

type fallbackParams struct {
	//inParams
	service string
	op      string
	r       *Req
	tm      time.Duration
	//outParams
	s       *Res
	errcode int32
	e       error
}
type command struct {
	sync.Mutex
	ticket      *struct{}
	start       time.Time
	circuit     *CircuitBreaker
	runDuration time.Duration
	events      []string
	fallback    HystrixFallback

	params *fallbackParams
}

func (c *command) reportEvent(eventType string) {
	c.Lock()
	defer c.Unlock()

	c.events = append(c.events, eventType)
}

func (c *command) tryFallback(ctx context.Context, err error) error {
	if c.fallback == nil {
		return err
	}

	c.params.s, c.params.errcode, c.params.e = c.fallback(
		ctx,
		err,
		c.params.service,
		c.params.op,
		c.params.r)
	if c.params.e != nil {
		c.reportEvent("fallback-failure")
		return fmt.Errorf("fallback failed with '%v'. run error was '%v'", c.params.e, err)
	}

	c.reportEvent("fallback-success")

	return nil
}

func (c *command) errorWithFallback(ctx context.Context, err error) error {

	eventType := "failure"
	if err == ErrCircuitOpen {
		eventType = "short-circuit"
	} else if err == ErrMaxConcurrency {
		eventType = "rejected"
	} else if err == ErrTimeout {
		eventType = "timeout"
	} else if err == context.Canceled {
		eventType = "context_canceled"
	} else if err == context.DeadlineExceeded {
		eventType = "context_deadline_exceeded"
	}

	c.reportEvent(eventType)
	fallbackErr := c.tryFallback(incrCtx(ctx, currentHierarchy), err)

	return fallbackErr
}

func (c *Caller) CallWrapper(service string, op string, r *Req, tm time.Duration) (*Res, int32, error) {

	ctx, cancel := context.WithTimeout(context.Background(), tm)
	defer cancel()

	return c.CallCtx(ctx, service, op, r)
}

func (c *Caller) CallCtx(ctx context.Context, service string, op string, r *Req) (*Res, int32, error) {
	c.cli.Log.Infow("CallWrapper", "service", service, "op", op)

	//由于递归降级的存在，此处判断超时是否失败
	select {
	case <-ctx.Done():
		return nil, UnknownCode, ErrTimeout
	default:
	}

	//参数异常
	if len(op) == 0 {
		return nil, INVAILDPARAM, EINAILDOP
	}
	r.SetOp(op)

	//旁路模式
	if c.cli.circuit.byPass(service, op) {
		return c.oneShortCall(c.apiVersion, service, r, extractTimeout(ctx))
	}

	//防止递归层级超过阈值
	if extractHierarchy(ctx) > RecursiveThreshold {
		return nil, ErrRecursiveHierarchyExceed.Code, ErrRecursiveHierarchyExceed
	}

	//hystrix代理对象
	cmd := &command{
		start: time.Now(),
		params: &fallbackParams{
			service: service,
			op:      op,
			r:       r,
		},
	}

	circuit, _, err := c.cli.circuit.getCircuitBreaker(service, op)
	if err != nil {
		panic("TODO") //TODO necessary
	}
	cmd.circuit = circuit
	cmd.fallback = c.cli.circuit.getFunc(service, op).fallbackFuncC

	returnTicket := func() {
		cmd.circuit.GetExecutorPool().Return(cmd.ticket)
	}
	reportAllEvent := func() {
		err := cmd.circuit.ReportEvent(cmd.events, cmd.start, cmd.runDuration)
		if err != nil {
			panic(err) //todo not complete
		}
	}
	if allowRequest, allowRequestErr := cmd.circuit.AllowRequest(); !allowRequest {
		c.cli.Log.Warnw("CallWrapper circuit open!!!", "service", service, "op", op)
		returnTicket()
		defer reportAllEvent()
		err := cmd.errorWithFallback(ctx, allowRequestErr)
		if err != nil {
			cmd.params.e = err
		}
		return cmd.params.s, cmd.params.errcode, cmd.params.e
	}
	c.cli.Log.Infow("CallWrapper circuit closed!!!", "service", service, "op", op)

	select {
	case cmd.ticket = <-circuit.executorPool.Tickets:
	default:
		{
			c.cli.Log.Warnw("CallWrapper lack of ticket", "service", service, "op", op)
			returnTicket()
			defer reportAllEvent()
			err := cmd.errorWithFallback(ctx, ErrCircuitOpenLackOfTicket)
			if err != nil {
				cmd.params.e = err
			}
			return cmd.params.s, cmd.params.errcode, cmd.params.e
		}
	}

	s, errcode, e := c.oneShortCall(c.apiVersion, service, r, extractTimeout(ctx))
	returnTicket()
	if errcode != 0 || e != nil {
		c.cli.Log.Warnw("CallWrapper call error", "service", service, "op", op)
		defer reportAllEvent()
		err := cmd.errorWithFallback(ctx, e)
		if err != nil {
			cmd.params.e = err
		}
		return cmd.params.s, cmd.params.errcode, cmd.params.e
	}
	cmd.runDuration = time.Since(cmd.start)
	cmd.reportEvent("success")
	reportAllEvent()
	return s, errcode, e
}

func incrCtx(ctx context.Context, key string) context.Context {
	val := ctx.Value(key)
	if val == nil {
		ctx = context.WithValue(ctx, key, 1)
		return ctx
	}
	switch val.(type) {
	case int:
		return context.WithValue(ctx, key, val.(int)+1)
	case int32:
		return context.WithValue(ctx, key, val.(int32)+1)
	case int64:
		return context.WithValue(ctx, key, val.(int64)+1)
	default:
		return context.WithValue(ctx, key, 1)
	}
}
func extractTimeout(ctx context.Context) time.Duration {
	deadLine, deadLineOk := ctx.Deadline()
	if !deadLineOk {
		return 0
	}
	return deadLine.Sub(time.Now())
}
func extractHierarchy(ctx context.Context) int32 {
	val := ctx.Value(currentHierarchy)
	if val == nil {
		return 0
	}
	switch val.(type) {
	case int:
		return int32(val.(int))
	case int32:
		return val.(int32)
	case int64:
		return int32(val.(int64))
	default:
		return 0
	}
}
