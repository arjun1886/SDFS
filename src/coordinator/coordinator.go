package coordinator

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/worker"
	context "context"
	sync "sync"

	grpc "google.golang.org/grpc"
)

type Coordinator struct {
	UnimplementedCoordinatorServiceServer
}

func (c *Coordinator) FetchCoordinatorOutput(ctx context.Context, coordinatorInput *CoordinatorInput) (*CoordinatorOutput, error) {

	workerOutputs, err := FetchWorkerOutputs(ctx, &worker.WorkerInput{Data: coordinatorInput.Data, Flag: coordinatorInput.Flag})
	if err != nil {
		return nil, err
	}
	coordinatorOutput := CoordinatorOutput{}
	fileNameList := []string{}
	matchesList := []string{}

	for i := 0; i < len(workerOutputs); i++ {
		fileNameList = append(fileNameList, workerOutputs[i].GetFileName())
		matchesList = append(matchesList, workerOutputs[i].GetMatches())
	}
	coordinatorOutput.FileName = fileNameList
	coordinatorOutput.Matches = matchesList
	return &coordinatorOutput, nil
}

func FetchWorkerOutputs(ctx context.Context, workerInput *worker.WorkerInput) ([]worker.WorkerOutput, error) {
	var workerOutputs []worker.WorkerOutput
	workerConfigs := conf.GetWorkerConfigs()
	workerOutputChan := make(chan worker.WorkerOutput, 10)
	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//	that are about to run concurrently.
	wg.Add(len(workerConfigs))
	for i := 0; i < len(workerConfigs); i++ {
		// Spawn a thread for each iteration in the loop.
		go func(ctx context.Context, workerInput *worker.WorkerInput, workerConfig conf.WorkerConfig, workerOutputChan chan worker.WorkerOutput) {
			// At the end of the goroutine, tell the WaitGroup
			//   that another thread has completed.
			defer wg.Done()
			var conn *grpc.ClientConn
			conn, err := grpc.Dial(workerConfig.Endpoint, grpc.WithInsecure())
			if err != nil {
				//log.Printf("Failed to connect to server %d: %s", i, err)
				//return
			}
			defer conn.Close()
			w := worker.NewWorkerServiceClient(conn)
			workerInput.LogFileName = workerConfig.LogFileName
			// take input from user to a list of args
			workerOutput, err := w.FetchWorkerOutput(ctx, workerInput)
			if err != nil {
				//log.Printf("Error fetching input from server %d: %s", i, err)
				//return
			}
			workerOutputChan <- *workerOutput
		}(ctx, workerInput, workerConfigs[i], workerOutputChan)
		workerOutputs = append(workerOutputs, <-workerOutputChan)
	}
	// Wait for `wg.Done()` to be executed the number of times
	//   specified in the `wg.Add()` call.
	// `wg.Done()` should be called the exact number of times
	//   that was specified in `wg.Add()`.
	wg.Wait()
	close(workerOutputChan)

	return workerOutputs, nil
}
