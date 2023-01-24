package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"internal/metrics"
)

func GetMetricsHandler(options *ServerConfig) *metrics.Handler {
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
	} else {
		return metrics.NewHandler()
	}
}

func StartSavingToDisk(options *ServerConfig, handler *metrics.Handler) {
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

type ServerConfig struct {
	Address       string
	Restore       bool
	StoreInterval int
	StoreFile     string
}

func GetServerConfig(config *metrics.EnvConfig, args *metrics.ServerCLIOptions) *ServerConfig {
	serverConfig := ServerConfig{}
	if config.Address == "127.0.0.1:8080" {
		serverConfig.Address = args.Address
	} else {
		serverConfig.Address = config.Address
	}
	if config.Restore {
		serverConfig.Restore = config.Restore
	} else {
		serverConfig.Restore = args.Restore
	}
	if config.StoreFile != "/tmp/devops-metrics-db.json" {
		serverConfig.StoreFile = config.StoreFile
	} else {
		serverConfig.StoreFile = args.StoreFile
	}
	if config.StoreInterval == "300s" {
		serverConfig.StoreInterval = func() int {
			store, err := args.GetNumericInterval("StoreInterval")
			if err != nil {
				panic(err)
			}
			return int(store)
		}()
	} else {
		serverConfig.StoreInterval = func() int {
			store, err := config.GetNumericInterval("StoreInterval")
			if err != nil {
				panic(err)
			}
			return int(store)
		}()
	}
	return &serverConfig
}

func main() {
	config, err := metrics.NewConfig()
	if err != nil {
		panic(err)
	}
	metrics.InfoLog.Println("Initializing web server...")
	var address = flag.String("a", "localhost:8080", "server address")
	var restore = flag.Bool("r", true, "restore metrics from file")
	var storeInterval = flag.String("i", "30s", "period between file save")
	var storeFile = flag.String("f", "/tmp/devops-metrics-db.json", "name of file where metrics stored")
	flag.Parse()

	args := &metrics.ServerCLIOptions{
		Address:       *address,
		Restore:       *restore,
		StoreInterval: *storeInterval,
		StoreFile:     *storeFile,
	}
	serverConfig := GetServerConfig(config, args)

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
