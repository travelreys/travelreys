package footprints

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	defaultPageSize = 10

	bsonKeyID          = "id"
	bsonKeyUserID      = "userID"
	bsonKeyPlaceID     = "placeID"
	footprintsItemColl = "footprints"
	storeLoggerName    = "footprints.store"
)

var (
	ErrFpNotFound           = errors.New("footprints.store.footprintNotFound")
	ErrUnexpectedStoreError = errors.New("footprints.store.unexpected-error")
)

type ListFootprintsFilter struct {
	UserID *string `json:"userID" bson:"userID"`
}

func (ff ListFootprintsFilter) toBSON() bson.M {
	bsonM := bson.M{}
	if ff.UserID != nil || *ff.UserID != "" {
		bsonM["userID"] = *ff.UserID
	}

	return bsonM
}

type Store interface {
	Read(ctx context.Context, userID, placeID string) (Footprint, error)
	Save(ctx context.Context, userID string, footprint Footprint) error
	List(ctx context.Context, ff ListFootprintsFilter) (FootprintList, error)
	Delete(ctx context.Context, id string) error
}

type store struct {
	db   *mongo.Database
	coll *mongo.Collection

	logger *zap.Logger
}

func NewStore(ctx context.Context, db *mongo.Database, logger *zap.Logger) Store {
	ctx, cancel := context.WithTimeout(ctx, common.DbReqTimeout)
	defer cancel()

	coll := db.Collection(footprintsItemColl)
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyUserID: 1, bsonKeyPlaceID: 1}},
	})
	return &store{db, coll, logger.Named(storeLoggerName)}
}

func (str *store) Read(ctx context.Context, userID, placeID string) (Footprint, error) {
	var fp Footprint
	err := str.coll.FindOne(ctx, bson.M{bsonKeyUserID: userID, bsonKeyPlaceID: placeID}).Decode(&fp)
	if err == mongo.ErrNoDocuments {
		return fp, ErrFpNotFound
	}
	if err != nil {
		str.logger.Error(
			"Read",
			zap.String("userID", userID), zap.String("placeID", placeID), zap.Error(err),
		)
		return fp, ErrUnexpectedStoreError
	}
	return fp, nil
}

func (str *store) Save(ctx context.Context, userID string, footprint Footprint) error {
	saveFF := bson.M{bsonKeyID: footprint.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := str.coll.ReplaceOne(ctx, saveFF, footprint, opts)
	if err != nil {
		str.logger.Error("Save", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (str *store) List(ctx context.Context, ff ListFootprintsFilter) (FootprintList, error) {
	list := FootprintList{}
	bsonFF := ff.toBSON()

	cursor, err := str.coll.Find(ctx, bsonFF, options.Find())
	if err != nil {
		str.logger.Error("List", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return list, err
	}

	cursor.All(ctx, &list)
	return list, err
}

func (str *store) Delete(ctx context.Context, id string) error {
	_, err := str.coll.DeleteOne(ctx, bson.M{"id": id})
	return err
}
