package collector

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nmramorov/gowatcher/internal/collector/metrics"
)

func TestExtraMetricsCollector(t *testing.T) {
	testCollector := NewExtraMetricsCollector()
	assert.Equal(t, metrics.Gauge(0.0), testCollector.Metrics["TotalMemory"])
	assert.Equal(t, metrics.Gauge(0.0), testCollector.Metrics["FreeMemory"])
	assert.Equal(t, metrics.Gauge(0.0), testCollector.Metrics["CPUutilization1"])
}

func TestExtraMetricsCollectorUpdate(t *testing.T) {
	testCollector := NewExtraMetricsCollector()
	testCollector.Update()
	fmt.Println(testCollector.Metrics)
	assert.NotEqual(t, metrics.Gauge(0.0), testCollector.Metrics["TotalMemory"])
	assert.NotEqual(t, metrics.Gauge(0.0), testCollector.Metrics["FreeMemory"])
	assert.NotEqual(t, metrics.Gauge(0.0), testCollector.Metrics["CPUutilization1"])
}
