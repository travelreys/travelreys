package trips

import "fmt"

type SyncOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
	From  string      `json:"from,omitempty"`
}

func MakeAddSyncOp(path string, val interface{}) SyncOp {
	return SyncOp{"add", path, val, ""}
}

func MakeRemoveSyncOp(path string, val interface{}) SyncOp {
	return SyncOp{"remove", path, val, ""}
}

func MakeRepSyncOp(path string, val interface{}) SyncOp {
	return SyncOp{"replace", path, val, ""}
}

// SyncMsgTOBUpdateOpUpdateTripMembers
func MakeSyncMsgTOBUpdateOpUpdateTripMembersOps(mem Member) []SyncOp {
	return []SyncOp{
		MakeAddSyncOp(fmt.Sprintf("/members/%s", mem.ID), mem),
		MakeAddSyncOp(fmt.Sprintf("/membersId/%s", mem.ID), mem.ID),
	}
}
