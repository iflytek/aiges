package xsf

import (
	"sync"
	"time"
)

var (
	DefaultTimeout = 1000

	DefaultMaxConcurrent = 10

	DefaultVolumeThreshold = 20

	DefaultSleepWindow = 5000

	DefaultErrorPercentThreshold = 50
)

type Settings struct {
	Timeout                time.Duration
	MaxConcurrentRequests  int
	RequestVolumeThreshold uint64
	SleepWindow            time.Duration
	ErrorPercentThreshold  int
}

type CommandConfig struct {
	Timeout                int             `json:"timeout"` /*暂时无用*/
	MaxConcurrentRequests  int             `json:"max_concurrent_requests"`
	RequestVolumeThreshold int             `json:"request_volume_threshold"`
	SleepWindow            int             `json:"sleep_window"`
	ErrorPercentThreshold  int             `json:"error_percent_threshold"`
	ErrorFallback          HystrixFallback `json:"error_fallback"`
}
type CircuitSettings struct {
	circuit              *Circuit
	circuitSettings      map[string]*Settings
	circuitSettingsMutex *sync.RWMutex
}

func (s *CircuitSettings) init(parent *Circuit) {
	s.circuitSettings = make(map[string]*Settings)
	s.circuitSettingsMutex = &sync.RWMutex{}
	s.circuit = parent
}

func (s *CircuitSettings) configureSetting(name string, config CommandConfig) {
	s.circuitSettingsMutex.Lock()
	defer s.circuitSettingsMutex.Unlock()

	timeout := DefaultTimeout
	if config.Timeout != 0 {
		timeout = config.Timeout
	}

	max := DefaultMaxConcurrent
	if config.MaxConcurrentRequests != 0 {
		max = config.MaxConcurrentRequests
	}

	volume := DefaultVolumeThreshold
	if config.RequestVolumeThreshold != 0 {
		volume = config.RequestVolumeThreshold
	}

	sleep := DefaultSleepWindow
	if config.SleepWindow != 0 {
		sleep = config.SleepWindow
	}

	errorPercent := DefaultErrorPercentThreshold
	if config.ErrorPercentThreshold != 0 {
		errorPercent = config.ErrorPercentThreshold
	}

	s.circuitSettings[name] = &Settings{
		Timeout:                time.Duration(timeout) * time.Millisecond,
		MaxConcurrentRequests:  max,
		RequestVolumeThreshold: uint64(volume),
		SleepWindow:            time.Duration(sleep) * time.Millisecond,
		ErrorPercentThreshold:  errorPercent,
	}
}

func (s *CircuitSettings) getSettings(name string) *Settings {
	s.circuitSettingsMutex.RLock()
	settings, exists := s.circuitSettings[name]
	s.circuitSettingsMutex.RUnlock()

	if !exists {
		s.configureSetting(name, CommandConfig{})
		settings = s.getSettings(name)
	}

	return settings
}

func (s *CircuitSettings) getCircuitSettings() map[string]*Settings {
	_copy := make(map[string]*Settings)

	s.circuitSettingsMutex.RLock()
	for key, val := range s.circuitSettings {
		_copy[key] = val
	}
	s.circuitSettingsMutex.RUnlock()

	return _copy
}
