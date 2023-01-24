package trips

import (
	"time"

	"github.com/google/uuid"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/flights"
	"github.com/tiinyplanet/tiinyplanet/pkg/images"
)

// Trip Plan

type TripPlan struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`

	CoverImage images.ImageMetadata `json:"coverImage" bson:"coverImage"`
	StartDate  time.Time            `json:"startDate" bson:"startDate"`
	EndDate    time.Time            `json:"endDate" bson:"endDate"`
	IsArchived bool                 `json:"isArchived" bson:"isArchived"`

	// Trip Members
	Creator TripMember            `json:"creator" bson:"creator"`
	Members map[string]TripMember `json:"members" bson:"members"`

	// Trip Logistics
	Flights  map[string]Flight      `json:"flights" bson:"flights"`
	Transits map[string]BaseTransit `json:"transits" bson:"transits"`
	Lodgings map[string]Lodging     `json:"lodgings" bson:"lodgings"`

	// Trip Contents
	Contents map[string]TripContentList `json:"contents" bson:"contents"` // Map of trip contents

	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	Labels    common.Labels `json:"labels" bson:"labels"`
	Tags      common.Tags   `json:"tags" bson:"tags"`
}

type TripPlansList []TripPlan

func NewTripPlan(creator TripMember, name string) TripPlan {
	return TripPlan{
		ID:         uuid.New().String(),
		Name:       name,
		CoverImage: images.ImageMetadata{},
		StartDate:  time.Time{},
		EndDate:    time.Time{},

		Creator: creator,
		Members: map[string]TripMember{},

		Contents: map[string]TripContentList{},
		Flights:  map[string]Flight{},
		Transits: map[string]BaseTransit{},
		Lodgings: map[string]Lodging{},

		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),

		IsArchived: false,
		Labels:     common.Labels{},
		Tags:       common.Tags{},
	}
}

func NewTripPlanWithDates(creator TripMember, name string, start, end time.Time) TripPlan {
	plan := NewTripPlan(creator, name)
	plan.StartDate = start
	plan.EndDate = end
	return plan
}

// Members

const (
	TripMemberPermCollaborator = "collaborator"
	TripMemberPermParticipant  = "participant"
)

type TripMember struct {
	MemberID    string `json:"memberID" bson:"memberID"`
	MemberEmail string `json:"memberEmail" bson:"memberEmail"`
}

type TripMembersList []TripMember

// Transit

const (
	TransitTypeFlight = "flight"
)

type TransitCenter struct {
	ID                string             `json:"id" bson:"id"`
	TransitCenterType string             `json:"transitCenterType" bson:"transitCenterType"`
	Name              string             `json:"name" bson:"name"`
	Positioning       common.Positioning `json:"positioning" bson:"positioning"`
	Labels            common.Labels      `json:"labels" bson:"labels"`
}

type BaseTransit struct {
	ID   string `json:"id" bson:"id"`
	Type string `json:"type"`

	ConfirmationID string               `json:"confirmationID" bson:"confirmationID"`
	Notes          string               `json:"notes" bson:"notes"`
	Price          common.PriceMetadata `json:"priceMetadata" bson:"priceMetadata"`
	Tags           common.Tags          `json:"tags" bson:"tags"`
	Labels         common.Labels        `json:"labels" bson:"labels"`
	Attachments    []common.FileObject  `json:"attachments" bson:"attachments"`
}

// Flights

type Flight struct {
	BaseTransit
	ItineraryType string         `json:"itineraryType" bson:"itineraryType"`
	Depart        flights.Flight `json:"depart" bson:"depart"`
	Return        flights.Flight `json:"return" bson:"return"`
}

// Lodging

type Lodging struct {
	ID string `json:"id" bson:"id"`

	NumGuests int32 `json:"numGuests" bson:"numGuests"`

	Cost           float64       `json:"cost" bson:"cost"`
	ConfirmationID string        `json:"confirmationID" bson:"confirmationID"`
	Notes          string        `json:"notes" bson:"notes"`
	Tags           common.Tags   `json:"tags" bson:"tags"`
	Labels         common.Labels `json:"labels" bson:"labels"`

	Positioning  common.Positioning `json:"positiioning" bson:"positiioning"`
	CheckinTime  time.Time          `json:"checkinTime" bson:"checkinTime"`
	CheckoutTime time.Time          `json:"checkoutTime" bson:"checkoutTime"`
}

// Trip Content

type TripContent struct {
	ID          string           `json:"id" bson:"id"`
	ContentType string           `json:"contentType" bson:"contentType"`
	Location    LocationContent  `json:"locationContent" bson:"locationContent"`
	Notes       NoteContent      `json:"noteContent" bson:"noteContent"`
	Checklist   ChecklistContent `json:"checklistContent" bson:"checklistContent"`

	Comments []TripContentComment `json:"comments" bson:"comments"`
	Labels   common.Labels        `json:"labels" bson:"labels"`
}

type TripContentList []TripContent

type LocationContent struct {
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Positioning common.Positioning `json:"positioning" bson:"positioning"`
}

type NoteContent struct {
	ImageURL string `json:"imageURL" bson:"imageURL"`
	Note     string `json:"note" bson:"note"`
}

type ChecklistContent struct {
	ImageURL string   `json:"imageURL" bson:"imageURL"`
	Items    []string `json:"items" bson:"items"`
}

type TripContentComment struct {
	ID      string `json:"id" bson:"id"`
	Comment string `json:"comment" bson:"comment"`

	Member TripMember `json:"member" bson:"member"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// Itinerary

type ItineraryList struct {
	ID       string               `json:"id" bson:"id"`
	Name     string               `json:"name" bson:"name"`
	Date     time.Time            `json:"date" bson:"date"`
	Contents ItineraryContentList `json:"contents" bson:"contents"`
}

type ItineraryContent struct {
	TripContent
	StartTime time.Time `json:"startTime" bson:"startTime"`
	EndTime   time.Time `json:"endTime" bson:"endTime"`
}

type ItineraryContentList []ItineraryContent
