package xsf

/*
# 组件延迟分布
module_delay{sub=iat,name=ats,idc=dx,cs=3s,type=max} 100
module_delay{sub=iat,name=ats,idc=dx,cs=3s,type=min} 100
module_delay{sub=iat,name=ats,idc=dx,cs=3s,type=avg} 100

# 组件错误分布
module_request{sub=iat,name=ats,idc=dx,cs=s,type=passive,code=10010} 80
module_request{sub=iat,name=ats,idc=dx,cs=s,type=passive,code=0} 20

# 实时引擎路数
ats_auth{sub=iat,name=ats,idc=dx,cs=5s,type=max,ent=sms-en} 100
ats_auth{sub=iat,name=ats,idc=dx,cs=5s,type=idle,ent=sms-en} 100

# 自定义部分
ats_auth{sub=iat,name=ats,idc=dx,cs=5s} 100
ats_auth{sub=iat,name=ats,idc=dx,cs=5s} 100

*/
import (
	"bytes"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MetricsTimeUnit        = "timeUnit"
	defaultMetricsTimeUnit = time.Second //todo 后续优化
)

var (
	//此处切片主要是考虑到可能出现的多server情况，存在数据冲突问题
	slidingErrCodeWindows []*errCodeWindow
	slidingDelayWindows   []*delayWindow
	ConcurrentStatistics  int64
	slidingRwMu           sync.RWMutex
)
var (
	MetricsInvalidParams = fmt.Errorf("invalid params")
	MetricsRepeatInit    = fmt.Errorf("metrics repeat initialization")
)

func incrConcurrentStatistics() {
	atomic.AddInt64(&ConcurrentStatistics, 1)
}
func decrConcurrentStatistics() {
	atomic.AddInt64(&ConcurrentStatistics, -1)
}
func AddSlidingErrCodeWindows(item *errCodeWindow) {
	slidingRwMu.Lock()
	defer slidingRwMu.Unlock()
	slidingErrCodeWindows = append(slidingErrCodeWindows, item)
}
func AddSlidingDelayWindows(item *delayWindow) {
	slidingRwMu.Lock()
	defer slidingRwMu.Unlock()
	slidingDelayWindows = append(slidingDelayWindows, item)
}

var (
	metricsLabelCollection = []metricsLabel{
		{
			MetricsName:   string(DELAY),
			MetricsHelp:   "module_delay",
			MetricsLabels: []string{"sub", "name", "idc", "cs", "type"},
		},
		{
			MetricsName:   string(CONCURRENT),
			MetricsHelp:   "concurrent_statistics",
			MetricsLabels: []string{"sub", "name", "idc", "cs"},
		},
		{
			MetricsName:   string(REQUEST),
			MetricsHelp:   "module_request",
			MetricsLabels: []string{"sub", "name", "idc", "cs", "type", "code"},
		},
		{
			MetricsName:   string(AUTH),
			MetricsHelp:   "eng_auth",
			MetricsLabels: []string{"sub", "name", "idc", "cs", "type", "ent"},
		},
	}

	registryInst Registry
)

type metricsClass string

const (
	//限制用户随便传
	DELAY      metricsClass = "module_delay"
	REQUEST    metricsClass = "module_request"
	AUTH       metricsClass = "eng_auth"
	CONCURRENT metricsClass = "concurrent_statistics"
)

type metricsLabel struct {
	MetricsName   string
	MetricsHelp   string
	MetricsLabels []string
}

//采集器
type Collector interface {
	prometheus.Collector
}

type Registry struct {
	sync.RWMutex
	ready        bool
	registry     *prometheus.Registry
	collectorMap map[string]Collector

	sub, svcName, idc, cs string
}

func (r *Registry) init() error {
	if r.ready {
		loggerStd.Println("metrics repeat initialization")
		return MetricsRepeatInit
	} else {
		r.ready = true
	}

	r.registry = prometheus.NewRegistry()
	r.collectorMap = make(map[string]Collector)
	return nil
}
func (r *Registry) initEx(sub, name, idc, cs string, opts ...MetricsOpt) error {
	if r.ready {
		loggerStd.Println("metrics repeat initialization")
		return MetricsRepeatInit
	} else {
		r.ready = true
	}

	r.Lock()
	defer r.Unlock()

	r.svcName = name
	r.idc = idc
	r.sub = sub
	r.cs = cs

	r.registry = prometheus.NewRegistry()
	r.collectorMap = make(map[string]Collector)

	for _, metricsLabelInst := range metricsLabelCollection {
		gaugeVecInst := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsLabelInst.MetricsName,
			Help: metricsLabelInst.MetricsHelp,
		}, metricsLabelInst.MetricsLabels)

		r.collectorMap[metricsLabelInst.MetricsName] = gaugeVecInst
		r.registry.MustRegister(gaugeVecInst)
	}

	return nil
}

func (r *Registry) withLabelValuesItem(name string, lvs []string, val float64) {
	if !r.ready {
		return
	}

	r.RLock()
	defer r.RUnlock()

	r.collectorMap[name].(*prometheus.GaugeVec).WithLabelValues(lvs...).Set(val)
}
func (r *Registry) WithLabelValues(val float64, opts ...MetricsOpt) {
	r.withLabelValues(val, opts...)
}
func (r *Registry) withLabelValues(val float64, opts ...MetricsOpt) {
	if !r.ready {
		return
	}

	labelValInst := metricsConfigure{}
	for _, opt := range opts {
		opt(&labelValInst)
	}

	switch labelValInst.metricsClass {
	case DELAY:
		{
			r.withLabelValuesItem(
				string(DELAY),
				[]string{
					r.sub,
					r.svcName,
					r.idc,
					r.cs,
					labelValInst.metricsType,
				},
				val,
			)
		}
	case REQUEST:
		{
			r.withLabelValuesItem(
				string(REQUEST),
				[]string{
					r.sub,
					r.svcName,
					r.idc,
					r.cs,
					labelValInst.metricsType,
					labelValInst.metricsCode,
				},
				val,
			)
		}
	case AUTH:
		{
			r.withLabelValuesItem(
				string(AUTH),
				[]string{
					r.sub,
					r.svcName,
					r.idc,
					r.cs,
					labelValInst.metricsType,
					labelValInst.metricsEnt,
				},
				val,
			)
		}
	case CONCURRENT:
		{
			r.withLabelValuesItem(
				string(CONCURRENT),
				[]string{
					r.sub,
					r.svcName,
					r.idc,
					r.cs,
				},
				val,
			)
		}
	default:
		{
			panic("param err") //todo 2019-04-10 15:19:51
		}
	}
}

