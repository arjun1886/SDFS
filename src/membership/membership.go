package membership

import (
	"CS425/cs-425-mp1/src/conf"
	"os"
	"strings"
	"sync"
)

var Members *[]conf.Member
var IncarnationNumber int = 1
var Self = os.Getenv("my_endpoint")

// Mutex is safe to use concurrently.
type Membership struct {
	mu sync.Mutex
}

func (c *Membership) UpdateMembers(responseMembershipList *[]conf.Member) {
	c.mu.Lock()
	// members := Members
	//flag = 1
	// var selfEndpoint[] = {}
	// var receivingEndpoint[]
	for i := 0; i < len(*Members); i++ {
		// selfEndpoint = strings.Split((*Members)[i].ProcessId, ":")
		/*if flag == 0 {
			flag = 1
			continue
		}*/
		for j := 0; j < len(*responseMembershipList); j++ {
			// receivingEndpoint = strings.Split(members[j].ProcessId, ":")
			if (*Members)[i] == (*responseMembershipList)[j] {
				//flag = 0
				break
			}
			if (*Members)[i].ProcessId != (*responseMembershipList)[j].ProcessId {
				continue
			} else {
				if (*Members)[i].IncarnationNumber > (*responseMembershipList)[j].IncarnationNumber {
					break
				} else if (*Members)[i].IncarnationNumber < (*responseMembershipList)[j].IncarnationNumber {
					(*Members)[i] = (*responseMembershipList)[j] // does inc number also get updated here?
					//flag = 0
					break
				} else {
					if (*Members)[i].State == "ACTIVE" && (*responseMembershipList)[j].State == "FAILED" {
						(*Members)[i].State = "FAILED"
					}
				}
			}
		}

	}

	if len(*responseMembershipList) > len(*Members) {

		for j := 0; j < len(*responseMembershipList); j++ {
			flag := 0
			for i := 0; i < len(*Members); i++ {
				if (*Members)[i].ProcessId == (*responseMembershipList)[j].ProcessId {
					flag = 1
					break
				}
			}
			if flag == 0 {
				*Members = append(*Members, (*responseMembershipList)[j])
			}
		}
	}
	c.mu.Unlock()
}

func (c *Membership) UpdateEntry(processId string, processState string) {
	endpoint := strings.Split(processId, ":")[0]
	c.mu.Lock()
	for i := 0; i < len(*Members); i++ {
		if endpoint == strings.Split((*Members)[i].ProcessId, ":")[0] {
			if processState == "FAILED" {
				(*Members)[i].State = "FAILED"
			}
			if processState == "DELETE" {
				for j = i; j < len(*Members)-1; j++ {
					(*Members)[j] = (*Members)[j+1]
				}
				*Members = (*Members)[:len(*Members)-1]
			}
			break
		}
	}
	c.mu.Unlock()
}

func (c *Membership) GetMembers() *[]conf.Member {
	c.mu.Lock()
	members := *Members
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split((members)[i].ProcessId, ":")[0]
		if endpoint == Self {
			(members)[i].IncarnationNumber += 1
		}
	}
	c.mu.Unlock()
	return &members
}

func GetTargets() []string {
	members := *Members
	targetsMap := make(map[string]interface{})
	targets := []string{}
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split((members)[i].ProcessId, ":")[0]
		if endpoint == Self {
			if i == 0 {
				targetsMap[members[len(members)-1].ProcessId] = nil
				targets = append(targets, members[len(members)-1].ProcessId)
			} else {
				targetsMap[members[i-1].ProcessId] = nil
			}
			targetsMap[members[i+1%(len(members))].ProcessId] = nil
			if i+3%(len(members)) > len(members)-1 {
				targetsMap[members[len(members)-1].ProcessId] = nil
			}
		}
	}
	for k, _ := range targetsMap {
		targets = append(targets, k)
	}
	return targets
}

/*

func main() {
 hostname, error := os.Hostname()
 if error != nil {
  panic(error)
 }
 fmt.Println("hostname returned from Environment : ", hostname)
 fmt.Println("error : ", error)

}
*/
