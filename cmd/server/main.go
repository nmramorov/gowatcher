package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"internal/metrics"
)

func GetMetricsHandler(config *metrics.EnvConfig) *metrics.Handler {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	if config.Restore {
		metrics.InfoLog.Println("Restoring configuration from file...")
		reader, err := metrics.NewFileReader(path + config.StoreFile)
		defer func() {
			err := reader.Close()
			if err != nil {
				metrics.ErrorLog.Printf("Error closing file during read operation: %e", err)
			}
		}()
		if err != nil {
			metrics.ErrorLog.Printf("Error happend creating File Reader: %e", err)
			panic(err)
		}
		storedMetrics, err := reader.ReadJson()
		if err != nil {
			metrics.ErrorLog.Printf("Error happend during JSON reading: %e", err)
			return metrics.NewHandler()
		}
		metricsHandler := metrics.NewHandlerFromSavedData(storedMetrics)
		metrics.InfoLog.Println("Configuration restored.")
		return metricsHandler
	} else {
		return metrics.NewHandler()
	}
}

func StartSavingToDisk(config *metrics.EnvConfig, handler *metrics.Handler) {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	writer, err := metrics.NewFileWriter(path + config.StoreFile)
	defer func() {
		err := writer.Close()
		if err != nil {
			metrics.ErrorLog.Printf("Error closing file during write operation: %e", err)
		}
	}()
	if err != nil {
		metrics.ErrorLog.Printf("Error with file writer: %e", err)
	}
	ticker := time.NewTicker(1 * time.Second)
	startTime := time.Now()
	for {
		tickedTime := <-ticker.C
		timeDiffSec := int64(tickedTime.Sub(startTime).Seconds())
		interval, err := config.GetNumericInterval("StoreInterval")
		if err != nil {
			panic(err)
		}
		if timeDiffSec%int64(interval) == 0 {
			fmt.Println()
			err := writer.WriteJson(handler.GetCurrentMetrics())
			if err != nil {
				metrics.ErrorLog.Printf("Error happened during saving metrics to JSON: %e", err)
			}
			metrics.InfoLog.Println("Metrics successfully saved to file")
		}
	}

}

func main() {
	metrics.InfoLog.Println("Initializing web server...")
	config, err := metrics.NewConfig()
	if err != nil {
		panic(err)
	}
	metricsHandler := GetMetricsHandler(config)
	if config.StoreFile != "" {
		go StartSavingToDisk(config, metricsHandler)
	}

	server := &http.Server{
		Addr:    config.Address,
		Handler: metricsHandler,
	}

	metrics.InfoLog.Println("Web server is ready to accept connections...")
	server.ListenAndServe()
}
