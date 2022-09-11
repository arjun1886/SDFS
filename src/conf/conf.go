package conf

var (
	CoordinatorEndpoint1  string = "simrita2@fa22-cs425-6601.cs.illinois.edu:8001"
	CoordinatorEndpoint2  string = "simrita2@fa22-cs425-6602.cs.illinois.edu:8001"
	CoordinatorEndpoint3  string = "simrita2@fa22-cs425-6603.cs.illinois.edu:8001"
	CoordinatorEndpoint4  string = "simrita2@fa22-cs425-6604.cs.illinois.edu:8001"
	CoordinatorEndpoint5  string = "simrita2@fa22-cs425-6605.cs.illinois.edu:8001"
	CoordinatorEndpoint6  string = "simrita2@fa22-cs425-6606.cs.illinois.edu:8001"
	CoordinatorEndpoint7  string = "simrita2@fa22-cs425-6607.cs.illinois.edu:8001"
	CoordinatorEndpoint8  string = "simrita2@fa22-cs425-6608.cs.illinois.edu:8001"
	CoordinatorEndpoint9  string = "simrita2@fa22-cs425-6609.cs.illinois.edu:8001"
	CoordinatorEndpoint10 string = "simrita2@fa22-cs425-6610.cs.illinois.edu:8001"

	WorkerEndpoint1  string = "simrita2@fa22-cs425-6601.cs.illinois.edu:8000"
	WorkerEndpoint2  string = "simrita2@fa22-cs425-6602.cs.illinois.edu:8000"
	WorkerEndpoint3  string = "simrita2@fa22-cs425-6603.cs.illinois.edu:8000"
	WorkerEndpoint4  string = "simrita2@fa22-cs425-6604.cs.illinois.edu:8000"
	WorkerEndpoint5  string = "simrita2@fa22-cs425-6605.cs.illinois.edu:8000"
	WorkerEndpoint6  string = "simrita2@fa22-cs425-6606.cs.illinois.edu:8000"
	WorkerEndpoint7  string = "simrita2@fa22-cs425-6607.cs.illinois.edu:8000"
	WorkerEndpoint8  string = "simrita2@fa22-cs425-6608.cs.illinois.edu:8000"
	WorkerEndpoint9  string = "simrita2@fa22-cs425-6609.cs.illinois.edu:8000"
	WorkerEndpoint10 string = "simrita2@fa22-cs425-6610.cs.illinois.edu:8000"

	LogFileName1  string = "vm1.log"
	LogFileName2  string = "vm2.log"
	LogFileName3  string = "vm3.log"
	LogFileName4  string = "vm4.log"
	LogFileName5  string = "vm5.log"
	LogFileName6  string = "vm6.log"
	LogFileName7  string = "vm7.log"
	LogFileName8  string = "vm8.log"
	LogFileName9  string = "vm9.log"
	LogFileName10 string = "vm10.log"
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
