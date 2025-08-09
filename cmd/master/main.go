package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"master-worker-system/internal/config"
	"master-worker-system/internal/master"
)

func main() {
	configFile := flag.String("config", "config.yml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	config, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Loaded configuration from %s", *configFile)
	log.Printf("Server will start on %s", config.GetServerAddress())
	log.Printf("Configured workers: %d", len(config.Workers))

	// Initialize worker pool
	workerPool, err := master.NewWorkerPool(config)
	if err != nil {
		log.Fatalf("Failed to initialize worker pool: %v", err)
	}
	defer workerPool.Close()

	// Setup HTTP server
	router := master.SetupRoutes(workerPool, config)

	// Start server
	log.Printf("Master starting HTTP server on %s", config.GetServerAddress())
	go func() {
		if err := router.Run(config.GetServerAddress()); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Master shutting down...")
}
