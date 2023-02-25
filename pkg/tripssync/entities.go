package tripssync

import (
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
)

// Connection represents a user's participation
// to the multiplayer collaboration session.
type Connection struct {
	ID     string
	TripID string
	Member trips.Member
}

// Session keeps track of all the participants in current session
type Session struct {
	// Members is a list of members in the current session
	Members trips.MembersList `json:"members"`
}

const (
	OpJoinSession  = "OpJoinSession"
	OpLeaveSession = "OpLeaveSession"
	OpPingSession  = "OpPingSession"
	OpMemberUpdate = "OpMemberUpdate"
	OpUpdateTrip   = "OpUpdateTrip"
)

func isValidMessageType(op string) bool {
	return common.StringContains([]string{
		OpJoinSession,
		OpLeaveSession,
		OpPingSession,
		OpMemberUpdate,
		OpUpdateTrip,
	}, op)
}

type Message struct {
	ConnID  string      `json:"connId"`
	TripID  string      `json:"tripId"`
	Op      string      `json:"op"`
	Counter uint64      `json:"counter"` // Should be monotonically increasing
	Data    MessageData `json:"data"`
}

type MessageData struct {
	JoinSession  MsgDataJoinSession  `json:"joinSession"`
	LeaveSession MsgDataLeaveSession `json:"leaveSession"`
	MemberUpdate MsgDataMemberUpdate `json:"memberUpdate"`
	Ping         MsgDataPing         `json:"ping"`
	UpdateTrip   MsgDataUpdateTrip   `json:"updateTrip"`
}

// MsgDataJoinSession contains the member that joins the session
type MsgDataJoinSession struct {
	trips.Member
}

// MsgDataLeaveSession contains the member that left the session
type MsgDataLeaveSession struct {
	trips.Member
}

func NewMsgLeaveSession(connID, tripID string) Message {
	return Message{
		ConnID: connID,
		TripID: tripID,
		Op:     OpLeaveSession,
		Data: MessageData{
			LeaveSession: MsgDataLeaveSession{},
		},
	}
}

type MsgDataMemberUpdate struct {
	trips.MembersList `json:"members"`
}

func NewMsgMemberUpdate(tripID string, members trips.MembersList) Message {
	return Message{
		TripID: tripID,
		Op:     OpMemberUpdate,
		Data:   MessageData{MemberUpdate: MsgDataMemberUpdate{members}},
	}
}

type MsgDataPing struct{}

type MsgDataUpdateTrip struct {
	Title string               `json:"title"`
	Ops   []common.JSONPatchOp `json:"ops"`
}
