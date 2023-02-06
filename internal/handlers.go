package metrics

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*chi.Mux
	collector *Collector
	secretkey string
}

func NewHandler(key string) *Handler {
	h := &Handler{
		Mux:       chi.NewMux(),
		collector: NewCollector(),
		secretkey: key,
	}
	h.Use(GzipHandle)
	h.Get("/", h.ListMetricsHTML)
	h.Get("/value/{type}/{name}", h.GetMetricByTypeAndName)
	h.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	h.Post("/update/", h.UpdateMetricsJson)
	h.Post("/value/", h.GetMetricByJson)

	return h
}

func NewHandlerFromSavedData(saved *Metrics) *Handler {
	h := &Handler{
		Mux:       chi.NewMux(),
		collector: NewCollectorFromSavedFile(saved),
	}
	h.Use(GzipHandle)
	h.Get("/", h.ListMetricsHTML)
	h.Get("/value/{type}/{name}", h.GetMetricByTypeAndName)
	h.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	h.Post("/update/", h.UpdateMetricsJson)
	h.Post("/value/", h.GetMetricByJson)

	return h
}

func (h *Handler) checkHash(rw http.ResponseWriter, metricData *JSONMetrics) {
	var hash string
	generator := NewHashGenerator(h.secretkey)
	switch metricData.MType {
	case "gauge":
		hash = generator.GenerateHash(metricData.MType, metricData.ID, *metricData.Value)
	case "counter":
		hash = generator.GenerateHash(metricData.MType, metricData.ID, *metricData.Delta)
	}
	if hash != metricData.Hash {
		ErrorLog.Printf("wrong hash for %s", metricData.ID)
		http.Error(rw, "wrong hash", http.StatusBadRequest)
		return
	}
	InfoLog.Printf("hash is valid for %s", metricData.ID)
}

func (h *Handler) getHash(metricData *JSONMetrics) string {
	var hash string
	generator := NewHashGenerator(h.secretkey)
	switch metricData.MType {
	case "gauge":
		value := &metricData.Value
		hash = generator.GenerateHash(metricData.MType, metricData.ID, value)
	case "counter":
		delta := &metricData.Delta
		hash = generator.GenerateHash(metricData.MType, metricData.ID, delta)
	}
	return hash
}

func (h *Handler) UpdateMetricsJson(rw http.ResponseWriter, r *http.Request) {
	metricData := JSONMetrics{}
	if err := json.NewDecoder(r.Body).Decode(&metricData); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	InfoLog.Println(metricData)
	if h.secretkey != "" {
		h.checkHash(rw, &metricData)
	}
	updatedData, err := h.collector.UpdateMetricFromJson(&metricData)
	InfoLog.Println(updatedData)
	if err != nil {
		panic("Error occured during metric update from json")
	}
	InfoLog.Println(metricData.Hash)
	InfoLog.Println(h.getHash(&metricData))
	InfoLog.Println(h.getHash(updatedData))
	updatedData.Hash = h.getHash(updatedData)
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.Encode(updatedData)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(buf.Bytes())
}

func (h *Handler) GetMetricByJson(rw http.ResponseWriter, r *http.Request) {
	metricData := JSONMetrics{}
	if err := json.NewDecoder(r.Body).Decode(&metricData); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	metric, err := h.collector.GetMetricJson(&metricData)
	if err != nil {
		panic("Error occured during metric getting from json")
	}
	var hash string
	if h.secretkey != "" {
		hash = h.getHash(metric)
	}
	metric.Hash = hash
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.Encode(metric)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(buf.Bytes())
}

func (h *Handler) GetMetricByTypeAndName(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	isValid := ValidateMetric(metricType, metricName, h.collector)
	if isValid {
		metric, err := h.collector.GetMetric(metricName)
		if err != nil {
			ErrorLog.Fatalf("No such metric %s of type %s: %e", metricName, metricType, err)
			http.Error(rw, "Metric not found", http.StatusNotFound)
		}
		payload, err := h.collector.String(metric)
		if err != nil {
			ErrorLog.Fatalf("Encoding error with metric %s of type %s: %e", metricName, metricType, err)
			http.Error(rw, "Decoding error", http.StatusInternalServerError)
		}
		rw.Write([]byte(payload))
	} else {
		http.Error(rw, "Metric not found", http.StatusNotFound)
	}
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	InfoLog.Println("Started handling metric...")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed for now", http.StatusMethodNotAllowed)
		return
	}
	InfoLog.Println("Method is valid.")
	w.Header().Set("Content-Type", "text/plain")

	path := r.URL.Path
	args := strings.Split(path, "/")
	operation := args[1]
	if strings.Compare(operation, "update") != 0 {
		http.Error(w, "Provide proper operation", http.StatusNotFound)
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
			http.Error(w, "Wrong Gauge value", http.StatusBadRequest)
			return
		}
		h.collector.metrics.GaugeMetrics[metricName] = Gauge(newMetricValue)
		InfoLog.Printf("Value %s is set to %f", metricName, newMetricValue)

	case "counter":
		newMetricValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong Counter value", http.StatusBadRequest)
			return
		}
		InfoLog.Printf("Value %s is set to %d", metricName, newMetricValue)
		newValue := h.collector.metrics.CounterMetrics[metricName] + Counter(newMetricValue)
		h.collector.metrics.CounterMetrics[metricName] = newValue
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) ListMetricsHTML(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("").Parse(`
	<strong>Gauge Metrics:</strong>\n {{range $key, $val := .GaugeMetrics}} {{$key}} = {{$val}}\n {{end}}
	<strong>Counter Metrics:</strong>\n {{range $key, $val := .CounterMetrics}} {{$key}} = {{$val}}\n {{end}}
	`))
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, h.collector.metrics)
}

func (h *Handler) GetCurrentMetrics() *Metrics {
	return h.collector.metrics
}
