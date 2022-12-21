package metrics

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type MetricsHandler struct {
	Metrics *Metrics
}

var InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

func (m *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	InfoLog.Println("Started handling metric...")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed for now", http.StatusMethodNotAllowed)
		return
	}
	// if r.Header.Get("Content-Type") != "text/plain" {
	// 	http.Error(w, "Provide proper header", http.StatusForbidden)
	// 	return
	// }
	InfoLog.Println("Method and Headers are valid.")
	w.Header().Set("Content-Type", "text/plain")

	path := r.URL.Path
	args := strings.Split(path, "/")
	operation := args[1]
	if strings.Compare(operation, "update") != 0 {
		http.Error(w, "Provide proper operation", http.StatusNotImplemented)
		return
	}
	InfoLog.Println(args)
	if len(args) != 5 {
		http.Error(w, "Wrong arguments in request", http.StatusNotFound)
		return
	}
	var metricType, metricName, metricValue = args[2], args[3], args[4]
	InfoLog.Printf("Received metric data:\nMetric type: %s\nMetric name: %s\nMetric value: %s", metricType, metricName, metricValue)
	switch metricType {
	case "gauge":
		newMetricValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Wrong Gauge value", http.StatusInternalServerError)
			return
		}
		m.Metrics.GaugeMetrics[metricName] = Gauge(newMetricValue)
		InfoLog.Printf("Value %s is set to %f", metricName, newMetricValue)

	case "counter":
		newMetricValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong Counter value", http.StatusInternalServerError)
			return
		}
		InfoLog.Printf("Value %s is set to %d", metricName, newMetricValue)
		newValue := m.Metrics.CounterMetrics[metricName] + Counter(newMetricValue)
		m.Metrics.CounterMetrics[metricName] = newValue
	default:
		http.Error(w, "Wrong metric type", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"status":"ok"}`))
}
