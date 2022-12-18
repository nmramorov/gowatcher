package metrics

import (
	"math/rand"
	"runtime"
)

type gauge float64
type counter int64

type Metrics struct {
	GaugeMetrics map[string]gauge
	CountMetrics map[string]counter
}

type MetricsCollector interface {
	CollectMetrics()
}

func NewMetrics(m runtime.MemStats) *Metrics {
	return &Metrics{
		GaugeMetrics: map[string]gauge{
			"Alloc":         gauge(m.Alloc),
			"BuckHashSys":   gauge(m.BuckHashSys),
			"Frees":         gauge(m.Frees),
			"GCCPUFraction": gauge(m.GCCPUFraction),
			"GCSys":         gauge(m.GCSys),
			"HeapAlloc":     gauge(m.HeapAlloc),
			"HeapIdle":      gauge(m.HeapIdle),
			"HeapInuse":     gauge(m.HeapInuse),
			"HeapObjects":   gauge(m.HeapObjects),
			"HeapReleased":  gauge(m.HeapReleased),
			"HeapSys":       gauge(m.HeapSys),
			"LastGC":        gauge(m.LastGC),
			"Lookups":       gauge(m.Lookups),
			"MCacheInuse":   gauge(m.MCacheInuse),
			"MCacheSys":     gauge(m.MCacheSys),
			"MSpanInuse":    gauge(m.MSpanInuse),
			"MSpanSys":      gauge(m.MSpanSys),
			"Mallocs":       gauge(m.Mallocs),
			"NextGC":        gauge(m.NextGC),
			"NumForcedGC":   gauge(m.NumForcedGC),
			"NumGC":         gauge(m.NumGC),
			"OtherSys":      gauge(m.OtherSys),
			"PauseTotalNs":  gauge(m.PauseTotalNs),
			"StackInuse":    gauge(m.StackInuse),
			"StackSys":      gauge(m.StackSys),
			"Sys":           gauge(m.Sys),
			"TotalAlloc":    gauge(m.TotalAlloc),
		},
		CountMetrics: map[string]counter{
			"PollCount":   0,
			"RandomValue": counter(rand.Int63()),
		},
	}
}

func GetMemStats() runtime.MemStats {
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)
	return memstats
}

// func main() {
// 	memst := GetMemStats()
// 	gaugeMetrics := NewMetrics(memst).GaugeMetrics
// 	fmt.Println(gaugeMetrics)
// }
