package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"internal/metrics"
)

type MetricsHandler struct {
	metrics *metrics.Metrics
}

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

func (m *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	infoLog.Println("Started handling metric...")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed for now", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Provide proper header", http.StatusForbidden)
		return
	}
	infoLog.Println("Method and Headers are valid.")
	path := r.URL.Path
	args := strings.Split(path, "/")
	if len(args) != 3 {
		http.Error(w, "Wrong arguments in request", http.StatusInternalServerError)
		return
	}
	var metricType, metricName, metricValue = args[0], args[1], args[2]
	infoLog.Printf("Received metric data:\nMetric type: %s\nMetric name: %s\nMetric value: %s", metricType, metricName, metricValue)

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
		infoLog.Printf("Value %s is set to %f", metricName, newMetricValue)

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
		infoLog.Printf("Value %s is set to %d", metricName, newMetricValue)
		m.metrics.CounterMetrics[metricName] = metrics.ToCounter(newMetricValue)
	default:
		http.Error(w, "Wrong metric type", http.StatusInternalServerError)
	}
	infoLog.Println("Request successfully handled")
}

func main() {
	infoLog.Println("Initializing web server...")
	metricsHandler := MetricsHandler{
		metrics: &metrics.Metrics{},
	}

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: &metricsHandler,
	}

	infoLog.Println("Web server is ready to accept connections...")
	server.ListenAndServe()
}
