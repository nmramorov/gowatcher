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

func NewCollectorFromSavedFile(saved *Metrics) *Collector {
	return &Collector{
		metrics: saved,
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
	result := JSONMetrics{}
	switch newMetric.MType {
	case "gauge":
		col.metrics.GaugeMetrics[newMetric.ID] = Gauge(*newMetric.Value)
		val := col.metrics.GaugeMetrics[newMetric.ID]
		result.Value = (*float64)(&val)
	case "counter":

		col.metrics.CounterMetrics[newMetric.ID] = col.metrics.CounterMetrics[newMetric.ID] + Counter(*newMetric.Delta)
		delta := col.metrics.CounterMetrics[newMetric.ID]
		result.Delta = (*int64)(&delta)
	default:
		return &result, ErrorMetricNotFound
	}
	result.MType = newMetric.MType
	result.ID = newMetric.ID
	return &result, nil
}

func (col *Collector) GetMetricJson(requestedMetric *JSONMetrics) (*JSONMetrics, error) {
	result := JSONMetrics{}
	switch requestedMetric.MType {
	case "gauge":
		res := col.metrics.GaugeMetrics[requestedMetric.ID]
		result.Value = (*float64)(&res)

	case "counter":
		res := col.metrics.CounterMetrics[requestedMetric.ID]
		result.Delta = (*int64)(&res)
	default:
		return requestedMetric, ErrorMetricNotFound
	}
	result.MType = requestedMetric.MType
	result.ID = requestedMetric.ID

	return &result, nil
}

func (col *Collector) UpdateBatch(metrics []*JSONMetrics) error {
	for _, metric := range metrics {
		_, err := col.UpdateMetricFromJson(metric)
		if err != nil {
			ErrorLog.Printf("could not update metric as batch part: %e", err)
			return err
		}
	}
	return nil
}
