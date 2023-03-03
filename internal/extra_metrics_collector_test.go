package metrics

import (
	// "fmt"
	// "os"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtraMetricsCollector(t *testing.T) {
	testCollector := NewExtraMetricsCollector()
	assert.Equal(t, Gauge(0.0), testCollector.Metrics["TotalMemory"])
	assert.Equal(t, Gauge(0.0), testCollector.Metrics["FreeMemory"])
	assert.Equal(t, Gauge(0.0), testCollector.Metrics["CPUutilization1"])
}

func TestExtraMetricsCollectorUpdate(t *testing.T) {
	testCollector := NewExtraMetricsCollector()
	testCollector.Update()
	fmt.Println(testCollector.Metrics)
	assert.NotEqual(t, Gauge(0.0), testCollector.Metrics["TotalMemory"])
	assert.NotEqual(t, Gauge(0.0), testCollector.Metrics["FreeMemory"])
	assert.NotEqual(t, Gauge(0.0), testCollector.Metrics["CPUutilization1"])
}
