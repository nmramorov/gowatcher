package handlers

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	_ "net/http/pprof"

	middleware "github.com/nmramorov/gowatcher/internal/api/middlewares"
	val "github.com/nmramorov/gowatcher/internal/api/validators"
	col "github.com/nmramorov/gowatcher/internal/collector"
	m "github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/db"
	"github.com/nmramorov/gowatcher/internal/hashgen"
	"github.com/nmramorov/gowatcher/internal/log"
)

// Базовый тип Handler, отвечающий за обработку запросов.
type Handler struct {
	*chi.Mux
	collector *col.Collector
	secretkey string
	cursor    *db.Cursor
}

// Конструктор для объектов типа Handler.
func NewHandler(key string, newCursor *db.Cursor) *Handler {
	h := &Handler{
		Mux:       chi.NewMux(),
		collector: col.NewCollector(),
		secretkey: key,
		cursor:    newCursor,
	}
	h.Use(middleware.GzipHandle)
	h.Get("/", h.ListMetricsHTML)
	h.Get("/ping", h.HandlePing)
	h.Get("/value/{type}/{name}", h.GetMetricByTypeAndName)
	h.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	h.Post("/update/", h.UpdateMetricsJSON)
	h.Post("/value/", h.GetMetricByJSON)
	h.Post("/updates/", h.UpdateJSONBatch)

	return h
}

// Конструктор Handler, который инициализируется с записанными ранее данными Metrics.
func NewHandlerFromSavedData(saved *m.Metrics, secretkey string, cursor *db.Cursor) *Handler {
	h := &Handler{
		Mux:       chi.NewMux(),
		collector: col.NewCollectorFromSavedFile(saved),
		secretkey: secretkey,
		cursor:    cursor,
	}
	h.Use(middleware.GzipHandle)
	h.Get("/", h.ListMetricsHTML)
	h.Get("/ping", h.HandlePing)
	h.Get("/value/{type}/{name}", h.GetMetricByTypeAndName)
	h.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	h.Post("/update/", h.UpdateMetricsJSON)
	h.Post("/value/", h.GetMetricByJSON)
	h.Post("/updates/", h.UpdateJSONBatch)

	return h
}

func (h *Handler) checkHash(rw http.ResponseWriter, metricData *m.JSONMetrics) {
	var hash string
	generator := hashgen.NewHashGenerator(h.secretkey)
	switch metricData.MType {
	case "gauge":
		hash = generator.GenerateHash(metricData.MType, metricData.ID, *metricData.Value)
	case "counter":
		hash = generator.GenerateHash(metricData.MType, metricData.ID, *metricData.Delta)
	}
	d, _ := hex.DecodeString(hash)
	if hmac.Equal(d, []byte(metricData.Hash)) {
		log.ErrorLog.Printf("wrong hash for %s", metricData.ID)
		http.Error(rw, "wrong hash", http.StatusBadRequest)
		return
	}
}

func (h *Handler) getHash(metricData *m.JSONMetrics) string {
	var hash string
	generator := hashgen.NewHashGenerator(h.secretkey)
	switch metricData.MType {
	case "gauge":
		hash = generator.GenerateHash(metricData.MType, metricData.ID, *metricData.Value)
	case "counter":
		hash = generator.GenerateHash(metricData.MType, metricData.ID, *metricData.Delta)
	}
	return hash
}

// Метод для обновления метрики, полученной в формате JSON.
func (h *Handler) UpdateMetricsJSON(rw http.ResponseWriter, r *http.Request) {
	metricData := m.JSONMetrics{}
	if err := json.NewDecoder(r.Body).Decode(&metricData); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if h.secretkey != "" {
		h.checkHash(rw, &metricData)
	}
	updatedData, err := h.collector.UpdateMetricFromJSON(&metricData)
	if h.cursor.IsValid {
		err = h.cursor.Add(updatedData)
		if err != nil {
			log.ErrorLog.Println("could not add data to db...")
		}
	}
	if err != nil {
		log.ErrorLog.Printf("Error occured during metric update from json: %e", err)
	}
	updatedData.Hash = h.getHash(updatedData)
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(updatedData)
	if err != nil {
		log.ErrorLog.Printf("error encoding updated data: %e", err)
	}
	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(buf.Bytes())
	if err != nil {
		log.ErrorLog.Printf("error writing data to update metrics request: %e", err)
	}
}

// Метод, позволяющий получить требуемую метрику в формате JSON.
// На вход требует JSON с заполненными полями id и mtype.
func (h *Handler) GetMetricByJSON(rw http.ResponseWriter, r *http.Request) {
	metricData := m.JSONMetrics{}
	if err := json.NewDecoder(r.Body).Decode(&metricData); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var metric *m.JSONMetrics
	var err error
	if h.cursor.IsValid {
		metric, err = h.cursor.Get(&metricData)
		if err != nil {
			log.ErrorLog.Println("could not get data from db...")
			metric, err = h.collector.GetMetricJSON(&metricData)
			if err != nil {
				log.ErrorLog.Printf("Error occured during metric getting from json: %e", err)
			}
		}
	} else {
		metric, err = h.collector.GetMetricJSON(&metricData)
		if err != nil {
			log.ErrorLog.Printf("Error occured during metric getting from json: %e", err)
		}
	}
	var hash string
	if h.secretkey != "" {
		hash = h.getHash(metric)
	}
	metric.Hash = hash
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(metric)
	if err != nil {
		log.ErrorLog.Printf("error encoding get metric: %e", err)
	}
	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(buf.Bytes())
	if err != nil {
		log.ErrorLog.Printf("error writing data to get metrics request: %e", err)
	}
}

