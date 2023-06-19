package metrics

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetrics(t *testing.T) {
	newMetr := NewMetrics()
	for key, metric := range newMetr.GaugeMetrics {
		if strings.Compare(key, "RandomValue") == 0 {
			assert.InDelta(t, 0.0, float64(metric), 1.0)
		} else {
			assert.Equal(t, metric, Gauge(0.0))
		}
	}
	for _, metric := range newMetr.CounterMetrics {
		assert.Equal(t, metric, Counter(0))
	}
}

func TestUpdateMetrics(t *testing.T) {
	m := GetMemStats()
	updatedMetrics := UpdateMetrics(m, 1)
	assert.Equal(t, Counter(1), updatedMetrics.CounterMetrics["PollCount"])

	expectedRes := map[string]Gauge{
		"Alloc":         Gauge(m.Alloc),
		"BuckHashSys":   Gauge(m.BuckHashSys),
		"Frees":         Gauge(m.Frees),
		"GCCPUFraction": Gauge(m.GCCPUFraction),
		"GCSys":         Gauge(m.GCSys),
		"HeapAlloc":     Gauge(m.HeapAlloc),
		"HeapIdle":      Gauge(m.HeapIdle),
		"HeapInuse":     Gauge(m.HeapInuse),
		"HeapObjects":   Gauge(m.HeapObjects),
		"HeapReleased":  Gauge(m.HeapReleased),
		"HeapSys":       Gauge(m.HeapSys),
		"LastGC":        Gauge(m.LastGC),
		"Lookups":       Gauge(m.Lookups),
		"MCacheInuse":   Gauge(m.MCacheInuse),
		"MCacheSys":     Gauge(m.MCacheSys),
		"MSpanInuse":    Gauge(m.MSpanInuse),
		"MSpanSys":      Gauge(m.MSpanSys),
		"Mallocs":       Gauge(m.Mallocs),
		"NextGC":        Gauge(m.NextGC),
		"NumForcedGC":   Gauge(m.NumForcedGC),
		"NumGC":         Gauge(m.NumGC),
		"OtherSys":      Gauge(m.OtherSys),
		"PauseTotalNs":  Gauge(m.PauseTotalNs),
		"StackInuse":    Gauge(m.StackInuse),
		"StackSys":      Gauge(m.StackSys),
		"Sys":           Gauge(m.Sys),
		"TotalAlloc":    Gauge(m.TotalAlloc),
	}
	// Iterating over Gauge metrics because there is "RandomValue" metric
	for k, v := range expectedRes {
		assert.Equal(t, v, updatedMetrics.GaugeMetrics[k])
	}
	assert.InDelta(t, 0.0, float64(updatedMetrics.GaugeMetrics["RandomValue"]), 1.0)
}
