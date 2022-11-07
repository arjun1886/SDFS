package main

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/membership"
	"CS425/cs-425-mp1/src/sdfs_server"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

func ping(targets []string) {
	for i := 0; i < len(targets); i++ {
		// fmt.Println("Targets:", targets)

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
	ticker := time.NewTicker(1000 * time.Millisecond)

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
			_ = sdfs_server.DeleteAllFiles()
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
			targetReplicas := sdfs_server.GetTargetReplicas()
			localFileName := ""
			SdfsFileName := ""
			// Channel used to store a max of 5 put outputs
			nodeOutputChan := make(chan sdfs_server.PutOutput, 5)
			nodeOutputList := []sdfs_server.PutOutput{}
			ctx := context.Background()
			var wg sync.WaitGroup
			// Tell the 'wg' WaitGroup how many threads/goroutines
			//	that are about to run concurrently.
			wg.Add(len(targetReplicas))
			for i := 0; i < len(targetReplicas); i++ {
				// Spawn a thread for each iteration in the loop.
				go func(ctx context.Context, localFileName string, target string, nodeOutputChan chan sdfs_server.PutOutput) {
					// At the end of the goroutine, tell the WaitGroup
					//   that another thread has completed.
					defer wg.Done()
					var conn *grpc.ClientConn
					putOutput := &sdfs_server.PutOutput{}
					conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithTimeout(time.Duration(2000)*time.Millisecond), grpc.WithBlock())
					if err != nil {
						putOutput.Success = false
					} else {
						defer conn.Close()
						s := sdfs_server.NewSdfsServerClient(conn)
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
							putInput := sdfs_server.PutInput{}
							putInput.FileName = SdfsFileName
							putInput.Chunk = buf[:num]
							if err := putClient.Send(&putInput); err != nil {
								putOutput.Success = false
							}
						}

						putOutput, err = putClient.CloseAndRecv()
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
				fmt.Println("Write Successful")
			}
		} else if arg == "get" {
			ctx := context.Background()
			var conn *grpc.ClientConn
			target := ""
			localFileName := ""
			sdfsFileName := ""
			getOutput := &sdfs_server.GetOutput{}
			conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithTimeout(time.Duration(2000)*time.Millisecond), grpc.WithBlock())
			if err != nil {
				fmt.Println("Failed to connect to SDFS to get file")
			} else {
				defer conn.Close()
				s := sdfs_server.NewSdfsServerClient(conn)
				getInput := sdfs_server.GetInput{}
				getInput.FileName = sdfsFileName
				stream, err := s.Get(ctx, &getInput)
				if err != nil {
					fmt.Println("Failed to make Get call to SDFS to get file")
				}

				for {
					getOutput, err = stream.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						fmt.Println("Failed to get file from stream")
					}

					f, err := os.OpenFile(localFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
					if err != nil {
						fmt.Println("Failed to open local file to write from stream")
					}

					defer f.Close()

					n, err := f.Write(getOutput.GetChunk())
					if err != nil {
						fmt.Println("Failed to write to local file from stream")
					}

					if n != len(getOutput.GetChunk()) {
						fmt.Println("Could not complete full write into file")
					}
				}

				fmt.Println("Successfully completed Get call to Sdfs")

			}
		} else if arg == "delete" {
			sdfsFileName := ""
			if val, ok := sdfs_server.FileToServerMapping[sdfsFileName]; ok {
				targetReplicas := val
				sdfsFileName := ""
				// Channel used to store a max of 5 delete outputs
				deleteOutputChan := make(chan sdfs_server.DeleteOutput, 5)
				deleteOutputList := []sdfs_server.DeleteOutput{}
				ctx := context.Background()
				var wg sync.WaitGroup
				// Tell the 'wg' WaitGroup how many threads/goroutines
				//	that are about to run concurrently.
				wg.Add(len(targetReplicas))
				for i := 0; i < len(targetReplicas); i++ {
					// Spawn a thread for each iteration in the loop.
					go func(ctx context.Context, target string, fileName, deleteOutputChan chan sdfs_server.DeleteOutput) {
						// At the end of the goroutine, tell the WaitGroup
						//   that another thread has completed.
						defer wg.Done()
						var conn *grpc.ClientConn
						deleteOutput := &sdfs_server.DeleteOutput{}
						conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithTimeout(time.Duration(2000)*time.Millisecond), grpc.WithBlock())
						if err != nil {
							deleteOutput.Success = false
						} else {
							defer conn.Close()
							s := sdfs_server.NewSdfsServerClient(conn)
							deleteInput := sdfs_server.DeleteInput{FileName: sdfsFileName}
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
					fmt.Println("Delete Successful")
				} else {
					fmt.Println("Delete Failed")
				}
			}
		} else if arg == "get-versions" {
			numVersions := 5
			sdfsFileName := ""
			localFileName := ""
			sdfs_server.GetNumVersionsFileNames(sdfsFileName, numVersions)
		} else if arg == "store" {
			fmt.Println(sdfs_server.Store())
		} else if arg == "ls" {
			sdfsFileName := ""
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
