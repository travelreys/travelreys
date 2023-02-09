package trips

import (
	"time"

	"github.com/google/uuid"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/flights"
	"github.com/tiinyplanet/tiinyplanet/pkg/images"
	"github.com/tiinyplanet/tiinyplanet/pkg/maps"
)

// Trip Plan

type TripPlan struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`

	CoverImage images.ImageMetadata `json:"coverImage" bson:"coverImage"`
	StartDate  time.Time            `json:"startDate" bson:"startDate"`
	EndDate    time.Time            `json:"endDate" bson:"endDate"`

	// Members
	Creator TripMember            `json:"creator" bson:"creator"`
	Members map[string]TripMember `json:"members" bson:"members"`

	// Logistics
	Notes    string                 `json:"notes" bson:"notes"`
	Flights  map[string]Flight      `json:"flights" bson:"flights"`
	Transits map[string]BaseTransit `json:"transits" bson:"transits"`
	Lodgings map[string]Lodging     `json:"lodgings" bson:"lodgings"`

	// Contents
	Contents  map[string]TripContentList `json:"contents" bson:"contents"`
	Itinerary []ItineraryList            `json:"itinerary" bson:"itinerary"`

	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	IsArchived bool      `json:"isArchived" bson:"isArchived"`

	Labels common.Labels `json:"labels" bson:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags"`
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

		Flights:  map[string]Flight{},
		Transits: map[string]BaseTransit{},
		Lodgings: map[string]Lodging{},

		Contents:  map[string]TripContentList{},
		Itinerary: []ItineraryList{},

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

type BaseTransit struct {
	ID   string `json:"id" bson:"id"`
	Type string `json:"type"`

	ConfirmationID string               `json:"confirmationID" bson:"confirmationID"`
	Notes          string               `json:"notes" bson:"notes"`
	Price          common.PriceMetadata `json:"priceMetadata" bson:"priceMetadata"`

	Tags        common.Tags         `json:"tags" bson:"tags"`
	Labels      common.Labels       `json:"labels" bson:"labels"`
	Attachments []common.FileObject `json:"attachments" bson:"attachments"`
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

	NumGuests      int32                `json:"numGuests" bson:"numGuests"`
	CheckinTime    time.Time            `json:"checkinTime" bson:"checkinTime"`
	CheckoutTime   time.Time            `json:"checkoutTime" bson:"checkoutTime"`
	Price          common.PriceMetadata `json:"priceMetadata" bson:"priceMetadata"`
	ConfirmationID string               `json:"confirmationID" bson:"confirmationID"`
	Notes          string               `json:"notes" bson:"notes"`
	Place          maps.Place           `json:"place" bson:"place"`

	Tags        common.Tags         `json:"tags" bson:"tags"`
	Labels      common.Labels       `json:"labels" bson:"labels"`
	Attachments []common.FileObject `json:"attachments" bson:"attachments"`
}

// Trip Content

type TripContent struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`

	Place maps.Place `json:"place" bson:"place"`
	Notes string     `json:"notes" bson:"notes"`

	Comments []TripContentComment `json:"comments" bson:"comments"`
	Labels   common.Labels        `json:"labels" bson:"labels"`
}

type TripContentList struct {
	ID       string        `json:"id" bson:"id"`
	Name     string        `json:"name" bson:"name"`
	Contents []TripContent `json:"contents" bson:"contents"`
}

func NewTripContentList(name string) TripContentList {
	return TripContentList{
		ID:       uuid.New().String(),
		Name:     name,
		Contents: []TripContent{},
	}
}

type TripContentComment struct {
	ID      string `json:"id" bson:"id"`
	Comment string `json:"comment" bson:"comment"`

	Member TripMember `json:"member" bson:"member"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// Itinerary

type ItineraryContent struct {
	ID                string `json:"id" bson:"id"`
	TripContentListID string `json:"tripContentListId" bson:"tripContentListId"`
	TripContentID     string `json:"tripContentId" bson:"tripContentId"`

	Price     common.PriceMetadata `json:"priceMetadata" bson:"priceMetadata"`
	StartTime time.Time            `json:"startTime" bson:"startTime"`
	EndTime   time.Time            `json:"endTime" bson:"endTime"`

	Labels common.Labels `json:"labels" bson:"labels"`
}

type ItineraryContentList []ItineraryContent

type ItineraryList struct {
	ID          string               `json:"id" bson:"id"`
	Date        time.Time            `json:"date" bson:"date"`
	Description string               `json:"desc" bson:"desc"`
	Contents    ItineraryContentList `json:"contents" bson:"contents"`

	Labels common.Labels `json:"labels" bson:"labels"`
}

func NewItineraryList(date time.Time) ItineraryList {
	return ItineraryList{
		ID:          uuid.New().String(),
		Date:        date,
		Description: "",
		Contents:    ItineraryContentList{},
		Labels:      common.Labels{},
	}
}
