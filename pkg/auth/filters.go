package auth

import (
	"fmt"

	"github.com/travelreys/travelreys/pkg/common"
	"go.mongodb.org/mongo-driver/bson"
)

type ReadFilter struct {
	ID    string `json:"id" bson:"id,omitempty"`
	Email string `json:"email" bson:"email,omitempty"`
}

type UpdateFilter struct {
	Email    *string        `json:"email" bson:"email"`
	Name     *string        `json:"name" bson:"name"`
	Username *string        `json:"username" bson:"username"`
	Labels   *common.Labels `json:"labels" bson:"labels"`
}

func (ff *UpdateFilter) toBSON() (bson.M, bool) {
	bsonM := bson.M{}
	isSet := false

	if ff.Email != nil && *ff.Email != "" {
		bsonM[bsonKeyEmail] = ff.Email
		isSet = true
	}
	if ff.Name != nil && *ff.Name != "" {
		bsonM[bsonKeyName] = ff.Name
		isSet = true
	}
	if ff.Username != nil && *ff.Username != "" {
		bsonM[bsonKeyUsername] = ff.Username
		isSet = true
	}
	if ff.Labels != nil {
		bsonM[bsonKeyLabels] = ff.Labels
		isSet = true
	}
	return bsonM, isSet
}

func (ff *UpdateFilter) Validate() error {
	if ff.Email != nil && !common.EmailRegexp.MatchString(*ff.Email) {
		return ErrInvalidFilter
	}
	if ff.Name != nil && *ff.Name == "" {
		return ErrInvalidFilter
	}
	if ff.Username != nil && !IsValidUsername(*ff.Username) {
		return ErrInvalidFilter
	}
	_, ok := ff.toBSON()
	if !ok {
		return ErrInvalidFilter
	}
	return nil
}

type ListFilter struct {
	IDs       []string `json:"ids" bson:"id"`
	Username  *string  `json:"username" bson:"username"`
	Email     *string  `json:"email" bson:"email"`
	PageCount *int     `json:"pageCount" bson:"pageCount"`
}

func (ff ListFilter) Validate() error {
	if ff.Username != nil && !UsernameRegexp.MatchString(*ff.Username) {
		return ErrInvalidFilter
	}
	if ff.Email != nil && !common.EmailRegexp.MatchString(*ff.Email) {
		return ErrInvalidFilter
	}
	_, ok := ff.toBSON()
	if !ok {
		return ErrInvalidFilter
	}
	return nil
}

func (ff ListFilter) toBSON() (bson.M, bool) {
	bsonM := bson.M{}
	isSet := false

	if len(ff.IDs) > 0 {
		bsonM[bsonKeyID] = bson.M{"$in": ff.IDs}
		isSet = true
	}
	if ff.Email != nil && *ff.Username != "" {
		bsonM[bsonKeyEmail] = ff.Email
		isSet = true
	}
	if ff.Username != nil && *ff.Username != "" {
		bsonM[bsonKeyUsername] = bson.M{
			"$regex": fmt.Sprintf("^%s", *ff.Username),
		}
		isSet = true
	}
	return bsonM, isSet
}
