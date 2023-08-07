package trips

import "github.com/travelreys/travelreys/pkg/jsonpatch"

// SessionContext represents a user's participation
// to the multiplayer collaboration session.
type SessionContext struct {
	// ConnID tracks the connection from an instance of
	// the travelreys client app.
	ConnID string `json:"connID" msgpack:"connID"`

	// TripID represents the trip currently being updated
	TripID string `json:"tripID" msgpack:"tripID"`

	// Participating member
	MemberID string `json:"memberID" msgpack:"memberID"`
}

type SessionContextList []SessionContext

func (l SessionContextList) ToMembers() []string {
	mList := []string{}
	for _, ctx := range l {
		mList = append(mList, ctx.MemberID)
	}
	return mList
}

const (
	SyncMsgTypeBroadcast = "broadcast"
	SyncMsgTypeTOB       = "tob"
)

type SyncMsg struct {
	Type     string `json:"type" msgpack:"type"`
	ConnID   string `json:"connID" msgpack:"connID"`
	TripID   string `json:"tripID" msgpack:"tripID"`
	MemberID string `json:"memberID" msgpack:"memberID"`
}

const (
	SyncMsgBroadcastTopicPing        = "SyncMsgBroadcastTopicPing"
	SyncMsgBroadcastTopiCursor       = "SyncMsgBroadcastTopicCursor"
	SyncMsgBroadcastTopiFormPresence = "SyncMsgBroadcastTopicFormPresence"
)

type SyncMsgBroadcast struct {
	SyncMsg `msgpack:",inline"`

	Topic string `json:"topic" msgpack:"topic"`

	Ping         *SyncMsgBroadcastPayloadPing         `json:"ping,omitempty" msgpack:"ping,omitempty"`
	Cursor       *SyncMsgBroadcastPayloadCursor       `json:"cursor,omitempty" msgpack:"cursor,omitempty"`
	FormPresence *SyncMsgBroadcastPayloadFormPresence `json:"formPresence,omitempty" msgpack:"formPresence,omitempty"`
}

type SyncMsgBroadcastPayloadPing struct {
}

type SyncMsgBroadcastPayloadCursor struct{}

type SyncMsgBroadcastPayloadFormPresence struct {
	IsActive bool   `json:"isActive" msgpack:"isActive"` // toggle on/off
	EditPath string `json:"editPath" msgpack:"editPath"`
}

func MakeSyncMsgBroadcastTopicPing(
	connID,
	tripID string,
	mem string,
) SyncMsgBroadcast {
	return SyncMsgBroadcast{
		SyncMsg: SyncMsg{
			Type:     SyncMsgTypeBroadcast,
			ConnID:   connID,
			TripID:   tripID,
			MemberID: mem,
		},
		Topic: SyncMsgBroadcastTopicPing,
	}
}

const (
	SyncMsgTOBTopicJoin   = "SyncMsgTOBTopicJoin"
	SyncMsgTOBTopicLeave  = "SyncMsgTOBTopicLeave"
	SyncMsgTOBTopicUpdate = "SyncMsgBroadcastTopicUpdate"

	// Trip
	SyncMsgTOBUpdateOpDeleteTrip        = "SyncMsgTOBTopicDeleteTrip"
	SyncMsgTOBUpdateOpUpdateTripDates   = "SyncMsgTOBTopicUpdateTripDates"
	SyncMsgTOBUpdateOpUpdateTripMembers = "SyncMsgTOBTopicUpdateTripMembers"

	// Lodgings
	SyncMsgTOBUpdateOpAddLodging    = "SyncMsgTOBTopicAddLodging"
	SyncMsgTOBUpdateOpDeleteLodging = "SyncMsgTOBTopicDeleteLodging"
	SyncMsgTOBUpdateOpUpdateLodging = "SyncMsgTOBTopicUpdateLodging"

	// Itinerary
	SyncMsgTOBUpdateOpDeleteActivity              = "SyncMsgTOBTopicDeleteActivity"
	SyncMsgTOBUpdateOpOptimizeItinerary           = "SyncMsgTOBTopicOptimizeItinerary"
	SyncMsgTOBUpdateOpReorderActivityToAnotherDay = "SyncMsgTOBTopicReorderActivityToAnotherDay"
	SyncMsgTOBUpdateOpReorderItinerary            = "SyncMsgTOBTopicReorderItinerary"
	SyncMsgTOBUpdateOpUpdateActivityPlace         = "SyncMsgTOBTopicUpdateActivityPlace"

	// Media
	SyncMsgTOBUpdateOpAddMediaItem = "SyncMsgTOBTopicAddMediaItem"
)

type SyncMsgTOB struct {
	SyncMsg

	Topic   string `json:"topic" msgpack:"topic"`
	Counter uint64 `json:"counter" msgpack:"counter"`

	Join   *SyncMsgTOBPayloadJoin   `json:"join,omitempty" msgpack:"join,omitempty"`
	Leave  *SyncMsgTOBPayloadLeave  `json:"leave,omitempty" msgpack:"leave,omitempty"`
	Update *SyncMsgTOBPayloadUpdate `json:"update,omitempty" msgpack:"update,omitempty"`
}

type SyncMsgTOBPayloadJoin struct {
	// Latest Snapshot of the trip
	Trip *Trip `json:"trip" msgpack:"trip"`

	// List of members updated (presence)
	Members []string `json:"members" msgpack:"members"`
}

type SyncMsgTOBPayloadLeave struct {
}

type SyncMsgTOBPayloadUpdate struct {
	Op  string         `json:"op" msgpack:"op"`
	Ops []jsonpatch.Op `json:"ops" msgpack:"ops"`
}

func MakeSyncMsgTOBTopicJoin(
	connID,
	tripID string,
	mem string,
) SyncMsgTOB {
	return SyncMsgTOB{
		SyncMsg: SyncMsg{
			Type:     SyncMsgTypeBroadcast,
			ConnID:   connID,
			TripID:   tripID,
			MemberID: mem,
		},
		Topic: SyncMsgTOBTopicJoin,
	}
}

func MakeSyncMsgTOBTopicLeave(
	connID,
	tripID string,
	mem string,
) SyncMsgTOB {
	return SyncMsgTOB{
		SyncMsg: SyncMsg{
			Type:     SyncMsgTypeBroadcast,
			ConnID:   connID,
			TripID:   tripID,
			MemberID: mem,
		},
		Topic: SyncMsgTOBTopicLeave,
	}
}

func MakeSyncMsgTOBTopicUpdate(
	connID,
	tripID,
	mem,
	op string,
	ops []jsonpatch.Op,
) SyncMsgTOB {
	return SyncMsgTOB{
		SyncMsg: SyncMsg{
			Type:     SyncMsgTypeTOB,
			ConnID:   connID,
			TripID:   tripID,
			MemberID: mem,
		},
		Topic: SyncMsgTOBTopicUpdate,
		Update: &SyncMsgTOBPayloadUpdate{
			Op:  op,
			Ops: ops,
		},
	}
}
