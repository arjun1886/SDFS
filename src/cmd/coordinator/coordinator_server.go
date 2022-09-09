package main

import (
	"fmt"
	"log"
	"net"

	"CS425/cs-425-mp1/src/coordinator"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	c := coordinator.Coordinator{}

	grpcServer := grpc.NewServer()

	coordinator.RegisterCoordinatorServiceServer(grpcServer, &c)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
