package metrics

import (
	"errors"
	"runtime"
)

type Collector struct {
	MetricsCollector
	metrics *Metrics
	updates int
}

var ErrorMetricNotFound error = errors.New("no such metric")

func NewCollector() *Collector {
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)

	return &Collector{
		metrics: NewMetrics(memstats),
		updates: 0,
	}
}

func (col *Collector) UpdateMetrics() {
	var newstats runtime.MemStats
	runtime.ReadMemStats(&newstats)

	col.metrics = NewMetrics(newstats)

	col.updates++
	col.metrics.CountMetrics["PollCount"] = counter(col.updates)
}

func (col *Collector) GetMetrics() Metrics {
	return *col.metrics
}

func (col *Collector) GetMetric(name string) (interface{}, error) {
	for k, v := range col.metrics.CountMetrics {
		if k == name {
			return v, nil
		}
	}
	for k, v := range col.metrics.GaugeMetrics {
		if k == name {
			return v, nil
		}
	}
	return 1, ErrorMetricNotFound
}
