package xsf

const circuitSeparator = `_`

type Circuit struct {
	CircuitFuncs
	CircuitSettings
	CircuitBreakers
}

func newCircuit() *Circuit {
	c := Circuit{}
	c.init()
	return &c
}
func (c *Circuit) init() {
	c.CircuitFuncs.init(c)
	c.CircuitSettings.init(c)
	c.CircuitBreakers.init(c)
}

//key由svc和op组成
func (c *Circuit) getCommandsKey(svc, op string) string {
	return svc + circuitSeparator + op
}
func (c *Circuit) configure(svc, op string, config CommandConfig) {
	c.CircuitSettings.configureSetting(c.getCommandsKey(svc, op), config)
	c.CircuitFuncs.configureFunc(c.getCommandsKey(svc, op), config)
	c.resetCircuitBreakers()
}
func (c *Circuit) getCircuitBreaker(svc, op string) (*CircuitBreaker, bool, error) {
	return c.CircuitBreakers.getCircuitBreaker(c.getCommandsKey(svc, op))
}
func (c *Circuit) getFunc(svc, op string) *Func {
	return c.CircuitFuncs.getFunc(c.getCommandsKey(svc, op))
}
func (c *Circuit) byPass(svc, op string) bool {
	if nil == c.getFunc(svc, op).fallbackFuncC {
		return true
	}
	return false
}
func (c *Circuit) resetCircuitBreakers() {
	c.flush()
}
