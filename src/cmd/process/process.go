package main

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/membership"
	"CS425/cs-425-mp1/src/sdfs_server"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func ping(targets []string) {
	for i := 0; i < len(targets); i++ {
		hostName := strings.Split(targets[i], ":")[0]
		portNum := "8001"
		service := hostName + ":" + portNum
		RemoteAddr, err := net.ResolveUDPAddr("udp", service)
		conn, err := net.DialUDP("udp", nil, RemoteAddr)
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		// write a message to server
		message := []byte("PING")

		_, err = conn.Write(message)

		if err != nil {
			log.Println(err)
		}

		// receive message from server
		buffer := make([]byte, 1024)
		// Wait 1 second for the Response from server
		err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, err := conn.Read(buffer)
		if err != nil {
			membershipStruct := membership.Membership{}
			membershipStruct.UpdateEntry(targets[i], "FAILED")
			go membershipStruct.Cleanup(targets[i])
			log.Println(err)
		} else {
			var members []conf.Member
			json.Unmarshal(buffer[:n], &members)
			if len(members) != 0 {
				membershipStruct := membership.Membership{}
				membershipStruct.UpdateMembers(&members)
			}
		}
	}
}

func main() {

	// UDP Server which listens for Ping
	go Server()
	go SdfsServer()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			sdfs_server.UpdateFileNames()
		}
	}()
	ticker := time.NewTicker(1000 * time.Millisecond)
	go sdfs_server.Replication()
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				membershipStruct := membership.Membership{}
				targets := membershipStruct.GetTargets()
				if len(targets) >= 1 {
					ping(targets)
				}
			}
		}
	}()

	for {
		var arg string
		fmt.Scanf("%s", &arg)
		if arg == "JOIN" {
			err := sdfs_server.DeleteAllFiles()
			if err != nil {
				fmt.Println("could not delete files before joining : ", err)
			}
			isPartOfNetwork := false
			membershipStruct := membership.Membership{}
			members := membershipStruct.GetMembers()
			for i := 0; i < len(*members); i++ {
				endpoint := strings.Split((*members)[i].ProcessId, ":")[0]
				if endpoint == membership.Self {
					isPartOfNetwork = true
					break
				}
			}
			if !isPartOfNetwork {
				request, _ := os.Hostname()
				servAddr := conf.IntroducerEndpoint
				tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
				if err != nil {
					println("ResolveTCPAddr failed:", err.Error())
					os.Exit(1)
				}

				conn, err := net.DialTCP("tcp", nil, tcpAddr)
				if err != nil {
					println("Dial failed:", err.Error())
					os.Exit(1)
				}

				_, err = conn.Write([]byte(request))
				if err != nil {
					println("Write to server failed:", err.Error())
					os.Exit(1)
				}

				println("write to server = ", request)
				log.Println("write to server=", request)

				reply := make([]byte, 1024)

				n, err := conn.Read(reply)
				if err != nil {
					println("Write to server failed:", err.Error())
					os.Exit(1)
				}

				err = json.Unmarshal(reply[:n], membership.Members)
				if err != nil {
					log.Println("hii", err)
				}

				membership.PrintMembershipList()
				conn.Close()
			} else {
				fmt.Println("Already in network")
			}
		} else if arg == "LEAVE" {
			membershipStruct := membership.Membership{}
			membershipStruct.LeaveNetwork()
		} else if arg == "LIST_MEM" {
			membership.PrintMembershipList()
		} else if arg == "LIST_SELF" {
			hostName, _ := os.Hostname()
			membership.PrintSelfId(hostName)
		} else if arg == "PUT" {
			var sdfsFileName string
			var command string
			var localFileName string
			fmt.Scanf("%s%s%s", &command, &localFileName, &sdfsFileName)
			targetReplicas := sdfs_server.GetReplicaTargets(sdfsFileName)
			err := sdfs_server.PutUtil(localFileName, sdfsFileName, targetReplicas)
			if err != nil {
				fmt.Println("Failed to perform put call : ", err)
			} else {
				fmt.Println("put call successful")
			}
		} else if arg == "GET" {
			var sdfsFileName string
			var command string
			var localFileName string
			fmt.Scanf("%s%s%s", &command, &sdfsFileName, &localFileName)
			readAck := 2
			nodeToFileArray := sdfs_server.GetReadTargetsInLatestOrder(sdfsFileName, 1)
			if len(nodeToFileArray) == 0 {
				fmt.Println("Length of nodeToFileArray 0")
			}
			flag := true
			for i := 0; i < readAck; i++ {
				fmt.Println("First file of nodeToFileArray:", nodeToFileArray[i].FileVersion[0])
				err := sdfs_server.GetUtil(nodeToFileArray[i].ProcessId, localFileName, nodeToFileArray[i].FileVersion[0])
				if err != nil {
					flag = false
					sdfs_server.ClearFile(localFileName)
				} else {
					flag = true
					break
				}
			}
			if flag == true {
				fmt.Println("get call successful")
			} else {
				fmt.Println("failed to perform get call")
			}
		} else if arg == "DELETE" {
			var sdfsFileName string
			var command string
			fmt.Scanf("%s%s", &command, &sdfsFileName)
			err := sdfs_server.DeleteUtil(sdfsFileName)
			if err != nil {
				fmt.Println("Failed to perform delete : ", err)
			} else {
				fmt.Println("Delete successful")
			}
		} else if arg == "GET_VERSIONS" {
			var sdfsFileName string
			var command string
			var localFileName string
			var numVersions int
			fmt.Scanf("%s%s%d%s", &command, &sdfsFileName, &numVersions, &localFileName)
			readAck := 2
			err := sdfs_server.GetNumVersionsUtil(sdfsFileName, numVersions, localFileName, readAck)
			if err != nil {
				fmt.Println("Failed to get file name versions : ", err)
			} else {
				fmt.Println("Get num versions successful")
			}
		} else if arg == "STORE" {
			fmt.Println(sdfs_server.Store())
		} else if arg == "LS" {
			var sdfsFileName string
			var command string
			fmt.Scanf("%s%s", &command, &sdfsFileName)
			hostNames := sdfs_server.Ls(sdfsFileName)
			fmt.Println(hostNames)
		}
	}

	//select {}

}

func handleUDPConnection(conn *net.UDPConn) {

	buffer := make([]byte, 1024)

	_, addr, err := conn.ReadFromUDP(buffer)

	if err != nil {
		log.Fatal(err)
	}

	membershipStruct := membership.Membership{}
	members := membershipStruct.GetMembers()
	membersByte, err := json.Marshal(members)
	if err != nil {
		fmt.Println(err)
	}
	message := membersByte
	_, err = conn.WriteToUDP(message, addr)

	if err != nil {
		log.Println(err)
	}

}

func Server() {
	hostName, _ := os.Hostname()
	portNum := "8001"
	service := hostName + ":" + portNum

	udpAddr, err := net.ResolveUDPAddr("udp", service)

	if err != nil {
		log.Fatal(err)
	}

	// setup listener for incoming UDP connection
	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	for {
		// wait for UDP client to connect
		handleUDPConnection(ln)
	}

}

func SdfsServer() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8003))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	sdfs_server.RegisterSdfsServerServer(grpcServer, sdfs_server.UnimplementedSdfsServerServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
