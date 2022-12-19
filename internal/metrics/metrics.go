package metrics

import (
	"math/rand"
	"runtime"
)

type Gauge float64
type Counter int64

type Metrics struct {
	GaugeMetrics   map[string]Gauge
	CounterMetrics map[string]Counter
}

type MetricsCollector interface {
	UpdateMetrics()
	GetMetrics() Metrics
	GetMetric(name string) (interface{}, error)
}

func UpdateMetrics(m runtime.MemStats, counter int) *Metrics {
	return &Metrics{
		GaugeMetrics: map[string]Gauge{
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
		},
		CounterMetrics: map[string]Counter{
			"PollCount":   Counter(counter),
			"RandomValue": Counter(rand.Int63()),
		},
	}
}

func NewMetrics() *Metrics {
	return &Metrics{
		GaugeMetrics: map[string]Gauge{
			"Alloc":         Gauge(0.0),
			"BuckHashSys":   Gauge(0.0),
			"Frees":         Gauge(0.0),
			"GCCPUFraction": Gauge(0.0),
			"GCSys":         Gauge(0.0),
			"HeapAlloc":     Gauge(0.0),
			"HeapIdle":      Gauge(0.0),
			"HeapInuse":     Gauge(0.0),
			"HeapObjects":   Gauge(0.0),
			"HeapReleased":  Gauge(0.0),
			"HeapSys":       Gauge(0.0),
			"LastGC":        Gauge(0.0),
			"Lookups":       Gauge(0.0),
			"MCacheInuse":   Gauge(0.0),
			"MCacheSys":     Gauge(0.0),
			"MSpanInuse":    Gauge(0.0),
			"MSpanSys":      Gauge(0.0),
			"Mallocs":       Gauge(0.0),
			"NextGC":        Gauge(0.0),
			"NumForcedGC":   Gauge(0.0),
			"NumGC":         Gauge(0.0),
			"OtherSys":      Gauge(0.0),
			"PauseTotalNs":  Gauge(0.0),
			"StackInuse":    Gauge(0.0),
			"StackSys":      Gauge(0.0),
			"Sys":           Gauge(0.0),
			"TotalAlloc":    Gauge(0.0),
		},
		CounterMetrics: map[string]Counter{
			"PollCount":   Counter(0),
			"RandomValue": Counter(rand.Int63()),
		},
	}
}

func GetMemStats() runtime.MemStats {
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)
	return memstats
}
