package main

import (
	"CS425/cs-425-mp1/src/coordinator"
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := coordinator.NewCoordinatorServiceClient(conn)

	// take input as grep -Ec arjun .log or as grep -c arjun .log
	clientInputFlag := os.Args[2]
	clientInputString := os.Args[3]

	// take input from user to a list of args
	coordinatorOutput, err := c.FetchCoordinatorOutput(context.Background(), &coordinator.CoordinatorInput{Data: clientInputString, Flag: clientInputFlag})
	if err != nil {
		log.Fatalf("Error when calling Distributed Log Querier: %s", err)
	}

	log.Printf("Response from server: %s", coordinatorOutput)

}
