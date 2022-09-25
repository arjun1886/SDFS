package main

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/membership"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func ping(targets []string) {
	for i := 0; i < len(targets); i++ {
		fmt.Println("Targets:", targets)

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
		request := "JOIN"
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

		// fmt.Println(membership.Members)
		membership.PrintMembershipList()
		conn.Close()
	}

	go Server()

	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				targets := membership.GetTargets()
				if len(targets) > 1 {
					ping(targets)
				}
			}
		}
	}()
	select {}

}

func handleUDPConnection(conn *net.UDPConn) {

	// here is where you want to do stuff like read or write to client

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

	fmt.Println("UDP server up and listening on port 8001")

	defer ln.Close()

	for {
		// wait for UDP client to connect
		handleUDPConnection(ln)
	}

}
