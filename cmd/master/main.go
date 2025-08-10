package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shariful-NomaD-Islam/ds-with-rest-grpc/internal/config"
	"github.com/Shariful-NomaD-Islam/ds-with-rest-grpc/internal/logger"
	"github.com/Shariful-NomaD-Islam/ds-with-rest-grpc/internal/master"
)

func main() {
	configFile := flag.String("config", "config.yml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		logger.GetLogger().Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger with configured level
	logger.Init(cfg.Logging.Level)

	logger.GetLogger().Infof("Loaded configuration from %s", *configFile)
	logger.GetLogger().Infof("Server will start on %s", cfg.GetServerAddress())
	logger.GetLogger().Infof("Configured workers: %d", len(cfg.Workers))

	// Initialize worker pool
	workerPool, err := master.NewWorkerPool(cfg)
	if err != nil {
		logger.GetLogger().Fatalf("Failed to initialize worker pool: %v", err)
	}
	defer workerPool.Close()

	// Setup HTTP server
	router := master.SetupRoutes(workerPool, cfg)

	// Start server
	logger.GetLogger().Infof("Master starting HTTP server on %s", cfg.GetServerAddress())
	go func() {
		if err := router.Run(cfg.GetServerAddress()); err != nil {
			logger.GetLogger().Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.GetLogger().Info("Master shutting down...")
}
