package membership

import (
	"CS425/cs-425-mp1/src/conf"
	"os"
	"strings"
)

var Members []*conf.Member
var IncarnationNumber int = 1
var Self = os.Getenv("my_endpoint")

func UpdateMembers(updatedMembers []*conf.Member) {
	Members = updatedMembers
}

func RespondToPing() []*conf.Member {
	return GetMembers()
}

func GetMembers() []*conf.Member {
	members := Members
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split(members[i].ProcessID, ":")[0]
		if endpoint == Self {
			members[i].IncarnationNumber += 1
		}
	}
	return members
}

func GetTargets() []string {
	members := Members
	targetsMap := make(map[string]interface{})
	targets := []string{}
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split(members[i].ProcessID, ":")[0]
		if endpoint == Self {
			if i == 0 {
				targetsMap[members[len(members)-1].ProcessID] = nil
				targets = append(targets, members[len(members)-1].ProcessID)
			} else {
				targetsMap[members[i-1].ProcessID] = nil
			}
			targetsMap[members[i+1%(len(members))].ProcessID] = nil
			if i+3%(len(members)) > len(members)-1 {
				targetsMap[members[len(members)-1].ProcessID] = nil
			}
		}
	}
	for k, _ := range targetsMap {
		targets = append(targets, k)
	}
	return targets
}
