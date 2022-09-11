package conf

var (
	CoordinatorEndpoint1  string = "localhost:9003"
	CoordinatorEndpoint2  string = "localhost:9004"
	CoordinatorEndpoint3  string = "localhost:9003"
	CoordinatorEndpoint4  string = "localhost:9004"
	CoordinatorEndpoint5  string = "localhost:9003"
	CoordinatorEndpoint6  string = "localhost:9004"
	CoordinatorEndpoint7  string = "localhost:9003"
	CoordinatorEndpoint8  string = "localhost:9004"
	CoordinatorEndpoint9  string = "localhost:9004"
	CoordinatorEndpoint10 string = "localhost:9004"

	WorkerEndpoint1  string = "localhost:9001"
	WorkerEndpoint2  string = "localhost:9002"
	WorkerEndpoint3  string = ""
	WorkerEndpoint4  string = ""
	WorkerEndpoint5  string = ""
	WorkerEndpoint6  string = ""
	WorkerEndpoint7  string = ""
	WorkerEndpoint8  string = ""
	WorkerEndpoint9  string = ""
	WorkerEndpoint10 string = ""

	LogFileName1  string = "machine.1.log"
	LogFileName2  string = "machine.2.log"
	LogFileName3  string = "machine.3.log"
	LogFileName4  string = "machine.4.log"
	LogFileName5  string = "machine.5.log"
	LogFileName6  string = "machine.6.log"
	LogFileName7  string = "machine.7.log"
	LogFileName8  string = "machine.8.log"
	LogFileName9  string = "machine.9.log"
	LogFileName10 string = "machine.10.log"
)

type CoordinatorConfigs struct {
	Endpoints []string
}

type WorkerConfig struct {
	Endpoint    string
	LogFileName string
}

func GetWorkerConfigs() []WorkerConfig {
	return []WorkerConfig{
		{
			Endpoint:    WorkerEndpoint1,
			LogFileName: LogFileName1,
		},
		{
			Endpoint:    WorkerEndpoint2,
			LogFileName: LogFileName1,
		},
	}
}

func GetCoordinatorConfigs() CoordinatorConfigs {
	return CoordinatorConfigs{
		Endpoints: []string{CoordinatorEndpoint1, CoordinatorEndpoint2},
	}
}

/*func GetWorkerConfigs() []WorkerConfig {
	return []WorkerConfig{
		{
			Endpoint:    WorkerEndpoint1,
			LogFileName: LogFileName1,
		},
		{
			Endpoint:    WorkerEndpoint2,
			LogFileName: LogFileName2,
		},
		{
			Endpoint:    WorkerEndpoint3,
			LogFileName: LogFileName3,
		},
		{
			Endpoint:    WorkerEndpoint4,
			LogFileName: LogFileName4,
		},
		{
			Endpoint:    WorkerEndpoint5,
			LogFileName: LogFileName5,
		},
		{
			Endpoint:    WorkerEndpoint6,
			LogFileName: LogFileName6,
		},
		{
			Endpoint:    WorkerEndpoint7,
			LogFileName: LogFileName7,
		},
		{
			Endpoint:    WorkerEndpoint8,
			LogFileName: LogFileName8,
		},
		{
			Endpoint:    WorkerEndpoint9,
			LogFileName: LogFileName9,
		},
		{
			Endpoint:    WorkerEndpoint10,
			LogFileName: LogFileName10,
		},
	}
}*/

/*func GetLogFileNames() []string {
	return []string{LogFileName1, LogFileName2, LogFileName3, LogFileName4, LogFileName5, LogFileName6, LogFileName7, LogFileName8, LogFileName9, LogFileName10}
}*/

/*func GetWorkerConfigs() []string {
	return []string{WorkerEndpoint1, WorkerEndpoint2, WorkerEndpoint3, WorkerEndpoint4, WorkerEndpoint5, WorkerEndpoint6, WorkerEndpoint7, WorkerEndpoint8, WorkerEndpoint9, WorkerEndpoint10}
}*/
