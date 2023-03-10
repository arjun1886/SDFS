package sdfs_server

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/membership"
	"bufio"
	context "context"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	sync "sync"
	"time"

	grpc "google.golang.org/grpc"
)

type SdfsServer struct {
	UnimplementedSdfsServerServer
}

type NodeToFiles struct {
	ProcessId   string
	FileVersion []string
}

func Store() []string {
	files, err := ioutil.ReadDir("../../sdfs_dir")
	if err != nil {
		log.Fatal(err)
	}

	fileNames := []string{}
	for _, file := range files {
		fileName := strings.Split(file.Name(), "_")[0]
		if !Contains(fileNames, fileName) {
			fileNames = append(fileNames, fileName)
		}
	}

	return fileNames
}

func UpdateFileNames() error {
	files, err := ioutil.ReadDir("../../sdfs_dir")
	if err != nil {
		return err
	}

	fileNames := []string{}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	membership.FileNames = &fileNames
	membership.UpdateFileNames()
	return nil
}

func (s *SdfsServer) Put(stream SdfsServer_PutServer) error {
	aggregated_data := []byte{}
	filename := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fileObject := conf.FileData{}
			fileObject.FileName = filename
			fileObject.Data = aggregated_data
			err := Put(fileObject)
			if err != nil {
				return nil
			}
			break
		}

		if filename == "" {
			filename = req.GetFileName()
		}
		aggregated_data = append(aggregated_data, req.GetChunk()...)

	}

	putOutput := PutOutput{}
	putOutput.Success = true
	return stream.SendAndClose(&putOutput)
}

func Put(fileObject conf.FileData) error {

	file := fileObject.FileName
	fileName := strings.Split(file, ".")[0]
	fileName = fileName + "_" + strconv.FormatInt(time.Now().Unix(), 10)

	f, err := os.OpenFile("../../sdfs_dir/"+fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	n, err := f.Write(fileObject.Data)
	if err != nil {
		return err
	}

	if n != len(fileObject.Data) {
		return errors.New("could not complete full write into file")
	}

	err = UpdateFileNames()
	return err
}

func (s *SdfsServer) Get(getInput *GetInput, srv SdfsServer_GetServer) error {

	f, err := os.Open("../../sdfs_dir/" + getInput.FileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	reader := bufio.NewReader(f)
	buf := make([]byte, 1024)

	for {
		n, err := reader.Read(buf)

		if err != nil {

			if err != io.EOF {

				log.Fatal(err)
			}

			break
		}

		getOutput := GetOutput{}
		getOutput.Chunk = buf[0:n]
		err = srv.Send(&getOutput)
		if err != nil {
			return err
		}

	}

	return nil
}

func Delete(fileName string) error {
	fileNameModified := strings.Split(fileName, ".")[0]
	files, err := ioutil.ReadDir("../../sdfs_dir")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), strings.Split(fileNameModified, "_")[0]) {
			err := os.Remove("../../sdfs_dir/" + file.Name())
			if err != nil {
				return err
			}
		}
	}

	UpdateFileNames()
	return nil
}

func (s *SdfsServer) Delete(ctx context.Context, deleteInput *DeleteInput) (*DeleteOutput, error) {
	err := Delete(deleteInput.GetFileName())
	deleteOutput := DeleteOutput{}
	if err != nil {
		deleteOutput.Success = false
		return &deleteOutput, errors.New("delete Failed")
	}
	deleteOutput.Success = true
	return &deleteOutput, nil
}

func DeleteAllFiles() error {
	files, err := ioutil.ReadDir("../../sdfs_dir")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		err := os.Remove("../../sdfs_dir/" + file.Name())
		if err != nil {
			return err
		}
	}
	UpdateFileNames()
	return nil
}

