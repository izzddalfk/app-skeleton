package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gosidekick/goconfig"
	"github.com/izzdalfk/app-skeleton/go/internal/core/service"
	"github.com/izzdalfk/app-skeleton/go/internal/driven/storage"
	"github.com/izzdalfk/app-skeleton/go/internal/driver/rest"
	"github.com/izzdalfk/app-skeleton/go/internal/shared/logger"
)

var cfg config

func main() {
	// parse app config
	err := goconfig.Parse(&cfg)
	if err != nil {
		log.Fatalf("failed to parse app config, due: %v", err)
	}
	// initialize logger
	zlogger, err := logger.NewZerologLogger(logger.Config{
		LogLevel:    "debug",
		ServiceName: "rest-api",
	})
	if err != nil {
		log.Fatalf("failed to initialize zerolog logger, due: %v", err)
	}

	// initialize dummy storage
	dummyStorage := storage.NewDummy()

	// initialize dummy service
	dummyService := service.NewDummy(service.Config{
		Storage: dummyStorage,
	})

	// initialize REST API driver
	api, err := rest.NewAPI(rest.APIConfig{
		Logger:       zlogger,
		DummyService: dummyService,
	})
	if err != nil {
		log.Fatalf("failed to initialize rest api, due: %v", err)
	}

	runServer(api.Handler())
}

func runServer(handler http.Handler) {
	// initialize cancel context
	svCtx, cancel := context.WithCancel(context.Background())

	// initialize the server instance
	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.Port),
		Handler:     handler,
		ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
	}

	// setup shutdown listener
	setShutdownListener(func(sig os.Signal) {
		log.Printf("[%v] shutdown signal received, terminating app...", sig)
		cancel()
		server.Shutdown(svCtx)
	})

	// run the http server
	log.Printf("server is listening on: %d", cfg.Port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to run http server, due: %v", err)
	}
}

func setShutdownListener(onSignalReceived func(os.Signal)) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		sig := <-signalCh
		onSignalReceived(sig)
	}()
}

type config struct {
	Port        int `cfg:"port" cfgDefault:"7100"`
	ReadTimeout int `cfg:"read_timeout" cfgDefault:"5"`
}
