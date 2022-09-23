package introducer

import (
	"CS425/cs-425-mp1/src/conf"
	"CS425/cs-425-mp1/src/membership"
)

func JoinNetwork(processID string) {
	addMember(processID)
}

func addMember(processID string) error {
	members := membership.Members
	newMember := conf.Member{
		ProcessID:         processID,
		State:             "Active",
		IncarnationNumber: 1,
	}
	members = append(members, &newMember)
}
