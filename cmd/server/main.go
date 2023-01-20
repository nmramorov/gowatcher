package main

import (
	"net/http"

	"internal/metrics"
)

func main() {
	metrics.InfoLog.Println("Initializing web server...")
	config, err := metrics.NewConfig()
	if err != nil {
		panic(err)
	}
	metricsHandler := metrics.NewHandler()

	server := &http.Server{
		Addr:    config.Address,
		Handler: metricsHandler,
	}

	metrics.InfoLog.Println("Web server is ready to accept connections...")
	server.ListenAndServe()
}
