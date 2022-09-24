package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

func main() {

	// CLIENT

	CONNECT := "localhost:8002"

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP client is %s\n", c.RemoteAddr().String())
	defer c.Close()

	//for {
	text := "PING"
	data := []byte(text + "\n")
	_, err = c.Write(data)

	if err != nil {
		fmt.Println(err)
		return
	}
	buffer := make([]byte, 1024)
	n, _, err := c.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error here:", err)
		//return
	}
	fmt.Println("Hello!!!")
	fmt.Printf("Reply: %s\n", string(buffer[0:n]))
	//}

	// SERVER

	//workerOutputChan := make([]byte, 1024)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		/*hostname, error := os.Hostname()
		if error != nil {
			panic(error)
		}
		PORT := ":8002"*/
		fmt.Printf("Entering go routine")

		s, err := net.ResolveUDPAddr("udp4", "localhost:8002")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("The UDP server is %s\n", s)
		connection, err := net.ListenUDP("udp4", s)
		//fmt.Printf("The UDP server 2 is %s\n", connection.RemoteAddr().String())
		if err != nil {
			fmt.Println(err)
			return
		}

		defer connection.Close()
		buffer := make([]byte, 1024)
		//rand.Seed(time.Now().Unix())
		fmt.Printf("connection is still active %s\n", connection)
		for {
			fmt.Printf("Entering for loop")
			n, addr, err := connection.ReadFromUDP(buffer)
			fmt.Printf("Read from UDP here")
			fmt.Print("-> ", string(buffer[0:n-1]))

			if strings.TrimSpace(string(buffer[0:n])) == "PING" {
				fmt.Printf("True")
				serverText := "PONG"
				serverData := []byte(serverText + "\n")
				fmt.Printf("data: %s\n", string(serverData))
				_, err = connection.WriteToUDP(serverData, addr)
				// write membership table
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()

}

/*func main() {
	isPartOfNetwork := false
	membershipStruct := membership.Membership{}
	members := membershipStruct.GetMembers()
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split(members[i].ProcessID, ":")[0]
		if endpoint == membership.Self {
			isPartOfNetwork = true
			break
		}
	}

	workerOutputChan := make([]byte, 1024)
		var wg sync.WaitGroup

		for {

			go func() {
				hostname, error := os.Hostname()
				if error != nil {
					panic(error)
				}
				PORT := ":8001"

				s, err := net.ResolveUDPAddr("udp4", PORT)
				if err != nil {
					fmt.Println(err)
					return
				}

				connection, err := net.ListenUDP("udp4", s)
				if err != nil {
					fmt.Println(err)
					return
				}

				defer connection.Close()
				buffer := make([]byte, 1024)
				rand.Seed(time.Now().Unix())

				for {
					n, addr, err := connection.ReadFromUDP(buffer)
					fmt.Print("-> ", string(buffer[0:n-1]))

					if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
						fmt.Println("Exiting UDP server!")
						return
					}

					data := []byte(strconv.Itoa(random(1, 1001)))
					fmt.Printf("data: %s\n", string(data))
					_, err = connection.WriteToUDP(data, addr)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}
	}
}
*/