// Deprecated: метод был создан для первых инкрементов, в настоящее время не используется.
func (h *Handler) GetMetricByTypeAndName(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	isValid := val.ValidateMetric(metricType, metricName, h.collector)
	if isValid {
		metric, err := h.collector.GetMetric(metricName)
		if err != nil {
			log.ErrorLog.Fatalf("No such metric %s of type %s: %e", metricName, metricType, err)
			http.Error(rw, "Metric not found", http.StatusNotFound)
		}
		payload, err := h.collector.String(metric)
		if err != nil {
			log.ErrorLog.Fatalf("Encoding error with metric %s of type %s: %e", metricName, metricType, err)
			http.Error(rw, "Decoding error", http.StatusInternalServerError)
		}

		_, err = rw.Write([]byte(payload))
		if err != nil {
			log.ErrorLog.Printf("error writing data to get metrics by type and name request: %e", err)
		}
	} else {
		http.Error(rw, "Metric not found", http.StatusNotFound)
	}
}

// Deprecated: метод был создан для первых инкрементов, в настоящее время не используется.
func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	log.InfoLog.Println("Started handling metric...")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed for now", http.StatusMethodNotAllowed)
		return
	}
	log.InfoLog.Println("Method is valid.")
	w.Header().Set("Content-Type", "text/plain")

	path := r.URL.Path
	args := strings.Split(path, "/")
	operation := args[1]
	if strings.Compare(operation, "update") != 0 {
		http.Error(w, "Provide proper operation", http.StatusNotFound)
		return
	}
	log.InfoLog.Println(args)
	if len(args) != 5 {
		http.Error(w, "Wrong arguments in request", http.StatusNotFound)
		return
	}
	metricType, metricName, metricValue := args[2], args[3], args[4]
	log.InfoLog.Printf("Received metric data:\nMetric type: %s\nMetric name: %s\nMetric value: %s", metricType, metricName, metricValue)
	switch metricType {
	case "gauge":
		newMetricValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Wrong Gauge value", http.StatusBadRequest)
			return
		}
		h.collector.Metrics.GaugeMetrics[metricName] = m.Gauge(newMetricValue)
		log.InfoLog.Printf("Value %s is set to %f", metricName, newMetricValue)

	case "counter":
		newMetricValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong Counter value", http.StatusBadRequest)
			return
		}
		log.InfoLog.Printf("Value %s is set to %d", metricName, newMetricValue)
		newValue := h.collector.Metrics.CounterMetrics[metricName] + m.Counter(newMetricValue)
		h.collector.Metrics.CounterMetrics[metricName] = newValue
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}
	_, err := w.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		log.ErrorLog.Printf("error writing data to non-JSON update request: %e", err)
	}
}

// Метод, позволяющий по корню посмотреть состояние текущих метрик на сервере.
func (h *Handler) ListMetricsHTML(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("").Parse(`
	<strong>Gauge Metrics:</strong>\n {{range $key, $val := .GaugeMetrics}} {{$key}} = {{$val}}\n {{end}}
	<strong>Counter Metrics:</strong>\n {{range $key, $val := .CounterMetrics}} {{$key}} = {{$val}}\n {{end}}
	`))
	w.Header().Set("Content-Type", "text/html")
	err := t.Execute(w, h.collector.Metrics)
	if err != nil {
		log.ErrorLog.Printf("error getting HTML list of metrics: %e", err)
	}
}

// Вспомогательный метод для получения метрик из коллектора.
func (h *Handler) GetCurrentMetrics() *m.Metrics {
	return h.collector.Metrics
}

// Метод для проверки соединения с БД.
func (h *Handler) HandlePing(w http.ResponseWriter, r *http.Request) {
	isAvailable := h.cursor.Ping()
	defer h.cursor.Close()
	switch isAvailable {
	case true:
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			log.ErrorLog.Printf("error handling ping request: %e", err)
		}
	case false:
		http.Error(w, "error with db", http.StatusInternalServerError)
		return
	}
}

// Метод, инициализирующий БД.
func (h *Handler) InitDB() error {
	return h.cursor.InitDB()
}

// Метод, позволяющий обновить несколько метрик за раз.
func (h *Handler) UpdateJSONBatch(rw http.ResponseWriter, r *http.Request) {
	var metricsBatch []*m.JSONMetrics
	if err := json.NewDecoder(r.Body).Decode(&metricsBatch); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.collector.UpdateBatch(metricsBatch)
	if err != nil {
		log.ErrorLog.Println("could not update batch...")
	}
	if h.cursor.IsValid {
		err = h.cursor.AddBatch(metricsBatch)
		if err != nil {
			log.ErrorLog.Println("could not add batch data to db...")
		}
	}
	if err != nil {
		log.ErrorLog.Printf("Error occured during metric update from json: %e", err)
	}
	log.InfoLog.Println("received and worker with metrics batch")
	_, err = rw.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		log.ErrorLog.Printf("error updating JSON batch request: %e", err)
	}
}
