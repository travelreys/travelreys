package trips

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type ListFilter struct {
	UserID     *string
	UserIDs    []string
	OnlyPublic bool
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
