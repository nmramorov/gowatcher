package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"internal/metrics"
)

func GetMetricsHandler(options *metrics.ServerConfig) *metrics.Handler {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
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
			panic(err)
		}
		storedMetrics, err := reader.ReadJson()
		metrics.InfoLog.Println(storedMetrics)
		if err != nil {
			metrics.ErrorLog.Printf("Error happend during JSON reading: %e", err)
			return metrics.NewHandler()
		}
		metricsHandler := metrics.NewHandlerFromSavedData(storedMetrics)
		metrics.InfoLog.Println("Configuration restored.")
		return metricsHandler
	}
	return metrics.NewHandler()
}

func StartSavingToDisk(options *metrics.ServerConfig, handler *metrics.Handler) {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
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
	// config, err := metrics.NewConfig()
	// if err != nil {
	// 	panic(err)
	// }
	// metrics.InfoLog.Println("Initializing web server...")
	// var address = flag.String("a", "localhost:8080", "server address")
	// var restore = flag.Bool("r", true, "restore metrics from file")
	// var storeInterval = flag.String("i", "30s", "period between file save")
	// var storeFile = flag.String("f", "/tmp/devops-metrics-db.json", "name of file where metrics stored")
	// flag.Parse()

	// args := &metrics.ServerCLIOptions{
	// 	Address:       *address,
	// 	Restore:       *restore,
	// 	StoreInterval: *storeInterval,
	// 	StoreFile:     *storeFile,
	// }
	serverConfig := metrics.GetServerConfig()

	metricsHandler := GetMetricsHandler(serverConfig)
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
