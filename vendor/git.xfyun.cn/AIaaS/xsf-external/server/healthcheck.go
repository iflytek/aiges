package xsf

import "sync"

// AdminCheckList holds health checker map
var AdminCheckList map[string]HealthChecker
var HealthMu sync.RWMutex
// HealthChecker health checker interface
type HealthChecker interface {
	Check() error
}

// AddHealthCheck add health checker with name string
func AddHealthCheck(name string, hc HealthChecker) {
	HealthMu.Lock()
	AdminCheckList[name] = hc
	HealthMu.Unlock()
}

func init() {
	HealthMu.Lock()
	AdminCheckList = make(map[string]HealthChecker)
	HealthMu.Unlock()
}
