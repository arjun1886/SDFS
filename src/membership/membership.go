package membership

import (
	"CS425/cs-425-mp1/src/conf"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var Members = &[]conf.Member{}
var IncarnationNumber int = 1
var Self, _ = os.Hostname()

// Mutex is safe to use concurrently.
type Membership struct {
	mu sync.Mutex
}

func areMembersEqual(member1, member2 conf.Member) bool {
	return member1.ProcessId == member2.ProcessId && member1.State == member2.State && member1.IncarnationNumber == member2.IncarnationNumber
}

func (c *Membership) UpdateMembers(responseMembershipList *[]conf.Member) {
	//c.mu.Lock()
	finalMembers := []conf.Member{}
	for i := 0; i < len(*Members); i++ {
		for j := 0; j < len(*responseMembershipList); j++ {
			if areMembersEqual((*Members)[i], (*responseMembershipList)[j]) {
				finalMembers = append(finalMembers, (*Members)[i])
				break
			}
			if (*Members)[i].ProcessId != (*responseMembershipList)[j].ProcessId {
				continue
			} else {
				if (*Members)[i].IncarnationNumber > (*responseMembershipList)[j].IncarnationNumber {
					finalMembers = append(finalMembers, (*Members)[i])
					break
				} else if (*Members)[i].IncarnationNumber < (*responseMembershipList)[j].IncarnationNumber {
					finalMembers = append(finalMembers, (*responseMembershipList)[j])
					//flag = 0
					break
				} else {
					if (*Members)[i].State == "ACTIVE" && (*responseMembershipList)[j].State == "FAILED" {
						(*Members)[i].State = "FAILED"
						finalMembers = append(finalMembers, (*Members)[i])
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
				finalMembers = append(finalMembers, (*responseMembershipList)[j])
			}
		}
	}
	Members = &finalMembers
	//c.mu.Unlock()
}

func (c *Membership) UpdateEntry(processId string, processState string) {
	fmt.Printf("Calling UpdateEntry with processId %s", processId)
	endpoint := strings.Split(processId, ":")[0]
	//c.mu.Lock()
	for i := 0; i < len(*Members); i++ {
		if endpoint == strings.Split((*Members)[i].ProcessId, ":")[0] {
			if processState == "FAILED" {
				(*Members)[i].State = "FAILED"
			}
			if processState == "DELETE" {
				for j := i; j < len(*Members)-1; j++ {
					(*Members)[j] = (*Members)[j+1]
				}
				*Members = (*Members)[:len(*Members)-1]
			}
			break
		}
	}
	//c.mu.Unlock()
}

func (c *Membership) Cleanup(processId string) {
	time.Sleep(3 * time.Second)
	membershipStruct := Membership{}
	membershipStruct.UpdateEntry(processId, "DELETE")
}

func (c *Membership) GetMembers() *[]conf.Member {
	//c.mu.Lock()
	members := *Members
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split((members)[i].ProcessId, ":")[0]
		if endpoint == Self {
			(members)[i].IncarnationNumber += 1
		}
	}
	//c.mu.Unlock()
	return &members
}

func GetTargets() []string {
	members := *Members
	targetsMap := make(map[string]interface{})
	var targets []string
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split((members)[i].ProcessId, ":")[0]
		if endpoint == Self && len(members) != 1 {
			var predecessor int
			successor := i + 1
			arbitraryTarget := i + 3
			if i == 0 {
				predecessor = len(members) - 1
			} else {
				predecessor = i - 1
			}
			for members[predecessor].State == "FAILED" {
				predecessor -= 1
				if predecessor < 0 {
					predecessor = len(targets) - 1
				}
			}
			targetsMap[members[predecessor].ProcessId] = nil
			for members[successor%(len(members))].State == "FAILED" {
				if members[successor%(len(members))].ProcessId == Self {
					break
				}
				successor += 1
			}
			targetsMap[members[successor%(len(members))].ProcessId] = nil
			for members[arbitraryTarget%(len(members))].State == "FAILED" {
				if members[arbitraryTarget%(len(members))].ProcessId == Self {
					break
				}
				arbitraryTarget += 1
			}
			targetsMap[members[arbitraryTarget%(len(members))].ProcessId] = nil
		}
	}
	for k := range targetsMap {
		endpoint := strings.Split(k, ":")[0]
		if endpoint != Self {
			targets = append(targets, k)
		}
	}
	return targets
}

func PrintSelfId(hostname string) {
	endpoint := strings.Split(hostname, ":")[0]
	//c.mu.Lock()
	for i := 0; i < len(*Members); i++ {
		if endpoint == strings.Split((*Members)[i].ProcessId, ":")[0] {
			fmt.Printf("Process ID: %s\n", (*Members)[i].ProcessId)
			log.Printf("Process ID: %s\n", (*Members)[i].ProcessId)
			break
		}
	}
}

func PrintMembershipList() {
	fmt.Printf("Process Id\t\tIncarnation Number\t\tState\n")
	log.Printf("Process Id\t\tIncarnation Number\t\tState\n")
	for i := 0; i < len(*Members); i++ {
		fmt.Printf("%s\t\t%d\t\t%s\n", (*Members)[i].ProcessId, (*Members)[i].IncarnationNumber,
			(*Members)[i].State)
		/*log.Printf("%s\t\t%d\t\t%s\n", (*Members)[i].ProcessId, (*Members)[i].IncarnationNumber,
		(*Members)[i].State)*/
	}
}
