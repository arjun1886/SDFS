package main

import (
	"CS425/cs-425-mp1/src/coordinator"
	"context"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to server: %s", err)
	}
	defer conn.Close()

	c := coordinator.NewCoordinatorServiceClient(conn)

	var sum int = 0

	// take input as grep -Ec arjun .log or as grep -c arjun .log
	clientInputFlag := os.Args[2]
	clientInputString := os.Args[3]

	// take input from user to a list of args
	coordinatorOutput, err := c.FetchCoordinatorOutput(context.Background(), &coordinator.CoordinatorInput{Data: clientInputString, Flag: clientInputFlag})
	if err != nil {
		log.Fatalf("Error when calling Distributed Log Querier: %s", err)
	}

	log.Printf("Response from server:\n\t\t\tFile Name\t\tMatches")

	for i := 0; i < len(coordinatorOutput.FileName); i++ {
		log.Printf("%s\t\t%s\n", coordinatorOutput.FileName[i], coordinatorOutput.Matches[i])
		intVar, err := strconv.Atoi(coordinatorOutput.Matches[i])
		if err != nil {
			log.Printf("Error from server %d; calculating remaining sum.", i)
		} else {
			sum = sum + intVar
		}
	}
	log.Printf("Total matches: %d", sum)

}
