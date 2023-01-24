package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"internal/metrics"
)

func GetMetricsHandler(options *metrics.ServerCLIOptions) *metrics.Handler {
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

func StartSavingToDisk(options *metrics.ServerCLIOptions, handler *metrics.Handler) {
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
		if timeDiffSec%int64(300) == 0 {
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
	var address = flag.String("a", "localhost:8080", "server address")
	var restore = flag.Bool("r", true, "restore metrics from file")
	var storeInterval = flag.String("i", "300s", "period between file save")
	var storeFile = flag.String("f", "/tmp/devops-metrics-db.json", "name of file where metrics stored")
	flag.Parse()

	args := &metrics.ServerCLIOptions{
		Address:       *address,
		Restore:       *restore,
		StoreInterval: *storeInterval,
		StoreFile:     *storeFile,
	}
	fmt.Println(args)

	metricsHandler := GetMetricsHandler(args)
	if args.StoreFile != "" {
		go StartSavingToDisk(args, metricsHandler)
	}

	server := &http.Server{
		Addr:    args.Address,
		Handler: metricsHandler,
	}

	metrics.InfoLog.Println("Web server is ready to accept connections...")
	server.ListenAndServe()
}