func GetNumVersionsUtil(fileName string, numVersions int, localFileName string, readAck int) error {
	nodeToFilesArray := GetReadTargetsInLatestOrder(fileName, numVersions)

	_ = ClearFile(localFileName)

	flag := true
	for i := 0; i < readAck; i++ {
		flag = true
		for j := 0; j < len(nodeToFilesArray[i].FileVersion); j++ {
			err := GetUtil(nodeToFilesArray[i].ProcessId, localFileName, nodeToFilesArray[i].FileVersion[j])
			if err != nil {
				flag = false
				break
			}

			f, err := os.OpenFile(localFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				flag = false
				break
			}

			defer f.Close()

			_, err = f.Write([]byte("_______________________________\n"))
			if err != nil {
				flag = false
				break
			}
		}
		if flag == true {
			return nil
		} else {
			err := ClearFile(localFileName)
			if err != nil {
				return errors.New("failed to get num versions")
			}
		}
	}

	if flag == true {
		return nil
	} else {
		return errors.New("failed to get num versions")
	}
}

/*func GetNumVersionsFileNames(fileName string, numVersions int) ([]string, error) {
	fileNameModified := strings.Split(fileName, ".")[0]
	files, err := ioutil.ReadDir("../../sdfs_dir")
	if err != nil {
		log.Fatal(err)
	}

	fileNames := []string{}
	for _, file := range files {
		if strings.Contains(file.Name(), fileNameModified+"_ver_") {
			fileNames = append(fileNames, file.Name())
		}
	}

	sort.Slice(fileNames, func(i, j int) bool {
		return fileNames[i] > fileNames[j]
	})
	finalFileNames := []string{}
	for i := 0; i < numVersions-1; i++ {
		finalFileNames = append(finalFileNames, fileNames[i])
	}
	return finalFileNames, nil
}*/

