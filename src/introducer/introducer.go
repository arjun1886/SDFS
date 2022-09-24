package introducer

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/membership"
)

func JoinNetwork(processId string) {
	addMember(processId)
}

func addMember(processId string) {
	members := membership.Members
	newMember := conf.Member{
		ProcessId:         processId,
		State:             "Active",
		IncarnationNumber: 1,
	}
	*members = append(*members, newMember)
}
