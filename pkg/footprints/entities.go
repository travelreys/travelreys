package footprints

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/trips"
)

var (
	ErrCheckinTripAlready = errors.New("footprint.entity.checkinTripAlready")
)

type Footprint struct {
	ID      string `json:"id" bson:"id"`
	UserID  string `json:"userID" bson:"userID"`
	PlaceID string `json:"placeID" bson:"placeID"`

	Country string     `json:"country"`
	City    string     `json:"city"`
	Place   maps.Place `json:"place" bson:"place"`

	// Map of TripID to ActivityIDs
	Trips map[string][]string `json:"trips" bson:"trips"`

	// Map of ActivityIDs to their checked in time
	CheckedIns    map[string]time.Time `json:"checkedIns" bson:"checkedIns"`
	LastCheckedIn time.Time            `json:"lastCheckedIn" bson:"lastCheckedIn"`

	Labels common.Labels `json:"labels" bson:"labels"`
}

type FootprintList []Footprint

func NewFootprint(userID, tripID string, activity trips.Activity) Footprint {
	return Footprint{
		ID:            uuid.NewString(),
		UserID:        userID,
		PlaceID:       activity.Place.Labels[maps.LabelPlaceID],
		Country:       activity.Place.Labels[maps.LabelCountry],
		City:          activity.Place.Labels[maps.LabelCity],
		Place:         activity.Place,
		LastCheckedIn: time.Now(),
		Trips: map[string][]string{
			tripID: {activity.ID},
		},
		CheckedIns: map[string]time.Time{
			activity.ID: time.Now(),
		},
		Labels: common.Labels{},
	}
}

func (fp *Footprint) AddNewCheckin(tripID string, activity trips.Activity) error {
	if _, ok := fp.CheckedIns[activity.ID]; ok {
		return ErrCheckinTripAlready
	}

	fp.LastCheckedIn = time.Now()
	fp.CheckedIns[activity.ID] = fp.LastCheckedIn
	fp.Trips[tripID] = append(fp.Trips[tripID], activity.ID)

	return nil
}