func Ls(fileName string) []string {
	membershipStruct := membership.Membership{}
	members := membershipStruct.GetMembers()
	hostNames := []string{}
	for i := 0; i < len(*members); i++ {
		fileNames := (*members)[i].FileNames
		for j := 0; j < len(fileNames); j++ {
			if strings.Split(fileNames[j], "_")[0] == fileName {
				hostNames = append(hostNames, (*members)[i].ProcessId)
				break
			}
		}
	}
	return hostNames
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func Replication() {

	for {
		time.Sleep(1 * time.Second)
		membershipStruct := membership.Membership{}
		members := membershipStruct.GetMembers()
		//fmt.Println("Before", membership.FileToServerMapping)
		newFileToServerMapping := map[string][]string{}

		for i := 0; i < len(*members); i++ {
			//fmt.Println("Entering for loop for members")
			fileNames := (*members)[i].FileNames

			for n := 0; n < len(fileNames); n++ {
				//fmt.Println("Entering fileNames loop with:", fileNames[n])
				fileName := strings.Split((fileNames)[n], "_")[0]
				if val, ok := newFileToServerMapping[fileName]; !ok {
					newVal := []string{}
					newVal = append(newVal, (*members)[i].ProcessId)
					//fmt.Println("inside if appending this to map:", (*members)[i].ProcessId)
					newFileToServerMapping[fileName] = newVal
					//fmt.Println("if map:", newFileToServerMapping[fileName])
				} else {
					if !Contains(val, (*members)[i].ProcessId) {
						val = append(val, (*members)[i].ProcessId)
						//fmt.Println("inside else appending this to map:", (*members)[i].ProcessId)
						newFileToServerMapping[fileName] = val
						//fmt.Println("else map:", newFileToServerMapping[fileName])
					}
				}
			}
		}

		membership.FileToServerMapping = newFileToServerMapping
		//fmt.Println("After", membership.FileToServerMapping)
		//fmt.Println("Hi", membership.FileToServerMapping)

		for fileName, value := range membership.FileToServerMapping {
			//fmt.Println("Entering loop")
			existingReplicas := value
			targets := GetReplicaTargets(fileName)

			for i := 0; i < len(targets); i++ {
				if !Contains(existingReplicas, targets[i]) {
					nodeToFileArray := GetReadTargetsInLatestOrder(fileName, 1)
					_ = ClearFile("dummy.txt")
					err := GetUtil(nodeToFileArray[0].ProcessId, "dummy.txt", nodeToFileArray[0].FileVersion[0])
					if err != nil {
						//fmt.Println("Error from GetUtil inside Replication")
					}
					err2 := PutUtil("dummy.txt", fileName, []string{targets[i]})
					if err2 != nil {
						//fmt.Println("Error from PutUtil inside Replication")
					}
				}
			}

			/*
				fmt.Println("Length of existing replicas:", len(existingReplicas))
				flag := 0
				numReplicas := len(existingReplicas)
				hostNames := existingReplicas

				if numReplicas == 5 {
					fmt.Println("Already 5 replicas")
					continue
				}
				requiredReplicas := 5 - numReplicas

				mainReplicaIndex := hash(fileName) % 10
				fmt.Println("mainReplicaIndex", mainReplicaIndex)
				// See if primary replica already exists
				for i := 0; i < len(existingReplicas); i++ {
					vm, _ := strconv.Atoi(string(existingReplicas[i][14]))
					if mainReplicaIndex == uint32(vm) {
						flag = 1
						fmt.Println("mainReplicaIndex already a replica")
						break

					}
				}

				// If it doesn't exist, find and put in primary replica

				if flag == 0 {
					fmt.Println("Finding primary replica to put in")
					found := 0
					// k := int(mainReplicaIndex)
					for found == 0 {
						for k := 0; k < len(*members); k++ {
							fmt.Println("Calculated mainReplicaIndex ", int(mainReplicaIndex))
							vm, _ := strconv.Atoi(string((*members)[k].ProcessId[14]))
							fmt.Println("VM no:", vm)
							if !Contains(hostNames, (*members)[k].ProcessId) &&
								int(mainReplicaIndex) == vm && (*members)[k].State == "ACTIVE" {
								nodeToFileArray := GetReadTargetsInLatestOrder(fileName, 1)
								_ = ClearFile("dummy.txt")
								err := GetUtil(nodeToFileArray[0].ProcessId, "dummy.txt", nodeToFileArray[0].FileVersion[0])
								if err != nil {
									fmt.Println("Error from GetUtil inside Replication")
									continue
								}
								err2 := PutUtil("dummy.txt", fileName, []string{(*members)[k].ProcessId})
								if err2 != nil {
									fmt.Println("Error from PutUtil inside Replication")
									continue
								}
								found = 1
								hostNames = append(hostNames, (*members)[k].ProcessId)
								requiredReplicas -= 1
								break
							}
							// k = (k + 1) % len(*members)
						}
						if found == 0 {
							mainReplicaIndex = (mainReplicaIndex + 1) % 10
						}
					}
				}

				/*
					if flag == 0 {
						fmt.Println("Finding primary replica to put in")
						found := 0
						k := int(mainReplicaIndex)
						for found == 0 {
							fmt.Println("Calculated mainReplicaIndex ", k)
							vm, _ := strconv.Atoi(string((*members)[k].ProcessId[14]))
							fmt.Println("VM no:", vm)
							if k == vm && (*members)[k].State == "ACTIVE" {
								nodeToFileArray := GetReadTargetsInLatestOrder(fileName, 1)
								err := GetUtil(nodeToFileArray[0].ProcessId, "dummy.txt", nodeToFileArray[0].FileVersion[0])
								if err != nil {
									fmt.Println("Error from GetUtil inside Replication")
								}
								err2 := PutUtil("dummy.txt", fileName, []string{(*members)[k].ProcessId})
								if err2 != nil {
									fmt.Println("Error from PutUtil inside Replication")
								}
								found = 1
							}
							k = (k + 1) % len(*members)
						}
					}*/
			/*for i := 0; i < len(*members); i++ {
				fmt.Println("Calculated mainReplicaIndex ", mainReplicaIndex)
				fmt.Println("VM no:", uint32((*members)[i].ProcessId[14]))
				if mainReplicaIndex == uint32((*members)[i].ProcessId[14]) {
					fmt.Println("Found matching VM")
					if (*members)[i].State == "ACTIVE" {
						fmt.Println("Found active primary replica")
						nodeToFileArray := GetReadTargetsInLatestOrder(fileName, 1)
						err := GetUtil(nodeToFileArray[0].ProcessId, "dummy.txt", nodeToFileArray[0].FileVersion[0])
						if err != nil {
							fmt.Println("Error from GetUtil inside Replication")
						}
						err2 := PutUtil("dummy.txt", fileName, []string{(*members)[i].ProcessId})
						if err2 != nil {
							fmt.Println("Error from PutUtil inside Replication")
						}
						break
					} else {
						fmt.Println("Not active")
						mainReplicaIndex = (mainReplicaIndex + 1) % uint32(len(*members))
					}
				} else {
					fmt.Println("Not matching")
					mainReplicaIndex = (mainReplicaIndex + 1) % uint32(len(*members))
				}
			}*/
			/*
				// Linearly scan to find and put in next replicas
				j := (int(mainReplicaIndex) + 1) % len(*members)
				for requiredReplicas > 0 {
					if !Contains(existingReplicas, (*members)[j].ProcessId) && (*members)[j].State == "ACTIVE" {
						fmt.Println("Found active replica")
						nodeToFileArray := GetReadTargetsInLatestOrder(fileName, 1)
						fmt.Println("Inside replication: length of nodeToFileArray: ", len(nodeToFileArray))
						err := GetUtil(nodeToFileArray[0].ProcessId, "dummy.txt", nodeToFileArray[0].FileVersion[0])
						if err != nil {
							fmt.Println("Error from GetUtil inside Replication")
						}
						err2 := PutUtil("dummy.txt", fileName, []string{(*members)[j].ProcessId})
						if err2 != nil {
							fmt.Println("Error from PutUtil inside Replication")
						}
						requiredReplicas -= 1
						j = (j + 1) % len(*members)
					}
				}
			*/
			/*
				j := (int(mainReplicaIndex) + 1) % 10
				for requiredReplicas > 0 {
					for k := 0; k < len(*members); k++ {
						fmt.Println("Calculated Index ", j)
						vm, _ := strconv.Atoi(string((*members)[k].ProcessId[14]))
						fmt.Println("VM no:", vm)
						if j == vm && !Contains(hostNames, (*members)[k].ProcessId) && (*members)[k].State == "ACTIVE" {
							fmt.Println("Found active replica")
							nodeToFileArray := GetReadTargetsInLatestOrder(fileName, 1)
							fmt.Println("Inside replication: length of nodeToFileArray: ", len(nodeToFileArray))
							_ = ClearFile("dummy.txt")
							err := GetUtil(nodeToFileArray[0].ProcessId, "dummy.txt", nodeToFileArray[0].FileVersion[0])
							if err != nil {
								fmt.Println("Error from GetUtil inside Replication")
								continue
							}
							err2 := PutUtil("dummy.txt", fileName, []string{(*members)[j].ProcessId})
							if err2 != nil {
								fmt.Println("Error from PutUtil inside Replication")
								continue
							}
							hostNames = append(hostNames, (*members)[k].ProcessId)
							requiredReplicas -= 1
							break
						}
					}
					j = (j + 1) % 10
				}

				//// Update global map
				//targets := GetReplicaTargets(fileNames[n])
				//if len(targets) == 0 {
				//	fmt.Println("Empty targets list")
				//} else {
				//	for j := 0; j < len(targets); j++ {
				//		fmt.Println(targets[j])
				//	}
				//}

				// newFileToServerMapping[fileName] = targets

			*/
		}
	}
}

//fmt.Println("Hello", newFileToServerMapping)
//membership.FileToServerMapping = newFileToServerMapping
//fmt.Println("After", membership.FileToServerMapping)

func GetReplicaTargets(file string) []string {
	membershipStruct := membership.Membership{}
	members := membershipStruct.GetMembers()
	// sdfsServerStruct := sdfsServer{}
	fileName := strings.Split(file, "_")[0]
	existingReplicas := membership.FileToServerMapping[fileName]
	fmt.Println("Inside get replica targets, map:", membership.FileToServerMapping[fileName])
	hostNames := existingReplicas
	numReplicas := len(hostNames)
	if numReplicas == 5 {
		return hostNames
	}
	mainIndex := 0
	if numReplicas == 0 {
		mainIndex = (int)(hash(fileName)) % len(*members)
	}
	requiredReplicas := 5 - numReplicas

	j := mainIndex
	for requiredReplicas > 0 {
		if !Contains(hostNames, (*members)[j].ProcessId) && (*members)[j].State == "ACTIVE" {
			hostNames = append(hostNames, (*members)[j].ProcessId)
			// fmt.Println("HostNames", hostNames)
			requiredReplicas -= 1
		}
		j = (j + 1) % len(*members)
	}
	/*
		requiredReplicas := 5 - numReplicas
		// mainReplicaIndex := hash(fileName) % uint32(len(*members))
		mainReplicaIndex := hash(fileName) % 10
		fmt.Println("mainReplicaIndex", mainReplicaIndex)
		flag := 0
		// See if primary replica already exists
		for i := 0; i < len(existingReplicas); i++ {
			vm, _ := strconv.Atoi(string(existingReplicas[i][14]))
			if mainReplicaIndex == uint32(vm) {
				flag = 1
				fmt.Println("mainReplicaIndex exists")
				break
			}
		}

		if flag == 0 {
			indexHost := ""
			newFlag := 0
			for newFlag == 0 {
				// Find new active primary  replica
				for i := 0; i < len(*members); i++ {
					vm, _ := strconv.Atoi(string((*members)[i].ProcessId[14]))
					if mainReplicaIndex == uint32(vm) && (*members)[i].State == "ACTIVE" {
						indexHost = (*members)[i].ProcessId
						newFlag = 1
						break
					}
				}
				mainReplicaIndex = (mainReplicaIndex + 1) % 10
			}

			if newFlag == 1 {
				requiredReplicas -= 1
				hostNames = append(hostNames, indexHost)
			}
		}

		// mainReplica guy is found

		// find the next guy and then just linearly scan

		nextGuyVM := (mainReplicaIndex + 1) % 10
		nextGuyIndex := 0
		nextGuyFlag := 0
		for nextGuyFlag == 0 {
			// Find new active primary  replica
			for i := 0; i < len(*members); i++ {
				vm, _ := strconv.Atoi(string((*members)[i].ProcessId[14]))
				if !Contains(existingReplicas, (*members)[i].ProcessId) &&
					nextGuyVM == uint32(vm) && (*members)[i].State == "ACTIVE" {
					nextGuyIndex = i
					nextGuyFlag = 1
					break
				}
			}
			nextGuyVM = (nextGuyVM + 1) % 10
		}

		j := nextGuyIndex
		for requiredReplicas > 0 {
			if !Contains(existingReplicas, (*members)[j].ProcessId) &&
				!Contains(hostNames, (*members)[j].ProcessId) && (*members)[j].State == "ACTIVE" {
				fmt.Println("Entering contains condition")
				// sdfsServerStruct.put(fileName, (*members)[j])
				hostNames = append(hostNames, (*members)[j].ProcessId)
				fmt.Println("HostNames", hostNames)

				requiredReplicas -= 1
				fmt.Println("Required replicas", requiredReplicas)

			}
			j = (j + 1) % len(*members)
		}

		// should we check for active main replica here?
		/*
			j := (int(mainReplicaIndex)) % len(*members)
			for requiredReplicas > 0 {
				fmt.Println("j", j)
				fmt.Println("Entering loop")
				if !Contains(existingReplicas, (*members)[j].ProcessId) && (*members)[j].State == "ACTIVE" {
					fmt.Println("Entering contains condition")
					// sdfsServerStruct.put(fileName, (*members)[j])
					hostNames = append(hostNames, (*members)[j].ProcessId)
					fmt.Println("HostNames", hostNames)

					requiredReplicas -= 1
					fmt.Println("Required replicas", requiredReplicas)
					j = (j + 1) % len(*members)
				}
			}
	*/

	fmt.Println("Returning replica targets ", hostNames)
	return hostNames
}

func Contains(list []string, element string) bool {
	var result bool = false
	for _, x := range list {
		if x == element {
			result = true
			break
		}
	}
	return result
}

func GetReadTargetsInLatestOrder(file string, num int) []NodeToFiles {
	fmt.Println("Inside get read targets")
	membershipStruct := membership.Membership{}
	members := membershipStruct.GetMembers()
	// hostNames := []string{}
	// highestFileVersionPerNode := []string{}
	var result []NodeToFiles

	for i := 0; i < len(*members); i++ {
		flag := 0
		fileNames := (*members)[i].FileNames
		for j := 0; j < len(fileNames); j++ {
			print("FileNames for member i :", i, ":", fileNames[j])
		}
		// minTimeStamp := 0
		var fileVersions []string
		// highestFileVersion := ""
		for j := 0; j < len(fileNames); j++ {
			if strings.Split(fileNames[j], "_")[0] == file {
				flag = 1
				fileVersions = append(fileVersions, fileNames[j])
				/*timeStamp, _ := strconv.Atoi(strings.Split(fileNames[j], "_")[1])
				if timeStamp > minTimeStamp {
					minTimeStamp = timeStamp
					highestFileVersion = fileNames[j]
				}*/
				//hostNames = append(hostNames, (*members)[i].ProcessId)
				//break
			}
		}
		if flag == 1 {

			sort.Slice(fileVersions[:], func(i, j int) bool {
				timestampI, _ := strconv.Atoi(strings.Split(fileVersions[i], "_")[1])
				timestampJ, _ := strconv.Atoi(strings.Split(fileVersions[j], "_")[1])
				return timestampI > timestampJ
			})

			nodeToFile := new(NodeToFiles)
			nodeToFile.ProcessId = (*members)[i].ProcessId
			// nodeToFile.FileVersion = fileVersions[:num]
			if num > len(fileVersions) {
				nodeToFile.FileVersion = fileVersions
			} else {
				nodeToFile.FileVersion = fileVersions[:num]
			}
			result = append(result, *nodeToFile)
		}
		// highestFileVersionPerNode = append(highestFileVersionPerNode, highestFileVersion)
	}
	sort.Slice(result[:], func(i, j int) bool {
		timestampI, _ := strconv.Atoi(strings.Split(result[i].FileVersion[0], "_")[1])
		timestampJ, _ := strconv.Atoi(strings.Split(result[j].FileVersion[0], "_")[1])
		return timestampI > timestampJ
	})
	return result
}

func ClearFile(file string) error {
	if err := os.Truncate(file, 0); err != nil {
		return errors.New("failed to truncate: %v")
	}
	return nil
}

func GetUtil(target string, localFileName string, sdfsFileName string) error {
	ctx := context.Background()
	var conn *grpc.ClientConn
	getOutput := &GetOutput{}
	target = strings.Split(target, ":")[0]
	conn, err := grpc.Dial(target+":8003", grpc.WithInsecure(), grpc.WithTimeout(time.Duration(2000)*time.Millisecond), grpc.WithBlock())
	if err != nil {
		return errors.New("failed to connect to SDFS to get file")
	} else {
		defer conn.Close()
		s := NewSdfsServerClient(conn)
		getInput := GetInput{}
		getInput.FileName = sdfsFileName
		stream, err := s.Get(ctx, &getInput)
		if err != nil {
			return errors.New("failed to make Get call to SDFS to get file")
		}

		for {
			getOutput, err = stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return errors.New("failed to get file from stream : " + err.Error())
			}

			f, err := os.OpenFile(localFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				return errors.New("failed to open local file to write from stream")
			}

			defer f.Close()

			n, err := f.Write(getOutput.GetChunk())
			if err != nil {
				return errors.New("failed to write to local file from stream")
			}

			if n != len(getOutput.GetChunk()) {
				return errors.New("could not complete full write into file")
			}
		}

		return nil
	}
}

