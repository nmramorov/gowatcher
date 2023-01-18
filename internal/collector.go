package metrics

import (
	"reflect"
	"runtime"
	"strconv"
)

type Collector struct {
	MetricsCollector
	metrics *Metrics
	updates int
}

func NewCollector() *Collector {
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)

	return &Collector{
		metrics: NewMetrics(),
		updates: 0,
	}
}

func (col *Collector) UpdateMetrics() {
	var newstats runtime.MemStats
	runtime.ReadMemStats(&newstats)

	col.updates++
	col.metrics = UpdateMetrics(newstats, col.updates)
}

func (col *Collector) GetMetrics() *Metrics {
	return col.metrics
}

func (col *Collector) GetMetric(name string) (interface{}, error) {
	for k, v := range col.metrics.CounterMetrics {
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

func (col *Collector) String(value interface{}) (string, error) {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil
	case reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	default:
		return "", ErrorWrongStringConvertion
	}
}

func (col *Collector) UpdateMetricFromJson(newMetric *JSONMetrics) (*JSONMetrics, error) {
	switch newMetric.MType {
	case "gauge":
		col.metrics.GaugeMetrics[newMetric.ID] = Gauge(*newMetric.Value)
	case "counter":
		col.metrics.CounterMetrics[newMetric.ID] = col.metrics.CounterMetrics[newMetric.ID] + Counter(*newMetric.Delta)
	default:
		return newMetric, ErrorMetricNotFound
	}
	return newMetric, nil
}

func (col *Collector) GetMetricJson(requestedMetric *JSONMetrics) (*JSONMetrics, error) {
	switch requestedMetric.MType {
	case "gauge":
		res := col.metrics.GaugeMetrics[requestedMetric.ID]
		requestedMetric.Value = (*float64)(&res)

	case "counter":
		res := col.metrics.CounterMetrics[requestedMetric.ID]
		requestedMetric.Delta = (*int64)(&res)
	default:
		return requestedMetric, ErrorMetricNotFound
	}
	return requestedMetric, nil
}
