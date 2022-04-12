package gauge

import (
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"math"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

type selfCollector struct {
	self prometheus.Metric
}

func (c *selfCollector) valid() bool {
	gaugeInst, gaugeInstOk := c.self.(*gauge)
	if !gaugeInstOk {
		return false
	}
	if !atomic.CompareAndSwapInt64(&gaugeInst.valid, 1, 0) {
		return false
	}
	return true
}

func (c *selfCollector) init(self prometheus.Metric) {
	c.self = self
}

func (c *selfCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.self.Desc()
}

func (c *selfCollector) Collect(ch chan<- prometheus.Metric) {
	if !c.valid() {
		return
	}
	ch <- c.self
}
func NewGauge(opts prometheus.GaugeOpts) prometheus.Gauge {
	desc := prometheus.NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	)
	result := &gauge{desc: desc, labelPairs: getConstLabelPairsFromDesc(desc)}
	result.init(result) // Init self-collection.
	return result
}

type gauge struct {
	valBits uint64

	selfCollector

	desc       *prometheus.Desc
	labelPairs []*dto.LabelPair
	valid      int64
}

func (g *gauge) Desc() *prometheus.Desc {
	return g.desc
}

func (g *gauge) Set(val float64) {
	atomic.StoreInt64(&g.valid, 1)
	atomic.StoreUint64(&g.valBits, math.Float64bits(val))
}

func (g *gauge) SetToCurrentTime() {
	g.Set(float64(time.Now().UnixNano()) / 1e9)
}

func (g *gauge) Inc() {
	g.Add(1)
}

func (g *gauge) Dec() {
	g.Add(-1)
}

func (g *gauge) Add(val float64) {
	for {
		oldBits := atomic.LoadUint64(&g.valBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + val)
		if atomic.CompareAndSwapUint64(&g.valBits, oldBits, newBits) {
			return
		}
	}
}

func (g *gauge) Sub(val float64) {
	g.Add(val * -1)
}

func (g *gauge) Write(out *dto.Metric) error {
	val := math.Float64frombits(atomic.LoadUint64(&g.valBits))
	out.Label = g.labelPairs
	out.Gauge = &dto.Gauge{Value: proto.Float64(val)}
	return nil
}

type GaugeVec struct {
	*metricVec
}

type labelPairSorter []*dto.LabelPair

func (s labelPairSorter) Len() int {
	return len(s)
}

func (s labelPairSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s labelPairSorter) Less(i, j int) bool {
	return s[i].GetName() < s[j].GetName()
}

func makeLabelPairs(desc *prometheus.Desc, labelValues []string) []*dto.LabelPair {
	var labelPairs []*dto.LabelPair
	for i, n := range getVariableLabelsFromDesc(desc) {
		labelPairs = append(labelPairs, &dto.LabelPair{
			Name:  proto.String(n),
			Value: proto.String(labelValues[i]),
		})
	}
	labelPairs = append(labelPairs, getConstLabelPairsFromDesc(desc)...)
	sort.Sort(labelPairSorter(labelPairs))
	return labelPairs
}
func BuildFQName(namespace, subsystem, name string) string {
	if name == "" {
		return ""
	}
	switch {
	case namespace != "" && subsystem != "":
		return strings.Join([]string{namespace, subsystem, name}, "_")
	case namespace != "":
		return strings.Join([]string{namespace, name}, "_")
	case subsystem != "":
		return strings.Join([]string{subsystem, name}, "_")
	}
	return name
}
func NewGaugeVec(opts prometheus.GaugeOpts, labelNames []string) *GaugeVec {
	desc := prometheus.NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		labelNames,
		opts.ConstLabels,
	)
	return &GaugeVec{
		metricVec: newMetricVec(desc, func(lvs ...string) prometheus.Metric {
			if len(lvs) != len(getVariableLabelsFromDesc(desc)) {
				panic(makeInconsistentCardinalityError(getFqNameFromDesc(desc), getVariableLabelsFromDesc(desc), lvs))
			}
			result := &gauge{desc: desc, labelPairs: makeLabelPairs(desc, lvs)}
			result.init(result) // Init self-collection.
			return result
		}),
	}
}

func (v *GaugeVec) GetMetricWithLabelValues(lvs ...string) (prometheus.Gauge, error) {
	metric, err := v.metricVec.getMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(prometheus.Gauge), err
	}
	return nil, err
}

func (v *GaugeVec) GetMetricWith(labels Labels) (prometheus.Gauge, error) {
	metric, err := v.metricVec.getMetricWith(labels)
	if metric != nil {
		return metric.(prometheus.Gauge), err
	}
	return nil, err
}

func (v *GaugeVec) WithLabelValues(lvs ...string) prometheus.Gauge {
	g, err := v.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return g
}

func (v *GaugeVec) With(labels Labels) prometheus.Gauge {
	g, err := v.GetMetricWith(labels)
	if err != nil {
		panic(err)
	}
	return g
}

func (v *GaugeVec) CurryWith(labels Labels) (*GaugeVec, error) {
	vec, err := v.curryWith(labels)
	if vec != nil {
		return &GaugeVec{vec}, err
	}
	return nil, err
}

func (v *GaugeVec) MustCurryWith(labels Labels) *GaugeVec {
	vec, err := v.CurryWith(labels)
	if err != nil {
		panic(err)
	}
	return vec
}
