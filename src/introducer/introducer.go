package introducer

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/membership"
)

func JoinNetwork(processId string) {
	addMember(processId)
}

// Each member contains process ID, state and incarnation number
func addMember(processId string) {
	members := membership.Members
	newMember := conf.Member{
		ProcessId:         processId,
		State:             "ACTIVE",
		IncarnationNumber: 1,
		FileNames:         []string{},
	}
	*members = append(*members, newMember)
}
