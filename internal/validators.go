package metrics

func validateMetricType(metricType string) bool {
	switch metricType {
	case "gauge":
		return true
	case "counter":
		return true
	default:
		return false
	}
}

func validateMetricName(metricName string, collector *Collector) bool {
	_, err := collector.GetMetric(metricName)
	return err == nil
}

func ValidateMetric(metricType, metricName string, col *Collector) bool {
	isValidType := validateMetricType(metricType)
	if isValidType {
		isValidName := validateMetricName(metricName, col)
		if isValidName {
			return true
		}
	}
	return false
}
