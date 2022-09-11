package main

import (
	"CS425/cs-425-mp1/src/coordinator"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"testing"
	"time"
)

func establishConnection(t *testing.T) (c coordinator.CoordinatorServiceClient, conn *grpc.ClientConn) {
	conn, err := grpc.Dial("fa22-cs425-6602.cs.illinois.edu:8001", grpc.WithInsecure(), grpc.WithTimeout(time.Duration(2000)*time.Millisecond), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Did not connect to coordinator: %s", err)
	}
	c = coordinator.NewCoordinatorServiceClient(conn)
	return c, conn
}

func TestFrequentPattern(t *testing.T) {

	c, conn := establishConnection(t)
	defer conn.Close()
	inputFlag := "-c"
	inputString := "Mozilla"
	expectedOutput := []int{269198, 254504, 255377, 257082, 257504, 255625, 254744, 260674, 256399, 251984}

	coordinatorOutput, duration := FetchOutput(inputFlag, inputString, c, 1)

	log.Printf("Duration of TestFrequentPattern: %d", duration)

	for i := 0; i < len(coordinatorOutput.FileName); i++ {
		intVar, err := strconv.Atoi(coordinatorOutput.Matches[i])
		if err != nil {
			t.Logf("Error converting matches to int")
		} else {
			if coordinatorOutput.FileName[i] != "" && intVar != expectedOutput[i] {
				t.Fatalf("Match of server %d was incorrect, got: %d, want: %d.", i+1,
					intVar, expectedOutput[i])
			}
		}
	}
}

func TestSomewhatFrequentPattern(t *testing.T) {

	c, conn := establishConnection(t)
	defer conn.Close()

	inputFlag := "-c"
	inputString := "iPod; U;"
	expectedOutput := []int{9348, 9018, 8823, 9161, 9070, 8934, 9006, 9033, 8996, 9012}

	coordinatorOutput, duration := FetchOutput(inputFlag, inputString, c, 1)

	log.Printf("Duration of TestSomewhatFrequentPattern: %d", duration)

	for i := 0; i < len(coordinatorOutput.FileName); i++ {
		intVar, err := strconv.Atoi(coordinatorOutput.Matches[i])
		if err != nil {
			t.Logf("Error converting matches to int")
		} else {
			if coordinatorOutput.FileName[i] != "" && intVar != expectedOutput[i] {
				t.Fatalf("Match of server %d was incorrect, got: %d, want: %d.", i+1,
					intVar, expectedOutput[i])
			}
		}
	}
}

func TestInfrequentPattern(t *testing.T) {

	c, conn := establishConnection(t)
	defer conn.Close()

	inputFlag := "-c"
	inputString := "GET /apps/cart.jsp?appID=5611 HTTP/1.0"
	expectedOutput := []int{0, 0, 4, 3, 0, 4, 3, 5, 2, 2}

	coordinatorOutput, duration := FetchOutput(inputFlag, inputString, c, 1)

	log.Printf("Duration of TestInfrequentPattern: %d", duration)

	for i := 0; i < len(coordinatorOutput.FileName); i++ {
		intVar, err := strconv.Atoi(coordinatorOutput.Matches[i])
		if err != nil {
			t.Logf("Error converting matches to int")
		} else {
			if coordinatorOutput.FileName[i] != "" && intVar != expectedOutput[i] {
				t.Fatalf("Match of server %d was incorrect, got: %d, want: %d.", i+1,
					intVar, expectedOutput[i])
			}
		}
	}
}

func TestRegexPattern(t *testing.T) {

	c, conn := establishConnection(t)
	defer conn.Close()

	inputFlag := "-Ec"
	inputString := "Ap+le|Linux"
	expectedOutput := []int{167943, 158574, 158953, 160474, 160538, 159126, 158859, 161943, 159587, 156940}

	coordinatorOutput, duration := FetchOutput(inputFlag, inputString, c, 1)

	log.Printf("Duration of TestRegexPattern: %d", duration)

	for i := 0; i < len(coordinatorOutput.FileName); i++ {
		intVar, err := strconv.Atoi(coordinatorOutput.Matches[i])
		if err != nil {
			t.Logf("Error converting matches to int")
		} else {
			if coordinatorOutput.FileName[i] != "" && intVar != expectedOutput[i] {
				t.Fatalf("Match of server %d was incorrect, got: %d, want: %d.", i+1,
					intVar, expectedOutput[i])
			}
		}
	}
}

func TestNonExistentPattern(t *testing.T) {

	c, conn := establishConnection(t)
	defer conn.Close()

	inputFlag := "-c"
	inputString := "Hi my name is Simrita"
	expectedOutput := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	coordinatorOutput, duration := FetchOutput(inputFlag, inputString, c, 1)

	log.Printf("Duration of TestNonExistentPattern: %d", duration)

	for i := 0; i < len(coordinatorOutput.FileName); i++ {
		intVar, err := strconv.Atoi(coordinatorOutput.Matches[i])
		if err != nil {
			t.Logf("Error converting matches to int")
		} else {
			if coordinatorOutput.FileName[i] != "" && intVar != expectedOutput[i] {
				t.Fatalf("Match of server %d was incorrect, got: %d, want: %d.", i+1,
					intVar, expectedOutput[i])
			}
		}
	}
}
