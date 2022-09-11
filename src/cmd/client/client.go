package main

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/coordinator"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
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
	log.Printf("Duration: %d", duration)
	return coordinatorOutput, duration
}

func PrintResults(coordinatorOutput coordinator.CoordinatorOutput, duration time.Duration) {
	sum := 0
	fmt.Printf("Time taken to fetch response: %d\n", duration)
	fmt.Printf("Response from server:\n\t\t\tFile Name\t\tMatches\n")

	for i := 0; i < len(coordinatorOutput.FileName); i++ {
		log.Printf("%s\t\t%s\n", coordinatorOutput.FileName[i], coordinatorOutput.Matches[i])
		intVar, err := strconv.Atoi(coordinatorOutput.Matches[i])
		if err != nil {
			log.Printf("Error from server %d; calculating remaining sum.", i+1)
		} else {
			sum = sum + intVar
		}
	}
	fmt.Printf("Total matches: %d", sum)
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
