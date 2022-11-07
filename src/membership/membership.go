package membership

import (
	"CS425/cs-425-mp1/src/conf"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var Members = &[]conf.Member{}
var IncarnationNumber int = 1
var Self, _ = os.Hostname()
var FileNames = &[]string{}
var FileToServerMapping = map[string][]string{}

// Membership :Using mutex lock for membership tables
type Membership struct {
	mu sync.Mutex
}

// areMembersEqual: function to check if values of struct Member are deeply equal
func areMembersEqual(member1, member2 conf.Member) bool {
	return member1.ProcessId == member2.ProcessId && member1.State == member2.State && member1.IncarnationNumber == member2.IncarnationNumber
}

func UpdateFileNames() {
	members := *Members
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split((members)[i].ProcessId, ":")[0]
		if endpoint == Self {
			(members)[i].FileNames = *FileNames
		}
	}
	return
}

func (c *Membership) UpdateMembers(responseMembershipList *[]conf.Member) {
	c.mu.Lock()
	// finalMembers is originally a copy of self's membership table
	// It will be updated and returned as final updated table for self
	finalMembers := []conf.Member{}
	// log.Printf("Length of own list: %d, length of incoming list: %d",
	//len(*Members), len(*responseMembershipList))
	flag := 0

	for i := 0; i < len(*Members); i++ {
		flag = 0
		for j := 0; j < len(*responseMembershipList); j++ {
			// If members are deeply equal, retain self's entry
			if areMembersEqual((*Members)[i], (*responseMembershipList)[j]) {
				// finalMembers = append(finalMembers, (*Members)[i])
				flag = 1
				break
			}
			// If process ID of incoming table's jth is not the same as self's ith, move to next j
			if (*Members)[i].ProcessId != (*responseMembershipList)[j].ProcessId {
				continue
			} else {
				// higher incarnation number takes higher precedence for every conflicting entry, update and move to next i
				if (*Members)[i].IncarnationNumber > (*responseMembershipList)[j].IncarnationNumber {
					// finalMembers = append(finalMembers, (*Members)[i])
					flag = 1
					break
				} else if (*Members)[i].IncarnationNumber < (*responseMembershipList)[j].IncarnationNumber {
					(*Members)[i].State = (*responseMembershipList)[j].State
					(*Members)[i].IncarnationNumber = (*responseMembershipList)[j].IncarnationNumber
					// finalMembers = append(finalMembers, (*Members)[i])
					flag = 1
					break
				} else {
					// in case of conflict but same incarnation number, only update self's entry when incoming entry says failed
					if (*Members)[i].State == "ACTIVE" && (*responseMembershipList)[j].State == "FAILED" {
						(*Members)[i].State = "FAILED"
						// finalMembers = append(finalMembers, (*Members)[i])
						flag = 1
						break
					}
				}
			}
		}

		// All self's entries which were not updated/affected by incoming membership table are retained as is
		if flag == 0 {
			finalMembers = append(finalMembers, (*Members)[i])
		}

	}

	// If incoming table contains less entries than self's table
	if len(*responseMembershipList) < len(*Members) {
		for i := 0; i < len(*Members); i++ {
			flag := 0
			// For every matching process ID, skip
			for j := 0; j < len(*responseMembershipList); j++ {
				if (*Members)[i].ProcessId == (*responseMembershipList)[j].ProcessId {
					flag = 1
					break
				}
			}
			// For every extra entry in self's table, if entry is marked as failed then delete it
			if flag == 0 && (*Members)[i].State == "FAILED" {
				for j := i; j < len(*Members)-1; j++ {
					(*Members)[j] = (*Members)[j+1]
				}
				*Members = (*Members)[:len(*Members)-1]
			}
		}
	} else {

		// if incoming table contains more entries than self's table
		if len(*responseMembershipList) > len(*Members) {
			for j := 0; j < len(*responseMembershipList); j++ {
				flag := 0
				// For every matching process ID, skip
				for i := 0; i < len(*Members); i++ {
					if (*Members)[i].ProcessId == (*responseMembershipList)[j].ProcessId {
						flag = 1
						break
					}
				}

				// For every extra entry in incoming table, if entry is marked as active then add it
				if flag == 0 && (*responseMembershipList)[j].State != "FAILED" {
					member := conf.Member{}
					member.ProcessId = (*responseMembershipList)[j].ProcessId
					member.State = (*responseMembershipList)[j].State
					member.IncarnationNumber = (*responseMembershipList)[j].IncarnationNumber
					*Members = append(*Members, member)
				}
			}
		}
	}
	// log.Printf("Final length of updated table:%d", len(*Members))
	// Members = &finalMembers
	c.mu.Unlock()
}

