package collector

type Collector struct {
	MetricsCollector
	metrics *Metrics
}

func NewCollector() *Collector {
	return &Collector{
		metrics: NewMetrics(),
	}
}

func (col *Collector) CollectMetrics() {

}
