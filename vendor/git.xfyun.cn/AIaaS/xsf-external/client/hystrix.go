package xsf

import "sync"

type xsfOpt interface {
	apply(interface{})
}

var circuitSettings map[string]*HystrixConfig

var settingsMutex *sync.RWMutex

type HystrixConfig struct {
	Timeout                int `json:"timeout"`                  //调用超时
	MaxConcurrentRequests  int `json:"max_concurrent_requests"`  //最大请求数
	RequestVolumeThreshold int `json:"request_volume_threshold"` //请求数阈值
	SleepWindow            int `json:"sleep_window"`             //窗口大小
	ErrorPercentThreshold  int `json:"error_percent_threshold"`  //错误阈值
}

type commandFunc func(*HystrixConfig)

func (f commandFunc) apply(l interface{}) {
	f(l.(*HystrixConfig))
}
func withCommandTimeout(Timeout int) xsfOpt                               { return commandFunc(func(in *HystrixConfig) { in.Timeout = Timeout }) }
func withCommandMaxConcurrentRequests(MaxConcurrentRequests int) xsfOpt   { return commandFunc(func(in *HystrixConfig) { in.MaxConcurrentRequests = MaxConcurrentRequests }) }
func withCommandRequestVolumeThreshold(RequestVolumeThreshold int) xsfOpt { return commandFunc(func(in *HystrixConfig) { in.RequestVolumeThreshold = RequestVolumeThreshold }) }
func withCommandSleepWindow(SleepWindow int) xsfOpt                       { return commandFunc(func(in *HystrixConfig) { in.SleepWindow = SleepWindow }) }
func withCommandErrorPercentThreshold(ErrorPercentThreshold int) xsfOpt   { return commandFunc(func(in *HystrixConfig) { in.ErrorPercentThreshold = ErrorPercentThreshold }) }

func configureHystrix(name string, opts ...xsfOpt) {

	cfgInst := &HystrixConfig{}
	for _, opt := range opts {
		opt.apply(cfgInst)
	}

	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	circuitSettings[name] = cfgInst
}
