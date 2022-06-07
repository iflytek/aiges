package gauge

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type customCollector struct {
	customDesc *prometheus.Desc
}

func newCustomCollector() *customCollector {
	return &customCollector{
		customDesc: prometheus.NewDesc(
			prometheus.BuildFQName("namespace", "subsystem", "name"),
			"help",
			[]string{"tag1", "tag2"},
			nil,
		),
	}
}

func (collector *customCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.customDesc
}

var c float64

func (collector *customCollector) Collect(ch chan<- prometheus.Metric) {
	c++
	ch <- prometheus.MustNewConstMetric(
		collector.customDesc,
		prometheus.GaugeValue,
		c,
		fmt.Sprintf("tag1_%v", c), fmt.Sprintf("tag2_%v", c),
	)

}

//func main() {
//	registry := prometheus.NewRegistry()
//	if err := registry.Register(newCustomCollector()); err != nil {
//		panic(err)
//	}
//
//	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
//	log.Println("Beginning to serve on port :8080")
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}
