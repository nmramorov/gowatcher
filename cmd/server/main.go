package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"internal/metrics"
)

func main() {
	metrics.InfoLog.Println("Initializing web server...")
	metricsHandler := metrics.MetricsHandler{
		Metrics: metrics.NewMetrics(),
	}

	router := metrics.NewRouter()

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: &metricsHandler,
	}

	metrics.InfoLog.Println("Web server is ready to accept connections...")
	server.ListenAndServe()
}
