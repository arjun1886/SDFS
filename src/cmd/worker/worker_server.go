package main

import (
	"fmt"
	"log"
	"net"

	"CS425/cs-425-mp1/src/worker"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9001))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	w := worker.Worker{}

	grpcServer := grpc.NewServer()

	worker.RegisterWorkerServiceServer(grpcServer, &w)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
