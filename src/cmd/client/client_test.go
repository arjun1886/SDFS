package main

import (
	"CS425/cs-425-mp1/src/coordinator"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"testing"
	"time"
)

func establishConnection(t *testing.T) (c coordinator.CoordinatorServiceClient) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("fa22-cs425-6602.cs.illinois.edu:8001", grpc.WithInsecure(), grpc.WithTimeout(time.Duration(2000)*time.Millisecond), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Did not connect to coordinator: %s", err)
	}
	defer conn.Close()
	c = coordinator.NewCoordinatorServiceClient(conn)
	return c
}

func TestFrequentPattern(t *testing.T) {

	c := establishConnection(t)

	inputFlag := "-c"
	inputString := "Mozilla"
	expectedOutput := []int{269198, 254504, 255377, 257082, 257504, 255625, 254744, 260674, 256399, 251984}

	coordinatorOutput, duration := FetchOutput(inputFlag, inputString, c, 1)

	t.Logf("Duration of TestFrequentPattern: %d", duration)

	for i := 0; i < len(coordinatorOutput.FileName); i++ {
		log.Printf("%s\t\t%s\n", coordinatorOutput.FileName[i], coordinatorOutput.Matches[i])
		intVar, err := strconv.Atoi(coordinatorOutput.Matches[i])
		if err != nil {
			t.Logf("Error converting matches to int")
		} else {
			if coordinatorOutput.FileName[i] != "" && intVar != expectedOutput[i] {
				t.Errorf("Match of server %d was incorrect, got: %d, want: %d.", i+1,
					intVar, expectedOutput[i])
			}
		}
	}
}

func TestSomewhatFrequentPattern(t *testing.T) {

	c := establishConnection(t)

	inputFlag := "-c"
	inputString := "iPod; U;" //GET /apps/cart.jsp?appID=5611 HTTP/1.0
	expectedOutput := []int{9348, 9018, 8823, 9161, 9070, 8934, 9006, 9033, 8996, 9012}

	coordinatorOutput, duration := FetchOutput(inputFlag, inputString, c, 1)

	t.Logf("Duration of TestSomewhatFrequentPattern: %d", duration)

	for i := 0; i < len(coordinatorOutput.FileName); i++ {
		log.Printf("%s\t\t%s\n", coordinatorOutput.FileName[i], coordinatorOutput.Matches[i])
		intVar, err := strconv.Atoi(coordinatorOutput.Matches[i])
		if err != nil {
			t.Logf("Error converting matches to int")
		} else {
			if coordinatorOutput.FileName[i] != "" && intVar != expectedOutput[i] {
				t.Errorf("Match of server %d was incorrect, got: %d, want: %d.", i+1,
					intVar, expectedOutput[i])
			}
		}
	}
}

func TestInfrequentPattern(t *testing.T) {

	c := establishConnection(t)

	inputFlag := "-c"
	inputString := "GET /apps/cart.jsp?appID=5611 HTTP/1.0"
	expectedOutput := []int{0, 0, 4, 3, 0, 4, 3, 5, 2, 2}

	coordinatorOutput, duration := FetchOutput(inputFlag, inputString, c, 1)

	t.Logf("Duration of InfrequentPattern: %d", duration)

	for i := 0; i < len(coordinatorOutput.FileName); i++ {
		log.Printf("%s\t\t%s\n", coordinatorOutput.FileName[i], coordinatorOutput.Matches[i])
		intVar, err := strconv.Atoi(coordinatorOutput.Matches[i])
		if err != nil {
			t.Logf("Error converting matches to int")
		} else {
			if coordinatorOutput.FileName[i] != "" && intVar != expectedOutput[i] {
				t.Errorf("Match of server %d was incorrect, got: %d, want: %d.", i+1,
					intVar, expectedOutput[i])
			}
		}
	}
}
