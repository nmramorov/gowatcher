package collector

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/nmramorov/gowatcher/internal/collector/metrics"
)

type ExtraMetricsCollector struct {
	Metrics map[string]metrics.Gauge
}

func NewExtraMetricsCollector() *ExtraMetricsCollector {
	return &ExtraMetricsCollector{
		Metrics: map[string]metrics.Gauge{
			"TotalMemory":     metrics.Gauge(0),
			"FreeMemory":      metrics.Gauge(0),
			"CPUutilization1": metrics.Gauge(0),
		},
	}
}

func (col *ExtraMetricsCollector) Update() {
	v, _ := mem.VirtualMemory()
	utilization, _ := cpu.Percent(time.Second, false)
	col.Metrics["TotalMemory"] = metrics.Gauge(v.Total)
	col.Metrics["FreeMemory"] = metrics.Gauge(v.Free)
	col.Metrics["CPUutilization1"] = metrics.Gauge(utilization[0])
}
