package trips

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/images"
	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/media"
	"github.com/travelreys/travelreys/pkg/ogp"
	"github.com/travelreys/travelreys/pkg/storage"
)

var (
	ErrInvalidTripImageKey = errors.New("trips.ErrInvalidTripImageKey")
)

const (
	LabelDelimeter = "|"

	LabelCreatedBy                = "createdBy"
	LabelFractionalIndex          = "fIndex"
	LabelUiColor                  = "ui|color"
	LabelActivityDisplayMediaItem = "displayMediaItem"

	LabelSharingAccess  = "sharing|access"
	SharingAccessViewer = "view"
)

const (
	MediaItemKeyTrip           = "trip"
	MediaItemKeyActivityPrefix = "activity"
)

type Trip struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`

	CoverImage *CoverImage `json:"coverImage" bson:"coverImage"`
	StartDate  time.Time   `json:"startDate" bson:"startDate"`
	EndDate    time.Time   `json:"endDate" bson:"endDate"`

	// Members
	Creator   Member            `json:"creator" bson:"creator"`
	Members   MembersMap        `json:"members" bson:"members"`
	MembersID map[string]string `json:"membersId" bson:"membersId"`

	// Logistics
	Notes    string      `json:"notes" bson:"notes"`
	Transits TransitsMap `json:"transits" bson:"transits"`
	Lodgings LodgingsMap `json:"lodgings" bson:"lodgings"`
	Budget   Budget      `json:"budget" bson:"budget"`
	Links    LinksMap    `json:"links" bson:"links"`

	Itineraries map[string]Itinerary `json:"itineraries" bson:"itineraries"`

	// Media, Attachements
	MediaItems map[string]media.MediaItemList `json:"mediaItems" bson:"mediaItems"`
	Files      FilesMap                       `json:"files" bson:"files"`

	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`

	Deleted bool          `json:"deleted" bson:"deleted"`
	Labels  common.Labels `json:"labels" bson:"labels"`
	Tags    common.Tags   `json:"tags" bson:"tags"`
}

type TripsList []*Trip

type TripOGP struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CoverImageURL string `json:"coverImageURL"`
	CreatorName   string `json:"creatorName"`
}

func NewTrip(creator Member, name string) *Trip {
	creator.Role = MemberRoleCreator

	return &Trip{
		ID:          uuid.New().String(),
		Name:        name,
		CoverImage:  &CoverImage{},
		StartDate:   time.Time{},
		EndDate:     time.Time{},
		Creator:     creator,
		Members:     MembersMap{},
		MembersID:   map[string]string{},
		Transits:    TransitsMap{},
		Lodgings:    LodgingsMap{},
		Itineraries: map[string]Itinerary{},
		Budget:      NewBudget(),
		Links:       LinksMap{},
		MediaItems: map[string]media.MediaItemList{
			MediaItemKeyTrip: {},
		},
		Files:     FilesMap{},
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		Deleted:   false,
		Labels:    common.Labels{},
		Tags:      common.Tags{},
	}
}

func NewTripWithDates(creator Member, name string, start, end time.Time) *Trip {
	trip := NewTrip(creator, name)
	trip.StartDate = start
	trip.EndDate = end

	return trip
}

func (trip Trip) ToOGP(creatorName, coverImageURL string) TripOGP {
	return TripOGP{
		ID:            trip.ID,
		Name:          trip.Name,
		CoverImageURL: coverImageURL,
		CreatorName:   creatorName,
	}
}

func (trip Trip) IsSharingEnabled() bool {
	value, ok := trip.Labels[LabelSharingAccess]
	if !ok {
		return false
	}
	return value == SharingAccessViewer
}

const (
	CoverImageSourceWeb  = "web"
	CoverImageSourceTrip = "trip"

	TripImageDelimiter = "@"
)

type CoverImage struct {
	Source   string               `json:"source" bson:"source"`
	WebImage images.ImageMetadata `json:"webImage" bson:"webImage"`

	// TripImage is the ID of mediaItem in the mediaItems["coverImage"]
	// e.g trip@<id> or activity@<id>
	TripImage string `json:"tripImage" bson:"tripImage"`
}

func (ci CoverImage) SplitTripImageKey() (string, string, error) {
	tkns := strings.Split(ci.TripImage, TripImageDelimiter)
	if len(tkns) != 2 {
		return "", "", ErrInvalidTripImageKey
	}
	return tkns[0], tkns[1], nil
}

// Members

