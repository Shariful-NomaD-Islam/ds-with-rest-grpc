package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shariful-NomaD-Islam/ds-with-rest-grpc/internal/logger"
	"github.com/Shariful-NomaD-Islam/ds-with-rest-grpc/internal/worker"
	pb "github.com/Shariful-NomaD-Islam/ds-with-rest-grpc/pb"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 50051, "gRPC server port")
	workerID := flag.String("id", "worker-1", "Worker ID")
	logLevel := flag.String("log-level", "info", "Logging level (debug, info, warn, error, fatal)")
	flag.Parse()

	// Initialize logger
	logger.Init(*logLevel)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.GetLogger().Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	workerServer := worker.NewWorkerServer(*workerID)

	pb.RegisterWorkerServiceServer(grpcServer, workerServer)

	logger.GetLogger().Infof("Worker %s starting gRPC server on port %d", *workerID, *port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.GetLogger().Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.GetLogger().Info("Worker shutting down...")
	grpcServer.GracefulStop()
}
