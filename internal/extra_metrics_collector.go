package metrics

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type ExtraMetricsCollector struct {
	Metrics map[string]Gauge
}

func NewExtraMetricsCollector() *ExtraMetricsCollector {
	return &ExtraMetricsCollector{
		Metrics: map[string]Gauge{
			"TotalMemory":     Gauge(0),
			"FreeMemory":      Gauge(0),
			"CPUutilization1": Gauge(0),
		},
	}
}

func (col *ExtraMetricsCollector) Update() {
	v, _ := mem.VirtualMemory()
	utilization, _ := cpu.Percent(time.Second, false)
	col.Metrics["TotalMemory"] = Gauge(v.Total)
	col.Metrics["FreeMemory"] = Gauge(v.Free)
	col.Metrics["CPUutilization1"] = Gauge(utilization[0])
}
