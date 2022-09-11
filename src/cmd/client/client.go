package main

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/coordinator"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

func FetchOutput(clientInputFlag string, clientInputString string, cc coordinator.CoordinatorServiceClient, index int) (coordinatorOutput *coordinator.CoordinatorOutput, duration time.Duration) {
	start := time.Now()
	// take input from user to a list of args
	coordinatorOutput, err := cc.FetchCoordinatorOutput(context.Background(), &coordinator.CoordinatorInput{Data: clientInputString, Flag: clientInputFlag})
	if err != nil {
		log.Printf("Error while querying distributed log querier through coordinator %d: %s", index+1, err)
	}
	duration = time.Since(start)
	return coordinatorOutput, duration
}

func PrintResults(coordinatorOutput coordinator.CoordinatorOutput, duration time.Duration) {
	fmt.Printf("Time taken to fetch response: %d\n", duration)
	fmt.Printf("Response from server:\n\n\t\tFile Name\t\tMatches\n")
	fmt.Printf("Total successful matches: %s", coordinatorOutput.TotalMatchCount)
}

func main() {

	coordinatorConfigs := conf.GetCoordinatorConfigs()

	// take input as grep -Ec arjun .log or as grep -c arjun .log
	clientInputFlag := os.Args[2]
	clientInputString := os.Args[3]

	// Loop over coordinator configs for fault handling

	for i := 0; i < len(coordinatorConfigs.Endpoints); i++ {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(coordinatorConfigs.Endpoints[i], grpc.WithInsecure(), grpc.WithTimeout(time.Duration(2000)*time.Millisecond), grpc.WithBlock())
		if err != nil {
			log.Printf("Error connecting to coordinator %d: %s", i+1, err)
		} else {
			defer conn.Close()
			c := coordinator.NewCoordinatorServiceClient(conn)

			// Fetch output from coordinator
			coordinatorOutput, duration := FetchOutput(clientInputFlag, clientInputString, c, i)

			// Print results
			PrintResults(*coordinatorOutput, duration)
			break
		}
	}
}
