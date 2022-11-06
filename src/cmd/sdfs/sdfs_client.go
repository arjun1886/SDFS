package main

import (
	"CS425/cs-425-mp1/src/sdfs_server"
	"fmt"
)

func main() {
	var arg string
	fmt.Scanf("%s", &arg)
	if arg == "store" {
		fmt.Println(sdfs_server.Store())
	}
}
