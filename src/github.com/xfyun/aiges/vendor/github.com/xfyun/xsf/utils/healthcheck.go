package utils

// AdminCheckList holds health checker map
var AdminCheckList map[string]HealthChecker

// HealthChecker health checker interface
type HealthChecker interface {
	Check() error
}

// AddHealthCheck add health checker with name string
func AddHealthCheck(name string, hc HealthChecker) {
	AdminCheckList[name] = hc
}

func init() {
	AdminCheckList = make(map[string]HealthChecker)
}
