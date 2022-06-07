package xsf

import (
	"github.com/xfyun/xsf/server/internal/gauge"
	"github.com/prometheus/client_golang/prometheus"
)

type GaugeEx interface {
	prometheus.Gauge
}
type GaugeVecEx struct {
	base *gauge.GaugeVec
}

func (v *GaugeVecEx) Describe(ch chan<- *prometheus.Desc) {
	v.base.Describe(ch)
}

func (v *GaugeVecEx) Collect(ch chan<- prometheus.Metric) {
	v.base.Collect(ch)
}
func (v *GaugeVecEx) WithLabelValues(lvs ...string) GaugeEx {
	g, err := v.base.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return g
}
func NewGaugeVecEx(opts GaugeOpts, labelNames []string) *GaugeVecEx {
	gaugeVecEx := gauge.NewGaugeVec(prometheus.GaugeOpts{
		Name: opts.Name,
		Help: opts.Help,
	}, labelNames)
	return &GaugeVecEx{gaugeVecEx}
}
