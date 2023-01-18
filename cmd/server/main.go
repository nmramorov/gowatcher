package main

import (
	"net/http"

	"internal/metrics"
)

func main() {
	metrics.InfoLog.Println("Initializing web server...")
	metricsHandler := metrics.NewHandler()

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: metricsHandler,
	}

	metrics.InfoLog.Println("Web server is ready to accept connections...")
	server.ListenAndServe()
}
