package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"internal/metrics"
)

type MetricsHandler struct {
	metrics *metrics.Metrics
}

func (m *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed for now", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Provide proper header", http.StatusForbidden)
		return
	}
	path := r.URL.Path
	args := strings.Split(path, "/")
	if len(args) != 3 {
		http.Error(w, "Wrong arguments in request", http.StatusInternalServerError)
		return
	}
	var metricType, metricName, metricValue = args[0], args[1], args[2]
	switch metricType {
	case "gauge":
		if _, ok := m.metrics.GaugeMetrics[metricName]; !ok {
			http.Error(w, "No such metric exists", http.StatusInternalServerError)
			return
		}
		newMetricValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Wrong Gauge value", http.StatusInternalServerError)
			return
		}
		m.metrics.GaugeMetrics[metricName] = metrics.ToGauge(newMetricValue)

	case "counter":
		if _, ok := m.metrics.CounterMetrics[metricName]; !ok {
			http.Error(w, "No such metric exists", http.StatusInternalServerError)
			return
		}
		newMetricValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong Gauge value", http.StatusInternalServerError)
			return
		}
		m.metrics.CounterMetrics[metricName] = metrics.ToCounter(newMetricValue)
	default:
		http.Error(w, "Wrong metric type", http.StatusInternalServerError)
	}
	fmt.Println(m.metrics.CounterMetrics)
	fmt.Println(m.metrics.GaugeMetrics)
}

func main() {
	metricsHandler := MetricsHandler{
		metrics: &metrics.Metrics{},
	}

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: &metricsHandler,
	}

	server.ListenAndServe()
}
