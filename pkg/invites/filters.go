package invites

import (
	"errors"
	"net/url"

	"github.com/travelreys/travelreys/pkg/common"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrInvalidFilter = errors.New("invites.ErrInvalidFilter")
)

type ListTripInvitesFilter struct {
	UserID *string
	TripID *string
}

func MakeListTripInvitesFilterFromURLParams(
	params url.Values,
) ListTripInvitesFilter {
	ff := ListTripInvitesFilter{}

	if params.Get("userID") != "" {
		ff.UserID = common.StringPtr(params.Get("userID"))
	}
	if params.Get("tripID") != "" {
		ff.TripID = common.StringPtr(params.Get("tripID"))
	}
	return ff
}

func (f ListTripInvitesFilter) Validate() error {
	_, ok := f.toBSON()
	if !ok {
		return ErrInvalidFilter
	}
	return nil
}

func (f ListTripInvitesFilter) toBSON() (bson.M, bool) {
	bsonM := bson.M{}
	isSet := false
	if f.TripID != nil && *f.TripID != "" {
		bsonM["tripID"] = f.TripID
		isSet = true
	}
	if f.UserID != nil && *f.UserID != "" {
		bsonM["userID"] = f.UserID
		isSet = true
	}
	return bsonM, isSet
}
