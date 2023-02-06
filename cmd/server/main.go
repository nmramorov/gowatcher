package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"internal/metrics"
)

func GetMetricsHandler(options *metrics.ServerConfig) (*metrics.Handler, error) {
	path, err := filepath.Abs(".")
	if err != nil {
		metrics.ErrorLog.Printf("no file to save exist: %e", err)
		return nil, err
	}
	if options.Restore {
		metrics.InfoLog.Println("Restoring configuration from file...")
		reader, err := metrics.NewFileReader(path + options.StoreFile)
		defer func() {
			err := reader.Close()
			if err != nil {
				metrics.ErrorLog.Printf("Error closing file during read operation: %e", err)
			}
		}()
		if err != nil {
			metrics.ErrorLog.Printf("Error happend creating File Reader: %e", err)
			return nil, err
		}
		storedMetrics, err := reader.ReadJson()
		if err != nil {
			metrics.ErrorLog.Printf("Error happend during JSON reading: %e", err)
			return metrics.NewHandler(), nil
		}
		metricsHandler := metrics.NewHandlerFromSavedData(storedMetrics)
		metrics.InfoLog.Println("Configuration restored.")
		return metricsHandler, nil
	}
	return metrics.NewHandler(), nil
}

func StartSavingToDisk(options *metrics.ServerConfig, handler *metrics.Handler) error {
	path, err := filepath.Abs(".")
	if err != nil {
		metrics.ErrorLog.Printf("no file to save exist: %e", err)
		return err
	}
	writer, err := metrics.NewFileWriter(path + options.StoreFile)
	defer func() {
		err := writer.Close()
		if err != nil {
			metrics.ErrorLog.Printf("Error closing file during write operation: %e", err)
		}
	}()
	if err != nil {
		metrics.ErrorLog.Printf("Error with file writer: %e", err)
		return err
	}
	ticker := time.NewTicker(1 * time.Second)
	startTime := time.Now()
	for {
		tickedTime := <-ticker.C
		timeDiffSec := int64(tickedTime.Sub(startTime).Seconds())
		if timeDiffSec%int64(options.StoreInterval) == 0 {
			fmt.Printf("Writing to file %s", path+options.StoreFile)
			err := writer.WriteJson(handler.GetCurrentMetrics())
			if err != nil {
				metrics.ErrorLog.Printf("Error happened during saving metrics to JSON: %e", err)
			}
			metrics.InfoLog.Println("Metrics successfully saved to file")
		}
	}
}

func main() {
	serverConfig := metrics.GetServerConfig()

	metricsHandler, _ := GetMetricsHandler(serverConfig)
	fmt.Println(serverConfig)
	if serverConfig.StoreFile != "" {
		go StartSavingToDisk(serverConfig, metricsHandler)
		metrics.InfoLog.Println("Initialized file saving")
	}

	server := &http.Server{
		Addr:    serverConfig.Address,
		Handler: metricsHandler,
	}

	metrics.InfoLog.Println("Web server is ready to accept connections...")
	server.ListenAndServe()
}