type metricsConfigure struct {
	metricsClass metricsClass
	metricsType  string
	metricsCode  string
	metricsEnt   string

	name string
	idc  string
	sub  string
	cs   string
}

func newRegistry() (*Registry, error) {
	registryInst := &Registry{}

	return registryInst, registryInst.init()
}
func newRegistryEx(sub, name, idc, cs string, opts ...MetricsOpt) (*Registry, error) {
	validCheck := func(param string) bool {
		if "" == param {
			return false
		}
		return true
	}
	if !(validCheck(sub) && validCheck(name) && validCheck(idc) && validCheck(cs)) {
		return nil, MetricsInvalidParams
	}
	registryInst := &Registry{}

	return registryInst, registryInst.initEx(sub, name, idc, cs, opts...)
}

//in: name&collector
func Register(n string, c Collector) error {
	return registryInst.Register(n, c)
}
func (r *Registry) Register(n string, c Collector) error {
	r.Lock()
	defer r.Unlock()

	if !r.ready {
		loggerStd.Printf("fn:Register,metrics is disable")
		return nil
	}

	if len(n) == 1 || c == nil {
		return fmt.Errorf("MetricsInvalidParams")
	}
	r.collectorMap[n] = c
	r.registry.MustRegister(c)
	return nil
}
func (r *Registry) Gather() ([]*dto.MetricFamily, error) {
	return r.registry.Gather()
}

func UnregisterAll() bool {
	return registryInst.UnregisterAll()
}
func (r *Registry) UnregisterAll() bool {
	r.Lock()
	defer r.Unlock()

	rst := false
	for k, v := range r.collectorMap {
		delete(r.collectorMap, k)
		rst = r.registry.Unregister(v)
		if !rst {
			return rst
		}
	}
	return rst
}

//name
func Unregister(n string) bool {
	return registryInst.Unregister(n)
}
func (r *Registry) Unregister(n string) bool {
	r.Lock()
	defer r.Unlock()

	collector, collectorOk := r.collectorMap[n]
	if !collectorOk {
		return false
	} else {
		delete(r.collectorMap, n)
	}

	return r.registry.Unregister(collector)
}
func (r *Registry) GetHttpHandler() http.Handler {
	//todo not complete
	panic("illegal")
}
func (r *Registry) toBuffer() *bytes.Buffer {
	if !r.ready {
		return nil
	}

	checkErr := func(err error) {
		if nil != err {
			panic(err)
		}
	}
	w := &bytes.Buffer{}

	enc := expfmt.NewEncoder(w, expfmt.FmtText)
	dataSet, dataSetErr := r.registry.Gather()
	checkErr(dataSetErr)

	for _, dataItem := range dataSet {
		checkErr(enc.Encode(dataItem))
	}
	return w
}
func (r *Registry) toBytes() []byte {
	if !r.ready {
		return nil
	}

	return r.toBuffer().Bytes()
}
func (r *Registry) toString() string {
	if !r.ready {
		return ""
	}

	return r.toBuffer().String()
}
func (r *Registry) String() string {
	if !r.ready {
		return ""
	}

	return r.toString()
}

func NewRegistry(opts ...MetricsOpt) (*Registry, error) {
	if len(opts) != 0 {
		validCheck := func(param string) bool {
			if "" == param {
				return false
			}
			return true
		}
		metricsConfigureInst := &metricsConfigure{}
		for _, opt := range opts {
			opt(metricsConfigureInst)
		}

		if !(validCheck(metricsConfigureInst.sub) && validCheck(metricsConfigureInst.name) && validCheck(metricsConfigureInst.idc) && validCheck(metricsConfigureInst.cs)) {
			return nil, MetricsInvalidParams
		}
		return newRegistryEx(metricsConfigureInst.sub, metricsConfigureInst.name, metricsConfigureInst.idc, metricsConfigureInst.cs, opts...)
	}
	return newRegistry()
}

type MetricsOpt func(in *metricsConfigure)

func WithMetricsClass(class metricsClass) MetricsOpt {
	return func(in *metricsConfigure) {
		in.metricsClass = class
	}
}
func WithMetricsType(_type string) MetricsOpt {
	return func(in *metricsConfigure) {
		in.metricsType = _type
	}
}
func WithMetricsCode(code string) MetricsOpt {
	return func(in *metricsConfigure) {
		in.metricsCode = code
	}
}
func WithMetricsSub(sub string) MetricsOpt {
	return func(in *metricsConfigure) {
		in.sub = sub
	}
}
func WithMetricsName(name string) MetricsOpt {
	return func(in *metricsConfigure) {
		in.name = name
	}
}
func WithMetricsIdc(idc string) MetricsOpt {
	return func(in *metricsConfigure) {
		in.idc = idc
	}
}
func WithMetricsCs(cs string) MetricsOpt {
	return func(in *metricsConfigure) {
		in.cs = cs
	}
}
