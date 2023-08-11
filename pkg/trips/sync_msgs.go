package trips

// SessionContext represents a user's participation
// to the multiplayer collaboration session.
type SessionContext struct {
	// ConnID tracks the connection from an instance of
	// the travelreys client app.
	ConnID string `json:"connID"`

	// TripID represents the trip currently being updated
	TripID string `json:"tripID"`

	// Participating member
	MemberID string `json:"memberID"`
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
	Type     string `json:"type"`
	ConnID   string `json:"connID"`
	TripID   string `json:"tripID"`
	MemberID string `json:"memberID"`
}

const (
	SyncMsgBroadcastTopicPing        = "SyncMsgBroadcastTopicPing"
	SyncMsgBroadcastTopiCursor       = "SyncMsgBroadcastTopicCursor"
	SyncMsgBroadcastTopiFormPresence = "SyncMsgBroadcastTopicFormPresence"
)

type SyncMsgBroadcast struct {
	SyncMsg `json:",inline"`

	Topic string `json:"topic"`

	Ping         *SyncMsgBroadcastPayloadPing         `json:"ping,omitempty"`
	Cursor       *SyncMsgBroadcastPayloadCursor       `json:"cursor,omitempty"`
	FormPresence *SyncMsgBroadcastPayloadFormPresence `json:"formPresence,omitempty"`
}

type SyncMsgBroadcastPayloadPing struct {
}

type SyncMsgBroadcastPayloadCursor struct{}

type SyncMsgBroadcastPayloadFormPresence struct {
	IsActive bool   `json:"isActive"`
	EditPath string `json:"editPath"`
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
	SyncMsgTOBTopicUpdate = "SyncMsgTOBTopicUpdate"

	// Trip
	SyncMsgTOBUpdateOpDeleteTrip        = "SyncMsgTOBUpdateOpDeleteTrip"
	SyncMsgTOBUpdateOpUpdateTripDates   = "SyncMsgTOBUpdateOpUpdateTripDates"
	SyncMsgTOBUpdateOpUpdateTripMembers = "SyncMsgTOBUpdateOpUpdateTripMembers"

	// Lodgings
	SyncMsgTOBUpdateOpAddLodging    = "SyncMsgTOBUpdateOpAddLodging"
	SyncMsgTOBUpdateOpDeleteLodging = "SyncMsgTOBUpdateOpDeleteLodging"
	SyncMsgTOBUpdateOpUpdateLodging = "SyncMsgTOBUpdateOpUpdateLodging"

	// Itinerary
	SyncMsgTOBUpdateOpDeleteActivity              = "SyncMsgTOBUpdateOpDeleteActivity"
	SyncMsgTOBUpdateOpOptimizeItinerary           = "SyncMsgTOBUpdateOpOptimizeItinerary"
	SyncMsgTOBUpdateOpReorderActivityToAnotherDay = "SyncMsgTOBUpdateOpReorderActivityToAnotherDay"
	SyncMsgTOBUpdateOpReorderItinerary            = "SyncMsgTOBUpdateOpReorderItinerary"
	SyncMsgTOBUpdateOpUpdateActivityPlace         = "SyncMsgTOBUpdateOpUpdateActivityPlace"

	// Media
	SyncMsgTOBUpdateOpAddMediaItem = "SyncMsgTOBUpdateOpAddMediaItem"
)

type SyncMsgTOB struct {
	SyncMsg

	Topic   string `json:"topic"`
	Counter uint64 `json:"counter"`

	Join   *SyncMsgTOBPayloadJoin   `json:"join,omitempty"`
	Leave  *SyncMsgTOBPayloadLeave  `json:"leave,omitempty"`
	Update *SyncMsgTOBPayloadUpdate `json:"update,omitempty"`
}

type SyncMsgTOBPayloadJoin struct {
	// Latest Snapshot of the trip
	Trip *Trip `json:"trip"`

	// List of members updated (presence)
	Members []string `json:"members"`
}

type SyncMsgTOBPayloadLeave struct {
}

type SyncMsgTOBPayloadUpdate struct {
	Op  string   `json:"op"`
	Ops []SyncOp `json:"ops"`
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
	ops []SyncOp,
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
