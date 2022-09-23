package main

import (
	"CS425/cs-425-mp1/src/membership"
	"strings"
)

func main() {
	isPartOfNetwork := false
	members := membership.GetMembers()
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split(members[i].ProcessID, ":")[0]
		if endpoint == membership.Self {
			isPartOfNetwork = true
			break
		}
	}
}