// UpdateEntry :function to update entry when a target is detected as failed or removed
func (c *Membership) UpdateEntry(processId string, processState string) {
	log.Printf("Calling UpdateEntry with processId %s and state %s", processId, processState)
	endpoint := strings.Split(processId, ":")[0]
	log.Println("UpdateEntry: Placing lock on membership table")
	c.mu.Lock()
	for i := 0; i < len(*Members); i++ {
		if endpoint == strings.Split((*Members)[i].ProcessId, ":")[0] {
			if processState == "FAILED" {
				log.Println("Updating entry with failed state")
				(*Members)[i].State = "FAILED"
			}
			if processState == "DELETE" {
				log.Println("Deleting entry")
				for j := i; j < len(*Members)-1; j++ {
					(*Members)[j] = (*Members)[j+1]
				}
				*Members = (*Members)[:len(*Members)-1]
			}
			break
		}
	}
	c.mu.Unlock()
	log.Println("UpdateEntry: Removing lock on membership table")
}

// Cleanup : to delete an entry 5 seconds after it is marked as failed
func (c *Membership) Cleanup(processId string) {
	time.Sleep(5 * time.Second)
	membershipStruct := Membership{}
	membershipStruct.UpdateEntry(processId, "DELETE")
}

// GetMembers : everytime self is pinged, increase own incarnation number and return self's membership table
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

// LeaveNetwork : voluntarily mark self's state as failed
func (c *Membership) LeaveNetwork() *[]conf.Member {
	c.mu.Lock()
	members := *Members
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split((members)[i].ProcessId, ":")[0]
		if endpoint == Self {
			(members)[i].State = "FAILED"
		}
	}
	c.mu.Unlock()
	return &members
}

// GetTargets : to get list of neighbours which the process pings - Ring as backbone with one extra target but calculate dynamically using hashing
func (c *Membership) GetTargets() []string {
	c.mu.Lock()
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
				if members[predecessor%(len(members))].ProcessId == Self {
					break
				}
				predecessor -= 1
				if predecessor < 0 {
					predecessor = len(members) - 1
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
	c.mu.Unlock()

	return targets
}

// PrintSelfId : to print self's process ID
func PrintSelfId(hostname string) {
	endpoint := strings.Split(hostname, ":")[0]
	// c.mu.Lock()
	for i := 0; i < len(*Members); i++ {
		if endpoint == strings.Split((*Members)[i].ProcessId, ":")[0] {
			// fmt.Printf("Process ID: %s\n", (*Members)[i].ProcessId)
			log.Printf("Process ID: %s\n", (*Members)[i].ProcessId)
			break
		}
	}
}

// PrintMembershipList : to print self's membership list

func PrintMembershipList() {
	// fmt.Printf("Process Id\t\tIncarnation Number\t\tState\n")
	log.Printf("\t\tProcess Id\t\tIncarnation Number\t\tState\n")
	for i := 0; i < len(*Members); i++ {
		/*fmt.Printf("%s\t\t%d\t\t%s\n", (*Members)[i].ProcessId, (*Members)[i].IncarnationNumber,
		(*Members)[i].State)*/
		log.Printf("%s\t\t%d\t\t%s\n", (*Members)[i].ProcessId, (*Members)[i].IncarnationNumber,
			(*Members)[i].State)
	}
}
