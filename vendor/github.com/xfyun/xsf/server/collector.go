package xsf

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

//gauge
type GaugeOpts struct {
	Name string
	Help string
}
type Gauge interface {
	prometheus.Gauge
}
type GaugeVec struct {
	base *prometheus.GaugeVec
}

func (v *GaugeVec) Describe(ch chan<- *prometheus.Desc) {
	v.base.Describe(ch)
}

func (v *GaugeVec) Collect(ch chan<- prometheus.Metric) {
	v.base.Collect(ch)
}

func (v *GaugeVec) WithLabelValues(lvs ...string) Gauge {
	g, err := v.base.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return g
}
func NewGaugeVec(opts GaugeOpts, labelNames []string) *GaugeVec {
	gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: opts.Name,
		Help: opts.Help,
	}, labelNames)
	return &GaugeVec{gaugeVec}
}

//counter
type CounterOpts struct {
	Name string
	Help string
}
type Counter interface {
	prometheus.Counter
}
type CounterVec struct {
	base *prometheus.CounterVec
}

func (v *CounterVec) Describe(ch chan<- *prometheus.Desc) {
	v.base.Describe(ch)
}

func (v *CounterVec) Collect(ch chan<- prometheus.Metric) {
	v.base.Collect(ch)
}

func (v *CounterVec) WithLabelValues(lvs ...string) Counter {
	g, err := v.base.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return g
}
func NewCounterVec(opts CounterOpts, labelNames []string) *CounterVec {
	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: opts.Name,
		Help: opts.Help,
	}, labelNames)
	return &CounterVec{counterVec}
}

//histogram
type HistogramOpts struct {
	Name    string
	Help    string
	Buckets []float64
}
type Observer interface {
	prometheus.Observer
}
type HistogramVec struct {
	base *prometheus.HistogramVec
}

func (v *HistogramVec) Describe(ch chan<- *prometheus.Desc) {
	v.base.Describe(ch)
}

func (v *HistogramVec) Collect(ch chan<- prometheus.Metric) {
	v.base.Collect(ch)
}

func (v *HistogramVec) WithLabelValues(lvs ...string) Observer {
	g, err := v.base.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return g
}
func NewHistogramVec(opts HistogramOpts, labelNames []string) *HistogramVec {
	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    opts.Name,
		Help:    opts.Help,
		Buckets: opts.Buckets,
	}, labelNames)
	return &HistogramVec{histogramVec}
}

//summary
type SummaryOpts struct {
	Name       string
	Help       string
	MaxAge     time.Duration
	AgeBuckets uint32
}

type SummaryVec struct {
	base *prometheus.SummaryVec
}

func (v *SummaryVec) Describe(ch chan<- *prometheus.Desc) {
	v.base.Describe(ch)
}

func (v *SummaryVec) Collect(ch chan<- prometheus.Metric) {
	v.base.Collect(ch)
}

func (v *SummaryVec) WithLabelValues(lvs ...string) Observer {
	g, err := v.base.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return g
}
func NewSummaryVec(opts SummaryOpts, labelNames []string) *SummaryVec {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       opts.Name,
		Help:       opts.Help,
		MaxAge:     opts.MaxAge,
		AgeBuckets: opts.AgeBuckets,
	}, labelNames)
	return &SummaryVec{summaryVec}
}
