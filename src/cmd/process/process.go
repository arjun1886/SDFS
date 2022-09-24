package main

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/membership"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func ping(targets []string) {
	for i := 0; i < len(targets); i++ {
		hostName := strings.Split(targets[i], ":")[0]
		portNum := "6000"
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
		_, err = conn.Read(buffer)
		if err != nil {
			membershipStruct := membership.Membership{}
			log.Println(err)
		}
		var members []conf.Member
		json.Unmarshal(buffer, &members)
		membershipStruct := membership.Membership{}
		membershipStruct.UpdateMembers(&members)
	}
}
func main() {

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

	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP client : ", addr)
	fmt.Println("Received from UDP client :  ", string(buffer[:n]))

	if err != nil {
		log.Fatal(err)
	}

	// NOTE : Need to specify client address in WriteToUDP() function
	//        otherwise, you will get this error message
	//        write udp : write: destination address required if you use Write() function instead of WriteToUDP()

	membershipStruct := membership.Membership{}
	members := membershipStruct.GetMembers()
	membersByte, err := json.Marshal(members)
	if err != nil {
		fmt.Println()
	}
	// write message back to client
	message := membersByte
	_, err = conn.WriteToUDP(message, addr)

	if err != nil {
		log.Println(err)
	}

}

func Server() {
	hostName := "localhost"
	portNum := "6000"
	service := hostName + ":" + portNum

	udpAddr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		log.Fatal(err)
	}

	// setup listener for incoming UDP connection
	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UDP server up and listening on port 6000")

	defer ln.Close()

	for {
		// wait for UDP client to connect
		handleUDPConnection(ln)
	}

}
