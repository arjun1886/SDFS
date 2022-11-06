package sdfs_server

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/membership"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var FileNames = &[]string{}

type sdfsServer struct {
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
			if fileNames[j] == fileName {
				hostNames = append(hostNames, (*members)[i].ProcessId)
				break
			}
		}
	}
	return hostNames
}
