package worker

import (
	"context"
	"log"
	"os/exec"
	"strings"
)

type Worker struct {
	UnimplementedWorkerServiceServer
}

func (w *Worker) FetchWorkerOutput(ctx context.Context, workerInput *WorkerInput) (*WorkerOutput, error) {
	app := "grep"
	arg0 := workerInput.Flag
	arg1 := workerInput.Data
	arg2 := workerInput.LogFileName

	log.Printf(arg0, arg1, arg2)

	cmd := exec.Command(app, arg0, arg1, arg2)

	stdout, _ := cmd.Output()
	matches := strings.TrimSuffix(string(stdout), "\n")
	workerOutput := WorkerOutput{}
	workerOutput.Matches = matches
	workerOutput.FileName = workerInput.LogFileName
	return &workerOutput, nil
}
