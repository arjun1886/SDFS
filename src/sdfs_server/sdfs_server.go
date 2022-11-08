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

func (s *SdfsServer) Put(ctx context.Context, stream SdfsServer_PutServer) error {
	req, err := stream.Recv()
	if err == io.EOF {
		fileObject := conf.FileData{}
		fileObject.FileName = req.GetFileName()
		fileObject.Data = req.GetChunk()
		err := Put(fileObject)
		if err != nil {
			return nil
		}
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

func (s *SdfsServer) Get(fileName string, srv SdfsServer_GetServer) error {

	f, err := os.Open("../../sdfs_dir/" + fileName)

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
		return &deleteOutput, errors.New("Delete Failed")
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

			_, err = f.Write([]byte("_______________________________"))
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
				return errors.New("Failed to get num versions")
			}
		}
	}

	if flag == true {
		return nil
	} else {
		return errors.New("Failed to get num versions")
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

func Replication() error {
	for {
		time.Sleep(1 * time.Second)
		membershipStruct := membership.Membership{}
		members := membershipStruct.GetMembers()

		newFileToServerMapping := map[string][]string{}

		for i := 0; i < len(*members); i++ {
			fileNames := (*members)[i].FileNames
			for n := 0; n < len(fileNames); n++ {
				fileName := strings.Split((fileNames)[n], "_")[0]
				if _, ok := newFileToServerMapping[fileName]; !ok {
					continue
				}
				existingReplicas := Ls(fileName)
				flag := 0
				numReplicas := len(existingReplicas)

				if numReplicas == 5 {
					continue
				}
				requiredReplicas := 5 - numReplicas

				mainReplicaIndex := hash(fileName) % uint32(len(*members))

				// See if primary replica already exists
				for i := 0; i < len(existingReplicas); i++ {
					if mainReplicaIndex == uint32(existingReplicas[i][14]) {
						flag = 1
					}
				}

				// If it doesn't exist, find and put in primary replica
				if flag == 0 {
					for i := 0; i < len(*members); i++ {
						if mainReplicaIndex == uint32((*members)[i].ProcessId[14]) {
							if (*members)[i].State == "ACTIVE" {
								nodeToFileArray := GetReadTargetsInLatestOrder(fileName, 1)
								_ = GetUtil(nodeToFileArray[0].ProcessId, "dummy.txt", nodeToFileArray[0].FileVersion[0])
								PutUtil("dummy.txt", fileName, []string{(*members)[i].ProcessId})
								break
							} else {
								mainReplicaIndex = (mainReplicaIndex + 1) % uint32(len(*members))
							}
						}
					}
				}

				// Linearly scan to find and put in next replicas
				j := (int(mainReplicaIndex) + 1) % len(*members)
				for requiredReplicas > 0 {
					if !Contains(existingReplicas, (*members)[j].ProcessId) && (*members)[j].State == "ACTIVE" {
						nodeToFileArray := GetReadTargetsInLatestOrder(fileName, 1)
						_ = GetUtil(nodeToFileArray[0].ProcessId, "dummy.txt", nodeToFileArray[0].FileVersion[0])
						PutUtil("dummy.txt", fileName, []string{(*members)[j].ProcessId})
						requiredReplicas -= 1
						j = (j + 1) % len(*members)
					}
				}

				// Update global map
				targets := GetReplicaTargets(fileNames[n])
				newFileToServerMapping[fileName] = targets
			}
		}
		membership.FileToServerMapping = newFileToServerMapping
		return nil
	}
}

func GetReplicaTargets(file string) []string {
	membershipStruct := membership.Membership{}
	members := membershipStruct.GetMembers()
	// sdfsServerStruct := sdfsServer{}
	fileName := strings.Split(file, "_")[0]
	existingReplicas := Ls(fileName)
	hostNames := existingReplicas
	numReplicas := len(hostNames)
	if numReplicas == 5 {
		return hostNames
	}

	requiredReplicas := 5 - numReplicas
	mainReplicaIndex := hash(fileName) % uint32(len(*members))

	// should we check for active main replica here?

	j := (int(mainReplicaIndex)) % len(*members)
	for requiredReplicas > 0 {
		if !Contains(existingReplicas, (*members)[j].ProcessId) && (*members)[j].State == "ACTIVE" {
			// sdfsServerStruct.put(fileName, (*members)[j])
			hostNames = append(hostNames, (*members)[j].ProcessId)
			requiredReplicas -= 1
			j = (j + 1) % len(*members)
		}
	}
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

	membershipStruct := membership.Membership{}
	members := membershipStruct.GetMembers()
	// hostNames := []string{}
	// highestFileVersionPerNode := []string{}
	var result []NodeToFiles

	for i := 0; i < len(*members); i++ {
		flag := 0
		fileNames := (*members)[i].FileNames
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
			nodeToFile.FileVersion = fileVersions[:num]
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
				return errors.New("failed to get file from stream")
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
			conn, err := grpc.Dial(target+":8003", grpc.WithInsecure(), grpc.WithTimeout(time.Duration(2000)*time.Millisecond), grpc.WithBlock())
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
