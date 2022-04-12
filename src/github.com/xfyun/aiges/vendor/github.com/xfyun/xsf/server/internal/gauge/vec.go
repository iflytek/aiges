package gauge

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"sync/atomic"

	"github.com/prometheus/common/model"
)

type metricVec struct {
	*metricMap

	curry []curriedLabelValue

	hashAdd     func(h uint64, s string) uint64
	hashAddByte func(h uint64, b byte) uint64
}

func newMetricVec(desc *prometheus.Desc, newMetric func(lvs ...string) prometheus.Metric) *metricVec {
	return &metricVec{
		metricMap: &metricMap{
			metrics:   map[uint64][]metricWithLabelValues{},
			desc:      desc,
			newMetric: newMetric,
		},
		hashAdd:     hashAdd,
		hashAddByte: hashAddByte,
	}
}

func (m *metricVec) DeleteLabelValues(lvs ...string) bool {
	h, err := m.hashLabelValues(lvs)
	if err != nil {
		return false
	}

	return m.metricMap.deleteByHashWithLabelValues(h, lvs, m.curry)
}

func (m *metricVec) Delete(labels Labels) bool {
	h, err := m.hashLabels(labels)
	if err != nil {
		return false
	}

	return m.metricMap.deleteByHashWithLabels(h, labels, m.curry)
}

func (m *metricVec) curryWith(labels Labels) (*metricVec, error) {
	var (
		newCurry []curriedLabelValue
		oldCurry = m.curry
		iCurry   int
	)
	for i, label := range getVariableLabelsFromDesc(m.desc) {
		val, ok := labels[label]
		if iCurry < len(oldCurry) && oldCurry[iCurry].index == i {
			if ok {
				return nil, fmt.Errorf("label name %q is already curried", label)
			}
			newCurry = append(newCurry, oldCurry[iCurry])
			iCurry++
		} else {
			if !ok {
				continue // Label stays uncurried.
			}
			newCurry = append(newCurry, curriedLabelValue{i, val})
		}
	}
	if l := len(oldCurry) + len(labels) - len(newCurry); l > 0 {
		return nil, fmt.Errorf("%d unknown label(s) found during currying", l)
	}

	return &metricVec{
		metricMap:   m.metricMap,
		curry:       newCurry,
		hashAdd:     m.hashAdd,
		hashAddByte: m.hashAddByte,
	}, nil
}

func (m *metricVec) getMetricWithLabelValues(lvs ...string) (prometheus.Metric, error) {
	h, err := m.hashLabelValues(lvs)
	if err != nil {
		return nil, err
	}

	return m.metricMap.getOrCreateMetricWithLabelValues(h, lvs, m.curry), nil
}

func (m *metricVec) getMetricWith(labels Labels) (prometheus.Metric, error) {
	h, err := m.hashLabels(labels)
	if err != nil {
		return nil, err
	}

	return m.metricMap.getOrCreateMetricWithLabels(h, labels, m.curry), nil
}

func (m *metricVec) hashLabelValues(vals []string) (uint64, error) {
	if err := validateLabelValues(vals, len(getVariableLabelsFromDesc(m.desc))-len(m.curry)); err != nil {
		return 0, err
	}

	var (
		h             = hashNew()
		curry         = m.curry
		iVals, iCurry int
	)
	for i := 0; i < len(getVariableLabelsFromDesc(m.desc)); i++ {
		if iCurry < len(curry) && curry[iCurry].index == i {
			h = m.hashAdd(h, curry[iCurry].value)
			iCurry++
		} else {
			h = m.hashAdd(h, vals[iVals])
			iVals++
		}
		h = m.hashAddByte(h, model.SeparatorByte)
	}
	return h, nil
}

func (m *metricVec) hashLabels(labels Labels) (uint64, error) {
	if err := validateValuesInLabels(labels, len(getVariableLabelsFromDesc(m.desc))-len(m.curry)); err != nil {
		return 0, err
	}

	var (
		h      = hashNew()
		curry  = m.curry
		iCurry int
	)
	for i, label := range getVariableLabelsFromDesc(m.desc) {
		val, ok := labels[label]
		if iCurry < len(curry) && curry[iCurry].index == i {
			if ok {
				return 0, fmt.Errorf("label name %q is already curried", label)
			}
			h = m.hashAdd(h, curry[iCurry].value)
			iCurry++
		} else {
			if !ok {
				return 0, fmt.Errorf("label name %q missing in label map", label)
			}
			h = m.hashAdd(h, val)
		}
		h = m.hashAddByte(h, model.SeparatorByte)
	}
	return h, nil
}

type metricWithLabelValues struct {
	values []string
	metric prometheus.Metric
}

type curriedLabelValue struct {
	index int
	value string
}

type metricMap struct {
	mtx       sync.RWMutex // Protects metrics.
	metrics   map[uint64][]metricWithLabelValues
	desc      *prometheus.Desc
	newMetric func(labelValues ...string) prometheus.Metric
}

func (m *metricMap) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.desc
}
func (m *metricMap) valid(metrics prometheus.Metric) bool {
	gaugeInst, gaugeInstOk := metrics.(*gauge)
	if !gaugeInstOk {
		return true
	}
	if !atomic.CompareAndSwapInt64(&gaugeInst.valid, 1, 0) {
		return false
	}
	return true
}
func (m *metricMap) Collect(ch chan<- prometheus.Metric) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	for _, metrics := range m.metrics {
		for _, metric := range metrics {
			if m.valid(metric.metric) {
				ch <- metric.metric
			}
		}
	}
}

func (m *metricMap) Reset() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for h := range m.metrics {
		delete(m.metrics, h)
	}
}

