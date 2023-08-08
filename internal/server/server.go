package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nmramorov/gowatcher/internal/api/handlers"
	"github.com/nmramorov/gowatcher/internal/config"
	"github.com/nmramorov/gowatcher/internal/db"
	"github.com/nmramorov/gowatcher/internal/file"
	"github.com/nmramorov/gowatcher/internal/log"
	pb "github.com/nmramorov/gowatcher/internal/proto"
	"google.golang.org/grpc"
)

func GetMetricsHandler(parent context.Context, options *config.ServerConfig) (*handlers.Handler, error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	cursor, err := db.NewCursor(ctx, options.Database, "pgx")
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
			return handlers.NewHandler(options.Key, options.PrivateKeyPath,
				options.TrustedSubnet, cursor), nil
		}
		metricsHandler := handlers.NewHandlerFromSavedData(storedMetrics, options.Key,
			options.PrivateKeyPath, options.TrustedSubnet, cursor)
		log.InfoLog.Println("Configuration restored.")
		return metricsHandler, nil
	}
	return handlers.NewHandler(options.Key, options.PrivateKeyPath, options.TrustedSubnet, cursor), nil
}

func StartSavingToDisk(killSig chan struct{}, options *config.ServerConfig, handler *handlers.Handler) error {
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
	ticker := time.NewTicker(time.Duration(options.StoreInterval) * time.Second)

	for {
		select {
		case <-killSig:
			log.InfoLog.Println("Stop saving file")
			return nil
		case <-ticker.C:
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

func (s *Server) Run(parent context.Context) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	idleConnsClosed := make(chan struct{})
	// канал для перенаправления прерываний
	// поскольку нужно отловить всего одно прерывание,
	// ёмкости 1 для канала будет достаточно
	sigint := make(chan os.Signal, 1)
	killFileSave := make(chan struct{}, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	var wg sync.WaitGroup

	serverConfig, err := config.GetServerConfig()
	if err != nil {
		log.ErrorLog.Printf("could not get server config: %e", err)
		return err
	}
	metricsHandler, err := GetMetricsHandler(ctx, serverConfig)
	if err != nil {
		log.ErrorLog.Printf("could not get metrics handler: %e", err)
		return err
	}

	if serverConfig.Database != "" {
		err = metricsHandler.InitDB(ctx)
		if err != nil {
			log.ErrorLog.Printf("error initializing db: %e", err)
			return err
		}
	}
	if serverConfig.StoreFile != "" {
		wg.Add(1)
		go func() {
			err = StartSavingToDisk(killFileSave, serverConfig, metricsHandler)
			if err != nil {
				log.ErrorLog.Printf("error starting saving file to disk: %e", err)
			}
			wg.Done()
		}()
		log.InfoLog.Println("Initialized file saving")
	}

	server := &http.Server{
		Addr:              serverConfig.Address,
		Handler:           metricsHandler,
		ReadHeaderTimeout: time.Duration(ServerReadHeaderTimeout) * time.Second,
	}

	// запускаем горутину обработки пойманных прерываний
	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err = server.Shutdown(ctx); err != nil {
			// ошибки закрытия Listener
			log.ErrorLog.Printf("HTTP server Shutdown: %v", err)
		}
		// Kill file save
		close(killFileSave)
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		log.InfoLog.Println("closing channels, shutting down server")
		close(idleConnsClosed)
	}()

	defer func() {
		err = server.Shutdown(ctx)
		if err != nil {
			log.ErrorLog.Printf("error shutting down server: %e", err)
		}
	}()

	if serverConfig.GRPC {
		// определяем порт для сервера
		listen, err := net.Listen("tcp", ":"+strings.Split(serverConfig.Address, ":")[1])
		if err != nil {
			log.ErrorLog.Fatal(err)
		}
		// создаём gRPC-сервер без зарегистрированной службы
		s := grpc.NewServer()
		// регистрируем сервис
		pb.RegisterMetricsServer(s, &MetricsServer{})

		log.InfoLog.Println("Сервер gRPC начал работу")
		// получаем запрос gRPC
		if err := s.Serve(listen); err != nil {
			log.ErrorLog.Fatal(err)
		}
	}

	log.InfoLog.Println("Web server is ready to accept connections...")
	err = server.ListenAndServe()
	if err != nil {
		log.ErrorLog.Printf("Unexpected error occurred: %e", err)
	}

	<-idleConnsClosed
	wg.Wait()
	return nil
}
