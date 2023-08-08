package trips

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/finance"
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

	LabelActivityDisplayMediaItem = "displayMediaItem"
	LabelCreatedBy                = "createdBy"
	LabelFractionalIndex          = "fIndex"
	LabelSharingAccess            = "sharing|access"
	LabelUiColor                  = "ui|color"

	// Sharing Access Permissions
	SharingAccessViewer = "view"

	// Media Items
	MediaItemKeyTrip           = "trip"
	MediaItemKeyActivityPrefix = "activity"

	// JSON Patch Paths
	JSONPathItineraryRoot = "/itineraries"
)

const ()

type Trip struct {
	ID   string `json:"id" bson:"id" msgpack:"id"`
	Name string `json:"name" bson:"name" msgpack:"name"`

	CoverImage *CoverImage `json:"coverImage" bson:"coverImage" msgpack:"coverImage"`
	StartDate  time.Time   `json:"startDate" bson:"startDate" msgpack:"startDate"`
	EndDate    time.Time   `json:"endDate" bson:"endDate" msgpack:"endDate"`

	// Members
	Creator   Member            `json:"creator" bson:"creator" msgpack:"creator"`
	Members   MembersMap        `json:"members" bson:"members" msgpack:"members"`
	MembersID map[string]string `json:"membersId" bson:"membersId" msgpack:"membersId"`

	// Logistics
	Notes    string      `json:"notes" bson:"notes" msgpack:"notes"`
	Transits TransitsMap `json:"transits" bson:"transits" msgpack:"transits"`
	Lodgings LodgingsMap `json:"lodgings" bson:"lodgings" msgpack:"lodgings"`
	Budget   Budget      `json:"budget" bson:"budget" msgpack:"budget"`
	Links    LinksMap    `json:"links" bson:"links" msgpack:"links"`

	Itineraries ItineraryMap `json:"itineraries" bson:"itineraries" msgpack:"itineraries"`

	// Media, Attachements
	MediaItems map[string]media.MediaItemList `json:"mediaItems" bson:"mediaItems" msgpack:"mediaItems"`
	Files      FilesMap                       `json:"files" bson:"files" msgpack:"files"`

	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt" msgpack:"updatedAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt" msgpack:"createdAt"`

	Deleted bool          `json:"deleted" bson:"deleted" msgpack:"deleted"`
	Labels  common.Labels `json:"labels" bson:"labels" msgpack:"labels"`
	Tags    common.Tags   `json:"tags" bson:"tags" msgpack:"tags"`
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
		Itineraries: ItineraryMap{},
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
	Source   string               `json:"source" bson:"source" msgpack:"source"`
	WebImage images.ImageMetadata `json:"webImage" bson:"webImage" msgpack:"webImage"`

	// TripImage is the ID of mediaItem in the mediaItems["coverImage"]
	// e.g trip@<id> or activity@<id>
	TripImage string `json:"tripImage" bson:"tripImage" msgpack:"tripImage"`
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
	ID     string            `json:"id" bson:"id" msgpack:"id"`
	Role   string            `json:"role" bson:"role" msgpack:"role"`
	Labels map[string]string `json:"labels" bson:"labels" msgpack:"labels"`

	User auth.User `json:"user" bsomsgpack" bson:"-"`
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
	ID              string            `json:"id" bson:"id" msgpack:"id"`
	Type            string            `json:"type"`
	DepartTime      time.Time         `json:"departTime" bson:"departTime" msgpack:"departTime"`
	DepartLocation  maps.Place        `json:"departLocation" bson:"departLocation" msgpack:"departLocation"`
	ArrivalTime     time.Time         `json:"arrivalTime" bson:"arrivalTime" msgpack:"arrivalTime"`
	ArrivalLocation maps.Place        `json:"arrivalLocation" bson:"arrivalLocation" msgpack:"arrivalLocation"`
	ConfirmationID  string            `json:"confirmationID" bson:"confirmationID" msgpack:"confirmationID"`
	Notes           string            `json:"notes" bson:"notes" msgpack:"notes"`
	PriceItem       finance.PriceItem `json:"price" bson:"price" msgpack:"price"`

	Tags   common.Tags   `json:"tags" bson:"tags" msgpack:"tags"`
	Labels common.Labels `json:"labels" bson:"labels" msgpack:"labels"`
}

type TransitsMap map[string]*BaseTransit

type Lodging struct {
	ID             string            `json:"id" bson:"id" msgpack:"id"`
	NumGuests      int32             `json:"numGuests" bson:"numGuests" msgpack:"numGuests"`
	CheckinTime    time.Time         `json:"checkinTime" bson:"checkinTime" msgpack:"checkinTime"`
	CheckoutTime   time.Time         `json:"checkoutTime" bson:"checkoutTime" msgpack:"checkoutTime"`
	PriceItem      finance.PriceItem `json:"price" bson:"price" msgpack:"price"`
	ConfirmationID string            `json:"confirmationID" bson:"confirmationID" msgpack:"confirmationID"`
	Notes          string            `json:"notes" bson:"notes" msgpack:"notes"`
	Place          maps.Place        `json:"place" bson:"place" msgpack:"place"`
	Tags           common.Tags       `json:"tags" bson:"tags" msgpack:"tags"`
	Labels         common.Labels     `json:"labels" bson:"labels" msgpack:"labels"`
}

type LodgingList []*Lodging
type LodgingsMap map[string]*Lodging

func (m LodgingsMap) GetLodgingsForDate(dt time.Time) LodgingsMap {
	results := LodgingsMap{}
	for _, l := range m {
		// TODO: fix this checkout and checking should be the same then show
		// else only show checkin or staying
		if (dt.Equal(l.CheckinTime) || dt.After(l.CheckinTime)) &&
			(dt.Before(l.CheckoutTime) || dt.Equal(l.CheckoutTime)) {
			results[l.ID] = l
		}
	}
	return results
}

type Budget struct {
	ID     string          `json:"id" bson:"id" msgpack:"id"`
	Amount finance.Price   `json:"amount" bson:"amount" msgpack:"amount"`
	Items  BudgetItemsList `json:"items" bson:"items" msgpack:"items"`

	Labels common.Labels `json:"labels" bson:"labels" msgpack:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags" msgpack:"tags"`
}

func NewBudget() Budget {
	return Budget{
		Amount: finance.Price{},
		Items:  BudgetItemsList{},
		Labels: common.Labels{},
		Tags:   common.Tags{},
	}
}

type BudgetItem struct {
	Title     string            `json:"title" bson:"title" msgpack:"title"`
	Desc      string            `json:"desc" bson:"desc" msgpack:"desc"`
	PriceItem finance.PriceItem `json:"price" bson:"price" msgpack:"price"`
	Labels    common.Labels     `json:"labels" bson:"labels" msgpack:"labels"`
	Tag       common.Tags       `json:"tags" bson:"tags" msgpack:"tags"`
}

type BudgetItemsList []*BudgetItem

type Link struct {
	ID     string        `json:"id" bson:"id" msgpack:"id"`
	Notes  string        `json:"notes" bson:"notes" msgpack:"notes"`
	OGP    ogp.Opengraph `json:"ogp"`
	Labels common.Labels `json:"labels" bson:"labels" msgpack:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags" msgpack:"tags"`
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
