package tripssync

import (
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/jsonpatch"
	"github.com/travelreys/travelreys/pkg/trips"
)

// SessionContext represents a user's participation
// to the multiplayer collaboration session.
type SessionContext struct {
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
	OpUpdateTrip   = "OpUpdateTrip"
)

func isValidMessageType(op string) bool {
	return common.StringContains([]string{
		OpJoinSession,
		OpLeaveSession,
		OpPingSession,
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
	Ping         MsgDataPing         `json:"ping"`
	UpdateTrip   MsgDataUpdateTrip   `json:"updateTrip"`
}

// MsgDataJoinSession contains the member that joins the session
type MsgDataJoinSession struct {
	trips.Member

	// Latest Snapshot of the trip
	Trip trips.Trip `json:"trip"`

	// Updated list of members
	Members trips.MembersList `json:"members"`
}

// MsgDataLeaveSession contains the member that left the session
type MsgDataLeaveSession struct {
	trips.Member
	Members trips.MembersList `json:"members"`
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

type MsgDataPing struct {
	trips.Member
}

func NewMsgPing(connID, tripID string) Message {
	return Message{
		ConnID: connID,
		TripID: tripID,
		Op:     OpPingSession,
		Data: MessageData{
			Ping: MsgDataPing{},
		},
	}
}

const (
	MsgUpdateTripTitleAddMediaItem                = "AddMediaItem"
	MsgUpdateTripTitleUpdateTripDates             = "UpdateTripDates"
	MsgUpdateTripTitleUpdateTripMembers           = "UpdateTripMembers"
	MsgUpdateTripTitleUpdateActivityPlace         = "UpdateActivityPlace"
	MsgUpdateTripTitleUpdateActivityDisplayImage  = "UpdateActivityDisplayImage"
	MsgUpdateTripTitleDeleteActivity              = "DeleteActivity"
	MsgUpdateTripTitleReorderItinerary            = "ReorderItinerary"
	MsgUpdateTripTitleOptimizeItinerary           = "OptimizeItinerary"
	MsgUpdateTripTitleReorderActivityToAnotherDay = "ReorderActivityToAnotherDay"
)

type MsgDataUpdateTrip struct {
	Title string         `json:"title"`
	Ops   []jsonpatch.Op `json:"ops"`
}
