package metrics

func validateMetricType(metricType string) bool {
	switch metricType {
	case "gauge":
		return true
	case "count":
		return true
	default:
		return false
	}
}

func validateMetricName(metricName string, collector *Collector) bool {
	_, err := collector.GetMetric(metricName)
	if err != nil {
		return false
	}
	return true
}
