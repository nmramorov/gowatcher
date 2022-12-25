package metrics

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update", func(r chi.Router) {
		r.Get("/{type}{name}", GetMetricHandler)
	})
	return r
}

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

func GetMetricHandler(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	isValidType := validateMetricType(metricType)
	if isValidType == true {
		isValidName := validateMetricName(metricName)
	}
}