func DeleteUtil(sdfsFileName string) error {
	if val, ok := membership.FileToServerMapping[sdfsFileName]; ok {
		targetReplicas := val
		fmt.Println("Inside Delete Util")
		fmt.Println("Target Replicas", targetReplicas)
		fmt.Println("file to server mapping", membership.FileToServerMapping[sdfsFileName])
		// Channel used to store a max of 5 delete outputs
		deleteOutputChan := make(chan DeleteOutput, 5)
		deleteOutputList := []DeleteOutput{}
		ctx := context.Background()
		var wg sync.WaitGroup
		// Tell the 'wg' WaitGroup how many threads/goroutines
		//	that are about to run concurrently.
		wg.Add(len(targetReplicas))
		for i := 0; i < len(targetReplicas); i++ {
			// Spawn a thread for each iteration in the loop.
			go func(ctx context.Context, target string, fileName string, deleteOutputChan chan DeleteOutput) {
				// At the end of the goroutine, tell the WaitGroup
				//   that another thread has completed.
				defer wg.Done()
				var conn *grpc.ClientConn
				deleteOutput := &DeleteOutput{}
				target = strings.Split(target, ":")[0]
				conn, err := grpc.Dial(target+":8003", grpc.WithInsecure(), grpc.WithTimeout(time.Duration(2000)*time.Millisecond), grpc.WithBlock())
				if err != nil {
					deleteOutput.Success = false
				} else {
					defer conn.Close()
					s := NewSdfsServerClient(conn)
					deleteInput := DeleteInput{FileName: sdfsFileName}
					deleteOutput, err = s.Delete(ctx, &deleteInput)
					if err != nil {
						deleteOutput.Success = false
					}
				}
				deleteOutputChan <- *deleteOutput
			}(ctx, targetReplicas[i], sdfsFileName, deleteOutputChan)
			deleteOutputList = append(deleteOutputList, <-deleteOutputChan)
		}
		// Wait for `wg.Done()` to be executed the number of times
		//   specified in the `wg.Add()` call.
		// `wg.Done()` should be called the exact number of times
		//   that was specified in `wg.Add()`.
		wg.Wait()
		close(deleteOutputChan)
		successCount := 0
		for i := 0; i < len(deleteOutputList); i++ {
			if deleteOutputList[i].Success == true {
				successCount += 1
			}
		}
		if successCount == len(targetReplicas) {
			return nil
		} else {
			return errors.New("delete Failed")
		}
	}
	return errors.New("no files to delete")
}