func (m *metricMap) deleteByHashWithLabelValues(
	h uint64, lvs []string, curry []curriedLabelValue,
) bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	metrics, ok := m.metrics[h]
	if !ok {
		return false
	}

	i := findMetricWithLabelValues(metrics, lvs, curry)
	if i >= len(metrics) {
		return false
	}

	if len(metrics) > 1 {
		m.metrics[h] = append(metrics[:i], metrics[i+1:]...)
	} else {
		delete(m.metrics, h)
	}
	return true
}

func (m *metricMap) deleteByHashWithLabels(
	h uint64, labels Labels, curry []curriedLabelValue,
) bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	metrics, ok := m.metrics[h]
	if !ok {
		return false
	}
	i := findMetricWithLabels(m.desc, metrics, labels, curry)
	if i >= len(metrics) {
		return false
	}

	if len(metrics) > 1 {
		m.metrics[h] = append(metrics[:i], metrics[i+1:]...)
	} else {
		delete(m.metrics, h)
	}
	return true
}

func (m *metricMap) getOrCreateMetricWithLabelValues(
	hash uint64, lvs []string, curry []curriedLabelValue,
) prometheus.Metric {
	m.mtx.RLock()
	metric, ok := m.getMetricWithHashAndLabelValues(hash, lvs, curry)
	m.mtx.RUnlock()
	if ok {
		return metric
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	metric, ok = m.getMetricWithHashAndLabelValues(hash, lvs, curry)
	if !ok {
		inlinedLVs := inlineLabelValues(lvs, curry)
		metric = m.newMetric(inlinedLVs...)
		m.metrics[hash] = append(m.metrics[hash], metricWithLabelValues{values: inlinedLVs, metric: metric})
	}
	return metric
}

func (m *metricMap) getOrCreateMetricWithLabels(
	hash uint64, labels Labels, curry []curriedLabelValue,
) prometheus.Metric {
	m.mtx.RLock()
	metric, ok := m.getMetricWithHashAndLabels(hash, labels, curry)
	m.mtx.RUnlock()
	if ok {
		return metric
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	metric, ok = m.getMetricWithHashAndLabels(hash, labels, curry)
	if !ok {
		lvs := extractLabelValues(m.desc, labels, curry)
		metric = m.newMetric(lvs...)
		m.metrics[hash] = append(m.metrics[hash], metricWithLabelValues{values: lvs, metric: metric})
	}
	return metric
}

func (m *metricMap) getMetricWithHashAndLabelValues(
	h uint64, lvs []string, curry []curriedLabelValue,
) (prometheus.Metric, bool) {
	metrics, ok := m.metrics[h]
	if ok {
		if i := findMetricWithLabelValues(metrics, lvs, curry); i < len(metrics) {
			return metrics[i].metric, true
		}
	}
	return nil, false
}

func (m *metricMap) getMetricWithHashAndLabels(
	h uint64, labels Labels, curry []curriedLabelValue,
) (prometheus.Metric, bool) {
	metrics, ok := m.metrics[h]
	if ok {
		if i := findMetricWithLabels(m.desc, metrics, labels, curry); i < len(metrics) {
			return metrics[i].metric, true
		}
	}
	return nil, false
}

func findMetricWithLabelValues(
	metrics []metricWithLabelValues, lvs []string, curry []curriedLabelValue,
) int {
	for i, metric := range metrics {
		if matchLabelValues(metric.values, lvs, curry) {
			return i
		}
	}
	return len(metrics)
}

func findMetricWithLabels(
	desc *prometheus.Desc, metrics []metricWithLabelValues, labels Labels, curry []curriedLabelValue,
) int {
	for i, metric := range metrics {
		if matchLabels(desc, metric.values, labels, curry) {
			return i
		}
	}
	return len(metrics)
}

func matchLabelValues(values []string, lvs []string, curry []curriedLabelValue) bool {
	if len(values) != len(lvs)+len(curry) {
		return false
	}
	var iLVs, iCurry int
	for i, v := range values {
		if iCurry < len(curry) && curry[iCurry].index == i {
			if v != curry[iCurry].value {
				return false
			}
			iCurry++
			continue
		}
		if v != lvs[iLVs] {
			return false
		}
		iLVs++
	}
	return true
}

func matchLabels(desc *prometheus.Desc, values []string, labels Labels, curry []curriedLabelValue) bool {
	if len(values) != len(labels)+len(curry) {
		return false
	}
	iCurry := 0
	for i, k := range getVariableLabelsFromDesc(desc) {
		if iCurry < len(curry) && curry[iCurry].index == i {
			if values[i] != curry[iCurry].value {
				return false
			}
			iCurry++
			continue
		}
		if values[i] != labels[k] {
			return false
		}
	}
	return true
}

func extractLabelValues(desc *prometheus.Desc, labels Labels, curry []curriedLabelValue) []string {
	labelValues := make([]string, len(labels)+len(curry))
	iCurry := 0
	for i, k := range getVariableLabelsFromDesc(desc) {
		if iCurry < len(curry) && curry[iCurry].index == i {
			labelValues[i] = curry[iCurry].value
			iCurry++
			continue
		}
		labelValues[i] = labels[k]
	}
	return labelValues
}

func inlineLabelValues(lvs []string, curry []curriedLabelValue) []string {
	labelValues := make([]string, len(lvs)+len(curry))
	var iCurry, iLVs int
	for i := range labelValues {
		if iCurry < len(curry) && curry[iCurry].index == i {
			labelValues[i] = curry[iCurry].value
			iCurry++
			continue
		}
		labelValues[i] = lvs[iLVs]
		iLVs++
	}
	return labelValues
}
