package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"master-worker-system/internal/worker"
	pb "master-worker-system/pb"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 50051, "gRPC server port")
	workerID := flag.String("id", "worker-1", "Worker ID")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	workerServer := worker.NewWorkerServer(*workerID)

	pb.RegisterWorkerServiceServer(grpcServer, workerServer)

	log.Printf("Worker %s starting gRPC server on port %d", *workerID, *port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
