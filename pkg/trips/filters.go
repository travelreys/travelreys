package trips

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrInvalidFilter = errors.New("trips.ErrInvalidFilter")
)

type ListFilter struct {
	UserID  *string
	UserIDs []string

	OnlyPublic bool
}

func (ff ListFilter) Validate() error {
	if ff.UserIDs != nil && len(ff.UserIDs) == 0 {
		return ErrInvalidFilter
	}
	if ff.UserID != nil && *ff.UserID == "" {
		return ErrInvalidFilter
	}
	_, ok := ff.toBSON()
	if !ok {
		return ErrInvalidFilter
	}
	return nil
}

func (f ListFilter) toBSON() (bson.M, bool) {
	bsonAnd := bson.A{}
	isSet := false

	if f.UserID != nil {
		bsonUserOr := bson.A{
			bson.M{bsonKeyCreatorId: *f.UserID},
			bson.M{fmt.Sprintf("membersId.%s", *f.UserID): *f.UserID},
		}
		bsonAnd = append(bsonAnd, bson.M{"$or": bsonUserOr})
		isSet = true
	}

	if f.UserIDs != nil && len(f.UserIDs) > 0 {
		bsonUserOr := bson.A{
			bson.M{bsonKeyCreatorId: bson.M{"$in": f.UserIDs}},
		}
		for _, userID := range f.UserIDs {
			bsonUserOr = append(bsonUserOr, bson.M{
				fmt.Sprintf("membersId.%s", userID): userID,
			})
		}
		bsonAnd = append(bsonAnd, bson.M{"$or": bsonUserOr})
		isSet = true
	}

	if f.OnlyPublic {
		sharingAccessLabel := "labels." + LabelSharingAccess
		bsonAnd = append(bsonAnd, bson.M{sharingAccessLabel: SharingAccessViewer})
		isSet = true
	}

	bsonAnd = append(bsonAnd, bson.M{bsonKeyDeleted: false})
	return bson.M{"$and": bsonAnd}, isSet
}