const (
	MemberRoleCreator      = "creator"
	MemberRoleCollaborator = "collaborator"
	MemberRoleParticipant  = "participant"
)

type Member struct {
	ID     string            `json:"id" bson:"id"`
	Role   string            `json:"role" bson:"role"`
	Labels map[string]string `json:"labels" bson:"labels"`

	User auth.User `json:"user" bson:"-"`
}

type MembersMap map[string]*Member
type MembersList []*Member

func NewCreator(id string) Member {
	return NewMember(id, MemberRoleCreator)
}

func NewMember(id, role string) Member {
	return Member{
		ID:     id,
		Role:   role,
		Labels: common.Labels{},
	}
}

func (t Trip) GetMemberIDs() []string {
	membersIDs := []string{t.Creator.ID}
	for id := range t.Members {
		membersIDs = append(membersIDs, id)
	}
	return membersIDs
}

const (
	TransitTypeFlight = "flight"
	TransitTypeTrain  = "train"
	TransitTypeBus    = "bus"
	TransitTypeOthers = "others"
)

type BaseTransit struct {
	ID              string           `json:"id" bson:"id"`
	Type            string           `json:"type"`
	DepartTime      time.Time        `json:"departTime" bson:"departTime"`
	DepartLocation  maps.Place       `json:"departLocation" bson:"departLocation"`
	ArrivalTime     time.Time        `json:"arrivalTime" bson:"arrivalTime"`
	ArrivalLocation maps.Place       `json:"arrivalLocation" bson:"arrivalLocation"`
	ConfirmationID  string           `json:"confirmationID" bson:"confirmationID"`
	Notes           string           `json:"notes" bson:"notes"`
	PriceItem       common.PriceItem `json:"price" bson:"price"`

	Tags   common.Tags   `json:"tags" bson:"tags"`
	Labels common.Labels `json:"labels" bson:"labels"`
}

type TransitsMap map[string]*BaseTransit

type Lodging struct {
	ID             string           `json:"id" bson:"id"`
	NumGuests      int32            `json:"numGuests" bson:"numGuests"`
	CheckinTime    time.Time        `json:"checkinTime" bson:"checkinTime"`
	CheckoutTime   time.Time        `json:"checkoutTime" bson:"checkoutTime"`
	PriceItem      common.PriceItem `json:"price" bson:"price"`
	ConfirmationID string           `json:"confirmationID" bson:"confirmationID"`
	Notes          string           `json:"notes" bson:"notes"`
	Place          maps.Place       `json:"place" bson:"place"`
	Tags           common.Tags      `json:"tags" bson:"tags"`
	Labels         common.Labels    `json:"labels" bson:"labels"`
}

type LodgingList []*Lodging
type LodgingsMap map[string]*Lodging

func (m LodgingsMap) GetLodgingsForDate(dt time.Time) LodgingsMap {
	results := LodgingsMap{}
	for _, l := range m {
		if (dt.Equal(l.CheckinTime) || dt.After(l.CheckinTime)) &&
			(dt.Before(l.CheckoutTime) || dt.Equal(l.CheckoutTime)) {
			results[l.ID] = l
		}
	}
	return results
}

type Budget struct {
	ID     string          `json:"id" bson:"id"`
	Amount common.Price    `json:"amount" bson:"amount"`
	Items  BudgetItemsList `json:"items" bson:"items"`

	Labels common.Labels `json:"labels" bson:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags"`
}

func NewBudget() Budget {
	return Budget{
		Amount: common.Price{},
		Items:  BudgetItemsList{},
		Labels: common.Labels{},
		Tags:   common.Tags{},
	}
}

type BudgetItem struct {
	Title     string           `json:"title" bson:"title"`
	Desc      string           `json:"desc" bson:"desc"`
	PriceItem common.PriceItem `json:"price" bson:"price"`
	Labels    common.Labels    `json:"labels" bson:"labels"`
	Tag       common.Tags      `json:"tags" bson:"tags"`
}

type BudgetItemsList []*BudgetItem

type Link struct {
	ID     string        `json:"id" bson:"id"`
	Notes  string        `json:"notes" bson:"notes"`
	OGP    ogp.Opengraph `json:"ogp"`
	Labels common.Labels `json:"labels" bson:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags"`
}

type LinksMap map[string]*Link

const (
	AttachmentTypeFile  = "file"
	AttachmentTypeMedia = "media"
)

func MakeActivityMediaItemsKey(actId string) string {
	return fmt.Sprintf(`%s|%s`, MediaItemKeyActivityPrefix, actId)
}

type FilesMap map[string]*storage.Object
