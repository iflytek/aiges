package xsf

type executorPool struct {
	parent  *CircuitBreaker
	Name    string
	Metrics *poolMetrics
	Max     int
	Tickets chan *struct{}
}

func newExecutorPool(parent *CircuitBreaker, name string) *executorPool {
	p := &executorPool{}
	p.Name = name
	p.parent = parent
	p.Metrics = newPoolMetrics(name, p.getSettings(name).SleepWindow)
	p.Max = p.getSettings(name).MaxConcurrentRequests
	p.Tickets = make(chan *struct{}, p.Max)
	for i := 0; i < p.Max; i++ {
		p.Tickets <- &struct{}{}
	}

	return p
}
func (p *executorPool) getSettings(name string) *Settings {
	return p.parent.parent.parent.getSettings(name)
}
func (p *executorPool) Return(ticket *struct{}) {
	if ticket == nil {
		return
	}

	p.Metrics.Updates <- poolMetricsUpdate{
		activeCount: p.ActiveCount(),
	}
	p.Tickets <- ticket
}

func (p *executorPool) ActiveCount() int {
	return p.Max - len(p.Tickets)
}
