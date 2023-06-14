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
		reader, err := file.NewFileReader(path + options.StoreFile)
		defer func() {
			err := reader.Close()
			if err != nil {
				log.ErrorLog.Printf("Error closing file during read operation: %e", err)
			}
		}()
		if err != nil {
			log.ErrorLog.Printf("Error happend creating File Reader: %e", err)
			return nil, err
		}
		storedMetrics, err := reader.ReadJson()
		if err != nil {
			log.ErrorLog.Printf("Error happend during JSON reading: %e", err)
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
	writer, err := file.NewFileWriter(path + options.StoreFile)
	defer func() {
		err := writer.Close()
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
			err := writer.WriteJson(handler.GetCurrentMetrics())
			if err != nil {
				log.ErrorLog.Printf("Error happened during saving metrics to JSON: %e", err)
			}
			log.InfoLog.Println("Metrics successfully saved to file")
		}
	}
}

type Server struct{}

func (s *Server) Run() {
	serverConfig := config.GetServerConfig()

	metricsHandler, _ := GetMetricsHandler(serverConfig)

	if serverConfig.Database != "" {
		metricsHandler.InitDb()
	}
	if serverConfig.StoreFile != "" {
		go StartSavingToDisk(serverConfig, metricsHandler)
		log.InfoLog.Println("Initialized file saving")
	}

	server := &http.Server{
		Addr:    serverConfig.Address,
		Handler: metricsHandler,
	}

	log.InfoLog.Println("Web server is ready to accept connections...")
	server.ListenAndServe()
}
