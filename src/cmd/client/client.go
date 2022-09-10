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

func main() {

	coordinatorConfigs := conf.GetCoordinatorConfigs()

	// take input as grep -Ec arjun .log or as grep -c arjun .log
	clientInputFlag := os.Args[2]
	clientInputString := os.Args[3]
	sum := 0
	for i := 0; i < len(coordinatorConfigs.Endpoints); i++ {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(coordinatorConfigs.Endpoints[i], grpc.WithInsecure(), grpc.WithTimeout(time.Duration(250)*time.Millisecond), grpc.WithBlock())
		if err != nil {
			fmt.Println(err)
		} else {
			defer conn.Close()
			c := coordinator.NewCoordinatorServiceClient(conn)
			// take input from user to a list of args
			start := time.Now()
			coordinatorOutput, err := c.FetchCoordinatorOutput(context.Background(), &coordinator.CoordinatorInput{Data: clientInputString, Flag: clientInputFlag})
			if err != nil {
				fmt.Println("Error while querying distributed log querier : ", err)
			}
			duration := time.Since(start)
			fmt.Println(duration)

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
			break
		}
	}
}
