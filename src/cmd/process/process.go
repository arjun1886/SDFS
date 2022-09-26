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
		}
		var members []conf.Member
		json.Unmarshal(buffer[:n], &members)
		membershipStruct := membership.Membership{}
		membershipStruct.UpdateMembers(&members)
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
				targets := membership.GetTargets()
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
		} else if arg == "LIST MEM" {
			membership.PrintMembershipList()
		} else {
			hostName, _ := os.Hostname()
			membership.PrintSelfId(hostName)
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