func PutUtil(localFileName, sdfsFileName string, targetReplicas []string) error {
	// Channel used to store a max of 5 put outputs
	nodeOutputChan := make(chan PutOutput, 5)
	nodeOutputList := []PutOutput{}
	ctx := context.Background()
	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//	that are about to run concurrently.
	wg.Add(len(targetReplicas))
	for i := 0; i < len(targetReplicas); i++ {
		// Spawn a thread for each iteration in the loop.
		go func(ctx context.Context, localFileName string, target string, nodeOutputChan chan PutOutput) {
			// At the end of the goroutine, tell the WaitGroup
			//   that another thread has completed.
			defer wg.Done()
			var conn *grpc.ClientConn
			putOutput := &PutOutput{}
			target = strings.Split(target, ":")[0]
			conn, err := grpc.Dial(target+":8003", grpc.WithInsecure(), grpc.WithTimeout(time.Duration(5000)*time.Millisecond), grpc.WithBlock())
			if err != nil {
				putOutput.Success = false
			} else {
				defer conn.Close()
				s := NewSdfsServerClient(conn)
				putClient, err := s.Put(ctx)

				fil, err := os.Open(localFileName)
				if err != nil {
					putOutput.Success = false
				}

				// Maximum 1KB size per stream.
				buf := make([]byte, 1024)

				for {
					num, err := fil.Read(buf)
					if err == io.EOF {
						break
					}
					if err != nil {
						putOutput.Success = false
					}
					putInput := PutInput{}
					putInput.FileName = sdfsFileName
					putInput.Chunk = buf[:num]
					if err := putClient.Send(&putInput); err != nil {
						putOutput.Success = false
					}
				}

				putOutput, err = putClient.CloseAndRecv()
				fmt.Println(putOutput)
				fmt.Println(err)
				if err != nil {
					putOutput.Success = false
				}
			}
			nodeOutputChan <- *putOutput
		}(ctx, localFileName, targetReplicas[i], nodeOutputChan)
		nodeOutputList = append(nodeOutputList, <-nodeOutputChan)
	}
	// Wait for `wg.Done()` to be executed the number of times
	//   specified in the `wg.Add()` call.
	// `wg.Done()` should be called the exact number of times
	//   that was specified in `wg.Add()`.
	wg.Wait()
	close(nodeOutputChan)
	successCount := 0
	for i := 0; i < len(nodeOutputList); i++ {
		if nodeOutputList[i].Success == true {
			successCount += 1
		}
	}
	if successCount >= 4 {
		return nil
	} else {
		return errors.New("write failed to be acked by W nodes")
	}
}
