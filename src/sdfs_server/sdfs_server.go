package sdfs_server

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/membership"
	"bufio"
	context "context"
	"errors"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var FileNames = &[]string{}

type SdfsServer struct {
	UnimplementedSdfsServerServer
}

func Store() []string {
	files, err := ioutil.ReadDir("../../sdfs_dir")
	if err != nil {
		log.Fatal(err)
	}

	fileNames := []string{}
	for _, file := range files {
		if !strings.Contains(file.Name(), "_ver_") {
			fileNames = append(fileNames, file.Name())
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

	FileNames = &fileNames
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
	fileName = fileName + "_ver_" + strconv.FormatInt(time.Now().Unix(), 10)

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
		return errors.New("Could not complete full write into file")
	}

	err = UpdateFileNames()
	return err
}

func (s *SdfsServer) Get(fileName string, srv SdfsServer_GetServer) error {

	f, err := os.Open("thermopylae.txt")

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
		if strings.Contains(file.Name(), fileNameModified+"_ver_") {
			err := os.Remove(file.Name())
			if err != nil {
				return err
			}
		}
	}
}

func DeleteAllFiles() error {
	files, err := ioutil.ReadDir("../../sdfs_dir")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		err := os.Remove(file.Name())
		if err != nil {
			return err
		}
	}
}

func GetNumVersionsFileNames(fileName string, numVersions int) ([]string, error) {
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
}

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

func Replication(file string) error {
	membershipStruct := membership.Membership{}
	members := membershipStruct.GetMembers()
	sdfsServerStruct := SdfsServer{}
	fileName := strings.Split(file, "_")[0]
	existingReplicas := Ls(fileName)
	flag := 0
	numReplicas := len(existingReplicas)

	if numReplicas == 5 {
		//for i := 0; i < 5; i++ {
		//	sdfsServerStruct.put(fileName, existingReplicas[i])
		//}
		return nil
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
			if mainReplicaIndex == uint32((*members)[i].processId[14]) {
				sdfsServerStruct.put(fileName, (*members)[i])
			}
		}
	}

	// Linearly scan to find and put in next replicas
	j := (int(mainReplicaIndex) + 1) % len(*members)
	for requiredReplicas > 0 {
		if !Contains(existingReplicas, members[j]) {
			sdfsServerStruct.put(fileName, members[j])
			requiredReplicas -= 1
			j = (j + 1) % len(*members)
		}
	}
	return nil
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
	j := (int(mainReplicaIndex)) % len(*members)
	for requiredReplicas > 0 {
		if !Contains(existingReplicas, *members[j]) {
			// sdfsServerStruct.put(fileName, members[j])
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
