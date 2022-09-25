package main

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/introducer"
	"CS425/cs-425-mp1/src/membership"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func ping(targets []string) {
	for i := 0; i < len(targets); i++ {
		hostName := strings.Split(targets[i], ":")[0]
		portNum := "8001"
		service := hostName + ":" + portNum
		//RemoteAddr, err := net.ResolveUDPAddr("udp", service)
		conn, err := net.DialTimeout("udp", service, 1*time.Second)

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
		}
		var members []conf.Member
		json.Unmarshal(buffer[:n], &members)
		membershipStruct := membership.Membership{}
		membershipStruct.UpdateMembers(&members)
	}
}

func main() {

	hostName, err := os.Hostname()

	if err != nil {
		log.Println(err)
	}

	go IntroducerServer()
	go Server()

	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				targets := membership.GetTargets()
				if len(targets) > 1 {
					ping(targets)
				} else {
					initMember := conf.Member{
						ProcessId:         hostName + ":" + strconv.FormatInt(time.Now().Unix(), 10),
						State:             "ACTIVE",
						IncarnationNumber: 1,
					}
					*membership.Members = append(*membership.Members, initMember)
				}
			}
		}
	}()
	select {}

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

	_, err := conn.Read(buffer)

	fmt.Println(membership.Members)

	hostName, err := os.Hostname()
	introducer.JoinNetwork(hostName + ":" + strconv.FormatInt(time.Now().Unix(), 10))
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
