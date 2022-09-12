package main

import (
	"fmt"
	"log"
	"net"

	"CS425/cs-425-mp1/src/worker"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8000))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Initialize worker struct that implements worker service interface
	w := worker.Worker{}

	grpcServer := grpc.NewServer()

	worker.RegisterWorkerServiceServer(grpcServer, &w)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
