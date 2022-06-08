package xsf

import (
	"github.com/xfyun/sonar"
)

type KV struct {
	Key   string
	Value interface{}
}
type sonarLogInterface interface {
	Infof(format string, params ...interface{})
	Debugf(format string, params ...interface{})
	Errorf(format string, params ...interface{}) error
}
type SonarAdapter struct {
	metricType  string
	endpoint    string
	serviceName string
	port        string
	ds          string
	able        bool

	logger             sonarLogInterface
	sonarDumpEnable    bool
	sonarDeliverEnable bool
	sonarHost          string
	sonarPort          string
	sonarBackend       int
}

type SonarAdapterOpt func(*SonarAdapter)

func WithSonarAdapterAble(able bool) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.able = able
	}
}
func WithSonarAdapterDs(ds string) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.ds = ds
	}
}
func WithSonarAdapterMetricEndpoint(metricEndpoint string) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.endpoint = metricEndpoint
	}
}
func WithSonarAdapterMetricServiceName(metricServiceName string) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.serviceName = metricServiceName
	}
}
func WithSonarAdapterMetricPort(metricPort string) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.port = metricPort
	}
}
func WithSonarAdapterLogger(logger sonarLogInterface) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.logger = logger
	}
}
func WithSonarAdapterSonarDumpEnable(sonarDumpEnable bool) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.sonarDumpEnable = sonarDumpEnable
	}
}
func WithSonarAdapterSonarDeliverEnable(sonarDeliverEnable bool) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.sonarDeliverEnable = sonarDeliverEnable
	}
}
func WithSonarAdapterSonarHost(sonarHost string) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.sonarHost = sonarHost
	}
}
func WithSonarAdapterSonarPort(sonarPort string) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.sonarPort = sonarPort
	}
}
func WithSonarAdapterSonarBackend(sonarBackend int) SonarAdapterOpt {
	return func(sa *SonarAdapter) {
		sa.sonarBackend = sonarBackend
	}
}

//initSonar(able bool, ds string, metricEndpoint string, metricServiceName string, metricPort string, logger sonarLogInterface, dumpenable bool, deliverenable bool, host string, port string, backend int)
func (s *SonarAdapter) initSonar(opts ...SonarAdapterOpt) error {
	for _, opt := range opts {
		opt(s)
	}
	if !s.able {
		return nil
	}

	sonar.Logger = s.logger
	sonar.DumpEnable = s.sonarDumpEnable
	sonar.DeliverEnable = s.sonarDeliverEnable
	return sonar.Init(s.sonarHost, s.sonarPort, s.sonarBackend)
}

func (s *SonarAdapter) newMetricWithNamePort(metricName string, metricNameVal float64, kv ...KV) error {
	if !s.able {
		return nil
	}
	metricData := sonar.NewMetricWithNamePort(s.metricType, metricName, s.endpoint, s.serviceName, s.port, s.ds).WithValue(metricNameVal)
	for _, v := range kv {
		metricData.Tag(sonar.KV{v.Key, v.Value})
	}
	return metricData.Flush()
}
func (s *SonarAdapter) NewMetric(metricName string, metricNameVal float64, kv ...KV) error {
	return s.newMetricWithNamePort(metricName, metricNameVal, kv...)
}
