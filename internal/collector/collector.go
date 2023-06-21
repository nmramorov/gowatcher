package collector

import (
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	m "github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/errors"
	"github.com/nmramorov/gowatcher/internal/log"
)

type Collector struct {
	m.MetricsCollector
	Metrics *m.Metrics
	Updates int
	mu      sync.Mutex
}

func NewCollector() *Collector {
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)

	return &Collector{
		Metrics: m.NewMetrics(),
		Updates: 0,
	}
}

func NewCollectorFromSavedFile(saved *m.Metrics) *Collector {
	return &Collector{
		Metrics: saved,
		Updates: 0,
	}
}

func (col *Collector) UpdateMetrics() {
	col.mu.Lock()
	defer col.mu.Unlock()
	var newstats runtime.MemStats
	runtime.ReadMemStats(&newstats)

	col.Updates++
	col.Metrics = m.UpdateMetrics(newstats, col.Updates)
}

func (col *Collector) GetMetrics() *m.Metrics {
	return col.Metrics
}

func (col *Collector) GetMetric(name string) (interface{}, error) {
	for k, v := range col.Metrics.CounterMetrics {
		if k == name {
			return v, nil
		}
	}
	for k, v := range col.Metrics.GaugeMetrics {
		if k == name {
			return v, nil
		}
	}
	return 1, errors.ErrorMetricNotFound
}

func (col *Collector) String(value interface{}) (string, error) {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil
	case reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	default:
		return "", errors.ErrorWrongStringConvertion
	}
}

func (col *Collector) UpdateMetricFromJSON(newMetric *m.JSONMetrics) (*m.JSONMetrics, error) {
	col.mu.Lock()
	defer col.mu.Unlock()
	result := m.JSONMetrics{}
	switch newMetric.MType {
	case "gauge":
		col.Metrics.GaugeMetrics[newMetric.ID] = m.Gauge(*newMetric.Value)
		val := col.Metrics.GaugeMetrics[newMetric.ID]
		result.Value = (*float64)(&val)
	case "counter":

		col.Metrics.CounterMetrics[newMetric.ID] = col.Metrics.CounterMetrics[newMetric.ID] + m.Counter(*newMetric.Delta)
		delta := col.Metrics.CounterMetrics[newMetric.ID]
		result.Delta = (*int64)(&delta)
	default:
		return &result, errors.ErrorMetricNotFound
	}
	result.MType = newMetric.MType
	result.ID = newMetric.ID
	return &result, nil
}

func (col *Collector) GetMetricJSON(requestedMetric *m.JSONMetrics) (*m.JSONMetrics, error) {
	result := m.JSONMetrics{}
	switch requestedMetric.MType {
	case "gauge":
		res := col.Metrics.GaugeMetrics[requestedMetric.ID]
		result.Value = (*float64)(&res)

	case "counter":
		res := col.Metrics.CounterMetrics[requestedMetric.ID]
		result.Delta = (*int64)(&res)
	default:
		return requestedMetric, errors.ErrorMetricNotFound
	}
	result.MType = requestedMetric.MType
	result.ID = requestedMetric.ID

	return &result, nil
}

func (col *Collector) UpdateBatch(metrics []*m.JSONMetrics) error {
	for _, metric := range metrics {
		_, err := col.UpdateMetricFromJSON(metric)
		if err != nil {
			log.ErrorLog.Printf("could not update metric as batch part: %e", err)
			return err
		}
	}
	return nil
}

func (col *Collector) UpdateExtraMetrics() {
	col.mu.Lock()
	defer col.mu.Unlock()
	v, _ := mem.VirtualMemory()
	utilization, _ := cpu.Percent(time.Second, false)
	col.Metrics.GaugeMetrics["TotalMemory"] = m.Gauge(v.Total)
	col.Metrics.GaugeMetrics["FreeMemory"] = m.Gauge(v.Free)
	col.Metrics.GaugeMetrics["CPUutilization1"] = m.Gauge(utilization[0])
}
