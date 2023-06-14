package validator

import (
	col "github.com/nmramorov/gowatcher/internal/collector"
)

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

func validateMetricName(metricName string, collector *col.Collector) bool {
	_, err := collector.GetMetric(metricName)
	return err == nil
}

func ValidateMetric(metricType, metricName string, col *col.Collector) bool {
	isValidType := validateMetricType(metricType)
	if isValidType {
		isValidName := validateMetricName(metricName, col)
		if isValidName {
			return true
		}
	}
	return false
}
