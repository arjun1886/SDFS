package worker

import (
	"context"
	"os/exec"
	"strings"
)

// Worker - Struct that implements the worker service methods from the pb file
type Worker struct {
	UnimplementedWorkerServiceServer
}

// FetchWorkerOutput - Executes the grep command on the worker and fetches its output
func (w *Worker) FetchWorkerOutput(ctx context.Context, workerInput *WorkerInput) (*WorkerOutput, error) {
	app := "grep"
	arg0 := workerInput.Flag
	arg1 := workerInput.Data
	arg2 := workerInput.LogFileName

	// Executes OS grep call
	cmd := exec.Command(app, arg0, arg1, arg2)

	stdout, _ := cmd.Output()
	matches := strings.TrimSuffix(string(stdout), "\n")
	workerOutput := WorkerOutput{}
	workerOutput.Matches = matches
	workerOutput.FileName = workerInput.LogFileName
	return &workerOutput, nil
}
