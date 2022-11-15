package main

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/introducer"
	"CS425/cs-425-mp1/src/membership"
	"CS425/cs-425-mp1/src/sdfs_server"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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

	hostName, err := os.Hostname()

	if err != nil {
		log.Println(err)
	}

	go IntroducerServer()
	go Server()
	go SdfsServer()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			sdfs_server.UpdateFileNames()
		}
	}()
	go sdfs_server.Replication()

	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				membershipStruct := membership.Membership{}
				targets := membershipStruct.GetTargets()
				if len(targets) >= 1 {
					ping(targets)
				} else if len(*membership.Members) == 0 {
					initMember := conf.Member{
						ProcessId:         hostName + ":" + strconv.FormatInt(time.Now().Unix(), 10),
						State:             "ACTIVE",
						IncarnationNumber: 1,
						FileNames:         []string{},
					}
					*membership.Members = append(*membership.Members, initMember)
				}
			}
		}
	}()

	for {
		var arg string
		fmt.Scanf("%s", &arg)
		if arg == "LEAVE" {
			membershipStruct := membership.Membership{}
			membershipStruct.LeaveNetwork()
		} else if arg == "LIST_MEM" {
			membership.PrintMembershipList()
		} else if arg == "LIST_SELF" {
			hostName, _ := os.Hostname()
			membership.PrintSelfId(hostName)
		}
	}
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
	// write message back to client
	message := membersByte
	_, err = conn.WriteToUDP(message, addr)

	if err != nil {
		log.Println(err)
	}

}

func handleTCPConnection(conn net.Conn) {

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	log.Println(string(buffer[:n]))
	// membership.PrintMembershipList()
	// Introducer allows the process to join the network
	introducer.JoinNetwork(string(buffer[:n]) + ":" + strconv.FormatInt(time.Now().Unix(), 10))
	// membership.PrintMembershipList()
	membersByte, err := json.Marshal(membership.Members)
	if err != nil {
		fmt.Println(err)
	}

	message := membersByte
	_, err = conn.Write(message)

	if err != nil {
		log.Println(err)
	}

}

func IntroducerServer() {
	hostName, _ := os.Hostname()
	portNum := "8002"
	service := hostName + ":" + portNum

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)

	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting TCP conn: ", err.Error())
			os.Exit(1)
		}
		handleTCPConnection(conn)
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

	sdfs_server.RegisterSdfsServerServer(grpcServer, &sdfs_server.SdfsServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
