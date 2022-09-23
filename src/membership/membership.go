package membership

import (
	"CS425/cs-425-mp1/src/conf"
	"os"
	"strings"
	"sync"
)

var Members []*conf.Member
var IncarnationNumber int = 1
var Self = os.Getenv("my_endpoint")

// Mutex is safe to use concurrently.
type Membership struct {
	mu sync.Mutex
}

func (c *Membership) UpdateMembers(updatedMembers []*conf.Member) {
	c.mu.Lock()
	Members = updatedMembers
	c.mu.Unlock()
}

func (c *Membership) GetMembers() []*conf.Member {
	c.mu.Lock()
	members := Members
	for i := 0; i < len(members); i++ {
		endpoint := strings.Split(members[i].ProcessID, ":")[0]
		if endpoint == Self {
			members[i].IncarnationNumber += 1
		}
	}
	c.mu.Unlock()
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
