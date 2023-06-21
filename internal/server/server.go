package server

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/nmramorov/gowatcher/internal/api/handlers"
	"github.com/nmramorov/gowatcher/internal/config"
	"github.com/nmramorov/gowatcher/internal/db"
	"github.com/nmramorov/gowatcher/internal/file"
	"github.com/nmramorov/gowatcher/internal/log"
)

func GetMetricsHandler(options *config.ServerConfig) (*handlers.Handler, error) {
	cursor, err := db.NewCursor(options.Database, "pgx")
	if err != nil {
		cursor.IsValid = false
	}
	path, err := filepath.Abs(".")
	if err != nil {
		log.ErrorLog.Printf("no file to save exist: %e", err)
		return nil, err
	}
	if options.Restore {
		log.InfoLog.Println("Restoring configuration from file...")
		reader, err := file.NewReader(path + options.StoreFile)
		defer func() {
			err = reader.Close()
			if err != nil {
				log.ErrorLog.Printf("Error closing file during read operation: %e", err)
			}
		}()
		if err != nil {
			log.ErrorLog.Printf("Error happened creating File Reader: %e", err)
			return nil, err
		}
		storedMetrics, err := reader.ReadJSON()
		if err != nil {
			log.ErrorLog.Printf("Error happened during JSON reading: %e", err)
			return handlers.NewHandler(options.Key, cursor), nil
		}
		metricsHandler := handlers.NewHandlerFromSavedData(storedMetrics, options.Key, cursor)
		log.InfoLog.Println("Configuration restored.")
		return metricsHandler, nil
	}
	return handlers.NewHandler(options.Key, cursor), nil
}

func StartSavingToDisk(options *config.ServerConfig, handler *handlers.Handler) error {
	path, err := filepath.Abs(".")
	if err != nil {
		log.ErrorLog.Printf("no file to save exist: %e", err)
		return err
	}
	writer, err := file.NewWriter(path + options.StoreFile)
	defer func() {
		err = writer.Close()
		if err != nil {
			log.ErrorLog.Printf("Error closing file during write operation: %e", err)
		}
	}()
	if err != nil {
		log.ErrorLog.Printf("Error with file writer: %e", err)
		return err
	}
	ticker := time.NewTicker(1 * time.Second)
	startTime := time.Now()
	for {
		tickedTime := <-ticker.C
		timeDiffSec := int64(tickedTime.Sub(startTime).Seconds())
		if timeDiffSec%int64(options.StoreInterval) == 0 {
			err := writer.WriteJSON(handler.GetCurrentMetrics())
			if err != nil {
				log.ErrorLog.Printf("Error happened during saving metrics to JSON: %e", err)
			}
			log.InfoLog.Println("Metrics successfully saved to file")
		}
	}
}

type Server struct{}

var ServerReadHeaderTimeout = 10

func (s *Server) Run() {
	serverConfig := config.GetServerConfig()

	metricsHandler, _ := GetMetricsHandler(serverConfig)

	if serverConfig.Database != "" {
		err := metricsHandler.InitDB()
		if err != nil {
			log.ErrorLog.Printf("error initializing db: %e", err)
		}
	}
	if serverConfig.StoreFile != "" {
		go func() {
			err := StartSavingToDisk(serverConfig, metricsHandler)
			if err != nil {
				log.ErrorLog.Printf("error starting saving file to disk: %e", err)
			}
		}()
		log.InfoLog.Println("Initialized file saving")
	}

	server := &http.Server{
		Addr:              serverConfig.Address,
		Handler:           metricsHandler,
		ReadHeaderTimeout: time.Duration(ServerReadHeaderTimeout) * time.Second,
	}

	log.InfoLog.Println("Web server is ready to accept connections...")
	err := server.ListenAndServe()
	if err != nil {
		log.ErrorLog.Printf("Unexpected error occurred: %e", err)
	}
}
