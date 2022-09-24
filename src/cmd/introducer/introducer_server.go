package main

import (
	"CS425/cs-425-mp1/src/introducer"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8001))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Initialize introducer structure that implements introducer service interface
	c := introducer.Introducer{}

	grpcServer := grpc.NewServer()

	introducer.RegisterIntroducerServiceServer(grpcServer, &c)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
