package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetrics(t *testing.T) {
	var tests = []struct {
		name  string
		value gauge
	}{
		{name: "Alloc", value: 0.0},
		{name: "BuckHashSys", value: 0.0},
		{name: "Frees", value: 0.0},
		{name: "GCCPUFraction", value: 0.0},
		{name: "GCSys", value: 0.0},
		{name: "HeapAlloc", value: 0.0},
		{name: "HeapIdle", value: 0.0},
		{name: "HeapInuse", value: 0.0},
		{name: "HeapInuse", value: 0.0},
		{name: "HeapObjects", value: 0.0},
		{name: "HeapReleased", value: 0.0},
		{name: "HeapSys", value: 0.0},
		{name: "LastGC", value: 0.0},
		{name: "Lookups", value: 0.0},
		{name: "MCacheInuse", value: 0.0},
		{name: "MCacheSys", value: 0.0},
		{name: "MSpanInuse", value: 0.0},
		{name: "MSpanSys", value: 0.0},
		{name: "Mallocs", value: 0.0},
		{name: "NextGC", value: 0.0},
		{name: "NumForcedGC", value: 0.0},
		{name: "NumGC", value: 0.0},
		{name: "OtherSys", value: 0.0},
		{name: "PauseTotalNs", value: 0.0},
		{name: "StackInuse", value: 0.0},
		{name: "StackSys", value: 0.0},
		{name: "Sys", value: 0.0},
		{name: "TotalAlloc", value: 0.0},
	}

	gaugeMetrics := Metrics{}.GaugeMetrics

	for _, v := range tests {
		assert.Equal(t, gaugeMetrics[v.name], v.value)
	}
}
